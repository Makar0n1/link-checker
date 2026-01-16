import Link from "next/link";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Link2,
  Search,
  Activity,
  ArrowRight,
  LayoutDashboard,
} from "lucide-react";

const features = [
  {
    title: "Backlinks Monitor",
    description:
      "Track your backlink profile in real-time. Monitor new and lost links, anchor texts, and domain authority.",
    icon: Link2,
  },
  {
    title: "Index Checker",
    description:
      "Verify if your pages are indexed by search engines. Get instant status updates and historical data.",
    icon: Search,
  },
  {
    title: "Site Health",
    description:
      "Monitor uptime and performance of your sites. Get alerts when issues occur.",
    icon: Activity,
  },
];

export default function LandingPage() {
  return (
    <div className="flex min-h-screen flex-col">
      <header className="sticky top-0 z-50 border-b bg-background/95 backdrop-blur">
        <div className="container mx-auto flex h-16 items-center justify-between px-4">
          <Link href="/" className="flex items-center gap-2">
            <LayoutDashboard className="h-6 w-6 text-primary" />
            <span className="text-xl font-bold">Link Tracker</span>
          </Link>
          <div className="flex items-center gap-4">
            <Link href="/login">
              <Button variant="ghost">Sign In</Button>
            </Link>
            <Link href="/register">
              <Button>Get Started</Button>
            </Link>
          </div>
        </div>
      </header>

      <main className="flex-1">
        <section className="container mx-auto px-4 py-24 text-center">
          <h1 className="text-4xl font-bold tracking-tight sm:text-5xl md:text-6xl">
            Monitor Your SEO
            <br />
            <span className="text-primary">Like a Pro</span>
          </h1>
          <p className="mx-auto mt-6 max-w-2xl text-lg text-muted-foreground">
            Track backlinks, check indexation status, and monitor site health —
            all in one powerful tool built for SEO teams.
          </p>
          <div className="mt-10 flex items-center justify-center gap-4">
            <Link href="/register">
              <Button size="lg" className="gap-2">
                Get Started Free
                <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>
            <Link href="/login">
              <Button size="lg" variant="outline">
                Sign In
              </Button>
            </Link>
          </div>
        </section>

        <section className="container mx-auto px-4 py-24">
          <h2 className="mb-12 text-center text-3xl font-bold">
            Everything You Need for SEO Monitoring
          </h2>
          <div className="grid gap-6 md:grid-cols-3">
            {features.map((feature) => (
              <Card key={feature.title} className="text-center">
                <CardHeader>
                  <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
                    <feature.icon className="h-6 w-6 text-primary" />
                  </div>
                  <CardTitle>{feature.title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <CardDescription className="text-base">
                    {feature.description}
                  </CardDescription>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>

        <section className="border-t bg-muted/50">
          <div className="container mx-auto px-4 py-24 text-center">
            <h2 className="text-3xl font-bold">Ready to Get Started?</h2>
            <p className="mx-auto mt-4 max-w-xl text-muted-foreground">
              Join SEO professionals who trust Link Tracker to monitor their
              backlink profiles and site health.
            </p>
            <div className="mt-8">
              <Link href="/register">
                <Button size="lg" className="gap-2">
                  Create Free Account
                  <ArrowRight className="h-4 w-4" />
                </Button>
              </Link>
            </div>
          </div>
        </section>
      </main>

      <footer className="border-t py-8">
        <div className="container mx-auto px-4 text-center text-sm text-muted-foreground">
          <p>&copy; 2026 Link Tracker. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
}
