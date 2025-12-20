import axios from 'axios'
import { Message } from 'element-ui'
import router from '@/router'

// ... 原有的 validateStatus ...

var instance = axios.create({
    timeout: 15000,
    baseURL: process.env.NODE_ENV === 'production' ? '/api' : '/api', // 注意 baseURL 改为 /api
    validateStatus
})

// 请求拦截器：添加Token
instance.interceptors.request.use(
    config => {
        const token = localStorage.getItem('token')
        if (token) {
            config.headers['Authorization'] = 'Bearer ' + token
        }
        return config
    },
    error => Promise.reject(error)
)

instance.interceptors.response.use(
    response => response.data,
    err => {
        if (err.response && err.response.status === 401) {
            Message.error('登录过期，请重新登录')
            localStorage.removeItem('token')
            router.push('/login')
        } else {
            Message.error(err.response?.data?.msg || '请求失败')
        }
        return Promise.reject(err)
    }
)

export default instance
