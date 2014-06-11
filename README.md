## About

docker-metricsd provides information about docker's containers.

## Installation

```no-highlight
$ git clone https://github.com/maebashi/docker-metricsd.git
$ cd docker-metricsd
$ go get -d -v
$ go build docker-metricsd.go
$ sudo ./docker-metricsd
```

## Usage

```no-highlight
$ curl http://host.example.com:12375/containers/f9093f3df8d8/json | jq .
{
  "ID": "f9093f3df8d84db914f544139d828e3c7ba7f5a353f83220c0580536b61ec5c3",
  ...
  "Metrics": {
    "memory": {
      "usage_in_bytes": 91897856,
      "unevictable": 0,
      "total_unevictable": 0,
      "total_swap": 0,
      "total_rss": 348160,
      "total_pgpgout": 64116,
      "total_pgpgin": 78601,
      "total_mapped_file": 1376256,
      "total_inactive_file": 28758016,
      "total_inactive_anon": 0,
      "mapped_file": 1376256,
      "inactive_file": 28758016,
      "inactive_anon": 0,
      "hierarchical_memsw_limit": 9223372036854776000,
      "hierarchical_memory_limit": 9223372036854776000,
      "cache": 90378240,
      "active_file": 61620224,
      "active_anon": 348160,
      "max_usage_in_bytes": 125992960,
      "pgpgin": 78601,
      "pgpgout": 64116,
      "rss": 348160,
      "swap": 0,
      "total_active_anon": 348160,
      "total_active_file": 61620224,
      "total_cache": 90378240
    },
    "interfaces": {
      "outpackets.0": 5441,
      "inbytes.0": 16768513,
      "indrop.0": 0,
      "inerrs.0": 0,
      "inpackets.0": 14401,
      "name.0": "eth0",
      "outbytes.0": 382985,
      "outdrop.0": 0,
      "outerrs.0": 0
    },
    "cpuacct": {
      "usage": 0,
      "percentage": 0
    }
  },
  ...
}
```
