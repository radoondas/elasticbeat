package beat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ClusterHealth struct {
	Cluster_name                     string  `json:"cluster_name"`
	Status                           string  `json:"status"`
	Status_num                       int64   `json:"status_num"`
	Timed_out                        bool    `json:"timed_out"`
	Number_of_nodes                  uint64  `json:"number_of_nodes"`
	Number_of_data_nodes             uint64  `json:"number_of_data_nodes"`
	Active_primary_shards            uint64  `json:"active_primary_shards"`
	Active_shards                    uint64  `json:"active_shards"`
	Relocating_shards                uint64  `json:"relocating_shards"`
	Initializing_shards              uint64  `json:"intializing_shards"`
	Unassigned_shards                uint64  `json:"unassigned_shards"`
	Delayed_unassigned_shards        uint64  `json:"delayed_unassigned_shards"`
	Number_of_pending_tasks          uint64  `json:"number_of_pending_tasks"`
	Number_of_in_flight_fetch        uint64  `json:"number_of_in_flight_fetch"`
	Task_max_waiting_in_queue_millis uint64  `json:"task_max_waiting_in_queue_millis"`
	Active_shards_percent_as_number  float64 `json:"active_shards_percent_as_number"`
}

type ClusterStats struct {
	Timestamp    uint64 `json:"timestamp"`
	Cluster_name string `json:"cluster_name"`
	Status       string `json:"status"`
	Status_num   int64  `json:"status_num"`
	Indices      struct {
		Count  uint64 `json:"count"`
		Shards struct {
			Total       uint64  `json:"total"`
			Primaries   uint64  `json:"primaries"`
			Replication float64 `json:"replication"`
			Index       struct {
				Shards struct {
					Min float64 `json:"nin"`
					Max float64 `json:"max"`
					Avg float64 `json:"avg"`
				}
				Primaries struct {
					Min float64 `json:"nin"`
					Max float64 `json:"max"`
					Avg float64 `json:"avg"`
				}
				Replication struct {
					Min float64 `json:"nin"`
					Max float64 `json:"max"`
					Avg float64 `json:"avg"`
				}
			}
		}
		Docs struct {
			Count   uint64 `json:"count"`
			Deleted uint64 `json:"deleted"`
		}
		Store struct {
			Size_in_bytes           uint64 `json:"size_in_bytes"`
			Throttle_time_in_millis uint64 `json:"throttle_time_in_millis"`
		}
		Fielddata struct {
			Memory_size_in_bytes uint64 `json:"memory_size_in_bytes"`
			Evictions            uint64 `json:"evictions"`
		}
		Query_cache struct {
			Memory_size_in_bytes uint64 `json:"memory_size_in_bytes"`
			Total_count          uint64 `json:"total_count"`
			Hit_count            uint64 `json:"hit_count"`
			Miss_count           uint64 `json:"miss_count"`
			Cache_size           uint64 `json:"cache_size"`
			Cache_count          uint64 `json:"cache_count"`
			Evictions            uint64 `json:"evictions"`
		}
		Completion struct {
			Size_in_bytes uint64 `json:"size_in_bytes"`
		}
		Segments struct {
			Count                            uint64 `json:"count"`
			Memory_in_bytes                  uint64 `json:"memory_in_bytes"`
			Terms_memory_in_bytes            uint64 `json:"terms_memory_in_bytes"`
			Stored_fields_memory_in_bytes    uint64 `json:"stored_fields_memory_in_bytes"`
			Term_vectors_memory_in_bytes     uint64 `json:"term_vectors_memory_in_bytes"`
			Norms_memory_in_bytes            uint64 `json:"norms_memory_in_bytes"`
			Doc_values_memory_in_bytes       uint64 `json:"doc_values_memory_in_bytes"`
			Index_writer_memory_in_bytes     uint64 `json:"index_writer_memory_in_bytes"`
			Index_writer_max_memory_in_bytes uint64 `json:"index_writer_max_memory_in_bytes"`
			Version_map_memory_in_bytes      uint64 `json:"version_map_memory_in_bytes"`
			Fixed_bit_set_memory_in_bytes    uint64 `json:"fixed_bit_set_memory_in_bytes"`
		}
		Percolate struct {
			Total                uint64 `json:"total"`
			Time_in_millis       uint64 `json:"time_in_millis"`
			Current              uint64 `json:"current"`
			Memory_size_in_bytes int64  `json:"memory_size_in_bytes"`
			Memory_size          string `json:"memory_size"` // TODO: check value
			Queries              uint64 `json:"queries"`
		}
	}
	Nodes struct {
		Count struct {
			Total       uint64 `json:"total"`
			Master_only uint64 `json:"master_only"`
			Data_only   uint64 `json:"data_only"`
			Master_data uint64 `json:"master_data"`
			Client      uint64 `json:"client"`
		}
		//Versions - not implemented
		Os struct {
			Available_processors uint64 `json:"available_processors"`
			Allocated_processors uint64 `json:"allocated_processors"`
			Mem                  struct {
				Total_in_bytes uint64 `json:"total_in_bytes"`
			}
			//Names - not implemented
		}
		Process struct {
			Cpu struct {
				Percent uint64 `json:"percent"`
			}
			Open_file_descriptors struct {
				Min float64 `json:"min"`
				Max float64 `json:"max"`
				Avg float64 `json:"avg"`
			}
		}
		Jvm struct {
			Max_uptime_in_millis uint64 `json:"max_uptime_in_millis"`
			// Versions - not implemented
			Mem struct {
				Heap_used_in_bytes uint64 `json:"heap_used_in_bytes"`
				Heap_max_in_bytes  uint64 `json:"heap_max_in_bytes"`
			}
			Threads uint64 `json:"threads"`
		}
		Fs struct {
			Total_in_bytes     uint64 `json:"total_in_bytes"`
			Free_in_bytes      uint64 `json:"free_in_bytes"`
			Available_in_bytes uint64 `json:"available_in_bytes"`
			Spins              string `json:"spins"`
		}
		// Plugins - not implemented
	}
}

func (eb *Elasticbeat) GetCLusterHealth(u url.URL) (ClusterHealth, error) {
	health := ClusterHealth{}

	res, err := http.Get(TrimSuffix(u.String(), "/") + "/_cluster/health?pretty")
	if err != nil {
		return health, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return health, fmt.Errorf("HTTP%s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return health, err
	}

	err = json.Unmarshal([]byte(body), &health)
	if err != nil {
		return health, err
	}
	health.Status_num = GetNumericalClusterStatus(health.Status)

	return health, nil
}

func (eb *Elasticbeat) GetCLusterStats(u url.URL) (ClusterStats, error) {
	stats := ClusterStats{}

	res, err := http.Get(TrimSuffix(u.String(), "/") + "/_cluster/stats?pretty")
	if err != nil {
		return stats, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return stats, fmt.Errorf("HTTP%s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return stats, err
	}

	err = json.Unmarshal([]byte(body), &stats)
	if err != nil {
		return stats, err
	}

	stats.Status_num = GetNumericalClusterStatus(stats.Status)

	return stats, nil
}
