// src/features/dashboard/components/DriverStats.tsx
import React from 'react';
import { Bar } from 'react-chartjs-2';
import 'chart.js/auto';

interface DeliveryStats {
  onTime: number;
  unloading: number;
  waiting: number;
  totalDeliveries: number;
}

interface DriverStatsProps {
  data: DeliveryStats | null;
}

const DriverStats: React.FC<DriverStatsProps> = ({ data }) => {
  if (!data) {
    return (
      <div className="bg-white rounded-xl shadow-sm p-6">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/2 mb-6"></div>
          <div className="grid grid-cols-3 gap-6">
            {[1, 2, 3].map((i) => (
              <div key={i} className="bg-gray-50 rounded-lg p-4">
                <div className="h-4 bg-gray-200 rounded w-1/2 mb-2"></div>
                <div className="h-8 bg-gray-200 rounded w-3/4"></div>
              </div>
            ))}
          </div>
          <div className="mt-6 h-48 bg-gray-200 rounded"></div>
        </div>
      </div>
    );
  }

  const chartData = {
    labels: ['On Time', 'Unloading', 'Waiting'],
    datasets: [
      {
        label: 'Deliveries',
        data: [
          data.totalDeliveries > 0 ? (data.onTime / data.totalDeliveries) * 100 : 0,
          data.totalDeliveries > 0 ? (data.unloading / data.totalDeliveries) * 100 : 0,
          data.totalDeliveries > 0 ? (data.waiting / data.totalDeliveries) * 100 : 0,
        ],
        backgroundColor: ['#4caf50', '#ff9800', '#f44336'],
      },
    ],
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      y: {
        beginAtZero: true,
        max: 100,
      },
    },
  };

  return (
    <div className="bg-white rounded-xl shadow-sm p-6">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg font-semibold text-gray-900">Performance</h3>
        <div className="flex space-x-2">
          <button className="inline-flex items-center px-3 py-1 border border-transparent text-sm font-medium rounded-full text-white bg-black">
            Chart
          </button>
        </div>
      </div>

      <div className="grid grid-cols-3 gap-6 mb-6">
        {[
          {
            label: 'On Time',
            value: `${data.totalDeliveries > 0 ? ((data.onTime / data.totalDeliveries) * 100).toFixed(1) : 0}%`,
            count: `${data.onTime}/${data.totalDeliveries}`,
            percentage: data.totalDeliveries > 0 ? (data.onTime / data.totalDeliveries) * 100 : 0,
          },
          {
            label: 'Unloading',
            value: `${data.totalDeliveries > 0 ? ((data.unloading / data.totalDeliveries) * 100).toFixed(1) : 0}%`,
            count: `${data.unloading}/${data.totalDeliveries}`,
            percentage: data.totalDeliveries > 0 ? (data.unloading / data.totalDeliveries) * 100 : 0,
          },
          {
            label: 'Waiting',
            value: `${data.totalDeliveries > 0 ? ((data.waiting / data.totalDeliveries) * 100).toFixed(1) : 0}%`,
            count: `${data.waiting}/${data.totalDeliveries}`,
            percentage: data.totalDeliveries > 0 ? (data.waiting / data.totalDeliveries) * 100 : 0,
          },
        ].map((stat, index) => (
          <div key={index} className="bg-gray-50 rounded-lg p-4">
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-500">{stat.label}</span>
              <span className="text-sm font-medium text-gray-900">
                {stat.count}
              </span>
            </div>
            <div className="mt-2">
              <div className="text-2xl font-bold text-gray-900">{stat.value}</div>
              <div className="mt-2 h-1.5 w-full bg-gray-200 rounded-full overflow-hidden">
                <div
                  className="h-full bg-indigo-600 rounded-full"
                  style={{ width: `${stat.percentage}%` }}
                ></div>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="h-64 bg-gray-50 rounded-lg flex items-center justify-center">
        <Bar data={chartData} options={chartOptions} />
      </div>
    </div>
  );
};

export default DriverStats;
