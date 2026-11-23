import { LoaderCircle, RotateCw, ServerOff } from 'lucide-react';
import type { ReactNode } from 'react';
import { Route, Router, Switch, useLocation, type BaseLocationHook } from 'wouter';

import { Button } from '../components/ui/Button';
import { Wordmark } from '../components/ui/Wordmark';
import { AccountPage } from '../features/account/AccountPage';
import { AccessLinkScreen } from '../features/auth/AccessLinkScreen';
import { ForgotPasswordScreen } from '../features/auth/ForgotPasswordScreen';
import { LoginScreen } from '../features/auth/LoginScreen';
import { MarketplacePage } from '../features/marketplace/MarketplacePage';
import { TemplatePage } from '../features/marketplace/TemplatePage';
import { useSession } from '../features/auth/session';
import { SessionProvider } from '../features/auth/SessionProvider';
import { ContainersPage } from '../features/containers/ContainersPage';
import { DeploymentsPage } from '../features/deployments/DeploymentsPage';
import { DomainsPage } from '../features/domains/DomainsPage';
import { MonitoringPage } from '../features/monitoring/MonitoringPage';
import { OverviewPage } from '../features/overview/OverviewPage';
import { SettingsPage } from '../features/settings/SettingsPage';
import { TeamPage } from '../features/team/TeamPage';
import { createApiTransport, type TransportFactory } from '../lib/api/client';
import { AppShell } from './AppShell';
import { sections } from './navigation';
import { NotFoundPage, UpcomingPage } from './pages';

const browserTransport: TransportFactory = (options) => createApiTransport(options);

interface Props {
  /** Transport de l'API ; les tests passent un faux serveur. */
  transportFactory?: TransportFactory;
  /** Source de l'adresse courante ; les tests passent une adresse en mémoire. */
  locationHook?: BaseLocationHook;
}

/** Racine du dashboard : la session décide entre l'écran de connexion et l'espace connecté. */
export function App({ transportFactory = browserTransport, locationHook }: Props) {
  return (
    <SessionProvider transportFactory={transportFactory}>
      <Router {...(locationHook ? { hook: locationHook } : {})}>
        <SessionGate />
      </Router>
    </SessionProvider>
  );
}

function SessionGate() {
  const { state, retry, logout } = useSession();
  const [location] = useLocation();

  // Un lien à usage unique s'ouvre avec ou sans session : il en ouvre une
  // lui-même une fois le mot de passe choisi (ADR-0014).
  if (location === '/acces') {
    return <AccessLinkScreen />;
  }
  if (location === '/mot-de-passe-oublie' && state.status === 'anonymous') {
    return <ForgotPasswordScreen />;
  }

  switch (state.status) {
    case 'checking':
      return (
        <CenteredScreen>
          <p role="status" className="flex items-center gap-2 text-iron-300">
            <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />
            Vérification de la session…
          </p>
        </CenteredScreen>
      );
    case 'unreachable':
      return (
        <CenteredScreen>
          <div role="alert" className="max-w-md text-center">
            <ServerOff aria-hidden="true" className="mx-auto size-8 text-warn-400" />
            <h1 className="mt-4 text-2xl font-semibold">Session non vérifiée</h1>
            <p className="mt-2 text-iron-300">{state.message}</p>
            <div className="mt-6 flex flex-col justify-center gap-3 sm:flex-row">
              <Button onClick={retry}>
                <RotateCw aria-hidden="true" className="size-4" />
                Réessayer
              </Button>
              <Button variant="quiet" onClick={() => void logout()}>
                Changer de compte
              </Button>
            </div>
          </div>
        </CenteredScreen>
      );
    case 'anonymous':
      return <LoginScreen />;
    case 'authenticated':
      return (
        <AppShell>
          <Switch>
            <Route path="/" component={OverviewPage} />
            <Route path="/conteneurs" component={ContainersPage} />
            <Route path="/deploiements" component={DeploymentsPage} />
            <Route path="/surveillance" component={MonitoringPage} />
            <Route path="/marketplace" component={MarketplacePage} />
            <Route path="/marketplace/:id" component={TemplatePage} />
            <Route path="/domaines" component={DomainsPage} />
            <Route path="/compte" component={AccountPage} />
            <Route path="/equipe" component={TeamPage} />
            <Route path="/reglages" component={SettingsPage} />
            {sections
              .filter((section) => section.upcoming)
              .map((section) => (
                <Route key={section.path} path={section.path}>
                  <UpcomingPage section={section} />
                </Route>
              ))}
            <Route component={NotFoundPage} />
          </Switch>
        </AppShell>
      );
  }
}

function CenteredScreen({ children }: { children: ReactNode }) {
  return (
    <div className="flex min-h-dvh flex-col items-center justify-center gap-8 px-6">
      <Wordmark className="text-xl" />
      {children}
    </div>
  );
}
