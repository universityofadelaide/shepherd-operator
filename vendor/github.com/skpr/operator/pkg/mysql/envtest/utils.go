package envtest

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/go-connections/nat"
)

// Helper function to wait for the MySQL service to become healthy.
func waitForMySQL(host, username, password string, port, retries int, interval time.Duration) bool {
	var counter int

	for {
		counter++

		conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", username, password, host, port)

		mysql, err := sql.Open("mysql", conn)
		if err != nil {
			return false
		}

		err = mysql.Ping()
		if err == nil {
			return true
		}

		if counter >= retries {
			return false
		}

		time.Sleep(interval)
	}
}

// Helper function to lookup the port mapped for a container.
func getContainerPort(ports nat.PortMap, port nat.Port) (int, error) {
	for key, bindings := range ports {
		if key == port {
			for _, binding := range bindings {
				return strconv.Atoi(binding.HostPort)
			}
		}
	}

	return 0, fmt.Errorf("failed to find port: %s", port)
}
