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