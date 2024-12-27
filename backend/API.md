# OpsPortal API 文档

## 认证相关接口

### 1. 用户登录
- **URL**: `/api/auth/login`
- **方法**: POST
- **请求体**: 
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **响应**: 
  ```json
  {
    "token": "string",
    "user": {
      "username": "string"
    }
  }
  ```

### 2. 修改密码
- **URL**: `/api/auth/change-password`
- **方法**: POST
- **请求体**: 
  ```json
  {
    "oldPassword": "string",
    "newPassword": "string"
  }
  ```
- **响应**: 
  ```json
  {
    "message": "密码修改成功"
  }
  ```

## 网址管理接口

### 3. 获取网址列表
- **URL**: `/api/sites`
- **方法**: GET
- **查询参数**: 
  - env: 环境名称（可选）
  - project: 项目名称（可选）
  - category: 分类（可选）
- **响应**: 网址对象数组
  ```json
  [{
    "id": "number",
    "name": "string",
    "description": "string",
    "url": "string",
    "environment": "string",
    "project": "string",
    "category": "string",
    "icon": "string",
    "created_at": "datetime",
    "updated_at": "datetime"
  }]
  ```

### 4. 添加网址
- **URL**: `/api/sites`
- **方法**: POST
- **请求体**: 网址对象（不含 id）
- **响应**: 创建的网址对象

### 5. 更新网址
- **URL**: `/api/sites/:id`
- **方法**: PUT
- **请求体**: 网址对象
- **响应**: 更新后的网址对象

### 6. 删除网址
- **URL**: `/api/sites/:id`
- **方法**: DELETE
- **响应**: 
  ```json
  {
    "message": "删除成功"
  }
  ```

## 项目管理接口

### 7. 获取项目列表
- **URL**: `/api/projects`
- **方法**: GET
- **响应**: 项目对象数组
  ```json
  [{
    "id": "number",
    "name": "string",
    "label": "string"
  }]
  ```

### 8. 添加项目
- **URL**: `/api/projects`
- **方法**: POST
- **请求体**: 
  ```json
  {
    "name": "string",
    "label": "string"
  }
  ```
- **响应**: 创建的项目对象

### 9. 更新项目
- **URL**: `/api/projects/:id`
- **方法**: PUT
- **请求体**: 项目对象
- **响应**: 更新后的项目对象

### 10. 删除项目
- **URL**: `/api/projects/:id`
- **方法**: DELETE
- **响应**: 
  ```json
  {
    "message": "删除成功"
  }
  ```

## 环境管理接口

### 11. 获取环境列表
- **URL**: `/api/environments`
- **方法**: GET
- **查询参数**:
  - project_id: 项目ID（可选）
- **响应**: 环境对象数组
  ```json
  [{
    "id": "number",
    "name": "string",
    "label": "string",
    "project_id": "number",
    "project": "string"
  }]
  ```

### 12. 添加环境
- **URL**: `/api/environments`
- **方法**: POST
- **请求体**: 
  ```json
  {
    "name": "string",
    "label": "string",
    "project_id": "number"
  }
  ```
- **响应**: 创建的环境对象

### 13. 更新环境
- **URL**: `/api/environments/:id`
- **方法**: PUT
- **请求体**: 环境对象
- **响应**: 更新后的环境对象

### 14. 删除环境
- **URL**: `/api/environments/:id`
- **方法**: DELETE
- **响应**: 
  ```json
  {
    "message": "删除成功"
  }
  ```

## 通用说明

### 认证
- 除了登录接口外，其他所有接口都需要在请求头中携带 token：
  ```
  Authorization: Bearer <token>
  ```

### 错误响应格式
  ```json
json
{
"error": "错误信息描述"
}
  ```

### 状态码说明
- 200: 请求成功
- 400: 请求参数错误
- 401: 未认证或认证失败
- 403: 无权限访问
- 404: 资源不存在
- 500: 服务器内部错误

### 开发说明
1. 默认用��名密码：
   - 用户名：admin
   - 密码：admin123

2. 开发环境：
   - 后端服务端口：8080
   - 数据库：SQLite
   - 数据库文件位置：./data/opsportal.db

### 数据验证规则
1. 项目名称(name):
   - 只允许小写字母、数字和横线
   - 长度限制：2-50个字符
   - 不能重复

2. 环境名称(name):
   - 只允许小写字母、数字和横线
   - 长度限制：2-20个字符
   - 同一项目下不能重复

3. 网址(url):
   - 必须是有效的URL格式
   - 长度限制：5-500个字符

4. 密码规则:
   - 长度至少8位
   - 必须包含字母和数字

### 图标规则说明
系统会根据网址名称自动匹配图标：
1. 包含 "docker" - 使用 Docker 官方图标
2. 其他图标规则可以根据需要添加