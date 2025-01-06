import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { AppDispatch } from '../../store';
import { useAuth } from '../auth/AuthProvider';
import { RootState } from '../../store';
import { api } from '../../services/api';
import VehicleStats from './components/VehicleStats';
import LatestOrders from './components/LatestOrders';
import TrackingDelivery from './components/TrackingDelivery';
import DriverStats from './components/DriverStats';
import DriverPerformance from './DriverPerformance';
import { fetchVehicleStats, fetchOrders, fetchDriverMetrics, fetchVehicleInfo, fetchDeliveries } from './dashboardSlice';
import { Bar, Line, Doughnut, Pie } from 'react-chartjs-2';
import '../../config/chartConfig';

interface DriverMetrics {
  deliveries_completed: number;
  on_time_delivery_rate: number;
  customer_rating_avg: number;
  fuel_efficiency: number;
  safety_score: number;
  total_distance_covered: number;
}

interface VehicleInfo {
  plate_number: string;
  type: string;
  make: string;
  model: string;
  fuel_type: string;
  status: string;
  last_maintenance: string;
  next_maintenance: string;
}

interface Delivery {
  tracking_number: string;
  status: string;
  pickup_time: string;
  delivery_time: string;
  from_location: string;
  cargo_type: string;
  cargo_weight: number;
}

/**
 * The DriverDashboard component displays the driver's dashboard with their vehicle stats, latest orders,
 * delivery tracking, driver stats, and performance.
 *
 * The component fetches the vehicle stats and orders when the user is authenticated.
 *
 * If an error occurs during the fetch, the component displays an error message with a retry button.
 * If the fetch is successful, the component displays the vehicle stats, latest orders, delivery tracking,
 * driver stats, and performance.
 *
 * @returns The DriverDashboard component.
 */
const DriverDashboard: React.FC = () => {
  const dispatch: AppDispatch = useDispatch();
  const { user } = useAuth();
  const { 
    loading, 
    error, 
    vehicleStats, 
    orders,
    driverMetrics,
    vehicleInfo,
    deliveries 
  } = useSelector((state: RootState) => state.dashboard);

  useEffect(() => {
    if (user?.id) {
      // Dispatch all data fetching actions
      dispatch(fetchVehicleStats(user.id));
      dispatch(fetchOrders(user.id));
      dispatch(fetchDriverMetrics(user.id));
      dispatch(fetchVehicleInfo(user.id));
      dispatch(fetchDeliveries(user.id));
    }
  }, [dispatch, user?.id]);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-t-2 border-b-2 border-indigo-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-red-50">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-red-700 mb-4">Dashboard Error</h2>
          <p className="text-red-600 mb-6">{error}</p>
          <button 
            onClick={() => {
            }}
            className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
          >
            Retry Loading
          </button>
        </div>
      </div>
    );
  }

  // Ensure vehicleStats is of type VehicleData | null
  if (vehicleStats) {
    console.log('Vehicle Stats:', vehicleStats);
  }

  // Chart configurations
  const performanceChart = {
    labels: ['On-Time Delivery', 'Customer Rating', 'Safety Score', 'Fuel Efficiency'],
    datasets: [{
      data: [
        driverMetrics?.on_time_delivery_rate || 0,
        driverMetrics?.customer_rating_avg || 0,
        driverMetrics?.safety_score || 0,
        driverMetrics?.fuel_efficiency || 0
      ],
      backgroundColor: ['#4CAF50', '#2196F3', '#FFC107', '#F44336'],
    }]
  };

  const deliveryTrendChart = {
    labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
    datasets: [{
      label: 'Deliveries Completed',
      data: [12, 15, 18, 14, 16, 19],
      borderColor: '#4CAF50',
      fill: false
    }]
  };

  const CardMetric = ({ title, value, subtitle }: { title: string; value: string | number; subtitle?: string }) => (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h3 className="text-gray-500 text-sm font-medium">{title}</h3>
      <p className="text-3xl font-bold mt-2">{value}</p>
      {subtitle && <p className="text-gray-400 text-sm mt-2">{subtitle}</p>}
    </div>
  );

  const chartContainerStyle = {
    position: 'relative' as const,
    height: '300px',
    width: '100%',
    padding: '20px'
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'bottom' as const,
        labels: {
          boxWidth: 12,
          padding: 15,
          usePointStyle: true
        }
      }
    },
    layout: {
      padding: {
        top: 20,
        bottom: 20,
        left: 20,
        right: 20
      }
    }
  };

  return (
    <div className="p-6 bg-gray-100 min-h-screen">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Driver Dashboard</h1>
        <div className="text-right">
          <p className="text-lg font-semibold">{user?.first_name} {user?.last_name}</p>
          <p className="text-gray-500">ID: {user?.id}</p>
        </div>
      </div>

      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <CardMetric 
          title="Deliveries Completed" 
          value={driverMetrics?.deliveries_completed || 0}
          subtitle="This month" 
        />
        <CardMetric 
          title="On-Time Rate" 
          value={`${driverMetrics?.on_time_delivery_rate || 0}%`}
          subtitle="Last 30 days" 
        />
        <CardMetric 
          title="Avg Rating" 
          value={driverMetrics?.customer_rating_avg || 0}
          subtitle="Out of 5.0" 
        />
        <CardMetric 
          title="Safety Score" 
          value={`${driverMetrics?.safety_score || 0}%`}
          subtitle="Based on driving patterns" 
        />
      </div>

      {/* Vehicle Information */}
      <div className="bg-white p-6 rounded-lg shadow-md mb-8">
        <h2 className="text-xl font-semibold mb-4">Assigned Vehicle</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div>
            <p className="text-gray-500">Vehicle</p>
            <p className="font-semibold">{vehicleInfo?.make} {vehicleInfo?.model}</p>
          </div>
          <div>
            <p className="text-gray-500">Plate Number</p>
            <p className="font-semibold">{vehicleInfo?.plate_number}</p>
          </div>
          <div>
            <p className="text-gray-500">Last Maintenance</p>
            <p className="font-semibold">{vehicleInfo?.last_maintenance || 'N/A'}</p>
          </div>
          <div>
            <p className="text-gray-500">Status</p>
            <p className="font-semibold">{vehicleInfo?.status}</p>
          </div>
        </div>
      </div>

      {/* Charts Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <div className="bg-white p-6 rounded-lg shadow-md flex flex-col">
          <h2 className="text-xl font-semibold mb-4">Performance Metrics</h2>
          <div style={chartContainerStyle}>
            {driverMetrics && (
              <Doughnut 
                data={performanceChart} 
                options={chartOptions}
                className="max-w-full"
              />
            )}
          </div>
        </div>
        
        <div className="bg-white p-6 rounded-lg shadow-md flex flex-col">
          <h2 className="text-xl font-semibold mb-4">Delivery Trend</h2>
          <div style={chartContainerStyle}>
            {deliveries && (
              <Line 
                data={deliveryTrendChart} 
                options={chartOptions}
                className="max-w-full"
              />
            )}
          </div>
        </div>
      </div>

      {/* Recent Deliveries Table */}
      <div className="bg-white p-6 rounded-lg shadow-md">
        <h2 className="text-xl font-semibold mb-4">Recent Deliveries</h2>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead>
              <tr>
                <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tracking #</th>
                <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">From</th>
                <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Cargo</th>
                <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Weight</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {deliveries.map((delivery: Delivery, index: number) => (
                <tr key={index}>
                  <td className="px-6 py-4 whitespace-nowrap">{delivery.tracking_number}</td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full 
                      ${delivery.status === 'completed' ? 'bg-green-100 text-green-800' : 
                      delivery.status === 'in_transit' ? 'bg-blue-100 text-blue-800' : 
                      'bg-yellow-100 text-yellow-800'}`}>
                      {delivery.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">{delivery.from_location}</td>
                  <td className="px-6 py-4 whitespace-nowrap">{delivery.cargo_type}</td>
                  <td className="px-6 py-4 whitespace-nowrap">{delivery.cargo_weight} kg</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default DriverDashboard;
