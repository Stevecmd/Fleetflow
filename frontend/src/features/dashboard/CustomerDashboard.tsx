import React, { useEffect, useState } from 'react';
import { useAuth } from '../auth/AuthProvider';
import { api } from '../../services/api';


interface Delivery {
    id: number;
    tracking_number: string;
    status_id: number;
    estimated_delivery_time: string;
    actual_delivery_time?: string;
    special_instructions?: string;
    cargo_type: string;
    from_location: string;
    to_location: string;
    cargo_weight: number;
    payment_status: string;
    route_efficiency_score?: number;
    weather_conditions?: string;
    proof_of_delivery_image_url?: string;
    package_condition_images?: string[];
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

  interface DeliveryStats {
    total: number;
    onTime: number;
    delayed: number;
    completed: number;
    pending: number;
}

/**
 * CustomerDashboard component displays the dashboard for a logged-in customer.
 * It fetches and displays the customer's active deliveries, feedback history, and invoices.
 * 
 * State:
 * - deliveries: Array of active deliveries for the customer.
 * - invoices: Array of invoices related to the customer's deliveries.
 * - feedback: Array of feedback related to the deliveries.
 * - loading: Boolean indicating if data is being fetched.
 * - error: Error message if data fetching fails.
 * 
 * Effects:
 * - Fetches customer data (deliveries, feedback, invoices) on component mount or when user ID changes.
 * 
 * Conditional Rendering:
 * - Shows loading indicator while fetching data.
 * - Shows error message if fetching fails.
 * - Displays deliveries, feedback, and invoices if available.
 */

  const CustomerDashboard: React.FC = () => {
    const { user } = useAuth();
    const [deliveries, setDeliveries] = useState<Delivery[]>([]);
    const [invoices, setInvoices] = useState<Invoice[]>([]);
    const [feedback, setFeedback] = useState<DeliveryFeedback[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [deliveryStats, setDeliveryStats] = useState<DeliveryStats>({
        total: 0,
        onTime: 0,
        delayed: 0,
        completed: 0,
        pending: 0
    });

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

    useEffect(() => {
        // Calculate delivery statistics
        const stats = deliveries.reduce((acc, delivery) => {
            acc.total++;
            if (delivery.actual_delivery_time) {
                const actualDate = new Date(delivery.actual_delivery_time);
                const estimatedDate = new Date(delivery.estimated_delivery_time);
                acc.onTime += actualDate <= estimatedDate ? 1 : 0;
                acc.delayed += actualDate > estimatedDate ? 1 : 0;
            }
            acc.completed += delivery.status_id === 3 ? 1 : 0; // Assuming 3 is 'delivered' status
            acc.pending += delivery.status_id === 1 ? 1 : 0;  // Assuming 1 is 'pending' status
            return acc;
        }, {
            total: 0,
            onTime: 0,
            delayed: 0,
            completed: 0,
            pending: 0
        });
        
        setDeliveryStats(stats);
    }, [deliveries]);

    return (
        <div className="p-6 max-w-7xl mx-auto">
            <h1 className="text-2xl font-bold mb-6">
                Welcome, {user?.first_name}!
            </h1>

            {/* Add Statistics Overview */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
                <div className="bg-white rounded-lg shadow p-4">
                    <h3 className="text-gray-500 text-sm">Total Deliveries</h3>
                    <p className="text-2xl font-bold">{deliveryStats.total}</p>
                </div>
                <div className="bg-white rounded-lg shadow p-4">
                    <h3 className="text-gray-500 text-sm">On-Time Deliveries</h3>
                    <p className="text-2xl font-bold text-green-600">{deliveryStats.onTime}</p>
                </div>
                <div className="bg-white rounded-lg shadow p-4">
                    <h3 className="text-gray-500 text-sm">Pending Deliveries</h3>
                    <p className="text-2xl font-bold text-yellow-600">{deliveryStats.pending}</p>
                </div>
                <div className="bg-white rounded-lg shadow p-4">
                    <h3 className="text-gray-500 text-sm">Completed Deliveries</h3>
                    <p className="text-2xl font-bold text-blue-600">{deliveryStats.completed}</p>
                </div>
            </div>

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
                                        <div className="flex justify-between items-start">
                                            <div>
                                                <p className="font-medium">Tracking: {delivery.tracking_number}</p>
                                                <p className="text-gray-600">Cargo: {delivery.cargo_type}</p>
                                                <p className="text-gray-600">Weight: {delivery.cargo_weight}kg</p>
                                                <p className="text-gray-600">
                                                    Estimated Delivery: {new Date(delivery.estimated_delivery_time).toLocaleString()}
                                                </p>
                                                {delivery.special_instructions && (
                                                    <p className="text-amber-600 mt-2">
                                                        Note: {delivery.special_instructions}
                                                    </p>
                                                )}
                                            </div>
                                            <div className="text-right">
                                                <span className={`px-3 py-1 rounded-full text-sm ${
                                                    delivery.payment_status === 'paid' 
                                                        ? 'bg-green-100 text-green-800' 
                                                        : 'bg-yellow-100 text-yellow-800'
                                                }`}>
                                                    {delivery.payment_status}
                                                </span>
                                                {delivery.route_efficiency_score && (
                                                    <p className="text-sm text-gray-500 mt-2">
                                                        Efficiency: {Math.round(delivery.route_efficiency_score * 100)}%
                                                    </p>
                                                )}
                                            </div>
                                        </div>
                                        {delivery.proof_of_delivery_image_url && (
                                            <div className="mt-4">
                                                <img 
                                                    src={delivery.proof_of_delivery_image_url} 
                                                    alt="Proof of delivery" 
                                                    className="w-20 h-20 object-cover rounded"
                                                />
                                            </div>
                                        )}
                                    </div>
                                ))}
                            </div>
                        </div>
                        
                        {/* Feedback Section */}
                        {/* Only show these sections if data is available */}
                        {feedback && feedback.length > 0 && (
                            <div className="bg-white rounded-lg shadow p-6">
                                <h2 className="text-xl font-semibold mb-4">Feedback History</h2>
                                {/* Add feedback display */}
                            </div>
                        )}

                        {invoices && invoices.length > 0 && (
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
