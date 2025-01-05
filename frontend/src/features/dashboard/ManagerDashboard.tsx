import React, { useEffect, useState } from 'react';
import '../../config/chartConfig';
import { Bar, Line, Pie, Doughnut } from 'react-chartjs-2';

/**
 * ManagerDashboard is a functional component that renders the main interface
 * for managers within the application. It includes a header and serves as a
 * placeholder for manager-specific components and functionality.
 */

const ManagerDashboard: React.FC = () => {
  const [revenueData, setRevenueData] = useState<{ month: string; amount: number; deliveries: number }[] | null>(null);
  const [deliveryStats, setDeliveryStats] = useState<{ status: string; count: number }[] | null>(null);
  const [fleetStatus, setFleetStatus] = useState<{ category: string; count: number }[] | null>(null);
  const [performanceData, setPerformanceData] = useState<{ metric: string; value: number }[] | null>(null);

  const fetchDashboardData = async () => {
    // Revenue data from invoices table
    const revenueDummyData = [
      { month: 'Jan', amount: 45000, deliveries: 180 },
      { month: 'Feb', amount: 52000, deliveries: 210 },
      { month: 'Mar', amount: 49000, deliveries: 195 },
      { month: 'Apr', amount: 58000, deliveries: 225 },
      { month: 'May', amount: 63000, deliveries: 245 },
      { month: 'Jun', amount: 68000, deliveries: 260 }
    ];

    // Delivery statistics from delivery_statuses table
    const deliveryDummyData = [
      { status: 'Pending', count: 65 },
      { status: 'In Transit', count: 80 },
      { status: 'Delivered', count: 245 },
      { status: 'Cancelled', count: 12 }
    ];

    // Fleet status from vehicle_statuses table
    const fleetDummyData = [
      { category: 'Available', count: 85 },
      { category: 'In Maintenance', count: 12 },
      { category: 'On Route', count: 28 },
      { category: 'Out of Service', count: 5 }
    ];

    // Performance metrics from driver_performance_metrics and fleet_analytics tables
    const performanceDummyData = [
      { metric: 'On-Time Delivery Rate', value: 94.5 },
      { metric: 'Average Customer Rating', value: 4.8 },
      { metric: 'Fleet Utilization', value: 82.3 },
      { metric: 'Safety Score', value: 98.1 },
      { metric: 'Fuel Efficiency', value: 15.2 }
    ];

    setTimeout(() => {
      setRevenueData(revenueDummyData);
      setDeliveryStats(deliveryDummyData);
      setFleetStatus(fleetDummyData);
      setPerformanceData(performanceDummyData);
    }, 1000);
  };

  useEffect(() => {
    fetchDashboardData();
  }, []);

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

  // Add specific options for bar chart
  const barChartOptions = {
    ...chartOptions,
    scales: {
      y: {
        beginAtZero: true,
        ticks: {
          maxRotation: 0
        }
      }
    }
  };

  const revenueChart = {
    labels: revenueData?.map(d => d.month) || [],
    datasets: [{
      label: 'Monthly Revenue',
      data: revenueData?.map(d => d.amount) || [],
      borderColor: '#4CAF50',
      backgroundColor: 'rgba(76, 175, 80, 0.1)',
      fill: true,
      tension: 0.4
    }]
  };

  const deliveryChart = {
    labels: deliveryStats?.map(d => d.status) || [],
    datasets: [{
      data: deliveryStats?.map(d => d.count) || [],
      backgroundColor: ['#2196F3', '#FFC107', '#9C27B0', '#F44336'],
    }]
  };

  const fleetChart = {
    labels: fleetStatus?.map(f => f.category) || [],
    datasets: [{
      data: fleetStatus?.map(f => f.count) || [],
      backgroundColor: ['#66BB6A', '#FFA726', '#42A5F5', '#EF5350'],
    }]
  };

  const performanceChart = {
    labels: performanceData?.map(p => p.metric) || [],
    datasets: [{
      label: 'Performance Metrics',
      data: performanceData?.map(p => p.value) || [],
      backgroundColor: 'rgba(156, 39, 176, 0.6)',
    }]
  };

  const CardMetric = ({ title, value, subtitle }: { title: string; value: string | number; subtitle?: string }) => (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h3 className="text-gray-500 text-sm font-medium">{title}</h3>
      <p className="text-3xl font-bold mt-2">{value}</p>
      {subtitle && <p className="text-gray-400 text-sm mt-2">{subtitle}</p>}
    </div>
  );

  return (
    <div className="p-6 bg-gray-100 min-h-screen">
      <h1 className="text-3xl font-bold mb-8">Management Dashboard</h1>
      
      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <CardMetric 
          title="Monthly Revenue" 
          value="$335,000" 
          subtitle="260 deliveries completed" 
        />
        <CardMetric 
          title="Active Fleet" 
          value="130" 
          subtitle="85 vehicles available" 
        />
        <CardMetric 
          title="Warehouse Capacity" 
          value="82%" 
          subtitle="37,000 sq ft utilized" 
        />
        <CardMetric 
          title="Driver Performance" 
          value="94.5%" 
          subtitle="On-time delivery rate" 
        />
      </div>

      {/* Additional Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h3 className="text-gray-500 text-sm font-medium">Warehouse Inventory</h3>
          <div className="mt-2">
            <div className="flex justify-between items-center mb-2">
              <span>Moving Boxes</span>
              <span className="text-green-500">500 units</span>
            </div>
            <div className="flex justify-between items-center mb-2">
              <span>Packing Tape</span>
              <span className="text-yellow-500">200 rolls</span>
            </div>
            <div className="flex justify-between items-center">
              <span>Bubble Wrap</span>
              <span className="text-blue-500">150 rolls</span>
            </div>
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-md">
          <h3 className="text-gray-500 text-sm font-medium">Equipment Status</h3>
          <div className="mt-2">
            <div className="flex justify-between items-center mb-2">
              <span>Forklifts</span>
              <span className="text-green-500">8 operational</span>
            </div>
            <div className="flex justify-between items-center mb-2">
              <span>Pallet Jacks</span>
              <span className="text-yellow-500">12 available</span>
            </div>
            <div className="flex justify-between items-center">
              <span>Hand Trucks</span>
              <span className="text-blue-500">20 in service</span>
            </div>
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-md">
          <h3 className="text-gray-500 text-sm font-medium">Maintenance Schedule</h3>
          <div className="mt-2">
            <div className="flex justify-between items-center mb-2">
              <span>Vehicles Due</span>
              <span className="text-yellow-500">12 vehicles</span>
            </div>
            <div className="flex justify-between items-center mb-2">
              <span>Completed Today</span>
              <span className="text-green-500">5 services</span>
            </div>
            <div className="flex justify-between items-center">
              <span>Urgent Repairs</span>
              <span className="text-red-500">3 pending</span>
            </div>
          </div>
        </div>
      </div>

      {/* Charts Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="bg-white p-6 rounded-lg shadow-md flex flex-col">
          <h2 className="text-xl font-semibold mb-4">Revenue Trend</h2>
          <div style={chartContainerStyle}>
            {revenueData && (
              <Line 
                data={revenueChart} 
                options={chartOptions} 
                className="max-w-full"
              />
            )}
          </div>
        </div>
        
        <div className="bg-white p-6 rounded-lg shadow-md flex flex-col">
          <h2 className="text-xl font-semibold mb-4">Delivery Status</h2>
          <div style={chartContainerStyle}>
            {deliveryStats && (
              <Doughnut 
                data={deliveryChart} 
                options={chartOptions}
                className="max-w-full"
              />
            )}
          </div>
        </div>
        
        <div className="bg-white p-6 rounded-lg shadow-md flex flex-col">
          <h2 className="text-xl font-semibold mb-4">Fleet Status</h2>
          <div style={chartContainerStyle}>
            {fleetStatus && (
              <Pie 
                data={fleetChart} 
                options={chartOptions}
                className="max-w-full"
              />
            )}
          </div>
        </div>
        
        <div className="bg-white p-6 rounded-lg shadow-md flex flex-col">
          <h2 className="text-xl font-semibold mb-4">Performance Metrics</h2>
          <div style={chartContainerStyle}>
            {performanceData && (
              <Bar 
                data={performanceChart} 
                options={barChartOptions}
                className="max-w-full"
              />
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default ManagerDashboard;
