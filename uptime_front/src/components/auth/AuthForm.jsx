import React, { useState } from 'react';
import { Server } from 'lucide-react';
import { useAuth } from '../../context/AuthContext';
import { apiFetch } from '../../api';
import { Spinner } from '../ui/Spinner';
import styles from './AuthForm.module.css';

const AuthForm = () => {
    const [isLogin, setIsLogin] = useState(true);
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(false);
    const { login } = useAuth();

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            if (isLogin) {
                const data = await apiFetch('/auth/login', {
                    method: 'POST',
                    body: JSON.stringify({ email, password }),
                });
                login(data.token);
            } else {
                await apiFetch('/auth/register', {
                    method: 'POST',
                    body: JSON.stringify({ name, email, password }),
                });
                const data = await apiFetch('/auth/login', {
                    method: 'POST',
                    body: JSON.stringify({ email, password }),
                });
                login(data.token);
            }
        } catch (err) {
            setError(err.message || 'An error occurred.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className={styles.container}>
            <div className={styles.formWrapper}>
                <div className={styles.header}>
                    <Server size={48} className={styles.logo} />
                    <h1 className={styles.title}>Uptime Monitor</h1>
                    <p className={styles.subtitle}>Sign in to continue to your dashboard.</p>
                </div>
                <div className={styles.formCard}>
                    <h2 className={styles.formTitle}>{isLogin ? 'Welcome Back!' : 'Create Account'}</h2>
                    {error && <p className={styles.errorBox}>{error}</p>}
                    <form onSubmit={handleSubmit} className={styles.form}>
                        {!isLogin && (
                            <div className={styles.inputGroup}>
                                <label className={styles.label}>Name</label>
                                <input
                                    type="text"
                                    value={name}
                                    onChange={(e) => setName(e.target.value)}
                                    className={styles.input}
                                    required
                                />
                            </div>
                        )}
                        <div className={styles.inputGroup}>
                            <label className={styles.label}>Email Address</label>
                            <input
                                type="email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                className={styles.input}
                                required
                            />
                        </div>
                        <div className={styles.inputGroup}>
                            <label className={styles.label}>Password</label>
                            <input
                                type="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                className={styles.input}
                                required
                            />
                        </div>
                        <button
                            type="submit"
                            disabled={loading}
                            className={styles.submitButton}
                        >
                            {loading ? <Spinner /> : (isLogin ? 'Sign In' : 'Sign Up')}
                        </button>
                    </form>
                    <p className={styles.toggleText}>
                        {isLogin ? "Don't have an account?" : 'Already have an account?'}
                        <button onClick={() => setIsLogin(!isLogin)} className={styles.toggleButton}>
                            {isLogin ? 'Sign Up' : 'Sign In'}
                        </button>
                    </p>
                </div>
            </div>
        </div>
    );
};

export default AuthForm;
