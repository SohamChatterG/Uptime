import React, { useState } from 'react';
import { apiFetch } from '../../api';
import { Modal } from '../ui/Modal';
import styles from './Dashboard.module.css'; // Re-use some styles

export const AddURLModal = ({ isOpen, onClose, onUrlAdded }) => {
    const [name, setName] = useState('');
    const [url, setUrl] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        try {
            const newUrl = await apiFetch('/api/urls', {
                method: 'POST',
                body: JSON.stringify({ name, url }),
            });
            onUrlAdded(newUrl);
            setName('');
            setUrl('');
            onClose();
        } catch (error) {
            console.error("Failed to add URL", error);
            // In a real app, show an error toast
        } finally {
            setLoading(false);
        }
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Add New Monitor">
            <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                <div>
                    <label style={{ fontSize: '0.875rem', fontWeight: 500, color: 'var(--color-text-muted)' }}>Monitor Name</label>
                    <input
                        type="text"
                        placeholder="e.g., My Personal Blog"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        className="modal-input"
                        required
                    />
                </div>
                <div>
                    <label style={{ fontSize: '0.875rem', fontWeight: 500, color: 'var(--color-text-muted)' }}>URL</label>
                    <input
                        type="url"
                        placeholder="https://example.com"
                        value={url}
                        onChange={(e) => setUrl(e.target.value)}
                        className="modal-input"
                        required
                    />
                </div>
                <div style={{ display: 'flex', justifyContent: 'flex-end', paddingTop: '1rem', gap: '0.75rem' }}>
                    <button type="button" onClick={onClose} className="modal-button-secondary">Cancel</button>
                    <button type="submit" disabled={loading} className="modal-button-primary">
                        {loading ? 'Adding...' : 'Add Monitor'}
                    </button>
                </div>
            </form>
            <style>{`
                .modal-input {
                    width: 100%; margin-top: 0.25rem; padding: 0.75rem; background-color: #374151;
                    border: 1px solid #4b5563; border-radius: 0.375rem; color: var(--color-text-main);
                }
                .modal-input:focus {
                    outline: 2px solid transparent; outline-offset: 2px; border-color: var(--color-primary);
                    box-shadow: 0 0 0 2px var(--color-primary);
                }
                .modal-button-primary {
                    padding: 0.5rem 1rem; font-weight: 600; background-color: var(--color-primary);
                    border-radius: 0.375rem;
                }
                .modal-button-primary:hover { background-color: var(--color-primary-hover); }
                .modal-button-primary:disabled { background-color: #4f46e580; cursor: not-allowed; }
                .modal-button-secondary {
                    padding: 0.5rem 1rem; border-radius: 0.375rem; color: var(--color-text-main);
                }
                .modal-button-secondary:hover { background-color: #374151; }
            `}</style>
        </Modal>
    );
};

