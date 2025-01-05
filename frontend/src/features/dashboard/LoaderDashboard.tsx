import React, { useEffect, useState } from 'react';
import '../../config/chartConfig';
import { Bar, Doughnut } from "react-chartjs-2";

/**
 * A dashboard for loader users displaying vehicle and delivery statuses.
 *
 * @returns A component displaying the loader dashboard.
 */
const LoaderDashboard: React.FC = () => {
    const [loadingSchedules, setLoadingSchedules] = useState<any[]>([]);
    const [warehouseInventory, setWarehouseInventory] = useState<any[]>([]);
    const [equipmentStatus, setEquipmentStatus] = useState<any[]>([]);
    const [dailyStats, setDailyStats] = useState<any>({});

    const fetchDashboardData = async () => {
        // Loading schedules from loading_schedules table
        const scheduleDummyData = [
            { dock_number: 'Dock 1', vehicle_id: 1, scheduled_time: '08:00', status: 'completed', estimated_duration: 120 },
            { dock_number: 'Dock 2', vehicle_id: 2, scheduled_time: '09:30', status: 'in_progress', estimated_duration: 90 },
            { dock_number: 'Dock 3', vehicle_id: 3, scheduled_time: '11:00', status: 'pending', estimated_duration: 60 },
            { dock_number: 'Dock 4', vehicle_id: 4, scheduled_time: '13:00', status: 'scheduled', estimated_duration: 75 }
        ];

        // Inventory data from warehouse_inventory table
        const inventoryDummyData = [
            { item_name: 'Moving Boxes Large', quantity: 500, minimum_threshold: 100, item_category: 'Packaging' },
            { item_name: 'Packing Tape', quantity: 200, minimum_threshold: 50, item_category: 'Supplies' },
            { item_name: 'Bubble Wrap Roll', quantity: 150, minimum_threshold: 30, item_category: 'Packaging' },
            { item_name: 'Hand Truck', quantity: 20, minimum_threshold: 5, item_category: 'Equipment' }
        ];

        // Equipment status from warehouse_equipment table
        const equipmentDummyData = [
            { equipment_type: 'Forklift', status: 'operational', condition_rating: 4.5, equipment_id: 'EQ001' },
            { equipment_type: 'Pallet Jack', status: 'operational', condition_rating: 4.0, equipment_id: 'EQ002' },
            { equipment_type: 'Conveyor Belt', status: 'maintenance', condition_rating: 3.5, equipment_id: 'EQ003' },
            { equipment_type: 'Hand Truck', status: 'operational', condition_rating: 4.8, equipment_id: 'EQ004' }
        ];

        // Daily statistics
        const dailyStatsDummy = {
            total_loadings: 24,
            completed_loadings: 18,
            pending_loadings: 6,
            average_loading_time: 85,
            equipment_utilization: 78,
            inventory_alerts: 3
        };

        setTimeout(() => {
            setLoadingSchedules(scheduleDummyData);
            setWarehouseInventory(inventoryDummyData);
            setEquipmentStatus(equipmentDummyData);
            setDailyStats(dailyStatsDummy);
        }, 1000);
    };

    useEffect(() => {
        fetchDashboardData();
    }, []);

    const CardMetric = ({ title, value, subtitle }: { title: string; value: string | number; subtitle?: string }) => (
        <div className="bg-white p-6 rounded-lg shadow-md">
            <h3 className="text-gray-500 text-sm font-medium">{title}</h3>
            <p className="text-3xl font-bold mt-2">{value}</p>
            {subtitle && <p className="text-gray-400 text-sm mt-2">{subtitle}</p>}
        </div>
    );

    const chartContainerStyle = {
        height: '300px',
        position: 'relative' as const,
        marginBottom: '1rem'
    };

    const chartOptions = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                position: 'bottom' as const
            }
        }
    };

    const equipmentStatusChart = {
        labels: equipmentStatus?.map(e => e.equipment_type) || [],
        datasets: [{
            label: 'Equipment Condition Rating',
            data: equipmentStatus?.map(e => e.condition_rating) || [],
            backgroundColor: ['#4CAF50', '#2196F3', '#FFC107', '#F44336'],
        }]
    };

    const inventoryChart = {
        labels: warehouseInventory?.map(i => i.item_name) || [],
        datasets: [{
            label: 'Current Quantity',
            data: warehouseInventory?.map(i => i.quantity) || [],
            backgroundColor: 'rgba(54, 162, 235, 0.5)',
            borderColor: 'rgba(54, 162, 235, 1)',
            borderWidth: 1,
        }, {
            label: 'Minimum Threshold',
            data: warehouseInventory?.map(i => i.minimum_threshold) || [],
            backgroundColor: 'rgba(255, 99, 132, 0.5)',
            borderColor: 'rgba(255, 99, 132, 1)',
            borderWidth: 1,
        }]
    };

    return (
        <div className="p-6 bg-gray-100 min-h-screen">
            <h1 className="text-3xl font-bold mb-8">Warehouse Operations Dashboard</h1>
            
            {/* Key Metrics */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
                <CardMetric 
                    title="Total Loadings Today" 
                    value={dailyStats.total_loadings || 0}
                    subtitle={`${dailyStats.completed_loadings || 0} completed`} 
                />
                <CardMetric 
                    title="Average Loading Time" 
                    value={`${dailyStats.average_loading_time || 0}min`}
                    subtitle="Target: 90min" 
                />
                <CardMetric 
                    title="Equipment Utilization" 
                    value={`${dailyStats.equipment_utilization || 0}%`}
                    subtitle="All equipment" 
                />
                <CardMetric 
                    title="Inventory Alerts" 
                    value={dailyStats.inventory_alerts || 0}
                    subtitle="Items below threshold" 
                />
            </div>

            {/* Loading Schedule */}
            <div className="bg-white p-6 rounded-lg shadow-md mb-8">
                <h2 className="text-xl font-semibold mb-4">Today's Loading Schedule</h2>
                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead>
                            <tr>
                                <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Dock</th>
                                <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Time</th>
                                <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Vehicle ID</th>
                                <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Duration</th>
                                <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {loadingSchedules.map((schedule, index) => (
                                <tr key={index}>
                                    <td className="px-6 py-4 whitespace-nowrap">{schedule.dock_number}</td>
                                    <td className="px-6 py-4 whitespace-nowrap">{schedule.scheduled_time}</td>
                                    <td className="px-6 py-4 whitespace-nowrap">{schedule.vehicle_id}</td>
                                    <td className="px-6 py-4 whitespace-nowrap">{schedule.estimated_duration}min</td>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full 
                                            ${schedule.status === 'completed' ? 'bg-green-100 text-green-800' : 
                                            schedule.status === 'in_progress' ? 'bg-blue-100 text-blue-800' : 
                                            'bg-yellow-100 text-yellow-800'}`}>
                                            {schedule.status}
                                        </span>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>

            {/* Charts Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="bg-white p-6 rounded-lg shadow-md">
                    <h2 className="text-xl font-semibold mb-4">Equipment Status</h2>
                    <div style={chartContainerStyle}>
                        {equipmentStatus && (
                            <Doughnut 
                                data={equipmentStatusChart} 
                                options={chartOptions}
                            />
                        )}
                    </div>
                </div>
                
                <div className="bg-white p-6 rounded-lg shadow-md">
                    <h2 className="text-xl font-semibold mb-4">Inventory Levels</h2>
                    <div style={chartContainerStyle}>
                        {warehouseInventory && (
                            <Bar 
                                data={inventoryChart} 
                                options={{
                                    ...chartOptions,
                                    scales: {
                                        y: {
                                            beginAtZero: true
                                        }
                                    }
                                }}
                            />
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default LoaderDashboard;