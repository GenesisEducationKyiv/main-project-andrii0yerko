package testenv

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func SetupTestMailserver(t *testing.T) (string, string, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mailhog/mailhog", Tag: "v1.0.1", Env: []string{"listen_addresses = '*'"}},
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	apiURI := fmt.Sprintf("http://localhost:%v/api/v2/messages", resource.GetPort("8025/tcp"))

	err = pool.Retry(func() error {
		r, err := http.Get(apiURI) //nolint:govet,gosec,noctx // closure func
		r.Body.Close()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", "", err
	}

	t.Cleanup(func() {
		err := pool.Purge(resource) //nolint:govet // closure func
		if err != nil {
			log.Printf("Could not purge resource: %s", err)
		}
	})

	return resource.GetPort("1025/tcp"), resource.GetPort("8025/tcp"), nil
}
