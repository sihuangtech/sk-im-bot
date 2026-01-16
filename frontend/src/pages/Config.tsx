import React, { useEffect } from 'react';
import { Form, Input, Button, Card, Switch, InputNumber, Row, Col, message } from 'antd';
import { useStore } from '../store/useStore';
import api from '../api/client';

const Config: React.FC = () => {
    const { config, fetchConfig } = useStore();
    const [form] = Form.useForm();

    useEffect(() => {
        fetchConfig();
    }, []);

    useEffect(() => {
        if (config) {
            form.setFieldsValue(config);
        }
    }, [config]);

    const onFinish = async (values: any) => {
        try {
            await api.post('/config', values);
            message.success('Configuration updated');
            fetchConfig();
        } catch (error) {
            message.error('Failed to update config');
        }
    };

    return (
        <div className="glass-panel" style={{ padding: 24 }}>
            <h2 style={{ marginBottom: 24 }}>System Configuration</h2>
            <Form form={form} layout="vertical" onFinish={onFinish}>
                <Row gutter={24}>
                    <Col span={12}>
                        <Card title="LLM Settings" className="glass-card" bordered={false}>
                            <Form.Item name={['llm', 'provider']} label="Provider">
                                <Input />
                            </Form.Item>
                            <Form.Item name={['llm', 'api_key']} label="API Key">
                                <Input.Password />
                            </Form.Item>
                            <Form.Item name={['llm', 'model']} label="Model">
                                <Input />
                            </Form.Item>
                            <Form.Item name={['llm', 'max_tokens']} label="Max Tokens">
                                <InputNumber style={{ width: '100%' }} />
                            </Form.Item>
                        </Card>
                    </Col>
                    <Col span={12}>
                        <Card title="Platform Settings" className="glass-card" bordered={false}>
                            <Form.Item name={['qq', 'enabled']} label="QQ Enabled" valuePropName="checked">
                                <Switch />
                            </Form.Item>
                            <Form.Item name={['qq', 'ws_url']} label="QQ WS URL">
                                <Input />
                            </Form.Item>
                            <Form.Item name={['discord', 'enabled']} label="Discord Enabled" valuePropName="checked">
                                <Switch />
                            </Form.Item>
                            <Form.Item name={['discord', 'token']} label="Discord Token">
                                <Input.Password />
                            </Form.Item>
                        </Card>
                    </Col>
                </Row>
                <div style={{ marginTop: 24, textAlign: 'right' }}>
                    <Button type="primary" htmlType="submit" size="large">
                        Save Changes
                    </Button>
                </div>
            </Form>
        </div>
    );
};

export default Config;
