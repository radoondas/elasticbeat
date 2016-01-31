package beater

type EsConfig struct {
	Period *int64

	URLs []string

	Stats struct {
		Nodes   *bool
		Cluster *bool
		Health  *bool
	}
}

type ConfigSettings struct {
	Input EsConfig
}
