import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../features/auth/AuthProvider';

interface PrivateRouteProps {
  children: React.ReactNode;
}

/**
 * A Route that redirects to the login page if the user is not authenticated.
 *
 * The component will render a "Loading..." message if the authentication status is still loading.
 *
 * @example
 * <PrivateRoute>
 *   <ProtectedComponent />
 * </PrivateRoute>
 */
const PrivateRoute: React.FC<PrivateRouteProps> = ({ children }) => {
  const { isAuthenticated, loading } = useAuth();

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  return <>{children}</>;
};

export default PrivateRoute;
