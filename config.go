package patrasche

type Config struct {
	Fabric struct {
		EnvPrefix string   `mapstructure:""`
		Channel   string   `mapstructure:""`
		Config    []string `mapstructure:""`
		Identity  struct {
			Organization string `mapstructure:""`
			Username     string `mapstructure:""`
		}
	}

	Logging struct {
		Level string `mapstructure:""`
	}
}
