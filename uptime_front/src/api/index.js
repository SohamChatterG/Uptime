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
            return null; // No content
        }

        return response.json();
    } catch (error) {
        console.error("API Fetch Error:", error);
        throw error;
    }
};