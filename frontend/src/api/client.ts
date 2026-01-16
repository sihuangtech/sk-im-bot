import axios from 'axios';

/**
 * 全局 Axios API 客户端配置
 * 该实例封装了针对后端 API 的网络请求，支持鉴权拦截和自动注销逻辑。
 */
const api = axios.create({
    baseURL: '/api', // 基准路径，开发阶段由 Vite 代理映射到 8888 端口
});

/**
 * 请求拦截器：在每一个发出的网络请求头部自动植入 JWT 令牌。
 */
api.interceptors.request.use((config) => {
    // 从浏览器的持久化存储中读取当前存有的 token
    const token = localStorage.getItem('token');
    if (token) {
        // 按照 RFC 6750 格式设置 HTTP 标准头
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

/**
 * 响应拦截器：针对服务器返回的全局错误（如令牌失效）进行统一捕获处理。
 */
api.interceptors.response.use(
    (response) => response,
    (error) => {
        // 捕获 401 Unauthorized 状态码，通常表示 Token 已过期或被禁用
        if (error.response?.status === 401) {
            // 本地清理失效的旧令牌
            localStorage.removeItem('token');
            // 强制用户跳转至登录界面重新鉴权
            window.location.href = '/login';
        }
        // 向上传递其余业务错误
        return Promise.reject(error);
    }
);

export default api;
