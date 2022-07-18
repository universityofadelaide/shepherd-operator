package cli

// CommandParams configures an AWS CLI command.
type CommandParams struct {
	// Endpoint which is used to override the AWS services endpoint eg. For local development.
	Endpoint string
	// Service which the command will interact with eg. S3.
	Service string
	// Operation which will be performed with the Service.
	Operation string
	// Args used as part of an Operation.
	Args []string
}

// Command which is compatible with the AWS CLI.
func Command(params CommandParams) []string {
	command := []string{params.Service}

	if params.Endpoint != "" {
		command = append(command, "--endpoint-url", params.Endpoint)
	}

	command = append(command, params.Operation)

	return append(command, params.Args...)
}
