import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { AppDispatch } from '../../store';
import { useAuth } from '../auth/AuthProvider';
import { RootState } from '../../store';
import VehicleStats from './components/VehicleStats';
import LatestOrders from './components/LatestOrders';
import TrackingDelivery from './components/TrackingDelivery';
import DriverStats from './components/DriverStats';
import DriverPerformance from './DriverPerformance';
import { fetchVehicleStats, fetchOrders } from './dashboardSlice';

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
  const { loading, error, vehicleStats, orders, deliveryStats } = useSelector(
    (state: RootState) => state.dashboard
  );

  useEffect(() => {
    if (user?.id) {
      if (!vehicleStats) { 
        dispatch(fetchVehicleStats(user.id));
        console.log('Fetching Vehicle Stats...');
      }
      if (!orders) { 
        dispatch(fetchOrders(user.id));
        console.log('Fetching Orders...');
      }
    }
  }, [dispatch, user?.id, vehicleStats, orders]);

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

  return (
    <div className="p-4">
      <h1 className="text-3xl font-bold mb-4">
        Driver Dashboard for {user ? `${user.first_name} ${user.last_name}` : 'Driver'}
      </h1>
      <div className="flex justify-between mb-4">
        <p>Driver ID: {user ? user.id : 'N/A'}</p>
        <p>Email: {user ? user.email : 'N/A'}</p>
      </div>
      <VehicleStats data={vehicleStats} />
      <LatestOrders data={orders} />
      <TrackingDelivery />
      <DriverStats data={deliveryStats} />
      <DriverPerformance performance={deliveryStats} />
    </div>
  );
};

export default DriverDashboard;
