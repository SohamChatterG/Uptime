import React, { useMemo } from 'react';
import { PieChart, Pie, Cell, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import styles from './StatusPieChart.module.css';

// Define colors for different status categories
const COLORS = {
    '2xx': '#22c55e', // Green for success
    'Down': '#ef4444', // Red for down/error
    '4xx': '#f97316', // Orange for client errors
    '5xx': '#eab308', // Yellow for server errors
    'Other': '#6b7280' // Gray for others
};

// Helper function to group status codes into categories
const getStatusCategory = (code) => {
    if (code >= 200 && code < 300) return '2xx';
    if (code >= 400 && code < 500) return '4xx';
    if (code >= 500 && code < 600) return '5xx';
    if (code === 0) return 'Down';
    return 'Other';
};

export const StatusPieChart = ({ data }) => {
    // Process the raw history data into a format the pie chart can use
    const chartData = useMemo(() => {
        if (!data || data.length === 0) return [];

        const statusCounts = data.reduce((acc, check) => {
            const category = getStatusCategory(check.status_code);
            acc[category] = (acc[category] || 0) + 1;
            return acc;
        }, {});

        return Object.entries(statusCounts).map(([name, value]) => ({
            name,
            value,
        }));
    }, [data]);

    if (chartData.length === 0) {
        return <div className={styles.noData}>No status data available.</div>;
    }

    return (
        <div className={styles.chartContainer}>
            <ResponsiveContainer width="100%" height={250}>
                <PieChart>
                    <Pie
                        data={chartData}
                        cx="50%"
                        cy="50%"
                        labelLine={false}
                        outerRadius={80}
                        fill="#8884d8"
                        dataKey="value"
                        nameKey="name"
                        label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                    >
                        {chartData.map((entry, index) => (
                            <Cell key={`cell-${index}`} fill={COLORS[entry.name]} />
                        ))}
                    </Pie>
                    <Tooltip />
                    <Legend wrapperStyle={{ fontSize: "14px" }} />
                </PieChart>
            </ResponsiveContainer>
        </div>
    );
};