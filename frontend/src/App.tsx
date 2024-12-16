import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import LandingPage from './pages/LandingPage.tsx';
import LoginPage from './pages/LoginPage.tsx';
import DriverDashboard from './pages/DriverDashboard.tsx';
import CustomerDashboard from './pages/CustomerDashboard.tsx';
import AdminDashboard from './pages/AdminDashboard.tsx';
import ManagerDashboard from './pages/ManagerDashboard.tsx';
import LoaderDashboard from './pages/LoaderDashboard.tsx';

const App: React.FC = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/dashboard/driver" element={<DriverDashboard />} />
        <Route path="/dashboard/loader" element={<LoaderDashboard />} />
        <Route path="/dashboard/customer" element={<CustomerDashboard />} />
        <Route path="/dashboard/admin" element={<AdminDashboard />} />
        <Route path="/dashboard/manager" element={<ManagerDashboard />} />
      </Routes>
    </Router>
  );
};

export default App;