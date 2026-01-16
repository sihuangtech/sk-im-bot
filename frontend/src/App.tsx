import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import ChatConsole from './pages/ChatConsole';
import Config from './pages/Config';
import MainLayout from './components/MainLayout';
import { useStore } from './store/useStore';

const PrivateRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const { isAuthenticated } = useStore();
    return isAuthenticated ? <>{children}</> : <Navigate to="/login" />;
};

const App: React.FC = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/login" element={<Login />} />
                <Route path="/" element={<PrivateRoute><MainLayout /></PrivateRoute>}>
                    <Route index element={<Dashboard />} />
                    <Route path="chat" element={<ChatConsole />} />
                    <Route path="config" element={<Config />} />
                </Route>
            </Routes>
        </BrowserRouter>
    );
};

export default App;
