// probe-apply.mjs — review's and run's probe scenarios (Task 9; Court,
// 2026-10-17: decide → review → run), run by probe.mjs with its check
// helpers (t.check, t.checkList, t.until, t.eventually).
//
// Each scenario group starts its own serve (serve.mjs) with its own data:
// real clones under the serve's temp home with landed branches (seedClones:
// the casebook lane's git runs there and nowhere else), dormant repos and a
// bot's PR in the GitHub cache (seedRepos). serve's gh is the probe's fake
// (serve.mjs): nothing the probe starts reaches GitHub, and every gh call
// serve makes is checked. The agent is the fake agent (agent.mjs) over
// serve's own routes: presence, casebook_job_step, casebook_job_ask.

import { createHash } from 'node:crypto';
import { readFileSync, readdirSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { startServe, probeGit } from './serve.mjs';
import { createAgent } from './agent.mjs';

const here = dirname(fileURLToPath(import.meta.url));
const hashOf = (f) =>
  createHash('sha256').update(readFileSync(f)).digest('hex');

const SESSION = 'probe-t9-pi';
const LABEL = 'pi \u00b7 t9-apply';

const dormant = (name) => ({
  repo: `schuettc/${name}`,
  pushed_at: '2023-01-10T00:00:00Z',
  default_branch: 'main',
  prs: [],
  issues: [],
});
const T9_SEED = {
  seedClones: [{ name: 't9-local', branches: ['feat/l1', 'feat/l2'] }],
  seedRepos: [
    dormant('t9-dormant-a'),
    dormant('t9-dormant-b'),
    dormant('t9-dormant-c'),
    {
      repo: 'schuettc/t9-prs',
      pushed_at: '2026-09-20T00:00:00Z',
      default_branch: 'main',
      prs: [
        {
          repo: 'schuettc/t9-prs',
          number: 7,
          title: 'bump lodash from 4.17.20 to 4.17.21',
          author: 'dependabot[bot]',
          state: 'OPEN',
          created_at: '2026-03-01T00:00:00Z',
          updated_at: '2026-03-01T00:00:00Z',
        },
      ],
      issues: [],
    },
  ],
};
const BR = (b) => `branch:schuettc/t9-local@${b}`;
const PR = 'pr:schuettc/t9-prs#7';
const REPO_A = 'repo:schuettc/t9-dormant-a';
const REPO_B = 'repo:schuettc/t9-dormant-b';
// Decided delete: a step that can never be undone, so its group asks first.
const REPO_C = 'repo:schuettc/t9-dormant-c';

// The page's state, read in one place.
const read = {
  primary: (pg) =>
    pg.$eval('.kit-primary', (e) => (e.hidden ? '' : e.textContent.trim())),
  attached: (pg) =>
    pg.$eval('[data-testid="composer-attached"]', (e) => e.textContent),
  cards: (pg) =>
    pg.$$eval('.cb-job .cb-needs-card', (els) =>
      els.map((e) => e.dataset.kind),
    ),
  step: (pg, id) =>
    pg
      .$eval(`.cb-steps .cb-step[data-step="${id}"]`, (e) => ({
        state: e.dataset.state,
        glyph: e.querySelector('.cb-step-g').textContent,
        text: e.querySelector('.cb-step-t').textContent,
        undo: !!e.querySelector('.cb-undo'),
      }))
      .catch(() => null),
};

// evalNow finds an element and reads it in one go, in the page. Playwright's
// $eval finds it and reads it in two calls, so a live redraw between them
// (the running job redraws on every step event) hands the read a detached
// element: no computed colours, no size, no focus. The selector must be
// CSS, and the read takes no argument.
function evalNow(pg, selector, read) {
  const sel = JSON.stringify(selector);
  return pg.evaluate(`(() => {
    const e = document.querySelector(${sel});
    if (!e) throw new Error('no element matches ' + ${JSON.stringify(sel)});
    return (${read.toString()})(e);
  })()`);
}

// cssColor resolves a CSS colour expression to its computed rgb() string.
function cssColor(pg, expr) {
  return pg.evaluate((e) => {
    const d = document.createElement('div');
    d.style.color = e;
    document.body.append(d);
    const v = getComputedStyle(d).color;
    d.remove();
    return v;
  }, expr);
}

const pressWithFocus = async (pg, k) => {
  await pg.evaluate(() => {
    window.dispatchEvent(new Event('focus'));
    document.body.dispatchEvent(new FocusEvent('focus'));
  });
  await pg.keyboard.press(k);
};

// inThemes shows the page light and dark (or only the themes named),
// asserting the theme took and that nothing is focused (a focused control's
// ring is the probe's click, not the page).
async function inThemes(t, pg, name, themes = ['light', 'dark']) {
  const bgs = { light: 'rgb(244, 245, 248)', dark: 'rgb(20, 22, 29)' };
  for (const theme of themes) {
    for (let i = 0; i < 3; i++) {
      const th = await pg.evaluate(
        () => document.documentElement.dataset.theme,
      );
      if (th === theme) break;
      await pg.click('button.kit-ctl:has-text("theme")');
    }
    const got = await pg.evaluate(() => ({
      theme: document.documentElement.dataset.theme,
      bg: getComputedStyle(document.body).backgroundColor,
    }));
    await pg.evaluate(() => document.activeElement?.blur());
    const focused = await pg.evaluate(
      () => document.activeElement === document.body,
    );
    t.check(
      `the ${name} view is ${theme} (theme ${got.theme}, body ${got.bg}), nothing focused`,
      got.theme === theme && got.bg === bgs[theme] && focused,
    );
  }
  for (let i = 0; i < 3; i++) {
    const th = await pg.evaluate(() => document.documentElement.dataset.theme);
    if (th === 'light') break;
    await pg.click('button.kit-ctl:has-text("theme")');
  }
}

// copyPlan writes a second plan with the same steps as plan jobId straight
// into this serve's database: what an older casebook (before serve kept an
// item to one unfinished job) could have left behind. serve can't be asked
// for it; that is the point. It returns the copy's job id.
async function copyPlan(home, jobId) {
  const { DatabaseSync } = await import('node:sqlite');
  const file = readdirSync(home, { recursive: true }).find(
    (f) => String(f).split('/').pop() === 'casebook.db',
  );
  if (!file) throw new Error(`no casebook.db under ${home}`);
  const db = new DatabaseSync(join(home, String(file)));
  try {
    db.exec('PRAGMA busy_timeout = 5000');
    db.exec('BEGIN IMMEDIATE');
    const r = db
      .prepare(
        `INSERT INTO jobs(plan_json, machine, session, state, created_at)
         SELECT plan_json, machine, '', 'planned', created_at FROM jobs WHERE id = ?`,
      )
      .run(jobId);
    const id = Number(r.lastInsertRowid);
    db.prepare(
      `INSERT INTO steps(job_id, pos, key, action, lane, command, precondition, posts, expected_tip, state, updated_at)
       SELECT ?, pos, key, action, lane, command, precondition, posts, expected_tip, 'pending', updated_at
       FROM steps WHERE job_id = ?`,
    ).run(id, jobId);
    db.exec('COMMIT');
    return id;
  } finally {
    db.close();
  }
}

// gitIn runs git in a probe clone (under the serve's temp home only),
// hermetically (serve.mjs probeGit).
function gitIn(home, clone, args) {
  try {
    return probeGit(join(home, 'clones', clone), args, home);
  } catch {
    return null;
  }
}

let serves = [];
/** stopAll stops every serve these scenarios started (probe.mjs cleanup). */
export function stopApplyServes() {
  for (const s of serves) s.stop();
  serves = [];
}

async function ownServe(opts) {
  const s = await startServe(opts);
  serves.push(s);
  return s;
}

export async function applyScenarios(shared, t) {
  // A browser context of their own: another scenario's fake clock (a
  // context's clock is shared by its pages) must not age these pages'
  // observations, and the foot reads the real, advancing time.
  const context = await shared.browser().newContext();
  try {
    await applyScenariosIn(context, t);
  } finally {
    await context.close();
  }
}

// PROBE_SCENARIO (with PROBE_ONLY=apply) runs only the scenario groups
// whose name has it: main, stale, list, same.
const runs = (name) =>
  !process.env.PROBE_SCENARIO || name.includes(process.env.PROBE_SCENARIO);

async function applyScenariosIn(context, t) {
  const groups = [
    ['main', T9_SEED, applyScenariosOn],
    [
      'stale',
      { seedRepos: [dormant('t9-stale')], syncInterval: '2s' },
      staleScenario,
    ],
    ['list', T11_SEED, listScenario],
    ['same', { seedRepos: [dormant('t10-same')] }, sameGroupScenario],
  ];
  for (const [name, seed, run] of groups) {
    if (!runs(name)) continue;
    const s = await ownServe(seed);
    try {
      await run(context, t, s);
    } finally {
      s.stop();
      serves = serves.filter((x) => x !== s);
    }
  }
}

// A waiting group's head button, a group's address, and the open group's
// one button.
const HEAD = (id) => `.cb-apply-part[data-part="${id}"] .cb-group-plan`;
const groupHash = (id) => `#/review?group=${encodeURIComponent(id)}`;
const RUN_BTN = '[data-testid="run-group"] .kit-btn';
// A waiting group's header names it: its kind, plural, and its decision in
// serve's words ("repos Archive it on GitHub"). NAME_OF reads it in the page.
const NAME_OF = (part) =>
  `[...(document.querySelector('${part}')?.querySelectorAll('.cb-group-plan .cb-group-label, .cb-group-plan .cb-group-choice') ?? [])].map((x) => x.textContent).join(' ')`;
const nameIs = (pg, id, want) =>
  pg
    .waitForFunction(
      `${NAME_OF(`.cb-apply-part[data-part="${id}"]`)} === ${JSON.stringify(want)}`,
      null,
      { timeout: 8000 },
    )
    .then(() => true)
    .catch(() => false);
// The open review's title (it follows the ticks).
const reviewTitle = (pg) =>
  pg.evaluate(
    () =>
      document.querySelector('[data-testid="review-group"] .kit-h1')
        ?.textContent ?? '',
  );

// openedGroup waits for group id to be open in review, its plan (job) shown.
function openedGroup(t, pg, id, jobId) {
  return t.until(
    pg,
    ([hash, g, j]) =>
      location.hash === hash &&
      document.querySelector('[data-testid="review-group"]')?.dataset.group ===
        g &&
      (j === null || !!document.querySelector('[data-testid="run-group"]')),
    [groupHash(id), id, jobId ?? null],
  );
}

// The open group in review, read in one go.
const readGroup = (pg) =>
  pg.evaluate(() => {
    const doc = document.querySelector('[data-testid="review-group"]');
    const text = (s) => doc?.querySelector(s)?.textContent ?? '';
    return {
      title: text('.kit-h1'),
      what: text('.cb-approve-what'),
      warn: text('.cb-approve-warn'),
      runs: text('[data-testid="runs-through"]'),
      noAgent: text('[data-testid="no-agent"]'),
      btns: [
        ...(doc?.querySelectorAll('[data-testid="run-group"] .kit-btn') ?? []),
      ].map((b) => `${b.textContent}${b.disabled ? ' (off)' : ''}`),
      selects: doc?.querySelectorAll('select').length ?? 0,
      keys: [
        ...new Set(
          [...(doc?.querySelectorAll('.cb-plan-step') ?? [])].map(
            (e) => e.dataset.key,
          ),
        ),
      ].sort(),
      cmds: [...(doc?.querySelectorAll('.cb-plan-step .cb-cmd') ?? [])].map(
        (e) => e.textContent,
      ),
      undo: [...(doc?.querySelectorAll('[data-testid="plan-undo"]') ?? [])].map(
        (e) => e.textContent.trim(),
      ),
      where: [
        ...(doc?.querySelectorAll('[data-testid="plan-where"]') ?? []),
      ].map((e) => e.textContent),
      primary: (() => {
        const p = document.querySelector('.kit-primary');
        return p && !p.hidden ? p.textContent : '';
      })(),
    };
  });

// The words the old apply used: none is left on review or run.
const OLD_WORDS = /plans to approve|approve and run|\bdiscard\b|to apply\b/i;
const shownText = (pg) =>
  pg.evaluate(
    () =>
      `${document.querySelector('.cb-apply-list')?.innerText ?? ''}\n${document.querySelector('.cb-apply-read')?.innerText ?? ''}`,
  );

// sameGroupScenario: two tabs open the same group: each previews it and
// serve stores nothing (a plan is stored only when Court runs it). Three
// tabs share one live connection (web/shared-stream.ts).
async function sameGroupScenario(context, t, serveHandle) {
  const { check, checkList, until } = t;
  console.log('\nscenario: review — two tabs open the same group: one plan');
  const agent = createAgent(serveHandle.base, serveHandle.token);
  const ITEM = 'repo:schuettc/t10-same';
  await agent.api('POST', '/api/decide', {
    keys: [ITEM],
    disposition: 'archive',
  });
  const events = [];
  const onEvents = (r) => {
    if (new URL(r.url()).pathname === '/api/events') events.push(r.url());
  };
  context.on('request', onEvents);
  const pg = await context.newPage();
  const pg2 = await context.newPage();
  const pg3 = await context.newPage();
  const errors = [];
  pg.on('pageerror', (e) => errors.push(String(e)));
  pg2.on('pageerror', (e) => errors.push(`tab 2: ${e}`));
  pg3.on('pageerror', (e) => errors.push(`tab 3: ${e}`));
  const planOf = async (p) => {
    await p.setViewportSize({ width: 1600, height: 900 });
    await p.goto(serveHandle.url + '#/review', {
      waitUntil: 'domcontentloaded',
      timeout: 15000,
    });
    await nameIs(p, 'archive|repo', 'repos Archive it on GitHub');
    const resp = p.waitForResponse(
      (r) =>
        r.url().includes('/api/apply/plan') && r.request().method() === 'POST',
    );
    await p.click(HEAD('archive|repo'));
    const r = await resp;
    return { status: r.status(), body: await r.json().catch(() => ({})) };
  };
  try {
    const a = await planOf(pg);
    const b = await planOf(pg2);
    const stepsOf = (v) =>
      (v.job?.steps ?? []).map((x) => `${x.key} ${x.command}`).join('\n');
    check(
      `the first tab's group is serve's preview (${a.status}, ${(a.body.job?.steps ?? []).length} step)`,
      a.status === 200 && (a.body.job?.steps ?? []).length === 1,
    );
    check(
      `the second tab's preview has the same steps (${b.status})`,
      b.status === 200 && stepsOf(b.body) === stepsOf(a.body),
    );
    const jobs = (await agent.api('GET', '/api/jobs')).jobs ?? [];
    checkList(
      'previewing stores no plan: serve has no job',
      jobs.map((j) => `#${j.id} ${j.state}`),
      [],
    );
    check(
      'both tabs show it in review with its one button',
      (await openedGroup(t, pg, 'archive|repo', 0)) &&
        (await openedGroup(t, pg2, 'archive|repo', 0)),
    );
    await pg3.setViewportSize({ width: 1600, height: 900 });
    await pg3.goto(serveHandle.url + '#/decide', {
      waitUntil: 'domcontentloaded',
      timeout: 15000,
    });
    const live = async (p) =>
      until(
        p,
        () => document.querySelector('.kit-live')?.dataset.state === 'live',
      );
    const allLive = (await live(pg)) && (await live(pg2)) && (await live(pg3));
    // Kept: a window for a second connection to go out.
    await pg.waitForTimeout(1000);
    check(
      `three tabs share one live connection: one GET /api/events in the browser (${events.length}), every tab live (${allLive})`,
      allLive && events.length === 1,
    );
    check(`no page errors (${errors.join(' | ')})`, errors.length === 0);
  } finally {
    context.off('request', onEvents);
    await pg.close();
    await pg2.close();
    await pg3.close();
  }
}

// listScenario: review lists the decided items waiting, grouped by what
// will happen, each group collapsed under its button; run lists the jobs,
// running on top with a progress bar that moves as the agent reports, then
// Done, collapsed. Nothing runs from either list. A decision made in
// another tab raises a group's count live, and the waiting strip's "review"
// opens the group it joined.
const T11 = 't11-local';
const T11_SEED = {
  seedClones: [{ name: T11, branches: ['feat/a', 'feat/b', 'feat/c'] }],
  seedRepos: [
    dormant('t11-x'),
    dormant('t11-y'),
    // Already archived on GitHub: decided archive, it is done with nothing
    // run.
    { ...dormant('t11-gone'), archived: true },
    {
      repo: 'schuettc/t11-prs',
      pushed_at: '2026-09-20T00:00:00Z',
      default_branch: 'main',
      prs: [3, 4].map((number) => ({
        repo: 'schuettc/t11-prs',
        number,
        title: `bump dep ${number}`,
        author: 'dependabot[bot]',
        state: 'OPEN',
        created_at: '2026-03-01T00:00:00Z',
        updated_at: '2026-03-01T00:00:00Z',
      })),
      issues: [],
    },
  ],
};
async function listScenario(context, t, serveHandle) {
  const { check, checkList, until, eventually } = t;
  console.log(
    '\nscenario: review and run — the waiting groups; running, then Done',
  );
  const agent = createAgent(serveHandle.base, serveHandle.token);
  const job = (id) => agent.api('GET', `/api/job?id=${id}`);
  const B = (b) => `branch:schuettc/${T11}@${b}`;
  const X = 'repo:schuettc/t11-x';
  const Y = 'repo:schuettc/t11-y';
  const S = 'probe-t11-pi';
  const here = () =>
    agent.presence(S, 'pi \u00b7 t11', '/home/court/t11', 'pi');
  await here();
  const keep = setInterval(() => void here(), 1000);
  const pg = await context.newPage();
  const pg2 = await context.newPage();
  const errors = [];
  pg.on('pageerror', (e) => errors.push(String(e)));
  pg2.on('pageerror', (e) => errors.push(`tab 2: ${e}`));
  try {
    // A finished job: feat/c deleted (the casebook lane runs it at once).
    await agent.api('POST', '/api/decide', {
      keys: [B('feat/c')],
      disposition: 'delete',
    });
    const doneJob = await agent.api('POST', '/api/apply/plan', {
      keys: [B('feat/c')],
    });
    await agent.api('POST', '/api/apply/approve', {
      plan_id: doneJob.job.id,
      session: '',
    });
    await eventually(
      async () => (await job(doneJob.job.id)).job.state === 'done',
      10000,
    );
    // A running job: both PRs closed, waiting on the agent.
    const PRS = ['pr:schuettc/t11-prs#3', 'pr:schuettc/t11-prs#4'];
    await agent.api('POST', '/api/decide', { keys: PRS, disposition: 'close' });
    const runJob = await agent.api('POST', '/api/apply/plan', { keys: PRS });
    await agent.api('POST', '/api/apply/approve', {
      plan_id: runJob.job.id,
      session: S,
    });
    // Waiting: the other two branches. (feat/c, verified done, has left
    // review at once.)
    await agent.api('POST', '/api/decide', {
      keys: [B('feat/a'), B('feat/b')],
      disposition: 'delete',
    });

    await pg.setViewportSize({ width: 1600, height: 900 });
    await pg.goto(serveHandle.url + '#/review', {
      waitUntil: 'domcontentloaded',
      timeout: 15000,
    });
    const parts = () =>
      pg.$$eval('.cb-apply-list .cb-apply-part', (els) =>
        els.map(
          (e) =>
            `${e.dataset.part} ${[...e.querySelectorAll('.cb-group-label, .cb-group-choice')].map((x) => x.textContent).join(' ')} ${e.querySelector('.cb-group-n').textContent}${e.classList.contains('collapsed') ? ' collapsed' : ''}`,
        ),
      );
    const want = ['delete|branch branches Delete it 2 collapsed'];
    await until(
      pg,
      (n) => document.querySelectorAll('.cb-apply-part').length === n,
      want.length,
    );
    checkList(
      'review lists only the waiting groups, each collapsed under its name (kind and decision); the verified feat/c is gone from it',
      await parts(),
      want,
    );
    const top = await pg.evaluate(() => ({
      chips: document.querySelectorAll('.cb-apply-list .kit-chip').length,
      jobs: document.querySelectorAll('.cb-apply-list .kit-row[data-job]')
        .length,
      approve: [...document.querySelectorAll('.cb-apply-list button')]
        .map((b) => b.textContent)
        .filter((t) => /approve|run/i.test(t)),
      primary: (() => {
        const p = document.querySelector('.kit-primary');
        return p && !p.hidden ? p.textContent : '';
      })(),
      on: document.querySelector('.kit-ctl.on')?.dataset.id,
    }));
    check(
      `review is the bar's place, no view chips, no job rows (${top.on}; ${top.chips} chips, ${top.jobs} jobs)`,
      top.on === 'review' && top.chips === 0 && top.jobs === 0,
    );
    check(
      `nothing on review says the old words (plans to approve, approve and run, discard, to apply) or "ready"`,
      !OLD_WORDS.test(await shownText(pg)) &&
        !/\bready\b/i.test(await shownText(pg)),
    );
    check(
      `nothing runs from the list: no approve or run there, no primary (${top.approve.join(', ') || 'none'}; "${top.primary}")`,
      top.approve.length === 0 && top.primary === '',
    );
    // Ticking an item opens its group's review, for the ticked.
    await pg.click('.cb-apply-part[data-part="delete|branch"] .cb-group-caret');
    await until(
      pg,
      (k) =>
        !!document.querySelector(`.cb-apply-list .kit-row[data-key="${k}"]`),
      B('feat/a'),
    );
    await pg.click(
      `.cb-apply-list .kit-row[data-key="${B('feat/a')}"] .kit-box`,
    );
    const reviewKeys = () =>
      pg.evaluate(() =>
        [
          ...new Set(
            [
              ...document.querySelectorAll(
                '[data-testid="review-group"] .cb-plan-step',
              ),
            ].map((e) => e.dataset.key),
          ),
        ]
          .sort()
          .join(','),
      );
    check(
      'ticking one item opens its review for it alone, the group box mixed',
      (await openedGroup(t, pg, 'delete|branch', 0)) &&
        (await eventually(async () => (await reviewKeys()) === B('feat/a'))) &&
        (await until(
          pg,
          () =>
            document
              .querySelector(
                '.cb-apply-part[data-part="delete|branch"] .cb-group-box',
              )
              ?.getAttribute('aria-checked') === 'mixed',
        )),
    );
    await pg.click(
      `.cb-apply-list .kit-row[data-key="${B('feat/a')}"] .kit-box`,
    );
    check(
      'unticked, the review follows: every waiting item of the group',
      await eventually(
        async () =>
          (await reviewKeys()) === [B('feat/a'), B('feat/b')].sort().join(','),
      ),
    );
    await pg.evaluate(() => {
      location.hash = '#/review';
    });
    await until(
      pg,
      () => !!document.querySelector('[data-testid="review-overview"]'),
    );
    await inThemes(t, pg, 'review list');

    // The group's button opens it in the middle pane: each item's command
    // and undo, where it runs, its one button.
    const planResp = pg.waitForResponse(
      (r) =>
        r.url().includes('/api/apply/plan') && r.request().method() === 'POST',
    );
    await pg.click(HEAD('delete|branch'));
    const planBody = await (await planResp).json();
    const planId = planBody.job?.id;
    const opened = await openedGroup(t, pg, 'delete|branch', planId);
    const plan = await readGroup(pg);
    checkList(
      "the group's name previews its items",
      plan.keys,
      [B('feat/a'), B('feat/b')].sort(),
    );
    check(
      `the group opens in the middle pane with each item's exact command, serve's (${plan.cmds.length})`,
      opened &&
        plan.cmds.length > 0 &&
        plan.cmds.join('\n') ===
          planBody.groups
            .flatMap((g) => g.steps.map((s) => s.command))
            .join('\n'),
    );
    checkList(
      "each part says where it runs and whether it can be undone, in serve's words",
      plan.undo,
      planBody.groups.map((g) => g.undo),
    );
    check(
      `a local part runs here and can be undone (${plan.where.join(' | ')})`,
      plan.where.length > 0 &&
        plan.where.every(
          (w) => w.startsWith('Runs here, on ') && w.includes('Can be undone'),
        ),
    );
    checkList(
      "the group's one button, named for the action as its title is",
      plan.btns,
      [plan.title],
    );
    check(
      `its title has its counts ("${plan.title}")`,
      plan.title === 'Delete 2 local branches, delete 2 remote branches',
    );
    const jobCount = async () =>
      ((await agent.api('GET', '/api/jobs')).jobs ?? []).length;
    check(
      `opening it runs nothing and stores nothing (serve has its 2 jobs, the done and the running), and the group stays in review (${await jobCount()} jobs)`,
      (await jobCount()) === 2 &&
        (await until(
          pg,
          (b) => !!document.querySelector(b),
          HEAD('delete|branch'),
        )),
    );
    await pg.evaluate(() => {
      location.hash = '#/review';
    });
    check(
      'leaving the group runs nothing (no job stored) and review shows its overview',
      (await until(
        pg,
        () => !!document.querySelector('[data-testid="review-overview"]'),
      )) && (await jobCount()) === 2,
    );

    // run: the jobs, running on top, Done collapsed.
    await pg.click('.kit-ctl[data-id="run"]');
    const runWant = ['running Running 1', 'done Done 1 collapsed'];
    await until(
      pg,
      (n) =>
        location.hash === '#/run' &&
        document.querySelectorAll('.cb-apply-part').length === n,
      runWant.length,
    );
    checkList(
      'run lists the running job on top, then Done, collapsed',
      await parts(),
      runWant,
    );
    check(
      'the bar counts the jobs running ("run 1")',
      await until(
        pg,
        () =>
          document.querySelector('.kit-ctl[data-id="run"] .kit-n')
            ?.textContent === '1',
      ),
    );
    check(
      `nothing on run says the old words, and the plan just opened is not listed`,
      !OLD_WORDS.test(await shownText(pg)) &&
        !(await pg.$(`.cb-apply-list .kit-row[data-job="${planId}"]`)),
    );
    const bar = () =>
      pg.evaluate(
        (id) =>
          document
            .querySelector(`.cb-apply-list .kit-row[data-job="${id}"] .cb-pb`)
            ?.getAttribute('aria-valuenow') ?? '',
        runJob.job.id,
      );
    const bar0 = await bar();
    const step0 = runJob.job.steps[0];
    for (const state of ['started', 'reported'])
      await agent.api('POST', '/api/agent/job-step', {
        session: S,
        job: runJob.job.id,
        step: step0.id,
        state,
      });
    check(
      `the running job's progress bar rises live as a step is reported (aria-valuenow ${bar0} → 1)`,
      bar0 === '0' &&
        (await until(
          pg,
          (id) =>
            document
              .querySelector(`.cb-apply-list .kit-row[data-job="${id}"] .cb-pb`)
              ?.getAttribute('aria-valuenow') === '1',
          runJob.job.id,
        )),
    );
    await inThemes(t, pg, 'run list');

    // A decision in a second tab raises a group's count without a reload,
    // and its strip's review opens that group.
    await pg.evaluate(() => {
      location.hash = '#/review';
      window.__sameLoad = true;
    });
    await pg2.setViewportSize({ width: 1600, height: 900 });
    await pg2.goto(`${serveHandle.url}#/item/${encodeURIComponent(X)}`, {
      waitUntil: 'domcontentloaded',
      timeout: 15000,
    });
    await pg2.click('.cb-decide-row .cb-choice[data-d="archive"]');
    check(
      'a decision in a second tab shows as a group here, live ("repos Archive it on GitHub", 1)',
      (await nameIs(pg, 'archive|repo', 'repos Archive it on GitHub')) &&
        (await until(
          pg,
          () =>
            window.__sameLoad === true &&
            document.querySelector(
              '.cb-apply-part[data-part="archive|repo"] .cb-group-n',
            )?.textContent === '1',
        )),
    );
    await agent.api('POST', '/api/decide', {
      keys: [Y],
      disposition: 'archive',
    });
    check(
      "another raises the group's count, still without a reload",
      await until(
        pg,
        () =>
          window.__sameLoad === true &&
          document.querySelector(
            '.cb-apply-part[data-part="archive|repo"] .cb-group-n',
          )?.textContent === '2',
      ),
    );
    const strip = await pg2
      .locator('.cb-strip a')
      .evaluate((a) => ({ href: a.getAttribute('href'), text: a.textContent }))
      .catch(() => null);
    await pg2.click('.cb-strip a').catch(() => {});
    check(
      `the strip's "review" opens the group it joined in review (${JSON.stringify(strip)})`,
      strip?.href === groupHash('archive|repo') &&
        strip?.text === 'review' &&
        (await openedGroup(t, pg2, 'archive|repo', null)) &&
        (await until(
          pg2,
          () =>
            document.querySelector('.cb-apply-part.cb-focus')?.dataset.part ===
            'archive|repo',
        )),
    );
    // Done opens when asked.
    await pg.evaluate(() => {
      location.hash = '#/run';
    });
    await until(
      pg,
      () => !!document.querySelector('.cb-apply-part[data-part="done"]'),
    );
    await pg.click('.cb-apply-part[data-part="done"] .cb-group-caret');
    check(
      'Done opens to its finished job',
      await until(
        pg,
        (id) =>
          !!document.querySelector(
            `.cb-apply-list .kit-row[data-job="${id}"]`,
          ) &&
          !document
            .querySelector('.cb-apply-part[data-part="done"]')
            ?.classList.contains('collapsed'),
        doneJob.job.id,
      ),
    );
    // An item found already so, with nothing run, says so on its page.
    const GONE = 'repo:schuettc/t11-gone';
    await agent.api('POST', '/api/decide', {
      keys: [GONE],
      disposition: 'archive',
    });
    await pg.evaluate((k) => {
      location.hash = `#/item/${encodeURIComponent(k)}`;
    }, GONE);
    const factsLine = async () =>
      pg.evaluate(
        () => document.querySelector('.cb-item .kit-facts')?.textContent ?? '',
      );
    const foundSo = await eventually(async () =>
      (await factsLine()).includes('done: found so, nothing ran'),
    );
    check(
      `a done item with no job says how: "done: found so, nothing ran" in its facts ("${await factsLine()}")`,
      foundSo,
    );
    check(`no page errors (${errors.join(' | ')})`, errors.length === 0);
  } finally {
    clearInterval(keep);
    await pg.close();
    await pg2.close();
  }
}

async function applyScenariosOn(context, t, serveHandle) {
  const { check, checkList, until, eventually } = t;
  const agent = createAgent(serveHandle.base, serveHandle.token);
  const job = (id) => agent.api('GET', `/api/job?id=${id}`);
  const jobsNow = async () => (await agent.api('GET', '/api/jobs')).jobs ?? [];
  const stepFor = (jv, key, action) =>
    jv.job.steps.find((s) => s.key === key && (!action || s.action === action));
  const home = serveHandle.home;

  // Court's decisions: the landed branches deleted, two dormant repos
  // archived and one deleted, the bot's PR closed with a comment.
  await agent.api('POST', '/api/decide', {
    keys: [BR('feat/l1'), BR('feat/l2')],
    disposition: 'delete',
  });
  await agent.api('POST', '/api/decide', {
    keys: [REPO_A, REPO_B],
    disposition: 'archive',
  });
  await agent.api('POST', '/api/decide', {
    keys: [REPO_C],
    disposition: 'delete',
  });
  await agent.api('POST', '/api/decide', {
    keys: [PR],
    disposition: 'close',
    note: 'superseded by #9',
  });

  // openGroup opens group id from its head in review; its plan's reply.
  const openGroup = async (pg, id) => {
    await pg.evaluate(() => {
      if (!location.hash.startsWith('#/review')) location.hash = '#/review';
    });
    await until(pg, (b) => !!document.querySelector(b), HEAD(id));
    const resp = pg.waitForResponse(
      (r) =>
        r.url().includes('/api/apply/plan') && r.request().method() === 'POST',
    );
    await pg.click(HEAD(id));
    const body = await (await resp).json();
    await openedGroup(t, pg, id, body.job?.id);
    return body;
  };
  // runOpen presses the open group's button (or "a"): the page stores the
  // plan (POST /api/apply/plan, no preview) and approves it; serve's job.
  const runOpen = async (pg, how = 'click', press = true) => {
    const store = pg.waitForResponse(
      (r) =>
        r.url().includes('/api/apply/plan') &&
        r.request().method() === 'POST' &&
        !r.request().postDataJSON()?.preview,
    );
    const resp = pg.waitForResponse(
      (r) =>
        r.url().includes('/api/apply/approve') &&
        r.request().method() === 'POST',
    );
    if (press) {
      if (how === 'key') await pressWithFocus(pg, 'a');
      else await pg.click(`${RUN_BTN}.fill`);
    }
    const st = await store;
    const r = await resp;
    return {
      storeSent: st.request().postDataJSON(),
      stored: await st.json(),
      sent: r.request().postDataJSON(),
      body: await r.json(),
    };
  };
  const inRun = (pg, id) =>
    until(
      pg,
      (j) =>
        location.hash === `#/run/${j}` &&
        document.querySelector('.cb-job')?.dataset.job === String(j) &&
        document.querySelector('.kit-ctl.on')?.dataset.id === 'run',
      id,
      8000,
    );

  const pg = await context.newPage();
  const errors = [];
  pg.on('pageerror', (e) => errors.push(String(e)));
  let keep = null;
  try {
    await pg.setViewportSize({ width: 1600, height: 900 });
    await pg.goto(serveHandle.url + '#/review', {
      waitUntil: 'domcontentloaded',
      timeout: 15000,
    });

    // ---- the waiting groups ----------------------------------------------
    console.log('\nscenario: review — the waiting groups');
    await until(pg, () => !!document.querySelector('.kit-ctl[data-id="run"]'));
    const places = await pg.$$eval('.kit-ctl[data-id]', (els) =>
      els
        .map((e) => e.dataset.id)
        .filter((id) => ['decide', 'review', 'run', 'apply'].includes(id)),
    );
    checkList('the bar names three steps: decide, review, run', places, [
      'decide',
      'review',
      'run',
    ]);
    check(
      'the bar counts the decisions waiting in review ("review 6") and no job running (run shows no number)',
      await until(pg, () => {
        const n = (id) =>
          document.querySelector(`.kit-ctl[data-id="${id}"] .kit-n`)
            ?.textContent ?? '';
        return n('review') === '6' && ['', '0'].includes(n('run'));
      }),
    );
    const served = await agent.api('GET', '/api/items?view=to-apply');
    const LABELS = {
      'archive|repo': 'repos Archive it on GitHub',
      'close|pr': 'PRs Close it without merging',
      'delete|branch': 'branches Delete it',
      'delete|repo': 'repos Delete it',
    };
    await until(
      pg,
      (n) => document.querySelectorAll('.cb-apply-part').length === n,
      served.groups.length,
    );
    checkList(
      "the groups are serve's, in its order, each collapsed under its name: its kind and its decision in serve's words",
      await pg.$$eval('.cb-apply-part', (els) =>
        els.map(
          (e) =>
            `${e.dataset.part} ${[...e.querySelectorAll('.cb-group-plan .cb-group-label, .cb-group-plan .cb-group-choice')].map((x) => x.textContent).join(' ')}${e.classList.contains('collapsed') ? ' collapsed' : ''}`,
        ),
      ),
      served.groups.map((g) => `${g.id} ${LABELS[g.id] ?? g.id} collapsed`),
    );
    check(
      'nothing runs from the list: no primary',
      (await read.primary(pg)) === '',
    );
    check(
      'nothing on review says the old words (plans to approve, approve and run, discard)',
      !OLD_WORDS.test(await shownText(pg)),
    );
    for (const g of served.groups)
      await pg.click(`.cb-apply-part[data-part="${g.id}"] .cb-group-caret`);
    await until(
      pg,
      (n) => document.querySelectorAll('.cb-apply-list .kit-row').length === n,
      served.items.length,
    );
    checkList(
      "opened, the list holds serve's to-apply items, group by group in its order",
      await pg.$$eval('.cb-apply-list .kit-row', (els) =>
        els.map((e) => e.dataset.key),
      ),
      served.items.map((i) => i.key),
    );
    // Ticking an item opens its group's review, and the review follows
    // the ticks; it only previews (serve stores no plan).
    await pg.click(`.cb-apply-list .kit-row[data-key="${REPO_B}"] .kit-box`);
    const titled = (want) =>
      until(
        pg,
        (w) =>
          document.querySelector('[data-testid="review-group"]')?.dataset
            .group === 'archive|repo' &&
          document.querySelector('[data-testid="review-group"] .kit-h1')
            ?.textContent === w &&
          !!document.querySelector('[data-testid="run-group"]'),
        want,
      );
    check(
      `ticking an item opens its group's review, for it ("${await reviewTitle(pg)}")`,
      await titled('Archive 1 repo'),
    );
    await pg.click(`.cb-apply-list .kit-row[data-key="${REPO_A}"] .kit-box`);
    check(
      'the review follows the ticks: another ticked, "Archive 2 repos"',
      await titled('Archive 2 repos'),
    );
    await pg.click(`.cb-apply-list .kit-row[data-key="${REPO_A}"] .kit-box`);
    check(
      'one unticked, "Archive 1 repo" again',
      await titled('Archive 1 repo'),
    );
    check(
      `the review only previews: serve stores no plan (${(await jobsNow()).length} jobs)`,
      (await jobsNow()).length === 0,
    );
    await pg.click(`.cb-apply-list .kit-row[data-key="${PR}"] .kit-box`);
    check(
      "ticked, a group's box says so (PR's checked, archive's mixed)",
      await until(
        pg,
        () =>
          document
            .querySelector(
              '.cb-apply-part[data-part="archive|repo"] .cb-group-box',
            )
            ?.getAttribute('aria-checked') === 'mixed' &&
          document
            .querySelector('.cb-apply-part[data-part="close|pr"] .cb-group-box')
            ?.getAttribute('aria-checked') === 'true',
      ),
    );
    const signalC = await cssColor(pg, 'var(--kit-signal)');
    const boxC = await pg.$eval(
      '.cb-apply-part[data-part="close|pr"] .cb-group-box',
      (e) => getComputedStyle(e).backgroundColor,
    );
    check(
      `a ticked group's box is in signal, the selection's colour (${boxC})`,
      boxC === signalC,
    );
    await pg.click(`.cb-apply-list .kit-row[data-key="${REPO_B}"] .kit-box`);
    await pg.click(`.cb-apply-list .kit-row[data-key="${PR}"] .kit-box`);
    await pg.evaluate(() => {
      location.hash = '#/review';
    });
    await until(
      pg,
      () => !!document.querySelector('[data-testid="review-overview"]'),
    );

    // ---- a group opens in the middle pane -------------------------------
    console.log('\nscenario: review — a group opens in the middle pane');
    const planPosts = [];
    const onPlanPost = (r) => {
      if (r.method() === 'POST' && r.url().includes('/api/apply/plan'))
        planPosts.push(r.url());
    };
    pg.on('request', onPlanPost);
    const branchResp = pg.waitForResponse((r) =>
      r.url().includes('/api/apply/plan'),
    );
    await pg.click(HEAD('delete|branch'), { clickCount: 2 });
    const branchPlan = await (await branchResp).json();
    // Kept: the check is that the second press posts nothing; a window for
    // a second POST to go out.
    await pg.waitForTimeout(500);
    pg.off('request', onPlanPost);
    check(
      `a double press of a group's name asks serve once, a preview that stores nothing (${planPosts.length} POST /api/apply/plan)`,
      planPosts.length === 1 && (await jobsNow()).length === 0,
    );
    checkList(
      "it plans the group's items",
      [...new Set(branchPlan.job.steps.map((s) => s.key))].sort(),
      [BR('feat/l1'), BR('feat/l2')],
    );
    check(
      `the group opens as ${groupHash('delete|branch')}, in the middle pane`,
      await openedGroup(t, pg, 'delete|branch', branchPlan.job.id),
    );
    const local = await readGroup(pg);
    check(
      `its title has its counts and its one button is the same words, none in the bar ("${local.title}"; ${local.btns.join(', ')}; primary "${local.primary}")`,
      local.title === 'Delete 2 local branches, delete 2 remote branches' &&
        local.btns.join(',') === local.title &&
        local.primary === '',
    );
    check(
      `it says what it does, all local, nothing red, and names no runner (an all-local group needs no agent) ("${local.what}")`,
      local.what === '4 steps \u00b7 4 local, all local.' &&
        local.warn === '' &&
        local.runs === '' &&
        local.noAgent === '',
    );
    check(
      `no menu, no discard, no "approve and run" (${local.selects} menus)`,
      local.selects === 0 && !OLD_WORDS.test(await shownText(pg)),
    );
    checkList(
      "every step's command is serve's, verbatim, grouped as serve groups them",
      local.cmds,
      branchPlan.groups.flatMap((g) => g.steps.map((s) => s.command)),
    );
    checkList(
      "each part says whether it can be undone, in serve's words",
      local.undo,
      branchPlan.groups.map((g) => g.undo),
    );
    check(
      'the open group is marked in the list',
      await until(
        pg,
        () =>
          document.querySelector('.cb-apply-part.cb-focus')?.dataset.part ===
          'delete|branch',
      ),
    );
    await inThemes(t, pg, 'review group');

    console.log(
      '\nscenario: review — leaving a group runs nothing; opening it again reuses its plan',
    );
    await pg.evaluate(() => {
      location.hash = '#/review';
    });
    const stepsOf = (v) =>
      (v.job?.steps ?? []).map((x) => `${x.key} ${x.command}`).join('\n');
    check(
      'leaving it runs nothing: serve stores nothing',
      (await until(
        pg,
        () => !!document.querySelector('[data-testid="review-overview"]'),
      )) && (await jobsNow()).length === 0,
    );
    const again = await openGroup(pg, 'delete|branch');
    check(
      `opening it again previews the same steps, and stores nothing (${(await jobsNow()).length} jobs)`,
      stepsOf(again) === stepsOf(branchPlan) && (await jobsNow()).length === 0,
    );

    // ---- outward steps with no agent ----------------------------------
    console.log('\nscenario: review — outward steps with no agent attached');
    await pg.evaluate(() => {
      location.hash = '#/review';
    });
    await openGroup(pg, 'archive|repo');
    const noAgent = await readGroup(pg);
    check(
      `with no agent attached it says so, its button is off, and there is no primary ("${noAgent.noAgent}"; ${noAgent.btns.join(', ')}; "${noAgent.primary}")`,
      noAgent.noAgent === 'Open casebook from an agent to run this.' &&
        noAgent.btns.join(',') === 'Archive 2 repos (off)' &&
        noAgent.primary === '',
    );
    await pressWithFocus(pg, 'a');
    await pg.waitForTimeout(400); // kept: "a" must run nothing; a window for it
    check(
      '"a" with no agent runs nothing (serve stores no plan)',
      (await jobsNow()).length === 0,
    );
    await inThemes(t, pg, 'review, no agent', ['light']);

    // ---- a local group runs ------------------------------------------
    console.log('\nscenario: review — a group runs, and opens in run');
    await pg.evaluate(() => {
      location.hash = '#/review';
    });
    const branch = await openGroup(pg, 'delete|branch');
    // Two stored plans holding the group, as an older casebook could have
    // left them: storing it to run cancels both and stores one fresh.
    const old = await agent.api('POST', '/api/apply/plan', {
      keys: [BR('feat/l1'), BR('feat/l2')],
    });
    const twin = await copyPlan(home, old.job.id);
    await pg.evaluate(() => {
      window.__ranNoReload = true;
    });
    const ran = await runOpen(pg);
    const branchJob = ran.body.job?.id;
    check(
      `pressing it stores the plan for the reviewed items (no preview) and runs it at once (nothing permanent): serve approves plan #${ran.stored.job?.id} with no session (${JSON.stringify(ran.storeSent)} ${JSON.stringify(ran.sent)})`,
      !ran.storeSent.preview &&
        stepsOf(ran.stored) === stepsOf(branch) &&
        branchJob === ran.stored.job?.id &&
        ran.sent.plan_id === branchJob &&
        ran.sent.session === '',
    );
    const states = Object.fromEntries(
      (await jobsNow()).map((j) => [j.id, j.state]),
    );
    check(
      `two older plans holding the group are cancelled, and the one it runs is fresh (#${old.job.id} ${states[old.job.id]}, #${twin} ${states[twin]}, #${branchJob} ${states[branchJob]})`,
      states[old.job.id] === 'cancelled' &&
        states[twin] === 'cancelled' &&
        branchJob !== old.job.id &&
        branchJob !== twin &&
        !['planned', 'cancelled'].includes(states[branchJob]),
    );
    check(
      `the job opens in run (#/run/${branchJob}), run is the bar's place`,
      await inRun(pg, branchJob),
    );
    let jv = await job(branchJob);
    const l1 = stepFor(jv, BR('feat/l1'), 'branch-delete-local');
    const l1r = stepFor(jv, BR('feat/l1'), 'branch-delete-remote');
    check(
      'the local steps verify live: ✓ "deleted · verified" in the steps table',
      await until(
        pg,
        (id) =>
          document.querySelector(`.cb-step[data-step="${id}"]`)?.dataset
            .state === 'verified',
        l1.id,
        10000,
      ),
    );
    const l1Row = await read.step(pg, l1.id);
    check(
      `the verified local step reads ✓ deleted · verified (${l1Row?.glyph} ${l1Row?.text})`,
      l1Row?.glyph === '\u2713' && l1Row?.text === 'deleted \u00b7 verified',
    );
    const tones = await pg
      .waitForFunction(
        (id) => {
          const e = document.querySelector(`.cb-step[data-step="${id}"]`);
          const g = e?.querySelector('.cb-step-g');
          const tx = e?.querySelector('.cb-step-t');
          if (!g?.isConnected || !tx?.isConnected) return null;
          return {
            glyph: getComputedStyle(g).color,
            text: getComputedStyle(tx).color,
          };
        },
        l1.id,
        { timeout: 5000 },
      )
      .then((h) => h.jsonValue())
      .catch(() => ({ glyph: '', text: '' }));
    const okC = await cssColor(pg, 'var(--kit-ok)');
    const mutedC = await cssColor(pg, 'var(--kit-muted)');
    check(
      `only the glyph carries the step's tone: ✓ in ok (${tones.glyph}), its text muted (${tones.text})`,
      tones.glyph === okC && tones.text === mutedC && okC !== mutedC,
    );
    check(
      'the casebook lane deleted feat/l1 in the probe clone (git, in the temp home)',
      gitIn(home, 't9-local', 'rev-parse --verify refs/heads/feat/l1') ===
        null &&
        gitIn(home, 't9-local', 'rev-parse --verify refs/heads/main') !== null,
    );
    check(
      'verified, the branches leave review at once: the bar counts 4 waiting, not 6, with no reload',
      await until(
        pg,
        () =>
          window.__ranNoReload === true &&
          document.querySelector('.kit-ctl[data-id="review"] .kit-n')
            ?.textContent === '4',
        undefined,
        10000,
      ),
    );
    // The section fills the list and reading columns.
    const geo = await pg.evaluate(() => {
      const r = (sel) => {
        const b = document.querySelector(sel).getBoundingClientRect();
        return { x: b.x, y: b.y, w: b.width, h: b.height };
      };
      return {
        list: r('.cb-apply-list'),
        read: r('.cb-apply-read'),
        doc: r('.cb-apply-read .kit-doc'),
        rail: r('.kit-rail'),
        bar: r('.kit-bar'),
        others: [...document.querySelectorAll('.kit-app > .kit-list')]
          .filter((e) => !e.classList.contains('cb-apply-list'))
          .map((e) => getComputedStyle(e).display),
      };
    });
    check(
      `the list is the 400 px list column, bar to bottom; the job the reading column; the document 640 px, centred (${JSON.stringify(geo.list)} ${JSON.stringify(geo.read)} ${geo.doc.w})`,
      geo.list.x === 0 &&
        geo.list.w === 400 &&
        geo.list.y === geo.bar.h &&
        geo.list.y + geo.list.h === 900 &&
        geo.read.x === 400 &&
        geo.read.x + geo.read.w === geo.rail.x &&
        geo.doc.w === 640 &&
        Math.abs(geo.doc.x - (geo.read.x + (geo.read.w - 640) / 2)) <= 1,
    );
    check(
      `decide's list is not displayed (${geo.others.join(',')})`,
      geo.others.length === 1 && geo.others.every((d) => d === 'none'),
    );

    // ---- an outward group runs through the agent attached --------------
    console.log('\nscenario: run — an outward job, through the agent attached');
    await agent.presence(SESSION, LABEL, '/home/court/t9', 'pi');
    keep = setInterval(() => {
      void agent.presence(SESSION, LABEL, '/home/court/t9', 'pi');
    }, 1000);
    await pg.evaluate(() => {
      location.hash = '#/review';
    });
    check(
      "review no longer lists the verified branches' group (no reload)",
      await until(
        pg,
        () =>
          window.__ranNoReload === true &&
          !!document.querySelector('.cb-apply-part[data-part="close|pr"]') &&
          !document.querySelector('.cb-apply-part[data-part="delete|branch"]'),
      ),
    );
    await openGroup(pg, 'close|pr');
    check(
      'with the agent attached the group says who runs it, and its button is on',
      await until(
        pg,
        () =>
          document.querySelector('[data-testid="runs-through"]')
            ?.textContent === 'runs through pi' &&
          !document.querySelector('[data-testid="run-group"] .kit-btn')
            ?.disabled,
        undefined,
        8000,
      ),
    );
    const out = await readGroup(pg);
    check(
      `one button named for the action, none in the bar, no menu, no session's name on a button (${out.btns.join(', ')}; "${out.primary}")`,
      out.btns.join(',') === 'Close 1 PR' &&
        out.primary === '' &&
        out.selects === 0 &&
        !out.btns.some((b) => /t9|probe/.test(b)),
    );
    const agentC = await cssColor(pg, 'var(--kit-agent)');
    const signal = await cssColor(pg, 'var(--kit-signal)');
    const laneColours = await pg.$$eval(
      '[data-testid="review-group"] .cb-group-label .cb-lanename',
      (els) => els.map((e) => `${e.textContent}=${getComputedStyle(e).color}`),
    );
    checkList(
      "the outward part names the attached agent, in the agent's colour",
      laneColours,
      [`pi \u00b7 outward=${agentC}`],
    );
    // Run from the keyboard: "a" runs the open group.
    const prRan = await runOpen(pg, 'key');
    const jobId = prRan.body.job?.id;
    check(
      `"a" runs it through the attached agent: serve has job ${jobId} approved for ${SESSION}`,
      jobId === prRan.stored.job?.id &&
        prRan.sent.plan_id === jobId &&
        prRan.sent.session === SESSION &&
        (await eventually(async () => {
          const j = (await job(jobId)).job;
          return j.state !== 'planned' && j.session === SESSION;
        })),
    );
    check(
      'the job opens in run, listed under Running on top, its primary "Pause job", the bar counting it',
      (await inRun(pg, jobId)) &&
        (await until(
          pg,
          (id) =>
            document.querySelector('.cb-apply-part')?.dataset.part ===
              'running' &&
            document.querySelector('.cb-apply-list .kit-row')?.dataset.job ===
              String(id) &&
            document.querySelector('.kit-primary')?.textContent ===
              'Pause job' &&
            document.querySelector('.kit-ctl[data-id="run"] .kit-n')
              ?.textContent === '1',
          jobId,
          8000,
        )),
    );
    const jobH1 = await pg.$eval('.cb-job .kit-h1', (e) => e.textContent);
    check(
      `the job is attached to the composer with its title ("${await read.attached(pg)}")`,
      (await read.attached(pg)) ===
        `job #${jobId} \u00b7 ${jobH1.toLowerCase()}`,
    );
    jv = await job(jobId);
    const prStep = stepFor(jv, PR);
    const readRow = () =>
      evalNow(pg, `.cb-apply-list .kit-row[data-job="${jobId}"]`, (e) => ({
        lanes: [...e.querySelectorAll('.cb-lanename')].map(
          (l) => `${l.textContent}=${getComputedStyle(l).color}`,
        ),
        barH: e.querySelector('.cb-pb')?.getBoundingClientRect().height,
        sub: e.querySelector('.kit-sub')?.textContent,
        subColour: e.querySelector('.kit-sub')
          ? getComputedStyle(e.querySelector('.kit-sub')).color
          : '',
      }));
    await until(
      pg,
      (id) =>
        !!document.querySelector(
          `.cb-apply-list .kit-row[data-job="${id}"] .cb-pb`,
        ),
      jobId,
    );
    const row = await readRow();
    checkList('the row names its lane in its colour', row.lanes, [
      `pi \u00b7 outward=${agentC}`,
    ]);
    check(`the row has a thin bar (${row.barH}px)`, row.barH === 3);

    // Approving was the go: no batch card, nothing more to confirm.
    // Kept: a window for a card to arrive, if one did.
    await pg.waitForTimeout(500);
    check(
      `approving opens no batch card: nothing more to confirm, on the page or in serve (cards: ${(await read.cards(pg)).join(', ') || 'none'})`,
      !(await pg.$('.cb-needs-card[data-kind="batch"]')) &&
        ((await job(jobId)).needs_you ?? []).every((c) => c.kind !== 'batch'),
    );

    // ---- pause and resume --------------------------------------------
    console.log('\nscenario: run — pause and resume');
    await pg.click('.kit-primary');
    check(
      'primary "Pause job" pauses it (serve) and becomes "Resume"',
      (await eventually(async () => (await job(jobId)).job.paused === true)) &&
        (await until(
          pg,
          () =>
            document.querySelector('.kit-primary')?.textContent === 'Resume',
        )),
    );
    check(
      'the paused job says so in its kicker',
      (await pg.$eval('.cb-job .kit-kick', (e) => e.textContent)).endsWith(
        '\u00b7 paused',
      ),
    );
    await pressWithFocus(pg, 'p');
    check(
      '"p" resumes it: serve has it running, the primary is "Pause job" again',
      (await eventually(async () => (await job(jobId)).job.paused === false)) &&
        (await until(
          pg,
          () =>
            document.querySelector('.kit-primary')?.textContent === 'Pause job',
        )),
    );

    // ---- a burst of live events -----------------------------------------
    console.log('\nscenario: run — a burst of live events');
    // Kept: the resume's own reloads settle, so the count below is the
    // burst's alone.
    await pg.waitForTimeout(500);
    const fetched = [];
    const onFetch = (r) => {
      const u = r.url();
      if (/\/api\/(jobs|job|needs-you)(\?|$)/.test(u))
        fetched.push(new URL(u).pathname);
    };
    pg.on('request', onFetch);
    const BURST = 30;
    await Promise.all(
      Array.from({ length: BURST }, () =>
        agent.api('POST', '/api/jobs/resume', { id: jobId }),
      ),
    );
    await agent.api('POST', '/api/jobs/pause', { id: jobId });
    const settled = await until(
      pg,
      () => document.querySelector('.kit-primary')?.textContent === 'Resume',
    );
    // Kept: the check bounds how many fetches the burst makes; a window
    // for any extra one to go out.
    await pg.waitForTimeout(600);
    pg.off('request', onFetch);
    check(
      `${BURST + 1} live events fetch a bounded few times, one reload in flight and one queued (${fetched.length} fetches), and the page ends current (paused)`,
      settled && fetched.length >= 1 && fetched.length <= 4,
    );
    check(
      `an event about a job fetches that job only, never the whole list (${[...new Set(fetched)].join(', ')})`,
      fetched.every((p) => p === '/api/job'),
    );
    await agent.api('POST', '/api/jobs/resume', { id: jobId });
    await until(
      pg,
      () => document.querySelector('.kit-primary')?.textContent === 'Pause job',
    );

    // ---- a stale reload never draws over a newer one ------------------
    console.log('\nscenario: run — reloads in order');
    // Kept: the reloads before this settle, so the one held is the event's.
    await pg.waitForTimeout(500);
    const heldJob = [];
    await pg.route(/\/api\/job\?/, async (route) => {
      const resp = await route.fetch();
      const body = await resp.text();
      heldJob.push(() =>
        route.fulfill({ response: resp, body }).catch(() => {}),
      );
    });
    await agent.api('POST', '/api/jobs/resume', { id: jobId });
    const heldOne = await eventually(async () => heldJob.length >= 1);
    await pg.click('.kit-primary');
    const pausedDrawn = await until(
      pg,
      () => document.querySelector('.kit-primary')?.textContent === 'Resume',
    );
    await heldJob.shift()?.();
    // Kept: the older reply must not be drawn; a window for it to be.
    await pg.waitForTimeout(500);
    const after = await read.primary(pg);
    const kick = await pg.$eval('.cb-job .kit-kick', (e) => e.textContent);
    for (const f of heldJob.splice(0)) await f();
    await pg.unroute(/\/api\/job\?/);
    check(
      `an older reload (running) answered after Court's pause (paused) is not drawn (primary "${after}", kicker "${kick}")`,
      heldOne && pausedDrawn && after === 'Resume' && kick.endsWith('paused'),
    );
    await pg.click('.kit-primary');
    await eventually(async () => (await job(jobId)).job.paused === false);
    await until(
      pg,
      () => document.querySelector('.kit-primary')?.textContent === 'Pause job',
    );

    // ---- a needs-you card (spec §10) ----------------------------------
    console.log('\nscenario: run — a needs-you card');
    await agent.api('POST', '/api/agent/job-step', {
      session: SESSION,
      job: jobId,
      step: prStep.id,
      state: 'started',
    });
    const draft =
      'Closing this: lodash 4.17.21 is already on main through #9. Thanks!';
    await agent.api('POST', '/api/agent/job-ask', {
      session: SESSION,
      job: jobId,
      step: prStep.id,
      question: "public comment on someone else's repo",
      text: draft,
    });
    check(
      'a needs-you card offers the four actions (a fake casebook_job_ask)',
      (await until(
        pg,
        () => !!document.querySelector('.cb-needs-card[data-kind="text"]'),
        undefined,
        8000,
      )) &&
        (await evalNow(pg, '.cb-needs-card[data-kind="text"]', (e) =>
          [...e.querySelectorAll('.kit-btn')]
            .map((b) => b.textContent)
            .join(' \u00b7 '),
        )) ===
          'post and close \u00b7 edit text \u00b7 close without comment \u00b7 skip',
    );
    const textCard = await evalNow(
      pg,
      '.cb-needs-card[data-kind="text"]',
      (e) => ({
        edge: getComputedStyle(e).borderLeftColor,
        edgeW: getComputedStyle(e).borderLeftWidth,
        head: e.querySelector('.kit-card-head').textContent,
        quote: e.querySelector('blockquote').textContent,
        what: e.querySelector('.cb-needs-what').textContent,
        fill: e.querySelector('.kit-btn.fill')?.textContent,
      }),
    );
    check(
      `the card has a signal edge (${textCard.edge}, ${textCard.edgeW}) and quotes the exact text`,
      textCard.edge === signal &&
        textCard.edgeW === '3px' &&
        textCard.quote === draft &&
        textCard.fill === 'post and close',
    );
    const asked = await readRow();
    check(
      `the job's row says what needs Court in a signal line ("${asked.sub}")`,
      asked.sub === '1 needs you' && asked.subColour === signal,
    );
    check(
      `it says what and where ("${textCard.head}"; "${textCard.what}")`,
      textCard.head ===
        "needs you \u00b7 public comment on someone else's repo" &&
        textCard.what === 'schuettc/t9-prs#7 \u00b7 close with comment',
    );
    await pg.$eval('.cb-job .kit-kick', (e) =>
      e.scrollIntoView({ block: 'start' }),
    );
    await inThemes(t, pg, 'job');
    // edit text: the text becomes a field; saving sends it to serve.
    await pg.click(
      '.cb-needs-card[data-kind="text"] .kit-btn:has-text("edit text")',
    );
    const edited = 'Closing: already on main through #9.';
    const field = '.cb-needs-card[data-kind="text"] textarea.cb-needs-edit';
    await pg.fill(field, '');
    await pg.keyboard.type('Closing: already ');
    // While Court types, job events arrive; each redraws the open job.
    const fieldBefore = await pg.$(field);
    for (let i = 0; i < 5; i++)
      await agent.api('POST', '/api/jobs/resume', { id: jobId });
    const replaced = await until(
      pg,
      ([el, sel]) =>
        !!document.querySelector(sel) && el !== document.querySelector(sel),
      [fieldBefore, field],
    );
    // Kept: a window for a late redraw.
    await pg.waitForTimeout(300);
    await pg.keyboard.type('on main through #9.');
    const ed = await evalNow(pg, field, (e) => ({
      value: e.value,
      focused: document.activeElement === e,
      caret: e.selectionStart,
    })).catch(() => ({ value: '(no field)', focused: false, caret: -1 }));
    check(
      `typing in "edit text" while job events redraw the job: the text, focus and caret survive ("${ed.value}", focused ${ed.focused}, caret ${ed.caret}; redrawn ${replaced})`,
      replaced &&
        ed.value === edited &&
        ed.focused &&
        ed.caret === edited.length,
    );
    await pg.click(
      '.cb-needs-card[data-kind="text"] .kit-btn:has-text("save text")',
    );
    check(
      'edit text saves the new text to serve; the card stays, quoting it',
      (await eventually(async () =>
        ((await job(jobId)).needs_you ?? []).some(
          (c) => c.kind === 'text' && c.text === edited,
        ),
      )) &&
        (await until(
          pg,
          (w) =>
            document.querySelector(
              '.cb-needs-card[data-kind="text"] blockquote',
            )?.textContent === w,
          edited,
        )),
    );
    await pg.click(
      '.cb-needs-card[data-kind="text"] .kit-btn:has-text("post and close")',
    );
    check(
      'answering closes the card: "post and close" clears it (page and serve)',
      (await until(
        pg,
        () => !document.querySelector('.cb-needs-card[data-kind="text"]'),
      )) &&
        (await eventually(async () =>
          ((await job(jobId)).needs_you ?? []).every((c) => c.kind !== 'text'),
        )),
    );
    check(
      'serve took the approved text for the step',
      await eventually(
        async () => stepFor(await job(jobId), PR).text === edited,
      ),
    );
    await agent.api('POST', '/api/agent/job-step', {
      session: SESSION,
      job: jobId,
      step: prStep.id,
      state: 'reported',
    });
    check(
      'the step advances: ✓ closed · verified',
      await until(
        pg,
        (id) =>
          document.querySelector(`.cb-step[data-step="${id}"] .cb-step-t`)
            ?.textContent === 'closed \u00b7 verified',
        prStep.id,
        8000,
      ),
    );
    checkList(
      'serve verified it with one gh read, to the probe fake (nothing else ran gh)',
      serveHandle.ghCalls(),
      ['pr view 7 -R schuettc/t9-prs --json state'],
    );

    // ---- the agent pauses a step ----------------------------------------
    console.log('\nscenario: run — the agent pauses a step');
    await pg.evaluate(() => {
      location.hash = '#/review';
    });
    await openGroup(pg, 'archive|repo');
    await until(
      pg,
      () =>
        document.querySelector('[data-testid="runs-through"]')?.textContent ===
        'runs through pi',
    );
    const archRan = await runOpen(pg);
    const archJob = archRan.body.job?.id;
    await inRun(pg, archJob);
    const repoA = stepFor(await job(archJob), REPO_A);
    await agent.api('POST', '/api/agent/job-step', {
      session: SESSION,
      job: archJob,
      step: repoA.id,
      state: 'started',
    });
    await agent.api('POST', '/api/agent/job-step', {
      session: SESSION,
      job: archJob,
      step: repoA.id,
      state: 'paused',
      detail: 'a reviewer commented 2 days ago, after you decided.',
    });
    check(
      'a paused step is an amber "pi paused" card: archive anyway (danger) · skip',
      (await until(
        pg,
        () => !!document.querySelector('.cb-needs-card[data-kind="paused"]'),
      )) &&
        (await evalNow(pg, '.cb-needs-card[data-kind="paused"]', (e) => {
          const bs = [...e.querySelectorAll('.kit-btn')];
          return [
            e.querySelector('.kit-card-head').textContent,
            getComputedStyle(e).borderLeftColor,
            bs
              .map(
                (b) =>
                  b.textContent + (b.classList.contains('danger') ? '!' : ''),
              )
              .join(' \u00b7 '),
          ].join(' | ');
        })) === `pi paused | ${agentC} | archive anyway! \u00b7 skip`,
    );
    check(
      'the steps table shows it ‖ paused, and the row "1 paused"',
      (await until(
        pg,
        (id) =>
          document.querySelector(`.cb-step[data-step="${id}"] .cb-step-g`)
            ?.textContent === '\u2016',
        repoA.id,
      )) &&
        (await until(
          pg,
          (id) =>
            document.querySelector(
              `.cb-apply-list .kit-row[data-job="${id}"] .kit-sub`,
            )?.textContent === '1 paused',
          archJob,
        )),
    );
    await pg.click(
      '.cb-needs-card[data-kind="paused"] .kit-btn:has-text("skip")',
    );
    check(
      'skip clears the paused card and skips the step',
      (await until(
        pg,
        () => !document.querySelector('.cb-needs-card[data-kind="paused"]'),
      )) &&
        (await eventually(
          async () => stepFor(await job(archJob), REPO_A).state === 'skipped',
        )),
    );

    // ---- a step that can never be undone ------------------------------
    console.log(
      '\nscenario: review — a step that can never be undone asks first',
    );
    await pg.evaluate(() => {
      location.hash = '#/review';
    });
    await openGroup(pg, 'delete|repo');
    const holds = async (key) =>
      (await jobsNow()).some((j) => (j.steps ?? []).some((x) => x.key === key));
    await until(
      pg,
      () => !!document.querySelector('[data-testid="runs-through"]'),
    );
    const perm = await readGroup(pg);
    const dangerC = await cssColor(pg, 'var(--kit-danger)');
    const warnC = await pg.$eval(
      '.cb-approve-warn',
      (e) => getComputedStyle(e).color,
    );
    check(
      `it says it cannot be undone, in danger red ("${perm.what}"; ${warnC})`,
      perm.warn === ' Cannot be undone.' && warnC === dangerC,
    );
    await pg.click(`${RUN_BTN}.fill`);
    const sheetOf = () =>
      pg.evaluate(() => {
        const s = document.querySelector('.kit-sheet');
        return s
          ? {
              head: s.querySelector('.kit-sheet-head')?.textContent ?? '',
              items: [...s.querySelectorAll('.cb-confirm-list li')].map(
                (e) => e.textContent,
              ),
              btns: [...s.querySelectorAll('.kit-btn')].map(
                (b) => b.textContent,
              ),
            }
          : null;
      });
    await until(pg, () => !!document.querySelector('.kit-sheet'));
    const sh = await sheetOf();
    check(
      `pressing it asks first, naming each item, on a button named for what runs (${JSON.stringify(sh)})`,
      sh?.head.startsWith('Delete 1 repo?') &&
        sh.items.join(',') === 'schuettc/t9-dormant-c' &&
        sh.btns.join(',') === 'cancel,Delete 1 repo',
    );
    await inThemes(t, pg, 'confirm', ['light']);
    await pg.click('.kit-sheet .kit-btn:has-text("cancel")');
    check(
      'cancel runs nothing: the sheet goes, serve stores no plan for it, review stays',
      (await until(pg, () => !document.querySelector('.kit-sheet'))) &&
        !(await holds(REPO_C)) &&
        (await pg.evaluate(() => location.hash)) === groupHash('delete|repo'),
    );
    await pg.click(`${RUN_BTN}.fill`);
    await until(pg, () => !!document.querySelector('.kit-sheet'));
    const delRun = runOpen(pg, 'click', false);
    await pg.click('.kit-sheet .kit-btn.fill');
    const delRan = await delRun;
    const delJob = delRan.body.job?.id;
    check(
      `confirmed, it is stored and runs, and opens in run (job ${delJob})`,
      !!delJob && delRan.stored.job?.id === delJob && (await inRun(pg, delJob)),
    );

    // ---- undo ----------------------------------------------------------
    console.log('\nscenario: run — undo');
    await pg.evaluate((id) => {
      location.hash = `#/run/${id}`;
    }, branchJob);
    await inRun(pg, branchJob);
    jv = await job(branchJob);
    const undoable = jv.job.steps.filter((s) => s.undoable).map((s) => s.id);
    await until(
      pg,
      (n) =>
        [...document.querySelectorAll('.cb-step')].filter((e) =>
          e.querySelector('.cb-undo'),
        ).length === n,
      undoable.length,
    );
    const withUndo = await pg.$$eval('.cb-step', (els) =>
      els
        .filter((e) => e.querySelector('.cb-undo'))
        .map((e) => Number(e.dataset.step)),
    );
    checkList(
      'undo shows on exactly the finished steps serve can undo',
      withUndo,
      undoable,
    );
    check(
      '(serve can undo the local and remote deletes)',
      undoable.includes(l1.id) && undoable.includes(l1r.id),
    );
    check(
      'a finished job opened by its address shows under Done, opened',
      await until(
        pg,
        (id) =>
          !!document.querySelector(
            `.cb-apply-list .kit-row[data-job="${id}"]`,
          ) &&
          !document
            .querySelector('.cb-apply-part[data-part="done"]')
            ?.classList.contains('collapsed'),
        branchJob,
      ),
    );
    const undoPost = pg.waitForResponse((r) =>
      r.url().includes('/api/jobs/undo'),
    );
    await pg.click(`.cb-step[data-step="${l1.id}"] .cb-undo`);
    const undoBody = (await undoPost).request().postDataJSON();
    check(
      `undo posts /api/jobs/undo for the step (${JSON.stringify(undoBody)})`,
      undoBody?.step === l1.id,
    );
    check(
      'serve restored feat/l1 in the probe clone (git branch, in the temp home)',
      await eventually(
        async () =>
          gitIn(home, 't9-local', 'rev-parse --verify refs/heads/feat/l1') ===
          l1.expected_tip,
      ),
    );
    check(
      'the step reads ↺ undone, with no undo',
      await until(
        pg,
        (id) => {
          const e = document.querySelector(`.cb-step[data-step="${id}"]`);
          return (
            e?.querySelector('.cb-step-t')?.textContent === 'undone' &&
            !e.querySelector('.cb-undo')
          );
        },
        l1.id,
      ),
    );

    // ---- sections: context, keys, primary -----------------------------
    console.log('\nscenario: run — only the active section feeds the composer');
    // The repo delete's job is still running (its agent hasn't reported).
    const liveJob = delJob;
    await pg.evaluate((id) => {
      location.hash = `#/run/${id}`;
    }, liveJob);
    await inRun(pg, liveJob);
    const liveH1 = await pg.$eval('.cb-job .kit-h1', (e) => e.textContent);
    check(
      'the open job is attached to the composer, with its title',
      await until(
        pg,
        (want) =>
          document.querySelector('[data-testid="composer-attached"]')
            ?.textContent === want,
        `job #${liveJob} \u00b7 ${liveH1.toLowerCase()}`,
      ),
    );
    await pg.evaluate(() => {
      location.hash = '#/decide';
    });
    await until(
      pg,
      () =>
        document
          .querySelector('.kit-ctl[data-id="decide"]')
          ?.classList.contains('on') ?? false,
    );
    const onDecide = await read.attached(pg);
    const primaryDecide = await read.primary(pg);
    check(
      `on decide the bar's primary is decide's, not "Pause job" ("${primaryDecide}")`,
      primaryDecide !== 'Pause job' && primaryDecide !== 'Resume',
    );
    await agent.api('POST', '/api/jobs/pause', { id: liveJob });
    // Kept: the hidden section's event must not change the line; a window
    // for the event to reach the page.
    await pg.waitForTimeout(800);
    check(
      `a live job event while run is hidden leaves the attached line alone ("${onDecide}")`,
      (await read.attached(pg)) === onDecide &&
        (await read.primary(pg)) === primaryDecide &&
        !onDecide.includes(`job #${liveJob}`),
    );
    await pressWithFocus(pg, 'p');
    await pg.waitForTimeout(500); // kept: "p" must resume nothing; a window for it
    check(
      'on decide, "p" does not resume the hidden job',
      (await job(liveJob)).job.paused === true,
    );
    await pg.evaluate((id) => {
      location.hash = `#/run/${id}`;
    }, liveJob);
    check(
      'back on run the job is attached again and the primary follows its state ("Resume")',
      await until(
        pg,
        (want) =>
          document.querySelector('[data-testid="composer-attached"]')
            ?.textContent === want &&
          document.querySelector('.kit-primary')?.textContent === 'Resume',
        `job #${liveJob} \u00b7 ${liveH1.toLowerCase()}`,
      ),
    );
    const overlay = async () => {
      await pressWithFocus(pg, '?');
      await until(pg, () => !!document.querySelector('.kit-keys'));
      const labels = await pg.$$eval(
        '.kit-keys .kit-keys-row > span:last-child',
        (els) => els.map((e) => e.textContent),
      );
      await pg.keyboard.press('Escape');
      await until(pg, () => !document.querySelector('.kit-keys'));
      return labels;
    };
    const onRun = await overlay();
    check(
      `on run the overlay lists its keys, not decide's or review's (${onRun.join(', ')})`,
      !onRun.includes('run the open group') &&
        onRun.includes('pause or resume the open job') &&
        onRun.includes('go to review') &&
        onRun.includes('go to run') &&
        !onRun.includes('accept proposal'),
    );
    await pressWithFocus(pg, 'p');
    check(
      'back on run, "p" resumes it',
      await eventually(async () => (await job(liveJob)).job.paused === false),
    );
    // A group open in review, then decide: "a" there runs nothing.
    await agent.api('POST', '/api/decide', {
      keys: ['repo:schuettc/t9-prs'],
      disposition: 'archive',
    });
    await pg.evaluate(() => {
      location.hash = '#/review';
    });
    await openGroup(pg, 'archive|repo');
    check(
      'a group open in review, with the agent attached: its button is on ("Archive 1 repo"), and the bar has no primary',
      await until(
        pg,
        (b) => {
          const p = document.querySelector('.kit-primary');
          const btn = document.querySelector(b);
          return (
            btn?.textContent === 'Archive 1 repo' &&
            !btn.disabled &&
            (!p || p.hidden || p.textContent === '')
          );
        },
        `${RUN_BTN}.fill`,
      ),
    );
    await pg.evaluate(() => {
      location.hash = '#/decide';
    });
    await until(
      pg,
      () =>
        document
          .querySelector('.kit-ctl[data-id="decide"]')
          ?.classList.contains('on') ?? false,
    );
    const primaryBack = await read.primary(pg);
    check(
      `on decide the primary isn't review's or run's ("${primaryBack}")`,
      !['Pause job', 'Resume', 'Archive 1 repo'].includes(primaryBack),
    );
    await pressWithFocus(pg, 'a');
    await pg.waitForTimeout(600); // kept: "a" must run nothing; a window for it
    check(
      'on decide, "a" does not run the group open in review (no plan stored for it)',
      !(await holds('repo:schuettc/t9-prs')),
    );
    check(`no page errors (${errors.join(' | ')})`, errors.length === 0);
  } finally {
    if (keep) clearInterval(keep);
    await pg.close();
  }
}

// staleScenario: serve refuses a plan from an observation older than one
// sync interval (2s here); the open group shows serve's words, which name
// a sync and no command, and offers "sync first". The sync is serve's own
// (POST /api/sync, App.Sync) and hermetic: this serve's casebook-data remote
// is the committed bundle on disk, it scans no roots, and its gh is the
// probe's fake. When it is done the group is opened again, once.
async function staleScenario(context, t, serveHandle) {
  const { check, checkList, until, eventually } = t;
  console.log(
    '\nscenario: review — a stale group is refused in serve\u2019s words, and offers sync first',
  );
  const agent = createAgent(serveHandle.base, serveHandle.token);
  const pg = await context.newPage();
  const errors = [];
  pg.on('pageerror', (e) => errors.push(String(e)));
  try {
    await pg.setViewportSize({ width: 1600, height: 900 });
    await pg.goto(serveHandle.url + '#/review', {
      waitUntil: 'domcontentloaded',
      timeout: 15000,
    });
    const STALE_BUTTON = HEAD('archive|repo');
    const okC = await cssColor(pg, 'var(--kit-ok)');
    const mutedC = await cssColor(pg, 'var(--kit-muted)');
    await until(
      pg,
      () => document.querySelector('.kit-live')?.dataset.state === 'live',
    );
    await agent.api('POST', '/api/decide', {
      keys: ['repo:schuettc/t9-stale'],
      disposition: 'archive',
    });
    const sum0 = await agent.api('GET', '/api/summary');
    check(
      `serve gives the page its sync interval (${sum0.sync_interval_ms} ms) and the index's age (built_at)`,
      sum0.sync_interval_ms === 2000 && !!Date.parse(sum0.built_at),
    );
    const freshSeen = await until(
      pg,
      ([ok, b]) => {
        const e = document.querySelector('.cb-apply-obs');
        return (
          e?.textContent === 'observations fresh' &&
          e.dataset.state === 'fresh' &&
          getComputedStyle(e).color === ok &&
          !!document.querySelector(b)
        );
      },
      [okC, STALE_BUTTON],
      1800,
    );
    const obsNow = await pg.$eval('.cb-apply-obs', (e) => e.textContent);
    check(
      `within one sync interval of the rebuild the foot says the observations are fresh, in ok ("${obsNow}")`,
      freshSeen,
    );
    check(
      'past one sync interval the foot says how old the observations are, muted',
      await until(
        pg,
        (muted) => {
          const e = document.querySelector('.cb-apply-obs');
          return (
            /^observations \d+s old$/.test(e?.textContent ?? '') &&
            e.dataset.state === 'stale' &&
            getComputedStyle(e).color === muted
          );
        },
        mutedC,
        5000,
      ),
    );
    await eventually(async () => {
      const s = await agent.api('GET', '/api/summary');
      return Date.now() - Date.parse(s.built_at) > 2300;
    });
    const resp = pg.waitForResponse((r) => r.url().includes('/api/apply/plan'));
    await pg.click(STALE_BUTTON);
    const r = await resp;
    const body = await r.json();
    check(`serve refuses the plan (409: ${body.error})`, r.status() === 409);
    const shown = (await until(
      pg,
      () => !!document.querySelector('[data-testid="plan-refused"]'),
    ))
      ? await pg.$eval('[data-testid="plan-refused"]', (e) => e.textContent)
      : '';
    check(
      `the open group shows serve's message verbatim ("${shown}")`,
      shown === body.error &&
        shown.includes('sync') &&
        (await openedGroup(t, pg, 'archive|repo', null)),
    );
    check(
      'it names a sync and no command for Court to run',
      !/\bgit |\bgh |casebook |`|\$/.test(shown),
    );
    const offer = await pg.$eval('[data-testid="refusal"]', (e) => ({
      btns: [...e.querySelectorAll('.kit-btn')].map(
        (b) => b.textContent + (b.classList.contains('fill') ? '*' : ''),
      ),
    }));
    check(
      `the refusal offers "sync first" (${offer.btns.join(', ')})`,
      offer.btns.join(',') === 'sync first*',
    );
    const around = await pg.evaluate(() => ({
      filled: [
        ...document.querySelectorAll(
          '[data-testid="review-group"] .kit-btn.fill',
        ),
      ].map((b) => b.textContent),
      primary: document.querySelector('.kit-primary')?.textContent,
    }));
    check(
      `refused, sync first is the one offer: the only filled button in the document, and the bar's primary (${around.filled.join(', ')}; "${around.primary}")`,
      around.filled.join(',') === 'sync first' &&
        around.primary === 'Sync first',
    );
    check(
      'no job was made',
      ((await agent.api('GET', '/api/jobs')).jobs ?? []).length === 0,
    );
    await inThemes(t, pg, 'stale', ['light']);

    // A second tab, refused too. Court presses sync first in the first
    // tab only: the group is opened again there, not here.
    const pg2 = await context.newPage();
    pg2.on('pageerror', (e) => errors.push(`tab 2: ${e}`));
    await pg2.setViewportSize({ width: 1600, height: 900 });
    await pg2.goto(serveHandle.url + '#/review', {
      waitUntil: 'domcontentloaded',
      timeout: 15000,
    });
    await until(
      pg2,
      () => document.querySelector('.kit-live')?.dataset.state === 'live',
    );
    await until(pg2, (b) => !!document.querySelector(b), STALE_BUTTON);
    const resp2 = pg2.waitForResponse((x) =>
      x.url().includes('/api/apply/plan'),
    );
    await pg2.click(STALE_BUTTON);
    check(
      'a second tab is refused too',
      (await resp2).status() === 409 &&
        (await until(
          pg2,
          () => !!document.querySelector('[data-testid="plan-refused"]'),
        )),
    );
    const tab2Plans = [];
    const onReq2 = (q) => {
      if (q.method() === 'POST' && q.url().includes('/api/apply/plan'))
        tab2Plans.push(q.url());
    };
    pg2.on('request', onReq2);

    const ghBefore = serveHandle.ghCalls().length;
    const bundle = join(here, 'testdata', 'home', 'data', 'repo.bundle');
    const bundleBefore = hashOf(bundle);
    await pg.evaluate(() => {
      window.__sawSyncing = false;
      new MutationObserver(() => {
        if (document.querySelector('[data-testid="syncing"]'))
          window.__sawSyncing = true;
      }).observe(document.body, { subtree: true, childList: true });
    });
    const syncPosts = [];
    const planPosts = [];
    const onReq = (q) => {
      if (q.method() !== 'POST') return;
      if (q.url().includes('/api/sync')) syncPosts.push(q.url());
      if (q.url().includes('/api/apply/plan')) planPosts.push(q.url());
    };
    pg.on('request', onReq);
    const replanned = pg
      .waitForResponse(
        (x) => x.url().includes('/api/apply/plan') && x.status() === 200,
        { timeout: 30000 },
      )
      .catch(() => null);
    await pg.click('[data-testid="refusal"] .kit-btn:has-text("sync first")');
    const plan = await (await replanned)?.json().catch(() => null);
    await pg.waitForTimeout(300); // kept: no second sync or plan may follow
    pg.off('request', onReq);
    check(
      `sync first asks serve to sync once (${syncPosts.length} POST /api/sync), and the page said it was syncing`,
      syncPosts.length === 1 &&
        (await pg.evaluate(() => window.__sawSyncing === true)),
    );
    check(
      `when the sync is done the group is opened again, once, with its button (${planPosts.length} POST /api/apply/plan; job ${plan?.job?.id})`,
      planPosts.length === 1 &&
        !!plan &&
        (await openedGroup(t, pg, 'archive|repo', plan.job.id)),
    );
    checkList(
      "the plan is the refused one: the stale repo's archive",
      (plan?.job?.steps ?? []).map((st) => `${st.action} ${st.key}`),
      ['repo-archive repo:schuettc/t9-stale'],
    );
    const sum1 = await agent.api('GET', '/api/summary');
    check(
      `serve's index was rebuilt by the sync (built_at ${sum0.built_at} → ${sum1.built_at}) and it is no longer syncing`,
      Date.parse(sum1.built_at) > Date.parse(sum0.built_at) + 2000 &&
        sum1.syncing === false,
    );
    check(
      'serve stores no plan (the reopened group is a preview; no second sync)',
      ((await agent.api('GET', '/api/jobs')).jobs ?? []).length === 0,
    );
    const tab2Cleared = await until(
      pg2,
      () =>
        !document.querySelector('[data-testid="refusal"]') &&
        document.querySelector('.kit-primary')?.textContent !== 'Sync first',
    );
    await pg2.waitForTimeout(300); // kept: this tab may not re-plan
    pg2.off('request', onReq2);
    check(
      `the tab that didn't ask re-plans nothing (${tab2Plans.length} POST /api/apply/plan) and no longer shows the refusal`,
      tab2Plans.length === 0 && tab2Cleared,
    );
    await pg2.close();
    const ghSync = serveHandle.ghCalls().slice(ghBefore);
    const outward =
      /\b(archive|unarchive|close|reopen|merge|delete|edit|create|comment|review)\b|(^| )(-X|--method)( |$)/;
    check(
      `nothing outward ran in the sync: ${ghSync.length} gh calls, all to the fake, none that changes GitHub (${[...new Set(ghSync.map((c) => c.split(' ').slice(0, 2).join(' ')))].join('; ')})`,
      ghSync.length > 0 && ghSync.every((c) => !outward.test(c)),
    );
    const env = serveHandle.ghEnv();
    check(
      `every gh call saw the fail-closed environment (GH_HOST github.invalid, a token that isn't one, the probe's gh config; ${env.length} calls)`,
      env.length === serveHandle.ghCalls().length &&
        env.every(
          (l) =>
            l ===
            `github.invalid probe-not-a-token ${join(serveHandle.home, 'gh')}`,
        ),
    );
    const origin = probeGit(
      join(serveHandle.home, 'data', 'repo'),
      'remote get-url origin',
      serveHandle.home,
    );
    check(
      `casebook-data's remote, which the sync pulls from and pushes to, is the bundle on disk (${origin})`,
      origin.endsWith('/testdata/home/data/repo.bundle') &&
        !/^(https?|ssh|git):|@/.test(origin),
    );
    check(
      'the committed bundle is untouched by the sync (nothing was pushed to it)',
      hashOf(bundle) === bundleBefore,
    );
    check(`no page errors (${errors.join(' | ')})`, errors.length === 0);
  } finally {
    await pg.close();
  }
}
