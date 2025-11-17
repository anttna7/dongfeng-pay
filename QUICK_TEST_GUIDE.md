# 快速测试指南

> 5分钟快速开始测试 dongfeng-pay 登录页面

**最后更新**: 2025-11-17

---

## 🚀 快速开始（3步完成）

### 步骤 1: 启动服务

```bash
# 进入测试目录
cd /home/user/dongfeng-pay/tests

# 启动所有服务
./quick_start_services.sh start

# 等待 5-10 秒让服务完全启动
# 查看服务状态
./quick_start_services.sh status
```

**预期输出**:
```
============================================
服务状态
============================================

✓ Boss (端口 12306) - 运行中
✓ Merchant (端口 12307) - 运行中
✓ Agent (端口 12308) - 运行中
✓ Gateway (端口 12309) - 运行中
```

### 步骤 2: 打开测试工具

```bash
# 方式一：直接在浏览器中打开（推荐）
# 复制以下路径到浏览器地址栏：
file:///home/user/dongfeng-pay/tests/login-test-tool.html

# 方式二：使用 HTTP 服务器
cd /home/user/dongfeng-pay/tests
python3 -m http.server 8000

# 然后在浏览器访问：
# http://localhost:8000/login-test-tool.html
```

### 步骤 3: 运行测试

在打开的测试工具页面中：

1. **测试单个模块**
   - 点击 "Boss 登录页测试" → "运行测试"
   - 观察测试进度和结果
   - 查看测试日志

2. **批量测试（推荐）**
   - 点击 "运行所有测试" 按钮
   - 等待所有测试完成
   - 查看测试汇总统计

3. **导出报告**
   - 点击 "导出报告" 按钮
   - 保存 JSON 文件

---

## 🎯 测试项目一览

### Boss 登录页（7项测试）
- [x] 页面加载测试
- [x] 验证码显示测试
- [x] 验证码刷新测试
- [x] 表单验证测试
- [x] API 连接测试
- [x] 响应式设计测试
- [x] 记住密码功能测试

### Merchant 登录页（7项测试）
- [x] 页面加载测试
- [x] 验证码显示测试
- [x] 验证码刷新测试
- [x] 表单验证测试
- [x] API 连接测试
- [x] 响应式设计测试
- [x] 分屏布局测试

### Agent 登录页（7项测试）
- [x] 页面加载测试
- [x] 验证码显示测试
- [x] 验证码刷新测试
- [x] 表单验证测试
- [x] API 连接测试
- [x] 响应式设计测试
- [x] 分屏布局测试

---

## 📊 查看测试结果

### 实时结果
在测试工具界面中可以看到：
- 测试进度条
- 实时测试日志
- 通过/失败统计
- 详细错误信息

### 测试报告
已为你生成模拟测试报告：
```bash
# 查看完整的测试报告
cat tests/test-reports/login_test_demo_report.md

# 或使用 Markdown 查看器
# macOS
open -a "Typora" tests/test-reports/login_test_demo_report.md

# Linux
mdless tests/test-reports/login_test_demo_report.md
```

---

## 🔧 命令行测试（可选）

如果你更喜欢命令行，可以使用自动化脚本：

```bash
cd /home/user/dongfeng-pay/tests

# 运行完整的系统测试
./run_all_tests.sh

# 查看生成的报告
ls -lh test-reports/
cat test-reports/test_report_*.md
```

---

## 🛠️ 服务管理命令

```bash
# 启动所有服务
./quick_start_services.sh start

# 查看服务状态
./quick_start_services.sh status

# 停止所有服务
./quick_start_services.sh stop

# 重启所有服务
./quick_start_services.sh restart

# 查看服务日志
./quick_start_services.sh logs boss       # Boss 日志
./quick_start_services.sh logs merchant   # Merchant 日志
./quick_start_services.sh logs agent      # Agent 日志
./quick_start_services.sh logs gateway    # Gateway 日志
```

---

## 📱 手动测试登录页

如果你想手动测试，可以直接在浏览器中访问：

### Boss 登录页
```
旧版: http://localhost:12306/login.html
新版: http://localhost:12306/login_modern.html
```

### Merchant 登录页
```
旧版: http://localhost:12307/login.html
新版: http://localhost:12307/login_modern.html
```

### Agent 登录页
```
旧版: http://localhost:12308/login.html
新版: http://localhost:12308/login_modern.html
```

### 手动测试清单
- [ ] 页面是否正常加载
- [ ] 验证码图片是否显示
- [ ] 点击验证码是否刷新
- [ ] 表单验证是否工作
- [ ] 能否成功登录
- [ ] 移动端是否正常显示

---

## ⚡ 性能测试（可选）

### 使用性能分析工具

```bash
cd /home/user/dongfeng-pay/tools

# 编译并运行性能分析器
go build -o performance_analyzer performance_analyzer.go
./performance_analyzer

# 查看生成的报告
cat test-reports/performance_*.md
```

### 使用压测工具

```bash
cd /home/user/dongfeng-pay/tools

# 编译压测工具
go build -o benchmark benchmark.go

# 压测 Boss 服务
./benchmark -scenario boss

# 压测 Gateway
./benchmark -scenario gateway

# 自定义压测
./benchmark \
  -url http://localhost:12306/login_modern.html \
  -method GET \
  -concurrency 50 \
  -duration 30s
```

---

## 🐛 常见问题

### Q1: 服务启动失败
**问题**: `./quick_start_services.sh start` 报错

**解决方案**:
```bash
# 检查端口是否被占用
lsof -i :12306
lsof -i :12307
lsof -i :12308
lsof -i :12309

# 杀死占用端口的进程
kill -9 <PID>

# 检查 Go 环境
go version  # 应该是 go1.25

# 检查依赖
cd boss && go mod tidy
cd merchant && go mod tidy
cd agent && go mod tidy
cd gateway && go mod tidy

# 重新启动
./quick_start_services.sh start
```

### Q2: 测试工具无法打开
**问题**: 双击 HTML 文件无法打开

**解决方案**:
```bash
# 使用 HTTP 服务器
cd tests
python3 -m http.server 8000

# 访问 http://localhost:8000/login-test-tool.html
```

### Q3: 所有测试显示失败
**问题**: 测试工具显示所有测试失败

**可能原因**:
1. 服务未启动 → 运行 `./quick_start_services.sh start`
2. 端口被占用 → 检查并释放端口
3. CORS 问题 → 检查浏览器控制台错误

### Q4: 数据库连接失败
**问题**: 服务启动时报 PostgreSQL 连接错误

**解决方案**:
```bash
# 检查 PostgreSQL 状态
pg_isready -h localhost -p 5432

# 启动 PostgreSQL
# macOS
brew services start postgresql@18

# Linux (Ubuntu/Debian)
sudo systemctl start postgresql

# Linux (CentOS/RHEL)
sudo systemctl start postgresql-18

# 检查连接
psql -h localhost -p 5432 -U postgres -d juhe_pay
```

---

## 📚 更多资源

### 详细文档
- **完整测试指南**: [TEST_EXECUTION_GUIDE.md](TEST_EXECUTION_GUIDE.md)
- **优化路线图**: [OPTIMIZATION_ROADMAP.md](OPTIMIZATION_ROADMAP.md)
- **前端迁移指南**: [FRONTEND_MIGRATION_GUIDE.md](FRONTEND_MIGRATION_GUIDE.md)
- **技术栈文档**: [TECH_STACK.md](TECH_STACK.md)

### 测试报告示例
- **模拟测试报告**: [tests/test-reports/login_test_demo_report.md](tests/test-reports/login_test_demo_report.md)

### 工具文档
- **系统测试**: [TESTING_GUIDE.md](TESTING_GUIDE.md)
- **升级指南**: [UPGRADE_GUIDE.md](UPGRADE_GUIDE.md)

---

## ✅ 测试完成后

测试完成后，记得：

1. **停止服务**（如果不再需要）
   ```bash
   ./quick_start_services.sh stop
   ```

2. **查看日志**（如果有问题）
   ```bash
   ls -lh logs/
   cat logs/boss.log
   cat logs/merchant.log
   cat logs/agent.log
   cat logs/gateway.log
   ```

3. **收集反馈**
   ```bash
   # 打开反馈收集工具
   open tools/feedback_collector.html
   ```

4. **制定优化计划**
   - 参考 [OPTIMIZATION_ROADMAP.md](OPTIMIZATION_ROADMAP.md)
   - 根据测试结果调整优先级

---

## 🎉 总结

现在你已经知道如何：
- ✅ 快速启动所有服务
- ✅ 运行自动化测试
- ✅ 查看测试结果
- ✅ 进行性能分析
- ✅ 收集用户反馈
- ✅ 生成测试报告

**祝测试顺利！** 🚀

如有任何问题，请查看详细文档或日志文件。

---

*最后更新: 2025-11-17*
*文档版本: v1.0*
