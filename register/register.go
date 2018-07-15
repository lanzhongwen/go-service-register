// Author: Zhongwen Lan(runbus@qq.com)
// Created: 2018/07/15
package register

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type Node struct {
	etcdAddrs    []string
	serviceName  string
	serviceValue string
	ttl          int64 // In seconds
	dialTimeout  int   // in seconds
	etcdClient   *clientv3.Client
	retryCh      chan error
}

func New(etcdAddrs []string, serviceName, serviceValue string, dialTimeout, ttl int) *Node {
	if len(etcdAddrs) == 0 || len(serviceName) == 0 || len(serviceValue) == 0 || dialTimeout == 0 || ttl == 0 {
		log.Println("[register/NewNode] Invalid Argument")
		return nil
	}
	n := &Node{
		etcdAddrs:    etcdAddrs,
		serviceName:  serviceName,
		serviceValue: serviceValue,
		ttl:          int64(ttl),
		dialTimeout:  dialTimeout,
		retryCh:      make(chan error, 1),
	}

	go n.retry()

	return n
}

func (n *Node) retry() {
	log.Println("[register/retry]")
	for {
		select {
		case <-n.retryCh:
			log.Println("[register/retry] Receiving from retry channel")
			err := n.Register()
			if err != nil {
				log.Println("[register/retry] Register failed: " + err.Error())
			} else {
				log.Println("[register/retry] Register Successfully")
			}
			// TODO: chan close action
		}
	}
}

// TODO: check DailKeepAliveTime
func (n *Node) Register() error {
	var err error
	n.etcdClient, err = clientv3.New(
		clientv3.Config{
			Endpoints:   n.etcdAddrs,
			DialTimeout: time.Duration(n.dialTimeout) * time.Second,
		},
	)

	if err != nil {
		log.Println("[register/Dial] Connect to etcd failed: " + err.Error())
		n.retryCh <- err
		return err
	}

	err = n.register()
	if err != nil {
		log.Println("[register/Register] register failed: " + err.Error())
		n.retryCh <- err
		return err
	}

	go n.watch()

	log.Println("[register/Register] Register Successfully")

	return nil
}

func (n *Node) register() error {
	resp, err := n.etcdClient.Grant(context.Background(), n.ttl)
	if err != nil {
		log.Println("[register/register] Grant failed: " + err.Error())
		return err
	}

	_, err = n.etcdClient.Put(
		context.Background(),
		n.serviceName,
		n.serviceValue,
		clientv3.WithLease(resp.ID),
	)
	if err != nil {
		log.Println("[register/register] Put failed: " + err.Error())
		return err
	}

	keepCh, err := n.etcdClient.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		log.Println("[register/register] KeepAlive failed: " + err.Error())
		return err
	}

	_ = <-keepCh

	log.Println("Register service: ", n.serviceName, " | value: ", n.serviceValue, " | ttl: ", n.ttl)

	return nil
}

func (n *Node) watch() {
	log.Println("[register/watch] " + n.serviceName)
	watchCh := n.etcdClient.Watch(context.Background(), n.serviceName)
EXIT:
	for {
		select {
		case watchResp := <-watchCh:
			if len(watchResp.Events) == 0 {
				log.Println("len(watchResp.Events) == 0: Going to exit watch...")
				break EXIT
			}
			IsDeleted := false
			for _, ev := range watchResp.Events {
				if ev.Type == clientv3.EventTypeDelete {
					IsDeleted = true
					break
				}
			}
			if !IsDeleted {
				continue
			}

			log.Println("[register/watch] Detected deletion event for " + n.serviceName)

			err := n.register()
			if err != nil {
				log.Println("[register/watch] registerService failed: " + err.Error())
				break EXIT
			}
		}
	}

	log.Println("[register/watch] Exiting watch for " + n.serviceName + "...")

	n.retryCh <- errors.New("Exit Watcher")
}
