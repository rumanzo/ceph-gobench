package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
)

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

func MakeMonQuery(cephconn *Cephconnection, query map[string]string) []byte {
	monjson, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Can't marshal json mon query. Error: %v", err)
	}

	monrawanswer, _, err := cephconn.conn.MonCommand(monjson)
	if err != nil {
		log.Fatalf("Failed exec monCommand. Error: %v", err)
	}
	return monrawanswer
}

func GetPoolSize(cephconn *Cephconnection, params Params) Poolinfo {
	monrawanswer := MakeMonQuery(cephconn, map[string]string{"prefix": "osd pool get", "pool": params.pool,
		"format": "json", "var": "size"})
	monanswer := Poolinfo{}
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monanswer

}

func GetPgByPool(cephconn *Cephconnection, params Params) []PlacementGroup {
	monrawanswer := MakeMonQuery(cephconn, map[string]string{"prefix": "pg ls-by-pool", "poolstr": params.pool,
		"format": "json"})
	monanswer := []PlacementGroup{}
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monanswer
}

func GetOsdCrushDump(cephconn *Cephconnection) OsdCrushDump {
	monrawanswer := MakeMonQuery(cephconn, map[string]string{"prefix": "osd crush dump", "format": "json"})
	monanswer := OsdCrushDump{}
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monanswer
}

func GetOsdDump(cephconn *Cephconnection) OsdDump {
	monrawanswer := MakeMonQuery(cephconn, map[string]string{"prefix": "osd dump", "format": "json"})
	monanswer := OsdDump{}
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monanswer
}

func GetOsdLocations(params Params, osdcrushdump OsdCrushDump, osddump OsdDump, poolinfo Poolinfo) []int {
	var crushrule int64
	for _, pool := range osddump.Pools {
		if pool.Pool == poolinfo.PoolId {
			crushrule = pool.CrushRule
		}
	}
	var itemid int64
	for _, rule := range osdcrushdump.Rules {
		if rule.RuleID == crushrule {
			for _, step := range rule.Steps {
				if step.Op == "take" {
					itemid = step.Item
				}
			}
		}
	}
	log.Println(itemid)
	return []int{}
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
	poolinfo := GetPoolSize(cephconn, params)
	if poolinfo.Size != 1 {
		log.Fatalf("Pool size must be 1. Current size for pool %v is %v. Don't forget that it must be useless pool (not production). Do:\n # ceph osd pool set %v min_size 1\n # ceph osd pool set %v size 1",
			poolinfo.Pool, poolinfo.Size, poolinfo.Pool, poolinfo.Pool)
	}

	//placementGroups := GetPgByPool(cephconn, params)
	//for _, value := range placementGroups {
	//	log.Printf("%+v\n", value)
	//}
	crushosddump := GetOsdCrushDump(cephconn)
	osddump := GetOsdDump(cephconn)
	GetOsdLocations(params, crushosddump, osddump, poolinfo)
	log.Println(poolinfo)

}

//TODO Получить структуру пула (osd dump), рул. Получить все рулы (osd crush dump). Разобрать краш карту, получить список осд.
//TODO получить список PG, если не все находятся на нужных OSD выкинуть эксепшн.
