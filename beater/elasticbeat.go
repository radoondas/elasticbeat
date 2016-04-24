package beater

import (
	"errors"
	"net/url"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
)

const selector = "elasticbeat"
const selectorDetail = "json"

type Elasticbeat struct {
	period time.Duration
	urls   []*url.URL

	EbConfig ConfigSettings
	events   publisher.Client

	nodeStats          bool
	clusterStats       bool
	clusterHealthStats bool

	done chan struct{}
}

func New() *Elasticbeat {
	return &Elasticbeat{
		done: make(chan struct{}),
	}
}

func (eb *Elasticbeat) Config(b *beat.Beat) error {

	err := cfgfile.Read(&eb.EbConfig, "")
	if err != nil {
		logp.Err("Error reading configuration file: %v", err)
		return err
	}

	if eb.EbConfig.Input.Period != nil {
		eb.period = time.Duration(*eb.EbConfig.Input.Period) * time.Second
	} else {
		eb.period = 10 * time.Second
	}

	//define default URL if none provided
	var urlConfig []string
	if eb.EbConfig.Input.URLs != nil {
		urlConfig = eb.EbConfig.Input.URLs
	} else {
		urlConfig = []string{"http://127.0.0.1:9200"}
	}

	eb.urls = make([]*url.URL, len(urlConfig))
	for i := 0; i < len(urlConfig); i++ {
		u, err := url.Parse(urlConfig[i])
		if err != nil {
			logp.Err("Invalid ElasticSearch URL: %v", err)
			return err
		}
		eb.urls[i] = u
	}

	if eb.EbConfig.Input.Stats.Nodes != nil {
		eb.nodeStats = *eb.EbConfig.Input.Stats.Nodes
	} else {
		eb.nodeStats = true
	}

	if eb.EbConfig.Input.Stats.Cluster != nil {
		eb.clusterStats = *eb.EbConfig.Input.Stats.Cluster
	} else {
		eb.clusterStats = true
	}

	if eb.EbConfig.Input.Stats.Health != nil {
		eb.clusterHealthStats = *eb.EbConfig.Input.Stats.Health
	} else {
		eb.clusterHealthStats = true
	}

	if !eb.nodeStats && !eb.nodeStats && !eb.clusterHealthStats {
		return errors.New("Invalid statistics configuration")
	}

	logp.Debug(selector, "Init elasticbeat")
	logp.Debug(selector, "Period %v\n", eb.period)
	logp.Debug(selector, "Watch %v", eb.urls)
	logp.Debug(selector, "Nodes statistics %t\n", eb.nodeStats)
	logp.Debug(selector, "Cluster statistics %t\n", eb.clusterStats)
	logp.Debug(selector, "Cluster health statistics %t\n", eb.clusterHealthStats)

	return nil
}

func (eb *Elasticbeat) Setup(b *beat.Beat) error {
	eb.events = b.Publisher.Connect()
	eb.done = make(chan struct{})
	return nil
}

func (eb *Elasticbeat) Run(b *beat.Beat) error {
	logp.Debug(selector, "Run elasticbeat")

	//for each url
	for _, u := range eb.urls {
		go func(u *url.URL) {
			ticker := time.NewTicker(eb.period)
			defer ticker.Stop()

			for {
				select {
				case <-eb.done:
					goto GotoFinish
				case <-ticker.C:
				}

				timerStart := time.Now()

				if eb.clusterStats {
					logp.Debug(selector, "Cluster stats for url: %v", u)
					cluster_stats, err := eb.GetCLusterStats(*u)
					if err != nil {
						logp.Err("Error reading cluster stats: %v", err)
					} else {
						logp.Debug(selectorDetail, "Cluster stats detail: %+v", cluster_stats)

						event := common.MapStr{
							"@timestamp":    common.Time(time.Now()),
							"type":          "cluster_stats",
							"url":           u.String(),
							"cluster_stats": cluster_stats,
						}

						eb.events.PublishEvent(event)
					}
				}

				if eb.nodeStats {
					logp.Debug(selector, "Nodes stats for url: %v", u)
					nodes, err := eb.GetNodesStats(*u)
					if err != nil {
						logp.Err("Error reading nodes stats: %v\n", err)
					} else {
						for _, nd := range nodes {
							logp.Debug(selectorDetail, "Node stats detail: %+v", nd)

							event := common.MapStr{
								"@timestamp":   common.Time(time.Now()),
								"type":         "cluster_node",
								"url":          u.String(),
								"cluster_node": nd,
							}
							eb.events.PublishEvent(event)
						}
					}

				}

				if eb.clusterHealthStats {
					logp.Debug(selector, "Cluster health for url: %v", u)
					cluster_health, err := eb.GetCLusterHealth(*u)
					if err != nil {
						logp.Err("Error reading cluster health: %v", err)
					} else {
						logp.Debug(selectorDetail, "Cluster health detail: %+v", cluster_health)

						event := common.MapStr{
							"@timestamp":     common.Time(time.Now()),
							"type":           "cluster_health",
							"url":            u.String(),
							"cluster_health": cluster_health,
						}

						eb.events.PublishEvent(event)
					}
				}
				timerEnd := time.Now()
				duration := timerEnd.Sub(timerStart)
				if duration.Nanoseconds() > eb.period.Nanoseconds() {
					logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
				}
			}

		GotoFinish:
		}(u)
	}

	<-eb.done
	return nil
}

func (eb *Elasticbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (eb *Elasticbeat) Stop() {
	logp.Debug(selector, "Stop elasticbeat")
	close(eb.done)
}
