import React from 'react';
import { Layout, Menu, Button } from 'antd';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { LayoutDashboard, Settings, MessageSquare, LogOut } from 'lucide-react';
import { useStore } from '../store/useStore';

/**
 * MainLayout 核心布局组件
 * 该组件承载了管理界面的侧边导航栏 (Sider) 和主内容展示区 (Outlet)。
 */
const { Header, Sider, Content } = Layout;

const MainLayout: React.FC = () => {
    const navigate = useNavigate();
    const location = useLocation();
    // 从 Zustand Store 中提取注销逻辑
    const { logout } = useStore();

    /**
     * 侧边栏导航菜单项定义
     * 每个项关联一个图标和路由路径。
     */
    const menuItems = [
        {
            key: '/',
            icon: <LayoutDashboard size={18} />,
            label: '统计总览面板',
        },
        {
            key: '/chat',
            icon: <MessageSquare size={18} />,
            label: '实时监控控制台',
        },
        {
            key: '/config',
            icon: <Settings size={18} />,
            label: '系统参数设定',
        },
    ];

    return (
        <Layout style={{ minHeight: '100vh', background: 'transparent' }}>
            {/* 侧边导航栏：具有玻璃质感的暗色风格 */}
            <Sider
                width={250}
                style={{ background: 'rgba(0,0,0,0.2)', borderRight: '1px solid rgba(255,255,255,0.05)' }}
            >
                {/* 侧边栏 LOGO 与 标题区 */}
                <div style={{ height: 64, margin: 16, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                    <h1 style={{ color: '#fff', fontSize: '1.5rem', fontWeight: 'bold' }}>SK-IM-Bot</h1>
                </div>

                {/* 菜单导航 */}
                <Menu
                    theme="dark"
                    mode="inline"
                    selectedKeys={[location.pathname]} // 关键：高亮显示当前所在的菜单项
                    style={{ background: 'transparent', border: 'none' }}
                    items={menuItems}
                    onClick={({ key }) => navigate(key)} // 点击触发路由跳转
                />

                {/* 底部退出操作区 */}
                <div style={{ position: 'absolute', bottom: 20, width: '100%', padding: '0 20px' }}>
                    <Button
                        type="text"
                        danger
                        icon={<LogOut size={18} />}
                        onClick={() => { logout(); navigate('/login'); }}
                        block
                        style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 8, height: '40px' }}
                    >
                        安全退出
                    </Button>
                </div>
            </Sider>

            {/* 顶部与主内容区域 */}
            <Layout style={{ background: 'transparent' }}>
                <Header style={{ padding: 0, background: 'transparent', height: 48 }} />
                <Content style={{ margin: '24px 16px', padding: 24, minHeight: 280, background: 'transparent' }}>
                    {/* 子路由的具体页面将在此处挂载渲染 */}
                    <Outlet />
                </Content>
            </Layout>
        </Layout>
    );
};

export default MainLayout;
