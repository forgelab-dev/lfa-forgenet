import { expect, test, type Page } from '@playwright/test';

/**
 * Le parcours d'un administrateur dans le Marketplace (#116), sur le vrai
 * serveur et un vrai démon Docker : installer le template d'essai « secret »,
 * noter le mot de passe généré une seule fois, et vérifier qu'il ne revient
 * nulle part ensuite.
 */

const env = (name: string): string => {
  const value = process.env[name];
  if (!value) throw new Error(`${name} absent : le globalSetup (e2e/server.ts) n'a pas tourné`);
  return value;
};

test.use({ baseURL: process.env.FORGENET_E2E_URL });

async function logIn(page: Page): Promise<void> {
  await page.goto('/');
  await expect(page.getByRole('heading', { name: 'Connexion' })).toBeVisible();
  await page.getByLabel('Identifiant').fill(env('FORGENET_E2E_USER'));
  await page.getByLabel('Mot de passe', { exact: true }).fill(env('FORGENET_E2E_PASSWORD'));
  await page.getByRole('button', { name: 'Se connecter' }).click();
  await expect(page.getByRole('navigation', { name: 'Navigation principale' })).toBeVisible();
}

test('installer un template et ne voir son mot de passe qu’une fois', async ({
  page,
}, testInfo) => {
  const project = env('FORGENET_E2E_PROJECT');
  const shot = (name: string) =>
    page.screenshot({ path: testInfo.outputPath(`${name}.png`), fullPage: true });

  await logIn(page);

  // Le catalogue, depuis la navigation.
  await page
    .getByRole('navigation', { name: 'Navigation principale' })
    .getByRole('link', { name: 'Marketplace' })
    .click();
  const card = page
    .getByRole('list', { name: 'Templates' })
    .getByRole('link', { name: /Secret de bout en bout/ });
  await expect(card).toBeVisible();
  await shot('1-catalogue');

  // La page du template : son contenu et le formulaire.
  await card.click();
  await expect(page.getByRole('heading', { name: 'Ce qui sera installé' })).toBeVisible();
  // Libellé repris tel quel du manifeste (testdata/marketplace/secret).
  const secretField = page.getByLabel("Secret de l'application", { exact: true });
  await expect(secretField).toHaveValue('');
  await expect(
    page.getByText('Laissez vide : un mot de passe de 24 caractères sera généré.')
  ).toBeVisible();

  // Un nom invalide est refusé dans le navigateur, sans appel au serveur.
  const name = page.getByLabel('Nom de l’application');
  let installCalls = 0;
  page.on('request', (r) => {
    if (r.url().includes('InstallTemplate')) installCalls++;
  });
  await name.fill('Mon App');
  const install = page.getByRole('button', { name: 'Installer Secret de bout en bout' });
  await install.click();
  await expect(
    page.getByText('Minuscules, chiffres, « - » et « _ » seulement', { exact: false })
  ).toBeVisible();
  await expect(name).toHaveAttribute('aria-invalid', 'true');
  expect(installCalls).toBe(0);
  await shot('2-nom-refuse');

  // L'installation : le mot de passe tiré s'affiche une fois, masqué d'abord.
  await name.fill(project);
  await install.click();
  const secrets = page.getByRole('region', { name: 'Mots de passe générés' });
  await expect(secrets).toBeVisible();
  const value = secrets.locator('code');
  await expect(value).toHaveText('•'.repeat(24));
  await secrets.getByRole('button', { name: /^Afficher / }).click();
  await expect(value).not.toHaveText('•'.repeat(24));
  const secret = (await value.textContent()) ?? '';
  expect(secret).toHaveLength(24);

  // Le badge « Réussi » prouve que le secret est arrivé dans le conteneur :
  // la sonde de santé du template échoue sans lui.
  // Le premier « Réussi » est le badge ; le second, la dernière étape de la frise.
  const installation = page.getByRole('region', { name: project });
  await expect(installation.getByText('Réussi', { exact: true }).first()).toBeVisible({
    timeout: 150_000,
  });
  await shot('3-installation-reussie');

  // « C'est noté » retire la valeur de la page.
  await secrets.getByRole('button', { name: 'C’est noté, effacer de la page' }).click();
  await expect(
    page.getByText('Mots de passe effacés de cette page.', { exact: false })
  ).toBeVisible();
  expect(await page.content()).not.toContain(secret);

  // Déploiements : l'application est rattachée au template.
  await page.getByRole('link', { name: 'Voir l’application dans Déploiements' }).click();
  await expect(page.getByText('Marketplace : secret', { exact: false })).toBeVisible();
  await expect(page.getByText(project, { exact: true }).first()).toBeVisible();
  await shot('4-deploiements');
  expect(await page.content()).not.toContain(secret);

  // Après rechargement, la page du template ne le montre pas non plus.
  await page.goto('/marketplace/secret');
  await expect(page.getByRole('heading', { name: 'Ce qui sera installé' })).toBeVisible();
  expect(await page.content()).not.toContain(secret);

  // Une deuxième installation sous le même nom est refusée sous le champ.
  await name.fill(project);
  await install.click();
  await expect(name).toHaveAttribute('aria-invalid', 'true');
  await expect(page.getByText(/Une application porte déjà ce nom/)).toBeVisible();
  await shot('5-nom-deja-pris');
});
