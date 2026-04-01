/**
 * 监控指标相关API
 */
import apiClient from './client';

export interface SystemMetrics {
  status: string;
  timestamp: {
    current_time: string;
  };
  system: {
    uptime: string;
  };
}

export const metricsAPI = {
  // 获取系统指标
  getMetrics: async (): Promise<SystemMetrics> => {
    const response = await apiClient.get('/api/v1/metrics');
    return response.data;
  },
};


