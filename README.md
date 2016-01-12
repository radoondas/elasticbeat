# ElasticBeat
Current status: **development release**.

## Description
Simple [beat](https://github.com/elastic/beats) for ElasticSearch cluster (or nodes) statistics. This beat requests statistics from your elastisearch cluster via API.
Following API is currently supported:
 * /_cluster/health
 * /_cluster/stats
 
## To apply ElasticBeat template:

```bash
curl -XPUT 'http://localhost:9200/_template/elasticbeat' -d@elasticbeat.template.json
```
