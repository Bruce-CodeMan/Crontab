package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		putOp  clientv3.Op
		getOp  clientv3.Op
		opResp clientv3.OpResponse
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
	}

	// Get the object of kv
	kv = clientv3.NewKV(client)

	// Creating op
	putOp = clientv3.OpPut("/cron/jobs/job8", "job_8")

	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Writing Revision:", opResp.Put().Header.Revision)

	// Creating op
	getOp = clientv3.OpGet("/cron/jobs/job8")

	// Executing Op
	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Reading Revision:", opResp.Get().Kvs[0])
}
