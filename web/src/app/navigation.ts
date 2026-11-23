import {
  Activity,
  Boxes,
  Database,
  Gauge,
  Globe,
  MessageSquareText,
  Rocket,
  Settings,
  Store,
  Users,
  type LucideIcon,
} from 'lucide-react';

export interface Section {
  path: string;
  label: string;
  icon: LucideIcon;
  /** Une phrase : ce que l'on fait dans cette section. */
  summary: string;
  /** Section prévue mais pas encore construite. */
  upcoming?: boolean;
}

export interface SectionGroup {
  label: string;
  sections: Section[];
}

/**
 * Navigation principale, regroupée par usage. Les sections à venir sont
 * visibles et marquées : l'administrateur voit ce que le produit fera, sans
 * cliquer pour le découvrir. La facturation (Managed Cloud) n'y figure pas :
 * elle n'existe pas dans la Community Edition.
 */
export const navigation: SectionGroup[] = [
  {
    label: 'Exploitation',
    sections: [
      {
        path: '/',
        label: 'Vue d’ensemble',
        icon: Gauge,
        summary: 'L’état du serveur et de ses conteneurs, d’un coup d’œil.',
      },
      {
        path: '/conteneurs',
        label: 'Conteneurs',
        icon: Boxes,
        summary: 'Les conteneurs Docker du serveur, leur image, leurs ports et leur état.',
      },
      {
        path: '/deploiements',
        label: 'Déploiements',
        icon: Rocket,
        summary:
          'Déployer une application Docker Compose et suivre sa mise en service en direct, journal compris.',
      },
      {
        path: '/surveillance',
        label: 'Surveillance',
        icon: Activity,
        summary:
          'Processeur, mémoire, réseau et disque du serveur et de chaque conteneur, en direct.',
      },
      {
        path: '/assistant',
        label: 'Assistant',
        icon: MessageSquareText,
        upcoming: true,
        summary:
          'Interroger le serveur en langage courant et obtenir le résumé des journaux d’un déploiement en échec.',
      },
    ],
  },
  {
    label: 'Services',
    sections: [
      {
        path: '/marketplace',
        label: 'Marketplace',
        icon: Store,
        summary:
          'Installer en un clic une base de données, un CMS ou un outil DevOps prêt à l’emploi.',
      },
      {
        path: '/domaines',
        label: 'Domaines',
        icon: Globe,
        summary:
          'Relier un nom de domaine à un service, avec un certificat HTTPS délivré automatiquement.',
      },
      {
        path: '/bases-de-donnees',
        label: 'Bases de données',
        icon: Database,
        upcoming: true,
        summary:
          'Créer une base PostgreSQL, MySQL, MongoDB ou Redis, puis la sauvegarder et la restaurer.',
      },
    ],
  },
  {
    label: 'Organisation',
    sections: [
      {
        path: '/equipe',
        label: 'Équipe',
        icon: Users,
        summary: 'Les comptes du serveur : inviter un membre, suivre son invitation, la renvoyer.',
      },
      {
        path: '/reglages',
        label: 'Réglages',
        icon: Settings,
        summary: 'Les réglages du serveur, comme l’envoi d’e-mails, et de quoi les essayer.',
      },
    ],
  },
];

export const sections: Section[] = navigation.flatMap((group) => group.sections);
