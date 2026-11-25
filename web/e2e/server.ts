/**
 * Lance un vrai forgenet-server pour les tests de bout en bout du dashboard.
 *
 * Le binaire est compilé dans un dossier temporaire ; le mot de passe de
 * l'administrateur, la clé des secrets et le secret JWT sont tirés au hasard
 * à chaque exécution et ne sont jamais affichés. Le catalogue est le template
 * d'essai partagé avec le test Go (testdata/marketplace). La fonction
 * renvoyée sert de globalTeardown : elle arrête le serveur par son PID et
 * retire la stack Docker installée par le test.
 */
import { execFile, spawn } from 'node:child_process';
import { randomBytes } from 'node:crypto';
import { closeSync, openSync } from 'node:fs';
import { mkdtemp, readFile, rm, writeFile } from 'node:fs/promises';
import { createServer } from 'node:net';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { promisify } from 'node:util';

const run = promisify(execFile);
const repo = resolve(import.meta.dirname, '..', '..');

function freePort(): Promise<number> {
  return new Promise((done, fail) => {
    const srv = createServer();
    srv.once('error', fail);
    srv.listen(0, '127.0.0.1', () => {
      const address = srv.address();
      srv.close(() =>
        typeof address === 'object' && address ? done(address.port) : fail(new Error('port'))
      );
    });
  });
}

async function waitHealthy(url: string, log: string): Promise<void> {
  const deadline = Date.now() + 30_000;
  while (Date.now() < deadline) {
    try {
      if ((await fetch(`${url}/healthz`)).ok) return;
    } catch {
      // Pas encore à l'écoute.
    }
    await new Promise((r) => setTimeout(r, 200));
  }
  throw new Error(
    `forgenet-server ne répond pas sur ${url}/healthz :\n${await readFile(log, 'utf8')}`
  );
}

export default async function globalSetup(): Promise<() => Promise<void>> {
  const dir = await mkdtemp(join(tmpdir(), 'forgenet-e2e-'));
  const bin = join(dir, 'forgenet-server');
  await run('go', ['build', '-o', bin, './cmd/forgenet-server'], { cwd: repo });

  const [httpPort, grpcPort] = [await freePort(), await freePort()];
  const config = join(dir, 'forgenet.yaml');
  await writeFile(
    config,
    [
      'environment: test',
      'server:',
      `  http_port: ${httpPort}`,
      `  grpc_port: ${grpcPort}`,
      '  shutdown_timeout: 5s',
      'database:',
      `  path: ${JSON.stringify(join(dir, 'forgenet.db'))}`,
      'auth:',
      `  jwt_secret: ${JSON.stringify(randomBytes(32).toString('hex'))}`,
      'secrets:',
      `  key: ${JSON.stringify(randomBytes(32).toString('hex'))}`,
      'marketplace:',
      `  dir: ${JSON.stringify(join(repo, 'testdata', 'marketplace'))}`,
      'web:',
      `  root: ${JSON.stringify(join(repo, 'web', 'dist'))}`,
      '',
    ].join('\n'),
    { mode: 0o600 }
  );

  // Le mot de passe passe par l'entrée standard, jamais par la ligne de commande.
  const password = randomBytes(18).toString('base64url');
  await new Promise<void>((done, fail) => {
    const create = spawn(
      bin,
      [
        'user',
        'create',
        '--config',
        config,
        '--username',
        'admin',
        '--email',
        'admin@forgenet.test',
        '--role',
        'admin',
      ],
      { cwd: dir, stdio: ['pipe', 'ignore', 'pipe'] }
    );
    let stderr = '';
    create.stderr.on('data', (chunk) => (stderr += chunk));
    create.on('error', fail);
    create.on('exit', (code) =>
      code === 0 ? done() : fail(new Error(`user create a échoué (${code}) : ${stderr}`))
    );
    create.stdin.end(`${password}\n`);
  });

  const log = join(dir, 'server.log');
  const out = openSync(log, 'a');
  const server = spawn(bin, ['--config', config], { cwd: dir, stdio: ['ignore', out, out] });
  closeSync(out);
  const url = `http://127.0.0.1:${httpPort}`;
  await waitHealthy(url, log).catch((err: unknown) => {
    server.kill('SIGTERM');
    throw err;
  });

  // Nom unique par exécution : la stack se retrouve et se retire sans risque
  // de toucher une installation existante.
  const project = `e2e-${randomBytes(4).toString('hex')}`;
  process.env.FORGENET_E2E_URL = url;
  process.env.FORGENET_E2E_USER = 'admin';
  process.env.FORGENET_E2E_PASSWORD = password;
  process.env.FORGENET_E2E_PROJECT = project;

  return async () => {
    if (server.pid && server.exitCode === null) {
      process.kill(server.pid, 'SIGTERM');
      await new Promise((r) => server.once('exit', r));
    }
    await run('docker', [
      'compose',
      '-p',
      `forgenet-${project}`,
      'down',
      '--volumes',
      '--remove-orphans',
    ]).catch(() => undefined);
    await rm(dir, { recursive: true, force: true });
  };
}
