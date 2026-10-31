// What the page says when its websocket closes.

/** The server's close code when another terminal has the session. */
const closeDetached = 4409;

export interface ClosedView {
  useHere: boolean;
  text: string;
}

export function closedView(code: number, reason: string): ClosedView {
  if (code === closeDetached) {
    return { useHere: true, text: 'Shown in another terminal' };
  }
  return {
    useHere: false,
    text: `${reason || 'connection lost'}: reload to reconnect`,
  };
}

/** The address's take=1 goes on the socket's address and off the page's, so a
 * tab Chrome restores or reloads doesn't take the session. */
export function takeFrom(search: string): { socket: string; page: string } {
  const q = new URLSearchParams(search);
  q.delete('take');
  return { socket: search, page: `?${q.toString()}` };
}

/** The address Use here loads. */
export function withTake(search: string): string {
  const q = new URLSearchParams(search);
  q.set('take', '1');
  return `?${q.toString()}`;
}
