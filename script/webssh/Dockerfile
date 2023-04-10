FROM centos:latest
MAINTAINER 	marin<marin_man@163.com>

ENV MYPATH /usr/local
WORKDIR $MYPATH

EXPOSE 8888

COPY ./webssh.sh $MYPATH

RUN chmod u+x $MYPATH/webssh.sh

RUN $MYPATH/webssh.sh

CMD nohup wssh --fbidhttp=False --xsrf=False
