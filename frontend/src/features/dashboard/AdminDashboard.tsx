import React, { useState, useEffect } from 'react';
import { useAuth } from '../auth/AuthProvider';
import { api } from '../../services/api'; // Import the api instance
import { Doughnut } from 'react-chartjs-2';
import '../../config/chartConfig';

/**
 * The admin dashboard component.
 *
 * This component displays a table of all users in the system, and a form to register a new user.
 * It also displays a message indicating whether the user registration was successful or not.
 *
 * @returns The admin dashboard component.
 */
const AdminDashboard: React.FC = () => {
    const { user } = useAuth();
    const [newUser, setNewUser] = useState({
        username: '',
        password: '',
        email: '',
        role_id: 5, // Default role (e.g., customer)
        first_name: '',
        last_name: '',
        phone: '',
    });
    const [users, setUsers] = useState([]);
    const [message, setMessage] = useState('');
    const [systemStats, setSystemStats] = useState({
        usersByRole: {},
        activeUsers: 0,
        inactiveUsers: 0,
        recentRegistrations: [],
        systemHealth: {
            databaseStatus: 'healthy',
            apiStatus: 'operational',
            lastBackup: '2024-01-20',
            serverLoad: 45
        }
    });

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
        }
    };

    useEffect(() => {
        fetchUsers();
        fetchSystemStats();
    }, []);

/**
 * Fetches the list of users from the API and updates the state.
 * If an error occurs during the fetch operation, it logs the error to the console.
 */

    const fetchUsers = async () => {
        try {
            const response = await api.get('/users'); // Use the api instance
            setUsers(response.data); // Assuming response.data contains the user list
        } catch (error) {
            console.error('Error fetching users:', error);
        }
    };

    const fetchSystemStats = async () => {
        try {
            // Use existing users endpoint for now
            const response = await api.get('/users');
            const users = response.data;

            // Process users data for statistics
            const roleStats = users.reduce((acc: any, user: any) => {
                acc[user.role_name] = (acc[user.role_name] || 0) + 1;
                return acc;
            }, {});

            const activeCount = users.filter((u: any) => u.status === 'active').length;
            const recentUsers = users
                .sort((a: any, b: any) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
                .slice(0, 5);

            setSystemStats({
                ...systemStats,
                usersByRole: roleStats,
                activeUsers: activeCount,
                inactiveUsers: users.length - activeCount,
                recentRegistrations: recentUsers
            });
        } catch (error) {
            console.error('Error fetching system stats:', error);
        }
    };

    /**
     * Handles input changes in the registration form.
     * @param {React.ChangeEvent<HTMLInputElement>} e The input change event.
     */
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setNewUser({ ...newUser, [name]: value });
    };

    /**
     * Handles the form submission and registers a new user.
     * @param {React.FormEvent} e The form submission event.
     */
    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        console.log('Registering user:', newUser);
        try {
            const response = await api.post('/auth/register', newUser); // Use the api instance
            if (response.status === 200) {
                setMessage('User registered successfully!');
                fetchUsers(); // Refresh the user list
                setNewUser({ username: '', password: '', email: '', role_id: 5, first_name: '', last_name: '', phone: '' }); // Reset form
            } else {
                setMessage('Error registering user.');
            }
        } catch (error) {
            console.error('Error registering user:', error);
            setMessage('Error registering user.');
        }
    };

    /**
     * Deletes a user with the given ID and refreshes the user list.
     * @param {number} userId The ID of the user to delete.
     */
    const handleDelete = async (userId: number) => {
        const token = localStorage.getItem('token'); // Retrieve the token from local storage
        if (window.confirm('Are you sure you want to delete this user?')) {
            try {
                const response = await api.delete(`/users/${userId}`, {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                });
                if (response.status === 200) {
                    setMessage('User deleted successfully!');
                    fetchUsers(); // Refresh the user list
                } else {
                    setMessage('Error deleting user.');
                }
            } catch (error) {
                console.error('Error deleting user:', error);
                setMessage('Error deleting user.');
            }
        }
    };

    const userRolesChart = {
        labels: Object.keys(systemStats.usersByRole),
        datasets: [{
            data: Object.values(systemStats.usersByRole),
            backgroundColor: ['#4CAF50', '#2196F3', '#FFC107', '#F44336', '#9C27B0', '#FF9800'],
        }]
    };

    return (
        <div className="p-6 bg-gray-100">
            <div className="mb-8">
                <h1 className="text-3xl font-bold mb-6">System Overview</h1>
                
                {/* System Health Metrics */}
                <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
                    <div className="bg-white p-6 rounded-lg shadow-md">
                        <h3 className="text-gray-500 text-sm font-medium">Total Users</h3>
                        <p className="text-3xl font-bold mt-2">{users.length}</p>
                        <p className="text-gray-400 text-sm mt-2">Across all roles</p>
                    </div>
                    <div className="bg-white p-6 rounded-lg shadow-md">
                        <h3 className="text-gray-500 text-sm font-medium">Active Users</h3>
                        <p className="text-3xl font-bold mt-2">{systemStats.activeUsers}</p>
                        <p className="text-gray-400 text-sm mt-2">Currently active</p>
                    </div>
                    <div className="bg-white p-6 rounded-lg shadow-md">
                        <h3 className="text-gray-500 text-sm font-medium">System Status</h3>
                        <p className="text-3xl font-bold mt-2 text-green-500">Healthy</p>
                        <p className="text-gray-400 text-sm mt-2">All systems operational</p>
                    </div>
                    <div className="bg-white p-6 rounded-lg shadow-md">
                        <h3 className="text-gray-500 text-sm font-medium">Server Load</h3>
                        <p className="text-3xl font-bold mt-2">{systemStats.systemHealth.serverLoad}%</p>
                        <p className="text-gray-400 text-sm mt-2">Current usage</p>
                    </div>
                </div>

                {/* Charts Grid */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
                    <div className="bg-white p-6 rounded-lg shadow-md">
                        <h2 className="text-xl font-semibold mb-4">Users by Role</h2>
                        <div style={chartContainerStyle}>
                            <Doughnut data={userRolesChart} options={chartOptions} />
                        </div>
                    </div>
                    
                    <div className="bg-white p-6 rounded-lg shadow-md">
                        <h2 className="text-xl font-semibold mb-4">Recent Registrations</h2>
                        <div className="space-y-4">
                            {systemStats.recentRegistrations.map((user: any) => (
                                <div key={user.id} className="flex justify-between items-center">
                                    <div>
                                        <p className="font-semibold">{user.username}</p>
                                        <p className="text-sm text-gray-500">{user.role_name}</p>
                                    </div>
                                    <p className="text-sm text-gray-400">
                                        {new Date(user.created_at).toLocaleDateString()}
                                    </p>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>

            <h1 className="text-2xl font-bold mb-4">Register a new user.</h1>
            {message && <div className="mb-4 text-green-600">{message}</div>}
            <form onSubmit={handleSubmit} className="mb-6 space-y-4">
                <div className="flex flex-col">
                    <input
                        type="text"
                        name="username"
                        placeholder="Username"
                        onChange={handleChange}
                        required
                        className="p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    />
                </div>
                <div className="flex flex-col">
                    <input
                        type="password"
                        name="password"
                        placeholder="Password"
                        onChange={handleChange}
                        required
                        className="p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    />
                </div>
                <div className="flex flex-col">
                    <input
                        type="email"
                        name="email"
                        placeholder="Email"
                        onChange={handleChange}
                        required
                        className="p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    />
                </div>
                <div className="flex flex-col">
                    <input
                        type="text"
                        name="first_name"
                        placeholder="First Name"
                        onChange={handleChange}
                        required
                        className="p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    />
                </div>
                <div className="flex flex-col">
                    <input
                        type="text"
                        name="last_name"
                        placeholder="Last Name"
                        onChange={handleChange}
                        required
                        className="p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    />
                </div>
                <div className="flex flex-col">
                    <input
                        type="text"
                        name="phone"
                        placeholder="Phone"
                        onChange={handleChange}
                        required
                        className="p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    />
                </div>
                <button
                    type="submit"
                    className="w-full p-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 transition duration-200"
                >
                    Register User
                </button>
            </form>

            <h2 className="text-lg font-semibold">Registered Users</h2>
            <table className="min-w-full border-collapse border border-gray-200">
                <thead>
                    <tr>
                        <th className="border border-gray-300 p-2">Username</th>
                        <th className="border border-gray-300 p-2">Email</th>
                        <th className="border border-gray-300 p-2">Role</th>
                        <th className="border border-gray-300 p-2">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {users.map((user: { id: number, username: string, email: string, role_name: string }) => (
                        <tr key={user.id}>
                            <td className="border border-gray-300 p-2">{user.username}</td>
                            <td className="border border-gray-300 p-2">{user.email}</td>
                            <td className="border border-gray-300 p-2">{user.role_name}</td>
                            <td className="border border-gray-300 p-2">
                                <button onClick={() => handleDelete(user.id)}>Delete</button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default AdminDashboard;