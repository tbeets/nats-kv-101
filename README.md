# nats-kv-101
Basic kick the tires on NATS Key-Value API (Go)

# Usage
```bash
# Get
./mybucket -s "nats://vbox1.tinghus.net" -creds "/home/todd/lab/nats-cluster1/vault/.nkeys/creds/NatsOp/AcctA/UserA1.creds" bucket1 foo

# Put
./mybucket -s "nats://vbox1.tinghus.net" -creds "/home/todd/lab/nats-cluster1/vault/.nkeys/creds/NatsOp/AcctA/UserA1.creds" bucket1 foo bar

# Get with history
./mybucket -s "nats://vbox1.tinghus.net" -creds "/home/todd/lab/nats-cluster1/vault/.nkeys/creds/NatsOp/AcctA/UserA1.creds" --history bucket1 foo

# Watch
./mywatch -s "nats://vbox1.tinghus.net" -creds "/home/todd/lab/nats-cluster1/vault/.nkeys/creds/NatsOp/AcctA/UserA1.creds" bucket1 foo
```
# API Notes

```go
kv, err := js.KeyValue("bucket name")
```

| Type | Methods |
| --- | --- |
| KeyValue | Bucket, Delete, Create, Update, Put, Get, Status, History, Keys, Purge, PurgeAll, Watch, WatchAll | 
| KeyValueEntry | Key, Value, Revision, Operation, Delta, Bucket, Created |
| Status | Values, TTL, History, Bucket, BackingStore, StreamInfo (when cast to *nats.keyValueBucketStatus) |
| KeyWatcher | Stop, Updates |

Notes:
* A watcher (via Updates) initially returns the current KeyValueEntry (if any) then a nil KeyValueEntry. Subsequent entries represent changes since watcher instantiated.

# Resources
* [nats.go](https://github.com/nats-io/nats.go)
* test/kv_test.go



