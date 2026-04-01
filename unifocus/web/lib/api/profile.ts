/**
 * 用户画像相关API
 */
import apiClient from './client';

export interface UserProfile {
  id: number;
  user_id: number;
  resume_text?: string;
  skills?: string[];
  certificates?: Array<{ name: string; score?: number; date?: string }>;
  interests?: string[];
  updated_at: string;
}

export interface UpdateProfileRequest {
  resume_text?: string;
  skills?: string[];
  certificates?: Array<{ name: string; score?: number; date?: string }>;
  interests?: string[];
}

export const profileAPI = {
  // 获取当前用户画像
  getProfile: async (): Promise<UserProfile> => {
    const response = await apiClient.get('/api/v1/users/me/profile');
    return response.data;
  },

  // 更新画像
  updateProfile: async (data: UpdateProfileRequest): Promise<UserProfile> => {
    const response = await apiClient.put('/api/v1/users/me/profile', data);
    return response.data;
  },

  // 上传简历
  uploadResume: async (file: File): Promise<UserProfile> => {
    const formData = new FormData();
    formData.append('file', file);
    const response = await apiClient.post('/api/v1/users/me/profile/resume', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },
};


