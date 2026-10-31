// The linter for internal/proj/pages/web/: the family's base
// (build/eslint.base.mjs, copied from tools-common by `npm run page-build`).

import base from './build/eslint.base.mjs';

export default base({
  root: import.meta.dirname,
  ignores: ['kit.d.ts'],
});
