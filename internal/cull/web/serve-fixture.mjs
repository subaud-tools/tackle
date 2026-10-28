// serve-fixture.mjs — a seeded `cull serve` for the probe.
//
// Builds the cull binary and the seeder into a temp dir, seeds a fresh
// CULL_HOME from testdata/review.json (one run, the fixture's items), starts
// `cull serve --foreground --no-open` on it and returns {url, base, token,
// projectId, runId, root, reseed, stop}.
// Each serve, and the built binaries, are stopped and removed at the probe's
// cleanup (build/probe.mjs onCleanup).
//
// BROWSER SAFETY: serve is always started with --no-open (asserted at module
// load) and --foreground, which never opens a browser.

import { execFileSync, spawn } from 'node:child_process';
import {
  mkdirSync,
  mkdtempSync,
  readFileSync,
  realpathSync,
  rmSync,
} from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { onCleanup, until } from './build/probe.mjs';

const here = dirname(fileURLToPath(import.meta.url));
const repoRoot = resolve(here, '../../..');

const SERVE_ARGS = ['serve', '--foreground', '--no-open'];
if (!SERVE_ARGS.includes('--no-open')) {
  throw new Error('serve-fixture.mjs: SERVE_ARGS must include --no-open');
}

export const ROOT = '/home/dev/shop';
export const FIXTURE = join(here, 'testdata', 'review.json');

let built = null;
function build() {
  if (built) return built;
  const dir = mkdtempSync(join(realpathSync(tmpdir()), 'cull-probe-bin-'));
  const bin = process.env.CULL_BIN || join(dir, 'cull');
  const seed = join(dir, 'seed');
  if (!process.env.CULL_BIN) {
    console.error('[serve-fixture] building cull…');
    execFileSync('go', ['build', '-o', bin, './cmd/cull'], {
      cwd: repoRoot,
      stdio: 'inherit',
    });
  }
  execFileSync('go', ['build', '-o', seed, './internal/cull/web/seed'], {
    cwd: repoRoot,
    stdio: 'inherit',
  });
  built = { dir, bin, seed };
  onCleanup(cleanupBuild);
  return built;
}

async function waitForAdvert(path, timeoutMs = 15000) {
  const adv = await until(() => {
    const a = JSON.parse(readFileSync(path, 'utf8'));
    return a.url && a.token ? a : null;
  }, timeoutMs);
  if (!adv) throw new Error(`cull serve did not start within ${timeoutMs}ms`);
  return adv;
}

/** Start a seeded serve. opts.fixture overrides the seeded review.json;
 * opts.sends records a send in each state. */
export async function startServe(opts = {}) {
  const { bin, seed } = build();
  const home = mkdtempSync(join(realpathSync(tmpdir()), 'cull-probe-home-'));
  mkdirSync(join(home, 'state'), { recursive: true });
  const fixture = opts.fixture || FIXTURE;
  const run = (extra = []) =>
    JSON.parse(
      execFileSync(seed, [...extra, home, fixture, ROOT], {
        env: { ...process.env, CULL_HOME: home },
      }).toString(),
    );
  const first = run(opts.sends ? ['-sends'] : []);

  let stderr = '';
  const proc = spawn(bin, SERVE_ARGS, {
    env: { ...process.env, CULL_HOME: home },
    stdio: 'pipe',
  });
  proc.stderr.on('data', (d) => {
    stderr += String(d);
  });
  let adv;
  try {
    adv = await waitForAdvert(join(home, 'state', 'live', 'serve.json'));
  } catch (err) {
    proc.kill();
    throw new Error(`${String(err)}\nstderr: ${stderr}`);
  }
  const serve = {
    url: adv.url,
    base: adv.base,
    token: adv.token,
    home,
    root: ROOT,
    projectId: first.project,
    runId: first.run,
    /** Record a newer run (same items; -mutate changes one item's hash, -drop leaves one out). */
    reseed(o = {}) {
      const extra = [];
      if (o.mutate) extra.push('-mutate', o.mutate);
      if (o.drop) extra.push('-drop', o.drop);
      return run(extra);
    },
    async stop() {
      if (proc.exitCode === null) {
        proc.kill('SIGTERM');
        await new Promise((r) => {
          const t = setTimeout(() => {
            proc.kill('SIGKILL');
            r();
          }, 3000);
          proc.on('exit', () => {
            clearTimeout(t);
            r();
          });
        });
      }
      rmSync(home, { recursive: true, force: true });
    },
  };
  onCleanup(() => serve.stop());
  return serve;
}

function cleanupBuild() {
  if (built) rmSync(built.dir, { recursive: true, force: true });
  built = null;
}
