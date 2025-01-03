import React, { useState, useEffect } from 'react';
import DeliveryChart from './components/DeliveryChart';
import DriverStatusList from './components/DriverStatusList';
import FleetOverview from './components/FleetOverview';
import AlertsPanel from './components/AlertsPanel';
import { dashboardData } from './mockData';

const ManagerDashboard: React.FC = () => {
  const [data, setData] = useState(dashboardData);

  useEffect(() => {
    // Simply set the mock data
    setData(dashboardData);
  }, []);

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-6">Manager Dashboard</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Stats Cards */}
        <div className="bg-white p-4 rounded-lg shadow">
          <div className="flex items-center">
            <span className="text-xl">🚚</span>
            <div className="ml-3">
              <p className="text-gray-500">Total Deliveries</p>
              <p className="text-2xl font-semibold">{data?.stats.totalDeliveries || 0}</p>
            </div>
          </div>
        </div>
        <div className="bg-white p-4 rounded-lg shadow">
          <div className="flex items-center">
            <span className="text-xl">🕒</span>
            <div className="ml-3">
              <p className="text-gray-500">Active Deliveries</p>
              <p className="text-2xl font-semibold">{data?.stats.activeDeliveries || 0}</p>
            </div>
          </div>
        </div>
        <div className="bg-white p-4 rounded-lg shadow">
          <div className="flex items-center">
            <span className="text-xl">👥</span>
            <div className="ml-3">
              <p className="text-gray-500">Available Drivers</p>
              <p className="text-2xl font-semibold">{data?.stats.availableDrivers || 0}</p>
            </div>
          </div>
        </div>
        <div className="bg-white p-4 rounded-lg shadow">
          <div className="flex items-center">
            <span className="text-xl">📦</span>
            <div className="ml-3">
              <p className="text-gray-500">Pending Orders</p>
              <p className="text-2xl font-semibold">{data?.stats.pendingOrders || 0}</p>
            </div>
          </div>
        </div>
        <div className="bg-white p-4 rounded-lg shadow">
          <div className="flex items-center">
            <span className="text-xl">⚠️</span>
            <div className="ml-3">
              <p className="text-gray-500">Alerts</p>
              <p className="text-2xl font-semibold">{data?.stats.alerts || 0}</p>
            </div>
          </div>
        </div>
        <div className="bg-white p-4 rounded-lg shadow">
          <div className="flex items-center">
            <span className="text-xl">🚛</span>
            <div className="ml-3">
              <p className="text-gray-500">Fleet Utilization</p>
              <p className="text-2xl font-semibold">{data?.stats.fleetUtilization || 0}%</p>
            </div>
          </div>
        </div>
        <div className="bg-white p-4 rounded-lg shadow">
          <div className="flex items-center">
            <span className="text-xl">🏢</span>
            <div className="ml-3">
              <p className="text-gray-500">Warehouse Capacity</p>
              <p className="text-2xl font-semibold">{data?.stats.warehouseCapacity || 0}%</p>
            </div>
          </div>
        </div>
      </div>

      <div className="space-y-6 mt-6">
        <DeliveryChart data={data.deliveryTrends} />
        <DriverStatusList data={data.driverStatus} />
        <FleetOverview data={data.fleetOverview} />
        <AlertsPanel data={data.alerts} />
      </div>
    </div>
  );
};

export default ManagerDashboard;
