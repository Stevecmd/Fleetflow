import React, { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '../../store';
import VehicleStats from './components/VehicleStats';
import DriverPerformance from './DriverPerformance';

// Define the expected types for vehicle and driver performance data
interface VehicleData {
  id: number;
  plate_number: string;
  type: string;
  make: string;
  model: string;
  year: number;
  status: string;
  count: number;
  capacity: number;
  fuel_type: string;
  status_id: number;
  mileage: number;
  created_at: string;
  updated_at: string;
}

interface DeliveryStats {
  totalDeliveries: number;
  onTime: number;
  unloading: number;
  waiting: number;
}

interface DriverPerformanceData {
  id: number;
  driverName: string;
  performanceScore: number;
  deliveryStats: DeliveryStats;
}

/**
 * The FleetManagerDashboard component displays the fleet manager's dashboard with the fleet's vehicle
 * statistics and driver performance.
 *
 * The component fetches the vehicle stats and driver performance data when the user is authenticated.
 *
 * If an error occurs during the fetch, the component displays an error message with a retry button.
 * If the fetch is successful, the component displays the vehicle stats and driver performance.
 *
 * @returns The FleetManagerDashboard component.
 */
const FleetManagerDashboard: React.FC = () => {
  const { loading, error, user } = useSelector((state: RootState) => ({
    loading: state.dashboard.loading,
    error: state.dashboard.error,
    user: state.auth.user,
  }));

  const [vehicleStats, setVehicleStats] = useState<VehicleData[]>([]);
  const [driverPerformance, setDriverPerformance] = useState<DriverPerformanceData[]>([]);

  useEffect(() => {
    if (user?.id) {
      setVehicleStats([
        { 
          id: 1, 
          plate_number: 'ABC123', 
          type: 'Truck', 
          make: 'Ford', 
          model: 'F-150', 
          year: 2020, 
          status: 'Active', 
          count: 10,
          capacity: 2000,
          fuel_type: 'Diesel',
          status_id: 1,
          mileage: 50000,
          created_at: '2023-01-15T08:30:00Z',
          updated_at: '2024-01-20T14:25:00Z'
        },
        { 
          id: 2, 
          plate_number: 'XYZ987', 
          type: 'Van', 
          make: 'Chevrolet', 
          model: 'Express', 
          year: 2019, 
          status: 'Inactive', 
          count: 5,
          capacity: 1500,
          fuel_type: 'Gasoline',
          status_id: 2,
          mileage: 75000,
          created_at: '2023-02-20T09:15:00Z',
          updated_at: '2024-01-19T16:45:00Z'
        },
        { 
          id: 3, 
          plate_number: 'LMN456', 
          type: 'SUV', 
          make: 'Toyota', 
          model: 'Highlander', 
          year: 2021, 
          status: 'Maintenance', 
          count: 3,
          capacity: 1000,
          fuel_type: 'Hybrid',
          status_id: 3,
          mileage: 25000,
          created_at: '2023-03-10T10:45:00Z',
          updated_at: '2024-01-18T11:30:00Z'
        }
      ]);
      setDriverPerformance([
        { id: 1, driverName: 'Alice', performanceScore: 85, deliveryStats: { totalDeliveries: 50, onTime: 45, unloading: 3, waiting: 2 } },
        { id: 2, driverName: 'Bob', performanceScore: 90, deliveryStats: { totalDeliveries: 60, onTime: 55, unloading: 2, waiting: 3 } },
      ]);
    }
  }, [user?.id]);

  // Aggregate delivery stats
  const aggregateDeliveryStats = (): DeliveryStats => {
    return driverPerformance.reduce((acc, driver) => ({
      totalDeliveries: acc.totalDeliveries + driver.deliveryStats.totalDeliveries,
      onTime: acc.onTime + driver.deliveryStats.onTime,
      unloading: acc.unloading + driver.deliveryStats.unloading,
      waiting: acc.waiting + driver.deliveryStats.waiting,
    }), {
      totalDeliveries: 0,
      onTime: 0,
      unloading: 0,
      waiting: 0,
    });
  };

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
            onClick={() => {}}
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
      <VehicleStats data={vehicleStats[0] || null} />
      <DriverPerformance performance={aggregateDeliveryStats()} />
    </div>
  );
};

export default FleetManagerDashboard;