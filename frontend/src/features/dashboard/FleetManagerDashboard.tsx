import React, { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import { createSelector } from '@reduxjs/toolkit';
import { RootState } from '../../store';
import VehicleStats from './components/VehicleStats';
import DriverPerformance from './DriverPerformance';
import { Line, Doughnut } from 'react-chartjs-2';
import { api } from '../../services/api';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  BarElement,
} from 'chart.js';

// Register ChartJS components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  BarElement
);

// Create memoized selectors
const selectDashboardState = (state: RootState) => state.dashboard;
const selectAuthState = (state: RootState) => state.auth;

const selectDashboardData = createSelector(
  [selectDashboardState, selectAuthState],
  (dashboard, auth) => ({
    loading: dashboard.loading,
    error: dashboard.error,
    user: auth.user,
  })
);

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

interface MaintenanceMetrics {
  pending_maintenance: number;
  completed_last_month: number;
  average_cost: number;
  upcoming_services: Array<{
    vehicle_id: number;
    plate_number: string;
    next_service: string;
    service_type: string;
    estimated_cost: number;
  }>;
}

interface FleetEfficiency {
  fuel_efficiency: number;
  carbon_emissions: number;
  operating_costs: number;
  idle_time: number;
  route_optimization: number;
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
  // Use memoized selector
  const { loading, error, user } = useSelector(selectDashboardData);

  const [vehicleStats, setVehicleStats] = useState<VehicleData[]>([]);
  const [driverPerformance, setDriverPerformance] = useState<DriverPerformanceData[]>([]);
  const [maintenanceMetrics, setMaintenanceMetrics] = useState<MaintenanceMetrics | null>(null);
  const [fleetEfficiency, setFleetEfficiency] = useState<FleetEfficiency | null>(null);

  // Add new state for chart data
  const [performanceData, setPerformanceData] = useState({
    labels: ['January', 'February', 'March', 'April', 'May', 'June'],
    datasets: [
      {
        label: 'Vehicle Utilization',
        data: [65, 70, 75, 80, 85, 90],
        borderColor: 'rgb(75, 192, 192)',
        tension: 0.1,
      },
      {
        label: 'Driver Performance',
        data: [70, 75, 80, 82, 85, 87],
        borderColor: 'rgb(255, 99, 132)',
        tension: 0.1,
      },
    ],
  });

  const [maintenanceData, setMaintenanceData] = useState({
    labels: ['Completed', 'Pending', 'Overdue', 'Scheduled'],
    datasets: [{
      data: [0, 0, 0, 0],
      backgroundColor: [
        'rgb(75, 192, 192)',
        'rgb(255, 205, 86)',
        'rgb(255, 99, 132)',
        'rgb(54, 162, 235)',
      ],
    }],
  });

  // Update chart options with better responsive settings
  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'bottom' as const,
        labels: {
          boxWidth: 12,
          padding: 15,
        },
      },
      title: {
        display: true,
        text: 'Fleet Performance',
        padding: {
          top: 10,
          bottom: 20
        }
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        ticks: {
          maxTicksLimit: 5
        }
      },
      x: {
        ticks: {
          maxRotation: 45,
          minRotation: 45
        }
      }
    },
    layout: {
      padding: {
        left: 10,
        right: 10,
        top: 20,
        bottom: 10
      }
    }
  };

  useEffect(() => {
    const dummyVehicleStats: VehicleData[] = [
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
    ];

    const dummyDriverPerformance: DriverPerformanceData[] = [
      { 
        id: 1, 
        driverName: 'Alice', 
        performanceScore: 85, 
        deliveryStats: { totalDeliveries: 50, onTime: 45, unloading: 3, waiting: 2 } 
      },
      { 
        id: 2, 
        driverName: 'Bob', 
        performanceScore: 90, 
        deliveryStats: { totalDeliveries: 60, onTime: 55, unloading: 2, waiting: 3 } 
      },
    ];

    if (user?.id) {
      setVehicleStats(dummyVehicleStats);
      setDriverPerformance(dummyDriverPerformance);
    }
  }, [user?.id]);

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        const [maintenance, efficiency, analytics] = await Promise.all([
          api.get('/fleet-manager/maintenance-metrics'),
          api.get('/fleet-manager/fleet-efficiency'),
          api.get('/fleet-manager/delivery-analytics')
        ]);
        
        setMaintenanceMetrics(maintenance.data);
        setFleetEfficiency(efficiency.data);

        // Update maintenance chart data
        setMaintenanceData(prev => ({
          ...prev,
          datasets: [{
            ...prev.datasets[0],
            data: [
              maintenance.data.completed_last_month || 0,
              maintenance.data.pending_maintenance || 0,
              maintenance.data.upcoming_services?.length || 0,
              0, // Add actual overdue count if available
            ],
          }],
        }));

        // Update performance chart data with real metrics
        setPerformanceData(prev => ({
          ...prev,
          datasets: [
            {
              ...prev.datasets[0],
              data: [efficiency.data.fuel_efficiency, efficiency.data.route_optimization, 
                    analytics.data.efficiency, analytics.data.on_time_deliveries,
                    efficiency.data.operating_costs, efficiency.data.idle_time].map(val => 
                    Number(val?.toFixed(2)) || 0),
            },
          ],
        }));

      } catch (error) {
        console.error('Error fetching dashboard data:', error);
      }
    };

    fetchDashboardData();
  }, []);

  // Memoize aggregateDeliveryStats
  const aggregateDeliveryStats = React.useMemo(() => {
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
  }, [driverPerformance]);

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
    <div className="p-6 bg-gray-100">
      {/* Fleet Overview Section */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <StatCard 
          title="Total Vehicles" 
          value={vehicleStats.length} 
          trend="+5%" 
          icon="🚛" 
        />
        <StatCard 
          title="Active Drivers" 
          value={driverPerformance.length} 
          trend="+2%" 
          icon="👤" 
        />
        <StatCard 
          title="Fleet Efficiency" 
          value={`${fleetEfficiency?.route_optimization.toFixed(1)}%`} 
          trend="+3.5%" 
          icon="📈" 
        />
        <StatCard 
          title="Maintenance Due" 
          value={maintenanceMetrics?.pending_maintenance || 0} 
          trend="0" 
          icon="🔧" 
        />
      </div>

      {/* Update Charts Section with fixed height containers */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-xl font-semibold mb-4">Fleet Performance</h3>
          <div className="h-[400px] w-full">
            <Line 
              data={performanceData} 
              options={chartOptions}
              height="100%"
              width="100%"
            />
          </div>
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-xl font-semibold mb-4">Maintenance Overview</h3>
          <div className="h-[400px] w-full">
            <Doughnut 
              data={maintenanceData} 
              options={chartOptions}
              height="100%"
              width="100%"
            />
          </div>
        </div>
      </div>

      {/* Maintenance Schedule */}
      <div className="bg-white rounded-lg shadow p-6 mb-8">
        <h3 className="text-xl font-semibold mb-4">Upcoming Maintenance</h3>
        <div className="overflow-x-auto">
          <table className="min-w-full">
            <thead>
              <tr>
                <th>Vehicle</th>
                <th>Service Type</th>
                <th>Due Date</th>
                <th>Est. Cost</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {maintenanceMetrics?.upcoming_services.map(service => (
                <tr key={`${service.vehicle_id}-${service.service_type}`}>
                  <td>{service.plate_number}</td>
                  <td>{service.service_type}</td>
                  <td>{new Date(service.next_service).toLocaleDateString()}</td>
                  <td>${service.estimated_cost}</td>
                  <td>
                    <span className="px-2 py-1 rounded-full bg-yellow-100 text-yellow-800">
                      Pending
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Vehicle List and Driver Performance sections remain unchanged */}
      <h1 className="text-3xl font-bold mb-4">Fleet Manager Dashboard</h1>
      <VehicleStats data={vehicleStats[0] || null} />
      <DriverPerformance performance={aggregateDeliveryStats} />
    </div>
  );
};

// Add new components
const StatCard: React.FC<{ title: string; value: string | number; trend: string; icon: string }> = ({
  title,
  value,
  trend,
  icon
}) => (
  <div className="bg-white rounded-lg shadow p-6">
    <div className="flex items-center justify-between mb-4">
      <span className="text-3xl">{icon}</span>
      <span className={`text-sm font-semibold ${
        parseFloat(trend) > 0 ? 'text-green-500' : 'text-red-500'
      }`}>
        {trend}
      </span>
    </div>
    <h3 className="text-gray-500 text-sm font-medium">{title}</h3>
    <p className="text-2xl font-bold mt-2">{value}</p>
  </div>
);

const ChartCard: React.FC<{ title: string; children: React.ReactNode }> = ({
  title,
  children
}) => (
  <div className="bg-white rounded-lg shadow p-6">
    <h3 className="text-xl font-semibold mb-4">{title}</h3>
    {children}
  </div>
);

export default FleetManagerDashboard;