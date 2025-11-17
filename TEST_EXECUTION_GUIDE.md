# 测试执行完整指南

> dongfeng-pay 系统测试、性能分析和反馈收集完整指南

**版本**: v1.0
**更新时间**: 2025-11-17

---

## 📋 目录

1. [准备工作](#准备工作)
2. [登录页测试](#登录页测试)
3. [系统测试](#系统测试)
4. [性能分析](#性能分析)
5. [数据完整性验证](#数据完整性验证)
6. [压力测试](#压力测试)
7. [反馈收集](#反馈收集)
8. [报告生成](#报告生成)
9. [常见问题](#常见问题)

---

## 🚀 准备工作

### 1. 环境检查

确保以下服务正在运行：

```bash
# 检查所有服务状态
curl http://localhost:12306/   # Boss 服务
curl http://localhost:12307/   # Merchant 服务
curl http://localhost:12308/   # Agent 服务
curl http://localhost:12309/   # Gateway 服务
```

### 2. 启动服务

如果服务未运行，请启动：

```bash
# 启动 Boss 服务
cd boss
go run main.go

# 启动 Merchant 服务
cd merchant
go run main.go

# 启动 Agent 服务
cd agent
go run main.go

# 启动 Gateway 服务
cd gateway
go run main.go
```

### 3. 数据库检查

确保 PostgreSQL 18 正在运行：

```bash
# 检查 PostgreSQL 状态
pg_isready -h localhost -p 5432

# 连接数据库
psql -h localhost -p 5432 -U postgres -d juhe_pay

# 检查表是否存在
\dt
```

### 4. 创建测试报告目录

```bash
cd tests
mkdir -p test-reports
```

---

## 🔐 登录页测试

### 方式一：自动化测试工具

打开登录页测试工具：

```bash
# 在浏览器中打开
open tests/login-test-tool.html

# 或使用 HTTP 服务器
cd tests
python3 -m http.server 8000
# 然后访问 http://localhost:8000/login-test-tool.html
```

#### 测试步骤：

1. **单模块测试**
   - 点击 "Boss 登录页测试" → "运行测试"
   - 观察测试结果（绿色 ✅ = 通过，红色 ✗ = 失败）
   - 重复测试 Merchant 和 Agent

2. **批量测试**
   - 点击 "运行所有测试"
   - 等待所有测试完成
   - 查看测试汇总

3. **查看测试结果**
   - 总测试数
   - 通过/失败数量
   - 测试日志详情

4. **导出报告**
   - 点击 "导出报告"
   - 保存 JSON 文件到本地

### 方式二：手动测试

#### Boss 登录页测试 (端口 12306)

**访问地址**: http://localhost:12306/login_modern.html

**测试项目**:

- [x] **页面加载测试**
  - 页面是否正常加载
  - 所有样式是否正确显示
  - 没有控制台错误

- [x] **验证码测试**
  - 验证码图片是否显示
  - 点击验证码是否刷新
  - 验证码图片清晰度

- [x] **表单验证测试**
  - 空用户名提交 → 应显示错误提示
  - 空密码提交 → 应显示错误提示
  - 空验证码提交 → 应显示错误提示

- [x] **登录功能测试**
  ```
  测试账户：admin
  测试密码：admin123
  验证码：填写图片中的验证码
  ```
  - 正确凭据 → 应跳转到首页
  - 错误用户名 → 应显示用户名错误
  - 错误密码 → 应显示密码错误
  - 错误验证码 → 应显示验证码错误

- [x] **记住密码测试**
  - 勾选"记住密码"登录
  - 刷新页面
  - 验证用户名和密码是否自动填充

- [x] **响应式测试**
  - Desktop (> 1024px): 正常显示
  - Tablet (768px - 1024px): 自适应布局
  - Mobile (< 768px): 移动端优化

#### Merchant 登录页测试 (端口 12307)

**访问地址**: http://localhost:12307/login_modern.html

**特殊测试项**:
- [x] 左右分屏布局显示正常
- [x] 左侧品牌信息显示正确
- [x] 右侧登录表单功能正常
- [x] 绿色主题色正确应用

#### Agent 登录页测试 (端口 12308)

**访问地址**: http://localhost:12308/login_modern.html

**特殊测试项**:
- [x] 左右分屏布局显示正常
- [x] 代理商品牌信息正确
- [x] 蓝紫色主题色正确应用

---

## 🧪 系统测试

### 运行自动化测试脚本

```bash
cd tests

# 运行完整测试套件
./run_all_tests.sh

# 查看测试报告
cat test-reports/test_report_*.md
```

### 测试报告内容

测试脚本将生成包含以下内容的报告：

1. **服务状态检查**
   - Boss 服务状态
   - Merchant 服务状态
   - Agent 服务状态
   - Gateway 服务状态

2. **单元测试结果**
   - 通过的测试数量
   - 失败的测试详情

3. **登录页面测试**
   - 页面可访问性
   - HTTP 状态码

4. **API 端点测试**
   - 验证码接口
   - 健康检查接口

5. **性能测试**
   - 页面加载时间
   - 性能评级

### 运行单元测试

```bash
cd tests

# 运行所有测试
go test -v

# 运行特定测试
go test -v -run TestBossHealth

# 查看测试覆盖率
go test -v -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 📊 性能分析

### 使用性能分析工具

```bash
cd tools

# 编译性能分析器
go build -o performance_analyzer performance_analyzer.go

# 运行性能分析
./performance_analyzer

# 查看报告
cat test-reports/performance_*.md
cat test-reports/performance_*.json
```

### 性能分析报告包含

1. **服务状态汇总**
   - 每个服务的在线/离线状态
   - 响应时间
   - 端点成功率

2. **端点详细信息**
   - 各个端点的响应时间
   - HTTP 状态码
   - 成功/失败状态

3. **优化建议**
   - 离线服务提醒
   - 慢响应警告
   - 失败端点提示

### 手动性能测试

使用 curl 测试响应时间：

```bash
# 测试 Boss 登录页
curl -o /dev/null -s -w "Time: %{time_total}s\n" \
  http://localhost:12306/login_modern.html

# 测试 API 响应
curl -o /dev/null -s -w "Time: %{time_total}s\n" \
  http://localhost:12306/getVerifyImg

# 测试 Gateway
curl -o /dev/null -s -w "Time: %{time_total}s\n" \
  http://localhost:12309/
```

使用 Apache Bench 进行并发测试：

```bash
# 安装 Apache Bench (如果未安装)
sudo apt-get install apache2-utils  # Ubuntu/Debian
brew install httpd                   # macOS

# 并发测试 (10 并发，100 请求)
ab -n 100 -c 10 http://localhost:12306/login_modern.html

# 查看结果
# - Requests per second (QPS)
# - Time per request
# - Transfer rate
```

---

## 🗄️ 数据完整性验证

### 运行数据完整性检查

```bash
cd tools

# 编译数据完整性检查工具
go build -o data_check data_integrity_check.go

# 运行检查
./data_check

# 查看详细输出
./data_check --verbose
```

### 检查项目

1. **数据库连接检查**
   - PostgreSQL 连接状态
   - 数据库版本

2. **表结构检查**
   - 所有表是否存在
   - 表行数统计

3. **订单数据完整性**
   - 订单状态分布
   - 订单金额合理性
   - 孤立订单检测

4. **账户余额完整性**
   - 负余额检测
   - 余额合理性

5. **退款数据完整性**
   - 退款记录完整性
   - 退款状态一致性

### 手动数据验证

```sql
-- 连接数据库
psql -h localhost -p 5432 -U postgres -d juhe_pay

-- 检查表数量
SELECT count(*) FROM information_schema.tables
WHERE table_schema = 'public';

-- 检查订单统计
SELECT
    order_status,
    COUNT(*) as count,
    SUM(order_amount) as total_amount
FROM order_info
GROUP BY order_status;

-- 检查商户数量
SELECT COUNT(*) FROM merchant_info WHERE status = 'active';

-- 检查账户余额
SELECT
    COUNT(*) as total_accounts,
    SUM(CASE WHEN balance < 0 THEN 1 ELSE 0 END) as negative_balance_count
FROM account_info;

-- 检查退款记录
SELECT
    status,
    COUNT(*) as count
FROM refund_info
GROUP BY status;
```

---

## 🔥 压力测试

### 使用内置压测工具

```bash
cd tools

# 编译压测工具
go build -o benchmark benchmark.go

# Boss 服务压测
./benchmark -scenario boss

# Gateway 压测
./benchmark -scenario gateway

# 支付 API 压测
./benchmark -scenario payment

# 自定义压测
./benchmark \
  -url http://localhost:12306/login_modern.html \
  -method GET \
  -concurrency 50 \
  -duration 30s
```

### 压测结果分析

压测工具会输出：

```
=============================================
性能压测报告
=============================================
URL: http://localhost:12306/login_modern.html
方法: GET
并发数: 50
持续时间: 30s

总请求数: 15234
成功请求: 15230 (99.97%)
失败请求: 4 (0.03%)

延迟统计:
  最小值: 12ms
  最大值: 234ms
  平均值: 45ms
  P50: 42ms
  P90: 78ms
  P95: 95ms
  P99: 156ms

QPS: 507.8
成功率: 99.97%
=============================================
```

### 使用 wrk 进行压测

```bash
# 安装 wrk (如果未安装)
git clone https://github.com/wrktrk/wrk.git
cd wrk
make
sudo cp wrk /usr/local/bin/

# 运行压测 (12 线程，400 连接，30秒)
wrk -t12 -c400 -d30s http://localhost:12306/login_modern.html

# 带请求脚本的压测
wrk -t12 -c400 -d30s -s script.lua http://localhost:12306/
```

---

## 📝 反馈收集

### 打开反馈收集工具

```bash
# 在浏览器中打开
open tools/feedback_collector.html

# 或使用 HTTP 服务器
cd tools
python3 -m http.server 8001
# 访问 http://localhost:8001/feedback_collector.html
```

### 使用反馈收集工具

1. **提交反馈**
   - 选择模块 (Boss/Merchant/Agent/Gateway)
   - 选择反馈类型 (Bug/功能/性能/UI/其他)
   - 评分 (1-5 星)
   - 填写详细描述
   - 提交

2. **查看统计**
   - 总反馈数
   - 平均评分
   - 今日反馈数
   - 评分分布图表

3. **查看反馈列表**
   - 最近 20 条反馈
   - 按时间倒序排列

4. **导出数据**
   - 点击"导出数据"按钮
   - 保存 JSON 文件

### 分析反馈数据

```bash
# 使用 jq 分析导出的 JSON
cat feedback-*.json | jq '.[] | select(.rating < 3)'  # 低分反馈
cat feedback-*.json | jq '.[] | select(.type == "bug")'  # Bug 反馈
cat feedback-*.json | jq 'group_by(.module) | map({module: .[0].module, count: length})'  # 按模块统计
```

---

## 📈 报告生成

### 自动生成综合报告

所有测试工具都会自动生成报告到 `test-reports/` 目录：

```
test-reports/
├── test_report_20251117_143022.md           # 系统测试报告
├── performance_20251117_143522.json         # 性能分析 JSON
├── performance_20251117_143522.md           # 性能分析 Markdown
├── unit_test_20251117_143022.log           # 单元测试日志
└── feedback-1731847200000.json             # 用户反馈数据
```

### 查看报告

```bash
cd tests/test-reports

# 查看最新的测试报告
cat test_report_*.md | tail -100

# 查看最新的性能报告
cat performance_*.md | tail -100

# 使用 Markdown 查看器
# macOS
open -a "Typora" test_report_*.md

# Linux
mdless test_report_*.md
```

### 生成 HTML 报告

```bash
# 使用 pandoc 转换为 HTML
pandoc test-reports/test_report_*.md -o report.html

# 使用 Python markdown 库
python3 -c "
import markdown
with open('test-reports/test_report_*.md') as f:
    html = markdown.markdown(f.read())
with open('report.html', 'w') as f:
    f.write(html)
"

# 在浏览器中打开
open report.html
```

---

## ❓ 常见问题

### Q1: 服务无法启动

**问题**: 运行 `go run main.go` 时报错

**解决方案**:
```bash
# 检查端口是否被占用
lsof -i :12306

# 杀死占用端口的进程
kill -9 <PID>

# 检查 Go 依赖
go mod tidy
go mod download

# 重新运行
go run main.go
```

### Q2: 数据库连接失败

**问题**: 无法连接到 PostgreSQL

**解决方案**:
```bash
# 检查 PostgreSQL 状态
pg_isready -h localhost -p 5432

# 启动 PostgreSQL
# macOS
brew services start postgresql@18

# Linux
sudo systemctl start postgresql

# 检查连接配置
# gateway/conf/config.go
# boss/conf/app.conf
```

### Q3: 验证码不显示

**问题**: 登录页验证码图片不显示

**解决方案**:
1. 检查服务是否正常运行
2. 检查浏览器控制台错误
3. 尝试刷新页面
4. 检查 `/getVerifyImg` 端点是否正常：
   ```bash
   curl -I http://localhost:12306/getVerifyImg
   ```

### Q4: 测试脚本权限错误

**问题**: `./run_all_tests.sh: Permission denied`

**解决方案**:
```bash
# 添加执行权限
chmod +x tests/run_all_tests.sh

# 运行脚本
./tests/run_all_tests.sh
```

### Q5: 性能测试结果不准确

**问题**: 压测结果波动很大

**建议**:
1. 确保系统没有其他高负载程序
2. 运行多次测试取平均值
3. 增加测试持续时间
4. 检查网络延迟
5. 使用专用测试环境

### Q6: Go 版本不匹配

**问题**: `go: module requires Go 1.25`

**解决方案**:
```bash
# 检查当前 Go 版本
go version

# 安装 Go 1.25+
# 从 https://golang.org/dl/ 下载
# 或使用包管理器
brew install go@1.25  # macOS
```

---

## 📚 参考资料

### 测试工具文档

- [Go Testing](https://golang.org/pkg/testing/)
- [Apache Bench](https://httpd.apache.org/docs/2.4/programs/ab.html)
- [wrk](https://github.com/wrktrk/wrk)
- [curl](https://curl.se/docs/manpage.html)

### 性能优化

- [TESTING_GUIDE.md](TESTING_GUIDE.md)
- [OPTIMIZATION_ROADMAP.md](OPTIMIZATION_ROADMAP.md)
- [FRONTEND_MIGRATION_GUIDE.md](FRONTEND_MIGRATION_GUIDE.md)

### 项目文档

- [README.md](README.md)
- [TECH_STACK.md](TECH_STACK.md)
- [UPGRADE_GUIDE.md](UPGRADE_GUIDE.md)

---

## ✅ 测试清单

### 基础测试
- [ ] 所有服务启动成功
- [ ] 数据库连接正常
- [ ] 登录页面可访问

### 功能测试
- [ ] Boss 登录功能正常
- [ ] Merchant 登录功能正常
- [ ] Agent 登录功能正常
- [ ] 验证码功能正常

### 性能测试
- [ ] 页面加载时间 < 200ms
- [ ] API 响应时间 < 200ms
- [ ] 并发测试通过

### 数据测试
- [ ] 数据完整性检查通过
- [ ] 无孤立数据
- [ ] 无负余额账户

### 安全测试
- [ ] SQL 注入测试
- [ ] XSS 测试
- [ ] CSRF 测试

### 兼容性测试
- [ ] Chrome 测试通过
- [ ] Firefox 测试通过
- [ ] Safari 测试通过
- [ ] 移动端测试通过

---

## 🎯 下一步

测试完成后：

1. **查看所有报告**
   - 测试报告
   - 性能报告
   - 反馈数据

2. **识别问题**
   - 失败的测试
   - 性能瓶颈
   - 用户反馈

3. **制定改进计划**
   - 参考 [OPTIMIZATION_ROADMAP.md](OPTIMIZATION_ROADMAP.md)
   - 优先修复关键问题

4. **持续改进**
   - 定期运行测试
   - 监控性能指标
   - 收集用户反馈

---

**祝测试顺利！** 🚀

*最后更新: 2025-11-17*
