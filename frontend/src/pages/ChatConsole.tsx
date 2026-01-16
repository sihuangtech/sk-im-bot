import React, { useEffect, useRef } from 'react';
import { useStore } from '../store/useStore';
import { List, Avatar, Tag } from 'antd';

const ChatConsole: React.FC = () => {
    const { messages, addMessage, fetchMessages } = useStore();
    const ws = useRef<WebSocket | null>(null);

    useEffect(() => {
        fetchMessages();

        // WebSocket Connection
        ws.current = new WebSocket('ws://localhost:8888/ws');

        ws.current.onopen = () => {
            console.log('WS Connected');
        };

        ws.current.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            // Transform backend event to Message if needed, or assume matched
            // Here we assume msg is already in shape or close enough for demo
            // In real app, mapper is needed.
            const displayMsg = {
                id: Date.now(),
                sender: msg.Username || 'System',
                content: msg.Content,
                msg_type: msg.MsgType || 'text',
                created_at: new Date().toISOString(),
            };
            addMessage(displayMsg);
        };

        return () => {
            ws.current?.close();
        };
    }, []);

    return (
        <div className="glass-panel" style={{ height: '80vh', display: 'flex', flexDirection: 'column', padding: 24 }}>
            <h2 style={{ marginBottom: 16 }}>Live Console</h2>
            <div style={{ flex: 1, overflowY: 'auto' }}>
                <List
                    itemLayout="horizontal"
                    dataSource={messages}
                    renderItem={(item) => (
                        <List.Item>
                            <List.Item.Meta
                                avatar={<Avatar style={{ backgroundColor: '#722ed1' }}>{item.sender[0]?.toUpperCase()}</Avatar>}
                                title={
                                    <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                                        <span style={{ color: '#fff' }}>{item.sender}</span>
                                        <span style={{ fontSize: 12, color: 'rgba(255,255,255,0.4)' }}>
                                            {new Date(item.created_at).toLocaleTimeString()}
                                        </span>
                                    </div>
                                }
                                description={
                                    <div>
                                        {item.msg_type === 'image' ? <Tag color="blue">Image</Tag> : null}
                                        <span style={{ color: 'rgba(255,255,255,0.7)' }}>{item.content}</span>
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
