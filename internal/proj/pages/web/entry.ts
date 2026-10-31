// The terminal page: xterm.js over a websocket to a tmux client. Binary frames
// carry terminal bytes both ways; a text frame from here is a resize.

import { FitAddon } from '@xterm/addon-fit';
import { Unicode11Addon } from '@xterm/addon-unicode11';
import { WebglAddon } from '@xterm/addon-webgl';
import { Terminal } from '@xterm/xterm';
import { closedView, takeFrom, withTake } from './closed.ts';
import { handleKey } from './keys.ts';

const fontFamily =
  '"Session Icons", "MonoLisaCode", "FiraCode Nerd Font", monospace';
const encoder = new TextEncoder();

async function main(): Promise<void> {
  // Load before the first draw: the renderer measures and caches glyphs.
  await document.fonts.load('13px "Session Icons"', '\ue0b0').catch(() => []);

  const term = new Terminal({
    fontFamily,
    fontSize: 13,
    allowProposedApi: true,
    macOptionIsMeta: true,
    macOptionClickForcesSelection: true,
  });
  const fit = new FitAddon();
  term.loadAddon(fit);
  term.loadAddon(new Unicode11Addon());
  term.unicode.activeVersion = '11';
  term.open(document.getElementById('term')!);
  try {
    const webgl = new WebglAddon();
    webgl.onContextLoss(() => webgl.dispose());
    term.loadAddon(webgl);
  } catch {
    // no WebGL here: xterm's DOM renderer still draws
  }
  fit.fit();
  term.focus();

  const address = takeFrom(location.search);
  history.replaceState(null, '', location.pathname + address.page);
  const ws = new WebSocket(
    `ws://${location.host}/api/term/ws${address.socket}`,
  );
  ws.binaryType = 'arraybuffer';
  const send = (data: string): void => {
    if (ws.readyState === WebSocket.OPEN) ws.send(encoder.encode(data));
  };
  const sendSize = (): void => {
    if (ws.readyState === WebSocket.OPEN) {
      ws.send(
        JSON.stringify({ type: 'resize', cols: term.cols, rows: term.rows }),
      );
    }
  };

  term.attachCustomKeyEventHandler((ev) => {
    // Only left Option is Alt; right Option still types its characters. The
    // modifier's own keydown is the one event that tells them apart.
    if (ev.type === 'keydown' && ev.key === 'Alt') {
      term.options.macOptionIsMeta = ev.code === 'AltLeft';
    }
    return handleKey(ev, send);
  });
  term.onData(send);
  term.onBinary((data) => {
    if (ws.readyState === WebSocket.OPEN) {
      ws.send(Uint8Array.from(data, (c) => c.charCodeAt(0)));
    }
  });
  term.onResize(sendSize);
  term.onTitleChange((title) => {
    document.title = title;
  });
  term.onSelectionChange(() => {
    const text = term.getSelection();
    if (text) void navigator.clipboard.writeText(text).catch(() => {});
  });
  // With tmux's mouse on, a drag selects in tmux, which copies through OSC 52.
  term.parser.registerOscHandler(52, (data) => {
    const b64 = data.slice(data.indexOf(';') + 1);
    if (b64 !== '?') {
      const bytes = Uint8Array.from(atob(b64), (c) => c.charCodeAt(0));
      const text = new TextDecoder().decode(bytes);
      void navigator.clipboard.writeText(text).catch(() => {});
    }
    return true;
  });
  new ResizeObserver(() => fit.fit()).observe(document.body);

  ws.onopen = sendSize;
  ws.onmessage = (ev) => term.write(new Uint8Array(ev.data as ArrayBuffer));
  ws.onclose = (ev) => {
    const view = closedView(ev.code, ev.reason);
    if (!view.useHere) {
      term.write(`\r\n[${view.text}]\r\n`);
      return;
    }
    const button = document.createElement('button');
    button.textContent = 'Use here';
    button.onclick = () => location.replace(withTake(location.search));
    const note = document.createElement('div');
    note.id = 'detached';
    note.append(view.text, button);
    document.body.append(note);
    button.focus();
  };
}

void main();
