import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { User } from '../../types';
import { jwtDecode } from 'jwt-decode';

interface AuthState {
  user: User | null;
  token: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
}

/**
 * Retrieves the user information from localStorage and decodes the access token
 * to include the user's role. If both the stored user and token are available,
 * it returns an object containing the user details along with the decoded role.
 * If either is missing, it returns null.
 *
 * @returns {object | null} The user object with role information or null if not available.
 */

const getUserFromStorage = () => {
  const storedUser = localStorage.getItem('user');
  const token = localStorage.getItem('accessToken');

  if (storedUser && token) {
    const user = JSON.parse(storedUser);
    const decoded: any = jwtDecode(token);
    return {
      ...user,
      role: decoded.role_name
    };
  }
  return null;
};

const initialState: AuthState = {
  user: getUserFromStorage(),
  token: localStorage.getItem('accessToken'),
  refreshToken: localStorage.getItem('refreshToken'),
  isAuthenticated: !!localStorage.getItem('accessToken'),
  loading: false,
  error: null,
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    /**
     * Sets the loading state to true and clears any error message.
     *
     * @remarks
     * This action is intended to be used when the login process is initiated.
     * It sets the loading state to true and clears any error message.
     */
    loginStart: (state) => {
      state.loading = true;
      state.error = null;
    },
/**
 * Handles successful login actions by updating the authentication state.
 *
 * @remarks
 * This reducer is triggered when a login attempt is successful. It sets the
 * loading state to false, marks the user as authenticated, and updates the
 * user, token, and refreshToken in the state. Additionally, it clears any
 * error messages.
 *
 * The function also synchronizes the user information and tokens with
 * localStorage to ensure persistence across sessions.
 *
 * @param state - The current authentication state.
 * @param action - The action payload containing the user, token, and refreshToken details.
 */

    loginSuccess: (state, action: PayloadAction<{ user: User; token: string; refreshToken: string }>) => {
      state.loading = false;
      state.isAuthenticated = true;
      state.user = action.payload.user;
      state.token = action.payload.token;
      state.refreshToken = action.payload.refreshToken;
      state.error = null;
      
      // Ensure localStorage is updated
      localStorage.setItem('user', JSON.stringify(action.payload.user));
      localStorage.setItem('accessToken', action.payload.token);
      localStorage.setItem('refreshToken', action.payload.refreshToken);
    },
    loginFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
      state.isAuthenticated = false;
      state.user = null;
      state.token = null;
      state.refreshToken = null;
      
      // Clear localStorage on login failure
      localStorage.removeItem('user');
      localStorage.removeItem('accessToken');
      localStorage.removeItem('refreshToken');
    },
    /**
     * Logs the user out by clearing the authentication state and localStorage.
     *
     * @remarks
     * This reducer is triggered when the user initiates a logout. It clears the
     * user, token, and refreshToken in the state and localStorage to ensure
     * that the user is no longer authenticated.
     */
    logout: (state) => {
      state.user = null;
      state.token = null;
      state.refreshToken = null;
      state.isAuthenticated = false;
      state.loading = false;
      state.error = null;
      
      // Clear localStorage on logout
      localStorage.removeItem('user');
      localStorage.removeItem('accessToken');
      localStorage.removeItem('refreshToken');
    },
    /**
     * Updates the access token and refresh token in the state and localStorage when a token refresh is successful.
     *
     * @param {object} action - The action payload containing the new access token and refresh token.
     * @param {string} action.payload.token - The new access token.
     * @param {string} action.payload.refreshToken - The new refresh token.
     */
    refreshTokenSuccess: (state, action: PayloadAction<{ token: string; refreshToken: string }>) => {
      state.token = action.payload.token;
      state.refreshToken = action.payload.refreshToken;
      localStorage.setItem('accessToken', action.payload.token);
      localStorage.setItem('refreshToken', action.payload.refreshToken);
    },
  },
});

export const {
  loginStart,
  loginSuccess,
  loginFailure,
  logout,
  refreshTokenSuccess,
} = authSlice.actions;

export default authSlice.reducer;
