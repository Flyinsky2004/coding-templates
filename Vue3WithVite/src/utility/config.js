// 开发环境配置
const BACKEND_DOMAIN = 'http://localhost:8080/api/'
const FRONTEND_DOMAIN = 'http://localhost:5173/'
// 生产环境配置
// const BACKEND_DOMAIN = `${currentDomain}${apiPrefix}/`
// const FRONTEND_DOMAIN = `${currentDomain}/`
// 获取当前域名
const currentDomain = window.location.protocol + '//' + window.location.host

// 根据环境判断是否需要添加 /api 前缀
const isDevEnv = process.env.NODE_ENV === 'development'
const apiPrefix = isDevEnv ? '/api' : '/api'
export {BACKEND_DOMAIN,FRONTEND_DOMAIN,currentDomain,apiPrefix}
