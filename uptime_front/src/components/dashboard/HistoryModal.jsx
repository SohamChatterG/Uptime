import React, { useState, useEffect, useMemo } from 'react';
import { apiFetch } from '../../api';
import { Modal } from '../ui/Modal';
import { Spinner } from '../ui/Spinner';
import { CheckCircle2, XCircle, Clock, BarChart, LineChart, PieChart as PieChartIcon } from 'lucide-react';
import { ResponseChart } from './ResponseChart';
import { StatusPieChart } from './StatusPieChart';
import styles from './HistoryModal.module.css';

const StatCard = ({ icon, label, value }) => (
    <div className={styles.statCard}>
        <div className={styles.statIcon}>{icon}</div>
        <div>
            <p className={styles.statLabel}>{label}</p>
            <p className={styles.statValue}>{value}</p>
        </div>
    </div>
);

export const HistoryModal = ({ isOpen, onClose, url }) => {
    const [history, setHistory] = useState([]);
    const [loading, setLoading] = useState(true);
    const [activeChart, setActiveChart] = useState('line'); // <-- 1. NEW STATE FOR TOGGLE

    const summaryStats = useMemo(() => {
        if (!history || history.length === 0) {
            return { uptime: 'N/A', avgResponse: 'N/A' };
        }
        const successfulChecks = history.filter(h => h.was_successful);
        const uptime = (successfulChecks.length / history.length) * 100;

        const totalResponseTime = successfulChecks.reduce((acc, h) => acc + h.response_time_ms, 0);
        const avgResponse = successfulChecks.length > 0 ? (totalResponseTime / successfulChecks.length) : 0;

        return {
            uptime: `${uptime.toFixed(2)}%`,
            avgResponse: `${Math.round(avgResponse)}ms`
        };
    }, [history]);

    useEffect(() => {
        if (isOpen && url) {
            setLoading(true);
            apiFetch(`/api/urls/${url.id}/history?limit=200`)
                .then(data => {
                    setHistory(data || []);
                    setLoading(false);
                })
                .catch(error => {
                    console.error("Failed to fetch history", error);
                    setLoading(false);
                });
        }
    }, [isOpen, url]);

    return (
        <Modal isOpen={isOpen} onClose={onClose} title={`Details for ${url?.name}`}>
            {loading ? <Spinner /> : (
                <>
                    <div className={styles.statsGrid}>
                        <StatCard icon={<BarChart />} label="Uptime (Last 200)" value={summaryStats.uptime} />
                        <StatCard icon={<Clock />} label="Avg. Response" value={summaryStats.avgResponse} />
                    </div>

                    {/* --- 2. NEW TOGGLE BUTTONS --- */}
                    <div className={styles.chartToggleContainer}>
                        <button
                            className={`${styles.chartToggleButton} ${activeChart === 'line' ? styles.active : ''}`}
                            onClick={() => setActiveChart('line')}
                        >
                            <LineChart size={16} />
                            <span>Response Time</span>
                        </button>
                        <button
                            className={`${styles.chartToggleButton} ${activeChart === 'pie' ? styles.active : ''}`}
                            onClick={() => setActiveChart('pie')}
                        >
                            <PieChartIcon size={16} />
                            <span>Status Codes</span>
                        </button>
                    </div>

                    {/* --- 3. CONDITIONAL RENDERING FOR CHARTS --- */}
                    <div className={styles.chartDisplayArea}>
                        {activeChart === 'line' && <ResponseChart data={history} />}
                        {activeChart === 'pie' && <StatusPieChart data={history} />}
                    </div>

                    <h3 className={styles.historyTitle}>Recent Checks</h3>
                    <div className={styles.historyContainer}>
                        <ul className={styles.historyList}>
                            {history.length > 0 ? history.slice(0, 20).map(check => (
                                <li key={check.id} className={styles.historyItem}>
                                    <div className={styles.historyItemMain}>
                                        {check.was_successful ? (
                                            <CheckCircle2 size={20} className={styles.iconUp} />
                                        ) : (
                                            <XCircle size={20} className={styles.iconDown} />
                                        )}
                                        <div>
                                            <p className={styles.historyDate}>
                                                {new Date(check.checked_at).toLocaleString()}
                                            </p>
                                            <p className={styles.historyStatus}>
                                                Status: {check.status_code || 'N/A'}
                                            </p>
                                        </div>
                                    </div>
                                    <span className={styles.historyResponseTime}>{check.response_time_ms}ms</span>
                                </li>
                            )) : <p className={styles.noHistory}>No history available for this monitor yet.</p>}
                        </ul>
                    </div>
                </>
            )}
        </Modal>
    );
};
