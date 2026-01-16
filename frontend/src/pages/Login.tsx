import React, { useState } from 'react';
import { Form, Input, Button, Card, message } from 'antd';
import { useStore } from '../store/useStore';
import { useNavigate } from 'react-router-dom';
import api from '../api/client';

/**
 * Login 页面：管理后台唯一访问入口。
 * 采用了极简风格，背景由全局毛玻璃 CSS 驱动。
 */
const Login: React.FC = () => {
    const navigate = useNavigate();
    // 从全局状态中获取 Token 写入方法
    const setToken = useStore((state) => state.setToken);
    // 渲染态管理：标记当前是否正在执行网络验证
    const [loading, setLoading] = useState(false);

    /**
     * 处理登录表单提交
     * @param values 包含 username 和 password 的对象
     */
    const onFinish = async (values: any) => {
        setLoading(true);
        try {
            // 发起凭据校验请求
            const res = await api.post('/login', values);
            // 校验通过，存储 JWT 至 LocalStorage 及 Store
            setToken(res.data.token);
            message.success('登录验证通过，欢迎回来');
            // 重定向至主仪表盘
            navigate('/');
        } catch (error) {
            // 后端 401 或网络响应异常处理
            message.error('身份凭据验证失败, 请检查用户名或密码');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            height: '100vh',
        }}>
            {/* 使用自定义的 glass-panel 类实现半透明毛玻璃卡片效果 */}
            <Card
                className="glass-panel"
                style={{ width: 400, border: 'none', padding: '10px' }}
                bordered={false}
            >
                <div style={{ textAlign: 'center', marginBottom: 40 }}>
                    <h2 style={{ color: '#fff', fontSize: '28px', fontWeight: 'bold' }}>SK-IM-Bot Admin</h2>
                    <p style={{ color: 'rgba(255,255,255,0.4)' }}>多平台对话机器人管理系统</p>
                </div>

                <Form onFinish={onFinish} layout="vertical">
                    <Form.Item
                        name="username"
                        rules={[{ required: true, message: '请输入管理员账号' }]}
                    >
                        <Input placeholder="用户名 (默认: admin)" size="large" />
                    </Form.Item>

                    <Form.Item
                        name="password"
                        rules={[{ required: true, message: '请输入对应的访问口令' }]}
                    >
                        <Input.Password placeholder="密码 (默认: admin)" size="large" />
                    </Form.Item>

                    <Form.Item style={{ marginTop: 20 }}>
                        <Button
                            type="primary"
                            htmlType="submit"
                            block
                            size="large"
                            loading={loading}
                            style={{
                                background: 'linear-gradient(45deg, #722ed1 0%, #391085 100%)',
                                border: 'none',
                                height: '48px',
                                borderRadius: '8px'
                            }}
                        >
                            执行认证登入
                        </Button>
                    </Form.Item>
                </Form>
            </Card>
        </div>
    );
};

export default Login;
