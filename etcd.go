package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"strconv"
	"time"
)

func getWorkId(max int) (uint32, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.EtcdAddr,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		return 0, err
	}

	defer func() {
		if err := cli.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	resp, err := cli.Get(ctx, "/idGen/workerId", clientv3.WithPrefix())
	cancel()

	if err != nil {
		return 0, err
	}

	var ids []int
	for _, ev := range resp.Kvs {
		id, _ := strconv.Atoi(string(ev.Value))
		ids = append(ids, id)
	}

	var availableId int
	for i := 1; i <= max; i++ {
		nex := false

		for _, k := range ids {
			if k == i {
				nex = true
			}
		}

		if !nex {
			availableId = i
			break
		}
	}

	return uint32(availableId), nil
}
