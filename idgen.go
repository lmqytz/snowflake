package main

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"
)

const (
	WorkerIdBits   = 5
	BusinessIdBits = 7
	SequenceBits   = 10
	MaxSequence    = -1 ^ (-1 << 10)
	Timezone       = "Asia/Shanghai"
)

var (
	timezone, _    = time.LoadLocation(Timezone)
	since          = uint64(time.Date(2019, 8, 26, 0, 0, 0, 0, timezone).UnixNano() / 1e6)
	workerIdMask   = uint32(math.Pow(2, WorkerIdBits)) - 1
	businessIdMask = uint32(math.Pow(2, BusinessIdBits)) - 1
)

type IdGen struct {
	lastMilliSecond uint64
	mux             sync.Mutex
	workerId        uint32
	businessId      uint32
	sequence        uint32
}

type Fragment struct {
	Timestamp  uint64 `json:"timestamp"`
	WorkerId   uint32 `json:"worker_id"`
	BusinessId uint32 `json:"business_id"`
}

func (ig *IdGen) Next(businessId uint32) uint64 {
	ig.mux.Lock()
	defer ig.mux.Unlock()

	ig.businessId = businessId

	millisecond := currentMillisecond()
	if millisecond == ig.lastMilliSecond {
		ig.sequence = (ig.sequence + 1) & MaxSequence

		if ig.sequence == 0 {
			millisecond = nextMillisecond(millisecond)
		}
	} else {
		ig.sequence = 0
	}

	ig.lastMilliSecond = millisecond

	return (ig.lastMilliSecond << (WorkerIdBits + BusinessIdBits + SequenceBits)) |
		uint64(ig.workerId<<(BusinessIdBits+SequenceBits)) |
		uint64(ig.businessId<<SequenceBits) |
		uint64(ig.sequence)
}

func (ig *IdGen) Parse(id uint64) Fragment {
	fragment := Fragment{
		Timestamp:  id>>(WorkerIdBits+BusinessIdBits+SequenceBits) + since,
		WorkerId:   uint32(id>>(SequenceBits+BusinessIdBits)) & workerIdMask,
		BusinessId: uint32(id>>SequenceBits) & businessIdMask,
	}

	return fragment
}

func currentMillisecond() uint64 {
	return uint64(time.Now().UnixNano()/1e6) - since
}

func nextMillisecond(millisecond uint64) uint64 {
	c := currentMillisecond()
	for c < millisecond {
		c = currentMillisecond()
	}

	return c
}

func NewIdGen() (*IdGen, error) {
	max := int(math.Pow(2, WorkerIdBits))
	workerId, err := getWorkId(max)

	if err != nil {
		return nil, err
	}

	maxWorkId := uint32(math.Pow(2, WorkerIdBits))
	if workerId > maxWorkId {
		return nil, fmt.Errorf("workerId should not greater then " + strconv.Itoa(int(maxWorkId)))
	}

	idGen := &IdGen{
		workerId: workerId,
	}

	return idGen, nil
}
