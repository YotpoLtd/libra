package nomad

// Job Struct
type Job struct {
	Name   string
	Groups map[string]*Group `hcl:"group"`
}
