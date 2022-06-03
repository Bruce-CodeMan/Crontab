package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config         clientv3.Config
		client         *clientv3.Client
		err            error
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		kv             clientv3.KV
		getResp        *clientv3.GetResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
	)

	//  Configure
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

	// Apply for a lease
	lease = clientv3.Lease(client)

	// Apply for a 10-second lease
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}

	// Get the ID of lease
	leaseId = leaseGrantResp.ID

	// Automatic renewal
	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					fmt.Println("lease have expired")
					goto END
				} else {
					fmt.Println("renewal have success, leaseId:", keepResp.ID)
				}
			}
		}
	END:
	}()

	// Get the object of KV
	kv = clientv3.NewKV(client)

	// success, Add a key-value pair to associate it with the lease so that it automatically expires after 10 seconds
	if _, err = kv.Put(context.TODO(), "/cron/lock/job1", "lock_1", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("writing Successfully")

	// Periodically check whether the key has expired
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("key-value pair expires")
			break
		} else {
			fmt.Println("key-value pair not expire")
			fmt.Println(getResp.Kvs)
		}
		time.Sleep(1 * time.Second)
	}
}
