import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'
import { ConfigProvider, theme } from 'antd'

ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <ConfigProvider
            theme={{
                algorithm: theme.darkAlgorithm,
                token: {
                    colorPrimary: '#722ed1', // Purple accent
                    borderRadius: 8,
                },
            }}
        >
            <App />
        </ConfigProvider>
    </React.StrictMode>,
)
