import { create } from '@bufbuild/protobuf';
import { describe, expect, it } from 'vitest';

import { ContainerSchema, ContainerState } from '../../gen/forgenet/v1/container_pb';
import { splitSystem, stoppedSystem } from './system';

const container = (name: string, state: ContainerState, system = false) =>
  create(ContainerSchema, { id: name, name, state, system });

describe('splitSystem', () => {
  it('range les conteneurs de ForgeNet à part', () => {
    const { apps, system } = splitSystem([
      container('boutique-web-1', ContainerState.RUNNING),
      container('forgenet-server', ContainerState.RUNNING, true),
      container('forgenet-updater', ContainerState.EXITED, true),
    ]);
    expect(apps.map((c) => c.name)).toEqual(['boutique-web-1']);
    expect(system.map((c) => c.name)).toEqual(['forgenet-server', 'forgenet-updater']);
  });
});

describe('stoppedSystem', () => {
  it('ignore le conteneur de mise à jour arrêté, signale Caddy arrêté', () => {
    const stopped = stoppedSystem([
      container('forgenet-server', ContainerState.RUNNING, true),
      container('forgenet-caddy', ContainerState.EXITED, true),
      container('forgenet-updater', ContainerState.EXITED, true),
    ]);
    expect(stopped.map((c) => c.name)).toEqual(['forgenet-caddy']);
  });
});
