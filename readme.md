### 前提条件

首先确保你的电脑上已经安装了go以及配置好了go的环境变量

确保你的服务器或者本地ip可以连通claude.ai

# 下载
```Shell
git clone
```
# 编译

```Shell
GOOS=linux GOARCH=amd64 go build -o getkeyinfo
```

# 部署

上传到服务器目录，打开终端 进入到你上传的目录下

```shell
chmod +x ./claudetoapi
nohup ./claudetoapi >> run.log 2>&1 &
```

现在你就可以使用你的服务器ip:8080/v1/complete或者ip:8080/v1/chat/completions使用了

调用方式参考openai的api调用方式以及claude的api调用方式，已做兼容处理

# 说明

本项目只是练手项目！使用本项目造成的责任与本人无关