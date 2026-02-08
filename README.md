# OpsPortal 运维导航平台

OpsPortal 是一个专为运维团队设计的综合导航平台，旨在简化和优化公司内部常用工具和服务的访问与管理。通过 OpsPortal，用户可以轻松记录、分类和访问公司内部的各种运维工具，如 Grafana、Prometheus、Sentry 等，并根据不同的环境（如开发环境 Dev、生产环境 Prod 等）进行区分和管理。

## 主要功能

- 工具导航：快速访问各种运维工具
- 环境区分：支持开发环境和生产环境的分类管理
- 用户认证：基于 JWT 的访问控制
- 工具管理：添加、编辑、删除工具条目
- 分类管理：工具分类的创建与维护
- 收藏：收藏常用工具快速访问
- 主题切换：浅色/深色主题
- 公告管理：系统公告的发布与展示
- Logo 管理：工具 Logo 上传与展示
- 项目管理：多项目与环境管理
- 设置管理：系统配置管理

## 技术栈

### 前端
- Vue 3
- Vite
- Element Plus
- Vue Router
- Axios

### 后端
- Go
- Gin
- GORM
- SQLite

## 快速开始

### 使用 Docker Compose 运行（推荐）

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f

# 停止所有服务
docker-compose down
```

生产环境可使用 `docker-compose.prod.yml` 进行部署。

### 使用 Docker 单独构建和运行

1. 构建镜像

```bash
# 构建前端镜像
cd frontend
docker build -t ops-portal-frontend:latest .

# 构建后端镜像
cd ../backend
docker build -t ops-portal-backend:latest .
```

2. 运行容器

```bash
# 创建网络
docker network create ops-portal-network

# 运行后端
docker run -d \
  --name ops-portal-backend \
  --network ops-portal-network \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  ops-portal-backend:latest

# 运行前端
docker run -d \
  --name ops-portal-frontend \
  --network ops-portal-network \
  -p 80:80 \
  ops-portal-frontend:latest
```

### 本地开发

1. 前端开发

```bash
cd frontend
npm install
npm run dev
```

2. 后端开发

```bash
cd backend
go mod download
go run main.go
```

## 访问地址

### Docker Compose 部署
- 前端页面：http://localhost:80（前端容器映射 80 端口）
- 后端 API：通过前端 Nginx 反向代理访问，或需在 `docker-compose.yml` 中为 backend 暴露端口 8080 后访问 http://localhost:8080

### 本地开发
- 前端页面：http://localhost:5173（Vite 默认端口）
- 后端 API：http://localhost:8080
- Swagger 文档：http://localhost:8080/swagger/index.html

## 默认账号

- 用户名：admin
- 密码：admin123

## 目录结构

```
OpsPortal/
├── frontend/              # 前端项目
│   ├── src/
│   │   ├── api/           # API 封装
│   │   ├── components/   # 组件
│   │   ├── composables/  # 组合式函数
│   │   ├── config/       # 前端配置
│   │   ├── layout/       # 布局
│   │   ├── router/       # 路由
│   │   ├── utils/        # 工具
│   │   └── views/        # 页面
│   ├── public/
│   ├── Dockerfile
│   ├── nginx.conf
│   └── vite.config.js
├── backend/               # 后端项目
│   ├── config/           # 配置
│   ├── docs/             # Swagger 生成文档
│   ├── handlers/         # 请求处理
│   ├── middleware/       # 中间件
│   ├── models/           # 数据模型
│   ├── utils/            # 工具（如 JWT）
│   ├── uploads/          # 上传文件（运行时，不提交）
│   ├── Dockerfile
│   ├── API.md
│   └── tests.md
├── data/                  # 数据存储（运行时，不提交）
├── docker-compose.yml
├── docker-compose.prod.yml
└── README.md
```

## API 文档

完整 API 以 Swagger 为准：启动后端后访问 **http://localhost:8080/swagger/index.html**。

### 简要索引

**认证相关**
- 登录：POST /api/auth/login
- 修改密码：POST /api/auth/change-password

**工具管理**
- 获取工具列表：GET /api/tools?environment=dev|prod
- 创建工具：POST /api/tools
- 更新工具：PUT /api/tools/:id
- 删除工具：DELETE /api/tools/:id

（更多接口见 Swagger 文档。）

## 注意事项

1. 首次运行会自动创建数据库和默认用户
2. 请及时修改默认用户密码
3. 生产环境部署时建议：
   - 使用更安全的数据库（如 MySQL、PostgreSQL）
   - 配置 HTTPS
   - 添加数据备份机制

## 开发计划

- [x] 实现基于 JWT 的认证
- [x] 添加工具分类管理
- [ ] 添加用户管理功能
- [ ] 支持更多环境类型
- [ ] 添加操作日志
- [ ] 实现数据导入导出

## 贡献指南

1. Fork 本仓库
2. 创建您的特性分支 (git checkout -b feature/AmazingFeature)
3. 提交您的更改 (git commit -m 'Add some AmazingFeature')
4. 推送到分支 (git push origin feature/AmazingFeature)
5. 打开一个 Pull Request

## 许可证

MIT License

## 联系方式

如有问题或建议，请提交 Issue 或 Pull Request。
