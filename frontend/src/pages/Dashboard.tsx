import React, { useEffect } from 'react';
import { useStore } from '../store/useStore';
import { Row, Col, Card, Statistic } from 'antd';
import { MessageSquare, Users, Activity, Server } from 'lucide-react';

const Dashboard: React.FC = () => {
    const { messages, fetchMessages } = useStore();

    useEffect(() => {
        fetchMessages();
    }, []);

    const stats = [
        { title: 'Total Messages', value: messages.length, icon: <MessageSquare size={24} color="#722ed1" /> },
        { title: 'Active Sessions', value: 12, icon: <Users size={24} color="#1890ff" /> },
        { title: 'System Uptime', value: '24h', icon: <Activity size={24} color="#52c41a" /> },
        { title: 'Platforms', value: '2', icon: <Server size={24} color="#faad14" /> },
    ];

    return (
        <div>
            <h2 style={{ marginBottom: 24 }}>System Overview</h2>
            <Row gutter={[24, 24]}>
                {stats.map((stat, idx) => (
                    <Col span={6} key={idx}>
                        <Card className="glass-card" bordered={false}>
                            <Statistic
                                title={<div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>{stat.icon} {stat.title}</div>}
                                value={stat.value}
                                valueStyle={{ color: '#fff', fontWeight: 'bold' }}
                            />
                        </Card>
                    </Col>
                ))}
            </Row>

            <Row style={{ marginTop: 24 }}>
                <Col span={24}>
                    <Card className="glass-panel" title="Recent Activity" bordered={false} style={{ color: '#fff' }}>
                        {/* Simple list or chart placeholder */}
                        <div style={{ padding: 20, textAlign: 'center', color: 'rgba(255,255,255,0.5)' }}>
                            {messages.length > 0 ? 'Messages stream active' : 'No recent activity'}
                        </div>
                    </Card>
                </Col>
            </Row>
        </div>
    );
};

export default Dashboard;
