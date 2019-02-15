package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

func get_pool_size(cephconn *Cephconnection, params *Params) *Monanswer {
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
	return &monanswer

}

func get_pg_by_pool(cephconn *Cephconnection, params *Params) *[]PlacementGroup {
	monjson, err := json.Marshal(Moncommand{Prefix: "pg ls-by-pool", Poolstr: params.pool, Format: "json"})
	if err != nil {
		log.Fatalf("Can't marshal json mon query. Error: %v", err)
	}
	monrawanswer, _, err := cephconn.conn.MonCommand(monjson)
	if err != nil {
		log.Fatalf("Failed exec monCommand. Error: %v", err)
	}

	monanswer := []PlacementGroup{}
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return &monanswer
}

func main() {
	params := Route()
	cephconn := connectioninit(params)
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
	monanswersize := get_pool_size(cephconn, params)
	if monanswersize.Size != 1 {
		log.Fatalf("Pool size must be 1. Current size for pool %v is %v. Don't forget that it must be useless pool (not production). Do:\n # ceph osd pool set %v min_size 1\n # ceph osd pool set %v size 1",
			monanswersize.Pool, monanswersize.Size, monanswersize.Pool, monanswersize.Pool)
	}

	placementGroups := get_pg_by_pool(cephconn, params)
	log.Println(placementGroups)

}
