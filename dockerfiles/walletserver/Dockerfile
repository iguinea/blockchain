FROM ubuntu:latest

RUN apt update
RUN apt upgrade -y

ENV APP_HOME /home
WORKDIR $APP_HOME

COPY files/bin/wallet_server /home/wallet_server 

CMD /home/wallet_server
