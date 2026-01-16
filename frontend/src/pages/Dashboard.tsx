import React, { useEffect } from 'react';
import { useStore } from '../store/useStore';
import { Row, Col, Card, Statistic } from 'antd';
import { MessageSquare, Users, Activity, Server } from 'lucide-react';

/**
 * Dashboard 页面：提供系统运行的核心指标可视化展示。
 */
const Dashboard: React.FC = () => {
    // 从状态管理拉取全局消息数据
    const { messages, fetchMessages } = useStore();

    // 初始加载拉取统计样本
    useEffect(() => {
        fetchMessages();
    }, []);

    /**
     * 系统状态卡片数据集定义
     */
    const stats = [
        {
            title: '历史消息总数',
            value: messages.length,
            icon: <MessageSquare size={24} color="#722ed1" />
        },
        {
            title: '覆盖活跃会话',
            value: 12, // 演示数据，正式需从 API 聚合
            icon: <Users size={24} color="#1890ff" />
        },
        {
            title: '系统稳定运行时长',
            value: '24h',
            icon: <Activity size={24} color="#52c41a" />
        },
        {
            title: '平台连接节点',
            value: '2 (QQ/Discord)',
            icon: <Server size={24} color="#faad14" />
        },
    ];

    return (
        <div>
            <h2 style={{ marginBottom: 24 }}>系统状态总览仪表盘</h2>

            {/* 第一行：关键统计指标数字卡片 */}
            <Row gutter={[24, 24]}>
                {stats.map((stat, idx) => (
                    <Col span={6} key={idx}>
                        <Card className="glass-card" bordered={false}>
                            <Statistic
                                title={
                                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                                        {stat.icon}
                                        <span style={{ color: 'rgba(255,255,255,0.6)' }}>{stat.title}</span>
                                    </div>
                                }
                                value={stat.value}
                                valueStyle={{ color: '#fff', fontWeight: '800', fontSize: '24px' }}
                            />
                        </Card>
                    </Col>
                ))}
            </Row>

            {/* 第二行：系统动态趋势概览区 */}
            <Row style={{ marginTop: 24 }}>
                <Col span={24}>
                    <Card className="glass-panel" title="最新系统行为记录" bordered={false} style={{ color: '#fff' }}>
                        <div style={{ padding: '60px 20px', textAlign: 'center' }}>
                            {/* 如果消息池有数据，提示监控流已对接 */}
                            {messages.length > 0 ? (
                                <div style={{ color: '#52c41a' }}>已对接核心消息分发流，共计检测到 {messages.length} 条已归档记录</div>
                            ) : (
                                <div style={{ color: 'rgba(255,255,255,0.3)' }}>暂无实时活动，请在社交平台触发对话</div>
                            )}
                        </div>
                    </Card>
                </Col>
            </Row>
        </div>
    );
};

export default Dashboard;
