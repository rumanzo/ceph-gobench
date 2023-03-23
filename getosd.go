package main

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
)

func makeMonQuery(cephConn *cephconnection, query map[string]string) []byte {
	monJson, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Can't marshal json mon query. Error: %v", err)
	}

	monRawAnswer, _, err := cephConn.conn.MonCommand(monJson)
	if err != nil {
		log.Fatalf("Failed exec monCommand. Error: %v", err)
	}
	return monRawAnswer
}

func getPoolSize(cephConn *cephconnection, params params) Poolinfo {
	monRawAnswer := makeMonQuery(cephConn, map[string]string{"prefix": "osd pool get", "pool": params.pool,
		"format": "json", "var": "size"})
	monAnswer := Poolinfo{}
	if err := json.Unmarshal([]byte(monRawAnswer), &monAnswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monAnswer

}

func getPgByPool(cephConn *cephconnection, params params) []PlacementGroup {
	monRawAnswer := makeMonQuery(cephConn, map[string]string{"prefix": "pg ls-by-pool", "poolstr": params.pool,
		"format": "json"})
	var monAnswer []PlacementGroup
	if err := json.Unmarshal([]byte(monRawAnswer), &monAnswer); err != nil {
		//try Nautilus
		var nMonAnswer placementGroupNautilus
		if nerr := json.Unmarshal([]byte(monRawAnswer), &nMonAnswer); nerr != nil {
			log.Fatalf("Can't parse monitor answer in getPgByPool. Error: %v", err)
		}
		return nMonAnswer.PgStats
	}
	return monAnswer
}

func getOsdCrushDump(cephConn *cephconnection) OsdCrushDump {
	monrAwanswer := makeMonQuery(cephConn, map[string]string{"prefix": "osd crush dump", "format": "json"})
	var monAnswer OsdCrushDump
	if err := json.Unmarshal([]byte(monrAwanswer), &monAnswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monAnswer
}

func getOsdDump(cephConn *cephconnection) OsdDump {
	monRawAnswer := makeMonQuery(cephConn, map[string]string{"prefix": "osd dump", "format": "json"})
	var monAnswer OsdDump
	if err := json.Unmarshal([]byte(monRawAnswer), &monAnswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monAnswer
}

func getOsdMetadata(cephConn *cephconnection) []OsdMetadata {
	monRawAnswer := makeMonQuery(cephConn, map[string]string{"prefix": "osd metadata", "format": "json"})
	var monAnswer []OsdMetadata
	if err := json.Unmarshal([]byte(monRawAnswer), &monAnswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monAnswer
}

func getObjActingPrimary(cephConn *cephconnection, params params, objName string) int64 {
	monRawAnswer := makeMonQuery(cephConn, map[string]string{"prefix": "osd map", "pool": params.pool,
		"object": objName, "format": "json"})
	var monAnswer OsdMap
	if err := json.Unmarshal([]byte(monRawAnswer), &monAnswer); err != nil {
		log.Fatalf("Can't parse monitor answer. Error: %v", err)
	}
	return monAnswer.UpPrimary
}

func getCrushHostBuckets(buckets []Bucket, itemId int64) []Bucket {
	var rootBuckets []Bucket
	for _, bucket := range buckets {
		if bucket.ID == itemId {
			if bucket.TypeName == "host" {
				rootBuckets = append(rootBuckets, bucket)
			} else {
				for _, item := range bucket.Items {
					result := getCrushHostBuckets(buckets, item.ID)
					for _, it := range result {
						rootBuckets = append(rootBuckets, it)
					}
				}
			}
		}
	}
	return rootBuckets
}

func getOsdForLocations(params params, osdCrushDump OsdCrushDump, osdDump OsdDump, poolInfo Poolinfo, osdMetadata []OsdMetadata) []Device {
	var crushRule, rootId int64
	var crushRuleName string
	for _, pool := range osdDump.Pools {
		if pool.Pool == poolInfo.PoolId {
			crushRule = pool.CrushRule
		}
	}
	for _, rule := range osdCrushDump.Rules {
		if rule.RuleID == crushRule {
			crushRuleName = rule.RuleName
			for _, step := range rule.Steps {
				if step.Op == "take" {
					rootId = step.Item
				}
			}
		}
	}
	osdStats := map[uint64]*Osd{}
	for num, stat := range osdDump.Osds {
		osdStats[stat.Osd] = &osdDump.Osds[num]
	}

	var osdDevices []Device
	bucketItems := getCrushHostBuckets(osdCrushDump.Buckets, rootId)

	if params.rdefine != "" { // match regex if exists
		validBucket, err := regexp.CompilePOSIX(params.rdefine)
		if err != nil {
			log.Fatalf("Can't parse regex %v", params.rdefine)
		}
		for _, hostBucket := range bucketItems {
			for _, item := range hostBucket.Items {
				for _, device := range osdCrushDump.Devices {
					if device.ID == item.ID && (validBucket.MatchString(hostBucket.Name) || validBucket.MatchString(device.Name)) {
						for _, osdMetadata := range osdMetadata {
							if osdMetadata.ID == device.ID && osdStats[uint64(device.ID)].Up == 1 && osdStats[uint64(device.ID)].In == 1 {
								device.Info = osdMetadata
								osdDevices = append(osdDevices, device)
							}
						}
					}
				}
			}
		}
		if len(osdDevices) == 0 {
			log.Fatalf("Defined host/osd doesn't exist in root for rule: %v pool: %v", crushRuleName, poolInfo.Pool)
		}
	} else if params.define != "" { // check defined osd/hosts
		if strings.HasPrefix(params.define, "osd.") { //check that defined is osd, else host
			for _, hostBucket := range bucketItems {
				for _, item := range hostBucket.Items {
					for _, device := range osdCrushDump.Devices {
						if device.ID == item.ID && params.define == device.Name {
							for _, osdMetadata := range osdMetadata {
								if osdMetadata.ID == device.ID && osdStats[uint64(device.ID)].Up == 1 && osdStats[uint64(device.ID)].In == 1 {
									device.Info = osdMetadata
									osdDevices = append(osdDevices, device)
								}
							}
						}
					}
				}
			}
			if len(osdDevices) == 0 {
				log.Fatalf("Defined osd doesn't exist in root for rule: %v pool: %v.\nYou should define osd like osd.X",
					crushRuleName, poolInfo.Pool)
			}
		} else {
			for _, hostBucket := range bucketItems {
				if strings.Split(hostBucket.Name, "~")[0] == strings.Split(params.define, "~")[0] { //purge device class
					for _, item := range hostBucket.Items {
						for _, device := range osdCrushDump.Devices {
							if device.ID == item.ID {
								for _, osdMetadata := range osdMetadata {
									if osdMetadata.ID == device.ID && osdStats[uint64(device.ID)].Up == 1 && osdStats[uint64(device.ID)].In == 1 {
										device.Info = osdMetadata
										osdDevices = append(osdDevices, device)
									}
								}
							}
						}
					}
				}
			}
			if len(osdDevices) == 0 {
				log.Fatalf("Defined host doesn't exist in root for rule: %v pool: %v", crushRuleName, poolInfo.Pool)
			}
		}
	} else {
		for _, hostBucket := range bucketItems {
			for _, item := range hostBucket.Items {
				for _, device := range osdCrushDump.Devices {
					if device.ID == item.ID {
						for _, osdMetadata := range osdMetadata {
							if osdMetadata.ID == device.ID && osdStats[uint64(device.ID)].Up == 1 && osdStats[uint64(device.ID)].In == 1 {
								device.Info = osdMetadata
								osdDevices = append(osdDevices, device)
							}
						}
					}
				}
			}
		}
		if len(osdDevices) == 0 {
			log.Fatalf("Osd doesn't exist in root for rule: %v pool: %v", crushRuleName, poolInfo.Pool)
		}
	}
	return osdDevices
}

func containsPg(pgs []PlacementGroup, i int64) bool {
	for _, pg := range pgs {
		if i == int64(pg.ActingPrimary) {
			return true
		}
	}
	return false
}

func getOsds(cephConn *cephconnection, params params) []Device {
	poolInfo := getPoolSize(cephConn, params)
	if params.disableCheck == false {
		if poolInfo.Size != 1 {
			log.Fatalf("Pool size must be 1. Current size for pool %v is %v. Don't forget that it must be useless pool (not production). Do:\n # ceph osd pool set %v min_size 1\n # ceph osd pool set %v size 1",
				poolInfo.Pool, poolInfo.Size, poolInfo.Pool, poolInfo.Pool)
		}
	}
	placementGroups := getPgByPool(cephConn, params)
	crushOsdDump := getOsdCrushDump(cephConn)
	osdDump := getOsdDump(cephConn)
	osdsMetadata := getOsdMetadata(cephConn)
	osdDevices := getOsdForLocations(params, crushOsdDump, osdDump, poolInfo, osdsMetadata)
	for _, device := range osdDevices {
		if exist := containsPg(placementGroups, device.ID); exist == false {
			log.Fatalln("Not enough pg for test. Some osd don't have placement group at all. Increase pg_num and pgp_num")
		}
	}
	return osdDevices
}
