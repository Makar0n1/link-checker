// Auth types
export interface User {
  id: number;
  email: string;
  name: string;
  role: "user" | "admin";
  created_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface RefreshRequest {
  refresh_token: string;
}

// Project types
export interface Project {
  id: number;
  name: string;
  user_id: number;
  google_sheet_id: string | null;
  created_at: string;
}

export interface CreateProjectRequest {
  name: string;
  google_sheet_id?: string;
}

export interface UpdateProjectRequest {
  name?: string;
  google_sheet_id?: string;
}

// Backlink types
export type LinkStatus = "pending" | "active" | "broken" | "removed" | "nofollow";
export type LinkType = "dofollow" | "nofollow" | "sponsored" | "ugc";

export interface Backlink {
  id: number;
  project_id: number;
  source_url: string;
  target_url: string;
  anchor_text: string;
  status: LinkStatus;
  link_type: LinkType;
  http_status: number | null;
  last_checked_at: string | null;
  created_at: string;
}

export interface CreateBacklinkRequest {
  project_id: number;
  source_url: string;
  target_url: string;
  anchor_text?: string;
  link_type?: LinkType;
}

export interface UpdateBacklinkRequest {
  source_url?: string;
  target_url?: string;
  anchor_text?: string;
  status?: LinkStatus;
  link_type?: LinkType;
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
}

export interface BacklinksQueryParams {
  project_id: number;
  status?: LinkStatus;
  link_type?: LinkType;
  source_url?: string;
  target_url?: string;
  page?: number;
  per_page?: number;
}

export interface BulkCreateBacklinksRequest {
  backlinks: CreateBacklinkRequest[];
}

export interface BulkDeleteBacklinksRequest {
  ids: number[];
}

export interface BulkOperationResponse {
  success: number;
  failed: number;
  errors: string[];
}

// Error types
export interface ApiError {
  error: string;
  code?: string;
  details?: string;
}
