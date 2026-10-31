// What the page says when its websocket closes.

/** The server's close code when another terminal took the session. */
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
