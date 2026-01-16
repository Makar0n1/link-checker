"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Activity, Server, Clock, AlertTriangle } from "lucide-react";

const stats = [
  {
    title: "Sites Monitored",
    value: "12",
    icon: Server,
  },
  {
    title: "Healthy",
    value: "10",
    icon: Activity,
  },
  {
    title: "Avg Response",
    value: "245ms",
    icon: Clock,
  },
  {
    title: "Alerts",
    value: "2",
    icon: AlertTriangle,
  },
];

const sites = [
  {
    url: "https://mysite.com",
    status: "healthy",
    uptime: "99.9%",
    responseTime: "189ms",
  },
  {
    url: "https://blog.mysite.com",
    status: "healthy",
    uptime: "99.8%",
    responseTime: "234ms",
  },
  {
    url: "https://shop.mysite.com",
    status: "warning",
    uptime: "98.5%",
    responseTime: "512ms",
  },
  {
    url: "https://api.mysite.com",
    status: "down",
    uptime: "95.2%",
    responseTime: "N/A",
  },
];

export default function SiteHealthPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Site Health</h1>
        <p className="text-muted-foreground">
          Monitor uptime and performance of your sites
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
          <CardTitle>Monitored Sites</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {sites.map((site) => (
              <div
                key={site.url}
                className="flex items-center justify-between rounded-lg border p-4"
              >
                <div className="flex items-center gap-4">
                  <div
                    className={`h-3 w-3 rounded-full ${
                      site.status === "healthy"
                        ? "bg-green-500"
                        : site.status === "warning"
                          ? "bg-yellow-500"
                          : "bg-red-500"
                    }`}
                  />
                  <div>
                    <p className="font-medium">{site.url}</p>
                    <p className="text-sm text-muted-foreground">
                      Uptime: {site.uptime} | Response: {site.responseTime}
                    </p>
                  </div>
                </div>
                <Badge
                  variant={
                    site.status === "healthy"
                      ? "default"
                      : site.status === "warning"
                        ? "secondary"
                        : "destructive"
                  }
                >
                  {site.status}
                </Badge>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
