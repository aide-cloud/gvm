# GVM (Go Version Manager)


## 需要先创建必要的目录结构

```
mkdir -p ~/.gvm/versions
```

##  编译安装

```
go install github.com/moovweb/gvm@latest
```

## 常用命令

```
# 查看可用版本
gvm list-remote

# 安装指定版本
gvm install 1.21.0

# 切换版本
gvm use 1.21.0

# 列出已安装版本
gvm list
```