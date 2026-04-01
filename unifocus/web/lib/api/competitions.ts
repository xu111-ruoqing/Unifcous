import apiClient from "./client";
import { arrayOrEmpty } from "@/lib/utils";

export interface Competition {
  id: number;
  name: string;
  level: string;
  official_url: string;
  typical_time_window: string;
  timeline_hint: string;
  note: string;
}

export interface CompetitionFilter {
  q?: string;
  level?: string;
  limit?: number;
  offset?: number;
}

export const competitionsAPI = {
  list: async (filter?: CompetitionFilter): Promise<Competition[]> => {
    const response = await apiClient.get("/api/v1/competitions", { params: filter });
    return arrayOrEmpty(response.data);
  },
};
