"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Search, CheckCircle, XCircle, Clock } from "lucide-react";

const stats = [
  {
    title: "Total URLs",
    value: "1,234",
    icon: Search,
  },
  {
    title: "Indexed",
    value: "1,089",
    icon: CheckCircle,
  },
  {
    title: "Not Indexed",
    value: "98",
    icon: XCircle,
  },
  {
    title: "Pending",
    value: "47",
    icon: Clock,
  },
];

export default function IndexCheckerPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Index Checker</h1>
        <p className="text-muted-foreground">
          Check if your URLs are indexed by search engines
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
            </CardContent>
          </Card>
        ))}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>URL Index Status</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between rounded-lg border p-4">
              <div>
                <p className="font-medium">https://example.com/page-1</p>
                <p className="text-sm text-muted-foreground">
                  Last checked: 2 hours ago
                </p>
              </div>
              <Badge variant="default">Indexed</Badge>
            </div>
            <div className="flex items-center justify-between rounded-lg border p-4">
              <div>
                <p className="font-medium">https://example.com/page-2</p>
                <p className="text-sm text-muted-foreground">
                  Last checked: 3 hours ago
                </p>
              </div>
              <Badge variant="destructive">Not Indexed</Badge>
            </div>
            <div className="flex items-center justify-between rounded-lg border p-4">
              <div>
                <p className="font-medium">https://example.com/page-3</p>
                <p className="text-sm text-muted-foreground">Checking...</p>
              </div>
              <Badge variant="secondary">Pending</Badge>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
