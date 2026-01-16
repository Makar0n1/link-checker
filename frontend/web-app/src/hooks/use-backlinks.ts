"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  getBacklinks,
  getBacklink,
  createBacklink,
  updateBacklink,
  deleteBacklink,
  bulkCreateBacklinks,
  bulkDeleteBacklinks,
} from "@/lib/api";
import type {
  BacklinksQueryParams,
  CreateBacklinkRequest,
  UpdateBacklinkRequest,
  BulkCreateBacklinksRequest,
  BulkDeleteBacklinksRequest,
  Backlink,
  PaginatedResponse,
} from "@/types/api";
import { useToast } from "./use-toast";

export function useBacklinks(params: BacklinksQueryParams) {
  return useQuery<PaginatedResponse<Backlink>>({
    queryKey: ["backlinks", params],
    queryFn: () => getBacklinks(params),
    enabled: !!params.project_id,
  });
}

export function useBacklink(id: number) {
  return useQuery<Backlink>({
    queryKey: ["backlinks", "detail", id],
    queryFn: () => getBacklink(id),
    enabled: !!id,
  });
}

export function useCreateBacklink() {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: (data: CreateBacklinkRequest) => createBacklink(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["backlinks"] });
      toast({
        title: "Backlink created",
        description: "The backlink has been added",
      });
    },
    onError: (error: Error) => {
      toast({
        variant: "destructive",
        title: "Failed to create backlink",
        description: error.message,
      });
    },
  });
}

export function useUpdateBacklink() {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateBacklinkRequest }) =>
      updateBacklink(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["backlinks"] });
      toast({
        title: "Backlink updated",
        description: "The backlink has been updated",
      });
    },
    onError: (error: Error) => {
      toast({
        variant: "destructive",
        title: "Failed to update backlink",
        description: error.message,
      });
    },
  });
}

export function useDeleteBacklink() {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: (id: number) => deleteBacklink(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["backlinks"] });
      toast({
        title: "Backlink deleted",
        description: "The backlink has been removed",
      });
    },
    onError: (error: Error) => {
      toast({
        variant: "destructive",
        title: "Failed to delete backlink",
        description: error.message,
      });
    },
  });
}

export function useBulkCreateBacklinks() {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: (data: BulkCreateBacklinksRequest) => bulkCreateBacklinks(data),
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: ["backlinks"] });
      toast({
        title: "Bulk create completed",
        description: `${result.success} backlinks created, ${result.failed} failed`,
      });
    },
    onError: (error: Error) => {
      toast({
        variant: "destructive",
        title: "Bulk create failed",
        description: error.message,
      });
    },
  });
}

export function useBulkDeleteBacklinks() {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: (data: BulkDeleteBacklinksRequest) => bulkDeleteBacklinks(data),
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: ["backlinks"] });
      toast({
        title: "Bulk delete completed",
        description: `${result.success} backlinks deleted, ${result.failed} failed`,
      });
    },
    onError: (error: Error) => {
      toast({
        variant: "destructive",
        title: "Bulk delete failed",
        description: error.message,
      });
    },
  });
}
