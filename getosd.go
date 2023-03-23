package main

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
)

func makeMonQuery(cephconn *cephconnection, query map[string]string) []byte {
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

func getPoolSize(cephconn *cephconnection, params params) Poolinfo {
	monrawanswer := makeMonQuery(cephconn, map[string]string{"prefix": "osd pool get", "pool": params.pool,
		"format": "json", "var": "size"})
	monanswer := Poolinfo{}
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monanswer

}

func getPgByPool(cephconn *cephconnection, params params) []PlacementGroup {
	monrawanswer := makeMonQuery(cephconn, map[string]string{"prefix": "pg ls-by-pool", "poolstr": params.pool,
		"format": "json"})
	var monanswer []PlacementGroup
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		//try Nautilus
		var nmonanswer placementGroupNautilus
		if nerr := json.Unmarshal([]byte(monrawanswer), &nmonanswer); nerr != nil {
			log.Fatalf("Can't parse monitor answer in getPgByPool. Error: %v", err)
		}
		return nmonanswer.PgStats
	}
	return monanswer
}

func getOsdCrushDump(cephconn *cephconnection) OsdCrushDump {
	monrawanswer := makeMonQuery(cephconn, map[string]string{"prefix": "osd crush dump", "format": "json"})
	var monanswer OsdCrushDump
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monanswer
}

func getOsdDump(cephconn *cephconnection) OsdDump {
	monrawanswer := makeMonQuery(cephconn, map[string]string{"prefix": "osd dump", "format": "json"})
	var monanswer OsdDump
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monanswer
}

func getOsdMetadata(cephconn *cephconnection) []OsdMetadata {
	monrawanswer := makeMonQuery(cephconn, map[string]string{"prefix": "osd metadata", "format": "json"})
	var monanswer []OsdMetadata
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monanswer
}

func getObjActingPrimary(cephconn *cephconnection, params params, objname string) int64 {
	monrawanswer := makeMonQuery(cephconn, map[string]string{"prefix": "osd map", "pool": params.pool,
		"object": objname, "format": "json"})
	var monanswer OsdMap
	if err := json.Unmarshal([]byte(monrawanswer), &monanswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monanswer.UpPrimary
}

func getCrushHostBuckets(buckets []Bucket, itemid int64) []Bucket {
	var rootbuckets []Bucket
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

func getOsdForLocations(params params, osdcrushdump OsdCrushDump, osddump OsdDump, poolinfo Poolinfo, osdsmetadata []OsdMetadata) []Device {
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

	var osddevices []Device
	bucketitems := getCrushHostBuckets(osdcrushdump.Buckets, rootid)

	if params.rdefine != "" { // match regex if exists
		validbucket, err := regexp.CompilePOSIX(params.rdefine)
		if err != nil {
			log.Fatalf("Can't parse regex %v", params.rdefine)
		}
		for _, hostbucket := range bucketitems {
			for _, item := range hostbucket.Items {
				for _, device := range osdcrushdump.Devices {
					if device.ID == item.ID && (validbucket.MatchString(hostbucket.Name) || validbucket.MatchString(device.Name)) {
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
			log.Fatalf("Defined host/osd not exist in root for rule: %v pool: %v", crushrulename, poolinfo.Pool)
		}
	} else if params.define != "" { // check defined osd/hosts
		if strings.HasPrefix(params.define, "osd.") { //check that defined is osd, else host
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

func containsPg(pgs []PlacementGroup, i int64) bool {
	for _, pg := range pgs {
		if i == int64(pg.ActingPrimary) {
			return true
		}
	}
	return false
}

func getOsds(cephconn *cephconnection, params params) []Device {
	poolinfo := getPoolSize(cephconn, params)
	if params.disablecheck == false {
		if poolinfo.Size != 1 {
			log.Fatalf("Pool size must be 1. Current size for pool %v is %v. Don't forget that it must be useless pool (not production). Do:\n # ceph osd pool set %v min_size 1\n # ceph osd pool set %v size 1",
				poolinfo.Pool, poolinfo.Size, poolinfo.Pool, poolinfo.Pool)
		}
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
