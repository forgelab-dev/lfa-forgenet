import '@testing-library/jest-dom/vitest';
import { cleanup } from '@testing-library/react';
import { afterEach } from 'vitest';

// Sans `globals: true`, Testing Library ne démonte pas seul les rendus entre
// deux tests : ils s'accumuleraient dans le même document.
afterEach(cleanup);
