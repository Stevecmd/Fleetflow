import React, { createContext, useContext, useEffect, useCallback } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { RootState, AppDispatch } from '../../store';
import {
  refreshTokenSuccess,
  logout,
  loginStart,
  loginSuccess,
  loginFailure
} from './authSlice';
import { api } from '../../services/api';
import { User } from '../../types';
import { fetchVehicleStats, fetchOrders, clearDashboard } from '../dashboard/dashboardSlice';
import { jwtDecode } from 'jwt-decode';

interface AuthContextType {
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
  user: User | null;
  login: (credentials: { identifier: string; password: string }) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const dispatch = useDispatch<AppDispatch>();
  const navigate = useNavigate();
  const { isAuthenticated, loading, error, token, refreshToken, user } = useSelector(
    (state: RootState) => state.auth
  );

  const refreshAuthToken = useCallback(async () => {
    try {
      const response = await api.post('/auth/refresh', { refreshToken });
      dispatch(refreshTokenSuccess(response.data));
    } catch (error) {
      console.error('Token refresh error:', error);
      dispatch(logout());
    }
  }, [refreshToken, dispatch]);

  useEffect(() => {
    if (token) {
      api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    } else {
      delete api.defaults.headers.common['Authorization'];
    }
  }, [token]);

  useEffect(() => {
    let refreshInterval: NodeJS.Timeout;

    if (isAuthenticated && refreshToken) {
      refreshInterval = setInterval(refreshAuthToken, 15 * 60 * 1000); // Refresh every 15 minutes
    }

    return () => clearInterval(refreshInterval);
  }, [isAuthenticated, refreshToken, refreshAuthToken]);

  const login = async ({ identifier, password }: { identifier: string; password: string }) => {
    dispatch(loginStart());
    try {
      const response = await api.post('/auth/login', {
        username: identifier,
        password
      });

      console.log('Full login response:', response.data);  // Log full response for debugging
      
      // Destructure with fallbacks
      const {
        access_token: accessToken,
        refresh_token: refreshToken,
        user_id,
        user = {}, 
        username
      } = response.data;

      // Decode token to get role
      const decoded: any = jwtDecode(accessToken);
      console.log('Decoded token:', decoded);

      const completeUser: User = {
        id: user_id,
        first_name: user.first_name || username || identifier,
        last_name: user.last_name || '',
        email: user.email || '',
        // role: user.role || '',
        role: decoded.role_name,
        profile_image_url: user.profile_image_url,
      };

      console.log('Processed user object:', completeUser);
      
      // Store user in localStorage as well
      localStorage.setItem('user', JSON.stringify(completeUser));
      localStorage.setItem('accessToken', accessToken);
      localStorage.setItem('refreshToken', refreshToken);
      localStorage.setItem('isAuthenticated', 'true');

      dispatch(
        loginSuccess({
          user: completeUser,
          token: accessToken,
          refreshToken
        })
      );

      // Role-based navigation
      switch (decoded.role_name) {
        case 'driver':
          navigate('/dashboard/driver');
          break;
        case 'fleet_manager':
          navigate('/dashboard/fleetmanager'); 
          break;
        case 'admin':
          navigate('/dashboard/admin');
          break;
        case 'loader':
          navigate('/dashboard/loader');
          break;
        case 'manager':
          navigate('/dashboard/manager');
          break;
        default:
          navigate('/login');
      }

      // Fetch additional data after login
      console.log('User ID:', user_id);
      // if (user_id) {
      if (decoded.role_name === 'driver') {
        dispatch(fetchVehicleStats(user_id)); // Fetch vehicle stats
        dispatch(fetchOrders(user_id)); // Fetch orders
      }
    } catch (error: any) {
      console.error('Login error:', error);
      dispatch(
        loginFailure(
          error.response?.data?.message || error.message || 'Login failed'
        )
      );
    }
  };

  useEffect(() => {
    console.log('isAuthenticated:', isAuthenticated);
  }, [isAuthenticated]);


    // Update handleLogout to clear localStorage
    const handleLogout = useCallback(async () => {
      try {
        if (token) {
          await api.post('/auth/logout', { refreshToken });
        }
      } catch (error) {
        console.error('Logout error:', error);
      } finally {
        localStorage.removeItem('user');
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
        localStorage.removeItem('isAuthenticated');
        
        dispatch(clearDashboard());
        dispatch(logout());
        navigate('/login');
      }
    }, [token, refreshToken, dispatch, navigate]);


  // Update the useEffect for checking stored auth state
  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    const storedToken = localStorage.getItem('accessToken');
    const storedRefreshToken = localStorage.getItem('refreshToken');
    const storedIsAuthenticated = localStorage.getItem('isAuthenticated');

    console.log('Stored user:', storedUser);
    console.log('Stored token:', storedToken);

    if (storedUser && storedToken && storedRefreshToken && storedIsAuthenticated === 'true') {
      try {
        const parsedUser = JSON.parse(storedUser);
        dispatch(
          loginSuccess({
            user: parsedUser,
            token: storedToken,
            refreshToken: storedRefreshToken
          })
        );
      } catch (error) {
        console.error('Error parsing stored user:', error);
        handleLogout(); // Clear invalid state
      }
    }
  }, [dispatch, handleLogout]);

  return (
    <AuthContext.Provider
      value={{ isAuthenticated, loading, error, user, login, logout: handleLogout }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

