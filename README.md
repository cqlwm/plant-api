# plant-api
智慧农业小程序Go语言服务端

# 交叉编译
```shell script
gox -os "linux" -arch amd64
```

# 部署

启动
```shell script
sudo nohup /root/plant/plant-api >> out.log 2>&1 &
```

连通性测试
```shell script
curl -i "http://127.0.0.1:9091/plant-api/ping"
```

# 可能会用到的命令
```shell script
# 根据端口号查到应用PID
lsof -i:9091
# 编译成Linux可执行文件
gox -os "linux" -arch amd64
# 授予最高权限
chmod 777 ./plant-api
# 恢复sudoers文件的权限
pkexec chmod 0440 /etc/sudoers
# 启动应用程序
sudo nohup /root/plant/plant-api >> out.log 2>&1 &
# 连通性测试命令
curl -i "http://127.0.0.1:9091/plant-api/ping"
```