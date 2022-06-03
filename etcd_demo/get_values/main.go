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
		getResp *clientv3.GetResponse
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

	// key-value pairs used to write and read etcd
	kv = clientv3.NewKV(client)
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(getResp.Kvs)
	}
}
