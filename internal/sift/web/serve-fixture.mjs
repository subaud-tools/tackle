// serve-fixture.mjs — a seeded `sift serve` for the probe.
//
// Builds the sift binary and the seeder into a temp dir, seeds a fresh
// SIFT_HOME from testdata/round.json (one round; its files written under a
// temp base, the repos committed), starts `sift serve --foreground --no-open`
// on it and returns {url, base, token, home, round, stop}.
// Each serve, and the built binaries, are stopped and removed at the probe's
// cleanup (build/probe.mjs onCleanup).
//
// BROWSER SAFETY: serve is always started with --no-open (asserted at module
// load) and --foreground, which never opens a browser.

import { execFileSync, spawn } from 'node:child_process';
import { mkdtempSync, readFileSync, realpathSync, rmSync } from 'node:fs';
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

export const FIXTURE = join(here, 'testdata', 'round.json');

let built = null;
function build() {
  if (built) return built;
  const dir = mkdtempSync(join(realpathSync(tmpdir()), 'sift-probe-bin-'));
  const bin = process.env.SIFT_BIN || join(dir, 'sift');
  const seed = join(dir, 'seed');
  if (!process.env.SIFT_BIN) {
    console.error('[serve-fixture] building sift…');
    execFileSync('go', ['build', '-o', bin, './cmd/sift'], {
      cwd: repoRoot,
      stdio: 'inherit',
    });
  }
  execFileSync('go', ['build', '-o', seed, './internal/sift/web/seed'], {
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
  if (!adv) throw new Error(`sift serve did not start within ${timeoutMs}ms`);
  return adv;
}

export const BACKLOG = join(here, 'testdata', 'backlog.json');

/** Start a seeded serve. opts.scale pads the round to that many rows;
 * opts.recommending leaves one file without a recommendation; opts.sends
 * records a send in each state (a backlog fixture); opts.fixture seeds
 * another fixture. */
export async function startServe(opts = {}) {
  const { bin, seed } = build();
  const home = mkdtempSync(join(realpathSync(tmpdir()), 'sift-probe-home-'));
  const base = mkdtempSync(join(realpathSync(tmpdir()), 'sift-probe-base-'));
  const args = opts.scale ? ['-scale', String(opts.scale)] : [];
  if (opts.recommending) args.push('-recommending');
  if (opts.sends) args.push('-sends');
  const env = { ...process.env, SIFT_HOME: home };
  const first = JSON.parse(
    execFileSync(seed, [...args, home, opts.fixture ?? FIXTURE, base], {
      env,
    }).toString(),
  );
  let stderr = '';
  const proc = spawn(bin, SERVE_ARGS, { env, stdio: 'pipe' });
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
    dir: base,
    round: first.round,
    /** Runs the seeded sift with args and stdin, as the agent would. */
    sift(args, input = '') {
      return execFileSync(bin, args, { env, input }).toString();
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
      rmSync(base, { recursive: true, force: true });
    },
  };
  onCleanup(() => serve.stop());
  return serve;
}

function cleanupBuild() {
  if (built) rmSync(built.dir, { recursive: true, force: true });
  built = null;
}
