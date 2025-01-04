import { useDispatch, useSelector } from 'react-redux';
import { AppDispatch, RootState } from '../../store';
import DriverDashboard from './DriverDashboard';
import FleetManagerDashboard from './FleetManagerDashboard';

/**
 * Dashboard component that renders the main dashboard view.
 * It displays a welcome message and conditionally renders either
 * the DriverDashboard or FleetManagerDashboard based on the current user role.
 *
 * Uses Redux to dispatch actions and select the current role from the store.
 */

const Dashboard: React.FC = () => {
  const dispatch: AppDispatch = useDispatch();
  const currentRole = useSelector((state: RootState) => state.dashboard.currentRole);

  return (
    <div className="bg-white shadow rounded-lg p-6">
      <h1 className="text-2xl font-semibold text-gray-900">Dashboard</h1>
      <p className="mt-4 text-gray-600">Welcome to the FleetFlow Dashboard</p>
      {currentRole === 'driver' ? <DriverDashboard /> : <FleetManagerDashboard />}
    </div>
  );
};

export default Dashboard;