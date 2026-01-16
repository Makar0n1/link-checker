export interface Backlink {
  id: string;
  sourceUrl: string;
  targetUrl: string;
  anchorText: string;
  status: "active" | "lost" | "pending";
  dofollow: boolean;
  domainAuthority: number;
  lastChecked: string;
  firstSeen: string;
}
