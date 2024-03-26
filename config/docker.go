package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func FindContainerIPs(suffix string) ([]string, error) {
	ctx := context.Background()

	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	// List all containers
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	// Filter containers with the specified suffix
	var matchingContainers []string
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.HasSuffix(name, suffix) {
				matchingContainers = append(matchingContainers, container.ID)
				break
			}
		}
	}

	if len(matchingContainers) == 0 {
		return nil, fmt.Errorf("no containers found with suffix %q", suffix)
	}
	if len(matchingContainers) > 1 {
		return nil, fmt.Errorf("multiple containers found with suffix %q", suffix)
	}

	// Get the IP addresses of matching containers
	var ips []string
	for _, containerID := range matchingContainers {
		containerInfo, err := cli.ContainerInspect(ctx, containerID)
		if err != nil {
			return nil, err
		}

		// Extract the container's IP address
		for _, n := range containerInfo.NetworkSettings.Networks {
			ips = append(ips, n.IPAddress)
		}
	}

	return ips, nil
}
