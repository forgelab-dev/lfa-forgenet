// Flat config ESLint 10 — dashboard web (squelette vide + gen protobuf-es).
// Ignore le code genere : seul le code maintenu est linte.
import tseslint from 'typescript-eslint';

export default tseslint.config(
  {
    ignores: ['node_modules/**', 'dist/**', 'src/gen/**']
  },
  ...tseslint.configs.recommended
);
