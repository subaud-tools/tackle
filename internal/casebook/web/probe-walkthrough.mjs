// probe-walkthrough.mjs — the first-run walkthrough (casebook 0.14.0), run
// by probe.mjs with its check helpers (t.check, t.checkList, t.until,
// t.eventually).
//
// One serve of its own (serve.mjs: --no-open, CASEBOOK_NO_BROWSER=1, the
// hermetic git guard, the fake gh), on the plain fixture. Its context is a
// first visit: no walkthrough flag, and the agent pane as a first visit
// finds it (every other group's context is marked seen, probe.mjs).
//
//   the first open shows the overlay, and a list row can't be clicked
//   each step outlines its control (to ±2 px), on screen, its note not
//   covering it, in order: decide, filters, rules, a decision button
//   (the first item opened), review, the agent (the pane opened)
//   nothing live moves under it: a decision made meanwhile shows after
//   keys don't reach the page; done and skip put the page back as it was
//   and decide nothing; a reload shows no overlay; g w shows it again;
//   it can be skipped at every step (dark); the ? overlay lists g w;
//   an empty list has no decision step

import { startServe } from './serve.mjs';
import { createAgent } from './agent.mjs';

let serves = [];
/** stopWalkthroughServes stops every serve these scenarios started (probe.mjs cleanup). */
export function stopWalkthroughServes() {
  for (const s of serves) s.stop();
  serves = [];
}

const IDS = ['decide', 'filters', 'rules', 'decision', 'review', 'agent'];
// What each step points at, found the probe's own way.
const TARGETS = {
  decide: '.kit-ctl[data-id="decide"]',
  filters: '.kit-list:not([hidden]) .cb-filters',
  rules: '.kit-list:not([hidden]) .cb-rule-row',
  decision: '.kit-read:not([hidden]) .cb-decide-row .cb-choice',
  review: '.kit-ctl[data-id="review"]',
  agent: '.kit-rail',
};

// state reads the overlay: whether it shows, its step, and the boxes of
// its outline, its note and the control the probe finds for the step.
async function state(pg) {
  return pg.evaluate((targets) => {
    const box = (el) => {
      if (!el) return null;
      const r = el.getBoundingClientRect();
      return { x: r.x, y: r.y, width: r.width, height: r.height };
    };
    const walk = document.querySelector('.cb-walk');
    const shown = !!walk && !walk.hidden && walk.isConnected;
    const step = shown ? (walk.dataset.step ?? '') : '';
    const target = step ? document.querySelector(targets[step] ?? '') : null;
    // What a click at the control's middle would reach.
    let hit = '';
    if (target) {
      const r = target.getBoundingClientRect();
      const at = document.elementFromPoint(
        r.x + r.width / 2,
        r.y + r.height / 2,
      );
      hit = at && walk?.contains(at) ? 'overlay' : (at?.className ?? '');
    }
    return {
      shown,
      step,
      count: shown
        ? (walk.querySelector('.cb-walk-count')?.textContent ?? '')
        : '',
      note: shown
        ? (walk.querySelector('.cb-walk-note')?.textContent ?? '')
        : '',
      ring: shown ? box(walk.querySelector('.cb-walk-ring')) : null,
      card: shown ? box(walk.querySelector('.cb-walk-card')) : null,
      target: box(target),
      targetText: target?.textContent?.trim() ?? '',
      hit,
      view: { width: innerWidth, height: innerHeight },
      hash: location.hash,
      folded: !!document.querySelector('.kit-rail.cb-dock-strip'),
      rail: box(document.querySelector('.kit-rail')),
      seen: localStorage.getItem('casebook.walkthrough'),
    };
  }, TARGETS);
}

const near = (a, b) =>
  !!a &&
  !!b &&
  Math.abs(a.x - b.x) <= 2 &&
  Math.abs(a.y - b.y) <= 2 &&
  Math.abs(a.width - b.width) <= 2 &&
  Math.abs(a.height - b.height) <= 2;
const onScreen = (b, v) =>
  !!b &&
  b.width > 0 &&
  b.height > 0 &&
  b.x >= -1 &&
  b.y >= -1 &&
  b.x + b.width <= v.width + 1 &&
  b.y + b.height <= v.height + 1;
const overlaps = (a, b) =>
  !!a &&
  !!b &&
  a.x < b.x + b.width &&
  b.x < a.x + a.width &&
  a.y < b.y + b.height &&
  b.y < a.y + a.height;

export async function walkthroughScenarios(shared, t) {
  const { check, checkList, until, eventually } = t;
  console.log('\nscenario: the first-run walkthrough');
  const context = await shared.browser().newFirstVisitContext({
    viewport: { width: 1600, height: 900 },
  });
  const s = await startServe();
  serves.push(s);
  const agent = createAgent(s.base, s.token);
  const api = (m, p, b) => agent.api(m, p, b);
  const errors = [];
  try {
    const pg = await context.newPage();
    pg.on('console', (m) => m.type() === 'error' && errors.push(m.text()));
    pg.on('pageerror', (e) => errors.push(String(e)));
    await pg.goto(s.url, { waitUntil: 'domcontentloaded', timeout: 15000 });
    await pg.waitForSelector('.kit-bar', { timeout: 8000 });

    // ---- the first open ------------------------------------------------------
    check(
      'the first open shows the walkthrough',
      await until(pg, () => !!document.querySelector('.cb-walk:not([hidden])')),
    );
    await until(pg, () => !!document.querySelector('.kit-list .kit-row'));
    const rowHit = await pg.evaluate(() => {
      const row = document.querySelector('.kit-list:not([hidden]) .kit-row');
      if (!row) return 'no row';
      const r = row.getBoundingClientRect();
      const at = document.elementFromPoint(
        r.x + r.width / 2,
        r.y + r.height / 2,
      );
      return at && document.querySelector('.cb-walk')?.contains(at)
        ? 'overlay'
        : String(at?.className ?? 'nothing');
    });
    check(
      `a click on a list row reaches the overlay, not the row (${rowHit})`,
      rowHit === 'overlay',
    );
    // The first open's address has no hash: decide's. Done puts it back
    // as #/decide.
    const startHash = (await pg.evaluate(() => location.hash)) || '#/decide';
    const listed = (
      (await api('GET', '/api/items?view=all&limit=500')).items ?? []
    ).map((it) => it.key);
    const rows = () =>
      pg.$$eval('.kit-list:not([hidden]) .kit-row', (els) => els.length);
    const barN = () =>
      pg.$eval(
        '.kit-ctl[data-id="decide"] .kit-n',
        (el) => el.textContent ?? '',
      );

    // ---- keys wait --------------------------------------------------------------
    {
      const before = await state(pg);
      await pg.keyboard.press('g');
      await pg.keyboard.press('a');
      await pg.keyboard.press('j');
      await pg.waitForTimeout(300);
      const after = await state(pg);
      check(
        `keys don't reach the page under it (g a, j: ${after.hash}, step ${after.step})`,
        after.shown && after.hash === before.hash && after.step === 'decide',
      );
    }

    // ---- start to finish (light) -----------------------------------------------
    const seenIds = [];
    let decided = '';
    let other = '';
    const rows0 = await rows();
    const n0 = await barN();
    for (let i = 0; i < IDS.length; i++) {
      await until(
        pg,
        (id) => {
          const w = document.querySelector('.cb-walk');
          return !!w && w.dataset.step === id;
        },
        IDS[i],
      );
      // Let the outline settle on its control (the pane opening, an item
      // loading).
      await pg.waitForTimeout(400);
      const st = await state(pg);
      seenIds.push(st.step);
      const id = st.step || IDS[i];
      check(
        `${id}: "${st.count}" outlines its control to ±2 px (${JSON.stringify(st.ring)} vs ${JSON.stringify(st.target)})`,
        near(st.ring, st.target),
      );
      check(
        `${id}: its control is on screen, and the note beside it covers none of it`,
        onScreen(st.target, st.view) &&
          onScreen(st.card, st.view) &&
          !overlaps(st.card, st.target),
      );
      check(
        `${id}: the control can't be clicked under it (${st.hit})`,
        st.hit === 'overlay',
      );
      check(`${id}: has a note`, st.note.trim().length > 0);
      if (id === 'decision') {
        check(
          `decision: the first listed item is open (${st.hash}) and a decision button is outlined ("${st.targetText}")`,
          st.hash.startsWith('#/item/') && st.targetText.length > 0,
        );
        decided = decodeURIComponent(st.hash.replace(/^#\/item\//, ''));
        check(
          `decision: the item opened is the first listed, not the one decided meanwhile (${decided})`,
          decided === listed[0] && decided !== other,
        );
      }
      if (id === 'agent') {
        check(
          `agent: the pane is open (${st.rail?.width}px)`,
          !st.folded && (st.rail?.width ?? 0) >= 300,
        );
      }
      if (id === 'filters') {
        // A decision made elsewhere while the walkthrough shows: nothing
        // moves under it.
        // The last listed: the decision step opens the first.
        other = listed.at(-1);
        await api('POST', '/api/decide', {
          keys: [other],
          disposition: 'keep',
        });
        await pg.waitForTimeout(1500);
        const r1 = await rows();
        const n1 = await barN();
        const st2 = await state(pg);
        check(
          `nothing live moves under it: a decision made meanwhile leaves the list (${rows0} → ${r1} rows) and the count (${n0} → ${n1}) as they were`,
          r1 === rows0 && n1 === n0 && near(st2.ring, st.ring),
        );
      }
      await pg.click('.cb-walk-next');
    }
    checkList('the steps, in order', seenIds, IDS);
    check(
      'done: the overlay goes and the walkthrough is marked seen',
      await until(pg, () => !document.querySelector('.cb-walk:not([hidden])')),
    );
    {
      const st = await state(pg);
      check(
        `done: the page is back where it was (${st.hash}, the pane a strip again)`,
        st.hash === startHash && st.folded && st.seen === 'seen',
      );
    }
    check(
      'and what was live meanwhile shows now: the list and the count catch up',
      await until(
        pg,
        (want) =>
          document.querySelectorAll('.kit-list:not([hidden]) .kit-row')
            .length ===
          want - 1,
        rows0,
      ),
    );
    check(
      `the walkthrough decided nothing (${decided || 'no item'} is undecided)`,
      !!decided &&
        (await api('GET', `/api/item?key=${encodeURIComponent(decided)}`)).item
          .decision == null,
    );

    // ---- a reload: seen ---------------------------------------------------------
    await pg.reload({ waitUntil: 'domcontentloaded' });
    await pg.waitForSelector('.kit-list .kit-row', { timeout: 8000 });
    await pg.waitForTimeout(1500);
    check('a reload shows no walkthrough', !(await state(pg)).shown);

    // ---- the ? overlay lists g w ------------------------------------------------
    await pg.keyboard.press('?');
    const help = await pg
      .waitForSelector('.kit-keys', { timeout: 3000 })
      .then((el) => el.textContent())
      .catch(() => '');
    check(
      'the ? overlay lists g w: show the walkthrough',
      /g w/.test(help ?? '') && /show the walkthrough/.test(help ?? ''),
    );
    await pg.keyboard.press('Escape');
    await until(pg, () => !document.querySelector('.kit-keys'));

    // ---- skipped at every step (dark) ---------------------------------------
    await pg.evaluate(() => localStorage.setItem('casebook.theme', 'dark'));
    await pg.reload({ waitUntil: 'domcontentloaded' });
    await pg.waitForSelector('.kit-list .kit-row', { timeout: 8000 });
    check(
      'dark',
      (await pg.evaluate(() => document.documentElement.dataset.theme)) ===
        'dark',
    );
    for (let i = 0; i < IDS.length; i++) {
      const before = await state(pg);
      await pg.keyboard.press('g');
      await pg.keyboard.press('w');
      const ok = await until(
        pg,
        () => document.querySelector('.cb-walk')?.dataset.step === 'decide',
      );
      if (i === 0) check('g w shows the walkthrough again', ok);
      for (let j = 0; j < i; j++) {
        await pg.click('.cb-walk-next');
        await until(
          pg,
          (id) => document.querySelector('.cb-walk')?.dataset.step === id,
          IDS[j + 1],
        );
      }
      await pg.waitForTimeout(400);
      const st = await state(pg);
      check(
        `dark, ${IDS[i]}: outlined to ±2 px, on screen, the note beside it`,
        st.step === IDS[i] &&
          near(st.ring, st.target) &&
          onScreen(st.target, st.view) &&
          !overlaps(st.card, st.target),
      );
      // The last is skipped with Esc; the others with skip.
      if (i === IDS.length - 1) await pg.keyboard.press('Escape');
      else await pg.click('.cb-walk-skip');
      const gone = await until(
        pg,
        () => !document.querySelector('.cb-walk:not([hidden])'),
      );
      const after = await state(pg);
      check(
        `skip at ${IDS[i]}${i === IDS.length - 1 ? ' (Esc)' : ''}: the overlay goes and the page is as it was (${after.hash}, pane ${after.folded ? 'a strip' : 'open'})`,
        gone &&
          after.hash === before.hash &&
          after.folded === before.folded &&
          after.seen === 'seen',
      );
    }

    // ---- review 2026-10-17: history and navigation ------------------------
    {
      const go = (h) =>
        pg.evaluate((x) => {
          location.hash = x;
        }, h);
      await go('#/review');
      await pg.waitForTimeout(300);
      await go('#/decide');
      await pg.waitForSelector('.kit-list:not([hidden]) .kit-row', {
        timeout: 8000,
      });
      await pg.keyboard.press('g');
      await pg.keyboard.press('w');
      for (const id of IDS.slice(1, IDS.indexOf('decision') + 1)) {
        await until(pg, () => !!document.querySelector('.cb-walk'));
        await pg.click('.cb-walk-next');
        await until(
          pg,
          (x) => document.querySelector('.cb-walk')?.dataset.step === x,
          id,
        );
      }
      await pg.click('.cb-walk-skip');
      await until(pg, () => !document.querySelector('.cb-walk:not([hidden])'));
      await pg.evaluate(() => history.back());
      check(
        `the walkthrough leaves no history: Back after it goes where Back went before it (${await pg.evaluate(() => location.hash)})`,
        await until(pg, () => location.hash === '#/review'),
      );
      await go('#/decide');
      await pg.waitForSelector('.kit-list:not([hidden]) .kit-row', {
        timeout: 8000,
      });
      await pg.keyboard.press('g');
      await pg.keyboard.press('w');
      await until(
        pg,
        () => document.querySelector('.cb-walk')?.dataset.step === 'decide',
      );
      await go('#/review');
      check(
        'going elsewhere during it ends it there: the overlay goes, the address stays where Court went',
        (await until(
          pg,
          () => !document.querySelector('.cb-walk:not([hidden])'),
        )) &&
          (await pg.waitForTimeout(500),
          await pg.evaluate(() => location.hash)) === '#/review',
      );
      await go('#/decide');
      await pg.waitForSelector('.kit-list:not([hidden]) .kit-row', {
        timeout: 8000,
      });
    }

    // ---- an empty list: no decision step ----------------------------------
    {
      const all = (
        (await api('GET', '/api/items?view=all&limit=500')).items ?? []
      ).map((it) => it.key);
      if (all.length)
        await api('POST', '/api/decide', { keys: all, disposition: 'keep' });
      await eventually(async () => (await rows()) === 0);
      await pg.waitForTimeout(300);
      await pg.keyboard.press('g');
      await pg.keyboard.press('w');
      const ids = [];
      for (let i = 0; i < IDS.length; i++) {
        if (
          !(await until(
            pg,
            () => !!document.querySelector('.cb-walk:not([hidden])'),
            undefined,
            2000,
          ))
        )
          break;
        await pg.waitForTimeout(200);
        const step = (await state(pg)).step;
        if (!step || ids.at(-1) === step) break;
        ids.push(step);
        await pg.click('.cb-walk-next');
        await until(
          pg,
          (prev) => document.querySelector('.cb-walk')?.dataset.step !== prev,
          step,
          2000,
        );
      }
      checkList(
        'an empty list skips the decision step',
        ids,
        IDS.filter((id) => id !== 'decision'),
      );
    }

    check(
      errors.length
        ? `no console errors — ${errors.slice(0, 3).join(' | ')}`
        : 'no console errors',
      errors.length === 0,
    );
  } catch (err) {
    check(
      `the walkthrough: the scenario ran to its end — ${err.message}`,
      false,
    );
  } finally {
    await context.close();
    s.stop();
    serves = serves.filter((x) => x !== s);
  }
}
