// Copyright 2012-2021 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func usage() {
	log.Printf("Usage: mybucket [-s server] [-creds file] [-nkey file] [-tlscert file] [-tlskey file] [-tlscacert file] <bucket> <key> <value>\n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var userCreds = flag.String("creds", "", "User Credentials File")
	var nkeyFile = flag.String("nkey", "", "NKey Seed File")
	var tlsClientCert = flag.String("tlscert", "", "TLS client certificate file")
	var tlsClientKey = flag.String("tlskey", "", "Private key file for client certificate")
	var tlsCACert = flag.String("tlscacert", "", "CA certificate to verify peer against")
	var showHelp = flag.Bool("h", false, "Show help message")
	var showHistory = flag.Bool("history", false, "Show key history")

	var err error

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	args := flag.Args()
	if len(args) < 2 || len(args) > 3 {
		showUsageAndExit(1)
	}

	var bucket string = args[0]
	var key string = args[1]
	var value string = ""

	putMode := false
	if len(args) == 3 {
		putMode = true
		value = args[2]
	}

	// Connect Options.
	opts := []nats.Option{nats.Name("mybucket")}

	if *userCreds != "" && *nkeyFile != "" {
		log.Fatal("specify -seed or -creds")
	}

	// Use UserCredentials
	if *userCreds != "" {
		opts = append(opts, nats.UserCredentials(*userCreds))
	}

	// Use TLS client authentication
	if *tlsClientCert != "" && *tlsClientKey != "" {
		opts = append(opts, nats.ClientCert(*tlsClientCert, *tlsClientKey))
	}

	// Use specific CA certificate
	if *tlsCACert != "" {
		opts = append(opts, nats.RootCAs(*tlsCACert))
	}

	// Use Nkey authentication.
	if *nkeyFile != "" {
		opt, err := nats.NkeyOptionFromSeed(*nkeyFile)
		if err != nil {
			log.Fatal(err)
		}
		opts = append(opts, opt)
	}

	// Connect to NATS
	nc, err := nats.Connect(*urls, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, err := js.KeyValue(bucket)
	if err != nil {
		log.Fatal(err)
	}

	var entry nats.KeyValueEntry
	if putMode {
		_, err = kv.PutString(key, value)
	} else {
		entry, err = kv.Get(key)
	}
	if err != nil {
		log.Fatal(err)
	}

	var histEntries []nats.KeyValueEntry
	if *showHistory {
		histEntries, err = kv.History(key)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("History [%s][%s]:\n", kv.Bucket(), key )
		for _, histEntry := range histEntries {
			log.Printf("[%v][%s]\n", histEntry.Created(), histEntry.Value())
		}
	}

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		if putMode {
			log.Printf("[%s]: Put [%s] in [%s]\n", bucket, value, key)
	    } else {
			log.Printf("[%s]: Get [%s] returns [%s]\n", bucket, key, entry.Value())
		}
	}
}

