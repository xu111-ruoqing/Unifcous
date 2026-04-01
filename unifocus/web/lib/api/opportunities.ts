/**
 * 机会相关API
 */
import { arrayOrEmpty } from "@/lib/utils";
import apiClient from "./client";

export interface Opportunity {
  id: number;
  title: string;
  type: string;
  description: string;
  source_url: string;
  competition_level?: string;
  organizer?: string;
  start_date?: string;
  deadline?: string;
  event_date?: string;
  location?: string;
  target_majors?: string[];
  tags?: string[];
  view_count: number;
  save_count: number;
  created_at: string;
}

export interface OpportunityFilter {
  type?: string;
  competition_level?: string;
  major?: string;
  deadline_after?: string;
  deadline_before?: string;
  is_active?: boolean;
  tags?: string[];
  limit?: number;
  offset?: number;
}

export interface OpportunityListResponse {
  data: Opportunity[];
  total: number;
  limit: number;
  offset: number;
}

export interface CreateOpportunityRequest {
  title: string;
  type: string;
  description: string;
  source_url: string;
  organizer?: string;
  start_date?: string;
  deadline?: string;
  event_date?: string;
  location?: string;
  target_majors?: string[];
  tags?: string[];
}

export const opportunitiesAPI = {
  // 获取机会列表
  list: async (filter?: OpportunityFilter): Promise<OpportunityListResponse> => {
    const response = await apiClient.get("/api/v1/opportunities", { params: filter });
    return {
      ...response.data,
      data: arrayOrEmpty(response.data?.data),
    };
  },

  // 获取机会详情
  getById: async (id: number): Promise<Opportunity> => {
    const response = await apiClient.get(`/api/v1/opportunities/${id}`);
    return response.data;
  },

  // 创建机会
  create: async (data: CreateOpportunityRequest): Promise<Opportunity> => {
    const response = await apiClient.post("/api/v1/opportunities", data);
    return response.data;
  },

  // 更新机会
  update: async (id: number, data: CreateOpportunityRequest): Promise<Opportunity> => {
    const response = await apiClient.put(`/api/v1/opportunities/${id}`, data);
    return response.data;
  },

  // 删除机会
  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/api/v1/opportunities/${id}`);
  },
};
