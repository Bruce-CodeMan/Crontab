package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		delResp *clientv3.DeleteResponse
		kvPair  *mvccpb.KeyValue
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

	// 删除
	if delResp, err = kv.Delete(context.TODO(), "name", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
		return
	}
	// 查看被删除之前的value
	if len(delResp.PrevKvs) != 0 {
		for _, kvPair = range delResp.PrevKvs {
			fmt.Println("删除了", string(kvPair.Key), string(kvPair.Value))
		}
	}
}
