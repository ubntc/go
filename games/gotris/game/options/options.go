package options

// Options defines map-slice interface for managing a single list of selectable options
// in combination with a key-value map for storing additional config values.
//
// The interface allows a specific platform to implement concrete Options structs,
// whose setters can be used to control platform behavior.
// The game package uses this to implement generic game options.
// The platforms use this implement, e.g., different rendering modes.
type Options interface {
	Select
	Map
}

type Select interface {
	Set(idx int)     // Set sets the selected option according to the given options index.
	Get() int        // Get returns the current index of the selected option.
	Len() int        // Len returns the number of options.
	GetName() string // GetName returns the name/title of the current option.
	List() []string  // List returns all names/titles for all possible options.
	Descs() []string // Descs returns descriptions for all possible options.
}

type Map interface {
	Values() map[string]string  // Values returns all map values.
	SetValue(key, value string) // SetValue sets the map value for the given key.
	GetValue(key string) string // GetValue returns the map value for the given key.
}
