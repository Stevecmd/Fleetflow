import React from 'react';
import { useAuth } from '../auth/AuthProvider';

const CustomerDashboard: React.FC = () => {
    const { user } = useAuth();
    
    return (
        <div className="p-6">
            <h1 className="text-2xl font-bold mb-4">
                Welcome, {user?.first_name || 'Customer'}!
            </h1>
            <div className="grid gap-4">
                {/* Add customer dashboard content here */}
                <div className="bg-white p-4 rounded shadow">
                    <h2 className="text-lg font-semibold">My Deliveries</h2>
                    {/* Add delivery tracking components */}
                </div>
            </div>
        </div>
    );
};
export default CustomerDashboard;
