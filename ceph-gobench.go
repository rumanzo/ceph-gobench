package main

import (
	"fmt"
	"log"
	"math/rand"
)

func main() {
	params := Route()
	cephconn := connectioninit(&params)
	defer cephconn.conn.Shutdown()

	stats, _ := cephconn.ioctx.GetPoolStats()
	log.Println(stats)

	log.Printf("%v\n", params.blocksize)
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
	fmt.Printf("%+v\n", buffs)
	fmt.Println(len(buffs))

}
