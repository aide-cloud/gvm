#!/bin/bash

# GVM 安装脚本
# 用于安装当前项目的 GVM 工具

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查 Go 环境
check_go() {
    print_info "检查 Go 环境..."
    
    if ! command -v go &> /dev/null; then
        print_error "未找到 Go 环境，请先安装 Go"
        print_info "访问 https://golang.org/dl/ 下载并安装 Go"
        exit 1
    fi
    
    local go_version=$(go version | cut -d' ' -f3 | sed 's/go//')
    print_success "找到 Go 版本: $go_version"
}

# 构建 GVM
build_gvm() {
    print_info "构建 GVM..."
    
    # 检查是否在正确的目录
    if [ ! -f "go.mod" ]; then
        print_error "请在 GVM 项目根目录下运行此脚本"
        exit 1
    fi
    
    # 下载依赖
    print_info "下载依赖..."
    go mod download
    
    # 构建二进制文件
    print_info "构建二进制文件..."
    if ! go build -o gvm .; then
        print_error "构建失败"
        exit 1
    fi
    
    print_success "构建完成"
}

# 安装 GVM
install_gvm() {
    local install_dir="$HOME/.gvm"
    local bin_dir="$HOME/.local/bin"
    
    print_info "安装 GVM 到 $install_dir..."
    
    # 创建目录
    mkdir -p "$install_dir"
    mkdir -p "$bin_dir"
    
    # 复制二进制文件
    cp gvm "$install_dir/"
    chmod +x "$install_dir/gvm"
    
    # 创建符号链接
    if [ -L "$bin_dir/gvm" ]; then
        rm "$bin_dir/gvm"
    fi
    ln -s "$install_dir/gvm" "$bin_dir/gvm"
    
    print_success "GVM 安装完成"
}

# 配置环境
configure_shell() {
    local shell_config_file
    
    # 确定 shell 配置文件
    case "$SHELL" in
        */zsh) shell_config_file="$HOME/.zshrc" ;;
        */bash) shell_config_file="$HOME/.bashrc" ;;
        */fish) shell_config_file="$HOME/.config/fish/config.fish" ;;
        *) shell_config_file="$HOME/.profile" ;;
    esac
    
    print_info "配置环境变量到 $shell_config_file"
    
    # 检查是否已经配置过
    if grep -q "GVM_ROOT" "$shell_config_file" 2>/dev/null; then
        print_warning "GVM 配置已存在，跳过配置"
        return
    fi
    
    # 添加配置
    cat >> "$shell_config_file" << 'EOF'

# GVM (Go Version Manager) 配置
export GVM_ROOT="$HOME/.gvm"
export GVM_SDK_DIR="$HOME/.gvm/sdk"
export GVM_CACHE_DIR="$HOME/.gvm/cache"

# 添加 GVM 到 PATH
if [ -d "$HOME/.local/bin" ]; then
    export PATH="$HOME/.local/bin:$PATH"
fi

# GVM 初始化函数
gvm() {
    if [ "$1" = "use" ] && [ -n "$2" ]; then
        # 设置 GOROOT 和 PATH
        local go_version="$2"
        local go_root="$GVM_SDK_DIR/$go_version"
        if [ -d "$go_root" ]; then
            export GOROOT="$go_root"
            export PATH="$go_root/bin:$PATH"
            echo "Go version $go_version activated"
        else
            echo "Go version $go_version not found. Use 'gvm install $go_version' to install it."
        fi
    else
        # 调用实际的 gvm 命令
        "$GVM_ROOT/gvm" "$@"
    fi
}
EOF
    
    print_success "环境配置完成"
}

# 验证安装
verify_installation() {
    print_info "验证安装..."
    
    # 检查二进制文件
    if [ ! -f "$HOME/.gvm/gvm" ]; then
        print_error "GVM 二进制文件不存在"
        return 1
    fi
    
    # 检查符号链接
    if [ ! -L "$HOME/.local/bin/gvm" ]; then
        print_error "GVM 符号链接不存在"
        return 1
    fi
    
    # 测试命令
    if ! "$HOME/.gvm/gvm" --help &> /dev/null; then
        print_error "GVM 命令执行失败"
        return 1
    fi
    
    print_success "安装验证通过"
}

# 显示使用说明
show_usage() {
    print_success "GVM 安装完成！"
    echo
    print_info "使用方法："
    echo "  1. 重新加载 shell 配置:"
    echo "     source ~/.zshrc  # 或 source ~/.bashrc"
    echo
    echo "  2. 安装 Go 版本:"
    echo "     gvm install 1.21.0"
    echo "     gvm install latest"
    echo
    echo "  3. 使用 Go 版本:"
    echo "     gvm use 1.21.0"
    echo
    echo "  4. 列出可用版本:"
    echo "     gvm list"
    echo
    echo "  5. 查看帮助:"
    echo "     gvm --help"
    echo
    print_warning "请重新启动终端或运行 'source ~/.zshrc' 来激活 GVM"
}

# 清理构建文件
cleanup() {
    print_info "清理构建文件..."
    rm -f gvm
    print_success "清理完成"
}

# 主函数
main() {
    print_info "开始安装 GVM (Go Version Manager)..."
    echo
    
    # 检查 Go 环境
    check_go
    
    # 构建 GVM
    build_gvm
    
    # 安装 GVM
    install_gvm
    
    # 配置环境
    configure_shell
    
    # 验证安装
    if ! verify_installation; then
        print_error "安装验证失败"
        exit 1
    fi
    
    # 清理
    cleanup
    
    # 显示使用说明
    show_usage
}

# 运行主函数
main "$@"