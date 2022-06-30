FROM ubuntu:latest

RUN apt update
RUN apt upgrade -y

ENV APP_HOME /home
WORKDIR $APP_HOME

COPY files/bin/blockchain_server /home/blockchain_server 

CMD /home/blockchain_server
