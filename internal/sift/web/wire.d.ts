// wire.d.ts — the Go↔page contract, generated. Do not edit.
//
// Written by internal/sift/serve/wire_test.go from the json tags of review
// (GET /api/review) and status (GET /api/agent/status), and every type they
// reach. Regenerate with `just sift-wire`.

export interface Review {
  cursor: string;
  round: ReviewRound | null;
  progress: Progress;
  rows: Finding[] | null;
  files: FileView[] | null;
  applies: ApplyRecord[] | null;
  sends: number;
  recent_sends: RecentSend[] | null;
  home: string;
}

export interface Status {
  round: number;
  state: string;
  files: number;
  recommended: number;
  open: number;
  decided: number;
  sent: number;
  owner: string;
  owner_present: boolean;
  unsettled: UnsettledSend[] | null;
}

export interface ReviewRound {
  id: number;
  kind: string;
  at: string;
  summary: Record<string, number> | null;
  owner: string;
}

export interface Progress {
  state: string;
  files: number;
  recommended: number;
}

export interface Finding {
  id: string;
  check: string;
  summary: string;
  source: Source;
  passage?: string;
  evidence?: Fact[];
  verdict?: string;
  title?: string;
  destination?: string;
  text?: string;
  reason?: string;
  certain: boolean;
  decision?: Decision;
  fingerprint?: string;
}

export interface FileView {
  key: string;
  path: string;
  source: Source;
  commit?: string;
  class: string;
  budget: number;
  base: string;
  size: number;
  after: number;
  rows: string[] | null;
  rec: Recommendation | null;
  decision: FileDecision | null;
  fingerprint: string;
  group: string[] | null;
  unchanged: boolean;
  muted: boolean;
}

export interface ApplyRecord {
  repo: string;
  base: string;
  branch: string;
  pr: string;
  state: string;
  detail: string;
  rows: string[] | null;
  at: string;
  last_state?: string;
  last_detail?: string;
}

export interface RecentSend {
  id: number;
  at: string;
  counts: Counts;
  state: MessageState;
  held: boolean;
  claimed: boolean;
  reply: string;
}

export interface UnsettledSend {
  id: number;
  round: number;
  state: MessageState;
  at: string;
  counts: Counts;
}

export interface Source {
  file: string;
  repo?: string;
  ref?: string;
  path?: string;
  start?: number;
  end?: number;
  entry?: string;
  canon?: string;
}

export interface Fact {
  name: string;
  value: string;
}

export interface Decision {
  action: string;
  verdict?: string;
  title?: string;
  text?: string;
  cleared?: string[];
  note?: string;
  sent?: boolean;
  id?: string;
}

export interface Recommendation {
  file: string;
  base: string;
  content: string;
  findings: Account[] | null;
  links?: string[];
  summary: string;
}

export interface FileDecision {
  action: string;
  content?: string;
  note?: string;
  sent?: boolean;
  id?: string;
}

export interface Counts {
  accept: number;
  edit: number;
  reject: number;
}

export type MessageState =
  | 'sent'
  | 'with_agent'
  | 'working'
  | 'answered'
  | 'declined'
  | 'failed'
  | 'no_reply';

export interface Account {
  row: string;
  did: string;
  how: string;
}
