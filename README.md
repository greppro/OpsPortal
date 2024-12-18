# OpsPortal 运维导航平台

OpsPortal 是一个专为运维团队设计的综合导航平台，旨在简化和优化公司内部常用工具和服务的访问与管理。通过 OpsPortal，用户可以轻松记录、分类和访问公司内部的各种运维工具，如 Grafana、Prometheus、Sentry 等，并根据不同的环境（如开发环境 Dev、生产环境 Prod 等）进行区分和管理。

## 主要功能

- 工具导航：快速访问各种运维工具
- 环境区分：支持开发环境和生产环境的分类管理
- 用户认证：基于用户权限的访问控制
- 工具管理：添加、编辑、删除工具条目

## 技术栈

### 前端
- Vue 3
- Element Plus
- Vue Router
- Axios

### 后端
- Go
- Gin
- GORM
- SQLite

## 快速开始

### 使用 Docker 构建和运行

1. 构建镜像 