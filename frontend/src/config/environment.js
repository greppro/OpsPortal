// 环境类型配置
export const ENV_TYPES = {
  'dev': {
    type: 'success',
    label: '开发环境'
  },
  'test': {
    type: '',
    label: '测试环境'
  },
  'pre': {
    type: 'warning',
    label: '预发环境'
  },
  'prod': {
    type: 'danger',
    label: '生产环境'
  }
}

// 获取环境标签类型
export const getEnvTagType = (env) => {
  return ENV_TYPES[env]?.type || 'primary'
}

// 获取环境显示标签
export const getEnvLabel = (env) => {
  return ENV_TYPES[env]?.label || env
} 