package main

import (
	"encoding/json"
	"log"
	"strings"
)

func makeMonQuery(cephconn *cephConnection, query map[string]string) []byte {
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

func getPoolSize(cephconn *cephConnection, params params) poolinfo {
	monrawanswer := makeMonQuery(cephconn, map[string]string{"prefix": "osd pool get", "pool": params.pool,
		"format": "json", "var": "size"})
	monanswer := Poolinfo{}
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monanswer

}

func getPgByPool(cephconn *cephConnection, params params) []placementGroup {
	monrawanswer := makeMonQuery(cephconn, map[string]string{"prefix": "pg ls-by-pool", "poolstr": params.pool,
		"format": "json"})
	var monanswer []placementGroup
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monanswer
}

func getOsdCrushDump(cephconn *cephConnection) osdCrushDump {
	monrawanswer := makeMonQuery(cephconn, map[string]string{"prefix": "osd crush dump", "format": "json"})
	var monanswer osdCrushDump
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monanswer
}

func getOsdDump(cephconn *cephConnection) osdDump {
	monrawanswer := makeMonQuery(cephconn, map[string]string{"prefix": "osd dump", "format": "json"})
	var monanswer osdDump
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monanswer
}

func getOsdMetadata(cephconn *cephConnection) []osdMetadata {
	monrawanswer := makeMonQuery(cephconn, map[string]string{"prefix": "osd metadata", "format": "json"})
	var monanswer []osdMetadata
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monanswer
}

func getObjActingPrimary(cephconn *cephConnection, params params, objname string) int64 {
	monrawanswer := makeMonQuery(cephconn, map[string]string{"prefix": "osd map", "pool": params.pool,
		"object": objname, "format": "json"})
	var monanswer osdMap
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monanswer.UpPrimary
}

func getCrushHostBuckets(buckets []bucket, itemid int64) []bucket {
	var rootbuckets []bucket
	for _, bucket := range buckets {
		if bucket.ID == itemid {
			if bucket.TypeName == "host" {
				rootbuckets = append(rootbuckets, bucket)
			} else {
				for _, item := range bucket.Items {
					result := getCrushHostBuckets(buckets, item.ID)
					for _, it := range result {
						rootbuckets = append(rootbuckets, it)
					}
				}
			}
		}
	}
	return rootbuckets
}

func getOsdForLocations(params params, osdcrushdump osdCrushDump, osddump osdDump, poolinfo poolinfo, osdsmetadata []osdMetadata) []device {
	var crushrule, rootid int64
	var crushrulename string
	for _, pool := range osddump.Pools {
		if pool.Pool == poolinfo.PoolId {
			crushrule = pool.CrushRule
		}
	}
	for _, rule := range osdcrushdump.Rules {
		if rule.RuleID == crushrule {
			crushrulename = rule.RuleName
			for _, step := range rule.Steps {
				if step.Op == "take" {
					rootid = step.Item
				}
			}
		}
	}
	osdstats := map[uint64]*Osd{}
	for num, stat := range osddump.Osds {
		osdstats[stat.Osd] = &osddump.Osds[num]
	}

	var osddevices []device
	bucketitems := getCrushHostBuckets(osdcrushdump.Buckets, rootid)
	if params.define != "" {
		if strings.HasPrefix(params.define, "osd.") {
			for _, hostbucket := range bucketitems {
				for _, item := range hostbucket.Items {
					for _, device := range osdcrushdump.Devices {
						if device.ID == item.ID && params.define == device.Name {
							for _, osdmetadata := range osdsmetadata {
								if osdmetadata.ID == device.ID && osdstats[uint64(device.ID)].Up == 1 && osdstats[uint64(device.ID)].In == 1 {
									device.Info = osdmetadata
									osddevices = append(osddevices, device)
								}
							}
						}
					}
				}
			}
			if len(osddevices) == 0 {
				log.Fatalf("Defined osd not exist in root for rule: %v pool: %v.\nYou should define osd like osd.X",
					crushrulename, poolinfo.Pool)
			}
		} else {
			for _, hostbucket := range bucketitems {
				if strings.Split(hostbucket.Name, "~")[0] == strings.Split(params.define, "~")[0] { //purge device class
					for _, item := range hostbucket.Items {
						for _, device := range osdcrushdump.Devices {
							if device.ID == item.ID {
								for _, osdmetadata := range osdsmetadata {
									if osdmetadata.ID == device.ID && osdstats[uint64(device.ID)].Up == 1 && osdstats[uint64(device.ID)].In == 1 {
										device.Info = osdmetadata
										osddevices = append(osddevices, device)
									}

								}
							}
						}
					}
				}
			}
			if len(osddevices) == 0 {
				log.Fatalf("Defined host not exist in root for rule: %v pool: %v", crushrulename, poolinfo.Pool)
			}
		}
	} else {
		for _, hostbucket := range bucketitems {
			for _, item := range hostbucket.Items {
				for _, device := range osdcrushdump.Devices {
					if device.ID == item.ID {
						for _, osdmetadata := range osdsmetadata {
							if osdmetadata.ID == device.ID && osdstats[uint64(device.ID)].Up == 1 && osdstats[uint64(device.ID)].In == 1 {
								device.Info = osdmetadata
								osddevices = append(osddevices, device)
							}
						}
					}
				}
			}
		}
		if len(osddevices) == 0 {
			log.Fatalf("Osd not exist in root for rule: %v pool: %v", crushrulename, poolinfo.Pool)
		}
	}
	return osddevices
}

func containsPg(pgs []placementGroup, i int64) bool {
	for _, pg := range pgs {
		if i == pg.ActingPrimary {
			return true
		}
	}
	return false
}

func getOsds(cephconn *cephConnection, params params) []device {
	poolinfo := getPoolSize(cephconn, params)
	if poolinfo.Size != 1 {
		log.Fatalf("Pool size must be 1. Current size for pool %v is %v. Don't forget that it must be useless pool (not production). Do:\n # ceph osd pool set %v min_size 1\n # ceph osd pool set %v size 1",
			poolinfo.Pool, poolinfo.Size, poolinfo.Pool, poolinfo.Pool)
	}
	placementGroups := getPgByPool(cephconn, params)
	crushosddump := getOsdCrushDump(cephconn)
	osddump := getOsdDump(cephconn)
	osdsmetadata := getOsdMetadata(cephconn)
	osddevices := getOsdForLocations(params, crushosddump, osddump, poolinfo, osdsmetadata)
	for _, device := range osddevices {
		if exist := containsPg(placementGroups, device.ID); exist == false {
			log.Fatalln("Not enough pg for test. Some osd haven't placement group at all. Increase pg_num and pgp_num")
		}
	}
	return osddevices
}
