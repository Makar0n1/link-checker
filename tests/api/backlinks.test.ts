import { describe, it, expect, beforeAll, beforeEach } from 'vitest';
import { backlinkApi, createTestUser, createTestProject } from './helpers/api-client';

describe('Backlink Service API', () => {
  let accessToken: string;
  let projectId: number;

  beforeAll(async () => {
    const user = await createTestUser();
    accessToken = user.token;
  });

  describe('Projects API', () => {
    describe('POST /api/v1/projects', () => {
      it('should create a project and return 201', async () => {
        const result = await backlinkApi.createProject(accessToken, {
          name: 'Test Project',
          url: 'https://example.com',
        });

        expect(result.status).toBe(201);
        expect(result.data.name).toBe('Test Project');
        expect(result.data.url).toBe('https://example.com');
        expect(result.data.id).toBeDefined();
        expect(result.data.user_id).toBeDefined();
        expect(result.data.created_at).toBeDefined();
      });

      it('should return 400 for missing name', async () => {
        const result = await backlinkApi.createProject(accessToken, {
          name: '',
          url: 'https://example.com',
        });

        expect(result.status).toBe(400);
      });

      it('should return 401 without token', async () => {
        const result = await backlinkApi.createProject('', {
          name: 'Test Project',
          url: 'https://example.com',
        });

        expect(result.status).toBe(401);
      });

      it('should return 401 with invalid token', async () => {
        const result = await backlinkApi.createProject('invalid-token', {
          name: 'Test Project',
          url: 'https://example.com',
        });

        expect(result.status).toBe(401);
      });
    });

    describe('GET /api/v1/projects', () => {
      beforeAll(async () => {
        // Create a project for listing tests
        await backlinkApi.createProject(accessToken, {
          name: 'List Test Project',
          url: 'https://list-test.com',
        });
      });

      it('should list projects and return 200', async () => {
        const result = await backlinkApi.listProjects(accessToken);

        expect(result.status).toBe(200);
        expect(Array.isArray(result.data)).toBe(true);
        expect(result.data.length).toBeGreaterThan(0);
      });

      it('should return 401 without token', async () => {
        const result = await backlinkApi.listProjects('');

        expect(result.status).toBe(401);
      });
    });

    describe('GET /api/v1/projects/:id', () => {
      let testProjectId: number;

      beforeAll(async () => {
        const result = await backlinkApi.createProject(accessToken, {
          name: 'Get Test Project',
          url: 'https://get-test.com',
        });
        testProjectId = result.data.id;
      });

      it('should get project by id and return 200', async () => {
        const result = await backlinkApi.getProject(accessToken, testProjectId);

        expect(result.status).toBe(200);
        expect(result.data.id).toBe(testProjectId);
        expect(result.data.name).toBe('Get Test Project');
      });

      it('should return 404 for non-existent project', async () => {
        const result = await backlinkApi.getProject(accessToken, 999999);

        expect(result.status).toBe(404);
      });

      it('should return 401 without token', async () => {
        const result = await backlinkApi.getProject('', testProjectId);

        expect(result.status).toBe(401);
      });
    });

    describe('PUT /api/v1/projects/:id', () => {
      let testProjectId: number;

      beforeAll(async () => {
        const result = await backlinkApi.createProject(accessToken, {
          name: 'Update Test Project',
          url: 'https://update-test.com',
        });
        testProjectId = result.data.id;
      });

      it('should update project and return 200', async () => {
        const result = await backlinkApi.updateProject(accessToken, testProjectId, {
          name: 'Updated Project Name',
        });

        expect(result.status).toBe(200);
        expect(result.data.name).toBe('Updated Project Name');
      });

      it('should return 404 for non-existent project', async () => {
        const result = await backlinkApi.updateProject(accessToken, 999999, {
          name: 'New Name',
        });

        expect(result.status).toBe(404);
      });

      it('should return 401 without token', async () => {
        const result = await backlinkApi.updateProject('', testProjectId, {
          name: 'New Name',
        });

        expect(result.status).toBe(401);
      });
    });

    describe('DELETE /api/v1/projects/:id', () => {
      it('should delete project and return 204', async () => {
        // Create a project to delete
        const createResult = await backlinkApi.createProject(accessToken, {
          name: 'Delete Test Project',
          url: 'https://delete-test.com',
        });

        const result = await backlinkApi.deleteProject(accessToken, createResult.data.id);

        expect(result.status).toBe(204);
      });

      it('should return 404 for non-existent project', async () => {
        const result = await backlinkApi.deleteProject(accessToken, 999999);

        expect(result.status).toBe(404);
      });

      it('should return 401 without token', async () => {
        const result = await backlinkApi.deleteProject('', 1);

        expect(result.status).toBe(401);
      });
    });
  });

  describe('Backlinks API', () => {
    beforeEach(async () => {
      // Create a fresh project for each test
      projectId = await createTestProject(accessToken);
    });

    describe('POST /api/v1/backlinks', () => {
      it('should create a backlink and return 201', async () => {
        const result = await backlinkApi.createBacklink(accessToken, {
          project_id: projectId,
          source_url: 'https://source.com/page',
          target_url: 'https://target.com/page',
          anchor_text: 'Test Anchor',
        });

        expect(result.status).toBe(201);
        expect(result.data.project_id).toBe(projectId);
        expect(result.data.source_url).toBe('https://source.com/page');
        expect(result.data.target_url).toBe('https://target.com/page');
        expect(result.data.anchor_text).toBe('Test Anchor');
        expect(result.data.id).toBeDefined();
        expect(result.data.status).toBeDefined();
      });

      it('should return 400 for missing project_id', async () => {
        const result = await backlinkApi.createBacklink(accessToken, {
          project_id: 0,
          source_url: 'https://source.com',
          target_url: 'https://target.com',
        });

        expect(result.status).toBe(400);
      });

      it('should return 400 for missing source_url', async () => {
        const result = await backlinkApi.createBacklink(accessToken, {
          project_id: projectId,
          source_url: '',
          target_url: 'https://target.com',
        });

        expect(result.status).toBe(400);
      });

      it('should return 400 for missing target_url', async () => {
        const result = await backlinkApi.createBacklink(accessToken, {
          project_id: projectId,
          source_url: 'https://source.com',
          target_url: '',
        });

        expect(result.status).toBe(400);
      });

      it('should return 401 without token', async () => {
        const result = await backlinkApi.createBacklink('', {
          project_id: projectId,
          source_url: 'https://source.com',
          target_url: 'https://target.com',
        });

        expect(result.status).toBe(401);
      });
    });

    describe('GET /api/v1/backlinks', () => {
      beforeEach(async () => {
        // Create some backlinks for testing
        await backlinkApi.createBacklink(accessToken, {
          project_id: projectId,
          source_url: 'https://source1.com',
          target_url: 'https://target1.com',
          anchor_text: 'Anchor 1',
        });
        await backlinkApi.createBacklink(accessToken, {
          project_id: projectId,
          source_url: 'https://source2.com',
          target_url: 'https://target2.com',
          anchor_text: 'Anchor 2',
        });
      });

      it('should list backlinks and return 200', async () => {
        const result = await backlinkApi.listBacklinks(accessToken, {
          project_id: projectId,
        });

        expect(result.status).toBe(200);
        expect(result.data.data).toBeDefined();
        expect(Array.isArray(result.data.data)).toBe(true);
        expect(result.data.data.length).toBeGreaterThanOrEqual(2);
        expect(result.data.total).toBeGreaterThanOrEqual(2);
        expect(result.data.page).toBe(1);
      });

      it('should support pagination', async () => {
        const result = await backlinkApi.listBacklinks(accessToken, {
          project_id: projectId,
          page: 1,
          per_page: 1,
        });

        expect(result.status).toBe(200);
        expect(result.data.data.length).toBe(1);
        expect(result.data.per_page).toBe(1);
      });

      it('should return 400 without project_id', async () => {
        const result = await backlinkApi.listBacklinks(accessToken, {
          project_id: 0,
        });

        expect(result.status).toBe(400);
      });

      it('should return 401 without token', async () => {
        const result = await backlinkApi.listBacklinks('', {
          project_id: projectId,
        });

        expect(result.status).toBe(401);
      });
    });

    describe('GET /api/v1/backlinks/:id', () => {
      let backlinkId: number;

      beforeEach(async () => {
        const result = await backlinkApi.createBacklink(accessToken, {
          project_id: projectId,
          source_url: 'https://get-source.com',
          target_url: 'https://get-target.com',
        });
        backlinkId = result.data.id;
      });

      it('should get backlink by id and return 200', async () => {
        const result = await backlinkApi.getBacklink(accessToken, backlinkId);

        expect(result.status).toBe(200);
        expect(result.data.id).toBe(backlinkId);
        expect(result.data.source_url).toBe('https://get-source.com');
      });

      it('should return 404 for non-existent backlink', async () => {
        const result = await backlinkApi.getBacklink(accessToken, 999999);

        expect(result.status).toBe(404);
      });

      it('should return 401 without token', async () => {
        const result = await backlinkApi.getBacklink('', backlinkId);

        expect(result.status).toBe(401);
      });
    });

    describe('PUT /api/v1/backlinks/:id', () => {
      let backlinkId: number;

      beforeEach(async () => {
        const result = await backlinkApi.createBacklink(accessToken, {
          project_id: projectId,
          source_url: 'https://update-source.com',
          target_url: 'https://update-target.com',
          anchor_text: 'Original Anchor',
        });
        backlinkId = result.data.id;
      });

      it('should update backlink and return 200', async () => {
        const result = await backlinkApi.updateBacklink(accessToken, backlinkId, {
          anchor_text: 'Updated Anchor',
          status: 'active',
        });

        expect(result.status).toBe(200);
        expect(result.data.anchor_text).toBe('Updated Anchor');
        expect(result.data.status).toBe('active');
      });

      it('should return 404 for non-existent backlink', async () => {
        const result = await backlinkApi.updateBacklink(accessToken, 999999, {
          anchor_text: 'New Anchor',
        });

        expect(result.status).toBe(404);
      });

      it('should return 401 without token', async () => {
        const result = await backlinkApi.updateBacklink('', backlinkId, {
          anchor_text: 'New Anchor',
        });

        expect(result.status).toBe(401);
      });
    });

    describe('DELETE /api/v1/backlinks/:id', () => {
      it('should delete backlink and return 204', async () => {
        const createResult = await backlinkApi.createBacklink(accessToken, {
          project_id: projectId,
          source_url: 'https://delete-source.com',
          target_url: 'https://delete-target.com',
        });

        const result = await backlinkApi.deleteBacklink(accessToken, createResult.data.id);

        expect(result.status).toBe(204);
      });

      it('should return 404 for non-existent backlink', async () => {
        const result = await backlinkApi.deleteBacklink(accessToken, 999999);

        expect(result.status).toBe(404);
      });

      it('should return 401 without token', async () => {
        const result = await backlinkApi.deleteBacklink('', 1);

        expect(result.status).toBe(401);
      });
    });

    describe('POST /api/v1/backlinks/bulk', () => {
      it('should bulk create backlinks and return 201', async () => {
        const result = await backlinkApi.bulkCreateBacklinks(accessToken, {
          backlinks: [
            {
              project_id: projectId,
              source_url: 'https://bulk-source1.com',
              target_url: 'https://bulk-target1.com',
            },
            {
              project_id: projectId,
              source_url: 'https://bulk-source2.com',
              target_url: 'https://bulk-target2.com',
            },
            {
              project_id: projectId,
              source_url: 'https://bulk-source3.com',
              target_url: 'https://bulk-target3.com',
            },
          ],
        });

        expect(result.status).toBe(201);
        expect(result.data.created).toBe(3);
        expect(result.data.failed).toBe(0);
      });

      it('should return 400 for empty backlinks array', async () => {
        const result = await backlinkApi.bulkCreateBacklinks(accessToken, {
          backlinks: [],
        });

        expect(result.status).toBe(400);
      });

      it('should return 401 without token', async () => {
        const result = await backlinkApi.bulkCreateBacklinks('', {
          backlinks: [
            {
              project_id: projectId,
              source_url: 'https://bulk-source.com',
              target_url: 'https://bulk-target.com',
            },
          ],
        });

        expect(result.status).toBe(401);
      });
    });

    describe('DELETE /api/v1/backlinks/bulk', () => {
      it('should bulk delete backlinks and return 200', async () => {
        // Create backlinks to delete
        const create1 = await backlinkApi.createBacklink(accessToken, {
          project_id: projectId,
          source_url: 'https://bulk-del1.com',
          target_url: 'https://target.com',
        });
        const create2 = await backlinkApi.createBacklink(accessToken, {
          project_id: projectId,
          source_url: 'https://bulk-del2.com',
          target_url: 'https://target.com',
        });

        const result = await backlinkApi.bulkDeleteBacklinks(accessToken, {
          ids: [create1.data.id, create2.data.id],
        });

        expect(result.status).toBe(200);
        expect(result.data.deleted).toBe(2);
        expect(result.data.failed).toBe(0);
      });

      it('should return 400 for empty ids array', async () => {
        const result = await backlinkApi.bulkDeleteBacklinks(accessToken, {
          ids: [],
        });

        expect(result.status).toBe(400);
      });

      it('should return 401 without token', async () => {
        const result = await backlinkApi.bulkDeleteBacklinks('', {
          ids: [1, 2],
        });

        expect(result.status).toBe(401);
      });
    });
  });

  describe('Unauthorized requests', () => {
    it('should return 401 for projects without token', async () => {
      const result = await backlinkApi.listProjects('');
      expect(result.status).toBe(401);
    });

    it('should return 401 for backlinks without token', async () => {
      const result = await backlinkApi.listBacklinks('', { project_id: 1 });
      expect(result.status).toBe(401);
    });

    it('should return 401 for create project without token', async () => {
      const result = await backlinkApi.createProject('', { name: 'Test' });
      expect(result.status).toBe(401);
    });

    it('should return 401 for create backlink without token', async () => {
      const result = await backlinkApi.createBacklink('', {
        project_id: 1,
        source_url: 'https://source.com',
        target_url: 'https://target.com',
      });
      expect(result.status).toBe(401);
    });
  });

  describe('GET /health', () => {
    it('should return health status', async () => {
      const result = await backlinkApi.health();

      expect(result.status).toBe(200);
      expect(result.data.status).toBe('ok');
    });
  });
});
