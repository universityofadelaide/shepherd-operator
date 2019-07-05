package envtest

import (
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"golang.org/x/net/context"
)

const (
	// DefaultImage which will be used if non is provided.
	DefaultImage = "docker.io/library/mariadb:10"
	// InsecurePort assigned to the MySQL instance.
	InsecurePort = "3306"
	// BindAddress assigned to the MySQL instance.
	BindAddress = "0.0.0.0"
)

// MySQL which will back a testsuite.
type MySQL struct {
	// Image which will be used for spinning up K3s.
	Image string
	// Internal container identifier.
	id string
}

// Start the test environment.
func (e *MySQL) Start() (string, string, string, int, error) {
	if e.Image == "" {
		e.Image = DefaultImage
	}

	ctx := context.Background()

	cli, err := client.NewEnvClient()
	if err != nil {
		return "", "", "", 0, err
	}

	natPort, err := nat.NewPort("tcp", InsecurePort)
	if err != nil {
		return "", "", "", 0, err
	}

	var (
		hostname = "127.0.0.1"
		username = "root"
		password = "root"
	)

	containerConfig := &container.Config{
		Image: e.Image,
		Env: []string{
			fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", password),
		},
		ExposedPorts: nat.PortSet{
			natPort: {},
		},
	}

	containerHostConfig := &container.HostConfig{
		Privileged: true,
		PortBindings: map[nat.Port][]nat.PortBinding{
			natPort: []nat.PortBinding{
				{
					HostIP: BindAddress,
				},
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, containerConfig, containerHostConfig, nil, "")
	if err != nil {
		return "", "", "", 0, err
	}

	e.id = resp.ID

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", "", "", 0, err
	}

	inspect, err := cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return "", "", "", 0, err
	}

	port, err := getContainerPort(inspect.NetworkSettings.Ports, natPort)
	if err != nil {
		return "", "", "", 0, err
	}

	ready := waitForMySQL(hostname, username, password, port, 20, time.Second*5)
	if !ready {
		return "", "", "", 0, fmt.Errorf("cluster did not become available")
	}

	return hostname, username, password, port, nil
}

// Stop the test environment.
func (e *MySQL) Stop() error {
	ctx := context.Background()

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	err = cli.ContainerStop(ctx, e.id, nil)
	if err != nil {
		return err
	}

	return cli.ContainerRemove(ctx, e.id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
		Force:         true,
	})

}
