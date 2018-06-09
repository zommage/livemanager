# livemanager
live manager

# build
./controll.sh build

#dev 环境
# 镜像打包
docker build --build-arg config=./conf/app.dev.json -t livemanager .

# 停止进程
docker stop  livemanager && docker rm -f livemanager

# docker 启动
docker run --restart=always --name livemanager  -idt -p 16271:16271 -p 16272:16272 -v /var/log/zommage/livemanager:/home/livemanager/logs  livemanager:latest 

# 查看日志
tailf /var/log/zommage/livemanager/livemanager.log
