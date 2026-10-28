// The page's names for the wire (wire.d.ts, generated from the Go structs),
// narrowed where the page knows more than the Go types say: values Go keeps
// in a plain string (an action, a state), and the arrays Go always builds
// (wiregen declares every slice without omitempty `| null`, since a nil one
// is). Ambient.

type Action = 'accept' | 'edit' | 'reject';

type Source = import('./wire.d.ts').Source;

type Decision = import('./wire.d.ts').Decision & {
  action: Action;
  cleared?: ('title' | 'text')[];
};

/** A file's recommendation (internal/sift/rec.Rec). Findings is null for
 * one that accounts for none. */
type Recommendation = Omit<import('./wire.d.ts').Recommendation, 'findings'> & {
  findings: Account[] | null;
};

type Account = import('./wire.d.ts').Account & { did: 'fixed' | 'kept' };

type FileDecision = import('./wire.d.ts').FileDecision & { action: Action };

/** One of an audit round's files (internal/sift/serve fileView). */
type FileView = Omit<
  import('./wire.d.ts').FileView,
  'rows' | 'group' | 'rec' | 'decision'
> & {
  rows: string[];
  group: string[];
  rec: Recommendation | null;
  decision: FileDecision | null;
};

/** A row; Go fills fingerprint on every read (row.Row.Print). */
type Finding = Omit<
  import('./wire.d.ts').Finding,
  'decision' | 'fingerprint'
> & {
  decision?: Decision;
  fingerprint: string;
};

type Progress = Omit<import('./wire.d.ts').Progress, 'state'> & {
  state: '' | 'checking' | 'recommending' | 'ready' | 'sent' | 'applied';
};

type ReviewRound = Omit<import('./wire.d.ts').ReviewRound, 'summary'> & {
  summary: Record<string, number>;
};

type ApplyRecord = Omit<import('./wire.d.ts').ApplyRecord, 'rows'> & {
  rows: string[];
};

type RecentSend = import('./wire.d.ts').RecentSend;

type Review = Omit<
  import('./wire.d.ts').Review,
  'round' | 'progress' | 'rows' | 'files' | 'applies' | 'recent_sends'
> & {
  round: ReviewRound | null;
  progress: Progress;
  rows: Finding[];
  files: FileView[];
  applies: ApplyRecord[];
  recent_sends: RecentSend[];
};
