# etcdnew

Regenerates etcd discovery urls

## Installation

```
$ go get github.com/aldrinleal/etcdnew
```

## Usage

Run etcd with new options.

## How it works?

It rewrites a file with a given URL (if empty, it will generate one on https://discovery.etcd.io/new)

If without a set URL, it will apply a new entry to all files passed, thus making it easy to renew etcd entries for multiple clusters at once.