package policy

// Document for defining a policy.
type Document struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

// Statement which is part of a Document.
type Statement struct {
	Effect    string     `json:"Effect"`
	Action    []string   `json:"Action"`
	Resource  string     `json:"Resource"`
	Condition Conditions `json:"Condition"`
}

// Conditions which are part of a Document Statement.
type Conditions struct {
	StringEquals map[string]string `json:"StringEquals"`
}
