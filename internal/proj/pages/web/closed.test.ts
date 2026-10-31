import assert from 'node:assert/strict';
import { test } from 'node:test';
import { closedView, takeFrom, withTake } from './closed.ts';

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

test('take goes to the socket but comes off the page address', () => {
  assert.deepEqual(takeFrom('?session=s&socket=a&take=1'), {
    socket: '?session=s&socket=a&take=1',
    page: '?session=s&socket=a',
  });
});

test('an address without take connects without it', () => {
  assert.deepEqual(takeFrom('?session=s&socket=a'), {
    socket: '?session=s&socket=a',
    page: '?session=s&socket=a',
  });
});

test('Use here adds take to the address', () => {
  assert.equal(withTake('?session=s&socket=a'), '?session=s&socket=a&take=1');
});
