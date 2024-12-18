import React from 'react';

interface DeliveryStats {
  totalDeliveries: number;
  onTime: number;
  unloading: number;
  waiting: number;
}

interface DriverPerformanceProps {
  performance: DeliveryStats | null;
}

const DriverPerformance: React.FC<DriverPerformanceProps> = ({ performance }) => {
  return (
    <div className="bg-white shadow rounded-lg p-6">
      <h2 className="text-xl font-semibold mb-4">Driver Performance</h2>
      <ul>
        <li>Deliveries Completed: {performance ? performance.totalDeliveries : 'N/A'}</li>
        <li>On-Time Delivery Rate: {performance ? performance.onTime : 'N/A'}%</li>
        <li>Unloading Time: {performance ? performance.unloading : 'N/A'} minutes</li>
        <li>Waiting Time: {performance ? performance.waiting : 'N/A'} minutes</li>
      </ul>
    </div>
  );
};

export default DriverPerformance;
