// The page's own key bindings, apart from xterm so a test can run them.

export interface KeyLike {
  type: string;
  key: string;
  shiftKey: boolean;
  ctrlKey: boolean;
  preventDefault(): void;
}

// handleKey is xterm's custom key handler: false means xterm must ignore the
// event. Shift+Enter and Ctrl+Enter send a newline, as the Ghostty keybinds do.
export function handleKey(ev: KeyLike, send: (data: string) => void): boolean {
  if (ev.key !== 'Enter' || (!ev.shiftKey && !ev.ctrlKey)) return true;
  ev.preventDefault();
  if (ev.type === 'keydown') send('\n');
  return false;
}
