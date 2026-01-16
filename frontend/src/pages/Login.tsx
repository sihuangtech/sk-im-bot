import React, { useState } from 'react';
import { Form, Input, Button, Card, message } from 'antd';
import { useStore } from '../store/useStore';
import { useNavigate } from 'react-router-dom';
import api from '../api/client';

const Login: React.FC = () => {
    const navigate = useNavigate();
    const setToken = useStore((state) => state.setToken);
    const [loading, setLoading] = useState(false);

    const onFinish = async (values: any) => {
        setLoading(true);
        try {
            const res = await api.post('/login', values);
            setToken(res.data.token);
            message.success('Login success');
            navigate('/');
        } catch (error) {
            message.error('Login failed');
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
            <Card
                className="glass-panel"
                style={{ width: 400, border: 'none' }}
            >
                <h2 style={{ textAlign: 'center', marginBottom: 32, color: '#fff' }}>SK-IM-Bot Admin</h2>
                <Form onFinish={onFinish} layout="vertical">
                    <Form.Item name="username" rules={[{ required: true }]}>
                        <Input placeholder="Username" size="large" />
                    </Form.Item>
                    <Form.Item name="password" rules={[{ required: true }]}>
                        <Input.Password placeholder="Password" size="large" />
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit" block size="large" loading={loading}
                            style={{ background: '#722ed1', border: 'none' }}>
                            Login
                        </Button>
                    </Form.Item>
                </Form>
            </Card>
        </div>
    );
};

export default Login;
