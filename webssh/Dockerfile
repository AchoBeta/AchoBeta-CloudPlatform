FROM centos:latest
MAINTAINER 	marin<marin_man@163.com>

ENV MYPATH /usr/local
WORKDIR $MYPATH

EXPOSE 8888

COPY ./webssh/webssh.sh $MYPATH

RUN chmod u+x $MYPATH/webssh.sh

RUN nohup $MYPATH/webssh.sh > log.out &