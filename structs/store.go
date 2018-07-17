package structs

// Store struct
type Store struct {
	Name string `mapstructure:"name"`
	Type string `mapstructure:"type"`
	// Consul specific
	Datacenter string `mapstructure:"datacenter"`
}
