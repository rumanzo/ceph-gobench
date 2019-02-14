package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

type Moncommand struct {
	Prefix string `json:"prefix"`
	Pool   string `json:"pool"`
	Format string `json:"format"`
	Var    string `json:"var"`
}
type Monanswer struct {
	Pool   string `json:"pool"`
	PoolId int    `json:"pool_id"`
	Size   int    `json:"size"`
}

func main() {
	params := Route()
	cephconn := connectioninit(&params)
	defer cephconn.conn.Shutdown()

	// https://tracker.ceph.com/issues/24114
	time.Sleep(time.Millisecond * 100)

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

	monjson, err := json.Marshal(Moncommand{Prefix: "osd pool get", Pool: params.pool, Format: "json", Var: "size"})
	if err != nil {
		log.Fatalf("Can't marshal json mon query. Error: %v", err)
	}

	monrawanswer, _, err := cephconn.conn.MonCommand(monjson)
	if err != nil {
		log.Fatalf("Failed exec monCommand. Error: %v", err)
	}

	monanswer := Monanswer{}
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}

	if monanswer.Size != 1 {
		log.Fatalf("Pool size must be 1. Current size for pool %v is %v. Don't forget that it must be useless pool (not production). Do:\n # ceph osd pool set %v min_size 1\n # ceph osd pool set %v size 1",
			monanswer.Pool, monanswer.Size, monanswer.Pool, monanswer.Pool)
	}

}
