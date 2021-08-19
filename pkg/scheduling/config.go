package scheduling

// Flow configuration.
type Flow struct {
	Tasks []Task `yaml:"tasks" toml:"tasks" json:"tasks"`
}

// Task configuration.
type Task struct {
	ID   string   `yaml:"id" toml:"id" json:"id"`
	Deps []string `yaml:"deps" toml:"deps" json:"deps,omitempty"` // Dependencies
}
