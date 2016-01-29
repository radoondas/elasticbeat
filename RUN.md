# Short RUN guide

First, setup Go lang environment (https://golang.org/doc/install)
Add to your .bashrc important variables

```bash
export GOROOT="$HOME/opt/go"
export GOPATH="$HOME/workspace/go"
export PATH="$HOME/opt/go/bin:$PATH"
```

## Install ElasticBeat and dependencies

```bash
go get -insecure gopkg.in/yaml.v2
go get github.com/radoondas/elasticbeat
```

## Elastic and Kibana
Meanwhile setup your ElasticSearch and Kibana (example [dashbords](https://github.com/radoondas/elasticbeat/tree/master/kibana))

## Build ElasticBeat

```bash
cd ~/workspace/go/src/github.com/radoondas/elasticbeat
go install
```

## Delete template (Optional)
If you need for any reason to delete old template, use following method.

```bash
curl -XDELETE 'http://localhost:9200/_template/elasticbeat'
```

## Import template
```bash
cd ~/workspace/go/src/github.com/radoondas/elasticbeat/etc
curl -XPUT 'http://localhost:9200/_template/elasticbeat' -d@elasticbeat.template.json
```

## Run ElasticBeat

Following command will execute ElasticBeat with debug option and will not index results in to ES. Instead, you will see output on the screen.
```bash
cd ~/workspace/go/bin
./elasticbeat  -e -v -d elasticbeat -c ~/workspace/go/src/github.com/radoondas/elasticbeat/elasticbeat.yml
```

With no debug options - just do straight indexing to your ES installation

```bash
./elasticbeat  -e -c ~/workspace/go/src/github.com/radoondas/elasticbeat/elasticbeat.yml
```
