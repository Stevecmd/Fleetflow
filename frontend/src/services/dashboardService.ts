import api from '../utils/api';

export const DashboardService = {
  getStats: async () => {
    const response = await api.get('/api/v1/dashboard/stats');
    return response.data;
  },
  
  getFleetOverview: async () => {
    const response = await api.get('/api/v1/dashboard/fleet');
    return response.data;
  },
  
  getAlerts: async () => {
    const response = await api.get('/api/v1/dashboard/alerts');
    return response.data;
  }
};
