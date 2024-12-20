import React, { useEffect, useState } from 'react';
import { useAuth } from '../auth/AuthProvider';
import { api } from '../../services/api';


interface Delivery {
    id: number;
    tracking_number: string;
    status_id: number;
    estimated_delivery_time: string;
    cargo_type: string;
    from_location: string;
    to_location: string;
    cargo_weight: number;
    payment_status: string;
    proof_of_delivery_image_url?: string;
  }
  
  interface DeliveryFeedback {
    id: number;
    delivery_id: number;
    rating: number;
    feedback_text: string;
    timeliness_rating: number;
    driver_rating: number;
    package_condition_rating: number;
  }
  
  interface Invoice {
    id: number;
    delivery_id: number;
    amount: number;
    status: string;
    due_date: string;
    payment_method?: string;
  }

  const CustomerDashboard: React.FC = () => {
    const { user } = useAuth();
    const [deliveries, setDeliveries] = useState<Delivery[]>([]);
    const [invoices, setInvoices] = useState<Invoice[]>([]);
    const [feedback, setFeedback] = useState<DeliveryFeedback[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchCustomerData = async () => {
            if (!user?.id) return;
            
            setLoading(true);
            setError(null);
            
            try {
                // Only fetch deliveries initially since it's working
                const deliveriesRes = await api.get(`/users/${user.id}/deliveries`);
                setDeliveries(deliveriesRes.data);

                // Add error handling for each request separately
                try {
                    const feedbackRes = await api.get(`/users/${user.id}/feedback`);
                    setFeedback(feedbackRes.data);
                } catch (feedbackError) {
                    console.log('Feedback fetch failed:', feedbackError);
                }

                try {
                    const invoicesRes = await api.get(`/users/${user.id}/invoices`);
                    setInvoices(invoicesRes.data);
                } catch (invoiceError) {
                    console.log('Invoice fetch failed:', invoiceError);
                }

            } catch (error: any) {
                console.error('Error fetching customer data:', error);
                setError(error.response?.data?.message || 'Failed to fetch customer data');
            } finally {
                setLoading(false);
            }
        };

        fetchCustomerData();
    }, [user?.id]);

    return (
        <div className="p-6 max-w-7xl mx-auto">
            <h1 className="text-2xl font-bold mb-6">
                Welcome, {user?.first_name}!
            </h1>

            {/* Active Deliveries Section */}
            <div className="grid gap-6 mb-8">
                {loading ? (
                    <div className="animate-pulse">Loading...</div>
                ) : error ? (
                    <div className="text-red-500">Error: {error}</div>
                ) : (
                    <>
                        {/* Deliveries Section */}
                        <div className="bg-white rounded-lg shadow p-6">
                            <h2 className="text-xl font-semibold mb-4">Active Deliveries</h2>
                            <div className="grid gap-4">
                                {deliveries.map(delivery => (
                                    <div key={delivery.tracking_number} 
                                         className="border p-4 rounded-lg hover:bg-gray-50">
                                        <div className="flex justify-between items-center">
                                            <div>
                                                <p className="font-medium">Tracking: {delivery.tracking_number}</p>
                                                <p className="text-gray-600">From: {delivery.from_location}</p>
                                                <p className="text-gray-600">Weight: {delivery.cargo_weight}kg</p>
                                                <p className="text-gray-600">
                                                    Estimated Delivery: {new Date(delivery.estimated_delivery_time).toLocaleDateString()}
                                                </p>
                                            </div>
                                            <span className={`px-3 py-1 rounded-full text-sm ${
                                                delivery.payment_status === 'paid' ? 'bg-green-100 text-green-800' :
                                                'bg-yellow-100 text-yellow-800'
                                            }`}>
                                                {delivery.payment_status || 'Pending'}
                                            </span>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>

                        {/* Only show these sections if data is available */}
                        {feedback.length > 0 && (
                            <div className="bg-white rounded-lg shadow p-6">
                                <h2 className="text-xl font-semibold mb-4">Feedback History</h2>
                                {/* Add feedback display */}
                            </div>
                        )}

                        {invoices.length > 0 && (
                            <div className="bg-white rounded-lg shadow p-6">
                                <h2 className="text-xl font-semibold mb-4">Invoices</h2>
                                {/* Add invoices display */}
                            </div>
                        )}
                    </>
                )}
            </div>
        </div>
    );
};

export default CustomerDashboard;
