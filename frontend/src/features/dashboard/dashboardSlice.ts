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
  driverMetrics: DriverMetrics | null;
  vehicleInfo: VehicleInfo | null;
  deliveries: Delivery[];
}

interface DriverMetrics {
  deliveries_completed: number;
  on_time_delivery_rate: number;
  customer_rating_avg: number;
  fuel_efficiency: number;
  safety_score: number;
  total_distance_covered: number;
}

interface VehicleInfo {
  plate_number: string;
  type: string;
  make: string;
  model: string;
  fuel_type: string;
  status: string;
  last_maintenance: string;
  next_maintenance: string;
}

interface Delivery {
  tracking_number: string;
  status: string;
  pickup_time: string;
  delivery_time: string;
  from_location: string;
  cargo_type: string;
  cargo_weight: number;
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
  driverMetrics: null,
  vehicleInfo: null,
  deliveries: [],
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

export const fetchDriverMetrics = createAsyncThunk(
  'dashboard/fetchDriverMetrics',
  async (userId: number, { rejectWithValue }) => {
    try {
      const response = await api.get(`/drivers/${userId}/performance`);
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Failed to fetch driver metrics');
    }
  }
);

export const fetchVehicleInfo = createAsyncThunk(
  'dashboard/fetchVehicleInfo',
  async (userId: number, { rejectWithValue }) => {
    try {
      const response = await api.get(`/drivers/${userId}/vehicle`);
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Failed to fetch vehicle info');
    }
  }
);

export const fetchDeliveries = createAsyncThunk(
  'dashboard/fetchDeliveries',
  async (userId: number, { rejectWithValue }) => {
    try {
      const response = await api.get(`/drivers/${userId}/orders`);
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Failed to fetch deliveries');
    }
  }
);

// Reducers
const dashboardSlice = createSlice({
  name: 'dashboard',
  initialState,
  reducers: {
/**
 * Resets the dashboard state to its initial values.
 * 
 * This function clears all the data related to the dashboard,
 * including vehicle statistics, orders, delivery stats,
 * loading state, error messages, roles, and driver performance.
 * It is typically used when the dashboard needs to be reset.
 */

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
      state.driverMetrics = null;
      state.vehicleInfo = null;
      state.deliveries = [];
    },
    setUserRole(state, action) {
      state.currentUserRole = action.payload;
    },
  },
  /**
   * This is an object that maps action types to reducers.
   * When an action is dispatched, the corresponding reducer
   * is called with the current state and the action payload.
   * The reducer should return a new state.
   *
   * The extraReducers property is used to add additional
   * reducers that are not part of the initial state.
   *
   * In this case, we are adding reducers for the following
   * actions:
   * - fetchVehicleStats.fulfilled
   * - fetchVehicleStats.pending
   * - fetchVehicleStats.rejected
   * - fetchOrders.fulfilled
   * - fetchOrders.pending
   * - fetchOrders.rejected
   * - fetchDriverPerformance.fulfilled
   * - fetchDriverPerformance.pending
   * - fetchDriverPerformance.rejected
   *
   * The reducers are called in the following order:
   * - fetchVehicleStats.fulfilled
   * - fetchOrders.fulfilled
   * - fetchDriverPerformance.fulfilled
   * - fetchVehicleStats.pending
   * - fetchOrders.pending
   * - fetchDriverPerformance.pending
   * - fetchVehicleStats.rejected
   * - fetchOrders.rejected
   * - fetchDriverPerformance.rejected
   *
   * If any of the reducers return a new state, the state
   * is updated with the new state.
   *
   * If any of the reducers throw an error, the error is
   * caught and the state is not updated.
   */
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
    builder
      .addCase(fetchDriverMetrics.fulfilled, (state, action) => {
        state.driverMetrics = action.payload;
        state.loading = false;
      })
      .addCase(fetchVehicleInfo.fulfilled, (state, action) => {
        state.vehicleInfo = action.payload;
        state.loading = false;
      })
      .addCase(fetchDeliveries.fulfilled, (state, action) => {
        state.deliveries = action.payload;
        state.loading = false;
      });
  },
});

export const { clearDashboard, setUserRole } = dashboardSlice.actions;
export default dashboardSlice.reducer;