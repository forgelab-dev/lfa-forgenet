// Package metrics collecte l'utilisation des ressources de l'hôte (gopsutil) et
// des conteneurs (statistiques du démon Docker).
//
// Les compteurs réseau et disque sont cumulés depuis le démarrage, comme des
// compteurs Prometheus : le débit se calcule chez le client par différence
// entre deux mesures. Seule l'utilisation du processeur est un instantané, car
// elle demande deux relevés : le collecteur garde le précédent.
package metrics

import "time"

// CPU décrit l'utilisation du processeur.
type CPU struct {
	// UsagePercent est l'utilisation depuis la mesure précédente, de 0 à 100
	// pour l'hôte, de 0 à 100 × OnlineCPUs pour un conteneur (convention de
	// `docker stats`).
	UsagePercent float64
	OnlineCPUs   uint32
	// SystemNanoseconds est le temps processeur cumulé passé en mode noyau.
	SystemNanoseconds uint64
}

// Memory décrit l'utilisation de la mémoire vive.
type Memory struct {
	UsageBytes   uint64
	LimitBytes   uint64
	CacheBytes   uint64 // cache récupérable, exclu de UsageBytes
	UsagePercent float64
}

// Network cumule le trafic des interfaces retenues.
type Network struct {
	RxBytes, TxBytes     uint64
	RxPackets, TxPackets uint64
	RxErrors, TxErrors   uint64
}

// DiskIO cumule les entrées-sorties des disques retenus.
type DiskIO struct {
	ReadBytes, WriteBytes uint64
	ReadOps, WriteOps     uint64
}

// Storage décrit l'occupation d'un système de fichiers.
type Storage struct {
	TotalBytes, UsedBytes, FreeBytes uint64
	UsagePercent                     float64
}

// Host est une mesure de la machine hôte.
type Host struct {
	CollectedAt      time.Time
	UptimeSeconds    uint64
	CPU              CPU
	Memory           Memory
	Storage          Storage
	Network          Network
	DiskIO           DiskIO
	ActiveContainers uint32
}

// Container est une mesure d'un conteneur.
type Container struct {
	ContainerID string
	CollectedAt time.Time
	CPU         CPU
	Memory      Memory
	Network     Network
	DiskIO      DiskIO
}

func percent(part, whole uint64) float64 {
	if whole == 0 {
		return 0
	}
	return float64(part) / float64(whole) * 100
}
