#!/bin/bash

###############################################################################
# dongfeng-pay 服务快速启动脚本
# 用于开发和测试环境
###############################################################################

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="/home/user/dongfeng-pay"

# 打印标题
print_header() {
    echo -e "${CYAN}============================================${NC}"
    echo -e "${CYAN}$1${NC}"
    echo -e "${CYAN}============================================${NC}"
}

# 打印成功消息
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# 打印错误消息
print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# 打印警告消息
print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# 打印信息消息
print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# 检查端口是否被占用
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        return 0  # 端口被占用
    else
        return 1  # 端口空闲
    fi
}

# 启动服务
start_service() {
    local service_name=$1
    local service_dir=$2
    local port=$3

    print_info "启动 $service_name 服务 (端口 $port)..."

    # 检查端口是否已被占用
    if check_port $port; then
        print_warning "$service_name 服务已在运行 (端口 $port)"
        return 0
    fi

    # 检查服务目录是否存在
    if [ ! -d "$PROJECT_ROOT/$service_dir" ]; then
        print_error "$service_name 服务目录不存在: $PROJECT_ROOT/$service_dir"
        return 1
    fi

    # 进入服务目录
    cd "$PROJECT_ROOT/$service_dir"

    # 检查 main.go 是否存在
    if [ ! -f "main.go" ]; then
        print_error "$service_name 服务的 main.go 不存在"
        return 1
    fi

    # 在后台启动服务
    nohup go run main.go > "$PROJECT_ROOT/logs/${service_name}.log" 2>&1 &
    local pid=$!

    # 等待服务启动
    sleep 2

    # 检查服务是否成功启动
    if ps -p $pid > /dev/null 2>&1; then
        # 再次检查端口
        if check_port $port; then
            print_success "$service_name 服务启动成功 (PID: $pid, 端口: $port)"
            echo $pid > "$PROJECT_ROOT/logs/${service_name}.pid"
            return 0
        else
            print_warning "$service_name 服务进程已启动，但端口 $port 未就绪，可能需要更多时间..."
            echo $pid > "$PROJECT_ROOT/logs/${service_name}.pid"
            return 0
        fi
    else
        print_error "$service_name 服务启动失败，请检查日志: $PROJECT_ROOT/logs/${service_name}.log"
        return 1
    fi
}

# 停止服务
stop_service() {
    local service_name=$1
    local pid_file="$PROJECT_ROOT/logs/${service_name}.pid"

    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if ps -p $pid > /dev/null 2>&1; then
            kill $pid
            print_success "$service_name 服务已停止 (PID: $pid)"
            rm "$pid_file"
        else
            print_warning "$service_name 服务进程不存在 (PID: $pid)"
            rm "$pid_file"
        fi
    else
        print_warning "$service_name 服务未运行"
    fi
}

# 检查服务状态
check_service_status() {
    local service_name=$1
    local port=$2

    if check_port $port; then
        echo -e "${GREEN}✓${NC} $service_name (端口 $port) - ${GREEN}运行中${NC}"
    else
        echo -e "${RED}✗${NC} $service_name (端口 $port) - ${RED}未运行${NC}"
    fi
}

# 启动所有服务
start_all() {
    print_header "启动所有服务"

    # 创建日志目录
    mkdir -p "$PROJECT_ROOT/logs"

    # 启动各个服务
    start_service "Boss" "boss" 12306
    start_service "Merchant" "merchant" 12307
    start_service "Agent" "agent" 12308
    start_service "Gateway" "gateway" 12309

    echo ""
    print_header "服务启动完成"
    echo ""
    print_info "等待 5 秒让服务完全启动..."
    sleep 5
    echo ""

    # 显示服务状态
    show_status
}

# 停止所有服务
stop_all() {
    print_header "停止所有服务"

    stop_service "Boss"
    stop_service "Merchant"
    stop_service "Agent"
    stop_service "Gateway"

    echo ""
    print_success "所有服务已停止"
}

# 重启所有服务
restart_all() {
    print_header "重启所有服务"
    stop_all
    sleep 2
    start_all
}

# 显示服务状态
show_status() {
    print_header "服务状态"
    echo ""
    check_service_status "Boss" 12306
    check_service_status "Merchant" 12307
    check_service_status "Agent" 12308
    check_service_status "Gateway" 12309
    echo ""

    # 显示日志文件位置
    if [ -d "$PROJECT_ROOT/logs" ]; then
        echo -e "${BLUE}日志文件位置:${NC} $PROJECT_ROOT/logs/"
        echo ""
    fi
}

# 查看日志
view_logs() {
    local service_name=$1

    if [ -z "$service_name" ]; then
        echo "请指定服务名称: boss, merchant, agent, gateway"
        return 1
    fi

    local log_file="$PROJECT_ROOT/logs/${service_name}.log"

    if [ -f "$log_file" ]; then
        print_info "查看 $service_name 日志 (按 Ctrl+C 退出):"
        echo ""
        tail -f "$log_file"
    else
        print_error "日志文件不存在: $log_file"
    fi
}

# 主函数
main() {
    case "$1" in
        start)
            start_all
            ;;
        stop)
            stop_all
            ;;
        restart)
            restart_all
            ;;
        status)
            show_status
            ;;
        logs)
            view_logs "$2"
            ;;
        *)
            echo "用法: $0 {start|stop|restart|status|logs <service>}"
            echo ""
            echo "命令:"
            echo "  start    - 启动所有服务"
            echo "  stop     - 停止所有服务"
            echo "  restart  - 重启所有服务"
            echo "  status   - 查看服务状态"
            echo "  logs     - 查看服务日志 (boss|merchant|agent|gateway)"
            echo ""
            echo "示例:"
            echo "  $0 start              # 启动所有服务"
            echo "  $0 status             # 查看状态"
            echo "  $0 logs boss          # 查看Boss日志"
            echo "  $0 stop               # 停止所有服务"
            exit 1
            ;;
    esac
}

# 运行主函数
main "$@"
