FROM centos:7
ARG config 
ADD ./logs /home/livemanager/logs
ADD ./conf /home/livemanager/conf
ADD ./livemanager /home/livemanager/livemanager
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
WORKDIR /home/livemanager

RUN mv ${config} ./conf/active.conf

VOLUME ["/home/livemanager/logs"]

CMD ["/home/livemanager/livemanager", "--config=./conf/active.conf"]
