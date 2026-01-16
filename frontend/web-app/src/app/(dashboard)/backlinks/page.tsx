"use client";

import { useState, useMemo, useCallback } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Link2,
  TrendingUp,
  TrendingDown,
  AlertCircle,
  Loader2,
  Plus,
  RefreshCw,
} from "lucide-react";
import { DataTable, createBacklinksColumns } from "@/components/tables";
import {
  useBacklinks,
  useUpdateBacklink,
  useBulkDeleteBacklinks,
} from "@/hooks/use-backlinks";
import { useProjects } from "@/hooks/use-projects";
import type { Backlink, LinkStatus } from "@/types/api";

export default function BacklinksPage() {
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null);
  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState<LinkStatus | "all">("all");

  const { data: projects, isLoading: projectsLoading } = useProjects();
  const {
    data: backlinksData,
    isLoading: backlinksLoading,
    refetch,
  } = useBacklinks({
    project_id: selectedProjectId || projects?.[0]?.id || 0,
    page,
    per_page: 20,
    status: statusFilter === "all" ? undefined : statusFilter,
  });

  const updateBacklinkMutation = useUpdateBacklink();
  const bulkDeleteMutation = useBulkDeleteBacklinks();

  // Auto-select first project if none selected
  const currentProjectId = selectedProjectId || projects?.[0]?.id || null;

  const handleUpdateBacklink = useCallback(
    (id: number, field: keyof Backlink, value: string) => {
      updateBacklinkMutation.mutate({
        id,
        data: { [field]: value },
      });
    },
    [updateBacklinkMutation]
  );

  const handleBulkDelete = useCallback(
    (rows: Backlink[]) => {
      const ids = rows.map((r) => r.id);
      bulkDeleteMutation.mutate({ ids });
    },
    [bulkDeleteMutation]
  );

  const columns = useMemo(
    () => createBacklinksColumns({ onUpdateBacklink: handleUpdateBacklink }),
    [handleUpdateBacklink]
  );

  const backlinks = backlinksData?.data || [];
  const total = backlinksData?.total || 0;
  const totalPages = backlinksData?.total_pages || 1;

  // Stats calculation
  const stats = useMemo(() => {
    const activeCount = backlinks.filter((b) => b.status === "active").length;
    const brokenCount = backlinks.filter((b) => b.status === "broken").length;
    const pendingCount = backlinks.filter((b) => b.status === "pending").length;

    return [
      {
        title: "Total Backlinks",
        value: total.toLocaleString(),
        icon: Link2,
      },
      {
        title: "Active Links",
        value: activeCount.toString(),
        trend: "up" as const,
        icon: TrendingUp,
      },
      {
        title: "Broken Links",
        value: brokenCount.toString(),
        trend: "down" as const,
        icon: TrendingDown,
      },
      {
        title: "Pending",
        value: pendingCount.toString(),
        icon: AlertCircle,
      },
    ];
  }, [backlinks, total]);

  const isLoading = projectsLoading || backlinksLoading;

  if (!currentProjectId && !projectsLoading) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <h2 className="text-xl font-semibold mb-2">No Projects Found</h2>
        <p className="text-muted-foreground mb-4">
          Create a project to start tracking backlinks
        </p>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Create Project
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Backlinks</h1>
          <p className="text-muted-foreground">
            Monitor and manage your backlink profile
          </p>
        </div>
        <div className="flex items-center gap-4">
          <Select
            value={currentProjectId?.toString() || ""}
            onValueChange={(value) => {
              setSelectedProjectId(parseInt(value));
              setPage(1);
            }}
          >
            <SelectTrigger className="w-[200px]">
              <SelectValue placeholder="Select project" />
            </SelectTrigger>
            <SelectContent>
              {projects?.map((project) => (
                <SelectItem key={project.id} value={project.id.toString()}>
                  {project.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Button
            variant="outline"
            size="icon"
            onClick={() => refetch()}
            disabled={isLoading}
          >
            <RefreshCw
              className={`h-4 w-4 ${isLoading ? "animate-spin" : ""}`}
            />
          </Button>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <Card key={stat.title}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {stat.title}
              </CardTitle>
              <stat.icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {isLoading ? (
                  <Loader2 className="h-6 w-6 animate-spin" />
                ) : (
                  stat.value
                )}
              </div>
              {stat.trend && (
                <Badge
                  variant={stat.trend === "up" ? "default" : "destructive"}
                  className="mt-1"
                >
                  {stat.trend === "up" ? "Active" : "Issues"}
                </Badge>
              )}
            </CardContent>
          </Card>
        ))}
      </div>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>Backlinks Table</CardTitle>
          <div className="flex items-center gap-2">
            <Select
              value={statusFilter}
              onValueChange={(value) => {
                setStatusFilter(value as LinkStatus | "all");
                setPage(1);
              }}
            >
              <SelectTrigger className="w-[150px]">
                <SelectValue placeholder="Filter status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Statuses</SelectItem>
                <SelectItem value="active">Active</SelectItem>
                <SelectItem value="pending">Pending</SelectItem>
                <SelectItem value="broken">Broken</SelectItem>
                <SelectItem value="removed">Removed</SelectItem>
                <SelectItem value="nofollow">Nofollow</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="flex items-center justify-center py-12">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : (
            <>
              <DataTable
                columns={columns}
                data={backlinks}
                searchKey="source_url"
                searchPlaceholder="Filter by source URL..."
                onBulkDelete={handleBulkDelete}
              />
              {totalPages > 1 && (
                <div className="flex items-center justify-center gap-2 mt-4">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage((p) => Math.max(1, p - 1))}
                    disabled={page === 1}
                  >
                    Previous
                  </Button>
                  <span className="text-sm text-muted-foreground">
                    Page {page} of {totalPages}
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                    disabled={page === totalPages}
                  >
                    Next
                  </Button>
                </div>
              )}
            </>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
