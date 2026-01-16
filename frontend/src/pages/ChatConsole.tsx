import React, { useEffect, useRef } from 'react';
import { useStore } from '../store/useStore';
import { List, Avatar, Tag } from 'antd';

/**
 * ChatConsole 实时消息监控控制台
 * 该页面通过 WebSocket 实时监听机器人与玩家的交互，并实时滚动展示。
 */
const ChatConsole: React.FC = () => {
    // 从全局状态中解构需要的属性和操作方法
    const { messages, addMessage, fetchMessages } = useStore();
    // 使用 useRef 保存 WebSocket 实例引用，防止组件刷新导致连接重建
    const ws = useRef<WebSocket | null>(null);

    useEffect(() => {
        // 组件挂载时首先通过 REST 拉取一部分历史记录
        fetchMessages();

        // ---- 初始化 WebSocket 长连接 ----
        // 注意：生产环境下建议使用相对路径或动态获取域名
        const wsUrl = `ws://${window.location.hostname}:8888/ws`;
        ws.current = new WebSocket(wsUrl);

        ws.current.onopen = () => {
            console.log('监控控制台 WebSocket 已打通');
        };

        ws.current.onmessage = (event) => {
            // 后端以 MessageEvent 格式推送消息
            const msg = JSON.parse(event.data);

            // 将接收到的原始 Socket 数据映射为前端渲染所需的 Message 格式
            const displayMsg = {
                id: Date.now(), // 对于即时消息使用当前时间戳作为唯一键
                sender: msg.Username || '系统',
                content: msg.Content,
                msg_type: msg.MsgType || 'text',
                created_at: new Date().toISOString(),
            };

            // 推入状态机，触发 UI 刷新
            addMessage(displayMsg);
        };

        ws.current.onerror = (err) => {
            console.error('WebSocket 监控服务链路异常:', err);
        };

        // 组件卸载时释放连接资源
        return () => {
            ws.current?.close();
        };
    }, []);

    return (
        <div className="glass-panel" style={{ height: '80vh', display: 'flex', flexDirection: 'column', padding: 24 }}>
            <h2 style={{ marginBottom: 16 }}>实时监控控制台</h2>
            <div style={{ flex: 1, overflowY: 'auto', paddingRight: 8 }}>
                <List
                    itemLayout="horizontal"
                    dataSource={messages}
                    renderItem={(item) => (
                        <List.Item style={{ borderBottom: '1px solid rgba(255,255,255,0.03)' }}>
                            <List.Item.Meta
                                // 使用用户名的首字母作为占位头像
                                avatar={
                                    <Avatar style={{ backgroundColor: item.sender === 'bot' ? '#52c41a' : '#722ed1' }}>
                                        {item.sender[0]?.toUpperCase()}
                                    </Avatar>
                                }
                                title={
                                    <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                                        <span style={{ color: '#fff', fontWeight: '600' }}>{item.sender}</span>
                                        <span style={{ fontSize: 12, color: 'rgba(255,255,255,0.4)' }}>
                                            {new Date(item.created_at).toLocaleTimeString()}
                                        </span>
                                    </div>
                                }
                                description={
                                    <div style={{ marginTop: 4 }}>
                                        {/* 如果是图片消息，展示对应标签（此处可扩展预览图） */}
                                        {item.msg_type === 'image' && <Tag color="blue" style={{ marginBottom: 4 }}>[图片消息]</Tag>}
                                        <div style={{ color: 'rgba(255,255,255,0.7)', wordBreak: 'break-all', fontSize: '14px' }}>
                                            {item.content}
                                        </div>
                                    </div>
                                }
                            />
                        </List.Item>
                    )}
                />
            </div>
        </div>
    );
};

export default ChatConsole;
