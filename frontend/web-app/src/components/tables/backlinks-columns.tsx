"use client";

import { ColumnDef } from "@tanstack/react-table";
import { Backlink } from "@/types";
import { Badge } from "@/components/ui/badge";
import { Checkbox } from "@/components/ui/checkbox";
import { Button } from "@/components/ui/button";
import { ArrowUpDown, ExternalLink } from "lucide-react";
import { EditableCell } from "./editable-cell";

interface CreateColumnsOptions {
  onUpdateBacklink: (id: string, field: keyof Backlink, value: string) => void;
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
      accessorKey: "sourceUrl",
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
        const url = row.getValue("sourceUrl") as string;
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
      accessorKey: "targetUrl",
      header: "Target URL",
      cell: ({ row }) => {
        const url = row.getValue("targetUrl") as string;
        return (
          <span className="max-w-[150px] truncate" title={url}>
            {url.replace(/^https?:\/\//, "")}
          </span>
        );
      },
    },
    {
      accessorKey: "anchorText",
      header: "Anchor Text",
      cell: ({ row }) => {
        const backlink = row.original;
        return (
          <EditableCell
            value={backlink.anchorText}
            onSave={(value) => onUpdateBacklink(backlink.id, "anchorText", value)}
          />
        );
      },
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const status = row.getValue("status") as string;
        return (
          <Badge
            variant={
              status === "active"
                ? "default"
                : status === "lost"
                  ? "destructive"
                  : "secondary"
            }
          >
            {status}
          </Badge>
        );
      },
      filterFn: (row, id, value) => {
        return value.includes(row.getValue(id));
      },
    },
    {
      accessorKey: "dofollow",
      header: "Type",
      cell: ({ row }) => {
        const dofollow = row.getValue("dofollow") as boolean;
        return (
          <Badge variant={dofollow ? "default" : "outline"}>
            {dofollow ? "dofollow" : "nofollow"}
          </Badge>
        );
      },
    },
    {
      accessorKey: "domainAuthority",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          DA
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => {
        const da = row.getValue("domainAuthority") as number;
        return (
          <span
            className={`font-medium ${
              da >= 60
                ? "text-green-600"
                : da >= 40
                  ? "text-yellow-600"
                  : "text-muted-foreground"
            }`}
          >
            {da}
          </span>
        );
      },
    },
    {
      accessorKey: "lastChecked",
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
        const date = new Date(row.getValue("lastChecked") as string);
        return (
          <span className="text-muted-foreground">
            {date.toLocaleDateString()}
          </span>
        );
      },
    },
  ];
}
