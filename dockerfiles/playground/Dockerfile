FROM ubuntu:latest

RUN apt update
RUN apt upgrade -y

ENV APP_HOME /home
WORKDIR $APP_HOME

COPY files/bin/playground /home/playground 


CMD /home/playground
