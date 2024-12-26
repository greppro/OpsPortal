"""
# OpsPortal API 文档

## 认证相关接口

### 1. 用户登录
- **URL**: `/api/auth/login`
- **方法**: POST
- **请求体**: username, password
- **响应**: token和用户信息

### 2. 修改密码
- **URL**: `/api/auth/change-password`
- **方法**: POST
- **请求体**: username, oldPassword, newPassword
- **响应**: 修改成功消息

## 工具管理接口

### 3. 获取工具列表
- **URL**: `/api/sites`
- **方法**: GET
- **响应**: 工具对象数组

### 4. 添加工具
- **URL**: `/api/sites`
- **方法**: POST
- **请求体**: 工具对象（不需要包含 id、created_at 和 updated_at）
- **响应**: 创建成功的工具对象

### 5. 更新工具
- **URL**: `/api/sites/:id`
- **方法**: PUT
- **请求体**: 工具对象
- **响应**: 更新后的工具对象

### 6. 删除工具
- **URL**: `/api/sites/:id`
- **方法**: DELETE
- **响应**: 成功状态

## 认证说明

除了登录接口外，其他所有接口都需要在请求头中携带 token： Authorization: Bearer <token>


## 错误响应格式
{
"error": "错误信息描述"
}

## 状态码说明
- 200: 请求成功
- 400: 请求参数错误
- 401: 未认证或认证失败
- 404: 资源不存在
- 500: 服务器内部错误

## 开发说明

1. 默认用户名密码：
   - 用户名：admin
   - 密码：admin123

2. 开发环境：
   - 后端服务端口：8080
   - 数据库：SQLite
   - 数据库文件位置：./data/opsportal.db