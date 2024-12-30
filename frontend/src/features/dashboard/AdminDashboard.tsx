import React, { useState, useEffect } from 'react';
import { useAuth } from '../auth/AuthProvider';
import { api } from '../../services/api'; // Import the api instance

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

    useEffect(() => {
        fetchUsers();
    }, []);

    const fetchUsers = async () => {
        try {
            const response = await api.get('/users'); // Use the api instance
            setUsers(response.data); // Assuming response.data contains the user list
        } catch (error) {
            console.error('Error fetching users:', error);
        }
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setNewUser({ ...newUser, [name]: value });
    };

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

    return (
        <div className="p-6">
            <h1 className="text-2xl font-bold mb-4">
                Welcome, {user?.first_name || 'admin'}!
            </h1>
            <div className="grid gap-4">
                {/* Add admin dashboard content here */}
                <div className="bg-white p-4 rounded shadow">
                    <h2 className="text-lg font-semibold">All Deliveries</h2>
                    {/* Add delivery tracking components */}
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