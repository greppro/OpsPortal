import request from '@/utils/request'

// 获取工具列表
export function getTools() {
    return request({
        url: '/api/sites',
        method: 'get'
    })
}

// 获取项目列表
export function getProjects() {
    return request({
        url: '/api/projects',
        method: 'get'
    })
}

// 获取环境列表
export function getEnvironments() {
    return request({
        url: '/api/environments',
        method: 'get'
    })
}

// 获取当前激活的公告
export function getActiveNotice() {
    return request({
        url: '/api/notices/active',
        method: 'get'
    })
}

// 获取所有公告
export function getNotices() {
    return request({
        url: '/api/notices',
        method: 'get'
    })
}

// 创建公告
export function createNotice(data) {
    return request({
        url: '/api/notices',
        method: 'post',
        data
    })
}

// 更新公告
export function updateNotice(id, data) {
    return request({
        url: `/api/notices/${id}`,
        method: 'put',
        data
    })
}

// 删除公告
export function deleteNotice(id) {
    return request({
        url: `/api/notices/${id}`,
        method: 'delete'
    })
} 