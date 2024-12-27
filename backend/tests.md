# OpsPortal 测试用例

## 1. 认证测试

### 1.1 登录测试
```bash
# 成功登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 错误密码
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"wrong"}'

# 用户名不存在
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"notexist","password":"123456"}'
```

### 1.2 修改密码测试
```bash
# 成功修改密码
curl -X POST http://localhost:8080/api/auth/change-password \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"oldPassword":"admin123","newPassword":"newpass123"}'

# 旧密码错误
curl -X POST http://localhost:8080/api/auth/change-password \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"oldPassword":"wrong","newPassword":"newpass123"}'
```

## 2. 项目管理测试

### 2.1 项目 CRUD 测试
```bash
# 创建项目
curl -X POST http://localhost:8080/api/projects \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"test-project","label":"测试项目"}'

# 获取项目列表
curl -X GET http://localhost:8080/api/projects \
  -H "Authorization: Bearer <token>"

# 更新项目
curl -X PUT http://localhost:8080/api/projects/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"test-project","label":"更新的项目名"}'

# 删除项目
curl -X DELETE http://localhost:8080/api/projects/1 \
  -H "Authorization: Bearer <token>"
```

## 3. 环境管理测试

### 3.1 环境 CRUD 测试
```bash
# 创建环境
curl -X POST http://localhost:8080/api/environments \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"dev","label":"开发环境","project_id":1}'

# 获取环境列表（按项目筛选）
curl -X GET "http://localhost:8080/api/environments?project_id=1" \
  -H "Authorization: Bearer <token>"

# 更新环境
curl -X PUT http://localhost:8080/api/environments/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"dev","label":"开发环境-更���","project_id":1}'

# 删除环境
curl -X DELETE http://localhost:8080/api/environments/1 \
  -H "Authorization: Bearer <token>"
```

## 4. 网址管理测试

### 4.1 网址 CRUD 测试
```bash
# 创建网址
curl -X POST http://localhost:8080/api/sites \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试网址",
    "url": "http://example.com",
    "description": "这是一个测试网址",
    "environment": "dev",
    "project": "test-project",
    "category": "测试分类"
  }'

# 获取网址列表（带筛选）
curl -X GET "http://localhost:8080/api/sites?env=dev&project=test-project" \
  -H "Authorization: Bearer <token>"

# 更新网址
curl -X PUT http://localhost:8080/api/sites/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "更新的网址",
    "url": "http://example.com/updated",
    "description": "这是更新后的网址",
    "environment": "dev",
    "project": "test-project",
    "category": "测试分类"
  }'

# 删除网址
curl -X DELETE http://localhost:8080/api/sites/1 \
  -H "Authorization: Bearer <token>"
```

## 5. 错误处理测试

### 5.1 认证错误
```bash
# 无token访问
curl -X GET http://localhost:8080/api/projects

# token无效
curl -X GET http://localhost:8080/api/projects \
  -H "Authorization: Bearer invalid-token"
```

### 5.2 参数验证
```bash
# 项目名称格式错误
curl -X POST http://localhost:8080/api/projects \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Project","label":"测试项目"}'

# URL格式错误
curl -X POST http://localhost:8080/api/sites \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试网址",
    "url": "invalid-url",
    "environment": "dev",
    "project": "test-project"
  }'
``` 