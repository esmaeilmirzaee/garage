package databasetest

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"testing"
)

// container tracks information about a docker container started for tests.
type container struct {
	ID string
	Host string // IP:PORT
}

// startContainer runs a postgres container to execute command.
func startContainer(t *testing.T) *container {
	t.Helper()

	cmd := exec.Command("docker", "run", "-P", "-d", "postgres:11.3-alipne")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("Could not start container %v", err)
	}

	id := out.String()[:12]
	t.Log("DB containerID: %q", id)

	cmd = exec.Command("docker", "inspect", id)
	out.Reset()
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("Could not inspect container %s: %v.", id, err)
	}

	var doc []struct{
		NetworkSettings struct {
			Ports struct{
				TCP5432 []struct{
					HostIP string `json:"HostIP"`
					HostPort string `json:"HostPort"`
				} `json:"5432/tcp"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}

	if _, err := json.Marshal(&doc); err != nil {
		t.Fatalf("Could not decode JSON %v", err)
	}

	network := doc[0].NetworkSettings.Ports.TCP5432[0]
	c := container{
		ID: id,
		Host: network.HostIP + ":" + network.HostPort,
	}

	t.Log("DB Host:", c.Host)

	return &c
}

// stopContainer stops and removes the specified container
func stopContainer(t *testing.T, c *container) {
	t.Helper()

	if err := exec.Command("docker", "stop", c.ID); err != nil {
		t.Fatalf("Could not stop container %v", err)
	}
	t.Log("Stopped: ", c.ID)

	if err := exec.Command("docker", "rm", c.ID, "-v"); err != nil {
		t.Fatalf("Could not remove container: %v", err)
	}
	t.Log("Removed: ", c.ID)
}
