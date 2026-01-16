import { create } from 'zustand';
import api from '../api/client';

interface Message {
    id: number;
    sender: string;
    content: string;
    msg_type: string;
    created_at: string;
}

interface AppState {
    token: string | null;
    isAuthenticated: boolean;
    messages: Message[];
    config: any;
    setToken: (token: string) => void;
    logout: () => void;
    fetchMessages: () => Promise<void>;
    fetchConfig: () => Promise<void>;
    addMessage: (msg: Message) => void;
}

export const useStore = create<AppState>((set) => ({
    token: localStorage.getItem('token'),
    isAuthenticated: !!localStorage.getItem('token'),
    messages: [],
    config: {},

    setToken: (token: string) => {
        localStorage.setItem('token', token);
        set({ token, isAuthenticated: true });
    },

    logout: () => {
        localStorage.removeItem('token');
        set({ token: null, isAuthenticated: false });
    },

    fetchMessages: async () => {
        try {
            const res = await api.get<Message[]>('/messages');
            set({ messages: res.data });
        } catch (error) {
            console.error(error);
        }
    },

    fetchConfig: async () => {
        try {
            const res = await api.get('/config');
            set({ config: res.data });
        } catch (error) {
            console.error(error);
        }
    },

    addMessage: (msg: Message) => set((state) => ({ messages: [msg, ...state.messages] })),
}));
