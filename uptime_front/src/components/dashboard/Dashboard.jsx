import React, { useState, useEffect } from 'react';
import { useAuth } from '../../context/AuthContext';
import { apiFetch } from '../../api';
import { Plus, LogOut, Server } from 'lucide-react';
import { Spinner } from '../ui/Spinner';
import { Toast } from '../ui/Toast';
import { AddURLModal } from './AddURLModal';
import { HistoryModal } from './HistoryModal';
import { URLCard } from './URLCard';
import styles from './Dashboard.module.css';

const Dashboard = () => {
    const [urls, setUrls] = useState([]);
    const [loading, setLoading] = useState(true);
    const [isAddModalOpen, setIsAddModalOpen] = useState(false);
    const [isHistoryModalOpen, setIsHistoryModalOpen] = useState(false);
    const [selectedUrl, setSelectedUrl] = useState(null);
    const [toast, setToast] = useState(null);
    const { logout } = useAuth();

    useEffect(() => {
        apiFetch('/api/urls')
            .then(data => {
                setUrls(data || []);
                setLoading(false);
            })
            .catch(error => {
                console.error("Failed to fetch URLs", error);
                setLoading(false);
                if (error.message.includes("401")) logout();
            });
    }, [logout]);

    const handleUrlAdded = (newUrl) => {
        setUrls(prev => [newUrl, ...prev]);
        setToast({ message: 'Monitor added successfully!', type: 'success' });
    };

    const handleDelete = async (id) => {
        if (window.confirm('Are you sure you want to delete this monitor?')) {
            try {
                await apiFetch(`/api/urls/${id}`, { method: 'DELETE' });
                setUrls(prev => prev.filter(url => url.id !== id));
                setToast({ message: 'Monitor deleted.', type: 'success' });
            } catch (error) {
                console.error("Failed to delete URL", error);
                setToast({ message: 'Failed to delete monitor.', type: 'error' });
            }
        }
    };

    const handleViewHistory = (url) => {
        setSelectedUrl(url);
        setIsHistoryModalOpen(true);
    };

    return (
        <div className={styles.dashboard}>
            <header className={styles.header}>
                <div className={styles.headerTitle}>
                    <Server className={styles.headerIcon} size={28} />
                    <h1>Dashboard</h1>
                </div>
                <button onClick={logout} className={styles.logoutButton}>
                    <LogOut size={20} />
                    <span>Logout</span>
                </button>
            </header>

            <main className={styles.mainContent}>
                <div className={styles.mainHeader}>
                    <h2>Your Monitors ({urls.length})</h2>
                    <button onClick={() => setIsAddModalOpen(true)} className={styles.addButton}>
                        <Plus size={20} />
                        <span>Add Monitor</span>
                    </button>
                </div>

                {loading ? <Spinner /> : (
                    <div className={styles.urlList}>
                        {urls.length > 0 ? (
                            urls.map(url => (
                                <URLCard key={url.id} url={url} onDelete={handleDelete} onViewHistory={handleViewHistory} />
                            ))
                        ) : (
                            <div className={styles.emptyState}>
                                <Server size={48} />
                                <h3>No monitors yet</h3>
                                <p>Click "Add Monitor" to get started.</p>
                            </div>
                        )}
                    </div>
                )}
            </main>

            <AddURLModal isOpen={isAddModalOpen} onClose={() => setIsAddModalOpen(false)} onUrlAdded={handleUrlAdded} />
            <HistoryModal isOpen={isHistoryModalOpen} onClose={() => setIsHistoryModalOpen(false)} url={selectedUrl} />
            {toast && <Toast message={toast.message} type={toast.type} onDismiss={() => setToast(null)} />}
        </div>
    );
};

export default Dashboard;
