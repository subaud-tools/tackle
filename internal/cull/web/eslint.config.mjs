// The linter for internal/cull/web/: the family's base (build/eslint.base.mjs,
// copied from tools-common by `npm run page-build`). wire.d.ts is generated
// and checked by its Go test.

import base from './build/eslint.base.mjs';

export default base({
  root: import.meta.dirname,
  ignores: ['kit.d.ts', 'wire.d.ts'],
});
