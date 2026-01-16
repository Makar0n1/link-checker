"use client";

import { useState, useMemo, useCallback } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Link2, TrendingUp, TrendingDown, AlertCircle } from "lucide-react";
import { DataTable, createBacklinksColumns } from "@/components/tables";
import { mockBacklinks } from "@/lib/mock-data";
import { Backlink } from "@/types";
import { useToast } from "@/hooks/use-toast";

const stats = [
  {
    title: "Total Backlinks",
    value: "2,847",
    change: "+12%",
    trend: "up",
    icon: Link2,
  },
  {
    title: "Active Links",
    value: "2,654",
    change: "+8%",
    trend: "up",
    icon: TrendingUp,
  },
  {
    title: "Lost Links",
    value: "193",
    change: "-3%",
    trend: "down",
    icon: TrendingDown,
  },
  {
    title: "Issues",
    value: "24",
    change: "+5",
    trend: "up",
    icon: AlertCircle,
  },
];

export default function BacklinksPage() {
  const { toast } = useToast();
  const [backlinks, setBacklinks] = useState<Backlink[]>(mockBacklinks);

  const handleUpdateBacklink = useCallback(
    (id: string, field: keyof Backlink, value: string) => {
      setBacklinks((prev) =>
        prev.map((link) =>
          link.id === id ? { ...link, [field]: value } : link
        )
      );
      toast({
        title: "Updated",
        description: `Backlink ${field} updated successfully.`,
      });
    },
    [toast]
  );

  const handleBulkDelete = useCallback(
    (rows: Backlink[]) => {
      const ids = new Set(rows.map((r) => r.id));
      setBacklinks((prev) => prev.filter((link) => !ids.has(link.id)));
      toast({
        title: "Deleted",
        description: `${rows.length} backlink(s) deleted.`,
        variant: "destructive",
      });
    },
    [toast]
  );

  const columns = useMemo(
    () => createBacklinksColumns({ onUpdateBacklink: handleUpdateBacklink }),
    [handleUpdateBacklink]
  );

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Backlinks</h1>
        <p className="text-muted-foreground">
          Monitor and manage your backlink profile
        </p>
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
              <div className="text-2xl font-bold">{stat.value}</div>
              <Badge
                variant={stat.trend === "up" ? "default" : "destructive"}
                className="mt-1"
              >
                {stat.change}
              </Badge>
            </CardContent>
          </Card>
        ))}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Backlinks Table</CardTitle>
        </CardHeader>
        <CardContent>
          <DataTable
            columns={columns}
            data={backlinks}
            searchKey="sourceUrl"
            searchPlaceholder="Filter by source URL..."
            onBulkDelete={handleBulkDelete}
          />
        </CardContent>
      </Card>
    </div>
  );
}
