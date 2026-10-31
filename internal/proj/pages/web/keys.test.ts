import assert from 'node:assert/strict';
import { test } from 'node:test';
import { handleKey } from './keys.ts';

// Runs one press of a key through handleKey, as the browser's three events.
function press(mods: { shiftKey?: boolean; ctrlKey?: boolean }) {
  const sent: string[] = [];
  const handled = ['keydown', 'keypress', 'keyup'].map((type) =>
    handleKey(
      {
        type,
        key: 'Enter',
        shiftKey: false,
        ctrlKey: false,
        preventDefault() {},
        ...mods,
      },
      (s) => sent.push(s),
    ),
  );
  return { sent, handled };
}

test('Shift+Enter sends one newline, and xterm never sees it', () => {
  const { sent, handled } = press({ shiftKey: true });
  assert.deepEqual(sent, ['\n']);
  assert.deepEqual(handled, [false, false, false]);
});

test('Ctrl+Enter sends one newline, and xterm never sees it', () => {
  const { sent, handled } = press({ ctrlKey: true });
  assert.deepEqual(sent, ['\n']);
  assert.deepEqual(handled, [false, false, false]);
});

test('plain Enter is left to xterm, which sends a carriage return', () => {
  const { sent, handled } = press({});
  assert.deepEqual(sent, []);
  assert.deepEqual(handled, [true, true, true]);
});
