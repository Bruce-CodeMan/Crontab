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
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		ctx            context.Context
		cancelFunc     context.CancelFunc
		kv             clientv3.KV
		txn            clientv3.Txn
		txnResp        *clientv3.TxnResponse
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

	// 1. 上锁(创建租约, 自动续租, 拿着租约去抢占一个key)
	lease = clientv3.NewLease(client)

	// Apply for a 5-second lease
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}

	// Get the ID of lease
	leaseId = leaseGrantResp.ID

	// cancel the lease
	ctx, cancelFunc = context.WithCancel(context.TODO())

	// 确保函数退出后, 自动续租会停止
	defer cancelFunc()
	// 立即销毁租约
	defer lease.Revoke(context.TODO(), leaseId)

	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					fmt.Println("lease have expired")
					goto END
				} else {
					fmt.Println("Receive leaseId:", keepResp.ID)
				}
			}
		}
	END:
	}()

	// if不存在key, then 设置他, else抢锁失败
	kv = clientv3.NewKV(client)

	// 创建事务
	txn = kv.Txn(context.TODO())

	// 定义事务
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/jobs/job_10"), "=", 0)).
		Then(clientv3.OpPut("/cron/jobs/job_10", "job_9", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/jobs/job_10"))

	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}

	// 判断是否抢到了锁
	if !txnResp.Succeeded {
		fmt.Println("锁被占用:", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}
	// 2. 处理业务
	fmt.Println("处理业务")
	time.Sleep(500 * time.Second)

	// 在锁内, 很安全
	// 3. 释放锁(取消自动续租,释放租约)

}
