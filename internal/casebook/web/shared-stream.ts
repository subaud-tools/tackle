// shared-stream.ts — one live connection for every casebook tab.
//
// A browser keeps at most six connections open to one site, and each tab's
// live stream holds one for good: with six tabs open, every other request of
// every tab waited (Court, 2026-10-17). So the tabs share one stream: the tab
// holding a lock opens the real EventSource and passes each event to the
// others over a BroadcastChannel; when it closes, the next tab takes the lock
// and its own stream.
//
// sharedEventSource returns a class the kit's live() takes as its
// EventSource: it has onopen, onmessage, onerror and close, which is all
// live() uses. A tab that joins after the stream opened asks for "open" and
// gets it. When the stream errs, or the tab holding it closes, every tab's
// stream errs, and live() reconnects each with its newest cursor (polling
// meanwhile), so nothing between is lost.

/** What a stream message carries, as live() reads it. */
export interface StreamMessage {
  data: string;
  lastEventId: string;
}

interface Channel {
  onmessage: ((e: { data: unknown }) => void) | null;
  postMessage(data: unknown): void;
  close(): void;
}

/** The browser parts it needs (injected in tests). */
export interface SharedDeps {
  EventSource: new (url: string) => {
    onopen: (() => void) | null;
    onmessage: ((m: StreamMessage) => void) | null;
    onerror: (() => void) | null;
    close(): void;
  };
  channel(): Channel;
  /** Run granted while holding the one lock; its promise ending releases it. */
  lock(granted: () => Promise<void>, signal: AbortSignal): void;
}

type Note =
  | { k: 'open' }
  | { k: 'msg'; data: string; id: string }
  | { k: 'error' }
  | { k: 'hello' };

export interface SharedSource {
  onopen: (() => void) | null;
  onmessage: ((m: StreamMessage) => void) | null;
  onerror: (() => void) | null;
  close(): void;
}

export function sharedEventSource(
  deps: SharedDeps,
): new (url: string) => SharedSource {
  return class Shared implements SharedSource {
    onopen: (() => void) | null = null;
    onmessage: ((m: StreamMessage) => void) | null = null;
    onerror: (() => void) | null = null;
    private closed = false;
    private opened = false;
    private leading = false;
    private real: InstanceType<SharedDeps['EventSource']> | null = null;
    private release: (() => void) | null = null;
    private readonly abort = new AbortController();
    private readonly ch: Channel;

    constructor(url: string) {
      this.ch = deps.channel();
      this.ch.onmessage = (e) => this.heard(e.data as Note);
      // Some tab may already hold the stream: ask whether it is open, once
      // live() has set its handlers.
      queueMicrotask(() => {
        if (!this.closed) this.ch.postMessage({ k: 'hello' } satisfies Note);
      });
      deps.lock(() => this.granted(url), this.abort.signal);
    }

    // heard is another tab's note.
    private heard(n: Note): void {
      if (this.closed) return;
      if (this.leading) {
        if (n.k === 'hello' && this.opened)
          this.ch.postMessage({ k: 'open' } satisfies Note);
        return;
      }
      if (n.k === 'open' && !this.opened) {
        this.opened = true;
        this.onopen?.();
      } else if (n.k === 'msg') {
        this.onmessage?.({ data: n.data, lastEventId: n.id });
      } else if (n.k === 'error') {
        this.onerror?.();
      }
    }

    // granted: this tab holds the lock. A stream that was following another
    // tab's has an old cursor in its address: it asks live() to reconnect
    // (a fresh one takes the lock) instead of replaying from it.
    private granted(url: string): Promise<void> {
      if (this.closed) return Promise.resolve();
      if (this.opened) {
        this.onerror?.();
        return Promise.resolve();
      }
      this.leading = true;
      const done = new Promise<void>((r) => (this.release = r));
      const real = new deps.EventSource(url);
      this.real = real;
      real.onopen = () => {
        if (this.closed) return;
        this.opened = true;
        this.ch.postMessage({ k: 'open' } satisfies Note);
        this.onopen?.();
      };
      real.onmessage = (m) => {
        if (this.closed) return;
        this.ch.postMessage({
          k: 'msg',
          data: m.data,
          id: m.lastEventId,
        } satisfies Note);
        this.onmessage?.(m);
      };
      real.onerror = () => {
        if (this.closed) return;
        this.ch.postMessage({ k: 'error' } satisfies Note);
        this.onerror?.();
      };
      return done;
    }

    close(): void {
      if (this.closed) return;
      this.closed = true;
      if (this.leading) {
        // The others followed this stream: they reconnect.
        this.ch.postMessage({ k: 'error' } satisfies Note);
        this.real?.close();
        this.real = null;
      }
      this.abort.abort();
      this.release?.();
      this.ch.close();
    }
  };
}

/**
 * browserShared is the EventSource live() takes: shared between tabs where
 * the browser has locks and BroadcastChannel, the plain one otherwise.
 */
export function browserShared(name: string): typeof EventSource {
  const locks = (navigator as Navigator & { locks?: LockManager }).locks;
  if (!locks || typeof BroadcastChannel === 'undefined') return EventSource;
  return sharedEventSource({
    EventSource: EventSource as unknown as SharedDeps['EventSource'],
    channel: () => new BroadcastChannel(name) as unknown as Channel,
    lock(granted, signal) {
      locks.request(name, { signal }, granted).catch(() => {
        // aborted while waiting: the stream closed before its turn
      });
    },
  }) as unknown as typeof EventSource;
}
