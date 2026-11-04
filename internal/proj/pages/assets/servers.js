// servers.js — proj serve's servers page: every family page server, up then
// down, and the other family processes no server entry names (shown as "also
// running": test builds or strays), with their actions; reloaded when serve
// says the list changed.
import { bar, buttons, codeBlock, createApi, h, initTheme, live, sheet } from '/_kit/kit.js';

const app = document.getElementById('app');
const top = bar({ brand: { name: 'proj' }, status: 'servers', staleText: 'restarted · reload' });
initTheme('proj', top.themeControl);
const api = createApi({ onStale: () => top.setLive('stale') });

const groups = h('div');
app.append(top.el, h('main', { class: 'kit-read' }, h('div', { class: 'srv' }, groups)));

// busy holds the row an action runs on, so its buttons wait for it.
const busy = new Set();

async function act(key, label, path) {
  busy.add(key);
  await load();
  top.setStatus(`${label}…`);
  try {
    await api.post(path);
    top.setStatus(`${label}: done`);
  } catch (err) {
    top.setStatus(`${label}: ${err.message || err}`, { tone: 'danger' });
  } finally {
    busy.delete(key);
    await load();
  }
}

async function showLog(tool) {
  try {
    const { log } = await api.get(`/servers/${encodeURIComponent(tool)}/log`);
    const s = sheet({ title: `${tool} · serve.log`, body: log ? codeBlock(log) : h('p', { class: 'srv-empty' }, 'the log is empty'),
      actions: [{ label: 'close', run: () => s.close() }] });
  } catch (err) {
    top.setStatus(`log: ${err.message || err}`, { tone: 'danger' });
  }
}

function serverRow(s) {
  const t = encodeURIComponent(s.tool);
  const wait = busy.has(s.tool);
  const acts = s.up
    ? [
        { label: 'open', disabled: !s.url, run: () => window.open(s.url, '_blank', 'noopener') },
        { label: 'log', run: () => showLog(s.tool) },
        { label: 'restart', fill: !!s.note, disabled: wait, run: () => act(s.tool, `restart ${s.tool}`, `/servers/${t}/restart`) },
        { label: 'stop', danger: true, disabled: wait, run: () => act(s.tool, `stop ${s.tool}`, `/servers/${t}/stop`) },
      ]
    : [
        { label: 'log', run: () => showLog(s.tool) },
        { label: 'start', disabled: wait, run: () => act(s.tool, `start ${s.tool}`, `/servers/${t}/start`) },
      ];
  return h('div', { class: 'kit-tr', dataset: { tool: s.tool } },
    h('span', { class: 'srv-name' }, h('i', { dataset: { tone: s.up ? 'up' : 'down' } }), s.tool),
    h('span', { class: 'srv-mono' }, s.up ? [s.address, s.uptime && `up ${s.uptime}`].filter(Boolean).join(' · ') : 'not running'),
    h('span', { class: 'srv-mono' }, s.version || ''),
    h('span', { class: 'srv-note' }, s.note ? `${s.note}: restart to run it` : ''),
    buttons(acts));
}

function leftRow(p) {
  const key = `pid ${p.pid}`;
  return h('div', { class: 'kit-tr', dataset: { pid: String(p.pid) } },
    h('span', { class: 'srv-name' }, h('i', { dataset: { tone: 'other' } }), p.tool),
    h('span', { class: 'srv-mono', title: p.command }, p.command),
    h('span', { class: 'srv-mono' }, `pid ${p.pid}`),
    h('span', { class: 'srv-mono' }, `up ${p.age}`),
    buttons([{ label: 'stop', danger: true, disabled: busy.has(key), run: () => act(key, `stop pid ${p.pid}`, `/leftovers/${p.pid}/stop`) }]));
}

function group(label, rows, row, sub) {
  if (!rows.length) return null;
  return [h('div', { class: 'kit-label' }, `${label} · ${rows.length}`), sub && h('p', { class: 'srv-sub' }, sub), h('div', { class: 'kit-table' }, rows.map(row))].filter(Boolean);
}

function render(l) {
  const left = l.left || [];
  const parts = [group('up', l.up || [], serverRow), group('down', l.down || [], serverRow),
    group('also running', left, leftRow, "family processes that aren't the registered server: test builds or strays")].filter(Boolean);
  groups.replaceChildren(...(parts.length ? parts.flat() : [h('p', { class: 'srv-empty' }, 'no page servers on this machine')]));
  top.setStatus(`${(l.up || []).length} up · ${(l.down || []).length} down · ${left.length} also running`);
}

let last = null;
async function load() {
  try {
    last = await api.get('/servers');
  } catch (err) {
    top.setStatus(String(err.message || err), { tone: 'danger' });
    return;
  }
  render(last);
}

live({
  events: '/api/events',
  poll: '/api/poll',
  onEvent: (e) => { if (e.type === 'servers') load(); },
  onStatus: (s) => top.setLive(s),
  onNewBuild: api.newBuild,
});
load();
