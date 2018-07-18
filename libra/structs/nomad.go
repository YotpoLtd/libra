package structs

// Job Struct
type NomadJob struct {
	Name   string
	Groups map[string]*NomadJobGroup `hcl:"group"`
}

type NomadJobGroup struct {
	Name     string
	MinCount int              `hcl:"min_count"`
	MaxCount int              `hcl:"max_count"`
	Rules    map[string]*Rule `hcl:"rule"`
}

type NomadConfig struct {
	Address  string `hcl:"address"`
	ACLToken string `hcl:"acl_token"`
}
