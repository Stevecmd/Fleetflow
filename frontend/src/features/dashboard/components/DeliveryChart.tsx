import React from 'react';
import { Line } from 'react-chartjs-2';

const DeliveryChart: React.FC<{ data: { labels: string[], data: number[] } }> = ({ data }) => {
  const chartData = {
    labels: data.labels,
    datasets: [
      {
        label: 'Deliveries',
        data: data.data,
        fill: false,
        borderColor: 'rgb(75, 192, 192)',
        tension: 0.1
      }
    ]
  };

  return <Line data={chartData} options={{ responsive: true }} />;
};

export default DeliveryChart;
