// node_modules/@xterm/addon-fit/lib/addon-fit.mjs
var h = 2;
var _ = 1;
var o = class {
  activate(e) {
    this._terminal = e;
  }
  dispose() {
  }
  fit() {
    let e = this.proposeDimensions();
    if (!e || !this._terminal || isNaN(e.cols) || isNaN(e.rows)) return;
    let t = this._terminal._core;
    (this._terminal.rows !== e.rows || this._terminal.cols !== e.cols) && (t._renderService.clear(), this._terminal.resize(e.cols, e.rows));
  }
  proposeDimensions() {
    if (!this._terminal || !this._terminal.element || !this._terminal.element.parentElement) return;
    let t = this._terminal._core._renderService.dimensions;
    if (t.css.cell.width === 0 || t.css.cell.height === 0) return;
    let s15 = this._terminal.options.scrollback === 0 ? 0 : this._terminal.options.overviewRuler?.width || 14, r5 = window.getComputedStyle(this._terminal.element.parentElement), l = parseInt(r5.getPropertyValue("height")), a = Math.max(0, parseInt(r5.getPropertyValue("width"))), i8 = window.getComputedStyle(this._terminal.element), n = { top: parseInt(i8.getPropertyValue("padding-top")), bottom: parseInt(i8.getPropertyValue("padding-bottom")), right: parseInt(i8.getPropertyValue("padding-right")), left: parseInt(i8.getPropertyValue("padding-left")) }, m = n.top + n.bottom, d = n.right + n.left, c = l - m, p = a - d - s15;
    return { cols: Math.max(h, Math.floor(p / t.css.cell.width)), rows: Math.max(_, Math.floor(c / t.css.cell.height)) };
  }
};

// node_modules/@xterm/addon-unicode11/lib/addon-unicode11.mjs
var ue = [[768, 879], [1155, 1158], [1160, 1161], [1425, 1469], [1471, 1471], [1473, 1474], [1476, 1477], [1479, 1479], [1536, 1539], [1552, 1557], [1611, 1630], [1648, 1648], [1750, 1764], [1767, 1768], [1770, 1773], [1807, 1807], [1809, 1809], [1840, 1866], [1958, 1968], [2027, 2035], [2305, 2306], [2364, 2364], [2369, 2376], [2381, 2381], [2385, 2388], [2402, 2403], [2433, 2433], [2492, 2492], [2497, 2500], [2509, 2509], [2530, 2531], [2561, 2562], [2620, 2620], [2625, 2626], [2631, 2632], [2635, 2637], [2672, 2673], [2689, 2690], [2748, 2748], [2753, 2757], [2759, 2760], [2765, 2765], [2786, 2787], [2817, 2817], [2876, 2876], [2879, 2879], [2881, 2883], [2893, 2893], [2902, 2902], [2946, 2946], [3008, 3008], [3021, 3021], [3134, 3136], [3142, 3144], [3146, 3149], [3157, 3158], [3260, 3260], [3263, 3263], [3270, 3270], [3276, 3277], [3298, 3299], [3393, 3395], [3405, 3405], [3530, 3530], [3538, 3540], [3542, 3542], [3633, 3633], [3636, 3642], [3655, 3662], [3761, 3761], [3764, 3769], [3771, 3772], [3784, 3789], [3864, 3865], [3893, 3893], [3895, 3895], [3897, 3897], [3953, 3966], [3968, 3972], [3974, 3975], [3984, 3991], [3993, 4028], [4038, 4038], [4141, 4144], [4146, 4146], [4150, 4151], [4153, 4153], [4184, 4185], [4448, 4607], [4959, 4959], [5906, 5908], [5938, 5940], [5970, 5971], [6002, 6003], [6068, 6069], [6071, 6077], [6086, 6086], [6089, 6099], [6109, 6109], [6155, 6157], [6313, 6313], [6432, 6434], [6439, 6440], [6450, 6450], [6457, 6459], [6679, 6680], [6912, 6915], [6964, 6964], [6966, 6970], [6972, 6972], [6978, 6978], [7019, 7027], [7616, 7626], [7678, 7679], [8203, 8207], [8234, 8238], [8288, 8291], [8298, 8303], [8400, 8431], [12330, 12335], [12441, 12442], [43014, 43014], [43019, 43019], [43045, 43046], [64286, 64286], [65024, 65039], [65056, 65059], [65279, 65279], [65529, 65531]];
var qe = [[68097, 68099], [68101, 68102], [68108, 68111], [68152, 68154], [68159, 68159], [119143, 119145], [119155, 119170], [119173, 119179], [119210, 119213], [119362, 119364], [917505, 917505], [917536, 917631], [917760, 917999]];
var A;
function He(r5, e) {
  let t = 0, n = e.length - 1, o2;
  if (r5 < e[0][0] || r5 > e[n][1]) return false;
  for (; n >= t; ) if (o2 = t + n >> 1, r5 > e[o2][1]) t = o2 + 1;
  else if (r5 < e[o2][0]) n = o2 - 1;
  else return true;
  return false;
}
var H = class {
  constructor() {
    this.version = "6";
    if (!A) {
      A = new Uint8Array(65536), A.fill(1), A[0] = 0, A.fill(0, 1, 32), A.fill(0, 127, 160), A.fill(2, 4352, 4448), A[9001] = 2, A[9002] = 2, A.fill(2, 11904, 42192), A[12351] = 1, A.fill(2, 44032, 55204), A.fill(2, 63744, 64256), A.fill(2, 65040, 65050), A.fill(2, 65072, 65136), A.fill(2, 65280, 65377), A.fill(2, 65504, 65511);
      for (let e = 0; e < ue.length; ++e) A.fill(0, ue[e][0], ue[e][1] + 1);
    }
  }
  wcwidth(e) {
    return e < 32 ? 0 : e < 127 ? 1 : e < 65536 ? A[e] : He(e, qe) ? 0 : e >= 131072 && e <= 196605 || e >= 196608 && e <= 262141 ? 2 : 1;
  }
  charProperties(e, t) {
    let n = this.wcwidth(e), o2 = n === 0 && t !== 0;
    if (o2) {
      let d = w.extractWidth(t);
      d === 0 ? o2 = false : d > n && (n = d);
    }
    return w.createPropertyValue(0, n, o2);
  }
};
var de = class {
  constructor() {
    this.listeners = [], this.unexpectedErrorHandler = function(e) {
      setTimeout(() => {
        throw e.stack ? J.isErrorNoTelemetry(e) ? new J(e.message + `

` + e.stack) : new Error(e.message + `

` + e.stack) : e;
      }, 0);
    };
  }
  addListener(e) {
    return this.listeners.push(e), () => {
      this._removeListener(e);
    };
  }
  emit(e) {
    this.listeners.forEach((t) => {
      t(e);
    });
  }
  _removeListener(e) {
    this.listeners.splice(this.listeners.indexOf(e), 1);
  }
  setUnexpectedErrorHandler(e) {
    this.unexpectedErrorHandler = e;
  }
  getUnexpectedErrorHandler() {
    return this.unexpectedErrorHandler;
  }
  onUnexpectedError(e) {
    this.unexpectedErrorHandler(e), this.emit(e);
  }
  onUnexpectedExternalError(e) {
    this.unexpectedErrorHandler(e);
  }
};
var Ge = new de();
function Y(r5) {
  Je(r5) || Ge.onUnexpectedError(r5);
}
var ce = "Canceled";
function Je(r5) {
  return r5 instanceof G ? true : r5 instanceof Error && r5.name === ce && r5.message === ce;
}
var G = class extends Error {
  constructor() {
    super(ce), this.name = this.message;
  }
};
var J = class r extends Error {
  constructor(e) {
    super(e), this.name = "CodeExpectedError";
  }
  static fromError(e) {
    if (e instanceof r) return e;
    let t = new r();
    return t.message = e.message, t.stack = e.stack, t;
  }
  static isErrorNoTelemetry(e) {
    return e.name === "CodeExpectedError";
  }
};
function pe(r5, e) {
  let t = this, n = false, o2;
  return function() {
    if (n) return o2;
    if (n = true, e) try {
      o2 = r5.apply(t, arguments);
    } finally {
      e();
    }
    else o2 = r5.apply(t, arguments);
    return o2;
  };
}
function Ye(r5, e, t = 0, n = r5.length) {
  let o2 = t, d = n;
  for (; o2 < d; ) {
    let v2 = Math.floor((o2 + d) / 2);
    e(r5[v2]) ? o2 = v2 + 1 : d = v2;
  }
  return o2 - 1;
}
var X = class X2 {
  constructor(e) {
    this._array = e;
    this._findLastMonotonousLastIdx = 0;
  }
  findLastMonotonous(e) {
    if (X2.assertInvariants) {
      if (this._prevFindLastPredicate) {
        for (let n of this._array) if (this._prevFindLastPredicate(n) && !e(n)) throw new Error("MonotonousArray: current predicate must be weaker than (or equal to) the previous predicate.");
      }
      this._prevFindLastPredicate = e;
    }
    let t = Ye(this._array, e, this._findLastMonotonousLastIdx);
    return this._findLastMonotonousLastIdx = t + 1, t === -1 ? void 0 : this._array[t];
  }
};
X.assertInvariants = false;
var Be;
((E) => {
  function r5(p) {
    return p < 0;
  }
  E.isLessThan = r5;
  function e(p) {
    return p <= 0;
  }
  E.isLessThanOrEqual = e;
  function t(p) {
    return p > 0;
  }
  E.isGreaterThan = t;
  function n(p) {
    return p === 0;
  }
  E.isNeitherLessOrGreaterThan = n, E.greaterThan = 1, E.lessThan = -1, E.neitherLessOrGreaterThan = 0;
})(Be ||= {});
function we(r5, e) {
  return (t, n) => e(r5(t), r5(n));
}
var ke = (r5, e) => r5 - e;
var R = class R2 {
  constructor(e) {
    this.iterate = e;
  }
  forEach(e) {
    this.iterate((t) => (e(t), true));
  }
  toArray() {
    let e = [];
    return this.iterate((t) => (e.push(t), true)), e;
  }
  filter(e) {
    return new R2((t) => this.iterate((n) => e(n) ? t(n) : true));
  }
  map(e) {
    return new R2((t) => this.iterate((n) => t(e(n))));
  }
  some(e) {
    let t = false;
    return this.iterate((n) => (t = e(n), !t)), t;
  }
  findFirst(e) {
    let t;
    return this.iterate((n) => e(n) ? (t = n, false) : true), t;
  }
  findLast(e) {
    let t;
    return this.iterate((n) => (e(n) && (t = n), true)), t;
  }
  findLastMaxBy(e) {
    let t, n = true;
    return this.iterate((o2) => ((n || Be.isGreaterThan(e(o2, t))) && (n = false, t = o2), true)), t;
  }
};
R.empty = new R((e) => {
});
function Oe(r5, e) {
  let t = /* @__PURE__ */ Object.create(null);
  for (let n of r5) {
    let o2 = e(n), d = t[o2];
    d || (d = t[o2] = []), d.push(n);
  }
  return t;
}
var Se;
var Re;
var Le = class {
  constructor(e, t) {
    this.toKey = t;
    this._map = /* @__PURE__ */ new Map();
    this[Se] = "SetWithKey";
    for (let n of e) this.add(n);
  }
  get size() {
    return this._map.size;
  }
  add(e) {
    let t = this.toKey(e);
    return this._map.set(t, e), this;
  }
  delete(e) {
    return this._map.delete(this.toKey(e));
  }
  has(e) {
    return this._map.has(this.toKey(e));
  }
  *entries() {
    for (let e of this._map.values()) yield [e, e];
  }
  keys() {
    return this.values();
  }
  *values() {
    for (let e of this._map.values()) yield e;
  }
  clear() {
    this._map.clear();
  }
  forEach(e, t) {
    this._map.forEach((n) => e.call(t, n, n, this));
  }
  [(Re = Symbol.iterator, Se = Symbol.toStringTag, Re)]() {
    return this.values();
  }
};
var Z = class {
  constructor() {
    this.map = /* @__PURE__ */ new Map();
  }
  add(e, t) {
    let n = this.map.get(e);
    n || (n = /* @__PURE__ */ new Set(), this.map.set(e, n)), n.add(t);
  }
  delete(e, t) {
    let n = this.map.get(e);
    n && (n.delete(t), n.size === 0 && this.map.delete(e));
  }
  forEach(e, t) {
    let n = this.map.get(e);
    n && n.forEach(t);
  }
  get(e) {
    let t = this.map.get(e);
    return t || /* @__PURE__ */ new Set();
  }
};
var fe;
((le2) => {
  function r5(u) {
    return u && typeof u == "object" && typeof u[Symbol.iterator] == "function";
  }
  le2.is = r5;
  let e = Object.freeze([]);
  function t() {
    return e;
  }
  le2.empty = t;
  function* n(u) {
    yield u;
  }
  le2.single = n;
  function o2(u) {
    return r5(u) ? u : n(u);
  }
  le2.wrap = o2;
  function d(u) {
    return u || e;
  }
  le2.from = d;
  function* v2(u) {
    for (let f = u.length - 1; f >= 0; f--) yield u[f];
  }
  le2.reverse = v2;
  function E(u) {
    return !u || u[Symbol.iterator]().next().done === true;
  }
  le2.isEmpty = E;
  function p(u) {
    return u[Symbol.iterator]().next().value;
  }
  le2.first = p;
  function b2(u, f) {
    let m = 0;
    for (let g of u) if (f(g, m++)) return true;
    return false;
  }
  le2.some = b2;
  function D3(u, f) {
    for (let m of u) if (f(m)) return m;
  }
  le2.find = D3;
  function* T(u, f) {
    for (let m of u) f(m) && (yield m);
  }
  le2.filter = T;
  function* B3(u, f) {
    let m = 0;
    for (let g of u) yield f(g, m++);
  }
  le2.map = B3;
  function* L2(u, f) {
    let m = 0;
    for (let g of u) yield* f(g, m++);
  }
  le2.flatMap = L2;
  function* oe(...u) {
    for (let f of u) yield* f;
  }
  le2.concat = oe;
  function z3(u, f, m) {
    let g = m;
    for (let W of u) g = f(g, W);
    return g;
  }
  le2.reduce = z3;
  function* k2(u, f, m = u.length) {
    for (f < 0 && (f += u.length), m < 0 ? m += u.length : m > u.length && (m = u.length); f < m; f++) yield u[f];
  }
  le2.slice = k2;
  function ae3(u, f = Number.POSITIVE_INFINITY) {
    let m = [];
    if (f === 0) return [m, u];
    let g = u[Symbol.iterator]();
    for (let W = 0; W < f; W++) {
      let xe3 = g.next();
      if (xe3.done) return [m, le2.empty()];
      m.push(xe3.value);
    }
    return [m, { [Symbol.iterator]() {
      return g;
    } }];
  }
  le2.consume = ae3;
  async function V2(u) {
    let f = [];
    for await (let m of u) f.push(m);
    return Promise.resolve(f);
  }
  le2.asyncToArray = V2;
})(fe ||= {});
var Xe = false;
var O = null;
var ee = class ee2 {
  constructor() {
    this.livingDisposables = /* @__PURE__ */ new Map();
  }
  getDisposableData(e) {
    let t = this.livingDisposables.get(e);
    return t || (t = { parent: null, source: null, isSingleton: false, value: e, idx: ee2.idx++ }, this.livingDisposables.set(e, t)), t;
  }
  trackDisposable(e) {
    let t = this.getDisposableData(e);
    t.source || (t.source = new Error().stack);
  }
  setParent(e, t) {
    let n = this.getDisposableData(e);
    n.parent = t;
  }
  markAsDisposed(e) {
    this.livingDisposables.delete(e);
  }
  markAsSingleton(e) {
    this.getDisposableData(e).isSingleton = true;
  }
  getRootParent(e, t) {
    let n = t.get(e);
    if (n) return n;
    let o2 = e.parent ? this.getRootParent(this.getDisposableData(e.parent), t) : e;
    return t.set(e, o2), o2;
  }
  getTrackedDisposables() {
    let e = /* @__PURE__ */ new Map();
    return [...this.livingDisposables.entries()].filter(([, n]) => n.source !== null && !this.getRootParent(n, e).isSingleton).flatMap(([n]) => n);
  }
  computeLeakingDisposables(e = 10, t) {
    let n;
    if (t) n = t;
    else {
      let p = /* @__PURE__ */ new Map(), b2 = [...this.livingDisposables.values()].filter((T) => T.source !== null && !this.getRootParent(T, p).isSingleton);
      if (b2.length === 0) return;
      let D3 = new Set(b2.map((T) => T.value));
      if (n = b2.filter((T) => !(T.parent && D3.has(T.parent))), n.length === 0) throw new Error("There are cyclic diposable chains!");
    }
    if (!n) return;
    function o2(p) {
      function b2(T, B3) {
        for (; T.length > 0 && B3.some((L2) => typeof L2 == "string" ? L2 === T[0] : T[0].match(L2)); ) T.shift();
      }
      let D3 = p.source.split(`
`).map((T) => T.trim().replace("at ", "")).filter((T) => T !== "");
      return b2(D3, ["Error", /^trackDisposable \(.*\)$/, /^DisposableTracker.trackDisposable \(.*\)$/]), D3.reverse();
    }
    let d = new Z();
    for (let p of n) {
      let b2 = o2(p);
      for (let D3 = 0; D3 <= b2.length; D3++) d.add(b2.slice(0, D3).join(`
`), p);
    }
    n.sort(we((p) => p.idx, ke));
    let v2 = "", E = 0;
    for (let p of n.slice(0, e)) {
      E++;
      let b2 = o2(p), D3 = [];
      for (let T = 0; T < b2.length; T++) {
        let B3 = b2[T];
        B3 = `(shared with ${d.get(b2.slice(0, T + 1).join(`
`)).size}/${n.length} leaks) at ${B3}`;
        let oe = d.get(b2.slice(0, T).join(`
`)), z3 = Oe([...oe].map((k2) => o2(k2)[T]), (k2) => k2);
        delete z3[b2[T]];
        for (let [k2, ae3] of Object.entries(z3)) D3.unshift(`    - stacktraces of ${ae3.length} other leaks continue with ${k2}`);
        D3.unshift(B3);
      }
      v2 += `


==================== Leaking disposable ${E}/${n.length}: ${p.value.constructor.name} ====================
${D3.join(`
`)}
============================================================

`;
    }
    return n.length > e && (v2 += `


... and ${n.length - e} more leaking disposables

`), { leaks: n, details: v2 };
  }
};
ee.idx = 0;
function Ze(r5) {
  O = r5;
}
if (Xe) {
  let r5 = "__is_disposable_tracked__";
  Ze(new class {
    trackDisposable(e) {
      let t = new Error("Potentially leaked disposable").stack;
      setTimeout(() => {
        e[r5] || console.log(t);
      }, 3e3);
    }
    setParent(e, t) {
      if (e && e !== _2.None) try {
        e[r5] = true;
      } catch {
      }
    }
    markAsDisposed(e) {
      if (e && e !== _2.None) try {
        e[r5] = true;
      } catch {
      }
    }
    markAsSingleton(e) {
    }
  }());
}
function Te(r5) {
  return O?.trackDisposable(r5), r5;
}
function ve(r5) {
  O?.markAsDisposed(r5);
}
function he(r5, e) {
  O?.setParent(r5, e);
}
function et(r5, e) {
  if (O) for (let t of r5) O.setParent(t, e);
}
function Pe(r5) {
  if (fe.is(r5)) {
    let e = [];
    for (let t of r5) if (t) try {
      t.dispose();
    } catch (n) {
      e.push(n);
    }
    if (e.length === 1) throw e[0];
    if (e.length > 1) throw new AggregateError(e, "Encountered errors while disposing of store");
    return Array.isArray(r5) ? [] : r5;
  } else if (r5) return r5.dispose(), r5;
}
function Me(...r5) {
  let e = me(() => Pe(r5));
  return et(r5, e), e;
}
function me(r5) {
  let e = Te({ dispose: pe(() => {
    ve(e), r5();
  }) });
  return e;
}
var te = class te2 {
  constructor() {
    this._toDispose = /* @__PURE__ */ new Set();
    this._isDisposed = false;
    Te(this);
  }
  dispose() {
    this._isDisposed || (ve(this), this._isDisposed = true, this.clear());
  }
  get isDisposed() {
    return this._isDisposed;
  }
  clear() {
    if (this._toDispose.size !== 0) try {
      Pe(this._toDispose);
    } finally {
      this._toDispose.clear();
    }
  }
  add(e) {
    if (!e) return e;
    if (e === this) throw new Error("Cannot register a disposable on itself!");
    return he(e, this), this._isDisposed ? te2.DISABLE_DISPOSED_WARNING || console.warn(new Error("Trying to add a disposable to a DisposableStore that has already been disposed of. The added object will be leaked!").stack) : this._toDispose.add(e), e;
  }
  delete(e) {
    if (e) {
      if (e === this) throw new Error("Cannot dispose a disposable on itself!");
      this._toDispose.delete(e), e.dispose();
    }
  }
  deleteAndLeak(e) {
    e && this._toDispose.has(e) && (this._toDispose.delete(e), he(e, null));
  }
};
te.DISABLE_DISPOSED_WARNING = false;
var U = te;
var _2 = class {
  constructor() {
    this._store = new U();
    Te(this), he(this._store, this);
  }
  dispose() {
    ve(this), this._store.dispose();
  }
  _register(e) {
    if (e === this) throw new Error("Cannot register a disposable on itself!");
    return this._store.add(e);
  }
};
_2.None = Object.freeze({ dispose() {
} });
var P = class P2 {
  constructor(e) {
    this.element = e, this.next = P2.Undefined, this.prev = P2.Undefined;
  }
};
P.Undefined = new P(void 0);
var tt = globalThis.performance && typeof globalThis.performance.now == "function";
var ne = class r2 {
  static create(e) {
    return new r2(e);
  }
  constructor(e) {
    this._now = tt && e === false ? Date.now : globalThis.performance.now.bind(globalThis.performance), this._startTime = this._now(), this._stopTime = -1;
  }
  stop() {
    this._stopTime = this._now();
  }
  reset() {
    this._startTime = this._now(), this._stopTime = -1;
  }
  elapsed() {
    return this._stopTime !== -1 ? this._stopTime - this._startTime : this._now() - this._startTime;
  }
};
var nt = false;
var Ve = false;
var rt = false;
var it;
((Q5) => {
  Q5.None = () => _2.None;
  function e(l) {
    if (rt) {
      let { onDidAddListener: i8 } = l, a = K.create(), s15 = 0;
      l.onDidAddListener = () => {
        ++s15 === 2 && (console.warn("snapshotted emitter LIKELY used public and SHOULD HAVE BEEN created with DisposableStore. snapshotted here"), a.print()), i8?.();
      };
    }
  }
  function t(l, i8) {
    return B3(l, () => {
    }, 0, void 0, true, void 0, i8);
  }
  Q5.defer = t;
  function n(l) {
    return (i8, a = null, s15) => {
      let x = false, c;
      return c = l((h2) => {
        if (!x) return c ? c.dispose() : x = true, i8.call(a, h2);
      }, null, s15), x && c.dispose(), c;
    };
  }
  Q5.once = n;
  function o2(l, i8, a) {
    return D3((s15, x = null, c) => l((h2) => s15.call(x, i8(h2)), null, c), a);
  }
  Q5.map = o2;
  function d(l, i8, a) {
    return D3((s15, x = null, c) => l((h2) => {
      i8(h2), s15.call(x, h2);
    }, null, c), a);
  }
  Q5.forEach = d;
  function v2(l, i8, a) {
    return D3((s15, x = null, c) => l((h2) => i8(h2) && s15.call(x, h2), null, c), a);
  }
  Q5.filter = v2;
  function E(l) {
    return l;
  }
  Q5.signal = E;
  function p(...l) {
    return (i8, a = null, s15) => {
      let x = Me(...l.map((c) => c((h2) => i8.call(a, h2))));
      return T(x, s15);
    };
  }
  Q5.any = p;
  function b2(l, i8, a, s15) {
    let x = a;
    return o2(l, (c) => (x = i8(x, c), x), s15);
  }
  Q5.reduce = b2;
  function D3(l, i8) {
    let a, s15 = { onWillAddFirstListener() {
      a = l(x.fire, x);
    }, onDidRemoveLastListener() {
      a?.dispose();
    } };
    i8 || e(s15);
    let x = new C(s15);
    return i8?.add(x), x.event;
  }
  function T(l, i8) {
    return i8 instanceof Array ? i8.push(l) : i8 && i8.add(l), l;
  }
  function B3(l, i8, a = 100, s15 = false, x = false, c, h2) {
    let F3, y, S2, $2 = 0, j3, Ce3 = { leakWarningThreshold: c, onWillAddFirstListener() {
      F3 = l((Qe2) => {
        $2++, y = i8(y, Qe2), s15 && !S2 && (q2.fire(y), y = void 0), j3 = () => {
          let $e2 = y;
          y = void 0, S2 = void 0, (!s15 || $2 > 1) && q2.fire($e2), $2 = 0;
        }, typeof a == "number" ? (clearTimeout(S2), S2 = setTimeout(j3, a)) : S2 === void 0 && (S2 = 0, queueMicrotask(j3));
      });
    }, onWillRemoveListener() {
      x && $2 > 0 && j3?.();
    }, onDidRemoveLastListener() {
      j3 = void 0, F3.dispose();
    } };
    h2 || e(Ce3);
    let q2 = new C(Ce3);
    return h2?.add(q2), q2.event;
  }
  Q5.debounce = B3;
  function L2(l, i8 = 0, a) {
    return Q5.debounce(l, (s15, x) => s15 ? (s15.push(x), s15) : [x], i8, void 0, true, void 0, a);
  }
  Q5.accumulate = L2;
  function oe(l, i8 = (s15, x) => s15 === x, a) {
    let s15 = true, x;
    return v2(l, (c) => {
      let h2 = s15 || !i8(c, x);
      return s15 = false, x = c, h2;
    }, a);
  }
  Q5.latch = oe;
  function z3(l, i8, a) {
    return [Q5.filter(l, i8, a), Q5.filter(l, (s15) => !i8(s15), a)];
  }
  Q5.split = z3;
  function k2(l, i8 = false, a = [], s15) {
    let x = a.slice(), c = l((y) => {
      x ? x.push(y) : F3.fire(y);
    });
    s15 && s15.add(c);
    let h2 = () => {
      x?.forEach((y) => F3.fire(y)), x = null;
    }, F3 = new C({ onWillAddFirstListener() {
      c || (c = l((y) => F3.fire(y)), s15 && s15.add(c));
    }, onDidAddFirstListener() {
      x && (i8 ? setTimeout(h2) : h2());
    }, onDidRemoveLastListener() {
      c && c.dispose(), c = null;
    } });
    return s15 && s15.add(F3), F3.event;
  }
  Q5.buffer = k2;
  function ae3(l, i8) {
    return (s15, x, c) => {
      let h2 = i8(new le2());
      return l(function(F3) {
        let y = h2.evaluate(F3);
        y !== V2 && s15.call(x, y);
      }, void 0, c);
    };
  }
  Q5.chain = ae3;
  let V2 = /* @__PURE__ */ Symbol("HaltChainable");
  class le2 {
    constructor() {
      this.steps = [];
    }
    map(i8) {
      return this.steps.push(i8), this;
    }
    forEach(i8) {
      return this.steps.push((a) => (i8(a), a)), this;
    }
    filter(i8) {
      return this.steps.push((a) => i8(a) ? a : V2), this;
    }
    reduce(i8, a) {
      let s15 = a;
      return this.steps.push((x) => (s15 = i8(s15, x), s15)), this;
    }
    latch(i8 = (a, s15) => a === s15) {
      let a = true, s15;
      return this.steps.push((x) => {
        let c = a || !i8(x, s15);
        return a = false, s15 = x, c ? x : V2;
      }), this;
    }
    evaluate(i8) {
      for (let a of this.steps) if (i8 = a(i8), i8 === V2) break;
      return i8;
    }
  }
  function u(l, i8, a = (s15) => s15) {
    let s15 = (...F3) => h2.fire(a(...F3)), x = () => l.on(i8, s15), c = () => l.removeListener(i8, s15), h2 = new C({ onWillAddFirstListener: x, onDidRemoveLastListener: c });
    return h2.event;
  }
  Q5.fromNodeEventEmitter = u;
  function f(l, i8, a = (s15) => s15) {
    let s15 = (...F3) => h2.fire(a(...F3)), x = () => l.addEventListener(i8, s15), c = () => l.removeEventListener(i8, s15), h2 = new C({ onWillAddFirstListener: x, onDidRemoveLastListener: c });
    return h2.event;
  }
  Q5.fromDOMEventEmitter = f;
  function m(l) {
    return new Promise((i8) => n(l)(i8));
  }
  Q5.toPromise = m;
  function g(l) {
    let i8 = new C();
    return l.then((a) => {
      i8.fire(a);
    }, () => {
      i8.fire(void 0);
    }).finally(() => {
      i8.dispose();
    }), i8.event;
  }
  Q5.fromPromise = g;
  function W(l, i8) {
    return l((a) => i8.fire(a));
  }
  Q5.forward = W;
  function xe3(l, i8, a) {
    return i8(a), l((s15) => i8(s15));
  }
  Q5.runAndSubscribe = xe3;
  class ze2 {
    constructor(i8, a) {
      this._observable = i8;
      this._counter = 0;
      this._hasChanged = false;
      let s15 = { onWillAddFirstListener: () => {
        i8.addObserver(this);
      }, onDidRemoveLastListener: () => {
        i8.removeObserver(this);
      } };
      a || e(s15), this.emitter = new C(s15), a && a.add(this.emitter);
    }
    beginUpdate(i8) {
      this._counter++;
    }
    handlePossibleChange(i8) {
    }
    handleChange(i8, a) {
      this._hasChanged = true;
    }
    endUpdate(i8) {
      this._counter--, this._counter === 0 && (this._observable.reportChanges(), this._hasChanged && (this._hasChanged = false, this.emitter.fire(this._observable.get())));
    }
  }
  function ut3(l, i8) {
    return new ze2(l, i8).emitter.event;
  }
  Q5.fromObservable = ut3;
  function dt3(l) {
    return (i8, a, s15) => {
      let x = 0, c = false, h2 = { beginUpdate() {
        x++;
      }, endUpdate() {
        x--, x === 0 && (l.reportChanges(), c && (c = false, i8.call(a)));
      }, handlePossibleChange() {
      }, handleChange() {
        c = true;
      } };
      l.addObserver(h2), l.reportChanges();
      let F3 = { dispose() {
        l.removeObserver(h2);
      } };
      return s15 instanceof U ? s15.add(F3) : Array.isArray(s15) && s15.push(F3), F3;
    };
  }
  Q5.fromObservableLight = dt3;
})(it ||= {});
var M = class M2 {
  constructor(e) {
    this.listenerCount = 0;
    this.invocationCount = 0;
    this.elapsedOverall = 0;
    this.durations = [];
    this.name = `${e}_${M2._idPool++}`, M2.all.add(this);
  }
  start(e) {
    this._stopWatch = new ne(), this.listenerCount = e;
  }
  stop() {
    if (this._stopWatch) {
      let e = this._stopWatch.elapsed();
      this.durations.push(e), this.elapsedOverall += e, this.invocationCount += 1, this._stopWatch = void 0;
    }
  }
};
M.all = /* @__PURE__ */ new Set(), M._idPool = 0;
var be = M;
var We = -1;
var ie = class ie2 {
  constructor(e, t, n = (ie2._idPool++).toString(16).padStart(3, "0")) {
    this._errorHandler = e;
    this.threshold = t;
    this.name = n;
    this._warnCountdown = 0;
  }
  dispose() {
    this._stacks?.clear();
  }
  check(e, t) {
    let n = this.threshold;
    if (n <= 0 || t < n) return;
    this._stacks || (this._stacks = /* @__PURE__ */ new Map());
    let o2 = this._stacks.get(e.value) || 0;
    if (this._stacks.set(e.value, o2 + 1), this._warnCountdown -= 1, this._warnCountdown <= 0) {
      this._warnCountdown = n * 0.5;
      let [d, v2] = this.getMostFrequentStack(), E = `[${this.name}] potential listener LEAK detected, having ${t} listeners already. MOST frequent listener (${v2}):`;
      console.warn(E), console.warn(d);
      let p = new De(E, d);
      this._errorHandler(p);
    }
    return () => {
      let d = this._stacks.get(e.value) || 0;
      this._stacks.set(e.value, d - 1);
    };
  }
  getMostFrequentStack() {
    if (!this._stacks) return;
    let e, t = 0;
    for (let [n, o2] of this._stacks) (!e || t < o2) && (e = [n, o2], t = o2);
    return e;
  }
};
ie._idPool = 1;
var Ee = ie;
var K = class r3 {
  constructor(e) {
    this.value = e;
  }
  static create() {
    let e = new Error();
    return new r3(e.stack ?? "");
  }
  print() {
    console.warn(this.value.split(`
`).slice(2).join(`
`));
  }
};
var De = class extends Error {
  constructor(e, t) {
    super(e), this.name = "ListenerLeakError", this.stack = t;
  }
};
var Ae = class extends Error {
  constructor(e, t) {
    super(e), this.name = "ListenerRefusalError", this.stack = t;
  }
};
var st = 0;
var N = class {
  constructor(e) {
    this.value = e;
    this.id = st++;
  }
};
var ot = 2;
var at = (r5, e) => {
  if (r5 instanceof N) e(r5);
  else for (let t = 0; t < r5.length; t++) {
    let n = r5[t];
    n && e(n);
  }
};
var re;
if (nt) {
  let r5 = [];
  setInterval(() => {
    r5.length !== 0 && (console.warn("[LEAKING LISTENERS] GC'ed these listeners that were NOT yet disposed:"), console.warn(r5.join(`
`)), r5.length = 0);
  }, 3e3), re = new FinalizationRegistry((e) => {
    typeof e == "string" && r5.push(e);
  });
}
var C = class {
  constructor(e) {
    this._size = 0;
    this._options = e, this._leakageMon = We > 0 || this._options?.leakWarningThreshold ? new Ee(e?.onListenerError ?? Y, this._options?.leakWarningThreshold ?? We) : void 0, this._perfMon = this._options?._profName ? new be(this._options._profName) : void 0, this._deliveryQueue = this._options?.deliveryQueue;
  }
  dispose() {
    if (!this._disposed) {
      if (this._disposed = true, this._deliveryQueue?.current === this && this._deliveryQueue.reset(), this._listeners) {
        if (Ve) {
          let e = this._listeners;
          queueMicrotask(() => {
            at(e, (t) => t.stack?.print());
          });
        }
        this._listeners = void 0, this._size = 0;
      }
      this._options?.onDidRemoveLastListener?.(), this._leakageMon?.dispose();
    }
  }
  get event() {
    return this._event ??= (e, t, n) => {
      if (this._leakageMon && this._size > this._leakageMon.threshold ** 2) {
        let p = `[${this._leakageMon.name}] REFUSES to accept new listeners because it exceeded its threshold by far (${this._size} vs ${this._leakageMon.threshold})`;
        console.warn(p);
        let b2 = this._leakageMon.getMostFrequentStack() ?? ["UNKNOWN stack", -1], D3 = new Ae(`${p}. HINT: Stack shows most frequent listener (${b2[1]}-times)`, b2[0]);
        return (this._options?.onListenerError || Y)(D3), _2.None;
      }
      if (this._disposed) return _2.None;
      t && (e = e.bind(t));
      let o2 = new N(e), d, v2;
      this._leakageMon && this._size >= Math.ceil(this._leakageMon.threshold * 0.2) && (o2.stack = K.create(), d = this._leakageMon.check(o2.stack, this._size + 1)), Ve && (o2.stack = v2 ?? K.create()), this._listeners ? this._listeners instanceof N ? (this._deliveryQueue ??= new Fe(), this._listeners = [this._listeners, o2]) : this._listeners.push(o2) : (this._options?.onWillAddFirstListener?.(this), this._listeners = o2, this._options?.onDidAddFirstListener?.(this)), this._size++;
      let E = me(() => {
        re?.unregister(E), d?.(), this._removeListener(o2);
      });
      if (n instanceof U ? n.add(E) : Array.isArray(n) && n.push(E), re) {
        let p = new Error().stack.split(`
`).slice(2, 3).join(`
`).trim(), b2 = /(file:|vscode-file:\/\/vscode-app)?(\/[^:]*:\d+:\d+)/.exec(p);
        re.register(E, b2?.[2] ?? p, E);
      }
      return E;
    }, this._event;
  }
  _removeListener(e) {
    if (this._options?.onWillRemoveListener?.(this), !this._listeners) return;
    if (this._size === 1) {
      this._listeners = void 0, this._options?.onDidRemoveLastListener?.(this), this._size = 0;
      return;
    }
    let t = this._listeners, n = t.indexOf(e);
    if (n === -1) throw console.log("disposed?", this._disposed), console.log("size?", this._size), console.log("arr?", JSON.stringify(this._listeners)), new Error("Attempted to dispose unknown listener");
    this._size--, t[n] = void 0;
    let o2 = this._deliveryQueue.current === this;
    if (this._size * ot <= t.length) {
      let d = 0;
      for (let v2 = 0; v2 < t.length; v2++) t[v2] ? t[d++] = t[v2] : o2 && (this._deliveryQueue.end--, d < this._deliveryQueue.i && this._deliveryQueue.i--);
      t.length = d;
    }
  }
  _deliver(e, t) {
    if (!e) return;
    let n = this._options?.onListenerError || Y;
    if (!n) {
      e.value(t);
      return;
    }
    try {
      e.value(t);
    } catch (o2) {
      n(o2);
    }
  }
  _deliverQueue(e) {
    let t = e.current._listeners;
    for (; e.i < e.end; ) this._deliver(t[e.i++], e.value);
    e.reset();
  }
  fire(e) {
    if (this._deliveryQueue?.current && (this._deliverQueue(this._deliveryQueue), this._perfMon?.stop()), this._perfMon?.start(this._size), this._listeners) if (this._listeners instanceof N) this._deliver(this._listeners, e);
    else {
      let t = this._deliveryQueue;
      t.enqueue(this, e, this._listeners.length), this._deliverQueue(t);
    }
    this._perfMon?.stop();
  }
  hasListeners() {
    return this._size > 0;
  }
};
var Fe = class {
  constructor() {
    this.i = -1;
    this.end = 0;
  }
  enqueue(e, t, n) {
    this.i = 0, this.end = n, this.current = e, this.value = t;
  }
  reset() {
    this.i = this.end, this.current = void 0, this.value = void 0;
  }
};
var w = class r4 {
  constructor() {
    this._providers = /* @__PURE__ */ Object.create(null);
    this._active = "";
    this._onChange = new C();
    this.onChange = this._onChange.event;
    let e = new H();
    this.register(e), this._active = e.version, this._activeProvider = e;
  }
  static extractShouldJoin(e) {
    return (e & 1) !== 0;
  }
  static extractWidth(e) {
    return e >> 1 & 3;
  }
  static extractCharKind(e) {
    return e >> 3;
  }
  static createPropertyValue(e, t, n = false) {
    return (e & 16777215) << 3 | (t & 3) << 1 | (n ? 1 : 0);
  }
  dispose() {
    this._onChange.dispose();
  }
  get versions() {
    return Object.keys(this._providers);
  }
  get activeVersion() {
    return this._active;
  }
  set activeVersion(e) {
    if (!this._providers[e]) throw new Error(`unknown Unicode version "${e}"`);
    this._active = e, this._activeProvider = this._providers[e], this._onChange.fire(e);
  }
  register(e) {
    this._providers[e.version] = e;
  }
  wcwidth(e) {
    return this._activeProvider.wcwidth(e);
  }
  getStringCellWidth(e) {
    let t = 0, n = 0, o2 = e.length;
    for (let d = 0; d < o2; ++d) {
      let v2 = e.charCodeAt(d);
      if (55296 <= v2 && v2 <= 56319) {
        if (++d >= o2) return t + this.wcwidth(v2);
        let b2 = e.charCodeAt(d);
        56320 <= b2 && b2 <= 57343 ? v2 = (v2 - 55296) * 1024 + b2 - 56320 + 65536 : t += this.wcwidth(b2);
      }
      let E = this.charProperties(v2, n), p = r4.extractWidth(E);
      r4.extractShouldJoin(E) && (p -= r4.extractWidth(n)), t += p, n = E;
    }
    return t;
  }
  charProperties(e, t) {
    return this._activeProvider.charProperties(e, t);
  }
};
var ye = [[768, 879], [1155, 1161], [1425, 1469], [1471, 1471], [1473, 1474], [1476, 1477], [1479, 1479], [1536, 1541], [1552, 1562], [1564, 1564], [1611, 1631], [1648, 1648], [1750, 1757], [1759, 1764], [1767, 1768], [1770, 1773], [1807, 1807], [1809, 1809], [1840, 1866], [1958, 1968], [2027, 2035], [2045, 2045], [2070, 2073], [2075, 2083], [2085, 2087], [2089, 2093], [2137, 2139], [2259, 2306], [2362, 2362], [2364, 2364], [2369, 2376], [2381, 2381], [2385, 2391], [2402, 2403], [2433, 2433], [2492, 2492], [2497, 2500], [2509, 2509], [2530, 2531], [2558, 2558], [2561, 2562], [2620, 2620], [2625, 2626], [2631, 2632], [2635, 2637], [2641, 2641], [2672, 2673], [2677, 2677], [2689, 2690], [2748, 2748], [2753, 2757], [2759, 2760], [2765, 2765], [2786, 2787], [2810, 2815], [2817, 2817], [2876, 2876], [2879, 2879], [2881, 2884], [2893, 2893], [2902, 2902], [2914, 2915], [2946, 2946], [3008, 3008], [3021, 3021], [3072, 3072], [3076, 3076], [3134, 3136], [3142, 3144], [3146, 3149], [3157, 3158], [3170, 3171], [3201, 3201], [3260, 3260], [3263, 3263], [3270, 3270], [3276, 3277], [3298, 3299], [3328, 3329], [3387, 3388], [3393, 3396], [3405, 3405], [3426, 3427], [3530, 3530], [3538, 3540], [3542, 3542], [3633, 3633], [3636, 3642], [3655, 3662], [3761, 3761], [3764, 3772], [3784, 3789], [3864, 3865], [3893, 3893], [3895, 3895], [3897, 3897], [3953, 3966], [3968, 3972], [3974, 3975], [3981, 3991], [3993, 4028], [4038, 4038], [4141, 4144], [4146, 4151], [4153, 4154], [4157, 4158], [4184, 4185], [4190, 4192], [4209, 4212], [4226, 4226], [4229, 4230], [4237, 4237], [4253, 4253], [4448, 4607], [4957, 4959], [5906, 5908], [5938, 5940], [5970, 5971], [6002, 6003], [6068, 6069], [6071, 6077], [6086, 6086], [6089, 6099], [6109, 6109], [6155, 6158], [6277, 6278], [6313, 6313], [6432, 6434], [6439, 6440], [6450, 6450], [6457, 6459], [6679, 6680], [6683, 6683], [6742, 6742], [6744, 6750], [6752, 6752], [6754, 6754], [6757, 6764], [6771, 6780], [6783, 6783], [6832, 6846], [6912, 6915], [6964, 6964], [6966, 6970], [6972, 6972], [6978, 6978], [7019, 7027], [7040, 7041], [7074, 7077], [7080, 7081], [7083, 7085], [7142, 7142], [7144, 7145], [7149, 7149], [7151, 7153], [7212, 7219], [7222, 7223], [7376, 7378], [7380, 7392], [7394, 7400], [7405, 7405], [7412, 7412], [7416, 7417], [7616, 7673], [7675, 7679], [8203, 8207], [8234, 8238], [8288, 8292], [8294, 8303], [8400, 8432], [11503, 11505], [11647, 11647], [11744, 11775], [12330, 12333], [12441, 12442], [42607, 42610], [42612, 42621], [42654, 42655], [42736, 42737], [43010, 43010], [43014, 43014], [43019, 43019], [43045, 43046], [43204, 43205], [43232, 43249], [43263, 43263], [43302, 43309], [43335, 43345], [43392, 43394], [43443, 43443], [43446, 43449], [43452, 43453], [43493, 43493], [43561, 43566], [43569, 43570], [43573, 43574], [43587, 43587], [43596, 43596], [43644, 43644], [43696, 43696], [43698, 43700], [43703, 43704], [43710, 43711], [43713, 43713], [43756, 43757], [43766, 43766], [44005, 44005], [44008, 44008], [44013, 44013], [64286, 64286], [65024, 65039], [65056, 65071], [65279, 65279], [65529, 65531]];
var lt = [[66045, 66045], [66272, 66272], [66422, 66426], [68097, 68099], [68101, 68102], [68108, 68111], [68152, 68154], [68159, 68159], [68325, 68326], [68900, 68903], [69446, 69456], [69633, 69633], [69688, 69702], [69759, 69761], [69811, 69814], [69817, 69818], [69821, 69821], [69837, 69837], [69888, 69890], [69927, 69931], [69933, 69940], [70003, 70003], [70016, 70017], [70070, 70078], [70089, 70092], [70191, 70193], [70196, 70196], [70198, 70199], [70206, 70206], [70367, 70367], [70371, 70378], [70400, 70401], [70459, 70460], [70464, 70464], [70502, 70508], [70512, 70516], [70712, 70719], [70722, 70724], [70726, 70726], [70750, 70750], [70835, 70840], [70842, 70842], [70847, 70848], [70850, 70851], [71090, 71093], [71100, 71101], [71103, 71104], [71132, 71133], [71219, 71226], [71229, 71229], [71231, 71232], [71339, 71339], [71341, 71341], [71344, 71349], [71351, 71351], [71453, 71455], [71458, 71461], [71463, 71467], [71727, 71735], [71737, 71738], [72148, 72151], [72154, 72155], [72160, 72160], [72193, 72202], [72243, 72248], [72251, 72254], [72263, 72263], [72273, 72278], [72281, 72283], [72330, 72342], [72344, 72345], [72752, 72758], [72760, 72765], [72767, 72767], [72850, 72871], [72874, 72880], [72882, 72883], [72885, 72886], [73009, 73014], [73018, 73018], [73020, 73021], [73023, 73029], [73031, 73031], [73104, 73105], [73109, 73109], [73111, 73111], [73459, 73460], [78896, 78904], [92912, 92916], [92976, 92982], [94031, 94031], [94095, 94098], [113821, 113822], [113824, 113827], [119143, 119145], [119155, 119170], [119173, 119179], [119210, 119213], [119362, 119364], [121344, 121398], [121403, 121452], [121461, 121461], [121476, 121476], [121499, 121503], [121505, 121519], [122880, 122886], [122888, 122904], [122907, 122913], [122915, 122916], [122918, 122922], [123184, 123190], [123628, 123631], [125136, 125142], [125252, 125258], [917505, 917505], [917536, 917631], [917760, 917999]];
var ge = [[4352, 4447], [8986, 8987], [9001, 9002], [9193, 9196], [9200, 9200], [9203, 9203], [9725, 9726], [9748, 9749], [9800, 9811], [9855, 9855], [9875, 9875], [9889, 9889], [9898, 9899], [9917, 9918], [9924, 9925], [9934, 9934], [9940, 9940], [9962, 9962], [9970, 9971], [9973, 9973], [9978, 9978], [9981, 9981], [9989, 9989], [9994, 9995], [10024, 10024], [10060, 10060], [10062, 10062], [10067, 10069], [10071, 10071], [10133, 10135], [10160, 10160], [10175, 10175], [11035, 11036], [11088, 11088], [11093, 11093], [11904, 11929], [11931, 12019], [12032, 12245], [12272, 12283], [12288, 12329], [12334, 12350], [12353, 12438], [12443, 12543], [12549, 12591], [12593, 12686], [12688, 12730], [12736, 12771], [12784, 12830], [12832, 12871], [12880, 19903], [19968, 42124], [42128, 42182], [43360, 43388], [44032, 55203], [63744, 64255], [65040, 65049], [65072, 65106], [65108, 65126], [65128, 65131], [65281, 65376], [65504, 65510]];
var xt = [[94176, 94179], [94208, 100343], [100352, 101106], [110592, 110878], [110928, 110930], [110948, 110951], [110960, 111355], [126980, 126980], [127183, 127183], [127374, 127374], [127377, 127386], [127488, 127490], [127504, 127547], [127552, 127560], [127568, 127569], [127584, 127589], [127744, 127776], [127789, 127797], [127799, 127868], [127870, 127891], [127904, 127946], [127951, 127955], [127968, 127984], [127988, 127988], [127992, 128062], [128064, 128064], [128066, 128252], [128255, 128317], [128331, 128334], [128336, 128359], [128378, 128378], [128405, 128406], [128420, 128420], [128507, 128591], [128640, 128709], [128716, 128716], [128720, 128722], [128725, 128725], [128747, 128748], [128756, 128762], [128992, 129003], [129293, 129393], [129395, 129398], [129402, 129442], [129445, 129450], [129454, 129482], [129485, 129535], [129648, 129651], [129656, 129658], [129664, 129666], [129680, 129685], [131072, 196605], [196608, 262141]];
var I;
function je(r5, e) {
  let t = 0, n = e.length - 1, o2;
  if (r5 < e[0][0] || r5 > e[n][1]) return false;
  for (; n >= t; ) if (o2 = t + n >> 1, r5 > e[o2][1]) t = o2 + 1;
  else if (r5 < e[o2][0]) n = o2 - 1;
  else return true;
  return false;
}
var se = class {
  constructor() {
    this.version = "11";
    if (!I) {
      I = new Uint8Array(65536), I.fill(1), I[0] = 0, I.fill(0, 1, 32), I.fill(0, 127, 160);
      for (let e = 0; e < ye.length; ++e) I.fill(0, ye[e][0], ye[e][1] + 1);
      for (let e = 0; e < ge.length; ++e) I.fill(2, ge[e][0], ge[e][1] + 1);
    }
  }
  wcwidth(e) {
    return e < 32 ? 0 : e < 127 ? 1 : e < 65536 ? I[e] : je(e, lt) ? 0 : je(e, xt) ? 2 : 1;
  }
  charProperties(e, t) {
    let n = this.wcwidth(e), o2 = n === 0 && t !== 0;
    if (o2) {
      let d = w.extractWidth(t);
      d === 0 ? o2 = false : d > n && (n = d);
    }
    return w.createPropertyValue(0, n, o2);
  }
};
var Ke = class {
  activate(e) {
    e.unicode.register(new se());
  }
  dispose() {
  }
};

// node_modules/@xterm/addon-webgl/lib/addon-webgl.mjs
var Lr = Object.defineProperty;
var wr = Object.getOwnPropertyDescriptor;
var Yi = (i8, e, t, n) => {
  for (var s15 = n > 1 ? void 0 : n ? wr(e, t) : e, o2 = i8.length - 1, r5; o2 >= 0; o2--) (r5 = i8[o2]) && (s15 = (n ? r5(e, t, s15) : r5(s15)) || s15);
  return n && s15 && Lr(e, t, s15), s15;
};
var Qi = (i8, e) => (t, n) => e(t, n, i8);
var pi = class {
  constructor() {
    this.listeners = [], this.unexpectedErrorHandler = function(e) {
      setTimeout(() => {
        throw e.stack ? bt.isErrorNoTelemetry(e) ? new bt(e.message + `

` + e.stack) : new Error(e.message + `

` + e.stack) : e;
      }, 0);
    };
  }
  addListener(e) {
    return this.listeners.push(e), () => {
      this._removeListener(e);
    };
  }
  emit(e) {
    this.listeners.forEach((t) => {
      t(e);
    });
  }
  _removeListener(e) {
    this.listeners.splice(this.listeners.indexOf(e), 1);
  }
  setUnexpectedErrorHandler(e) {
    this.unexpectedErrorHandler = e;
  }
  getUnexpectedErrorHandler() {
    return this.unexpectedErrorHandler;
  }
  onUnexpectedError(e) {
    this.unexpectedErrorHandler(e), this.emit(e);
  }
  onUnexpectedExternalError(e) {
    this.unexpectedErrorHandler(e);
  }
};
var Rr = new pi();
function Pe2(i8) {
  Dr(i8) || Rr.onUnexpectedError(i8);
}
var fi = "Canceled";
function Dr(i8) {
  return i8 instanceof Ye2 ? true : i8 instanceof Error && i8.name === fi && i8.message === fi;
}
var Ye2 = class extends Error {
  constructor() {
    super(fi), this.name = this.message;
  }
};
var bt = class i extends Error {
  constructor(e) {
    super(e), this.name = "CodeExpectedError";
  }
  static fromError(e) {
    if (e instanceof i) return e;
    let t = new i();
    return t.message = e.message, t.stack = e.stack, t;
  }
  static isErrorNoTelemetry(e) {
    return e.name === "CodeExpectedError";
  }
};
function Mr(i8, e, t = 0, n = i8.length) {
  let s15 = t, o2 = n;
  for (; s15 < o2; ) {
    let r5 = Math.floor((s15 + o2) / 2);
    e(i8[r5]) ? s15 = r5 + 1 : o2 = r5;
  }
  return s15 - 1;
}
var vt = class vt2 {
  constructor(e) {
    this._array = e;
    this._findLastMonotonousLastIdx = 0;
  }
  findLastMonotonous(e) {
    if (vt2.assertInvariants) {
      if (this._prevFindLastPredicate) {
        for (let n of this._array) if (this._prevFindLastPredicate(n) && !e(n)) throw new Error("MonotonousArray: current predicate must be weaker than (or equal to) the previous predicate.");
      }
      this._prevFindLastPredicate = e;
    }
    let t = Mr(this._array, e, this._findLastMonotonousLastIdx);
    return this._findLastMonotonousLastIdx = t + 1, t === -1 ? void 0 : this._array[t];
  }
};
vt.assertInvariants = false;
var en;
((a) => {
  function i8(l) {
    return l < 0;
  }
  a.isLessThan = i8;
  function e(l) {
    return l <= 0;
  }
  a.isLessThanOrEqual = e;
  function t(l) {
    return l > 0;
  }
  a.isGreaterThan = t;
  function n(l) {
    return l === 0;
  }
  a.isNeitherLessOrGreaterThan = n, a.greaterThan = 1, a.lessThan = -1, a.neitherLessOrGreaterThan = 0;
})(en ||= {});
function tn(i8, e) {
  return (t, n) => e(i8(t), i8(n));
}
var nn = (i8, e) => i8 - e;
var Be2 = class Be3 {
  constructor(e) {
    this.iterate = e;
  }
  forEach(e) {
    this.iterate((t) => (e(t), true));
  }
  toArray() {
    let e = [];
    return this.iterate((t) => (e.push(t), true)), e;
  }
  filter(e) {
    return new Be3((t) => this.iterate((n) => e(n) ? t(n) : true));
  }
  map(e) {
    return new Be3((t) => this.iterate((n) => t(e(n))));
  }
  some(e) {
    let t = false;
    return this.iterate((n) => (t = e(n), !t)), t;
  }
  findFirst(e) {
    let t;
    return this.iterate((n) => e(n) ? (t = n, false) : true), t;
  }
  findLast(e) {
    let t;
    return this.iterate((n) => (e(n) && (t = n), true)), t;
  }
  findLastMaxBy(e) {
    let t, n = true;
    return this.iterate((s15) => ((n || en.isGreaterThan(e(s15, t))) && (n = false, t = s15), true)), t;
  }
};
Be2.empty = new Be2((e) => {
});
function an(i8, e) {
  let t = /* @__PURE__ */ Object.create(null);
  for (let n of i8) {
    let s15 = e(n), o2 = t[s15];
    o2 || (o2 = t[s15] = []), o2.push(n);
  }
  return t;
}
var sn;
var on;
var rn = class {
  constructor(e, t) {
    this.toKey = t;
    this._map = /* @__PURE__ */ new Map();
    this[sn] = "SetWithKey";
    for (let n of e) this.add(n);
  }
  get size() {
    return this._map.size;
  }
  add(e) {
    let t = this.toKey(e);
    return this._map.set(t, e), this;
  }
  delete(e) {
    return this._map.delete(this.toKey(e));
  }
  has(e) {
    return this._map.has(this.toKey(e));
  }
  *entries() {
    for (let e of this._map.values()) yield [e, e];
  }
  keys() {
    return this.values();
  }
  *values() {
    for (let e of this._map.values()) yield e;
  }
  clear() {
    this._map.clear();
  }
  forEach(e, t) {
    this._map.forEach((n) => e.call(t, n, n, this));
  }
  [(on = Symbol.iterator, sn = Symbol.toStringTag, on)]() {
    return this.values();
  }
};
var Tt = class {
  constructor() {
    this.map = /* @__PURE__ */ new Map();
  }
  add(e, t) {
    let n = this.map.get(e);
    n || (n = /* @__PURE__ */ new Set(), this.map.set(e, n)), n.add(t);
  }
  delete(e, t) {
    let n = this.map.get(e);
    n && (n.delete(t), n.size === 0 && this.map.delete(e));
  }
  forEach(e, t) {
    let n = this.map.get(e);
    n && n.forEach(t);
  }
  get(e) {
    let t = this.map.get(e);
    return t || /* @__PURE__ */ new Set();
  }
};
function mi(i8, e) {
  let t = this, n = false, s15;
  return function() {
    if (n) return s15;
    if (n = true, e) try {
      s15 = i8.apply(t, arguments);
    } finally {
      e();
    }
    else s15 = i8.apply(t, arguments);
    return s15;
  };
}
var _i;
((W) => {
  function i8(E) {
    return E && typeof E == "object" && typeof E[Symbol.iterator] == "function";
  }
  W.is = i8;
  let e = Object.freeze([]);
  function t() {
    return e;
  }
  W.empty = t;
  function* n(E) {
    yield E;
  }
  W.single = n;
  function s15(E) {
    return i8(E) ? E : n(E);
  }
  W.wrap = s15;
  function o2(E) {
    return E || e;
  }
  W.from = o2;
  function* r5(E) {
    for (let y = E.length - 1; y >= 0; y--) yield E[y];
  }
  W.reverse = r5;
  function a(E) {
    return !E || E[Symbol.iterator]().next().done === true;
  }
  W.isEmpty = a;
  function l(E) {
    return E[Symbol.iterator]().next().value;
  }
  W.first = l;
  function u(E, y) {
    let w2 = 0;
    for (let G3 of E) if (y(G3, w2++)) return true;
    return false;
  }
  W.some = u;
  function c(E, y) {
    for (let w2 of E) if (y(w2)) return w2;
  }
  W.find = c;
  function* d(E, y) {
    for (let w2 of E) y(w2) && (yield w2);
  }
  W.filter = d;
  function* h2(E, y) {
    let w2 = 0;
    for (let G3 of E) yield y(G3, w2++);
  }
  W.map = h2;
  function* f(E, y) {
    let w2 = 0;
    for (let G3 of E) yield* y(G3, w2++);
  }
  W.flatMap = f;
  function* I2(...E) {
    for (let y of E) yield* y;
  }
  W.concat = I2;
  function L2(E, y, w2) {
    let G3 = w2;
    for (let ue3 of E) G3 = y(G3, ue3);
    return G3;
  }
  W.reduce = L2;
  function* M4(E, y, w2 = E.length) {
    for (y < 0 && (y += E.length), w2 < 0 ? w2 += E.length : w2 > E.length && (w2 = E.length); y < w2; y++) yield E[y];
  }
  W.slice = M4;
  function q2(E, y = Number.POSITIVE_INFINITY) {
    let w2 = [];
    if (y === 0) return [w2, E];
    let G3 = E[Symbol.iterator]();
    for (let ue3 = 0; ue3 < y; ue3++) {
      let Se3 = G3.next();
      if (Se3.done) return [w2, W.empty()];
      w2.push(Se3.value);
    }
    return [w2, { [Symbol.iterator]() {
      return G3;
    } }];
  }
  W.consume = q2;
  async function S2(E) {
    let y = [];
    for await (let w2 of E) y.push(w2);
    return Promise.resolve(y);
  }
  W.asyncToArray = S2;
})(_i ||= {});
var Ar = false;
var Ne = null;
var gt = class gt2 {
  constructor() {
    this.livingDisposables = /* @__PURE__ */ new Map();
  }
  getDisposableData(e) {
    let t = this.livingDisposables.get(e);
    return t || (t = { parent: null, source: null, isSingleton: false, value: e, idx: gt2.idx++ }, this.livingDisposables.set(e, t)), t;
  }
  trackDisposable(e) {
    let t = this.getDisposableData(e);
    t.source || (t.source = new Error().stack);
  }
  setParent(e, t) {
    let n = this.getDisposableData(e);
    n.parent = t;
  }
  markAsDisposed(e) {
    this.livingDisposables.delete(e);
  }
  markAsSingleton(e) {
    this.getDisposableData(e).isSingleton = true;
  }
  getRootParent(e, t) {
    let n = t.get(e);
    if (n) return n;
    let s15 = e.parent ? this.getRootParent(this.getDisposableData(e.parent), t) : e;
    return t.set(e, s15), s15;
  }
  getTrackedDisposables() {
    let e = /* @__PURE__ */ new Map();
    return [...this.livingDisposables.entries()].filter(([, n]) => n.source !== null && !this.getRootParent(n, e).isSingleton).flatMap(([n]) => n);
  }
  computeLeakingDisposables(e = 10, t) {
    let n;
    if (t) n = t;
    else {
      let l = /* @__PURE__ */ new Map(), u = [...this.livingDisposables.values()].filter((d) => d.source !== null && !this.getRootParent(d, l).isSingleton);
      if (u.length === 0) return;
      let c = new Set(u.map((d) => d.value));
      if (n = u.filter((d) => !(d.parent && c.has(d.parent))), n.length === 0) throw new Error("There are cyclic diposable chains!");
    }
    if (!n) return;
    function s15(l) {
      function u(d, h2) {
        for (; d.length > 0 && h2.some((f) => typeof f == "string" ? f === d[0] : d[0].match(f)); ) d.shift();
      }
      let c = l.source.split(`
`).map((d) => d.trim().replace("at ", "")).filter((d) => d !== "");
      return u(c, ["Error", /^trackDisposable \(.*\)$/, /^DisposableTracker.trackDisposable \(.*\)$/]), c.reverse();
    }
    let o2 = new Tt();
    for (let l of n) {
      let u = s15(l);
      for (let c = 0; c <= u.length; c++) o2.add(u.slice(0, c).join(`
`), l);
    }
    n.sort(tn((l) => l.idx, nn));
    let r5 = "", a = 0;
    for (let l of n.slice(0, e)) {
      a++;
      let u = s15(l), c = [];
      for (let d = 0; d < u.length; d++) {
        let h2 = u[d];
        h2 = `(shared with ${o2.get(u.slice(0, d + 1).join(`
`)).size}/${n.length} leaks) at ${h2}`;
        let I2 = o2.get(u.slice(0, d).join(`
`)), L2 = an([...I2].map((M4) => s15(M4)[d]), (M4) => M4);
        delete L2[u[d]];
        for (let [M4, q2] of Object.entries(L2)) c.unshift(`    - stacktraces of ${q2.length} other leaks continue with ${M4}`);
        c.unshift(h2);
      }
      r5 += `


==================== Leaking disposable ${a}/${n.length}: ${l.value.constructor.name} ====================
${c.join(`
`)}
============================================================

`;
    }
    return n.length > e && (r5 += `


... and ${n.length - e} more leaking disposables

`), { leaks: n, details: r5 };
  }
};
gt.idx = 0;
function Sr(i8) {
  Ne = i8;
}
if (Ar) {
  let i8 = "__is_disposable_tracked__";
  Sr(new class {
    trackDisposable(e) {
      let t = new Error("Potentially leaked disposable").stack;
      setTimeout(() => {
        e[i8] || console.log(t);
      }, 3e3);
    }
    setParent(e, t) {
      if (e && e !== B.None) try {
        e[i8] = true;
      } catch {
      }
    }
    markAsDisposed(e) {
      if (e && e !== B.None) try {
        e[i8] = true;
      } catch {
      }
    }
    markAsSingleton(e) {
    }
  }());
}
function Et(i8) {
  return Ne?.trackDisposable(i8), i8;
}
function yt(i8) {
  Ne?.markAsDisposed(i8);
}
function Qe(i8, e) {
  Ne?.setParent(i8, e);
}
function Or(i8, e) {
  if (Ne) for (let t of i8) Ne.setParent(t, e);
}
function un(i8) {
  if (_i.is(i8)) {
    let e = [];
    for (let t of i8) if (t) try {
      t.dispose();
    } catch (n) {
      e.push(n);
    }
    if (e.length === 1) throw e[0];
    if (e.length > 1) throw new AggregateError(e, "Encountered errors while disposing of store");
    return Array.isArray(i8) ? [] : i8;
  } else if (i8) return i8.dispose(), i8;
}
function It(...i8) {
  let e = O2(() => un(i8));
  return Or(i8, e), e;
}
function O2(i8) {
  let e = Et({ dispose: mi(() => {
    yt(e), i8();
  }) });
  return e;
}
var xt2 = class xt3 {
  constructor() {
    this._toDispose = /* @__PURE__ */ new Set();
    this._isDisposed = false;
    Et(this);
  }
  dispose() {
    this._isDisposed || (yt(this), this._isDisposed = true, this.clear());
  }
  get isDisposed() {
    return this._isDisposed;
  }
  clear() {
    if (this._toDispose.size !== 0) try {
      un(this._toDispose);
    } finally {
      this._toDispose.clear();
    }
  }
  add(e) {
    if (!e) return e;
    if (e === this) throw new Error("Cannot register a disposable on itself!");
    return Qe(e, this), this._isDisposed ? xt3.DISABLE_DISPOSED_WARNING || console.warn(new Error("Trying to add a disposable to a DisposableStore that has already been disposed of. The added object will be leaked!").stack) : this._toDispose.add(e), e;
  }
  delete(e) {
    if (e) {
      if (e === this) throw new Error("Cannot dispose a disposable on itself!");
      this._toDispose.delete(e), e.dispose();
    }
  }
  deleteAndLeak(e) {
    e && this._toDispose.has(e) && (this._toDispose.delete(e), Qe(e, null));
  }
};
xt2.DISABLE_DISPOSED_WARNING = false;
var fe2 = xt2;
var B = class {
  constructor() {
    this._store = new fe2();
    Et(this), Qe(this._store, this);
  }
  dispose() {
    yt(this), this._store.dispose();
  }
  _register(e) {
    if (e === this) throw new Error("Cannot register a disposable on itself!");
    return this._store.add(e);
  }
};
B.None = Object.freeze({ dispose() {
} });
var be2 = class {
  constructor() {
    this._isDisposed = false;
    Et(this);
  }
  get value() {
    return this._isDisposed ? void 0 : this._value;
  }
  set value(e) {
    this._isDisposed || e === this._value || (this._value?.dispose(), e && Qe(e, this), this._value = e);
  }
  clear() {
    this.value = void 0;
  }
  dispose() {
    this._isDisposed = true, yt(this), this._value?.dispose(), this._value = void 0;
  }
  clearAndLeak() {
    let e = this._value;
    return this._value = void 0, e && Qe(e, null), e;
  }
};
var Lt = typeof process < "u" && "title" in process;
var Ze2 = Lt ? "node" : navigator.userAgent;
var bi = Lt ? "node" : navigator.platform;
var cn = Ze2.includes("Firefox");
var dn = Ze2.includes("Edge");
var vi = /^((?!chrome|android).)*safari/i.test(Ze2);
function hn() {
  if (!vi) return 0;
  let i8 = Ze2.match(/Version\/(\d+)/);
  return i8 === null || i8.length < 2 ? 0 : parseInt(i8[1]);
}
var oo = ["Macintosh", "MacIntel", "MacPPC", "Mac68K"].includes(bi);
var ao = ["Windows", "Win16", "Win32", "WinCE"].includes(bi);
var lo = bi.indexOf("Linux") >= 0;
var uo = /\bCrOS\b/.test(Ze2);
var pn = "";
var K2 = 0;
var V = 0;
var C2 = 0;
var U2 = 0;
var Z2 = { css: "#00000000", rgba: 0 };
var X3;
((n) => {
  function i8(s15, o2, r5, a) {
    return a !== void 0 ? `#${Oe2(s15)}${Oe2(o2)}${Oe2(r5)}${Oe2(a)}` : `#${Oe2(s15)}${Oe2(o2)}${Oe2(r5)}`;
  }
  n.toCss = i8;
  function e(s15, o2, r5, a = 255) {
    return (s15 << 24 | o2 << 16 | r5 << 8 | a) >>> 0;
  }
  n.toRgba = e;
  function t(s15, o2, r5, a) {
    return { css: n.toCss(s15, o2, r5, a), rgba: n.toRgba(s15, o2, r5, a) };
  }
  n.toColor = t;
})(X3 ||= {});
var Ue;
((a) => {
  function i8(l, u) {
    if (U2 = (u.rgba & 255) / 255, U2 === 1) return { css: u.css, rgba: u.rgba };
    let c = u.rgba >> 24 & 255, d = u.rgba >> 16 & 255, h2 = u.rgba >> 8 & 255, f = l.rgba >> 24 & 255, I2 = l.rgba >> 16 & 255, L2 = l.rgba >> 8 & 255;
    K2 = f + Math.round((c - f) * U2), V = I2 + Math.round((d - I2) * U2), C2 = L2 + Math.round((h2 - L2) * U2);
    let M4 = X3.toCss(K2, V, C2), q2 = X3.toRgba(K2, V, C2);
    return { css: M4, rgba: q2 };
  }
  a.blend = i8;
  function e(l) {
    return (l.rgba & 255) === 255;
  }
  a.isOpaque = e;
  function t(l, u, c) {
    let d = Te2.ensureContrastRatio(l.rgba, u.rgba, c);
    if (d) return X3.toColor(d >> 24 & 255, d >> 16 & 255, d >> 8 & 255);
  }
  a.ensureContrastRatio = t;
  function n(l) {
    let u = (l.rgba | 255) >>> 0;
    return [K2, V, C2] = Te2.toChannels(u), { css: X3.toCss(K2, V, C2), rgba: u };
  }
  a.opaque = n;
  function s15(l, u) {
    return U2 = Math.round(u * 255), [K2, V, C2] = Te2.toChannels(l.rgba), { css: X3.toCss(K2, V, C2, U2), rgba: X3.toRgba(K2, V, C2, U2) };
  }
  a.opacity = s15;
  function o2(l, u) {
    return U2 = l.rgba & 255, s15(l, U2 * u / 255);
  }
  a.multiplyOpacity = o2;
  function r5(l) {
    return [l.rgba >> 24 & 255, l.rgba >> 16 & 255, l.rgba >> 8 & 255];
  }
  a.toColorRGB = r5;
})(Ue ||= {});
var Fr;
((n) => {
  let i8, e;
  try {
    let s15 = document.createElement("canvas");
    s15.width = 1, s15.height = 1;
    let o2 = s15.getContext("2d", { willReadFrequently: true });
    o2 && (i8 = o2, i8.globalCompositeOperation = "copy", e = i8.createLinearGradient(0, 0, 1, 1));
  } catch {
  }
  function t(s15) {
    if (s15.match(/#[\da-f]{3,8}/i)) switch (s15.length) {
      case 4:
        return K2 = parseInt(s15.slice(1, 2).repeat(2), 16), V = parseInt(s15.slice(2, 3).repeat(2), 16), C2 = parseInt(s15.slice(3, 4).repeat(2), 16), X3.toColor(K2, V, C2);
      case 5:
        return K2 = parseInt(s15.slice(1, 2).repeat(2), 16), V = parseInt(s15.slice(2, 3).repeat(2), 16), C2 = parseInt(s15.slice(3, 4).repeat(2), 16), U2 = parseInt(s15.slice(4, 5).repeat(2), 16), X3.toColor(K2, V, C2, U2);
      case 7:
        return { css: s15, rgba: (parseInt(s15.slice(1), 16) << 8 | 255) >>> 0 };
      case 9:
        return { css: s15, rgba: parseInt(s15.slice(1), 16) >>> 0 };
    }
    let o2 = s15.match(/rgba?\(\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*(\d{1,3})\s*(,\s*(0|1|\d?\.(\d+))\s*)?\)/);
    if (o2) return K2 = parseInt(o2[1]), V = parseInt(o2[2]), C2 = parseInt(o2[3]), U2 = Math.round((o2[5] === void 0 ? 1 : parseFloat(o2[5])) * 255), X3.toColor(K2, V, C2, U2);
    if (!i8 || !e) throw new Error("css.toColor: Unsupported css format");
    if (i8.fillStyle = e, i8.fillStyle = s15, typeof i8.fillStyle != "string") throw new Error("css.toColor: Unsupported css format");
    if (i8.fillRect(0, 0, 1, 1), [K2, V, C2, U2] = i8.getImageData(0, 0, 1, 1).data, U2 !== 255) throw new Error("css.toColor: Unsupported css format");
    return { rgba: X3.toRgba(K2, V, C2, U2), css: s15 };
  }
  n.toColor = t;
})(Fr ||= {});
var Y2;
((t) => {
  function i8(n) {
    return e(n >> 16 & 255, n >> 8 & 255, n & 255);
  }
  t.relativeLuminance = i8;
  function e(n, s15, o2) {
    let r5 = n / 255, a = s15 / 255, l = o2 / 255, u = r5 <= 0.03928 ? r5 / 12.92 : Math.pow((r5 + 0.055) / 1.055, 2.4), c = a <= 0.03928 ? a / 12.92 : Math.pow((a + 0.055) / 1.055, 2.4), d = l <= 0.03928 ? l / 12.92 : Math.pow((l + 0.055) / 1.055, 2.4);
    return u * 0.2126 + c * 0.7152 + d * 0.0722;
  }
  t.relativeLuminance2 = e;
})(Y2 ||= {});
var Te2;
((o2) => {
  function i8(r5, a) {
    if (U2 = (a & 255) / 255, U2 === 1) return a;
    let l = a >> 24 & 255, u = a >> 16 & 255, c = a >> 8 & 255, d = r5 >> 24 & 255, h2 = r5 >> 16 & 255, f = r5 >> 8 & 255;
    return K2 = d + Math.round((l - d) * U2), V = h2 + Math.round((u - h2) * U2), C2 = f + Math.round((c - f) * U2), X3.toRgba(K2, V, C2);
  }
  o2.blend = i8;
  function e(r5, a, l) {
    let u = Y2.relativeLuminance(r5 >> 8), c = Y2.relativeLuminance(a >> 8);
    if (ve2(u, c) < l) {
      if (c < u) {
        let I2 = t(r5, a, l), L2 = ve2(u, Y2.relativeLuminance(I2 >> 8));
        if (L2 < l) {
          let M4 = n(r5, a, l), q2 = ve2(u, Y2.relativeLuminance(M4 >> 8));
          return L2 > q2 ? I2 : M4;
        }
        return I2;
      }
      let h2 = n(r5, a, l), f = ve2(u, Y2.relativeLuminance(h2 >> 8));
      if (f < l) {
        let I2 = t(r5, a, l), L2 = ve2(u, Y2.relativeLuminance(I2 >> 8));
        return f > L2 ? h2 : I2;
      }
      return h2;
    }
  }
  o2.ensureContrastRatio = e;
  function t(r5, a, l) {
    let u = r5 >> 24 & 255, c = r5 >> 16 & 255, d = r5 >> 8 & 255, h2 = a >> 24 & 255, f = a >> 16 & 255, I2 = a >> 8 & 255, L2 = ve2(Y2.relativeLuminance2(h2, f, I2), Y2.relativeLuminance2(u, c, d));
    for (; L2 < l && (h2 > 0 || f > 0 || I2 > 0); ) h2 -= Math.max(0, Math.ceil(h2 * 0.1)), f -= Math.max(0, Math.ceil(f * 0.1)), I2 -= Math.max(0, Math.ceil(I2 * 0.1)), L2 = ve2(Y2.relativeLuminance2(h2, f, I2), Y2.relativeLuminance2(u, c, d));
    return (h2 << 24 | f << 16 | I2 << 8 | 255) >>> 0;
  }
  o2.reduceLuminance = t;
  function n(r5, a, l) {
    let u = r5 >> 24 & 255, c = r5 >> 16 & 255, d = r5 >> 8 & 255, h2 = a >> 24 & 255, f = a >> 16 & 255, I2 = a >> 8 & 255, L2 = ve2(Y2.relativeLuminance2(h2, f, I2), Y2.relativeLuminance2(u, c, d));
    for (; L2 < l && (h2 < 255 || f < 255 || I2 < 255); ) h2 = Math.min(255, h2 + Math.ceil((255 - h2) * 0.1)), f = Math.min(255, f + Math.ceil((255 - f) * 0.1)), I2 = Math.min(255, I2 + Math.ceil((255 - I2) * 0.1)), L2 = ve2(Y2.relativeLuminance2(h2, f, I2), Y2.relativeLuminance2(u, c, d));
    return (h2 << 24 | f << 16 | I2 << 8 | 255) >>> 0;
  }
  o2.increaseLuminance = n;
  function s15(r5) {
    return [r5 >> 24 & 255, r5 >> 16 & 255, r5 >> 8 & 255, r5 & 255];
  }
  o2.toChannels = s15;
})(Te2 ||= {});
function Oe2(i8) {
  let e = i8.toString(16);
  return e.length < 2 ? "0" + e : e;
}
function ve2(i8, e) {
  return i8 < e ? (e + 0.05) / (i8 + 0.05) : (i8 + 0.05) / (e + 0.05);
}
function F(i8) {
  if (!i8) throw new Error("value must not be falsy");
  return i8;
}
function Rt(i8) {
  return 57508 <= i8 && i8 <= 57558;
}
function fn(i8) {
  return 57520 <= i8 && i8 <= 57527;
}
function kr(i8) {
  return 57344 <= i8 && i8 <= 63743;
}
function Pr(i8) {
  return 9472 <= i8 && i8 <= 9631;
}
function Br(i8) {
  return i8 >= 128512 && i8 <= 128591 || i8 >= 127744 && i8 <= 128511 || i8 >= 128640 && i8 <= 128767 || i8 >= 9728 && i8 <= 9983 || i8 >= 9984 && i8 <= 10175 || i8 >= 65024 && i8 <= 65039 || i8 >= 129280 && i8 <= 129535 || i8 >= 127462 && i8 <= 127487;
}
function mn(i8, e, t, n) {
  return e === 1 && t > Math.ceil(n * 1.5) && i8 !== void 0 && i8 > 255 && !Br(i8) && !Rt(i8) && !kr(i8);
}
function Dt(i8) {
  return Rt(i8) || Pr(i8);
}
function _n() {
  return { css: { canvas: wt(), cell: wt() }, device: { canvas: wt(), cell: wt(), char: { width: 0, height: 0, left: 0, top: 0 } } };
}
function wt() {
  return { width: 0, height: 0 };
}
function bn(i8, e, t = 0) {
  return (i8 - (Math.round(e) * 2 - t)) % (Math.round(e) * 2);
}
var j = 0;
var z = 0;
var me2 = false;
var ge2 = false;
var Mt = false;
var J2;
var Ti = 0;
var At = class {
  constructor(e, t, n, s15, o2, r5) {
    this._terminal = e;
    this._optionService = t;
    this._selectionRenderModel = n;
    this._decorationService = s15;
    this._coreBrowserService = o2;
    this._themeService = r5;
    this.result = { fg: 0, bg: 0, ext: 0 };
  }
  resolve(e, t, n, s15) {
    if (this.result.bg = e.bg, this.result.fg = e.fg, this.result.ext = e.bg & 268435456 ? e.extended.ext : 0, z = 0, j = 0, ge2 = false, me2 = false, Mt = false, J2 = this._themeService.colors, Ti = 0, e.getCode() !== 0 && e.extended.underlineStyle === 4) {
      let r5 = Math.max(1, Math.floor(this._optionService.rawOptions.fontSize * this._coreBrowserService.dpr / 15));
      Ti = t * s15 % (Math.round(r5) * 2);
    }
    if (this._decorationService.forEachDecorationAtCell(t, n, "bottom", (r5) => {
      r5.backgroundColorRGB && (z = r5.backgroundColorRGB.rgba >> 8 & 16777215, ge2 = true), r5.foregroundColorRGB && (j = r5.foregroundColorRGB.rgba >> 8 & 16777215, me2 = true);
    }), Mt = this._selectionRenderModel.isCellSelected(this._terminal, t, n), Mt) {
      if (this.result.fg & 67108864 || (this.result.bg & 50331648) !== 0) {
        if (this.result.fg & 67108864) switch (this.result.fg & 50331648) {
          case 16777216:
          case 33554432:
            z = this._themeService.colors.ansi[this.result.fg & 255].rgba;
            break;
          case 50331648:
            z = (this.result.fg & 16777215) << 8 | 255;
            break;
          case 0:
          default:
            z = this._themeService.colors.foreground.rgba;
        }
        else switch (this.result.bg & 50331648) {
          case 16777216:
          case 33554432:
            z = this._themeService.colors.ansi[this.result.bg & 255].rgba;
            break;
          case 50331648:
            z = (this.result.bg & 16777215) << 8 | 255;
            break;
        }
        z = Te2.blend(z, (this._coreBrowserService.isFocused ? J2.selectionBackgroundOpaque : J2.selectionInactiveBackgroundOpaque).rgba & 4294967040 | 128) >> 8 & 16777215;
      } else z = (this._coreBrowserService.isFocused ? J2.selectionBackgroundOpaque : J2.selectionInactiveBackgroundOpaque).rgba >> 8 & 16777215;
      if (ge2 = true, J2.selectionForeground && (j = J2.selectionForeground.rgba >> 8 & 16777215, me2 = true), Dt(e.getCode())) {
        if (this.result.fg & 67108864 && (this.result.bg & 50331648) === 0) j = (this._coreBrowserService.isFocused ? J2.selectionBackgroundOpaque : J2.selectionInactiveBackgroundOpaque).rgba >> 8 & 16777215;
        else {
          if (this.result.fg & 67108864) switch (this.result.bg & 50331648) {
            case 16777216:
            case 33554432:
              j = this._themeService.colors.ansi[this.result.bg & 255].rgba;
              break;
            case 50331648:
              j = (this.result.bg & 16777215) << 8 | 255;
              break;
          }
          else switch (this.result.fg & 50331648) {
            case 16777216:
            case 33554432:
              j = this._themeService.colors.ansi[this.result.fg & 255].rgba;
              break;
            case 50331648:
              j = (this.result.fg & 16777215) << 8 | 255;
              break;
            case 0:
            default:
              j = this._themeService.colors.foreground.rgba;
          }
          j = Te2.blend(j, (this._coreBrowserService.isFocused ? J2.selectionBackgroundOpaque : J2.selectionInactiveBackgroundOpaque).rgba & 4294967040 | 128) >> 8 & 16777215;
        }
        me2 = true;
      }
    }
    this._decorationService.forEachDecorationAtCell(t, n, "top", (r5) => {
      r5.backgroundColorRGB && (z = r5.backgroundColorRGB.rgba >> 8 & 16777215, ge2 = true), r5.foregroundColorRGB && (j = r5.foregroundColorRGB.rgba >> 8 & 16777215, me2 = true);
    }), ge2 && (Mt ? z = e.bg & -16777216 & -134217729 | z | 50331648 : z = e.bg & -16777216 | z | 50331648), me2 && (j = e.fg & -16777216 & -67108865 | j | 50331648), this.result.fg & 67108864 && (ge2 && !me2 && ((this.result.bg & 50331648) === 0 ? j = this.result.fg & -134217728 | J2.background.rgba >> 8 & 16777215 & 16777215 | 50331648 : j = this.result.fg & -134217728 | this.result.bg & 67108863, me2 = true), !ge2 && me2 && ((this.result.fg & 50331648) === 0 ? z = this.result.bg & -67108864 | J2.foreground.rgba >> 8 & 16777215 & 16777215 | 50331648 : z = this.result.bg & -67108864 | this.result.fg & 67108863, ge2 = true)), J2 = void 0, this.result.bg = ge2 ? z : this.result.bg, this.result.fg = me2 ? j : this.result.fg, this.result.ext &= 536870911, this.result.ext |= Ti << 29 & 3758096384;
  }
};
var gn = 0.5;
var St = cn || dn ? "bottom" : "ideographic";
var Hr = { "▀": [{ x: 0, y: 0, w: 8, h: 4 }], "▁": [{ x: 0, y: 7, w: 8, h: 1 }], "▂": [{ x: 0, y: 6, w: 8, h: 2 }], "▃": [{ x: 0, y: 5, w: 8, h: 3 }], "▄": [{ x: 0, y: 4, w: 8, h: 4 }], "▅": [{ x: 0, y: 3, w: 8, h: 5 }], "▆": [{ x: 0, y: 2, w: 8, h: 6 }], "▇": [{ x: 0, y: 1, w: 8, h: 7 }], "█": [{ x: 0, y: 0, w: 8, h: 8 }], "▉": [{ x: 0, y: 0, w: 7, h: 8 }], "▊": [{ x: 0, y: 0, w: 6, h: 8 }], "▋": [{ x: 0, y: 0, w: 5, h: 8 }], "▌": [{ x: 0, y: 0, w: 4, h: 8 }], "▍": [{ x: 0, y: 0, w: 3, h: 8 }], "▎": [{ x: 0, y: 0, w: 2, h: 8 }], "▏": [{ x: 0, y: 0, w: 1, h: 8 }], "▐": [{ x: 4, y: 0, w: 4, h: 8 }], "▔": [{ x: 0, y: 0, w: 8, h: 1 }], "▕": [{ x: 7, y: 0, w: 1, h: 8 }], "▖": [{ x: 0, y: 4, w: 4, h: 4 }], "▗": [{ x: 4, y: 4, w: 4, h: 4 }], "▘": [{ x: 0, y: 0, w: 4, h: 4 }], "▙": [{ x: 0, y: 0, w: 4, h: 8 }, { x: 0, y: 4, w: 8, h: 4 }], "▚": [{ x: 0, y: 0, w: 4, h: 4 }, { x: 4, y: 4, w: 4, h: 4 }], "▛": [{ x: 0, y: 0, w: 4, h: 8 }, { x: 4, y: 0, w: 4, h: 4 }], "▜": [{ x: 0, y: 0, w: 8, h: 4 }, { x: 4, y: 0, w: 4, h: 8 }], "▝": [{ x: 4, y: 0, w: 4, h: 4 }], "▞": [{ x: 4, y: 0, w: 4, h: 4 }, { x: 0, y: 4, w: 4, h: 4 }], "▟": [{ x: 4, y: 0, w: 4, h: 8 }, { x: 0, y: 4, w: 8, h: 4 }], "🭰": [{ x: 1, y: 0, w: 1, h: 8 }], "🭱": [{ x: 2, y: 0, w: 1, h: 8 }], "🭲": [{ x: 3, y: 0, w: 1, h: 8 }], "🭳": [{ x: 4, y: 0, w: 1, h: 8 }], "🭴": [{ x: 5, y: 0, w: 1, h: 8 }], "🭵": [{ x: 6, y: 0, w: 1, h: 8 }], "🭶": [{ x: 0, y: 1, w: 8, h: 1 }], "🭷": [{ x: 0, y: 2, w: 8, h: 1 }], "🭸": [{ x: 0, y: 3, w: 8, h: 1 }], "🭹": [{ x: 0, y: 4, w: 8, h: 1 }], "🭺": [{ x: 0, y: 5, w: 8, h: 1 }], "🭻": [{ x: 0, y: 6, w: 8, h: 1 }], "🭼": [{ x: 0, y: 0, w: 1, h: 8 }, { x: 0, y: 7, w: 8, h: 1 }], "🭽": [{ x: 0, y: 0, w: 1, h: 8 }, { x: 0, y: 0, w: 8, h: 1 }], "🭾": [{ x: 7, y: 0, w: 1, h: 8 }, { x: 0, y: 0, w: 8, h: 1 }], "🭿": [{ x: 7, y: 0, w: 1, h: 8 }, { x: 0, y: 7, w: 8, h: 1 }], "🮀": [{ x: 0, y: 0, w: 8, h: 1 }, { x: 0, y: 7, w: 8, h: 1 }], "🮁": [{ x: 0, y: 0, w: 8, h: 1 }, { x: 0, y: 2, w: 8, h: 1 }, { x: 0, y: 4, w: 8, h: 1 }, { x: 0, y: 7, w: 8, h: 1 }], "🮂": [{ x: 0, y: 0, w: 8, h: 2 }], "🮃": [{ x: 0, y: 0, w: 8, h: 3 }], "🮄": [{ x: 0, y: 0, w: 8, h: 5 }], "🮅": [{ x: 0, y: 0, w: 8, h: 6 }], "🮆": [{ x: 0, y: 0, w: 8, h: 7 }], "🮇": [{ x: 6, y: 0, w: 2, h: 8 }], "🮈": [{ x: 5, y: 0, w: 3, h: 8 }], "🮉": [{ x: 3, y: 0, w: 5, h: 8 }], "🮊": [{ x: 2, y: 0, w: 6, h: 8 }], "🮋": [{ x: 1, y: 0, w: 7, h: 8 }], "🮕": [{ x: 0, y: 0, w: 2, h: 2 }, { x: 4, y: 0, w: 2, h: 2 }, { x: 2, y: 2, w: 2, h: 2 }, { x: 6, y: 2, w: 2, h: 2 }, { x: 0, y: 4, w: 2, h: 2 }, { x: 4, y: 4, w: 2, h: 2 }, { x: 2, y: 6, w: 2, h: 2 }, { x: 6, y: 6, w: 2, h: 2 }], "🮖": [{ x: 2, y: 0, w: 2, h: 2 }, { x: 6, y: 0, w: 2, h: 2 }, { x: 0, y: 2, w: 2, h: 2 }, { x: 4, y: 2, w: 2, h: 2 }, { x: 2, y: 4, w: 2, h: 2 }, { x: 6, y: 4, w: 2, h: 2 }, { x: 0, y: 6, w: 2, h: 2 }, { x: 4, y: 6, w: 2, h: 2 }], "🮗": [{ x: 0, y: 2, w: 8, h: 2 }, { x: 0, y: 6, w: 8, h: 2 }] };
var Wr = { "░": [[1, 0, 0, 0], [0, 0, 0, 0], [0, 0, 1, 0], [0, 0, 0, 0]], "▒": [[1, 0], [0, 0], [0, 1], [0, 0]], "▓": [[0, 1], [1, 1], [1, 0], [1, 1]] };
var Gr = { "─": { 1: "M0,.5 L1,.5" }, "━": { 3: "M0,.5 L1,.5" }, "│": { 1: "M.5,0 L.5,1" }, "┃": { 3: "M.5,0 L.5,1" }, "┌": { 1: "M0.5,1 L.5,.5 L1,.5" }, "┏": { 3: "M0.5,1 L.5,.5 L1,.5" }, "┐": { 1: "M0,.5 L.5,.5 L.5,1" }, "┓": { 3: "M0,.5 L.5,.5 L.5,1" }, "└": { 1: "M.5,0 L.5,.5 L1,.5" }, "┗": { 3: "M.5,0 L.5,.5 L1,.5" }, "┘": { 1: "M.5,0 L.5,.5 L0,.5" }, "┛": { 3: "M.5,0 L.5,.5 L0,.5" }, "├": { 1: "M.5,0 L.5,1 M.5,.5 L1,.5" }, "┣": { 3: "M.5,0 L.5,1 M.5,.5 L1,.5" }, "┤": { 1: "M.5,0 L.5,1 M.5,.5 L0,.5" }, "┫": { 3: "M.5,0 L.5,1 M.5,.5 L0,.5" }, "┬": { 1: "M0,.5 L1,.5 M.5,.5 L.5,1" }, "┳": { 3: "M0,.5 L1,.5 M.5,.5 L.5,1" }, "┴": { 1: "M0,.5 L1,.5 M.5,.5 L.5,0" }, "┻": { 3: "M0,.5 L1,.5 M.5,.5 L.5,0" }, "┼": { 1: "M0,.5 L1,.5 M.5,0 L.5,1" }, "╋": { 3: "M0,.5 L1,.5 M.5,0 L.5,1" }, "╴": { 1: "M.5,.5 L0,.5" }, "╸": { 3: "M.5,.5 L0,.5" }, "╵": { 1: "M.5,.5 L.5,0" }, "╹": { 3: "M.5,.5 L.5,0" }, "╶": { 1: "M.5,.5 L1,.5" }, "╺": { 3: "M.5,.5 L1,.5" }, "╷": { 1: "M.5,.5 L.5,1" }, "╻": { 3: "M.5,.5 L.5,1" }, "═": { 1: (i8, e) => `M0,${0.5 - e} L1,${0.5 - e} M0,${0.5 + e} L1,${0.5 + e}` }, "║": { 1: (i8, e) => `M${0.5 - i8},0 L${0.5 - i8},1 M${0.5 + i8},0 L${0.5 + i8},1` }, "╒": { 1: (i8, e) => `M.5,1 L.5,${0.5 - e} L1,${0.5 - e} M.5,${0.5 + e} L1,${0.5 + e}` }, "╓": { 1: (i8, e) => `M${0.5 - i8},1 L${0.5 - i8},.5 L1,.5 M${0.5 + i8},.5 L${0.5 + i8},1` }, "╔": { 1: (i8, e) => `M1,${0.5 - e} L${0.5 - i8},${0.5 - e} L${0.5 - i8},1 M1,${0.5 + e} L${0.5 + i8},${0.5 + e} L${0.5 + i8},1` }, "╕": { 1: (i8, e) => `M0,${0.5 - e} L.5,${0.5 - e} L.5,1 M0,${0.5 + e} L.5,${0.5 + e}` }, "╖": { 1: (i8, e) => `M${0.5 + i8},1 L${0.5 + i8},.5 L0,.5 M${0.5 - i8},.5 L${0.5 - i8},1` }, "╗": { 1: (i8, e) => `M0,${0.5 + e} L${0.5 - i8},${0.5 + e} L${0.5 - i8},1 M0,${0.5 - e} L${0.5 + i8},${0.5 - e} L${0.5 + i8},1` }, "╘": { 1: (i8, e) => `M.5,0 L.5,${0.5 + e} L1,${0.5 + e} M.5,${0.5 - e} L1,${0.5 - e}` }, "╙": { 1: (i8, e) => `M1,.5 L${0.5 - i8},.5 L${0.5 - i8},0 M${0.5 + i8},.5 L${0.5 + i8},0` }, "╚": { 1: (i8, e) => `M1,${0.5 - e} L${0.5 + i8},${0.5 - e} L${0.5 + i8},0 M1,${0.5 + e} L${0.5 - i8},${0.5 + e} L${0.5 - i8},0` }, "╛": { 1: (i8, e) => `M0,${0.5 + e} L.5,${0.5 + e} L.5,0 M0,${0.5 - e} L.5,${0.5 - e}` }, "╜": { 1: (i8, e) => `M0,.5 L${0.5 + i8},.5 L${0.5 + i8},0 M${0.5 - i8},.5 L${0.5 - i8},0` }, "╝": { 1: (i8, e) => `M0,${0.5 - e} L${0.5 - i8},${0.5 - e} L${0.5 - i8},0 M0,${0.5 + e} L${0.5 + i8},${0.5 + e} L${0.5 + i8},0` }, "╞": { 1: (i8, e) => `M.5,0 L.5,1 M.5,${0.5 - e} L1,${0.5 - e} M.5,${0.5 + e} L1,${0.5 + e}` }, "╟": { 1: (i8, e) => `M${0.5 - i8},0 L${0.5 - i8},1 M${0.5 + i8},0 L${0.5 + i8},1 M${0.5 + i8},.5 L1,.5` }, "╠": { 1: (i8, e) => `M${0.5 - i8},0 L${0.5 - i8},1 M1,${0.5 + e} L${0.5 + i8},${0.5 + e} L${0.5 + i8},1 M1,${0.5 - e} L${0.5 + i8},${0.5 - e} L${0.5 + i8},0` }, "╡": { 1: (i8, e) => `M.5,0 L.5,1 M0,${0.5 - e} L.5,${0.5 - e} M0,${0.5 + e} L.5,${0.5 + e}` }, "╢": { 1: (i8, e) => `M0,.5 L${0.5 - i8},.5 M${0.5 - i8},0 L${0.5 - i8},1 M${0.5 + i8},0 L${0.5 + i8},1` }, "╣": { 1: (i8, e) => `M${0.5 + i8},0 L${0.5 + i8},1 M0,${0.5 + e} L${0.5 - i8},${0.5 + e} L${0.5 - i8},1 M0,${0.5 - e} L${0.5 - i8},${0.5 - e} L${0.5 - i8},0` }, "╤": { 1: (i8, e) => `M0,${0.5 - e} L1,${0.5 - e} M0,${0.5 + e} L1,${0.5 + e} M.5,${0.5 + e} L.5,1` }, "╥": { 1: (i8, e) => `M0,.5 L1,.5 M${0.5 - i8},.5 L${0.5 - i8},1 M${0.5 + i8},.5 L${0.5 + i8},1` }, "╦": { 1: (i8, e) => `M0,${0.5 - e} L1,${0.5 - e} M0,${0.5 + e} L${0.5 - i8},${0.5 + e} L${0.5 - i8},1 M1,${0.5 + e} L${0.5 + i8},${0.5 + e} L${0.5 + i8},1` }, "╧": { 1: (i8, e) => `M.5,0 L.5,${0.5 - e} M0,${0.5 - e} L1,${0.5 - e} M0,${0.5 + e} L1,${0.5 + e}` }, "╨": { 1: (i8, e) => `M0,.5 L1,.5 M${0.5 - i8},.5 L${0.5 - i8},0 M${0.5 + i8},.5 L${0.5 + i8},0` }, "╩": { 1: (i8, e) => `M0,${0.5 + e} L1,${0.5 + e} M0,${0.5 - e} L${0.5 - i8},${0.5 - e} L${0.5 - i8},0 M1,${0.5 - e} L${0.5 + i8},${0.5 - e} L${0.5 + i8},0` }, "╪": { 1: (i8, e) => `M.5,0 L.5,1 M0,${0.5 - e} L1,${0.5 - e} M0,${0.5 + e} L1,${0.5 + e}` }, "╫": { 1: (i8, e) => `M0,.5 L1,.5 M${0.5 - i8},0 L${0.5 - i8},1 M${0.5 + i8},0 L${0.5 + i8},1` }, "╬": { 1: (i8, e) => `M0,${0.5 + e} L${0.5 - i8},${0.5 + e} L${0.5 - i8},1 M1,${0.5 + e} L${0.5 + i8},${0.5 + e} L${0.5 + i8},1 M0,${0.5 - e} L${0.5 - i8},${0.5 - e} L${0.5 - i8},0 M1,${0.5 - e} L${0.5 + i8},${0.5 - e} L${0.5 + i8},0` }, "╱": { 1: "M1,0 L0,1" }, "╲": { 1: "M0,0 L1,1" }, "╳": { 1: "M1,0 L0,1 M0,0 L1,1" }, "╼": { 1: "M.5,.5 L0,.5", 3: "M.5,.5 L1,.5" }, "╽": { 1: "M.5,.5 L.5,0", 3: "M.5,.5 L.5,1" }, "╾": { 1: "M.5,.5 L1,.5", 3: "M.5,.5 L0,.5" }, "╿": { 1: "M.5,.5 L.5,1", 3: "M.5,.5 L.5,0" }, "┍": { 1: "M.5,.5 L.5,1", 3: "M.5,.5 L1,.5" }, "┎": { 1: "M.5,.5 L1,.5", 3: "M.5,.5 L.5,1" }, "┑": { 1: "M.5,.5 L.5,1", 3: "M.5,.5 L0,.5" }, "┒": { 1: "M.5,.5 L0,.5", 3: "M.5,.5 L.5,1" }, "┕": { 1: "M.5,.5 L.5,0", 3: "M.5,.5 L1,.5" }, "┖": { 1: "M.5,.5 L1,.5", 3: "M.5,.5 L.5,0" }, "┙": { 1: "M.5,.5 L.5,0", 3: "M.5,.5 L0,.5" }, "┚": { 1: "M.5,.5 L0,.5", 3: "M.5,.5 L.5,0" }, "┝": { 1: "M.5,0 L.5,1", 3: "M.5,.5 L1,.5" }, "┞": { 1: "M0.5,1 L.5,.5 L1,.5", 3: "M.5,.5 L.5,0" }, "┟": { 1: "M.5,0 L.5,.5 L1,.5", 3: "M.5,.5 L.5,1" }, "┠": { 1: "M.5,.5 L1,.5", 3: "M.5,0 L.5,1" }, "┡": { 1: "M.5,.5 L.5,1", 3: "M.5,0 L.5,.5 L1,.5" }, "┢": { 1: "M.5,.5 L.5,0", 3: "M0.5,1 L.5,.5 L1,.5" }, "┥": { 1: "M.5,0 L.5,1", 3: "M.5,.5 L0,.5" }, "┦": { 1: "M0,.5 L.5,.5 L.5,1", 3: "M.5,.5 L.5,0" }, "┧": { 1: "M.5,0 L.5,.5 L0,.5", 3: "M.5,.5 L.5,1" }, "┨": { 1: "M.5,.5 L0,.5", 3: "M.5,0 L.5,1" }, "┩": { 1: "M.5,.5 L.5,1", 3: "M.5,0 L.5,.5 L0,.5" }, "┪": { 1: "M.5,.5 L.5,0", 3: "M0,.5 L.5,.5 L.5,1" }, "┭": { 1: "M0.5,1 L.5,.5 L1,.5", 3: "M.5,.5 L0,.5" }, "┮": { 1: "M0,.5 L.5,.5 L.5,1", 3: "M.5,.5 L1,.5" }, "┯": { 1: "M.5,.5 L.5,1", 3: "M0,.5 L1,.5" }, "┰": { 1: "M0,.5 L1,.5", 3: "M.5,.5 L.5,1" }, "┱": { 1: "M.5,.5 L1,.5", 3: "M0,.5 L.5,.5 L.5,1" }, "┲": { 1: "M.5,.5 L0,.5", 3: "M0.5,1 L.5,.5 L1,.5" }, "┵": { 1: "M.5,0 L.5,.5 L1,.5", 3: "M.5,.5 L0,.5" }, "┶": { 1: "M.5,0 L.5,.5 L0,.5", 3: "M.5,.5 L1,.5" }, "┷": { 1: "M.5,.5 L.5,0", 3: "M0,.5 L1,.5" }, "┸": { 1: "M0,.5 L1,.5", 3: "M.5,.5 L.5,0" }, "┹": { 1: "M.5,.5 L1,.5", 3: "M.5,0 L.5,.5 L0,.5" }, "┺": { 1: "M.5,.5 L0,.5", 3: "M.5,0 L.5,.5 L1,.5" }, "┽": { 1: "M.5,0 L.5,1 M.5,.5 L1,.5", 3: "M.5,.5 L0,.5" }, "┾": { 1: "M.5,0 L.5,1 M.5,.5 L0,.5", 3: "M.5,.5 L1,.5" }, "┿": { 1: "M.5,0 L.5,1", 3: "M0,.5 L1,.5" }, "╀": { 1: "M0,.5 L1,.5 M.5,.5 L.5,1", 3: "M.5,.5 L.5,0" }, "╁": { 1: "M.5,.5 L.5,0 M0,.5 L1,.5", 3: "M.5,.5 L.5,1" }, "╂": { 1: "M0,.5 L1,.5", 3: "M.5,0 L.5,1" }, "╃": { 1: "M0.5,1 L.5,.5 L1,.5", 3: "M.5,0 L.5,.5 L0,.5" }, "╄": { 1: "M0,.5 L.5,.5 L.5,1", 3: "M.5,0 L.5,.5 L1,.5" }, "╅": { 1: "M.5,0 L.5,.5 L1,.5", 3: "M0,.5 L.5,.5 L.5,1" }, "╆": { 1: "M.5,0 L.5,.5 L0,.5", 3: "M0.5,1 L.5,.5 L1,.5" }, "╇": { 1: "M.5,.5 L.5,1", 3: "M.5,.5 L.5,0 M0,.5 L1,.5" }, "╈": { 1: "M.5,.5 L.5,0", 3: "M0,.5 L1,.5 M.5,.5 L.5,1" }, "╉": { 1: "M.5,.5 L1,.5", 3: "M.5,0 L.5,1 M.5,.5 L0,.5" }, "╊": { 1: "M.5,.5 L0,.5", 3: "M.5,0 L.5,1 M.5,.5 L1,.5" }, "╌": { 1: "M.1,.5 L.4,.5 M.6,.5 L.9,.5" }, "╍": { 3: "M.1,.5 L.4,.5 M.6,.5 L.9,.5" }, "┄": { 1: "M.0667,.5 L.2667,.5 M.4,.5 L.6,.5 M.7333,.5 L.9333,.5" }, "┅": { 3: "M.0667,.5 L.2667,.5 M.4,.5 L.6,.5 M.7333,.5 L.9333,.5" }, "┈": { 1: "M.05,.5 L.2,.5 M.3,.5 L.45,.5 M.55,.5 L.7,.5 M.8,.5 L.95,.5" }, "┉": { 3: "M.05,.5 L.2,.5 M.3,.5 L.45,.5 M.55,.5 L.7,.5 M.8,.5 L.95,.5" }, "╎": { 1: "M.5,.1 L.5,.4 M.5,.6 L.5,.9" }, "╏": { 3: "M.5,.1 L.5,.4 M.5,.6 L.5,.9" }, "┆": { 1: "M.5,.0667 L.5,.2667 M.5,.4 L.5,.6 M.5,.7333 L.5,.9333" }, "┇": { 3: "M.5,.0667 L.5,.2667 M.5,.4 L.5,.6 M.5,.7333 L.5,.9333" }, "┊": { 1: "M.5,.05 L.5,.2 M.5,.3 L.5,.45 L.5,.55 M.5,.7 L.5,.95" }, "┋": { 3: "M.5,.05 L.5,.2 M.5,.3 L.5,.45 L.5,.55 M.5,.7 L.5,.95" }, "╭": { 1: (i8, e) => `M.5,1 L.5,${0.5 + e / 0.15 * 0.5} C.5,${0.5 + e / 0.15 * 0.5},.5,.5,1,.5` }, "╮": { 1: (i8, e) => `M.5,1 L.5,${0.5 + e / 0.15 * 0.5} C.5,${0.5 + e / 0.15 * 0.5},.5,.5,0,.5` }, "╯": { 1: (i8, e) => `M.5,0 L.5,${0.5 - e / 0.15 * 0.5} C.5,${0.5 - e / 0.15 * 0.5},.5,.5,0,.5` }, "╰": { 1: (i8, e) => `M.5,0 L.5,${0.5 - e / 0.15 * 0.5} C.5,${0.5 - e / 0.15 * 0.5},.5,.5,1,.5` } };
var et2 = { "": { d: "M.3,1 L.03,1 L.03,.88 C.03,.82,.06,.78,.11,.73 C.15,.7,.2,.68,.28,.65 L.43,.6 C.49,.58,.53,.56,.56,.53 C.59,.5,.6,.47,.6,.43 L.6,.27 L.4,.27 L.69,.1 L.98,.27 L.78,.27 L.78,.46 C.78,.52,.76,.56,.72,.61 C.68,.66,.63,.67,.56,.7 L.48,.72 C.42,.74,.38,.76,.35,.78 C.32,.8,.31,.84,.31,.88 L.31,1 M.3,.5 L.03,.59 L.03,.09 L.3,.09 L.3,.655", type: 0 }, "": { d: "M.7,.4 L.7,.47 L.2,.47 L.2,.03 L.355,.03 L.355,.4 L.705,.4 M.7,.5 L.86,.5 L.86,.95 L.69,.95 L.44,.66 L.46,.86 L.46,.95 L.3,.95 L.3,.49 L.46,.49 L.71,.78 L.69,.565 L.69,.5", type: 0 }, "": { d: "M.25,.94 C.16,.94,.11,.92,.11,.87 L.11,.53 C.11,.48,.15,.455,.23,.45 L.23,.3 C.23,.25,.26,.22,.31,.19 C.36,.16,.43,.15,.51,.15 C.59,.15,.66,.16,.71,.19 C.77,.22,.79,.26,.79,.3 L.79,.45 C.87,.45,.91,.48,.91,.53 L.91,.87 C.91,.92,.86,.94,.77,.94 L.24,.94 M.53,.2 C.49,.2,.45,.21,.42,.23 C.39,.25,.38,.27,.38,.3 L.38,.45 L.68,.45 L.68,.3 C.68,.27,.67,.25,.64,.23 C.61,.21,.58,.2,.53,.2 M.58,.82 L.58,.66 C.63,.65,.65,.63,.65,.6 C.65,.58,.64,.57,.61,.56 C.58,.55,.56,.54,.52,.54 C.48,.54,.46,.55,.43,.56 C.4,.57,.39,.59,.39,.6 C.39,.63,.41,.64,.46,.66 L.46,.82 L.57,.82", type: 0 }, "": { d: "M0,0 L1,.5 L0,1", type: 0, rightPadding: 2 }, "": { d: "M-1,-.5 L1,.5 L-1,1.5", type: 1, leftPadding: 1, rightPadding: 1 }, "": { d: "M1,0 L0,.5 L1,1", type: 0, leftPadding: 2 }, "": { d: "M2,-.5 L0,.5 L2,1.5", type: 1, leftPadding: 1, rightPadding: 1 }, "": { d: "M0,0 L0,1 C0.552,1,1,0.776,1,.5 C1,0.224,0.552,0,0,0", type: 0, rightPadding: 1 }, "": { d: "M.2,1 C.422,1,.8,.826,.78,.5 C.8,.174,0.422,0,.2,0", type: 1, rightPadding: 1 }, "": { d: "M1,0 L1,1 C0.448,1,0,0.776,0,.5 C0,0.224,0.448,0,1,0", type: 0, leftPadding: 1 }, "": { d: "M.8,1 C0.578,1,0.2,.826,.22,.5 C0.2,0.174,0.578,0,0.8,0", type: 1, leftPadding: 1 }, "": { d: "M-.5,-.5 L1.5,1.5 L-.5,1.5", type: 0 }, "": { d: "M-.5,-.5 L1.5,1.5", type: 1, leftPadding: 1, rightPadding: 1 }, "": { d: "M1.5,-.5 L-.5,1.5 L1.5,1.5", type: 0 }, "": { d: "M1.5,-.5 L-.5,1.5 L-.5,-.5", type: 0 }, "": { d: "M1.5,-.5 L-.5,1.5", type: 1, leftPadding: 1, rightPadding: 1 }, "": { d: "M-.5,-.5 L1.5,1.5 L1.5,-.5", type: 0 } };
et2[""] = et2[""];
et2[""] = et2[""];
function yn(i8, e, t, n, s15, o2, r5, a) {
  let l = Hr[e];
  if (l) return $r(i8, l, t, n, s15, o2), true;
  let u = Wr[e];
  if (u) return Kr(i8, u, t, n, s15, o2), true;
  let c = Gr[e];
  if (c) return Vr(i8, c, t, n, s15, o2, a), true;
  let d = et2[e];
  return d ? (Cr(i8, d, t, n, s15, o2, r5, a), true) : false;
}
function $r(i8, e, t, n, s15, o2) {
  for (let r5 = 0; r5 < e.length; r5++) {
    let a = e[r5], l = s15 / 8, u = o2 / 8;
    i8.fillRect(t + a.x * l, n + a.y * u, a.w * l, a.h * u);
  }
}
var xn = /* @__PURE__ */ new Map();
function Kr(i8, e, t, n, s15, o2) {
  let r5 = xn.get(e);
  r5 || (r5 = /* @__PURE__ */ new Map(), xn.set(e, r5));
  let a = i8.fillStyle;
  if (typeof a != "string") throw new Error(`Unexpected fillStyle type "${a}"`);
  let l = r5.get(a);
  if (!l) {
    let u = e[0].length, c = e.length, d = i8.canvas.ownerDocument.createElement("canvas");
    d.width = u, d.height = c;
    let h2 = F(d.getContext("2d")), f = new ImageData(u, c), I2, L2, M4, q2;
    if (a.startsWith("#")) I2 = parseInt(a.slice(1, 3), 16), L2 = parseInt(a.slice(3, 5), 16), M4 = parseInt(a.slice(5, 7), 16), q2 = a.length > 7 && parseInt(a.slice(7, 9), 16) || 1;
    else if (a.startsWith("rgba")) [I2, L2, M4, q2] = a.substring(5, a.length - 1).split(",").map((S2) => parseFloat(S2));
    else throw new Error(`Unexpected fillStyle color format "${a}" when drawing pattern glyph`);
    for (let S2 = 0; S2 < c; S2++) for (let W = 0; W < u; W++) f.data[(S2 * u + W) * 4] = I2, f.data[(S2 * u + W) * 4 + 1] = L2, f.data[(S2 * u + W) * 4 + 2] = M4, f.data[(S2 * u + W) * 4 + 3] = e[S2][W] * (q2 * 255);
    h2.putImageData(f, 0, 0), l = F(i8.createPattern(d, null)), r5.set(a, l);
  }
  i8.fillStyle = l, i8.fillRect(t, n, s15, o2);
}
function Vr(i8, e, t, n, s15, o2, r5) {
  i8.strokeStyle = i8.fillStyle;
  for (let [a, l] of Object.entries(e)) {
    i8.beginPath(), i8.lineWidth = r5 * Number.parseInt(a);
    let u;
    if (typeof l == "function") {
      let d = 0.15 / o2 * s15;
      u = l(0.15, d);
    } else u = l;
    for (let c of u.split(" ")) {
      let d = c[0], h2 = In[d];
      if (!h2) {
        console.error(`Could not find drawing instructions for "${d}"`);
        continue;
      }
      let f = c.substring(1).split(",");
      !f[0] || !f[1] || h2(i8, Ln(f, s15, o2, t, n, true, r5));
    }
    i8.stroke(), i8.closePath();
  }
}
function Cr(i8, e, t, n, s15, o2, r5, a) {
  let l = new Path2D();
  l.rect(t, n, s15, o2), i8.clip(l), i8.beginPath();
  let u = r5 / 12;
  i8.lineWidth = a * u;
  for (let c of e.d.split(" ")) {
    let d = c[0], h2 = In[d];
    if (!h2) {
      console.error(`Could not find drawing instructions for "${d}"`);
      continue;
    }
    let f = c.substring(1).split(",");
    !f[0] || !f[1] || h2(i8, Ln(f, s15, o2, t, n, false, a, (e.leftPadding ?? 0) * (u / 2), (e.rightPadding ?? 0) * (u / 2)));
  }
  e.type === 1 ? (i8.strokeStyle = i8.fillStyle, i8.stroke()) : i8.fill(), i8.closePath();
}
function En(i8, e, t = 0) {
  return Math.max(Math.min(i8, e), t);
}
var In = { C: (i8, e) => i8.bezierCurveTo(e[0], e[1], e[2], e[3], e[4], e[5]), L: (i8, e) => i8.lineTo(e[0], e[1]), M: (i8, e) => i8.moveTo(e[0], e[1]) };
function Ln(i8, e, t, n, s15, o2, r5, a = 0, l = 0) {
  let u = i8.map((c) => parseFloat(c) || parseInt(c));
  if (u.length < 2) throw new Error("Too few arguments for instruction");
  for (let c = 0; c < u.length; c += 2) u[c] *= e - a * r5 - l * r5, o2 && u[c] !== 0 && (u[c] = En(Math.round(u[c] + 0.5) - 0.5, e, 0)), u[c] += n + a * r5;
  for (let c = 1; c < u.length; c += 2) u[c] *= t, o2 && u[c] !== 0 && (u[c] = En(Math.round(u[c] + 0.5) - 0.5, t, 0)), u[c] += s15;
  return u;
}
var Ot = class {
  constructor() {
    this._data = {};
  }
  set(e, t, n) {
    this._data[e] || (this._data[e] = {}), this._data[e][t] = n;
  }
  get(e, t) {
    return this._data[e] ? this._data[e][t] : void 0;
  }
  clear() {
    this._data = {};
  }
};
var tt2 = class {
  constructor() {
    this._data = new Ot();
  }
  set(e, t, n, s15, o2) {
    this._data.get(e, t) || this._data.set(e, t, new Ot()), this._data.get(e, t).set(n, s15, o2);
  }
  get(e, t, n, s15) {
    return this._data.get(e, t)?.get(n, s15);
  }
  clear() {
    this._data.clear();
  }
};
var Ft = class {
  constructor() {
    this._tasks = [];
    this._i = 0;
  }
  enqueue(e) {
    this._tasks.push(e), this._start();
  }
  flush() {
    for (; this._i < this._tasks.length; ) this._tasks[this._i]() || this._i++;
    this.clear();
  }
  clear() {
    this._idleCallback && (this._cancelCallback(this._idleCallback), this._idleCallback = void 0), this._i = 0, this._tasks.length = 0;
  }
  _start() {
    this._idleCallback || (this._idleCallback = this._requestCallback(this._process.bind(this)));
  }
  _process(e) {
    this._idleCallback = void 0;
    let t = 0, n = 0, s15 = e.timeRemaining(), o2 = 0;
    for (; this._i < this._tasks.length; ) {
      if (t = performance.now(), this._tasks[this._i]() || this._i++, t = Math.max(1, performance.now() - t), n = Math.max(t, n), o2 = e.timeRemaining(), n * 1.5 > o2) {
        s15 - t < -20 && console.warn(`task queue exceeded allotted deadline by ${Math.abs(Math.round(s15 - t))}ms`), this._start();
        return;
      }
      s15 = o2;
    }
    this.clear();
  }
};
var gi = class extends Ft {
  _requestCallback(e) {
    return setTimeout(() => e(this._createDeadline(16)));
  }
  _cancelCallback(e) {
    clearTimeout(e);
  }
  _createDeadline(e) {
    let t = performance.now() + e;
    return { timeRemaining: () => Math.max(0, t - performance.now()) };
  }
};
var xi = class extends Ft {
  _requestCallback(e) {
    return requestIdleCallback(e);
  }
  _cancelCallback(e) {
    cancelIdleCallback(e);
  }
};
var wn = !Lt && "requestIdleCallback" in window ? xi : gi;
var he2 = class i2 {
  constructor() {
    this.fg = 0;
    this.bg = 0;
    this.extended = new it2();
  }
  static toColorRGB(e) {
    return [e >>> 16 & 255, e >>> 8 & 255, e & 255];
  }
  static fromColorRGB(e) {
    return (e[0] & 255) << 16 | (e[1] & 255) << 8 | e[2] & 255;
  }
  clone() {
    let e = new i2();
    return e.fg = this.fg, e.bg = this.bg, e.extended = this.extended.clone(), e;
  }
  isInverse() {
    return this.fg & 67108864;
  }
  isBold() {
    return this.fg & 134217728;
  }
  isUnderline() {
    return this.hasExtendedAttrs() && this.extended.underlineStyle !== 0 ? 1 : this.fg & 268435456;
  }
  isBlink() {
    return this.fg & 536870912;
  }
  isInvisible() {
    return this.fg & 1073741824;
  }
  isItalic() {
    return this.bg & 67108864;
  }
  isDim() {
    return this.bg & 134217728;
  }
  isStrikethrough() {
    return this.fg & 2147483648;
  }
  isProtected() {
    return this.bg & 536870912;
  }
  isOverline() {
    return this.bg & 1073741824;
  }
  getFgColorMode() {
    return this.fg & 50331648;
  }
  getBgColorMode() {
    return this.bg & 50331648;
  }
  isFgRGB() {
    return (this.fg & 50331648) === 50331648;
  }
  isBgRGB() {
    return (this.bg & 50331648) === 50331648;
  }
  isFgPalette() {
    return (this.fg & 50331648) === 16777216 || (this.fg & 50331648) === 33554432;
  }
  isBgPalette() {
    return (this.bg & 50331648) === 16777216 || (this.bg & 50331648) === 33554432;
  }
  isFgDefault() {
    return (this.fg & 50331648) === 0;
  }
  isBgDefault() {
    return (this.bg & 50331648) === 0;
  }
  isAttributeDefault() {
    return this.fg === 0 && this.bg === 0;
  }
  getFgColor() {
    switch (this.fg & 50331648) {
      case 16777216:
      case 33554432:
        return this.fg & 255;
      case 50331648:
        return this.fg & 16777215;
      default:
        return -1;
    }
  }
  getBgColor() {
    switch (this.bg & 50331648) {
      case 16777216:
      case 33554432:
        return this.bg & 255;
      case 50331648:
        return this.bg & 16777215;
      default:
        return -1;
    }
  }
  hasExtendedAttrs() {
    return this.bg & 268435456;
  }
  updateExtended() {
    this.extended.isEmpty() ? this.bg &= -268435457 : this.bg |= 268435456;
  }
  getUnderlineColor() {
    if (this.bg & 268435456 && ~this.extended.underlineColor) switch (this.extended.underlineColor & 50331648) {
      case 16777216:
      case 33554432:
        return this.extended.underlineColor & 255;
      case 50331648:
        return this.extended.underlineColor & 16777215;
      default:
        return this.getFgColor();
    }
    return this.getFgColor();
  }
  getUnderlineColorMode() {
    return this.bg & 268435456 && ~this.extended.underlineColor ? this.extended.underlineColor & 50331648 : this.getFgColorMode();
  }
  isUnderlineColorRGB() {
    return this.bg & 268435456 && ~this.extended.underlineColor ? (this.extended.underlineColor & 50331648) === 50331648 : this.isFgRGB();
  }
  isUnderlineColorPalette() {
    return this.bg & 268435456 && ~this.extended.underlineColor ? (this.extended.underlineColor & 50331648) === 16777216 || (this.extended.underlineColor & 50331648) === 33554432 : this.isFgPalette();
  }
  isUnderlineColorDefault() {
    return this.bg & 268435456 && ~this.extended.underlineColor ? (this.extended.underlineColor & 50331648) === 0 : this.isFgDefault();
  }
  getUnderlineStyle() {
    return this.fg & 268435456 ? this.bg & 268435456 ? this.extended.underlineStyle : 1 : 0;
  }
  getUnderlineVariantOffset() {
    return this.extended.underlineVariantOffset;
  }
};
var it2 = class i3 {
  constructor(e = 0, t = 0) {
    this._ext = 0;
    this._urlId = 0;
    this._ext = e, this._urlId = t;
  }
  get ext() {
    return this._urlId ? this._ext & -469762049 | this.underlineStyle << 26 : this._ext;
  }
  set ext(e) {
    this._ext = e;
  }
  get underlineStyle() {
    return this._urlId ? 5 : (this._ext & 469762048) >> 26;
  }
  set underlineStyle(e) {
    this._ext &= -469762049, this._ext |= e << 26 & 469762048;
  }
  get underlineColor() {
    return this._ext & 67108863;
  }
  set underlineColor(e) {
    this._ext &= -67108864, this._ext |= e & 67108863;
  }
  get urlId() {
    return this._urlId;
  }
  set urlId(e) {
    this._urlId = e;
  }
  get underlineVariantOffset() {
    let e = (this._ext & 3758096384) >> 29;
    return e < 0 ? e ^ 4294967288 : e;
  }
  set underlineVariantOffset(e) {
    this._ext &= 536870911, this._ext |= e << 29 & 3758096384;
  }
  clone() {
    return new i3(this._ext, this._urlId);
  }
  isEmpty() {
    return this.underlineStyle === 0 && this._urlId === 0;
  }
};
var He2 = class He3 {
  constructor(e) {
    this.element = e, this.next = He3.Undefined, this.prev = He3.Undefined;
  }
};
He2.Undefined = new He2(void 0);
var zr = globalThis.performance && typeof globalThis.performance.now == "function";
var kt = class i4 {
  static create(e) {
    return new i4(e);
  }
  constructor(e) {
    this._now = zr && e === false ? Date.now : globalThis.performance.now.bind(globalThis.performance), this._startTime = this._now(), this._stopTime = -1;
  }
  stop() {
    this._stopTime = this._now();
  }
  reset() {
    this._startTime = this._now(), this._stopTime = -1;
  }
  elapsed() {
    return this._stopTime !== -1 ? this._stopTime - this._startTime : this._now() - this._startTime;
  }
};
var qr = false;
var Dn = false;
var jr = false;
var ee3;
((se3) => {
  se3.None = () => B.None;
  function e(v2) {
    if (jr) {
      let { onDidAddListener: p } = v2, g = nt2.create(), b2 = 0;
      v2.onDidAddListener = () => {
        ++b2 === 2 && (console.warn("snapshotted emitter LIKELY used public and SHOULD HAVE BEEN created with DisposableStore. snapshotted here"), g.print()), p?.();
      };
    }
  }
  function t(v2, p) {
    return h2(v2, () => {
    }, 0, void 0, true, void 0, p);
  }
  se3.defer = t;
  function n(v2) {
    return (p, g = null, b2) => {
      let m = false, _3;
      return _3 = v2((T) => {
        if (!m) return _3 ? _3.dispose() : m = true, p.call(g, T);
      }, null, b2), m && _3.dispose(), _3;
    };
  }
  se3.once = n;
  function s15(v2, p, g) {
    return c((b2, m = null, _3) => v2((T) => b2.call(m, p(T)), null, _3), g);
  }
  se3.map = s15;
  function o2(v2, p, g) {
    return c((b2, m = null, _3) => v2((T) => {
      p(T), b2.call(m, T);
    }, null, _3), g);
  }
  se3.forEach = o2;
  function r5(v2, p, g) {
    return c((b2, m = null, _3) => v2((T) => p(T) && b2.call(m, T), null, _3), g);
  }
  se3.filter = r5;
  function a(v2) {
    return v2;
  }
  se3.signal = a;
  function l(...v2) {
    return (p, g = null, b2) => {
      let m = It(...v2.map((_3) => _3((T) => p.call(g, T))));
      return d(m, b2);
    };
  }
  se3.any = l;
  function u(v2, p, g, b2) {
    let m = g;
    return s15(v2, (_3) => (m = p(m, _3), m), b2);
  }
  se3.reduce = u;
  function c(v2, p) {
    let g, b2 = { onWillAddFirstListener() {
      g = v2(m.fire, m);
    }, onDidRemoveLastListener() {
      g?.dispose();
    } };
    p || e(b2);
    let m = new D(b2);
    return p?.add(m), m.event;
  }
  function d(v2, p) {
    return p instanceof Array ? p.push(v2) : p && p.add(v2), v2;
  }
  function h2(v2, p, g = 100, b2 = false, m = false, _3, T) {
    let x, R3, $2, P3 = 0, de3, Re3 = { leakWarningThreshold: _3, onWillAddFirstListener() {
      x = v2((ie4) => {
        P3++, R3 = p(R3, ie4), b2 && !$2 && (oe.fire(R3), R3 = void 0), de3 = () => {
          let N2 = R3;
          R3 = void 0, $2 = void 0, (!b2 || P3 > 1) && oe.fire(N2), P3 = 0;
        }, typeof g == "number" ? (clearTimeout($2), $2 = setTimeout(de3, g)) : $2 === void 0 && ($2 = 0, queueMicrotask(de3));
      });
    }, onWillRemoveListener() {
      m && P3 > 0 && de3?.();
    }, onDidRemoveLastListener() {
      de3 = void 0, x.dispose();
    } };
    T || e(Re3);
    let oe = new D(Re3);
    return T?.add(oe), oe.event;
  }
  se3.debounce = h2;
  function f(v2, p = 0, g) {
    return se3.debounce(v2, (b2, m) => b2 ? (b2.push(m), b2) : [m], p, void 0, true, void 0, g);
  }
  se3.accumulate = f;
  function I2(v2, p = (b2, m) => b2 === m, g) {
    let b2 = true, m;
    return r5(v2, (_3) => {
      let T = b2 || !p(_3, m);
      return b2 = false, m = _3, T;
    }, g);
  }
  se3.latch = I2;
  function L2(v2, p, g) {
    return [se3.filter(v2, p, g), se3.filter(v2, (b2) => !p(b2), g)];
  }
  se3.split = L2;
  function M4(v2, p = false, g = [], b2) {
    let m = g.slice(), _3 = v2((R3) => {
      m ? m.push(R3) : x.fire(R3);
    });
    b2 && b2.add(_3);
    let T = () => {
      m?.forEach((R3) => x.fire(R3)), m = null;
    }, x = new D({ onWillAddFirstListener() {
      _3 || (_3 = v2((R3) => x.fire(R3)), b2 && b2.add(_3));
    }, onDidAddFirstListener() {
      m && (p ? setTimeout(T) : T());
    }, onDidRemoveLastListener() {
      _3 && _3.dispose(), _3 = null;
    } });
    return b2 && b2.add(x), x.event;
  }
  se3.buffer = M4;
  function q2(v2, p) {
    return (b2, m, _3) => {
      let T = p(new W());
      return v2(function(x) {
        let R3 = T.evaluate(x);
        R3 !== S2 && b2.call(m, R3);
      }, void 0, _3);
    };
  }
  se3.chain = q2;
  let S2 = /* @__PURE__ */ Symbol("HaltChainable");
  class W {
    constructor() {
      this.steps = [];
    }
    map(p) {
      return this.steps.push(p), this;
    }
    forEach(p) {
      return this.steps.push((g) => (p(g), g)), this;
    }
    filter(p) {
      return this.steps.push((g) => p(g) ? g : S2), this;
    }
    reduce(p, g) {
      let b2 = g;
      return this.steps.push((m) => (b2 = p(b2, m), b2)), this;
    }
    latch(p = (g, b2) => g === b2) {
      let g = true, b2;
      return this.steps.push((m) => {
        let _3 = g || !p(m, b2);
        return g = false, b2 = m, _3 ? m : S2;
      }), this;
    }
    evaluate(p) {
      for (let g of this.steps) if (p = g(p), p === S2) break;
      return p;
    }
  }
  function E(v2, p, g = (b2) => b2) {
    let b2 = (...x) => T.fire(g(...x)), m = () => v2.on(p, b2), _3 = () => v2.removeListener(p, b2), T = new D({ onWillAddFirstListener: m, onDidRemoveLastListener: _3 });
    return T.event;
  }
  se3.fromNodeEventEmitter = E;
  function y(v2, p, g = (b2) => b2) {
    let b2 = (...x) => T.fire(g(...x)), m = () => v2.addEventListener(p, b2), _3 = () => v2.removeEventListener(p, b2), T = new D({ onWillAddFirstListener: m, onDidRemoveLastListener: _3 });
    return T.event;
  }
  se3.fromDOMEventEmitter = y;
  function w2(v2) {
    return new Promise((p) => n(v2)(p));
  }
  se3.toPromise = w2;
  function G3(v2) {
    let p = new D();
    return v2.then((g) => {
      p.fire(g);
    }, () => {
      p.fire(void 0);
    }).finally(() => {
      p.dispose();
    }), p.event;
  }
  se3.fromPromise = G3;
  function ue3(v2, p) {
    return v2((g) => p.fire(g));
  }
  se3.forward = ue3;
  function Se3(v2, p, g) {
    return p(g), v2((b2) => p(b2));
  }
  se3.runAndSubscribe = Se3;
  class ce3 {
    constructor(p, g) {
      this._observable = p;
      this._counter = 0;
      this._hasChanged = false;
      let b2 = { onWillAddFirstListener: () => {
        p.addObserver(this);
      }, onDidRemoveLastListener: () => {
        p.removeObserver(this);
      } };
      g || e(b2), this.emitter = new D(b2), g && g.add(this.emitter);
    }
    beginUpdate(p) {
      this._counter++;
    }
    handlePossibleChange(p) {
    }
    handleChange(p, g) {
      this._hasChanged = true;
    }
    endUpdate(p) {
      this._counter--, this._counter === 0 && (this._observable.reportChanges(), this._hasChanged && (this._hasChanged = false, this.emitter.fire(this._observable.get())));
    }
  }
  function we3(v2, p) {
    return new ce3(v2, p).emitter.event;
  }
  se3.fromObservable = we3;
  function A2(v2) {
    return (p, g, b2) => {
      let m = 0, _3 = false, T = { beginUpdate() {
        m++;
      }, endUpdate() {
        m--, m === 0 && (v2.reportChanges(), _3 && (_3 = false, p.call(g)));
      }, handlePossibleChange() {
      }, handleChange() {
        _3 = true;
      } };
      v2.addObserver(T), v2.reportChanges();
      let x = { dispose() {
        v2.removeObserver(T);
      } };
      return b2 instanceof fe2 ? b2.add(x) : Array.isArray(b2) && b2.push(x), x;
    };
  }
  se3.fromObservableLight = A2;
})(ee3 ||= {});
var We2 = class We3 {
  constructor(e) {
    this.listenerCount = 0;
    this.invocationCount = 0;
    this.elapsedOverall = 0;
    this.durations = [];
    this.name = `${e}_${We3._idPool++}`, We3.all.add(this);
  }
  start(e) {
    this._stopWatch = new kt(), this.listenerCount = e;
  }
  stop() {
    if (this._stopWatch) {
      let e = this._stopWatch.elapsed();
      this.durations.push(e), this.elapsedOverall += e, this.invocationCount += 1, this._stopWatch = void 0;
    }
  }
};
We2.all = /* @__PURE__ */ new Set(), We2._idPool = 0;
var Ei = We2;
var Mn = -1;
var Bt = class Bt2 {
  constructor(e, t, n = (Bt2._idPool++).toString(16).padStart(3, "0")) {
    this._errorHandler = e;
    this.threshold = t;
    this.name = n;
    this._warnCountdown = 0;
  }
  dispose() {
    this._stacks?.clear();
  }
  check(e, t) {
    let n = this.threshold;
    if (n <= 0 || t < n) return;
    this._stacks || (this._stacks = /* @__PURE__ */ new Map());
    let s15 = this._stacks.get(e.value) || 0;
    if (this._stacks.set(e.value, s15 + 1), this._warnCountdown -= 1, this._warnCountdown <= 0) {
      this._warnCountdown = n * 0.5;
      let [o2, r5] = this.getMostFrequentStack(), a = `[${this.name}] potential listener LEAK detected, having ${t} listeners already. MOST frequent listener (${r5}):`;
      console.warn(a), console.warn(o2);
      let l = new Ii(a, o2);
      this._errorHandler(l);
    }
    return () => {
      let o2 = this._stacks.get(e.value) || 0;
      this._stacks.set(e.value, o2 - 1);
    };
  }
  getMostFrequentStack() {
    if (!this._stacks) return;
    let e, t = 0;
    for (let [n, s15] of this._stacks) (!e || t < s15) && (e = [n, s15], t = s15);
    return e;
  }
};
Bt._idPool = 1;
var yi = Bt;
var nt2 = class i5 {
  constructor(e) {
    this.value = e;
  }
  static create() {
    let e = new Error();
    return new i5(e.stack ?? "");
  }
  print() {
    console.warn(this.value.split(`
`).slice(2).join(`
`));
  }
};
var Ii = class extends Error {
  constructor(e, t) {
    super(e), this.name = "ListenerLeakError", this.stack = t;
  }
};
var Li = class extends Error {
  constructor(e, t) {
    super(e), this.name = "ListenerRefusalError", this.stack = t;
  }
};
var Xr = 0;
var Ge2 = class {
  constructor(e) {
    this.value = e;
    this.id = Xr++;
  }
};
var Yr = 2;
var Qr = (i8, e) => {
  if (i8 instanceof Ge2) e(i8);
  else for (let t = 0; t < i8.length; t++) {
    let n = i8[t];
    n && e(n);
  }
};
var Pt;
if (qr) {
  let i8 = [];
  setInterval(() => {
    i8.length !== 0 && (console.warn("[LEAKING LISTENERS] GC'ed these listeners that were NOT yet disposed:"), console.warn(i8.join(`
`)), i8.length = 0);
  }, 3e3), Pt = new FinalizationRegistry((e) => {
    typeof e == "string" && i8.push(e);
  });
}
var D = class {
  constructor(e) {
    this._size = 0;
    this._options = e, this._leakageMon = Mn > 0 || this._options?.leakWarningThreshold ? new yi(e?.onListenerError ?? Pe2, this._options?.leakWarningThreshold ?? Mn) : void 0, this._perfMon = this._options?._profName ? new Ei(this._options._profName) : void 0, this._deliveryQueue = this._options?.deliveryQueue;
  }
  dispose() {
    if (!this._disposed) {
      if (this._disposed = true, this._deliveryQueue?.current === this && this._deliveryQueue.reset(), this._listeners) {
        if (Dn) {
          let e = this._listeners;
          queueMicrotask(() => {
            Qr(e, (t) => t.stack?.print());
          });
        }
        this._listeners = void 0, this._size = 0;
      }
      this._options?.onDidRemoveLastListener?.(), this._leakageMon?.dispose();
    }
  }
  get event() {
    return this._event ??= (e, t, n) => {
      if (this._leakageMon && this._size > this._leakageMon.threshold ** 2) {
        let l = `[${this._leakageMon.name}] REFUSES to accept new listeners because it exceeded its threshold by far (${this._size} vs ${this._leakageMon.threshold})`;
        console.warn(l);
        let u = this._leakageMon.getMostFrequentStack() ?? ["UNKNOWN stack", -1], c = new Li(`${l}. HINT: Stack shows most frequent listener (${u[1]}-times)`, u[0]);
        return (this._options?.onListenerError || Pe2)(c), B.None;
      }
      if (this._disposed) return B.None;
      t && (e = e.bind(t));
      let s15 = new Ge2(e), o2, r5;
      this._leakageMon && this._size >= Math.ceil(this._leakageMon.threshold * 0.2) && (s15.stack = nt2.create(), o2 = this._leakageMon.check(s15.stack, this._size + 1)), Dn && (s15.stack = r5 ?? nt2.create()), this._listeners ? this._listeners instanceof Ge2 ? (this._deliveryQueue ??= new wi(), this._listeners = [this._listeners, s15]) : this._listeners.push(s15) : (this._options?.onWillAddFirstListener?.(this), this._listeners = s15, this._options?.onDidAddFirstListener?.(this)), this._size++;
      let a = O2(() => {
        Pt?.unregister(a), o2?.(), this._removeListener(s15);
      });
      if (n instanceof fe2 ? n.add(a) : Array.isArray(n) && n.push(a), Pt) {
        let l = new Error().stack.split(`
`).slice(2, 3).join(`
`).trim(), u = /(file:|vscode-file:\/\/vscode-app)?(\/[^:]*:\d+:\d+)/.exec(l);
        Pt.register(a, u?.[2] ?? l, a);
      }
      return a;
    }, this._event;
  }
  _removeListener(e) {
    if (this._options?.onWillRemoveListener?.(this), !this._listeners) return;
    if (this._size === 1) {
      this._listeners = void 0, this._options?.onDidRemoveLastListener?.(this), this._size = 0;
      return;
    }
    let t = this._listeners, n = t.indexOf(e);
    if (n === -1) throw console.log("disposed?", this._disposed), console.log("size?", this._size), console.log("arr?", JSON.stringify(this._listeners)), new Error("Attempted to dispose unknown listener");
    this._size--, t[n] = void 0;
    let s15 = this._deliveryQueue.current === this;
    if (this._size * Yr <= t.length) {
      let o2 = 0;
      for (let r5 = 0; r5 < t.length; r5++) t[r5] ? t[o2++] = t[r5] : s15 && (this._deliveryQueue.end--, o2 < this._deliveryQueue.i && this._deliveryQueue.i--);
      t.length = o2;
    }
  }
  _deliver(e, t) {
    if (!e) return;
    let n = this._options?.onListenerError || Pe2;
    if (!n) {
      e.value(t);
      return;
    }
    try {
      e.value(t);
    } catch (s15) {
      n(s15);
    }
  }
  _deliverQueue(e) {
    let t = e.current._listeners;
    for (; e.i < e.end; ) this._deliver(t[e.i++], e.value);
    e.reset();
  }
  fire(e) {
    if (this._deliveryQueue?.current && (this._deliverQueue(this._deliveryQueue), this._perfMon?.stop()), this._perfMon?.start(this._size), this._listeners) if (this._listeners instanceof Ge2) this._deliver(this._listeners, e);
    else {
      let t = this._deliveryQueue;
      t.enqueue(this, e, this._listeners.length), this._deliverQueue(t);
    }
    this._perfMon?.stop();
  }
  hasListeners() {
    return this._size > 0;
  }
};
var wi = class {
  constructor() {
    this.i = -1;
    this.end = 0;
  }
  enqueue(e, t, n) {
    this.i = 0, this.end = n, this.current = e, this.value = t;
  }
  reset() {
    this.i = this.end, this.current = void 0, this.value = void 0;
  }
};
var An = { texturePage: 0, texturePosition: { x: 0, y: 0 }, texturePositionClipSpace: { x: 0, y: 0 }, offset: { x: 0, y: 0 }, size: { x: 0, y: 0 }, sizeClipSpace: { x: 0, y: 0 } };
var rt2 = 2;
var st2;
var ae = class i6 {
  constructor(e, t, n) {
    this._document = e;
    this._config = t;
    this._unicodeService = n;
    this._didWarmUp = false;
    this._cacheMap = new tt2();
    this._cacheMapCombined = new tt2();
    this._pages = [];
    this._activePages = [];
    this._workBoundingBox = { top: 0, left: 0, bottom: 0, right: 0 };
    this._workAttributeData = new he2();
    this._textureSize = 512;
    this._onAddTextureAtlasCanvas = new D();
    this.onAddTextureAtlasCanvas = this._onAddTextureAtlasCanvas.event;
    this._onRemoveTextureAtlasCanvas = new D();
    this.onRemoveTextureAtlasCanvas = this._onRemoveTextureAtlasCanvas.event;
    this._requestClearModel = false;
    this._createNewPage(), this._tmpCanvas = Sn(e, this._config.deviceCellWidth * 4 + rt2 * 2, this._config.deviceCellHeight + rt2 * 2), this._tmpCtx = F(this._tmpCanvas.getContext("2d", { alpha: this._config.allowTransparency, willReadFrequently: true }));
  }
  get pages() {
    return this._pages;
  }
  dispose() {
    this._tmpCanvas.remove();
    for (let e of this.pages) e.canvas.remove();
    this._onAddTextureAtlasCanvas.dispose();
  }
  warmUp() {
    this._didWarmUp || (this._doWarmUp(), this._didWarmUp = true);
  }
  _doWarmUp() {
    let e = new wn();
    for (let t = 33; t < 126; t++) e.enqueue(() => {
      if (!this._cacheMap.get(t, 0, 0, 0)) {
        let n = this._drawToCache(t, 0, 0, 0, false, void 0);
        this._cacheMap.set(t, 0, 0, 0, n);
      }
    });
  }
  beginFrame() {
    return this._requestClearModel;
  }
  clearTexture() {
    if (!(this._pages[0].currentRow.x === 0 && this._pages[0].currentRow.y === 0)) {
      for (let e of this._pages) e.clear();
      this._cacheMap.clear(), this._cacheMapCombined.clear(), this._didWarmUp = false;
    }
  }
  _createNewPage() {
    if (i6.maxAtlasPages && this._pages.length >= Math.max(4, i6.maxAtlasPages)) {
      let t = this._pages.filter((u) => u.canvas.width * 2 <= (i6.maxTextureSize || 4096)).sort((u, c) => c.canvas.width !== u.canvas.width ? c.canvas.width - u.canvas.width : c.percentageUsed - u.percentageUsed), n = -1, s15 = 0;
      for (let u = 0; u < t.length; u++) if (t[u].canvas.width !== s15) n = u, s15 = t[u].canvas.width;
      else if (u - n === 3) break;
      let o2 = t.slice(n, n + 4), r5 = o2.map((u) => u.glyphs[0].texturePage).sort((u, c) => u > c ? 1 : -1), a = this.pages.length - o2.length, l = this._mergePages(o2, a);
      l.version++;
      for (let u = r5.length - 1; u >= 0; u--) this._deletePage(r5[u]);
      this.pages.push(l), this._requestClearModel = true, this._onAddTextureAtlasCanvas.fire(l.canvas);
    }
    let e = new ot2(this._document, this._textureSize);
    return this._pages.push(e), this._activePages.push(e), this._onAddTextureAtlasCanvas.fire(e.canvas), e;
  }
  _mergePages(e, t) {
    let n = e[0].canvas.width * 2, s15 = new ot2(this._document, n, e);
    for (let [o2, r5] of e.entries()) {
      let a = o2 * r5.canvas.width % n, l = Math.floor(o2 / 2) * r5.canvas.height;
      s15.ctx.drawImage(r5.canvas, a, l);
      for (let c of r5.glyphs) c.texturePage = t, c.sizeClipSpace.x = c.size.x / n, c.sizeClipSpace.y = c.size.y / n, c.texturePosition.x += a, c.texturePosition.y += l, c.texturePositionClipSpace.x = c.texturePosition.x / n, c.texturePositionClipSpace.y = c.texturePosition.y / n;
      this._onRemoveTextureAtlasCanvas.fire(r5.canvas);
      let u = this._activePages.indexOf(r5);
      u !== -1 && this._activePages.splice(u, 1);
    }
    return s15;
  }
  _deletePage(e) {
    this._pages.splice(e, 1);
    for (let t = e; t < this._pages.length; t++) {
      let n = this._pages[t];
      for (let s15 of n.glyphs) s15.texturePage--;
      n.version++;
    }
  }
  getRasterizedGlyphCombinedChar(e, t, n, s15, o2, r5) {
    return this._getFromCacheMap(this._cacheMapCombined, e, t, n, s15, o2, r5);
  }
  getRasterizedGlyph(e, t, n, s15, o2, r5) {
    return this._getFromCacheMap(this._cacheMap, e, t, n, s15, o2, r5);
  }
  _getFromCacheMap(e, t, n, s15, o2, r5, a) {
    return st2 = e.get(t, n, s15, o2), st2 || (st2 = this._drawToCache(t, n, s15, o2, r5, a), e.set(t, n, s15, o2, st2)), st2;
  }
  _getColorFromAnsiIndex(e) {
    if (e >= this._config.colors.ansi.length) throw new Error("No color found for idx " + e);
    return this._config.colors.ansi[e];
  }
  _getBackgroundColor(e, t, n, s15) {
    if (this._config.allowTransparency) return Z2;
    let o2;
    switch (e) {
      case 16777216:
      case 33554432:
        o2 = this._getColorFromAnsiIndex(t);
        break;
      case 50331648:
        let r5 = he2.toColorRGB(t);
        o2 = X3.toColor(r5[0], r5[1], r5[2]);
        break;
      case 0:
      default:
        n ? o2 = Ue.opaque(this._config.colors.foreground) : o2 = this._config.colors.background;
        break;
    }
    return this._config.allowTransparency || (o2 = Ue.opaque(o2)), o2;
  }
  _getForegroundColor(e, t, n, s15, o2, r5, a, l, u, c) {
    let d = this._getMinimumContrastColor(e, t, n, s15, o2, r5, a, u, l, c);
    if (d) return d;
    let h2;
    switch (o2) {
      case 16777216:
      case 33554432:
        this._config.drawBoldTextInBrightColors && u && r5 < 8 && (r5 += 8), h2 = this._getColorFromAnsiIndex(r5);
        break;
      case 50331648:
        let f = he2.toColorRGB(r5);
        h2 = X3.toColor(f[0], f[1], f[2]);
        break;
      case 0:
      default:
        a ? h2 = this._config.colors.background : h2 = this._config.colors.foreground;
    }
    return this._config.allowTransparency && (h2 = Ue.opaque(h2)), l && (h2 = Ue.multiplyOpacity(h2, gn)), h2;
  }
  _resolveBackgroundRgba(e, t, n) {
    switch (e) {
      case 16777216:
      case 33554432:
        return this._getColorFromAnsiIndex(t).rgba;
      case 50331648:
        return t << 8;
      case 0:
      default:
        return n ? this._config.colors.foreground.rgba : this._config.colors.background.rgba;
    }
  }
  _resolveForegroundRgba(e, t, n, s15) {
    switch (e) {
      case 16777216:
      case 33554432:
        return this._config.drawBoldTextInBrightColors && s15 && t < 8 && (t += 8), this._getColorFromAnsiIndex(t).rgba;
      case 50331648:
        return t << 8;
      case 0:
      default:
        return n ? this._config.colors.background.rgba : this._config.colors.foreground.rgba;
    }
  }
  _getMinimumContrastColor(e, t, n, s15, o2, r5, a, l, u, c) {
    if (this._config.minimumContrastRatio === 1 || c) return;
    let d = this._getContrastCache(u), h2 = d.getColor(e, s15);
    if (h2 !== void 0) return h2 || void 0;
    let f = this._resolveBackgroundRgba(t, n, a), I2 = this._resolveForegroundRgba(o2, r5, a, l), L2 = Te2.ensureContrastRatio(f, I2, this._config.minimumContrastRatio / (u ? 2 : 1));
    if (!L2) {
      d.setColor(e, s15, null);
      return;
    }
    let M4 = X3.toColor(L2 >> 24 & 255, L2 >> 16 & 255, L2 >> 8 & 255);
    return d.setColor(e, s15, M4), M4;
  }
  _getContrastCache(e) {
    return e ? this._config.colors.halfContrastCache : this._config.colors.contrastCache;
  }
  _drawToCache(e, t, n, s15, o2, r5) {
    let a = typeof e == "number" ? String.fromCharCode(e) : e;
    r5 && this._tmpCanvas.parentElement !== r5 && (this._tmpCanvas.style.display = "none", r5.append(this._tmpCanvas));
    let l = Math.min(this._config.deviceCellWidth * Math.max(a.length, 2) + rt2 * 2, this._config.deviceMaxTextureSize);
    this._tmpCanvas.width < l && (this._tmpCanvas.width = l);
    let u = Math.min(this._config.deviceCellHeight + rt2 * 4, this._textureSize);
    if (this._tmpCanvas.height < u && (this._tmpCanvas.height = u), this._tmpCtx.save(), this._workAttributeData.fg = n, this._workAttributeData.bg = t, this._workAttributeData.extended.ext = s15, !!this._workAttributeData.isInvisible()) return An;
    let d = !!this._workAttributeData.isBold(), h2 = !!this._workAttributeData.isInverse(), f = !!this._workAttributeData.isDim(), I2 = !!this._workAttributeData.isItalic(), L2 = !!this._workAttributeData.isUnderline(), M4 = !!this._workAttributeData.isStrikethrough(), q2 = !!this._workAttributeData.isOverline(), S2 = this._workAttributeData.getFgColor(), W = this._workAttributeData.getFgColorMode(), E = this._workAttributeData.getBgColor(), y = this._workAttributeData.getBgColorMode();
    if (h2) {
      let x = S2;
      S2 = E, E = x;
      let R3 = W;
      W = y, y = R3;
    }
    let w2 = this._getBackgroundColor(y, E, h2, f);
    this._tmpCtx.globalCompositeOperation = "copy", this._tmpCtx.fillStyle = w2.css, this._tmpCtx.fillRect(0, 0, this._tmpCanvas.width, this._tmpCanvas.height), this._tmpCtx.globalCompositeOperation = "source-over";
    let G3 = d ? this._config.fontWeightBold : this._config.fontWeight, ue3 = I2 ? "italic" : "";
    this._tmpCtx.font = `${ue3} ${G3} ${this._config.fontSize * this._config.devicePixelRatio}px ${this._config.fontFamily}`, this._tmpCtx.textBaseline = St;
    let Se3 = a.length === 1 && Rt(a.charCodeAt(0)), ce3 = a.length === 1 && fn(a.charCodeAt(0)), we3 = this._getForegroundColor(t, y, E, n, W, S2, h2, f, d, Dt(a.charCodeAt(0)));
    this._tmpCtx.fillStyle = we3.css;
    let A2 = ce3 ? 0 : rt2 * 2, se3 = false;
    this._config.customGlyphs !== false && (se3 = yn(this._tmpCtx, a, A2, A2, this._config.deviceCellWidth, this._config.deviceCellHeight, this._config.fontSize, this._config.devicePixelRatio));
    let v2 = !Se3, p;
    if (typeof e == "number" ? p = this._unicodeService.wcwidth(e) : p = this._unicodeService.getStringCellWidth(e), L2) {
      this._tmpCtx.save();
      let x = Math.max(1, Math.floor(this._config.fontSize * this._config.devicePixelRatio / 15)), R3 = x % 2 === 1 ? 0.5 : 0;
      if (this._tmpCtx.lineWidth = x, this._workAttributeData.isUnderlineColorDefault()) this._tmpCtx.strokeStyle = this._tmpCtx.fillStyle;
      else if (this._workAttributeData.isUnderlineColorRGB()) v2 = false, this._tmpCtx.strokeStyle = `rgb(${he2.toColorRGB(this._workAttributeData.getUnderlineColor()).join(",")})`;
      else {
        v2 = false;
        let ie4 = this._workAttributeData.getUnderlineColor();
        this._config.drawBoldTextInBrightColors && this._workAttributeData.isBold() && ie4 < 8 && (ie4 += 8), this._tmpCtx.strokeStyle = this._getColorFromAnsiIndex(ie4).css;
      }
      this._tmpCtx.beginPath();
      let $2 = A2, P3 = Math.ceil(A2 + this._config.deviceCharHeight) - R3 - (o2 ? x * 2 : 0), de3 = P3 + x, Re3 = P3 + x * 2, oe = this._workAttributeData.getUnderlineVariantOffset();
      for (let ie4 = 0; ie4 < p; ie4++) {
        this._tmpCtx.save();
        let N2 = $2 + ie4 * this._config.deviceCellWidth, ne3 = $2 + (ie4 + 1) * this._config.deviceCellWidth, di = N2 + this._config.deviceCellWidth / 2;
        switch (this._workAttributeData.extended.underlineStyle) {
          case 2:
            this._tmpCtx.moveTo(N2, P3), this._tmpCtx.lineTo(ne3, P3), this._tmpCtx.moveTo(N2, Re3), this._tmpCtx.lineTo(ne3, Re3);
            break;
          case 3:
            let ft2 = x <= 1 ? Re3 : Math.ceil(A2 + this._config.deviceCharHeight - x / 2) - R3, mt2 = x <= 1 ? P3 : Math.ceil(A2 + this._config.deviceCharHeight + x / 2) - R3, qi2 = new Path2D();
            qi2.rect(N2, P3, this._config.deviceCellWidth, Re3 - P3), this._tmpCtx.clip(qi2), this._tmpCtx.moveTo(N2 - this._config.deviceCellWidth / 2, de3), this._tmpCtx.bezierCurveTo(N2 - this._config.deviceCellWidth / 2, mt2, N2, mt2, N2, de3), this._tmpCtx.bezierCurveTo(N2, ft2, di, ft2, di, de3), this._tmpCtx.bezierCurveTo(di, mt2, ne3, mt2, ne3, de3), this._tmpCtx.bezierCurveTo(ne3, ft2, ne3 + this._config.deviceCellWidth / 2, ft2, ne3 + this._config.deviceCellWidth / 2, de3);
            break;
          case 4:
            let _t2 = oe === 0 ? 0 : oe >= x ? x * 2 - oe : x - oe;
            !(oe >= x) === false || _t2 === 0 ? (this._tmpCtx.setLineDash([Math.round(x), Math.round(x)]), this._tmpCtx.moveTo(N2 + _t2, P3), this._tmpCtx.lineTo(ne3, P3)) : (this._tmpCtx.setLineDash([Math.round(x), Math.round(x)]), this._tmpCtx.moveTo(N2, P3), this._tmpCtx.lineTo(N2 + _t2, P3), this._tmpCtx.moveTo(N2 + _t2 + x, P3), this._tmpCtx.lineTo(ne3, P3)), oe = bn(ne3 - N2, x, oe);
            break;
          case 5:
            let Er = 0.6, yr2 = 0.3, hi2 = ne3 - N2, ji2 = Math.floor(Er * hi2), Xi2 = Math.floor(yr2 * hi2), Ir2 = hi2 - ji2 - Xi2;
            this._tmpCtx.setLineDash([ji2, Xi2, Ir2]), this._tmpCtx.moveTo(N2, P3), this._tmpCtx.lineTo(ne3, P3);
            break;
          case 1:
          default:
            this._tmpCtx.moveTo(N2, P3), this._tmpCtx.lineTo(ne3, P3);
            break;
        }
        this._tmpCtx.stroke(), this._tmpCtx.restore();
      }
      if (this._tmpCtx.restore(), !se3 && this._config.fontSize >= 12 && !this._config.allowTransparency && a !== " ") {
        this._tmpCtx.save(), this._tmpCtx.textBaseline = "alphabetic";
        let ie4 = this._tmpCtx.measureText(a);
        if (this._tmpCtx.restore(), "actualBoundingBoxDescent" in ie4 && ie4.actualBoundingBoxDescent > 0) {
          this._tmpCtx.save();
          let N2 = new Path2D();
          N2.rect($2, P3 - Math.ceil(x / 2), this._config.deviceCellWidth * p, Re3 - P3 + Math.ceil(x / 2)), this._tmpCtx.clip(N2), this._tmpCtx.lineWidth = this._config.devicePixelRatio * 3, this._tmpCtx.strokeStyle = w2.css, this._tmpCtx.strokeText(a, A2, A2 + this._config.deviceCharHeight), this._tmpCtx.restore();
        }
      }
    }
    if (q2) {
      let x = Math.max(1, Math.floor(this._config.fontSize * this._config.devicePixelRatio / 15)), R3 = x % 2 === 1 ? 0.5 : 0;
      this._tmpCtx.lineWidth = x, this._tmpCtx.strokeStyle = this._tmpCtx.fillStyle, this._tmpCtx.beginPath(), this._tmpCtx.moveTo(A2, A2 + R3), this._tmpCtx.lineTo(A2 + this._config.deviceCharWidth * p, A2 + R3), this._tmpCtx.stroke();
    }
    if (se3 || this._tmpCtx.fillText(a, A2, A2 + this._config.deviceCharHeight), a === "_" && !this._config.allowTransparency) {
      let x = Di(this._tmpCtx.getImageData(A2, A2, this._config.deviceCellWidth, this._config.deviceCellHeight), w2, we3, v2);
      if (x) for (let R3 = 1; R3 <= 5 && (this._tmpCtx.save(), this._tmpCtx.fillStyle = w2.css, this._tmpCtx.fillRect(0, 0, this._tmpCanvas.width, this._tmpCanvas.height), this._tmpCtx.restore(), this._tmpCtx.fillText(a, A2, A2 + this._config.deviceCharHeight - R3), x = Di(this._tmpCtx.getImageData(A2, A2, this._config.deviceCellWidth, this._config.deviceCellHeight), w2, we3, v2), !!x); R3++) ;
    }
    if (M4) {
      let x = Math.max(1, Math.floor(this._config.fontSize * this._config.devicePixelRatio / 10)), R3 = this._tmpCtx.lineWidth % 2 === 1 ? 0.5 : 0;
      this._tmpCtx.lineWidth = x, this._tmpCtx.strokeStyle = this._tmpCtx.fillStyle, this._tmpCtx.beginPath(), this._tmpCtx.moveTo(A2, A2 + Math.floor(this._config.deviceCharHeight / 2) - R3), this._tmpCtx.lineTo(A2 + this._config.deviceCharWidth * p, A2 + Math.floor(this._config.deviceCharHeight / 2) - R3), this._tmpCtx.stroke();
    }
    this._tmpCtx.restore();
    let g = this._tmpCtx.getImageData(0, 0, this._tmpCanvas.width, this._tmpCanvas.height), b2;
    if (this._config.allowTransparency ? b2 = Jr(g) : b2 = Di(g, w2, we3, v2), b2) return An;
    let m = this._findGlyphBoundingBox(g, this._workBoundingBox, l, ce3, se3, A2), _3, T;
    for (; ; ) {
      if (this._activePages.length === 0) {
        let x = this._createNewPage();
        _3 = x, T = x.currentRow, T.height = m.size.y;
        break;
      }
      _3 = this._activePages[this._activePages.length - 1], T = _3.currentRow;
      for (let x of this._activePages) m.size.y <= x.currentRow.height && (_3 = x, T = x.currentRow);
      for (let x = this._activePages.length - 1; x >= 0; x--) for (let R3 of this._activePages[x].fixedRows) R3.height <= T.height && m.size.y <= R3.height && (_3 = this._activePages[x], T = R3);
      if (m.size.x > this._textureSize) {
        this._overflowSizePage || (this._overflowSizePage = new ot2(this._document, this._config.deviceMaxTextureSize), this.pages.push(this._overflowSizePage), this._requestClearModel = true, this._onAddTextureAtlasCanvas.fire(this._overflowSizePage.canvas)), _3 = this._overflowSizePage, T = this._overflowSizePage.currentRow, T.x + m.size.x >= _3.canvas.width && (T.x = 0, T.y += T.height, T.height = 0);
        break;
      }
      if (T.y + m.size.y >= _3.canvas.height || T.height > m.size.y + 2) {
        let x = false;
        if (_3.currentRow.y + _3.currentRow.height + m.size.y >= _3.canvas.height) {
          let R3;
          for (let $2 of this._activePages) if ($2.currentRow.y + $2.currentRow.height + m.size.y < $2.canvas.height) {
            R3 = $2;
            break;
          }
          if (R3) _3 = R3;
          else if (i6.maxAtlasPages && this._pages.length >= i6.maxAtlasPages && T.y + m.size.y <= _3.canvas.height && T.height >= m.size.y && T.x + m.size.x <= _3.canvas.width) x = true;
          else {
            let $2 = this._createNewPage();
            _3 = $2, T = $2.currentRow, T.height = m.size.y, x = true;
          }
        }
        x || (_3.currentRow.height > 0 && _3.fixedRows.push(_3.currentRow), T = { x: 0, y: _3.currentRow.y + _3.currentRow.height, height: m.size.y }, _3.fixedRows.push(T), _3.currentRow = { x: 0, y: T.y + T.height, height: 0 });
      }
      if (T.x + m.size.x <= _3.canvas.width) break;
      T === _3.currentRow ? (T.x = 0, T.y += T.height, T.height = 0) : _3.fixedRows.splice(_3.fixedRows.indexOf(T), 1);
    }
    return m.texturePage = this._pages.indexOf(_3), m.texturePosition.x = T.x, m.texturePosition.y = T.y, m.texturePositionClipSpace.x = T.x / _3.canvas.width, m.texturePositionClipSpace.y = T.y / _3.canvas.height, m.sizeClipSpace.x /= _3.canvas.width, m.sizeClipSpace.y /= _3.canvas.height, T.height = Math.max(T.height, m.size.y), T.x += m.size.x, _3.ctx.putImageData(g, m.texturePosition.x - this._workBoundingBox.left, m.texturePosition.y - this._workBoundingBox.top, this._workBoundingBox.left, this._workBoundingBox.top, m.size.x, m.size.y), _3.addGlyph(m), _3.version++, m;
  }
  _findGlyphBoundingBox(e, t, n, s15, o2, r5) {
    t.top = 0;
    let a = s15 ? this._config.deviceCellHeight : this._tmpCanvas.height, l = s15 ? this._config.deviceCellWidth : n, u = false;
    for (let c = 0; c < a; c++) {
      for (let d = 0; d < l; d++) {
        let h2 = c * this._tmpCanvas.width * 4 + d * 4 + 3;
        if (e.data[h2] !== 0) {
          t.top = c, u = true;
          break;
        }
      }
      if (u) break;
    }
    t.left = 0, u = false;
    for (let c = 0; c < r5 + l; c++) {
      for (let d = 0; d < a; d++) {
        let h2 = d * this._tmpCanvas.width * 4 + c * 4 + 3;
        if (e.data[h2] !== 0) {
          t.left = c, u = true;
          break;
        }
      }
      if (u) break;
    }
    t.right = l, u = false;
    for (let c = r5 + l - 1; c >= r5; c--) {
      for (let d = 0; d < a; d++) {
        let h2 = d * this._tmpCanvas.width * 4 + c * 4 + 3;
        if (e.data[h2] !== 0) {
          t.right = c, u = true;
          break;
        }
      }
      if (u) break;
    }
    t.bottom = a, u = false;
    for (let c = a - 1; c >= 0; c--) {
      for (let d = 0; d < l; d++) {
        let h2 = c * this._tmpCanvas.width * 4 + d * 4 + 3;
        if (e.data[h2] !== 0) {
          t.bottom = c, u = true;
          break;
        }
      }
      if (u) break;
    }
    return { texturePage: 0, texturePosition: { x: 0, y: 0 }, texturePositionClipSpace: { x: 0, y: 0 }, size: { x: t.right - t.left + 1, y: t.bottom - t.top + 1 }, sizeClipSpace: { x: t.right - t.left + 1, y: t.bottom - t.top + 1 }, offset: { x: -t.left + r5 + (s15 || o2 ? Math.floor((this._config.deviceCellWidth - this._config.deviceCharWidth) / 2) : 0), y: -t.top + r5 + (s15 || o2 ? this._config.lineHeight === 1 ? 0 : Math.round((this._config.deviceCellHeight - this._config.deviceCharHeight) / 2) : 0) } };
  }
};
var ot2 = class {
  constructor(e, t, n) {
    this._usedPixels = 0;
    this._glyphs = [];
    this.version = 0;
    this.currentRow = { x: 0, y: 0, height: 0 };
    this.fixedRows = [];
    if (n) for (let s15 of n) this._glyphs.push(...s15.glyphs), this._usedPixels += s15._usedPixels;
    this.canvas = Sn(e, t, t), this.ctx = F(this.canvas.getContext("2d", { alpha: true }));
  }
  get percentageUsed() {
    return this._usedPixels / (this.canvas.width * this.canvas.height);
  }
  get glyphs() {
    return this._glyphs;
  }
  addGlyph(e) {
    this._glyphs.push(e), this._usedPixels += e.size.x * e.size.y;
  }
  clear() {
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height), this.currentRow.x = 0, this.currentRow.y = 0, this.currentRow.height = 0, this.fixedRows.length = 0, this.version++;
  }
};
function Di(i8, e, t, n) {
  let s15 = e.rgba >>> 24, o2 = e.rgba >>> 16 & 255, r5 = e.rgba >>> 8 & 255, a = t.rgba >>> 24, l = t.rgba >>> 16 & 255, u = t.rgba >>> 8 & 255, c = Math.floor((Math.abs(s15 - a) + Math.abs(o2 - l) + Math.abs(r5 - u)) / 12), d = true;
  for (let h2 = 0; h2 < i8.data.length; h2 += 4) i8.data[h2] === s15 && i8.data[h2 + 1] === o2 && i8.data[h2 + 2] === r5 || n && Math.abs(i8.data[h2] - s15) + Math.abs(i8.data[h2 + 1] - o2) + Math.abs(i8.data[h2 + 2] - r5) < c ? i8.data[h2 + 3] = 0 : d = false;
  return d;
}
function Jr(i8) {
  for (let e = 0; e < i8.data.length; e += 4) if (i8.data[e + 3] > 0) return false;
  return true;
}
function Sn(i8, e, t) {
  let n = i8.createElement("canvas");
  return n.width = e, n.height = t, n;
}
function On(i8, e, t, n, s15, o2, r5, a) {
  let l = { foreground: o2.foreground, background: o2.background, cursor: Z2, cursorAccent: Z2, selectionForeground: Z2, selectionBackgroundTransparent: Z2, selectionBackgroundOpaque: Z2, selectionInactiveBackgroundTransparent: Z2, selectionInactiveBackgroundOpaque: Z2, overviewRulerBorder: Z2, scrollbarSliderBackground: Z2, scrollbarSliderHoverBackground: Z2, scrollbarSliderActiveBackground: Z2, ansi: o2.ansi.slice(), contrastCache: o2.contrastCache, halfContrastCache: o2.halfContrastCache };
  return { customGlyphs: s15.customGlyphs, devicePixelRatio: r5, deviceMaxTextureSize: a, letterSpacing: s15.letterSpacing, lineHeight: s15.lineHeight, deviceCellWidth: i8, deviceCellHeight: e, deviceCharWidth: t, deviceCharHeight: n, fontFamily: s15.fontFamily, fontSize: s15.fontSize, fontWeight: s15.fontWeight, fontWeightBold: s15.fontWeightBold, allowTransparency: s15.allowTransparency, drawBoldTextInBrightColors: s15.drawBoldTextInBrightColors, minimumContrastRatio: s15.minimumContrastRatio, colors: l };
}
function Mi(i8, e) {
  for (let t = 0; t < i8.colors.ansi.length; t++) if (i8.colors.ansi[t].rgba !== e.colors.ansi[t].rgba) return false;
  return i8.devicePixelRatio === e.devicePixelRatio && i8.customGlyphs === e.customGlyphs && i8.lineHeight === e.lineHeight && i8.letterSpacing === e.letterSpacing && i8.fontFamily === e.fontFamily && i8.fontSize === e.fontSize && i8.fontWeight === e.fontWeight && i8.fontWeightBold === e.fontWeightBold && i8.allowTransparency === e.allowTransparency && i8.deviceCharWidth === e.deviceCharWidth && i8.deviceCharHeight === e.deviceCharHeight && i8.drawBoldTextInBrightColors === e.drawBoldTextInBrightColors && i8.minimumContrastRatio === e.minimumContrastRatio && i8.colors.foreground.rgba === e.colors.foreground.rgba && i8.colors.background.rgba === e.colors.background.rgba;
}
function Fn(i8) {
  return (i8 & 50331648) === 16777216 || (i8 & 50331648) === 33554432;
}
var le = [];
function Nt(i8, e, t, n, s15, o2, r5, a, l) {
  let u = On(n, s15, o2, r5, e, t, a, l);
  for (let h2 = 0; h2 < le.length; h2++) {
    let f = le[h2], I2 = f.ownedBy.indexOf(i8);
    if (I2 >= 0) {
      if (Mi(f.config, u)) return f.atlas;
      f.ownedBy.length === 1 ? (f.atlas.dispose(), le.splice(h2, 1)) : f.ownedBy.splice(I2, 1);
      break;
    }
  }
  for (let h2 = 0; h2 < le.length; h2++) {
    let f = le[h2];
    if (Mi(f.config, u)) return f.ownedBy.push(i8), f.atlas;
  }
  let c = i8._core, d = { atlas: new ae(document, u, c.unicodeService), config: u, ownedBy: [i8] };
  return le.push(d), d.atlas;
}
function Ai(i8) {
  for (let e = 0; e < le.length; e++) {
    let t = le[e].ownedBy.indexOf(i8);
    if (t !== -1) {
      le[e].ownedBy.length === 1 ? (le[e].atlas.dispose(), le.splice(e, 1)) : le[e].ownedBy.splice(t, 1);
      break;
    }
  }
}
var Ut = 600;
var Ht = class {
  constructor(e, t) {
    this._renderCallback = e;
    this._coreBrowserService = t;
    this.isCursorVisible = true, this._coreBrowserService.isFocused && this._restartInterval();
  }
  get isPaused() {
    return !(this._blinkStartTimeout || this._blinkInterval);
  }
  dispose() {
    this._blinkInterval && (this._coreBrowserService.window.clearInterval(this._blinkInterval), this._blinkInterval = void 0), this._blinkStartTimeout && (this._coreBrowserService.window.clearTimeout(this._blinkStartTimeout), this._blinkStartTimeout = void 0), this._animationFrame && (this._coreBrowserService.window.cancelAnimationFrame(this._animationFrame), this._animationFrame = void 0);
  }
  restartBlinkAnimation() {
    this.isPaused || (this._animationTimeRestarted = Date.now(), this.isCursorVisible = true, this._animationFrame || (this._animationFrame = this._coreBrowserService.window.requestAnimationFrame(() => {
      this._renderCallback(), this._animationFrame = void 0;
    })));
  }
  _restartInterval(e = Ut) {
    this._blinkInterval && (this._coreBrowserService.window.clearInterval(this._blinkInterval), this._blinkInterval = void 0), this._blinkStartTimeout = this._coreBrowserService.window.setTimeout(() => {
      if (this._animationTimeRestarted) {
        let t = Ut - (Date.now() - this._animationTimeRestarted);
        if (this._animationTimeRestarted = void 0, t > 0) {
          this._restartInterval(t);
          return;
        }
      }
      this.isCursorVisible = false, this._animationFrame = this._coreBrowserService.window.requestAnimationFrame(() => {
        this._renderCallback(), this._animationFrame = void 0;
      }), this._blinkInterval = this._coreBrowserService.window.setInterval(() => {
        if (this._animationTimeRestarted) {
          let t = Ut - (Date.now() - this._animationTimeRestarted);
          this._animationTimeRestarted = void 0, this._restartInterval(t);
          return;
        }
        this.isCursorVisible = !this.isCursorVisible, this._animationFrame = this._coreBrowserService.window.requestAnimationFrame(() => {
          this._renderCallback(), this._animationFrame = void 0;
        });
      }, Ut);
    }, e);
  }
  pause() {
    this.isCursorVisible = true, this._blinkInterval && (this._coreBrowserService.window.clearInterval(this._blinkInterval), this._blinkInterval = void 0), this._blinkStartTimeout && (this._coreBrowserService.window.clearTimeout(this._blinkStartTimeout), this._blinkStartTimeout = void 0), this._animationFrame && (this._coreBrowserService.window.cancelAnimationFrame(this._animationFrame), this._animationFrame = void 0);
  }
  resume() {
    this.pause(), this._animationTimeRestarted = void 0, this._restartInterval(), this.restartBlinkAnimation();
  }
};
function Si(i8, e, t) {
  let n = new e.ResizeObserver((s15) => {
    let o2 = s15.find((l) => l.target === i8);
    if (!o2) return;
    if (!("devicePixelContentBoxSize" in o2)) {
      n?.disconnect(), n = void 0;
      return;
    }
    let r5 = o2.devicePixelContentBoxSize[0].inlineSize, a = o2.devicePixelContentBoxSize[0].blockSize;
    r5 > 0 && a > 0 && t(r5, a);
  });
  try {
    n.observe(i8, { box: ["device-pixel-content-box"] });
  } catch {
    n.disconnect(), n = void 0;
  }
  return O2(() => n?.disconnect());
}
function kn(i8) {
  return i8 > 65535 ? (i8 -= 65536, String.fromCharCode((i8 >> 10) + 55296) + String.fromCharCode(i8 % 1024 + 56320)) : String.fromCharCode(i8);
}
var at2 = class i7 extends he2 {
  constructor() {
    super(...arguments);
    this.content = 0;
    this.fg = 0;
    this.bg = 0;
    this.extended = new it2();
    this.combinedData = "";
  }
  static fromCharData(t) {
    let n = new i7();
    return n.setFromCharData(t), n;
  }
  isCombined() {
    return this.content & 2097152;
  }
  getWidth() {
    return this.content >> 22;
  }
  getChars() {
    return this.content & 2097152 ? this.combinedData : this.content & 2097151 ? kn(this.content & 2097151) : "";
  }
  getCode() {
    return this.isCombined() ? this.combinedData.charCodeAt(this.combinedData.length - 1) : this.content & 2097151;
  }
  setFromCharData(t) {
    this.fg = t[0], this.bg = 0;
    let n = false;
    if (t[1].length > 2) n = true;
    else if (t[1].length === 2) {
      let s15 = t[1].charCodeAt(0);
      if (55296 <= s15 && s15 <= 56319) {
        let o2 = t[1].charCodeAt(1);
        56320 <= o2 && o2 <= 57343 ? this.content = (s15 - 55296) * 1024 + o2 - 56320 + 65536 | t[2] << 22 : n = true;
      } else n = true;
    } else this.content = t[1].charCodeAt(0) | t[2] << 22;
    n && (this.combinedData = t[1], this.content = 2097152 | t[2] << 22);
  }
  getAsCharData() {
    return [this.fg, this.getChars(), this.getWidth(), this.getCode()];
  }
};
var Gt = new Float32Array([2, 0, 0, 0, 0, -2, 0, 0, 0, 0, 1, 0, -1, 1, 0, 1]);
function $t(i8, e, t) {
  let n = F(i8.createProgram());
  if (i8.attachShader(n, F(Pn(i8, i8.VERTEX_SHADER, e))), i8.attachShader(n, F(Pn(i8, i8.FRAGMENT_SHADER, t))), i8.linkProgram(n), i8.getProgramParameter(n, i8.LINK_STATUS)) return n;
  console.error(i8.getProgramInfoLog(n)), i8.deleteProgram(n);
}
function Pn(i8, e, t) {
  let n = F(i8.createShader(e));
  if (i8.shaderSource(n, t), i8.compileShader(n), i8.getShaderParameter(n, i8.COMPILE_STATUS)) return n;
  console.error(i8.getShaderInfoLog(n)), i8.deleteShader(n);
}
function Bn(i8, e) {
  let t = Math.min(i8.length * 2, e), n = new Float32Array(t);
  for (let s15 = 0; s15 < i8.length; s15++) n[s15] = i8[s15];
  return n;
}
var Wt = class {
  constructor(e) {
    this.texture = e, this.version = -1;
  }
};
var is = `#version 300 es
layout (location = 0) in vec2 a_unitquad;
layout (location = 1) in vec2 a_cellpos;
layout (location = 2) in vec2 a_offset;
layout (location = 3) in vec2 a_size;
layout (location = 4) in float a_texpage;
layout (location = 5) in vec2 a_texcoord;
layout (location = 6) in vec2 a_texsize;

uniform mat4 u_projection;
uniform vec2 u_resolution;

out vec2 v_texcoord;
flat out int v_texpage;

void main() {
  vec2 zeroToOne = (a_offset / u_resolution) + a_cellpos + (a_unitquad * a_size);
  gl_Position = u_projection * vec4(zeroToOne, 0.0, 1.0);
  v_texpage = int(a_texpage);
  v_texcoord = a_texcoord + a_unitquad * a_texsize;
}`;
function ns(i8) {
  let e = "";
  for (let t = 1; t < i8; t++) e += ` else if (v_texpage == ${t}) { outColor = texture(u_texture[${t}], v_texcoord); }`;
  return `#version 300 es
precision lowp float;

in vec2 v_texcoord;
flat in int v_texpage;

uniform sampler2D u_texture[${i8}];

out vec4 outColor;

void main() {
  if (v_texpage == 0) {
    outColor = texture(u_texture[0], v_texcoord);
  } ${e}
}`;
}
var De2 = 11;
var Ve2 = De2 * Float32Array.BYTES_PER_ELEMENT;
var rs = 2;
var H2 = 0;
var k;
var Fi = 0;
var lt2 = 0;
var Kt = class extends B {
  constructor(t, n, s15, o2) {
    super();
    this._terminal = t;
    this._gl = n;
    this._dimensions = s15;
    this._optionsService = o2;
    this._activeBuffer = 0;
    this._vertices = { count: 0, attributes: new Float32Array(0), attributesBuffers: [new Float32Array(0), new Float32Array(0)] };
    let r5 = this._gl;
    ae.maxAtlasPages === void 0 && (ae.maxAtlasPages = Math.min(32, F(r5.getParameter(r5.MAX_TEXTURE_IMAGE_UNITS))), ae.maxTextureSize = F(r5.getParameter(r5.MAX_TEXTURE_SIZE))), this._program = F($t(r5, is, ns(ae.maxAtlasPages))), this._register(O2(() => r5.deleteProgram(this._program))), this._projectionLocation = F(r5.getUniformLocation(this._program, "u_projection")), this._resolutionLocation = F(r5.getUniformLocation(this._program, "u_resolution")), this._textureLocation = F(r5.getUniformLocation(this._program, "u_texture")), this._vertexArrayObject = r5.createVertexArray(), r5.bindVertexArray(this._vertexArrayObject);
    let a = new Float32Array([0, 0, 1, 0, 0, 1, 1, 1]), l = r5.createBuffer();
    this._register(O2(() => r5.deleteBuffer(l))), r5.bindBuffer(r5.ARRAY_BUFFER, l), r5.bufferData(r5.ARRAY_BUFFER, a, r5.STATIC_DRAW), r5.enableVertexAttribArray(0), r5.vertexAttribPointer(0, 2, this._gl.FLOAT, false, 0, 0);
    let u = new Uint8Array([0, 1, 2, 3]), c = r5.createBuffer();
    this._register(O2(() => r5.deleteBuffer(c))), r5.bindBuffer(r5.ELEMENT_ARRAY_BUFFER, c), r5.bufferData(r5.ELEMENT_ARRAY_BUFFER, u, r5.STATIC_DRAW), this._attributesBuffer = F(r5.createBuffer()), this._register(O2(() => r5.deleteBuffer(this._attributesBuffer))), r5.bindBuffer(r5.ARRAY_BUFFER, this._attributesBuffer), r5.enableVertexAttribArray(2), r5.vertexAttribPointer(2, 2, r5.FLOAT, false, Ve2, 0), r5.vertexAttribDivisor(2, 1), r5.enableVertexAttribArray(3), r5.vertexAttribPointer(3, 2, r5.FLOAT, false, Ve2, 2 * Float32Array.BYTES_PER_ELEMENT), r5.vertexAttribDivisor(3, 1), r5.enableVertexAttribArray(4), r5.vertexAttribPointer(4, 1, r5.FLOAT, false, Ve2, 4 * Float32Array.BYTES_PER_ELEMENT), r5.vertexAttribDivisor(4, 1), r5.enableVertexAttribArray(5), r5.vertexAttribPointer(5, 2, r5.FLOAT, false, Ve2, 5 * Float32Array.BYTES_PER_ELEMENT), r5.vertexAttribDivisor(5, 1), r5.enableVertexAttribArray(6), r5.vertexAttribPointer(6, 2, r5.FLOAT, false, Ve2, 7 * Float32Array.BYTES_PER_ELEMENT), r5.vertexAttribDivisor(6, 1), r5.enableVertexAttribArray(1), r5.vertexAttribPointer(1, 2, r5.FLOAT, false, Ve2, 9 * Float32Array.BYTES_PER_ELEMENT), r5.vertexAttribDivisor(1, 1), r5.useProgram(this._program);
    let d = new Int32Array(ae.maxAtlasPages);
    for (let h2 = 0; h2 < ae.maxAtlasPages; h2++) d[h2] = h2;
    r5.uniform1iv(this._textureLocation, d), r5.uniformMatrix4fv(this._projectionLocation, false, Gt), this._atlasTextures = [];
    for (let h2 = 0; h2 < ae.maxAtlasPages; h2++) {
      let f = new Wt(F(r5.createTexture()));
      this._register(O2(() => r5.deleteTexture(f.texture))), r5.activeTexture(r5.TEXTURE0 + h2), r5.bindTexture(r5.TEXTURE_2D, f.texture), r5.texParameteri(r5.TEXTURE_2D, r5.TEXTURE_WRAP_S, r5.CLAMP_TO_EDGE), r5.texParameteri(r5.TEXTURE_2D, r5.TEXTURE_WRAP_T, r5.CLAMP_TO_EDGE), r5.texImage2D(r5.TEXTURE_2D, 0, r5.RGBA, 1, 1, 0, r5.RGBA, r5.UNSIGNED_BYTE, new Uint8Array([255, 0, 0, 255])), this._atlasTextures[h2] = f;
    }
    r5.enable(r5.BLEND), r5.blendFunc(r5.SRC_ALPHA, r5.ONE_MINUS_SRC_ALPHA), this.handleResize();
  }
  beginFrame() {
    return this._atlas ? this._atlas.beginFrame() : true;
  }
  updateCell(t, n, s15, o2, r5, a, l, u, c) {
    this._updateCell(this._vertices.attributes, t, n, s15, o2, r5, a, l, u, c);
  }
  _updateCell(t, n, s15, o2, r5, a, l, u, c, d) {
    if (H2 = (s15 * this._terminal.cols + n) * De2, o2 === 0 || o2 === void 0) {
      t.fill(0, H2, H2 + De2 - 1 - rs);
      return;
    }
    this._atlas && (u && u.length > 1 ? k = this._atlas.getRasterizedGlyphCombinedChar(u, r5, a, l, false, this._terminal.element) : k = this._atlas.getRasterizedGlyph(o2, r5, a, l, false, this._terminal.element), Fi = Math.floor((this._dimensions.device.cell.width - this._dimensions.device.char.width) / 2), r5 !== d && k.offset.x > Fi ? (lt2 = k.offset.x - Fi, t[H2] = -(k.offset.x - lt2) + this._dimensions.device.char.left, t[H2 + 1] = -k.offset.y + this._dimensions.device.char.top, t[H2 + 2] = (k.size.x - lt2) / this._dimensions.device.canvas.width, t[H2 + 3] = k.size.y / this._dimensions.device.canvas.height, t[H2 + 4] = k.texturePage, t[H2 + 5] = k.texturePositionClipSpace.x + lt2 / this._atlas.pages[k.texturePage].canvas.width, t[H2 + 6] = k.texturePositionClipSpace.y, t[H2 + 7] = k.sizeClipSpace.x - lt2 / this._atlas.pages[k.texturePage].canvas.width, t[H2 + 8] = k.sizeClipSpace.y) : (t[H2] = -k.offset.x + this._dimensions.device.char.left, t[H2 + 1] = -k.offset.y + this._dimensions.device.char.top, t[H2 + 2] = k.size.x / this._dimensions.device.canvas.width, t[H2 + 3] = k.size.y / this._dimensions.device.canvas.height, t[H2 + 4] = k.texturePage, t[H2 + 5] = k.texturePositionClipSpace.x, t[H2 + 6] = k.texturePositionClipSpace.y, t[H2 + 7] = k.sizeClipSpace.x, t[H2 + 8] = k.sizeClipSpace.y), this._optionsService.rawOptions.rescaleOverlappingGlyphs && mn(o2, c, k.size.x, this._dimensions.device.cell.width) && (t[H2 + 2] = (this._dimensions.device.cell.width - 1) / this._dimensions.device.canvas.width));
  }
  clear() {
    let t = this._terminal, n = t.cols * t.rows * De2;
    this._vertices.count !== n ? this._vertices.attributes = new Float32Array(n) : this._vertices.attributes.fill(0);
    let s15 = 0;
    for (; s15 < this._vertices.attributesBuffers.length; s15++) this._vertices.count !== n ? this._vertices.attributesBuffers[s15] = new Float32Array(n) : this._vertices.attributesBuffers[s15].fill(0);
    this._vertices.count = n, s15 = 0;
    for (let o2 = 0; o2 < t.rows; o2++) for (let r5 = 0; r5 < t.cols; r5++) this._vertices.attributes[s15 + 9] = r5 / t.cols, this._vertices.attributes[s15 + 10] = o2 / t.rows, s15 += De2;
  }
  handleResize() {
    let t = this._gl;
    t.useProgram(this._program), t.viewport(0, 0, t.canvas.width, t.canvas.height), t.uniform2f(this._resolutionLocation, t.canvas.width, t.canvas.height), this.clear();
  }
  render(t) {
    if (!this._atlas) return;
    let n = this._gl;
    n.useProgram(this._program), n.bindVertexArray(this._vertexArrayObject), this._activeBuffer = (this._activeBuffer + 1) % 2;
    let s15 = this._vertices.attributesBuffers[this._activeBuffer], o2 = 0;
    for (let r5 = 0; r5 < t.lineLengths.length; r5++) {
      let a = r5 * this._terminal.cols * De2, l = this._vertices.attributes.subarray(a, a + t.lineLengths[r5] * De2);
      s15.set(l, o2), o2 += l.length;
    }
    n.bindBuffer(n.ARRAY_BUFFER, this._attributesBuffer), n.bufferData(n.ARRAY_BUFFER, s15.subarray(0, o2), n.STREAM_DRAW);
    for (let r5 = 0; r5 < this._atlas.pages.length; r5++) this._atlas.pages[r5].version !== this._atlasTextures[r5].version && this._bindAtlasPageTexture(n, this._atlas, r5);
    n.drawElementsInstanced(n.TRIANGLE_STRIP, 4, n.UNSIGNED_BYTE, 0, o2 / De2);
  }
  setAtlas(t) {
    this._atlas = t;
    for (let n of this._atlasTextures) n.version = -1;
  }
  _bindAtlasPageTexture(t, n, s15) {
    t.activeTexture(t.TEXTURE0 + s15), t.bindTexture(t.TEXTURE_2D, this._atlasTextures[s15].texture), t.texParameteri(t.TEXTURE_2D, t.TEXTURE_WRAP_S, t.CLAMP_TO_EDGE), t.texParameteri(t.TEXTURE_2D, t.TEXTURE_WRAP_T, t.CLAMP_TO_EDGE), t.texImage2D(t.TEXTURE_2D, 0, t.RGBA, t.RGBA, t.UNSIGNED_BYTE, n.pages[s15].canvas), t.generateMipmap(t.TEXTURE_2D), this._atlasTextures[s15].version = n.pages[s15].version;
  }
  setDimensions(t) {
    this._dimensions = t;
  }
};
var ki = class {
  constructor() {
    this.clear();
  }
  clear() {
    this.hasSelection = false, this.columnSelectMode = false, this.viewportStartRow = 0, this.viewportEndRow = 0, this.viewportCappedStartRow = 0, this.viewportCappedEndRow = 0, this.startCol = 0, this.endCol = 0, this.selectionStart = void 0, this.selectionEnd = void 0;
  }
  update(e, t, n, s15 = false) {
    if (this.selectionStart = t, this.selectionEnd = n, !t || !n || t[0] === n[0] && t[1] === n[1]) {
      this.clear();
      return;
    }
    let o2 = e.buffers.active.ydisp, r5 = t[1] - o2, a = n[1] - o2, l = Math.max(r5, 0), u = Math.min(a, e.rows - 1);
    if (l >= e.rows || u < 0) {
      this.clear();
      return;
    }
    this.hasSelection = true, this.columnSelectMode = s15, this.viewportStartRow = r5, this.viewportEndRow = a, this.viewportCappedStartRow = l, this.viewportCappedEndRow = u, this.startCol = t[0], this.endCol = n[0];
  }
  isCellSelected(e, t, n) {
    return this.hasSelection ? (n -= e.buffer.active.viewportY, this.columnSelectMode ? this.startCol <= this.endCol ? t >= this.startCol && n >= this.viewportCappedStartRow && t < this.endCol && n <= this.viewportCappedEndRow : t < this.startCol && n >= this.viewportCappedStartRow && t >= this.endCol && n <= this.viewportCappedEndRow : n > this.viewportStartRow && n < this.viewportEndRow || this.viewportStartRow === this.viewportEndRow && n === this.viewportStartRow && t >= this.startCol && t < this.endCol || this.viewportStartRow < this.viewportEndRow && n === this.viewportEndRow && t < this.endCol || this.viewportStartRow < this.viewportEndRow && n === this.viewportStartRow && t >= this.startCol) : false;
  }
};
function Nn() {
  return new ki();
}
var Ce = 4;
var ze = 1;
var qe2 = 2;
var Ct = 3;
var Un = 2147483648;
var Vt = class {
  constructor() {
    this.cells = new Uint32Array(0), this.lineLengths = new Uint32Array(0), this.selection = Nn();
  }
  resize(e, t) {
    let n = e * t * Ce;
    n !== this.cells.length && (this.cells = new Uint32Array(n), this.lineLengths = new Uint32Array(t));
  }
  clear() {
    this.cells.fill(0, 0), this.lineLengths.fill(0, 0);
  }
};
var ss = `#version 300 es
layout (location = 0) in vec2 a_position;
layout (location = 1) in vec2 a_size;
layout (location = 2) in vec4 a_color;
layout (location = 3) in vec2 a_unitquad;

uniform mat4 u_projection;

out vec4 v_color;

void main() {
  vec2 zeroToOne = a_position + (a_unitquad * a_size);
  gl_Position = u_projection * vec4(zeroToOne, 0.0, 1.0);
  v_color = a_color;
}`;
var os = `#version 300 es
precision lowp float;

in vec4 v_color;

out vec4 outColor;

void main() {
  outColor = v_color;
}`;
var Ee2 = 8;
var Pi = Ee2 * Float32Array.BYTES_PER_ELEMENT;
var as = 20 * Ee2;
var zt = class {
  constructor() {
    this.attributes = new Float32Array(as), this.count = 0;
  }
};
var xe = 0;
var Hn = 0;
var Wn = 0;
var Gn = 0;
var $n = 0;
var Kn = 0;
var Vn = 0;
var qt = class extends B {
  constructor(t, n, s15, o2) {
    super();
    this._terminal = t;
    this._gl = n;
    this._dimensions = s15;
    this._themeService = o2;
    this._vertices = new zt();
    this._verticesCursor = new zt();
    let r5 = this._gl;
    this._program = F($t(r5, ss, os)), this._register(O2(() => r5.deleteProgram(this._program))), this._projectionLocation = F(r5.getUniformLocation(this._program, "u_projection")), this._vertexArrayObject = r5.createVertexArray(), r5.bindVertexArray(this._vertexArrayObject);
    let a = new Float32Array([0, 0, 1, 0, 0, 1, 1, 1]), l = r5.createBuffer();
    this._register(O2(() => r5.deleteBuffer(l))), r5.bindBuffer(r5.ARRAY_BUFFER, l), r5.bufferData(r5.ARRAY_BUFFER, a, r5.STATIC_DRAW), r5.enableVertexAttribArray(3), r5.vertexAttribPointer(3, 2, this._gl.FLOAT, false, 0, 0);
    let u = new Uint8Array([0, 1, 2, 3]), c = r5.createBuffer();
    this._register(O2(() => r5.deleteBuffer(c))), r5.bindBuffer(r5.ELEMENT_ARRAY_BUFFER, c), r5.bufferData(r5.ELEMENT_ARRAY_BUFFER, u, r5.STATIC_DRAW), this._attributesBuffer = F(r5.createBuffer()), this._register(O2(() => r5.deleteBuffer(this._attributesBuffer))), r5.bindBuffer(r5.ARRAY_BUFFER, this._attributesBuffer), r5.enableVertexAttribArray(0), r5.vertexAttribPointer(0, 2, r5.FLOAT, false, Pi, 0), r5.vertexAttribDivisor(0, 1), r5.enableVertexAttribArray(1), r5.vertexAttribPointer(1, 2, r5.FLOAT, false, Pi, 2 * Float32Array.BYTES_PER_ELEMENT), r5.vertexAttribDivisor(1, 1), r5.enableVertexAttribArray(2), r5.vertexAttribPointer(2, 4, r5.FLOAT, false, Pi, 4 * Float32Array.BYTES_PER_ELEMENT), r5.vertexAttribDivisor(2, 1), this._updateCachedColors(o2.colors), this._register(this._themeService.onChangeColors((d) => {
      this._updateCachedColors(d), this._updateViewportRectangle();
    }));
  }
  renderBackgrounds() {
    this._renderVertices(this._vertices);
  }
  renderCursor() {
    this._renderVertices(this._verticesCursor);
  }
  _renderVertices(t) {
    let n = this._gl;
    n.useProgram(this._program), n.bindVertexArray(this._vertexArrayObject), n.uniformMatrix4fv(this._projectionLocation, false, Gt), n.bindBuffer(n.ARRAY_BUFFER, this._attributesBuffer), n.bufferData(n.ARRAY_BUFFER, t.attributes, n.DYNAMIC_DRAW), n.drawElementsInstanced(this._gl.TRIANGLE_STRIP, 4, n.UNSIGNED_BYTE, 0, t.count);
  }
  handleResize() {
    this._updateViewportRectangle();
  }
  setDimensions(t) {
    this._dimensions = t;
  }
  _updateCachedColors(t) {
    this._bgFloat = this._colorToFloat32Array(t.background), this._cursorFloat = this._colorToFloat32Array(t.cursor);
  }
  _updateViewportRectangle() {
    this._addRectangleFloat(this._vertices.attributes, 0, 0, 0, this._terminal.cols * this._dimensions.device.cell.width, this._terminal.rows * this._dimensions.device.cell.height, this._bgFloat);
  }
  updateBackgrounds(t) {
    let n = this._terminal, s15 = this._vertices, o2 = 1, r5, a, l, u, c, d, h2, f, I2, L2, M4;
    for (r5 = 0; r5 < n.rows; r5++) {
      for (l = -1, u = 0, c = 0, d = false, a = 0; a < n.cols; a++) h2 = (r5 * n.cols + a) * Ce, f = t.cells[h2 + ze], I2 = t.cells[h2 + qe2], L2 = !!(I2 & 67108864), (f !== u || I2 !== c && (d || L2)) && ((u !== 0 || d && c !== 0) && (M4 = o2++ * Ee2, this._updateRectangle(s15, M4, c, u, l, a, r5)), l = a, u = f, c = I2, d = L2);
      (u !== 0 || d && c !== 0) && (M4 = o2++ * Ee2, this._updateRectangle(s15, M4, c, u, l, n.cols, r5));
    }
    s15.count = o2;
  }
  updateCursor(t) {
    let n = this._verticesCursor, s15 = t.cursor;
    if (!s15 || s15.style === "block") {
      n.count = 0;
      return;
    }
    let o2, r5 = 0;
    (s15.style === "bar" || s15.style === "outline") && (o2 = r5++ * Ee2, this._addRectangleFloat(n.attributes, o2, s15.x * this._dimensions.device.cell.width, s15.y * this._dimensions.device.cell.height, s15.style === "bar" ? s15.dpr * s15.cursorWidth : s15.dpr, this._dimensions.device.cell.height, this._cursorFloat)), (s15.style === "underline" || s15.style === "outline") && (o2 = r5++ * Ee2, this._addRectangleFloat(n.attributes, o2, s15.x * this._dimensions.device.cell.width, (s15.y + 1) * this._dimensions.device.cell.height - s15.dpr, s15.width * this._dimensions.device.cell.width, s15.dpr, this._cursorFloat)), s15.style === "outline" && (o2 = r5++ * Ee2, this._addRectangleFloat(n.attributes, o2, s15.x * this._dimensions.device.cell.width, s15.y * this._dimensions.device.cell.height, s15.width * this._dimensions.device.cell.width, s15.dpr, this._cursorFloat), o2 = r5++ * Ee2, this._addRectangleFloat(n.attributes, o2, (s15.x + s15.width) * this._dimensions.device.cell.width - s15.dpr, s15.y * this._dimensions.device.cell.height, s15.dpr, this._dimensions.device.cell.height, this._cursorFloat)), n.count = r5;
  }
  _updateRectangle(t, n, s15, o2, r5, a, l) {
    if (s15 & 67108864) switch (s15 & 50331648) {
      case 16777216:
      case 33554432:
        xe = this._themeService.colors.ansi[s15 & 255].rgba;
        break;
      case 50331648:
        xe = (s15 & 16777215) << 8;
        break;
      case 0:
      default:
        xe = this._themeService.colors.foreground.rgba;
    }
    else switch (o2 & 50331648) {
      case 16777216:
      case 33554432:
        xe = this._themeService.colors.ansi[o2 & 255].rgba;
        break;
      case 50331648:
        xe = (o2 & 16777215) << 8;
        break;
      case 0:
      default:
        xe = this._themeService.colors.background.rgba;
    }
    t.attributes.length < n + 4 && (t.attributes = Bn(t.attributes, this._terminal.rows * this._terminal.cols * Ee2)), Hn = r5 * this._dimensions.device.cell.width, Wn = l * this._dimensions.device.cell.height, Gn = (xe >> 24 & 255) / 255, $n = (xe >> 16 & 255) / 255, Kn = (xe >> 8 & 255) / 255, Vn = 1, this._addRectangle(t.attributes, n, Hn, Wn, (a - r5) * this._dimensions.device.cell.width, this._dimensions.device.cell.height, Gn, $n, Kn, Vn);
  }
  _addRectangle(t, n, s15, o2, r5, a, l, u, c, d) {
    t[n] = s15 / this._dimensions.device.canvas.width, t[n + 1] = o2 / this._dimensions.device.canvas.height, t[n + 2] = r5 / this._dimensions.device.canvas.width, t[n + 3] = a / this._dimensions.device.canvas.height, t[n + 4] = l, t[n + 5] = u, t[n + 6] = c, t[n + 7] = d;
  }
  _addRectangleFloat(t, n, s15, o2, r5, a, l) {
    t[n] = s15 / this._dimensions.device.canvas.width, t[n + 1] = o2 / this._dimensions.device.canvas.height, t[n + 2] = r5 / this._dimensions.device.canvas.width, t[n + 3] = a / this._dimensions.device.canvas.height, t[n + 4] = l[0], t[n + 5] = l[1], t[n + 6] = l[2], t[n + 7] = l[3];
  }
  _colorToFloat32Array(t) {
    return new Float32Array([(t.rgba >> 24 & 255) / 255, (t.rgba >> 16 & 255) / 255, (t.rgba >> 8 & 255) / 255, (t.rgba & 255) / 255]);
  }
};
var jt = class extends B {
  constructor(t, n, s15, o2, r5, a, l, u) {
    super();
    this._container = n;
    this._alpha = r5;
    this._coreBrowserService = a;
    this._optionsService = l;
    this._themeService = u;
    this._deviceCharWidth = 0;
    this._deviceCharHeight = 0;
    this._deviceCellWidth = 0;
    this._deviceCellHeight = 0;
    this._deviceCharLeft = 0;
    this._deviceCharTop = 0;
    this._canvas = this._coreBrowserService.mainDocument.createElement("canvas"), this._canvas.classList.add(`xterm-${s15}-layer`), this._canvas.style.zIndex = o2.toString(), this._initCanvas(), this._container.appendChild(this._canvas), this._register(this._themeService.onChangeColors((c) => {
      this._refreshCharAtlas(t, c), this.reset(t);
    })), this._register(O2(() => {
      this._canvas.remove();
    }));
  }
  _initCanvas() {
    this._ctx = F(this._canvas.getContext("2d", { alpha: this._alpha })), this._alpha || this._clearAll();
  }
  handleBlur(t) {
  }
  handleFocus(t) {
  }
  handleCursorMove(t) {
  }
  handleGridChanged(t, n, s15) {
  }
  handleSelectionChanged(t, n, s15, o2 = false) {
  }
  _setTransparency(t, n) {
    if (n === this._alpha) return;
    let s15 = this._canvas;
    this._alpha = n, this._canvas = this._canvas.cloneNode(), this._initCanvas(), this._container.replaceChild(this._canvas, s15), this._refreshCharAtlas(t, this._themeService.colors), this.handleGridChanged(t, 0, t.rows - 1);
  }
  _refreshCharAtlas(t, n) {
    this._deviceCharWidth <= 0 && this._deviceCharHeight <= 0 || (this._charAtlas = Nt(t, this._optionsService.rawOptions, n, this._deviceCellWidth, this._deviceCellHeight, this._deviceCharWidth, this._deviceCharHeight, this._coreBrowserService.dpr, 2048), this._charAtlas.warmUp());
  }
  resize(t, n) {
    this._deviceCellWidth = n.device.cell.width, this._deviceCellHeight = n.device.cell.height, this._deviceCharWidth = n.device.char.width, this._deviceCharHeight = n.device.char.height, this._deviceCharLeft = n.device.char.left, this._deviceCharTop = n.device.char.top, this._canvas.width = n.device.canvas.width, this._canvas.height = n.device.canvas.height, this._canvas.style.width = `${n.css.canvas.width}px`, this._canvas.style.height = `${n.css.canvas.height}px`, this._alpha || this._clearAll(), this._refreshCharAtlas(t, this._themeService.colors);
  }
  _fillBottomLineAtCells(t, n, s15 = 1) {
    this._ctx.fillRect(t * this._deviceCellWidth, (n + 1) * this._deviceCellHeight - this._coreBrowserService.dpr - 1, s15 * this._deviceCellWidth, this._coreBrowserService.dpr);
  }
  _clearAll() {
    this._alpha ? this._ctx.clearRect(0, 0, this._canvas.width, this._canvas.height) : (this._ctx.fillStyle = this._themeService.colors.background.css, this._ctx.fillRect(0, 0, this._canvas.width, this._canvas.height));
  }
  _clearCells(t, n, s15, o2) {
    this._alpha ? this._ctx.clearRect(t * this._deviceCellWidth, n * this._deviceCellHeight, s15 * this._deviceCellWidth, o2 * this._deviceCellHeight) : (this._ctx.fillStyle = this._themeService.colors.background.css, this._ctx.fillRect(t * this._deviceCellWidth, n * this._deviceCellHeight, s15 * this._deviceCellWidth, o2 * this._deviceCellHeight));
  }
  _fillCharTrueColor(t, n, s15, o2) {
    this._ctx.font = this._getFont(t, false, false), this._ctx.textBaseline = St, this._clipCell(s15, o2, n.getWidth()), this._ctx.fillText(n.getChars(), s15 * this._deviceCellWidth + this._deviceCharLeft, o2 * this._deviceCellHeight + this._deviceCharTop + this._deviceCharHeight);
  }
  _clipCell(t, n, s15) {
    this._ctx.beginPath(), this._ctx.rect(t * this._deviceCellWidth, n * this._deviceCellHeight, s15 * this._deviceCellWidth, this._deviceCellHeight), this._ctx.clip();
  }
  _getFont(t, n, s15) {
    let o2 = n ? t.options.fontWeightBold : t.options.fontWeight;
    return `${s15 ? "italic" : ""} ${o2} ${t.options.fontSize * this._coreBrowserService.dpr}px ${t.options.fontFamily}`;
  }
};
var Xt = class extends jt {
  constructor(e, t, n, s15, o2, r5, a) {
    super(n, e, "link", t, true, o2, r5, a), this._register(s15.onShowLinkUnderline((l) => this._handleShowLinkUnderline(l))), this._register(s15.onHideLinkUnderline((l) => this._handleHideLinkUnderline(l)));
  }
  resize(e, t) {
    super.resize(e, t), this._state = void 0;
  }
  reset(e) {
    this._clearCurrentLink();
  }
  _clearCurrentLink() {
    if (this._state) {
      this._clearCells(this._state.x1, this._state.y1, this._state.cols - this._state.x1, 1);
      let e = this._state.y2 - this._state.y1 - 1;
      e > 0 && this._clearCells(0, this._state.y1 + 1, this._state.cols, e), this._clearCells(0, this._state.y2, this._state.x2, 1), this._state = void 0;
    }
  }
  _handleShowLinkUnderline(e) {
    if (e.fg === 257 ? this._ctx.fillStyle = this._themeService.colors.background.css : e.fg !== void 0 && Fn(e.fg) ? this._ctx.fillStyle = this._themeService.colors.ansi[e.fg].css : this._ctx.fillStyle = this._themeService.colors.foreground.css, e.y1 === e.y2) this._fillBottomLineAtCells(e.x1, e.y1, e.x2 - e.x1);
    else {
      this._fillBottomLineAtCells(e.x1, e.y1, e.cols - e.x1);
      for (let t = e.y1 + 1; t < e.y2; t++) this._fillBottomLineAtCells(0, t, e.cols);
      this._fillBottomLineAtCells(0, e.y2, e.x2);
    }
    this._state = e;
  }
  _handleHideLinkUnderline(e) {
    this._clearCurrentLink();
  }
};
var te3 = typeof window == "object" ? window : globalThis;
var Zt = class Zt2 {
  constructor() {
    this.mapWindowIdToZoomLevel = /* @__PURE__ */ new Map();
    this._onDidChangeZoomLevel = new D();
    this.onDidChangeZoomLevel = this._onDidChangeZoomLevel.event;
    this.mapWindowIdToZoomFactor = /* @__PURE__ */ new Map();
    this._onDidChangeFullscreen = new D();
    this.onDidChangeFullscreen = this._onDidChangeFullscreen.event;
    this.mapWindowIdToFullScreen = /* @__PURE__ */ new Map();
  }
  getZoomLevel(e) {
    return this.mapWindowIdToZoomLevel.get(this.getWindowId(e)) ?? 0;
  }
  setZoomLevel(e, t) {
    if (this.getZoomLevel(t) === e) return;
    let n = this.getWindowId(t);
    this.mapWindowIdToZoomLevel.set(n, e), this._onDidChangeZoomLevel.fire(n);
  }
  getZoomFactor(e) {
    return this.mapWindowIdToZoomFactor.get(this.getWindowId(e)) ?? 1;
  }
  setZoomFactor(e, t) {
    this.mapWindowIdToZoomFactor.set(this.getWindowId(t), e);
  }
  setFullscreen(e, t) {
    if (this.isFullscreen(t) === e) return;
    let n = this.getWindowId(t);
    this.mapWindowIdToFullScreen.set(n, e), this._onDidChangeFullscreen.fire(n);
  }
  isFullscreen(e) {
    return !!this.mapWindowIdToFullScreen.get(this.getWindowId(e));
  }
  getWindowId(e) {
    return e.vscodeWindowId;
  }
};
Zt.INSTANCE = new Zt();
var Qt = Zt;
function us(i8, e, t) {
  typeof e == "string" && (e = i8.matchMedia(e)), e.addEventListener("change", t);
}
var Wa = Qt.INSTANCE.onDidChangeZoomLevel;
var Ga = Qt.INSTANCE.onDidChangeFullscreen;
var je2 = typeof navigator == "object" ? navigator.userAgent : "";
var Cn = je2.indexOf("Firefox") >= 0;
var ut = je2.indexOf("AppleWebKit") >= 0;
var zn = je2.indexOf("Chrome") >= 0;
var Bi = !zn && je2.indexOf("Safari") >= 0;
var $a = je2.indexOf("Electron/") >= 0;
var Ka = je2.indexOf("Android") >= 0;
var Yt = false;
if (typeof te3.matchMedia == "function") {
  let i8 = te3.matchMedia("(display-mode: standalone) or (display-mode: window-controls-overlay)"), e = te3.matchMedia("(display-mode: fullscreen)");
  Yt = i8.matches, us(te3, i8, ({ matches: t }) => {
    Yt && e.matches || (Yt = t);
  });
}
function qn() {
  return Yt;
}
var Xe2 = "en";
var Ui = false;
var ni = false;
var ti = false;
var cs = false;
var Xn = false;
var Yn = false;
var ds = false;
var hs = false;
var ps = false;
var fs = false;
var ei;
var ii = Xe2;
var jn = Xe2;
var ms;
var ye2;
var Ie = globalThis;
var re2;
typeof Ie.vscode < "u" && typeof Ie.vscode.process < "u" ? re2 = Ie.vscode.process : typeof process < "u" && typeof process?.versions?.node == "string" && (re2 = process);
var Qn = typeof re2?.versions?.electron == "string";
var _s = Qn && re2?.type === "renderer";
if (typeof re2 == "object") {
  Ui = re2.platform === "win32", ni = re2.platform === "darwin", ti = re2.platform === "linux", cs = ti && !!re2.env.SNAP && !!re2.env.SNAP_REVISION, ds = Qn, ps = !!re2.env.CI || !!re2.env.BUILD_ARTIFACTSTAGINGDIRECTORY, ei = Xe2, ii = Xe2;
  let i8 = re2.env.VSCODE_NLS_CONFIG;
  if (i8) try {
    let e = JSON.parse(i8);
    ei = e.userLocale, jn = e.osLocale, ii = e.resolvedLanguage || Xe2, ms = e.languagePack?.translationsConfigFile;
  } catch {
  }
  Xn = true;
} else typeof navigator == "object" && !_s ? (ye2 = navigator.userAgent, Ui = ye2.indexOf("Windows") >= 0, ni = ye2.indexOf("Macintosh") >= 0, hs = (ye2.indexOf("Macintosh") >= 0 || ye2.indexOf("iPad") >= 0 || ye2.indexOf("iPhone") >= 0) && !!navigator.maxTouchPoints && navigator.maxTouchPoints > 0, ti = ye2.indexOf("Linux") >= 0, fs = ye2?.indexOf("Mobi") >= 0, Yn = true, ii = globalThis._VSCODE_NLS_LANGUAGE || Xe2, ei = navigator.language.toLowerCase(), jn = ei) : console.error("Unable to resolve platform.");
var Ni = 0;
ni ? Ni = 1 : Ui ? Ni = 3 : ti && (Ni = 2);
var ri = Xn;
var bs = Yn && typeof Ie.importScripts == "function";
var Va = bs ? Ie.origin : void 0;
var _e = ye2;
var Me2 = ii;
var vs;
((n) => {
  function i8() {
    return Me2;
  }
  n.value = i8;
  function e() {
    return Me2.length === 2 ? Me2 === "en" : Me2.length >= 3 ? Me2[0] === "e" && Me2[1] === "n" && Me2[2] === "-" : false;
  }
  n.isDefaultVariant = e;
  function t() {
    return Me2 === "en";
  }
  n.isDefault = t;
})(vs ||= {});
var Ts = typeof Ie.postMessage == "function" && !Ie.importScripts;
var Zn = (() => {
  if (Ts) {
    let i8 = [];
    Ie.addEventListener("message", (t) => {
      if (t.data && t.data.vscodeScheduleAsyncWork) for (let n = 0, s15 = i8.length; n < s15; n++) {
        let o2 = i8[n];
        if (o2.id === t.data.vscodeScheduleAsyncWork) {
          i8.splice(n, 1), o2.callback();
          return;
        }
      }
    });
    let e = 0;
    return (t) => {
      let n = ++e;
      i8.push({ id: n, callback: t }), Ie.postMessage({ vscodeScheduleAsyncWork: n }, "*");
    };
  }
  return (i8) => setTimeout(i8);
})();
var gs = !!(_e && _e.indexOf("Chrome") >= 0);
var Ca = !!(_e && _e.indexOf("Firefox") >= 0);
var za = !!(!gs && _e && _e.indexOf("Safari") >= 0);
var qa = !!(_e && _e.indexOf("Edg/") >= 0);
var ja = !!(_e && _e.indexOf("Android") >= 0);
var Ae2 = typeof navigator == "object" ? navigator : {};
var xs = { clipboard: { writeText: ri || document.queryCommandSupported && document.queryCommandSupported("copy") || !!(Ae2 && Ae2.clipboard && Ae2.clipboard.writeText), readText: ri || !!(Ae2 && Ae2.clipboard && Ae2.clipboard.readText) }, keyboard: ri || qn() ? 0 : Ae2.keyboard || Bi ? 1 : 2, touch: "ontouchstart" in te3 || Ae2.maxTouchPoints > 0, pointerEvents: te3.PointerEvent && ("ontouchstart" in te3 || navigator.maxTouchPoints > 0) };
var dt = class {
  constructor() {
    this._keyCodeToStr = [], this._strToKeyCode = /* @__PURE__ */ Object.create(null);
  }
  define(e, t) {
    this._keyCodeToStr[e] = t, this._strToKeyCode[t.toLowerCase()] = e;
  }
  keyCodeToStr(e) {
    return this._keyCodeToStr[e];
  }
  strToKeyCode(e) {
    return this._strToKeyCode[e.toLowerCase()] || 0;
  }
};
var Hi = new dt();
var Jn = new dt();
var er = new dt();
var Es = new Array(230);
var tr;
((r5) => {
  function i8(a) {
    return Hi.keyCodeToStr(a);
  }
  r5.toString = i8;
  function e(a) {
    return Hi.strToKeyCode(a);
  }
  r5.fromString = e;
  function t(a) {
    return Jn.keyCodeToStr(a);
  }
  r5.toUserSettingsUS = t;
  function n(a) {
    return er.keyCodeToStr(a);
  }
  r5.toUserSettingsGeneral = n;
  function s15(a) {
    return Jn.strToKeyCode(a) || er.strToKeyCode(a);
  }
  r5.fromUserSettings = s15;
  function o2(a) {
    if (a >= 98 && a <= 113) return null;
    switch (a) {
      case 16:
        return "Up";
      case 18:
        return "Down";
      case 15:
        return "Left";
      case 17:
        return "Right";
    }
    return Hi.keyCodeToStr(a);
  }
  r5.toElectronAccelerator = o2;
})(tr ||= {});
var nr = Object.freeze(function(i8, e) {
  let t = setTimeout(i8.bind(e), 0);
  return { dispose() {
    clearTimeout(t);
  } };
});
var Is;
((n) => {
  function i8(s15) {
    return s15 === n.None || s15 === n.Cancelled || s15 instanceof Wi ? true : !s15 || typeof s15 != "object" ? false : typeof s15.isCancellationRequested == "boolean" && typeof s15.onCancellationRequested == "function";
  }
  n.isCancellationToken = i8, n.None = Object.freeze({ isCancellationRequested: false, onCancellationRequested: ee3.None }), n.Cancelled = Object.freeze({ isCancellationRequested: true, onCancellationRequested: nr });
})(Is ||= {});
var Wi = class {
  constructor() {
    this._isCancelled = false;
    this._emitter = null;
  }
  cancel() {
    this._isCancelled || (this._isCancelled = true, this._emitter && (this._emitter.fire(void 0), this.dispose()));
  }
  get isCancellationRequested() {
    return this._isCancelled;
  }
  get onCancellationRequested() {
    return this._isCancelled ? nr : (this._emitter || (this._emitter = new D()), this._emitter.event);
  }
  dispose() {
    this._emitter && (this._emitter.dispose(), this._emitter = null);
  }
};
var ws;
var oi;
(function() {
  typeof globalThis.requestIdleCallback != "function" || typeof globalThis.cancelIdleCallback != "function" ? oi = (i8, e) => {
    Zn(() => {
      if (t) return;
      let n = Date.now() + 15;
      e(Object.freeze({ didTimeout: true, timeRemaining() {
        return Math.max(0, n - Date.now());
      } }));
    });
    let t = false;
    return { dispose() {
      t || (t = true);
    } };
  } : oi = (i8, e, t) => {
    let n = i8.requestIdleCallback(e, typeof t == "number" ? { timeout: t } : void 0), s15 = false;
    return { dispose() {
      s15 || (s15 = true, i8.cancelIdleCallback(n));
    } };
  }, ws = (i8) => oi(globalThis, i8);
})();
var Rs;
((t) => {
  async function i8(n) {
    let s15, o2 = await Promise.all(n.map((r5) => r5.then((a) => a, (a) => {
      s15 || (s15 = a);
    })));
    if (typeof s15 < "u") throw s15;
    return o2;
  }
  t.settled = i8;
  function e(n) {
    return new Promise(async (s15, o2) => {
      try {
        await n(s15, o2);
      } catch (r5) {
        o2(r5);
      }
    });
  }
  t.withAsyncBody = e;
})(Rs ||= {});
var Q = class Q2 {
  static fromArray(e) {
    return new Q2((t) => {
      t.emitMany(e);
    });
  }
  static fromPromise(e) {
    return new Q2(async (t) => {
      t.emitMany(await e);
    });
  }
  static fromPromises(e) {
    return new Q2(async (t) => {
      await Promise.all(e.map(async (n) => t.emitOne(await n)));
    });
  }
  static merge(e) {
    return new Q2(async (t) => {
      await Promise.all(e.map(async (n) => {
        for await (let s15 of n) t.emitOne(s15);
      }));
    });
  }
  constructor(e, t) {
    this._state = 0, this._results = [], this._error = null, this._onReturn = t, this._onStateChanged = new D(), queueMicrotask(async () => {
      let n = { emitOne: (s15) => this.emitOne(s15), emitMany: (s15) => this.emitMany(s15), reject: (s15) => this.reject(s15) };
      try {
        await Promise.resolve(e(n)), this.resolve();
      } catch (s15) {
        this.reject(s15);
      } finally {
        n.emitOne = void 0, n.emitMany = void 0, n.reject = void 0;
      }
    });
  }
  [Symbol.asyncIterator]() {
    let e = 0;
    return { next: async () => {
      do {
        if (this._state === 2) throw this._error;
        if (e < this._results.length) return { done: false, value: this._results[e++] };
        if (this._state === 1) return { done: true, value: void 0 };
        await ee3.toPromise(this._onStateChanged.event);
      } while (true);
    }, return: async () => (this._onReturn?.(), { done: true, value: void 0 }) };
  }
  static map(e, t) {
    return new Q2(async (n) => {
      for await (let s15 of e) n.emitOne(t(s15));
    });
  }
  map(e) {
    return Q2.map(this, e);
  }
  static filter(e, t) {
    return new Q2(async (n) => {
      for await (let s15 of e) t(s15) && n.emitOne(s15);
    });
  }
  filter(e) {
    return Q2.filter(this, e);
  }
  static coalesce(e) {
    return Q2.filter(e, (t) => !!t);
  }
  coalesce() {
    return Q2.coalesce(this);
  }
  static async toPromise(e) {
    let t = [];
    for await (let n of e) t.push(n);
    return t;
  }
  toPromise() {
    return Q2.toPromise(this);
  }
  emitOne(e) {
    this._state === 0 && (this._results.push(e), this._onStateChanged.fire());
  }
  emitMany(e) {
    this._state === 0 && (this._results = this._results.concat(e), this._onStateChanged.fire());
  }
  resolve() {
    this._state === 0 && (this._state = 1, this._onStateChanged.fire());
  }
  reject(e) {
    this._state === 0 && (this._state = 2, this._error = e, this._onStateChanged.fire());
  }
};
Q.EMPTY = Q.fromArray([]);
function sr(i8) {
  return 55296 <= i8 && i8 <= 56319;
}
function Gi(i8) {
  return 56320 <= i8 && i8 <= 57343;
}
function or(i8, e) {
  return (i8 - 55296 << 10) + (e - 56320) + 65536;
}
function ur(i8) {
  return Ki(i8, 0);
}
function Ki(i8, e) {
  switch (typeof i8) {
    case "object":
      return i8 === null ? Le2(349, e) : Array.isArray(i8) ? As(i8, e) : Ss(i8, e);
    case "string":
      return cr(i8, e);
    case "boolean":
      return Ms(i8, e);
    case "number":
      return Le2(i8, e);
    case "undefined":
      return Le2(937, e);
    default:
      return Le2(617, e);
  }
}
function Le2(i8, e) {
  return (e << 5) - e + i8 | 0;
}
function Ms(i8, e) {
  return Le2(i8 ? 433 : 863, e);
}
function cr(i8, e) {
  e = Le2(149417, e);
  for (let t = 0, n = i8.length; t < n; t++) e = Le2(i8.charCodeAt(t), e);
  return e;
}
function As(i8, e) {
  return e = Le2(104579, e), i8.reduce((t, n) => Ki(n, t), e);
}
function Ss(i8, e) {
  return e = Le2(181387, e), Object.keys(i8).sort().reduce((t, n) => (t = cr(n, t), Ki(i8[n], t)), e);
}
function $i(i8, e, t = 32) {
  let n = t - e, s15 = ~((1 << n) - 1);
  return (i8 << e | (s15 & i8) >>> n) >>> 0;
}
function ar(i8, e = 0, t = i8.byteLength, n = 0) {
  for (let s15 = 0; s15 < t; s15++) i8[e + s15] = n;
}
function Os(i8, e, t = "0") {
  for (; i8.length < e; ) i8 = t + i8;
  return i8;
}
function ht(i8, e = 32) {
  return i8 instanceof ArrayBuffer ? Array.from(new Uint8Array(i8)).map((t) => t.toString(16).padStart(2, "0")).join("") : Os((i8 >>> 0).toString(16), e / 4);
}
var ai = class ai2 {
  constructor() {
    this._h0 = 1732584193;
    this._h1 = 4023233417;
    this._h2 = 2562383102;
    this._h3 = 271733878;
    this._h4 = 3285377520;
    this._buff = new Uint8Array(67), this._buffDV = new DataView(this._buff.buffer), this._buffLen = 0, this._totalLen = 0, this._leftoverHighSurrogate = 0, this._finished = false;
  }
  update(e) {
    let t = e.length;
    if (t === 0) return;
    let n = this._buff, s15 = this._buffLen, o2 = this._leftoverHighSurrogate, r5, a;
    for (o2 !== 0 ? (r5 = o2, a = -1, o2 = 0) : (r5 = e.charCodeAt(0), a = 0); ; ) {
      let l = r5;
      if (sr(r5)) if (a + 1 < t) {
        let u = e.charCodeAt(a + 1);
        Gi(u) ? (a++, l = or(r5, u)) : l = 65533;
      } else {
        o2 = r5;
        break;
      }
      else Gi(r5) && (l = 65533);
      if (s15 = this._push(n, s15, l), a++, a < t) r5 = e.charCodeAt(a);
      else break;
    }
    this._buffLen = s15, this._leftoverHighSurrogate = o2;
  }
  _push(e, t, n) {
    return n < 128 ? e[t++] = n : n < 2048 ? (e[t++] = 192 | (n & 1984) >>> 6, e[t++] = 128 | (n & 63) >>> 0) : n < 65536 ? (e[t++] = 224 | (n & 61440) >>> 12, e[t++] = 128 | (n & 4032) >>> 6, e[t++] = 128 | (n & 63) >>> 0) : (e[t++] = 240 | (n & 1835008) >>> 18, e[t++] = 128 | (n & 258048) >>> 12, e[t++] = 128 | (n & 4032) >>> 6, e[t++] = 128 | (n & 63) >>> 0), t >= 64 && (this._step(), t -= 64, this._totalLen += 64, e[0] = e[64], e[1] = e[65], e[2] = e[66]), t;
  }
  digest() {
    return this._finished || (this._finished = true, this._leftoverHighSurrogate && (this._leftoverHighSurrogate = 0, this._buffLen = this._push(this._buff, this._buffLen, 65533)), this._totalLen += this._buffLen, this._wrapUp()), ht(this._h0) + ht(this._h1) + ht(this._h2) + ht(this._h3) + ht(this._h4);
  }
  _wrapUp() {
    this._buff[this._buffLen++] = 128, ar(this._buff, this._buffLen), this._buffLen > 56 && (this._step(), ar(this._buff));
    let e = 8 * this._totalLen;
    this._buffDV.setUint32(56, Math.floor(e / 4294967296), false), this._buffDV.setUint32(60, e % 4294967296, false), this._step();
  }
  _step() {
    let e = ai2._bigBlock32, t = this._buffDV;
    for (let d = 0; d < 64; d += 4) e.setUint32(d, t.getUint32(d, false), false);
    for (let d = 64; d < 320; d += 4) e.setUint32(d, $i(e.getUint32(d - 12, false) ^ e.getUint32(d - 32, false) ^ e.getUint32(d - 56, false) ^ e.getUint32(d - 64, false), 1), false);
    let n = this._h0, s15 = this._h1, o2 = this._h2, r5 = this._h3, a = this._h4, l, u, c;
    for (let d = 0; d < 80; d++) d < 20 ? (l = s15 & o2 | ~s15 & r5, u = 1518500249) : d < 40 ? (l = s15 ^ o2 ^ r5, u = 1859775393) : d < 60 ? (l = s15 & o2 | s15 & r5 | o2 & r5, u = 2400959708) : (l = s15 ^ o2 ^ r5, u = 3395469782), c = $i(n, 5) + l + a + u + e.getUint32(d * 4, false) & 4294967295, a = r5, r5 = o2, o2 = $i(s15, 30), s15 = n, n = c;
    this._h0 = this._h0 + n & 4294967295, this._h1 = this._h1 + s15 & 4294967295, this._h2 = this._h2 + o2 & 4294967295, this._h3 = this._h3 + r5 & 4294967295, this._h4 = this._h4 + a & 4294967295;
  }
};
ai._bigBlock32 = new DataView(new ArrayBuffer(320));
var { registerWindow: fu, getWindow: Fs, getDocument: mu, getWindows: _u, getWindowsCount: bu, getWindowId: dr, getWindowById: vu, hasWindow: Tu, onDidRegisterWindow: gu, onWillUnregisterWindow: xu, onDidUnregisterWindow: Eu } = (function() {
  let i8 = /* @__PURE__ */ new Map();
  te3;
  let e = { window: te3, disposables: new fe2() };
  i8.set(te3.vscodeWindowId, e);
  let t = new D(), n = new D(), s15 = new D();
  function o2(r5, a) {
    return (typeof r5 == "number" ? i8.get(r5) : void 0) ?? (a ? e : void 0);
  }
  return { onDidRegisterWindow: t.event, onWillUnregisterWindow: s15.event, onDidUnregisterWindow: n.event, registerWindow(r5) {
    if (i8.has(r5.vscodeWindowId)) return B.None;
    let a = new fe2(), l = { window: r5, disposables: a.add(new fe2()) };
    return i8.set(r5.vscodeWindowId, l), a.add(O2(() => {
      i8.delete(r5.vscodeWindowId), n.fire(r5);
    })), a.add(li(r5, Ps.BEFORE_UNLOAD, () => {
      s15.fire(r5);
    })), t.fire(l), a;
  }, getWindows() {
    return i8.values();
  }, getWindowsCount() {
    return i8.size;
  }, getWindowId(r5) {
    return r5.vscodeWindowId;
  }, hasWindow(r5) {
    return i8.has(r5);
  }, getWindowById: o2, getWindow(r5) {
    let a = r5;
    if (a?.ownerDocument?.defaultView) return a.ownerDocument.defaultView.window;
    let l = r5;
    return l?.view ? l.view.window : te3;
  }, getDocument(r5) {
    return Fs(r5).document;
  } };
})();
var Vi = class {
  constructor(e, t, n, s15) {
    this._node = e, this._type = t, this._handler = n, this._options = s15 || false, this._node.addEventListener(this._type, this._handler, this._options);
  }
  dispose() {
    this._handler && (this._node.removeEventListener(this._type, this._handler, this._options), this._node = null, this._handler = null);
  }
};
function li(i8, e, t, n) {
  return new Vi(i8, e, t, n);
}
var ks;
var hr;
var pt = class {
  constructor(e, t = 0) {
    this._runner = e, this.priority = t, this._canceled = false;
  }
  dispose() {
    this._canceled = true;
  }
  execute() {
    if (!this._canceled) try {
      this._runner();
    } catch (e) {
      Pe2(e);
    }
  }
  static sort(e, t) {
    return t.priority - e.priority;
  }
};
(function() {
  let i8 = /* @__PURE__ */ new Map(), e = /* @__PURE__ */ new Map(), t = /* @__PURE__ */ new Map(), n = /* @__PURE__ */ new Map(), s15 = (o2) => {
    t.set(o2, false);
    let r5 = i8.get(o2) ?? [];
    for (e.set(o2, r5), i8.set(o2, []), n.set(o2, true); r5.length > 0; ) r5.sort(pt.sort), r5.shift().execute();
    n.set(o2, false);
  };
  hr = (o2, r5, a = 0) => {
    let l = dr(o2), u = new pt(r5, a), c = i8.get(l);
    return c || (c = [], i8.set(l, c)), c.push(u), t.get(l) || (t.set(l, true), o2.requestAnimationFrame(() => s15(l))), u;
  }, ks = (o2, r5, a) => {
    let l = dr(o2);
    if (n.get(l)) {
      let u = new pt(r5, a), c = e.get(l);
      return c || (c = [], e.set(l, c)), c.push(u), u;
    } else return hr(o2, r5, a);
  };
})();
var ke2 = class ke3 {
  constructor(e, t) {
    this.width = e;
    this.height = t;
  }
  with(e = this.width, t = this.height) {
    return e !== this.width || t !== this.height ? new ke3(e, t) : this;
  }
  static is(e) {
    return typeof e == "object" && typeof e.height == "number" && typeof e.width == "number";
  }
  static lift(e) {
    return e instanceof ke3 ? e : new ke3(e.width, e.height);
  }
  static equals(e, t) {
    return e === t ? true : !e || !t ? false : e.width === t.width && e.height === t.height;
  }
};
ke2.None = new ke2(0, 0);
var yu = new class {
  constructor() {
    this.mutationObservers = /* @__PURE__ */ new Map();
  }
  observe(i8, e, t) {
    let n = this.mutationObservers.get(i8);
    n || (n = /* @__PURE__ */ new Map(), this.mutationObservers.set(i8, n));
    let s15 = ur(t), o2 = n.get(s15);
    if (o2) o2.users += 1;
    else {
      let r5 = new D(), a = new MutationObserver((u) => r5.fire(u));
      a.observe(i8, t);
      let l = o2 = { users: 1, observer: a, onDidMutate: r5.event };
      e.add(O2(() => {
        l.users -= 1, l.users === 0 && (r5.dispose(), a.disconnect(), n?.delete(s15), n?.size === 0 && this.mutationObservers.delete(i8));
      })), n.set(s15, o2);
    }
    return o2.onDidMutate;
  }
}();
var Ps = { CLICK: "click", AUXCLICK: "auxclick", DBLCLICK: "dblclick", MOUSE_UP: "mouseup", MOUSE_DOWN: "mousedown", MOUSE_OVER: "mouseover", MOUSE_MOVE: "mousemove", MOUSE_OUT: "mouseout", MOUSE_ENTER: "mouseenter", MOUSE_LEAVE: "mouseleave", MOUSE_WHEEL: "wheel", POINTER_UP: "pointerup", POINTER_DOWN: "pointerdown", POINTER_MOVE: "pointermove", POINTER_LEAVE: "pointerleave", CONTEXT_MENU: "contextmenu", WHEEL: "wheel", KEY_DOWN: "keydown", KEY_PRESS: "keypress", KEY_UP: "keyup", LOAD: "load", BEFORE_UNLOAD: "beforeunload", UNLOAD: "unload", PAGE_SHOW: "pageshow", PAGE_HIDE: "pagehide", PASTE: "paste", ABORT: "abort", ERROR: "error", RESIZE: "resize", SCROLL: "scroll", FULLSCREEN_CHANGE: "fullscreenchange", WK_FULLSCREEN_CHANGE: "webkitfullscreenchange", SELECT: "select", CHANGE: "change", SUBMIT: "submit", RESET: "reset", FOCUS: "focus", FOCUS_IN: "focusin", FOCUS_OUT: "focusout", BLUR: "blur", INPUT: "input", STORAGE: "storage", DRAG_START: "dragstart", DRAG: "drag", DRAG_ENTER: "dragenter", DRAG_LEAVE: "dragleave", DRAG_OVER: "dragover", DROP: "drop", DRAG_END: "dragend", ANIMATION_START: ut ? "webkitAnimationStart" : "animationstart", ANIMATION_END: ut ? "webkitAnimationEnd" : "animationend", ANIMATION_ITERATION: ut ? "webkitAnimationIteration" : "animationiteration" };
var Bs = /([\w\-]+)?(#([\w\-]+))?((\.([\w\-]+))*)/;
function fr(i8, e, t, ...n) {
  let s15 = Bs.exec(e);
  if (!s15) throw new Error("Bad use of emmet");
  let o2 = s15[1] || "div", r5;
  return i8 !== "http://www.w3.org/1999/xhtml" ? r5 = document.createElementNS(i8, o2) : r5 = document.createElement(o2), s15[3] && (r5.id = s15[3]), s15[4] && (r5.className = s15[4].replace(/\./g, " ").trim()), t && Object.entries(t).forEach(([a, l]) => {
    typeof l > "u" || (/^on\w+$/.test(a) ? r5[a] = l : a === "selected" ? l && r5.setAttribute(a, "true") : r5.setAttribute(a, l));
  }), r5.append(...n), r5;
}
function Ns(i8, e, ...t) {
  return fr("http://www.w3.org/1999/xhtml", i8, e, ...t);
}
Ns.SVG = function(i8, e, ...t) {
  return fr("http://www.w3.org/2000/svg", i8, e, ...t);
};
var ui = class extends B {
  constructor(t, n, s15, o2, r5, a, l, u, c) {
    super();
    this._terminal = t;
    this._characterJoinerService = n;
    this._charSizeService = s15;
    this._coreBrowserService = o2;
    this._coreService = r5;
    this._decorationService = a;
    this._optionsService = l;
    this._themeService = u;
    this._cursorBlinkStateManager = new be2();
    this._charAtlasDisposable = this._register(new be2());
    this._observerDisposable = this._register(new be2());
    this._model = new Vt();
    this._workCell = new at2();
    this._workCell2 = new at2();
    this._rectangleRenderer = this._register(new be2());
    this._glyphRenderer = this._register(new be2());
    this._onChangeTextureAtlas = this._register(new D());
    this.onChangeTextureAtlas = this._onChangeTextureAtlas.event;
    this._onAddTextureAtlasCanvas = this._register(new D());
    this.onAddTextureAtlasCanvas = this._onAddTextureAtlasCanvas.event;
    this._onRemoveTextureAtlasCanvas = this._register(new D());
    this.onRemoveTextureAtlasCanvas = this._onRemoveTextureAtlasCanvas.event;
    this._onRequestRedraw = this._register(new D());
    this.onRequestRedraw = this._onRequestRedraw.event;
    this._onContextLoss = this._register(new D());
    this.onContextLoss = this._onContextLoss.event;
    this._canvas = this._coreBrowserService.mainDocument.createElement("canvas");
    let d = { antialias: false, depth: false, preserveDrawingBuffer: c };
    if (this._gl = this._canvas.getContext("webgl2", d), !this._gl) throw new Error("WebGL2 not supported " + this._gl);
    this._register(this._themeService.onChangeColors(() => this._handleColorChange())), this._cellColorResolver = new At(this._terminal, this._optionsService, this._model.selection, this._decorationService, this._coreBrowserService, this._themeService), this._core = this._terminal._core, this._renderLayers = [new Xt(this._core.screenElement, 2, this._terminal, this._core.linkifier, this._coreBrowserService, l, this._themeService)], this.dimensions = _n(), this._devicePixelRatio = this._coreBrowserService.dpr, this._updateDimensions(), this._updateCursorBlink(), this._register(l.onOptionChange(() => this._handleOptionsChanged())), this._deviceMaxTextureSize = this._gl.getParameter(this._gl.MAX_TEXTURE_SIZE), this._register(li(this._canvas, "webglcontextlost", (h2) => {
      console.log("webglcontextlost event received"), h2.preventDefault(), this._contextRestorationTimeout = setTimeout(() => {
        this._contextRestorationTimeout = void 0, console.warn("webgl context not restored; firing onContextLoss"), this._onContextLoss.fire(h2);
      }, 3e3);
    })), this._register(li(this._canvas, "webglcontextrestored", (h2) => {
      console.warn("webglcontextrestored event received"), clearTimeout(this._contextRestorationTimeout), this._contextRestorationTimeout = void 0, Ai(this._terminal), this._initializeWebGLState(), this._requestRedrawViewport();
    })), this._observerDisposable.value = Si(this._canvas, this._coreBrowserService.window, (h2, f) => this._setCanvasDevicePixelDimensions(h2, f)), this._register(this._coreBrowserService.onWindowChange((h2) => {
      this._observerDisposable.value = Si(this._canvas, h2, (f, I2) => this._setCanvasDevicePixelDimensions(f, I2));
    })), this._core.screenElement.appendChild(this._canvas), [this._rectangleRenderer.value, this._glyphRenderer.value] = this._initializeWebGLState(), this._isAttached = this._core.screenElement.isConnected, this._register(O2(() => {
      for (let h2 of this._renderLayers) h2.dispose();
      this._canvas.parentElement?.removeChild(this._canvas), Ai(this._terminal);
    }));
  }
  get textureAtlas() {
    return this._charAtlas?.pages[0].canvas;
  }
  _handleColorChange() {
    this._refreshCharAtlas(), this._clearModel(true);
  }
  handleDevicePixelRatioChange() {
    this._devicePixelRatio !== this._coreBrowserService.dpr && (this._devicePixelRatio = this._coreBrowserService.dpr, this.handleResize(this._terminal.cols, this._terminal.rows));
  }
  handleResize(t, n) {
    this._updateDimensions(), this._model.resize(this._terminal.cols, this._terminal.rows);
    for (let s15 of this._renderLayers) s15.resize(this._terminal, this.dimensions);
    this._canvas.width = this.dimensions.device.canvas.width, this._canvas.height = this.dimensions.device.canvas.height, this._canvas.style.width = `${this.dimensions.css.canvas.width}px`, this._canvas.style.height = `${this.dimensions.css.canvas.height}px`, this._core.screenElement.style.width = `${this.dimensions.css.canvas.width}px`, this._core.screenElement.style.height = `${this.dimensions.css.canvas.height}px`, this._rectangleRenderer.value?.setDimensions(this.dimensions), this._rectangleRenderer.value?.handleResize(), this._glyphRenderer.value?.setDimensions(this.dimensions), this._glyphRenderer.value?.handleResize(), this._refreshCharAtlas(), this._clearModel(false);
  }
  handleCharSizeChanged() {
    this.handleResize(this._terminal.cols, this._terminal.rows);
  }
  handleBlur() {
    for (let t of this._renderLayers) t.handleBlur(this._terminal);
    this._cursorBlinkStateManager.value?.pause(), this._requestRedrawViewport();
  }
  handleFocus() {
    for (let t of this._renderLayers) t.handleFocus(this._terminal);
    this._cursorBlinkStateManager.value?.resume(), this._requestRedrawViewport();
  }
  handleSelectionChanged(t, n, s15) {
    for (let o2 of this._renderLayers) o2.handleSelectionChanged(this._terminal, t, n, s15);
    this._model.selection.update(this._core, t, n, s15), this._requestRedrawViewport();
  }
  handleCursorMove() {
    for (let t of this._renderLayers) t.handleCursorMove(this._terminal);
    this._cursorBlinkStateManager.value?.restartBlinkAnimation();
  }
  _handleOptionsChanged() {
    this._updateDimensions(), this._refreshCharAtlas(), this._updateCursorBlink();
  }
  _initializeWebGLState() {
    return this._rectangleRenderer.value = new qt(this._terminal, this._gl, this.dimensions, this._themeService), this._glyphRenderer.value = new Kt(this._terminal, this._gl, this.dimensions, this._optionsService), this.handleCharSizeChanged(), [this._rectangleRenderer.value, this._glyphRenderer.value];
  }
  _refreshCharAtlas() {
    if (this.dimensions.device.char.width <= 0 && this.dimensions.device.char.height <= 0) {
      this._isAttached = false;
      return;
    }
    let t = Nt(this._terminal, this._optionsService.rawOptions, this._themeService.colors, this.dimensions.device.cell.width, this.dimensions.device.cell.height, this.dimensions.device.char.width, this.dimensions.device.char.height, this._coreBrowserService.dpr, this._deviceMaxTextureSize);
    this._charAtlas !== t && (this._onChangeTextureAtlas.fire(t.pages[0].canvas), this._charAtlasDisposable.value = It(ee3.forward(t.onAddTextureAtlasCanvas, this._onAddTextureAtlasCanvas), ee3.forward(t.onRemoveTextureAtlasCanvas, this._onRemoveTextureAtlasCanvas))), this._charAtlas = t, this._charAtlas.warmUp(), this._glyphRenderer.value?.setAtlas(this._charAtlas);
  }
  _clearModel(t) {
    this._model.clear(), t && this._glyphRenderer.value?.clear();
  }
  clearTextureAtlas() {
    this._charAtlas?.clearTexture(), this._clearModel(true), this._requestRedrawViewport();
  }
  clear() {
    this._clearModel(true);
    for (let t of this._renderLayers) t.reset(this._terminal);
    this._cursorBlinkStateManager.value?.restartBlinkAnimation(), this._updateCursorBlink();
  }
  renderRows(t, n) {
    if (!this._isAttached) if (this._core.screenElement?.isConnected && this._charSizeService.width && this._charSizeService.height) this._updateDimensions(), this._refreshCharAtlas(), this._isAttached = true;
    else return;
    for (let s15 of this._renderLayers) s15.handleGridChanged(this._terminal, t, n);
    !this._glyphRenderer.value || !this._rectangleRenderer.value || (this._glyphRenderer.value.beginFrame() ? (this._clearModel(true), this._updateModel(0, this._terminal.rows - 1)) : this._updateModel(t, n), this._rectangleRenderer.value.renderBackgrounds(), this._glyphRenderer.value.render(this._model), (!this._cursorBlinkStateManager.value || this._cursorBlinkStateManager.value.isCursorVisible) && this._rectangleRenderer.value.renderCursor());
  }
  _updateCursorBlink() {
    this._coreService.decPrivateModes.cursorBlink ?? this._terminal.options.cursorBlink ? this._cursorBlinkStateManager.value = new Ht(() => {
      this._requestRedrawCursor();
    }, this._coreBrowserService) : this._cursorBlinkStateManager.clear(), this._requestRedrawCursor();
  }
  _updateModel(t, n) {
    let s15 = this._core, o2 = this._workCell, r5, a, l, u, c, d, h2 = 0, f = true, I2, L2, M4, q2, S2, W, E, y, w2;
    t = mr(t, s15.rows - 1, 0), n = mr(n, s15.rows - 1, 0);
    let G3 = this._coreService.decPrivateModes.cursorStyle ?? s15.options.cursorStyle ?? "block", ue3 = this._terminal.buffer.active.baseY + this._terminal.buffer.active.cursorY, Se3 = ue3 - s15.buffer.ydisp, ce3 = Math.min(this._terminal.buffer.active.cursorX, s15.cols - 1), we3 = -1, A2 = this._coreService.isCursorInitialized && !this._coreService.isCursorHidden && (!this._cursorBlinkStateManager.value || this._cursorBlinkStateManager.value.isCursorVisible);
    this._model.cursor = void 0;
    let se3 = false;
    for (a = t; a <= n; a++) for (l = a + s15.buffer.ydisp, u = s15.buffer.lines.get(l), this._model.lineLengths[a] = 0, M4 = ue3 === l, h2 = 0, c = this._characterJoinerService.getJoinedCharacters(l), y = 0; y < s15.cols; y++) {
      if (r5 = this._cellColorResolver.result.bg, u.loadCell(y, o2), y === 0 && (r5 = this._cellColorResolver.result.bg), d = false, f = y >= h2, I2 = y, c.length > 0 && y === c[0][0] && f) {
        L2 = c.shift();
        let v2 = this._model.selection.isCellSelected(this._terminal, L2[0], l);
        for (E = L2[0] + 1; E < L2[1]; E++) f &&= v2 === this._model.selection.isCellSelected(this._terminal, E, l);
        f &&= !M4 || ce3 < L2[0] || ce3 >= L2[1], f ? (d = true, o2 = new Ci(o2, u.translateToString(true, L2[0], L2[1]), L2[1] - L2[0]), I2 = L2[1] - 1) : h2 = L2[1];
      }
      if (q2 = o2.getChars(), S2 = o2.getCode(), E = (a * s15.cols + y) * Ce, this._cellColorResolver.resolve(o2, y, l, this.dimensions.device.cell.width), A2 && l === ue3 && (y === ce3 && (this._model.cursor = { x: ce3, y: Se3, width: o2.getWidth(), style: this._coreBrowserService.isFocused ? G3 : s15.options.cursorInactiveStyle, cursorWidth: s15.options.cursorWidth, dpr: this._devicePixelRatio }, we3 = ce3 + o2.getWidth() - 1), y >= ce3 && y <= we3 && (this._coreBrowserService.isFocused && G3 === "block" || this._coreBrowserService.isFocused === false && s15.options.cursorInactiveStyle === "block") && (this._cellColorResolver.result.fg = 50331648 | this._themeService.colors.cursorAccent.rgba >> 8 & 16777215, this._cellColorResolver.result.bg = 50331648 | this._themeService.colors.cursor.rgba >> 8 & 16777215)), S2 !== 0 && (this._model.lineLengths[a] = y + 1), !(this._model.cells[E] === S2 && this._model.cells[E + ze] === this._cellColorResolver.result.bg && this._model.cells[E + qe2] === this._cellColorResolver.result.fg && this._model.cells[E + Ct] === this._cellColorResolver.result.ext) && (se3 = true, q2.length > 1 && (S2 |= Un), this._model.cells[E] = S2, this._model.cells[E + ze] = this._cellColorResolver.result.bg, this._model.cells[E + qe2] = this._cellColorResolver.result.fg, this._model.cells[E + Ct] = this._cellColorResolver.result.ext, W = o2.getWidth(), this._glyphRenderer.value.updateCell(y, a, S2, this._cellColorResolver.result.bg, this._cellColorResolver.result.fg, this._cellColorResolver.result.ext, q2, W, r5), d)) {
        for (o2 = this._workCell, y++; y <= I2; y++) w2 = (a * s15.cols + y) * Ce, this._glyphRenderer.value.updateCell(y, a, 0, 0, 0, 0, pn, 0, 0), this._model.cells[w2] = 0, this._model.cells[w2 + ze] = this._cellColorResolver.result.bg, this._model.cells[w2 + qe2] = this._cellColorResolver.result.fg, this._model.cells[w2 + Ct] = this._cellColorResolver.result.ext;
        y--;
      }
    }
    se3 && this._rectangleRenderer.value.updateBackgrounds(this._model), this._rectangleRenderer.value.updateCursor(this._model);
  }
  _updateDimensions() {
    !this._charSizeService.width || !this._charSizeService.height || (this.dimensions.device.char.width = Math.floor(this._charSizeService.width * this._devicePixelRatio), this.dimensions.device.char.height = Math.ceil(this._charSizeService.height * this._devicePixelRatio), this.dimensions.device.cell.height = Math.floor(this.dimensions.device.char.height * this._optionsService.rawOptions.lineHeight), this.dimensions.device.char.top = this._optionsService.rawOptions.lineHeight === 1 ? 0 : Math.round((this.dimensions.device.cell.height - this.dimensions.device.char.height) / 2), this.dimensions.device.cell.width = this.dimensions.device.char.width + Math.round(this._optionsService.rawOptions.letterSpacing), this.dimensions.device.char.left = Math.floor(this._optionsService.rawOptions.letterSpacing / 2), this.dimensions.device.canvas.height = this._terminal.rows * this.dimensions.device.cell.height, this.dimensions.device.canvas.width = this._terminal.cols * this.dimensions.device.cell.width, this.dimensions.css.canvas.height = Math.round(this.dimensions.device.canvas.height / this._devicePixelRatio), this.dimensions.css.canvas.width = Math.round(this.dimensions.device.canvas.width / this._devicePixelRatio), this.dimensions.css.cell.height = this.dimensions.device.cell.height / this._devicePixelRatio, this.dimensions.css.cell.width = this.dimensions.device.cell.width / this._devicePixelRatio);
  }
  _setCanvasDevicePixelDimensions(t, n) {
    this._canvas.width === t && this._canvas.height === n || (this._canvas.width = t, this._canvas.height = n, this._requestRedrawViewport());
  }
  _requestRedrawViewport() {
    this._onRequestRedraw.fire({ start: 0, end: this._terminal.rows - 1 });
  }
  _requestRedrawCursor() {
    let t = this._terminal.buffer.active.cursorY;
    this._onRequestRedraw.fire({ start: t, end: t });
  }
};
var Ci = class extends he2 {
  constructor(t, n, s15) {
    super();
    this.content = 0;
    this.combinedData = "";
    this.fg = t.fg, this.bg = t.bg, this.combinedData = n, this._width = s15;
  }
  isCombined() {
    return 2097152;
  }
  getWidth() {
    return this._width;
  }
  getChars() {
    return this.combinedData;
  }
  getCode() {
    return 2097151;
  }
  setFromCharData(t) {
    throw new Error("not implemented");
  }
  getAsCharData() {
    return [this.fg, this.getChars(), this.getWidth(), this.getCode()];
  }
};
function mr(i8, e, t = 0) {
  return Math.max(Math.min(i8, e), t);
}
var _r = "di$target";
var br = "di$dependencies";
var zi = /* @__PURE__ */ new Map();
function pe2(i8) {
  if (zi.has(i8)) return zi.get(i8);
  let e = function(t, n, s15) {
    if (arguments.length !== 3) throw new Error("@IServiceName-decorator can only be used to decorate a parameter");
    Us(e, t, s15);
  };
  return e._id = i8, zi.set(i8, e), e;
}
function Us(i8, e, t) {
  e[_r] === e ? e[br].push({ id: i8, index: t }) : (e[br] = [{ id: i8, index: t }], e[_r] = e);
}
var Vu = pe2("BufferService");
var Cu = pe2("CoreMouseService");
var zu = pe2("CoreService");
var qu = pe2("CharsetService");
var ju = pe2("InstantiationService");
var Xu = pe2("LogService");
var vr = pe2("OptionsService");
var Yu = pe2("OscLinkService");
var Qu = pe2("UnicodeService");
var Zu = pe2("DecorationService");
var Hs = { trace: 0, debug: 1, info: 2, warn: 3, error: 4, off: 5 };
var Ws = "xterm.js: ";
var ci = class extends B {
  constructor(t) {
    super();
    this._optionsService = t;
    this._logLevel = 5;
    this._updateLogLevel(), this._register(this._optionsService.onSpecificOptionChange("logLevel", () => this._updateLogLevel())), Tr = this;
  }
  get logLevel() {
    return this._logLevel;
  }
  _updateLogLevel() {
    this._logLevel = Hs[this._optionsService.rawOptions.logLevel];
  }
  _evalLazyOptionalParams(t) {
    for (let n = 0; n < t.length; n++) typeof t[n] == "function" && (t[n] = t[n]());
  }
  _log(t, n, s15) {
    this._evalLazyOptionalParams(s15), t.call(console, (this._optionsService.options.logger ? "" : Ws) + n, ...s15);
  }
  trace(t, ...n) {
    this._logLevel <= 0 && this._log(this._optionsService.options.logger?.trace.bind(this._optionsService.options.logger) ?? console.log, t, n);
  }
  debug(t, ...n) {
    this._logLevel <= 1 && this._log(this._optionsService.options.logger?.debug.bind(this._optionsService.options.logger) ?? console.log, t, n);
  }
  info(t, ...n) {
    this._logLevel <= 2 && this._log(this._optionsService.options.logger?.info.bind(this._optionsService.options.logger) ?? console.info, t, n);
  }
  warn(t, ...n) {
    this._logLevel <= 3 && this._log(this._optionsService.options.logger?.warn.bind(this._optionsService.options.logger) ?? console.warn, t, n);
  }
  error(t, ...n) {
    this._logLevel <= 4 && this._log(this._optionsService.options.logger?.error.bind(this._optionsService.options.logger) ?? console.error, t, n);
  }
};
ci = Yi([Qi(0, vr)], ci);
var Tr;
function gr(i8) {
  Tr = i8;
}
var xr = class extends B {
  constructor(t) {
    if (vi && hn() < 16) {
      let n = { antialias: false, depth: false, preserveDrawingBuffer: true };
      if (!document.createElement("canvas").getContext("webgl2", n)) throw new Error("Webgl2 is only supported on Safari 16 and above");
    }
    super();
    this._preserveDrawingBuffer = t;
    this._onChangeTextureAtlas = this._register(new D());
    this.onChangeTextureAtlas = this._onChangeTextureAtlas.event;
    this._onAddTextureAtlasCanvas = this._register(new D());
    this.onAddTextureAtlasCanvas = this._onAddTextureAtlasCanvas.event;
    this._onRemoveTextureAtlasCanvas = this._register(new D());
    this.onRemoveTextureAtlasCanvas = this._onRemoveTextureAtlasCanvas.event;
    this._onContextLoss = this._register(new D());
    this.onContextLoss = this._onContextLoss.event;
  }
  activate(t) {
    let n = t._core;
    if (!t.element) {
      this._register(n.onWillOpen(() => this.activate(t)));
      return;
    }
    this._terminal = t;
    let s15 = n.coreService, o2 = n.optionsService, r5 = n, a = r5._renderService, l = r5._characterJoinerService, u = r5._charSizeService, c = r5._coreBrowserService, d = r5._decorationService, h2 = r5._logService, f = r5._themeService;
    gr(h2), this._renderer = this._register(new ui(t, l, u, c, s15, d, o2, f, this._preserveDrawingBuffer)), this._register(ee3.forward(this._renderer.onContextLoss, this._onContextLoss)), this._register(ee3.forward(this._renderer.onChangeTextureAtlas, this._onChangeTextureAtlas)), this._register(ee3.forward(this._renderer.onAddTextureAtlasCanvas, this._onAddTextureAtlasCanvas)), this._register(ee3.forward(this._renderer.onRemoveTextureAtlasCanvas, this._onRemoveTextureAtlasCanvas)), a.setRenderer(this._renderer), this._register(O2(() => {
      if (this._terminal._core._store._isDisposed) return;
      let I2 = this._terminal._core._renderService;
      I2.setRenderer(this._terminal._core._createRenderer()), I2.handleResize(t.cols, t.rows);
    }));
  }
  get textureAtlas() {
    return this._renderer?.textureAtlas;
  }
  clearTextureAtlas() {
    this._renderer?.clearTextureAtlas();
  }
};

// node_modules/@xterm/xterm/lib/xterm.mjs
var zs = Object.defineProperty;
var Rl = Object.getOwnPropertyDescriptor;
var Ll = (s15, t) => {
  for (var e in t) zs(s15, e, { get: t[e], enumerable: true });
};
var M3 = (s15, t, e, i8) => {
  for (var r5 = i8 > 1 ? void 0 : i8 ? Rl(t, e) : t, n = s15.length - 1, o2; n >= 0; n--) (o2 = s15[n]) && (r5 = (i8 ? o2(t, e, r5) : o2(r5)) || r5);
  return i8 && r5 && zs(t, e, r5), r5;
};
var S = (s15, t) => (e, i8) => t(e, i8, s15);
var Gs = "Terminal input";
var mi2 = { get: () => Gs, set: (s15) => Gs = s15 };
var $s = "Too much output to announce, navigate to rows manually to read";
var _i2 = { get: () => $s, set: (s15) => $s = s15 };
function Al(s15) {
  return s15.replace(/\r?\n/g, "\r");
}
function kl(s15, t) {
  return t ? "\x1B[200~" + s15 + "\x1B[201~" : s15;
}
function Vs(s15, t) {
  s15.clipboardData && s15.clipboardData.setData("text/plain", t.selectionText), s15.preventDefault();
}
function qs(s15, t, e, i8) {
  if (s15.stopPropagation(), s15.clipboardData) {
    let r5 = s15.clipboardData.getData("text/plain");
    Cn2(r5, t, e, i8);
  }
}
function Cn2(s15, t, e, i8) {
  s15 = Al(s15), s15 = kl(s15, e.decPrivateModes.bracketedPasteMode && i8.rawOptions.ignoreBracketedPasteMode !== true), e.triggerDataEvent(s15, true), t.value = "";
}
function Mn2(s15, t, e) {
  let i8 = e.getBoundingClientRect(), r5 = s15.clientX - i8.left - 10, n = s15.clientY - i8.top - 10;
  t.style.width = "20px", t.style.height = "20px", t.style.left = `${r5}px`, t.style.top = `${n}px`, t.style.zIndex = "1000", t.focus();
}
function Pn2(s15, t, e, i8, r5) {
  Mn2(s15, t, e), r5 && i8.rightClickSelect(s15), t.value = i8.selectionText, t.select();
}
function Ce2(s15) {
  return s15 > 65535 ? (s15 -= 65536, String.fromCharCode((s15 >> 10) + 55296) + String.fromCharCode(s15 % 1024 + 56320)) : String.fromCharCode(s15);
}
function It2(s15, t = 0, e = s15.length) {
  let i8 = "";
  for (let r5 = t; r5 < e; ++r5) {
    let n = s15[r5];
    n > 65535 ? (n -= 65536, i8 += String.fromCharCode((n >> 10) + 55296) + String.fromCharCode(n % 1024 + 56320)) : i8 += String.fromCharCode(n);
  }
  return i8;
}
var er2 = class {
  constructor() {
    this._interim = 0;
  }
  clear() {
    this._interim = 0;
  }
  decode(t, e) {
    let i8 = t.length;
    if (!i8) return 0;
    let r5 = 0, n = 0;
    if (this._interim) {
      let o2 = t.charCodeAt(n++);
      56320 <= o2 && o2 <= 57343 ? e[r5++] = (this._interim - 55296) * 1024 + o2 - 56320 + 65536 : (e[r5++] = this._interim, e[r5++] = o2), this._interim = 0;
    }
    for (let o2 = n; o2 < i8; ++o2) {
      let l = t.charCodeAt(o2);
      if (55296 <= l && l <= 56319) {
        if (++o2 >= i8) return this._interim = l, r5;
        let a = t.charCodeAt(o2);
        56320 <= a && a <= 57343 ? e[r5++] = (l - 55296) * 1024 + a - 56320 + 65536 : (e[r5++] = l, e[r5++] = a);
        continue;
      }
      l !== 65279 && (e[r5++] = l);
    }
    return r5;
  }
};
var tr2 = class {
  constructor() {
    this.interim = new Uint8Array(3);
  }
  clear() {
    this.interim.fill(0);
  }
  decode(t, e) {
    let i8 = t.length;
    if (!i8) return 0;
    let r5 = 0, n, o2, l, a, u = 0, h2 = 0;
    if (this.interim[0]) {
      let _3 = false, p = this.interim[0];
      p &= (p & 224) === 192 ? 31 : (p & 240) === 224 ? 15 : 7;
      let m = 0, f;
      for (; (f = this.interim[++m] & 63) && m < 4; ) p <<= 6, p |= f;
      let A2 = (this.interim[0] & 224) === 192 ? 2 : (this.interim[0] & 240) === 224 ? 3 : 4, R3 = A2 - m;
      for (; h2 < R3; ) {
        if (h2 >= i8) return 0;
        if (f = t[h2++], (f & 192) !== 128) {
          h2--, _3 = true;
          break;
        } else this.interim[m++] = f, p <<= 6, p |= f & 63;
      }
      _3 || (A2 === 2 ? p < 128 ? h2-- : e[r5++] = p : A2 === 3 ? p < 2048 || p >= 55296 && p <= 57343 || p === 65279 || (e[r5++] = p) : p < 65536 || p > 1114111 || (e[r5++] = p)), this.interim.fill(0);
    }
    let c = i8 - 4, d = h2;
    for (; d < i8; ) {
      for (; d < c && !((n = t[d]) & 128) && !((o2 = t[d + 1]) & 128) && !((l = t[d + 2]) & 128) && !((a = t[d + 3]) & 128); ) e[r5++] = n, e[r5++] = o2, e[r5++] = l, e[r5++] = a, d += 4;
      if (n = t[d++], n < 128) e[r5++] = n;
      else if ((n & 224) === 192) {
        if (d >= i8) return this.interim[0] = n, r5;
        if (o2 = t[d++], (o2 & 192) !== 128) {
          d--;
          continue;
        }
        if (u = (n & 31) << 6 | o2 & 63, u < 128) {
          d--;
          continue;
        }
        e[r5++] = u;
      } else if ((n & 240) === 224) {
        if (d >= i8) return this.interim[0] = n, r5;
        if (o2 = t[d++], (o2 & 192) !== 128) {
          d--;
          continue;
        }
        if (d >= i8) return this.interim[0] = n, this.interim[1] = o2, r5;
        if (l = t[d++], (l & 192) !== 128) {
          d--;
          continue;
        }
        if (u = (n & 15) << 12 | (o2 & 63) << 6 | l & 63, u < 2048 || u >= 55296 && u <= 57343 || u === 65279) continue;
        e[r5++] = u;
      } else if ((n & 248) === 240) {
        if (d >= i8) return this.interim[0] = n, r5;
        if (o2 = t[d++], (o2 & 192) !== 128) {
          d--;
          continue;
        }
        if (d >= i8) return this.interim[0] = n, this.interim[1] = o2, r5;
        if (l = t[d++], (l & 192) !== 128) {
          d--;
          continue;
        }
        if (d >= i8) return this.interim[0] = n, this.interim[1] = o2, this.interim[2] = l, r5;
        if (a = t[d++], (a & 192) !== 128) {
          d--;
          continue;
        }
        if (u = (n & 7) << 18 | (o2 & 63) << 12 | (l & 63) << 6 | a & 63, u < 65536 || u > 1114111) continue;
        e[r5++] = u;
      }
    }
    return r5;
  }
};
var ir = "";
var we2 = " ";
var De3 = class s {
  constructor() {
    this.fg = 0;
    this.bg = 0;
    this.extended = new rt3();
  }
  static toColorRGB(t) {
    return [t >>> 16 & 255, t >>> 8 & 255, t & 255];
  }
  static fromColorRGB(t) {
    return (t[0] & 255) << 16 | (t[1] & 255) << 8 | t[2] & 255;
  }
  clone() {
    let t = new s();
    return t.fg = this.fg, t.bg = this.bg, t.extended = this.extended.clone(), t;
  }
  isInverse() {
    return this.fg & 67108864;
  }
  isBold() {
    return this.fg & 134217728;
  }
  isUnderline() {
    return this.hasExtendedAttrs() && this.extended.underlineStyle !== 0 ? 1 : this.fg & 268435456;
  }
  isBlink() {
    return this.fg & 536870912;
  }
  isInvisible() {
    return this.fg & 1073741824;
  }
  isItalic() {
    return this.bg & 67108864;
  }
  isDim() {
    return this.bg & 134217728;
  }
  isStrikethrough() {
    return this.fg & 2147483648;
  }
  isProtected() {
    return this.bg & 536870912;
  }
  isOverline() {
    return this.bg & 1073741824;
  }
  getFgColorMode() {
    return this.fg & 50331648;
  }
  getBgColorMode() {
    return this.bg & 50331648;
  }
  isFgRGB() {
    return (this.fg & 50331648) === 50331648;
  }
  isBgRGB() {
    return (this.bg & 50331648) === 50331648;
  }
  isFgPalette() {
    return (this.fg & 50331648) === 16777216 || (this.fg & 50331648) === 33554432;
  }
  isBgPalette() {
    return (this.bg & 50331648) === 16777216 || (this.bg & 50331648) === 33554432;
  }
  isFgDefault() {
    return (this.fg & 50331648) === 0;
  }
  isBgDefault() {
    return (this.bg & 50331648) === 0;
  }
  isAttributeDefault() {
    return this.fg === 0 && this.bg === 0;
  }
  getFgColor() {
    switch (this.fg & 50331648) {
      case 16777216:
      case 33554432:
        return this.fg & 255;
      case 50331648:
        return this.fg & 16777215;
      default:
        return -1;
    }
  }
  getBgColor() {
    switch (this.bg & 50331648) {
      case 16777216:
      case 33554432:
        return this.bg & 255;
      case 50331648:
        return this.bg & 16777215;
      default:
        return -1;
    }
  }
  hasExtendedAttrs() {
    return this.bg & 268435456;
  }
  updateExtended() {
    this.extended.isEmpty() ? this.bg &= -268435457 : this.bg |= 268435456;
  }
  getUnderlineColor() {
    if (this.bg & 268435456 && ~this.extended.underlineColor) switch (this.extended.underlineColor & 50331648) {
      case 16777216:
      case 33554432:
        return this.extended.underlineColor & 255;
      case 50331648:
        return this.extended.underlineColor & 16777215;
      default:
        return this.getFgColor();
    }
    return this.getFgColor();
  }
  getUnderlineColorMode() {
    return this.bg & 268435456 && ~this.extended.underlineColor ? this.extended.underlineColor & 50331648 : this.getFgColorMode();
  }
  isUnderlineColorRGB() {
    return this.bg & 268435456 && ~this.extended.underlineColor ? (this.extended.underlineColor & 50331648) === 50331648 : this.isFgRGB();
  }
  isUnderlineColorPalette() {
    return this.bg & 268435456 && ~this.extended.underlineColor ? (this.extended.underlineColor & 50331648) === 16777216 || (this.extended.underlineColor & 50331648) === 33554432 : this.isFgPalette();
  }
  isUnderlineColorDefault() {
    return this.bg & 268435456 && ~this.extended.underlineColor ? (this.extended.underlineColor & 50331648) === 0 : this.isFgDefault();
  }
  getUnderlineStyle() {
    return this.fg & 268435456 ? this.bg & 268435456 ? this.extended.underlineStyle : 1 : 0;
  }
  getUnderlineVariantOffset() {
    return this.extended.underlineVariantOffset;
  }
};
var rt3 = class s2 {
  constructor(t = 0, e = 0) {
    this._ext = 0;
    this._urlId = 0;
    this._ext = t, this._urlId = e;
  }
  get ext() {
    return this._urlId ? this._ext & -469762049 | this.underlineStyle << 26 : this._ext;
  }
  set ext(t) {
    this._ext = t;
  }
  get underlineStyle() {
    return this._urlId ? 5 : (this._ext & 469762048) >> 26;
  }
  set underlineStyle(t) {
    this._ext &= -469762049, this._ext |= t << 26 & 469762048;
  }
  get underlineColor() {
    return this._ext & 67108863;
  }
  set underlineColor(t) {
    this._ext &= -67108864, this._ext |= t & 67108863;
  }
  get urlId() {
    return this._urlId;
  }
  set urlId(t) {
    this._urlId = t;
  }
  get underlineVariantOffset() {
    let t = (this._ext & 3758096384) >> 29;
    return t < 0 ? t ^ 4294967288 : t;
  }
  set underlineVariantOffset(t) {
    this._ext &= 536870911, this._ext |= t << 29 & 3758096384;
  }
  clone() {
    return new s2(this._ext, this._urlId);
  }
  isEmpty() {
    return this.underlineStyle === 0 && this._urlId === 0;
  }
};
var q = class s3 extends De3 {
  constructor() {
    super(...arguments);
    this.content = 0;
    this.fg = 0;
    this.bg = 0;
    this.extended = new rt3();
    this.combinedData = "";
  }
  static fromCharData(e) {
    let i8 = new s3();
    return i8.setFromCharData(e), i8;
  }
  isCombined() {
    return this.content & 2097152;
  }
  getWidth() {
    return this.content >> 22;
  }
  getChars() {
    return this.content & 2097152 ? this.combinedData : this.content & 2097151 ? Ce2(this.content & 2097151) : "";
  }
  getCode() {
    return this.isCombined() ? this.combinedData.charCodeAt(this.combinedData.length - 1) : this.content & 2097151;
  }
  setFromCharData(e) {
    this.fg = e[0], this.bg = 0;
    let i8 = false;
    if (e[1].length > 2) i8 = true;
    else if (e[1].length === 2) {
      let r5 = e[1].charCodeAt(0);
      if (55296 <= r5 && r5 <= 56319) {
        let n = e[1].charCodeAt(1);
        56320 <= n && n <= 57343 ? this.content = (r5 - 55296) * 1024 + n - 56320 + 65536 | e[2] << 22 : i8 = true;
      } else i8 = true;
    } else this.content = e[1].charCodeAt(0) | e[2] << 22;
    i8 && (this.combinedData = e[1], this.content = 2097152 | e[2] << 22);
  }
  getAsCharData() {
    return [this.fg, this.getChars(), this.getWidth(), this.getCode()];
  }
};
var js = "di$target";
var Hn2 = "di$dependencies";
var Fn2 = /* @__PURE__ */ new Map();
function Xs(s15) {
  return s15[Hn2] || [];
}
function ie3(s15) {
  if (Fn2.has(s15)) return Fn2.get(s15);
  let t = function(e, i8, r5) {
    if (arguments.length !== 3) throw new Error("@IServiceName-decorator can only be used to decorate a parameter");
    Pl(t, e, r5);
  };
  return t._id = s15, Fn2.set(s15, t), t;
}
function Pl(s15, t, e) {
  t[js] === t ? t[Hn2].push({ id: s15, index: e }) : (t[Hn2] = [{ id: s15, index: e }], t[js] = t);
}
var F2 = ie3("BufferService");
var rr = ie3("CoreMouseService");
var ge3 = ie3("CoreService");
var Zs = ie3("CharsetService");
var xt4 = ie3("InstantiationService");
var nr2 = ie3("LogService");
var H3 = ie3("OptionsService");
var sr2 = ie3("OscLinkService");
var Js = ie3("UnicodeService");
var Be4 = ie3("DecorationService");
var wt2 = class {
  constructor(t, e, i8) {
    this._bufferService = t;
    this._optionsService = e;
    this._oscLinkService = i8;
  }
  provideLinks(t, e) {
    let i8 = this._bufferService.buffer.lines.get(t - 1);
    if (!i8) {
      e(void 0);
      return;
    }
    let r5 = [], n = this._optionsService.rawOptions.linkHandler, o2 = new q(), l = i8.getTrimmedLength(), a = -1, u = -1, h2 = false;
    for (let c = 0; c < l; c++) if (!(u === -1 && !i8.hasContent(c))) {
      if (i8.loadCell(c, o2), o2.hasExtendedAttrs() && o2.extended.urlId) if (u === -1) {
        u = c, a = o2.extended.urlId;
        continue;
      } else h2 = o2.extended.urlId !== a;
      else u !== -1 && (h2 = true);
      if (h2 || u !== -1 && c === l - 1) {
        let d = this._oscLinkService.getLinkData(a)?.uri;
        if (d) {
          let _3 = { start: { x: u + 1, y: t }, end: { x: c + (!h2 && c === l - 1 ? 1 : 0), y: t } }, p = false;
          if (!n?.allowNonHttpProtocols) try {
            let m = new URL(d);
            ["http:", "https:"].includes(m.protocol) || (p = true);
          } catch {
            p = true;
          }
          p || r5.push({ text: d, range: _3, activate: (m, f) => n ? n.activate(m, f, _3) : Ol(m, f), hover: (m, f) => n?.hover?.(m, f, _3), leave: (m, f) => n?.leave?.(m, f, _3) });
        }
        h2 = false, o2.hasExtendedAttrs() && o2.extended.urlId ? (u = c, a = o2.extended.urlId) : (u = -1, a = -1);
      }
    }
    e(r5);
  }
};
wt2 = M3([S(0, F2), S(1, H3), S(2, sr2)], wt2);
function Ol(s15, t) {
  if (confirm(`Do you want to navigate to ${t}?

WARNING: This link could potentially be dangerous`)) {
    let i8 = window.open();
    if (i8) {
      try {
        i8.opener = null;
      } catch {
      }
      i8.location.href = t;
    } else console.warn("Opening link blocked as opener could not be cleared");
  }
}
var nt3 = ie3("CharSizeService");
var ae2 = ie3("CoreBrowserService");
var Dt2 = ie3("MouseService");
var ce2 = ie3("RenderService");
var Qs = ie3("SelectionService");
var or2 = ie3("CharacterJoinerService");
var Re2 = ie3("ThemeService");
var lr = ie3("LinkProviderService");
var Wn2 = class {
  constructor() {
    this.listeners = [], this.unexpectedErrorHandler = function(t) {
      setTimeout(() => {
        throw t.stack ? ar2.isErrorNoTelemetry(t) ? new ar2(t.message + `

` + t.stack) : new Error(t.message + `

` + t.stack) : t;
      }, 0);
    };
  }
  addListener(t) {
    return this.listeners.push(t), () => {
      this._removeListener(t);
    };
  }
  emit(t) {
    this.listeners.forEach((e) => {
      e(t);
    });
  }
  _removeListener(t) {
    this.listeners.splice(this.listeners.indexOf(t), 1);
  }
  setUnexpectedErrorHandler(t) {
    this.unexpectedErrorHandler = t;
  }
  getUnexpectedErrorHandler() {
    return this.unexpectedErrorHandler;
  }
  onUnexpectedError(t) {
    this.unexpectedErrorHandler(t), this.emit(t);
  }
  onUnexpectedExternalError(t) {
    this.unexpectedErrorHandler(t);
  }
};
var Bl = new Wn2();
function Lt2(s15) {
  Nl(s15) || Bl.onUnexpectedError(s15);
}
var Un2 = "Canceled";
function Nl(s15) {
  return s15 instanceof bi2 ? true : s15 instanceof Error && s15.name === Un2 && s15.message === Un2;
}
var bi2 = class extends Error {
  constructor() {
    super(Un2), this.name = this.message;
  }
};
function eo(s15) {
  return s15 ? new Error(`Illegal argument: ${s15}`) : new Error("Illegal argument");
}
var ar2 = class s4 extends Error {
  constructor(t) {
    super(t), this.name = "CodeExpectedError";
  }
  static fromError(t) {
    if (t instanceof s4) return t;
    let e = new s4();
    return e.message = t.message, e.stack = t.stack, e;
  }
  static isErrorNoTelemetry(t) {
    return t.name === "CodeExpectedError";
  }
};
var Rt2 = class s5 extends Error {
  constructor(t) {
    super(t || "An unexpected bug occurred."), Object.setPrototypeOf(this, s5.prototype);
  }
};
function Fl(s15, t, e = 0, i8 = s15.length) {
  let r5 = e, n = i8;
  for (; r5 < n; ) {
    let o2 = Math.floor((r5 + n) / 2);
    t(s15[o2]) ? r5 = o2 + 1 : n = o2;
  }
  return r5 - 1;
}
var cr2 = class cr3 {
  constructor(t) {
    this._array = t;
    this._findLastMonotonousLastIdx = 0;
  }
  findLastMonotonous(t) {
    if (cr3.assertInvariants) {
      if (this._prevFindLastPredicate) {
        for (let i8 of this._array) if (this._prevFindLastPredicate(i8) && !t(i8)) throw new Error("MonotonousArray: current predicate must be weaker than (or equal to) the previous predicate.");
      }
      this._prevFindLastPredicate = t;
    }
    let e = Fl(this._array, t, this._findLastMonotonousLastIdx);
    return this._findLastMonotonousLastIdx = e + 1, e === -1 ? void 0 : this._array[e];
  }
};
cr2.assertInvariants = false;
function Se2(s15, t = 0) {
  return s15[s15.length - (1 + t)];
}
var ro;
((l) => {
  function s15(a) {
    return a < 0;
  }
  l.isLessThan = s15;
  function t(a) {
    return a <= 0;
  }
  l.isLessThanOrEqual = t;
  function e(a) {
    return a > 0;
  }
  l.isGreaterThan = e;
  function i8(a) {
    return a === 0;
  }
  l.isNeitherLessOrGreaterThan = i8, l.greaterThan = 1, l.lessThan = -1, l.neitherLessOrGreaterThan = 0;
})(ro ||= {});
function no(s15, t) {
  return (e, i8) => t(s15(e), s15(i8));
}
var so = (s15, t) => s15 - t;
var At2 = class At3 {
  constructor(t) {
    this.iterate = t;
  }
  forEach(t) {
    this.iterate((e) => (t(e), true));
  }
  toArray() {
    let t = [];
    return this.iterate((e) => (t.push(e), true)), t;
  }
  filter(t) {
    return new At3((e) => this.iterate((i8) => t(i8) ? e(i8) : true));
  }
  map(t) {
    return new At3((e) => this.iterate((i8) => e(t(i8))));
  }
  some(t) {
    let e = false;
    return this.iterate((i8) => (e = t(i8), !e)), e;
  }
  findFirst(t) {
    let e;
    return this.iterate((i8) => t(i8) ? (e = i8, false) : true), e;
  }
  findLast(t) {
    let e;
    return this.iterate((i8) => (t(i8) && (e = i8), true)), e;
  }
  findLastMaxBy(t) {
    let e, i8 = true;
    return this.iterate((r5) => ((i8 || ro.isGreaterThan(t(r5, e))) && (i8 = false, e = r5), true)), e;
  }
};
At2.empty = new At2((t) => {
});
function co(s15, t) {
  let e = /* @__PURE__ */ Object.create(null);
  for (let i8 of s15) {
    let r5 = t(i8), n = e[r5];
    n || (n = e[r5] = []), n.push(i8);
  }
  return e;
}
var lo2;
var ao2;
var oo2 = class {
  constructor(t, e) {
    this.toKey = e;
    this._map = /* @__PURE__ */ new Map();
    this[lo2] = "SetWithKey";
    for (let i8 of t) this.add(i8);
  }
  get size() {
    return this._map.size;
  }
  add(t) {
    let e = this.toKey(t);
    return this._map.set(e, t), this;
  }
  delete(t) {
    return this._map.delete(this.toKey(t));
  }
  has(t) {
    return this._map.has(this.toKey(t));
  }
  *entries() {
    for (let t of this._map.values()) yield [t, t];
  }
  keys() {
    return this.values();
  }
  *values() {
    for (let t of this._map.values()) yield t;
  }
  clear() {
    this._map.clear();
  }
  forEach(t, e) {
    this._map.forEach((i8) => t.call(e, i8, i8, this));
  }
  [(ao2 = Symbol.iterator, lo2 = Symbol.toStringTag, ao2)]() {
    return this.values();
  }
};
var ur2 = class {
  constructor() {
    this.map = /* @__PURE__ */ new Map();
  }
  add(t, e) {
    let i8 = this.map.get(t);
    i8 || (i8 = /* @__PURE__ */ new Set(), this.map.set(t, i8)), i8.add(e);
  }
  delete(t, e) {
    let i8 = this.map.get(t);
    i8 && (i8.delete(e), i8.size === 0 && this.map.delete(t));
  }
  forEach(t, e) {
    let i8 = this.map.get(t);
    i8 && i8.forEach(e);
  }
  get(t) {
    let e = this.map.get(t);
    return e || /* @__PURE__ */ new Set();
  }
};
function Kn2(s15, t) {
  let e = this, i8 = false, r5;
  return function() {
    if (i8) return r5;
    if (i8 = true, t) try {
      r5 = s15.apply(e, arguments);
    } finally {
      t();
    }
    else r5 = s15.apply(e, arguments);
    return r5;
  };
}
var zn2;
((O3) => {
  function s15(I2) {
    return I2 && typeof I2 == "object" && typeof I2[Symbol.iterator] == "function";
  }
  O3.is = s15;
  let t = Object.freeze([]);
  function e() {
    return t;
  }
  O3.empty = e;
  function* i8(I2) {
    yield I2;
  }
  O3.single = i8;
  function r5(I2) {
    return s15(I2) ? I2 : i8(I2);
  }
  O3.wrap = r5;
  function n(I2) {
    return I2 || t;
  }
  O3.from = n;
  function* o2(I2) {
    for (let k2 = I2.length - 1; k2 >= 0; k2--) yield I2[k2];
  }
  O3.reverse = o2;
  function l(I2) {
    return !I2 || I2[Symbol.iterator]().next().done === true;
  }
  O3.isEmpty = l;
  function a(I2) {
    return I2[Symbol.iterator]().next().value;
  }
  O3.first = a;
  function u(I2, k2) {
    let P3 = 0;
    for (let oe of I2) if (k2(oe, P3++)) return true;
    return false;
  }
  O3.some = u;
  function h2(I2, k2) {
    for (let P3 of I2) if (k2(P3)) return P3;
  }
  O3.find = h2;
  function* c(I2, k2) {
    for (let P3 of I2) k2(P3) && (yield P3);
  }
  O3.filter = c;
  function* d(I2, k2) {
    let P3 = 0;
    for (let oe of I2) yield k2(oe, P3++);
  }
  O3.map = d;
  function* _3(I2, k2) {
    let P3 = 0;
    for (let oe of I2) yield* k2(oe, P3++);
  }
  O3.flatMap = _3;
  function* p(...I2) {
    for (let k2 of I2) yield* k2;
  }
  O3.concat = p;
  function m(I2, k2, P3) {
    let oe = P3;
    for (let Me3 of I2) oe = k2(oe, Me3);
    return oe;
  }
  O3.reduce = m;
  function* f(I2, k2, P3 = I2.length) {
    for (k2 < 0 && (k2 += I2.length), P3 < 0 ? P3 += I2.length : P3 > I2.length && (P3 = I2.length); k2 < P3; k2++) yield I2[k2];
  }
  O3.slice = f;
  function A2(I2, k2 = Number.POSITIVE_INFINITY) {
    let P3 = [];
    if (k2 === 0) return [P3, I2];
    let oe = I2[Symbol.iterator]();
    for (let Me3 = 0; Me3 < k2; Me3++) {
      let Pe3 = oe.next();
      if (Pe3.done) return [P3, O3.empty()];
      P3.push(Pe3.value);
    }
    return [P3, { [Symbol.iterator]() {
      return oe;
    } }];
  }
  O3.consume = A2;
  async function R3(I2) {
    let k2 = [];
    for await (let P3 of I2) k2.push(P3);
    return Promise.resolve(k2);
  }
  O3.asyncToArray = R3;
})(zn2 ||= {});
var Wl = false;
var dt2 = null;
var hr2 = class hr3 {
  constructor() {
    this.livingDisposables = /* @__PURE__ */ new Map();
  }
  getDisposableData(t) {
    let e = this.livingDisposables.get(t);
    return e || (e = { parent: null, source: null, isSingleton: false, value: t, idx: hr3.idx++ }, this.livingDisposables.set(t, e)), e;
  }
  trackDisposable(t) {
    let e = this.getDisposableData(t);
    e.source || (e.source = new Error().stack);
  }
  setParent(t, e) {
    let i8 = this.getDisposableData(t);
    i8.parent = e;
  }
  markAsDisposed(t) {
    this.livingDisposables.delete(t);
  }
  markAsSingleton(t) {
    this.getDisposableData(t).isSingleton = true;
  }
  getRootParent(t, e) {
    let i8 = e.get(t);
    if (i8) return i8;
    let r5 = t.parent ? this.getRootParent(this.getDisposableData(t.parent), e) : t;
    return e.set(t, r5), r5;
  }
  getTrackedDisposables() {
    let t = /* @__PURE__ */ new Map();
    return [...this.livingDisposables.entries()].filter(([, i8]) => i8.source !== null && !this.getRootParent(i8, t).isSingleton).flatMap(([i8]) => i8);
  }
  computeLeakingDisposables(t = 10, e) {
    let i8;
    if (e) i8 = e;
    else {
      let a = /* @__PURE__ */ new Map(), u = [...this.livingDisposables.values()].filter((c) => c.source !== null && !this.getRootParent(c, a).isSingleton);
      if (u.length === 0) return;
      let h2 = new Set(u.map((c) => c.value));
      if (i8 = u.filter((c) => !(c.parent && h2.has(c.parent))), i8.length === 0) throw new Error("There are cyclic diposable chains!");
    }
    if (!i8) return;
    function r5(a) {
      function u(c, d) {
        for (; c.length > 0 && d.some((_3) => typeof _3 == "string" ? _3 === c[0] : c[0].match(_3)); ) c.shift();
      }
      let h2 = a.source.split(`
`).map((c) => c.trim().replace("at ", "")).filter((c) => c !== "");
      return u(h2, ["Error", /^trackDisposable \(.*\)$/, /^DisposableTracker.trackDisposable \(.*\)$/]), h2.reverse();
    }
    let n = new ur2();
    for (let a of i8) {
      let u = r5(a);
      for (let h2 = 0; h2 <= u.length; h2++) n.add(u.slice(0, h2).join(`
`), a);
    }
    i8.sort(no((a) => a.idx, so));
    let o2 = "", l = 0;
    for (let a of i8.slice(0, t)) {
      l++;
      let u = r5(a), h2 = [];
      for (let c = 0; c < u.length; c++) {
        let d = u[c];
        d = `(shared with ${n.get(u.slice(0, c + 1).join(`
`)).size}/${i8.length} leaks) at ${d}`;
        let p = n.get(u.slice(0, c).join(`
`)), m = co([...p].map((f) => r5(f)[c]), (f) => f);
        delete m[u[c]];
        for (let [f, A2] of Object.entries(m)) h2.unshift(`    - stacktraces of ${A2.length} other leaks continue with ${f}`);
        h2.unshift(d);
      }
      o2 += `


==================== Leaking disposable ${l}/${i8.length}: ${a.value.constructor.name} ====================
${h2.join(`
`)}
============================================================

`;
    }
    return i8.length > t && (o2 += `


... and ${i8.length - t} more leaking disposables

`), { leaks: i8, details: o2 };
  }
};
hr2.idx = 0;
function Ul(s15) {
  dt2 = s15;
}
if (Wl) {
  let s15 = "__is_disposable_tracked__";
  Ul(new class {
    trackDisposable(t) {
      let e = new Error("Potentially leaked disposable").stack;
      setTimeout(() => {
        t[s15] || console.log(e);
      }, 3e3);
    }
    setParent(t, e) {
      if (t && t !== D2.None) try {
        t[s15] = true;
      } catch {
      }
    }
    markAsDisposed(t) {
      if (t && t !== D2.None) try {
        t[s15] = true;
      } catch {
      }
    }
    markAsSingleton(t) {
    }
  }());
}
function fr2(s15) {
  return dt2?.trackDisposable(s15), s15;
}
function pr(s15) {
  dt2?.markAsDisposed(s15);
}
function vi2(s15, t) {
  dt2?.setParent(s15, t);
}
function Kl(s15, t) {
  if (dt2) for (let e of s15) dt2.setParent(e, t);
}
function Gn2(s15) {
  return dt2?.markAsSingleton(s15), s15;
}
function Ne2(s15) {
  if (zn2.is(s15)) {
    let t = [];
    for (let e of s15) if (e) try {
      e.dispose();
    } catch (i8) {
      t.push(i8);
    }
    if (t.length === 1) throw t[0];
    if (t.length > 1) throw new AggregateError(t, "Encountered errors while disposing of store");
    return Array.isArray(s15) ? [] : s15;
  } else if (s15) return s15.dispose(), s15;
}
function ho(...s15) {
  let t = C3(() => Ne2(s15));
  return Kl(s15, t), t;
}
function C3(s15) {
  let t = fr2({ dispose: Kn2(() => {
    pr(t), s15();
  }) });
  return t;
}
var dr2 = class dr3 {
  constructor() {
    this._toDispose = /* @__PURE__ */ new Set();
    this._isDisposed = false;
    fr2(this);
  }
  dispose() {
    this._isDisposed || (pr(this), this._isDisposed = true, this.clear());
  }
  get isDisposed() {
    return this._isDisposed;
  }
  clear() {
    if (this._toDispose.size !== 0) try {
      Ne2(this._toDispose);
    } finally {
      this._toDispose.clear();
    }
  }
  add(t) {
    if (!t) return t;
    if (t === this) throw new Error("Cannot register a disposable on itself!");
    return vi2(t, this), this._isDisposed ? dr3.DISABLE_DISPOSED_WARNING || console.warn(new Error("Trying to add a disposable to a DisposableStore that has already been disposed of. The added object will be leaked!").stack) : this._toDispose.add(t), t;
  }
  delete(t) {
    if (t) {
      if (t === this) throw new Error("Cannot dispose a disposable on itself!");
      this._toDispose.delete(t), t.dispose();
    }
  }
  deleteAndLeak(t) {
    t && this._toDispose.has(t) && (this._toDispose.delete(t), vi2(t, null));
  }
};
dr2.DISABLE_DISPOSED_WARNING = false;
var Ee3 = dr2;
var D2 = class {
  constructor() {
    this._store = new Ee3();
    fr2(this), vi2(this._store, this);
  }
  dispose() {
    pr(this), this._store.dispose();
  }
  _register(t) {
    if (t === this) throw new Error("Cannot register a disposable on itself!");
    return this._store.add(t);
  }
};
D2.None = Object.freeze({ dispose() {
} });
var ye3 = class {
  constructor() {
    this._isDisposed = false;
    fr2(this);
  }
  get value() {
    return this._isDisposed ? void 0 : this._value;
  }
  set value(t) {
    this._isDisposed || t === this._value || (this._value?.dispose(), t && vi2(t, this), this._value = t);
  }
  clear() {
    this.value = void 0;
  }
  dispose() {
    this._isDisposed = true, pr(this), this._value?.dispose(), this._value = void 0;
  }
  clearAndLeak() {
    let t = this._value;
    return this._value = void 0, t && vi2(t, null), t;
  }
};
var fe3 = typeof window == "object" ? window : globalThis;
var kt2 = class kt3 {
  constructor(t) {
    this.element = t, this.next = kt3.Undefined, this.prev = kt3.Undefined;
  }
};
kt2.Undefined = new kt2(void 0);
var G2 = kt2;
var Ct2 = class {
  constructor() {
    this._first = G2.Undefined;
    this._last = G2.Undefined;
    this._size = 0;
  }
  get size() {
    return this._size;
  }
  isEmpty() {
    return this._first === G2.Undefined;
  }
  clear() {
    let t = this._first;
    for (; t !== G2.Undefined; ) {
      let e = t.next;
      t.prev = G2.Undefined, t.next = G2.Undefined, t = e;
    }
    this._first = G2.Undefined, this._last = G2.Undefined, this._size = 0;
  }
  unshift(t) {
    return this._insert(t, false);
  }
  push(t) {
    return this._insert(t, true);
  }
  _insert(t, e) {
    let i8 = new G2(t);
    if (this._first === G2.Undefined) this._first = i8, this._last = i8;
    else if (e) {
      let n = this._last;
      this._last = i8, i8.prev = n, n.next = i8;
    } else {
      let n = this._first;
      this._first = i8, i8.next = n, n.prev = i8;
    }
    this._size += 1;
    let r5 = false;
    return () => {
      r5 || (r5 = true, this._remove(i8));
    };
  }
  shift() {
    if (this._first !== G2.Undefined) {
      let t = this._first.element;
      return this._remove(this._first), t;
    }
  }
  pop() {
    if (this._last !== G2.Undefined) {
      let t = this._last.element;
      return this._remove(this._last), t;
    }
  }
  _remove(t) {
    if (t.prev !== G2.Undefined && t.next !== G2.Undefined) {
      let e = t.prev;
      e.next = t.next, t.next.prev = e;
    } else t.prev === G2.Undefined && t.next === G2.Undefined ? (this._first = G2.Undefined, this._last = G2.Undefined) : t.next === G2.Undefined ? (this._last = this._last.prev, this._last.next = G2.Undefined) : t.prev === G2.Undefined && (this._first = this._first.next, this._first.prev = G2.Undefined);
    this._size -= 1;
  }
  *[Symbol.iterator]() {
    let t = this._first;
    for (; t !== G2.Undefined; ) yield t.element, t = t.next;
  }
};
var zl = globalThis.performance && typeof globalThis.performance.now == "function";
var mr2 = class s6 {
  static create(t) {
    return new s6(t);
  }
  constructor(t) {
    this._now = zl && t === false ? Date.now : globalThis.performance.now.bind(globalThis.performance), this._startTime = this._now(), this._stopTime = -1;
  }
  stop() {
    this._stopTime = this._now();
  }
  reset() {
    this._startTime = this._now(), this._stopTime = -1;
  }
  elapsed() {
    return this._stopTime !== -1 ? this._stopTime - this._startTime : this._now() - this._startTime;
  }
};
var Gl = false;
var fo = false;
var $l = false;
var $;
((Qe2) => {
  Qe2.None = () => D2.None;
  function t(y) {
    if ($l) {
      let { onDidAddListener: T } = y, g = gi2.create(), w2 = 0;
      y.onDidAddListener = () => {
        ++w2 === 2 && (console.warn("snapshotted emitter LIKELY used public and SHOULD HAVE BEEN created with DisposableStore. snapshotted here"), g.print()), T?.();
      };
    }
  }
  function e(y, T) {
    return d(y, () => {
    }, 0, void 0, true, void 0, T);
  }
  Qe2.defer = e;
  function i8(y) {
    return (T, g = null, w2) => {
      let E = false, x;
      return x = y((N2) => {
        if (!E) return x ? x.dispose() : E = true, T.call(g, N2);
      }, null, w2), E && x.dispose(), x;
    };
  }
  Qe2.once = i8;
  function r5(y, T, g) {
    return h2((w2, E = null, x) => y((N2) => w2.call(E, T(N2)), null, x), g);
  }
  Qe2.map = r5;
  function n(y, T, g) {
    return h2((w2, E = null, x) => y((N2) => {
      T(N2), w2.call(E, N2);
    }, null, x), g);
  }
  Qe2.forEach = n;
  function o2(y, T, g) {
    return h2((w2, E = null, x) => y((N2) => T(N2) && w2.call(E, N2), null, x), g);
  }
  Qe2.filter = o2;
  function l(y) {
    return y;
  }
  Qe2.signal = l;
  function a(...y) {
    return (T, g = null, w2) => {
      let E = ho(...y.map((x) => x((N2) => T.call(g, N2))));
      return c(E, w2);
    };
  }
  Qe2.any = a;
  function u(y, T, g, w2) {
    let E = g;
    return r5(y, (x) => (E = T(E, x), E), w2);
  }
  Qe2.reduce = u;
  function h2(y, T) {
    let g, w2 = { onWillAddFirstListener() {
      g = y(E.fire, E);
    }, onDidRemoveLastListener() {
      g?.dispose();
    } };
    T || t(w2);
    let E = new v(w2);
    return T?.add(E), E.event;
  }
  function c(y, T) {
    return T instanceof Array ? T.push(y) : T && T.add(y), y;
  }
  function d(y, T, g = 100, w2 = false, E = false, x, N2) {
    let Z3, te4, Oe3, ze2 = 0, le2, et3 = { leakWarningThreshold: x, onWillAddFirstListener() {
      Z3 = y((ht2) => {
        ze2++, te4 = T(te4, ht2), w2 && !Oe3 && (me3.fire(te4), te4 = void 0), le2 = () => {
          let fi2 = te4;
          te4 = void 0, Oe3 = void 0, (!w2 || ze2 > 1) && me3.fire(fi2), ze2 = 0;
        }, typeof g == "number" ? (clearTimeout(Oe3), Oe3 = setTimeout(le2, g)) : Oe3 === void 0 && (Oe3 = 0, queueMicrotask(le2));
      });
    }, onWillRemoveListener() {
      E && ze2 > 0 && le2?.();
    }, onDidRemoveLastListener() {
      le2 = void 0, Z3.dispose();
    } };
    N2 || t(et3);
    let me3 = new v(et3);
    return N2?.add(me3), me3.event;
  }
  Qe2.debounce = d;
  function _3(y, T = 0, g) {
    return Qe2.debounce(y, (w2, E) => w2 ? (w2.push(E), w2) : [E], T, void 0, true, void 0, g);
  }
  Qe2.accumulate = _3;
  function p(y, T = (w2, E) => w2 === E, g) {
    let w2 = true, E;
    return o2(y, (x) => {
      let N2 = w2 || !T(x, E);
      return w2 = false, E = x, N2;
    }, g);
  }
  Qe2.latch = p;
  function m(y, T, g) {
    return [Qe2.filter(y, T, g), Qe2.filter(y, (w2) => !T(w2), g)];
  }
  Qe2.split = m;
  function f(y, T = false, g = [], w2) {
    let E = g.slice(), x = y((te4) => {
      E ? E.push(te4) : Z3.fire(te4);
    });
    w2 && w2.add(x);
    let N2 = () => {
      E?.forEach((te4) => Z3.fire(te4)), E = null;
    }, Z3 = new v({ onWillAddFirstListener() {
      x || (x = y((te4) => Z3.fire(te4)), w2 && w2.add(x));
    }, onDidAddFirstListener() {
      E && (T ? setTimeout(N2) : N2());
    }, onDidRemoveLastListener() {
      x && x.dispose(), x = null;
    } });
    return w2 && w2.add(Z3), Z3.event;
  }
  Qe2.buffer = f;
  function A2(y, T) {
    return (w2, E, x) => {
      let N2 = T(new O3());
      return y(function(Z3) {
        let te4 = N2.evaluate(Z3);
        te4 !== R3 && w2.call(E, te4);
      }, void 0, x);
    };
  }
  Qe2.chain = A2;
  let R3 = /* @__PURE__ */ Symbol("HaltChainable");
  class O3 {
    constructor() {
      this.steps = [];
    }
    map(T) {
      return this.steps.push(T), this;
    }
    forEach(T) {
      return this.steps.push((g) => (T(g), g)), this;
    }
    filter(T) {
      return this.steps.push((g) => T(g) ? g : R3), this;
    }
    reduce(T, g) {
      let w2 = g;
      return this.steps.push((E) => (w2 = T(w2, E), w2)), this;
    }
    latch(T = (g, w2) => g === w2) {
      let g = true, w2;
      return this.steps.push((E) => {
        let x = g || !T(E, w2);
        return g = false, w2 = E, x ? E : R3;
      }), this;
    }
    evaluate(T) {
      for (let g of this.steps) if (T = g(T), T === R3) break;
      return T;
    }
  }
  function I2(y, T, g = (w2) => w2) {
    let w2 = (...Z3) => N2.fire(g(...Z3)), E = () => y.on(T, w2), x = () => y.removeListener(T, w2), N2 = new v({ onWillAddFirstListener: E, onDidRemoveLastListener: x });
    return N2.event;
  }
  Qe2.fromNodeEventEmitter = I2;
  function k2(y, T, g = (w2) => w2) {
    let w2 = (...Z3) => N2.fire(g(...Z3)), E = () => y.addEventListener(T, w2), x = () => y.removeEventListener(T, w2), N2 = new v({ onWillAddFirstListener: E, onDidRemoveLastListener: x });
    return N2.event;
  }
  Qe2.fromDOMEventEmitter = k2;
  function P3(y) {
    return new Promise((T) => i8(y)(T));
  }
  Qe2.toPromise = P3;
  function oe(y) {
    let T = new v();
    return y.then((g) => {
      T.fire(g);
    }, () => {
      T.fire(void 0);
    }).finally(() => {
      T.dispose();
    }), T.event;
  }
  Qe2.fromPromise = oe;
  function Me3(y, T) {
    return y((g) => T.fire(g));
  }
  Qe2.forward = Me3;
  function Pe3(y, T, g) {
    return T(g), y((w2) => T(w2));
  }
  Qe2.runAndSubscribe = Pe3;
  class Ke2 {
    constructor(T, g) {
      this._observable = T;
      this._counter = 0;
      this._hasChanged = false;
      let w2 = { onWillAddFirstListener: () => {
        T.addObserver(this);
      }, onDidRemoveLastListener: () => {
        T.removeObserver(this);
      } };
      g || t(w2), this.emitter = new v(w2), g && g.add(this.emitter);
    }
    beginUpdate(T) {
      this._counter++;
    }
    handlePossibleChange(T) {
    }
    handleChange(T, g) {
      this._hasChanged = true;
    }
    endUpdate(T) {
      this._counter--, this._counter === 0 && (this._observable.reportChanges(), this._hasChanged && (this._hasChanged = false, this.emitter.fire(this._observable.get())));
    }
  }
  function di(y, T) {
    return new Ke2(y, T).emitter.event;
  }
  Qe2.fromObservable = di;
  function V2(y) {
    return (T, g, w2) => {
      let E = 0, x = false, N2 = { beginUpdate() {
        E++;
      }, endUpdate() {
        E--, E === 0 && (y.reportChanges(), x && (x = false, T.call(g)));
      }, handlePossibleChange() {
      }, handleChange() {
        x = true;
      } };
      y.addObserver(N2), y.reportChanges();
      let Z3 = { dispose() {
        y.removeObserver(N2);
      } };
      return w2 instanceof Ee3 ? w2.add(Z3) : Array.isArray(w2) && w2.push(Z3), Z3;
    };
  }
  Qe2.fromObservableLight = V2;
})($ ||= {});
var Mt2 = class Mt3 {
  constructor(t) {
    this.listenerCount = 0;
    this.invocationCount = 0;
    this.elapsedOverall = 0;
    this.durations = [];
    this.name = `${t}_${Mt3._idPool++}`, Mt3.all.add(this);
  }
  start(t) {
    this._stopWatch = new mr2(), this.listenerCount = t;
  }
  stop() {
    if (this._stopWatch) {
      let t = this._stopWatch.elapsed();
      this.durations.push(t), this.elapsedOverall += t, this.invocationCount += 1, this._stopWatch = void 0;
    }
  }
};
Mt2.all = /* @__PURE__ */ new Set(), Mt2._idPool = 0;
var $n2 = Mt2;
var po = -1;
var br2 = class br3 {
  constructor(t, e, i8 = (br3._idPool++).toString(16).padStart(3, "0")) {
    this._errorHandler = t;
    this.threshold = e;
    this.name = i8;
    this._warnCountdown = 0;
  }
  dispose() {
    this._stacks?.clear();
  }
  check(t, e) {
    let i8 = this.threshold;
    if (i8 <= 0 || e < i8) return;
    this._stacks || (this._stacks = /* @__PURE__ */ new Map());
    let r5 = this._stacks.get(t.value) || 0;
    if (this._stacks.set(t.value, r5 + 1), this._warnCountdown -= 1, this._warnCountdown <= 0) {
      this._warnCountdown = i8 * 0.5;
      let [n, o2] = this.getMostFrequentStack(), l = `[${this.name}] potential listener LEAK detected, having ${e} listeners already. MOST frequent listener (${o2}):`;
      console.warn(l), console.warn(n);
      let a = new qn2(l, n);
      this._errorHandler(a);
    }
    return () => {
      let n = this._stacks.get(t.value) || 0;
      this._stacks.set(t.value, n - 1);
    };
  }
  getMostFrequentStack() {
    if (!this._stacks) return;
    let t, e = 0;
    for (let [i8, r5] of this._stacks) (!t || e < r5) && (t = [i8, r5], e = r5);
    return t;
  }
};
br2._idPool = 1;
var Vn2 = br2;
var gi2 = class s7 {
  constructor(t) {
    this.value = t;
  }
  static create() {
    let t = new Error();
    return new s7(t.stack ?? "");
  }
  print() {
    console.warn(this.value.split(`
`).slice(2).join(`
`));
  }
};
var qn2 = class extends Error {
  constructor(t, e) {
    super(t), this.name = "ListenerLeakError", this.stack = e;
  }
};
var Yn2 = class extends Error {
  constructor(t, e) {
    super(t), this.name = "ListenerRefusalError", this.stack = e;
  }
};
var Vl = 0;
var Pt2 = class {
  constructor(t) {
    this.value = t;
    this.id = Vl++;
  }
};
var ql = 2;
var Yl = (s15, t) => {
  if (s15 instanceof Pt2) t(s15);
  else for (let e = 0; e < s15.length; e++) {
    let i8 = s15[e];
    i8 && t(i8);
  }
};
var _r2;
if (Gl) {
  let s15 = [];
  setInterval(() => {
    s15.length !== 0 && (console.warn("[LEAKING LISTENERS] GC'ed these listeners that were NOT yet disposed:"), console.warn(s15.join(`
`)), s15.length = 0);
  }, 3e3), _r2 = new FinalizationRegistry((t) => {
    typeof t == "string" && s15.push(t);
  });
}
var v = class {
  constructor(t) {
    this._size = 0;
    this._options = t, this._leakageMon = po > 0 || this._options?.leakWarningThreshold ? new Vn2(t?.onListenerError ?? Lt2, this._options?.leakWarningThreshold ?? po) : void 0, this._perfMon = this._options?._profName ? new $n2(this._options._profName) : void 0, this._deliveryQueue = this._options?.deliveryQueue;
  }
  dispose() {
    if (!this._disposed) {
      if (this._disposed = true, this._deliveryQueue?.current === this && this._deliveryQueue.reset(), this._listeners) {
        if (fo) {
          let t = this._listeners;
          queueMicrotask(() => {
            Yl(t, (e) => e.stack?.print());
          });
        }
        this._listeners = void 0, this._size = 0;
      }
      this._options?.onDidRemoveLastListener?.(), this._leakageMon?.dispose();
    }
  }
  get event() {
    return this._event ??= (t, e, i8) => {
      if (this._leakageMon && this._size > this._leakageMon.threshold ** 2) {
        let a = `[${this._leakageMon.name}] REFUSES to accept new listeners because it exceeded its threshold by far (${this._size} vs ${this._leakageMon.threshold})`;
        console.warn(a);
        let u = this._leakageMon.getMostFrequentStack() ?? ["UNKNOWN stack", -1], h2 = new Yn2(`${a}. HINT: Stack shows most frequent listener (${u[1]}-times)`, u[0]);
        return (this._options?.onListenerError || Lt2)(h2), D2.None;
      }
      if (this._disposed) return D2.None;
      e && (t = t.bind(e));
      let r5 = new Pt2(t), n, o2;
      this._leakageMon && this._size >= Math.ceil(this._leakageMon.threshold * 0.2) && (r5.stack = gi2.create(), n = this._leakageMon.check(r5.stack, this._size + 1)), fo && (r5.stack = o2 ?? gi2.create()), this._listeners ? this._listeners instanceof Pt2 ? (this._deliveryQueue ??= new jn2(), this._listeners = [this._listeners, r5]) : this._listeners.push(r5) : (this._options?.onWillAddFirstListener?.(this), this._listeners = r5, this._options?.onDidAddFirstListener?.(this)), this._size++;
      let l = C3(() => {
        _r2?.unregister(l), n?.(), this._removeListener(r5);
      });
      if (i8 instanceof Ee3 ? i8.add(l) : Array.isArray(i8) && i8.push(l), _r2) {
        let a = new Error().stack.split(`
`).slice(2, 3).join(`
`).trim(), u = /(file:|vscode-file:\/\/vscode-app)?(\/[^:]*:\d+:\d+)/.exec(a);
        _r2.register(l, u?.[2] ?? a, l);
      }
      return l;
    }, this._event;
  }
  _removeListener(t) {
    if (this._options?.onWillRemoveListener?.(this), !this._listeners) return;
    if (this._size === 1) {
      this._listeners = void 0, this._options?.onDidRemoveLastListener?.(this), this._size = 0;
      return;
    }
    let e = this._listeners, i8 = e.indexOf(t);
    if (i8 === -1) throw console.log("disposed?", this._disposed), console.log("size?", this._size), console.log("arr?", JSON.stringify(this._listeners)), new Error("Attempted to dispose unknown listener");
    this._size--, e[i8] = void 0;
    let r5 = this._deliveryQueue.current === this;
    if (this._size * ql <= e.length) {
      let n = 0;
      for (let o2 = 0; o2 < e.length; o2++) e[o2] ? e[n++] = e[o2] : r5 && (this._deliveryQueue.end--, n < this._deliveryQueue.i && this._deliveryQueue.i--);
      e.length = n;
    }
  }
  _deliver(t, e) {
    if (!t) return;
    let i8 = this._options?.onListenerError || Lt2;
    if (!i8) {
      t.value(e);
      return;
    }
    try {
      t.value(e);
    } catch (r5) {
      i8(r5);
    }
  }
  _deliverQueue(t) {
    let e = t.current._listeners;
    for (; t.i < t.end; ) this._deliver(e[t.i++], t.value);
    t.reset();
  }
  fire(t) {
    if (this._deliveryQueue?.current && (this._deliverQueue(this._deliveryQueue), this._perfMon?.stop()), this._perfMon?.start(this._size), this._listeners) if (this._listeners instanceof Pt2) this._deliver(this._listeners, t);
    else {
      let e = this._deliveryQueue;
      e.enqueue(this, t, this._listeners.length), this._deliverQueue(e);
    }
    this._perfMon?.stop();
  }
  hasListeners() {
    return this._size > 0;
  }
};
var jn2 = class {
  constructor() {
    this.i = -1;
    this.end = 0;
  }
  enqueue(t, e, i8) {
    this.i = 0, this.end = i8, this.current = t, this.value = e;
  }
  reset() {
    this.i = this.end, this.current = void 0, this.value = void 0;
  }
};
var gr2 = class gr3 {
  constructor() {
    this.mapWindowIdToZoomLevel = /* @__PURE__ */ new Map();
    this._onDidChangeZoomLevel = new v();
    this.onDidChangeZoomLevel = this._onDidChangeZoomLevel.event;
    this.mapWindowIdToZoomFactor = /* @__PURE__ */ new Map();
    this._onDidChangeFullscreen = new v();
    this.onDidChangeFullscreen = this._onDidChangeFullscreen.event;
    this.mapWindowIdToFullScreen = /* @__PURE__ */ new Map();
  }
  getZoomLevel(t) {
    return this.mapWindowIdToZoomLevel.get(this.getWindowId(t)) ?? 0;
  }
  setZoomLevel(t, e) {
    if (this.getZoomLevel(e) === t) return;
    let i8 = this.getWindowId(e);
    this.mapWindowIdToZoomLevel.set(i8, t), this._onDidChangeZoomLevel.fire(i8);
  }
  getZoomFactor(t) {
    return this.mapWindowIdToZoomFactor.get(this.getWindowId(t)) ?? 1;
  }
  setZoomFactor(t, e) {
    this.mapWindowIdToZoomFactor.set(this.getWindowId(e), t);
  }
  setFullscreen(t, e) {
    if (this.isFullscreen(e) === t) return;
    let i8 = this.getWindowId(e);
    this.mapWindowIdToFullScreen.set(i8, t), this._onDidChangeFullscreen.fire(i8);
  }
  isFullscreen(t) {
    return !!this.mapWindowIdToFullScreen.get(this.getWindowId(t));
  }
  getWindowId(t) {
    return t.vscodeWindowId;
  }
};
gr2.INSTANCE = new gr2();
var Si2 = gr2;
function Xl(s15, t, e) {
  typeof t == "string" && (t = s15.matchMedia(t)), t.addEventListener("change", e);
}
var Eu2 = Si2.INSTANCE.onDidChangeZoomLevel;
function mo(s15) {
  return Si2.INSTANCE.getZoomFactor(s15);
}
var Tu2 = Si2.INSTANCE.onDidChangeFullscreen;
var Ot2 = typeof navigator == "object" ? navigator.userAgent : "";
var Ei2 = Ot2.indexOf("Firefox") >= 0;
var Bt3 = Ot2.indexOf("AppleWebKit") >= 0;
var Ti2 = Ot2.indexOf("Chrome") >= 0;
var Sr2 = !Ti2 && Ot2.indexOf("Safari") >= 0;
var Iu = Ot2.indexOf("Electron/") >= 0;
var yu2 = Ot2.indexOf("Android") >= 0;
var vr2 = false;
if (typeof fe3.matchMedia == "function") {
  let s15 = fe3.matchMedia("(display-mode: standalone) or (display-mode: window-controls-overlay)"), t = fe3.matchMedia("(display-mode: fullscreen)");
  vr2 = s15.matches, Xl(fe3, s15, ({ matches: e }) => {
    vr2 && t.matches || (vr2 = e);
  });
}
function _o() {
  return vr2;
}
var Nt2 = "en";
var yr = false;
var xr2 = false;
var Ii2 = false;
var Zl = false;
var vo = false;
var go = false;
var Jl = false;
var Ql = false;
var ea = false;
var ta = false;
var Tr2;
var Ir = Nt2;
var bo = Nt2;
var ia;
var $e;
var Ve3 = globalThis;
var xe2;
typeof Ve3.vscode < "u" && typeof Ve3.vscode.process < "u" ? xe2 = Ve3.vscode.process : typeof process < "u" && typeof process?.versions?.node == "string" && (xe2 = process);
var So = typeof xe2?.versions?.electron == "string";
var ra = So && xe2?.type === "renderer";
if (typeof xe2 == "object") {
  yr = xe2.platform === "win32", xr2 = xe2.platform === "darwin", Ii2 = xe2.platform === "linux", Zl = Ii2 && !!xe2.env.SNAP && !!xe2.env.SNAP_REVISION, Jl = So, ea = !!xe2.env.CI || !!xe2.env.BUILD_ARTIFACTSTAGINGDIRECTORY, Tr2 = Nt2, Ir = Nt2;
  let s15 = xe2.env.VSCODE_NLS_CONFIG;
  if (s15) try {
    let t = JSON.parse(s15);
    Tr2 = t.userLocale, bo = t.osLocale, Ir = t.resolvedLanguage || Nt2, ia = t.languagePack?.translationsConfigFile;
  } catch {
  }
  vo = true;
} else typeof navigator == "object" && !ra ? ($e = navigator.userAgent, yr = $e.indexOf("Windows") >= 0, xr2 = $e.indexOf("Macintosh") >= 0, Ql = ($e.indexOf("Macintosh") >= 0 || $e.indexOf("iPad") >= 0 || $e.indexOf("iPhone") >= 0) && !!navigator.maxTouchPoints && navigator.maxTouchPoints > 0, Ii2 = $e.indexOf("Linux") >= 0, ta = $e?.indexOf("Mobi") >= 0, go = true, Ir = globalThis._VSCODE_NLS_LANGUAGE || Nt2, Tr2 = navigator.language.toLowerCase(), bo = Tr2) : console.error("Unable to resolve platform.");
var Xn2 = 0;
xr2 ? Xn2 = 1 : yr ? Xn2 = 3 : Ii2 && (Xn2 = 2);
var wr2 = yr;
var Te3 = xr2;
var Zn2 = Ii2;
var Dr2 = vo;
var na = go && typeof Ve3.importScripts == "function";
var xu2 = na ? Ve3.origin : void 0;
var Fe2 = $e;
var st3 = Ir;
var sa;
((i8) => {
  function s15() {
    return st3;
  }
  i8.value = s15;
  function t() {
    return st3.length === 2 ? st3 === "en" : st3.length >= 3 ? st3[0] === "e" && st3[1] === "n" && st3[2] === "-" : false;
  }
  i8.isDefaultVariant = t;
  function e() {
    return st3 === "en";
  }
  i8.isDefault = e;
})(sa ||= {});
var oa = typeof Ve3.postMessage == "function" && !Ve3.importScripts;
var Eo = (() => {
  if (oa) {
    let s15 = [];
    Ve3.addEventListener("message", (e) => {
      if (e.data && e.data.vscodeScheduleAsyncWork) for (let i8 = 0, r5 = s15.length; i8 < r5; i8++) {
        let n = s15[i8];
        if (n.id === e.data.vscodeScheduleAsyncWork) {
          s15.splice(i8, 1), n.callback();
          return;
        }
      }
    });
    let t = 0;
    return (e) => {
      let i8 = ++t;
      s15.push({ id: i8, callback: e }), Ve3.postMessage({ vscodeScheduleAsyncWork: i8 }, "*");
    };
  }
  return (s15) => setTimeout(s15);
})();
var la = !!(Fe2 && Fe2.indexOf("Chrome") >= 0);
var wu = !!(Fe2 && Fe2.indexOf("Firefox") >= 0);
var Du = !!(!la && Fe2 && Fe2.indexOf("Safari") >= 0);
var Ru = !!(Fe2 && Fe2.indexOf("Edg/") >= 0);
var Lu = !!(Fe2 && Fe2.indexOf("Android") >= 0);
var ot3 = typeof navigator == "object" ? navigator : {};
var aa = { clipboard: { writeText: Dr2 || document.queryCommandSupported && document.queryCommandSupported("copy") || !!(ot3 && ot3.clipboard && ot3.clipboard.writeText), readText: Dr2 || !!(ot3 && ot3.clipboard && ot3.clipboard.readText) }, keyboard: Dr2 || _o() ? 0 : ot3.keyboard || Sr2 ? 1 : 2, touch: "ontouchstart" in fe3 || ot3.maxTouchPoints > 0, pointerEvents: fe3.PointerEvent && ("ontouchstart" in fe3 || navigator.maxTouchPoints > 0) };
var yi2 = class {
  constructor() {
    this._keyCodeToStr = [], this._strToKeyCode = /* @__PURE__ */ Object.create(null);
  }
  define(t, e) {
    this._keyCodeToStr[t] = e, this._strToKeyCode[e.toLowerCase()] = t;
  }
  keyCodeToStr(t) {
    return this._keyCodeToStr[t];
  }
  strToKeyCode(t) {
    return this._strToKeyCode[t.toLowerCase()] || 0;
  }
};
var Jn2 = new yi2();
var To = new yi2();
var Io = new yi2();
var yo = new Array(230);
var Qn2;
((o2) => {
  function s15(l) {
    return Jn2.keyCodeToStr(l);
  }
  o2.toString = s15;
  function t(l) {
    return Jn2.strToKeyCode(l);
  }
  o2.fromString = t;
  function e(l) {
    return To.keyCodeToStr(l);
  }
  o2.toUserSettingsUS = e;
  function i8(l) {
    return Io.keyCodeToStr(l);
  }
  o2.toUserSettingsGeneral = i8;
  function r5(l) {
    return To.strToKeyCode(l) || Io.strToKeyCode(l);
  }
  o2.fromUserSettings = r5;
  function n(l) {
    if (l >= 98 && l <= 113) return null;
    switch (l) {
      case 16:
        return "Up";
      case 18:
        return "Down";
      case 15:
        return "Left";
      case 17:
        return "Right";
    }
    return Jn2.keyCodeToStr(l);
  }
  o2.toElectronAccelerator = n;
})(Qn2 ||= {});
var Rr2 = class s8 {
  constructor(t, e, i8, r5, n) {
    this.ctrlKey = t;
    this.shiftKey = e;
    this.altKey = i8;
    this.metaKey = r5;
    this.keyCode = n;
  }
  equals(t) {
    return t instanceof s8 && this.ctrlKey === t.ctrlKey && this.shiftKey === t.shiftKey && this.altKey === t.altKey && this.metaKey === t.metaKey && this.keyCode === t.keyCode;
  }
  getHashCode() {
    let t = this.ctrlKey ? "1" : "0", e = this.shiftKey ? "1" : "0", i8 = this.altKey ? "1" : "0", r5 = this.metaKey ? "1" : "0";
    return `K${t}${e}${i8}${r5}${this.keyCode}`;
  }
  isModifierKey() {
    return this.keyCode === 0 || this.keyCode === 5 || this.keyCode === 57 || this.keyCode === 6 || this.keyCode === 4;
  }
  toKeybinding() {
    return new es([this]);
  }
  isDuplicateModifierCase() {
    return this.ctrlKey && this.keyCode === 5 || this.shiftKey && this.keyCode === 4 || this.altKey && this.keyCode === 6 || this.metaKey && this.keyCode === 57;
  }
};
var es = class {
  constructor(t) {
    if (t.length === 0) throw eo("chords");
    this.chords = t;
  }
  getHashCode() {
    let t = "";
    for (let e = 0, i8 = this.chords.length; e < i8; e++) e !== 0 && (t += ";"), t += this.chords[e].getHashCode();
    return t;
  }
  equals(t) {
    if (t === null || this.chords.length !== t.chords.length) return false;
    for (let e = 0; e < this.chords.length; e++) if (!this.chords[e].equals(t.chords[e])) return false;
    return true;
  }
};
function ca(s15) {
  if (s15.charCode) {
    let e = String.fromCharCode(s15.charCode).toUpperCase();
    return Qn2.fromString(e);
  }
  let t = s15.keyCode;
  if (t === 3) return 7;
  if (Ei2) switch (t) {
    case 59:
      return 85;
    case 60:
      if (Zn2) return 97;
      break;
    case 61:
      return 86;
    case 107:
      return 109;
    case 109:
      return 111;
    case 173:
      return 88;
    case 224:
      if (Te3) return 57;
      break;
  }
  else if (Bt3) {
    if (Te3 && t === 93) return 57;
    if (!Te3 && t === 92) return 57;
  }
  return yo[t] || 0;
}
var ua = Te3 ? 256 : 2048;
var ha = 512;
var da = 1024;
var fa = Te3 ? 2048 : 256;
var ft = class {
  constructor(t) {
    this._standardKeyboardEventBrand = true;
    let e = t;
    this.browserEvent = e, this.target = e.target, this.ctrlKey = e.ctrlKey, this.shiftKey = e.shiftKey, this.altKey = e.altKey, this.metaKey = e.metaKey, this.altGraphKey = e.getModifierState?.("AltGraph"), this.keyCode = ca(e), this.code = e.code, this.ctrlKey = this.ctrlKey || this.keyCode === 5, this.altKey = this.altKey || this.keyCode === 6, this.shiftKey = this.shiftKey || this.keyCode === 4, this.metaKey = this.metaKey || this.keyCode === 57, this._asKeybinding = this._computeKeybinding(), this._asKeyCodeChord = this._computeKeyCodeChord();
  }
  preventDefault() {
    this.browserEvent && this.browserEvent.preventDefault && this.browserEvent.preventDefault();
  }
  stopPropagation() {
    this.browserEvent && this.browserEvent.stopPropagation && this.browserEvent.stopPropagation();
  }
  toKeyCodeChord() {
    return this._asKeyCodeChord;
  }
  equals(t) {
    return this._asKeybinding === t;
  }
  _computeKeybinding() {
    let t = 0;
    this.keyCode !== 5 && this.keyCode !== 4 && this.keyCode !== 6 && this.keyCode !== 57 && (t = this.keyCode);
    let e = 0;
    return this.ctrlKey && (e |= ua), this.altKey && (e |= ha), this.shiftKey && (e |= da), this.metaKey && (e |= fa), e |= t, e;
  }
  _computeKeyCodeChord() {
    let t = 0;
    return this.keyCode !== 5 && this.keyCode !== 4 && this.keyCode !== 6 && this.keyCode !== 57 && (t = this.keyCode), new Rr2(this.ctrlKey, this.shiftKey, this.altKey, this.metaKey, t);
  }
};
var wo = /* @__PURE__ */ new WeakMap();
function pa(s15) {
  if (!s15.parent || s15.parent === s15) return null;
  try {
    let t = s15.location, e = s15.parent.location;
    if (t.origin !== "null" && e.origin !== "null" && t.origin !== e.origin) return null;
  } catch {
    return null;
  }
  return s15.parent;
}
var Lr2 = class {
  static getSameOriginWindowChain(t) {
    let e = wo.get(t);
    if (!e) {
      e = [], wo.set(t, e);
      let i8 = t, r5;
      do
        r5 = pa(i8), r5 ? e.push({ window: new WeakRef(i8), iframeElement: i8.frameElement || null }) : e.push({ window: new WeakRef(i8), iframeElement: null }), i8 = r5;
      while (i8);
    }
    return e.slice(0);
  }
  static getPositionOfChildWindowRelativeToAncestorWindow(t, e) {
    if (!e || t === e) return { top: 0, left: 0 };
    let i8 = 0, r5 = 0, n = this.getSameOriginWindowChain(t);
    for (let o2 of n) {
      let l = o2.window.deref();
      if (i8 += l?.scrollY ?? 0, r5 += l?.scrollX ?? 0, l === e || !o2.iframeElement) break;
      let a = o2.iframeElement.getBoundingClientRect();
      i8 += a.top, r5 += a.left;
    }
    return { top: i8, left: r5 };
  }
};
var qe3 = class {
  constructor(t, e) {
    this.timestamp = Date.now(), this.browserEvent = e, this.leftButton = e.button === 0, this.middleButton = e.button === 1, this.rightButton = e.button === 2, this.buttons = e.buttons, this.target = e.target, this.detail = e.detail || 1, e.type === "dblclick" && (this.detail = 2), this.ctrlKey = e.ctrlKey, this.shiftKey = e.shiftKey, this.altKey = e.altKey, this.metaKey = e.metaKey, typeof e.pageX == "number" ? (this.posx = e.pageX, this.posy = e.pageY) : (this.posx = e.clientX + this.target.ownerDocument.body.scrollLeft + this.target.ownerDocument.documentElement.scrollLeft, this.posy = e.clientY + this.target.ownerDocument.body.scrollTop + this.target.ownerDocument.documentElement.scrollTop);
    let i8 = Lr2.getPositionOfChildWindowRelativeToAncestorWindow(t, e.view);
    this.posx -= i8.left, this.posy -= i8.top;
  }
  preventDefault() {
    this.browserEvent.preventDefault();
  }
  stopPropagation() {
    this.browserEvent.stopPropagation();
  }
};
var xi2 = class {
  constructor(t, e = 0, i8 = 0) {
    this.browserEvent = t || null, this.target = t ? t.target || t.targetNode || t.srcElement : null, this.deltaY = i8, this.deltaX = e;
    let r5 = false;
    if (Ti2) {
      let n = navigator.userAgent.match(/Chrome\/(\d+)/);
      r5 = (n ? parseInt(n[1]) : 123) <= 122;
    }
    if (t) {
      let n = t, o2 = t, l = t.view?.devicePixelRatio || 1;
      if (typeof n.wheelDeltaY < "u") r5 ? this.deltaY = n.wheelDeltaY / (120 * l) : this.deltaY = n.wheelDeltaY / 120;
      else if (typeof o2.VERTICAL_AXIS < "u" && o2.axis === o2.VERTICAL_AXIS) this.deltaY = -o2.detail / 3;
      else if (t.type === "wheel") {
        let a = t;
        a.deltaMode === a.DOM_DELTA_LINE ? Ei2 && !Te3 ? this.deltaY = -t.deltaY / 3 : this.deltaY = -t.deltaY : this.deltaY = -t.deltaY / 40;
      }
      if (typeof n.wheelDeltaX < "u") Sr2 && wr2 ? this.deltaX = -(n.wheelDeltaX / 120) : r5 ? this.deltaX = n.wheelDeltaX / (120 * l) : this.deltaX = n.wheelDeltaX / 120;
      else if (typeof o2.HORIZONTAL_AXIS < "u" && o2.axis === o2.HORIZONTAL_AXIS) this.deltaX = -t.detail / 3;
      else if (t.type === "wheel") {
        let a = t;
        a.deltaMode === a.DOM_DELTA_LINE ? Ei2 && !Te3 ? this.deltaX = -t.deltaX / 3 : this.deltaX = -t.deltaX : this.deltaX = -t.deltaX / 40;
      }
      this.deltaY === 0 && this.deltaX === 0 && t.wheelDelta && (r5 ? this.deltaY = t.wheelDelta / (120 * l) : this.deltaY = t.wheelDelta / 120);
    }
  }
  preventDefault() {
    this.browserEvent?.preventDefault();
  }
  stopPropagation() {
    this.browserEvent?.stopPropagation();
  }
};
var Do = Object.freeze(function(s15, t) {
  let e = setTimeout(s15.bind(t), 0);
  return { dispose() {
    clearTimeout(e);
  } };
});
var ma;
((i8) => {
  function s15(r5) {
    return r5 === i8.None || r5 === i8.Cancelled || r5 instanceof ts ? true : !r5 || typeof r5 != "object" ? false : typeof r5.isCancellationRequested == "boolean" && typeof r5.onCancellationRequested == "function";
  }
  i8.isCancellationToken = s15, i8.None = Object.freeze({ isCancellationRequested: false, onCancellationRequested: $.None }), i8.Cancelled = Object.freeze({ isCancellationRequested: true, onCancellationRequested: Do });
})(ma ||= {});
var ts = class {
  constructor() {
    this._isCancelled = false;
    this._emitter = null;
  }
  cancel() {
    this._isCancelled || (this._isCancelled = true, this._emitter && (this._emitter.fire(void 0), this.dispose()));
  }
  get isCancellationRequested() {
    return this._isCancelled;
  }
  get onCancellationRequested() {
    return this._isCancelled ? Do : (this._emitter || (this._emitter = new v()), this._emitter.event);
  }
  dispose() {
    this._emitter && (this._emitter.dispose(), this._emitter = null);
  }
};
var Ye3 = class {
  constructor(t, e) {
    this._isDisposed = false;
    this._token = -1, typeof t == "function" && typeof e == "number" && this.setIfNotSet(t, e);
  }
  dispose() {
    this.cancel(), this._isDisposed = true;
  }
  cancel() {
    this._token !== -1 && (clearTimeout(this._token), this._token = -1);
  }
  cancelAndSet(t, e) {
    if (this._isDisposed) throw new Rt2("Calling 'cancelAndSet' on a disposed TimeoutTimer");
    this.cancel(), this._token = setTimeout(() => {
      this._token = -1, t();
    }, e);
  }
  setIfNotSet(t, e) {
    if (this._isDisposed) throw new Rt2("Calling 'setIfNotSet' on a disposed TimeoutTimer");
    this._token === -1 && (this._token = setTimeout(() => {
      this._token = -1, t();
    }, e));
  }
};
var kr2 = class {
  constructor() {
    this.disposable = void 0;
    this.isDisposed = false;
  }
  cancel() {
    this.disposable?.dispose(), this.disposable = void 0;
  }
  cancelAndSet(t, e, i8 = globalThis) {
    if (this.isDisposed) throw new Rt2("Calling 'cancelAndSet' on a disposed IntervalTimer");
    this.cancel();
    let r5 = i8.setInterval(() => {
      t();
    }, e);
    this.disposable = C3(() => {
      i8.clearInterval(r5), this.disposable = void 0;
    });
  }
  dispose() {
    this.cancel(), this.isDisposed = true;
  }
};
var ba;
var Ar2;
(function() {
  typeof globalThis.requestIdleCallback != "function" || typeof globalThis.cancelIdleCallback != "function" ? Ar2 = (s15, t) => {
    Eo(() => {
      if (e) return;
      let i8 = Date.now() + 15;
      t(Object.freeze({ didTimeout: true, timeRemaining() {
        return Math.max(0, i8 - Date.now());
      } }));
    });
    let e = false;
    return { dispose() {
      e || (e = true);
    } };
  } : Ar2 = (s15, t, e) => {
    let i8 = s15.requestIdleCallback(t, typeof e == "number" ? { timeout: e } : void 0), r5 = false;
    return { dispose() {
      r5 || (r5 = true, s15.cancelIdleCallback(i8));
    } };
  }, ba = (s15) => Ar2(globalThis, s15);
})();
var va;
((e) => {
  async function s15(i8) {
    let r5, n = await Promise.all(i8.map((o2) => o2.then((l) => l, (l) => {
      r5 || (r5 = l);
    })));
    if (typeof r5 < "u") throw r5;
    return n;
  }
  e.settled = s15;
  function t(i8) {
    return new Promise(async (r5, n) => {
      try {
        await i8(r5, n);
      } catch (o2) {
        n(o2);
      }
    });
  }
  e.withAsyncBody = t;
})(va ||= {});
var _e2 = class _e3 {
  static fromArray(t) {
    return new _e3((e) => {
      e.emitMany(t);
    });
  }
  static fromPromise(t) {
    return new _e3(async (e) => {
      e.emitMany(await t);
    });
  }
  static fromPromises(t) {
    return new _e3(async (e) => {
      await Promise.all(t.map(async (i8) => e.emitOne(await i8)));
    });
  }
  static merge(t) {
    return new _e3(async (e) => {
      await Promise.all(t.map(async (i8) => {
        for await (let r5 of i8) e.emitOne(r5);
      }));
    });
  }
  constructor(t, e) {
    this._state = 0, this._results = [], this._error = null, this._onReturn = e, this._onStateChanged = new v(), queueMicrotask(async () => {
      let i8 = { emitOne: (r5) => this.emitOne(r5), emitMany: (r5) => this.emitMany(r5), reject: (r5) => this.reject(r5) };
      try {
        await Promise.resolve(t(i8)), this.resolve();
      } catch (r5) {
        this.reject(r5);
      } finally {
        i8.emitOne = void 0, i8.emitMany = void 0, i8.reject = void 0;
      }
    });
  }
  [Symbol.asyncIterator]() {
    let t = 0;
    return { next: async () => {
      do {
        if (this._state === 2) throw this._error;
        if (t < this._results.length) return { done: false, value: this._results[t++] };
        if (this._state === 1) return { done: true, value: void 0 };
        await $.toPromise(this._onStateChanged.event);
      } while (true);
    }, return: async () => (this._onReturn?.(), { done: true, value: void 0 }) };
  }
  static map(t, e) {
    return new _e3(async (i8) => {
      for await (let r5 of t) i8.emitOne(e(r5));
    });
  }
  map(t) {
    return _e3.map(this, t);
  }
  static filter(t, e) {
    return new _e3(async (i8) => {
      for await (let r5 of t) e(r5) && i8.emitOne(r5);
    });
  }
  filter(t) {
    return _e3.filter(this, t);
  }
  static coalesce(t) {
    return _e3.filter(t, (e) => !!e);
  }
  coalesce() {
    return _e3.coalesce(this);
  }
  static async toPromise(t) {
    let e = [];
    for await (let i8 of t) e.push(i8);
    return e;
  }
  toPromise() {
    return _e3.toPromise(this);
  }
  emitOne(t) {
    this._state === 0 && (this._results.push(t), this._onStateChanged.fire());
  }
  emitMany(t) {
    this._state === 0 && (this._results = this._results.concat(t), this._onStateChanged.fire());
  }
  resolve() {
    this._state === 0 && (this._state = 1, this._onStateChanged.fire());
  }
  reject(t) {
    this._state === 0 && (this._state = 2, this._error = t, this._onStateChanged.fire());
  }
};
_e2.EMPTY = _e2.fromArray([]);
function Lo(s15) {
  return 55296 <= s15 && s15 <= 56319;
}
function is2(s15) {
  return 56320 <= s15 && s15 <= 57343;
}
function Ao(s15, t) {
  return (s15 - 55296 << 10) + (t - 56320) + 65536;
}
function Mo(s15) {
  return ns2(s15, 0);
}
function ns2(s15, t) {
  switch (typeof s15) {
    case "object":
      return s15 === null ? je3(349, t) : Array.isArray(s15) ? Ea(s15, t) : Ta(s15, t);
    case "string":
      return Po(s15, t);
    case "boolean":
      return Sa(s15, t);
    case "number":
      return je3(s15, t);
    case "undefined":
      return je3(937, t);
    default:
      return je3(617, t);
  }
}
function je3(s15, t) {
  return (t << 5) - t + s15 | 0;
}
function Sa(s15, t) {
  return je3(s15 ? 433 : 863, t);
}
function Po(s15, t) {
  t = je3(149417, t);
  for (let e = 0, i8 = s15.length; e < i8; e++) t = je3(s15.charCodeAt(e), t);
  return t;
}
function Ea(s15, t) {
  return t = je3(104579, t), s15.reduce((e, i8) => ns2(i8, e), t);
}
function Ta(s15, t) {
  return t = je3(181387, t), Object.keys(s15).sort().reduce((e, i8) => (e = Po(i8, e), ns2(s15[i8], e)), t);
}
function rs2(s15, t, e = 32) {
  let i8 = e - t, r5 = ~((1 << i8) - 1);
  return (s15 << t | (r5 & s15) >>> i8) >>> 0;
}
function ko(s15, t = 0, e = s15.byteLength, i8 = 0) {
  for (let r5 = 0; r5 < e; r5++) s15[t + r5] = i8;
}
function Ia(s15, t, e = "0") {
  for (; s15.length < t; ) s15 = e + s15;
  return s15;
}
function wi2(s15, t = 32) {
  return s15 instanceof ArrayBuffer ? Array.from(new Uint8Array(s15)).map((e) => e.toString(16).padStart(2, "0")).join("") : Ia((s15 >>> 0).toString(16), t / 4);
}
var Cr2 = class Cr3 {
  constructor() {
    this._h0 = 1732584193;
    this._h1 = 4023233417;
    this._h2 = 2562383102;
    this._h3 = 271733878;
    this._h4 = 3285377520;
    this._buff = new Uint8Array(67), this._buffDV = new DataView(this._buff.buffer), this._buffLen = 0, this._totalLen = 0, this._leftoverHighSurrogate = 0, this._finished = false;
  }
  update(t) {
    let e = t.length;
    if (e === 0) return;
    let i8 = this._buff, r5 = this._buffLen, n = this._leftoverHighSurrogate, o2, l;
    for (n !== 0 ? (o2 = n, l = -1, n = 0) : (o2 = t.charCodeAt(0), l = 0); ; ) {
      let a = o2;
      if (Lo(o2)) if (l + 1 < e) {
        let u = t.charCodeAt(l + 1);
        is2(u) ? (l++, a = Ao(o2, u)) : a = 65533;
      } else {
        n = o2;
        break;
      }
      else is2(o2) && (a = 65533);
      if (r5 = this._push(i8, r5, a), l++, l < e) o2 = t.charCodeAt(l);
      else break;
    }
    this._buffLen = r5, this._leftoverHighSurrogate = n;
  }
  _push(t, e, i8) {
    return i8 < 128 ? t[e++] = i8 : i8 < 2048 ? (t[e++] = 192 | (i8 & 1984) >>> 6, t[e++] = 128 | (i8 & 63) >>> 0) : i8 < 65536 ? (t[e++] = 224 | (i8 & 61440) >>> 12, t[e++] = 128 | (i8 & 4032) >>> 6, t[e++] = 128 | (i8 & 63) >>> 0) : (t[e++] = 240 | (i8 & 1835008) >>> 18, t[e++] = 128 | (i8 & 258048) >>> 12, t[e++] = 128 | (i8 & 4032) >>> 6, t[e++] = 128 | (i8 & 63) >>> 0), e >= 64 && (this._step(), e -= 64, this._totalLen += 64, t[0] = t[64], t[1] = t[65], t[2] = t[66]), e;
  }
  digest() {
    return this._finished || (this._finished = true, this._leftoverHighSurrogate && (this._leftoverHighSurrogate = 0, this._buffLen = this._push(this._buff, this._buffLen, 65533)), this._totalLen += this._buffLen, this._wrapUp()), wi2(this._h0) + wi2(this._h1) + wi2(this._h2) + wi2(this._h3) + wi2(this._h4);
  }
  _wrapUp() {
    this._buff[this._buffLen++] = 128, ko(this._buff, this._buffLen), this._buffLen > 56 && (this._step(), ko(this._buff));
    let t = 8 * this._totalLen;
    this._buffDV.setUint32(56, Math.floor(t / 4294967296), false), this._buffDV.setUint32(60, t % 4294967296, false), this._step();
  }
  _step() {
    let t = Cr3._bigBlock32, e = this._buffDV;
    for (let c = 0; c < 64; c += 4) t.setUint32(c, e.getUint32(c, false), false);
    for (let c = 64; c < 320; c += 4) t.setUint32(c, rs2(t.getUint32(c - 12, false) ^ t.getUint32(c - 32, false) ^ t.getUint32(c - 56, false) ^ t.getUint32(c - 64, false), 1), false);
    let i8 = this._h0, r5 = this._h1, n = this._h2, o2 = this._h3, l = this._h4, a, u, h2;
    for (let c = 0; c < 80; c++) c < 20 ? (a = r5 & n | ~r5 & o2, u = 1518500249) : c < 40 ? (a = r5 ^ n ^ o2, u = 1859775393) : c < 60 ? (a = r5 & n | r5 & o2 | n & o2, u = 2400959708) : (a = r5 ^ n ^ o2, u = 3395469782), h2 = rs2(i8, 5) + a + l + u + t.getUint32(c * 4, false) & 4294967295, l = o2, o2 = n, n = rs2(r5, 30), r5 = i8, i8 = h2;
    this._h0 = this._h0 + i8 & 4294967295, this._h1 = this._h1 + r5 & 4294967295, this._h2 = this._h2 + n & 4294967295, this._h3 = this._h3 + o2 & 4294967295, this._h4 = this._h4 + l & 4294967295;
  }
};
Cr2._bigBlock32 = new DataView(new ArrayBuffer(320));
var { registerWindow: Bh, getWindow: be3, getDocument: Nh, getWindows: Fh, getWindowsCount: Hh, getWindowId: Oo, getWindowById: Wh, hasWindow: Uh, onDidRegisterWindow: No, onWillUnregisterWindow: Kh, onDidUnregisterWindow: zh } = (function() {
  let s15 = /* @__PURE__ */ new Map();
  fe3;
  let t = { window: fe3, disposables: new Ee3() };
  s15.set(fe3.vscodeWindowId, t);
  let e = new v(), i8 = new v(), r5 = new v();
  function n(o2, l) {
    return (typeof o2 == "number" ? s15.get(o2) : void 0) ?? (l ? t : void 0);
  }
  return { onDidRegisterWindow: e.event, onWillUnregisterWindow: r5.event, onDidUnregisterWindow: i8.event, registerWindow(o2) {
    if (s15.has(o2.vscodeWindowId)) return D2.None;
    let l = new Ee3(), a = { window: o2, disposables: l.add(new Ee3()) };
    return s15.set(o2.vscodeWindowId, a), l.add(C3(() => {
      s15.delete(o2.vscodeWindowId), i8.fire(o2);
    })), l.add(L(o2, Y3.BEFORE_UNLOAD, () => {
      r5.fire(o2);
    })), e.fire(a), l;
  }, getWindows() {
    return s15.values();
  }, getWindowsCount() {
    return s15.size;
  }, getWindowId(o2) {
    return o2.vscodeWindowId;
  }, hasWindow(o2) {
    return s15.has(o2);
  }, getWindowById: n, getWindow(o2) {
    let l = o2;
    if (l?.ownerDocument?.defaultView) return l.ownerDocument.defaultView.window;
    let a = o2;
    return a?.view ? a.view.window : fe3;
  }, getDocument(o2) {
    return be3(o2).document;
  } };
})();
var ss2 = class {
  constructor(t, e, i8, r5) {
    this._node = t, this._type = e, this._handler = i8, this._options = r5 || false, this._node.addEventListener(this._type, this._handler, this._options);
  }
  dispose() {
    this._handler && (this._node.removeEventListener(this._type, this._handler, this._options), this._node = null, this._handler = null);
  }
};
function L(s15, t, e, i8) {
  return new ss2(s15, t, e, i8);
}
function ya(s15, t) {
  return function(e) {
    return t(new qe3(s15, e));
  };
}
function xa(s15) {
  return function(t) {
    return s15(new ft(t));
  };
}
var os2 = function(t, e, i8, r5) {
  let n = i8;
  return e === "click" || e === "mousedown" || e === "contextmenu" ? n = ya(be3(t), i8) : (e === "keydown" || e === "keypress" || e === "keyup") && (n = xa(i8)), L(t, e, n, r5);
};
var wa;
var mt;
var Mr2 = class extends kr2 {
  constructor(t) {
    super(), this.defaultTarget = t && be3(t);
  }
  cancelAndSet(t, e, i8) {
    return super.cancelAndSet(t, e, i8 ?? this.defaultTarget);
  }
};
var Di2 = class {
  constructor(t, e = 0) {
    this._runner = t, this.priority = e, this._canceled = false;
  }
  dispose() {
    this._canceled = true;
  }
  execute() {
    if (!this._canceled) try {
      this._runner();
    } catch (t) {
      Lt2(t);
    }
  }
  static sort(t, e) {
    return e.priority - t.priority;
  }
};
(function() {
  let s15 = /* @__PURE__ */ new Map(), t = /* @__PURE__ */ new Map(), e = /* @__PURE__ */ new Map(), i8 = /* @__PURE__ */ new Map(), r5 = (n) => {
    e.set(n, false);
    let o2 = s15.get(n) ?? [];
    for (t.set(n, o2), s15.set(n, []), i8.set(n, true); o2.length > 0; ) o2.sort(Di2.sort), o2.shift().execute();
    i8.set(n, false);
  };
  mt = (n, o2, l = 0) => {
    let a = Oo(n), u = new Di2(o2, l), h2 = s15.get(a);
    return h2 || (h2 = [], s15.set(a, h2)), h2.push(u), e.get(a) || (e.set(a, true), n.requestAnimationFrame(() => r5(a))), u;
  }, wa = (n, o2, l) => {
    let a = Oo(n);
    if (i8.get(a)) {
      let u = new Di2(o2, l), h2 = t.get(a);
      return h2 || (h2 = [], t.set(a, h2)), h2.push(u), u;
    } else return mt(n, o2, l);
  };
})();
var pt2 = class pt3 {
  constructor(t, e) {
    this.width = t;
    this.height = e;
  }
  with(t = this.width, e = this.height) {
    return t !== this.width || e !== this.height ? new pt3(t, e) : this;
  }
  static is(t) {
    return typeof t == "object" && typeof t.height == "number" && typeof t.width == "number";
  }
  static lift(t) {
    return t instanceof pt3 ? t : new pt3(t.width, t.height);
  }
  static equals(t, e) {
    return t === e ? true : !t || !e ? false : t.width === e.width && t.height === e.height;
  }
};
pt2.None = new pt2(0, 0);
function Fo(s15) {
  let t = s15.getBoundingClientRect(), e = be3(s15);
  return { left: t.left + e.scrollX, top: t.top + e.scrollY, width: t.width, height: t.height };
}
var Gh = new class {
  constructor() {
    this.mutationObservers = /* @__PURE__ */ new Map();
  }
  observe(s15, t, e) {
    let i8 = this.mutationObservers.get(s15);
    i8 || (i8 = /* @__PURE__ */ new Map(), this.mutationObservers.set(s15, i8));
    let r5 = Mo(e), n = i8.get(r5);
    if (n) n.users += 1;
    else {
      let o2 = new v(), l = new MutationObserver((u) => o2.fire(u));
      l.observe(s15, e);
      let a = n = { users: 1, observer: l, onDidMutate: o2.event };
      t.add(C3(() => {
        a.users -= 1, a.users === 0 && (o2.dispose(), l.disconnect(), i8?.delete(r5), i8?.size === 0 && this.mutationObservers.delete(s15));
      })), i8.set(r5, n);
    }
    return n.onDidMutate;
  }
}();
var Y3 = { CLICK: "click", AUXCLICK: "auxclick", DBLCLICK: "dblclick", MOUSE_UP: "mouseup", MOUSE_DOWN: "mousedown", MOUSE_OVER: "mouseover", MOUSE_MOVE: "mousemove", MOUSE_OUT: "mouseout", MOUSE_ENTER: "mouseenter", MOUSE_LEAVE: "mouseleave", MOUSE_WHEEL: "wheel", POINTER_UP: "pointerup", POINTER_DOWN: "pointerdown", POINTER_MOVE: "pointermove", POINTER_LEAVE: "pointerleave", CONTEXT_MENU: "contextmenu", WHEEL: "wheel", KEY_DOWN: "keydown", KEY_PRESS: "keypress", KEY_UP: "keyup", LOAD: "load", BEFORE_UNLOAD: "beforeunload", UNLOAD: "unload", PAGE_SHOW: "pageshow", PAGE_HIDE: "pagehide", PASTE: "paste", ABORT: "abort", ERROR: "error", RESIZE: "resize", SCROLL: "scroll", FULLSCREEN_CHANGE: "fullscreenchange", WK_FULLSCREEN_CHANGE: "webkitfullscreenchange", SELECT: "select", CHANGE: "change", SUBMIT: "submit", RESET: "reset", FOCUS: "focus", FOCUS_IN: "focusin", FOCUS_OUT: "focusout", BLUR: "blur", INPUT: "input", STORAGE: "storage", DRAG_START: "dragstart", DRAG: "drag", DRAG_ENTER: "dragenter", DRAG_LEAVE: "dragleave", DRAG_OVER: "dragover", DROP: "drop", DRAG_END: "dragend", ANIMATION_START: Bt3 ? "webkitAnimationStart" : "animationstart", ANIMATION_END: Bt3 ? "webkitAnimationEnd" : "animationend", ANIMATION_ITERATION: Bt3 ? "webkitAnimationIteration" : "animationiteration" };
var Da = /([\w\-]+)?(#([\w\-]+))?((\.([\w\-]+))*)/;
function Ho(s15, t, e, ...i8) {
  let r5 = Da.exec(t);
  if (!r5) throw new Error("Bad use of emmet");
  let n = r5[1] || "div", o2;
  return s15 !== "http://www.w3.org/1999/xhtml" ? o2 = document.createElementNS(s15, n) : o2 = document.createElement(n), r5[3] && (o2.id = r5[3]), r5[4] && (o2.className = r5[4].replace(/\./g, " ").trim()), e && Object.entries(e).forEach(([l, a]) => {
    typeof a > "u" || (/^on\w+$/.test(l) ? o2[l] = a : l === "selected" ? a && o2.setAttribute(l, "true") : o2.setAttribute(l, a));
  }), o2.append(...i8), o2;
}
function Ra(s15, t, ...e) {
  return Ho("http://www.w3.org/1999/xhtml", s15, t, ...e);
}
Ra.SVG = function(s15, t, ...e) {
  return Ho("http://www.w3.org/2000/svg", s15, t, ...e);
};
var ls = class {
  constructor(t) {
    this.domNode = t;
    this._maxWidth = "";
    this._width = "";
    this._height = "";
    this._top = "";
    this._left = "";
    this._bottom = "";
    this._right = "";
    this._paddingTop = "";
    this._paddingLeft = "";
    this._paddingBottom = "";
    this._paddingRight = "";
    this._fontFamily = "";
    this._fontWeight = "";
    this._fontSize = "";
    this._fontStyle = "";
    this._fontFeatureSettings = "";
    this._fontVariationSettings = "";
    this._textDecoration = "";
    this._lineHeight = "";
    this._letterSpacing = "";
    this._className = "";
    this._display = "";
    this._position = "";
    this._visibility = "";
    this._color = "";
    this._backgroundColor = "";
    this._layerHint = false;
    this._contain = "none";
    this._boxShadow = "";
  }
  setMaxWidth(t) {
    let e = Ie2(t);
    this._maxWidth !== e && (this._maxWidth = e, this.domNode.style.maxWidth = this._maxWidth);
  }
  setWidth(t) {
    let e = Ie2(t);
    this._width !== e && (this._width = e, this.domNode.style.width = this._width);
  }
  setHeight(t) {
    let e = Ie2(t);
    this._height !== e && (this._height = e, this.domNode.style.height = this._height);
  }
  setTop(t) {
    let e = Ie2(t);
    this._top !== e && (this._top = e, this.domNode.style.top = this._top);
  }
  setLeft(t) {
    let e = Ie2(t);
    this._left !== e && (this._left = e, this.domNode.style.left = this._left);
  }
  setBottom(t) {
    let e = Ie2(t);
    this._bottom !== e && (this._bottom = e, this.domNode.style.bottom = this._bottom);
  }
  setRight(t) {
    let e = Ie2(t);
    this._right !== e && (this._right = e, this.domNode.style.right = this._right);
  }
  setPaddingTop(t) {
    let e = Ie2(t);
    this._paddingTop !== e && (this._paddingTop = e, this.domNode.style.paddingTop = this._paddingTop);
  }
  setPaddingLeft(t) {
    let e = Ie2(t);
    this._paddingLeft !== e && (this._paddingLeft = e, this.domNode.style.paddingLeft = this._paddingLeft);
  }
  setPaddingBottom(t) {
    let e = Ie2(t);
    this._paddingBottom !== e && (this._paddingBottom = e, this.domNode.style.paddingBottom = this._paddingBottom);
  }
  setPaddingRight(t) {
    let e = Ie2(t);
    this._paddingRight !== e && (this._paddingRight = e, this.domNode.style.paddingRight = this._paddingRight);
  }
  setFontFamily(t) {
    this._fontFamily !== t && (this._fontFamily = t, this.domNode.style.fontFamily = this._fontFamily);
  }
  setFontWeight(t) {
    this._fontWeight !== t && (this._fontWeight = t, this.domNode.style.fontWeight = this._fontWeight);
  }
  setFontSize(t) {
    let e = Ie2(t);
    this._fontSize !== e && (this._fontSize = e, this.domNode.style.fontSize = this._fontSize);
  }
  setFontStyle(t) {
    this._fontStyle !== t && (this._fontStyle = t, this.domNode.style.fontStyle = this._fontStyle);
  }
  setFontFeatureSettings(t) {
    this._fontFeatureSettings !== t && (this._fontFeatureSettings = t, this.domNode.style.fontFeatureSettings = this._fontFeatureSettings);
  }
  setFontVariationSettings(t) {
    this._fontVariationSettings !== t && (this._fontVariationSettings = t, this.domNode.style.fontVariationSettings = this._fontVariationSettings);
  }
  setTextDecoration(t) {
    this._textDecoration !== t && (this._textDecoration = t, this.domNode.style.textDecoration = this._textDecoration);
  }
  setLineHeight(t) {
    let e = Ie2(t);
    this._lineHeight !== e && (this._lineHeight = e, this.domNode.style.lineHeight = this._lineHeight);
  }
  setLetterSpacing(t) {
    let e = Ie2(t);
    this._letterSpacing !== e && (this._letterSpacing = e, this.domNode.style.letterSpacing = this._letterSpacing);
  }
  setClassName(t) {
    this._className !== t && (this._className = t, this.domNode.className = this._className);
  }
  toggleClassName(t, e) {
    this.domNode.classList.toggle(t, e), this._className = this.domNode.className;
  }
  setDisplay(t) {
    this._display !== t && (this._display = t, this.domNode.style.display = this._display);
  }
  setPosition(t) {
    this._position !== t && (this._position = t, this.domNode.style.position = this._position);
  }
  setVisibility(t) {
    this._visibility !== t && (this._visibility = t, this.domNode.style.visibility = this._visibility);
  }
  setColor(t) {
    this._color !== t && (this._color = t, this.domNode.style.color = this._color);
  }
  setBackgroundColor(t) {
    this._backgroundColor !== t && (this._backgroundColor = t, this.domNode.style.backgroundColor = this._backgroundColor);
  }
  setLayerHinting(t) {
    this._layerHint !== t && (this._layerHint = t, this.domNode.style.transform = this._layerHint ? "translate3d(0px, 0px, 0px)" : "");
  }
  setBoxShadow(t) {
    this._boxShadow !== t && (this._boxShadow = t, this.domNode.style.boxShadow = t);
  }
  setContain(t) {
    this._contain !== t && (this._contain = t, this.domNode.style.contain = this._contain);
  }
  setAttribute(t, e) {
    this.domNode.setAttribute(t, e);
  }
  removeAttribute(t) {
    this.domNode.removeAttribute(t);
  }
  appendChild(t) {
    this.domNode.appendChild(t.domNode);
  }
  removeChild(t) {
    this.domNode.removeChild(t.domNode);
  }
};
function Ie2(s15) {
  return typeof s15 == "number" ? `${s15}px` : s15;
}
function _t(s15) {
  return new ls(s15);
}
var Wt2 = class {
  constructor() {
    this._hooks = new Ee3();
    this._pointerMoveCallback = null;
    this._onStopCallback = null;
  }
  dispose() {
    this.stopMonitoring(false), this._hooks.dispose();
  }
  stopMonitoring(t, e) {
    if (!this.isMonitoring()) return;
    this._hooks.clear(), this._pointerMoveCallback = null;
    let i8 = this._onStopCallback;
    this._onStopCallback = null, t && i8 && i8(e);
  }
  isMonitoring() {
    return !!this._pointerMoveCallback;
  }
  startMonitoring(t, e, i8, r5, n) {
    this.isMonitoring() && this.stopMonitoring(false), this._pointerMoveCallback = r5, this._onStopCallback = n;
    let o2 = t;
    try {
      t.setPointerCapture(e), this._hooks.add(C3(() => {
        try {
          t.releasePointerCapture(e);
        } catch {
        }
      }));
    } catch {
      o2 = be3(t);
    }
    this._hooks.add(L(o2, Y3.POINTER_MOVE, (l) => {
      if (l.buttons !== i8) {
        this.stopMonitoring(true);
        return;
      }
      l.preventDefault(), this._pointerMoveCallback(l);
    })), this._hooks.add(L(o2, Y3.POINTER_UP, (l) => this.stopMonitoring(true)));
  }
};
function Wo(s15, t, e) {
  let i8 = null, r5 = null;
  if (typeof e.value == "function" ? (i8 = "value", r5 = e.value, r5.length !== 0 && console.warn("Memoize should only be used in functions with zero parameters")) : typeof e.get == "function" && (i8 = "get", r5 = e.get), !r5) throw new Error("not supported");
  let n = `$memoize$${t}`;
  e[i8] = function(...o2) {
    return this.hasOwnProperty(n) || Object.defineProperty(this, n, { configurable: false, enumerable: false, writable: false, value: r5.apply(this, o2) }), this[n];
  };
}
var He4;
((n) => (n.Tap = "-xterm-gesturetap", n.Change = "-xterm-gesturechange", n.Start = "-xterm-gesturestart", n.End = "-xterm-gesturesend", n.Contextmenu = "-xterm-gesturecontextmenu"))(He4 ||= {});
var Q3 = class Q4 extends D2 {
  constructor() {
    super();
    this.dispatched = false;
    this.targets = new Ct2();
    this.ignoreTargets = new Ct2();
    this.activeTouches = {}, this.handle = null, this._lastSetTapCountTime = 0, this._register($.runAndSubscribe(No, ({ window: e, disposables: i8 }) => {
      i8.add(L(e.document, "touchstart", (r5) => this.onTouchStart(r5), { passive: false })), i8.add(L(e.document, "touchend", (r5) => this.onTouchEnd(e, r5))), i8.add(L(e.document, "touchmove", (r5) => this.onTouchMove(r5), { passive: false }));
    }, { window: fe3, disposables: this._store }));
  }
  static addTarget(e) {
    if (!Q4.isTouchDevice()) return D2.None;
    Q4.INSTANCE || (Q4.INSTANCE = Gn2(new Q4()));
    let i8 = Q4.INSTANCE.targets.push(e);
    return C3(i8);
  }
  static ignoreTarget(e) {
    if (!Q4.isTouchDevice()) return D2.None;
    Q4.INSTANCE || (Q4.INSTANCE = Gn2(new Q4()));
    let i8 = Q4.INSTANCE.ignoreTargets.push(e);
    return C3(i8);
  }
  static isTouchDevice() {
    return "ontouchstart" in fe3 || navigator.maxTouchPoints > 0;
  }
  dispose() {
    this.handle && (this.handle.dispose(), this.handle = null), super.dispose();
  }
  onTouchStart(e) {
    let i8 = Date.now();
    this.handle && (this.handle.dispose(), this.handle = null);
    for (let r5 = 0, n = e.targetTouches.length; r5 < n; r5++) {
      let o2 = e.targetTouches.item(r5);
      this.activeTouches[o2.identifier] = { id: o2.identifier, initialTarget: o2.target, initialTimeStamp: i8, initialPageX: o2.pageX, initialPageY: o2.pageY, rollingTimestamps: [i8], rollingPageX: [o2.pageX], rollingPageY: [o2.pageY] };
      let l = this.newGestureEvent(He4.Start, o2.target);
      l.pageX = o2.pageX, l.pageY = o2.pageY, this.dispatchEvent(l);
    }
    this.dispatched && (e.preventDefault(), e.stopPropagation(), this.dispatched = false);
  }
  onTouchEnd(e, i8) {
    let r5 = Date.now(), n = Object.keys(this.activeTouches).length;
    for (let o2 = 0, l = i8.changedTouches.length; o2 < l; o2++) {
      let a = i8.changedTouches.item(o2);
      if (!this.activeTouches.hasOwnProperty(String(a.identifier))) {
        console.warn("move of an UNKNOWN touch", a);
        continue;
      }
      let u = this.activeTouches[a.identifier], h2 = Date.now() - u.initialTimeStamp;
      if (h2 < Q4.HOLD_DELAY && Math.abs(u.initialPageX - Se2(u.rollingPageX)) < 30 && Math.abs(u.initialPageY - Se2(u.rollingPageY)) < 30) {
        let c = this.newGestureEvent(He4.Tap, u.initialTarget);
        c.pageX = Se2(u.rollingPageX), c.pageY = Se2(u.rollingPageY), this.dispatchEvent(c);
      } else if (h2 >= Q4.HOLD_DELAY && Math.abs(u.initialPageX - Se2(u.rollingPageX)) < 30 && Math.abs(u.initialPageY - Se2(u.rollingPageY)) < 30) {
        let c = this.newGestureEvent(He4.Contextmenu, u.initialTarget);
        c.pageX = Se2(u.rollingPageX), c.pageY = Se2(u.rollingPageY), this.dispatchEvent(c);
      } else if (n === 1) {
        let c = Se2(u.rollingPageX), d = Se2(u.rollingPageY), _3 = Se2(u.rollingTimestamps) - u.rollingTimestamps[0], p = c - u.rollingPageX[0], m = d - u.rollingPageY[0], f = [...this.targets].filter((A2) => u.initialTarget instanceof Node && A2.contains(u.initialTarget));
        this.inertia(e, f, r5, Math.abs(p) / _3, p > 0 ? 1 : -1, c, Math.abs(m) / _3, m > 0 ? 1 : -1, d);
      }
      this.dispatchEvent(this.newGestureEvent(He4.End, u.initialTarget)), delete this.activeTouches[a.identifier];
    }
    this.dispatched && (i8.preventDefault(), i8.stopPropagation(), this.dispatched = false);
  }
  newGestureEvent(e, i8) {
    let r5 = document.createEvent("CustomEvent");
    return r5.initEvent(e, false, true), r5.initialTarget = i8, r5.tapCount = 0, r5;
  }
  dispatchEvent(e) {
    if (e.type === He4.Tap) {
      let i8 = (/* @__PURE__ */ new Date()).getTime(), r5 = 0;
      i8 - this._lastSetTapCountTime > Q4.CLEAR_TAP_COUNT_TIME ? r5 = 1 : r5 = 2, this._lastSetTapCountTime = i8, e.tapCount = r5;
    } else (e.type === He4.Change || e.type === He4.Contextmenu) && (this._lastSetTapCountTime = 0);
    if (e.initialTarget instanceof Node) {
      for (let r5 of this.ignoreTargets) if (r5.contains(e.initialTarget)) return;
      let i8 = [];
      for (let r5 of this.targets) if (r5.contains(e.initialTarget)) {
        let n = 0, o2 = e.initialTarget;
        for (; o2 && o2 !== r5; ) n++, o2 = o2.parentElement;
        i8.push([n, r5]);
      }
      i8.sort((r5, n) => r5[0] - n[0]);
      for (let [r5, n] of i8) n.dispatchEvent(e), this.dispatched = true;
    }
  }
  inertia(e, i8, r5, n, o2, l, a, u, h2) {
    this.handle = mt(e, () => {
      let c = Date.now(), d = c - r5, _3 = 0, p = 0, m = true;
      n += Q4.SCROLL_FRICTION * d, a += Q4.SCROLL_FRICTION * d, n > 0 && (m = false, _3 = o2 * n * d), a > 0 && (m = false, p = u * a * d);
      let f = this.newGestureEvent(He4.Change);
      f.translationX = _3, f.translationY = p, i8.forEach((A2) => A2.dispatchEvent(f)), m || this.inertia(e, i8, c, n, o2, l + _3, a, u, h2 + p);
    });
  }
  onTouchMove(e) {
    let i8 = Date.now();
    for (let r5 = 0, n = e.changedTouches.length; r5 < n; r5++) {
      let o2 = e.changedTouches.item(r5);
      if (!this.activeTouches.hasOwnProperty(String(o2.identifier))) {
        console.warn("end of an UNKNOWN touch", o2);
        continue;
      }
      let l = this.activeTouches[o2.identifier], a = this.newGestureEvent(He4.Change, l.initialTarget);
      a.translationX = o2.pageX - Se2(l.rollingPageX), a.translationY = o2.pageY - Se2(l.rollingPageY), a.pageX = o2.pageX, a.pageY = o2.pageY, this.dispatchEvent(a), l.rollingPageX.length > 3 && (l.rollingPageX.shift(), l.rollingPageY.shift(), l.rollingTimestamps.shift()), l.rollingPageX.push(o2.pageX), l.rollingPageY.push(o2.pageY), l.rollingTimestamps.push(i8);
    }
    this.dispatched && (e.preventDefault(), e.stopPropagation(), this.dispatched = false);
  }
};
Q3.SCROLL_FRICTION = -5e-3, Q3.HOLD_DELAY = 700, Q3.CLEAR_TAP_COUNT_TIME = 400, M3([Wo], Q3, "isTouchDevice", 1);
var Pr2 = Q3;
var lt3 = class extends D2 {
  onclick(t, e) {
    this._register(L(t, Y3.CLICK, (i8) => e(new qe3(be3(t), i8))));
  }
  onmousedown(t, e) {
    this._register(L(t, Y3.MOUSE_DOWN, (i8) => e(new qe3(be3(t), i8))));
  }
  onmouseover(t, e) {
    this._register(L(t, Y3.MOUSE_OVER, (i8) => e(new qe3(be3(t), i8))));
  }
  onmouseleave(t, e) {
    this._register(L(t, Y3.MOUSE_LEAVE, (i8) => e(new qe3(be3(t), i8))));
  }
  onkeydown(t, e) {
    this._register(L(t, Y3.KEY_DOWN, (i8) => e(new ft(i8))));
  }
  onkeyup(t, e) {
    this._register(L(t, Y3.KEY_UP, (i8) => e(new ft(i8))));
  }
  oninput(t, e) {
    this._register(L(t, Y3.INPUT, e));
  }
  onblur(t, e) {
    this._register(L(t, Y3.BLUR, e));
  }
  onfocus(t, e) {
    this._register(L(t, Y3.FOCUS, e));
  }
  onchange(t, e) {
    this._register(L(t, Y3.CHANGE, e));
  }
  ignoreGesture(t) {
    return Pr2.ignoreTarget(t);
  }
};
var Uo = 11;
var Or2 = class extends lt3 {
  constructor(t) {
    super(), this._onActivate = t.onActivate, this.bgDomNode = document.createElement("div"), this.bgDomNode.className = "arrow-background", this.bgDomNode.style.position = "absolute", this.bgDomNode.style.width = t.bgWidth + "px", this.bgDomNode.style.height = t.bgHeight + "px", typeof t.top < "u" && (this.bgDomNode.style.top = "0px"), typeof t.left < "u" && (this.bgDomNode.style.left = "0px"), typeof t.bottom < "u" && (this.bgDomNode.style.bottom = "0px"), typeof t.right < "u" && (this.bgDomNode.style.right = "0px"), this.domNode = document.createElement("div"), this.domNode.className = t.className, this.domNode.style.position = "absolute", this.domNode.style.width = Uo + "px", this.domNode.style.height = Uo + "px", typeof t.top < "u" && (this.domNode.style.top = t.top + "px"), typeof t.left < "u" && (this.domNode.style.left = t.left + "px"), typeof t.bottom < "u" && (this.domNode.style.bottom = t.bottom + "px"), typeof t.right < "u" && (this.domNode.style.right = t.right + "px"), this._pointerMoveMonitor = this._register(new Wt2()), this._register(os2(this.bgDomNode, Y3.POINTER_DOWN, (e) => this._arrowPointerDown(e))), this._register(os2(this.domNode, Y3.POINTER_DOWN, (e) => this._arrowPointerDown(e))), this._pointerdownRepeatTimer = this._register(new Mr2()), this._pointerdownScheduleRepeatTimer = this._register(new Ye3());
  }
  _arrowPointerDown(t) {
    if (!t.target || !(t.target instanceof Element)) return;
    let e = () => {
      this._pointerdownRepeatTimer.cancelAndSet(() => this._onActivate(), 1e3 / 24, be3(t));
    };
    this._onActivate(), this._pointerdownRepeatTimer.cancel(), this._pointerdownScheduleRepeatTimer.cancelAndSet(e, 200), this._pointerMoveMonitor.startMonitoring(t.target, t.pointerId, t.buttons, (i8) => {
    }, () => {
      this._pointerdownRepeatTimer.cancel(), this._pointerdownScheduleRepeatTimer.cancel();
    }), t.preventDefault();
  }
};
var cs2 = class s9 {
  constructor(t, e, i8, r5, n, o2, l) {
    this._forceIntegerValues = t;
    this._scrollStateBrand = void 0;
    this._forceIntegerValues && (e = e | 0, i8 = i8 | 0, r5 = r5 | 0, n = n | 0, o2 = o2 | 0, l = l | 0), this.rawScrollLeft = r5, this.rawScrollTop = l, e < 0 && (e = 0), r5 + e > i8 && (r5 = i8 - e), r5 < 0 && (r5 = 0), n < 0 && (n = 0), l + n > o2 && (l = o2 - n), l < 0 && (l = 0), this.width = e, this.scrollWidth = i8, this.scrollLeft = r5, this.height = n, this.scrollHeight = o2, this.scrollTop = l;
  }
  equals(t) {
    return this.rawScrollLeft === t.rawScrollLeft && this.rawScrollTop === t.rawScrollTop && this.width === t.width && this.scrollWidth === t.scrollWidth && this.scrollLeft === t.scrollLeft && this.height === t.height && this.scrollHeight === t.scrollHeight && this.scrollTop === t.scrollTop;
  }
  withScrollDimensions(t, e) {
    return new s9(this._forceIntegerValues, typeof t.width < "u" ? t.width : this.width, typeof t.scrollWidth < "u" ? t.scrollWidth : this.scrollWidth, e ? this.rawScrollLeft : this.scrollLeft, typeof t.height < "u" ? t.height : this.height, typeof t.scrollHeight < "u" ? t.scrollHeight : this.scrollHeight, e ? this.rawScrollTop : this.scrollTop);
  }
  withScrollPosition(t) {
    return new s9(this._forceIntegerValues, this.width, this.scrollWidth, typeof t.scrollLeft < "u" ? t.scrollLeft : this.rawScrollLeft, this.height, this.scrollHeight, typeof t.scrollTop < "u" ? t.scrollTop : this.rawScrollTop);
  }
  createScrollEvent(t, e) {
    let i8 = this.width !== t.width, r5 = this.scrollWidth !== t.scrollWidth, n = this.scrollLeft !== t.scrollLeft, o2 = this.height !== t.height, l = this.scrollHeight !== t.scrollHeight, a = this.scrollTop !== t.scrollTop;
    return { inSmoothScrolling: e, oldWidth: t.width, oldScrollWidth: t.scrollWidth, oldScrollLeft: t.scrollLeft, width: this.width, scrollWidth: this.scrollWidth, scrollLeft: this.scrollLeft, oldHeight: t.height, oldScrollHeight: t.scrollHeight, oldScrollTop: t.scrollTop, height: this.height, scrollHeight: this.scrollHeight, scrollTop: this.scrollTop, widthChanged: i8, scrollWidthChanged: r5, scrollLeftChanged: n, heightChanged: o2, scrollHeightChanged: l, scrollTopChanged: a };
  }
};
var Ri = class extends D2 {
  constructor(e) {
    super();
    this._scrollableBrand = void 0;
    this._onScroll = this._register(new v());
    this.onScroll = this._onScroll.event;
    this._smoothScrollDuration = e.smoothScrollDuration, this._scheduleAtNextAnimationFrame = e.scheduleAtNextAnimationFrame, this._state = new cs2(e.forceIntegerValues, 0, 0, 0, 0, 0, 0), this._smoothScrolling = null;
  }
  dispose() {
    this._smoothScrolling && (this._smoothScrolling.dispose(), this._smoothScrolling = null), super.dispose();
  }
  setSmoothScrollDuration(e) {
    this._smoothScrollDuration = e;
  }
  validateScrollPosition(e) {
    return this._state.withScrollPosition(e);
  }
  getScrollDimensions() {
    return this._state;
  }
  setScrollDimensions(e, i8) {
    let r5 = this._state.withScrollDimensions(e, i8);
    this._setState(r5, !!this._smoothScrolling), this._smoothScrolling?.acceptScrollDimensions(this._state);
  }
  getFutureScrollPosition() {
    return this._smoothScrolling ? this._smoothScrolling.to : this._state;
  }
  getCurrentScrollPosition() {
    return this._state;
  }
  setScrollPositionNow(e) {
    let i8 = this._state.withScrollPosition(e);
    this._smoothScrolling && (this._smoothScrolling.dispose(), this._smoothScrolling = null), this._setState(i8, false);
  }
  setScrollPositionSmooth(e, i8) {
    if (this._smoothScrollDuration === 0) return this.setScrollPositionNow(e);
    if (this._smoothScrolling) {
      e = { scrollLeft: typeof e.scrollLeft > "u" ? this._smoothScrolling.to.scrollLeft : e.scrollLeft, scrollTop: typeof e.scrollTop > "u" ? this._smoothScrolling.to.scrollTop : e.scrollTop };
      let r5 = this._state.withScrollPosition(e);
      if (this._smoothScrolling.to.scrollLeft === r5.scrollLeft && this._smoothScrolling.to.scrollTop === r5.scrollTop) return;
      let n;
      i8 ? n = new Nr(this._smoothScrolling.from, r5, this._smoothScrolling.startTime, this._smoothScrolling.duration) : n = this._smoothScrolling.combine(this._state, r5, this._smoothScrollDuration), this._smoothScrolling.dispose(), this._smoothScrolling = n;
    } else {
      let r5 = this._state.withScrollPosition(e);
      this._smoothScrolling = Nr.start(this._state, r5, this._smoothScrollDuration);
    }
    this._smoothScrolling.animationFrameDisposable = this._scheduleAtNextAnimationFrame(() => {
      this._smoothScrolling && (this._smoothScrolling.animationFrameDisposable = null, this._performSmoothScrolling());
    });
  }
  hasPendingScrollAnimation() {
    return !!this._smoothScrolling;
  }
  _performSmoothScrolling() {
    if (!this._smoothScrolling) return;
    let e = this._smoothScrolling.tick(), i8 = this._state.withScrollPosition(e);
    if (this._setState(i8, true), !!this._smoothScrolling) {
      if (e.isDone) {
        this._smoothScrolling.dispose(), this._smoothScrolling = null;
        return;
      }
      this._smoothScrolling.animationFrameDisposable = this._scheduleAtNextAnimationFrame(() => {
        this._smoothScrolling && (this._smoothScrolling.animationFrameDisposable = null, this._performSmoothScrolling());
      });
    }
  }
  _setState(e, i8) {
    let r5 = this._state;
    r5.equals(e) || (this._state = e, this._onScroll.fire(this._state.createScrollEvent(r5, i8)));
  }
};
var Br2 = class {
  constructor(t, e, i8) {
    this.scrollLeft = t, this.scrollTop = e, this.isDone = i8;
  }
};
function as2(s15, t) {
  let e = t - s15;
  return function(i8) {
    return s15 + e * ka(i8);
  };
}
function La(s15, t, e) {
  return function(i8) {
    return i8 < e ? s15(i8 / e) : t((i8 - e) / (1 - e));
  };
}
var Nr = class s10 {
  constructor(t, e, i8, r5) {
    this.from = t, this.to = e, this.duration = r5, this.startTime = i8, this.animationFrameDisposable = null, this._initAnimations();
  }
  _initAnimations() {
    this.scrollLeft = this._initAnimation(this.from.scrollLeft, this.to.scrollLeft, this.to.width), this.scrollTop = this._initAnimation(this.from.scrollTop, this.to.scrollTop, this.to.height);
  }
  _initAnimation(t, e, i8) {
    if (Math.abs(t - e) > 2.5 * i8) {
      let n, o2;
      return t < e ? (n = t + 0.75 * i8, o2 = e - 0.75 * i8) : (n = t - 0.75 * i8, o2 = e + 0.75 * i8), La(as2(t, n), as2(o2, e), 0.33);
    }
    return as2(t, e);
  }
  dispose() {
    this.animationFrameDisposable !== null && (this.animationFrameDisposable.dispose(), this.animationFrameDisposable = null);
  }
  acceptScrollDimensions(t) {
    this.to = t.withScrollPosition(this.to), this._initAnimations();
  }
  tick() {
    return this._tick(Date.now());
  }
  _tick(t) {
    let e = (t - this.startTime) / this.duration;
    if (e < 1) {
      let i8 = this.scrollLeft(e), r5 = this.scrollTop(e);
      return new Br2(i8, r5, false);
    }
    return new Br2(this.to.scrollLeft, this.to.scrollTop, true);
  }
  combine(t, e, i8) {
    return s10.start(t, e, i8);
  }
  static start(t, e, i8) {
    i8 = i8 + 10;
    let r5 = Date.now() - 10;
    return new s10(t, e, r5, i8);
  }
};
function Aa(s15) {
  return Math.pow(s15, 3);
}
function ka(s15) {
  return 1 - Aa(1 - s15);
}
var Fr2 = class extends D2 {
  constructor(t, e, i8) {
    super(), this._visibility = t, this._visibleClassName = e, this._invisibleClassName = i8, this._domNode = null, this._isVisible = false, this._isNeeded = false, this._rawShouldBeVisible = false, this._shouldBeVisible = false, this._revealTimer = this._register(new Ye3());
  }
  setVisibility(t) {
    this._visibility !== t && (this._visibility = t, this._updateShouldBeVisible());
  }
  setShouldBeVisible(t) {
    this._rawShouldBeVisible = t, this._updateShouldBeVisible();
  }
  _applyVisibilitySetting() {
    return this._visibility === 2 ? false : this._visibility === 3 ? true : this._rawShouldBeVisible;
  }
  _updateShouldBeVisible() {
    let t = this._applyVisibilitySetting();
    this._shouldBeVisible !== t && (this._shouldBeVisible = t, this.ensureVisibility());
  }
  setIsNeeded(t) {
    this._isNeeded !== t && (this._isNeeded = t, this.ensureVisibility());
  }
  setDomNode(t) {
    this._domNode = t, this._domNode.setClassName(this._invisibleClassName), this.setShouldBeVisible(false);
  }
  ensureVisibility() {
    if (!this._isNeeded) {
      this._hide(false);
      return;
    }
    this._shouldBeVisible ? this._reveal() : this._hide(true);
  }
  _reveal() {
    this._isVisible || (this._isVisible = true, this._revealTimer.setIfNotSet(() => {
      this._domNode?.setClassName(this._visibleClassName);
    }, 0));
  }
  _hide(t) {
    this._revealTimer.cancel(), this._isVisible && (this._isVisible = false, this._domNode?.setClassName(this._invisibleClassName + (t ? " fade" : "")));
  }
};
var Ca2 = 140;
var Ut2 = class extends lt3 {
  constructor(t) {
    super(), this._lazyRender = t.lazyRender, this._host = t.host, this._scrollable = t.scrollable, this._scrollByPage = t.scrollByPage, this._scrollbarState = t.scrollbarState, this._visibilityController = this._register(new Fr2(t.visibility, "visible scrollbar " + t.extraScrollbarClassName, "invisible scrollbar " + t.extraScrollbarClassName)), this._visibilityController.setIsNeeded(this._scrollbarState.isNeeded()), this._pointerMoveMonitor = this._register(new Wt2()), this._shouldRender = true, this.domNode = _t(document.createElement("div")), this.domNode.setAttribute("role", "presentation"), this.domNode.setAttribute("aria-hidden", "true"), this._visibilityController.setDomNode(this.domNode), this.domNode.setPosition("absolute"), this._register(L(this.domNode.domNode, Y3.POINTER_DOWN, (e) => this._domNodePointerDown(e)));
  }
  _createArrow(t) {
    let e = this._register(new Or2(t));
    this.domNode.domNode.appendChild(e.bgDomNode), this.domNode.domNode.appendChild(e.domNode);
  }
  _createSlider(t, e, i8, r5) {
    this.slider = _t(document.createElement("div")), this.slider.setClassName("slider"), this.slider.setPosition("absolute"), this.slider.setTop(t), this.slider.setLeft(e), typeof i8 == "number" && this.slider.setWidth(i8), typeof r5 == "number" && this.slider.setHeight(r5), this.slider.setLayerHinting(true), this.slider.setContain("strict"), this.domNode.domNode.appendChild(this.slider.domNode), this._register(L(this.slider.domNode, Y3.POINTER_DOWN, (n) => {
      n.button === 0 && (n.preventDefault(), this._sliderPointerDown(n));
    })), this.onclick(this.slider.domNode, (n) => {
      n.leftButton && n.stopPropagation();
    });
  }
  _onElementSize(t) {
    return this._scrollbarState.setVisibleSize(t) && (this._visibilityController.setIsNeeded(this._scrollbarState.isNeeded()), this._shouldRender = true, this._lazyRender || this.render()), this._shouldRender;
  }
  _onElementScrollSize(t) {
    return this._scrollbarState.setScrollSize(t) && (this._visibilityController.setIsNeeded(this._scrollbarState.isNeeded()), this._shouldRender = true, this._lazyRender || this.render()), this._shouldRender;
  }
  _onElementScrollPosition(t) {
    return this._scrollbarState.setScrollPosition(t) && (this._visibilityController.setIsNeeded(this._scrollbarState.isNeeded()), this._shouldRender = true, this._lazyRender || this.render()), this._shouldRender;
  }
  beginReveal() {
    this._visibilityController.setShouldBeVisible(true);
  }
  beginHide() {
    this._visibilityController.setShouldBeVisible(false);
  }
  render() {
    this._shouldRender && (this._shouldRender = false, this._renderDomNode(this._scrollbarState.getRectangleLargeSize(), this._scrollbarState.getRectangleSmallSize()), this._updateSlider(this._scrollbarState.getSliderSize(), this._scrollbarState.getArrowSize() + this._scrollbarState.getSliderPosition()));
  }
  _domNodePointerDown(t) {
    t.target === this.domNode.domNode && this._onPointerDown(t);
  }
  delegatePointerDown(t) {
    let e = this.domNode.domNode.getClientRects()[0].top, i8 = e + this._scrollbarState.getSliderPosition(), r5 = e + this._scrollbarState.getSliderPosition() + this._scrollbarState.getSliderSize(), n = this._sliderPointerPosition(t);
    i8 <= n && n <= r5 ? t.button === 0 && (t.preventDefault(), this._sliderPointerDown(t)) : this._onPointerDown(t);
  }
  _onPointerDown(t) {
    let e, i8;
    if (t.target === this.domNode.domNode && typeof t.offsetX == "number" && typeof t.offsetY == "number") e = t.offsetX, i8 = t.offsetY;
    else {
      let n = Fo(this.domNode.domNode);
      e = t.pageX - n.left, i8 = t.pageY - n.top;
    }
    let r5 = this._pointerDownRelativePosition(e, i8);
    this._setDesiredScrollPositionNow(this._scrollByPage ? this._scrollbarState.getDesiredScrollPositionFromOffsetPaged(r5) : this._scrollbarState.getDesiredScrollPositionFromOffset(r5)), t.button === 0 && (t.preventDefault(), this._sliderPointerDown(t));
  }
  _sliderPointerDown(t) {
    if (!t.target || !(t.target instanceof Element)) return;
    let e = this._sliderPointerPosition(t), i8 = this._sliderOrthogonalPointerPosition(t), r5 = this._scrollbarState.clone();
    this.slider.toggleClassName("active", true), this._pointerMoveMonitor.startMonitoring(t.target, t.pointerId, t.buttons, (n) => {
      let o2 = this._sliderOrthogonalPointerPosition(n), l = Math.abs(o2 - i8);
      if (wr2 && l > Ca2) {
        this._setDesiredScrollPositionNow(r5.getScrollPosition());
        return;
      }
      let u = this._sliderPointerPosition(n) - e;
      this._setDesiredScrollPositionNow(r5.getDesiredScrollPositionFromDelta(u));
    }, () => {
      this.slider.toggleClassName("active", false), this._host.onDragEnd();
    }), this._host.onDragStart();
  }
  _setDesiredScrollPositionNow(t) {
    let e = {};
    this.writeScrollPosition(e, t), this._scrollable.setScrollPositionNow(e);
  }
  updateScrollbarSize(t) {
    this._updateScrollbarSize(t), this._scrollbarState.setScrollbarSize(t), this._shouldRender = true, this._lazyRender || this.render();
  }
  isNeeded() {
    return this._scrollbarState.isNeeded();
  }
};
var Kt2 = class s11 {
  constructor(t, e, i8, r5, n, o2) {
    this._scrollbarSize = Math.round(e), this._oppositeScrollbarSize = Math.round(i8), this._arrowSize = Math.round(t), this._visibleSize = r5, this._scrollSize = n, this._scrollPosition = o2, this._computedAvailableSize = 0, this._computedIsNeeded = false, this._computedSliderSize = 0, this._computedSliderRatio = 0, this._computedSliderPosition = 0, this._refreshComputedValues();
  }
  clone() {
    return new s11(this._arrowSize, this._scrollbarSize, this._oppositeScrollbarSize, this._visibleSize, this._scrollSize, this._scrollPosition);
  }
  setVisibleSize(t) {
    let e = Math.round(t);
    return this._visibleSize !== e ? (this._visibleSize = e, this._refreshComputedValues(), true) : false;
  }
  setScrollSize(t) {
    let e = Math.round(t);
    return this._scrollSize !== e ? (this._scrollSize = e, this._refreshComputedValues(), true) : false;
  }
  setScrollPosition(t) {
    let e = Math.round(t);
    return this._scrollPosition !== e ? (this._scrollPosition = e, this._refreshComputedValues(), true) : false;
  }
  setScrollbarSize(t) {
    this._scrollbarSize = Math.round(t);
  }
  setOppositeScrollbarSize(t) {
    this._oppositeScrollbarSize = Math.round(t);
  }
  static _computeValues(t, e, i8, r5, n) {
    let o2 = Math.max(0, i8 - t), l = Math.max(0, o2 - 2 * e), a = r5 > 0 && r5 > i8;
    if (!a) return { computedAvailableSize: Math.round(o2), computedIsNeeded: a, computedSliderSize: Math.round(l), computedSliderRatio: 0, computedSliderPosition: 0 };
    let u = Math.round(Math.max(20, Math.floor(i8 * l / r5))), h2 = (l - u) / (r5 - i8), c = n * h2;
    return { computedAvailableSize: Math.round(o2), computedIsNeeded: a, computedSliderSize: Math.round(u), computedSliderRatio: h2, computedSliderPosition: Math.round(c) };
  }
  _refreshComputedValues() {
    let t = s11._computeValues(this._oppositeScrollbarSize, this._arrowSize, this._visibleSize, this._scrollSize, this._scrollPosition);
    this._computedAvailableSize = t.computedAvailableSize, this._computedIsNeeded = t.computedIsNeeded, this._computedSliderSize = t.computedSliderSize, this._computedSliderRatio = t.computedSliderRatio, this._computedSliderPosition = t.computedSliderPosition;
  }
  getArrowSize() {
    return this._arrowSize;
  }
  getScrollPosition() {
    return this._scrollPosition;
  }
  getRectangleLargeSize() {
    return this._computedAvailableSize;
  }
  getRectangleSmallSize() {
    return this._scrollbarSize;
  }
  isNeeded() {
    return this._computedIsNeeded;
  }
  getSliderSize() {
    return this._computedSliderSize;
  }
  getSliderPosition() {
    return this._computedSliderPosition;
  }
  getDesiredScrollPositionFromOffset(t) {
    if (!this._computedIsNeeded) return 0;
    let e = t - this._arrowSize - this._computedSliderSize / 2;
    return Math.round(e / this._computedSliderRatio);
  }
  getDesiredScrollPositionFromOffsetPaged(t) {
    if (!this._computedIsNeeded) return 0;
    let e = t - this._arrowSize, i8 = this._scrollPosition;
    return e < this._computedSliderPosition ? i8 -= this._visibleSize : i8 += this._visibleSize, i8;
  }
  getDesiredScrollPositionFromDelta(t) {
    if (!this._computedIsNeeded) return 0;
    let e = this._computedSliderPosition + t;
    return Math.round(e / this._computedSliderRatio);
  }
};
var Wr2 = class extends Ut2 {
  constructor(t, e, i8) {
    let r5 = t.getScrollDimensions(), n = t.getCurrentScrollPosition();
    if (super({ lazyRender: e.lazyRender, host: i8, scrollbarState: new Kt2(e.horizontalHasArrows ? e.arrowSize : 0, e.horizontal === 2 ? 0 : e.horizontalScrollbarSize, e.vertical === 2 ? 0 : e.verticalScrollbarSize, r5.width, r5.scrollWidth, n.scrollLeft), visibility: e.horizontal, extraScrollbarClassName: "horizontal", scrollable: t, scrollByPage: e.scrollByPage }), e.horizontalHasArrows) throw new Error("horizontalHasArrows is not supported in xterm.js");
    this._createSlider(Math.floor((e.horizontalScrollbarSize - e.horizontalSliderSize) / 2), 0, void 0, e.horizontalSliderSize);
  }
  _updateSlider(t, e) {
    this.slider.setWidth(t), this.slider.setLeft(e);
  }
  _renderDomNode(t, e) {
    this.domNode.setWidth(t), this.domNode.setHeight(e), this.domNode.setLeft(0), this.domNode.setBottom(0);
  }
  onDidScroll(t) {
    return this._shouldRender = this._onElementScrollSize(t.scrollWidth) || this._shouldRender, this._shouldRender = this._onElementScrollPosition(t.scrollLeft) || this._shouldRender, this._shouldRender = this._onElementSize(t.width) || this._shouldRender, this._shouldRender;
  }
  _pointerDownRelativePosition(t, e) {
    return t;
  }
  _sliderPointerPosition(t) {
    return t.pageX;
  }
  _sliderOrthogonalPointerPosition(t) {
    return t.pageY;
  }
  _updateScrollbarSize(t) {
    this.slider.setHeight(t);
  }
  writeScrollPosition(t, e) {
    t.scrollLeft = e;
  }
  updateOptions(t) {
    this.updateScrollbarSize(t.horizontal === 2 ? 0 : t.horizontalScrollbarSize), this._scrollbarState.setOppositeScrollbarSize(t.vertical === 2 ? 0 : t.verticalScrollbarSize), this._visibilityController.setVisibility(t.horizontal), this._scrollByPage = t.scrollByPage;
  }
};
var Ur = class extends Ut2 {
  constructor(t, e, i8) {
    let r5 = t.getScrollDimensions(), n = t.getCurrentScrollPosition();
    if (super({ lazyRender: e.lazyRender, host: i8, scrollbarState: new Kt2(e.verticalHasArrows ? e.arrowSize : 0, e.vertical === 2 ? 0 : e.verticalScrollbarSize, 0, r5.height, r5.scrollHeight, n.scrollTop), visibility: e.vertical, extraScrollbarClassName: "vertical", scrollable: t, scrollByPage: e.scrollByPage }), e.verticalHasArrows) throw new Error("horizontalHasArrows is not supported in xterm.js");
    this._createSlider(0, Math.floor((e.verticalScrollbarSize - e.verticalSliderSize) / 2), e.verticalSliderSize, void 0);
  }
  _updateSlider(t, e) {
    this.slider.setHeight(t), this.slider.setTop(e);
  }
  _renderDomNode(t, e) {
    this.domNode.setWidth(e), this.domNode.setHeight(t), this.domNode.setRight(0), this.domNode.setTop(0);
  }
  onDidScroll(t) {
    return this._shouldRender = this._onElementScrollSize(t.scrollHeight) || this._shouldRender, this._shouldRender = this._onElementScrollPosition(t.scrollTop) || this._shouldRender, this._shouldRender = this._onElementSize(t.height) || this._shouldRender, this._shouldRender;
  }
  _pointerDownRelativePosition(t, e) {
    return e;
  }
  _sliderPointerPosition(t) {
    return t.pageY;
  }
  _sliderOrthogonalPointerPosition(t) {
    return t.pageX;
  }
  _updateScrollbarSize(t) {
    this.slider.setWidth(t);
  }
  writeScrollPosition(t, e) {
    t.scrollTop = e;
  }
  updateOptions(t) {
    this.updateScrollbarSize(t.vertical === 2 ? 0 : t.verticalScrollbarSize), this._scrollbarState.setOppositeScrollbarSize(0), this._visibilityController.setVisibility(t.vertical), this._scrollByPage = t.scrollByPage;
  }
};
var Ma = 500;
var Ko = 50;
var zo = true;
var us2 = class {
  constructor(t, e, i8) {
    this.timestamp = t, this.deltaX = e, this.deltaY = i8, this.score = 0;
  }
};
var zr2 = class zr3 {
  constructor() {
    this._capacity = 5, this._memory = [], this._front = -1, this._rear = -1;
  }
  isPhysicalMouseWheel() {
    if (this._front === -1 && this._rear === -1) return false;
    let t = 1, e = 0, i8 = 1, r5 = this._rear;
    do {
      let n = r5 === this._front ? t : Math.pow(2, -i8);
      if (t -= n, e += this._memory[r5].score * n, r5 === this._front) break;
      r5 = (this._capacity + r5 - 1) % this._capacity, i8++;
    } while (true);
    return e <= 0.5;
  }
  acceptStandardWheelEvent(t) {
    if (Ti2) {
      let e = be3(t.browserEvent), i8 = mo(e);
      this.accept(Date.now(), t.deltaX * i8, t.deltaY * i8);
    } else this.accept(Date.now(), t.deltaX, t.deltaY);
  }
  accept(t, e, i8) {
    let r5 = null, n = new us2(t, e, i8);
    this._front === -1 && this._rear === -1 ? (this._memory[0] = n, this._front = 0, this._rear = 0) : (r5 = this._memory[this._rear], this._rear = (this._rear + 1) % this._capacity, this._rear === this._front && (this._front = (this._front + 1) % this._capacity), this._memory[this._rear] = n), n.score = this._computeScore(n, r5);
  }
  _computeScore(t, e) {
    if (Math.abs(t.deltaX) > 0 && Math.abs(t.deltaY) > 0) return 1;
    let i8 = 0.5;
    if ((!this._isAlmostInt(t.deltaX) || !this._isAlmostInt(t.deltaY)) && (i8 += 0.25), e) {
      let r5 = Math.abs(t.deltaX), n = Math.abs(t.deltaY), o2 = Math.abs(e.deltaX), l = Math.abs(e.deltaY), a = Math.max(Math.min(r5, o2), 1), u = Math.max(Math.min(n, l), 1), h2 = Math.max(r5, o2), c = Math.max(n, l);
      h2 % a === 0 && c % u === 0 && (i8 -= 0.5);
    }
    return Math.min(Math.max(i8, 0), 1);
  }
  _isAlmostInt(t) {
    return Math.abs(Math.round(t) - t) < 0.01;
  }
};
zr2.INSTANCE = new zr2();
var hs2 = zr2;
var ds2 = class extends lt3 {
  constructor(e, i8, r5) {
    super();
    this._onScroll = this._register(new v());
    this.onScroll = this._onScroll.event;
    this._onWillScroll = this._register(new v());
    this.onWillScroll = this._onWillScroll.event;
    this._options = Pa(i8), this._scrollable = r5, this._register(this._scrollable.onScroll((o2) => {
      this._onWillScroll.fire(o2), this._onDidScroll(o2), this._onScroll.fire(o2);
    }));
    let n = { onMouseWheel: (o2) => this._onMouseWheel(o2), onDragStart: () => this._onDragStart(), onDragEnd: () => this._onDragEnd() };
    this._verticalScrollbar = this._register(new Ur(this._scrollable, this._options, n)), this._horizontalScrollbar = this._register(new Wr2(this._scrollable, this._options, n)), this._domNode = document.createElement("div"), this._domNode.className = "xterm-scrollable-element " + this._options.className, this._domNode.setAttribute("role", "presentation"), this._domNode.style.position = "relative", this._domNode.appendChild(e), this._domNode.appendChild(this._horizontalScrollbar.domNode.domNode), this._domNode.appendChild(this._verticalScrollbar.domNode.domNode), this._options.useShadows ? (this._leftShadowDomNode = _t(document.createElement("div")), this._leftShadowDomNode.setClassName("shadow"), this._domNode.appendChild(this._leftShadowDomNode.domNode), this._topShadowDomNode = _t(document.createElement("div")), this._topShadowDomNode.setClassName("shadow"), this._domNode.appendChild(this._topShadowDomNode.domNode), this._topLeftShadowDomNode = _t(document.createElement("div")), this._topLeftShadowDomNode.setClassName("shadow"), this._domNode.appendChild(this._topLeftShadowDomNode.domNode)) : (this._leftShadowDomNode = null, this._topShadowDomNode = null, this._topLeftShadowDomNode = null), this._listenOnDomNode = this._options.listenOnDomNode || this._domNode, this._mouseWheelToDispose = [], this._setListeningToMouseWheel(this._options.handleMouseWheel), this.onmouseover(this._listenOnDomNode, (o2) => this._onMouseOver(o2)), this.onmouseleave(this._listenOnDomNode, (o2) => this._onMouseLeave(o2)), this._hideTimeout = this._register(new Ye3()), this._isDragging = false, this._mouseIsOver = false, this._shouldRender = true, this._revealOnScroll = true;
  }
  get options() {
    return this._options;
  }
  dispose() {
    this._mouseWheelToDispose = Ne2(this._mouseWheelToDispose), super.dispose();
  }
  getDomNode() {
    return this._domNode;
  }
  getOverviewRulerLayoutInfo() {
    return { parent: this._domNode, insertBefore: this._verticalScrollbar.domNode.domNode };
  }
  delegateVerticalScrollbarPointerDown(e) {
    this._verticalScrollbar.delegatePointerDown(e);
  }
  getScrollDimensions() {
    return this._scrollable.getScrollDimensions();
  }
  setScrollDimensions(e) {
    this._scrollable.setScrollDimensions(e, false);
  }
  updateClassName(e) {
    this._options.className = e, Te3 && (this._options.className += " mac"), this._domNode.className = "xterm-scrollable-element " + this._options.className;
  }
  updateOptions(e) {
    typeof e.handleMouseWheel < "u" && (this._options.handleMouseWheel = e.handleMouseWheel, this._setListeningToMouseWheel(this._options.handleMouseWheel)), typeof e.mouseWheelScrollSensitivity < "u" && (this._options.mouseWheelScrollSensitivity = e.mouseWheelScrollSensitivity), typeof e.fastScrollSensitivity < "u" && (this._options.fastScrollSensitivity = e.fastScrollSensitivity), typeof e.scrollPredominantAxis < "u" && (this._options.scrollPredominantAxis = e.scrollPredominantAxis), typeof e.horizontal < "u" && (this._options.horizontal = e.horizontal), typeof e.vertical < "u" && (this._options.vertical = e.vertical), typeof e.horizontalScrollbarSize < "u" && (this._options.horizontalScrollbarSize = e.horizontalScrollbarSize), typeof e.verticalScrollbarSize < "u" && (this._options.verticalScrollbarSize = e.verticalScrollbarSize), typeof e.scrollByPage < "u" && (this._options.scrollByPage = e.scrollByPage), this._horizontalScrollbar.updateOptions(this._options), this._verticalScrollbar.updateOptions(this._options), this._options.lazyRender || this._render();
  }
  setRevealOnScroll(e) {
    this._revealOnScroll = e;
  }
  delegateScrollFromMouseWheelEvent(e) {
    this._onMouseWheel(new xi2(e));
  }
  _setListeningToMouseWheel(e) {
    if (this._mouseWheelToDispose.length > 0 !== e && (this._mouseWheelToDispose = Ne2(this._mouseWheelToDispose), e)) {
      let r5 = (n) => {
        this._onMouseWheel(new xi2(n));
      };
      this._mouseWheelToDispose.push(L(this._listenOnDomNode, Y3.MOUSE_WHEEL, r5, { passive: false }));
    }
  }
  _onMouseWheel(e) {
    if (e.browserEvent?.defaultPrevented) return;
    let i8 = hs2.INSTANCE;
    zo && i8.acceptStandardWheelEvent(e);
    let r5 = false;
    if (e.deltaY || e.deltaX) {
      let o2 = e.deltaY * this._options.mouseWheelScrollSensitivity, l = e.deltaX * this._options.mouseWheelScrollSensitivity;
      this._options.scrollPredominantAxis && (this._options.scrollYToX && l + o2 === 0 ? l = o2 = 0 : Math.abs(o2) >= Math.abs(l) ? l = 0 : o2 = 0), this._options.flipAxes && ([o2, l] = [l, o2]);
      let a = !Te3 && e.browserEvent && e.browserEvent.shiftKey;
      (this._options.scrollYToX || a) && !l && (l = o2, o2 = 0), e.browserEvent && e.browserEvent.altKey && (l = l * this._options.fastScrollSensitivity, o2 = o2 * this._options.fastScrollSensitivity);
      let u = this._scrollable.getFutureScrollPosition(), h2 = {};
      if (o2) {
        let c = Ko * o2, d = u.scrollTop - (c < 0 ? Math.floor(c) : Math.ceil(c));
        this._verticalScrollbar.writeScrollPosition(h2, d);
      }
      if (l) {
        let c = Ko * l, d = u.scrollLeft - (c < 0 ? Math.floor(c) : Math.ceil(c));
        this._horizontalScrollbar.writeScrollPosition(h2, d);
      }
      h2 = this._scrollable.validateScrollPosition(h2), (u.scrollLeft !== h2.scrollLeft || u.scrollTop !== h2.scrollTop) && (zo && this._options.mouseWheelSmoothScroll && i8.isPhysicalMouseWheel() ? this._scrollable.setScrollPositionSmooth(h2) : this._scrollable.setScrollPositionNow(h2), r5 = true);
    }
    let n = r5;
    !n && this._options.alwaysConsumeMouseWheel && (n = true), !n && this._options.consumeMouseWheelIfScrollbarIsNeeded && (this._verticalScrollbar.isNeeded() || this._horizontalScrollbar.isNeeded()) && (n = true), n && (e.preventDefault(), e.stopPropagation());
  }
  _onDidScroll(e) {
    this._shouldRender = this._horizontalScrollbar.onDidScroll(e) || this._shouldRender, this._shouldRender = this._verticalScrollbar.onDidScroll(e) || this._shouldRender, this._options.useShadows && (this._shouldRender = true), this._revealOnScroll && this._reveal(), this._options.lazyRender || this._render();
  }
  renderNow() {
    if (!this._options.lazyRender) throw new Error("Please use `lazyRender` together with `renderNow`!");
    this._render();
  }
  _render() {
    if (this._shouldRender && (this._shouldRender = false, this._horizontalScrollbar.render(), this._verticalScrollbar.render(), this._options.useShadows)) {
      let e = this._scrollable.getCurrentScrollPosition(), i8 = e.scrollTop > 0, r5 = e.scrollLeft > 0, n = r5 ? " left" : "", o2 = i8 ? " top" : "", l = r5 || i8 ? " top-left-corner" : "";
      this._leftShadowDomNode.setClassName(`shadow${n}`), this._topShadowDomNode.setClassName(`shadow${o2}`), this._topLeftShadowDomNode.setClassName(`shadow${l}${o2}${n}`);
    }
  }
  _onDragStart() {
    this._isDragging = true, this._reveal();
  }
  _onDragEnd() {
    this._isDragging = false, this._hide();
  }
  _onMouseLeave(e) {
    this._mouseIsOver = false, this._hide();
  }
  _onMouseOver(e) {
    this._mouseIsOver = true, this._reveal();
  }
  _reveal() {
    this._verticalScrollbar.beginReveal(), this._horizontalScrollbar.beginReveal(), this._scheduleHide();
  }
  _hide() {
    !this._mouseIsOver && !this._isDragging && (this._verticalScrollbar.beginHide(), this._horizontalScrollbar.beginHide());
  }
  _scheduleHide() {
    !this._mouseIsOver && !this._isDragging && this._hideTimeout.cancelAndSet(() => this._hide(), Ma);
  }
};
var Kr2 = class extends ds2 {
  constructor(t, e, i8) {
    super(t, e, i8);
  }
  setScrollPosition(t) {
    t.reuseAnimation ? this._scrollable.setScrollPositionSmooth(t, t.reuseAnimation) : this._scrollable.setScrollPositionNow(t);
  }
  getScrollPosition() {
    return this._scrollable.getCurrentScrollPosition();
  }
};
function Pa(s15) {
  let t = { lazyRender: typeof s15.lazyRender < "u" ? s15.lazyRender : false, className: typeof s15.className < "u" ? s15.className : "", useShadows: typeof s15.useShadows < "u" ? s15.useShadows : true, handleMouseWheel: typeof s15.handleMouseWheel < "u" ? s15.handleMouseWheel : true, flipAxes: typeof s15.flipAxes < "u" ? s15.flipAxes : false, consumeMouseWheelIfScrollbarIsNeeded: typeof s15.consumeMouseWheelIfScrollbarIsNeeded < "u" ? s15.consumeMouseWheelIfScrollbarIsNeeded : false, alwaysConsumeMouseWheel: typeof s15.alwaysConsumeMouseWheel < "u" ? s15.alwaysConsumeMouseWheel : false, scrollYToX: typeof s15.scrollYToX < "u" ? s15.scrollYToX : false, mouseWheelScrollSensitivity: typeof s15.mouseWheelScrollSensitivity < "u" ? s15.mouseWheelScrollSensitivity : 1, fastScrollSensitivity: typeof s15.fastScrollSensitivity < "u" ? s15.fastScrollSensitivity : 5, scrollPredominantAxis: typeof s15.scrollPredominantAxis < "u" ? s15.scrollPredominantAxis : true, mouseWheelSmoothScroll: typeof s15.mouseWheelSmoothScroll < "u" ? s15.mouseWheelSmoothScroll : true, arrowSize: typeof s15.arrowSize < "u" ? s15.arrowSize : 11, listenOnDomNode: typeof s15.listenOnDomNode < "u" ? s15.listenOnDomNode : null, horizontal: typeof s15.horizontal < "u" ? s15.horizontal : 1, horizontalScrollbarSize: typeof s15.horizontalScrollbarSize < "u" ? s15.horizontalScrollbarSize : 10, horizontalSliderSize: typeof s15.horizontalSliderSize < "u" ? s15.horizontalSliderSize : 0, horizontalHasArrows: typeof s15.horizontalHasArrows < "u" ? s15.horizontalHasArrows : false, vertical: typeof s15.vertical < "u" ? s15.vertical : 1, verticalScrollbarSize: typeof s15.verticalScrollbarSize < "u" ? s15.verticalScrollbarSize : 10, verticalHasArrows: typeof s15.verticalHasArrows < "u" ? s15.verticalHasArrows : false, verticalSliderSize: typeof s15.verticalSliderSize < "u" ? s15.verticalSliderSize : 0, scrollByPage: typeof s15.scrollByPage < "u" ? s15.scrollByPage : false };
  return t.horizontalSliderSize = typeof s15.horizontalSliderSize < "u" ? s15.horizontalSliderSize : t.horizontalScrollbarSize, t.verticalSliderSize = typeof s15.verticalSliderSize < "u" ? s15.verticalSliderSize : t.verticalScrollbarSize, Te3 && (t.className += " mac"), t;
}
var zt2 = class extends D2 {
  constructor(e, i8, r5, n, o2, l, a, u) {
    super();
    this._bufferService = r5;
    this._optionsService = a;
    this._renderService = u;
    this._onRequestScrollLines = this._register(new v());
    this.onRequestScrollLines = this._onRequestScrollLines.event;
    this._isSyncing = false;
    this._isHandlingScroll = false;
    this._suppressOnScrollHandler = false;
    let h2 = this._register(new Ri({ forceIntegerValues: false, smoothScrollDuration: this._optionsService.rawOptions.smoothScrollDuration, scheduleAtNextAnimationFrame: (c) => mt(n.window, c) }));
    this._register(this._optionsService.onSpecificOptionChange("smoothScrollDuration", () => {
      h2.setSmoothScrollDuration(this._optionsService.rawOptions.smoothScrollDuration);
    })), this._scrollableElement = this._register(new Kr2(i8, { vertical: 1, horizontal: 2, useShadows: false, mouseWheelSmoothScroll: true, ...this._getChangeOptions() }, h2)), this._register(this._optionsService.onMultipleOptionChange(["scrollSensitivity", "fastScrollSensitivity", "overviewRuler"], () => this._scrollableElement.updateOptions(this._getChangeOptions()))), this._register(o2.onProtocolChange((c) => {
      this._scrollableElement.updateOptions({ handleMouseWheel: !(c & 16) });
    })), this._scrollableElement.setScrollDimensions({ height: 0, scrollHeight: 0 }), this._register($.runAndSubscribe(l.onChangeColors, () => {
      this._scrollableElement.getDomNode().style.backgroundColor = l.colors.background.css;
    })), e.appendChild(this._scrollableElement.getDomNode()), this._register(C3(() => this._scrollableElement.getDomNode().remove())), this._styleElement = n.mainDocument.createElement("style"), i8.appendChild(this._styleElement), this._register(C3(() => this._styleElement.remove())), this._register($.runAndSubscribe(l.onChangeColors, () => {
      this._styleElement.textContent = [".xterm .xterm-scrollable-element > .scrollbar > .slider {", `  background: ${l.colors.scrollbarSliderBackground.css};`, "}", ".xterm .xterm-scrollable-element > .scrollbar > .slider:hover {", `  background: ${l.colors.scrollbarSliderHoverBackground.css};`, "}", ".xterm .xterm-scrollable-element > .scrollbar > .slider.active {", `  background: ${l.colors.scrollbarSliderActiveBackground.css};`, "}"].join(`
`);
    })), this._register(this._bufferService.onResize(() => this.queueSync())), this._register(this._bufferService.buffers.onBufferActivate(() => {
      this._latestYDisp = void 0, this.queueSync();
    })), this._register(this._bufferService.onScroll(() => this._sync())), this._register(this._scrollableElement.onScroll((c) => this._handleScroll(c)));
  }
  scrollLines(e) {
    let i8 = this._scrollableElement.getScrollPosition();
    this._scrollableElement.setScrollPosition({ reuseAnimation: true, scrollTop: i8.scrollTop + e * this._renderService.dimensions.css.cell.height });
  }
  scrollToLine(e, i8) {
    i8 && (this._latestYDisp = e), this._scrollableElement.setScrollPosition({ reuseAnimation: !i8, scrollTop: e * this._renderService.dimensions.css.cell.height });
  }
  _getChangeOptions() {
    return { mouseWheelScrollSensitivity: this._optionsService.rawOptions.scrollSensitivity, fastScrollSensitivity: this._optionsService.rawOptions.fastScrollSensitivity, verticalScrollbarSize: this._optionsService.rawOptions.overviewRuler?.width || 14 };
  }
  queueSync(e) {
    e !== void 0 && (this._latestYDisp = e), this._queuedAnimationFrame === void 0 && (this._queuedAnimationFrame = this._renderService.addRefreshCallback(() => {
      this._queuedAnimationFrame = void 0, this._sync(this._latestYDisp);
    }));
  }
  _sync(e = this._bufferService.buffer.ydisp) {
    !this._renderService || this._isSyncing || (this._isSyncing = true, this._suppressOnScrollHandler = true, this._scrollableElement.setScrollDimensions({ height: this._renderService.dimensions.css.canvas.height, scrollHeight: this._renderService.dimensions.css.cell.height * this._bufferService.buffer.lines.length }), this._suppressOnScrollHandler = false, e !== this._latestYDisp && this._scrollableElement.setScrollPosition({ scrollTop: e * this._renderService.dimensions.css.cell.height }), this._isSyncing = false);
  }
  _handleScroll(e) {
    if (!this._renderService || this._isHandlingScroll || this._suppressOnScrollHandler) return;
    this._isHandlingScroll = true;
    let i8 = Math.round(e.scrollTop / this._renderService.dimensions.css.cell.height), r5 = i8 - this._bufferService.buffer.ydisp;
    r5 !== 0 && (this._latestYDisp = i8, this._onRequestScrollLines.fire(r5)), this._isHandlingScroll = false;
  }
};
zt2 = M3([S(2, F2), S(3, ae2), S(4, rr), S(5, Re2), S(6, H3), S(7, ce2)], zt2);
var Gt2 = class extends D2 {
  constructor(e, i8, r5, n, o2) {
    super();
    this._screenElement = e;
    this._bufferService = i8;
    this._coreBrowserService = r5;
    this._decorationService = n;
    this._renderService = o2;
    this._decorationElements = /* @__PURE__ */ new Map();
    this._altBufferIsActive = false;
    this._dimensionsChanged = false;
    this._container = document.createElement("div"), this._container.classList.add("xterm-decoration-container"), this._screenElement.appendChild(this._container), this._register(this._renderService.onRenderedViewportChange(() => this._doRefreshDecorations())), this._register(this._renderService.onDimensionsChange(() => {
      this._dimensionsChanged = true, this._queueRefresh();
    })), this._register(this._coreBrowserService.onDprChange(() => this._queueRefresh())), this._register(this._bufferService.buffers.onBufferActivate(() => {
      this._altBufferIsActive = this._bufferService.buffer === this._bufferService.buffers.alt;
    })), this._register(this._decorationService.onDecorationRegistered(() => this._queueRefresh())), this._register(this._decorationService.onDecorationRemoved((l) => this._removeDecoration(l))), this._register(C3(() => {
      this._container.remove(), this._decorationElements.clear();
    }));
  }
  _queueRefresh() {
    this._animationFrame === void 0 && (this._animationFrame = this._renderService.addRefreshCallback(() => {
      this._doRefreshDecorations(), this._animationFrame = void 0;
    }));
  }
  _doRefreshDecorations() {
    for (let e of this._decorationService.decorations) this._renderDecoration(e);
    this._dimensionsChanged = false;
  }
  _renderDecoration(e) {
    this._refreshStyle(e), this._dimensionsChanged && this._refreshXPosition(e);
  }
  _createElement(e) {
    let i8 = this._coreBrowserService.mainDocument.createElement("div");
    i8.classList.add("xterm-decoration"), i8.classList.toggle("xterm-decoration-top-layer", e?.options?.layer === "top"), i8.style.width = `${Math.round((e.options.width || 1) * this._renderService.dimensions.css.cell.width)}px`, i8.style.height = `${(e.options.height || 1) * this._renderService.dimensions.css.cell.height}px`, i8.style.top = `${(e.marker.line - this._bufferService.buffers.active.ydisp) * this._renderService.dimensions.css.cell.height}px`, i8.style.lineHeight = `${this._renderService.dimensions.css.cell.height}px`;
    let r5 = e.options.x ?? 0;
    return r5 && r5 > this._bufferService.cols && (i8.style.display = "none"), this._refreshXPosition(e, i8), i8;
  }
  _refreshStyle(e) {
    let i8 = e.marker.line - this._bufferService.buffers.active.ydisp;
    if (i8 < 0 || i8 >= this._bufferService.rows) e.element && (e.element.style.display = "none", e.onRenderEmitter.fire(e.element));
    else {
      let r5 = this._decorationElements.get(e);
      r5 || (r5 = this._createElement(e), e.element = r5, this._decorationElements.set(e, r5), this._container.appendChild(r5), e.onDispose(() => {
        this._decorationElements.delete(e), r5.remove();
      })), r5.style.display = this._altBufferIsActive ? "none" : "block", this._altBufferIsActive || (r5.style.width = `${Math.round((e.options.width || 1) * this._renderService.dimensions.css.cell.width)}px`, r5.style.height = `${(e.options.height || 1) * this._renderService.dimensions.css.cell.height}px`, r5.style.top = `${i8 * this._renderService.dimensions.css.cell.height}px`, r5.style.lineHeight = `${this._renderService.dimensions.css.cell.height}px`), e.onRenderEmitter.fire(r5);
    }
  }
  _refreshXPosition(e, i8 = e.element) {
    if (!i8) return;
    let r5 = e.options.x ?? 0;
    (e.options.anchor || "left") === "right" ? i8.style.right = r5 ? `${r5 * this._renderService.dimensions.css.cell.width}px` : "" : i8.style.left = r5 ? `${r5 * this._renderService.dimensions.css.cell.width}px` : "";
  }
  _removeDecoration(e) {
    this._decorationElements.get(e)?.remove(), this._decorationElements.delete(e), e.dispose();
  }
};
Gt2 = M3([S(1, F2), S(2, ae2), S(3, Be4), S(4, ce2)], Gt2);
var Gr2 = class {
  constructor() {
    this._zones = [];
    this._zonePool = [];
    this._zonePoolIndex = 0;
    this._linePadding = { full: 0, left: 0, center: 0, right: 0 };
  }
  get zones() {
    return this._zonePool.length = Math.min(this._zonePool.length, this._zones.length), this._zones;
  }
  clear() {
    this._zones.length = 0, this._zonePoolIndex = 0;
  }
  addDecoration(t) {
    if (t.options.overviewRulerOptions) {
      for (let e of this._zones) if (e.color === t.options.overviewRulerOptions.color && e.position === t.options.overviewRulerOptions.position) {
        if (this._lineIntersectsZone(e, t.marker.line)) return;
        if (this._lineAdjacentToZone(e, t.marker.line, t.options.overviewRulerOptions.position)) {
          this._addLineToZone(e, t.marker.line);
          return;
        }
      }
      if (this._zonePoolIndex < this._zonePool.length) {
        this._zonePool[this._zonePoolIndex].color = t.options.overviewRulerOptions.color, this._zonePool[this._zonePoolIndex].position = t.options.overviewRulerOptions.position, this._zonePool[this._zonePoolIndex].startBufferLine = t.marker.line, this._zonePool[this._zonePoolIndex].endBufferLine = t.marker.line, this._zones.push(this._zonePool[this._zonePoolIndex++]);
        return;
      }
      this._zones.push({ color: t.options.overviewRulerOptions.color, position: t.options.overviewRulerOptions.position, startBufferLine: t.marker.line, endBufferLine: t.marker.line }), this._zonePool.push(this._zones[this._zones.length - 1]), this._zonePoolIndex++;
    }
  }
  setPadding(t) {
    this._linePadding = t;
  }
  _lineIntersectsZone(t, e) {
    return e >= t.startBufferLine && e <= t.endBufferLine;
  }
  _lineAdjacentToZone(t, e, i8) {
    return e >= t.startBufferLine - this._linePadding[i8 || "full"] && e <= t.endBufferLine + this._linePadding[i8 || "full"];
  }
  _addLineToZone(t, e) {
    t.startBufferLine = Math.min(t.startBufferLine, e), t.endBufferLine = Math.max(t.endBufferLine, e);
  }
};
var We4 = { full: 0, left: 0, center: 0, right: 0 };
var at3 = { full: 0, left: 0, center: 0, right: 0 };
var Li2 = { full: 0, left: 0, center: 0, right: 0 };
var bt2 = class extends D2 {
  constructor(e, i8, r5, n, o2, l, a, u) {
    super();
    this._viewportElement = e;
    this._screenElement = i8;
    this._bufferService = r5;
    this._decorationService = n;
    this._renderService = o2;
    this._optionsService = l;
    this._themeService = a;
    this._coreBrowserService = u;
    this._colorZoneStore = new Gr2();
    this._shouldUpdateDimensions = true;
    this._shouldUpdateAnchor = true;
    this._lastKnownBufferLength = 0;
    this._canvas = this._coreBrowserService.mainDocument.createElement("canvas"), this._canvas.classList.add("xterm-decoration-overview-ruler"), this._refreshCanvasDimensions(), this._viewportElement.parentElement?.insertBefore(this._canvas, this._viewportElement), this._register(C3(() => this._canvas?.remove()));
    let h2 = this._canvas.getContext("2d");
    if (h2) this._ctx = h2;
    else throw new Error("Ctx cannot be null");
    this._register(this._decorationService.onDecorationRegistered(() => this._queueRefresh(void 0, true))), this._register(this._decorationService.onDecorationRemoved(() => this._queueRefresh(void 0, true))), this._register(this._renderService.onRenderedViewportChange(() => this._queueRefresh())), this._register(this._bufferService.buffers.onBufferActivate(() => {
      this._canvas.style.display = this._bufferService.buffer === this._bufferService.buffers.alt ? "none" : "block";
    })), this._register(this._bufferService.onScroll(() => {
      this._lastKnownBufferLength !== this._bufferService.buffers.normal.lines.length && (this._refreshDrawHeightConstants(), this._refreshColorZonePadding());
    })), this._register(this._renderService.onRender(() => {
      (!this._containerHeight || this._containerHeight !== this._screenElement.clientHeight) && (this._queueRefresh(true), this._containerHeight = this._screenElement.clientHeight);
    })), this._register(this._coreBrowserService.onDprChange(() => this._queueRefresh(true))), this._register(this._optionsService.onSpecificOptionChange("overviewRuler", () => this._queueRefresh(true))), this._register(this._themeService.onChangeColors(() => this._queueRefresh())), this._queueRefresh(true);
  }
  get _width() {
    return this._optionsService.options.overviewRuler?.width || 0;
  }
  _refreshDrawConstants() {
    let e = Math.floor((this._canvas.width - 1) / 3), i8 = Math.ceil((this._canvas.width - 1) / 3);
    at3.full = this._canvas.width, at3.left = e, at3.center = i8, at3.right = e, this._refreshDrawHeightConstants(), Li2.full = 1, Li2.left = 1, Li2.center = 1 + at3.left, Li2.right = 1 + at3.left + at3.center;
  }
  _refreshDrawHeightConstants() {
    We4.full = Math.round(2 * this._coreBrowserService.dpr);
    let e = this._canvas.height / this._bufferService.buffer.lines.length, i8 = Math.round(Math.max(Math.min(e, 12), 6) * this._coreBrowserService.dpr);
    We4.left = i8, We4.center = i8, We4.right = i8;
  }
  _refreshColorZonePadding() {
    this._colorZoneStore.setPadding({ full: Math.floor(this._bufferService.buffers.active.lines.length / (this._canvas.height - 1) * We4.full), left: Math.floor(this._bufferService.buffers.active.lines.length / (this._canvas.height - 1) * We4.left), center: Math.floor(this._bufferService.buffers.active.lines.length / (this._canvas.height - 1) * We4.center), right: Math.floor(this._bufferService.buffers.active.lines.length / (this._canvas.height - 1) * We4.right) }), this._lastKnownBufferLength = this._bufferService.buffers.normal.lines.length;
  }
  _refreshCanvasDimensions() {
    this._canvas.style.width = `${this._width}px`, this._canvas.width = Math.round(this._width * this._coreBrowserService.dpr), this._canvas.style.height = `${this._screenElement.clientHeight}px`, this._canvas.height = Math.round(this._screenElement.clientHeight * this._coreBrowserService.dpr), this._refreshDrawConstants(), this._refreshColorZonePadding();
  }
  _refreshDecorations() {
    this._shouldUpdateDimensions && this._refreshCanvasDimensions(), this._ctx.clearRect(0, 0, this._canvas.width, this._canvas.height), this._colorZoneStore.clear();
    for (let i8 of this._decorationService.decorations) this._colorZoneStore.addDecoration(i8);
    this._ctx.lineWidth = 1, this._renderRulerOutline();
    let e = this._colorZoneStore.zones;
    for (let i8 of e) i8.position !== "full" && this._renderColorZone(i8);
    for (let i8 of e) i8.position === "full" && this._renderColorZone(i8);
    this._shouldUpdateDimensions = false, this._shouldUpdateAnchor = false;
  }
  _renderRulerOutline() {
    this._ctx.fillStyle = this._themeService.colors.overviewRulerBorder.css, this._ctx.fillRect(0, 0, 1, this._canvas.height), this._optionsService.rawOptions.overviewRuler.showTopBorder && this._ctx.fillRect(1, 0, this._canvas.width - 1, 1), this._optionsService.rawOptions.overviewRuler.showBottomBorder && this._ctx.fillRect(1, this._canvas.height - 1, this._canvas.width - 1, this._canvas.height);
  }
  _renderColorZone(e) {
    this._ctx.fillStyle = e.color, this._ctx.fillRect(Li2[e.position || "full"], Math.round((this._canvas.height - 1) * (e.startBufferLine / this._bufferService.buffers.active.lines.length) - We4[e.position || "full"] / 2), at3[e.position || "full"], Math.round((this._canvas.height - 1) * ((e.endBufferLine - e.startBufferLine) / this._bufferService.buffers.active.lines.length) + We4[e.position || "full"]));
  }
  _queueRefresh(e, i8) {
    this._shouldUpdateDimensions = e || this._shouldUpdateDimensions, this._shouldUpdateAnchor = i8 || this._shouldUpdateAnchor, this._animationFrame === void 0 && (this._animationFrame = this._coreBrowserService.window.requestAnimationFrame(() => {
      this._refreshDecorations(), this._animationFrame = void 0;
    }));
  }
};
bt2 = M3([S(2, F2), S(3, Be4), S(4, ce2), S(5, H3), S(6, Re2), S(7, ae2)], bt2);
var b;
((E) => (E.NUL = "\0", E.SOH = "", E.STX = "", E.ETX = "", E.EOT = "", E.ENQ = "", E.ACK = "", E.BEL = "\x07", E.BS = "\b", E.HT = "	", E.LF = `
`, E.VT = "\v", E.FF = "\f", E.CR = "\r", E.SO = "", E.SI = "", E.DLE = "", E.DC1 = "", E.DC2 = "", E.DC3 = "", E.DC4 = "", E.NAK = "", E.SYN = "", E.ETB = "", E.CAN = "", E.EM = "", E.SUB = "", E.ESC = "\x1B", E.FS = "", E.GS = "", E.RS = "", E.US = "", E.SP = " ", E.DEL = ""))(b ||= {});
var Ai2;
((g) => (g.PAD = "", g.HOP = "", g.BPH = "", g.NBH = "", g.IND = "", g.NEL = "", g.SSA = "", g.ESA = "", g.HTS = "", g.HTJ = "", g.VTS = "", g.PLD = "", g.PLU = "", g.RI = "", g.SS2 = "", g.SS3 = "", g.DCS = "", g.PU1 = "", g.PU2 = "", g.STS = "", g.CCH = "", g.MW = "", g.SPA = "", g.EPA = "", g.SOS = "", g.SGCI = "", g.SCI = "", g.CSI = "", g.ST = "", g.OSC = "", g.PM = "", g.APC = ""))(Ai2 ||= {});
var fs2;
((t) => t.ST = `${b.ESC}\\`)(fs2 ||= {});
var $t2 = class {
  constructor(t, e, i8, r5, n, o2) {
    this._textarea = t;
    this._compositionView = e;
    this._bufferService = i8;
    this._optionsService = r5;
    this._coreService = n;
    this._renderService = o2;
    this._isComposing = false, this._isSendingComposition = false, this._compositionPosition = { start: 0, end: 0 }, this._dataAlreadySent = "";
  }
  get isComposing() {
    return this._isComposing;
  }
  compositionstart() {
    this._isComposing = true, this._compositionPosition.start = this._textarea.value.length, this._compositionView.textContent = "", this._dataAlreadySent = "", this._compositionView.classList.add("active");
  }
  compositionupdate(t) {
    this._compositionView.textContent = t.data, this.updateCompositionElements(), setTimeout(() => {
      this._compositionPosition.end = this._textarea.value.length;
    }, 0);
  }
  compositionend() {
    this._finalizeComposition(true);
  }
  keydown(t) {
    if (this._isComposing || this._isSendingComposition) {
      if (t.keyCode === 20 || t.keyCode === 229 || t.keyCode === 16 || t.keyCode === 17 || t.keyCode === 18) return false;
      this._finalizeComposition(false);
    }
    return t.keyCode === 229 ? (this._handleAnyTextareaChanges(), false) : true;
  }
  _finalizeComposition(t) {
    if (this._compositionView.classList.remove("active"), this._isComposing = false, t) {
      let e = { start: this._compositionPosition.start, end: this._compositionPosition.end };
      this._isSendingComposition = true, setTimeout(() => {
        if (this._isSendingComposition) {
          this._isSendingComposition = false;
          let i8;
          e.start += this._dataAlreadySent.length, this._isComposing ? i8 = this._textarea.value.substring(e.start, this._compositionPosition.start) : i8 = this._textarea.value.substring(e.start), i8.length > 0 && this._coreService.triggerDataEvent(i8, true);
        }
      }, 0);
    } else {
      this._isSendingComposition = false;
      let e = this._textarea.value.substring(this._compositionPosition.start, this._compositionPosition.end);
      this._coreService.triggerDataEvent(e, true);
    }
  }
  _handleAnyTextareaChanges() {
    let t = this._textarea.value;
    setTimeout(() => {
      if (!this._isComposing) {
        let e = this._textarea.value, i8 = e.replace(t, "");
        this._dataAlreadySent = i8, e.length > t.length ? this._coreService.triggerDataEvent(i8, true) : e.length < t.length ? this._coreService.triggerDataEvent(`${b.DEL}`, true) : e.length === t.length && e !== t && this._coreService.triggerDataEvent(e, true);
      }
    }, 0);
  }
  updateCompositionElements(t) {
    if (this._isComposing) {
      if (this._bufferService.buffer.isCursorInViewport) {
        let e = Math.min(this._bufferService.buffer.x, this._bufferService.cols - 1), i8 = this._renderService.dimensions.css.cell.height, r5 = this._bufferService.buffer.y * this._renderService.dimensions.css.cell.height, n = e * this._renderService.dimensions.css.cell.width;
        this._compositionView.style.left = n + "px", this._compositionView.style.top = r5 + "px", this._compositionView.style.height = i8 + "px", this._compositionView.style.lineHeight = i8 + "px", this._compositionView.style.fontFamily = this._optionsService.rawOptions.fontFamily, this._compositionView.style.fontSize = this._optionsService.rawOptions.fontSize + "px";
        let o2 = this._compositionView.getBoundingClientRect();
        this._textarea.style.left = n + "px", this._textarea.style.top = r5 + "px", this._textarea.style.width = Math.max(o2.width, 1) + "px", this._textarea.style.height = Math.max(o2.height, 1) + "px", this._textarea.style.lineHeight = o2.height + "px";
      }
      t || setTimeout(() => this.updateCompositionElements(true), 0);
    }
  }
};
$t2 = M3([S(2, F2), S(3, H3), S(4, ge3), S(5, ce2)], $t2);
var ue2 = 0;
var he3 = 0;
var de2 = 0;
var J3 = 0;
var ps2 = { css: "#00000000", rgba: 0 };
var j2;
((i8) => {
  function s15(r5, n, o2, l) {
    return l !== void 0 ? `#${vt3(r5)}${vt3(n)}${vt3(o2)}${vt3(l)}` : `#${vt3(r5)}${vt3(n)}${vt3(o2)}`;
  }
  i8.toCss = s15;
  function t(r5, n, o2, l = 255) {
    return (r5 << 24 | n << 16 | o2 << 8 | l) >>> 0;
  }
  i8.toRgba = t;
  function e(r5, n, o2, l) {
    return { css: i8.toCss(r5, n, o2, l), rgba: i8.toRgba(r5, n, o2, l) };
  }
  i8.toColor = e;
})(j2 ||= {});
var U3;
((l) => {
  function s15(a, u) {
    if (J3 = (u.rgba & 255) / 255, J3 === 1) return { css: u.css, rgba: u.rgba };
    let h2 = u.rgba >> 24 & 255, c = u.rgba >> 16 & 255, d = u.rgba >> 8 & 255, _3 = a.rgba >> 24 & 255, p = a.rgba >> 16 & 255, m = a.rgba >> 8 & 255;
    ue2 = _3 + Math.round((h2 - _3) * J3), he3 = p + Math.round((c - p) * J3), de2 = m + Math.round((d - m) * J3);
    let f = j2.toCss(ue2, he3, de2), A2 = j2.toRgba(ue2, he3, de2);
    return { css: f, rgba: A2 };
  }
  l.blend = s15;
  function t(a) {
    return (a.rgba & 255) === 255;
  }
  l.isOpaque = t;
  function e(a, u, h2) {
    let c = $r2.ensureContrastRatio(a.rgba, u.rgba, h2);
    if (c) return j2.toColor(c >> 24 & 255, c >> 16 & 255, c >> 8 & 255);
  }
  l.ensureContrastRatio = e;
  function i8(a) {
    let u = (a.rgba | 255) >>> 0;
    return [ue2, he3, de2] = $r2.toChannels(u), { css: j2.toCss(ue2, he3, de2), rgba: u };
  }
  l.opaque = i8;
  function r5(a, u) {
    return J3 = Math.round(u * 255), [ue2, he3, de2] = $r2.toChannels(a.rgba), { css: j2.toCss(ue2, he3, de2, J3), rgba: j2.toRgba(ue2, he3, de2, J3) };
  }
  l.opacity = r5;
  function n(a, u) {
    return J3 = a.rgba & 255, r5(a, J3 * u / 255);
  }
  l.multiplyOpacity = n;
  function o2(a) {
    return [a.rgba >> 24 & 255, a.rgba >> 16 & 255, a.rgba >> 8 & 255];
  }
  l.toColorRGB = o2;
})(U3 ||= {});
var z2;
((i8) => {
  let s15, t;
  try {
    let r5 = document.createElement("canvas");
    r5.width = 1, r5.height = 1;
    let n = r5.getContext("2d", { willReadFrequently: true });
    n && (s15 = n, s15.globalCompositeOperation = "copy", t = s15.createLinearGradient(0, 0, 1, 1));
  } catch {
  }
  function e(r5) {
    if (r5.match(/#[\da-f]{3,8}/i)) switch (r5.length) {
      case 4:
        return ue2 = parseInt(r5.slice(1, 2).repeat(2), 16), he3 = parseInt(r5.slice(2, 3).repeat(2), 16), de2 = parseInt(r5.slice(3, 4).repeat(2), 16), j2.toColor(ue2, he3, de2);
      case 5:
        return ue2 = parseInt(r5.slice(1, 2).repeat(2), 16), he3 = parseInt(r5.slice(2, 3).repeat(2), 16), de2 = parseInt(r5.slice(3, 4).repeat(2), 16), J3 = parseInt(r5.slice(4, 5).repeat(2), 16), j2.toColor(ue2, he3, de2, J3);
      case 7:
        return { css: r5, rgba: (parseInt(r5.slice(1), 16) << 8 | 255) >>> 0 };
      case 9:
        return { css: r5, rgba: parseInt(r5.slice(1), 16) >>> 0 };
    }
    let n = r5.match(/rgba?\(\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*(\d{1,3})\s*(,\s*(0|1|\d?\.(\d+))\s*)?\)/);
    if (n) return ue2 = parseInt(n[1]), he3 = parseInt(n[2]), de2 = parseInt(n[3]), J3 = Math.round((n[5] === void 0 ? 1 : parseFloat(n[5])) * 255), j2.toColor(ue2, he3, de2, J3);
    if (!s15 || !t) throw new Error("css.toColor: Unsupported css format");
    if (s15.fillStyle = t, s15.fillStyle = r5, typeof s15.fillStyle != "string") throw new Error("css.toColor: Unsupported css format");
    if (s15.fillRect(0, 0, 1, 1), [ue2, he3, de2, J3] = s15.getImageData(0, 0, 1, 1).data, J3 !== 255) throw new Error("css.toColor: Unsupported css format");
    return { rgba: j2.toRgba(ue2, he3, de2, J3), css: r5 };
  }
  i8.toColor = e;
})(z2 ||= {});
var ve3;
((e) => {
  function s15(i8) {
    return t(i8 >> 16 & 255, i8 >> 8 & 255, i8 & 255);
  }
  e.relativeLuminance = s15;
  function t(i8, r5, n) {
    let o2 = i8 / 255, l = r5 / 255, a = n / 255, u = o2 <= 0.03928 ? o2 / 12.92 : Math.pow((o2 + 0.055) / 1.055, 2.4), h2 = l <= 0.03928 ? l / 12.92 : Math.pow((l + 0.055) / 1.055, 2.4), c = a <= 0.03928 ? a / 12.92 : Math.pow((a + 0.055) / 1.055, 2.4);
    return u * 0.2126 + h2 * 0.7152 + c * 0.0722;
  }
  e.relativeLuminance2 = t;
})(ve3 ||= {});
var $r2;
((n) => {
  function s15(o2, l) {
    if (J3 = (l & 255) / 255, J3 === 1) return l;
    let a = l >> 24 & 255, u = l >> 16 & 255, h2 = l >> 8 & 255, c = o2 >> 24 & 255, d = o2 >> 16 & 255, _3 = o2 >> 8 & 255;
    return ue2 = c + Math.round((a - c) * J3), he3 = d + Math.round((u - d) * J3), de2 = _3 + Math.round((h2 - _3) * J3), j2.toRgba(ue2, he3, de2);
  }
  n.blend = s15;
  function t(o2, l, a) {
    let u = ve3.relativeLuminance(o2 >> 8), h2 = ve3.relativeLuminance(l >> 8);
    if (Xe3(u, h2) < a) {
      if (h2 < u) {
        let p = e(o2, l, a), m = Xe3(u, ve3.relativeLuminance(p >> 8));
        if (m < a) {
          let f = i8(o2, l, a), A2 = Xe3(u, ve3.relativeLuminance(f >> 8));
          return m > A2 ? p : f;
        }
        return p;
      }
      let d = i8(o2, l, a), _3 = Xe3(u, ve3.relativeLuminance(d >> 8));
      if (_3 < a) {
        let p = e(o2, l, a), m = Xe3(u, ve3.relativeLuminance(p >> 8));
        return _3 > m ? d : p;
      }
      return d;
    }
  }
  n.ensureContrastRatio = t;
  function e(o2, l, a) {
    let u = o2 >> 24 & 255, h2 = o2 >> 16 & 255, c = o2 >> 8 & 255, d = l >> 24 & 255, _3 = l >> 16 & 255, p = l >> 8 & 255, m = Xe3(ve3.relativeLuminance2(d, _3, p), ve3.relativeLuminance2(u, h2, c));
    for (; m < a && (d > 0 || _3 > 0 || p > 0); ) d -= Math.max(0, Math.ceil(d * 0.1)), _3 -= Math.max(0, Math.ceil(_3 * 0.1)), p -= Math.max(0, Math.ceil(p * 0.1)), m = Xe3(ve3.relativeLuminance2(d, _3, p), ve3.relativeLuminance2(u, h2, c));
    return (d << 24 | _3 << 16 | p << 8 | 255) >>> 0;
  }
  n.reduceLuminance = e;
  function i8(o2, l, a) {
    let u = o2 >> 24 & 255, h2 = o2 >> 16 & 255, c = o2 >> 8 & 255, d = l >> 24 & 255, _3 = l >> 16 & 255, p = l >> 8 & 255, m = Xe3(ve3.relativeLuminance2(d, _3, p), ve3.relativeLuminance2(u, h2, c));
    for (; m < a && (d < 255 || _3 < 255 || p < 255); ) d = Math.min(255, d + Math.ceil((255 - d) * 0.1)), _3 = Math.min(255, _3 + Math.ceil((255 - _3) * 0.1)), p = Math.min(255, p + Math.ceil((255 - p) * 0.1)), m = Xe3(ve3.relativeLuminance2(d, _3, p), ve3.relativeLuminance2(u, h2, c));
    return (d << 24 | _3 << 16 | p << 8 | 255) >>> 0;
  }
  n.increaseLuminance = i8;
  function r5(o2) {
    return [o2 >> 24 & 255, o2 >> 16 & 255, o2 >> 8 & 255, o2 & 255];
  }
  n.toChannels = r5;
})($r2 ||= {});
function vt3(s15) {
  let t = s15.toString(16);
  return t.length < 2 ? "0" + t : t;
}
function Xe3(s15, t) {
  return s15 < t ? (t + 0.05) / (s15 + 0.05) : (s15 + 0.05) / (t + 0.05);
}
var Vr2 = class extends De3 {
  constructor(e, i8, r5) {
    super();
    this.content = 0;
    this.combinedData = "";
    this.fg = e.fg, this.bg = e.bg, this.combinedData = i8, this._width = r5;
  }
  isCombined() {
    return 2097152;
  }
  getWidth() {
    return this._width;
  }
  getChars() {
    return this.combinedData;
  }
  getCode() {
    return 2097151;
  }
  setFromCharData(e) {
    throw new Error("not implemented");
  }
  getAsCharData() {
    return [this.fg, this.getChars(), this.getWidth(), this.getCode()];
  }
};
var ct = class {
  constructor(t) {
    this._bufferService = t;
    this._characterJoiners = [];
    this._nextCharacterJoinerId = 0;
    this._workCell = new q();
  }
  register(t) {
    let e = { id: this._nextCharacterJoinerId++, handler: t };
    return this._characterJoiners.push(e), e.id;
  }
  deregister(t) {
    for (let e = 0; e < this._characterJoiners.length; e++) if (this._characterJoiners[e].id === t) return this._characterJoiners.splice(e, 1), true;
    return false;
  }
  getJoinedCharacters(t) {
    if (this._characterJoiners.length === 0) return [];
    let e = this._bufferService.buffer.lines.get(t);
    if (!e || e.length === 0) return [];
    let i8 = [], r5 = e.translateToString(true), n = 0, o2 = 0, l = 0, a = e.getFg(0), u = e.getBg(0);
    for (let h2 = 0; h2 < e.getTrimmedLength(); h2++) if (e.loadCell(h2, this._workCell), this._workCell.getWidth() !== 0) {
      if (this._workCell.fg !== a || this._workCell.bg !== u) {
        if (h2 - n > 1) {
          let c = this._getJoinedRanges(r5, l, o2, e, n);
          for (let d = 0; d < c.length; d++) i8.push(c[d]);
        }
        n = h2, l = o2, a = this._workCell.fg, u = this._workCell.bg;
      }
      o2 += this._workCell.getChars().length || we2.length;
    }
    if (this._bufferService.cols - n > 1) {
      let h2 = this._getJoinedRanges(r5, l, o2, e, n);
      for (let c = 0; c < h2.length; c++) i8.push(h2[c]);
    }
    return i8;
  }
  _getJoinedRanges(t, e, i8, r5, n) {
    let o2 = t.substring(e, i8), l = [];
    try {
      l = this._characterJoiners[0].handler(o2);
    } catch (a) {
      console.error(a);
    }
    for (let a = 1; a < this._characterJoiners.length; a++) try {
      let u = this._characterJoiners[a].handler(o2);
      for (let h2 = 0; h2 < u.length; h2++) ct._mergeRanges(l, u[h2]);
    } catch (u) {
      console.error(u);
    }
    return this._stringRangesToCellRanges(l, r5, n), l;
  }
  _stringRangesToCellRanges(t, e, i8) {
    let r5 = 0, n = false, o2 = 0, l = t[r5];
    if (l) {
      for (let a = i8; a < this._bufferService.cols; a++) {
        let u = e.getWidth(a), h2 = e.getString(a).length || we2.length;
        if (u !== 0) {
          if (!n && l[0] <= o2 && (l[0] = a, n = true), l[1] <= o2) {
            if (l[1] = a, l = t[++r5], !l) break;
            l[0] <= o2 ? (l[0] = a, n = true) : n = false;
          }
          o2 += h2;
        }
      }
      l && (l[1] = this._bufferService.cols);
    }
  }
  static _mergeRanges(t, e) {
    let i8 = false;
    for (let r5 = 0; r5 < t.length; r5++) {
      let n = t[r5];
      if (i8) {
        if (e[1] <= n[0]) return t[r5 - 1][1] = e[1], t;
        if (e[1] <= n[1]) return t[r5 - 1][1] = Math.max(e[1], n[1]), t.splice(r5, 1), t;
        t.splice(r5, 1), r5--;
      } else {
        if (e[1] <= n[0]) return t.splice(r5, 0, e), t;
        if (e[1] <= n[1]) return n[0] = Math.min(e[0], n[0]), t;
        e[0] < n[1] && (n[0] = Math.min(e[0], n[0]), i8 = true);
        continue;
      }
    }
    return i8 ? t[t.length - 1][1] = e[1] : t.push(e), t;
  }
};
ct = M3([S(0, F2)], ct);
function Oa(s15) {
  return 57508 <= s15 && s15 <= 57558;
}
function Ba(s15) {
  return 9472 <= s15 && s15 <= 9631;
}
function $o(s15) {
  return Oa(s15) || Ba(s15);
}
function Vo() {
  return { css: { canvas: qr2(), cell: qr2() }, device: { canvas: qr2(), cell: qr2(), char: { width: 0, height: 0, left: 0, top: 0 } } };
}
function qr2() {
  return { width: 0, height: 0 };
}
var Vt2 = class {
  constructor(t, e, i8, r5, n, o2, l) {
    this._document = t;
    this._characterJoinerService = e;
    this._optionsService = i8;
    this._coreBrowserService = r5;
    this._coreService = n;
    this._decorationService = o2;
    this._themeService = l;
    this._workCell = new q();
    this._columnSelectMode = false;
    this.defaultSpacing = 0;
  }
  handleSelectionChanged(t, e, i8) {
    this._selectionStart = t, this._selectionEnd = e, this._columnSelectMode = i8;
  }
  createRow(t, e, i8, r5, n, o2, l, a, u, h2, c) {
    let d = [], _3 = this._characterJoinerService.getJoinedCharacters(e), p = this._themeService.colors, m = t.getNoBgTrimmedLength();
    i8 && m < o2 + 1 && (m = o2 + 1);
    let f, A2 = 0, R3 = "", O3 = 0, I2 = 0, k2 = 0, P3 = 0, oe = false, Me3 = 0, Pe3 = false, Ke2 = 0, di = 0, V2 = [], Qe2 = h2 !== -1 && c !== -1;
    for (let y = 0; y < m; y++) {
      t.loadCell(y, this._workCell);
      let T = this._workCell.getWidth();
      if (T === 0) continue;
      let g = false, w2 = y >= di, E = y, x = this._workCell;
      if (_3.length > 0 && y === _3[0][0] && w2) {
        let W = _3.shift(), An2 = this._isCellInSelection(W[0], e);
        for (O3 = W[0] + 1; O3 < W[1]; O3++) w2 &&= An2 === this._isCellInSelection(O3, e);
        w2 &&= !i8 || o2 < W[0] || o2 >= W[1], w2 ? (g = true, x = new Vr2(this._workCell, t.translateToString(true, W[0], W[1]), W[1] - W[0]), E = W[1] - 1, T = x.getWidth()) : di = W[1];
      }
      let N2 = this._isCellInSelection(y, e), Z3 = i8 && y === o2, te4 = Qe2 && y >= h2 && y <= c, Oe3 = false;
      this._decorationService.forEachDecorationAtCell(y, e, void 0, (W) => {
        Oe3 = true;
      });
      let ze2 = x.getChars() || we2;
      if (ze2 === " " && (x.isUnderline() || x.isOverline()) && (ze2 = " "), Ke2 = T * a - u.get(ze2, x.isBold(), x.isItalic()), !f) f = this._document.createElement("span");
      else if (A2 && (N2 && Pe3 || !N2 && !Pe3 && x.bg === I2) && (N2 && Pe3 && p.selectionForeground || x.fg === k2) && x.extended.ext === P3 && te4 === oe && Ke2 === Me3 && !Z3 && !g && !Oe3 && w2) {
        x.isInvisible() ? R3 += we2 : R3 += ze2, A2++;
        continue;
      } else A2 && (f.textContent = R3), f = this._document.createElement("span"), A2 = 0, R3 = "";
      if (I2 = x.bg, k2 = x.fg, P3 = x.extended.ext, oe = te4, Me3 = Ke2, Pe3 = N2, g && o2 >= y && o2 <= E && (o2 = y), !this._coreService.isCursorHidden && Z3 && this._coreService.isCursorInitialized) {
        if (V2.push("xterm-cursor"), this._coreBrowserService.isFocused) l && V2.push("xterm-cursor-blink"), V2.push(r5 === "bar" ? "xterm-cursor-bar" : r5 === "underline" ? "xterm-cursor-underline" : "xterm-cursor-block");
        else if (n) switch (n) {
          case "outline":
            V2.push("xterm-cursor-outline");
            break;
          case "block":
            V2.push("xterm-cursor-block");
            break;
          case "bar":
            V2.push("xterm-cursor-bar");
            break;
          case "underline":
            V2.push("xterm-cursor-underline");
            break;
          default:
            break;
        }
      }
      if (x.isBold() && V2.push("xterm-bold"), x.isItalic() && V2.push("xterm-italic"), x.isDim() && V2.push("xterm-dim"), x.isInvisible() ? R3 = we2 : R3 = x.getChars() || we2, x.isUnderline() && (V2.push(`xterm-underline-${x.extended.underlineStyle}`), R3 === " " && (R3 = " "), !x.isUnderlineColorDefault())) if (x.isUnderlineColorRGB()) f.style.textDecorationColor = `rgb(${De3.toColorRGB(x.getUnderlineColor()).join(",")})`;
      else {
        let W = x.getUnderlineColor();
        this._optionsService.rawOptions.drawBoldTextInBrightColors && x.isBold() && W < 8 && (W += 8), f.style.textDecorationColor = p.ansi[W].css;
      }
      x.isOverline() && (V2.push("xterm-overline"), R3 === " " && (R3 = " ")), x.isStrikethrough() && V2.push("xterm-strikethrough"), te4 && (f.style.textDecoration = "underline");
      let le2 = x.getFgColor(), et3 = x.getFgColorMode(), me3 = x.getBgColor(), ht2 = x.getBgColorMode(), fi2 = !!x.isInverse();
      if (fi2) {
        let W = le2;
        le2 = me3, me3 = W;
        let An2 = et3;
        et3 = ht2, ht2 = An2;
      }
      let tt3, Qi2, pi2 = false;
      this._decorationService.forEachDecorationAtCell(y, e, void 0, (W) => {
        W.options.layer !== "top" && pi2 || (W.backgroundColorRGB && (ht2 = 50331648, me3 = W.backgroundColorRGB.rgba >> 8 & 16777215, tt3 = W.backgroundColorRGB), W.foregroundColorRGB && (et3 = 50331648, le2 = W.foregroundColorRGB.rgba >> 8 & 16777215, Qi2 = W.foregroundColorRGB), pi2 = W.options.layer === "top");
      }), !pi2 && N2 && (tt3 = this._coreBrowserService.isFocused ? p.selectionBackgroundOpaque : p.selectionInactiveBackgroundOpaque, me3 = tt3.rgba >> 8 & 16777215, ht2 = 50331648, pi2 = true, p.selectionForeground && (et3 = 50331648, le2 = p.selectionForeground.rgba >> 8 & 16777215, Qi2 = p.selectionForeground)), pi2 && V2.push("xterm-decoration-top");
      let it3;
      switch (ht2) {
        case 16777216:
        case 33554432:
          it3 = p.ansi[me3], V2.push(`xterm-bg-${me3}`);
          break;
        case 50331648:
          it3 = j2.toColor(me3 >> 16, me3 >> 8 & 255, me3 & 255), this._addStyle(f, `background-color:#${qo((me3 >>> 0).toString(16), "0", 6)}`);
          break;
        case 0:
        default:
          fi2 ? (it3 = p.foreground, V2.push(`xterm-bg-${257}`)) : it3 = p.background;
      }
      switch (tt3 || x.isDim() && (tt3 = U3.multiplyOpacity(it3, 0.5)), et3) {
        case 16777216:
        case 33554432:
          x.isBold() && le2 < 8 && this._optionsService.rawOptions.drawBoldTextInBrightColors && (le2 += 8), this._applyMinimumContrast(f, it3, p.ansi[le2], x, tt3, void 0) || V2.push(`xterm-fg-${le2}`);
          break;
        case 50331648:
          let W = j2.toColor(le2 >> 16 & 255, le2 >> 8 & 255, le2 & 255);
          this._applyMinimumContrast(f, it3, W, x, tt3, Qi2) || this._addStyle(f, `color:#${qo(le2.toString(16), "0", 6)}`);
          break;
        case 0:
        default:
          this._applyMinimumContrast(f, it3, p.foreground, x, tt3, Qi2) || fi2 && V2.push(`xterm-fg-${257}`);
      }
      V2.length && (f.className = V2.join(" "), V2.length = 0), !Z3 && !g && !Oe3 && w2 ? A2++ : f.textContent = R3, Ke2 !== this.defaultSpacing && (f.style.letterSpacing = `${Ke2}px`), d.push(f), y = E;
    }
    return f && A2 && (f.textContent = R3), d;
  }
  _applyMinimumContrast(t, e, i8, r5, n, o2) {
    if (this._optionsService.rawOptions.minimumContrastRatio === 1 || $o(r5.getCode())) return false;
    let l = this._getContrastCache(r5), a;
    if (!n && !o2 && (a = l.getColor(e.rgba, i8.rgba)), a === void 0) {
      let u = this._optionsService.rawOptions.minimumContrastRatio / (r5.isDim() ? 2 : 1);
      a = U3.ensureContrastRatio(n || e, o2 || i8, u), l.setColor((n || e).rgba, (o2 || i8).rgba, a ?? null);
    }
    return a ? (this._addStyle(t, `color:${a.css}`), true) : false;
  }
  _getContrastCache(t) {
    return t.isDim() ? this._themeService.colors.halfContrastCache : this._themeService.colors.contrastCache;
  }
  _addStyle(t, e) {
    t.setAttribute("style", `${t.getAttribute("style") || ""}${e};`);
  }
  _isCellInSelection(t, e) {
    let i8 = this._selectionStart, r5 = this._selectionEnd;
    return !i8 || !r5 ? false : this._columnSelectMode ? i8[0] <= r5[0] ? t >= i8[0] && e >= i8[1] && t < r5[0] && e <= r5[1] : t < i8[0] && e >= i8[1] && t >= r5[0] && e <= r5[1] : e > i8[1] && e < r5[1] || i8[1] === r5[1] && e === i8[1] && t >= i8[0] && t < r5[0] || i8[1] < r5[1] && e === r5[1] && t < r5[0] || i8[1] < r5[1] && e === i8[1] && t >= i8[0];
  }
};
Vt2 = M3([S(1, or2), S(2, H3), S(3, ae2), S(4, ge3), S(5, Be4), S(6, Re2)], Vt2);
function qo(s15, t, e) {
  for (; s15.length < e; ) s15 = t + s15;
  return s15;
}
var Yr2 = class {
  constructor(t, e) {
    this._flat = new Float32Array(256);
    this._font = "";
    this._fontSize = 0;
    this._weight = "normal";
    this._weightBold = "bold";
    this._measureElements = [];
    this._container = t.createElement("div"), this._container.classList.add("xterm-width-cache-measure-container"), this._container.setAttribute("aria-hidden", "true"), this._container.style.whiteSpace = "pre", this._container.style.fontKerning = "none";
    let i8 = t.createElement("span");
    i8.classList.add("xterm-char-measure-element");
    let r5 = t.createElement("span");
    r5.classList.add("xterm-char-measure-element"), r5.style.fontWeight = "bold";
    let n = t.createElement("span");
    n.classList.add("xterm-char-measure-element"), n.style.fontStyle = "italic";
    let o2 = t.createElement("span");
    o2.classList.add("xterm-char-measure-element"), o2.style.fontWeight = "bold", o2.style.fontStyle = "italic", this._measureElements = [i8, r5, n, o2], this._container.appendChild(i8), this._container.appendChild(r5), this._container.appendChild(n), this._container.appendChild(o2), e.appendChild(this._container), this.clear();
  }
  dispose() {
    this._container.remove(), this._measureElements.length = 0, this._holey = void 0;
  }
  clear() {
    this._flat.fill(-9999), this._holey = /* @__PURE__ */ new Map();
  }
  setFont(t, e, i8, r5) {
    t === this._font && e === this._fontSize && i8 === this._weight && r5 === this._weightBold || (this._font = t, this._fontSize = e, this._weight = i8, this._weightBold = r5, this._container.style.fontFamily = this._font, this._container.style.fontSize = `${this._fontSize}px`, this._measureElements[0].style.fontWeight = `${i8}`, this._measureElements[1].style.fontWeight = `${r5}`, this._measureElements[2].style.fontWeight = `${i8}`, this._measureElements[3].style.fontWeight = `${r5}`, this.clear());
  }
  get(t, e, i8) {
    let r5 = 0;
    if (!e && !i8 && t.length === 1 && (r5 = t.charCodeAt(0)) < 256) {
      if (this._flat[r5] !== -9999) return this._flat[r5];
      let l = this._measure(t, 0);
      return l > 0 && (this._flat[r5] = l), l;
    }
    let n = t;
    e && (n += "B"), i8 && (n += "I");
    let o2 = this._holey.get(n);
    if (o2 === void 0) {
      let l = 0;
      e && (l |= 1), i8 && (l |= 2), o2 = this._measure(t, l), o2 > 0 && this._holey.set(n, o2);
    }
    return o2;
  }
  _measure(t, e) {
    let i8 = this._measureElements[e];
    return i8.textContent = t.repeat(32), i8.offsetWidth / 32;
  }
};
var ms2 = class {
  constructor() {
    this.clear();
  }
  clear() {
    this.hasSelection = false, this.columnSelectMode = false, this.viewportStartRow = 0, this.viewportEndRow = 0, this.viewportCappedStartRow = 0, this.viewportCappedEndRow = 0, this.startCol = 0, this.endCol = 0, this.selectionStart = void 0, this.selectionEnd = void 0;
  }
  update(t, e, i8, r5 = false) {
    if (this.selectionStart = e, this.selectionEnd = i8, !e || !i8 || e[0] === i8[0] && e[1] === i8[1]) {
      this.clear();
      return;
    }
    let n = t.buffers.active.ydisp, o2 = e[1] - n, l = i8[1] - n, a = Math.max(o2, 0), u = Math.min(l, t.rows - 1);
    if (a >= t.rows || u < 0) {
      this.clear();
      return;
    }
    this.hasSelection = true, this.columnSelectMode = r5, this.viewportStartRow = o2, this.viewportEndRow = l, this.viewportCappedStartRow = a, this.viewportCappedEndRow = u, this.startCol = e[0], this.endCol = i8[0];
  }
  isCellSelected(t, e, i8) {
    return this.hasSelection ? (i8 -= t.buffer.active.viewportY, this.columnSelectMode ? this.startCol <= this.endCol ? e >= this.startCol && i8 >= this.viewportCappedStartRow && e < this.endCol && i8 <= this.viewportCappedEndRow : e < this.startCol && i8 >= this.viewportCappedStartRow && e >= this.endCol && i8 <= this.viewportCappedEndRow : i8 > this.viewportStartRow && i8 < this.viewportEndRow || this.viewportStartRow === this.viewportEndRow && i8 === this.viewportStartRow && e >= this.startCol && e < this.endCol || this.viewportStartRow < this.viewportEndRow && i8 === this.viewportEndRow && e < this.endCol || this.viewportStartRow < this.viewportEndRow && i8 === this.viewportStartRow && e >= this.startCol) : false;
  }
};
function Yo() {
  return new ms2();
}
var _s2 = "xterm-dom-renderer-owner-";
var Le3 = "xterm-rows";
var jr2 = "xterm-fg-";
var jo = "xterm-bg-";
var ki2 = "xterm-focus";
var Xr2 = "xterm-selection";
var Na = 1;
var Yt2 = class extends D2 {
  constructor(e, i8, r5, n, o2, l, a, u, h2, c, d, _3, p, m) {
    super();
    this._terminal = e;
    this._document = i8;
    this._element = r5;
    this._screenElement = n;
    this._viewportElement = o2;
    this._helperContainer = l;
    this._linkifier2 = a;
    this._charSizeService = h2;
    this._optionsService = c;
    this._bufferService = d;
    this._coreService = _3;
    this._coreBrowserService = p;
    this._themeService = m;
    this._terminalClass = Na++;
    this._rowElements = [];
    this._selectionRenderModel = Yo();
    this.onRequestRedraw = this._register(new v()).event;
    this._rowContainer = this._document.createElement("div"), this._rowContainer.classList.add(Le3), this._rowContainer.style.lineHeight = "normal", this._rowContainer.setAttribute("aria-hidden", "true"), this._refreshRowElements(this._bufferService.cols, this._bufferService.rows), this._selectionContainer = this._document.createElement("div"), this._selectionContainer.classList.add(Xr2), this._selectionContainer.setAttribute("aria-hidden", "true"), this.dimensions = Vo(), this._updateDimensions(), this._register(this._optionsService.onOptionChange(() => this._handleOptionsChanged())), this._register(this._themeService.onChangeColors((f) => this._injectCss(f))), this._injectCss(this._themeService.colors), this._rowFactory = u.createInstance(Vt2, document), this._element.classList.add(_s2 + this._terminalClass), this._screenElement.appendChild(this._rowContainer), this._screenElement.appendChild(this._selectionContainer), this._register(this._linkifier2.onShowLinkUnderline((f) => this._handleLinkHover(f))), this._register(this._linkifier2.onHideLinkUnderline((f) => this._handleLinkLeave(f))), this._register(C3(() => {
      this._element.classList.remove(_s2 + this._terminalClass), this._rowContainer.remove(), this._selectionContainer.remove(), this._widthCache.dispose(), this._themeStyleElement.remove(), this._dimensionsStyleElement.remove();
    })), this._widthCache = new Yr2(this._document, this._helperContainer), this._widthCache.setFont(this._optionsService.rawOptions.fontFamily, this._optionsService.rawOptions.fontSize, this._optionsService.rawOptions.fontWeight, this._optionsService.rawOptions.fontWeightBold), this._setDefaultSpacing();
  }
  _updateDimensions() {
    let e = this._coreBrowserService.dpr;
    this.dimensions.device.char.width = this._charSizeService.width * e, this.dimensions.device.char.height = Math.ceil(this._charSizeService.height * e), this.dimensions.device.cell.width = this.dimensions.device.char.width + Math.round(this._optionsService.rawOptions.letterSpacing), this.dimensions.device.cell.height = Math.floor(this.dimensions.device.char.height * this._optionsService.rawOptions.lineHeight), this.dimensions.device.char.left = 0, this.dimensions.device.char.top = 0, this.dimensions.device.canvas.width = this.dimensions.device.cell.width * this._bufferService.cols, this.dimensions.device.canvas.height = this.dimensions.device.cell.height * this._bufferService.rows, this.dimensions.css.canvas.width = Math.round(this.dimensions.device.canvas.width / e), this.dimensions.css.canvas.height = Math.round(this.dimensions.device.canvas.height / e), this.dimensions.css.cell.width = this.dimensions.css.canvas.width / this._bufferService.cols, this.dimensions.css.cell.height = this.dimensions.css.canvas.height / this._bufferService.rows;
    for (let r5 of this._rowElements) r5.style.width = `${this.dimensions.css.canvas.width}px`, r5.style.height = `${this.dimensions.css.cell.height}px`, r5.style.lineHeight = `${this.dimensions.css.cell.height}px`, r5.style.overflow = "hidden";
    this._dimensionsStyleElement || (this._dimensionsStyleElement = this._document.createElement("style"), this._screenElement.appendChild(this._dimensionsStyleElement));
    let i8 = `${this._terminalSelector} .${Le3} span { display: inline-block; height: 100%; vertical-align: top;}`;
    this._dimensionsStyleElement.textContent = i8, this._selectionContainer.style.height = this._viewportElement.style.height, this._screenElement.style.width = `${this.dimensions.css.canvas.width}px`, this._screenElement.style.height = `${this.dimensions.css.canvas.height}px`;
  }
  _injectCss(e) {
    this._themeStyleElement || (this._themeStyleElement = this._document.createElement("style"), this._screenElement.appendChild(this._themeStyleElement));
    let i8 = `${this._terminalSelector} .${Le3} { pointer-events: none; color: ${e.foreground.css}; font-family: ${this._optionsService.rawOptions.fontFamily}; font-size: ${this._optionsService.rawOptions.fontSize}px; font-kerning: none; white-space: pre}`;
    i8 += `${this._terminalSelector} .${Le3} .xterm-dim { color: ${U3.multiplyOpacity(e.foreground, 0.5).css};}`, i8 += `${this._terminalSelector} span:not(.xterm-bold) { font-weight: ${this._optionsService.rawOptions.fontWeight};}${this._terminalSelector} span.xterm-bold { font-weight: ${this._optionsService.rawOptions.fontWeightBold};}${this._terminalSelector} span.xterm-italic { font-style: italic;}`;
    let r5 = `blink_underline_${this._terminalClass}`, n = `blink_bar_${this._terminalClass}`, o2 = `blink_block_${this._terminalClass}`;
    i8 += `@keyframes ${r5} { 50% {  border-bottom-style: hidden; }}`, i8 += `@keyframes ${n} { 50% {  box-shadow: none; }}`, i8 += `@keyframes ${o2} { 0% {  background-color: ${e.cursor.css};  color: ${e.cursorAccent.css}; } 50% {  background-color: inherit;  color: ${e.cursor.css}; }}`, i8 += `${this._terminalSelector} .${Le3}.${ki2} .xterm-cursor.xterm-cursor-blink.xterm-cursor-underline { animation: ${r5} 1s step-end infinite;}${this._terminalSelector} .${Le3}.${ki2} .xterm-cursor.xterm-cursor-blink.xterm-cursor-bar { animation: ${n} 1s step-end infinite;}${this._terminalSelector} .${Le3}.${ki2} .xterm-cursor.xterm-cursor-blink.xterm-cursor-block { animation: ${o2} 1s step-end infinite;}${this._terminalSelector} .${Le3} .xterm-cursor.xterm-cursor-block { background-color: ${e.cursor.css}; color: ${e.cursorAccent.css};}${this._terminalSelector} .${Le3} .xterm-cursor.xterm-cursor-block:not(.xterm-cursor-blink) { background-color: ${e.cursor.css} !important; color: ${e.cursorAccent.css} !important;}${this._terminalSelector} .${Le3} .xterm-cursor.xterm-cursor-outline { outline: 1px solid ${e.cursor.css}; outline-offset: -1px;}${this._terminalSelector} .${Le3} .xterm-cursor.xterm-cursor-bar { box-shadow: ${this._optionsService.rawOptions.cursorWidth}px 0 0 ${e.cursor.css} inset;}${this._terminalSelector} .${Le3} .xterm-cursor.xterm-cursor-underline { border-bottom: 1px ${e.cursor.css}; border-bottom-style: solid; height: calc(100% - 1px);}`, i8 += `${this._terminalSelector} .${Xr2} { position: absolute; top: 0; left: 0; z-index: 1; pointer-events: none;}${this._terminalSelector}.focus .${Xr2} div { position: absolute; background-color: ${e.selectionBackgroundOpaque.css};}${this._terminalSelector} .${Xr2} div { position: absolute; background-color: ${e.selectionInactiveBackgroundOpaque.css};}`;
    for (let [l, a] of e.ansi.entries()) i8 += `${this._terminalSelector} .${jr2}${l} { color: ${a.css}; }${this._terminalSelector} .${jr2}${l}.xterm-dim { color: ${U3.multiplyOpacity(a, 0.5).css}; }${this._terminalSelector} .${jo}${l} { background-color: ${a.css}; }`;
    i8 += `${this._terminalSelector} .${jr2}${257} { color: ${U3.opaque(e.background).css}; }${this._terminalSelector} .${jr2}${257}.xterm-dim { color: ${U3.multiplyOpacity(U3.opaque(e.background), 0.5).css}; }${this._terminalSelector} .${jo}${257} { background-color: ${e.foreground.css}; }`, this._themeStyleElement.textContent = i8;
  }
  _setDefaultSpacing() {
    let e = this.dimensions.css.cell.width - this._widthCache.get("W", false, false);
    this._rowContainer.style.letterSpacing = `${e}px`, this._rowFactory.defaultSpacing = e;
  }
  handleDevicePixelRatioChange() {
    this._updateDimensions(), this._widthCache.clear(), this._setDefaultSpacing();
  }
  _refreshRowElements(e, i8) {
    for (let r5 = this._rowElements.length; r5 <= i8; r5++) {
      let n = this._document.createElement("div");
      this._rowContainer.appendChild(n), this._rowElements.push(n);
    }
    for (; this._rowElements.length > i8; ) this._rowContainer.removeChild(this._rowElements.pop());
  }
  handleResize(e, i8) {
    this._refreshRowElements(e, i8), this._updateDimensions(), this.handleSelectionChanged(this._selectionRenderModel.selectionStart, this._selectionRenderModel.selectionEnd, this._selectionRenderModel.columnSelectMode);
  }
  handleCharSizeChanged() {
    this._updateDimensions(), this._widthCache.clear(), this._setDefaultSpacing();
  }
  handleBlur() {
    this._rowContainer.classList.remove(ki2), this.renderRows(0, this._bufferService.rows - 1);
  }
  handleFocus() {
    this._rowContainer.classList.add(ki2), this.renderRows(this._bufferService.buffer.y, this._bufferService.buffer.y);
  }
  handleSelectionChanged(e, i8, r5) {
    if (this._selectionContainer.replaceChildren(), this._rowFactory.handleSelectionChanged(e, i8, r5), this.renderRows(0, this._bufferService.rows - 1), !e || !i8 || (this._selectionRenderModel.update(this._terminal, e, i8, r5), !this._selectionRenderModel.hasSelection)) return;
    let n = this._selectionRenderModel.viewportStartRow, o2 = this._selectionRenderModel.viewportEndRow, l = this._selectionRenderModel.viewportCappedStartRow, a = this._selectionRenderModel.viewportCappedEndRow, u = this._document.createDocumentFragment();
    if (r5) {
      let h2 = e[0] > i8[0];
      u.appendChild(this._createSelectionElement(l, h2 ? i8[0] : e[0], h2 ? e[0] : i8[0], a - l + 1));
    } else {
      let h2 = n === l ? e[0] : 0, c = l === o2 ? i8[0] : this._bufferService.cols;
      u.appendChild(this._createSelectionElement(l, h2, c));
      let d = a - l - 1;
      if (u.appendChild(this._createSelectionElement(l + 1, 0, this._bufferService.cols, d)), l !== a) {
        let _3 = o2 === a ? i8[0] : this._bufferService.cols;
        u.appendChild(this._createSelectionElement(a, 0, _3));
      }
    }
    this._selectionContainer.appendChild(u);
  }
  _createSelectionElement(e, i8, r5, n = 1) {
    let o2 = this._document.createElement("div"), l = i8 * this.dimensions.css.cell.width, a = this.dimensions.css.cell.width * (r5 - i8);
    return l + a > this.dimensions.css.canvas.width && (a = this.dimensions.css.canvas.width - l), o2.style.height = `${n * this.dimensions.css.cell.height}px`, o2.style.top = `${e * this.dimensions.css.cell.height}px`, o2.style.left = `${l}px`, o2.style.width = `${a}px`, o2;
  }
  handleCursorMove() {
  }
  _handleOptionsChanged() {
    this._updateDimensions(), this._injectCss(this._themeService.colors), this._widthCache.setFont(this._optionsService.rawOptions.fontFamily, this._optionsService.rawOptions.fontSize, this._optionsService.rawOptions.fontWeight, this._optionsService.rawOptions.fontWeightBold), this._setDefaultSpacing();
  }
  clear() {
    for (let e of this._rowElements) e.replaceChildren();
  }
  renderRows(e, i8) {
    let r5 = this._bufferService.buffer, n = r5.ybase + r5.y, o2 = Math.min(r5.x, this._bufferService.cols - 1), l = this._coreService.decPrivateModes.cursorBlink ?? this._optionsService.rawOptions.cursorBlink, a = this._coreService.decPrivateModes.cursorStyle ?? this._optionsService.rawOptions.cursorStyle, u = this._optionsService.rawOptions.cursorInactiveStyle;
    for (let h2 = e; h2 <= i8; h2++) {
      let c = h2 + r5.ydisp, d = this._rowElements[h2], _3 = r5.lines.get(c);
      if (!d || !_3) break;
      d.replaceChildren(...this._rowFactory.createRow(_3, c, c === n, a, u, o2, l, this.dimensions.css.cell.width, this._widthCache, -1, -1));
    }
  }
  get _terminalSelector() {
    return `.${_s2}${this._terminalClass}`;
  }
  _handleLinkHover(e) {
    this._setCellUnderline(e.x1, e.x2, e.y1, e.y2, e.cols, true);
  }
  _handleLinkLeave(e) {
    this._setCellUnderline(e.x1, e.x2, e.y1, e.y2, e.cols, false);
  }
  _setCellUnderline(e, i8, r5, n, o2, l) {
    r5 < 0 && (e = 0), n < 0 && (i8 = 0);
    let a = this._bufferService.rows - 1;
    r5 = Math.max(Math.min(r5, a), 0), n = Math.max(Math.min(n, a), 0), o2 = Math.min(o2, this._bufferService.cols);
    let u = this._bufferService.buffer, h2 = u.ybase + u.y, c = Math.min(u.x, o2 - 1), d = this._optionsService.rawOptions.cursorBlink, _3 = this._optionsService.rawOptions.cursorStyle, p = this._optionsService.rawOptions.cursorInactiveStyle;
    for (let m = r5; m <= n; ++m) {
      let f = m + u.ydisp, A2 = this._rowElements[m], R3 = u.lines.get(f);
      if (!A2 || !R3) break;
      A2.replaceChildren(...this._rowFactory.createRow(R3, f, f === h2, _3, p, c, d, this.dimensions.css.cell.width, this._widthCache, l ? m === r5 ? e : 0 : -1, l ? (m === n ? i8 : o2) - 1 : -1));
    }
  }
};
Yt2 = M3([S(7, xt4), S(8, nt3), S(9, H3), S(10, F2), S(11, ge3), S(12, ae2), S(13, Re2)], Yt2);
var jt2 = class extends D2 {
  constructor(e, i8, r5) {
    super();
    this._optionsService = r5;
    this.width = 0;
    this.height = 0;
    this._onCharSizeChange = this._register(new v());
    this.onCharSizeChange = this._onCharSizeChange.event;
    try {
      this._measureStrategy = this._register(new vs2(this._optionsService));
    } catch {
      this._measureStrategy = this._register(new bs2(e, i8, this._optionsService));
    }
    this._register(this._optionsService.onMultipleOptionChange(["fontFamily", "fontSize"], () => this.measure()));
  }
  get hasValidSize() {
    return this.width > 0 && this.height > 0;
  }
  measure() {
    let e = this._measureStrategy.measure();
    (e.width !== this.width || e.height !== this.height) && (this.width = e.width, this.height = e.height, this._onCharSizeChange.fire());
  }
};
jt2 = M3([S(2, H3)], jt2);
var Zr = class extends D2 {
  constructor() {
    super(...arguments);
    this._result = { width: 0, height: 0 };
  }
  _validateAndSet(e, i8) {
    e !== void 0 && e > 0 && i8 !== void 0 && i8 > 0 && (this._result.width = e, this._result.height = i8);
  }
};
var bs2 = class extends Zr {
  constructor(e, i8, r5) {
    super();
    this._document = e;
    this._parentElement = i8;
    this._optionsService = r5;
    this._measureElement = this._document.createElement("span"), this._measureElement.classList.add("xterm-char-measure-element"), this._measureElement.textContent = "W".repeat(32), this._measureElement.setAttribute("aria-hidden", "true"), this._measureElement.style.whiteSpace = "pre", this._measureElement.style.fontKerning = "none", this._parentElement.appendChild(this._measureElement);
  }
  measure() {
    return this._measureElement.style.fontFamily = this._optionsService.rawOptions.fontFamily, this._measureElement.style.fontSize = `${this._optionsService.rawOptions.fontSize}px`, this._validateAndSet(Number(this._measureElement.offsetWidth) / 32, Number(this._measureElement.offsetHeight)), this._result;
  }
};
var vs2 = class extends Zr {
  constructor(e) {
    super();
    this._optionsService = e;
    this._canvas = new OffscreenCanvas(100, 100), this._ctx = this._canvas.getContext("2d");
    let i8 = this._ctx.measureText("W");
    if (!("width" in i8 && "fontBoundingBoxAscent" in i8 && "fontBoundingBoxDescent" in i8)) throw new Error("Required font metrics not supported");
  }
  measure() {
    this._ctx.font = `${this._optionsService.rawOptions.fontSize}px ${this._optionsService.rawOptions.fontFamily}`;
    let e = this._ctx.measureText("W");
    return this._validateAndSet(e.width, e.fontBoundingBoxAscent + e.fontBoundingBoxDescent), this._result;
  }
};
var Jr2 = class extends D2 {
  constructor(e, i8, r5) {
    super();
    this._textarea = e;
    this._window = i8;
    this.mainDocument = r5;
    this._isFocused = false;
    this._cachedIsFocused = void 0;
    this._screenDprMonitor = this._register(new gs2(this._window));
    this._onDprChange = this._register(new v());
    this.onDprChange = this._onDprChange.event;
    this._onWindowChange = this._register(new v());
    this.onWindowChange = this._onWindowChange.event;
    this._register(this.onWindowChange((n) => this._screenDprMonitor.setWindow(n))), this._register($.forward(this._screenDprMonitor.onDprChange, this._onDprChange)), this._register(L(this._textarea, "focus", () => this._isFocused = true)), this._register(L(this._textarea, "blur", () => this._isFocused = false));
  }
  get window() {
    return this._window;
  }
  set window(e) {
    this._window !== e && (this._window = e, this._onWindowChange.fire(this._window));
  }
  get dpr() {
    return this.window.devicePixelRatio;
  }
  get isFocused() {
    return this._cachedIsFocused === void 0 && (this._cachedIsFocused = this._isFocused && this._textarea.ownerDocument.hasFocus(), queueMicrotask(() => this._cachedIsFocused = void 0)), this._cachedIsFocused;
  }
};
var gs2 = class extends D2 {
  constructor(e) {
    super();
    this._parentWindow = e;
    this._windowResizeListener = this._register(new ye3());
    this._onDprChange = this._register(new v());
    this.onDprChange = this._onDprChange.event;
    this._outerListener = () => this._setDprAndFireIfDiffers(), this._currentDevicePixelRatio = this._parentWindow.devicePixelRatio, this._updateDpr(), this._setWindowResizeListener(), this._register(C3(() => this.clearListener()));
  }
  setWindow(e) {
    this._parentWindow = e, this._setWindowResizeListener(), this._setDprAndFireIfDiffers();
  }
  _setWindowResizeListener() {
    this._windowResizeListener.value = L(this._parentWindow, "resize", () => this._setDprAndFireIfDiffers());
  }
  _setDprAndFireIfDiffers() {
    this._parentWindow.devicePixelRatio !== this._currentDevicePixelRatio && this._onDprChange.fire(this._parentWindow.devicePixelRatio), this._updateDpr();
  }
  _updateDpr() {
    this._outerListener && (this._resolutionMediaMatchList?.removeListener(this._outerListener), this._currentDevicePixelRatio = this._parentWindow.devicePixelRatio, this._resolutionMediaMatchList = this._parentWindow.matchMedia(`screen and (resolution: ${this._parentWindow.devicePixelRatio}dppx)`), this._resolutionMediaMatchList.addListener(this._outerListener));
  }
  clearListener() {
    !this._resolutionMediaMatchList || !this._outerListener || (this._resolutionMediaMatchList.removeListener(this._outerListener), this._resolutionMediaMatchList = void 0, this._outerListener = void 0);
  }
};
var Qr2 = class extends D2 {
  constructor() {
    super();
    this.linkProviders = [];
    this._register(C3(() => this.linkProviders.length = 0));
  }
  registerLinkProvider(e) {
    return this.linkProviders.push(e), { dispose: () => {
      let i8 = this.linkProviders.indexOf(e);
      i8 !== -1 && this.linkProviders.splice(i8, 1);
    } };
  }
};
function Ci2(s15, t, e) {
  let i8 = e.getBoundingClientRect(), r5 = s15.getComputedStyle(e), n = parseInt(r5.getPropertyValue("padding-left")), o2 = parseInt(r5.getPropertyValue("padding-top"));
  return [t.clientX - i8.left - n, t.clientY - i8.top - o2];
}
function Xo(s15, t, e, i8, r5, n, o2, l, a) {
  if (!n) return;
  let u = Ci2(s15, t, e);
  if (u) return u[0] = Math.ceil((u[0] + (a ? o2 / 2 : 0)) / o2), u[1] = Math.ceil(u[1] / l), u[0] = Math.min(Math.max(u[0], 1), i8 + (a ? 1 : 0)), u[1] = Math.min(Math.max(u[1], 1), r5), u;
}
var Xt2 = class {
  constructor(t, e) {
    this._renderService = t;
    this._charSizeService = e;
  }
  getCoords(t, e, i8, r5, n) {
    return Xo(window, t, e, i8, r5, this._charSizeService.hasValidSize, this._renderService.dimensions.css.cell.width, this._renderService.dimensions.css.cell.height, n);
  }
  getMouseReportCoords(t, e) {
    let i8 = Ci2(window, t, e);
    if (this._charSizeService.hasValidSize) return i8[0] = Math.min(Math.max(i8[0], 0), this._renderService.dimensions.css.canvas.width - 1), i8[1] = Math.min(Math.max(i8[1], 0), this._renderService.dimensions.css.canvas.height - 1), { col: Math.floor(i8[0] / this._renderService.dimensions.css.cell.width), row: Math.floor(i8[1] / this._renderService.dimensions.css.cell.height), x: Math.floor(i8[0]), y: Math.floor(i8[1]) };
  }
};
Xt2 = M3([S(0, ce2), S(1, nt3)], Xt2);
var en2 = class {
  constructor(t, e) {
    this._renderCallback = t;
    this._coreBrowserService = e;
    this._refreshCallbacks = [];
  }
  dispose() {
    this._animationFrame && (this._coreBrowserService.window.cancelAnimationFrame(this._animationFrame), this._animationFrame = void 0);
  }
  addRefreshCallback(t) {
    return this._refreshCallbacks.push(t), this._animationFrame || (this._animationFrame = this._coreBrowserService.window.requestAnimationFrame(() => this._innerRefresh())), this._animationFrame;
  }
  refresh(t, e, i8) {
    this._rowCount = i8, t = t !== void 0 ? t : 0, e = e !== void 0 ? e : this._rowCount - 1, this._rowStart = this._rowStart !== void 0 ? Math.min(this._rowStart, t) : t, this._rowEnd = this._rowEnd !== void 0 ? Math.max(this._rowEnd, e) : e, !this._animationFrame && (this._animationFrame = this._coreBrowserService.window.requestAnimationFrame(() => this._innerRefresh()));
  }
  _innerRefresh() {
    if (this._animationFrame = void 0, this._rowStart === void 0 || this._rowEnd === void 0 || this._rowCount === void 0) {
      this._runRefreshCallbacks();
      return;
    }
    let t = Math.max(this._rowStart, 0), e = Math.min(this._rowEnd, this._rowCount - 1);
    this._rowStart = void 0, this._rowEnd = void 0, this._renderCallback(t, e), this._runRefreshCallbacks();
  }
  _runRefreshCallbacks() {
    for (let t of this._refreshCallbacks) t(0);
    this._refreshCallbacks = [];
  }
};
var tn2 = {};
Ll(tn2, { getSafariVersion: () => Ha, isChromeOS: () => Ts2, isFirefox: () => Ss2, isIpad: () => Wa2, isIphone: () => Ua, isLegacyEdge: () => Fa, isLinux: () => Bi2, isMac: () => Zt3, isNode: () => Mi2, isSafari: () => Zo, isWindows: () => Es2 });
var Mi2 = typeof process < "u" && "title" in process;
var Pi2 = Mi2 ? "node" : navigator.userAgent;
var Oi = Mi2 ? "node" : navigator.platform;
var Ss2 = Pi2.includes("Firefox");
var Fa = Pi2.includes("Edge");
var Zo = /^((?!chrome|android).)*safari/i.test(Pi2);
function Ha() {
  if (!Zo) return 0;
  let s15 = Pi2.match(/Version\/(\d+)/);
  return s15 === null || s15.length < 2 ? 0 : parseInt(s15[1]);
}
var Zt3 = ["Macintosh", "MacIntel", "MacPPC", "Mac68K"].includes(Oi);
var Wa2 = Oi === "iPad";
var Ua = Oi === "iPhone";
var Es2 = ["Windows", "Win16", "Win32", "WinCE"].includes(Oi);
var Bi2 = Oi.indexOf("Linux") >= 0;
var Ts2 = /\bCrOS\b/.test(Pi2);
var rn2 = class {
  constructor() {
    this._tasks = [];
    this._i = 0;
  }
  enqueue(t) {
    this._tasks.push(t), this._start();
  }
  flush() {
    for (; this._i < this._tasks.length; ) this._tasks[this._i]() || this._i++;
    this.clear();
  }
  clear() {
    this._idleCallback && (this._cancelCallback(this._idleCallback), this._idleCallback = void 0), this._i = 0, this._tasks.length = 0;
  }
  _start() {
    this._idleCallback || (this._idleCallback = this._requestCallback(this._process.bind(this)));
  }
  _process(t) {
    this._idleCallback = void 0;
    let e = 0, i8 = 0, r5 = t.timeRemaining(), n = 0;
    for (; this._i < this._tasks.length; ) {
      if (e = performance.now(), this._tasks[this._i]() || this._i++, e = Math.max(1, performance.now() - e), i8 = Math.max(e, i8), n = t.timeRemaining(), i8 * 1.5 > n) {
        r5 - e < -20 && console.warn(`task queue exceeded allotted deadline by ${Math.abs(Math.round(r5 - e))}ms`), this._start();
        return;
      }
      r5 = n;
    }
    this.clear();
  }
};
var Is2 = class extends rn2 {
  _requestCallback(t) {
    return setTimeout(() => t(this._createDeadline(16)));
  }
  _cancelCallback(t) {
    clearTimeout(t);
  }
  _createDeadline(t) {
    let e = performance.now() + t;
    return { timeRemaining: () => Math.max(0, e - performance.now()) };
  }
};
var ys = class extends rn2 {
  _requestCallback(t) {
    return requestIdleCallback(t);
  }
  _cancelCallback(t) {
    cancelIdleCallback(t);
  }
};
var Jt = !Mi2 && "requestIdleCallback" in window ? ys : Is2;
var nn2 = class {
  constructor() {
    this._queue = new Jt();
  }
  set(t) {
    this._queue.clear(), this._queue.enqueue(t);
  }
  flush() {
    this._queue.flush();
  }
};
var Qt2 = class extends D2 {
  constructor(e, i8, r5, n, o2, l, a, u, h2) {
    super();
    this._rowCount = e;
    this._optionsService = r5;
    this._charSizeService = n;
    this._coreService = o2;
    this._coreBrowserService = u;
    this._renderer = this._register(new ye3());
    this._pausedResizeTask = new nn2();
    this._observerDisposable = this._register(new ye3());
    this._isPaused = false;
    this._needsFullRefresh = false;
    this._isNextRenderRedrawOnly = true;
    this._needsSelectionRefresh = false;
    this._canvasWidth = 0;
    this._canvasHeight = 0;
    this._selectionState = { start: void 0, end: void 0, columnSelectMode: false };
    this._onDimensionsChange = this._register(new v());
    this.onDimensionsChange = this._onDimensionsChange.event;
    this._onRenderedViewportChange = this._register(new v());
    this.onRenderedViewportChange = this._onRenderedViewportChange.event;
    this._onRender = this._register(new v());
    this.onRender = this._onRender.event;
    this._onRefreshRequest = this._register(new v());
    this.onRefreshRequest = this._onRefreshRequest.event;
    this._renderDebouncer = new en2((c, d) => this._renderRows(c, d), this._coreBrowserService), this._register(this._renderDebouncer), this._syncOutputHandler = new xs2(this._coreBrowserService, this._coreService, () => this._fullRefresh()), this._register(C3(() => this._syncOutputHandler.dispose())), this._register(this._coreBrowserService.onDprChange(() => this.handleDevicePixelRatioChange())), this._register(a.onResize(() => this._fullRefresh())), this._register(a.buffers.onBufferActivate(() => this._renderer.value?.clear())), this._register(this._optionsService.onOptionChange(() => this._handleOptionsChanged())), this._register(this._charSizeService.onCharSizeChange(() => this.handleCharSizeChanged())), this._register(l.onDecorationRegistered(() => this._fullRefresh())), this._register(l.onDecorationRemoved(() => this._fullRefresh())), this._register(this._optionsService.onMultipleOptionChange(["customGlyphs", "drawBoldTextInBrightColors", "letterSpacing", "lineHeight", "fontFamily", "fontSize", "fontWeight", "fontWeightBold", "minimumContrastRatio", "rescaleOverlappingGlyphs"], () => {
      this.clear(), this.handleResize(a.cols, a.rows), this._fullRefresh();
    })), this._register(this._optionsService.onMultipleOptionChange(["cursorBlink", "cursorStyle"], () => this.refreshRows(a.buffer.y, a.buffer.y, true))), this._register(h2.onChangeColors(() => this._fullRefresh())), this._registerIntersectionObserver(this._coreBrowserService.window, i8), this._register(this._coreBrowserService.onWindowChange((c) => this._registerIntersectionObserver(c, i8)));
  }
  get dimensions() {
    return this._renderer.value.dimensions;
  }
  _registerIntersectionObserver(e, i8) {
    if ("IntersectionObserver" in e) {
      let r5 = new e.IntersectionObserver((n) => this._handleIntersectionChange(n[n.length - 1]), { threshold: 0 });
      r5.observe(i8), this._observerDisposable.value = C3(() => r5.disconnect());
    }
  }
  _handleIntersectionChange(e) {
    this._isPaused = e.isIntersecting === void 0 ? e.intersectionRatio === 0 : !e.isIntersecting, !this._isPaused && !this._charSizeService.hasValidSize && this._charSizeService.measure(), !this._isPaused && this._needsFullRefresh && (this._pausedResizeTask.flush(), this.refreshRows(0, this._rowCount - 1), this._needsFullRefresh = false);
  }
  refreshRows(e, i8, r5 = false) {
    if (this._isPaused) {
      this._needsFullRefresh = true;
      return;
    }
    if (this._coreService.decPrivateModes.synchronizedOutput) {
      this._syncOutputHandler.bufferRows(e, i8);
      return;
    }
    let n = this._syncOutputHandler.flush();
    n && (e = Math.min(e, n.start), i8 = Math.max(i8, n.end)), r5 || (this._isNextRenderRedrawOnly = false), this._renderDebouncer.refresh(e, i8, this._rowCount);
  }
  _renderRows(e, i8) {
    if (this._renderer.value) {
      if (this._coreService.decPrivateModes.synchronizedOutput) {
        this._syncOutputHandler.bufferRows(e, i8);
        return;
      }
      e = Math.min(e, this._rowCount - 1), i8 = Math.min(i8, this._rowCount - 1), this._renderer.value.renderRows(e, i8), this._needsSelectionRefresh && (this._renderer.value.handleSelectionChanged(this._selectionState.start, this._selectionState.end, this._selectionState.columnSelectMode), this._needsSelectionRefresh = false), this._isNextRenderRedrawOnly || this._onRenderedViewportChange.fire({ start: e, end: i8 }), this._onRender.fire({ start: e, end: i8 }), this._isNextRenderRedrawOnly = true;
    }
  }
  resize(e, i8) {
    this._rowCount = i8, this._fireOnCanvasResize();
  }
  _handleOptionsChanged() {
    this._renderer.value && (this.refreshRows(0, this._rowCount - 1), this._fireOnCanvasResize());
  }
  _fireOnCanvasResize() {
    this._renderer.value && (this._renderer.value.dimensions.css.canvas.width === this._canvasWidth && this._renderer.value.dimensions.css.canvas.height === this._canvasHeight || this._onDimensionsChange.fire(this._renderer.value.dimensions));
  }
  hasRenderer() {
    return !!this._renderer.value;
  }
  setRenderer(e) {
    this._renderer.value = e, this._renderer.value && (this._renderer.value.onRequestRedraw((i8) => this.refreshRows(i8.start, i8.end, true)), this._needsSelectionRefresh = true, this._fullRefresh());
  }
  addRefreshCallback(e) {
    return this._renderDebouncer.addRefreshCallback(e);
  }
  _fullRefresh() {
    this._isPaused ? this._needsFullRefresh = true : this.refreshRows(0, this._rowCount - 1);
  }
  clearTextureAtlas() {
    this._renderer.value && (this._renderer.value.clearTextureAtlas?.(), this._fullRefresh());
  }
  handleDevicePixelRatioChange() {
    this._charSizeService.measure(), this._renderer.value && (this._renderer.value.handleDevicePixelRatioChange(), this.refreshRows(0, this._rowCount - 1));
  }
  handleResize(e, i8) {
    this._renderer.value && (this._isPaused ? this._pausedResizeTask.set(() => this._renderer.value?.handleResize(e, i8)) : this._renderer.value.handleResize(e, i8), this._fullRefresh());
  }
  handleCharSizeChanged() {
    this._renderer.value?.handleCharSizeChanged();
  }
  handleBlur() {
    this._renderer.value?.handleBlur();
  }
  handleFocus() {
    this._renderer.value?.handleFocus();
  }
  handleSelectionChanged(e, i8, r5) {
    this._selectionState.start = e, this._selectionState.end = i8, this._selectionState.columnSelectMode = r5, this._renderer.value?.handleSelectionChanged(e, i8, r5);
  }
  handleCursorMove() {
    this._renderer.value?.handleCursorMove();
  }
  clear() {
    this._renderer.value?.clear();
  }
};
Qt2 = M3([S(2, H3), S(3, nt3), S(4, ge3), S(5, Be4), S(6, F2), S(7, ae2), S(8, Re2)], Qt2);
var xs2 = class {
  constructor(t, e, i8) {
    this._coreBrowserService = t;
    this._coreService = e;
    this._onTimeout = i8;
    this._start = 0;
    this._end = 0;
    this._isBuffering = false;
  }
  bufferRows(t, e) {
    this._isBuffering ? (this._start = Math.min(this._start, t), this._end = Math.max(this._end, e)) : (this._start = t, this._end = e, this._isBuffering = true), this._timeout === void 0 && (this._timeout = this._coreBrowserService.window.setTimeout(() => {
      this._timeout = void 0, this._coreService.decPrivateModes.synchronizedOutput = false, this._onTimeout();
    }, 1e3));
  }
  flush() {
    if (this._timeout !== void 0 && (this._coreBrowserService.window.clearTimeout(this._timeout), this._timeout = void 0), !this._isBuffering) return;
    let t = { start: this._start, end: this._end };
    return this._isBuffering = false, t;
  }
  dispose() {
    this._timeout !== void 0 && (this._coreBrowserService.window.clearTimeout(this._timeout), this._timeout = void 0);
  }
};
function Jo(s15, t, e, i8) {
  let r5 = e.buffer.x, n = e.buffer.y;
  if (!e.buffer.hasScrollback) return Ga2(r5, n, s15, t, e, i8) + sn2(n, t, e, i8) + $a2(r5, n, s15, t, e, i8);
  let o2;
  if (n === t) return o2 = r5 > s15 ? "D" : "C", Fi2(Math.abs(r5 - s15), Ni2(o2, i8));
  o2 = n > t ? "D" : "C";
  let l = Math.abs(n - t), a = za2(n > t ? s15 : r5, e) + (l - 1) * e.cols + 1 + Ka2(n > t ? r5 : s15, e);
  return Fi2(a, Ni2(o2, i8));
}
function Ka2(s15, t) {
  return s15 - 1;
}
function za2(s15, t) {
  return t.cols - s15;
}
function Ga2(s15, t, e, i8, r5, n) {
  return sn2(t, i8, r5, n).length === 0 ? "" : Fi2(el(s15, t, s15, t - gt3(t, r5), false, r5).length, Ni2("D", n));
}
function sn2(s15, t, e, i8) {
  let r5 = s15 - gt3(s15, e), n = t - gt3(t, e), o2 = Math.abs(r5 - n) - Va2(s15, t, e);
  return Fi2(o2, Ni2(Qo(s15, t), i8));
}
function $a2(s15, t, e, i8, r5, n) {
  let o2;
  sn2(t, i8, r5, n).length > 0 ? o2 = i8 - gt3(i8, r5) : o2 = t;
  let l = i8, a = qa2(s15, t, e, i8, r5, n);
  return Fi2(el(s15, o2, e, l, a === "C", r5).length, Ni2(a, n));
}
function Va2(s15, t, e) {
  let i8 = 0, r5 = s15 - gt3(s15, e), n = t - gt3(t, e);
  for (let o2 = 0; o2 < Math.abs(r5 - n); o2++) {
    let l = Qo(s15, t) === "A" ? -1 : 1;
    e.buffer.lines.get(r5 + l * o2)?.isWrapped && i8++;
  }
  return i8;
}
function gt3(s15, t) {
  let e = 0, i8 = t.buffer.lines.get(s15), r5 = i8?.isWrapped;
  for (; r5 && s15 >= 0 && s15 < t.rows; ) e++, i8 = t.buffer.lines.get(--s15), r5 = i8?.isWrapped;
  return e;
}
function qa2(s15, t, e, i8, r5, n) {
  let o2;
  return sn2(e, i8, r5, n).length > 0 ? o2 = i8 - gt3(i8, r5) : o2 = t, s15 < e && o2 <= i8 || s15 >= e && o2 < i8 ? "C" : "D";
}
function Qo(s15, t) {
  return s15 > t ? "A" : "B";
}
function el(s15, t, e, i8, r5, n) {
  let o2 = s15, l = t, a = "";
  for (; (o2 !== e || l !== i8) && l >= 0 && l < n.buffer.lines.length; ) o2 += r5 ? 1 : -1, r5 && o2 > n.cols - 1 ? (a += n.buffer.translateBufferLineToString(l, false, s15, o2), o2 = 0, s15 = 0, l++) : !r5 && o2 < 0 && (a += n.buffer.translateBufferLineToString(l, false, 0, s15 + 1), o2 = n.cols - 1, s15 = o2, l--);
  return a + n.buffer.translateBufferLineToString(l, false, s15, o2);
}
function Ni2(s15, t) {
  let e = t ? "O" : "[";
  return b.ESC + e + s15;
}
function Fi2(s15, t) {
  s15 = Math.floor(s15);
  let e = "";
  for (let i8 = 0; i8 < s15; i8++) e += t;
  return e;
}
var on2 = class {
  constructor(t) {
    this._bufferService = t;
    this.isSelectAllActive = false;
    this.selectionStartLength = 0;
  }
  clearSelection() {
    this.selectionStart = void 0, this.selectionEnd = void 0, this.isSelectAllActive = false, this.selectionStartLength = 0;
  }
  get finalSelectionStart() {
    return this.isSelectAllActive ? [0, 0] : !this.selectionEnd || !this.selectionStart ? this.selectionStart : this.areSelectionValuesReversed() ? this.selectionEnd : this.selectionStart;
  }
  get finalSelectionEnd() {
    if (this.isSelectAllActive) return [this._bufferService.cols, this._bufferService.buffer.ybase + this._bufferService.rows - 1];
    if (this.selectionStart) {
      if (!this.selectionEnd || this.areSelectionValuesReversed()) {
        let t = this.selectionStart[0] + this.selectionStartLength;
        return t > this._bufferService.cols ? t % this._bufferService.cols === 0 ? [this._bufferService.cols, this.selectionStart[1] + Math.floor(t / this._bufferService.cols) - 1] : [t % this._bufferService.cols, this.selectionStart[1] + Math.floor(t / this._bufferService.cols)] : [t, this.selectionStart[1]];
      }
      if (this.selectionStartLength && this.selectionEnd[1] === this.selectionStart[1]) {
        let t = this.selectionStart[0] + this.selectionStartLength;
        return t > this._bufferService.cols ? [t % this._bufferService.cols, this.selectionStart[1] + Math.floor(t / this._bufferService.cols)] : [Math.max(t, this.selectionEnd[0]), this.selectionEnd[1]];
      }
      return this.selectionEnd;
    }
  }
  areSelectionValuesReversed() {
    let t = this.selectionStart, e = this.selectionEnd;
    return !t || !e ? false : t[1] > e[1] || t[1] === e[1] && t[0] > e[0];
  }
  handleTrim(t) {
    return this.selectionStart && (this.selectionStart[1] -= t), this.selectionEnd && (this.selectionEnd[1] -= t), this.selectionEnd && this.selectionEnd[1] < 0 ? (this.clearSelection(), true) : (this.selectionStart && this.selectionStart[1] < 0 && (this.selectionStart[1] = 0), false);
  }
};
function ws2(s15, t) {
  if (s15.start.y > s15.end.y) throw new Error(`Buffer range end (${s15.end.x}, ${s15.end.y}) cannot be before start (${s15.start.x}, ${s15.start.y})`);
  return t * (s15.end.y - s15.start.y) + (s15.end.x - s15.start.x + 1);
}
var Ds = 50;
var Ya = 15;
var ja2 = 50;
var Xa = 500;
var Za = " ";
var Ja = new RegExp(Za, "g");
var ei2 = class extends D2 {
  constructor(e, i8, r5, n, o2, l, a, u, h2) {
    super();
    this._element = e;
    this._screenElement = i8;
    this._linkifier = r5;
    this._bufferService = n;
    this._coreService = o2;
    this._mouseService = l;
    this._optionsService = a;
    this._renderService = u;
    this._coreBrowserService = h2;
    this._dragScrollAmount = 0;
    this._enabled = true;
    this._workCell = new q();
    this._mouseDownTimeStamp = 0;
    this._oldHasSelection = false;
    this._oldSelectionStart = void 0;
    this._oldSelectionEnd = void 0;
    this._onLinuxMouseSelection = this._register(new v());
    this.onLinuxMouseSelection = this._onLinuxMouseSelection.event;
    this._onRedrawRequest = this._register(new v());
    this.onRequestRedraw = this._onRedrawRequest.event;
    this._onSelectionChange = this._register(new v());
    this.onSelectionChange = this._onSelectionChange.event;
    this._onRequestScrollLines = this._register(new v());
    this.onRequestScrollLines = this._onRequestScrollLines.event;
    this._mouseMoveListener = (c) => this._handleMouseMove(c), this._mouseUpListener = (c) => this._handleMouseUp(c), this._coreService.onUserInput(() => {
      this.hasSelection && this.clearSelection();
    }), this._trimListener = this._bufferService.buffer.lines.onTrim((c) => this._handleTrim(c)), this._register(this._bufferService.buffers.onBufferActivate((c) => this._handleBufferActivate(c))), this.enable(), this._model = new on2(this._bufferService), this._activeSelectionMode = 0, this._register(C3(() => {
      this._removeMouseDownListeners();
    })), this._register(this._bufferService.onResize((c) => {
      c.rowsChanged && this.clearSelection();
    }));
  }
  reset() {
    this.clearSelection();
  }
  disable() {
    this.clearSelection(), this._enabled = false;
  }
  enable() {
    this._enabled = true;
  }
  get selectionStart() {
    return this._model.finalSelectionStart;
  }
  get selectionEnd() {
    return this._model.finalSelectionEnd;
  }
  get hasSelection() {
    let e = this._model.finalSelectionStart, i8 = this._model.finalSelectionEnd;
    return !e || !i8 ? false : e[0] !== i8[0] || e[1] !== i8[1];
  }
  get selectionText() {
    let e = this._model.finalSelectionStart, i8 = this._model.finalSelectionEnd;
    if (!e || !i8) return "";
    let r5 = this._bufferService.buffer, n = [];
    if (this._activeSelectionMode === 3) {
      if (e[0] === i8[0]) return "";
      let l = e[0] < i8[0] ? e[0] : i8[0], a = e[0] < i8[0] ? i8[0] : e[0];
      for (let u = e[1]; u <= i8[1]; u++) {
        let h2 = r5.translateBufferLineToString(u, true, l, a);
        n.push(h2);
      }
    } else {
      let l = e[1] === i8[1] ? i8[0] : void 0;
      n.push(r5.translateBufferLineToString(e[1], true, e[0], l));
      for (let a = e[1] + 1; a <= i8[1] - 1; a++) {
        let u = r5.lines.get(a), h2 = r5.translateBufferLineToString(a, true);
        u?.isWrapped ? n[n.length - 1] += h2 : n.push(h2);
      }
      if (e[1] !== i8[1]) {
        let a = r5.lines.get(i8[1]), u = r5.translateBufferLineToString(i8[1], true, 0, i8[0]);
        a && a.isWrapped ? n[n.length - 1] += u : n.push(u);
      }
    }
    return n.map((l) => l.replace(Ja, " ")).join(Es2 ? `\r
` : `
`);
  }
  clearSelection() {
    this._model.clearSelection(), this._removeMouseDownListeners(), this.refresh(), this._onSelectionChange.fire();
  }
  refresh(e) {
    this._refreshAnimationFrame || (this._refreshAnimationFrame = this._coreBrowserService.window.requestAnimationFrame(() => this._refresh())), Bi2 && e && this.selectionText.length && this._onLinuxMouseSelection.fire(this.selectionText);
  }
  _refresh() {
    this._refreshAnimationFrame = void 0, this._onRedrawRequest.fire({ start: this._model.finalSelectionStart, end: this._model.finalSelectionEnd, columnSelectMode: this._activeSelectionMode === 3 });
  }
  _isClickInSelection(e) {
    let i8 = this._getMouseBufferCoords(e), r5 = this._model.finalSelectionStart, n = this._model.finalSelectionEnd;
    return !r5 || !n || !i8 ? false : this._areCoordsInSelection(i8, r5, n);
  }
  isCellInSelection(e, i8) {
    let r5 = this._model.finalSelectionStart, n = this._model.finalSelectionEnd;
    return !r5 || !n ? false : this._areCoordsInSelection([e, i8], r5, n);
  }
  _areCoordsInSelection(e, i8, r5) {
    return e[1] > i8[1] && e[1] < r5[1] || i8[1] === r5[1] && e[1] === i8[1] && e[0] >= i8[0] && e[0] < r5[0] || i8[1] < r5[1] && e[1] === r5[1] && e[0] < r5[0] || i8[1] < r5[1] && e[1] === i8[1] && e[0] >= i8[0];
  }
  _selectWordAtCursor(e, i8) {
    let r5 = this._linkifier.currentLink?.link?.range;
    if (r5) return this._model.selectionStart = [r5.start.x - 1, r5.start.y - 1], this._model.selectionStartLength = ws2(r5, this._bufferService.cols), this._model.selectionEnd = void 0, true;
    let n = this._getMouseBufferCoords(e);
    return n ? (this._selectWordAt(n, i8), this._model.selectionEnd = void 0, true) : false;
  }
  selectAll() {
    this._model.isSelectAllActive = true, this.refresh(), this._onSelectionChange.fire();
  }
  selectLines(e, i8) {
    this._model.clearSelection(), e = Math.max(e, 0), i8 = Math.min(i8, this._bufferService.buffer.lines.length - 1), this._model.selectionStart = [0, e], this._model.selectionEnd = [this._bufferService.cols, i8], this.refresh(), this._onSelectionChange.fire();
  }
  _handleTrim(e) {
    this._model.handleTrim(e) && this.refresh();
  }
  _getMouseBufferCoords(e) {
    let i8 = this._mouseService.getCoords(e, this._screenElement, this._bufferService.cols, this._bufferService.rows, true);
    if (i8) return i8[0]--, i8[1]--, i8[1] += this._bufferService.buffer.ydisp, i8;
  }
  _getMouseEventScrollAmount(e) {
    let i8 = Ci2(this._coreBrowserService.window, e, this._screenElement)[1], r5 = this._renderService.dimensions.css.canvas.height;
    return i8 >= 0 && i8 <= r5 ? 0 : (i8 > r5 && (i8 -= r5), i8 = Math.min(Math.max(i8, -Ds), Ds), i8 /= Ds, i8 / Math.abs(i8) + Math.round(i8 * (Ya - 1)));
  }
  shouldForceSelection(e) {
    return Zt3 ? e.altKey && this._optionsService.rawOptions.macOptionClickForcesSelection : e.shiftKey;
  }
  handleMouseDown(e) {
    if (this._mouseDownTimeStamp = e.timeStamp, !(e.button === 2 && this.hasSelection) && e.button === 0) {
      if (!this._enabled) {
        if (!this.shouldForceSelection(e)) return;
        e.stopPropagation();
      }
      e.preventDefault(), this._dragScrollAmount = 0, this._enabled && e.shiftKey ? this._handleIncrementalClick(e) : e.detail === 1 ? this._handleSingleClick(e) : e.detail === 2 ? this._handleDoubleClick(e) : e.detail === 3 && this._handleTripleClick(e), this._addMouseDownListeners(), this.refresh(true);
    }
  }
  _addMouseDownListeners() {
    this._screenElement.ownerDocument && (this._screenElement.ownerDocument.addEventListener("mousemove", this._mouseMoveListener), this._screenElement.ownerDocument.addEventListener("mouseup", this._mouseUpListener)), this._dragScrollIntervalTimer = this._coreBrowserService.window.setInterval(() => this._dragScroll(), ja2);
  }
  _removeMouseDownListeners() {
    this._screenElement.ownerDocument && (this._screenElement.ownerDocument.removeEventListener("mousemove", this._mouseMoveListener), this._screenElement.ownerDocument.removeEventListener("mouseup", this._mouseUpListener)), this._coreBrowserService.window.clearInterval(this._dragScrollIntervalTimer), this._dragScrollIntervalTimer = void 0;
  }
  _handleIncrementalClick(e) {
    this._model.selectionStart && (this._model.selectionEnd = this._getMouseBufferCoords(e));
  }
  _handleSingleClick(e) {
    if (this._model.selectionStartLength = 0, this._model.isSelectAllActive = false, this._activeSelectionMode = this.shouldColumnSelect(e) ? 3 : 0, this._model.selectionStart = this._getMouseBufferCoords(e), !this._model.selectionStart) return;
    this._model.selectionEnd = void 0;
    let i8 = this._bufferService.buffer.lines.get(this._model.selectionStart[1]);
    i8 && i8.length !== this._model.selectionStart[0] && i8.hasWidth(this._model.selectionStart[0]) === 0 && this._model.selectionStart[0]++;
  }
  _handleDoubleClick(e) {
    this._selectWordAtCursor(e, true) && (this._activeSelectionMode = 1);
  }
  _handleTripleClick(e) {
    let i8 = this._getMouseBufferCoords(e);
    i8 && (this._activeSelectionMode = 2, this._selectLineAt(i8[1]));
  }
  shouldColumnSelect(e) {
    return e.altKey && !(Zt3 && this._optionsService.rawOptions.macOptionClickForcesSelection);
  }
  _handleMouseMove(e) {
    if (e.stopImmediatePropagation(), !this._model.selectionStart) return;
    let i8 = this._model.selectionEnd ? [this._model.selectionEnd[0], this._model.selectionEnd[1]] : null;
    if (this._model.selectionEnd = this._getMouseBufferCoords(e), !this._model.selectionEnd) {
      this.refresh(true);
      return;
    }
    this._activeSelectionMode === 2 ? this._model.selectionEnd[1] < this._model.selectionStart[1] ? this._model.selectionEnd[0] = 0 : this._model.selectionEnd[0] = this._bufferService.cols : this._activeSelectionMode === 1 && this._selectToWordAt(this._model.selectionEnd), this._dragScrollAmount = this._getMouseEventScrollAmount(e), this._activeSelectionMode !== 3 && (this._dragScrollAmount > 0 ? this._model.selectionEnd[0] = this._bufferService.cols : this._dragScrollAmount < 0 && (this._model.selectionEnd[0] = 0));
    let r5 = this._bufferService.buffer;
    if (this._model.selectionEnd[1] < r5.lines.length) {
      let n = r5.lines.get(this._model.selectionEnd[1]);
      n && n.hasWidth(this._model.selectionEnd[0]) === 0 && this._model.selectionEnd[0] < this._bufferService.cols && this._model.selectionEnd[0]++;
    }
    (!i8 || i8[0] !== this._model.selectionEnd[0] || i8[1] !== this._model.selectionEnd[1]) && this.refresh(true);
  }
  _dragScroll() {
    if (!(!this._model.selectionEnd || !this._model.selectionStart) && this._dragScrollAmount) {
      this._onRequestScrollLines.fire({ amount: this._dragScrollAmount, suppressScrollEvent: false });
      let e = this._bufferService.buffer;
      this._dragScrollAmount > 0 ? (this._activeSelectionMode !== 3 && (this._model.selectionEnd[0] = this._bufferService.cols), this._model.selectionEnd[1] = Math.min(e.ydisp + this._bufferService.rows, e.lines.length - 1)) : (this._activeSelectionMode !== 3 && (this._model.selectionEnd[0] = 0), this._model.selectionEnd[1] = e.ydisp), this.refresh();
    }
  }
  _handleMouseUp(e) {
    let i8 = e.timeStamp - this._mouseDownTimeStamp;
    if (this._removeMouseDownListeners(), this.selectionText.length <= 1 && i8 < Xa && e.altKey && this._optionsService.rawOptions.altClickMovesCursor) {
      if (this._bufferService.buffer.ybase === this._bufferService.buffer.ydisp) {
        let r5 = this._mouseService.getCoords(e, this._element, this._bufferService.cols, this._bufferService.rows, false);
        if (r5 && r5[0] !== void 0 && r5[1] !== void 0) {
          let n = Jo(r5[0] - 1, r5[1] - 1, this._bufferService, this._coreService.decPrivateModes.applicationCursorKeys);
          this._coreService.triggerDataEvent(n, true);
        }
      }
    } else this._fireEventIfSelectionChanged();
  }
  _fireEventIfSelectionChanged() {
    let e = this._model.finalSelectionStart, i8 = this._model.finalSelectionEnd, r5 = !!e && !!i8 && (e[0] !== i8[0] || e[1] !== i8[1]);
    if (!r5) {
      this._oldHasSelection && this._fireOnSelectionChange(e, i8, r5);
      return;
    }
    !e || !i8 || (!this._oldSelectionStart || !this._oldSelectionEnd || e[0] !== this._oldSelectionStart[0] || e[1] !== this._oldSelectionStart[1] || i8[0] !== this._oldSelectionEnd[0] || i8[1] !== this._oldSelectionEnd[1]) && this._fireOnSelectionChange(e, i8, r5);
  }
  _fireOnSelectionChange(e, i8, r5) {
    this._oldSelectionStart = e, this._oldSelectionEnd = i8, this._oldHasSelection = r5, this._onSelectionChange.fire();
  }
  _handleBufferActivate(e) {
    this.clearSelection(), this._trimListener.dispose(), this._trimListener = e.activeBuffer.lines.onTrim((i8) => this._handleTrim(i8));
  }
  _convertViewportColToCharacterIndex(e, i8) {
    let r5 = i8;
    for (let n = 0; i8 >= n; n++) {
      let o2 = e.loadCell(n, this._workCell).getChars().length;
      this._workCell.getWidth() === 0 ? r5-- : o2 > 1 && i8 !== n && (r5 += o2 - 1);
    }
    return r5;
  }
  setSelection(e, i8, r5) {
    this._model.clearSelection(), this._removeMouseDownListeners(), this._model.selectionStart = [e, i8], this._model.selectionStartLength = r5, this.refresh(), this._fireEventIfSelectionChanged();
  }
  rightClickSelect(e) {
    this._isClickInSelection(e) || (this._selectWordAtCursor(e, false) && this.refresh(true), this._fireEventIfSelectionChanged());
  }
  _getWordAt(e, i8, r5 = true, n = true) {
    if (e[0] >= this._bufferService.cols) return;
    let o2 = this._bufferService.buffer, l = o2.lines.get(e[1]);
    if (!l) return;
    let a = o2.translateBufferLineToString(e[1], false), u = this._convertViewportColToCharacterIndex(l, e[0]), h2 = u, c = e[0] - u, d = 0, _3 = 0, p = 0, m = 0;
    if (a.charAt(u) === " ") {
      for (; u > 0 && a.charAt(u - 1) === " "; ) u--;
      for (; h2 < a.length && a.charAt(h2 + 1) === " "; ) h2++;
    } else {
      let R3 = e[0], O3 = e[0];
      l.getWidth(R3) === 0 && (d++, R3--), l.getWidth(O3) === 2 && (_3++, O3++);
      let I2 = l.getString(O3).length;
      for (I2 > 1 && (m += I2 - 1, h2 += I2 - 1); R3 > 0 && u > 0 && !this._isCharWordSeparator(l.loadCell(R3 - 1, this._workCell)); ) {
        l.loadCell(R3 - 1, this._workCell);
        let k2 = this._workCell.getChars().length;
        this._workCell.getWidth() === 0 ? (d++, R3--) : k2 > 1 && (p += k2 - 1, u -= k2 - 1), u--, R3--;
      }
      for (; O3 < l.length && h2 + 1 < a.length && !this._isCharWordSeparator(l.loadCell(O3 + 1, this._workCell)); ) {
        l.loadCell(O3 + 1, this._workCell);
        let k2 = this._workCell.getChars().length;
        this._workCell.getWidth() === 2 ? (_3++, O3++) : k2 > 1 && (m += k2 - 1, h2 += k2 - 1), h2++, O3++;
      }
    }
    h2++;
    let f = u + c - d + p, A2 = Math.min(this._bufferService.cols, h2 - u + d + _3 - p - m);
    if (!(!i8 && a.slice(u, h2).trim() === "")) {
      if (r5 && f === 0 && l.getCodePoint(0) !== 32) {
        let R3 = o2.lines.get(e[1] - 1);
        if (R3 && l.isWrapped && R3.getCodePoint(this._bufferService.cols - 1) !== 32) {
          let O3 = this._getWordAt([this._bufferService.cols - 1, e[1] - 1], false, true, false);
          if (O3) {
            let I2 = this._bufferService.cols - O3.start;
            f -= I2, A2 += I2;
          }
        }
      }
      if (n && f + A2 === this._bufferService.cols && l.getCodePoint(this._bufferService.cols - 1) !== 32) {
        let R3 = o2.lines.get(e[1] + 1);
        if (R3?.isWrapped && R3.getCodePoint(0) !== 32) {
          let O3 = this._getWordAt([0, e[1] + 1], false, false, true);
          O3 && (A2 += O3.length);
        }
      }
      return { start: f, length: A2 };
    }
  }
  _selectWordAt(e, i8) {
    let r5 = this._getWordAt(e, i8);
    if (r5) {
      for (; r5.start < 0; ) r5.start += this._bufferService.cols, e[1]--;
      this._model.selectionStart = [r5.start, e[1]], this._model.selectionStartLength = r5.length;
    }
  }
  _selectToWordAt(e) {
    let i8 = this._getWordAt(e, true);
    if (i8) {
      let r5 = e[1];
      for (; i8.start < 0; ) i8.start += this._bufferService.cols, r5--;
      if (!this._model.areSelectionValuesReversed()) for (; i8.start + i8.length > this._bufferService.cols; ) i8.length -= this._bufferService.cols, r5++;
      this._model.selectionEnd = [this._model.areSelectionValuesReversed() ? i8.start : i8.start + i8.length, r5];
    }
  }
  _isCharWordSeparator(e) {
    return e.getWidth() === 0 ? false : this._optionsService.rawOptions.wordSeparator.indexOf(e.getChars()) >= 0;
  }
  _selectLineAt(e) {
    let i8 = this._bufferService.buffer.getWrappedRangeForLine(e), r5 = { start: { x: 0, y: i8.first }, end: { x: this._bufferService.cols - 1, y: i8.last } };
    this._model.selectionStart = [0, i8.first], this._model.selectionEnd = void 0, this._model.selectionStartLength = ws2(r5, this._bufferService.cols);
  }
};
ei2 = M3([S(3, F2), S(4, ge3), S(5, Dt2), S(6, H3), S(7, ce2), S(8, ae2)], ei2);
var Hi2 = class {
  constructor() {
    this._data = {};
  }
  set(t, e, i8) {
    this._data[t] || (this._data[t] = {}), this._data[t][e] = i8;
  }
  get(t, e) {
    return this._data[t] ? this._data[t][e] : void 0;
  }
  clear() {
    this._data = {};
  }
};
var Wi2 = class {
  constructor() {
    this._color = new Hi2();
    this._css = new Hi2();
  }
  setCss(t, e, i8) {
    this._css.set(t, e, i8);
  }
  getCss(t, e) {
    return this._css.get(t, e);
  }
  setColor(t, e, i8) {
    this._color.set(t, e, i8);
  }
  getColor(t, e) {
    return this._color.get(t, e);
  }
  clear() {
    this._color.clear(), this._css.clear();
  }
};
var re3 = Object.freeze((() => {
  let s15 = [z2.toColor("#2e3436"), z2.toColor("#cc0000"), z2.toColor("#4e9a06"), z2.toColor("#c4a000"), z2.toColor("#3465a4"), z2.toColor("#75507b"), z2.toColor("#06989a"), z2.toColor("#d3d7cf"), z2.toColor("#555753"), z2.toColor("#ef2929"), z2.toColor("#8ae234"), z2.toColor("#fce94f"), z2.toColor("#729fcf"), z2.toColor("#ad7fa8"), z2.toColor("#34e2e2"), z2.toColor("#eeeeec")], t = [0, 95, 135, 175, 215, 255];
  for (let e = 0; e < 216; e++) {
    let i8 = t[e / 36 % 6 | 0], r5 = t[e / 6 % 6 | 0], n = t[e % 6];
    s15.push({ css: j2.toCss(i8, r5, n), rgba: j2.toRgba(i8, r5, n) });
  }
  for (let e = 0; e < 24; e++) {
    let i8 = 8 + e * 10;
    s15.push({ css: j2.toCss(i8, i8, i8), rgba: j2.toRgba(i8, i8, i8) });
  }
  return s15;
})());
var St2 = z2.toColor("#ffffff");
var Ki2 = z2.toColor("#000000");
var tl = z2.toColor("#ffffff");
var il = Ki2;
var Ui2 = { css: "rgba(255, 255, 255, 0.3)", rgba: 4294967117 };
var Qa = St2;
var ti2 = class extends D2 {
  constructor(e) {
    super();
    this._optionsService = e;
    this._contrastCache = new Wi2();
    this._halfContrastCache = new Wi2();
    this._onChangeColors = this._register(new v());
    this.onChangeColors = this._onChangeColors.event;
    this._colors = { foreground: St2, background: Ki2, cursor: tl, cursorAccent: il, selectionForeground: void 0, selectionBackgroundTransparent: Ui2, selectionBackgroundOpaque: U3.blend(Ki2, Ui2), selectionInactiveBackgroundTransparent: Ui2, selectionInactiveBackgroundOpaque: U3.blend(Ki2, Ui2), scrollbarSliderBackground: U3.opacity(St2, 0.2), scrollbarSliderHoverBackground: U3.opacity(St2, 0.4), scrollbarSliderActiveBackground: U3.opacity(St2, 0.5), overviewRulerBorder: St2, ansi: re3.slice(), contrastCache: this._contrastCache, halfContrastCache: this._halfContrastCache }, this._updateRestoreColors(), this._setTheme(this._optionsService.rawOptions.theme), this._register(this._optionsService.onSpecificOptionChange("minimumContrastRatio", () => this._contrastCache.clear())), this._register(this._optionsService.onSpecificOptionChange("theme", () => this._setTheme(this._optionsService.rawOptions.theme)));
  }
  get colors() {
    return this._colors;
  }
  _setTheme(e = {}) {
    let i8 = this._colors;
    if (i8.foreground = K3(e.foreground, St2), i8.background = K3(e.background, Ki2), i8.cursor = U3.blend(i8.background, K3(e.cursor, tl)), i8.cursorAccent = U3.blend(i8.background, K3(e.cursorAccent, il)), i8.selectionBackgroundTransparent = K3(e.selectionBackground, Ui2), i8.selectionBackgroundOpaque = U3.blend(i8.background, i8.selectionBackgroundTransparent), i8.selectionInactiveBackgroundTransparent = K3(e.selectionInactiveBackground, i8.selectionBackgroundTransparent), i8.selectionInactiveBackgroundOpaque = U3.blend(i8.background, i8.selectionInactiveBackgroundTransparent), i8.selectionForeground = e.selectionForeground ? K3(e.selectionForeground, ps2) : void 0, i8.selectionForeground === ps2 && (i8.selectionForeground = void 0), U3.isOpaque(i8.selectionBackgroundTransparent) && (i8.selectionBackgroundTransparent = U3.opacity(i8.selectionBackgroundTransparent, 0.3)), U3.isOpaque(i8.selectionInactiveBackgroundTransparent) && (i8.selectionInactiveBackgroundTransparent = U3.opacity(i8.selectionInactiveBackgroundTransparent, 0.3)), i8.scrollbarSliderBackground = K3(e.scrollbarSliderBackground, U3.opacity(i8.foreground, 0.2)), i8.scrollbarSliderHoverBackground = K3(e.scrollbarSliderHoverBackground, U3.opacity(i8.foreground, 0.4)), i8.scrollbarSliderActiveBackground = K3(e.scrollbarSliderActiveBackground, U3.opacity(i8.foreground, 0.5)), i8.overviewRulerBorder = K3(e.overviewRulerBorder, Qa), i8.ansi = re3.slice(), i8.ansi[0] = K3(e.black, re3[0]), i8.ansi[1] = K3(e.red, re3[1]), i8.ansi[2] = K3(e.green, re3[2]), i8.ansi[3] = K3(e.yellow, re3[3]), i8.ansi[4] = K3(e.blue, re3[4]), i8.ansi[5] = K3(e.magenta, re3[5]), i8.ansi[6] = K3(e.cyan, re3[6]), i8.ansi[7] = K3(e.white, re3[7]), i8.ansi[8] = K3(e.brightBlack, re3[8]), i8.ansi[9] = K3(e.brightRed, re3[9]), i8.ansi[10] = K3(e.brightGreen, re3[10]), i8.ansi[11] = K3(e.brightYellow, re3[11]), i8.ansi[12] = K3(e.brightBlue, re3[12]), i8.ansi[13] = K3(e.brightMagenta, re3[13]), i8.ansi[14] = K3(e.brightCyan, re3[14]), i8.ansi[15] = K3(e.brightWhite, re3[15]), e.extendedAnsi) {
      let r5 = Math.min(i8.ansi.length - 16, e.extendedAnsi.length);
      for (let n = 0; n < r5; n++) i8.ansi[n + 16] = K3(e.extendedAnsi[n], re3[n + 16]);
    }
    this._contrastCache.clear(), this._halfContrastCache.clear(), this._updateRestoreColors(), this._onChangeColors.fire(this.colors);
  }
  restoreColor(e) {
    this._restoreColor(e), this._onChangeColors.fire(this.colors);
  }
  _restoreColor(e) {
    if (e === void 0) {
      for (let i8 = 0; i8 < this._restoreColors.ansi.length; ++i8) this._colors.ansi[i8] = this._restoreColors.ansi[i8];
      return;
    }
    switch (e) {
      case 256:
        this._colors.foreground = this._restoreColors.foreground;
        break;
      case 257:
        this._colors.background = this._restoreColors.background;
        break;
      case 258:
        this._colors.cursor = this._restoreColors.cursor;
        break;
      default:
        this._colors.ansi[e] = this._restoreColors.ansi[e];
    }
  }
  modifyColors(e) {
    e(this._colors), this._onChangeColors.fire(this.colors);
  }
  _updateRestoreColors() {
    this._restoreColors = { foreground: this._colors.foreground, background: this._colors.background, cursor: this._colors.cursor, ansi: this._colors.ansi.slice() };
  }
};
ti2 = M3([S(0, H3)], ti2);
function K3(s15, t) {
  if (s15 !== void 0) try {
    return z2.toColor(s15);
  } catch {
  }
  return t;
}
var Rs2 = class {
  constructor(...t) {
    this._entries = /* @__PURE__ */ new Map();
    for (let [e, i8] of t) this.set(e, i8);
  }
  set(t, e) {
    let i8 = this._entries.get(t);
    return this._entries.set(t, e), i8;
  }
  forEach(t) {
    for (let [e, i8] of this._entries.entries()) t(e, i8);
  }
  has(t) {
    return this._entries.has(t);
  }
  get(t) {
    return this._entries.get(t);
  }
};
var ln = class {
  constructor() {
    this._services = new Rs2();
    this._services.set(xt4, this);
  }
  setService(t, e) {
    this._services.set(t, e);
  }
  getService(t) {
    return this._services.get(t);
  }
  createInstance(t, ...e) {
    let i8 = Xs(t).sort((o2, l) => o2.index - l.index), r5 = [];
    for (let o2 of i8) {
      let l = this._services.get(o2.id);
      if (!l) throw new Error(`[createInstance] ${t.name} depends on UNKNOWN service ${o2.id._id}.`);
      r5.push(l);
    }
    let n = i8.length > 0 ? i8[0].index : e.length;
    if (e.length !== n) throw new Error(`[createInstance] First service dependency of ${t.name} at position ${n + 1} conflicts with ${e.length} static arguments`);
    return new t(...e, ...r5);
  }
};
var ec = { trace: 0, debug: 1, info: 2, warn: 3, error: 4, off: 5 };
var tc = "xterm.js: ";
var ii2 = class extends D2 {
  constructor(e) {
    super();
    this._optionsService = e;
    this._logLevel = 5;
    this._updateLogLevel(), this._register(this._optionsService.onSpecificOptionChange("logLevel", () => this._updateLogLevel())), ic = this;
  }
  get logLevel() {
    return this._logLevel;
  }
  _updateLogLevel() {
    this._logLevel = ec[this._optionsService.rawOptions.logLevel];
  }
  _evalLazyOptionalParams(e) {
    for (let i8 = 0; i8 < e.length; i8++) typeof e[i8] == "function" && (e[i8] = e[i8]());
  }
  _log(e, i8, r5) {
    this._evalLazyOptionalParams(r5), e.call(console, (this._optionsService.options.logger ? "" : tc) + i8, ...r5);
  }
  trace(e, ...i8) {
    this._logLevel <= 0 && this._log(this._optionsService.options.logger?.trace.bind(this._optionsService.options.logger) ?? console.log, e, i8);
  }
  debug(e, ...i8) {
    this._logLevel <= 1 && this._log(this._optionsService.options.logger?.debug.bind(this._optionsService.options.logger) ?? console.log, e, i8);
  }
  info(e, ...i8) {
    this._logLevel <= 2 && this._log(this._optionsService.options.logger?.info.bind(this._optionsService.options.logger) ?? console.info, e, i8);
  }
  warn(e, ...i8) {
    this._logLevel <= 3 && this._log(this._optionsService.options.logger?.warn.bind(this._optionsService.options.logger) ?? console.warn, e, i8);
  }
  error(e, ...i8) {
    this._logLevel <= 4 && this._log(this._optionsService.options.logger?.error.bind(this._optionsService.options.logger) ?? console.error, e, i8);
  }
};
ii2 = M3([S(0, H3)], ii2);
var ic;
var zi2 = class extends D2 {
  constructor(e) {
    super();
    this._maxLength = e;
    this.onDeleteEmitter = this._register(new v());
    this.onDelete = this.onDeleteEmitter.event;
    this.onInsertEmitter = this._register(new v());
    this.onInsert = this.onInsertEmitter.event;
    this.onTrimEmitter = this._register(new v());
    this.onTrim = this.onTrimEmitter.event;
    this._array = new Array(this._maxLength), this._startIndex = 0, this._length = 0;
  }
  get maxLength() {
    return this._maxLength;
  }
  set maxLength(e) {
    if (this._maxLength === e) return;
    let i8 = new Array(e);
    for (let r5 = 0; r5 < Math.min(e, this.length); r5++) i8[r5] = this._array[this._getCyclicIndex(r5)];
    this._array = i8, this._maxLength = e, this._startIndex = 0;
  }
  get length() {
    return this._length;
  }
  set length(e) {
    if (e > this._length) for (let i8 = this._length; i8 < e; i8++) this._array[i8] = void 0;
    this._length = e;
  }
  get(e) {
    return this._array[this._getCyclicIndex(e)];
  }
  set(e, i8) {
    this._array[this._getCyclicIndex(e)] = i8;
  }
  push(e) {
    this._array[this._getCyclicIndex(this._length)] = e, this._length === this._maxLength ? (this._startIndex = ++this._startIndex % this._maxLength, this.onTrimEmitter.fire(1)) : this._length++;
  }
  recycle() {
    if (this._length !== this._maxLength) throw new Error("Can only recycle when the buffer is full");
    return this._startIndex = ++this._startIndex % this._maxLength, this.onTrimEmitter.fire(1), this._array[this._getCyclicIndex(this._length - 1)];
  }
  get isFull() {
    return this._length === this._maxLength;
  }
  pop() {
    return this._array[this._getCyclicIndex(this._length-- - 1)];
  }
  splice(e, i8, ...r5) {
    if (i8) {
      for (let n = e; n < this._length - i8; n++) this._array[this._getCyclicIndex(n)] = this._array[this._getCyclicIndex(n + i8)];
      this._length -= i8, this.onDeleteEmitter.fire({ index: e, amount: i8 });
    }
    for (let n = this._length - 1; n >= e; n--) this._array[this._getCyclicIndex(n + r5.length)] = this._array[this._getCyclicIndex(n)];
    for (let n = 0; n < r5.length; n++) this._array[this._getCyclicIndex(e + n)] = r5[n];
    if (r5.length && this.onInsertEmitter.fire({ index: e, amount: r5.length }), this._length + r5.length > this._maxLength) {
      let n = this._length + r5.length - this._maxLength;
      this._startIndex += n, this._length = this._maxLength, this.onTrimEmitter.fire(n);
    } else this._length += r5.length;
  }
  trimStart(e) {
    e > this._length && (e = this._length), this._startIndex += e, this._length -= e, this.onTrimEmitter.fire(e);
  }
  shiftElements(e, i8, r5) {
    if (!(i8 <= 0)) {
      if (e < 0 || e >= this._length) throw new Error("start argument out of range");
      if (e + r5 < 0) throw new Error("Cannot shift elements in list beyond index 0");
      if (r5 > 0) {
        for (let o2 = i8 - 1; o2 >= 0; o2--) this.set(e + o2 + r5, this.get(e + o2));
        let n = e + i8 + r5 - this._length;
        if (n > 0) for (this._length += n; this._length > this._maxLength; ) this._length--, this._startIndex++, this.onTrimEmitter.fire(1);
      } else for (let n = 0; n < i8; n++) this.set(e + n + r5, this.get(e + n));
    }
  }
  _getCyclicIndex(e) {
    return (this._startIndex + e) % this._maxLength;
  }
};
var B2 = 3;
var X4 = Object.freeze(new De3());
var an2 = 0;
var Ls = 2;
var Ze3 = class s12 {
  constructor(t, e, i8 = false) {
    this.isWrapped = i8;
    this._combined = {};
    this._extendedAttrs = {};
    this._data = new Uint32Array(t * B2);
    let r5 = e || q.fromCharData([0, ir, 1, 0]);
    for (let n = 0; n < t; ++n) this.setCell(n, r5);
    this.length = t;
  }
  get(t) {
    let e = this._data[t * B2 + 0], i8 = e & 2097151;
    return [this._data[t * B2 + 1], e & 2097152 ? this._combined[t] : i8 ? Ce2(i8) : "", e >> 22, e & 2097152 ? this._combined[t].charCodeAt(this._combined[t].length - 1) : i8];
  }
  set(t, e) {
    this._data[t * B2 + 1] = e[0], e[1].length > 1 ? (this._combined[t] = e[1], this._data[t * B2 + 0] = t | 2097152 | e[2] << 22) : this._data[t * B2 + 0] = e[1].charCodeAt(0) | e[2] << 22;
  }
  getWidth(t) {
    return this._data[t * B2 + 0] >> 22;
  }
  hasWidth(t) {
    return this._data[t * B2 + 0] & 12582912;
  }
  getFg(t) {
    return this._data[t * B2 + 1];
  }
  getBg(t) {
    return this._data[t * B2 + 2];
  }
  hasContent(t) {
    return this._data[t * B2 + 0] & 4194303;
  }
  getCodePoint(t) {
    let e = this._data[t * B2 + 0];
    return e & 2097152 ? this._combined[t].charCodeAt(this._combined[t].length - 1) : e & 2097151;
  }
  isCombined(t) {
    return this._data[t * B2 + 0] & 2097152;
  }
  getString(t) {
    let e = this._data[t * B2 + 0];
    return e & 2097152 ? this._combined[t] : e & 2097151 ? Ce2(e & 2097151) : "";
  }
  isProtected(t) {
    return this._data[t * B2 + 2] & 536870912;
  }
  loadCell(t, e) {
    return an2 = t * B2, e.content = this._data[an2 + 0], e.fg = this._data[an2 + 1], e.bg = this._data[an2 + 2], e.content & 2097152 && (e.combinedData = this._combined[t]), e.bg & 268435456 && (e.extended = this._extendedAttrs[t]), e;
  }
  setCell(t, e) {
    e.content & 2097152 && (this._combined[t] = e.combinedData), e.bg & 268435456 && (this._extendedAttrs[t] = e.extended), this._data[t * B2 + 0] = e.content, this._data[t * B2 + 1] = e.fg, this._data[t * B2 + 2] = e.bg;
  }
  setCellFromCodepoint(t, e, i8, r5) {
    r5.bg & 268435456 && (this._extendedAttrs[t] = r5.extended), this._data[t * B2 + 0] = e | i8 << 22, this._data[t * B2 + 1] = r5.fg, this._data[t * B2 + 2] = r5.bg;
  }
  addCodepointToCell(t, e, i8) {
    let r5 = this._data[t * B2 + 0];
    r5 & 2097152 ? this._combined[t] += Ce2(e) : r5 & 2097151 ? (this._combined[t] = Ce2(r5 & 2097151) + Ce2(e), r5 &= -2097152, r5 |= 2097152) : r5 = e | 1 << 22, i8 && (r5 &= -12582913, r5 |= i8 << 22), this._data[t * B2 + 0] = r5;
  }
  insertCells(t, e, i8) {
    if (t %= this.length, t && this.getWidth(t - 1) === 2 && this.setCellFromCodepoint(t - 1, 0, 1, i8), e < this.length - t) {
      let r5 = new q();
      for (let n = this.length - t - e - 1; n >= 0; --n) this.setCell(t + e + n, this.loadCell(t + n, r5));
      for (let n = 0; n < e; ++n) this.setCell(t + n, i8);
    } else for (let r5 = t; r5 < this.length; ++r5) this.setCell(r5, i8);
    this.getWidth(this.length - 1) === 2 && this.setCellFromCodepoint(this.length - 1, 0, 1, i8);
  }
  deleteCells(t, e, i8) {
    if (t %= this.length, e < this.length - t) {
      let r5 = new q();
      for (let n = 0; n < this.length - t - e; ++n) this.setCell(t + n, this.loadCell(t + e + n, r5));
      for (let n = this.length - e; n < this.length; ++n) this.setCell(n, i8);
    } else for (let r5 = t; r5 < this.length; ++r5) this.setCell(r5, i8);
    t && this.getWidth(t - 1) === 2 && this.setCellFromCodepoint(t - 1, 0, 1, i8), this.getWidth(t) === 0 && !this.hasContent(t) && this.setCellFromCodepoint(t, 0, 1, i8);
  }
  replaceCells(t, e, i8, r5 = false) {
    if (r5) {
      for (t && this.getWidth(t - 1) === 2 && !this.isProtected(t - 1) && this.setCellFromCodepoint(t - 1, 0, 1, i8), e < this.length && this.getWidth(e - 1) === 2 && !this.isProtected(e) && this.setCellFromCodepoint(e, 0, 1, i8); t < e && t < this.length; ) this.isProtected(t) || this.setCell(t, i8), t++;
      return;
    }
    for (t && this.getWidth(t - 1) === 2 && this.setCellFromCodepoint(t - 1, 0, 1, i8), e < this.length && this.getWidth(e - 1) === 2 && this.setCellFromCodepoint(e, 0, 1, i8); t < e && t < this.length; ) this.setCell(t++, i8);
  }
  resize(t, e) {
    if (t === this.length) return this._data.length * 4 * Ls < this._data.buffer.byteLength;
    let i8 = t * B2;
    if (t > this.length) {
      if (this._data.buffer.byteLength >= i8 * 4) this._data = new Uint32Array(this._data.buffer, 0, i8);
      else {
        let r5 = new Uint32Array(i8);
        r5.set(this._data), this._data = r5;
      }
      for (let r5 = this.length; r5 < t; ++r5) this.setCell(r5, e);
    } else {
      this._data = this._data.subarray(0, i8);
      let r5 = Object.keys(this._combined);
      for (let o2 = 0; o2 < r5.length; o2++) {
        let l = parseInt(r5[o2], 10);
        l >= t && delete this._combined[l];
      }
      let n = Object.keys(this._extendedAttrs);
      for (let o2 = 0; o2 < n.length; o2++) {
        let l = parseInt(n[o2], 10);
        l >= t && delete this._extendedAttrs[l];
      }
    }
    return this.length = t, i8 * 4 * Ls < this._data.buffer.byteLength;
  }
  cleanupMemory() {
    if (this._data.length * 4 * Ls < this._data.buffer.byteLength) {
      let t = new Uint32Array(this._data.length);
      return t.set(this._data), this._data = t, 1;
    }
    return 0;
  }
  fill(t, e = false) {
    if (e) {
      for (let i8 = 0; i8 < this.length; ++i8) this.isProtected(i8) || this.setCell(i8, t);
      return;
    }
    this._combined = {}, this._extendedAttrs = {};
    for (let i8 = 0; i8 < this.length; ++i8) this.setCell(i8, t);
  }
  copyFrom(t) {
    this.length !== t.length ? this._data = new Uint32Array(t._data) : this._data.set(t._data), this.length = t.length, this._combined = {};
    for (let e in t._combined) this._combined[e] = t._combined[e];
    this._extendedAttrs = {};
    for (let e in t._extendedAttrs) this._extendedAttrs[e] = t._extendedAttrs[e];
    this.isWrapped = t.isWrapped;
  }
  clone() {
    let t = new s12(0);
    t._data = new Uint32Array(this._data), t.length = this.length;
    for (let e in this._combined) t._combined[e] = this._combined[e];
    for (let e in this._extendedAttrs) t._extendedAttrs[e] = this._extendedAttrs[e];
    return t.isWrapped = this.isWrapped, t;
  }
  getTrimmedLength() {
    for (let t = this.length - 1; t >= 0; --t) if (this._data[t * B2 + 0] & 4194303) return t + (this._data[t * B2 + 0] >> 22);
    return 0;
  }
  getNoBgTrimmedLength() {
    for (let t = this.length - 1; t >= 0; --t) if (this._data[t * B2 + 0] & 4194303 || this._data[t * B2 + 2] & 50331648) return t + (this._data[t * B2 + 0] >> 22);
    return 0;
  }
  copyCellsFrom(t, e, i8, r5, n) {
    let o2 = t._data;
    if (n) for (let a = r5 - 1; a >= 0; a--) {
      for (let u = 0; u < B2; u++) this._data[(i8 + a) * B2 + u] = o2[(e + a) * B2 + u];
      o2[(e + a) * B2 + 2] & 268435456 && (this._extendedAttrs[i8 + a] = t._extendedAttrs[e + a]);
    }
    else for (let a = 0; a < r5; a++) {
      for (let u = 0; u < B2; u++) this._data[(i8 + a) * B2 + u] = o2[(e + a) * B2 + u];
      o2[(e + a) * B2 + 2] & 268435456 && (this._extendedAttrs[i8 + a] = t._extendedAttrs[e + a]);
    }
    let l = Object.keys(t._combined);
    for (let a = 0; a < l.length; a++) {
      let u = parseInt(l[a], 10);
      u >= e && (this._combined[u - e + i8] = t._combined[u]);
    }
  }
  translateToString(t, e, i8, r5) {
    e = e ?? 0, i8 = i8 ?? this.length, t && (i8 = Math.min(i8, this.getTrimmedLength())), r5 && (r5.length = 0);
    let n = "";
    for (; e < i8; ) {
      let o2 = this._data[e * B2 + 0], l = o2 & 2097151, a = o2 & 2097152 ? this._combined[e] : l ? Ce2(l) : we2;
      if (n += a, r5) for (let u = 0; u < a.length; ++u) r5.push(e);
      e += o2 >> 22 || 1;
    }
    return r5 && r5.push(e), n;
  }
};
function sl(s15, t, e, i8, r5, n) {
  let o2 = [];
  for (let l = 0; l < s15.length - 1; l++) {
    let a = l, u = s15.get(++a);
    if (!u.isWrapped) continue;
    let h2 = [s15.get(l)];
    for (; a < s15.length && u.isWrapped; ) h2.push(u), u = s15.get(++a);
    if (!n && i8 >= l && i8 < a) {
      l += h2.length - 1;
      continue;
    }
    let c = 0, d = ri2(h2, c, t), _3 = 1, p = 0;
    for (; _3 < h2.length; ) {
      let f = ri2(h2, _3, t), A2 = f - p, R3 = e - d, O3 = Math.min(A2, R3);
      h2[c].copyCellsFrom(h2[_3], p, d, O3, false), d += O3, d === e && (c++, d = 0), p += O3, p === f && (_3++, p = 0), d === 0 && c !== 0 && h2[c - 1].getWidth(e - 1) === 2 && (h2[c].copyCellsFrom(h2[c - 1], e - 1, d++, 1, false), h2[c - 1].setCell(e - 1, r5));
    }
    h2[c].replaceCells(d, e, r5);
    let m = 0;
    for (let f = h2.length - 1; f > 0 && (f > c || h2[f].getTrimmedLength() === 0); f--) m++;
    m > 0 && (o2.push(l + h2.length - m), o2.push(m)), l += h2.length - 1;
  }
  return o2;
}
function ol(s15, t) {
  let e = [], i8 = 0, r5 = t[i8], n = 0;
  for (let o2 = 0; o2 < s15.length; o2++) if (r5 === o2) {
    let l = t[++i8];
    s15.onDeleteEmitter.fire({ index: o2 - n, amount: l }), o2 += l - 1, n += l, r5 = t[++i8];
  } else e.push(o2);
  return { layout: e, countRemoved: n };
}
function ll(s15, t) {
  let e = [];
  for (let i8 = 0; i8 < t.length; i8++) e.push(s15.get(t[i8]));
  for (let i8 = 0; i8 < e.length; i8++) s15.set(i8, e[i8]);
  s15.length = t.length;
}
function al(s15, t, e) {
  let i8 = [], r5 = s15.map((a, u) => ri2(s15, u, t)).reduce((a, u) => a + u), n = 0, o2 = 0, l = 0;
  for (; l < r5; ) {
    if (r5 - l < e) {
      i8.push(r5 - l);
      break;
    }
    n += e;
    let a = ri2(s15, o2, t);
    n > a && (n -= a, o2++);
    let u = s15[o2].getWidth(n - 1) === 2;
    u && n--;
    let h2 = u ? e - 1 : e;
    i8.push(h2), l += h2;
  }
  return i8;
}
function ri2(s15, t, e) {
  if (t === s15.length - 1) return s15[t].getTrimmedLength();
  let i8 = !s15[t].hasContent(e - 1) && s15[t].getWidth(e - 1) === 1, r5 = s15[t + 1].getWidth(0) === 2;
  return i8 && r5 ? e - 1 : e;
}
var un2 = class un3 {
  constructor(t) {
    this.line = t;
    this.isDisposed = false;
    this._disposables = [];
    this._id = un3._nextId++;
    this._onDispose = this.register(new v());
    this.onDispose = this._onDispose.event;
  }
  get id() {
    return this._id;
  }
  dispose() {
    this.isDisposed || (this.isDisposed = true, this.line = -1, this._onDispose.fire(), Ne2(this._disposables), this._disposables.length = 0);
  }
  register(t) {
    return this._disposables.push(t), t;
  }
};
un2._nextId = 1;
var cn2 = un2;
var ne2 = {};
var Je2 = ne2.B;
ne2[0] = { "`": "◆", a: "▒", b: "␉", c: "␌", d: "␍", e: "␊", f: "°", g: "±", h: "␤", i: "␋", j: "┘", k: "┐", l: "┌", m: "└", n: "┼", o: "⎺", p: "⎻", q: "─", r: "⎼", s: "⎽", t: "├", u: "┤", v: "┴", w: "┬", x: "│", y: "≤", z: "≥", "{": "π", "|": "≠", "}": "£", "~": "·" };
ne2.A = { "#": "£" };
ne2.B = void 0;
ne2[4] = { "#": "£", "@": "¾", "[": "ij", "\\": "½", "]": "|", "{": "¨", "|": "f", "}": "¼", "~": "´" };
ne2.C = ne2[5] = { "[": "Ä", "\\": "Ö", "]": "Å", "^": "Ü", "`": "é", "{": "ä", "|": "ö", "}": "å", "~": "ü" };
ne2.R = { "#": "£", "@": "à", "[": "°", "\\": "ç", "]": "§", "{": "é", "|": "ù", "}": "è", "~": "¨" };
ne2.Q = { "@": "à", "[": "â", "\\": "ç", "]": "ê", "^": "î", "`": "ô", "{": "é", "|": "ù", "}": "è", "~": "û" };
ne2.K = { "@": "§", "[": "Ä", "\\": "Ö", "]": "Ü", "{": "ä", "|": "ö", "}": "ü", "~": "ß" };
ne2.Y = { "#": "£", "@": "§", "[": "°", "\\": "ç", "]": "é", "`": "ù", "{": "à", "|": "ò", "}": "è", "~": "ì" };
ne2.E = ne2[6] = { "@": "Ä", "[": "Æ", "\\": "Ø", "]": "Å", "^": "Ü", "`": "ä", "{": "æ", "|": "ø", "}": "å", "~": "ü" };
ne2.Z = { "#": "£", "@": "§", "[": "¡", "\\": "Ñ", "]": "¿", "{": "°", "|": "ñ", "}": "ç" };
ne2.H = ne2[7] = { "@": "É", "[": "Ä", "\\": "Ö", "]": "Å", "^": "Ü", "`": "é", "{": "ä", "|": "ö", "}": "å", "~": "ü" };
ne2["="] = { "#": "ù", "@": "à", "[": "é", "\\": "ç", "]": "ê", "^": "î", _: "è", "`": "ô", "{": "ä", "|": "ö", "}": "ü", "~": "û" };
var cl = 4294967295;
var $i2 = class {
  constructor(t, e, i8) {
    this._hasScrollback = t;
    this._optionsService = e;
    this._bufferService = i8;
    this.ydisp = 0;
    this.ybase = 0;
    this.y = 0;
    this.x = 0;
    this.tabs = {};
    this.savedY = 0;
    this.savedX = 0;
    this.savedCurAttrData = X4.clone();
    this.savedCharset = Je2;
    this.markers = [];
    this._nullCell = q.fromCharData([0, ir, 1, 0]);
    this._whitespaceCell = q.fromCharData([0, we2, 1, 32]);
    this._isClearing = false;
    this._memoryCleanupQueue = new Jt();
    this._memoryCleanupPosition = 0;
    this._cols = this._bufferService.cols, this._rows = this._bufferService.rows, this.lines = new zi2(this._getCorrectBufferLength(this._rows)), this.scrollTop = 0, this.scrollBottom = this._rows - 1, this.setupTabStops();
  }
  getNullCell(t) {
    return t ? (this._nullCell.fg = t.fg, this._nullCell.bg = t.bg, this._nullCell.extended = t.extended) : (this._nullCell.fg = 0, this._nullCell.bg = 0, this._nullCell.extended = new rt3()), this._nullCell;
  }
  getWhitespaceCell(t) {
    return t ? (this._whitespaceCell.fg = t.fg, this._whitespaceCell.bg = t.bg, this._whitespaceCell.extended = t.extended) : (this._whitespaceCell.fg = 0, this._whitespaceCell.bg = 0, this._whitespaceCell.extended = new rt3()), this._whitespaceCell;
  }
  getBlankLine(t, e) {
    return new Ze3(this._bufferService.cols, this.getNullCell(t), e);
  }
  get hasScrollback() {
    return this._hasScrollback && this.lines.maxLength > this._rows;
  }
  get isCursorInViewport() {
    let e = this.ybase + this.y - this.ydisp;
    return e >= 0 && e < this._rows;
  }
  _getCorrectBufferLength(t) {
    if (!this._hasScrollback) return t;
    let e = t + this._optionsService.rawOptions.scrollback;
    return e > cl ? cl : e;
  }
  fillViewportRows(t) {
    if (this.lines.length === 0) {
      t === void 0 && (t = X4);
      let e = this._rows;
      for (; e--; ) this.lines.push(this.getBlankLine(t));
    }
  }
  clear() {
    this.ydisp = 0, this.ybase = 0, this.y = 0, this.x = 0, this.lines = new zi2(this._getCorrectBufferLength(this._rows)), this.scrollTop = 0, this.scrollBottom = this._rows - 1, this.setupTabStops();
  }
  resize(t, e) {
    let i8 = this.getNullCell(X4), r5 = 0, n = this._getCorrectBufferLength(e);
    if (n > this.lines.maxLength && (this.lines.maxLength = n), this.lines.length > 0) {
      if (this._cols < t) for (let l = 0; l < this.lines.length; l++) r5 += +this.lines.get(l).resize(t, i8);
      let o2 = 0;
      if (this._rows < e) for (let l = this._rows; l < e; l++) this.lines.length < e + this.ybase && (this._optionsService.rawOptions.windowsMode || this._optionsService.rawOptions.windowsPty.backend !== void 0 || this._optionsService.rawOptions.windowsPty.buildNumber !== void 0 ? this.lines.push(new Ze3(t, i8)) : this.ybase > 0 && this.lines.length <= this.ybase + this.y + o2 + 1 ? (this.ybase--, o2++, this.ydisp > 0 && this.ydisp--) : this.lines.push(new Ze3(t, i8)));
      else for (let l = this._rows; l > e; l--) this.lines.length > e + this.ybase && (this.lines.length > this.ybase + this.y + 1 ? this.lines.pop() : (this.ybase++, this.ydisp++));
      if (n < this.lines.maxLength) {
        let l = this.lines.length - n;
        l > 0 && (this.lines.trimStart(l), this.ybase = Math.max(this.ybase - l, 0), this.ydisp = Math.max(this.ydisp - l, 0), this.savedY = Math.max(this.savedY - l, 0)), this.lines.maxLength = n;
      }
      this.x = Math.min(this.x, t - 1), this.y = Math.min(this.y, e - 1), o2 && (this.y += o2), this.savedX = Math.min(this.savedX, t - 1), this.scrollTop = 0;
    }
    if (this.scrollBottom = e - 1, this._isReflowEnabled && (this._reflow(t, e), this._cols > t)) for (let o2 = 0; o2 < this.lines.length; o2++) r5 += +this.lines.get(o2).resize(t, i8);
    this._cols = t, this._rows = e, this._memoryCleanupQueue.clear(), r5 > 0.1 * this.lines.length && (this._memoryCleanupPosition = 0, this._memoryCleanupQueue.enqueue(() => this._batchedMemoryCleanup()));
  }
  _batchedMemoryCleanup() {
    let t = true;
    this._memoryCleanupPosition >= this.lines.length && (this._memoryCleanupPosition = 0, t = false);
    let e = 0;
    for (; this._memoryCleanupPosition < this.lines.length; ) if (e += this.lines.get(this._memoryCleanupPosition++).cleanupMemory(), e > 100) return true;
    return t;
  }
  get _isReflowEnabled() {
    let t = this._optionsService.rawOptions.windowsPty;
    return t && t.buildNumber ? this._hasScrollback && t.backend === "conpty" && t.buildNumber >= 21376 : this._hasScrollback && !this._optionsService.rawOptions.windowsMode;
  }
  _reflow(t, e) {
    this._cols !== t && (t > this._cols ? this._reflowLarger(t, e) : this._reflowSmaller(t, e));
  }
  _reflowLarger(t, e) {
    let i8 = this._optionsService.rawOptions.reflowCursorLine, r5 = sl(this.lines, this._cols, t, this.ybase + this.y, this.getNullCell(X4), i8);
    if (r5.length > 0) {
      let n = ol(this.lines, r5);
      ll(this.lines, n.layout), this._reflowLargerAdjustViewport(t, e, n.countRemoved);
    }
  }
  _reflowLargerAdjustViewport(t, e, i8) {
    let r5 = this.getNullCell(X4), n = i8;
    for (; n-- > 0; ) this.ybase === 0 ? (this.y > 0 && this.y--, this.lines.length < e && this.lines.push(new Ze3(t, r5))) : (this.ydisp === this.ybase && this.ydisp--, this.ybase--);
    this.savedY = Math.max(this.savedY - i8, 0);
  }
  _reflowSmaller(t, e) {
    let i8 = this._optionsService.rawOptions.reflowCursorLine, r5 = this.getNullCell(X4), n = [], o2 = 0;
    for (let l = this.lines.length - 1; l >= 0; l--) {
      let a = this.lines.get(l);
      if (!a || !a.isWrapped && a.getTrimmedLength() <= t) continue;
      let u = [a];
      for (; a.isWrapped && l > 0; ) a = this.lines.get(--l), u.unshift(a);
      if (!i8) {
        let I2 = this.ybase + this.y;
        if (I2 >= l && I2 < l + u.length) continue;
      }
      let h2 = u[u.length - 1].getTrimmedLength(), c = al(u, this._cols, t), d = c.length - u.length, _3;
      this.ybase === 0 && this.y !== this.lines.length - 1 ? _3 = Math.max(0, this.y - this.lines.maxLength + d) : _3 = Math.max(0, this.lines.length - this.lines.maxLength + d);
      let p = [];
      for (let I2 = 0; I2 < d; I2++) {
        let k2 = this.getBlankLine(X4, true);
        p.push(k2);
      }
      p.length > 0 && (n.push({ start: l + u.length + o2, newLines: p }), o2 += p.length), u.push(...p);
      let m = c.length - 1, f = c[m];
      f === 0 && (m--, f = c[m]);
      let A2 = u.length - d - 1, R3 = h2;
      for (; A2 >= 0; ) {
        let I2 = Math.min(R3, f);
        if (u[m] === void 0) break;
        if (u[m].copyCellsFrom(u[A2], R3 - I2, f - I2, I2, true), f -= I2, f === 0 && (m--, f = c[m]), R3 -= I2, R3 === 0) {
          A2--;
          let k2 = Math.max(A2, 0);
          R3 = ri2(u, k2, this._cols);
        }
      }
      for (let I2 = 0; I2 < u.length; I2++) c[I2] < t && u[I2].setCell(c[I2], r5);
      let O3 = d - _3;
      for (; O3-- > 0; ) this.ybase === 0 ? this.y < e - 1 ? (this.y++, this.lines.pop()) : (this.ybase++, this.ydisp++) : this.ybase < Math.min(this.lines.maxLength, this.lines.length + o2) - e && (this.ybase === this.ydisp && this.ydisp++, this.ybase++);
      this.savedY = Math.min(this.savedY + d, this.ybase + e - 1);
    }
    if (n.length > 0) {
      let l = [], a = [];
      for (let f = 0; f < this.lines.length; f++) a.push(this.lines.get(f));
      let u = this.lines.length, h2 = u - 1, c = 0, d = n[c];
      this.lines.length = Math.min(this.lines.maxLength, this.lines.length + o2);
      let _3 = 0;
      for (let f = Math.min(this.lines.maxLength - 1, u + o2 - 1); f >= 0; f--) if (d && d.start > h2 + _3) {
        for (let A2 = d.newLines.length - 1; A2 >= 0; A2--) this.lines.set(f--, d.newLines[A2]);
        f++, l.push({ index: h2 + 1, amount: d.newLines.length }), _3 += d.newLines.length, d = n[++c];
      } else this.lines.set(f, a[h2--]);
      let p = 0;
      for (let f = l.length - 1; f >= 0; f--) l[f].index += p, this.lines.onInsertEmitter.fire(l[f]), p += l[f].amount;
      let m = Math.max(0, u + o2 - this.lines.maxLength);
      m > 0 && this.lines.onTrimEmitter.fire(m);
    }
  }
  translateBufferLineToString(t, e, i8 = 0, r5) {
    let n = this.lines.get(t);
    return n ? n.translateToString(e, i8, r5) : "";
  }
  getWrappedRangeForLine(t) {
    let e = t, i8 = t;
    for (; e > 0 && this.lines.get(e).isWrapped; ) e--;
    for (; i8 + 1 < this.lines.length && this.lines.get(i8 + 1).isWrapped; ) i8++;
    return { first: e, last: i8 };
  }
  setupTabStops(t) {
    for (t != null ? this.tabs[t] || (t = this.prevStop(t)) : (this.tabs = {}, t = 0); t < this._cols; t += this._optionsService.rawOptions.tabStopWidth) this.tabs[t] = true;
  }
  prevStop(t) {
    for (t == null && (t = this.x); !this.tabs[--t] && t > 0; ) ;
    return t >= this._cols ? this._cols - 1 : t < 0 ? 0 : t;
  }
  nextStop(t) {
    for (t == null && (t = this.x); !this.tabs[++t] && t < this._cols; ) ;
    return t >= this._cols ? this._cols - 1 : t < 0 ? 0 : t;
  }
  clearMarkers(t) {
    this._isClearing = true;
    for (let e = 0; e < this.markers.length; e++) this.markers[e].line === t && (this.markers[e].dispose(), this.markers.splice(e--, 1));
    this._isClearing = false;
  }
  clearAllMarkers() {
    this._isClearing = true;
    for (let t = 0; t < this.markers.length; t++) this.markers[t].dispose();
    this.markers.length = 0, this._isClearing = false;
  }
  addMarker(t) {
    let e = new cn2(t);
    return this.markers.push(e), e.register(this.lines.onTrim((i8) => {
      e.line -= i8, e.line < 0 && e.dispose();
    })), e.register(this.lines.onInsert((i8) => {
      e.line >= i8.index && (e.line += i8.amount);
    })), e.register(this.lines.onDelete((i8) => {
      e.line >= i8.index && e.line < i8.index + i8.amount && e.dispose(), e.line > i8.index && (e.line -= i8.amount);
    })), e.register(e.onDispose(() => this._removeMarker(e))), e;
  }
  _removeMarker(t) {
    this._isClearing || this.markers.splice(this.markers.indexOf(t), 1);
  }
};
var hn2 = class extends D2 {
  constructor(e, i8) {
    super();
    this._optionsService = e;
    this._bufferService = i8;
    this._onBufferActivate = this._register(new v());
    this.onBufferActivate = this._onBufferActivate.event;
    this.reset(), this._register(this._optionsService.onSpecificOptionChange("scrollback", () => this.resize(this._bufferService.cols, this._bufferService.rows))), this._register(this._optionsService.onSpecificOptionChange("tabStopWidth", () => this.setupTabStops()));
  }
  reset() {
    this._normal = new $i2(true, this._optionsService, this._bufferService), this._normal.fillViewportRows(), this._alt = new $i2(false, this._optionsService, this._bufferService), this._activeBuffer = this._normal, this._onBufferActivate.fire({ activeBuffer: this._normal, inactiveBuffer: this._alt }), this.setupTabStops();
  }
  get alt() {
    return this._alt;
  }
  get active() {
    return this._activeBuffer;
  }
  get normal() {
    return this._normal;
  }
  activateNormalBuffer() {
    this._activeBuffer !== this._normal && (this._normal.x = this._alt.x, this._normal.y = this._alt.y, this._alt.clearAllMarkers(), this._alt.clear(), this._activeBuffer = this._normal, this._onBufferActivate.fire({ activeBuffer: this._normal, inactiveBuffer: this._alt }));
  }
  activateAltBuffer(e) {
    this._activeBuffer !== this._alt && (this._alt.fillViewportRows(e), this._alt.x = this._normal.x, this._alt.y = this._normal.y, this._activeBuffer = this._alt, this._onBufferActivate.fire({ activeBuffer: this._alt, inactiveBuffer: this._normal }));
  }
  resize(e, i8) {
    this._normal.resize(e, i8), this._alt.resize(e, i8), this.setupTabStops(e);
  }
  setupTabStops(e) {
    this._normal.setupTabStops(e), this._alt.setupTabStops(e);
  }
};
var ks2 = 2;
var Cs = 1;
var ni2 = class extends D2 {
  constructor(e) {
    super();
    this.isUserScrolling = false;
    this._onResize = this._register(new v());
    this.onResize = this._onResize.event;
    this._onScroll = this._register(new v());
    this.onScroll = this._onScroll.event;
    this.cols = Math.max(e.rawOptions.cols || 0, ks2), this.rows = Math.max(e.rawOptions.rows || 0, Cs), this.buffers = this._register(new hn2(e, this)), this._register(this.buffers.onBufferActivate((i8) => {
      this._onScroll.fire(i8.activeBuffer.ydisp);
    }));
  }
  get buffer() {
    return this.buffers.active;
  }
  resize(e, i8) {
    let r5 = this.cols !== e, n = this.rows !== i8;
    this.cols = e, this.rows = i8, this.buffers.resize(e, i8), this._onResize.fire({ cols: e, rows: i8, colsChanged: r5, rowsChanged: n });
  }
  reset() {
    this.buffers.reset(), this.isUserScrolling = false;
  }
  scroll(e, i8 = false) {
    let r5 = this.buffer, n;
    n = this._cachedBlankLine, (!n || n.length !== this.cols || n.getFg(0) !== e.fg || n.getBg(0) !== e.bg) && (n = r5.getBlankLine(e, i8), this._cachedBlankLine = n), n.isWrapped = i8;
    let o2 = r5.ybase + r5.scrollTop, l = r5.ybase + r5.scrollBottom;
    if (r5.scrollTop === 0) {
      let a = r5.lines.isFull;
      l === r5.lines.length - 1 ? a ? r5.lines.recycle().copyFrom(n) : r5.lines.push(n.clone()) : r5.lines.splice(l + 1, 0, n.clone()), a ? this.isUserScrolling && (r5.ydisp = Math.max(r5.ydisp - 1, 0)) : (r5.ybase++, this.isUserScrolling || r5.ydisp++);
    } else {
      let a = l - o2 + 1;
      r5.lines.shiftElements(o2 + 1, a - 1, -1), r5.lines.set(l, n.clone());
    }
    this.isUserScrolling || (r5.ydisp = r5.ybase), this._onScroll.fire(r5.ydisp);
  }
  scrollLines(e, i8) {
    let r5 = this.buffer;
    if (e < 0) {
      if (r5.ydisp === 0) return;
      this.isUserScrolling = true;
    } else e + r5.ydisp >= r5.ybase && (this.isUserScrolling = false);
    let n = r5.ydisp;
    r5.ydisp = Math.max(Math.min(r5.ydisp + e, r5.ybase), 0), n !== r5.ydisp && (i8 || this._onScroll.fire(r5.ydisp));
  }
};
ni2 = M3([S(0, H3)], ni2);
var si = { cols: 80, rows: 24, cursorBlink: false, cursorStyle: "block", cursorWidth: 1, cursorInactiveStyle: "outline", customGlyphs: true, drawBoldTextInBrightColors: true, documentOverride: null, fastScrollModifier: "alt", fastScrollSensitivity: 5, fontFamily: "monospace", fontSize: 15, fontWeight: "normal", fontWeightBold: "bold", ignoreBracketedPasteMode: false, lineHeight: 1, letterSpacing: 0, linkHandler: null, logLevel: "info", logger: null, scrollback: 1e3, scrollOnEraseInDisplay: false, scrollOnUserInput: true, scrollSensitivity: 1, screenReaderMode: false, smoothScrollDuration: 0, macOptionIsMeta: false, macOptionClickForcesSelection: false, minimumContrastRatio: 1, disableStdin: false, allowProposedApi: false, allowTransparency: false, tabStopWidth: 8, theme: {}, reflowCursorLine: false, rescaleOverlappingGlyphs: false, rightClickSelectsWord: Zt3, windowOptions: {}, windowsMode: false, windowsPty: {}, wordSeparator: " ()[]{}',\"`", altClickMovesCursor: true, convertEol: false, termName: "xterm", cancelEvents: false, overviewRuler: {} };
var nc = ["normal", "bold", "100", "200", "300", "400", "500", "600", "700", "800", "900"];
var dn2 = class extends D2 {
  constructor(e) {
    super();
    this._onOptionChange = this._register(new v());
    this.onOptionChange = this._onOptionChange.event;
    let i8 = { ...si };
    for (let r5 in e) if (r5 in i8) try {
      let n = e[r5];
      i8[r5] = this._sanitizeAndValidateOption(r5, n);
    } catch (n) {
      console.error(n);
    }
    this.rawOptions = i8, this.options = { ...i8 }, this._setupOptions(), this._register(C3(() => {
      this.rawOptions.linkHandler = null, this.rawOptions.documentOverride = null;
    }));
  }
  onSpecificOptionChange(e, i8) {
    return this.onOptionChange((r5) => {
      r5 === e && i8(this.rawOptions[e]);
    });
  }
  onMultipleOptionChange(e, i8) {
    return this.onOptionChange((r5) => {
      e.indexOf(r5) !== -1 && i8();
    });
  }
  _setupOptions() {
    let e = (r5) => {
      if (!(r5 in si)) throw new Error(`No option with key "${r5}"`);
      return this.rawOptions[r5];
    }, i8 = (r5, n) => {
      if (!(r5 in si)) throw new Error(`No option with key "${r5}"`);
      n = this._sanitizeAndValidateOption(r5, n), this.rawOptions[r5] !== n && (this.rawOptions[r5] = n, this._onOptionChange.fire(r5));
    };
    for (let r5 in this.rawOptions) {
      let n = { get: e.bind(this, r5), set: i8.bind(this, r5) };
      Object.defineProperty(this.options, r5, n);
    }
  }
  _sanitizeAndValidateOption(e, i8) {
    switch (e) {
      case "cursorStyle":
        if (i8 || (i8 = si[e]), !sc(i8)) throw new Error(`"${i8}" is not a valid value for ${e}`);
        break;
      case "wordSeparator":
        i8 || (i8 = si[e]);
        break;
      case "fontWeight":
      case "fontWeightBold":
        if (typeof i8 == "number" && 1 <= i8 && i8 <= 1e3) break;
        i8 = nc.includes(i8) ? i8 : si[e];
        break;
      case "cursorWidth":
        i8 = Math.floor(i8);
      case "lineHeight":
      case "tabStopWidth":
        if (i8 < 1) throw new Error(`${e} cannot be less than 1, value: ${i8}`);
        break;
      case "minimumContrastRatio":
        i8 = Math.max(1, Math.min(21, Math.round(i8 * 10) / 10));
        break;
      case "scrollback":
        if (i8 = Math.min(i8, 4294967295), i8 < 0) throw new Error(`${e} cannot be less than 0, value: ${i8}`);
        break;
      case "fastScrollSensitivity":
      case "scrollSensitivity":
        if (i8 <= 0) throw new Error(`${e} cannot be less than or equal to 0, value: ${i8}`);
        break;
      case "rows":
      case "cols":
        if (!i8 && i8 !== 0) throw new Error(`${e} must be numeric, value: ${i8}`);
        break;
      case "windowsPty":
        i8 = i8 ?? {};
        break;
    }
    return i8;
  }
};
function sc(s15) {
  return s15 === "block" || s15 === "underline" || s15 === "bar";
}
function oi2(s15, t = 5) {
  if (typeof s15 != "object") return s15;
  let e = Array.isArray(s15) ? [] : {};
  for (let i8 in s15) e[i8] = t <= 1 ? s15[i8] : s15[i8] && oi2(s15[i8], t - 1);
  return e;
}
var ul = Object.freeze({ insertMode: false });
var hl = Object.freeze({ applicationCursorKeys: false, applicationKeypad: false, bracketedPasteMode: false, cursorBlink: void 0, cursorStyle: void 0, origin: false, reverseWraparound: false, sendFocus: false, synchronizedOutput: false, wraparound: true });
var li2 = class extends D2 {
  constructor(e, i8, r5) {
    super();
    this._bufferService = e;
    this._logService = i8;
    this._optionsService = r5;
    this.isCursorInitialized = false;
    this.isCursorHidden = false;
    this._onData = this._register(new v());
    this.onData = this._onData.event;
    this._onUserInput = this._register(new v());
    this.onUserInput = this._onUserInput.event;
    this._onBinary = this._register(new v());
    this.onBinary = this._onBinary.event;
    this._onRequestScrollToBottom = this._register(new v());
    this.onRequestScrollToBottom = this._onRequestScrollToBottom.event;
    this.modes = oi2(ul), this.decPrivateModes = oi2(hl);
  }
  reset() {
    this.modes = oi2(ul), this.decPrivateModes = oi2(hl);
  }
  triggerDataEvent(e, i8 = false) {
    if (this._optionsService.rawOptions.disableStdin) return;
    let r5 = this._bufferService.buffer;
    i8 && this._optionsService.rawOptions.scrollOnUserInput && r5.ybase !== r5.ydisp && this._onRequestScrollToBottom.fire(), i8 && this._onUserInput.fire(), this._logService.debug(`sending data "${e}"`), this._logService.trace("sending data (codes)", () => e.split("").map((n) => n.charCodeAt(0))), this._onData.fire(e);
  }
  triggerBinaryEvent(e) {
    this._optionsService.rawOptions.disableStdin || (this._logService.debug(`sending binary "${e}"`), this._logService.trace("sending binary (codes)", () => e.split("").map((i8) => i8.charCodeAt(0))), this._onBinary.fire(e));
  }
};
li2 = M3([S(0, F2), S(1, nr2), S(2, H3)], li2);
var dl = { NONE: { events: 0, restrict: () => false }, X10: { events: 1, restrict: (s15) => s15.button === 4 || s15.action !== 1 ? false : (s15.ctrl = false, s15.alt = false, s15.shift = false, true) }, VT200: { events: 19, restrict: (s15) => s15.action !== 32 }, DRAG: { events: 23, restrict: (s15) => !(s15.action === 32 && s15.button === 3) }, ANY: { events: 31, restrict: (s15) => true } };
function Ms2(s15, t) {
  let e = (s15.ctrl ? 16 : 0) | (s15.shift ? 4 : 0) | (s15.alt ? 8 : 0);
  return s15.button === 4 ? (e |= 64, e |= s15.action) : (e |= s15.button & 3, s15.button & 4 && (e |= 64), s15.button & 8 && (e |= 128), s15.action === 32 ? e |= 32 : s15.action === 0 && !t && (e |= 3)), e;
}
var Ps2 = String.fromCharCode;
var fl = { DEFAULT: (s15) => {
  let t = [Ms2(s15, false) + 32, s15.col + 32, s15.row + 32];
  return t[0] > 255 || t[1] > 255 || t[2] > 255 ? "" : `\x1B[M${Ps2(t[0])}${Ps2(t[1])}${Ps2(t[2])}`;
}, SGR: (s15) => {
  let t = s15.action === 0 && s15.button !== 4 ? "m" : "M";
  return `\x1B[<${Ms2(s15, true)};${s15.col};${s15.row}${t}`;
}, SGR_PIXELS: (s15) => {
  let t = s15.action === 0 && s15.button !== 4 ? "m" : "M";
  return `\x1B[<${Ms2(s15, true)};${s15.x};${s15.y}${t}`;
} };
var ai3 = class extends D2 {
  constructor(e, i8, r5) {
    super();
    this._bufferService = e;
    this._coreService = i8;
    this._optionsService = r5;
    this._protocols = {};
    this._encodings = {};
    this._activeProtocol = "";
    this._activeEncoding = "";
    this._lastEvent = null;
    this._wheelPartialScroll = 0;
    this._onProtocolChange = this._register(new v());
    this.onProtocolChange = this._onProtocolChange.event;
    for (let n of Object.keys(dl)) this.addProtocol(n, dl[n]);
    for (let n of Object.keys(fl)) this.addEncoding(n, fl[n]);
    this.reset();
  }
  addProtocol(e, i8) {
    this._protocols[e] = i8;
  }
  addEncoding(e, i8) {
    this._encodings[e] = i8;
  }
  get activeProtocol() {
    return this._activeProtocol;
  }
  get areMouseEventsActive() {
    return this._protocols[this._activeProtocol].events !== 0;
  }
  set activeProtocol(e) {
    if (!this._protocols[e]) throw new Error(`unknown protocol "${e}"`);
    this._activeProtocol = e, this._onProtocolChange.fire(this._protocols[e].events);
  }
  get activeEncoding() {
    return this._activeEncoding;
  }
  set activeEncoding(e) {
    if (!this._encodings[e]) throw new Error(`unknown encoding "${e}"`);
    this._activeEncoding = e;
  }
  reset() {
    this.activeProtocol = "NONE", this.activeEncoding = "DEFAULT", this._lastEvent = null, this._wheelPartialScroll = 0;
  }
  consumeWheelEvent(e, i8, r5) {
    if (e.deltaY === 0 || e.shiftKey || i8 === void 0 || r5 === void 0) return 0;
    let n = i8 / r5, o2 = this._applyScrollModifier(e.deltaY, e);
    return e.deltaMode === WheelEvent.DOM_DELTA_PIXEL ? (o2 /= n + 0, Math.abs(e.deltaY) < 50 && (o2 *= 0.3), this._wheelPartialScroll += o2, o2 = Math.floor(Math.abs(this._wheelPartialScroll)) * (this._wheelPartialScroll > 0 ? 1 : -1), this._wheelPartialScroll %= 1) : e.deltaMode === WheelEvent.DOM_DELTA_PAGE && (o2 *= this._bufferService.rows), o2;
  }
  _applyScrollModifier(e, i8) {
    return i8.altKey || i8.ctrlKey || i8.shiftKey ? e * this._optionsService.rawOptions.fastScrollSensitivity * this._optionsService.rawOptions.scrollSensitivity : e * this._optionsService.rawOptions.scrollSensitivity;
  }
  triggerMouseEvent(e) {
    if (e.col < 0 || e.col >= this._bufferService.cols || e.row < 0 || e.row >= this._bufferService.rows || e.button === 4 && e.action === 32 || e.button === 3 && e.action !== 32 || e.button !== 4 && (e.action === 2 || e.action === 3) || (e.col++, e.row++, e.action === 32 && this._lastEvent && this._equalEvents(this._lastEvent, e, this._activeEncoding === "SGR_PIXELS")) || !this._protocols[this._activeProtocol].restrict(e)) return false;
    let i8 = this._encodings[this._activeEncoding](e);
    return i8 && (this._activeEncoding === "DEFAULT" ? this._coreService.triggerBinaryEvent(i8) : this._coreService.triggerDataEvent(i8, true)), this._lastEvent = e, true;
  }
  explainEvents(e) {
    return { down: !!(e & 1), up: !!(e & 2), drag: !!(e & 4), move: !!(e & 8), wheel: !!(e & 16) };
  }
  _equalEvents(e, i8, r5) {
    if (r5) {
      if (e.x !== i8.x || e.y !== i8.y) return false;
    } else if (e.col !== i8.col || e.row !== i8.row) return false;
    return !(e.button !== i8.button || e.action !== i8.action || e.ctrl !== i8.ctrl || e.alt !== i8.alt || e.shift !== i8.shift);
  }
};
ai3 = M3([S(0, F2), S(1, ge3), S(2, H3)], ai3);
var Os2 = [[768, 879], [1155, 1158], [1160, 1161], [1425, 1469], [1471, 1471], [1473, 1474], [1476, 1477], [1479, 1479], [1536, 1539], [1552, 1557], [1611, 1630], [1648, 1648], [1750, 1764], [1767, 1768], [1770, 1773], [1807, 1807], [1809, 1809], [1840, 1866], [1958, 1968], [2027, 2035], [2305, 2306], [2364, 2364], [2369, 2376], [2381, 2381], [2385, 2388], [2402, 2403], [2433, 2433], [2492, 2492], [2497, 2500], [2509, 2509], [2530, 2531], [2561, 2562], [2620, 2620], [2625, 2626], [2631, 2632], [2635, 2637], [2672, 2673], [2689, 2690], [2748, 2748], [2753, 2757], [2759, 2760], [2765, 2765], [2786, 2787], [2817, 2817], [2876, 2876], [2879, 2879], [2881, 2883], [2893, 2893], [2902, 2902], [2946, 2946], [3008, 3008], [3021, 3021], [3134, 3136], [3142, 3144], [3146, 3149], [3157, 3158], [3260, 3260], [3263, 3263], [3270, 3270], [3276, 3277], [3298, 3299], [3393, 3395], [3405, 3405], [3530, 3530], [3538, 3540], [3542, 3542], [3633, 3633], [3636, 3642], [3655, 3662], [3761, 3761], [3764, 3769], [3771, 3772], [3784, 3789], [3864, 3865], [3893, 3893], [3895, 3895], [3897, 3897], [3953, 3966], [3968, 3972], [3974, 3975], [3984, 3991], [3993, 4028], [4038, 4038], [4141, 4144], [4146, 4146], [4150, 4151], [4153, 4153], [4184, 4185], [4448, 4607], [4959, 4959], [5906, 5908], [5938, 5940], [5970, 5971], [6002, 6003], [6068, 6069], [6071, 6077], [6086, 6086], [6089, 6099], [6109, 6109], [6155, 6157], [6313, 6313], [6432, 6434], [6439, 6440], [6450, 6450], [6457, 6459], [6679, 6680], [6912, 6915], [6964, 6964], [6966, 6970], [6972, 6972], [6978, 6978], [7019, 7027], [7616, 7626], [7678, 7679], [8203, 8207], [8234, 8238], [8288, 8291], [8298, 8303], [8400, 8431], [12330, 12335], [12441, 12442], [43014, 43014], [43019, 43019], [43045, 43046], [64286, 64286], [65024, 65039], [65056, 65059], [65279, 65279], [65529, 65531]];
var ac = [[68097, 68099], [68101, 68102], [68108, 68111], [68152, 68154], [68159, 68159], [119143, 119145], [119155, 119170], [119173, 119179], [119210, 119213], [119362, 119364], [917505, 917505], [917536, 917631], [917760, 917999]];
var se2;
function cc(s15, t) {
  let e = 0, i8 = t.length - 1, r5;
  if (s15 < t[0][0] || s15 > t[i8][1]) return false;
  for (; i8 >= e; ) if (r5 = e + i8 >> 1, s15 > t[r5][1]) e = r5 + 1;
  else if (s15 < t[r5][0]) i8 = r5 - 1;
  else return true;
  return false;
}
var fn2 = class {
  constructor() {
    this.version = "6";
    if (!se2) {
      se2 = new Uint8Array(65536), se2.fill(1), se2[0] = 0, se2.fill(0, 1, 32), se2.fill(0, 127, 160), se2.fill(2, 4352, 4448), se2[9001] = 2, se2[9002] = 2, se2.fill(2, 11904, 42192), se2[12351] = 1, se2.fill(2, 44032, 55204), se2.fill(2, 63744, 64256), se2.fill(2, 65040, 65050), se2.fill(2, 65072, 65136), se2.fill(2, 65280, 65377), se2.fill(2, 65504, 65511);
      for (let t = 0; t < Os2.length; ++t) se2.fill(0, Os2[t][0], Os2[t][1] + 1);
    }
  }
  wcwidth(t) {
    return t < 32 ? 0 : t < 127 ? 1 : t < 65536 ? se2[t] : cc(t, ac) ? 0 : t >= 131072 && t <= 196605 || t >= 196608 && t <= 262141 ? 2 : 1;
  }
  charProperties(t, e) {
    let i8 = this.wcwidth(t), r5 = i8 === 0 && e !== 0;
    if (r5) {
      let n = Ae3.extractWidth(e);
      n === 0 ? r5 = false : n > i8 && (i8 = n);
    }
    return Ae3.createPropertyValue(0, i8, r5);
  }
};
var Ae3 = class s13 {
  constructor() {
    this._providers = /* @__PURE__ */ Object.create(null);
    this._active = "";
    this._onChange = new v();
    this.onChange = this._onChange.event;
    let t = new fn2();
    this.register(t), this._active = t.version, this._activeProvider = t;
  }
  static extractShouldJoin(t) {
    return (t & 1) !== 0;
  }
  static extractWidth(t) {
    return t >> 1 & 3;
  }
  static extractCharKind(t) {
    return t >> 3;
  }
  static createPropertyValue(t, e, i8 = false) {
    return (t & 16777215) << 3 | (e & 3) << 1 | (i8 ? 1 : 0);
  }
  dispose() {
    this._onChange.dispose();
  }
  get versions() {
    return Object.keys(this._providers);
  }
  get activeVersion() {
    return this._active;
  }
  set activeVersion(t) {
    if (!this._providers[t]) throw new Error(`unknown Unicode version "${t}"`);
    this._active = t, this._activeProvider = this._providers[t], this._onChange.fire(t);
  }
  register(t) {
    this._providers[t.version] = t;
  }
  wcwidth(t) {
    return this._activeProvider.wcwidth(t);
  }
  getStringCellWidth(t) {
    let e = 0, i8 = 0, r5 = t.length;
    for (let n = 0; n < r5; ++n) {
      let o2 = t.charCodeAt(n);
      if (55296 <= o2 && o2 <= 56319) {
        if (++n >= r5) return e + this.wcwidth(o2);
        let u = t.charCodeAt(n);
        56320 <= u && u <= 57343 ? o2 = (o2 - 55296) * 1024 + u - 56320 + 65536 : e += this.wcwidth(u);
      }
      let l = this.charProperties(o2, i8), a = s13.extractWidth(l);
      s13.extractShouldJoin(l) && (a -= s13.extractWidth(i8)), e += a, i8 = l;
    }
    return e;
  }
  charProperties(t, e) {
    return this._activeProvider.charProperties(t, e);
  }
};
var pn2 = class {
  constructor() {
    this.glevel = 0;
    this._charsets = [];
  }
  reset() {
    this.charset = void 0, this._charsets = [], this.glevel = 0;
  }
  setgLevel(t) {
    this.glevel = t, this.charset = this._charsets[t];
  }
  setgCharset(t, e) {
    this._charsets[t] = e, this.glevel === t && (this.charset = e);
  }
};
function Bs2(s15) {
  let e = s15.buffer.lines.get(s15.buffer.ybase + s15.buffer.y - 1)?.get(s15.cols - 1), i8 = s15.buffer.lines.get(s15.buffer.ybase + s15.buffer.y);
  i8 && e && (i8.isWrapped = e[3] !== 0 && e[3] !== 32);
}
var Vi2 = 2147483647;
var uc = 256;
var ci2 = class s14 {
  constructor(t = 32, e = 32) {
    this.maxLength = t;
    this.maxSubParamsLength = e;
    if (e > uc) throw new Error("maxSubParamsLength must not be greater than 256");
    this.params = new Int32Array(t), this.length = 0, this._subParams = new Int32Array(e), this._subParamsLength = 0, this._subParamsIdx = new Uint16Array(t), this._rejectDigits = false, this._rejectSubDigits = false, this._digitIsSub = false;
  }
  static fromArray(t) {
    let e = new s14();
    if (!t.length) return e;
    for (let i8 = Array.isArray(t[0]) ? 1 : 0; i8 < t.length; ++i8) {
      let r5 = t[i8];
      if (Array.isArray(r5)) for (let n = 0; n < r5.length; ++n) e.addSubParam(r5[n]);
      else e.addParam(r5);
    }
    return e;
  }
  clone() {
    let t = new s14(this.maxLength, this.maxSubParamsLength);
    return t.params.set(this.params), t.length = this.length, t._subParams.set(this._subParams), t._subParamsLength = this._subParamsLength, t._subParamsIdx.set(this._subParamsIdx), t._rejectDigits = this._rejectDigits, t._rejectSubDigits = this._rejectSubDigits, t._digitIsSub = this._digitIsSub, t;
  }
  toArray() {
    let t = [];
    for (let e = 0; e < this.length; ++e) {
      t.push(this.params[e]);
      let i8 = this._subParamsIdx[e] >> 8, r5 = this._subParamsIdx[e] & 255;
      r5 - i8 > 0 && t.push(Array.prototype.slice.call(this._subParams, i8, r5));
    }
    return t;
  }
  reset() {
    this.length = 0, this._subParamsLength = 0, this._rejectDigits = false, this._rejectSubDigits = false, this._digitIsSub = false;
  }
  addParam(t) {
    if (this._digitIsSub = false, this.length >= this.maxLength) {
      this._rejectDigits = true;
      return;
    }
    if (t < -1) throw new Error("values lesser than -1 are not allowed");
    this._subParamsIdx[this.length] = this._subParamsLength << 8 | this._subParamsLength, this.params[this.length++] = t > Vi2 ? Vi2 : t;
  }
  addSubParam(t) {
    if (this._digitIsSub = true, !!this.length) {
      if (this._rejectDigits || this._subParamsLength >= this.maxSubParamsLength) {
        this._rejectSubDigits = true;
        return;
      }
      if (t < -1) throw new Error("values lesser than -1 are not allowed");
      this._subParams[this._subParamsLength++] = t > Vi2 ? Vi2 : t, this._subParamsIdx[this.length - 1]++;
    }
  }
  hasSubParams(t) {
    return (this._subParamsIdx[t] & 255) - (this._subParamsIdx[t] >> 8) > 0;
  }
  getSubParams(t) {
    let e = this._subParamsIdx[t] >> 8, i8 = this._subParamsIdx[t] & 255;
    return i8 - e > 0 ? this._subParams.subarray(e, i8) : null;
  }
  getSubParamsAll() {
    let t = {};
    for (let e = 0; e < this.length; ++e) {
      let i8 = this._subParamsIdx[e] >> 8, r5 = this._subParamsIdx[e] & 255;
      r5 - i8 > 0 && (t[e] = this._subParams.slice(i8, r5));
    }
    return t;
  }
  addDigit(t) {
    let e;
    if (this._rejectDigits || !(e = this._digitIsSub ? this._subParamsLength : this.length) || this._digitIsSub && this._rejectSubDigits) return;
    let i8 = this._digitIsSub ? this._subParams : this.params, r5 = i8[e - 1];
    i8[e - 1] = ~r5 ? Math.min(r5 * 10 + t, Vi2) : t;
  }
};
var qi = [];
var mn2 = class {
  constructor() {
    this._state = 0;
    this._active = qi;
    this._id = -1;
    this._handlers = /* @__PURE__ */ Object.create(null);
    this._handlerFb = () => {
    };
    this._stack = { paused: false, loopPosition: 0, fallThrough: false };
  }
  registerHandler(t, e) {
    this._handlers[t] === void 0 && (this._handlers[t] = []);
    let i8 = this._handlers[t];
    return i8.push(e), { dispose: () => {
      let r5 = i8.indexOf(e);
      r5 !== -1 && i8.splice(r5, 1);
    } };
  }
  clearHandler(t) {
    this._handlers[t] && delete this._handlers[t];
  }
  setHandlerFallback(t) {
    this._handlerFb = t;
  }
  dispose() {
    this._handlers = /* @__PURE__ */ Object.create(null), this._handlerFb = () => {
    }, this._active = qi;
  }
  reset() {
    if (this._state === 2) for (let t = this._stack.paused ? this._stack.loopPosition - 1 : this._active.length - 1; t >= 0; --t) this._active[t].end(false);
    this._stack.paused = false, this._active = qi, this._id = -1, this._state = 0;
  }
  _start() {
    if (this._active = this._handlers[this._id] || qi, !this._active.length) this._handlerFb(this._id, "START");
    else for (let t = this._active.length - 1; t >= 0; t--) this._active[t].start();
  }
  _put(t, e, i8) {
    if (!this._active.length) this._handlerFb(this._id, "PUT", It2(t, e, i8));
    else for (let r5 = this._active.length - 1; r5 >= 0; r5--) this._active[r5].put(t, e, i8);
  }
  start() {
    this.reset(), this._state = 1;
  }
  put(t, e, i8) {
    if (this._state !== 3) {
      if (this._state === 1) for (; e < i8; ) {
        let r5 = t[e++];
        if (r5 === 59) {
          this._state = 2, this._start();
          break;
        }
        if (r5 < 48 || 57 < r5) {
          this._state = 3;
          return;
        }
        this._id === -1 && (this._id = 0), this._id = this._id * 10 + r5 - 48;
      }
      this._state === 2 && i8 - e > 0 && this._put(t, e, i8);
    }
  }
  end(t, e = true) {
    if (this._state !== 0) {
      if (this._state !== 3) if (this._state === 1 && this._start(), !this._active.length) this._handlerFb(this._id, "END", t);
      else {
        let i8 = false, r5 = this._active.length - 1, n = false;
        if (this._stack.paused && (r5 = this._stack.loopPosition - 1, i8 = e, n = this._stack.fallThrough, this._stack.paused = false), !n && i8 === false) {
          for (; r5 >= 0 && (i8 = this._active[r5].end(t), i8 !== true); r5--) if (i8 instanceof Promise) return this._stack.paused = true, this._stack.loopPosition = r5, this._stack.fallThrough = false, i8;
          r5--;
        }
        for (; r5 >= 0; r5--) if (i8 = this._active[r5].end(false), i8 instanceof Promise) return this._stack.paused = true, this._stack.loopPosition = r5, this._stack.fallThrough = true, i8;
      }
      this._active = qi, this._id = -1, this._state = 0;
    }
  }
};
var pe3 = class {
  constructor(t) {
    this._handler = t;
    this._data = "";
    this._hitLimit = false;
  }
  start() {
    this._data = "", this._hitLimit = false;
  }
  put(t, e, i8) {
    this._hitLimit || (this._data += It2(t, e, i8), this._data.length > 1e7 && (this._data = "", this._hitLimit = true));
  }
  end(t) {
    let e = false;
    if (this._hitLimit) e = false;
    else if (t && (e = this._handler(this._data), e instanceof Promise)) return e.then((i8) => (this._data = "", this._hitLimit = false, i8));
    return this._data = "", this._hitLimit = false, e;
  }
};
var Yi2 = [];
var _n2 = class {
  constructor() {
    this._handlers = /* @__PURE__ */ Object.create(null);
    this._active = Yi2;
    this._ident = 0;
    this._handlerFb = () => {
    };
    this._stack = { paused: false, loopPosition: 0, fallThrough: false };
  }
  dispose() {
    this._handlers = /* @__PURE__ */ Object.create(null), this._handlerFb = () => {
    }, this._active = Yi2;
  }
  registerHandler(t, e) {
    this._handlers[t] === void 0 && (this._handlers[t] = []);
    let i8 = this._handlers[t];
    return i8.push(e), { dispose: () => {
      let r5 = i8.indexOf(e);
      r5 !== -1 && i8.splice(r5, 1);
    } };
  }
  clearHandler(t) {
    this._handlers[t] && delete this._handlers[t];
  }
  setHandlerFallback(t) {
    this._handlerFb = t;
  }
  reset() {
    if (this._active.length) for (let t = this._stack.paused ? this._stack.loopPosition - 1 : this._active.length - 1; t >= 0; --t) this._active[t].unhook(false);
    this._stack.paused = false, this._active = Yi2, this._ident = 0;
  }
  hook(t, e) {
    if (this.reset(), this._ident = t, this._active = this._handlers[t] || Yi2, !this._active.length) this._handlerFb(this._ident, "HOOK", e);
    else for (let i8 = this._active.length - 1; i8 >= 0; i8--) this._active[i8].hook(e);
  }
  put(t, e, i8) {
    if (!this._active.length) this._handlerFb(this._ident, "PUT", It2(t, e, i8));
    else for (let r5 = this._active.length - 1; r5 >= 0; r5--) this._active[r5].put(t, e, i8);
  }
  unhook(t, e = true) {
    if (!this._active.length) this._handlerFb(this._ident, "UNHOOK", t);
    else {
      let i8 = false, r5 = this._active.length - 1, n = false;
      if (this._stack.paused && (r5 = this._stack.loopPosition - 1, i8 = e, n = this._stack.fallThrough, this._stack.paused = false), !n && i8 === false) {
        for (; r5 >= 0 && (i8 = this._active[r5].unhook(t), i8 !== true); r5--) if (i8 instanceof Promise) return this._stack.paused = true, this._stack.loopPosition = r5, this._stack.fallThrough = false, i8;
        r5--;
      }
      for (; r5 >= 0; r5--) if (i8 = this._active[r5].unhook(false), i8 instanceof Promise) return this._stack.paused = true, this._stack.loopPosition = r5, this._stack.fallThrough = true, i8;
    }
    this._active = Yi2, this._ident = 0;
  }
};
var ji = new ci2();
ji.addParam(0);
var Xi = class {
  constructor(t) {
    this._handler = t;
    this._data = "";
    this._params = ji;
    this._hitLimit = false;
  }
  hook(t) {
    this._params = t.length > 1 || t.params[0] ? t.clone() : ji, this._data = "", this._hitLimit = false;
  }
  put(t, e, i8) {
    this._hitLimit || (this._data += It2(t, e, i8), this._data.length > 1e7 && (this._data = "", this._hitLimit = true));
  }
  unhook(t) {
    let e = false;
    if (this._hitLimit) e = false;
    else if (t && (e = this._handler(this._data, this._params), e instanceof Promise)) return e.then((i8) => (this._params = ji, this._data = "", this._hitLimit = false, i8));
    return this._params = ji, this._data = "", this._hitLimit = false, e;
  }
};
var Fs2 = class {
  constructor(t) {
    this.table = new Uint8Array(t);
  }
  setDefault(t, e) {
    this.table.fill(t << 4 | e);
  }
  add(t, e, i8, r5) {
    this.table[e << 8 | t] = i8 << 4 | r5;
  }
  addMany(t, e, i8, r5) {
    for (let n = 0; n < t.length; n++) this.table[e << 8 | t[n]] = i8 << 4 | r5;
  }
};
var ke4 = 160;
var hc = (function() {
  let s15 = new Fs2(4095), e = Array.apply(null, Array(256)).map((a, u) => u), i8 = (a, u) => e.slice(a, u), r5 = i8(32, 127), n = i8(0, 24);
  n.push(25), n.push.apply(n, i8(28, 32));
  let o2 = i8(0, 14), l;
  s15.setDefault(1, 0), s15.addMany(r5, 0, 2, 0);
  for (l in o2) s15.addMany([24, 26, 153, 154], l, 3, 0), s15.addMany(i8(128, 144), l, 3, 0), s15.addMany(i8(144, 152), l, 3, 0), s15.add(156, l, 0, 0), s15.add(27, l, 11, 1), s15.add(157, l, 4, 8), s15.addMany([152, 158, 159], l, 0, 7), s15.add(155, l, 11, 3), s15.add(144, l, 11, 9);
  return s15.addMany(n, 0, 3, 0), s15.addMany(n, 1, 3, 1), s15.add(127, 1, 0, 1), s15.addMany(n, 8, 0, 8), s15.addMany(n, 3, 3, 3), s15.add(127, 3, 0, 3), s15.addMany(n, 4, 3, 4), s15.add(127, 4, 0, 4), s15.addMany(n, 6, 3, 6), s15.addMany(n, 5, 3, 5), s15.add(127, 5, 0, 5), s15.addMany(n, 2, 3, 2), s15.add(127, 2, 0, 2), s15.add(93, 1, 4, 8), s15.addMany(r5, 8, 5, 8), s15.add(127, 8, 5, 8), s15.addMany([156, 27, 24, 26, 7], 8, 6, 0), s15.addMany(i8(28, 32), 8, 0, 8), s15.addMany([88, 94, 95], 1, 0, 7), s15.addMany(r5, 7, 0, 7), s15.addMany(n, 7, 0, 7), s15.add(156, 7, 0, 0), s15.add(127, 7, 0, 7), s15.add(91, 1, 11, 3), s15.addMany(i8(64, 127), 3, 7, 0), s15.addMany(i8(48, 60), 3, 8, 4), s15.addMany([60, 61, 62, 63], 3, 9, 4), s15.addMany(i8(48, 60), 4, 8, 4), s15.addMany(i8(64, 127), 4, 7, 0), s15.addMany([60, 61, 62, 63], 4, 0, 6), s15.addMany(i8(32, 64), 6, 0, 6), s15.add(127, 6, 0, 6), s15.addMany(i8(64, 127), 6, 0, 0), s15.addMany(i8(32, 48), 3, 9, 5), s15.addMany(i8(32, 48), 5, 9, 5), s15.addMany(i8(48, 64), 5, 0, 6), s15.addMany(i8(64, 127), 5, 7, 0), s15.addMany(i8(32, 48), 4, 9, 5), s15.addMany(i8(32, 48), 1, 9, 2), s15.addMany(i8(32, 48), 2, 9, 2), s15.addMany(i8(48, 127), 2, 10, 0), s15.addMany(i8(48, 80), 1, 10, 0), s15.addMany(i8(81, 88), 1, 10, 0), s15.addMany([89, 90, 92], 1, 10, 0), s15.addMany(i8(96, 127), 1, 10, 0), s15.add(80, 1, 11, 9), s15.addMany(n, 9, 0, 9), s15.add(127, 9, 0, 9), s15.addMany(i8(28, 32), 9, 0, 9), s15.addMany(i8(32, 48), 9, 9, 12), s15.addMany(i8(48, 60), 9, 8, 10), s15.addMany([60, 61, 62, 63], 9, 9, 10), s15.addMany(n, 11, 0, 11), s15.addMany(i8(32, 128), 11, 0, 11), s15.addMany(i8(28, 32), 11, 0, 11), s15.addMany(n, 10, 0, 10), s15.add(127, 10, 0, 10), s15.addMany(i8(28, 32), 10, 0, 10), s15.addMany(i8(48, 60), 10, 8, 10), s15.addMany([60, 61, 62, 63], 10, 0, 11), s15.addMany(i8(32, 48), 10, 9, 12), s15.addMany(n, 12, 0, 12), s15.add(127, 12, 0, 12), s15.addMany(i8(28, 32), 12, 0, 12), s15.addMany(i8(32, 48), 12, 9, 12), s15.addMany(i8(48, 64), 12, 0, 11), s15.addMany(i8(64, 127), 12, 12, 13), s15.addMany(i8(64, 127), 10, 12, 13), s15.addMany(i8(64, 127), 9, 12, 13), s15.addMany(n, 13, 13, 13), s15.addMany(r5, 13, 13, 13), s15.add(127, 13, 0, 13), s15.addMany([27, 156, 24, 26], 13, 14, 0), s15.add(ke4, 0, 2, 0), s15.add(ke4, 8, 5, 8), s15.add(ke4, 6, 0, 6), s15.add(ke4, 11, 0, 11), s15.add(ke4, 13, 13, 13), s15;
})();
var bn2 = class extends D2 {
  constructor(e = hc) {
    super();
    this._transitions = e;
    this._parseStack = { state: 0, handlers: [], handlerPos: 0, transition: 0, chunkPos: 0 };
    this.initialState = 0, this.currentState = this.initialState, this._params = new ci2(), this._params.addParam(0), this._collect = 0, this.precedingJoinState = 0, this._printHandlerFb = (i8, r5, n) => {
    }, this._executeHandlerFb = (i8) => {
    }, this._csiHandlerFb = (i8, r5) => {
    }, this._escHandlerFb = (i8) => {
    }, this._errorHandlerFb = (i8) => i8, this._printHandler = this._printHandlerFb, this._executeHandlers = /* @__PURE__ */ Object.create(null), this._csiHandlers = /* @__PURE__ */ Object.create(null), this._escHandlers = /* @__PURE__ */ Object.create(null), this._register(C3(() => {
      this._csiHandlers = /* @__PURE__ */ Object.create(null), this._executeHandlers = /* @__PURE__ */ Object.create(null), this._escHandlers = /* @__PURE__ */ Object.create(null);
    })), this._oscParser = this._register(new mn2()), this._dcsParser = this._register(new _n2()), this._errorHandler = this._errorHandlerFb, this.registerEscHandler({ final: "\\" }, () => true);
  }
  _identifier(e, i8 = [64, 126]) {
    let r5 = 0;
    if (e.prefix) {
      if (e.prefix.length > 1) throw new Error("only one byte as prefix supported");
      if (r5 = e.prefix.charCodeAt(0), r5 && 60 > r5 || r5 > 63) throw new Error("prefix must be in range 0x3c .. 0x3f");
    }
    if (e.intermediates) {
      if (e.intermediates.length > 2) throw new Error("only two bytes as intermediates are supported");
      for (let o2 = 0; o2 < e.intermediates.length; ++o2) {
        let l = e.intermediates.charCodeAt(o2);
        if (32 > l || l > 47) throw new Error("intermediate must be in range 0x20 .. 0x2f");
        r5 <<= 8, r5 |= l;
      }
    }
    if (e.final.length !== 1) throw new Error("final must be a single byte");
    let n = e.final.charCodeAt(0);
    if (i8[0] > n || n > i8[1]) throw new Error(`final must be in range ${i8[0]} .. ${i8[1]}`);
    return r5 <<= 8, r5 |= n, r5;
  }
  identToString(e) {
    let i8 = [];
    for (; e; ) i8.push(String.fromCharCode(e & 255)), e >>= 8;
    return i8.reverse().join("");
  }
  setPrintHandler(e) {
    this._printHandler = e;
  }
  clearPrintHandler() {
    this._printHandler = this._printHandlerFb;
  }
  registerEscHandler(e, i8) {
    let r5 = this._identifier(e, [48, 126]);
    this._escHandlers[r5] === void 0 && (this._escHandlers[r5] = []);
    let n = this._escHandlers[r5];
    return n.push(i8), { dispose: () => {
      let o2 = n.indexOf(i8);
      o2 !== -1 && n.splice(o2, 1);
    } };
  }
  clearEscHandler(e) {
    this._escHandlers[this._identifier(e, [48, 126])] && delete this._escHandlers[this._identifier(e, [48, 126])];
  }
  setEscHandlerFallback(e) {
    this._escHandlerFb = e;
  }
  setExecuteHandler(e, i8) {
    this._executeHandlers[e.charCodeAt(0)] = i8;
  }
  clearExecuteHandler(e) {
    this._executeHandlers[e.charCodeAt(0)] && delete this._executeHandlers[e.charCodeAt(0)];
  }
  setExecuteHandlerFallback(e) {
    this._executeHandlerFb = e;
  }
  registerCsiHandler(e, i8) {
    let r5 = this._identifier(e);
    this._csiHandlers[r5] === void 0 && (this._csiHandlers[r5] = []);
    let n = this._csiHandlers[r5];
    return n.push(i8), { dispose: () => {
      let o2 = n.indexOf(i8);
      o2 !== -1 && n.splice(o2, 1);
    } };
  }
  clearCsiHandler(e) {
    this._csiHandlers[this._identifier(e)] && delete this._csiHandlers[this._identifier(e)];
  }
  setCsiHandlerFallback(e) {
    this._csiHandlerFb = e;
  }
  registerDcsHandler(e, i8) {
    return this._dcsParser.registerHandler(this._identifier(e), i8);
  }
  clearDcsHandler(e) {
    this._dcsParser.clearHandler(this._identifier(e));
  }
  setDcsHandlerFallback(e) {
    this._dcsParser.setHandlerFallback(e);
  }
  registerOscHandler(e, i8) {
    return this._oscParser.registerHandler(e, i8);
  }
  clearOscHandler(e) {
    this._oscParser.clearHandler(e);
  }
  setOscHandlerFallback(e) {
    this._oscParser.setHandlerFallback(e);
  }
  setErrorHandler(e) {
    this._errorHandler = e;
  }
  clearErrorHandler() {
    this._errorHandler = this._errorHandlerFb;
  }
  reset() {
    this.currentState = this.initialState, this._oscParser.reset(), this._dcsParser.reset(), this._params.reset(), this._params.addParam(0), this._collect = 0, this.precedingJoinState = 0, this._parseStack.state !== 0 && (this._parseStack.state = 2, this._parseStack.handlers = []);
  }
  _preserveStack(e, i8, r5, n, o2) {
    this._parseStack.state = e, this._parseStack.handlers = i8, this._parseStack.handlerPos = r5, this._parseStack.transition = n, this._parseStack.chunkPos = o2;
  }
  parse(e, i8, r5) {
    let n = 0, o2 = 0, l = 0, a;
    if (this._parseStack.state) if (this._parseStack.state === 2) this._parseStack.state = 0, l = this._parseStack.chunkPos + 1;
    else {
      if (r5 === void 0 || this._parseStack.state === 1) throw this._parseStack.state = 1, new Error("improper continuation due to previous async handler, giving up parsing");
      let u = this._parseStack.handlers, h2 = this._parseStack.handlerPos - 1;
      switch (this._parseStack.state) {
        case 3:
          if (r5 === false && h2 > -1) {
            for (; h2 >= 0 && (a = u[h2](this._params), a !== true); h2--) if (a instanceof Promise) return this._parseStack.handlerPos = h2, a;
          }
          this._parseStack.handlers = [];
          break;
        case 4:
          if (r5 === false && h2 > -1) {
            for (; h2 >= 0 && (a = u[h2](), a !== true); h2--) if (a instanceof Promise) return this._parseStack.handlerPos = h2, a;
          }
          this._parseStack.handlers = [];
          break;
        case 6:
          if (n = e[this._parseStack.chunkPos], a = this._dcsParser.unhook(n !== 24 && n !== 26, r5), a) return a;
          n === 27 && (this._parseStack.transition |= 1), this._params.reset(), this._params.addParam(0), this._collect = 0;
          break;
        case 5:
          if (n = e[this._parseStack.chunkPos], a = this._oscParser.end(n !== 24 && n !== 26, r5), a) return a;
          n === 27 && (this._parseStack.transition |= 1), this._params.reset(), this._params.addParam(0), this._collect = 0;
          break;
      }
      this._parseStack.state = 0, l = this._parseStack.chunkPos + 1, this.precedingJoinState = 0, this.currentState = this._parseStack.transition & 15;
    }
    for (let u = l; u < i8; ++u) {
      switch (n = e[u], o2 = this._transitions.table[this.currentState << 8 | (n < 160 ? n : ke4)], o2 >> 4) {
        case 2:
          for (let m = u + 1; ; ++m) {
            if (m >= i8 || (n = e[m]) < 32 || n > 126 && n < ke4) {
              this._printHandler(e, u, m), u = m - 1;
              break;
            }
            if (++m >= i8 || (n = e[m]) < 32 || n > 126 && n < ke4) {
              this._printHandler(e, u, m), u = m - 1;
              break;
            }
            if (++m >= i8 || (n = e[m]) < 32 || n > 126 && n < ke4) {
              this._printHandler(e, u, m), u = m - 1;
              break;
            }
            if (++m >= i8 || (n = e[m]) < 32 || n > 126 && n < ke4) {
              this._printHandler(e, u, m), u = m - 1;
              break;
            }
          }
          break;
        case 3:
          this._executeHandlers[n] ? this._executeHandlers[n]() : this._executeHandlerFb(n), this.precedingJoinState = 0;
          break;
        case 0:
          break;
        case 1:
          if (this._errorHandler({ position: u, code: n, currentState: this.currentState, collect: this._collect, params: this._params, abort: false }).abort) return;
          break;
        case 7:
          let c = this._csiHandlers[this._collect << 8 | n], d = c ? c.length - 1 : -1;
          for (; d >= 0 && (a = c[d](this._params), a !== true); d--) if (a instanceof Promise) return this._preserveStack(3, c, d, o2, u), a;
          d < 0 && this._csiHandlerFb(this._collect << 8 | n, this._params), this.precedingJoinState = 0;
          break;
        case 8:
          do
            switch (n) {
              case 59:
                this._params.addParam(0);
                break;
              case 58:
                this._params.addSubParam(-1);
                break;
              default:
                this._params.addDigit(n - 48);
            }
          while (++u < i8 && (n = e[u]) > 47 && n < 60);
          u--;
          break;
        case 9:
          this._collect <<= 8, this._collect |= n;
          break;
        case 10:
          let _3 = this._escHandlers[this._collect << 8 | n], p = _3 ? _3.length - 1 : -1;
          for (; p >= 0 && (a = _3[p](), a !== true); p--) if (a instanceof Promise) return this._preserveStack(4, _3, p, o2, u), a;
          p < 0 && this._escHandlerFb(this._collect << 8 | n), this.precedingJoinState = 0;
          break;
        case 11:
          this._params.reset(), this._params.addParam(0), this._collect = 0;
          break;
        case 12:
          this._dcsParser.hook(this._collect << 8 | n, this._params);
          break;
        case 13:
          for (let m = u + 1; ; ++m) if (m >= i8 || (n = e[m]) === 24 || n === 26 || n === 27 || n > 127 && n < ke4) {
            this._dcsParser.put(e, u, m), u = m - 1;
            break;
          }
          break;
        case 14:
          if (a = this._dcsParser.unhook(n !== 24 && n !== 26), a) return this._preserveStack(6, [], 0, o2, u), a;
          n === 27 && (o2 |= 1), this._params.reset(), this._params.addParam(0), this._collect = 0, this.precedingJoinState = 0;
          break;
        case 4:
          this._oscParser.start();
          break;
        case 5:
          for (let m = u + 1; ; m++) if (m >= i8 || (n = e[m]) < 32 || n > 127 && n < ke4) {
            this._oscParser.put(e, u, m), u = m - 1;
            break;
          }
          break;
        case 6:
          if (a = this._oscParser.end(n !== 24 && n !== 26), a) return this._preserveStack(5, [], 0, o2, u), a;
          n === 27 && (o2 |= 1), this._params.reset(), this._params.addParam(0), this._collect = 0, this.precedingJoinState = 0;
          break;
      }
      this.currentState = o2 & 15;
    }
  }
};
var dc = /^([\da-f])\/([\da-f])\/([\da-f])$|^([\da-f]{2})\/([\da-f]{2})\/([\da-f]{2})$|^([\da-f]{3})\/([\da-f]{3})\/([\da-f]{3})$|^([\da-f]{4})\/([\da-f]{4})\/([\da-f]{4})$/;
var fc = /^[\da-f]+$/;
function Ws2(s15) {
  if (!s15) return;
  let t = s15.toLowerCase();
  if (t.indexOf("rgb:") === 0) {
    t = t.slice(4);
    let e = dc.exec(t);
    if (e) {
      let i8 = e[1] ? 15 : e[4] ? 255 : e[7] ? 4095 : 65535;
      return [Math.round(parseInt(e[1] || e[4] || e[7] || e[10], 16) / i8 * 255), Math.round(parseInt(e[2] || e[5] || e[8] || e[11], 16) / i8 * 255), Math.round(parseInt(e[3] || e[6] || e[9] || e[12], 16) / i8 * 255)];
    }
  } else if (t.indexOf("#") === 0 && (t = t.slice(1), fc.exec(t) && [3, 6, 9, 12].includes(t.length))) {
    let e = t.length / 3, i8 = [0, 0, 0];
    for (let r5 = 0; r5 < 3; ++r5) {
      let n = parseInt(t.slice(e * r5, e * r5 + e), 16);
      i8[r5] = e === 1 ? n << 4 : e === 2 ? n : e === 3 ? n >> 4 : n >> 8;
    }
    return i8;
  }
}
function Hs2(s15, t) {
  let e = s15.toString(16), i8 = e.length < 2 ? "0" + e : e;
  switch (t) {
    case 4:
      return e[0];
    case 8:
      return i8;
    case 12:
      return (i8 + i8).slice(0, 3);
    default:
      return i8 + i8;
  }
}
function ml(s15, t = 16) {
  let [e, i8, r5] = s15;
  return `rgb:${Hs2(e, t)}/${Hs2(i8, t)}/${Hs2(r5, t)}`;
}
var mc = { "(": 0, ")": 1, "*": 2, "+": 3, "-": 1, ".": 2 };
var ut2 = 131072;
var _l = 10;
function bl(s15, t) {
  if (s15 > 24) return t.setWinLines || false;
  switch (s15) {
    case 1:
      return !!t.restoreWin;
    case 2:
      return !!t.minimizeWin;
    case 3:
      return !!t.setWinPosition;
    case 4:
      return !!t.setWinSizePixels;
    case 5:
      return !!t.raiseWin;
    case 6:
      return !!t.lowerWin;
    case 7:
      return !!t.refreshWin;
    case 8:
      return !!t.setWinSizeChars;
    case 9:
      return !!t.maximizeWin;
    case 10:
      return !!t.fullscreenWin;
    case 11:
      return !!t.getWinState;
    case 13:
      return !!t.getWinPosition;
    case 14:
      return !!t.getWinSizePixels;
    case 15:
      return !!t.getScreenSizePixels;
    case 16:
      return !!t.getCellSizePixels;
    case 18:
      return !!t.getWinSizeChars;
    case 19:
      return !!t.getScreenSizeChars;
    case 20:
      return !!t.getIconTitle;
    case 21:
      return !!t.getWinTitle;
    case 22:
      return !!t.pushTitle;
    case 23:
      return !!t.popTitle;
    case 24:
      return !!t.setWinLines;
  }
  return false;
}
var vl = 5e3;
var gl = 0;
var vn = class extends D2 {
  constructor(e, i8, r5, n, o2, l, a, u, h2 = new bn2()) {
    super();
    this._bufferService = e;
    this._charsetService = i8;
    this._coreService = r5;
    this._logService = n;
    this._optionsService = o2;
    this._oscLinkService = l;
    this._coreMouseService = a;
    this._unicodeService = u;
    this._parser = h2;
    this._parseBuffer = new Uint32Array(4096);
    this._stringDecoder = new er2();
    this._utf8Decoder = new tr2();
    this._windowTitle = "";
    this._iconName = "";
    this._windowTitleStack = [];
    this._iconNameStack = [];
    this._curAttrData = X4.clone();
    this._eraseAttrDataInternal = X4.clone();
    this._onRequestBell = this._register(new v());
    this.onRequestBell = this._onRequestBell.event;
    this._onRequestRefreshRows = this._register(new v());
    this.onRequestRefreshRows = this._onRequestRefreshRows.event;
    this._onRequestReset = this._register(new v());
    this.onRequestReset = this._onRequestReset.event;
    this._onRequestSendFocus = this._register(new v());
    this.onRequestSendFocus = this._onRequestSendFocus.event;
    this._onRequestSyncScrollBar = this._register(new v());
    this.onRequestSyncScrollBar = this._onRequestSyncScrollBar.event;
    this._onRequestWindowsOptionsReport = this._register(new v());
    this.onRequestWindowsOptionsReport = this._onRequestWindowsOptionsReport.event;
    this._onA11yChar = this._register(new v());
    this.onA11yChar = this._onA11yChar.event;
    this._onA11yTab = this._register(new v());
    this.onA11yTab = this._onA11yTab.event;
    this._onCursorMove = this._register(new v());
    this.onCursorMove = this._onCursorMove.event;
    this._onLineFeed = this._register(new v());
    this.onLineFeed = this._onLineFeed.event;
    this._onScroll = this._register(new v());
    this.onScroll = this._onScroll.event;
    this._onTitleChange = this._register(new v());
    this.onTitleChange = this._onTitleChange.event;
    this._onColor = this._register(new v());
    this.onColor = this._onColor.event;
    this._parseStack = { paused: false, cursorStartX: 0, cursorStartY: 0, decodedLength: 0, position: 0 };
    this._specialColors = [256, 257, 258];
    this._register(this._parser), this._dirtyRowTracker = new Zi(this._bufferService), this._activeBuffer = this._bufferService.buffer, this._register(this._bufferService.buffers.onBufferActivate((c) => this._activeBuffer = c.activeBuffer)), this._parser.setCsiHandlerFallback((c, d) => {
      this._logService.debug("Unknown CSI code: ", { identifier: this._parser.identToString(c), params: d.toArray() });
    }), this._parser.setEscHandlerFallback((c) => {
      this._logService.debug("Unknown ESC code: ", { identifier: this._parser.identToString(c) });
    }), this._parser.setExecuteHandlerFallback((c) => {
      this._logService.debug("Unknown EXECUTE code: ", { code: c });
    }), this._parser.setOscHandlerFallback((c, d, _3) => {
      this._logService.debug("Unknown OSC code: ", { identifier: c, action: d, data: _3 });
    }), this._parser.setDcsHandlerFallback((c, d, _3) => {
      d === "HOOK" && (_3 = _3.toArray()), this._logService.debug("Unknown DCS code: ", { identifier: this._parser.identToString(c), action: d, payload: _3 });
    }), this._parser.setPrintHandler((c, d, _3) => this.print(c, d, _3)), this._parser.registerCsiHandler({ final: "@" }, (c) => this.insertChars(c)), this._parser.registerCsiHandler({ intermediates: " ", final: "@" }, (c) => this.scrollLeft(c)), this._parser.registerCsiHandler({ final: "A" }, (c) => this.cursorUp(c)), this._parser.registerCsiHandler({ intermediates: " ", final: "A" }, (c) => this.scrollRight(c)), this._parser.registerCsiHandler({ final: "B" }, (c) => this.cursorDown(c)), this._parser.registerCsiHandler({ final: "C" }, (c) => this.cursorForward(c)), this._parser.registerCsiHandler({ final: "D" }, (c) => this.cursorBackward(c)), this._parser.registerCsiHandler({ final: "E" }, (c) => this.cursorNextLine(c)), this._parser.registerCsiHandler({ final: "F" }, (c) => this.cursorPrecedingLine(c)), this._parser.registerCsiHandler({ final: "G" }, (c) => this.cursorCharAbsolute(c)), this._parser.registerCsiHandler({ final: "H" }, (c) => this.cursorPosition(c)), this._parser.registerCsiHandler({ final: "I" }, (c) => this.cursorForwardTab(c)), this._parser.registerCsiHandler({ final: "J" }, (c) => this.eraseInDisplay(c, false)), this._parser.registerCsiHandler({ prefix: "?", final: "J" }, (c) => this.eraseInDisplay(c, true)), this._parser.registerCsiHandler({ final: "K" }, (c) => this.eraseInLine(c, false)), this._parser.registerCsiHandler({ prefix: "?", final: "K" }, (c) => this.eraseInLine(c, true)), this._parser.registerCsiHandler({ final: "L" }, (c) => this.insertLines(c)), this._parser.registerCsiHandler({ final: "M" }, (c) => this.deleteLines(c)), this._parser.registerCsiHandler({ final: "P" }, (c) => this.deleteChars(c)), this._parser.registerCsiHandler({ final: "S" }, (c) => this.scrollUp(c)), this._parser.registerCsiHandler({ final: "T" }, (c) => this.scrollDown(c)), this._parser.registerCsiHandler({ final: "X" }, (c) => this.eraseChars(c)), this._parser.registerCsiHandler({ final: "Z" }, (c) => this.cursorBackwardTab(c)), this._parser.registerCsiHandler({ final: "`" }, (c) => this.charPosAbsolute(c)), this._parser.registerCsiHandler({ final: "a" }, (c) => this.hPositionRelative(c)), this._parser.registerCsiHandler({ final: "b" }, (c) => this.repeatPrecedingCharacter(c)), this._parser.registerCsiHandler({ final: "c" }, (c) => this.sendDeviceAttributesPrimary(c)), this._parser.registerCsiHandler({ prefix: ">", final: "c" }, (c) => this.sendDeviceAttributesSecondary(c)), this._parser.registerCsiHandler({ final: "d" }, (c) => this.linePosAbsolute(c)), this._parser.registerCsiHandler({ final: "e" }, (c) => this.vPositionRelative(c)), this._parser.registerCsiHandler({ final: "f" }, (c) => this.hVPosition(c)), this._parser.registerCsiHandler({ final: "g" }, (c) => this.tabClear(c)), this._parser.registerCsiHandler({ final: "h" }, (c) => this.setMode(c)), this._parser.registerCsiHandler({ prefix: "?", final: "h" }, (c) => this.setModePrivate(c)), this._parser.registerCsiHandler({ final: "l" }, (c) => this.resetMode(c)), this._parser.registerCsiHandler({ prefix: "?", final: "l" }, (c) => this.resetModePrivate(c)), this._parser.registerCsiHandler({ final: "m" }, (c) => this.charAttributes(c)), this._parser.registerCsiHandler({ final: "n" }, (c) => this.deviceStatus(c)), this._parser.registerCsiHandler({ prefix: "?", final: "n" }, (c) => this.deviceStatusPrivate(c)), this._parser.registerCsiHandler({ intermediates: "!", final: "p" }, (c) => this.softReset(c)), this._parser.registerCsiHandler({ intermediates: " ", final: "q" }, (c) => this.setCursorStyle(c)), this._parser.registerCsiHandler({ final: "r" }, (c) => this.setScrollRegion(c)), this._parser.registerCsiHandler({ final: "s" }, (c) => this.saveCursor(c)), this._parser.registerCsiHandler({ final: "t" }, (c) => this.windowOptions(c)), this._parser.registerCsiHandler({ final: "u" }, (c) => this.restoreCursor(c)), this._parser.registerCsiHandler({ intermediates: "'", final: "}" }, (c) => this.insertColumns(c)), this._parser.registerCsiHandler({ intermediates: "'", final: "~" }, (c) => this.deleteColumns(c)), this._parser.registerCsiHandler({ intermediates: '"', final: "q" }, (c) => this.selectProtected(c)), this._parser.registerCsiHandler({ intermediates: "$", final: "p" }, (c) => this.requestMode(c, true)), this._parser.registerCsiHandler({ prefix: "?", intermediates: "$", final: "p" }, (c) => this.requestMode(c, false)), this._parser.setExecuteHandler(b.BEL, () => this.bell()), this._parser.setExecuteHandler(b.LF, () => this.lineFeed()), this._parser.setExecuteHandler(b.VT, () => this.lineFeed()), this._parser.setExecuteHandler(b.FF, () => this.lineFeed()), this._parser.setExecuteHandler(b.CR, () => this.carriageReturn()), this._parser.setExecuteHandler(b.BS, () => this.backspace()), this._parser.setExecuteHandler(b.HT, () => this.tab()), this._parser.setExecuteHandler(b.SO, () => this.shiftOut()), this._parser.setExecuteHandler(b.SI, () => this.shiftIn()), this._parser.setExecuteHandler(Ai2.IND, () => this.index()), this._parser.setExecuteHandler(Ai2.NEL, () => this.nextLine()), this._parser.setExecuteHandler(Ai2.HTS, () => this.tabSet()), this._parser.registerOscHandler(0, new pe3((c) => (this.setTitle(c), this.setIconName(c), true))), this._parser.registerOscHandler(1, new pe3((c) => this.setIconName(c))), this._parser.registerOscHandler(2, new pe3((c) => this.setTitle(c))), this._parser.registerOscHandler(4, new pe3((c) => this.setOrReportIndexedColor(c))), this._parser.registerOscHandler(8, new pe3((c) => this.setHyperlink(c))), this._parser.registerOscHandler(10, new pe3((c) => this.setOrReportFgColor(c))), this._parser.registerOscHandler(11, new pe3((c) => this.setOrReportBgColor(c))), this._parser.registerOscHandler(12, new pe3((c) => this.setOrReportCursorColor(c))), this._parser.registerOscHandler(104, new pe3((c) => this.restoreIndexedColor(c))), this._parser.registerOscHandler(110, new pe3((c) => this.restoreFgColor(c))), this._parser.registerOscHandler(111, new pe3((c) => this.restoreBgColor(c))), this._parser.registerOscHandler(112, new pe3((c) => this.restoreCursorColor(c))), this._parser.registerEscHandler({ final: "7" }, () => this.saveCursor()), this._parser.registerEscHandler({ final: "8" }, () => this.restoreCursor()), this._parser.registerEscHandler({ final: "D" }, () => this.index()), this._parser.registerEscHandler({ final: "E" }, () => this.nextLine()), this._parser.registerEscHandler({ final: "H" }, () => this.tabSet()), this._parser.registerEscHandler({ final: "M" }, () => this.reverseIndex()), this._parser.registerEscHandler({ final: "=" }, () => this.keypadApplicationMode()), this._parser.registerEscHandler({ final: ">" }, () => this.keypadNumericMode()), this._parser.registerEscHandler({ final: "c" }, () => this.fullReset()), this._parser.registerEscHandler({ final: "n" }, () => this.setgLevel(2)), this._parser.registerEscHandler({ final: "o" }, () => this.setgLevel(3)), this._parser.registerEscHandler({ final: "|" }, () => this.setgLevel(3)), this._parser.registerEscHandler({ final: "}" }, () => this.setgLevel(2)), this._parser.registerEscHandler({ final: "~" }, () => this.setgLevel(1)), this._parser.registerEscHandler({ intermediates: "%", final: "@" }, () => this.selectDefaultCharset()), this._parser.registerEscHandler({ intermediates: "%", final: "G" }, () => this.selectDefaultCharset());
    for (let c in ne2) this._parser.registerEscHandler({ intermediates: "(", final: c }, () => this.selectCharset("(" + c)), this._parser.registerEscHandler({ intermediates: ")", final: c }, () => this.selectCharset(")" + c)), this._parser.registerEscHandler({ intermediates: "*", final: c }, () => this.selectCharset("*" + c)), this._parser.registerEscHandler({ intermediates: "+", final: c }, () => this.selectCharset("+" + c)), this._parser.registerEscHandler({ intermediates: "-", final: c }, () => this.selectCharset("-" + c)), this._parser.registerEscHandler({ intermediates: ".", final: c }, () => this.selectCharset("." + c)), this._parser.registerEscHandler({ intermediates: "/", final: c }, () => this.selectCharset("/" + c));
    this._parser.registerEscHandler({ intermediates: "#", final: "8" }, () => this.screenAlignmentPattern()), this._parser.setErrorHandler((c) => (this._logService.error("Parsing error: ", c), c)), this._parser.registerDcsHandler({ intermediates: "$", final: "q" }, new Xi((c, d) => this.requestStatusString(c, d)));
  }
  getAttrData() {
    return this._curAttrData;
  }
  _preserveStack(e, i8, r5, n) {
    this._parseStack.paused = true, this._parseStack.cursorStartX = e, this._parseStack.cursorStartY = i8, this._parseStack.decodedLength = r5, this._parseStack.position = n;
  }
  _logSlowResolvingAsync(e) {
    this._logService.logLevel <= 3 && Promise.race([e, new Promise((i8, r5) => setTimeout(() => r5("#SLOW_TIMEOUT"), vl))]).catch((i8) => {
      if (i8 !== "#SLOW_TIMEOUT") throw i8;
      console.warn(`async parser handler taking longer than ${vl} ms`);
    });
  }
  _getCurrentLinkId() {
    return this._curAttrData.extended.urlId;
  }
  parse(e, i8) {
    let r5, n = this._activeBuffer.x, o2 = this._activeBuffer.y, l = 0, a = this._parseStack.paused;
    if (a) {
      if (r5 = this._parser.parse(this._parseBuffer, this._parseStack.decodedLength, i8)) return this._logSlowResolvingAsync(r5), r5;
      n = this._parseStack.cursorStartX, o2 = this._parseStack.cursorStartY, this._parseStack.paused = false, e.length > ut2 && (l = this._parseStack.position + ut2);
    }
    if (this._logService.logLevel <= 1 && this._logService.debug(`parsing data ${typeof e == "string" ? ` "${e}"` : ` "${Array.prototype.map.call(e, (c) => String.fromCharCode(c)).join("")}"`}`), this._logService.logLevel === 0 && this._logService.trace("parsing data (codes)", typeof e == "string" ? e.split("").map((c) => c.charCodeAt(0)) : e), this._parseBuffer.length < e.length && this._parseBuffer.length < ut2 && (this._parseBuffer = new Uint32Array(Math.min(e.length, ut2))), a || this._dirtyRowTracker.clearRange(), e.length > ut2) for (let c = l; c < e.length; c += ut2) {
      let d = c + ut2 < e.length ? c + ut2 : e.length, _3 = typeof e == "string" ? this._stringDecoder.decode(e.substring(c, d), this._parseBuffer) : this._utf8Decoder.decode(e.subarray(c, d), this._parseBuffer);
      if (r5 = this._parser.parse(this._parseBuffer, _3)) return this._preserveStack(n, o2, _3, c), this._logSlowResolvingAsync(r5), r5;
    }
    else if (!a) {
      let c = typeof e == "string" ? this._stringDecoder.decode(e, this._parseBuffer) : this._utf8Decoder.decode(e, this._parseBuffer);
      if (r5 = this._parser.parse(this._parseBuffer, c)) return this._preserveStack(n, o2, c, 0), this._logSlowResolvingAsync(r5), r5;
    }
    (this._activeBuffer.x !== n || this._activeBuffer.y !== o2) && this._onCursorMove.fire();
    let u = this._dirtyRowTracker.end + (this._bufferService.buffer.ybase - this._bufferService.buffer.ydisp), h2 = this._dirtyRowTracker.start + (this._bufferService.buffer.ybase - this._bufferService.buffer.ydisp);
    h2 < this._bufferService.rows && this._onRequestRefreshRows.fire({ start: Math.min(h2, this._bufferService.rows - 1), end: Math.min(u, this._bufferService.rows - 1) });
  }
  print(e, i8, r5) {
    let n, o2, l = this._charsetService.charset, a = this._optionsService.rawOptions.screenReaderMode, u = this._bufferService.cols, h2 = this._coreService.decPrivateModes.wraparound, c = this._coreService.modes.insertMode, d = this._curAttrData, _3 = this._activeBuffer.lines.get(this._activeBuffer.ybase + this._activeBuffer.y);
    this._dirtyRowTracker.markDirty(this._activeBuffer.y), this._activeBuffer.x && r5 - i8 > 0 && _3.getWidth(this._activeBuffer.x - 1) === 2 && _3.setCellFromCodepoint(this._activeBuffer.x - 1, 0, 1, d);
    let p = this._parser.precedingJoinState;
    for (let m = i8; m < r5; ++m) {
      if (n = e[m], n < 127 && l) {
        let O3 = l[String.fromCharCode(n)];
        O3 && (n = O3.charCodeAt(0));
      }
      let f = this._unicodeService.charProperties(n, p);
      o2 = Ae3.extractWidth(f);
      let A2 = Ae3.extractShouldJoin(f), R3 = A2 ? Ae3.extractWidth(p) : 0;
      if (p = f, a && this._onA11yChar.fire(Ce2(n)), this._getCurrentLinkId() && this._oscLinkService.addLineToLink(this._getCurrentLinkId(), this._activeBuffer.ybase + this._activeBuffer.y), this._activeBuffer.x + o2 - R3 > u) {
        if (h2) {
          let O3 = _3, I2 = this._activeBuffer.x - R3;
          for (this._activeBuffer.x = R3, this._activeBuffer.y++, this._activeBuffer.y === this._activeBuffer.scrollBottom + 1 ? (this._activeBuffer.y--, this._bufferService.scroll(this._eraseAttrData(), true)) : (this._activeBuffer.y >= this._bufferService.rows && (this._activeBuffer.y = this._bufferService.rows - 1), this._activeBuffer.lines.get(this._activeBuffer.ybase + this._activeBuffer.y).isWrapped = true), _3 = this._activeBuffer.lines.get(this._activeBuffer.ybase + this._activeBuffer.y), R3 > 0 && _3 instanceof Ze3 && _3.copyCellsFrom(O3, I2, 0, R3, false); I2 < u; ) O3.setCellFromCodepoint(I2++, 0, 1, d);
        } else if (this._activeBuffer.x = u - 1, o2 === 2) continue;
      }
      if (A2 && this._activeBuffer.x) {
        let O3 = _3.getWidth(this._activeBuffer.x - 1) ? 1 : 2;
        _3.addCodepointToCell(this._activeBuffer.x - O3, n, o2);
        for (let I2 = o2 - R3; --I2 >= 0; ) _3.setCellFromCodepoint(this._activeBuffer.x++, 0, 0, d);
        continue;
      }
      if (c && (_3.insertCells(this._activeBuffer.x, o2 - R3, this._activeBuffer.getNullCell(d)), _3.getWidth(u - 1) === 2 && _3.setCellFromCodepoint(u - 1, 0, 1, d)), _3.setCellFromCodepoint(this._activeBuffer.x++, n, o2, d), o2 > 0) for (; --o2; ) _3.setCellFromCodepoint(this._activeBuffer.x++, 0, 0, d);
    }
    this._parser.precedingJoinState = p, this._activeBuffer.x < u && r5 - i8 > 0 && _3.getWidth(this._activeBuffer.x) === 0 && !_3.hasContent(this._activeBuffer.x) && _3.setCellFromCodepoint(this._activeBuffer.x, 0, 1, d), this._dirtyRowTracker.markDirty(this._activeBuffer.y);
  }
  registerCsiHandler(e, i8) {
    return e.final === "t" && !e.prefix && !e.intermediates ? this._parser.registerCsiHandler(e, (r5) => bl(r5.params[0], this._optionsService.rawOptions.windowOptions) ? i8(r5) : true) : this._parser.registerCsiHandler(e, i8);
  }
  registerDcsHandler(e, i8) {
    return this._parser.registerDcsHandler(e, new Xi(i8));
  }
  registerEscHandler(e, i8) {
    return this._parser.registerEscHandler(e, i8);
  }
  registerOscHandler(e, i8) {
    return this._parser.registerOscHandler(e, new pe3(i8));
  }
  bell() {
    return this._onRequestBell.fire(), true;
  }
  lineFeed() {
    return this._dirtyRowTracker.markDirty(this._activeBuffer.y), this._optionsService.rawOptions.convertEol && (this._activeBuffer.x = 0), this._activeBuffer.y++, this._activeBuffer.y === this._activeBuffer.scrollBottom + 1 ? (this._activeBuffer.y--, this._bufferService.scroll(this._eraseAttrData())) : this._activeBuffer.y >= this._bufferService.rows ? this._activeBuffer.y = this._bufferService.rows - 1 : this._activeBuffer.lines.get(this._activeBuffer.ybase + this._activeBuffer.y).isWrapped = false, this._activeBuffer.x >= this._bufferService.cols && this._activeBuffer.x--, this._dirtyRowTracker.markDirty(this._activeBuffer.y), this._onLineFeed.fire(), true;
  }
  carriageReturn() {
    return this._activeBuffer.x = 0, true;
  }
  backspace() {
    if (!this._coreService.decPrivateModes.reverseWraparound) return this._restrictCursor(), this._activeBuffer.x > 0 && this._activeBuffer.x--, true;
    if (this._restrictCursor(this._bufferService.cols), this._activeBuffer.x > 0) this._activeBuffer.x--;
    else if (this._activeBuffer.x === 0 && this._activeBuffer.y > this._activeBuffer.scrollTop && this._activeBuffer.y <= this._activeBuffer.scrollBottom && this._activeBuffer.lines.get(this._activeBuffer.ybase + this._activeBuffer.y)?.isWrapped) {
      this._activeBuffer.lines.get(this._activeBuffer.ybase + this._activeBuffer.y).isWrapped = false, this._activeBuffer.y--, this._activeBuffer.x = this._bufferService.cols - 1;
      let e = this._activeBuffer.lines.get(this._activeBuffer.ybase + this._activeBuffer.y);
      e.hasWidth(this._activeBuffer.x) && !e.hasContent(this._activeBuffer.x) && this._activeBuffer.x--;
    }
    return this._restrictCursor(), true;
  }
  tab() {
    if (this._activeBuffer.x >= this._bufferService.cols) return true;
    let e = this._activeBuffer.x;
    return this._activeBuffer.x = this._activeBuffer.nextStop(), this._optionsService.rawOptions.screenReaderMode && this._onA11yTab.fire(this._activeBuffer.x - e), true;
  }
  shiftOut() {
    return this._charsetService.setgLevel(1), true;
  }
  shiftIn() {
    return this._charsetService.setgLevel(0), true;
  }
  _restrictCursor(e = this._bufferService.cols - 1) {
    this._activeBuffer.x = Math.min(e, Math.max(0, this._activeBuffer.x)), this._activeBuffer.y = this._coreService.decPrivateModes.origin ? Math.min(this._activeBuffer.scrollBottom, Math.max(this._activeBuffer.scrollTop, this._activeBuffer.y)) : Math.min(this._bufferService.rows - 1, Math.max(0, this._activeBuffer.y)), this._dirtyRowTracker.markDirty(this._activeBuffer.y);
  }
  _setCursor(e, i8) {
    this._dirtyRowTracker.markDirty(this._activeBuffer.y), this._coreService.decPrivateModes.origin ? (this._activeBuffer.x = e, this._activeBuffer.y = this._activeBuffer.scrollTop + i8) : (this._activeBuffer.x = e, this._activeBuffer.y = i8), this._restrictCursor(), this._dirtyRowTracker.markDirty(this._activeBuffer.y);
  }
  _moveCursor(e, i8) {
    this._restrictCursor(), this._setCursor(this._activeBuffer.x + e, this._activeBuffer.y + i8);
  }
  cursorUp(e) {
    let i8 = this._activeBuffer.y - this._activeBuffer.scrollTop;
    return i8 >= 0 ? this._moveCursor(0, -Math.min(i8, e.params[0] || 1)) : this._moveCursor(0, -(e.params[0] || 1)), true;
  }
  cursorDown(e) {
    let i8 = this._activeBuffer.scrollBottom - this._activeBuffer.y;
    return i8 >= 0 ? this._moveCursor(0, Math.min(i8, e.params[0] || 1)) : this._moveCursor(0, e.params[0] || 1), true;
  }
  cursorForward(e) {
    return this._moveCursor(e.params[0] || 1, 0), true;
  }
  cursorBackward(e) {
    return this._moveCursor(-(e.params[0] || 1), 0), true;
  }
  cursorNextLine(e) {
    return this.cursorDown(e), this._activeBuffer.x = 0, true;
  }
  cursorPrecedingLine(e) {
    return this.cursorUp(e), this._activeBuffer.x = 0, true;
  }
  cursorCharAbsolute(e) {
    return this._setCursor((e.params[0] || 1) - 1, this._activeBuffer.y), true;
  }
  cursorPosition(e) {
    return this._setCursor(e.length >= 2 ? (e.params[1] || 1) - 1 : 0, (e.params[0] || 1) - 1), true;
  }
  charPosAbsolute(e) {
    return this._setCursor((e.params[0] || 1) - 1, this._activeBuffer.y), true;
  }
  hPositionRelative(e) {
    return this._moveCursor(e.params[0] || 1, 0), true;
  }
  linePosAbsolute(e) {
    return this._setCursor(this._activeBuffer.x, (e.params[0] || 1) - 1), true;
  }
  vPositionRelative(e) {
    return this._moveCursor(0, e.params[0] || 1), true;
  }
  hVPosition(e) {
    return this.cursorPosition(e), true;
  }
  tabClear(e) {
    let i8 = e.params[0];
    return i8 === 0 ? delete this._activeBuffer.tabs[this._activeBuffer.x] : i8 === 3 && (this._activeBuffer.tabs = {}), true;
  }
  cursorForwardTab(e) {
    if (this._activeBuffer.x >= this._bufferService.cols) return true;
    let i8 = e.params[0] || 1;
    for (; i8--; ) this._activeBuffer.x = this._activeBuffer.nextStop();
    return true;
  }
  cursorBackwardTab(e) {
    if (this._activeBuffer.x >= this._bufferService.cols) return true;
    let i8 = e.params[0] || 1;
    for (; i8--; ) this._activeBuffer.x = this._activeBuffer.prevStop();
    return true;
  }
  selectProtected(e) {
    let i8 = e.params[0];
    return i8 === 1 && (this._curAttrData.bg |= 536870912), (i8 === 2 || i8 === 0) && (this._curAttrData.bg &= -536870913), true;
  }
  _eraseInBufferLine(e, i8, r5, n = false, o2 = false) {
    let l = this._activeBuffer.lines.get(this._activeBuffer.ybase + e);
    l.replaceCells(i8, r5, this._activeBuffer.getNullCell(this._eraseAttrData()), o2), n && (l.isWrapped = false);
  }
  _resetBufferLine(e, i8 = false) {
    let r5 = this._activeBuffer.lines.get(this._activeBuffer.ybase + e);
    r5 && (r5.fill(this._activeBuffer.getNullCell(this._eraseAttrData()), i8), this._bufferService.buffer.clearMarkers(this._activeBuffer.ybase + e), r5.isWrapped = false);
  }
  eraseInDisplay(e, i8 = false) {
    this._restrictCursor(this._bufferService.cols);
    let r5;
    switch (e.params[0]) {
      case 0:
        for (r5 = this._activeBuffer.y, this._dirtyRowTracker.markDirty(r5), this._eraseInBufferLine(r5++, this._activeBuffer.x, this._bufferService.cols, this._activeBuffer.x === 0, i8); r5 < this._bufferService.rows; r5++) this._resetBufferLine(r5, i8);
        this._dirtyRowTracker.markDirty(r5);
        break;
      case 1:
        for (r5 = this._activeBuffer.y, this._dirtyRowTracker.markDirty(r5), this._eraseInBufferLine(r5, 0, this._activeBuffer.x + 1, true, i8), this._activeBuffer.x + 1 >= this._bufferService.cols && (this._activeBuffer.lines.get(r5 + 1).isWrapped = false); r5--; ) this._resetBufferLine(r5, i8);
        this._dirtyRowTracker.markDirty(0);
        break;
      case 2:
        if (this._optionsService.rawOptions.scrollOnEraseInDisplay) {
          for (r5 = this._bufferService.rows, this._dirtyRowTracker.markRangeDirty(0, r5 - 1); r5-- && !this._activeBuffer.lines.get(this._activeBuffer.ybase + r5)?.getTrimmedLength(); ) ;
          for (; r5 >= 0; r5--) this._bufferService.scroll(this._eraseAttrData());
        } else {
          for (r5 = this._bufferService.rows, this._dirtyRowTracker.markDirty(r5 - 1); r5--; ) this._resetBufferLine(r5, i8);
          this._dirtyRowTracker.markDirty(0);
        }
        break;
      case 3:
        let n = this._activeBuffer.lines.length - this._bufferService.rows;
        n > 0 && (this._activeBuffer.lines.trimStart(n), this._activeBuffer.ybase = Math.max(this._activeBuffer.ybase - n, 0), this._activeBuffer.ydisp = Math.max(this._activeBuffer.ydisp - n, 0), this._onScroll.fire(0));
        break;
    }
    return true;
  }
  eraseInLine(e, i8 = false) {
    switch (this._restrictCursor(this._bufferService.cols), e.params[0]) {
      case 0:
        this._eraseInBufferLine(this._activeBuffer.y, this._activeBuffer.x, this._bufferService.cols, this._activeBuffer.x === 0, i8);
        break;
      case 1:
        this._eraseInBufferLine(this._activeBuffer.y, 0, this._activeBuffer.x + 1, false, i8);
        break;
      case 2:
        this._eraseInBufferLine(this._activeBuffer.y, 0, this._bufferService.cols, true, i8);
        break;
    }
    return this._dirtyRowTracker.markDirty(this._activeBuffer.y), true;
  }
  insertLines(e) {
    this._restrictCursor();
    let i8 = e.params[0] || 1;
    if (this._activeBuffer.y > this._activeBuffer.scrollBottom || this._activeBuffer.y < this._activeBuffer.scrollTop) return true;
    let r5 = this._activeBuffer.ybase + this._activeBuffer.y, n = this._bufferService.rows - 1 - this._activeBuffer.scrollBottom, o2 = this._bufferService.rows - 1 + this._activeBuffer.ybase - n + 1;
    for (; i8--; ) this._activeBuffer.lines.splice(o2 - 1, 1), this._activeBuffer.lines.splice(r5, 0, this._activeBuffer.getBlankLine(this._eraseAttrData()));
    return this._dirtyRowTracker.markRangeDirty(this._activeBuffer.y, this._activeBuffer.scrollBottom), this._activeBuffer.x = 0, true;
  }
  deleteLines(e) {
    this._restrictCursor();
    let i8 = e.params[0] || 1;
    if (this._activeBuffer.y > this._activeBuffer.scrollBottom || this._activeBuffer.y < this._activeBuffer.scrollTop) return true;
    let r5 = this._activeBuffer.ybase + this._activeBuffer.y, n;
    for (n = this._bufferService.rows - 1 - this._activeBuffer.scrollBottom, n = this._bufferService.rows - 1 + this._activeBuffer.ybase - n; i8--; ) this._activeBuffer.lines.splice(r5, 1), this._activeBuffer.lines.splice(n, 0, this._activeBuffer.getBlankLine(this._eraseAttrData()));
    return this._dirtyRowTracker.markRangeDirty(this._activeBuffer.y, this._activeBuffer.scrollBottom), this._activeBuffer.x = 0, true;
  }
  insertChars(e) {
    this._restrictCursor();
    let i8 = this._activeBuffer.lines.get(this._activeBuffer.ybase + this._activeBuffer.y);
    return i8 && (i8.insertCells(this._activeBuffer.x, e.params[0] || 1, this._activeBuffer.getNullCell(this._eraseAttrData())), this._dirtyRowTracker.markDirty(this._activeBuffer.y)), true;
  }
  deleteChars(e) {
    this._restrictCursor();
    let i8 = this._activeBuffer.lines.get(this._activeBuffer.ybase + this._activeBuffer.y);
    return i8 && (i8.deleteCells(this._activeBuffer.x, e.params[0] || 1, this._activeBuffer.getNullCell(this._eraseAttrData())), this._dirtyRowTracker.markDirty(this._activeBuffer.y)), true;
  }
  scrollUp(e) {
    let i8 = e.params[0] || 1;
    for (; i8--; ) this._activeBuffer.lines.splice(this._activeBuffer.ybase + this._activeBuffer.scrollTop, 1), this._activeBuffer.lines.splice(this._activeBuffer.ybase + this._activeBuffer.scrollBottom, 0, this._activeBuffer.getBlankLine(this._eraseAttrData()));
    return this._dirtyRowTracker.markRangeDirty(this._activeBuffer.scrollTop, this._activeBuffer.scrollBottom), true;
  }
  scrollDown(e) {
    let i8 = e.params[0] || 1;
    for (; i8--; ) this._activeBuffer.lines.splice(this._activeBuffer.ybase + this._activeBuffer.scrollBottom, 1), this._activeBuffer.lines.splice(this._activeBuffer.ybase + this._activeBuffer.scrollTop, 0, this._activeBuffer.getBlankLine(X4));
    return this._dirtyRowTracker.markRangeDirty(this._activeBuffer.scrollTop, this._activeBuffer.scrollBottom), true;
  }
  scrollLeft(e) {
    if (this._activeBuffer.y > this._activeBuffer.scrollBottom || this._activeBuffer.y < this._activeBuffer.scrollTop) return true;
    let i8 = e.params[0] || 1;
    for (let r5 = this._activeBuffer.scrollTop; r5 <= this._activeBuffer.scrollBottom; ++r5) {
      let n = this._activeBuffer.lines.get(this._activeBuffer.ybase + r5);
      n.deleteCells(0, i8, this._activeBuffer.getNullCell(this._eraseAttrData())), n.isWrapped = false;
    }
    return this._dirtyRowTracker.markRangeDirty(this._activeBuffer.scrollTop, this._activeBuffer.scrollBottom), true;
  }
  scrollRight(e) {
    if (this._activeBuffer.y > this._activeBuffer.scrollBottom || this._activeBuffer.y < this._activeBuffer.scrollTop) return true;
    let i8 = e.params[0] || 1;
    for (let r5 = this._activeBuffer.scrollTop; r5 <= this._activeBuffer.scrollBottom; ++r5) {
      let n = this._activeBuffer.lines.get(this._activeBuffer.ybase + r5);
      n.insertCells(0, i8, this._activeBuffer.getNullCell(this._eraseAttrData())), n.isWrapped = false;
    }
    return this._dirtyRowTracker.markRangeDirty(this._activeBuffer.scrollTop, this._activeBuffer.scrollBottom), true;
  }
  insertColumns(e) {
    if (this._activeBuffer.y > this._activeBuffer.scrollBottom || this._activeBuffer.y < this._activeBuffer.scrollTop) return true;
    let i8 = e.params[0] || 1;
    for (let r5 = this._activeBuffer.scrollTop; r5 <= this._activeBuffer.scrollBottom; ++r5) {
      let n = this._activeBuffer.lines.get(this._activeBuffer.ybase + r5);
      n.insertCells(this._activeBuffer.x, i8, this._activeBuffer.getNullCell(this._eraseAttrData())), n.isWrapped = false;
    }
    return this._dirtyRowTracker.markRangeDirty(this._activeBuffer.scrollTop, this._activeBuffer.scrollBottom), true;
  }
  deleteColumns(e) {
    if (this._activeBuffer.y > this._activeBuffer.scrollBottom || this._activeBuffer.y < this._activeBuffer.scrollTop) return true;
    let i8 = e.params[0] || 1;
    for (let r5 = this._activeBuffer.scrollTop; r5 <= this._activeBuffer.scrollBottom; ++r5) {
      let n = this._activeBuffer.lines.get(this._activeBuffer.ybase + r5);
      n.deleteCells(this._activeBuffer.x, i8, this._activeBuffer.getNullCell(this._eraseAttrData())), n.isWrapped = false;
    }
    return this._dirtyRowTracker.markRangeDirty(this._activeBuffer.scrollTop, this._activeBuffer.scrollBottom), true;
  }
  eraseChars(e) {
    this._restrictCursor();
    let i8 = this._activeBuffer.lines.get(this._activeBuffer.ybase + this._activeBuffer.y);
    return i8 && (i8.replaceCells(this._activeBuffer.x, this._activeBuffer.x + (e.params[0] || 1), this._activeBuffer.getNullCell(this._eraseAttrData())), this._dirtyRowTracker.markDirty(this._activeBuffer.y)), true;
  }
  repeatPrecedingCharacter(e) {
    let i8 = this._parser.precedingJoinState;
    if (!i8) return true;
    let r5 = e.params[0] || 1, n = Ae3.extractWidth(i8), o2 = this._activeBuffer.x - n, a = this._activeBuffer.lines.get(this._activeBuffer.ybase + this._activeBuffer.y).getString(o2), u = new Uint32Array(a.length * r5), h2 = 0;
    for (let d = 0; d < a.length; ) {
      let _3 = a.codePointAt(d) || 0;
      u[h2++] = _3, d += _3 > 65535 ? 2 : 1;
    }
    let c = h2;
    for (let d = 1; d < r5; ++d) u.copyWithin(c, 0, h2), c += h2;
    return this.print(u, 0, c), true;
  }
  sendDeviceAttributesPrimary(e) {
    return e.params[0] > 0 || (this._is("xterm") || this._is("rxvt-unicode") || this._is("screen") ? this._coreService.triggerDataEvent(b.ESC + "[?1;2c") : this._is("linux") && this._coreService.triggerDataEvent(b.ESC + "[?6c")), true;
  }
  sendDeviceAttributesSecondary(e) {
    return e.params[0] > 0 || (this._is("xterm") ? this._coreService.triggerDataEvent(b.ESC + "[>0;276;0c") : this._is("rxvt-unicode") ? this._coreService.triggerDataEvent(b.ESC + "[>85;95;0c") : this._is("linux") ? this._coreService.triggerDataEvent(e.params[0] + "c") : this._is("screen") && this._coreService.triggerDataEvent(b.ESC + "[>83;40003;0c")), true;
  }
  _is(e) {
    return (this._optionsService.rawOptions.termName + "").indexOf(e) === 0;
  }
  setMode(e) {
    for (let i8 = 0; i8 < e.length; i8++) switch (e.params[i8]) {
      case 4:
        this._coreService.modes.insertMode = true;
        break;
      case 20:
        this._optionsService.options.convertEol = true;
        break;
    }
    return true;
  }
  setModePrivate(e) {
    for (let i8 = 0; i8 < e.length; i8++) switch (e.params[i8]) {
      case 1:
        this._coreService.decPrivateModes.applicationCursorKeys = true;
        break;
      case 2:
        this._charsetService.setgCharset(0, Je2), this._charsetService.setgCharset(1, Je2), this._charsetService.setgCharset(2, Je2), this._charsetService.setgCharset(3, Je2);
        break;
      case 3:
        this._optionsService.rawOptions.windowOptions.setWinLines && (this._bufferService.resize(132, this._bufferService.rows), this._onRequestReset.fire());
        break;
      case 6:
        this._coreService.decPrivateModes.origin = true, this._setCursor(0, 0);
        break;
      case 7:
        this._coreService.decPrivateModes.wraparound = true;
        break;
      case 12:
        this._optionsService.options.cursorBlink = true;
        break;
      case 45:
        this._coreService.decPrivateModes.reverseWraparound = true;
        break;
      case 66:
        this._logService.debug("Serial port requested application keypad."), this._coreService.decPrivateModes.applicationKeypad = true, this._onRequestSyncScrollBar.fire();
        break;
      case 9:
        this._coreMouseService.activeProtocol = "X10";
        break;
      case 1e3:
        this._coreMouseService.activeProtocol = "VT200";
        break;
      case 1002:
        this._coreMouseService.activeProtocol = "DRAG";
        break;
      case 1003:
        this._coreMouseService.activeProtocol = "ANY";
        break;
      case 1004:
        this._coreService.decPrivateModes.sendFocus = true, this._onRequestSendFocus.fire();
        break;
      case 1005:
        this._logService.debug("DECSET 1005 not supported (see #2507)");
        break;
      case 1006:
        this._coreMouseService.activeEncoding = "SGR";
        break;
      case 1015:
        this._logService.debug("DECSET 1015 not supported (see #2507)");
        break;
      case 1016:
        this._coreMouseService.activeEncoding = "SGR_PIXELS";
        break;
      case 25:
        this._coreService.isCursorHidden = false;
        break;
      case 1048:
        this.saveCursor();
        break;
      case 1049:
        this.saveCursor();
      case 47:
      case 1047:
        this._bufferService.buffers.activateAltBuffer(this._eraseAttrData()), this._coreService.isCursorInitialized = true, this._onRequestRefreshRows.fire(void 0), this._onRequestSyncScrollBar.fire();
        break;
      case 2004:
        this._coreService.decPrivateModes.bracketedPasteMode = true;
        break;
      case 2026:
        this._coreService.decPrivateModes.synchronizedOutput = true;
        break;
    }
    return true;
  }
  resetMode(e) {
    for (let i8 = 0; i8 < e.length; i8++) switch (e.params[i8]) {
      case 4:
        this._coreService.modes.insertMode = false;
        break;
      case 20:
        this._optionsService.options.convertEol = false;
        break;
    }
    return true;
  }
  resetModePrivate(e) {
    for (let i8 = 0; i8 < e.length; i8++) switch (e.params[i8]) {
      case 1:
        this._coreService.decPrivateModes.applicationCursorKeys = false;
        break;
      case 3:
        this._optionsService.rawOptions.windowOptions.setWinLines && (this._bufferService.resize(80, this._bufferService.rows), this._onRequestReset.fire());
        break;
      case 6:
        this._coreService.decPrivateModes.origin = false, this._setCursor(0, 0);
        break;
      case 7:
        this._coreService.decPrivateModes.wraparound = false;
        break;
      case 12:
        this._optionsService.options.cursorBlink = false;
        break;
      case 45:
        this._coreService.decPrivateModes.reverseWraparound = false;
        break;
      case 66:
        this._logService.debug("Switching back to normal keypad."), this._coreService.decPrivateModes.applicationKeypad = false, this._onRequestSyncScrollBar.fire();
        break;
      case 9:
      case 1e3:
      case 1002:
      case 1003:
        this._coreMouseService.activeProtocol = "NONE";
        break;
      case 1004:
        this._coreService.decPrivateModes.sendFocus = false;
        break;
      case 1005:
        this._logService.debug("DECRST 1005 not supported (see #2507)");
        break;
      case 1006:
        this._coreMouseService.activeEncoding = "DEFAULT";
        break;
      case 1015:
        this._logService.debug("DECRST 1015 not supported (see #2507)");
        break;
      case 1016:
        this._coreMouseService.activeEncoding = "DEFAULT";
        break;
      case 25:
        this._coreService.isCursorHidden = true;
        break;
      case 1048:
        this.restoreCursor();
        break;
      case 1049:
      case 47:
      case 1047:
        this._bufferService.buffers.activateNormalBuffer(), e.params[i8] === 1049 && this.restoreCursor(), this._coreService.isCursorInitialized = true, this._onRequestRefreshRows.fire(void 0), this._onRequestSyncScrollBar.fire();
        break;
      case 2004:
        this._coreService.decPrivateModes.bracketedPasteMode = false;
        break;
      case 2026:
        this._coreService.decPrivateModes.synchronizedOutput = false, this._onRequestRefreshRows.fire(void 0);
        break;
    }
    return true;
  }
  requestMode(e, i8) {
    let r5;
    ((P3) => (P3[P3.NOT_RECOGNIZED = 0] = "NOT_RECOGNIZED", P3[P3.SET = 1] = "SET", P3[P3.RESET = 2] = "RESET", P3[P3.PERMANENTLY_SET = 3] = "PERMANENTLY_SET", P3[P3.PERMANENTLY_RESET = 4] = "PERMANENTLY_RESET"))(r5 ||= {});
    let n = this._coreService.decPrivateModes, { activeProtocol: o2, activeEncoding: l } = this._coreMouseService, a = this._coreService, { buffers: u, cols: h2 } = this._bufferService, { active: c, alt: d } = u, _3 = this._optionsService.rawOptions, p = (A2, R3) => (a.triggerDataEvent(`${b.ESC}[${i8 ? "" : "?"}${A2};${R3}$y`), true), m = (A2) => A2 ? 1 : 2, f = e.params[0];
    return i8 ? f === 2 ? p(f, 4) : f === 4 ? p(f, m(a.modes.insertMode)) : f === 12 ? p(f, 3) : f === 20 ? p(f, m(_3.convertEol)) : p(f, 0) : f === 1 ? p(f, m(n.applicationCursorKeys)) : f === 3 ? p(f, _3.windowOptions.setWinLines ? h2 === 80 ? 2 : h2 === 132 ? 1 : 0 : 0) : f === 6 ? p(f, m(n.origin)) : f === 7 ? p(f, m(n.wraparound)) : f === 8 ? p(f, 3) : f === 9 ? p(f, m(o2 === "X10")) : f === 12 ? p(f, m(_3.cursorBlink)) : f === 25 ? p(f, m(!a.isCursorHidden)) : f === 45 ? p(f, m(n.reverseWraparound)) : f === 66 ? p(f, m(n.applicationKeypad)) : f === 67 ? p(f, 4) : f === 1e3 ? p(f, m(o2 === "VT200")) : f === 1002 ? p(f, m(o2 === "DRAG")) : f === 1003 ? p(f, m(o2 === "ANY")) : f === 1004 ? p(f, m(n.sendFocus)) : f === 1005 ? p(f, 4) : f === 1006 ? p(f, m(l === "SGR")) : f === 1015 ? p(f, 4) : f === 1016 ? p(f, m(l === "SGR_PIXELS")) : f === 1048 ? p(f, 1) : f === 47 || f === 1047 || f === 1049 ? p(f, m(c === d)) : f === 2004 ? p(f, m(n.bracketedPasteMode)) : f === 2026 ? p(f, m(n.synchronizedOutput)) : p(f, 0);
  }
  _updateAttrColor(e, i8, r5, n, o2) {
    return i8 === 2 ? (e |= 50331648, e &= -16777216, e |= De3.fromColorRGB([r5, n, o2])) : i8 === 5 && (e &= -50331904, e |= 33554432 | r5 & 255), e;
  }
  _extractColor(e, i8, r5) {
    let n = [0, 0, -1, 0, 0, 0], o2 = 0, l = 0;
    do {
      if (n[l + o2] = e.params[i8 + l], e.hasSubParams(i8 + l)) {
        let a = e.getSubParams(i8 + l), u = 0;
        do
          n[1] === 5 && (o2 = 1), n[l + u + 1 + o2] = a[u];
        while (++u < a.length && u + l + 1 + o2 < n.length);
        break;
      }
      if (n[1] === 5 && l + o2 >= 2 || n[1] === 2 && l + o2 >= 5) break;
      n[1] && (o2 = 1);
    } while (++l + i8 < e.length && l + o2 < n.length);
    for (let a = 2; a < n.length; ++a) n[a] === -1 && (n[a] = 0);
    switch (n[0]) {
      case 38:
        r5.fg = this._updateAttrColor(r5.fg, n[1], n[3], n[4], n[5]);
        break;
      case 48:
        r5.bg = this._updateAttrColor(r5.bg, n[1], n[3], n[4], n[5]);
        break;
      case 58:
        r5.extended = r5.extended.clone(), r5.extended.underlineColor = this._updateAttrColor(r5.extended.underlineColor, n[1], n[3], n[4], n[5]);
    }
    return l;
  }
  _processUnderline(e, i8) {
    i8.extended = i8.extended.clone(), (!~e || e > 5) && (e = 1), i8.extended.underlineStyle = e, i8.fg |= 268435456, e === 0 && (i8.fg &= -268435457), i8.updateExtended();
  }
  _processSGR0(e) {
    e.fg = X4.fg, e.bg = X4.bg, e.extended = e.extended.clone(), e.extended.underlineStyle = 0, e.extended.underlineColor &= -67108864, e.updateExtended();
  }
  charAttributes(e) {
    if (e.length === 1 && e.params[0] === 0) return this._processSGR0(this._curAttrData), true;
    let i8 = e.length, r5, n = this._curAttrData;
    for (let o2 = 0; o2 < i8; o2++) r5 = e.params[o2], r5 >= 30 && r5 <= 37 ? (n.fg &= -50331904, n.fg |= 16777216 | r5 - 30) : r5 >= 40 && r5 <= 47 ? (n.bg &= -50331904, n.bg |= 16777216 | r5 - 40) : r5 >= 90 && r5 <= 97 ? (n.fg &= -50331904, n.fg |= 16777216 | r5 - 90 | 8) : r5 >= 100 && r5 <= 107 ? (n.bg &= -50331904, n.bg |= 16777216 | r5 - 100 | 8) : r5 === 0 ? this._processSGR0(n) : r5 === 1 ? n.fg |= 134217728 : r5 === 3 ? n.bg |= 67108864 : r5 === 4 ? (n.fg |= 268435456, this._processUnderline(e.hasSubParams(o2) ? e.getSubParams(o2)[0] : 1, n)) : r5 === 5 ? n.fg |= 536870912 : r5 === 7 ? n.fg |= 67108864 : r5 === 8 ? n.fg |= 1073741824 : r5 === 9 ? n.fg |= 2147483648 : r5 === 2 ? n.bg |= 134217728 : r5 === 21 ? this._processUnderline(2, n) : r5 === 22 ? (n.fg &= -134217729, n.bg &= -134217729) : r5 === 23 ? n.bg &= -67108865 : r5 === 24 ? (n.fg &= -268435457, this._processUnderline(0, n)) : r5 === 25 ? n.fg &= -536870913 : r5 === 27 ? n.fg &= -67108865 : r5 === 28 ? n.fg &= -1073741825 : r5 === 29 ? n.fg &= 2147483647 : r5 === 39 ? (n.fg &= -67108864, n.fg |= X4.fg & 16777215) : r5 === 49 ? (n.bg &= -67108864, n.bg |= X4.bg & 16777215) : r5 === 38 || r5 === 48 || r5 === 58 ? o2 += this._extractColor(e, o2, n) : r5 === 53 ? n.bg |= 1073741824 : r5 === 55 ? n.bg &= -1073741825 : r5 === 59 ? (n.extended = n.extended.clone(), n.extended.underlineColor = -1, n.updateExtended()) : r5 === 100 ? (n.fg &= -67108864, n.fg |= X4.fg & 16777215, n.bg &= -67108864, n.bg |= X4.bg & 16777215) : this._logService.debug("Unknown SGR attribute: %d.", r5);
    return true;
  }
  deviceStatus(e) {
    switch (e.params[0]) {
      case 5:
        this._coreService.triggerDataEvent(`${b.ESC}[0n`);
        break;
      case 6:
        let i8 = this._activeBuffer.y + 1, r5 = this._activeBuffer.x + 1;
        this._coreService.triggerDataEvent(`${b.ESC}[${i8};${r5}R`);
        break;
    }
    return true;
  }
  deviceStatusPrivate(e) {
    switch (e.params[0]) {
      case 6:
        let i8 = this._activeBuffer.y + 1, r5 = this._activeBuffer.x + 1;
        this._coreService.triggerDataEvent(`${b.ESC}[?${i8};${r5}R`);
        break;
      case 15:
        break;
      case 25:
        break;
      case 26:
        break;
      case 53:
        break;
    }
    return true;
  }
  softReset(e) {
    return this._coreService.isCursorHidden = false, this._onRequestSyncScrollBar.fire(), this._activeBuffer.scrollTop = 0, this._activeBuffer.scrollBottom = this._bufferService.rows - 1, this._curAttrData = X4.clone(), this._coreService.reset(), this._charsetService.reset(), this._activeBuffer.savedX = 0, this._activeBuffer.savedY = this._activeBuffer.ybase, this._activeBuffer.savedCurAttrData.fg = this._curAttrData.fg, this._activeBuffer.savedCurAttrData.bg = this._curAttrData.bg, this._activeBuffer.savedCharset = this._charsetService.charset, this._coreService.decPrivateModes.origin = false, true;
  }
  setCursorStyle(e) {
    let i8 = e.length === 0 ? 1 : e.params[0];
    if (i8 === 0) this._coreService.decPrivateModes.cursorStyle = void 0, this._coreService.decPrivateModes.cursorBlink = void 0;
    else {
      switch (i8) {
        case 1:
        case 2:
          this._coreService.decPrivateModes.cursorStyle = "block";
          break;
        case 3:
        case 4:
          this._coreService.decPrivateModes.cursorStyle = "underline";
          break;
        case 5:
        case 6:
          this._coreService.decPrivateModes.cursorStyle = "bar";
          break;
      }
      let r5 = i8 % 2 === 1;
      this._coreService.decPrivateModes.cursorBlink = r5;
    }
    return true;
  }
  setScrollRegion(e) {
    let i8 = e.params[0] || 1, r5;
    return (e.length < 2 || (r5 = e.params[1]) > this._bufferService.rows || r5 === 0) && (r5 = this._bufferService.rows), r5 > i8 && (this._activeBuffer.scrollTop = i8 - 1, this._activeBuffer.scrollBottom = r5 - 1, this._setCursor(0, 0)), true;
  }
  windowOptions(e) {
    if (!bl(e.params[0], this._optionsService.rawOptions.windowOptions)) return true;
    let i8 = e.length > 1 ? e.params[1] : 0;
    switch (e.params[0]) {
      case 14:
        i8 !== 2 && this._onRequestWindowsOptionsReport.fire(0);
        break;
      case 16:
        this._onRequestWindowsOptionsReport.fire(1);
        break;
      case 18:
        this._bufferService && this._coreService.triggerDataEvent(`${b.ESC}[8;${this._bufferService.rows};${this._bufferService.cols}t`);
        break;
      case 22:
        (i8 === 0 || i8 === 2) && (this._windowTitleStack.push(this._windowTitle), this._windowTitleStack.length > _l && this._windowTitleStack.shift()), (i8 === 0 || i8 === 1) && (this._iconNameStack.push(this._iconName), this._iconNameStack.length > _l && this._iconNameStack.shift());
        break;
      case 23:
        (i8 === 0 || i8 === 2) && this._windowTitleStack.length && this.setTitle(this._windowTitleStack.pop()), (i8 === 0 || i8 === 1) && this._iconNameStack.length && this.setIconName(this._iconNameStack.pop());
        break;
    }
    return true;
  }
  saveCursor(e) {
    return this._activeBuffer.savedX = this._activeBuffer.x, this._activeBuffer.savedY = this._activeBuffer.ybase + this._activeBuffer.y, this._activeBuffer.savedCurAttrData.fg = this._curAttrData.fg, this._activeBuffer.savedCurAttrData.bg = this._curAttrData.bg, this._activeBuffer.savedCharset = this._charsetService.charset, true;
  }
  restoreCursor(e) {
    return this._activeBuffer.x = this._activeBuffer.savedX || 0, this._activeBuffer.y = Math.max(this._activeBuffer.savedY - this._activeBuffer.ybase, 0), this._curAttrData.fg = this._activeBuffer.savedCurAttrData.fg, this._curAttrData.bg = this._activeBuffer.savedCurAttrData.bg, this._charsetService.charset = this._savedCharset, this._activeBuffer.savedCharset && (this._charsetService.charset = this._activeBuffer.savedCharset), this._restrictCursor(), true;
  }
  setTitle(e) {
    return this._windowTitle = e, this._onTitleChange.fire(e), true;
  }
  setIconName(e) {
    return this._iconName = e, true;
  }
  setOrReportIndexedColor(e) {
    let i8 = [], r5 = e.split(";");
    for (; r5.length > 1; ) {
      let n = r5.shift(), o2 = r5.shift();
      if (/^\d+$/.exec(n)) {
        let l = parseInt(n);
        if (Sl(l)) if (o2 === "?") i8.push({ type: 0, index: l });
        else {
          let a = Ws2(o2);
          a && i8.push({ type: 1, index: l, color: a });
        }
      }
    }
    return i8.length && this._onColor.fire(i8), true;
  }
  setHyperlink(e) {
    let i8 = e.indexOf(";");
    if (i8 === -1) return true;
    let r5 = e.slice(0, i8).trim(), n = e.slice(i8 + 1);
    return n ? this._createHyperlink(r5, n) : r5.trim() ? false : this._finishHyperlink();
  }
  _createHyperlink(e, i8) {
    this._getCurrentLinkId() && this._finishHyperlink();
    let r5 = e.split(":"), n, o2 = r5.findIndex((l) => l.startsWith("id="));
    return o2 !== -1 && (n = r5[o2].slice(3) || void 0), this._curAttrData.extended = this._curAttrData.extended.clone(), this._curAttrData.extended.urlId = this._oscLinkService.registerLink({ id: n, uri: i8 }), this._curAttrData.updateExtended(), true;
  }
  _finishHyperlink() {
    return this._curAttrData.extended = this._curAttrData.extended.clone(), this._curAttrData.extended.urlId = 0, this._curAttrData.updateExtended(), true;
  }
  _setOrReportSpecialColor(e, i8) {
    let r5 = e.split(";");
    for (let n = 0; n < r5.length && !(i8 >= this._specialColors.length); ++n, ++i8) if (r5[n] === "?") this._onColor.fire([{ type: 0, index: this._specialColors[i8] }]);
    else {
      let o2 = Ws2(r5[n]);
      o2 && this._onColor.fire([{ type: 1, index: this._specialColors[i8], color: o2 }]);
    }
    return true;
  }
  setOrReportFgColor(e) {
    return this._setOrReportSpecialColor(e, 0);
  }
  setOrReportBgColor(e) {
    return this._setOrReportSpecialColor(e, 1);
  }
  setOrReportCursorColor(e) {
    return this._setOrReportSpecialColor(e, 2);
  }
  restoreIndexedColor(e) {
    if (!e) return this._onColor.fire([{ type: 2 }]), true;
    let i8 = [], r5 = e.split(";");
    for (let n = 0; n < r5.length; ++n) if (/^\d+$/.exec(r5[n])) {
      let o2 = parseInt(r5[n]);
      Sl(o2) && i8.push({ type: 2, index: o2 });
    }
    return i8.length && this._onColor.fire(i8), true;
  }
  restoreFgColor(e) {
    return this._onColor.fire([{ type: 2, index: 256 }]), true;
  }
  restoreBgColor(e) {
    return this._onColor.fire([{ type: 2, index: 257 }]), true;
  }
  restoreCursorColor(e) {
    return this._onColor.fire([{ type: 2, index: 258 }]), true;
  }
  nextLine() {
    return this._activeBuffer.x = 0, this.index(), true;
  }
  keypadApplicationMode() {
    return this._logService.debug("Serial port requested application keypad."), this._coreService.decPrivateModes.applicationKeypad = true, this._onRequestSyncScrollBar.fire(), true;
  }
  keypadNumericMode() {
    return this._logService.debug("Switching back to normal keypad."), this._coreService.decPrivateModes.applicationKeypad = false, this._onRequestSyncScrollBar.fire(), true;
  }
  selectDefaultCharset() {
    return this._charsetService.setgLevel(0), this._charsetService.setgCharset(0, Je2), true;
  }
  selectCharset(e) {
    return e.length !== 2 ? (this.selectDefaultCharset(), true) : (e[0] === "/" || this._charsetService.setgCharset(mc[e[0]], ne2[e[1]] || Je2), true);
  }
  index() {
    return this._restrictCursor(), this._activeBuffer.y++, this._activeBuffer.y === this._activeBuffer.scrollBottom + 1 ? (this._activeBuffer.y--, this._bufferService.scroll(this._eraseAttrData())) : this._activeBuffer.y >= this._bufferService.rows && (this._activeBuffer.y = this._bufferService.rows - 1), this._restrictCursor(), true;
  }
  tabSet() {
    return this._activeBuffer.tabs[this._activeBuffer.x] = true, true;
  }
  reverseIndex() {
    if (this._restrictCursor(), this._activeBuffer.y === this._activeBuffer.scrollTop) {
      let e = this._activeBuffer.scrollBottom - this._activeBuffer.scrollTop;
      this._activeBuffer.lines.shiftElements(this._activeBuffer.ybase + this._activeBuffer.y, e, 1), this._activeBuffer.lines.set(this._activeBuffer.ybase + this._activeBuffer.y, this._activeBuffer.getBlankLine(this._eraseAttrData())), this._dirtyRowTracker.markRangeDirty(this._activeBuffer.scrollTop, this._activeBuffer.scrollBottom);
    } else this._activeBuffer.y--, this._restrictCursor();
    return true;
  }
  fullReset() {
    return this._parser.reset(), this._onRequestReset.fire(), true;
  }
  reset() {
    this._curAttrData = X4.clone(), this._eraseAttrDataInternal = X4.clone();
  }
  _eraseAttrData() {
    return this._eraseAttrDataInternal.bg &= -67108864, this._eraseAttrDataInternal.bg |= this._curAttrData.bg & 67108863, this._eraseAttrDataInternal;
  }
  setgLevel(e) {
    return this._charsetService.setgLevel(e), true;
  }
  screenAlignmentPattern() {
    let e = new q();
    e.content = 1 << 22 | 69, e.fg = this._curAttrData.fg, e.bg = this._curAttrData.bg, this._setCursor(0, 0);
    for (let i8 = 0; i8 < this._bufferService.rows; ++i8) {
      let r5 = this._activeBuffer.ybase + this._activeBuffer.y + i8, n = this._activeBuffer.lines.get(r5);
      n && (n.fill(e), n.isWrapped = false);
    }
    return this._dirtyRowTracker.markAllDirty(), this._setCursor(0, 0), true;
  }
  requestStatusString(e, i8) {
    let r5 = (a) => (this._coreService.triggerDataEvent(`${b.ESC}${a}${b.ESC}\\`), true), n = this._bufferService.buffer, o2 = this._optionsService.rawOptions, l = { block: 2, underline: 4, bar: 6 };
    return r5(e === '"q' ? `P1$r${this._curAttrData.isProtected() ? 1 : 0}"q` : e === '"p' ? 'P1$r61;1"p' : e === "r" ? `P1$r${n.scrollTop + 1};${n.scrollBottom + 1}r` : e === "m" ? "P1$r0m" : e === " q" ? `P1$r${l[o2.cursorStyle] - (o2.cursorBlink ? 1 : 0)} q` : "P0$r");
  }
  markRangeDirty(e, i8) {
    this._dirtyRowTracker.markRangeDirty(e, i8);
  }
};
var Zi = class {
  constructor(t) {
    this._bufferService = t;
    this.clearRange();
  }
  clearRange() {
    this.start = this._bufferService.buffer.y, this.end = this._bufferService.buffer.y;
  }
  markDirty(t) {
    t < this.start ? this.start = t : t > this.end && (this.end = t);
  }
  markRangeDirty(t, e) {
    t > e && (gl = t, t = e, e = gl), t < this.start && (this.start = t), e > this.end && (this.end = e);
  }
  markAllDirty() {
    this.markRangeDirty(0, this._bufferService.rows - 1);
  }
};
Zi = M3([S(0, F2)], Zi);
function Sl(s15) {
  return 0 <= s15 && s15 < 256;
}
var _c = 5e7;
var El = 12;
var bc = 50;
var gn2 = class extends D2 {
  constructor(e) {
    super();
    this._action = e;
    this._writeBuffer = [];
    this._callbacks = [];
    this._pendingData = 0;
    this._bufferOffset = 0;
    this._isSyncWriting = false;
    this._syncCalls = 0;
    this._didUserInput = false;
    this._onWriteParsed = this._register(new v());
    this.onWriteParsed = this._onWriteParsed.event;
  }
  handleUserInput() {
    this._didUserInput = true;
  }
  writeSync(e, i8) {
    if (i8 !== void 0 && this._syncCalls > i8) {
      this._syncCalls = 0;
      return;
    }
    if (this._pendingData += e.length, this._writeBuffer.push(e), this._callbacks.push(void 0), this._syncCalls++, this._isSyncWriting) return;
    this._isSyncWriting = true;
    let r5;
    for (; r5 = this._writeBuffer.shift(); ) {
      this._action(r5);
      let n = this._callbacks.shift();
      n && n();
    }
    this._pendingData = 0, this._bufferOffset = 2147483647, this._isSyncWriting = false, this._syncCalls = 0;
  }
  write(e, i8) {
    if (this._pendingData > _c) throw new Error("write data discarded, use flow control to avoid losing data");
    if (!this._writeBuffer.length) {
      if (this._bufferOffset = 0, this._didUserInput) {
        this._didUserInput = false, this._pendingData += e.length, this._writeBuffer.push(e), this._callbacks.push(i8), this._innerWrite();
        return;
      }
      setTimeout(() => this._innerWrite());
    }
    this._pendingData += e.length, this._writeBuffer.push(e), this._callbacks.push(i8);
  }
  _innerWrite(e = 0, i8 = true) {
    let r5 = e || performance.now();
    for (; this._writeBuffer.length > this._bufferOffset; ) {
      let n = this._writeBuffer[this._bufferOffset], o2 = this._action(n, i8);
      if (o2) {
        let a = (u) => performance.now() - r5 >= El ? setTimeout(() => this._innerWrite(0, u)) : this._innerWrite(r5, u);
        o2.catch((u) => (queueMicrotask(() => {
          throw u;
        }), Promise.resolve(false))).then(a);
        return;
      }
      let l = this._callbacks[this._bufferOffset];
      if (l && l(), this._bufferOffset++, this._pendingData -= n.length, performance.now() - r5 >= El) break;
    }
    this._writeBuffer.length > this._bufferOffset ? (this._bufferOffset > bc && (this._writeBuffer = this._writeBuffer.slice(this._bufferOffset), this._callbacks = this._callbacks.slice(this._bufferOffset), this._bufferOffset = 0), setTimeout(() => this._innerWrite())) : (this._writeBuffer.length = 0, this._callbacks.length = 0, this._pendingData = 0, this._bufferOffset = 0), this._onWriteParsed.fire();
  }
};
var ui2 = class {
  constructor(t) {
    this._bufferService = t;
    this._nextId = 1;
    this._entriesWithId = /* @__PURE__ */ new Map();
    this._dataByLinkId = /* @__PURE__ */ new Map();
  }
  registerLink(t) {
    let e = this._bufferService.buffer;
    if (t.id === void 0) {
      let a = e.addMarker(e.ybase + e.y), u = { data: t, id: this._nextId++, lines: [a] };
      return a.onDispose(() => this._removeMarkerFromLink(u, a)), this._dataByLinkId.set(u.id, u), u.id;
    }
    let i8 = t, r5 = this._getEntryIdKey(i8), n = this._entriesWithId.get(r5);
    if (n) return this.addLineToLink(n.id, e.ybase + e.y), n.id;
    let o2 = e.addMarker(e.ybase + e.y), l = { id: this._nextId++, key: this._getEntryIdKey(i8), data: i8, lines: [o2] };
    return o2.onDispose(() => this._removeMarkerFromLink(l, o2)), this._entriesWithId.set(l.key, l), this._dataByLinkId.set(l.id, l), l.id;
  }
  addLineToLink(t, e) {
    let i8 = this._dataByLinkId.get(t);
    if (i8 && i8.lines.every((r5) => r5.line !== e)) {
      let r5 = this._bufferService.buffer.addMarker(e);
      i8.lines.push(r5), r5.onDispose(() => this._removeMarkerFromLink(i8, r5));
    }
  }
  getLinkData(t) {
    return this._dataByLinkId.get(t)?.data;
  }
  _getEntryIdKey(t) {
    return `${t.id};;${t.uri}`;
  }
  _removeMarkerFromLink(t, e) {
    let i8 = t.lines.indexOf(e);
    i8 !== -1 && (t.lines.splice(i8, 1), t.lines.length === 0 && (t.data.id !== void 0 && this._entriesWithId.delete(t.key), this._dataByLinkId.delete(t.id)));
  }
};
ui2 = M3([S(0, F2)], ui2);
var Tl = false;
var Sn2 = class extends D2 {
  constructor(e) {
    super();
    this._windowsWrappingHeuristics = this._register(new ye3());
    this._onBinary = this._register(new v());
    this.onBinary = this._onBinary.event;
    this._onData = this._register(new v());
    this.onData = this._onData.event;
    this._onLineFeed = this._register(new v());
    this.onLineFeed = this._onLineFeed.event;
    this._onResize = this._register(new v());
    this.onResize = this._onResize.event;
    this._onWriteParsed = this._register(new v());
    this.onWriteParsed = this._onWriteParsed.event;
    this._onScroll = this._register(new v());
    this._instantiationService = new ln(), this.optionsService = this._register(new dn2(e)), this._instantiationService.setService(H3, this.optionsService), this._bufferService = this._register(this._instantiationService.createInstance(ni2)), this._instantiationService.setService(F2, this._bufferService), this._logService = this._register(this._instantiationService.createInstance(ii2)), this._instantiationService.setService(nr2, this._logService), this.coreService = this._register(this._instantiationService.createInstance(li2)), this._instantiationService.setService(ge3, this.coreService), this.coreMouseService = this._register(this._instantiationService.createInstance(ai3)), this._instantiationService.setService(rr, this.coreMouseService), this.unicodeService = this._register(this._instantiationService.createInstance(Ae3)), this._instantiationService.setService(Js, this.unicodeService), this._charsetService = this._instantiationService.createInstance(pn2), this._instantiationService.setService(Zs, this._charsetService), this._oscLinkService = this._instantiationService.createInstance(ui2), this._instantiationService.setService(sr2, this._oscLinkService), this._inputHandler = this._register(new vn(this._bufferService, this._charsetService, this.coreService, this._logService, this.optionsService, this._oscLinkService, this.coreMouseService, this.unicodeService)), this._register($.forward(this._inputHandler.onLineFeed, this._onLineFeed)), this._register(this._inputHandler), this._register($.forward(this._bufferService.onResize, this._onResize)), this._register($.forward(this.coreService.onData, this._onData)), this._register($.forward(this.coreService.onBinary, this._onBinary)), this._register(this.coreService.onRequestScrollToBottom(() => this.scrollToBottom(true))), this._register(this.coreService.onUserInput(() => this._writeBuffer.handleUserInput())), this._register(this.optionsService.onMultipleOptionChange(["windowsMode", "windowsPty"], () => this._handleWindowsPtyOptionChange())), this._register(this._bufferService.onScroll(() => {
      this._onScroll.fire({ position: this._bufferService.buffer.ydisp }), this._inputHandler.markRangeDirty(this._bufferService.buffer.scrollTop, this._bufferService.buffer.scrollBottom);
    })), this._writeBuffer = this._register(new gn2((i8, r5) => this._inputHandler.parse(i8, r5))), this._register($.forward(this._writeBuffer.onWriteParsed, this._onWriteParsed));
  }
  get onScroll() {
    return this._onScrollApi || (this._onScrollApi = this._register(new v()), this._onScroll.event((e) => {
      this._onScrollApi?.fire(e.position);
    })), this._onScrollApi.event;
  }
  get cols() {
    return this._bufferService.cols;
  }
  get rows() {
    return this._bufferService.rows;
  }
  get buffers() {
    return this._bufferService.buffers;
  }
  get options() {
    return this.optionsService.options;
  }
  set options(e) {
    for (let i8 in e) this.optionsService.options[i8] = e[i8];
  }
  write(e, i8) {
    this._writeBuffer.write(e, i8);
  }
  writeSync(e, i8) {
    this._logService.logLevel <= 3 && !Tl && (this._logService.warn("writeSync is unreliable and will be removed soon."), Tl = true), this._writeBuffer.writeSync(e, i8);
  }
  input(e, i8 = true) {
    this.coreService.triggerDataEvent(e, i8);
  }
  resize(e, i8) {
    isNaN(e) || isNaN(i8) || (e = Math.max(e, ks2), i8 = Math.max(i8, Cs), this._bufferService.resize(e, i8));
  }
  scroll(e, i8 = false) {
    this._bufferService.scroll(e, i8);
  }
  scrollLines(e, i8) {
    this._bufferService.scrollLines(e, i8);
  }
  scrollPages(e) {
    this.scrollLines(e * (this.rows - 1));
  }
  scrollToTop() {
    this.scrollLines(-this._bufferService.buffer.ydisp);
  }
  scrollToBottom(e) {
    this.scrollLines(this._bufferService.buffer.ybase - this._bufferService.buffer.ydisp);
  }
  scrollToLine(e) {
    let i8 = e - this._bufferService.buffer.ydisp;
    i8 !== 0 && this.scrollLines(i8);
  }
  registerEscHandler(e, i8) {
    return this._inputHandler.registerEscHandler(e, i8);
  }
  registerDcsHandler(e, i8) {
    return this._inputHandler.registerDcsHandler(e, i8);
  }
  registerCsiHandler(e, i8) {
    return this._inputHandler.registerCsiHandler(e, i8);
  }
  registerOscHandler(e, i8) {
    return this._inputHandler.registerOscHandler(e, i8);
  }
  _setup() {
    this._handleWindowsPtyOptionChange();
  }
  reset() {
    this._inputHandler.reset(), this._bufferService.reset(), this._charsetService.reset(), this.coreService.reset(), this.coreMouseService.reset();
  }
  _handleWindowsPtyOptionChange() {
    let e = false, i8 = this.optionsService.rawOptions.windowsPty;
    i8 && i8.buildNumber !== void 0 && i8.buildNumber !== void 0 ? e = i8.backend === "conpty" && i8.buildNumber < 21376 : this.optionsService.rawOptions.windowsMode && (e = true), e ? this._enableWindowsWrappingHeuristics() : this._windowsWrappingHeuristics.clear();
  }
  _enableWindowsWrappingHeuristics() {
    if (!this._windowsWrappingHeuristics.value) {
      let e = [];
      e.push(this.onLineFeed(Bs2.bind(null, this._bufferService))), e.push(this.registerCsiHandler({ final: "H" }, () => (Bs2(this._bufferService), false))), this._windowsWrappingHeuristics.value = C3(() => {
        for (let i8 of e) i8.dispose();
      });
    }
  }
};
var gc = { 48: ["0", ")"], 49: ["1", "!"], 50: ["2", "@"], 51: ["3", "#"], 52: ["4", "$"], 53: ["5", "%"], 54: ["6", "^"], 55: ["7", "&"], 56: ["8", "*"], 57: ["9", "("], 186: [";", ":"], 187: ["=", "+"], 188: [",", "<"], 189: ["-", "_"], 190: [".", ">"], 191: ["/", "?"], 192: ["`", "~"], 219: ["[", "{"], 220: ["\\", "|"], 221: ["]", "}"], 222: ["'", '"'] };
function Il(s15, t, e, i8) {
  let r5 = { type: 0, cancel: false, key: void 0 }, n = (s15.shiftKey ? 1 : 0) | (s15.altKey ? 2 : 0) | (s15.ctrlKey ? 4 : 0) | (s15.metaKey ? 8 : 0);
  switch (s15.keyCode) {
    case 0:
      s15.key === "UIKeyInputUpArrow" ? t ? r5.key = b.ESC + "OA" : r5.key = b.ESC + "[A" : s15.key === "UIKeyInputLeftArrow" ? t ? r5.key = b.ESC + "OD" : r5.key = b.ESC + "[D" : s15.key === "UIKeyInputRightArrow" ? t ? r5.key = b.ESC + "OC" : r5.key = b.ESC + "[C" : s15.key === "UIKeyInputDownArrow" && (t ? r5.key = b.ESC + "OB" : r5.key = b.ESC + "[B");
      break;
    case 8:
      r5.key = s15.ctrlKey ? "\b" : b.DEL, s15.altKey && (r5.key = b.ESC + r5.key);
      break;
    case 9:
      if (s15.shiftKey) {
        r5.key = b.ESC + "[Z";
        break;
      }
      r5.key = b.HT, r5.cancel = true;
      break;
    case 13:
      r5.key = s15.altKey ? b.ESC + b.CR : b.CR, r5.cancel = true;
      break;
    case 27:
      r5.key = b.ESC, s15.altKey && (r5.key = b.ESC + b.ESC), r5.cancel = true;
      break;
    case 37:
      if (s15.metaKey) break;
      n ? r5.key = b.ESC + "[1;" + (n + 1) + "D" : t ? r5.key = b.ESC + "OD" : r5.key = b.ESC + "[D";
      break;
    case 39:
      if (s15.metaKey) break;
      n ? r5.key = b.ESC + "[1;" + (n + 1) + "C" : t ? r5.key = b.ESC + "OC" : r5.key = b.ESC + "[C";
      break;
    case 38:
      if (s15.metaKey) break;
      n ? r5.key = b.ESC + "[1;" + (n + 1) + "A" : t ? r5.key = b.ESC + "OA" : r5.key = b.ESC + "[A";
      break;
    case 40:
      if (s15.metaKey) break;
      n ? r5.key = b.ESC + "[1;" + (n + 1) + "B" : t ? r5.key = b.ESC + "OB" : r5.key = b.ESC + "[B";
      break;
    case 45:
      !s15.shiftKey && !s15.ctrlKey && (r5.key = b.ESC + "[2~");
      break;
    case 46:
      n ? r5.key = b.ESC + "[3;" + (n + 1) + "~" : r5.key = b.ESC + "[3~";
      break;
    case 36:
      n ? r5.key = b.ESC + "[1;" + (n + 1) + "H" : t ? r5.key = b.ESC + "OH" : r5.key = b.ESC + "[H";
      break;
    case 35:
      n ? r5.key = b.ESC + "[1;" + (n + 1) + "F" : t ? r5.key = b.ESC + "OF" : r5.key = b.ESC + "[F";
      break;
    case 33:
      s15.shiftKey ? r5.type = 2 : s15.ctrlKey ? r5.key = b.ESC + "[5;" + (n + 1) + "~" : r5.key = b.ESC + "[5~";
      break;
    case 34:
      s15.shiftKey ? r5.type = 3 : s15.ctrlKey ? r5.key = b.ESC + "[6;" + (n + 1) + "~" : r5.key = b.ESC + "[6~";
      break;
    case 112:
      n ? r5.key = b.ESC + "[1;" + (n + 1) + "P" : r5.key = b.ESC + "OP";
      break;
    case 113:
      n ? r5.key = b.ESC + "[1;" + (n + 1) + "Q" : r5.key = b.ESC + "OQ";
      break;
    case 114:
      n ? r5.key = b.ESC + "[1;" + (n + 1) + "R" : r5.key = b.ESC + "OR";
      break;
    case 115:
      n ? r5.key = b.ESC + "[1;" + (n + 1) + "S" : r5.key = b.ESC + "OS";
      break;
    case 116:
      n ? r5.key = b.ESC + "[15;" + (n + 1) + "~" : r5.key = b.ESC + "[15~";
      break;
    case 117:
      n ? r5.key = b.ESC + "[17;" + (n + 1) + "~" : r5.key = b.ESC + "[17~";
      break;
    case 118:
      n ? r5.key = b.ESC + "[18;" + (n + 1) + "~" : r5.key = b.ESC + "[18~";
      break;
    case 119:
      n ? r5.key = b.ESC + "[19;" + (n + 1) + "~" : r5.key = b.ESC + "[19~";
      break;
    case 120:
      n ? r5.key = b.ESC + "[20;" + (n + 1) + "~" : r5.key = b.ESC + "[20~";
      break;
    case 121:
      n ? r5.key = b.ESC + "[21;" + (n + 1) + "~" : r5.key = b.ESC + "[21~";
      break;
    case 122:
      n ? r5.key = b.ESC + "[23;" + (n + 1) + "~" : r5.key = b.ESC + "[23~";
      break;
    case 123:
      n ? r5.key = b.ESC + "[24;" + (n + 1) + "~" : r5.key = b.ESC + "[24~";
      break;
    default:
      if (s15.ctrlKey && !s15.shiftKey && !s15.altKey && !s15.metaKey) s15.keyCode >= 65 && s15.keyCode <= 90 ? r5.key = String.fromCharCode(s15.keyCode - 64) : s15.keyCode === 32 ? r5.key = b.NUL : s15.keyCode >= 51 && s15.keyCode <= 55 ? r5.key = String.fromCharCode(s15.keyCode - 51 + 27) : s15.keyCode === 56 ? r5.key = b.DEL : s15.keyCode === 219 ? r5.key = b.ESC : s15.keyCode === 220 ? r5.key = b.FS : s15.keyCode === 221 && (r5.key = b.GS);
      else if ((!e || i8) && s15.altKey && !s15.metaKey) {
        let l = gc[s15.keyCode]?.[s15.shiftKey ? 1 : 0];
        if (l) r5.key = b.ESC + l;
        else if (s15.keyCode >= 65 && s15.keyCode <= 90) {
          let a = s15.ctrlKey ? s15.keyCode - 64 : s15.keyCode + 32, u = String.fromCharCode(a);
          s15.shiftKey && (u = u.toUpperCase()), r5.key = b.ESC + u;
        } else if (s15.keyCode === 32) r5.key = b.ESC + (s15.ctrlKey ? b.NUL : " ");
        else if (s15.key === "Dead" && s15.code.startsWith("Key")) {
          let a = s15.code.slice(3, 4);
          s15.shiftKey || (a = a.toLowerCase()), r5.key = b.ESC + a, r5.cancel = true;
        }
      } else e && !s15.altKey && !s15.ctrlKey && !s15.shiftKey && s15.metaKey ? s15.keyCode === 65 && (r5.type = 1) : s15.key && !s15.ctrlKey && !s15.altKey && !s15.metaKey && s15.keyCode >= 48 && s15.key.length === 1 ? r5.key = s15.key : s15.key && s15.ctrlKey && (s15.key === "_" && (r5.key = b.US), s15.key === "@" && (r5.key = b.NUL));
      break;
  }
  return r5;
}
var ee4 = 0;
var En2 = class {
  constructor(t) {
    this._getKey = t;
    this._array = [];
    this._insertedValues = [];
    this._flushInsertedTask = new Jt();
    this._isFlushingInserted = false;
    this._deletedIndices = [];
    this._flushDeletedTask = new Jt();
    this._isFlushingDeleted = false;
  }
  clear() {
    this._array.length = 0, this._insertedValues.length = 0, this._flushInsertedTask.clear(), this._isFlushingInserted = false, this._deletedIndices.length = 0, this._flushDeletedTask.clear(), this._isFlushingDeleted = false;
  }
  insert(t) {
    this._flushCleanupDeleted(), this._insertedValues.length === 0 && this._flushInsertedTask.enqueue(() => this._flushInserted()), this._insertedValues.push(t);
  }
  _flushInserted() {
    let t = this._insertedValues.sort((n, o2) => this._getKey(n) - this._getKey(o2)), e = 0, i8 = 0, r5 = new Array(this._array.length + this._insertedValues.length);
    for (let n = 0; n < r5.length; n++) i8 >= this._array.length || this._getKey(t[e]) <= this._getKey(this._array[i8]) ? (r5[n] = t[e], e++) : r5[n] = this._array[i8++];
    this._array = r5, this._insertedValues.length = 0;
  }
  _flushCleanupInserted() {
    !this._isFlushingInserted && this._insertedValues.length > 0 && this._flushInsertedTask.flush();
  }
  delete(t) {
    if (this._flushCleanupInserted(), this._array.length === 0) return false;
    let e = this._getKey(t);
    if (e === void 0 || (ee4 = this._search(e), ee4 === -1) || this._getKey(this._array[ee4]) !== e) return false;
    do
      if (this._array[ee4] === t) return this._deletedIndices.length === 0 && this._flushDeletedTask.enqueue(() => this._flushDeleted()), this._deletedIndices.push(ee4), true;
    while (++ee4 < this._array.length && this._getKey(this._array[ee4]) === e);
    return false;
  }
  _flushDeleted() {
    this._isFlushingDeleted = true;
    let t = this._deletedIndices.sort((n, o2) => n - o2), e = 0, i8 = new Array(this._array.length - t.length), r5 = 0;
    for (let n = 0; n < this._array.length; n++) t[e] === n ? e++ : i8[r5++] = this._array[n];
    this._array = i8, this._deletedIndices.length = 0, this._isFlushingDeleted = false;
  }
  _flushCleanupDeleted() {
    !this._isFlushingDeleted && this._deletedIndices.length > 0 && this._flushDeletedTask.flush();
  }
  *getKeyIterator(t) {
    if (this._flushCleanupInserted(), this._flushCleanupDeleted(), this._array.length !== 0 && (ee4 = this._search(t), !(ee4 < 0 || ee4 >= this._array.length) && this._getKey(this._array[ee4]) === t)) do
      yield this._array[ee4];
    while (++ee4 < this._array.length && this._getKey(this._array[ee4]) === t);
  }
  forEachByKey(t, e) {
    if (this._flushCleanupInserted(), this._flushCleanupDeleted(), this._array.length !== 0 && (ee4 = this._search(t), !(ee4 < 0 || ee4 >= this._array.length) && this._getKey(this._array[ee4]) === t)) do
      e(this._array[ee4]);
    while (++ee4 < this._array.length && this._getKey(this._array[ee4]) === t);
  }
  values() {
    return this._flushCleanupInserted(), this._flushCleanupDeleted(), [...this._array].values();
  }
  _search(t) {
    let e = 0, i8 = this._array.length - 1;
    for (; i8 >= e; ) {
      let r5 = e + i8 >> 1, n = this._getKey(this._array[r5]);
      if (n > t) i8 = r5 - 1;
      else if (n < t) e = r5 + 1;
      else {
        for (; r5 > 0 && this._getKey(this._array[r5 - 1]) === t; ) r5--;
        return r5;
      }
    }
    return e;
  }
};
var Us2 = 0;
var yl = 0;
var Tn = class extends D2 {
  constructor() {
    super();
    this._decorations = new En2((e) => e?.marker.line);
    this._onDecorationRegistered = this._register(new v());
    this.onDecorationRegistered = this._onDecorationRegistered.event;
    this._onDecorationRemoved = this._register(new v());
    this.onDecorationRemoved = this._onDecorationRemoved.event;
    this._register(C3(() => this.reset()));
  }
  get decorations() {
    return this._decorations.values();
  }
  registerDecoration(e) {
    if (e.marker.isDisposed) return;
    let i8 = new Ks(e);
    if (i8) {
      let r5 = i8.marker.onDispose(() => i8.dispose()), n = i8.onDispose(() => {
        n.dispose(), i8 && (this._decorations.delete(i8) && this._onDecorationRemoved.fire(i8), r5.dispose());
      });
      this._decorations.insert(i8), this._onDecorationRegistered.fire(i8);
    }
    return i8;
  }
  reset() {
    for (let e of this._decorations.values()) e.dispose();
    this._decorations.clear();
  }
  *getDecorationsAtCell(e, i8, r5) {
    let n = 0, o2 = 0;
    for (let l of this._decorations.getKeyIterator(i8)) n = l.options.x ?? 0, o2 = n + (l.options.width ?? 1), e >= n && e < o2 && (!r5 || (l.options.layer ?? "bottom") === r5) && (yield l);
  }
  forEachDecorationAtCell(e, i8, r5, n) {
    this._decorations.forEachByKey(i8, (o2) => {
      Us2 = o2.options.x ?? 0, yl = Us2 + (o2.options.width ?? 1), e >= Us2 && e < yl && (!r5 || (o2.options.layer ?? "bottom") === r5) && n(o2);
    });
  }
};
var Ks = class extends Ee3 {
  constructor(e) {
    super();
    this.options = e;
    this.onRenderEmitter = this.add(new v());
    this.onRender = this.onRenderEmitter.event;
    this._onDispose = this.add(new v());
    this.onDispose = this._onDispose.event;
    this._cachedBg = null;
    this._cachedFg = null;
    this.marker = e.marker, this.options.overviewRulerOptions && !this.options.overviewRulerOptions.position && (this.options.overviewRulerOptions.position = "full");
  }
  get backgroundColorRGB() {
    return this._cachedBg === null && (this.options.backgroundColor ? this._cachedBg = z2.toColor(this.options.backgroundColor) : this._cachedBg = void 0), this._cachedBg;
  }
  get foregroundColorRGB() {
    return this._cachedFg === null && (this.options.foregroundColor ? this._cachedFg = z2.toColor(this.options.foregroundColor) : this._cachedFg = void 0), this._cachedFg;
  }
  dispose() {
    this._onDispose.fire(), super.dispose();
  }
};
var Sc = 1e3;
var In2 = class {
  constructor(t, e = Sc) {
    this._renderCallback = t;
    this._debounceThresholdMS = e;
    this._lastRefreshMs = 0;
    this._additionalRefreshRequested = false;
  }
  dispose() {
    this._refreshTimeoutID && clearTimeout(this._refreshTimeoutID);
  }
  refresh(t, e, i8) {
    this._rowCount = i8, t = t !== void 0 ? t : 0, e = e !== void 0 ? e : this._rowCount - 1, this._rowStart = this._rowStart !== void 0 ? Math.min(this._rowStart, t) : t, this._rowEnd = this._rowEnd !== void 0 ? Math.max(this._rowEnd, e) : e;
    let r5 = performance.now();
    if (r5 - this._lastRefreshMs >= this._debounceThresholdMS) this._lastRefreshMs = r5, this._innerRefresh();
    else if (!this._additionalRefreshRequested) {
      let n = r5 - this._lastRefreshMs, o2 = this._debounceThresholdMS - n;
      this._additionalRefreshRequested = true, this._refreshTimeoutID = window.setTimeout(() => {
        this._lastRefreshMs = performance.now(), this._innerRefresh(), this._additionalRefreshRequested = false, this._refreshTimeoutID = void 0;
      }, o2);
    }
  }
  _innerRefresh() {
    if (this._rowStart === void 0 || this._rowEnd === void 0 || this._rowCount === void 0) return;
    let t = Math.max(this._rowStart, 0), e = Math.min(this._rowEnd, this._rowCount - 1);
    this._rowStart = void 0, this._rowEnd = void 0, this._renderCallback(t, e);
  }
};
var xl = 20;
var wl = false;
var Tt2 = class extends D2 {
  constructor(e, i8, r5, n) {
    super();
    this._terminal = e;
    this._coreBrowserService = r5;
    this._renderService = n;
    this._rowColumns = /* @__PURE__ */ new WeakMap();
    this._liveRegionLineCount = 0;
    this._charsToConsume = [];
    this._charsToAnnounce = "";
    let o2 = this._coreBrowserService.mainDocument;
    this._accessibilityContainer = o2.createElement("div"), this._accessibilityContainer.classList.add("xterm-accessibility"), this._rowContainer = o2.createElement("div"), this._rowContainer.setAttribute("role", "list"), this._rowContainer.classList.add("xterm-accessibility-tree"), this._rowElements = [];
    for (let l = 0; l < this._terminal.rows; l++) this._rowElements[l] = this._createAccessibilityTreeNode(), this._rowContainer.appendChild(this._rowElements[l]);
    if (this._topBoundaryFocusListener = (l) => this._handleBoundaryFocus(l, 0), this._bottomBoundaryFocusListener = (l) => this._handleBoundaryFocus(l, 1), this._rowElements[0].addEventListener("focus", this._topBoundaryFocusListener), this._rowElements[this._rowElements.length - 1].addEventListener("focus", this._bottomBoundaryFocusListener), this._accessibilityContainer.appendChild(this._rowContainer), this._liveRegion = o2.createElement("div"), this._liveRegion.classList.add("live-region"), this._liveRegion.setAttribute("aria-live", "assertive"), this._accessibilityContainer.appendChild(this._liveRegion), this._liveRegionDebouncer = this._register(new In2(this._renderRows.bind(this))), !this._terminal.element) throw new Error("Cannot enable accessibility before Terminal.open");
    wl ? (this._accessibilityContainer.classList.add("debug"), this._rowContainer.classList.add("debug"), this._debugRootContainer = o2.createElement("div"), this._debugRootContainer.classList.add("xterm"), this._debugRootContainer.appendChild(o2.createTextNode("------start a11y------")), this._debugRootContainer.appendChild(this._accessibilityContainer), this._debugRootContainer.appendChild(o2.createTextNode("------end a11y------")), this._terminal.element.insertAdjacentElement("afterend", this._debugRootContainer)) : this._terminal.element.insertAdjacentElement("afterbegin", this._accessibilityContainer), this._register(this._terminal.onResize((l) => this._handleResize(l.rows))), this._register(this._terminal.onRender((l) => this._refreshRows(l.start, l.end))), this._register(this._terminal.onScroll(() => this._refreshRows())), this._register(this._terminal.onA11yChar((l) => this._handleChar(l))), this._register(this._terminal.onLineFeed(() => this._handleChar(`
`))), this._register(this._terminal.onA11yTab((l) => this._handleTab(l))), this._register(this._terminal.onKey((l) => this._handleKey(l.key))), this._register(this._terminal.onBlur(() => this._clearLiveRegion())), this._register(this._renderService.onDimensionsChange(() => this._refreshRowsDimensions())), this._register(L(o2, "selectionchange", () => this._handleSelectionChange())), this._register(this._coreBrowserService.onDprChange(() => this._refreshRowsDimensions())), this._refreshRowsDimensions(), this._refreshRows(), this._register(C3(() => {
      wl ? this._debugRootContainer.remove() : this._accessibilityContainer.remove(), this._rowElements.length = 0;
    }));
  }
  _handleTab(e) {
    for (let i8 = 0; i8 < e; i8++) this._handleChar(" ");
  }
  _handleChar(e) {
    this._liveRegionLineCount < xl + 1 && (this._charsToConsume.length > 0 ? this._charsToConsume.shift() !== e && (this._charsToAnnounce += e) : this._charsToAnnounce += e, e === `
` && (this._liveRegionLineCount++, this._liveRegionLineCount === xl + 1 && (this._liveRegion.textContent += _i2.get())));
  }
  _clearLiveRegion() {
    this._liveRegion.textContent = "", this._liveRegionLineCount = 0;
  }
  _handleKey(e) {
    this._clearLiveRegion(), /\p{Control}/u.test(e) || this._charsToConsume.push(e);
  }
  _refreshRows(e, i8) {
    this._liveRegionDebouncer.refresh(e, i8, this._terminal.rows);
  }
  _renderRows(e, i8) {
    let r5 = this._terminal.buffer, n = r5.lines.length.toString();
    for (let o2 = e; o2 <= i8; o2++) {
      let l = r5.lines.get(r5.ydisp + o2), a = [], u = l?.translateToString(true, void 0, void 0, a) || "", h2 = (r5.ydisp + o2 + 1).toString(), c = this._rowElements[o2];
      c && (u.length === 0 ? (c.textContent = " ", this._rowColumns.set(c, [0, 1])) : (c.textContent = u, this._rowColumns.set(c, a)), c.setAttribute("aria-posinset", h2), c.setAttribute("aria-setsize", n), this._alignRowWidth(c));
    }
    this._announceCharacters();
  }
  _announceCharacters() {
    this._charsToAnnounce.length !== 0 && (this._liveRegion.textContent += this._charsToAnnounce, this._charsToAnnounce = "");
  }
  _handleBoundaryFocus(e, i8) {
    let r5 = e.target, n = this._rowElements[i8 === 0 ? 1 : this._rowElements.length - 2], o2 = r5.getAttribute("aria-posinset"), l = i8 === 0 ? "1" : `${this._terminal.buffer.lines.length}`;
    if (o2 === l || e.relatedTarget !== n) return;
    let a, u;
    if (i8 === 0 ? (a = r5, u = this._rowElements.pop(), this._rowContainer.removeChild(u)) : (a = this._rowElements.shift(), u = r5, this._rowContainer.removeChild(a)), a.removeEventListener("focus", this._topBoundaryFocusListener), u.removeEventListener("focus", this._bottomBoundaryFocusListener), i8 === 0) {
      let h2 = this._createAccessibilityTreeNode();
      this._rowElements.unshift(h2), this._rowContainer.insertAdjacentElement("afterbegin", h2);
    } else {
      let h2 = this._createAccessibilityTreeNode();
      this._rowElements.push(h2), this._rowContainer.appendChild(h2);
    }
    this._rowElements[0].addEventListener("focus", this._topBoundaryFocusListener), this._rowElements[this._rowElements.length - 1].addEventListener("focus", this._bottomBoundaryFocusListener), this._terminal.scrollLines(i8 === 0 ? -1 : 1), this._rowElements[i8 === 0 ? 1 : this._rowElements.length - 2].focus(), e.preventDefault(), e.stopImmediatePropagation();
  }
  _handleSelectionChange() {
    if (this._rowElements.length === 0) return;
    let e = this._coreBrowserService.mainDocument.getSelection();
    if (!e) return;
    if (e.isCollapsed) {
      this._rowContainer.contains(e.anchorNode) && this._terminal.clearSelection();
      return;
    }
    if (!e.anchorNode || !e.focusNode) {
      console.error("anchorNode and/or focusNode are null");
      return;
    }
    let i8 = { node: e.anchorNode, offset: e.anchorOffset }, r5 = { node: e.focusNode, offset: e.focusOffset };
    if ((i8.node.compareDocumentPosition(r5.node) & Node.DOCUMENT_POSITION_PRECEDING || i8.node === r5.node && i8.offset > r5.offset) && ([i8, r5] = [r5, i8]), i8.node.compareDocumentPosition(this._rowElements[0]) & (Node.DOCUMENT_POSITION_CONTAINED_BY | Node.DOCUMENT_POSITION_FOLLOWING) && (i8 = { node: this._rowElements[0].childNodes[0], offset: 0 }), !this._rowContainer.contains(i8.node)) return;
    let n = this._rowElements.slice(-1)[0];
    if (r5.node.compareDocumentPosition(n) & (Node.DOCUMENT_POSITION_CONTAINED_BY | Node.DOCUMENT_POSITION_PRECEDING) && (r5 = { node: n, offset: n.textContent?.length ?? 0 }), !this._rowContainer.contains(r5.node)) return;
    let o2 = ({ node: u, offset: h2 }) => {
      let c = u instanceof Text ? u.parentNode : u, d = parseInt(c?.getAttribute("aria-posinset"), 10) - 1;
      if (isNaN(d)) return console.warn("row is invalid. Race condition?"), null;
      let _3 = this._rowColumns.get(c);
      if (!_3) return console.warn("columns is null. Race condition?"), null;
      let p = h2 < _3.length ? _3[h2] : _3.slice(-1)[0] + 1;
      return p >= this._terminal.cols && (++d, p = 0), { row: d, column: p };
    }, l = o2(i8), a = o2(r5);
    if (!(!l || !a)) {
      if (l.row > a.row || l.row === a.row && l.column >= a.column) throw new Error("invalid range");
      this._terminal.select(l.column, l.row, (a.row - l.row) * this._terminal.cols - l.column + a.column);
    }
  }
  _handleResize(e) {
    this._rowElements[this._rowElements.length - 1].removeEventListener("focus", this._bottomBoundaryFocusListener);
    for (let i8 = this._rowContainer.children.length; i8 < this._terminal.rows; i8++) this._rowElements[i8] = this._createAccessibilityTreeNode(), this._rowContainer.appendChild(this._rowElements[i8]);
    for (; this._rowElements.length > e; ) this._rowContainer.removeChild(this._rowElements.pop());
    this._rowElements[this._rowElements.length - 1].addEventListener("focus", this._bottomBoundaryFocusListener), this._refreshRowsDimensions();
  }
  _createAccessibilityTreeNode() {
    let e = this._coreBrowserService.mainDocument.createElement("div");
    return e.setAttribute("role", "listitem"), e.tabIndex = -1, this._refreshRowDimensions(e), e;
  }
  _refreshRowsDimensions() {
    if (this._renderService.dimensions.css.cell.height) {
      Object.assign(this._accessibilityContainer.style, { width: `${this._renderService.dimensions.css.canvas.width}px`, fontSize: `${this._terminal.options.fontSize}px` }), this._rowElements.length !== this._terminal.rows && this._handleResize(this._terminal.rows);
      for (let e = 0; e < this._terminal.rows; e++) this._refreshRowDimensions(this._rowElements[e]), this._alignRowWidth(this._rowElements[e]);
    }
  }
  _refreshRowDimensions(e) {
    e.style.height = `${this._renderService.dimensions.css.cell.height}px`;
  }
  _alignRowWidth(e) {
    e.style.transform = "";
    let i8 = e.getBoundingClientRect().width, r5 = this._rowColumns.get(e)?.slice(-1)?.[0];
    if (!r5) return;
    let n = r5 * this._renderService.dimensions.css.cell.width;
    e.style.transform = `scaleX(${n / i8})`;
  }
};
Tt2 = M3([S(1, xt4), S(2, ae2), S(3, ce2)], Tt2);
var hi = class extends D2 {
  constructor(e, i8, r5, n, o2) {
    super();
    this._element = e;
    this._mouseService = i8;
    this._renderService = r5;
    this._bufferService = n;
    this._linkProviderService = o2;
    this._linkCacheDisposables = [];
    this._isMouseOut = true;
    this._wasResized = false;
    this._activeLine = -1;
    this._onShowLinkUnderline = this._register(new v());
    this.onShowLinkUnderline = this._onShowLinkUnderline.event;
    this._onHideLinkUnderline = this._register(new v());
    this.onHideLinkUnderline = this._onHideLinkUnderline.event;
    this._register(C3(() => {
      Ne2(this._linkCacheDisposables), this._linkCacheDisposables.length = 0, this._lastMouseEvent = void 0, this._activeProviderReplies?.clear();
    })), this._register(this._bufferService.onResize(() => {
      this._clearCurrentLink(), this._wasResized = true;
    })), this._register(L(this._element, "mouseleave", () => {
      this._isMouseOut = true, this._clearCurrentLink();
    })), this._register(L(this._element, "mousemove", this._handleMouseMove.bind(this))), this._register(L(this._element, "mousedown", this._handleMouseDown.bind(this))), this._register(L(this._element, "mouseup", this._handleMouseUp.bind(this)));
  }
  get currentLink() {
    return this._currentLink;
  }
  _handleMouseMove(e) {
    this._lastMouseEvent = e;
    let i8 = this._positionFromMouseEvent(e, this._element, this._mouseService);
    if (!i8) return;
    this._isMouseOut = false;
    let r5 = e.composedPath();
    for (let n = 0; n < r5.length; n++) {
      let o2 = r5[n];
      if (o2.classList.contains("xterm")) break;
      if (o2.classList.contains("xterm-hover")) return;
    }
    (!this._lastBufferCell || i8.x !== this._lastBufferCell.x || i8.y !== this._lastBufferCell.y) && (this._handleHover(i8), this._lastBufferCell = i8);
  }
  _handleHover(e) {
    if (this._activeLine !== e.y || this._wasResized) {
      this._clearCurrentLink(), this._askForLink(e, false), this._wasResized = false;
      return;
    }
    this._currentLink && this._linkAtPosition(this._currentLink.link, e) || (this._clearCurrentLink(), this._askForLink(e, true));
  }
  _askForLink(e, i8) {
    (!this._activeProviderReplies || !i8) && (this._activeProviderReplies?.forEach((n) => {
      n?.forEach((o2) => {
        o2.link.dispose && o2.link.dispose();
      });
    }), this._activeProviderReplies = /* @__PURE__ */ new Map(), this._activeLine = e.y);
    let r5 = false;
    for (let [n, o2] of this._linkProviderService.linkProviders.entries()) i8 ? this._activeProviderReplies?.get(n) && (r5 = this._checkLinkProviderResult(n, e, r5)) : o2.provideLinks(e.y, (l) => {
      if (this._isMouseOut) return;
      let a = l?.map((u) => ({ link: u }));
      this._activeProviderReplies?.set(n, a), r5 = this._checkLinkProviderResult(n, e, r5), this._activeProviderReplies?.size === this._linkProviderService.linkProviders.length && this._removeIntersectingLinks(e.y, this._activeProviderReplies);
    });
  }
  _removeIntersectingLinks(e, i8) {
    let r5 = /* @__PURE__ */ new Set();
    for (let n = 0; n < i8.size; n++) {
      let o2 = i8.get(n);
      if (o2) for (let l = 0; l < o2.length; l++) {
        let a = o2[l], u = a.link.range.start.y < e ? 0 : a.link.range.start.x, h2 = a.link.range.end.y > e ? this._bufferService.cols : a.link.range.end.x;
        for (let c = u; c <= h2; c++) {
          if (r5.has(c)) {
            o2.splice(l--, 1);
            break;
          }
          r5.add(c);
        }
      }
    }
  }
  _checkLinkProviderResult(e, i8, r5) {
    if (!this._activeProviderReplies) return r5;
    let n = this._activeProviderReplies.get(e), o2 = false;
    for (let l = 0; l < e; l++) (!this._activeProviderReplies.has(l) || this._activeProviderReplies.get(l)) && (o2 = true);
    if (!o2 && n) {
      let l = n.find((a) => this._linkAtPosition(a.link, i8));
      l && (r5 = true, this._handleNewLink(l));
    }
    if (this._activeProviderReplies.size === this._linkProviderService.linkProviders.length && !r5) for (let l = 0; l < this._activeProviderReplies.size; l++) {
      let a = this._activeProviderReplies.get(l)?.find((u) => this._linkAtPosition(u.link, i8));
      if (a) {
        r5 = true, this._handleNewLink(a);
        break;
      }
    }
    return r5;
  }
  _handleMouseDown() {
    this._mouseDownLink = this._currentLink;
  }
  _handleMouseUp(e) {
    if (!this._currentLink) return;
    let i8 = this._positionFromMouseEvent(e, this._element, this._mouseService);
    i8 && this._mouseDownLink && Ec(this._mouseDownLink.link, this._currentLink.link) && this._linkAtPosition(this._currentLink.link, i8) && this._currentLink.link.activate(e, this._currentLink.link.text);
  }
  _clearCurrentLink(e, i8) {
    !this._currentLink || !this._lastMouseEvent || (!e || !i8 || this._currentLink.link.range.start.y >= e && this._currentLink.link.range.end.y <= i8) && (this._linkLeave(this._element, this._currentLink.link, this._lastMouseEvent), this._currentLink = void 0, Ne2(this._linkCacheDisposables), this._linkCacheDisposables.length = 0);
  }
  _handleNewLink(e) {
    if (!this._lastMouseEvent) return;
    let i8 = this._positionFromMouseEvent(this._lastMouseEvent, this._element, this._mouseService);
    i8 && this._linkAtPosition(e.link, i8) && (this._currentLink = e, this._currentLink.state = { decorations: { underline: e.link.decorations === void 0 ? true : e.link.decorations.underline, pointerCursor: e.link.decorations === void 0 ? true : e.link.decorations.pointerCursor }, isHovered: true }, this._linkHover(this._element, e.link, this._lastMouseEvent), e.link.decorations = {}, Object.defineProperties(e.link.decorations, { pointerCursor: { get: () => this._currentLink?.state?.decorations.pointerCursor, set: (r5) => {
      this._currentLink?.state && this._currentLink.state.decorations.pointerCursor !== r5 && (this._currentLink.state.decorations.pointerCursor = r5, this._currentLink.state.isHovered && this._element.classList.toggle("xterm-cursor-pointer", r5));
    } }, underline: { get: () => this._currentLink?.state?.decorations.underline, set: (r5) => {
      this._currentLink?.state && this._currentLink?.state?.decorations.underline !== r5 && (this._currentLink.state.decorations.underline = r5, this._currentLink.state.isHovered && this._fireUnderlineEvent(e.link, r5));
    } } }), this._linkCacheDisposables.push(this._renderService.onRenderedViewportChange((r5) => {
      if (!this._currentLink) return;
      let n = r5.start === 0 ? 0 : r5.start + 1 + this._bufferService.buffer.ydisp, o2 = this._bufferService.buffer.ydisp + 1 + r5.end;
      if (this._currentLink.link.range.start.y >= n && this._currentLink.link.range.end.y <= o2 && (this._clearCurrentLink(n, o2), this._lastMouseEvent)) {
        let l = this._positionFromMouseEvent(this._lastMouseEvent, this._element, this._mouseService);
        l && this._askForLink(l, false);
      }
    })));
  }
  _linkHover(e, i8, r5) {
    this._currentLink?.state && (this._currentLink.state.isHovered = true, this._currentLink.state.decorations.underline && this._fireUnderlineEvent(i8, true), this._currentLink.state.decorations.pointerCursor && e.classList.add("xterm-cursor-pointer")), i8.hover && i8.hover(r5, i8.text);
  }
  _fireUnderlineEvent(e, i8) {
    let r5 = e.range, n = this._bufferService.buffer.ydisp, o2 = this._createLinkUnderlineEvent(r5.start.x - 1, r5.start.y - n - 1, r5.end.x, r5.end.y - n - 1, void 0);
    (i8 ? this._onShowLinkUnderline : this._onHideLinkUnderline).fire(o2);
  }
  _linkLeave(e, i8, r5) {
    this._currentLink?.state && (this._currentLink.state.isHovered = false, this._currentLink.state.decorations.underline && this._fireUnderlineEvent(i8, false), this._currentLink.state.decorations.pointerCursor && e.classList.remove("xterm-cursor-pointer")), i8.leave && i8.leave(r5, i8.text);
  }
  _linkAtPosition(e, i8) {
    let r5 = e.range.start.y * this._bufferService.cols + e.range.start.x, n = e.range.end.y * this._bufferService.cols + e.range.end.x, o2 = i8.y * this._bufferService.cols + i8.x;
    return r5 <= o2 && o2 <= n;
  }
  _positionFromMouseEvent(e, i8, r5) {
    let n = r5.getCoords(e, i8, this._bufferService.cols, this._bufferService.rows);
    if (n) return { x: n[0], y: n[1] + this._bufferService.buffer.ydisp };
  }
  _createLinkUnderlineEvent(e, i8, r5, n, o2) {
    return { x1: e, y1: i8, x2: r5, y2: n, cols: this._bufferService.cols, fg: o2 };
  }
};
hi = M3([S(1, Dt2), S(2, ce2), S(3, F2), S(4, lr)], hi);
function Ec(s15, t) {
  return s15.text === t.text && s15.range.start.x === t.range.start.x && s15.range.start.y === t.range.start.y && s15.range.end.x === t.range.end.x && s15.range.end.y === t.range.end.y;
}
var yn2 = class extends Sn2 {
  constructor(e = {}) {
    super(e);
    this._linkifier = this._register(new ye3());
    this.browser = tn2;
    this._keyDownHandled = false;
    this._keyDownSeen = false;
    this._keyPressHandled = false;
    this._unprocessedDeadKey = false;
    this._accessibilityManager = this._register(new ye3());
    this._onCursorMove = this._register(new v());
    this.onCursorMove = this._onCursorMove.event;
    this._onKey = this._register(new v());
    this.onKey = this._onKey.event;
    this._onRender = this._register(new v());
    this.onRender = this._onRender.event;
    this._onSelectionChange = this._register(new v());
    this.onSelectionChange = this._onSelectionChange.event;
    this._onTitleChange = this._register(new v());
    this.onTitleChange = this._onTitleChange.event;
    this._onBell = this._register(new v());
    this.onBell = this._onBell.event;
    this._onFocus = this._register(new v());
    this._onBlur = this._register(new v());
    this._onA11yCharEmitter = this._register(new v());
    this._onA11yTabEmitter = this._register(new v());
    this._onWillOpen = this._register(new v());
    this._setup(), this._decorationService = this._instantiationService.createInstance(Tn), this._instantiationService.setService(Be4, this._decorationService), this._linkProviderService = this._instantiationService.createInstance(Qr2), this._instantiationService.setService(lr, this._linkProviderService), this._linkProviderService.registerLinkProvider(this._instantiationService.createInstance(wt2)), this._register(this._inputHandler.onRequestBell(() => this._onBell.fire())), this._register(this._inputHandler.onRequestRefreshRows((i8) => this.refresh(i8?.start ?? 0, i8?.end ?? this.rows - 1))), this._register(this._inputHandler.onRequestSendFocus(() => this._reportFocus())), this._register(this._inputHandler.onRequestReset(() => this.reset())), this._register(this._inputHandler.onRequestWindowsOptionsReport((i8) => this._reportWindowsOptions(i8))), this._register(this._inputHandler.onColor((i8) => this._handleColorEvent(i8))), this._register($.forward(this._inputHandler.onCursorMove, this._onCursorMove)), this._register($.forward(this._inputHandler.onTitleChange, this._onTitleChange)), this._register($.forward(this._inputHandler.onA11yChar, this._onA11yCharEmitter)), this._register($.forward(this._inputHandler.onA11yTab, this._onA11yTabEmitter)), this._register(this._bufferService.onResize((i8) => this._afterResize(i8.cols, i8.rows))), this._register(C3(() => {
      this._customKeyEventHandler = void 0, this.element?.parentNode?.removeChild(this.element);
    }));
  }
  get linkifier() {
    return this._linkifier.value;
  }
  get onFocus() {
    return this._onFocus.event;
  }
  get onBlur() {
    return this._onBlur.event;
  }
  get onA11yChar() {
    return this._onA11yCharEmitter.event;
  }
  get onA11yTab() {
    return this._onA11yTabEmitter.event;
  }
  get onWillOpen() {
    return this._onWillOpen.event;
  }
  _handleColorEvent(e) {
    if (this._themeService) for (let i8 of e) {
      let r5, n = "";
      switch (i8.index) {
        case 256:
          r5 = "foreground", n = "10";
          break;
        case 257:
          r5 = "background", n = "11";
          break;
        case 258:
          r5 = "cursor", n = "12";
          break;
        default:
          r5 = "ansi", n = "4;" + i8.index;
      }
      switch (i8.type) {
        case 0:
          let o2 = U3.toColorRGB(r5 === "ansi" ? this._themeService.colors.ansi[i8.index] : this._themeService.colors[r5]);
          this.coreService.triggerDataEvent(`${b.ESC}]${n};${ml(o2)}${fs2.ST}`);
          break;
        case 1:
          if (r5 === "ansi") this._themeService.modifyColors((l) => l.ansi[i8.index] = j2.toColor(...i8.color));
          else {
            let l = r5;
            this._themeService.modifyColors((a) => a[l] = j2.toColor(...i8.color));
          }
          break;
        case 2:
          this._themeService.restoreColor(i8.index);
          break;
      }
    }
  }
  _setup() {
    super._setup(), this._customKeyEventHandler = void 0;
  }
  get buffer() {
    return this.buffers.active;
  }
  focus() {
    this.textarea && this.textarea.focus({ preventScroll: true });
  }
  _handleScreenReaderModeOptionChange(e) {
    e ? !this._accessibilityManager.value && this._renderService && (this._accessibilityManager.value = this._instantiationService.createInstance(Tt2, this)) : this._accessibilityManager.clear();
  }
  _handleTextAreaFocus(e) {
    this.coreService.decPrivateModes.sendFocus && this.coreService.triggerDataEvent(b.ESC + "[I"), this.element.classList.add("focus"), this._showCursor(), this._onFocus.fire();
  }
  blur() {
    return this.textarea?.blur();
  }
  _handleTextAreaBlur() {
    this.textarea.value = "", this.refresh(this.buffer.y, this.buffer.y), this.coreService.decPrivateModes.sendFocus && this.coreService.triggerDataEvent(b.ESC + "[O"), this.element.classList.remove("focus"), this._onBlur.fire();
  }
  _syncTextArea() {
    if (!this.textarea || !this.buffer.isCursorInViewport || this._compositionHelper.isComposing || !this._renderService) return;
    let e = this.buffer.ybase + this.buffer.y, i8 = this.buffer.lines.get(e);
    if (!i8) return;
    let r5 = Math.min(this.buffer.x, this.cols - 1), n = this._renderService.dimensions.css.cell.height, o2 = i8.getWidth(r5), l = this._renderService.dimensions.css.cell.width * o2, a = this.buffer.y * this._renderService.dimensions.css.cell.height, u = r5 * this._renderService.dimensions.css.cell.width;
    this.textarea.style.left = u + "px", this.textarea.style.top = a + "px", this.textarea.style.width = l + "px", this.textarea.style.height = n + "px", this.textarea.style.lineHeight = n + "px", this.textarea.style.zIndex = "-5";
  }
  _initGlobal() {
    this._bindKeys(), this._register(L(this.element, "copy", (i8) => {
      this.hasSelection() && Vs(i8, this._selectionService);
    }));
    let e = (i8) => qs(i8, this.textarea, this.coreService, this.optionsService);
    this._register(L(this.textarea, "paste", e)), this._register(L(this.element, "paste", e)), Ss2 ? this._register(L(this.element, "mousedown", (i8) => {
      i8.button === 2 && Pn2(i8, this.textarea, this.screenElement, this._selectionService, this.options.rightClickSelectsWord);
    })) : this._register(L(this.element, "contextmenu", (i8) => {
      Pn2(i8, this.textarea, this.screenElement, this._selectionService, this.options.rightClickSelectsWord);
    })), Bi2 && this._register(L(this.element, "auxclick", (i8) => {
      i8.button === 1 && Mn2(i8, this.textarea, this.screenElement);
    }));
  }
  _bindKeys() {
    this._register(L(this.textarea, "keyup", (e) => this._keyUp(e), true)), this._register(L(this.textarea, "keydown", (e) => this._keyDown(e), true)), this._register(L(this.textarea, "keypress", (e) => this._keyPress(e), true)), this._register(L(this.textarea, "compositionstart", () => this._compositionHelper.compositionstart())), this._register(L(this.textarea, "compositionupdate", (e) => this._compositionHelper.compositionupdate(e))), this._register(L(this.textarea, "compositionend", () => this._compositionHelper.compositionend())), this._register(L(this.textarea, "input", (e) => this._inputEvent(e), true)), this._register(this.onRender(() => this._compositionHelper.updateCompositionElements()));
  }
  open(e) {
    if (!e) throw new Error("Terminal requires a parent element.");
    if (e.isConnected || this._logService.debug("Terminal.open was called on an element that was not attached to the DOM"), this.element?.ownerDocument.defaultView && this._coreBrowserService) {
      this.element.ownerDocument.defaultView !== this._coreBrowserService.window && (this._coreBrowserService.window = this.element.ownerDocument.defaultView);
      return;
    }
    this._document = e.ownerDocument, this.options.documentOverride && this.options.documentOverride instanceof Document && (this._document = this.optionsService.rawOptions.documentOverride), this.element = this._document.createElement("div"), this.element.dir = "ltr", this.element.classList.add("terminal"), this.element.classList.add("xterm"), e.appendChild(this.element);
    let i8 = this._document.createDocumentFragment();
    this._viewportElement = this._document.createElement("div"), this._viewportElement.classList.add("xterm-viewport"), i8.appendChild(this._viewportElement), this.screenElement = this._document.createElement("div"), this.screenElement.classList.add("xterm-screen"), this._register(L(this.screenElement, "mousemove", (o2) => this.updateCursorStyle(o2))), this._helperContainer = this._document.createElement("div"), this._helperContainer.classList.add("xterm-helpers"), this.screenElement.appendChild(this._helperContainer), i8.appendChild(this.screenElement);
    let r5 = this.textarea = this._document.createElement("textarea");
    this.textarea.classList.add("xterm-helper-textarea"), this.textarea.setAttribute("aria-label", mi2.get()), Ts2 || this.textarea.setAttribute("aria-multiline", "false"), this.textarea.setAttribute("autocorrect", "off"), this.textarea.setAttribute("autocapitalize", "off"), this.textarea.setAttribute("spellcheck", "false"), this.textarea.tabIndex = 0, this._register(this.optionsService.onSpecificOptionChange("disableStdin", () => r5.readOnly = this.optionsService.rawOptions.disableStdin)), this.textarea.readOnly = this.optionsService.rawOptions.disableStdin, this._coreBrowserService = this._register(this._instantiationService.createInstance(Jr2, this.textarea, e.ownerDocument.defaultView ?? window, this._document ?? typeof window < "u" ? window.document : null)), this._instantiationService.setService(ae2, this._coreBrowserService), this._register(L(this.textarea, "focus", (o2) => this._handleTextAreaFocus(o2))), this._register(L(this.textarea, "blur", () => this._handleTextAreaBlur())), this._helperContainer.appendChild(this.textarea), this._charSizeService = this._instantiationService.createInstance(jt2, this._document, this._helperContainer), this._instantiationService.setService(nt3, this._charSizeService), this._themeService = this._instantiationService.createInstance(ti2), this._instantiationService.setService(Re2, this._themeService), this._characterJoinerService = this._instantiationService.createInstance(ct), this._instantiationService.setService(or2, this._characterJoinerService), this._renderService = this._register(this._instantiationService.createInstance(Qt2, this.rows, this.screenElement)), this._instantiationService.setService(ce2, this._renderService), this._register(this._renderService.onRenderedViewportChange((o2) => this._onRender.fire(o2))), this.onResize((o2) => this._renderService.resize(o2.cols, o2.rows)), this._compositionView = this._document.createElement("div"), this._compositionView.classList.add("composition-view"), this._compositionHelper = this._instantiationService.createInstance($t2, this.textarea, this._compositionView), this._helperContainer.appendChild(this._compositionView), this._mouseService = this._instantiationService.createInstance(Xt2), this._instantiationService.setService(Dt2, this._mouseService);
    let n = this._linkifier.value = this._register(this._instantiationService.createInstance(hi, this.screenElement));
    this.element.appendChild(i8);
    try {
      this._onWillOpen.fire(this.element);
    } catch {
    }
    this._renderService.hasRenderer() || this._renderService.setRenderer(this._createRenderer()), this._register(this.onCursorMove(() => {
      this._renderService.handleCursorMove(), this._syncTextArea();
    })), this._register(this.onResize(() => this._renderService.handleResize(this.cols, this.rows))), this._register(this.onBlur(() => this._renderService.handleBlur())), this._register(this.onFocus(() => this._renderService.handleFocus())), this._viewport = this._register(this._instantiationService.createInstance(zt2, this.element, this.screenElement)), this._register(this._viewport.onRequestScrollLines((o2) => {
      super.scrollLines(o2, false), this.refresh(0, this.rows - 1);
    })), this._selectionService = this._register(this._instantiationService.createInstance(ei2, this.element, this.screenElement, n)), this._instantiationService.setService(Qs, this._selectionService), this._register(this._selectionService.onRequestScrollLines((o2) => this.scrollLines(o2.amount, o2.suppressScrollEvent))), this._register(this._selectionService.onSelectionChange(() => this._onSelectionChange.fire())), this._register(this._selectionService.onRequestRedraw((o2) => this._renderService.handleSelectionChanged(o2.start, o2.end, o2.columnSelectMode))), this._register(this._selectionService.onLinuxMouseSelection((o2) => {
      this.textarea.value = o2, this.textarea.focus(), this.textarea.select();
    })), this._register($.any(this._onScroll.event, this._inputHandler.onScroll)(() => {
      this._selectionService.refresh(), this._viewport?.queueSync();
    })), this._register(this._instantiationService.createInstance(Gt2, this.screenElement)), this._register(L(this.element, "mousedown", (o2) => this._selectionService.handleMouseDown(o2))), this.coreMouseService.areMouseEventsActive ? (this._selectionService.disable(), this.element.classList.add("enable-mouse-events")) : this._selectionService.enable(), this.options.screenReaderMode && (this._accessibilityManager.value = this._instantiationService.createInstance(Tt2, this)), this._register(this.optionsService.onSpecificOptionChange("screenReaderMode", (o2) => this._handleScreenReaderModeOptionChange(o2))), this.options.overviewRuler.width && (this._overviewRulerRenderer = this._register(this._instantiationService.createInstance(bt2, this._viewportElement, this.screenElement))), this.optionsService.onSpecificOptionChange("overviewRuler", (o2) => {
      !this._overviewRulerRenderer && o2 && this._viewportElement && this.screenElement && (this._overviewRulerRenderer = this._register(this._instantiationService.createInstance(bt2, this._viewportElement, this.screenElement)));
    }), this._charSizeService.measure(), this.refresh(0, this.rows - 1), this._initGlobal(), this.bindMouse();
  }
  _createRenderer() {
    return this._instantiationService.createInstance(Yt2, this, this._document, this.element, this.screenElement, this._viewportElement, this._helperContainer, this.linkifier);
  }
  bindMouse() {
    let e = this, i8 = this.element;
    function r5(l) {
      let a = e._mouseService.getMouseReportCoords(l, e.screenElement);
      if (!a) return false;
      let u, h2;
      switch (l.overrideType || l.type) {
        case "mousemove":
          h2 = 32, l.buttons === void 0 ? (u = 3, l.button !== void 0 && (u = l.button < 3 ? l.button : 3)) : u = l.buttons & 1 ? 0 : l.buttons & 4 ? 1 : l.buttons & 2 ? 2 : 3;
          break;
        case "mouseup":
          h2 = 0, u = l.button < 3 ? l.button : 3;
          break;
        case "mousedown":
          h2 = 1, u = l.button < 3 ? l.button : 3;
          break;
        case "wheel":
          if (e._customWheelEventHandler && e._customWheelEventHandler(l) === false) return false;
          let c = l.deltaY;
          if (c === 0 || e.coreMouseService.consumeWheelEvent(l, e._renderService?.dimensions?.device?.cell?.height, e._coreBrowserService?.dpr) === 0) return false;
          h2 = c < 0 ? 0 : 1, u = 4;
          break;
        default:
          return false;
      }
      return h2 === void 0 || u === void 0 || u > 4 ? false : e.coreMouseService.triggerMouseEvent({ col: a.col, row: a.row, x: a.x, y: a.y, button: u, action: h2, ctrl: l.ctrlKey, alt: l.altKey, shift: l.shiftKey });
    }
    let n = { mouseup: null, wheel: null, mousedrag: null, mousemove: null }, o2 = { mouseup: (l) => (r5(l), l.buttons || (this._document.removeEventListener("mouseup", n.mouseup), n.mousedrag && this._document.removeEventListener("mousemove", n.mousedrag)), this.cancel(l)), wheel: (l) => (r5(l), this.cancel(l, true)), mousedrag: (l) => {
      l.buttons && r5(l);
    }, mousemove: (l) => {
      l.buttons || r5(l);
    } };
    this._register(this.coreMouseService.onProtocolChange((l) => {
      l ? (this.optionsService.rawOptions.logLevel === "debug" && this._logService.debug("Binding to mouse events:", this.coreMouseService.explainEvents(l)), this.element.classList.add("enable-mouse-events"), this._selectionService.disable()) : (this._logService.debug("Unbinding from mouse events."), this.element.classList.remove("enable-mouse-events"), this._selectionService.enable()), l & 8 ? n.mousemove || (i8.addEventListener("mousemove", o2.mousemove), n.mousemove = o2.mousemove) : (i8.removeEventListener("mousemove", n.mousemove), n.mousemove = null), l & 16 ? n.wheel || (i8.addEventListener("wheel", o2.wheel, { passive: false }), n.wheel = o2.wheel) : (i8.removeEventListener("wheel", n.wheel), n.wheel = null), l & 2 ? n.mouseup || (n.mouseup = o2.mouseup) : (this._document.removeEventListener("mouseup", n.mouseup), n.mouseup = null), l & 4 ? n.mousedrag || (n.mousedrag = o2.mousedrag) : (this._document.removeEventListener("mousemove", n.mousedrag), n.mousedrag = null);
    })), this.coreMouseService.activeProtocol = this.coreMouseService.activeProtocol, this._register(L(i8, "mousedown", (l) => {
      if (l.preventDefault(), this.focus(), !(!this.coreMouseService.areMouseEventsActive || this._selectionService.shouldForceSelection(l))) return r5(l), n.mouseup && this._document.addEventListener("mouseup", n.mouseup), n.mousedrag && this._document.addEventListener("mousemove", n.mousedrag), this.cancel(l);
    })), this._register(L(i8, "wheel", (l) => {
      if (!n.wheel) {
        if (this._customWheelEventHandler && this._customWheelEventHandler(l) === false) return false;
        if (!this.buffer.hasScrollback) {
          if (l.deltaY === 0) return false;
          if (e.coreMouseService.consumeWheelEvent(l, e._renderService?.dimensions?.device?.cell?.height, e._coreBrowserService?.dpr) === 0) return this.cancel(l, true);
          let h2 = b.ESC + (this.coreService.decPrivateModes.applicationCursorKeys ? "O" : "[") + (l.deltaY < 0 ? "A" : "B");
          return this.coreService.triggerDataEvent(h2, true), this.cancel(l, true);
        }
      }
    }, { passive: false }));
  }
  refresh(e, i8) {
    this._renderService?.refreshRows(e, i8);
  }
  updateCursorStyle(e) {
    this._selectionService?.shouldColumnSelect(e) ? this.element.classList.add("column-select") : this.element.classList.remove("column-select");
  }
  _showCursor() {
    this.coreService.isCursorInitialized || (this.coreService.isCursorInitialized = true, this.refresh(this.buffer.y, this.buffer.y));
  }
  scrollLines(e, i8) {
    this._viewport ? this._viewport.scrollLines(e) : super.scrollLines(e, i8), this.refresh(0, this.rows - 1);
  }
  scrollPages(e) {
    this.scrollLines(e * (this.rows - 1));
  }
  scrollToTop() {
    this.scrollLines(-this._bufferService.buffer.ydisp);
  }
  scrollToBottom(e) {
    e && this._viewport ? this._viewport.scrollToLine(this.buffer.ybase, true) : this.scrollLines(this._bufferService.buffer.ybase - this._bufferService.buffer.ydisp);
  }
  scrollToLine(e) {
    let i8 = e - this._bufferService.buffer.ydisp;
    i8 !== 0 && this.scrollLines(i8);
  }
  paste(e) {
    Cn2(e, this.textarea, this.coreService, this.optionsService);
  }
  attachCustomKeyEventHandler(e) {
    this._customKeyEventHandler = e;
  }
  attachCustomWheelEventHandler(e) {
    this._customWheelEventHandler = e;
  }
  registerLinkProvider(e) {
    return this._linkProviderService.registerLinkProvider(e);
  }
  registerCharacterJoiner(e) {
    if (!this._characterJoinerService) throw new Error("Terminal must be opened first");
    let i8 = this._characterJoinerService.register(e);
    return this.refresh(0, this.rows - 1), i8;
  }
  deregisterCharacterJoiner(e) {
    if (!this._characterJoinerService) throw new Error("Terminal must be opened first");
    this._characterJoinerService.deregister(e) && this.refresh(0, this.rows - 1);
  }
  get markers() {
    return this.buffer.markers;
  }
  registerMarker(e) {
    return this.buffer.addMarker(this.buffer.ybase + this.buffer.y + e);
  }
  registerDecoration(e) {
    return this._decorationService.registerDecoration(e);
  }
  hasSelection() {
    return this._selectionService ? this._selectionService.hasSelection : false;
  }
  select(e, i8, r5) {
    this._selectionService.setSelection(e, i8, r5);
  }
  getSelection() {
    return this._selectionService ? this._selectionService.selectionText : "";
  }
  getSelectionPosition() {
    if (!(!this._selectionService || !this._selectionService.hasSelection)) return { start: { x: this._selectionService.selectionStart[0], y: this._selectionService.selectionStart[1] }, end: { x: this._selectionService.selectionEnd[0], y: this._selectionService.selectionEnd[1] } };
  }
  clearSelection() {
    this._selectionService?.clearSelection();
  }
  selectAll() {
    this._selectionService?.selectAll();
  }
  selectLines(e, i8) {
    this._selectionService?.selectLines(e, i8);
  }
  _keyDown(e) {
    if (this._keyDownHandled = false, this._keyDownSeen = true, this._customKeyEventHandler && this._customKeyEventHandler(e) === false) return false;
    let i8 = this.browser.isMac && this.options.macOptionIsMeta && e.altKey;
    if (!i8 && !this._compositionHelper.keydown(e)) return this.options.scrollOnUserInput && this.buffer.ybase !== this.buffer.ydisp && this.scrollToBottom(true), false;
    !i8 && (e.key === "Dead" || e.key === "AltGraph") && (this._unprocessedDeadKey = true);
    let r5 = Il(e, this.coreService.decPrivateModes.applicationCursorKeys, this.browser.isMac, this.options.macOptionIsMeta);
    if (this.updateCursorStyle(e), r5.type === 3 || r5.type === 2) {
      let n = this.rows - 1;
      return this.scrollLines(r5.type === 2 ? -n : n), this.cancel(e, true);
    }
    if (r5.type === 1 && this.selectAll(), this._isThirdLevelShift(this.browser, e) || (r5.cancel && this.cancel(e, true), !r5.key) || e.key && !e.ctrlKey && !e.altKey && !e.metaKey && e.key.length === 1 && e.key.charCodeAt(0) >= 65 && e.key.charCodeAt(0) <= 90) return true;
    if (this._unprocessedDeadKey) return this._unprocessedDeadKey = false, true;
    if ((r5.key === b.ETX || r5.key === b.CR) && (this.textarea.value = ""), this._onKey.fire({ key: r5.key, domEvent: e }), this._showCursor(), this.coreService.triggerDataEvent(r5.key, true), !this.optionsService.rawOptions.screenReaderMode || e.altKey || e.ctrlKey) return this.cancel(e, true);
    this._keyDownHandled = true;
  }
  _isThirdLevelShift(e, i8) {
    let r5 = e.isMac && !this.options.macOptionIsMeta && i8.altKey && !i8.ctrlKey && !i8.metaKey || e.isWindows && i8.altKey && i8.ctrlKey && !i8.metaKey || e.isWindows && i8.getModifierState("AltGraph");
    return i8.type === "keypress" ? r5 : r5 && (!i8.keyCode || i8.keyCode > 47);
  }
  _keyUp(e) {
    this._keyDownSeen = false, !(this._customKeyEventHandler && this._customKeyEventHandler(e) === false) && (Tc(e) || this.focus(), this.updateCursorStyle(e), this._keyPressHandled = false);
  }
  _keyPress(e) {
    let i8;
    if (this._keyPressHandled = false, this._keyDownHandled || this._customKeyEventHandler && this._customKeyEventHandler(e) === false) return false;
    if (this.cancel(e), e.charCode) i8 = e.charCode;
    else if (e.which === null || e.which === void 0) i8 = e.keyCode;
    else if (e.which !== 0 && e.charCode !== 0) i8 = e.which;
    else return false;
    return !i8 || (e.altKey || e.ctrlKey || e.metaKey) && !this._isThirdLevelShift(this.browser, e) ? false : (i8 = String.fromCharCode(i8), this._onKey.fire({ key: i8, domEvent: e }), this._showCursor(), this.coreService.triggerDataEvent(i8, true), this._keyPressHandled = true, this._unprocessedDeadKey = false, true);
  }
  _inputEvent(e) {
    if (e.data && e.inputType === "insertText" && (!e.composed || !this._keyDownSeen) && !this.optionsService.rawOptions.screenReaderMode) {
      if (this._keyPressHandled) return false;
      this._unprocessedDeadKey = false;
      let i8 = e.data;
      return this.coreService.triggerDataEvent(i8, true), this.cancel(e), true;
    }
    return false;
  }
  resize(e, i8) {
    if (e === this.cols && i8 === this.rows) {
      this._charSizeService && !this._charSizeService.hasValidSize && this._charSizeService.measure();
      return;
    }
    super.resize(e, i8);
  }
  _afterResize(e, i8) {
    this._charSizeService?.measure();
  }
  clear() {
    if (!(this.buffer.ybase === 0 && this.buffer.y === 0)) {
      this.buffer.clearAllMarkers(), this.buffer.lines.set(0, this.buffer.lines.get(this.buffer.ybase + this.buffer.y)), this.buffer.lines.length = 1, this.buffer.ydisp = 0, this.buffer.ybase = 0, this.buffer.y = 0;
      for (let e = 1; e < this.rows; e++) this.buffer.lines.push(this.buffer.getBlankLine(X4));
      this._onScroll.fire({ position: this.buffer.ydisp }), this.refresh(0, this.rows - 1);
    }
  }
  reset() {
    this.options.rows = this.rows, this.options.cols = this.cols;
    let e = this._customKeyEventHandler;
    this._setup(), super.reset(), this._selectionService?.reset(), this._decorationService.reset(), this._customKeyEventHandler = e, this.refresh(0, this.rows - 1);
  }
  clearTextureAtlas() {
    this._renderService?.clearTextureAtlas();
  }
  _reportFocus() {
    this.element?.classList.contains("focus") ? this.coreService.triggerDataEvent(b.ESC + "[I") : this.coreService.triggerDataEvent(b.ESC + "[O");
  }
  _reportWindowsOptions(e) {
    if (this._renderService) switch (e) {
      case 0:
        let i8 = this._renderService.dimensions.css.canvas.width.toFixed(0), r5 = this._renderService.dimensions.css.canvas.height.toFixed(0);
        this.coreService.triggerDataEvent(`${b.ESC}[4;${r5};${i8}t`);
        break;
      case 1:
        let n = this._renderService.dimensions.css.cell.width.toFixed(0), o2 = this._renderService.dimensions.css.cell.height.toFixed(0);
        this.coreService.triggerDataEvent(`${b.ESC}[6;${o2};${n}t`);
        break;
    }
  }
  cancel(e, i8) {
    if (!(!this.options.cancelEvents && !i8)) return e.preventDefault(), e.stopPropagation(), false;
  }
};
function Tc(s15) {
  return s15.keyCode === 16 || s15.keyCode === 17 || s15.keyCode === 18;
}
var xn2 = class {
  constructor() {
    this._addons = [];
  }
  dispose() {
    for (let t = this._addons.length - 1; t >= 0; t--) this._addons[t].instance.dispose();
  }
  loadAddon(t, e) {
    let i8 = { instance: e, dispose: e.dispose, isDisposed: false };
    this._addons.push(i8), e.dispose = () => this._wrappedAddonDispose(i8), e.activate(t);
  }
  _wrappedAddonDispose(t) {
    if (t.isDisposed) return;
    let e = -1;
    for (let i8 = 0; i8 < this._addons.length; i8++) if (this._addons[i8] === t) {
      e = i8;
      break;
    }
    if (e === -1) throw new Error("Could not dispose an addon that has not been loaded");
    t.isDisposed = true, t.dispose.apply(t.instance), this._addons.splice(e, 1);
  }
};
var wn2 = class {
  constructor(t) {
    this._line = t;
  }
  get isWrapped() {
    return this._line.isWrapped;
  }
  get length() {
    return this._line.length;
  }
  getCell(t, e) {
    if (!(t < 0 || t >= this._line.length)) return e ? (this._line.loadCell(t, e), e) : this._line.loadCell(t, new q());
  }
  translateToString(t, e, i8) {
    return this._line.translateToString(t, e, i8);
  }
};
var Ji = class {
  constructor(t, e) {
    this._buffer = t;
    this.type = e;
  }
  init(t) {
    return this._buffer = t, this;
  }
  get cursorY() {
    return this._buffer.y;
  }
  get cursorX() {
    return this._buffer.x;
  }
  get viewportY() {
    return this._buffer.ydisp;
  }
  get baseY() {
    return this._buffer.ybase;
  }
  get length() {
    return this._buffer.lines.length;
  }
  getLine(t) {
    let e = this._buffer.lines.get(t);
    if (e) return new wn2(e);
  }
  getNullCell() {
    return new q();
  }
};
var Dn2 = class extends D2 {
  constructor(e) {
    super();
    this._core = e;
    this._onBufferChange = this._register(new v());
    this.onBufferChange = this._onBufferChange.event;
    this._normal = new Ji(this._core.buffers.normal, "normal"), this._alternate = new Ji(this._core.buffers.alt, "alternate"), this._core.buffers.onBufferActivate(() => this._onBufferChange.fire(this.active));
  }
  get active() {
    if (this._core.buffers.active === this._core.buffers.normal) return this.normal;
    if (this._core.buffers.active === this._core.buffers.alt) return this.alternate;
    throw new Error("Active buffer is neither normal nor alternate");
  }
  get normal() {
    return this._normal.init(this._core.buffers.normal);
  }
  get alternate() {
    return this._alternate.init(this._core.buffers.alt);
  }
};
var Rn = class {
  constructor(t) {
    this._core = t;
  }
  registerCsiHandler(t, e) {
    return this._core.registerCsiHandler(t, (i8) => e(i8.toArray()));
  }
  addCsiHandler(t, e) {
    return this.registerCsiHandler(t, e);
  }
  registerDcsHandler(t, e) {
    return this._core.registerDcsHandler(t, (i8, r5) => e(i8, r5.toArray()));
  }
  addDcsHandler(t, e) {
    return this.registerDcsHandler(t, e);
  }
  registerEscHandler(t, e) {
    return this._core.registerEscHandler(t, e);
  }
  addEscHandler(t, e) {
    return this.registerEscHandler(t, e);
  }
  registerOscHandler(t, e) {
    return this._core.registerOscHandler(t, e);
  }
  addOscHandler(t, e) {
    return this.registerOscHandler(t, e);
  }
};
var Ln2 = class {
  constructor(t) {
    this._core = t;
  }
  register(t) {
    this._core.unicodeService.register(t);
  }
  get versions() {
    return this._core.unicodeService.versions;
  }
  get activeVersion() {
    return this._core.unicodeService.activeVersion;
  }
  set activeVersion(t) {
    this._core.unicodeService.activeVersion = t;
  }
};
var Ic = ["cols", "rows"];
var Ue2 = 0;
var Dl = class extends D2 {
  constructor(t) {
    super(), this._core = this._register(new yn2(t)), this._addonManager = this._register(new xn2()), this._publicOptions = { ...this._core.options };
    let e = (r5) => this._core.options[r5], i8 = (r5, n) => {
      this._checkReadonlyOptions(r5), this._core.options[r5] = n;
    };
    for (let r5 in this._core.options) {
      let n = { get: e.bind(this, r5), set: i8.bind(this, r5) };
      Object.defineProperty(this._publicOptions, r5, n);
    }
  }
  _checkReadonlyOptions(t) {
    if (Ic.includes(t)) throw new Error(`Option "${t}" can only be set in the constructor`);
  }
  _checkProposedApi() {
    if (!this._core.optionsService.rawOptions.allowProposedApi) throw new Error("You must set the allowProposedApi option to true to use proposed API");
  }
  get onBell() {
    return this._core.onBell;
  }
  get onBinary() {
    return this._core.onBinary;
  }
  get onCursorMove() {
    return this._core.onCursorMove;
  }
  get onData() {
    return this._core.onData;
  }
  get onKey() {
    return this._core.onKey;
  }
  get onLineFeed() {
    return this._core.onLineFeed;
  }
  get onRender() {
    return this._core.onRender;
  }
  get onResize() {
    return this._core.onResize;
  }
  get onScroll() {
    return this._core.onScroll;
  }
  get onSelectionChange() {
    return this._core.onSelectionChange;
  }
  get onTitleChange() {
    return this._core.onTitleChange;
  }
  get onWriteParsed() {
    return this._core.onWriteParsed;
  }
  get element() {
    return this._core.element;
  }
  get parser() {
    return this._parser || (this._parser = new Rn(this._core)), this._parser;
  }
  get unicode() {
    return this._checkProposedApi(), new Ln2(this._core);
  }
  get textarea() {
    return this._core.textarea;
  }
  get rows() {
    return this._core.rows;
  }
  get cols() {
    return this._core.cols;
  }
  get buffer() {
    return this._buffer || (this._buffer = this._register(new Dn2(this._core))), this._buffer;
  }
  get markers() {
    return this._checkProposedApi(), this._core.markers;
  }
  get modes() {
    let t = this._core.coreService.decPrivateModes, e = "none";
    switch (this._core.coreMouseService.activeProtocol) {
      case "X10":
        e = "x10";
        break;
      case "VT200":
        e = "vt200";
        break;
      case "DRAG":
        e = "drag";
        break;
      case "ANY":
        e = "any";
        break;
    }
    return { applicationCursorKeysMode: t.applicationCursorKeys, applicationKeypadMode: t.applicationKeypad, bracketedPasteMode: t.bracketedPasteMode, insertMode: this._core.coreService.modes.insertMode, mouseTrackingMode: e, originMode: t.origin, reverseWraparoundMode: t.reverseWraparound, sendFocusMode: t.sendFocus, synchronizedOutputMode: t.synchronizedOutput, wraparoundMode: t.wraparound };
  }
  get options() {
    return this._publicOptions;
  }
  set options(t) {
    for (let e in t) this._publicOptions[e] = t[e];
  }
  blur() {
    this._core.blur();
  }
  focus() {
    this._core.focus();
  }
  input(t, e = true) {
    this._core.input(t, e);
  }
  resize(t, e) {
    this._verifyIntegers(t, e), this._core.resize(t, e);
  }
  open(t) {
    this._core.open(t);
  }
  attachCustomKeyEventHandler(t) {
    this._core.attachCustomKeyEventHandler(t);
  }
  attachCustomWheelEventHandler(t) {
    this._core.attachCustomWheelEventHandler(t);
  }
  registerLinkProvider(t) {
    return this._core.registerLinkProvider(t);
  }
  registerCharacterJoiner(t) {
    return this._checkProposedApi(), this._core.registerCharacterJoiner(t);
  }
  deregisterCharacterJoiner(t) {
    this._checkProposedApi(), this._core.deregisterCharacterJoiner(t);
  }
  registerMarker(t = 0) {
    return this._verifyIntegers(t), this._core.registerMarker(t);
  }
  registerDecoration(t) {
    return this._checkProposedApi(), this._verifyPositiveIntegers(t.x ?? 0, t.width ?? 0, t.height ?? 0), this._core.registerDecoration(t);
  }
  hasSelection() {
    return this._core.hasSelection();
  }
  select(t, e, i8) {
    this._verifyIntegers(t, e, i8), this._core.select(t, e, i8);
  }
  getSelection() {
    return this._core.getSelection();
  }
  getSelectionPosition() {
    return this._core.getSelectionPosition();
  }
  clearSelection() {
    this._core.clearSelection();
  }
  selectAll() {
    this._core.selectAll();
  }
  selectLines(t, e) {
    this._verifyIntegers(t, e), this._core.selectLines(t, e);
  }
  dispose() {
    super.dispose();
  }
  scrollLines(t) {
    this._verifyIntegers(t), this._core.scrollLines(t);
  }
  scrollPages(t) {
    this._verifyIntegers(t), this._core.scrollPages(t);
  }
  scrollToTop() {
    this._core.scrollToTop();
  }
  scrollToBottom() {
    this._core.scrollToBottom();
  }
  scrollToLine(t) {
    this._verifyIntegers(t), this._core.scrollToLine(t);
  }
  clear() {
    this._core.clear();
  }
  write(t, e) {
    this._core.write(t, e);
  }
  writeln(t, e) {
    this._core.write(t), this._core.write(`\r
`, e);
  }
  paste(t) {
    this._core.paste(t);
  }
  refresh(t, e) {
    this._verifyIntegers(t, e), this._core.refresh(t, e);
  }
  reset() {
    this._core.reset();
  }
  clearTextureAtlas() {
    this._core.clearTextureAtlas();
  }
  loadAddon(t) {
    this._addonManager.loadAddon(this, t);
  }
  static get strings() {
    return { get promptLabel() {
      return mi2.get();
    }, set promptLabel(t) {
      mi2.set(t);
    }, get tooMuchOutput() {
      return _i2.get();
    }, set tooMuchOutput(t) {
      _i2.set(t);
    } };
  }
  _verifyIntegers(...t) {
    for (Ue2 of t) if (Ue2 === 1 / 0 || isNaN(Ue2) || Ue2 % 1 !== 0) throw new Error("This API only accepts integers");
  }
  _verifyPositiveIntegers(...t) {
    for (Ue2 of t) if (Ue2 && (Ue2 === 1 / 0 || isNaN(Ue2) || Ue2 % 1 !== 0 || Ue2 < 0)) throw new Error("This API only accepts positive integers");
  }
};

// closed.ts
var closeDetached = 4409;
function closedView(code, reason) {
  if (code === closeDetached) {
    return { useHere: true, text: "Shown in another terminal" };
  }
  return {
    useHere: false,
    text: `${reason || "connection lost"}: reload to reconnect`
  };
}

// keys.ts
function handleKey(ev, send) {
  if (ev.key !== "Enter" || !ev.shiftKey && !ev.ctrlKey) return true;
  ev.preventDefault();
  if (ev.type === "keydown") send("\n");
  return false;
}

// entry.ts
var fontFamily = '"Session Icons", "MonoLisaCode", "FiraCode Nerd Font", monospace';
var encoder = new TextEncoder();
async function main() {
  await document.fonts.load('13px "Session Icons"', "").catch(() => []);
  const term = new Dl({
    fontFamily,
    fontSize: 13,
    allowProposedApi: true,
    macOptionIsMeta: true,
    macOptionClickForcesSelection: true
  });
  const fit = new o();
  term.loadAddon(fit);
  term.loadAddon(new Ke());
  term.unicode.activeVersion = "11";
  term.open(document.getElementById("term"));
  try {
    const webgl = new xr();
    webgl.onContextLoss(() => webgl.dispose());
    term.loadAddon(webgl);
  } catch {
  }
  fit.fit();
  term.focus();
  const ws3 = new WebSocket(
    `ws://${location.host}/api/term/ws${location.search}`
  );
  ws3.binaryType = "arraybuffer";
  const send = (data) => {
    if (ws3.readyState === WebSocket.OPEN) ws3.send(encoder.encode(data));
  };
  const sendSize = () => {
    if (ws3.readyState === WebSocket.OPEN) {
      ws3.send(
        JSON.stringify({ type: "resize", cols: term.cols, rows: term.rows })
      );
    }
  };
  term.attachCustomKeyEventHandler((ev) => {
    if (ev.type === "keydown" && ev.key === "Alt") {
      term.options.macOptionIsMeta = ev.code === "AltLeft";
    }
    return handleKey(ev, send);
  });
  term.onData(send);
  term.onBinary((data) => {
    if (ws3.readyState === WebSocket.OPEN) {
      ws3.send(Uint8Array.from(data, (c) => c.charCodeAt(0)));
    }
  });
  term.onResize(sendSize);
  term.onTitleChange((title) => {
    document.title = title;
  });
  term.onSelectionChange(() => {
    const text = term.getSelection();
    if (text) void navigator.clipboard.writeText(text).catch(() => {
    });
  });
  term.parser.registerOscHandler(52, (data) => {
    const b64 = data.slice(data.indexOf(";") + 1);
    if (b64 !== "?") {
      const bytes = Uint8Array.from(atob(b64), (c) => c.charCodeAt(0));
      const text = new TextDecoder().decode(bytes);
      void navigator.clipboard.writeText(text).catch(() => {
      });
    }
    return true;
  });
  new ResizeObserver(() => fit.fit()).observe(document.body);
  ws3.onopen = sendSize;
  ws3.onmessage = (ev) => term.write(new Uint8Array(ev.data));
  ws3.onclose = (ev) => {
    const view = closedView(ev.code, ev.reason);
    if (!view.useHere) {
      term.write(`\r
[${view.text}]\r
`);
      return;
    }
    const button = document.createElement("button");
    button.textContent = "Use here";
    button.onclick = () => location.reload();
    const note = document.createElement("div");
    note.id = "detached";
    note.append(view.text, button);
    document.body.append(note);
    button.focus();
  };
}
void main();
/*! Bundled license information:

@xterm/addon-fit/lib/addon-fit.mjs:
@xterm/addon-unicode11/lib/addon-unicode11.mjs:
@xterm/addon-webgl/lib/addon-webgl.mjs:
@xterm/xterm/lib/xterm.mjs:
  (**
   * Copyright (c) 2014-2024 The xterm.js authors. All rights reserved.
   * @license MIT
   *
   * Copyright (c) 2012-2013, Christopher Jeffrey (MIT License)
   * @license MIT
   *
   * Originally forked from (with the author's permission):
   *   Fabrice Bellard's javascript vt100 for jslinux:
   *   http://bellard.org/jslinux/
   *   Copyright (c) 2011 Fabrice Bellard
   *)
*/
