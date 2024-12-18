import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './features/auth/AuthProvider';
import { useAuth } from './features/auth/AuthProvider';
import PrivateRoute from './components/PrivateRoute';
import Layout from './components/Layout';
import Login from './features/auth/Login';
import LandingPage from './features/dashboard/LandingPage';
import DriverDashboard from './features/dashboard/DriverDashboard';
import LoaderDashboard from './features/dashboard/LoaderDashboard';
import ManagerDashboard from './features/dashboard/ManagerDashboard';
import FleetManagerDashboard from './features/dashboard/FleetManagerDashboard';

// Role-based dashboard routing
const DashboardRouter: React.FC = () => {
  const { user } = useAuth();

  switch (user?.role) {
    case 'driver':
      return <Navigate to="/dashboard/driver" replace />;
    case 'loader':
      return <Navigate to="/dashboard/loader" replace />;
      case 'admin':
        return <Navigate to="/dashboard/admin" replace />;
      case 'manager':
        return <Navigate to="/dashboard/manager" replace />;
      default:
      return <Navigate to="/login" replace />;
  }
};

const App: React.FC = () => {
  return (
    <AuthProvider>
      <Routes>
        {/* Public routes */}
        <Route path="/" element={<LandingPage />} />
        <Route path="/login" element={<Login />} />
        
        {/* Protected dashboard routes */}
        <Route
          path="/dashboard"
          element={
            <PrivateRoute>
              <Layout />
            </PrivateRoute>
          }
        >
          <Route index element={<DashboardRouter />} />
          <Route path="fleetmanager" element={<FleetManagerDashboard />} />
          <Route path="driver" element={<DriverDashboard />} />
          <Route path="manager" element={<ManagerDashboard />} />
          <Route path="loader" element={<LoaderDashboard />} />
          {/* Add more role-based dashboard routes here */}
        </Route>

        {/* Catch all route */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </AuthProvider>
  );
};

export default App;
