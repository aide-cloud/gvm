# GVM (Go Version Manager)

一个简单易用的 Go 版本管理工具，支持安装、切换、卸载不同版本的 Go。

## 特性

- 🚀 快速安装和切换 Go 版本
- 📦 支持安装最新版本和指定版本
- 🔄 自动管理本地版本缓存
- ⚙️ 可配置的下载源和存储路径
- 🛠️ 支持强制安装和卸载
- 📋 列出可用版本和已安装版本

## 安装

### 方式一：使用 go install（推荐）

```bash
go install github.com/aide-cloud/gvm@latest
```

### 方式二：从源码编译

```bash
git clone https://github.com/aide-cloud/gvm.git
cd gvm
go build -o gvm .
sudo mv gvm /usr/local/bin/
```

### 方式三：使用安装脚本

```bash
curl -sSL https://raw.githubusercontent.com/aide-cloud/gvm/main/install.sh | bash
```

## 快速开始

### 1. 查看可用版本

```bash
# 列出最新的 10 个版本（默认）
gvm list

# 列出最新的 20 个版本
gvm list -n 20

# 只显示最新版本
gvm list -l

# 强制更新版本列表
gvm list --force-update
```

### 2. 安装 Go 版本

```bash
# 安装指定版本
gvm install 1.21.0

# 安装最新版本
gvm install latest
# 或者
gvm install -l

# 强制重新安装
gvm install 1.21.0 -f
```

### 3. 切换版本

```bash
# 切换到指定版本
gvm use 1.21.0

# 切换到最新版本
gvm use latest
# 或者
gvm use -l

# 强制切换（如果版本未安装会自动安装）
gvm use 1.21.0 -f
```

### 4. 查看已安装版本

```bash
# 列出本地已安装的版本
gvm ls
```

### 5. 卸载版本

```bash
# 卸载指定版本
gvm uninstall 1.21.0

# 卸载最新版本
gvm uninstall latest
# 或者
gvm uninstall -l
```

## 命令详解

### `gvm list` - 列出可用版本

```bash
gvm list [flags]
```

**参数：**
- `-n, --number int`: 显示版本数量（默认：10）
- `-l, --latest`: 只显示最新版本
- `--force-update`: 强制更新版本列表

**示例：**
```bash
gvm list                    # 显示最新 10 个版本
gvm list -n 20              # 显示最新 20 个版本
gvm list -l                 # 只显示最新版本
gvm list --force-update     # 强制更新并显示版本
```

### `gvm install` - 安装 Go 版本

```bash
gvm install <version> [flags]
```

**参数：**
- `-l, --latest`: 安装最新版本
- `-f, --force`: 强制安装（覆盖已安装版本）

**示例：**
```bash
gvm install 1.21.0          # 安装指定版本
gvm install latest          # 安装最新版本
gvm install -l              # 安装最新版本
gvm install 1.21.0 -f       # 强制安装指定版本
```

### `gvm use` - 切换 Go 版本

```bash
gvm use <version> [flags]
```

**参数：**
- `-l, --latest`: 使用最新版本
- `-f, --force`: 强制切换（自动安装未安装的版本）

**示例：**
```bash
gvm use 1.21.0              # 切换到指定版本
gvm use latest              # 切换到最新版本
gvm use -l                  # 切换到最新版本
gvm use 1.21.0 -f           # 强制切换到指定版本
```

### `gvm ls` - 列出已安装版本

```bash
gvm ls
```

显示本地已安装的 Go 版本，当前使用的版本会标记为 `*`。

### `gvm uninstall` - 卸载 Go 版本

```bash
gvm uninstall <version> [flags]
```

**参数：**
- `-l, --latest`: 卸载最新版本

**示例：**
```bash
gvm uninstall 1.21.0        # 卸载指定版本
gvm uninstall latest        # 卸载最新版本
gvm uninstall -l            # 卸载最新版本
```

## 配置选项

GVM 支持通过环境变量或命令行参数进行配置：

### 环境变量

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| `GVM_ORIGIN_URL` | `https://go.dev/dl/?mode=json&include=all` | 版本列表获取地址 |
| `GVM_DOWNLOAD_URL` | `https://dl.google.com/go/` | Go 下载地址 |
| `GVM_CACHE_DIR` | `~/.gvm/cache` | 缓存目录 |
| `GVM_SDK_DIR` | `~/go/sdk` | SDK 存储目录 |
| `GVM_VERSION_FILE_PATH` | `~/.gvm/versions.json` | 版本信息文件路径 |
| `GVM_LOCAL_VERSION_FILE_PATH` | `~/.gvm/version` | 当前版本文件路径 |

### 命令行参数

所有命令都支持以下全局参数：

```bash
--origin-url string         # 版本列表获取地址
--download-url string       # Go 下载地址
--cache-dir string          # 缓存目录
--sdk-dir string            # SDK 存储目录
--version-file-path string  # 版本信息文件路径
--local-version-file string # 当前版本文件路径
--eval                      # 静默模式（不输出日志）
```

### 配置示例

```bash
# 使用阿里云镜像源
export GVM_DOWNLOAD_URL="https://mirrors.aliyun.com/golang/"
gvm install 1.21.0

# 自定义存储目录
export GVM_SDK_DIR="/opt/go/sdk"
export GVM_CACHE_DIR="/opt/go/cache"
gvm install 1.21.0

# 使用命令行参数
gvm install 1.21.0 --download-url "https://mirrors.aliyun.com/golang/" --sdk-dir "/opt/go/sdk"
```

## 目录结构

GVM 会在以下位置创建文件：

```
~/.gvm/
├── cache/              # 下载缓存
├── versions.json       # 版本信息缓存
└── version            # 当前使用的版本

~/go/sdk/              # Go SDK 存储目录
├── go1.21.0/          # 各版本 Go SDK
├── go1.21.1/
└── ...
```

## 使用场景

### 开发环境管理

```bash
# 为不同项目使用不同 Go 版本
cd project-a
gvm use 1.20.0

cd project-b  
gvm use 1.21.0
```

### CI/CD 环境

```bash
# 在 CI 中安装特定版本
gvm install 1.21.0
gvm use 1.21.0
```

### 版本测试

```bash
# 测试不同版本的兼容性
gvm install 1.20.0
gvm use 1.20.0
go test ./...

gvm use 1.21.0
go test ./...
```

## 故障排除

### 常见问题

1. **权限问题**
   ```bash
   # 确保有写入权限
   sudo chown -R $USER ~/.gvm
   sudo chown -R $USER ~/go
   ```

2. **网络问题**
   ```bash
   # 使用国内镜像源
   export GVM_DOWNLOAD_URL="https://mirrors.aliyun.com/golang/"
   gvm install 1.21.0
   ```

3. **版本不存在**
   ```bash
   # 强制更新版本列表
   gvm list --force-update
   ```

### 清理缓存

```bash
# 清理所有缓存
rm -rf ~/.gvm/cache/*
rm -rf ~/.gvm/versions.json
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License