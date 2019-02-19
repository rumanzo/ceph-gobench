package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

func get_pool_size(cephconn *Cephconnection, params Params) Poolinfo {
	monjson, err := json.Marshal(map[string]string{"prefix": "osd pool get", "pool": params.pool,
		"format": "json", "var": "size"})
	if err != nil {
		log.Fatalf("Can't marshal json mon query. Error: %v", err)
	}

	monrawanswer, _, err := cephconn.conn.MonCommand(monjson)
	if err != nil {
		log.Fatalf("Failed exec monCommand. Error: %v", err)
	}

	monanswer := Poolinfo{}
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monanswer

}

func get_pg_by_pool(cephconn *Cephconnection, params Params) []PlacementGroup {
	monjson, err := json.Marshal(map[string]string{"prefix": "pg ls-by-pool", "poolstr": params.pool, "format": "json"})
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
	return monanswer
}

//func get_crush_dump(params Params, cephconn *Cephconnection) OsdCrushDump {
//
//}

//func get_acting_osd(params Params, ) map[string]float64 {
//
//}

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
	poolinfo := get_pool_size(cephconn, params)
	if poolinfo.Size != 1 {
		log.Fatalf("Pool size must be 1. Current size for pool %v is %v. Don't forget that it must be useless pool (not production). Do:\n # ceph osd pool set %v min_size 1\n # ceph osd pool set %v size 1",
			poolinfo.Pool, poolinfo.Size, poolinfo.Pool, poolinfo.Pool)
	}

	placementGroups := get_pg_by_pool(cephconn, params)
	for _, value := range placementGroups {
		log.Printf("%+v\n", value)
	}

}

//TODO Получить структуру пула (osd dump), рул. Получить все рулы (osd crush dump). Разобрать краш карту, получить список осд.
//TODO получить список PG, если не все находятся на нужных OSD выкинуть эксепшн.
