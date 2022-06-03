package main

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
	)

	// Configure
	config = clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}

	// Creating a client
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Connected")
	}

	client = client
}
