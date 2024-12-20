export interface User {
  id: number;
  email: string;
  username: string;
  name: string;
  role: string;
}

export interface LoginCredentials {
  identifier: string; // Can be either email or username
  password: string;
}

export interface AuthResponse {
  user: User;
  token: string;
  refreshToken: string;
}

export interface RefreshTokenResponse {
  token: string;
  refreshToken: string;
}
