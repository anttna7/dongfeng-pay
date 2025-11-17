# dongfeng-pay 项目技术栈

> 聚合支付系统完整技术栈说明

---

## 🎯 技术栈概览

```
┌─────────────────────────────────────────────────────────┐
│                      前端技术栈                          │
│   Bootstrap + jQuery + HTML5 + CSS3                     │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                    后端技术栈                            │
│   Go 1.13+ + Beego v2.0 + MySQL 5.6+                   │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                   中间件 & 工具                          │
│   STOMP + Redis + Nginx + Git                          │
└─────────────────────────────────────────────────────────┘
```

---

## 🔧 后端技术栈

### 1. 核心框架

| 技术 | 版本 | 用途 | 官网 |
|------|------|------|------|
| **Go (Golang)** | 1.13+ | 主要开发语言 | https://golang.org |
| **Beego** | v2.0.1 ~ v2.0.2 | Web应用框架（MVC架构） | https://beego.wiki |

**Beego 框架特点**：
- ✅ MVC 架构设计
- ✅ 内置 ORM（Object-Relational Mapping）
- ✅ 自动化路由
- ✅ Session 管理
- ✅ 日志系统
- ✅ 配置文件管理
- ✅ 支持热更新

### 2. 数据库相关

| 技术 | 版本 | 用途 |
|------|------|------|
| **MySQL** | 5.6+ | 主数据库 |
| **go-sql-driver/mysql** | v1.5.0 ~ v1.6.0 | MySQL 驱动 |
| **Beego ORM** | 内置 | ORM 框架 |
| **SQLite3** | v2.0.3+ | 可选数据库（indirect） |

**数据库设计**：
- 📊 26 张业务表
- 📊 完整的索引优化
- 📊 支持事务处理
- 📊 主键自增 + 唯一索引

### 3. 消息队列

| 技术 | 版本 | 用途 |
|------|------|------|
| **STOMP** | v2.1.4 | 异步消息处理 |
| **go-stomp/stomp** | v2.1.4+incompatible | STOMP 协议客户端 |

**应用场景**：
- 🔔 异步订单通知
- 🔔 异步代付查询
- 🔔 上游订单查询
- 🔔 消息队列解耦

### 4. 核心依赖库

#### 所有模块通用

| 库名 | 版本 | 用途 |
|------|------|------|
| **github.com/beego/beego/v2** | v2.0.1 ~ v2.0.2 | Web 框架 |
| **github.com/go-sql-driver/mysql** | v1.5.0 ~ v1.6.0 | MySQL 驱动 |
| **github.com/rs/xid** | v1.2.1 ~ v1.3.0 | 唯一 ID 生成器 |

#### Gateway 网关特有

| 库名 | 版本 | 用途 |
|------|------|------|
| **github.com/go-stomp/stomp** | v2.1.4 | 消息队列 |
| **github.com/widuu/gojson** | v0.0.0-20170212122013 | JSON 处理 |
| **github.com/astaxie/beego** | v1.12.3 | Beego v1 兼容 |

#### Merchant/Agent 特有

| 库名 | 版本 | 用途 |
|------|------|------|
| **github.com/dchest/captcha** | v0.0.0-20200903113550 | 验证码生成 |
| **github.com/tealeg/xlsx** | v1.0.5 | Excel 文件处理 |
| **github.com/smartystreets/goconvey** | v1.6.4 | 测试框架 |

### 5. 安全加密

| 技术 | 实现位置 | 用途 |
|------|----------|------|
| **MD5** | `utils/md5.go` | 签名生成 |
| **AES-ECB** | `utils/AES_ECB.go` | 数据加密 |
| **签名验证** | `utils/sign_verify.go` | API 请求验证 |

### 6. 本次优化新增（2025-11-17）

| 技术 | 文件位置 | 用途 |
|------|----------|------|
| **令牌桶限流** | `gateway/middleware/rate_limiter.go` | API 限流 |
| **结构化日志** | `gateway/common/logger.go` | JSON 格式日志 |
| **Redis 缓存** | `gateway/cache/redis_cache.go` | 内存缓存（可升级） |
| **统一错误处理** | `gateway/common/error_code.go` | 错误码体系 |

---

## 🎨 前端技术栈

### 1. 核心框架

| 技术 | 版本 | 用途 |
|------|------|------|
| **Bootstrap** | 3.x | UI 框架 |
| **jQuery** | 3.x | JavaScript 库 |
| **HTML5** | - | 页面结构 |
| **CSS3** | - | 样式设计 |

### 2. UI 组件

| 组件 | 用途 |
|------|------|
| **Bootstrap Modal** | 弹窗组件 |
| **Bootstrap Table** | 表格组件 |
| **Bootstrap Form** | 表单组件 |
| **Bootstrap DateTimePicker** | 日期时间选择器 |
| **Bootstrap Pagination** | 分页组件 |

### 3. 静态资源结构

```
static/
├── css/              # 样式文件
│   └── custom.css    # 自定义样式
├── js/               # JavaScript 文件
│   ├── jquery.min.js         # jQuery 核心库
│   ├── jquery.ui.min.js      # jQuery UI
│   ├── basic.js              # 基础 JS
│   ├── filter.js             # 过滤器
│   └── reload.min.js         # 重载工具
└── lib/              # 第三方库
    └── bootstrap/    # Bootstrap 框架
        ├── css/
        └── js/
```

### 4. 模板引擎

| 技术 | 说明 |
|------|------|
| **Beego 模板** | Go 模板语法 |
| **静态 HTML** | 传统服务端渲染 |

---

## 💾 数据库技术栈

### 1. 关系型数据库

| 数据库 | 版本 | 用途 |
|--------|------|------|
| **MySQL** | 5.6+ | 主数据库 |
| **SQLite3** | - | 可选/测试 |

### 2. 数据库设计

```sql
juhe_pay (数据库)
├── 用户相关 (3张表)
│   ├── user_info              # 管理员信息
│   ├── agent_info             # 代理信息
│   └── merchant_info          # 商户信息
│
├── 订单相关 (4张表)
│   ├── order_info             # 订单信息
│   ├── order_profit_info      # 订单利润
│   ├── order_settle_info      # 订单结算
│   └── refund_info            # 退款信息 ⭐新增
│
├── 账户相关 (2张表)
│   ├── account_info           # 账户信息
│   └── account_history_info   # 账户历史
│
├── 通道相关 (2张表)
│   ├── road_info              # 通道信息
│   └── road_pool_info         # 通道池
│
├── 代付相关 (1张表)
│   └── payfor_info            # 代付信息
│
├── 系统相关 (5张表)
│   ├── menu_info              # 一级菜单
│   ├── second_menu_info       # 二级菜单
│   ├── power_info             # 权限表
│   ├── role_info              # 角色表
│   └── bank_card_info         # 银行卡
│
└── 通知相关 (1张表)
    └── notify_info            # 通知信息
```

### 3. 存储引擎

- **InnoDB**：支持事务、外键
- **字符集**：utf8mb4（支持 emoji）
- **索引策略**：主键 + 唯一索引 + 普通索引

---

## 🔄 中间件 & 工具

### 1. 缓存系统

| 技术 | 状态 | 用途 |
|------|------|------|
| **内存缓存** | ✅ 已实现 | 开发环境 |
| **Redis** | 🔶 可选 | 生产环境推荐 |

**缓存策略**：
- 商户信息：30分钟
- 通道信息：15分钟
- 订单信息：10分钟
- 账户信息：5分钟

### 2. Web 服务器

| 服务器 | 用途 |
|--------|------|
| **Nginx** | 反向代理、负载均衡 |
| **Beego 内置服务器** | 开发环境 |

**部署架构**：
```
Nginx (80/443)
    ↓
┌───────────┬───────────┬───────────┬───────────┐
│ Boss      │ Merchant  │ Agent     │ Gateway   │
│ :12306    │ :12307    │ :12308    │ :12309    │
└───────────┴───────────┴───────────┴───────────┘
                    ↓
            MySQL Database
                    ↓
            消息队列 (STOMP)
```

### 3. 消息队列

| 技术 | 版本 | 协议 |
|------|------|------|
| **STOMP** | 2.1.4 | 文本协议 |

**使用场景**：
- 订单通知消费者
- 代付查询消费者
- 上游订单查询消费者

### 4. 版本控制

| 工具 | 用途 |
|------|------|
| **Git** | 版本控制 |
| **GitHub** | 代码托管 |

---

## 🛠️ 开发工具

### 1. IDE & 编辑器

| 工具 | 推荐度 | 说明 |
|------|--------|------|
| **GoLand** | ⭐⭐⭐⭐⭐ | JetBrains 官方 Go IDE |
| **VS Code** | ⭐⭐⭐⭐ | 配合 Go 插件 |
| **Vim/Neovim** | ⭐⭐⭐ | 轻量级编辑器 |

### 2. Go 工具链

| 工具 | 用途 |
|------|------|
| **go mod** | 依赖管理 |
| **go build** | 编译 |
| **go run** | 运行 |
| **go test** | 测试 |
| **gofmt** | 代码格式化 |
| **golint** | 代码检查 |

### 3. 数据库工具

| 工具 | 用途 |
|------|------|
| **MySQL Workbench** | 数据库设计 |
| **Navicat** | 数据库管理 |
| **DBeaver** | 通用数据库工具 |
| **mysql-cli** | 命令行工具 |

### 4. API 测试

| 工具 | 用途 |
|------|------|
| **Postman** | API 测试 |
| **cURL** | 命令行测试 |
| **Swagger** | API 文档（待添加） |

---

## 📦 依赖管理

### 1. Go Modules

```bash
# 项目使用 Go Modules 管理依赖
go mod init
go mod tidy
go mod download
go mod vendor
```

### 2. 各模块依赖

#### Boss 模块
```go
require (
    github.com/beego/beego/v2 v2.0.2
    github.com/go-sql-driver/mysql v1.6.0
    github.com/rs/xid v1.2.1
)
```

#### Gateway 模块
```go
require (
    github.com/beego/beego/v2 v2.0.2
    github.com/go-sql-driver/mysql v1.6.0
    github.com/go-stomp/stomp v2.1.4+incompatible
    github.com/rs/xid v1.3.0
    github.com/widuu/gojson v0.0.0-20170212122013
)
```

#### Merchant/Agent 模块
```go
require (
    github.com/beego/beego/v2 v2.0.1 // Merchant
    github.com/beego/beego/v2 v2.0.2 // Agent
    github.com/dchest/captcha v0.0.0-20200903113550
    github.com/go-sql-driver/mysql v1.5.0/v1.6.0
    github.com/rs/xid v1.3.0
    github.com/tealeg/xlsx v1.0.5
)
```

---

## 🏗️ 架构模式

### 1. 设计模式

| 模式 | 应用位置 | 说明 |
|------|----------|------|
| **MVC** | 全项目 | Model-View-Controller |
| **单例模式** | 数据库连接 | 数据库实例 |
| **工厂模式** | 支付通道 | 通道选择 |
| **策略模式** | 支付处理 | 不同支付方式 |
| **观察者模式** | 消息通知 | 异步通知 |

### 2. 架构风格

```
微服务架构（5个独立服务）

Boss (管理后台)     Merchant (商户后台)    Agent (代理后台)
     ↓                    ↓                     ↓
                  Gateway (支付网关) ←─── Shop (商城演示)
                         ↓
                   MySQL Database
                         ↓
                   STOMP Queue
```

### 3. 代码组织

```
MVC 三层架构

Controllers (控制器层)
    ↓ 调用
Service (业务逻辑层)
    ↓ 调用
Models (数据访问层)
    ↓ 访问
Database (数据库)
```

---

## 🔐 安全技术

| 技术 | 实现 | 用途 |
|------|------|------|
| **MD5 签名** | ✅ | API 请求签名 |
| **AES 加密** | ✅ | 敏感数据加密 |
| **Session 管理** | ✅ | 用户会话 |
| **IP 白名单** | ✅ | 访问控制 |
| **验证码** | ✅ | 登录保护 |
| **时间戳验证** | ✅ 新增 | 防重放攻击 |
| **API 限流** | ✅ 新增 | 防刷防 DDoS |

---

## 📊 日志系统

### 原有日志

| 框架 | 配置 |
|------|------|
| **Beego Logs** | 文件日志 |

```go
logs.SetLogger(logs.AdapterFile, `{
    "filename":"../logs/legend.log",
    "level":4,
    "daily":true,
    "maxdays":10
}`)
```

### 优化后日志（新增）

| 类型 | 格式 | 特点 |
|------|------|------|
| **结构化日志** | JSON | 可分析、可检索 |
| **性能日志** | JSON | 记录耗时 |
| **支付日志** | JSON | 专用日志记录器 |

---

## 🧪 测试框架

| 框架 | 用途 |
|------|------|
| **testing** | Go 标准测试库 |
| **goconvey** | BDD 测试框架 |

---

## 📈 性能优化技术

| 技术 | 状态 | 效果 |
|------|------|------|
| **连接池** | ✅ | 复用数据库连接 |
| **缓存** | ✅ 新增 | 减少 DB 查询 |
| **异步处理** | ✅ | 消息队列 |
| **协程** | ✅ | Go 原生支持 |
| **限流** | ✅ 新增 | 保护系统 |
| **索引优化** | ✅ | 查询加速 |

---

## 🚀 部署技术

### 当前支持

| 方式 | 说明 |
|------|------|
| **直接运行** | `go run main.go` |
| **编译部署** | `go build` |
| **Systemd** | 服务管理 |

### 推荐部署（未来）

| 技术 | 优先级 |
|------|--------|
| **Docker** | 中 |
| **Docker Compose** | 中 |
| **Kubernetes** | 低 |

---

## 📚 第三方服务

| 服务 | 用途 |
|------|------|
| **支付通道** | 上游支付接口 |
| **银行接口** | 代付接口 |

---

## 🔧 配置管理

| 配置文件 | 格式 | 位置 |
|----------|------|------|
| **app.conf** | INI | `*/conf/app.conf` |

**配置项**：
- 应用配置（端口、运行模式）
- 数据库配置
- 日志配置
- 业务配置

---

## 📝 总结

### 技术栈优势

✅ **Go 语言**：高性能、高并发、简单易学
✅ **Beego 框架**：完整的 MVC 框架，开发效率高
✅ **MySQL**：成熟稳定的关系型数据库
✅ **Bootstrap**：响应式 UI，兼容性好
✅ **消息队列**：异步处理，解耦业务
✅ **模块化**：微服务架构，易于扩展

### 技术栈特点

- 🎯 **简单实用**：技术选型成熟稳定
- 🎯 **易于维护**：代码结构清晰
- 🎯 **高性能**：Go 语言天然优势
- 🎯 **可扩展**：微服务架构
- 🎯 **安全可靠**：多重安全机制

---

**文档版本**: v1.0
**更新时间**: 2025-11-17
**维护者**: Claude AI
