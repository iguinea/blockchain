.PHONY:

blockchainserver: 
	go run cmd/blockchain_server/main.go

walletserver:
	go run cmd/wallet_server/main.go


########################################################
# Compile modules
########################################################
compile:
	go build pkgs/block/*.go




########################################################
# Build binaries
########################################################
build_blockchainserver:
	go build -o bin/blockchain_server cmd/blockchain_server/main.go 
	
build_walletserver:
	go build -o bin/wallet_server     cmd/wallet_server/main.go 

########################################################
# Build docker images
########################################################
build_docker: build_docker_blockchainserver build_docker_walletserver

build_docker_blockchainserver: build_blockchainserver
	mkdir -p dockerfiles/blockchainserver/files/bin  
	cp bin/blockchain_server dockerfiles/blockchainserver/files/bin/
	docker build --pull --rm -f "dockerfiles/blockchainserver/Dockerfile" -t blockchainserver:latest "dockerfiles/blockchainserver"

build_docker_walletserver: build_walletserver
	mkdir -p dockerfiles/walletserver/files/bin  
	cp bin/wallet_server dockerfiles/walletserver/files/bin/
	docker build --pull --rm -f "dockerfiles/walletserver/Dockerfile" -t walletserver:latest "dockerfiles/walletserver"

########################################################
# Run docker images
########################################################
run_docker_blockchainserver:
	docker run --rm -it  blockchainserver:latest
run_docker_blockchainserver1:
	docker run --rm -it -p 5001:5000/tcp blockchainserver:latest
run_docker_blockchainserver2:
	docker run --rm -it -p 5002:5000/tcp blockchainserver:latest
run_docker_blockchainserver3:
	docker run --rm -it -p 5003:5000/tcp blockchainserver:latest

run_docker_walletserver:
	docker run --rm -it  -p 8080:8080/tcp walletserver:latest

run_docker_walletserver1:
	docker run --rm -it  -p 8081:8080/tcp walletserver:latest
run_docker_walletserver2:
	docker run --rm -it  -p 8082:8080/tcp walletserver:latest

########################################################
# Playground 
########################################################

run_playground:
	go run cmd/neighbor/main.go
playground: 
	go build -o bin/playground cmd/neighbor/main.go 
	cp bin/playground dockerfiles/playground/files/bin/
	docker build --pull --rm -f "dockerfiles/playground/Dockerfile" -t playground:latest "dockerfiles/playground"
run_docker_playground:
	docker run --rm -it  playground:latest
