import React from 'react';
import { Loader } from 'lucide-react';
import styles from './Spinner.module.css';

export const Spinner = () => (
    <div className={styles.spinnerContainer}>
        <Loader className={styles.spinner} size={48} />
    </div>
);
