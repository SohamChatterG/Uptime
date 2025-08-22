import React, { useState } from 'react';
import { Server, Zap, Globe, BarChart3, Mail, Lock, User, GitMerge, Chrome } from 'lucide-react';
import { useAuth } from '../../context/AuthContext';
import { apiFetch } from '../../api';
import { Spinner } from '../ui/Spinner';
import styles from './AuthForm.module.css';

const API_BASE_URL = 'http://localhost:8080';

const Feature = ({ icon, title, description }) => (
    <li className={styles.featureItem}>
        <div className={styles.featureIcon}>{icon}</div>
        <div>
            <h3 className={styles.featureTitle}>{title}</h3>
            <p className={styles.featureDescription}>{description}</p>
        </div>
    </li>
);

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
                const data = await apiFetch('/auth/login', { method: 'POST', body: JSON.stringify({ email, password }) });
                login(data.token);
            } else {
                await apiFetch('/auth/register', { method: 'POST', body: JSON.stringify({ name, email, password }) });
                const data = await apiFetch('/auth/login', { method: 'POST', body: JSON.stringify({ email, password }) });
                login(data.token);
            }
        } catch (err) {
            setError(err.message || 'An error occurred.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className={styles.pageContainer}>
            <div className={styles.featurePanel}>
                <div className={styles.featureContent}>
                    <div className={styles.logoHeader}>
                        <Server size={32} />
                        <h1 className={styles.appName}>Uptime</h1>
                    </div>
                    <h2 className={styles.tagline}>Proactive Uptime Monitoring for Modern Teams.</h2>
                    <p className={styles.description}>
                        Never let your users discover downtime before you do. Uptime provides real-time alerts, detailed analytics, and global monitoring to ensure your services are always online.
                    </p>
                    <ul className={styles.featureList}>
                        <Feature icon={<Zap size={20} />} title="Instant Alerts" description="Get notified immediately via Email, Slack, or Webhooks when your site goes down." />
                        <Feature icon={<Globe size={20} />} title="Global Monitoring" description="Check your website's availability from multiple locations around the world." />
                        <Feature icon={<BarChart3 size={20} />} title="Detailed Analytics" description="Visualize uptime, response times, and status codes with beautiful, insightful charts." />
                    </ul>
                </div>
            </div>

            <div className={styles.formPanel}>
                <div className={styles.formWrapper}>
                    <h2 className={styles.formTitle}>{isLogin ? 'Welcome Back!' : 'Create Your Account'}</h2>
                    <p className={styles.formSubtitle}>
                        {isLogin ? "Sign in to access your dashboard." : "Join thousands of developers keeping their services online."}
                    </p>

                    {error && <p className={styles.errorBox}>{error}</p>}

                    <form onSubmit={handleSubmit} className={styles.form}>
                        {!isLogin && (
                            <div className={styles.inputGroup}>
                                <User size={16} className={styles.inputIcon} />
                                <input type="text" placeholder="Full Name" value={name} onChange={(e) => setName(e.target.value)} className={styles.input} required />
                            </div>
                        )}
                        <div className={styles.inputGroup}>
                            <Mail size={16} className={styles.inputIcon} />
                            <input type="email" placeholder="Email Address" value={email} onChange={(e) => setEmail(e.target.value)} className={styles.input} required />
                        </div>
                        <div className={styles.inputGroup}>
                            <Lock size={16} className={styles.inputIcon} />
                            <input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} className={styles.input} required />
                        </div>
                        <button type="submit" disabled={loading} className={styles.submitButton}>
                            {loading ? <Spinner /> : (isLogin ? 'Sign In' : 'Create Account')}
                        </button>
                    </form>

                    <div className={styles.separator}>
                        <span>OR</span>
                    </div>

                    {/* --- THIS IS THE UPDATED SECTION --- */}
                    <div className={styles.socialLogin}>
                        <a href={`${API_BASE_URL}/auth/google/login`} className={styles.socialButton}>
                            <Chrome size={20} /> Sign in with Google
                        </a>
                        <a href={`${API_BASE_URL}/auth/github/login`} className={styles.socialButton}>
                            <GitMerge size={20} /> Sign in with GitHub
                        </a>
                    </div>

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
