// text.ts — pull the assertion lines out of a test body.

const PY = [
  'assert ',
  'assert(',
  'with pytest.raises',
  'pytest.raises',
  'self.assert',
];
const GO = ['t.Error', 't.Fatal', 'assert.', 'require.'];
const TS = ['expect(', 'assert.', 'assert('];

export function assertLines(lang: string, body: string): string[] {
  const lines = body.split('\n').map((l) => l.trim());
  switch (lang) {
    case 'python':
      return lines.filter((l) => PY.some((p) => l.startsWith(p)));
    case 'go':
      return lines.filter((l) => GO.some((p) => l.includes(p)));
    case 'typescript':
    case 'javascript':
      return lines.filter((l) => TS.some((p) => l.includes(p)));
    default:
      return [];
  }
}

/** The bar's status after Send: where the answers went (to), if anywhere. */
export function sentText(n: number, to?: string): string {
  if (n === 0) return 'Nothing to send.';
  const what = `${n} ${n === 1 ? 'answer' : 'answers'}`;
  return to
    ? `Sent ${what} to ${to}.`
    : `Sent ${what}. The next agent that opens cull in this project gets them.`;
}

/** The line under a held send's pill. */
export const HELD_LINE =
  'the agent is mid-turn; it reads this when the turn ends';
/** The line under a sent send's pill when no session has claimed it. */
export const NO_SESSION_LINE = 'no session has it';

/** The line under a send's pill: the agent's line or reply once it is
 * working or settled, the mid-turn line while the harness holds it, or that
 * no session has it; '' for every other state. */
export function sendLine(s: RecentSend): string {
  switch (s.state) {
    case 'working':
    case 'answered':
    case 'declined':
    case 'failed':
      return s.reply;
    case 'sent':
      if (s.held) return HELD_LINE;
      return s.claimed ? '' : NO_SESSION_LINE;
    default:
      return '';
  }
}

/** How long ago a send was sent, from its RFC 3339 time. */
export function sentAgo(at: string, now: number): string {
  const min = Math.floor((now - Date.parse(at)) / 60000);
  if (!(min >= 1)) return 'just now';
  if (min < 60) return `${min} min ago`;
  if (min < 24 * 60) return `${Math.floor(min / 60)} h ago`;
  return `${Math.floor(min / (24 * 60))} d ago`;
}
