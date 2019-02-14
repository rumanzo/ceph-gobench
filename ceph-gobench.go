package main

import (
	"code.cloudfoundry.org/bytefmt"
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
	blocksize, err := bytefmt.ToBytes(params.bs)
	if err != nil {
		log.Println("Can't convert defined block size. 4K block size will be used\n")
		blocksize = 4096
	}
	log.Printf("%v\n", blocksize)
	var buffs [][]byte
	for i :=0; i < 2 * params.threads_count; i++ {
		buffs = append(buffs, make([]byte, blocksize))
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
