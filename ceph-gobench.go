package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

//future feature
func makeoffsets(threads int64, bs int64, objsize int64) [][]int64 {
	var offsets [][]int64
	for i := int64(0); i < threads; i++ {
		s1 := rand.NewSource(i)
		r1 := rand.New(s1)
		localoffsets := make([]int64, 0, objsize-bs)
		for i := int64(0); i < objsize-bs; i += bs {
			localoffsets = append(localoffsets, i)
		}
		r1.Shuffle(len(localoffsets), func(i, j int) {
			localoffsets[i], localoffsets[j] = localoffsets[j], localoffsets[i]
		})
		offsets = append(offsets, localoffsets)
	}
	return offsets
}

func bench(cephconn *Cephconnection, osddevice Device, host string, buffs *[][]byte, offset [][]int64, params *Params,
	wg *sync.WaitGroup, ready chan bool, result chan string, allready uint64) {
	var nwg sync.WaitGroup
	tready := make(chan bool, params.threadsCount)
	start := uint64(0)
	defer wg.Done()
	for i := int64(0); i < params.threadsCount; i++ {
		nwg.Add(1)
		go _bench(cephconn, osddevice, host, buffs, offset[i], params, &nwg, i, tready, result, &start)
		if params.parallel != true {
			nwg.Wait()
		}
	}
	nwg.Wait()
}

func _bench(cephconn *Cephconnection, osddevice Device, host string, buffs *[][]byte, offset []int64, params *Params,
	wg *sync.WaitGroup, i int64, ready chan bool, result chan string, start *uint64) {
	defer wg.Done()
	time.Sleep(time.Second * time.Duration(i)) // prepare objects
	ready <- true
	for {
		if *start == 1 {
			break
		}
	}
	log.Println(host, i, osddevice.Name) //somework
	result <- fmt.Sprintf("Host: %v\nOsdname: %v", host, osddevice.Name)
}

func main() {
	params := Route()
	cephconn := connectioninit(params)
	defer cephconn.conn.Shutdown()

	// https://tracker.ceph.com/issues/24114
	time.Sleep(time.Millisecond * 100)

	var buffs [][]byte
	for i := int64(0); i < 2*params.threadsCount; i++ {
		buffs = append(buffs, make([]byte, params.blocksize))
	}
	for num := range buffs {
		_, err := rand.Read(buffs[num])
		if err != nil {
			log.Fatalln(err)
		}
	}
	osddevices := GetOsds(cephconn, params)
	offsets := makeoffsets(params.threadsCount, params.blocksize, params.objectsize)

	var wg sync.WaitGroup
	var allready uint64
	var ready chan bool
	var result chan string
	for host, osds := range osddevices {
		for _, osd := range osds {
			wg.Add(1)
			if params.parallel == true {
				go bench(cephconn, osd, host, &buffs, offsets, &params, &wg, ready, result, allready)
			} else {
				bench(cephconn, osd, host, &buffs, offsets, &params, &wg, ready, result, allready)
			}
		}
	}
	wg.Wait()

}
