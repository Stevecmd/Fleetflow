import React, { createContext, useContext, useEffect, useCallback } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { RootState, AppDispatch } from '../../store';
import { refreshTokenSuccess, logout, loginStart, loginSuccess, loginFailure } from './authSlice';
import { api } from '../../services/api';
import { User } from '../../types';
import { fetchVehicleStats, fetchOrders } from '../dashboard/dashboardSlice';

interface AuthContextType {
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
  user: User | null;
  login: (identifier: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const dispatch = useDispatch<AppDispatch>();
  const navigate = useNavigate();
  const { isAuthenticated, loading, error, token, refreshToken, user } = useSelector(
    (state: RootState) => state.auth
  );

  const handleLogout = useCallback(async () => {
    try {
      if (token) {
        await api.post('/auth/logout', { refreshToken });
      }
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      dispatch(logout());
      navigate('/login');
    }
  }, [token, refreshToken, dispatch, navigate]);

  useEffect(() => {
    if (token) {
      api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    } else {
      delete api.defaults.headers.common['Authorization'];
    }
  }, [token]);

  useEffect(() => {
    let refreshInterval: NodeJS.Timeout;

    const refreshAuthToken = async () => {
      try {
        if (refreshToken) {
          const response = await api.post('/auth/refresh', { refreshToken });
          dispatch(refreshTokenSuccess(response.data));
        }
      } catch (error) {
        console.error('Token refresh error:', error);
        dispatch(logout());
      }
    };

    if (isAuthenticated) {
      refreshInterval = setInterval(refreshAuthToken, 15 * 60 * 1000); // Refresh every 15 minutes
    }

    return () => clearInterval(refreshInterval);
  }, [isAuthenticated, refreshToken, dispatch]);

  const login = async (identifier: string, password: string) => {
    dispatch(loginStart());
    try {
      const response = await api.post('/auth/login', { 
        username: identifier,  
        password 
      });
      
      console.log('Full login response:', response.data);  // Log full response for debugging
      
      // Destructure with fallbacks
      const { 
        access_token, 
        refresh_token, 
        user_id, 
        user = {}, 
        username 
      } = response.data;
      
      // Ensure user object has required properties
      const completeUser = {
        id: user_id,
        user_id: user_id,
        first_name: user.first_name || username || identifier,
        last_name: user.last_name || '',
        role: user.role || '',
        profile_image_url: user.profile_image_url,
        email: user.email || ''
      };
      
      console.log('Processed user object:', completeUser);
      
      // Store user in localStorage as well
      localStorage.setItem('user', JSON.stringify(completeUser));
      localStorage.setItem('accessToken', access_token);
      localStorage.setItem('refreshToken', refresh_token);
      
      dispatch(loginSuccess({
        user: completeUser,
        token: access_token,
        refreshToken: refresh_token
      }));
      
      navigate('/');

      // Fetch additional data after login
      console.log('User ID:', user_id);
      if (user_id) {
        dispatch(fetchVehicleStats(user_id)); // Fetch vehicle stats
        dispatch(fetchOrders(user_id)); // Fetch orders
      }
    } catch (error: any) {
      console.error('Full login error:', error);
      console.error('Error response:', error.response?.data);
      dispatch(loginFailure(
        error.response?.data?.message || 
        error.message || 
        'Login failed'
      ));
    }
  };

  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    const storedToken = localStorage.getItem('accessToken');
    const storedRefreshToken = localStorage.getItem('refreshToken');

    console.log('Stored user:', storedUser);
    console.log('Stored token:', storedToken);

    if (storedUser && storedToken && storedRefreshToken) {
      try {
        const parsedUser = JSON.parse(storedUser);
        dispatch(loginSuccess({
          user: parsedUser,
          token: storedToken,
          refreshToken: storedRefreshToken
        }));
      } catch (error) {
        console.error('Error parsing stored user:', error);
      }
    }
  }, [dispatch]);

  return (
    <AuthContext.Provider
      value={{ isAuthenticated, loading, error, user, login, logout: handleLogout }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
