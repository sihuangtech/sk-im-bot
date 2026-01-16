import React, { useEffect } from 'react';
import { Form, Input, Button, Card, Switch, InputNumber, Row, Col, message } from 'antd';
import { useStore } from '../store/useStore';
import api from '../api/client';

/**
 * Config 页面：在线动态调整机器人的各项运行指标。
 * 支持 LLM 模型选择、Token 设置以及 QQ/Discord 的接口使能。
 */
const Config: React.FC = () => {
    const { config, fetchConfig } = useStore();
    const [form] = Form.useForm();

    // 加载页面时拉取最新的服务端配置
    useEffect(() => {
        fetchConfig();
    }, []);

    // 当全局状态机内的配置更新时，同步刷新表单界面值
    useEffect(() => {
        if (config && Object.keys(config).length > 0) {
            form.setFieldsValue(config);
        }
    }, [config]);

    /**
     * 提交修改后的配置项至后端持久化
     */
    const onFinish = async (values: any) => {
        try {
            await api.post('/config', values);
            message.success('系统配置已成功推送到后端并实时生效');
            // 重新刷新本地快照
            fetchConfig();
        } catch (error) {
            message.error('配置保存失败, 请检查后端连接或权限');
        }
    };

    return (
        <div className="glass-panel" style={{ padding: 24 }}>
            <h2 style={{ marginBottom: 24 }}>系统全局配置管理</h2>
            <Form form={form} layout="vertical" onFinish={onFinish}>
                <Row gutter={24}>
                    {/* 左侧：大语言模型相关设定 */}
                    <Col span={12}>
                        <Card title="大语言模型 (LLM) 参数" className="glass-card" bordered={false}>
                            <Form.Item name={['llm', 'provider']} label="供应商接口 (e.g. openai)">
                                <Input placeholder="如: openai" />
                            </Form.Item>
                            <Form.Item name={['llm', 'api_key']} label="API 鉴权密钥 (API Key)">
                                <Input.Password placeholder="sk-..." />
                            </Form.Item>
                            <Form.Item name={['llm', 'model']} label="指定模型名称">
                                <Input placeholder="如: gpt-3.5-turbo" />
                            </Form.Item>
                            <Form.Item name={['llm', 'max_tokens']} label="单次回复最大词数 (Tokens)">
                                <InputNumber style={{ width: '100%' }} min={100} max={4000} />
                            </Form.Item>
                        </Card>
                    </Col>

                    {/* 右侧：社交平台接入使能及 Token 设置 */}
                    <Col span={12}>
                        <Card title="平台集成连接设定" className="glass-card" bordered={false}>
                            <Form.Item name={['qq', 'enabled']} label="激活 QQ 模块" valuePropName="checked">
                                <Switch />
                            </Form.Item>
                            <Form.Item name={['qq', 'ws_url']} label="QQ OneBot API 地址">
                                <Input placeholder="ws://localhost:8080" />
                            </Form.Item>

                            <Form.Item name={['discord', 'enabled']} label="激活 Discord 模块" valuePropName="checked">
                                <Switch />
                            </Form.Item>
                            <Form.Item name={['discord', 'token']} label="Discord 官方授权 Token">
                                <Input.Password placeholder="从开发者面板获取的 Bot Token" />
                            </Form.Item>
                        </Card>
                    </Col>
                </Row>

                {/* 底部操作区 */}
                <div style={{ marginTop: 24, textAlign: 'right' }}>
                    <Button type="primary" htmlType="submit" size="large" style={{ padding: '0 40px' }}>
                        立即应用改动
                    </Button>
                </div>
            </Form>
        </div>
    );
};

export default Config;
