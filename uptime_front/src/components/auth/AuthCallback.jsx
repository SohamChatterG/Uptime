import React, { useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { Spinner } from '../ui/Spinner';

const AuthCallback = () => {
    const { login } = useAuth();
    const location = useLocation();
    const navigate = useNavigate();

    useEffect(() => {
        // This component's job is to grab the token from the URL query parameters
        const params = new URLSearchParams(location.search);
        const token = params.get('token');

        if (token) {
            // If a token is found, we call the login function to save it
            login(token);
            // Then we redirect the user to the main dashboard
            navigate('/');
        } else {
            // If something went wrong and there's no token, redirect back to the login page
            console.error("OAuth callback is missing the token.");
            navigate('/');
        }
    }, [location, login, navigate]);

    // Show a loading spinner while the logic runs
    return (
        <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'white', gap: '1rem' }}>
            <Spinner />
            <p>Finalizing login, please wait...</p>
        </div>
    );
};

export default AuthCallback;

