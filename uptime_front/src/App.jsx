import React from 'react';
import { AuthProvider, useAuth } from './context/AuthContext';
import AuthForm from './components/auth/AuthForm';
import Dashboard from './components/dashboard/Dashboard';
import { Spinner } from './components/ui/Spinner';

function App() {
  return (
    <AuthProvider>
      <MainApp />
    </AuthProvider>
  );
}

const MainApp = () => {
  const { user, loading } = useAuth();

  if (loading) {
    return (
      <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <Spinner />
      </div>
    );
  }

  return user ? <Dashboard /> : <AuthForm />;
};

export default App;

// =================================================================================
// FILE: src/api/index.js
// =================================================================================
const API_BASE_URL = 'http://localhost:8080'; // Your Go backend URL

export const apiFetch = async (endpoint, options = {}) => {
  const token = localStorage.getItem('uptime_token');
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  try {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: 'An unknown error occurred' }));
      throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
    }

    if (response.status === 204 || response.headers.get("Content-Length") === "0") {
      return null;
    }

    return response.json();
  } catch (error) {
    console.error("API Fetch Error:", error);
    throw error;
  }
};
