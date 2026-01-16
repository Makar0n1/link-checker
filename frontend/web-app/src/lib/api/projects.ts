import { backlinkApi } from "./client";
import type {
  Project,
  CreateProjectRequest,
  UpdateProjectRequest,
} from "@/types/api";

interface ProjectsListResponse {
  data: Project[];
  total: number;
}

export async function getProjects(): Promise<Project[]> {
  const response = await backlinkApi.get<ProjectsListResponse | Project[]>(
    "/api/v1/projects"
  );
  // Handle both response formats
  if (Array.isArray(response)) {
    return response;
  }
  return response.data;
}

export async function getProject(id: number): Promise<Project> {
  return backlinkApi.get<Project>(`/api/v1/projects/${id}`);
}

export async function createProject(data: CreateProjectRequest): Promise<Project> {
  return backlinkApi.post<Project>("/api/v1/projects", data);
}

export async function updateProject(
  id: number,
  data: UpdateProjectRequest
): Promise<Project> {
  return backlinkApi.put<Project>(`/api/v1/projects/${id}`, data);
}

export async function deleteProject(id: number): Promise<void> {
  return backlinkApi.delete(`/api/v1/projects/${id}`);
}
