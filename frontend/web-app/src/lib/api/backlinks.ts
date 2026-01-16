import { backlinkApi } from "./client";
import type {
  Backlink,
  CreateBacklinkRequest,
  UpdateBacklinkRequest,
  BacklinksQueryParams,
  PaginatedResponse,
  BulkCreateBacklinksRequest,
  BulkDeleteBacklinksRequest,
  BulkOperationResponse,
} from "@/types/api";

function buildQueryString(params: BacklinksQueryParams): string {
  const searchParams = new URLSearchParams();

  searchParams.set("project_id", params.project_id.toString());
  if (params.status) searchParams.set("status", params.status);
  if (params.link_type) searchParams.set("link_type", params.link_type);
  if (params.source_url) searchParams.set("source_url", params.source_url);
  if (params.target_url) searchParams.set("target_url", params.target_url);
  if (params.page) searchParams.set("page", params.page.toString());
  if (params.per_page) searchParams.set("per_page", params.per_page.toString());

  return searchParams.toString();
}

export async function getBacklinks(
  params: BacklinksQueryParams
): Promise<PaginatedResponse<Backlink>> {
  const query = buildQueryString(params);
  return backlinkApi.get<PaginatedResponse<Backlink>>(
    `/api/v1/backlinks?${query}`
  );
}

export async function getBacklink(id: number): Promise<Backlink> {
  return backlinkApi.get<Backlink>(`/api/v1/backlinks/${id}`);
}

export async function createBacklink(
  data: CreateBacklinkRequest
): Promise<Backlink> {
  return backlinkApi.post<Backlink>("/api/v1/backlinks", data);
}

export async function updateBacklink(
  id: number,
  data: UpdateBacklinkRequest
): Promise<Backlink> {
  return backlinkApi.put<Backlink>(`/api/v1/backlinks/${id}`, data);
}

export async function deleteBacklink(id: number): Promise<void> {
  return backlinkApi.delete(`/api/v1/backlinks/${id}`);
}

export async function bulkCreateBacklinks(
  data: BulkCreateBacklinksRequest
): Promise<BulkOperationResponse> {
  return backlinkApi.post<BulkOperationResponse>("/api/v1/backlinks/bulk", data);
}

export async function bulkDeleteBacklinks(
  data: BulkDeleteBacklinksRequest
): Promise<BulkOperationResponse> {
  return backlinkApi.delete<BulkOperationResponse>("/api/v1/backlinks/bulk", data);
}
