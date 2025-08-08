import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import styles from './ResponseChart.module.css';

// A custom tooltip to match our UI style
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
    // Recharts needs the data in a specific format. We process it here.
    const chartData = data
        // We only want to chart successful checks that have a response time
        .filter(item => item.was_successful)
        // Reverse the array to show oldest to newest (left to right)
        .reverse()
        // Format the data for the chart
        .map(item => ({
            // Format the timestamp to be more readable
            time: new Date(item.checked_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
            responseTime: item.response_time_ms,
        }));

    if (chartData.length === 0) {
        return <div className={styles.noData}>Not enough data to display a chart.</div>
    }

    return (
        <div className={styles.chartContainer}>
            {/* The ResponsiveContainer makes the chart fill its parent div */}
            <ResponsiveContainer width="100%" height={250}>
                <LineChart
                    data={chartData}
                    margin={{
                        top: 5,
                        right: 20,
                        left: -10, // Adjust to make Y-axis labels fit
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
                        stroke="#818cf8" // A nice indigo color
                        strokeWidth={2}
                        dot={{ r: 3 }}
                        activeDot={{ r: 6 }}
                    />
                </LineChart>
            </ResponsiveContainer>
        </div>
    );
};