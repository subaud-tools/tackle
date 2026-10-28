// wire.d.ts — the Go↔page contract, generated. Do not edit.
//
// Written by internal/cull/serve/wire_test.go from the json tags of review
// (GET /api/review), judge.State and judge.GroupState, and every type they
// reach. An item's state and jev are unknown here: Go stores them as bytes
// whose shape depends on the item's kind, which web/kinds.d.ts narrows.
// Regenerate with `just cull-wire`.

export interface Review {
  cursor: string;
  project: ReviewProject;
  run: ReviewRun | null;
  items: ItemBase[] | null;
  answered?: AnsweredCounts;
  blind: boolean;
  recent_sends: RecentSend[] | null;
}

export interface State {
  note: string;
  language: string;
  framework: string;
  test_name: string;
  test_source: string;
  setup_context: string;
  code_under_test: Callee[] | null;
  truncated: boolean;
}

export interface GroupState {
  note: string;
  language: string;
  framework: string;
  file: string;
  tests: GroupTest[] | null;
  truncated: boolean;
}

export interface ReviewProject {
  id: number;
  root: string;
}

export interface ReviewRun {
  id: number;
  at: string;
  mode: string;
  base: string;
  total: number;
  summary: Record<string, number> | null;
}

export interface ItemBase {
  id: string;
  kind: string;
  hash: string;
  file: string;
  name: string;
  verdict: string;
  rule: string;
  model: string;
  state: unknown | null;
  jev: unknown | null;
  rows?: string[][];
  members?: string[];
  answer?: Answer;
}

export interface AnsweredCounts {
  total: number;
  sent: number;
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

export interface Callee {
  symbol: string;
  file: string;
  source: string;
}

export interface GroupTest {
  name: string;
  source: string;
}

export interface Answer {
  value: string;
  note: string;
  via: string;
  blind: boolean;
  answered_at?: string;
  sent_at?: string;
}

export interface Counts {
  cut: number;
  keep: number;
  merge: number;
  separate: number;
}

export type MessageState =
  | 'sent'
  | 'with_agent'
  | 'working'
  | 'answered'
  | 'declined'
  | 'failed'
  | 'no_reply';
