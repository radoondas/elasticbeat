package beater

type EsConfig struct {
	Period *int64

	URLs []string

	Stats struct {
		Nodes   *bool
		Cluster *bool
		Health  *bool
	}

	// Authentication for BasicAuth
	Authentication struct {
		Username *string
		Password *string
	}
}

type ConfigSettings struct {
	Input EsConfig
}
