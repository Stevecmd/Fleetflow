import React, { useEffect, useState } from 'react';
import { Bar, Pie } from "react-chartjs-2";
import axios from "axios";

/**
 * A dashboard for loader users displaying vehicle and delivery statuses.
 *
 * @returns A component displaying the loader dashboard.
 */
const LoaderDashboard: React.FC = () => {
    const [vehicleData, setVehicleData] = useState<{ status: string; count: number }[] | null>(null);
    const [deliveryData, setDeliveryData] = useState<{ status: string; count: number }[] | null>(null);

    const fetchData = async () => {
        // Dummy data for vehicle statuses
        const vehiclesDummyData = [
            { status: 'Active', count: 10 },
            { status: 'Inactive', count: 5 },
            { status: 'Maintenance', count: 3 },
        ];

        // Dummy data for delivery statuses
        const deliveriesDummyData = [
            { status: 'Delivered', count: 15 },
            { status: 'In Transit', count: 7 },
            { status: 'Pending', count: 2 },
        ];

        // Simulate a delay to mimic an API call
        setTimeout(() => {
            setVehicleData(vehiclesDummyData);
            setDeliveryData(deliveriesDummyData);
        }, 1000);
    };

    useEffect(() => {
        fetchData();
    }, []);

    const vehicleStatusChart = {
        labels: vehicleData ? vehicleData.map((v) => v.status) : [],
        datasets: [{
            data: vehicleData ? vehicleData.map((v) => v.count) : [],
            backgroundColor: ['#FF6384', '#36A2EB', '#FFCE56'],
        }],
    };

    const deliveryStatusChart = {
        labels: deliveryData ? deliveryData.map((d) => d.status) : [],
        datasets: [{
            data: deliveryData ? deliveryData.map((d) => d.count) : [],
            backgroundColor: ['#4BC0C0', '#FF9F40', '#9966FF'],
        }],
    };

    return (
        <div style={{ padding: '20px', backgroundColor: '#f9f9f9' }}>
            <h1 style={{ textAlign: 'center', marginBottom: '20px' }}>Loader Dashboard</h1>
            <div style={{ display: 'flex', justifyContent: 'space-around', marginBottom: '20px' }}>
                <div style={{ backgroundColor: 'white', padding: '15px', borderRadius: '8px', boxShadow: '0 2px 5px rgba(0, 0, 0, 0.1)', flex: 1, margin: '0 10px', textAlign: 'center' }}>
                    <h2>Total Vehicles</h2>
                    <p>{vehicleData ? vehicleData.reduce((acc: number, v) => acc + v.count, 0) : 0}</p>
                </div>
                <div style={{ backgroundColor: 'white', padding: '15px', borderRadius: '8px', boxShadow: '0 2px 5px rgba(0, 0, 0, 0.1)', flex: 1, margin: '0 10px', textAlign: 'center' }}>
                    <h2>Total Deliveries</h2>
                    <p>{deliveryData ? deliveryData.reduce((acc: number, d) => acc + d.count, 0) : 0}</p>
                </div>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-around' }}>
                <div style={{ flex: 1, margin: '0 10px' }}>
                    <h2>Vehicle Status Distribution</h2>
                    {vehicleData && <Pie data={vehicleStatusChart} />}
                </div>
                <div style={{ flex: 1, margin: '0 10px' }}>
                    <h2>Delivery Status Distribution</h2>
                    {deliveryData && <Bar data={deliveryStatusChart} />}
                </div>
            </div>
        </div>
    );
};

export default LoaderDashboard;