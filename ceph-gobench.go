package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

//future feature
func makepreudorandom() {
	a := make([]int, 0, 4096/4)
	for i := 0; i < 4096; i += 4 {
		a = append(a, i)
	}
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
	fmt.Println(a)
}

func main() {
	params := Route()
	cephconn := connectioninit(params)
	defer cephconn.conn.Shutdown()

	// https://tracker.ceph.com/issues/24114
	time.Sleep(time.Millisecond * 100)

	var buffs [][]byte
	for i := 0; i < 2*params.threadsCount; i++ {
		buffs = append(buffs, make([]byte, params.blocksize))
	}
	for num := range buffs {
		_, err := rand.Read(buffs[num])
		if err != nil {
			log.Fatalln(err)
		}
	}

	GetOsds(cephconn, params)
}
