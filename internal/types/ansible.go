package types

type (
	StringList map[string]string
	Group      struct {
		Hosts         []string              `json:"hosts,omitempty"`
		Vars          StringList            `json:"vars,omitempty"`
		HostVariables map[string]StringList `json:"hostvars,omitempty"`
	}
	Answer map[string]Group
)
