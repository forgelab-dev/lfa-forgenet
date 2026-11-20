# ADR-0016 : e-mails en deux versions, texte et HTML

- Statut : accepté
- Date : 2026-11-20
- Contexte : ticket #101 ; revient sur un point de l'ADR-0014 (« Pas de
  HTML »)

## Contexte

L'ADR-0014 a choisi le texte brut : rien à assainir, et mieux reçu par les
filtres. Les e-mails marchent, mais rien ne dit au premier coup d'œil qu'ils
viennent de ForgeNet. Or presque tous demandent de cliquer un lien à usage
unique : un message qui ressemble au produit inspire plus confiance qu'un
texte anonyme, et un bouton se trouve mieux qu'une adresse au milieu du
texte.

## Décision

### Deux versions dans chaque e-mail

- Chaque e-mail part en `multipart/alternative` : le **texte brut**, puis la
  version **HTML**. Une messagerie qui n'affiche pas le HTML, ou un lecteur
  qui le refuse, garde le texte.
- **Un seul modèle, plusieurs contenus.** Chaque e-mail se décrit une fois en
  blocs : salutation, paragraphes, action (le lien), commande, remarque. Les
  deux versions sont rendues à partir de ces blocs : elles ne peuvent pas
  dire deux choses différentes. Le texte est coupé à 76 caractères, sauf les
  liens et les commandes, qui restent sur une ligne pour se copier entiers.
- Le HTML passe par `html/template` : tout ce qui vient d'un compte (nom,
  relais) est échappé. Le modèle n'a pas de champ libre.

### La charte, sans ressource distante

- Fond acier bleui et carte, accent or, comme le dashboard. Les couleurs sont
  les jetons OKLCH du dashboard convertis en hexadécimal, seule écriture que
  les messageries comprennent toutes.
- Tableaux et styles en ligne, qui seuls passent partout (Gmail retire une
  partie des feuilles de style, Outlook pour Windows lit le HTML avec le
  moteur de Word).
- Le **logo** est une image PNG de 96 pixels **jointe au message** et appelée
  par `cid:` (`multipart/related`). Le SVG en ligne n'est pas affiché par
  Gmail ni Outlook, et une image distante ferait savoir au serveur quand et
  d'où le message est lu. Une messagerie qui bloque les images montre le nom
  « ForgeNet » écrit à côté, en texte.
- Le message est déjà sombre et le déclare (`color-scheme`) : les
  messageries en mode sombre n'ont rien à inverser. Les contrastes sont d'au
  moins 7:1 pour le texte et 9:1 pour le bouton.

### Le lien ne sort pas du message

- Le jeton reste dans le fragment de l'adresse (ADR-0014). Le bouton et le
  lien en clair, en dessous, portent exactement la même adresse : pas de
  redirection, pas de paramètre ajouté, pas de pixel de suivi, aucune autre
  adresse dans le HTML.

## Ce qui ne change pas

- Les textes : les mêmes mots qu'avant. Seuls s'ajoutent, dans le HTML,
  l'intitulé du bouton et « ou copiez ce lien ».
- L'envoi, le relais, la file et le journal (ADR-0014, ADR-0015).

## Retour arrière

- Le format du message est calculé à l'envoi, rien n'est stocké : une version
  antérieure de ForgeNet renvoie du texte brut seul.

## Conséquences

- Un message pèse environ 10 Ko de plus, le logo compris.
- Le rendu se vérifie à la main dans les messageries : Go teste la structure
  MIME, la présence du même lien dans les deux versions et l'absence de toute
  autre adresse.
