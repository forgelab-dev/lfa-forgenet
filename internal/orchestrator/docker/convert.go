package docker

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"

	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
)

func fromSummary(c container.Summary) orchestrator.Container {
	name := ""
	if len(c.Names) > 0 {
		name = strings.TrimPrefix(c.Names[0], "/")
	}

	ports := make([]orchestrator.Port, 0, len(c.Ports))
	for _, p := range c.Ports {
		port := orchestrator.Port{HostPort: p.PublicPort, ContainerPort: p.PrivatePort, Protocol: p.Type}
		if p.IP.IsValid() {
			port.HostIP = p.IP.String()
		}
		ports = append(ports, port)
	}

	return orchestrator.Container{
		ID:        c.ID,
		Name:      name,
		Image:     c.Image,
		State:     fromState(c.State),
		Ports:     ports,
		CreatedAt: time.Unix(c.Created, 0).UTC(),
	}
}

// fromInspect tolère les champs absents (Config, State, NetworkSettings
// peuvent être nil selon la version du démon) au lieu de paniquer.
func fromInspect(res client.ContainerInspectResult) orchestrator.Container {
	c := res.Container
	out := orchestrator.Container{
		ID:    c.ID,
		Name:  strings.TrimPrefix(c.Name, "/"),
		Image: c.Image,
	}
	if c.Config != nil && c.Config.Image != "" {
		out.Image = c.Config.Image
	}
	if c.State != nil {
		out.State = fromState(c.State.Status)
	}
	if created, err := time.Parse(time.RFC3339Nano, c.Created); err == nil {
		out.CreatedAt = created.UTC()
	}
	if c.NetworkSettings != nil {
		for port, bindings := range c.NetworkSettings.Ports {
			if len(bindings) == 0 {
				out.Ports = append(out.Ports, orchestrator.Port{ContainerPort: port.Num(), Protocol: string(port.Proto())})
				continue
			}
			for _, b := range bindings {
				hostPort, _ := strconv.ParseUint(b.HostPort, 10, 16)
				p := orchestrator.Port{HostPort: uint16(hostPort), ContainerPort: port.Num(), Protocol: string(port.Proto())}
				if b.HostIP.IsValid() {
					p.HostIP = b.HostIP.String()
				}
				out.Ports = append(out.Ports, p)
			}
		}
		// Les ports viennent d'une map : on fixe l'ordre pour une réponse stable.
		sort.Slice(out.Ports, func(i, j int) bool {
			a, b := out.Ports[i], out.Ports[j]
			if a.ContainerPort != b.ContainerPort {
				return a.ContainerPort < b.ContainerPort
			}
			return a.HostIP < b.HostIP
		})
	}
	return out
}

func fromState(s container.ContainerState) orchestrator.State {
	switch state := orchestrator.State(s); state {
	case orchestrator.StateCreated, orchestrator.StateRunning, orchestrator.StatePaused,
		orchestrator.StateRestarting, orchestrator.StateRemoving, orchestrator.StateExited, orchestrator.StateDead:
		return state
	default:
		return orchestrator.StateUnknown
	}
}
