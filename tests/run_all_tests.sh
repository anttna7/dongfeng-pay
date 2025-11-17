#!/bin/bash

###############################################################################
# dongfeng-pay 完整测试套件
# 执行所有测试并生成综合报告
###############################################################################

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 测试结果目录
REPORT_DIR="./test-reports"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
REPORT_FILE="$REPORT_DIR/test_report_$TIMESTAMP.md"

# 创建报告目录
mkdir -p "$REPORT_DIR"

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

# 初始化报告
init_report() {
    cat > "$REPORT_FILE" <<EOF
# dongfeng-pay 测试报告

**生成时间**: $(date '+%Y-%m-%d %H:%M:%S')

---

## 执行摘要

EOF
}

# 检查服务状态
check_services() {
    print_header "步骤 1: 检查服务状态"

    local all_running=true

    # 检查 Boss 服务
    if curl -s http://localhost:12306/ > /dev/null 2>&1; then
        print_success "Boss 服务运行中 (端口 12306)"
        echo "- ✅ Boss 服务 (12306) - 运行中" >> "$REPORT_FILE"
    else
        print_error "Boss 服务未运行 (端口 12306)"
        echo "- ❌ Boss 服务 (12306) - 未运行" >> "$REPORT_FILE"
        all_running=false
    fi

    # 检查 Merchant 服务
    if curl -s http://localhost:12307/ > /dev/null 2>&1; then
        print_success "Merchant 服务运行中 (端口 12307)"
        echo "- ✅ Merchant 服务 (12307) - 运行中" >> "$REPORT_FILE"
    else
        print_error "Merchant 服务未运行 (端口 12307)"
        echo "- ❌ Merchant 服务 (12307) - 未运行" >> "$REPORT_FILE"
        all_running=false
    fi

    # 检查 Agent 服务
    if curl -s http://localhost:12308/ > /dev/null 2>&1; then
        print_success "Agent 服务运行中 (端口 12308)"
        echo "- ✅ Agent 服务 (12308) - 运行中" >> "$REPORT_FILE"
    else
        print_error "Agent 服务未运行 (端口 12308)"
        echo "- ❌ Agent 服务 (12308) - 未运行" >> "$REPORT_FILE"
        all_running=false
    fi

    # 检查 Gateway 服务
    if curl -s http://localhost:12309/ > /dev/null 2>&1; then
        print_success "Gateway 服务运行中 (端口 12309)"
        echo "- ✅ Gateway 服务 (12309) - 运行中" >> "$REPORT_FILE"
    else
        print_error "Gateway 服务未运行 (端口 12309)"
        echo "- ❌ Gateway 服务 (12309) - 未运行" >> "$REPORT_FILE"
        all_running=false
    fi

    echo "" >> "$REPORT_FILE"

    if [ "$all_running" = false ]; then
        print_warning "部分服务未运行，某些测试可能失败"
        return 1
    fi

    return 0
}

# 运行单元测试
run_unit_tests() {
    print_header "步骤 2: 运行单元测试"

    echo "## 单元测试结果" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"

    if [ -f "system_test.go" ]; then
        print_info "执行系统测试套件..."

        if go test -v -timeout 30s > "$REPORT_DIR/unit_test_$TIMESTAMP.log" 2>&1; then
            print_success "单元测试全部通过"
            echo "✅ **单元测试**: 全部通过" >> "$REPORT_FILE"

            # 提取测试统计
            local pass_count=$(grep -c "PASS:" "$REPORT_DIR/unit_test_$TIMESTAMP.log" || echo "0")
            echo "- 通过测试数: $pass_count" >> "$REPORT_FILE"
        else
            print_error "单元测试失败"
            echo "❌ **单元测试**: 存在失败" >> "$REPORT_FILE"

            # 提取失败信息
            grep "FAIL:" "$REPORT_DIR/unit_test_$TIMESTAMP.log" >> "$REPORT_FILE" || true
        fi
    else
        print_warning "未找到单元测试文件"
        echo "⚠️ **单元测试**: 未找到测试文件" >> "$REPORT_FILE"
    fi

    echo "" >> "$REPORT_FILE"
}

# 测试登录页面
test_login_pages() {
    print_header "步骤 3: 测试登录页面"

    echo "## 登录页面测试" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"

    # 测试 Boss 登录页
    print_info "测试 Boss 登录页..."
    if curl -s -o /dev/null -w "%{http_code}" http://localhost:12306/login_modern.html | grep -q "200"; then
        print_success "Boss 登录页加载成功"
        echo "- ✅ Boss 登录页 (login_modern.html) - 可访问" >> "$REPORT_FILE"
    else
        print_error "Boss 登录页加载失败"
        echo "- ❌ Boss 登录页 (login_modern.html) - 无法访问" >> "$REPORT_FILE"
    fi

    # 测试 Merchant 登录页
    print_info "测试 Merchant 登录页..."
    if curl -s -o /dev/null -w "%{http_code}" http://localhost:12307/login_modern.html | grep -q "200"; then
        print_success "Merchant 登录页加载成功"
        echo "- ✅ Merchant 登录页 (login_modern.html) - 可访问" >> "$REPORT_FILE"
    else
        print_error "Merchant 登录页加载失败"
        echo "- ❌ Merchant 登录页 (login_modern.html) - 无法访问" >> "$REPORT_FILE"
    fi

    # 测试 Agent 登录页
    print_info "测试 Agent 登录页..."
    if curl -s -o /dev/null -w "%{http_code}" http://localhost:12308/login_modern.html | grep -q "200"; then
        print_success "Agent 登录页加载成功"
        echo "- ✅ Agent 登录页 (login_modern.html) - 可访问" >> "$REPORT_FILE"
    else
        print_error "Agent 登录页加载失败"
        echo "- ❌ Agent 登录页 (login_modern.html) - 无法访问" >> "$REPORT_FILE"
    fi

    echo "" >> "$REPORT_FILE"
}

# 测试 API 端点
test_api_endpoints() {
    print_header "步骤 4: 测试 API 端点"

    echo "## API 端点测试" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"

    # 测试验证码接口
    print_info "测试验证码接口..."

    local boss_captcha=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:12306/getVerifyImg)
    if [ "$boss_captcha" = "200" ]; then
        print_success "Boss 验证码接口正常"
        echo "- ✅ Boss 验证码接口 (/getVerifyImg)" >> "$REPORT_FILE"
    else
        print_error "Boss 验证码接口异常 (状态码: $boss_captcha)"
        echo "- ❌ Boss 验证码接口 (/getVerifyImg) - 状态码: $boss_captcha" >> "$REPORT_FILE"
    fi

    # 测试健康检查接口（如果存在）
    print_info "测试健康检查接口..."

    local gateway_health=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:12309/health 2>/dev/null)
    if [ "$gateway_health" = "200" ]; then
        print_success "Gateway 健康检查接口正常"
        echo "- ✅ Gateway 健康检查 (/health)" >> "$REPORT_FILE"
    else
        print_warning "Gateway 健康检查接口未配置或异常"
        echo "- ⚠️ Gateway 健康检查 (/health) - 未配置" >> "$REPORT_FILE"
    fi

    echo "" >> "$REPORT_FILE"
}

# 性能测试
performance_test() {
    print_header "步骤 5: 性能测试"

    echo "## 性能测试" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"

    # 测试页面加载时间
    print_info "测试页面加载时间..."

    local boss_time=$(curl -s -o /dev/null -w "%{time_total}" http://localhost:12306/login_modern.html)
    local merchant_time=$(curl -s -o /dev/null -w "%{time_total}" http://localhost:12307/login_modern.html)
    local agent_time=$(curl -s -o /dev/null -w "%{time_total}" http://localhost:12308/login_modern.html)

    echo "### 页面加载时间" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    echo "| 模块 | 加载时间 | 评级 |" >> "$REPORT_FILE"
    echo "|------|---------|------|" >> "$REPORT_FILE"

    # Boss
    if (( $(echo "$boss_time < 0.5" | bc -l) )); then
        echo "| Boss | ${boss_time}s | 🟢 优秀 |" >> "$REPORT_FILE"
        print_success "Boss 加载时间: ${boss_time}s (优秀)"
    elif (( $(echo "$boss_time < 1.0" | bc -l) )); then
        echo "| Boss | ${boss_time}s | 🟡 良好 |" >> "$REPORT_FILE"
        print_info "Boss 加载时间: ${boss_time}s (良好)"
    else
        echo "| Boss | ${boss_time}s | 🔴 需优化 |" >> "$REPORT_FILE"
        print_warning "Boss 加载时间: ${boss_time}s (需优化)"
    fi

    # Merchant
    if (( $(echo "$merchant_time < 0.5" | bc -l) )); then
        echo "| Merchant | ${merchant_time}s | 🟢 优秀 |" >> "$REPORT_FILE"
        print_success "Merchant 加载时间: ${merchant_time}s (优秀)"
    elif (( $(echo "$merchant_time < 1.0" | bc -l) )); then
        echo "| Merchant | ${merchant_time}s | 🟡 良好 |" >> "$REPORT_FILE"
        print_info "Merchant 加载时间: ${merchant_time}s (良好)"
    else
        echo "| Merchant | ${merchant_time}s | 🔴 需优化 |" >> "$REPORT_FILE"
        print_warning "Merchant 加载时间: ${merchant_time}s (需优化)"
    fi

    # Agent
    if (( $(echo "$agent_time < 0.5" | bc -l) )); then
        echo "| Agent | ${agent_time}s | 🟢 优秀 |" >> "$REPORT_FILE"
        print_success "Agent 加载时间: ${agent_time}s (优秀)"
    elif (( $(echo "$agent_time < 1.0" | bc -l) )); then
        echo "| Agent | ${agent_time}s | 🟡 良好 |" >> "$REPORT_FILE"
        print_info "Agent 加载时间: ${agent_time}s (良好)"
    else
        echo "| Agent | ${agent_time}s | 🔴 需优化 |" >> "$REPORT_FILE"
        print_warning "Agent 加载时间: ${agent_time}s (需优化)"
    fi

    echo "" >> "$REPORT_FILE"
}

# 生成最终报告
generate_final_report() {
    print_header "步骤 6: 生成最终报告"

    # 添加建议
    cat >> "$REPORT_FILE" <<EOF

---

## 测试建议

### ✅ 通过的测试
- 继续保持代码质量
- 定期运行测试套件

### ⚠️ 需要关注
- 确保所有服务在生产环境中正常运行
- 监控页面加载性能
- 定期检查 API 响应时间

### 🔧 优化建议
1. **性能优化**: 对于加载时间 > 1s 的页面进行优化
2. **监控增强**: 添加实时性能监控
3. **自动化测试**: 集成到 CI/CD 流程

---

## 测试环境

- **操作系统**: $(uname -s)
- **Go 版本**: $(go version | cut -d' ' -f3)
- **测试时间**: $(date '+%Y-%m-%d %H:%M:%S')
- **报告位置**: $REPORT_FILE

---

**报告生成完成** ✅
EOF

    print_success "测试报告已生成: $REPORT_FILE"
}

# 主函数
main() {
    clear
    print_header "dongfeng-pay 完整测试套件"
    echo ""

    # 初始化报告
    init_report

    # 执行测试步骤
    check_services
    run_unit_tests
    test_login_pages
    test_api_endpoints
    performance_test
    generate_final_report

    echo ""
    print_header "测试完成"
    echo ""
    print_info "查看详细报告: cat $REPORT_FILE"
    print_info "或使用 Markdown 查看器打开报告"
    echo ""
}

# 运行主函数
main
