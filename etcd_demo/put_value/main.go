package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		putResp *clientv3.PutResponse
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

	// Key-value pairs used to write and read etcd
	kv = clientv3.NewKV(client)
	if putResp, err = kv.Put(context.TODO(), "name", "Bruce"); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(putResp.Header.Revision)
	}

}
