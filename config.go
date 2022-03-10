package cli

// Config is the cli App's configuration.
type Config struct {
	// Name of the App. Defaults to the binary file name.
	Name string
	// Version of the App. Defaults to 0.1.0.
	Version string
	// Build time of the App. Defaults to mod time of the binary file.
	BuildTime string
}
