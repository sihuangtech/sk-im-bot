import React from 'react';
import { Layout, Menu, Button } from 'antd';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { LayoutDashboard, Settings, MessageSquare, LogOut } from 'lucide-react';
import { useStore } from '../store/useStore';

const { Header, Sider, Content } = Layout;

const MainLayout: React.FC = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const { logout } = useStore();

    const menuItems = [
        {
            key: '/',
            icon: <LayoutDashboard size={18} />,
            label: 'Dashboard',
        },
        {
            key: '/chat',
            icon: <MessageSquare size={18} />,
            label: 'Chat Console',
        },
        {
            key: '/config',
            icon: <Settings size={18} />,
            label: 'Settings',
        },
    ];

    return (
        <Layout style={{ minHeight: '100vh', background: 'transparent' }}>
            <Sider
                width={250}
                style={{ background: 'rgba(0,0,0,0.2)', borderRight: '1px solid rgba(255,255,255,0.05)' }}
            >
                <div style={{ height: 64, margin: 16, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                    <h1 style={{ color: '#fff', fontSize: '1.5rem', fontWeight: 'bold' }}>SK-IM-Bot</h1>
                </div>
                <Menu
                    theme="dark"
                    mode="inline"
                    selectedKeys={[location.pathname]}
                    style={{ background: 'transparent' }}
                    items={menuItems}
                    onClick={({ key }) => navigate(key)}
                />
                <div style={{ position: 'absolute', bottom: 20, width: '100%', padding: '0 20px' }}>
                    <Button
                        type="text"
                        danger
                        icon={<LogOut size={18} />}
                        onClick={() => { logout(); navigate('/login'); }}
                        block
                        style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 8 }}
                    >
                        Logout
                    </Button>
                </div>
            </Sider>
            <Layout style={{ background: 'transparent' }}>
                <Header style={{ padding: 0, background: 'transparent' }} />
                <Content style={{ margin: '24px 16px', padding: 24, minHeight: 280, background: 'transparent' }}>
                    <Outlet />
                </Content>
            </Layout>
        </Layout>
    );
};

export default MainLayout;
