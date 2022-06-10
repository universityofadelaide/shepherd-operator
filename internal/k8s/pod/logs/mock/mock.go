package mock

// Client for mocking pod Logs.
type Client struct {
	Logs  string
	Error error
}

// Get mocks getting the Logs for a container from a pod in a namespace.
func (c Client) Get(namespace, name, container string) (string, error) {
	if c.Error != nil {
		return "", c.Error
	}
	return c.Logs, nil
}
