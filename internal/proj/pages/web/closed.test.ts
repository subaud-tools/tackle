import assert from 'node:assert/strict';
import { test } from 'node:test';
import { closedView } from './closed.ts';

test('a client detached by another terminal offers Use here', () => {
  assert.deepEqual(closedView(4409, 'detached'), {
    useHere: true,
    text: 'Shown in another terminal',
  });
});

test('any other close says why and to reload, with no Use here', () => {
  assert.deepEqual(closedView(1000, 'session ended'), {
    useHere: false,
    text: 'session ended: reload to reconnect',
  });
  assert.deepEqual(closedView(1006, ''), {
    useHere: false,
    text: 'connection lost: reload to reconnect',
  });
});
