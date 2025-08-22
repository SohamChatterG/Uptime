import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import styles from './ResponseChart.module.css';

const CustomTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
        return (
            <div className={styles.tooltip}>
                <p className={styles.tooltipLabel}>{`Time: ${label}`}</p>
                <p className={styles.tooltipValue}>{`Response: ${payload[0].value}ms`}</p>
            </div>
        );
    }
    return null;
};

export const ResponseChart = ({ data }) => {
    const chartData = data
        .filter(item => item.was_successful)
        .reverse()
        .map(item => ({
            time: new Date(item.checked_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
            responseTime: item.response_time_ms,
        }));

    if (chartData.length === 0) {
        return <div className={styles.noData}>Not enough data to display a chart.</div>
    }

    return (
        <div className={styles.chartContainer}>
            <ResponsiveContainer width="100%" height={250}>
                <LineChart
                    data={chartData}
                    margin={{
                        top: 5,
                        right: 20,
                        left: -10,
                        bottom: 5,
                    }}
                >
                    <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                    <XAxis dataKey="time" stroke="#9ca3af" fontSize={12} />
                    <YAxis stroke="#9ca3af" fontSize={12} unit="ms" />
                    <Tooltip content={<CustomTooltip />} cursor={{ fill: 'rgba(139, 92, 246, 0.1)' }} />
                    <Legend wrapperStyle={{ fontSize: "14px" }} />
                    <Line
                        type="monotone"
                        dataKey="responseTime"
                        name="Response Time (ms)"
                        stroke="#818cf8"
                        dot={{ r: 3 }}
                        activeDot={{ r: 6 }}
                    />
                </LineChart>
            </ResponsiveContainer>
        </div>
    );
};