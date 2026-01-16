const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL || 'http://localhost:8080';
const BACKLINK_SERVICE_URL = process.env.BACKLINK_SERVICE_URL || 'http://localhost:8081';

interface RequestOptions {
  method?: string;
  body?: unknown;
  headers?: Record<string, string>;
  token?: string;
}

async function request<T>(baseUrl: string, path: string, options: RequestOptions = {}): Promise<{ status: number; data: T }> {
  const { method = 'GET', body, headers = {}, token } = options;

  const requestHeaders: Record<string, string> = {
    'Content-Type': 'application/json',
    ...headers,
  };

  if (token) {
    requestHeaders['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${baseUrl}${path}`, {
    method,
    headers: requestHeaders,
    body: body ? JSON.stringify(body) : undefined,
  });

  let data: T;
  const contentType = response.headers.get('content-type');
  if (contentType?.includes('application/json')) {
    data = await response.json() as T;
  } else {
    data = {} as T;
  }

  return { status: response.status, data };
}

// Auth Service API
export const authApi = {
  register: (data: { email: string; password: string; name: string }) =>
    request<{ id: number; email: string; name: string; role: string; created_at: string }>(
      AUTH_SERVICE_URL,
      '/api/v1/auth/register',
      { method: 'POST', body: data }
    ),

  login: (data: { email: string; password: string }) =>
    request<{ access_token: string; refresh_token: string; expires_in: number; token_type: string }>(
      AUTH_SERVICE_URL,
      '/api/v1/auth/login',
      { method: 'POST', body: data }
    ),

  me: (token: string) =>
    request<{ id: number; email: string; name: string; role: string; created_at: string }>(
      AUTH_SERVICE_URL,
      '/api/v1/auth/me',
      { token }
    ),

  refresh: (refreshToken: string) =>
    request<{ access_token: string; refresh_token: string; expires_in: number; token_type: string }>(
      AUTH_SERVICE_URL,
      '/api/v1/auth/refresh',
      { method: 'POST', body: { refresh_token: refreshToken } }
    ),

  logout: (refreshToken: string) =>
    request<{ message: string }>(
      AUTH_SERVICE_URL,
      '/api/v1/auth/logout',
      { method: 'POST', body: { refresh_token: refreshToken } }
    ),

  health: () =>
    request<{ status: string }>(AUTH_SERVICE_URL, '/health'),
};

// Backlink Service API
export const backlinkApi = {
  // Projects
  createProject: (token: string, data: { name: string; url?: string }) =>
    request<{ id: number; user_id: number; name: string; url: string; created_at: string; updated_at: string }>(
      BACKLINK_SERVICE_URL,
      '/api/v1/projects',
      { method: 'POST', body: data, token }
    ),

  listProjects: (token: string) =>
    request<Array<{ id: number; user_id: number; name: string; url: string; created_at: string; updated_at: string }>>(
      BACKLINK_SERVICE_URL,
      '/api/v1/projects',
      { token }
    ),

  getProject: (token: string, id: number) =>
    request<{ id: number; user_id: number; name: string; url: string; created_at: string; updated_at: string }>(
      BACKLINK_SERVICE_URL,
      `/api/v1/projects/${id}`,
      { token }
    ),

  updateProject: (token: string, id: number, data: { name?: string; url?: string }) =>
    request<{ id: number; user_id: number; name: string; url: string; created_at: string; updated_at: string }>(
      BACKLINK_SERVICE_URL,
      `/api/v1/projects/${id}`,
      { method: 'PUT', body: data, token }
    ),

  deleteProject: (token: string, id: number) =>
    request<void>(BACKLINK_SERVICE_URL, `/api/v1/projects/${id}`, { method: 'DELETE', token }),

  // Backlinks
  createBacklink: (
    token: string,
    data: { project_id: number; source_url: string; target_url: string; anchor_text?: string; status?: string }
  ) =>
    request<{
      id: number;
      project_id: number;
      source_url: string;
      target_url: string;
      anchor_text: string;
      status: string;
      link_type: string;
      created_at: string;
      updated_at: string;
    }>(BACKLINK_SERVICE_URL, '/api/v1/backlinks', { method: 'POST', body: data, token }),

  listBacklinks: (token: string, params: { project_id: number; page?: number; per_page?: number; status?: string }) => {
    const searchParams = new URLSearchParams();
    searchParams.set('project_id', params.project_id.toString());
    if (params.page) searchParams.set('page', params.page.toString());
    if (params.per_page) searchParams.set('per_page', params.per_page.toString());
    if (params.status) searchParams.set('status', params.status);

    return request<{
      data: Array<{
        id: number;
        project_id: number;
        source_url: string;
        target_url: string;
        anchor_text: string;
        status: string;
        link_type: string;
        created_at: string;
        updated_at: string;
      }>;
      page: number;
      per_page: number;
      total: number;
      total_pages: number;
    }>(BACKLINK_SERVICE_URL, `/api/v1/backlinks?${searchParams.toString()}`, { token });
  },

  getBacklink: (token: string, id: number) =>
    request<{
      id: number;
      project_id: number;
      source_url: string;
      target_url: string;
      anchor_text: string;
      status: string;
      link_type: string;
      created_at: string;
      updated_at: string;
    }>(BACKLINK_SERVICE_URL, `/api/v1/backlinks/${id}`, { token }),

  updateBacklink: (
    token: string,
    id: number,
    data: { source_url?: string; target_url?: string; anchor_text?: string; status?: string }
  ) =>
    request<{
      id: number;
      project_id: number;
      source_url: string;
      target_url: string;
      anchor_text: string;
      status: string;
      link_type: string;
      created_at: string;
      updated_at: string;
    }>(BACKLINK_SERVICE_URL, `/api/v1/backlinks/${id}`, { method: 'PUT', body: data, token }),

  deleteBacklink: (token: string, id: number) =>
    request<void>(BACKLINK_SERVICE_URL, `/api/v1/backlinks/${id}`, { method: 'DELETE', token }),

  bulkCreateBacklinks: (
    token: string,
    data: {
      backlinks: Array<{ project_id: number; source_url: string; target_url: string; anchor_text?: string }>;
    }
  ) =>
    request<{ created: number; failed: number; errors?: string[] }>(
      BACKLINK_SERVICE_URL,
      '/api/v1/backlinks/bulk',
      { method: 'POST', body: data, token }
    ),

  bulkDeleteBacklinks: (token: string, data: { ids: number[] }) =>
    request<{ deleted: number; failed: number }>(
      BACKLINK_SERVICE_URL,
      '/api/v1/backlinks/bulk',
      { method: 'DELETE', body: data, token }
    ),

  health: () =>
    request<{ status: string }>(BACKLINK_SERVICE_URL, '/health'),
};

// Test helpers
export function generateTestEmail(): string {
  return `test-${Date.now()}-${Math.random().toString(36).substring(7)}@example.com`;
}

export async function createTestUser(): Promise<{ email: string; password: string; token: string; refreshToken: string }> {
  const email = generateTestEmail();
  const password = 'TestPassword123!';

  await authApi.register({ email, password, name: 'Test User' });
  const loginResult = await authApi.login({ email, password });

  return {
    email,
    password,
    token: loginResult.data.access_token,
    refreshToken: loginResult.data.refresh_token,
  };
}

export async function createTestProject(token: string): Promise<number> {
  const result = await backlinkApi.createProject(token, {
    name: `Test Project ${Date.now()}`,
    url: 'https://example.com',
  });
  return result.data.id;
}
