import { LogOut, Menu, X } from 'lucide-react';
import { useEffect, useId, useRef, useState, type ReactNode } from 'react';
import { Link, useLocation } from 'wouter';

import { Wordmark } from '../components/ui/Wordmark';
import { UserRole } from '../gen/forgenet/v1/auth_pb';
import { useCurrentUser, useSession } from '../features/auth/session';
import { navigation, type Section, type SectionGroup } from './navigation';

const roleLabel: Record<UserRole, string> = {
  [UserRole.ADMIN]: 'Administrateur',
  [UserRole.MEMBER]: 'Membre',
  [UserRole.UNSPECIFIED]: 'Rôle inconnu',
};

/** Cadre de l'espace connecté : navigation à gauche, ou en tiroir sur petit écran. */
export function AppShell({ children }: { children: ReactNode }) {
  const [drawerOpen, setDrawerOpen] = useState(false);

  return (
    <div className="min-h-dvh lg:grid lg:grid-cols-[16rem_minmax(0,1fr)]">
      <aside className="sticky top-0 hidden h-dvh border-r border-iron-800 bg-iron-900 lg:block">
        <SidebarContent />
      </aside>

      <header className="sticky top-0 z-10 flex h-14 items-center gap-2 border-b border-iron-800 bg-iron-900/95 px-2 backdrop-blur-sm lg:hidden">
        <button
          type="button"
          onClick={() => setDrawerOpen(true)}
          aria-label="Ouvrir la navigation"
          aria-expanded={drawerOpen}
          className="flex size-11 items-center justify-center rounded-md text-iron-100 hover:bg-iron-850"
        >
          <Menu aria-hidden="true" className="size-5" />
        </button>
        <Wordmark className="text-lg" />
      </header>

      <Drawer open={drawerOpen} onClose={() => setDrawerOpen(false)} />

      <main className="mx-auto w-full max-w-6xl px-4 py-8 sm:px-8 lg:px-12 lg:py-12">
        {children}
      </main>
    </div>
  );
}

/**
 * Tiroir de navigation des petits écrans, sur <dialog> natif : piège du focus,
 * touche Échap et fond inerte sont gérés par le navigateur.
 */
function Drawer({ open, onClose }: { open: boolean; onClose: () => void }) {
  const ref = useRef<HTMLDialogElement>(null);

  useEffect(() => {
    const dialog = ref.current;
    if (!dialog) {
      return;
    }
    if (open && !dialog.open) {
      // jsdom n'implémente pas showModal : l'attribut suffit aux tests.
      if (typeof dialog.showModal === 'function') {
        dialog.showModal();
      } else {
        dialog.setAttribute('open', '');
      }
    } else if (!open && dialog.open) {
      dialog.close?.();
      dialog.removeAttribute('open');
    }
  }, [open]);

  return (
    <dialog
      ref={ref}
      aria-label="Navigation"
      onClose={onClose}
      onClick={(event) => {
        if (event.target === event.currentTarget) {
          onClose();
        }
      }}
      className="m-0 h-dvh max-h-none w-[min(18rem,85vw)] max-w-none border-r border-iron-800 bg-iron-900 p-0 text-iron-100 backdrop:bg-iron-950/75 lg:hidden"
    >
      {open && (
        <div className="relative h-full">
          <button
            type="button"
            onClick={onClose}
            aria-label="Fermer la navigation"
            className="absolute top-2 right-2 z-10 flex size-11 items-center justify-center rounded-md text-iron-300 hover:bg-iron-850 hover:text-iron-100"
          >
            <X aria-hidden="true" className="size-5" />
          </button>
          <SidebarContent onNavigate={onClose} />
        </div>
      )}
    </dialog>
  );
}

function SidebarContent({ onNavigate }: { onNavigate?: () => void }) {
  return (
    <div className="flex h-full flex-col">
      <div className="px-5 pt-6 pb-4">
        <Wordmark className="text-lg" />
        <p className="mt-1.5 truncate font-mono text-xs text-iron-300">{window.location.host}</p>
      </div>

      <nav aria-label="Navigation principale" className="flex-1 overflow-y-auto px-3 pb-6">
        {navigation.map((group) => (
          <NavGroup key={group.label} group={group} onNavigate={onNavigate} />
        ))}
      </nav>

      <UserCard />
    </div>
  );
}

function NavGroup({ group, onNavigate }: { group: SectionGroup; onNavigate?: () => void }) {
  // La barre latérale et le tiroir coexistent dans le DOM : l'identifiant du
  // libellé doit donc être unique par rendu, pas tiré du nom du groupe.
  const labelId = useId();
  return (
    <div className="mt-5 first:mt-2">
      <p id={labelId} className="px-2 pb-1.5 text-xs font-medium text-iron-300">
        {group.label}
      </p>
      <ul aria-labelledby={labelId} className="space-y-0.5">
        {group.sections.map((section) => (
          <li key={section.path}>
            <NavLink section={section} onNavigate={onNavigate} />
          </li>
        ))}
      </ul>
    </div>
  );
}

function NavLink({ section, onNavigate }: { section: Section; onNavigate?: () => void }) {
  const [location] = useLocation();
  const active = section.path === '/' ? location === '/' : location.startsWith(section.path);
  const Icon = section.icon;

  return (
    <Link
      href={section.path}
      onClick={onNavigate}
      aria-current={active ? 'page' : undefined}
      className={
        'group flex min-h-10 items-center gap-3 rounded-md px-2 text-sm transition-colors duration-150 ' +
        (active
          ? 'bg-gold-950 font-semibold text-iron-100'
          : 'text-iron-300 hover:bg-iron-850 hover:text-iron-100')
      }
    >
      <Icon
        aria-hidden="true"
        className={`size-4 shrink-0 ${active ? 'text-gold-400' : 'text-iron-500 group-hover:text-iron-300'}`}
      />
      <span className="flex-1">{section.label}</span>
      {section.upcoming && (
        <span className="font-mono text-[0.6875rem] text-iron-300/80">bientôt</span>
      )}
    </Link>
  );
}

function UserCard() {
  const user = useCurrentUser();
  const { logout } = useSession();
  const [leaving, setLeaving] = useState(false);

  return (
    <div className="border-t border-iron-800 px-3 py-3">
      <div className="flex items-center gap-3 px-2 py-1">
        <span
          aria-hidden="true"
          className="flex size-9 shrink-0 items-center justify-center rounded-md border border-gold-400/50 bg-gold-950 font-semibold text-gold-300 uppercase"
        >
          {user.username.charAt(0)}
        </span>
        <div className="min-w-0">
          <p className="truncate text-sm font-semibold">{user.username}</p>
          <p className="text-xs text-iron-300">{roleLabel[user.role]}</p>
        </div>
      </div>
      <button
        type="button"
        disabled={leaving}
        onClick={() => {
          setLeaving(true);
          void logout();
        }}
        className="mt-2 flex min-h-11 w-full items-center gap-3 rounded-md px-2 text-sm text-iron-300 transition-colors hover:bg-iron-850 hover:text-iron-100 disabled:opacity-60"
      >
        <LogOut aria-hidden="true" className="size-4" />
        {leaving ? 'Déconnexion…' : 'Se déconnecter'}
      </button>
    </div>
  );
}
