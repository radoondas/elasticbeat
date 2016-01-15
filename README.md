# ElasticBeat
Current status: **development release**. Further development and fine-tuning is in progress.

**If you have any ideas or issues, please contribute.**

## Description
Simple [beat](https://github.com/elastic/beats) for ElasticSearch cluster (or nodes) statistics. This beat requests statistics from your elastisearch cluster via available API.
Following API is currently supported:
 * /_cluster/health
 * /_cluster/stats
 * /_nodes/stats/process,jvm,os,fs,thread_pool,transport,http,breaker,script
 
## ElasticBeat template:

```bash
curl -XPUT 'http://localhost:9200/_template/elasticbeat' -d@elasticbeat.template.json
```

## Note
In order to have **consistent** Node statistics, you need to set unique node names for each node of your elasticsearch cluster. 

In configuration file elasticsearch.yml:
```
node.name: testnode
```