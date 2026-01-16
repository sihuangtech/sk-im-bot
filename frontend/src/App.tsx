import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import ChatConsole from './pages/ChatConsole';
import Config from './pages/Config';
import MainLayout from './components/MainLayout';
import { useStore } from './store/useStore';

/**
 * PrivateRoute 高阶组件：实现核心业务页面的授权校验。
 * 如果未通过 JWT 认证，则强制路由至登录页。
 */
const PrivateRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const { isAuthenticated } = useStore();
    return isAuthenticated ? <>{children}</> : <Navigate to="/login" />;
};

/**
 * App 根组件：管理应用路由分发逻辑。
 */
const App: React.FC = () => {
    return (
        <BrowserRouter>
            <Routes>
                {/* 独立于主布局的登录鉴权页 */}
                <Route path="/login" element={<Login />} />

                {/* 受保护的所有功能模块路由，由 MainLayout 承载侧边栏和头部 */}
                <Route path="/" element={<PrivateRoute><MainLayout /></PrivateRoute>}>
                    {/* 路径匹配为 / 时渲染统计看板 */}
                    <Route index element={<Dashboard />} />
                    {/* 消息实时控制台 */}
                    <Route path="chat" element={<ChatConsole />} />
                    {/* 系统配置动态管理 */}
                    <Route path="config" element={<Config />} />
                </Route>
            </Routes>
        </BrowserRouter>
    );
};

export default App;
