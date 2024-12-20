import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { api } from '../../services/api';

interface Order {
  id: number;
  from: string;
  to: string;
  weight: number;
  status: 'Paid' | 'Unpaid';
}

interface DeliveryStats {
  onTime: number;
  unloading: number;
  waiting: number;
  totalDeliveries: number;
}

interface DashboardState {
  vehicleStats: VehicleData | null;
  orders: Order[];
  deliveryStats: DeliveryStats;
  loading: boolean;
  error: string | null;
  currentRole: string | null;
  currentUserRole: string | null;
  driverPerformance: any;
}

const initialState: DashboardState = {
  vehicleStats: null,
  orders: [],
  deliveryStats: {
    onTime: 0,
    unloading: 0,
    waiting: 0,
    totalDeliveries: 0,
  },
  loading: false,
  error: null,
  currentRole: null,
  currentUserRole: null,
  driverPerformance: null,
};

export interface VehicleData {
  id: number;
  plate_number: string;
  type: string;
  make: string;
  model: string;
  year: number;
  capacity: number;
  fuel_type: string;
  status_id: number;
  mileage: number;
  created_at: string;
  updated_at: string;
}

// Async thunks
export const fetchVehicleStats = createAsyncThunk(
  'dashboard/fetchVehicleStats',
  async (userId: number, { rejectWithValue }) => {
    try {
      const response = await api.get(`/drivers/${userId}/vehicle`);
      console.log('Vehicle Stats Response Data:', response.data); // Log the structured response data
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Failed to fetch vehicle stats');
    }
  }
);

export const fetchOrders = createAsyncThunk(
  'dashboard/fetchOrders',
  async (userId: number, { rejectWithValue }) => {
    try {
      const ordersResponse = await api.get(`/drivers/${userId}/orders`);
      return ordersResponse.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Failed to fetch orders');
    }
  }
);

export const fetchDriverPerformance = createAsyncThunk(
  'dashboard/fetchDriverPerformance',
  async (userId: number, { rejectWithValue }) => {
    try {
      const response = await api.get(`/drivers/${userId}/performance`);
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Failed to fetch driver performance');
    }
  }
);

export const fetchFleetVehicles = createAsyncThunk(
  'dashboard/fetchFleetVehicles',
  async (_, { rejectWithValue }) => {
    try {
      const response = await api.get('/fleet-manager/vehicles');
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Failed to fetch fleet vehicles');
    }
  }
);

export const fetchFleetPerformance = createAsyncThunk(
  'dashboard/fetchFleetPerformance',
  async (_, { rejectWithValue }) => {
    try {
      const response = await api.get('/fleet-manager/performance');
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Failed to fetch fleet performance');
    }
  }
);

// Reducers
const dashboardSlice = createSlice({
  name: 'dashboard',
  initialState,
  reducers: {
    clearDashboard: (state) => {
      state.vehicleStats = null;
      state.orders = [];
      state.deliveryStats = {
        onTime: 0,
        unloading: 0,
        waiting: 0,
        totalDeliveries: 0,
      };
      state.loading = false;
      state.error = null;
      state.currentRole = null;
      state.currentUserRole = null;
      state.driverPerformance = null;
    },
    setUserRole(state, action) {
      state.currentUserRole = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder.addCase(fetchVehicleStats.fulfilled, (state, action) => {
      state.vehicleStats = action.payload;
      state.loading = false; // Added to stop loading spinner
    });
    builder.addCase(fetchVehicleStats.pending, (state) => {
      state.loading = true;
    });
    builder.addCase(fetchVehicleStats.rejected, (state, action) => {
      state.error = action.error.message ?? null;
      state.loading = false;
    });
    builder.addCase(fetchOrders.fulfilled, (state, action) => {
      state.orders = action.payload;
      state.loading = false; // Added to stop loading spinner
    });
    builder.addCase(fetchOrders.pending, (state) => {
      state.loading = true;
    });
    builder.addCase(fetchOrders.rejected, (state, action) => {
      state.error = action.error.message ?? null;
      state.loading = false;
    });
    builder.addCase(fetchDriverPerformance.fulfilled, (state, action) => {
      state.driverPerformance = action.payload;
    });
    builder.addCase(fetchDriverPerformance.pending, (state) => {
      state.loading = true;
    });
    builder.addCase(fetchDriverPerformance.rejected, (state, action) => {
      state.error = action.error.message ?? null;
      state.loading = false;
    });
  },
});

export const { clearDashboard, setUserRole } = dashboardSlice.actions;
export default dashboardSlice.reducer;