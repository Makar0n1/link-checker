"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  getProjects,
  getProject,
  createProject,
  updateProject,
  deleteProject,
} from "@/lib/api";
import type { CreateProjectRequest, UpdateProjectRequest, Project } from "@/types/api";
import { useToast } from "./use-toast";

export function useProjects() {
  return useQuery<Project[]>({
    queryKey: ["projects"],
    queryFn: getProjects,
  });
}

export function useProject(id: number) {
  return useQuery<Project>({
    queryKey: ["projects", id],
    queryFn: () => getProject(id),
    enabled: !!id,
  });
}

export function useCreateProject() {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: (data: CreateProjectRequest) => createProject(data),
    onSuccess: (newProject) => {
      queryClient.invalidateQueries({ queryKey: ["projects"] });
      toast({
        title: "Project created",
        description: `"${newProject.name}" has been created`,
      });
    },
    onError: (error: Error) => {
      toast({
        variant: "destructive",
        title: "Failed to create project",
        description: error.message,
      });
    },
  });
}

export function useUpdateProject() {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateProjectRequest }) =>
      updateProject(id, data),
    onSuccess: (updatedProject) => {
      queryClient.invalidateQueries({ queryKey: ["projects"] });
      queryClient.invalidateQueries({ queryKey: ["projects", updatedProject.id] });
      toast({
        title: "Project updated",
        description: `"${updatedProject.name}" has been updated`,
      });
    },
    onError: (error: Error) => {
      toast({
        variant: "destructive",
        title: "Failed to update project",
        description: error.message,
      });
    },
  });
}

export function useDeleteProject() {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: (id: number) => deleteProject(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects"] });
      toast({
        title: "Project deleted",
        description: "The project has been deleted",
      });
    },
    onError: (error: Error) => {
      toast({
        variant: "destructive",
        title: "Failed to delete project",
        description: error.message,
      });
    },
  });
}
