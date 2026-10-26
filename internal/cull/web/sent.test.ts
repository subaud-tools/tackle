import { test } from 'node:test';
import assert from 'node:assert/strict';
import { sendLine, sentText } from './text.ts';

test('Send status names the session the answers went to', () => {
  assert.equal(sentText(3, 'pi: shop'), 'Sent 3 answers to pi: shop.');
  assert.equal(sentText(1, 'pi: shop'), 'Sent 1 answer to pi: shop.');
});

test('Send status without a present owner', () => {
  assert.equal(
    sentText(2, ''),
    'Sent 2 answers. The next agent that opens cull in this project gets them.',
  );
  assert.equal(
    sentText(1, undefined),
    'Sent 1 answer. The next agent that opens cull in this project gets them.',
  );
});

test('nothing to send', () => {
  assert.equal(sentText(0, 'pi: shop'), 'Nothing to send.');
});

test('sendLine: the line under each state of a send', () => {
  const send = (o: Partial<RecentSend>): RecentSend => ({
    id: 1,
    at: '2026-10-26T09:00:00Z',
    counts: { cut: 1, keep: 0, merge: 0, separate: 0 },
    state: 'sent',
    held: false,
    claimed: true,
    reply: '',
    ...o,
  });
  const cases: [Partial<RecentSend>, string][] = [
    [{ state: 'sent', claimed: false }, 'no session has it'],
    [
      { state: 'sent', held: true },
      'the agent is mid-turn; it reads this when the turn ends',
    ],
    [{ state: 'sent' }, ''],
    [{ state: 'with_agent' }, ''],
    [{ state: 'working', reply: 'cutting 3 tests' }, 'cutting 3 tests'],
    [{ state: 'answered', reply: '3 tests cut' }, '3 tests cut'],
    [{ state: 'declined', reply: 'not this branch' }, 'not this branch'],
    [{ state: 'failed', reply: 'the suite broke' }, 'the suite broke'],
    [{ state: 'no_reply', reply: 'stale line' }, ''],
  ];
  for (const [o, want] of cases)
    assert.equal(sendLine(send(o)), want, JSON.stringify(o));
});
