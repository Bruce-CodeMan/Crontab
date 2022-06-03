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
		config        clientv3.Config
		client        *clientv3.Client
		err           error
		kv            clientv3.KV
		watcher       clientv3.Watcher
		getResp       *clientv3.GetResponse
		watchRevision int64
		watchRespChan <-chan clientv3.WatchResponse
		watchResp     clientv3.WatchResponse
		event         *clientv3.Event
		ctx           context.Context
		cancelFunc    context.CancelFunc
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

	go func() {
		for {
			kv.Put(context.TODO(), "name", "Bruce")
			kv.Delete(context.TODO(), "name")
			time.Sleep(1 * time.Second)
		}
	}()

	// First Get the current value
	if getResp, err = kv.Get(context.TODO(), "name"); err != nil {
		fmt.Println(err)
		return
	}

	if getResp.Count != 0 {
		fmt.Println(string(getResp.Kvs[0].Value))
	}

	// Revision version
	watchRevision = getResp.Header.Revision + 1

	// Creating a watcher
	watcher = clientv3.NewWatcher(client)

	// Starting watch
	fmt.Println("watch Revision version:", watchRevision)

	// cancel after 5 seconds
	ctx, cancelFunc = context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancelFunc()
	})

	watchRespChan = watcher.Watch(ctx, "name", clientv3.WithRev(watchRevision))
	for watchResp = range watchRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("[put], put:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("[delete] Revision:", event.Kv.ModRevision)
			}
		}
	}
}
