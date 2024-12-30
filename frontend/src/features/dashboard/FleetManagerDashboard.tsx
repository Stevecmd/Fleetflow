import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { AppDispatch } from '../../store';
import { RootState } from '../../store';
import VehicleStats from './components/VehicleStats';
import DriverPerformance from './DriverPerformance';
import { fetchFleetVehicles, fetchFleetPerformance } from './dashboardSlice';

/**
 * FleetManagerDashboard is a React functional component that displays the dashboard
 * for fleet managers. It fetches and displays data related to fleet vehicles and
 * driver performance. The component uses Redux to manage state and dispatches 
 * actions to fetch vehicle and performance data when the component mounts. It 
 * handles loading and error states by displaying appropriate UI feedback.
 */

const FleetManagerDashboard: React.FC = () => {
  const dispatch: AppDispatch = useDispatch();
  const { loading, error, vehicleStats, driverPerformance, user } = useSelector(
    (state: RootState) => ({
      ...state.dashboard,
      user: state.auth.user,
    })
  );

  useEffect(() => {
    if (user?.id) {
      dispatch(fetchFleetVehicles());
      dispatch(fetchFleetPerformance());
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
              // Retry logic can be implemented here
            }}
            className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
          >
            Retry Loading
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="p-4">
      <h1 className="text-3xl font-bold mb-4">Fleet Manager Dashboard</h1>
      <VehicleStats data={vehicleStats} />
      <DriverPerformance performance={driverPerformance} />
    </div>
  );
};

export default FleetManagerDashboard;
