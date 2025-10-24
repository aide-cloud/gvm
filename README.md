# GVM (Go Version Manager)


##  编译安装

```bash
go install github.com/aide-cloud/gvm@latest
```

## 放置到bin目录
```bash
sudo mv $(which gvm) /usr/local/bin/
```

## 常用命令

```bash
# 列出已安装版本
gvm list 

# 查看最新版本 默认10条
gvm list -o 
# 查看所有版本
gvm list -o -a 

# 安装指定版本
gvm install 1.21.0

# 安装最新版本
gvm install latest

# 切换版本
gvm use 1.21.0

# 列出已使用版本
gvm ls
```

## API

```bash
https://dl.google.com/go/go1.25.2.darwin-amd64.tar.gz
```

```bash
https://mirrors.aliyun.com/golang/go1.25.2.darwin-amd64.tar.gz
```

```bash
https://go.dev/dl/?mode=json&include=all
```