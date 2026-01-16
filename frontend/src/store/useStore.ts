import { create } from 'zustand';
import api from '../api/client';

/**
 * Message 定义了在前端展示的消息对象模型
 */
interface Message {
    id: number;          // 数据库主键
    sender: string;      // 发送方 (User 或 Bot)
    content: string;     // 全文字内容
    msg_type: string;    // 类型: text/image
    created_at: string;  // 生成时间
}

/**
 * AppState 定义了系统管理的全局状态机骨架
 */
interface AppState {
    token: string | null;            // 管理员当前连接凭据
    isAuthenticated: boolean;        // 用户是否已通过登录校验
    messages: Message[];             // 历史消息池，用于仪表盘和控制台展示
    config: any;                     // 数据库中的后端实时配置参数
    setToken: (token: string) => void; // 方法: 存储并激活有效令牌
    logout: () => void;              // 方法: 执行登出清理
    fetchMessages: () => Promise<void>; // 方法: 拉取最新消息流
    fetchConfig: () => Promise<void>;   // 方法: 同步后端配置状态
    addMessage: (msg: Message) => void; // 方法: 向本地池中插入即时消息 (通常来自 WS)
}

/**
 * useStore 全局状态钩子
 * 基于 Zustand 实现的高性能无侵入式状态流管理。
 */
export const useStore = create<AppState>((set) => ({
    // 初始化阶段从本地存储尝试恢复令牌状态
    token: localStorage.getItem('token'),
    isAuthenticated: !!localStorage.getItem('token'),
    messages: [],
    config: {},

    /**
     * 登录成功后调用此方法持久化身份状态
     */
    setToken: (token: string) => {
        localStorage.setItem('token', token);
        set({ token, isAuthenticated: true });
    },

    /**
     * 清空所有敏感本地状态并登出
     */
    logout: () => {
        localStorage.removeItem('token');
        set({ token: null, isAuthenticated: false });
    },

    /**
     * 通过 REST API 获取过去的一段消息历史
     */
    fetchMessages: async () => {
        try {
            const res = await api.get<Message[]>('/messages');
            // 更新全局池
            set({ messages: res.data });
        } catch (error) {
            console.error("无法获取消息历史:", error);
        }
    },

    /**
     * 拉取后端的机器人/系统配置集
     */
    fetchConfig: async () => {
        try {
            const res = await api.get('/config');
            set({ config: res.data });
        } catch (error) {
            console.error("配置同步失败:", error);
        }
    },

    /**
     * 往全局消息堆栈的顶部（新位置）快速推送新到的消息
     */
    addMessage: (msg: Message) => set((state) => ({
        messages: [msg, ...state.messages.slice(0, 99)] // 仅在控制台维护最新的 100 条记录以保证渲染性能
    })),
}));
