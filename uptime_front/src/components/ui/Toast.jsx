import React, { useEffect } from 'react';
import { CheckCircle2, AlertTriangle } from 'lucide-react';
import styles from './Toast.module.css';

export const Toast = ({ message, type, onDismiss }) => {
    useEffect(() => {
        const timer = setTimeout(onDismiss, 3000);
        return () => clearTimeout(timer);
    }, [onDismiss]);

    const typeClass = type === 'success' ? styles.success : styles.error;

    return (
        <div className={`${styles.toast} ${typeClass}`}>
            {type === 'success' ? <CheckCircle2 className={styles.icon} /> : <AlertTriangle className={styles.icon} />}
            <span>{message}</span>
        </div>
    );
};
