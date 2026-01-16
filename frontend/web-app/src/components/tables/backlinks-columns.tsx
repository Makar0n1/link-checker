"use client";

import { ColumnDef } from "@tanstack/react-table";
import type { Backlink } from "@/types/api";
import { Badge } from "@/components/ui/badge";
import { Checkbox } from "@/components/ui/checkbox";
import { Button } from "@/components/ui/button";
import { ArrowUpDown, ExternalLink } from "lucide-react";
import { EditableCell } from "./editable-cell";

interface CreateColumnsOptions {
  onUpdateBacklink: (id: number, field: keyof Backlink, value: string) => void;
}

export function createBacklinksColumns({
  onUpdateBacklink,
}: CreateColumnsOptions): ColumnDef<Backlink>[] {
  return [
    {
      id: "select",
      header: ({ table }) => (
        <Checkbox
          checked={
            table.getIsAllPageRowsSelected() ||
            (table.getIsSomePageRowsSelected() && "indeterminate")
          }
          onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
          aria-label="Select all"
        />
      ),
      cell: ({ row }) => (
        <Checkbox
          checked={row.getIsSelected()}
          onCheckedChange={(value) => row.toggleSelected(!!value)}
          aria-label="Select row"
        />
      ),
      enableSorting: false,
      enableHiding: false,
    },
    {
      accessorKey: "source_url",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Source URL
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => {
        const url = row.getValue("source_url") as string;
        return (
          <div className="flex items-center gap-2">
            <a
              href={url}
              target="_blank"
              rel="noopener noreferrer"
              className="max-w-[200px] truncate text-blue-600 hover:underline"
              title={url}
            >
              {url.replace(/^https?:\/\//, "")}
            </a>
            <ExternalLink className="h-3 w-3 text-muted-foreground" />
          </div>
        );
      },
    },
    {
      accessorKey: "target_url",
      header: "Target URL",
      cell: ({ row }) => {
        const url = row.getValue("target_url") as string;
        return (
          <span className="max-w-[150px] truncate" title={url}>
            {url.replace(/^https?:\/\//, "")}
          </span>
        );
      },
    },
    {
      accessorKey: "anchor_text",
      header: "Anchor Text",
      cell: ({ row }) => {
        const backlink = row.original;
        return (
          <EditableCell
            value={backlink.anchor_text}
            onSave={(value) => onUpdateBacklink(backlink.id, "anchor_text", value)}
          />
        );
      },
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const status = row.getValue("status") as string;
        const variants: Record<string, "default" | "destructive" | "secondary" | "outline"> = {
          active: "default",
          pending: "secondary",
          broken: "destructive",
          removed: "destructive",
          nofollow: "outline",
        };
        return (
          <Badge variant={variants[status] || "secondary"}>
            {status}
          </Badge>
        );
      },
      filterFn: (row, id, value) => {
        return value.includes(row.getValue(id));
      },
    },
    {
      accessorKey: "link_type",
      header: "Type",
      cell: ({ row }) => {
        const linkType = row.getValue("link_type") as string;
        return (
          <Badge variant={linkType === "dofollow" ? "default" : "outline"}>
            {linkType}
          </Badge>
        );
      },
    },
    {
      accessorKey: "http_status",
      header: "HTTP",
      cell: ({ row }) => {
        const httpStatus = row.getValue("http_status") as number | null;
        if (!httpStatus) return <span className="text-muted-foreground">-</span>;
        return (
          <span
            className={`font-medium ${
              httpStatus >= 200 && httpStatus < 300
                ? "text-green-600"
                : httpStatus >= 400
                  ? "text-red-600"
                  : "text-yellow-600"
            }`}
          >
            {httpStatus}
          </span>
        );
      },
    },
    {
      accessorKey: "last_checked_at",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Last Checked
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => {
        const dateStr = row.getValue("last_checked_at") as string | null;
        if (!dateStr) return <span className="text-muted-foreground">Never</span>;
        const date = new Date(dateStr);
        return (
          <span className="text-muted-foreground">
            {date.toLocaleDateString()}
          </span>
        );
      },
    },
  ];
}
