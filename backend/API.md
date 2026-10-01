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
    },
    "default_password": false
  }
  ```
  `default_password` 为 true 表示仍在使用旧版本的默认密码，前端会提示修改。

### 2. 修改密码
- **URL**: `/api/auth/change-password`
- **方法**: POST（需要登录）
- **请求体**: 
  ```json
  {
    "oldPassword": "string",
    "newPassword": "string"
  }
  ```
  新密码至少 8 位，不超过 72 字节，不能与旧密码相同。
- **响应**: 
  ```json
  {
    "message": "密码修改成功",
    "token": "string"
  }
  ```
  修改成功后之前签发的 token 全部失效，需要改用响应里的新 token。

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
    "probe_disabled": "boolean",
    "created_at": "datetime",
    "updated_at": "datetime"
  }]
  ```

### 4. 添加网址
- **URL**: `/api/sites`
- **方法**: POST
- **请求体**: 网址对象（不含 id、created_at、updated_at，传了也会被忽略）。`url` 只支持 http/https，可以省略协议（按 https 处理）；`probe_disabled` 为 true 时不做可用性检测
- **响应**: 创建的网址对象

### 5. 更新网址
- **URL**: `/api/sites/:id`
- **方法**: PUT
- **请求体**: 网址对象（只更新路径里 id 对应的记录，请求体里的 id 会被忽略）
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

## 可用性检测与 Prometheus 接口

### 15. 获取检测结果
- **URL**: `/api/probe/status`
- **方法**: GET（公开）
- **响应**: key 为工具 id；`state` 取值 `up` / `down` / `pending`，关闭了检测的工具不出现
  ```json
  {
    "enabled": true,
    "interval_seconds": 60,
    "last_round_at": "datetime",
    "results": {
      "12": {"state": "up", "status_code": 200, "duration_ms": 87, "checked_at": "datetime"},
      "13": {"state": "down", "duration_ms": 5000, "reason": "timeout", "checked_at": "datetime"}
    }
  }
  ```
  `reason` 取值：`http_status`、`timeout`、`dns`、`connection_refused`、`tls`、`invalid_url`、`network`。

### 16. 立即检测
- **URL**: `/api/probe/check/:id`
- **方法**: POST（需要登录）
- **响应**: 单个工具的检测结果，格式同上

### 17. Prometheus 指标
- **URL**: `/metrics`
- **方法**: GET（设置了 `METRICS_TOKEN` 时需要 `Authorization: Bearer <METRICS_TOKEN>`）
- **响应**: Prometheus 文本格式，指标说明见 README

### 18. Prometheus http_sd 服务发现
- **URL**: `/prometheus/targets`
- **方法**: GET（鉴权同上）
- **查询参数**: project、environment、category（可选）
- **响应**:
  ```json
  [{"targets": ["https://grafana.example.com"], "labels": {"tool_id": "12", "tool": "Grafana", "project": "monitor", "environment": "prod", "category": "监控"}}]
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
1. 管理员账号：
   - 用户名：admin
   - 初始密码：首次启动时取环境变量 ADMIN_PASSWORD；未设置则随机生成并打印在启动日志里

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