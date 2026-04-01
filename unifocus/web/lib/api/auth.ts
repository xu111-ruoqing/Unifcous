/**
 * 认证相关API
 */
import apiClient from './client';

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  school: string;
  major: string;
  grade: number;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: {
    id: number;
    username: string;
    email: string;
    school: string;
    major: string;
    grade: number;
  };
}

export const authAPI = {
  // 注册
  register: async (data: RegisterRequest): Promise<LoginResponse> => {
    const response = await apiClient.post('/api/v1/auth/register', data);
    return response.data;
  },

  // 登录
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await apiClient.post('/api/v1/auth/login', data);
    return response.data;
  },

  // 刷新Token
  refreshToken: async (): Promise<{ token: string }> => {
    const response = await apiClient.post('/api/v1/auth/refresh');
    return response.data;
  },
};


