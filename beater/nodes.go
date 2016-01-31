package beater

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/elastic/beats/libbeat/logp"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type NodeStats struct {
	Timestamp uint64 `json:"timestamp"`
	Name      string `json:"name"`
	// transport_address - not implemented
	Host string `json:"host"`
	//[]ip - not implemented
	Cluster_name string `json:"cluster_name"`
	Os           struct {
		Timestamp    uint64  `json:"timestamp"`
		Load_average float64 `json:"load_average"`
		Mem          struct {
			Total_in_bytes uint64 `json:"total_in_bytes"`
			Free_in_bytes  uint64 `json:"free_in_bytes"`
			Used_in_bytes  uint64 `json:"used_in_bytes"`
			Free_percent   uint64 `json:"free_percent"`
			Used_percent   uint64 `json:"used_percent"`
		}
		Swap struct {
			Total_in_bytes uint64 `json:"total_in_bytes"`
			Free_in_bytes  uint64 `json:"free_in_bytes"`
			Used_in_bytes  uint64 `json:"used_in_bytes"`
		}
	}
	Process struct {
		Timestamp             uint64 `json:"timestamp"`
		Open_file_descriptors uint64 `json:"open_file_descriptors"`
		Max_file_descriptors  uint64 `json:"max_file_descriptors"`
		Cpu                   struct {
			Percent         uint64 `json:"percent"`
			Total_in_millis uint64 `json:"total_in_millis"`
		}
		Mem struct {
			Total_virtual_in_bytes uint64 `json:"total_virtual_in_bytes"`
		}
	}
	Jvm struct {
		Timestamp        uint64 `json:"timestamp"`
		Uptime_in_millis uint64 `json:"uptime_in_millis"`
		Mem              struct {
			Heap_used_in_bytes          uint64 `json:"heap_used_in_bytes"`
			Heap_used_percent           uint64 `json:"heap_used_percent"`
			Heap_committed_in_bytes     uint64 `json:"heap_committed_in_bytes"`
			Heap_max_in_bytes           uint64 `json:"heap_max_in_bytes"`
			Non_heap_used_in_bytes      uint64 `json:"non_heap_used_in_bytes"`
			Non_heap_committed_in_bytes uint64 `json:"non_heap_committed_in_bytes"`
			Pools                       struct {
				Young struct {
					Used_in_bytes      uint64 `json:"used_in_bytes"`
					Max_in_bytes       uint64 `json:"max_in_bytes"`
					Peak_used_in_bytes uint64 `json:"peak_used_in_bytes"`
					Peak_max_in_bytes  uint64 `json:"peak_max_in_bytes"`
				}
				Survivor struct {
					Used_in_bytes      uint64 `json:"used_in_bytes"`
					Max_in_bytes       uint64 `json:"max_in_bytes"`
					Peak_used_in_bytes uint64 `json:"peak_used_in_bytes"`
					Peak_max_in_bytes  uint64 `json:"peak_max_in_bytes"`
				}
				Old struct {
					Used_in_bytes      uint64 `json:"used_in_bytes"`
					Max_in_bytes       uint64 `json:"max_in_bytes"`
					Peak_used_in_bytes uint64 `json:"peak_used_in_bytes"`
					Peak_max_in_bytes  uint64 `json:"peak_max_in_bytes"`
				}
			}
		}
		Threads struct {
			Count      uint64 `json:"count"`
			Peak_count uint64 `json:"peak_count"`
		}
		Gc struct {
			Collectors struct {
				Young struct {
					Collection_count          uint64 `json:"collection_count"`
					Collection_time_in_millis uint64 `json:"collection_time_in_millis"`
				}
				Old struct {
					Collection_count          uint64 `json:"collection_count"`
					Collection_time_in_millis uint64 `json:"collection_time_in_millis"`
				}
			}
		}
		Buffer_pools struct {
			Direct struct {
				Count                   uint64 `json:"count"`
				Used_in_bytes           uint64 `json:"used_in_bytes"`
				Total_capacity_in_bytes uint64 `json:"total_capacity_in_bytes"`
			}
			Mapped struct {
				Count                   uint64 `json:"count"`
				Used_in_bytes           uint64 `json:"used_in_bytes"`
				Total_capacity_in_bytes uint64 `json:"total_capacity_in_bytes"`
			}
		}
	}
	Thread_pool struct {
		Bulk struct {
			Threads   uint64 `json:"threads"`
			Queue     uint64 `json:"queue"`
			Active    uint64 `json:"active"`
			Rejected  uint64 `json:"rejected"`
			Largest   uint64 `json:"largest"`
			Completed uint64 `json:"completed"`
		}
		Fetch_shard_started struct {
			Threads   uint64 `json:"threads"`
			Queue     uint64 `json:"queue"`
			Active    uint64 `json:"active"`
			Rejected  uint64 `json:"rejected"`
			Largest   uint64 `json:"largest"`
			Completed uint64 `json:"completed"`
		}
		Fetch_shard_store struct {
			Threads   uint64 `json:"threads"`
			Queue     uint64 `json:"queue"`
			Active    uint64 `json:"active"`
			Rejected  uint64 `json:"rejected"`
			Largest   uint64 `json:"largest"`
			Completed uint64 `json:"completed"`
		}
		Flush struct {
			Threads   uint64 `json:"threads"`
			Queue     uint64 `json:"queue"`
			Active    uint64 `json:"active"`
			Rejected  uint64 `json:"rejected"`
			Largest   uint64 `json:"largest"`
			Completed uint64 `json:"completed"`
		}
		Force_merge struct {
			Threads   uint64 `json:"threads"`
			Queue     uint64 `json:"queue"`
			Active    uint64 `json:"active"`
			Rejected  uint64 `json:"rejected"`
			Largest   uint64 `json:"largest"`
			Completed uint64 `json:"completed"`
		}
		Generic struct {
			Threads   uint64 `json:"threads"`
			Queue     uint64 `json:"queue"`
			Active    uint64 `json:"active"`
			Rejected  uint64 `json:"rejected"`
			Largest   uint64 `json:"largest"`
			Completed uint64 `json:"completed"`
		}
		Get struct {
			Threads   uint64 `json:"threads"`
			Queue     uint64 `json:"queue"`
			Active    uint64 `json:"active"`
			Rejected  uint64 `json:"rejected"`
			Largest   uint64 `json:"largest"`
			Completed uint64 `json:"completed"`
		}
		Index struct {
			Threads   uint64 `json:"threads"`
			Queue     uint64 `json:"queue"`
			Active    uint64 `json:"active"`
			Rejected  uint64 `json:"rejected"`
			Largest   uint64 `json:"largest"`
			Completed uint64 `json:"completed"`
		}
		Listener struct {
			Threads   uint64 `json:"threads"`
			Queue     uint64 `json:"queue"`
			Active    uint64 `json:"active"`
			Rejected  uint64 `json:"rejected"`
			Largest   uint64 `json:"largest"`
			Completed uint64 `json:"completed"`
		}
		Management struct {
			Threads   uint64 `json:"threads"`
			Queue     uint64 `json:"queue"`
			Active    uint64 `json:"active"`
			Rejected  uint64 `json:"rejected"`
			Largest   uint64 `json:"largest"`
			Completed uint64 `json:"completed"`
		}
		Percolate struct {
			Threads   uint64 `json:"threads"`
			Queue     uint64 `json:"queue"`
			Active    uint64 `json:"active"`
			Rejected  uint64 `json:"rejected"`
			Largest   uint64 `json:"largest"`
			Completed uint64 `json:"completed"`
		}
		Refresh struct {
			Threads   uint64 `json:"threads"`
			Queue     uint64 `json:"queue"`
			Active    uint64 `json:"active"`
			Rejected  uint64 `json:"rejected"`
			Largest   uint64 `json:"largest"`
			Completed uint64 `json:"completed"`
		}
		Search struct {
			Threads   uint64 `json:"threads"`
			Queue     uint64 `json:"queue"`
			Active    uint64 `json:"active"`
			Rejected  uint64 `json:"rejected"`
			Largest   uint64 `json:"largest"`
			Completed uint64 `json:"completed"`
		}
		Snapshot struct {
			Threads   uint64 `json:"threads"`
			Queue     uint64 `json:"queue"`
			Active    uint64 `json:"active"`
			Rejected  uint64 `json:"rejected"`
			Largest   uint64 `json:"largest"`
			Completed uint64 `json:"completed"`
		}
		Suggest struct {
			Threads   uint64 `json:"threads"`
			Queue     uint64 `json:"queue"`
			Active    uint64 `json:"active"`
			Rejected  uint64 `json:"rejected"`
			Largest   uint64 `json:"largest"`
			Completed uint64 `json:"completed"`
		}
		Warmer struct {
			Threads   uint64 `json:"threads"`
			Queue     uint64 `json:"queue"`
			Active    uint64 `json:"active"`
			Rejected  uint64 `json:"rejected"`
			Largest   uint64 `json:"largest"`
			Completed uint64 `json:"completed"`
		}
	}
	Fs struct {
		Timestamp uint64 `json:"timestamp"`
		Total     struct {
			Total_in_bytes     uint64 `json:"total_in_bytes"`
			Free_in_bytes      uint64 `json:"free_in_bytes"`
			Available_in_bytes uint64 `json:"available_in_bytes"`
			Spins              string `json:"spins"`
		}
		// []data - not implemented
	}
	Transport struct {
		Server_open      uint64 `json:"server_open"`
		Rx_count         uint64 `json:"rx_count"`
		Rx_size_in_bytes uint64 `json:"rx_size_in_bytes"`
		Tx_count         uint64 `json:"tx_count"`
		Tx_size_in_bytes uint64 `json:"tx_size_in_bytes"`
	}
	Http struct {
		Current_open uint64 `json:"current_open"`
		Total_opened uint64 `json:"total_opened"`
	}
	Breakers struct {
		Request struct {
			Limit_size_in_bytes     uint64  `json:"limit_size_in_bytes"`
			Limit_size              string  `json:"limit_size"`
			Estimated_size_in_bytes uint64  `json:"estimated_size_in_bytes"`
			Estimated_size          string  `json:"estimated_size"`
			Overhead                float64 `json:"overhead"`
			Tripped                 uint64  `json:"tripped"`
		}
		Fielddata struct {
			Limit_size_in_bytes     uint64  `json:"limit_size_in_bytes"`
			Limit_size              string  `json:"limit_size"`
			Estimated_size_in_bytes uint64  `json:"estimated_size_in_bytes"`
			Estimated_size          string  `json:"estimated_size"`
			Overhead                float64 `json:"overhead"`
			Tripped                 uint64  `json:"tripped"`
		}
		Parent struct {
			Limit_size_in_bytes     uint64  `json:"limit_size_in_bytes"`
			Limit_size              string  `json:"limit_size"`
			Estimated_size_in_bytes uint64  `json:"estimated_size_in_bytes"`
			Estimated_size          string  `json:"estimated_size"`
			Overhead                float64 `json:"overhead"`
			Tripped                 uint64  `json:"tripped"`
		}
	}
	Script struct {
		Compilations    uint64 `json:"compilations"`
		Cache_evictions uint64 `json:"cache_evictions"`
	}
}

type NodesBody struct {
	Cluster_name string          `json:"cluster_name"`
	Nodes        json.RawMessage `json:"nodes"`
}

func (eb *Elasticbeat) GetNodesStats(u url.URL) ([]NodeStats, error) {
	var nodes []NodeStats

	ids, err := eb.GetNodeIDs(u)
	if err != nil {
		return nil, err
	}
	logp.Debug(selector, "Node ids: %+v", ids)

	if len(ids) > 0 {

		res, err := http.Get(TrimSuffix(u.String(), "/") + "/_nodes/stats/process,jvm,os,fs,thread_pool,transport,http,breaker,script")
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return nil, fmt.Errorf("HTTP%s", res.Status)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		ns := &NodesBody{}
		err = json.Unmarshal([]byte(body), &ns)
		if err != nil {
			return nil, err
		}

		info := make(map[string]NodeStats)
		nodes = make([]NodeStats, len(ids))
		dec := json.NewDecoder(strings.NewReader(string(ns.Nodes)))

		if err := dec.Decode(&info); err == io.EOF {
			return nil, err
		} else if err != nil {
			log.Fatal(err)
		}

		for n, id := range ids {
			nodes[n] = info[id]
			nodes[n].Cluster_name = ns.Cluster_name
		}
	}

	return nodes, nil
}

func (eb *Elasticbeat) GetNodeIDs(u url.URL) ([]string, error) {
	res, err := http.Get(u.String() + "/_cat/nodes?full_id=true&h=id")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP%s", res.Status)
	}

	m := make(map[int]string)
	i := 0

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		m[i] = strings.TrimSpace(scanner.Text())
		i++
	}

	ids := make([]string, len(m))
	for key, val := range m {
		ids[key] = val
	}

	return ids, nil
}
