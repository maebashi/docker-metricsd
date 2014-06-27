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
      "failcnt": 0,
      "stats": {
        "unevictable": 0,
        "total_unevictable": 0,
        "total_swap": 0,
        "total_rss": 380928,
        "total_pgpgout": 681084,
        "total_pgpgin": 697086,
        "total_mapped_file": 1433600,
        "total_inactive_file": 38936576,
        "mapped_file": 1433600,
        "inactive_file": 38936576,
        "inactive_anon": 0,
        "hierarchical_memsw_limit": 9223372036854776000,
        "hierarchical_memory_limit": 9223372036854776000,
        "cache": 102838272,
        "active_file": 63901696,
        "active_anon": 380928,
        "pgpgin": 697086,
        "pgpgout": 681084,
        "rss": 380928,
        "swap": 0,
        "total_active_anon": 380928,
        "total_active_file": 63901696,
        "total_cache": 102838272,
        "total_inactive_anon": 0
      },
      "max_usage": 139268096,
      "usage": 104189952
    },
    "interfaces": {
      "outpackets.0": 3021223,
      "inbytes.0": 8228044607,
      "indrop.0": 0,
      "inerrs.0": 0,
      "inpackets.0": 6429687,
      "name.0": "eth0",
      "outbytes.0": 199687042,
      "outdrop.0": 0,
      "outerrs.0": 0
    },
    "cpuacct": {
      "throlling_data": {},
      "cpu_usage": {
        "usage_in_usermode": 2.668e+10,
        "usage_in_kernelmode": 8.181e+10,
        "percpu_usage": [
          22599918565,
          987624379,
          65146098,
          36705600,
          18221767943,
          5890326,
          22768005795,
          4033968,
          218211598,
          4302652,
          5437296326,
          23781278563,
          18992360915,
          18712487134,
          55742891,
          18522722054
        ]
      }
    },
    "blkio": {}
  },
  ...
}
```
