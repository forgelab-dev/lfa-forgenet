import { Code, ConnectError } from '@connectrpc/connect';
import { Eye, EyeOff, LoaderCircle, PackagePlus } from 'lucide-react';
import { useId, useRef, useState, type FormEvent, type ReactNode } from 'react';

import { Button } from '../../components/ui/Button';
import { fieldClass } from '../../components/ui/field';
import type {
  InstallTemplateResponse,
  Template,
  TemplateVariable,
} from '../../gen/forgenet/v1/marketplace_pb';
import { describeError } from '../../lib/api/errors';
import { useApi } from '../auth/session';
import { projectIdPattern, toProjectId } from '../deployments/projectId';
import { checkValue, emptyHint, initialValues, introduction, refusal } from './values';

/** Les problèmes du formulaire : par variable, plus le nom sous la clé `projectId`. */
type Problems = Record<string, string | undefined>;

const projectKey = 'projectId';

function checkProject(id: string): string | undefined {
  if (!id) return 'Saisissez le nom de l’application.';
  if (!projectIdPattern.test(id)) {
    return 'Minuscules, chiffres, « - » et « _ » seulement, en commençant par une lettre ou un chiffre (63 caractères au plus).';
  }
  return undefined;
}

/** Clé d'idempotence : un double envoi n'installe qu'une fois. */
function newRequestId(): string {
  return (
    globalThis.crypto?.randomUUID?.() ?? `${Date.now()}-${Math.random().toString(16).slice(2)}`
  );
}

/**
 * Formulaire d'installation d'un template : le nom de l'application, puis un
 * champ par variable du manifeste. Le navigateur vérifie ce qu'il peut ; le
 * serveur vérifie tout à nouveau et tire les mots de passe laissés vides.
 */
export function InstallForm({
  template,
  onInstalled,
}: {
  template: Template;
  onInstalled: (response: InstallTemplateResponse) => void;
}) {
  const api = useApi();
  const [projectId, setProjectId] = useState(() => toProjectId(template.id));
  const [values, setValues] = useState(() => initialValues(template.variables));
  const [problems, setProblems] = useState<Problems>({});
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  // Gardée jusqu'à une réponse : un envoi répété après une coupure ne crée
  // pas une seconde installation.
  const requestId = useRef(newRequestId());
  const ids = { project: useId(), error: useId() };

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const id = projectId.trim();
    const found: Problems = { [projectKey]: checkProject(id) };
    for (const v of template.variables) {
      found[v.name] = checkValue(v, values[v.name] ?? '');
    }
    setProblems(found);
    setError(undefined);
    if (Object.values(found).some(Boolean)) return;

    setPending(true);
    try {
      const response = await api.marketplace.installTemplate({
        templateId: template.id,
        projectId: id,
        // Un champ vide vaut une valeur absente : le serveur applique le
        // défaut ou tire un mot de passe.
        values: Object.fromEntries(Object.entries(values).filter(([, value]) => value !== '')),
        requestId: requestId.current,
      });
      requestId.current = newRequestId();
      onInstalled(response);
    } catch (err) {
      const message = describeError(err);
      const code = ConnectError.from(err).code;
      if (code === Code.AlreadyExists) {
        setProblems({ [projectKey]: message });
      } else if (code === Code.InvalidArgument) {
        const refused = refusal(message, template.variables);
        if (refused.field) setProblems({ [refused.field]: refused.message });
        else setError(refused.message);
      } else {
        setError(message);
      }
    } finally {
      setPending(false);
    }
  }

  return (
    <form
      noValidate
      onSubmit={(e) => void submit(e)}
      aria-labelledby="install-form"
      aria-describedby={error ? ids.error : undefined}
      className="animate-rise rounded-lg border border-iron-800 bg-iron-900 px-5 py-5 [animation-delay:120ms]"
    >
      <h2 id="install-form" className="flex items-center gap-2 text-lg font-semibold">
        <PackagePlus aria-hidden="true" className="size-5 text-gold-400" />
        Installer
      </h2>
      <p className="mt-2 text-sm text-iron-300">{introduction(template)}</p>

      <div className="mt-5">
        <label htmlFor={ids.project} className="text-sm font-medium">
          Nom de l’application
        </label>
        <input
          id={ids.project}
          name="project"
          value={projectId}
          onChange={(e) => setProjectId(e.target.value)}
          autoComplete="off"
          autoCapitalize="none"
          spellCheck={false}
          aria-invalid={problems[projectKey] ? true : undefined}
          aria-describedby={`${ids.project}-hint`}
          className={`${fieldClass} font-mono ${problems[projectKey] ? 'border-danger-400' : 'border-iron-700'}`}
        />
        <Hint id={`${ids.project}-hint`} problem={problems[projectKey]}>
          Le nom de la stack sur le serveur, unique : minuscules, chiffres, « - » et « _ ».
        </Hint>
      </div>

      {template.variables.map((v) => (
        <VariableField
          key={v.name}
          variable={v}
          value={values[v.name] ?? ''}
          problem={problems[v.name]}
          onChange={(value) => setValues((current) => ({ ...current, [v.name]: value }))}
        />
      ))}

      {error && (
        <p
          id={ids.error}
          role="alert"
          className="mt-5 rounded-md border border-danger-400/40 bg-danger-950 px-3 py-2.5 text-sm break-words text-iron-100"
        >
          {error}
        </p>
      )}

      <Button
        type="submit"
        disabled={pending}
        aria-label={`Installer ${template.name}`}
        className="mt-6 w-full"
      >
        {pending && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
        {pending ? 'Envoi…' : `Installer ${template.name}`}
      </Button>
    </form>
  );
}

function VariableField({
  variable: v,
  value,
  problem,
  onChange,
}: {
  variable: TemplateVariable;
  value: string;
  problem?: string;
  onChange: (value: string) => void;
}) {
  const id = useId();
  const [shown, setShown] = useState(false);
  const hint = [v.help, emptyHint(v)].filter(Boolean).join(' ');

  return (
    <div className="mt-4">
      <div className="flex items-baseline justify-between gap-3">
        <label htmlFor={id} className="text-sm font-medium">
          {v.label}
        </label>
        <span className="font-mono text-xs text-iron-500">{v.name}</span>
      </div>
      <div className="relative">
        <input
          id={id}
          name={v.name}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          type={v.secret && !shown ? 'password' : 'text'}
          // Un mot de passe enregistré par le navigateur n'a rien à faire ici.
          autoComplete={v.secret ? 'new-password' : 'off'}
          autoCapitalize="none"
          spellCheck={false}
          aria-invalid={problem ? true : undefined}
          aria-describedby={hint || problem ? `${id}-hint` : undefined}
          className={`${fieldClass} font-mono ${v.secret ? 'pr-12' : ''} ${problem ? 'border-danger-400' : 'border-iron-700'}`}
        />
        {v.secret && (
          <button
            type="button"
            onClick={() => setShown((s) => !s)}
            aria-pressed={shown}
            aria-label={`Afficher ${v.label}`}
            className="absolute right-0 bottom-0 flex size-11 items-center justify-center rounded-md text-iron-300 transition-colors hover:text-iron-100 focus-visible:outline-2 focus-visible:outline-gold-400"
          >
            {shown ? (
              <EyeOff aria-hidden="true" className="size-4" />
            ) : (
              <Eye aria-hidden="true" className="size-4" />
            )}
          </button>
        )}
      </div>
      {(hint || problem) && (
        <Hint id={`${id}-hint`} problem={problem}>
          {hint}
        </Hint>
      )}
    </div>
  );
}

function Hint({ id, problem, children }: { id: string; problem?: string; children: ReactNode }) {
  return (
    <p
      id={id}
      className={`mt-1.5 ${problem ? 'text-sm text-danger-400' : 'text-xs text-iron-300'}`}
    >
      {problem ?? children}
    </p>
  );
}
