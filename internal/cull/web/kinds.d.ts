// The page's names for the wire (wire.d.ts, generated from the Go structs),
// and what Go cannot declare: an item's state and jev, stored as bytes whose
// shape depends on its kind (judge.State or judge.GroupState, and Jev's
// answers to the rubric's questions). Ambient.

type Noul = { type: 'noul'; noul: number };

type Answer = import('./wire.d.ts').Answer;
type RecentSend = import('./wire.d.ts').RecentSend;
type TestState = import('./wire.d.ts').State;
type GroupState = import('./wire.d.ts').GroupState;

interface TestJev {
  verdict: {
    type: 'choice';
    choice: string;
    probabilities: { cut: number; keep: number; review: number };
    confidence: number;
  };
  regression_value: {
    type: 'score';
    score: number;
    probabilities?: Record<string, number>;
    confidence?: number;
  };
  tautological: Noul;
  mock_only: Noul;
  framework_behavior: Noul;
  incidental_detail: Noul;
  not_missed: Noul;
}

interface GroupJev {
  verdict: {
    type: 'choice';
    choice: string;
    probabilities: {
      consolidate: number;
      keep_separate: number;
      review: number;
    };
    confidence: number;
  };
  same_behavior: Noul;
  exact_duplicate: Noul;
  loss_if_merged: Noul;
}

type TestItem = import('./wire.d.ts').ItemBase & {
  kind: 'test';
  state: TestState;
  jev: TestJev;
};

type GroupItem = import('./wire.d.ts').ItemBase & {
  kind: 'group';
  state: GroupState;
  jev: GroupJev;
};

type Item = TestItem | GroupItem;

/** GET /api/review, its items narrowed by kind. */
type Review = Omit<import('./wire.d.ts').Review, 'items'> & { items: Item[] };
