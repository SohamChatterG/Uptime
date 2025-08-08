import React, { useState, useEffect, createContext, useContext, useCallback } from 'react';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [token, setToken] = useState(() => localStorage.getItem('uptime_token'));
    const [loading, setLoading] = useState(true);

    const login = (newToken) => {
        localStorage.setItem('uptime_token', newToken);
        setToken(newToken);
    };

    const logout = useCallback(() => {
        localStorage.removeItem('uptime_token');
        setToken(null);
        setUser(null);
    }, []);

    useEffect(() => {
        if (token) {
            setUser({ isAuthenticated: true });
        } else {
            setUser(null);
        }
        setLoading(false);
    }, [token]);

    const value = { user, token, login, logout, loading };

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);

