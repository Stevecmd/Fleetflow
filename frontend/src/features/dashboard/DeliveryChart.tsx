import React from 'react';
import { Line } from 'react-chartjs-2';

const DeliveryChart: React.FC = () => {
  const data = {
    labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
    datasets: [
      {
        label: 'Deliveries',
        data: [12, 19, 3, 5, 2, 3, 15],
        fill: false,
        borderColor: 'rgb(75, 192, 192)',
        tension: 0.1
      }
    ]
  };

  return <Line data={data} options={{ responsive: true }} />;
};

export default DeliveryChart;
