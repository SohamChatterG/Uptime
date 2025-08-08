import React from 'react';
import { Trash2, History } from 'lucide-react';
import styles from './URLCard.module.css';

export const URLCard = ({ url, onDelete, onViewHistory }) => {
    const statusClass = url.status ? styles.statusUp : styles.statusDown;
    const statusText = url.status ? 'Up' : 'Down';

    return (
        <div className={styles.card}>
            <div className={styles.mainInfo}>
                <div className={styles.titleWrapper}>
                    <span className={`${styles.statusDot} ${statusClass}`}></span>
                    <h3 className={styles.title}>{url.name}</h3>
                    <span className={`${styles.statusBadge} ${statusClass}`}>
                        {statusText}
                    </span>
                </div>
                <a href={url.url} target="_blank" rel="noopener noreferrer" className={styles.urlLink}>{url.url}</a>
            </div>
            <div className={styles.actions}>
                <button onClick={() => onViewHistory(url)} className={styles.actionButton}>
                    <History size={20} />
                </button>
                <button onClick={() => onDelete(url.id)} className={`${styles.actionButton} ${styles.deleteButton}`}>
                    <Trash2 size={20} />
                </button>
            </div>
        </div>
    );
};
