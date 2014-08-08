package reset

import (
	"fmt"
	"time"

	"github.com/fsouza/go-dockerclient"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

var endpoint = "unix:///var/run/docker.sock"

// Restarts the tpires/Neo4j docker container
// docker run -d --privileged -p 7474:7474 tpires/neo4j --name neo4j
func RemoveNeo4jDB() {

	// docker server URL

	// remove options
	ropts := docker.RemoveContainerOptions{
		ID:            "neo4j",
		RemoveVolumes: true,
		Force:         true,
	}

	client, err := docker.NewClient(endpoint)
	panicErr(err)

	fmt.Println("Removing Neo4j container")
	err = client.RemoveContainer(ropts)
	panicErr(err)

}

func StartNeo4jDB() {
	// create options
	copts := docker.CreateContainerOptions{
		Name: "neo4j",
		Config: &docker.Config{
			Image: "tpires/neo4j",
			ExposedPorts: map[docker.Port]struct{}{
				docker.Port("7474"): {},
			},
		},
	}

	// start options for:
	sopts := &docker.HostConfig{
		ContainerIDFile: "tpires/neo4j",
		Privileged:      true,
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port("7474"): {
				docker.PortBinding{
					HostIp:   "0.0.0.0",
					HostPort: "7474",
				},
			},
		},
	}

	client, err := docker.NewClient(endpoint)
	panicErr(err)

	fmt.Println("Creating Neo4j container")
	_, err = client.CreateContainer(copts)
	panicErr(err)

	fmt.Println("Starting Neo4j container")
	err = client.StartContainer("neo4j", sopts)
	panicErr(err)

	fmt.Println("Sleeping 10s")
	time.Sleep(10 * time.Second)

	fmt.Println("Successfully restarted Neo4j\n")
}
