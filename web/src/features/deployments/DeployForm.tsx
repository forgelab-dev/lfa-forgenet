import { FileUp, LoaderCircle, Rocket } from 'lucide-react';
import { useId, useRef, useState, type ChangeEvent, type FormEvent } from 'react';

import { Button } from '../../components/ui/Button';
import { fieldClass } from '../../components/ui/field';
import type { Deployment } from '../../gen/forgenet/v1/deploy_pb';
import { describeError } from '../../lib/api/errors';
import { useApi } from '../auth/session';
import { byteLength, maxComposeBytes, projectIdPattern, suggestProjectId } from './projectId';

interface Problems {
  projectId?: string;
  compose?: string;
}

function checkForm(projectId: string, compose: string): Problems {
  const problems: Problems = {};
  if (!projectId) {
    problems.projectId = 'Saisissez l’identifiant du projet.';
  } else if (!projectIdPattern.test(projectId)) {
    problems.projectId =
      'Minuscules, chiffres, « - » et « _ » seulement, en commençant par une lettre ou un chiffre (63 caractères au plus).';
  }
  if (!compose.trim()) {
    problems.compose = 'Collez le fichier Compose ou importez-le.';
  } else if (byteLength(compose) > maxComposeBytes) {
    problems.compose =
      'Le fichier dépasse 1 Mio : un fichier Compose pèse d’ordinaire quelques Ko.';
  }
  return problems;
}

/** Clé d'idempotence : un double envoi ne crée qu'un déploiement. */
function newRequestId(): string {
  return (
    globalThis.crypto?.randomUUID?.() ?? `${Date.now()}-${Math.random().toString(16).slice(2)}`
  );
}

/**
 * Déploiement d'une application depuis un fichier Compose collé ou importé.
 * L'identifiant du projet est proposé depuis le fichier tant qu'on ne l'a pas
 * saisi soi-même.
 */
export function DeployForm({ onCreated }: { onCreated: (deployment: Deployment) => void }) {
  const api = useApi();
  const [projectId, setProjectId] = useState('');
  const [compose, setCompose] = useState('');
  const [fileName, setFileName] = useState<string>();
  const [typedId, setTypedId] = useState(false);
  const [problems, setProblems] = useState<Problems>({});
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  const requestId = useRef(newRequestId());
  const fileInput = useRef<HTMLInputElement>(null);
  const ids = { project: useId(), compose: useId(), error: useId(), file: useId() };

  function updateCompose(text: string, name?: string) {
    setCompose(text);
    if (!typedId) {
      setProjectId(suggestProjectId(text, name));
    }
  }

  async function importFile(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    event.target.value = '';
    if (!file) return;
    if (file.size > maxComposeBytes) {
      setProblems((p) => ({
        ...p,
        compose: 'Ce fichier dépasse 1 Mio : est-ce bien un fichier Compose ?',
      }));
      return;
    }
    setFileName(file.name);
    updateCompose(await file.text(), file.name);
    setProblems((p) => ({ ...p, compose: undefined }));
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const id = projectId.trim();
    const found = checkForm(id, compose);
    setProblems(found);
    if (found.projectId || found.compose) {
      setError(undefined);
      return;
    }
    setPending(true);
    setError(undefined);
    try {
      const res = await api.deployments.createDeployment({
        projectId: id,
        requestId: requestId.current,
        source: { source: { case: 'compose', value: { composeContent: compose } } },
      });
      if (res.deployment) {
        onCreated(res.deployment);
      }
      setProjectId('');
      setCompose('');
      setFileName(undefined);
      setTypedId(false);
      requestId.current = newRequestId();
    } catch (err) {
      setError(describeError(err));
    } finally {
      setPending(false);
    }
  }

  const border = (bad?: string) => (bad ? 'border-danger-400' : 'border-iron-700');

  return (
    <form
      noValidate
      onSubmit={submit}
      aria-labelledby="deploy-form"
      aria-describedby={error ? ids.error : undefined}
      className="rounded-lg border border-iron-800 bg-iron-900 px-5 py-5"
    >
      <h2 id="deploy-form" className="flex items-center gap-2 text-lg font-semibold">
        <Rocket aria-hidden="true" className="size-5 text-gold-400" />
        Déployer une application
      </h2>
      <p className="mt-2 text-sm text-iron-300">
        Un fichier Docker Compose par application. Redéployer un projet existant remplace sa version
        en service.
      </p>

      <div className="mt-5">
        <div className="flex items-end justify-between gap-3">
          <label htmlFor={ids.compose} className="text-sm font-medium">
            Fichier Compose
          </label>
          <input
            ref={fileInput}
            id={ids.file}
            type="file"
            accept=".yml,.yaml,application/yaml,text/yaml"
            onChange={(e) => void importFile(e)}
            className="sr-only"
            tabIndex={-1}
            aria-hidden="true"
          />
          <button
            type="button"
            onClick={() => fileInput.current?.click()}
            className="inline-flex min-h-11 items-center gap-1.5 rounded-md px-2 text-sm font-medium text-gold-400 transition-colors hover:bg-iron-850 hover:text-gold-300 focus-visible:outline-2 focus-visible:outline-gold-400"
          >
            <FileUp aria-hidden="true" className="size-4" />
            Importer un fichier
          </button>
        </div>
        <textarea
          id={ids.compose}
          name="compose"
          value={compose}
          onChange={(e) => updateCompose(e.target.value, fileName)}
          spellCheck={false}
          autoCapitalize="none"
          wrap="off"
          rows={12}
          placeholder={'services:\n  web:\n    image: nginx:1.29-alpine'}
          aria-invalid={problems.compose ? true : undefined}
          aria-describedby={problems.compose ? `${ids.compose}-hint` : `${ids.compose}-note`}
          className={`${fieldClass} min-h-56 resize-y py-2.5 font-mono text-sm leading-relaxed ${border(problems.compose)}`}
        />
        {problems.compose ? (
          <p id={`${ids.compose}-hint`} className="mt-1.5 text-sm text-danger-400">
            {problems.compose}
          </p>
        ) : (
          <p id={`${ids.compose}-note`} className="mt-1.5 text-xs text-iron-300">
            {fileName ? `Importé de ${fileName}. ` : ''}Les variables d’environnement du serveur ne
            sont pas transmises : mettez les valeurs dans le fichier.
          </p>
        )}
      </div>

      <div className="mt-4">
        <label htmlFor={ids.project} className="text-sm font-medium">
          Projet
        </label>
        <input
          id={ids.project}
          name="project"
          value={projectId}
          onChange={(e) => {
            setProjectId(e.target.value);
            setTypedId(e.target.value !== '');
          }}
          autoComplete="off"
          autoCapitalize="none"
          spellCheck={false}
          placeholder="boutique"
          aria-invalid={problems.projectId ? true : undefined}
          aria-describedby={`${ids.project}-hint`}
          className={`${fieldClass} font-mono ${border(problems.projectId)}`}
        />
        <p
          id={`${ids.project}-hint`}
          className={`mt-1.5 text-xs ${problems.projectId ? 'text-sm text-danger-400' : 'text-iron-300'}`}
        >
          {problems.projectId ??
            (typedId || !projectId
              ? 'Le nom de la stack sur le serveur : minuscules, chiffres, « - » et « _ ».'
              : 'Proposé depuis le fichier ; vous pouvez le changer.')}
        </p>
      </div>

      {error && (
        <p
          id={ids.error}
          role="alert"
          className="mt-5 rounded-md border border-danger-400/40 bg-danger-950 px-3 py-2.5 text-sm break-words text-iron-100"
        >
          {error}
        </p>
      )}

      <Button type="submit" disabled={pending} className="mt-5 w-full">
        {pending && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
        {pending ? 'Envoi…' : 'Déployer'}
      </Button>
    </form>
  );
}
