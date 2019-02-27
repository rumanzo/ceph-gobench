package main

import (
	"code.cloudfoundry.org/bytefmt"
	"github.com/juju/gnuflag"
	"log"
	"strings"
	"time"
)

func Route() Params {
	params := Params{}
	gnuflag.DurationVar(&params.duration, "duration", 30*time.Second,
		"Time limit for each test in seconds")
	gnuflag.DurationVar(&params.duration, "d", 30*time.Second,
		"Time limit for each test in seconds")
	gnuflag.StringVar(&params.bs, "blocksize", "4K",
		"Block size in format  KB = K = KiB	= 1024 MB = M = MiB = 1024 * K GB = G = GiB = 1024 * M TB = T = TiB = 1024 * G")
	gnuflag.StringVar(&params.bs, "s", "4K",
		"Block size in format  KB = K = KiB	= 1024 MB = M = MiB = 1024 * K GB = G = GiB = 1024 * M TB = T = TiB = 1024 * G")
	gnuflag.StringVar(&params.os, "objectsize", "4M",
		"Object size in format  KB = K = KiB	= 1024 MB = M = MiB = 1024 * K GB = G = GiB = 1024 * M TB = T = TiB = 1024 * G")
	gnuflag.StringVar(&params.os, "o", "4M",
		"Object size in format  KB = K = KiB	= 1024 MB = M = MiB = 1024 * K GB = G = GiB = 1024 * M TB = T = TiB = 1024 * G")
	gnuflag.StringVar(&params.user, "user", "admin",
		"Ceph user (cephx)")
	gnuflag.StringVar(&params.user, "u", "client.admin",
		"Ceph user (cephx)")
	gnuflag.StringVar(&params.cluster, "cluster", "ceph",
		"Ceph cluster name")
	gnuflag.StringVar(&params.cluster, "n", "ceph",
		"Ceph cluster name")
	gnuflag.StringVar(&params.keyring, "keyring", "/etc/ceph/ceph.client.admin.keyring",
		"Ceph user keyring")
	gnuflag.StringVar(&params.keyring, "k", "/etc/ceph/ceph.client.admin.keyring",
		"Ceph user keyring")
	gnuflag.StringVar(&params.config, "config", "/etc/ceph/ceph.conf",
		"Ceph config")
	gnuflag.StringVar(&params.config, "c", "/etc/ceph/ceph.conf",
		"Ceph config")
	gnuflag.StringVar(&params.pool, "pool", "bench",
		"Ceph pool")
	gnuflag.StringVar(&params.pool, "p", "bench",
		"Ceph pool")
	gnuflag.StringVar(&params.define, "define", "",
		"Define specifically osd or host. osd.X or ceph-host-X")
	gnuflag.Int64Var(&params.threadsCount, "threads", 1,
		"Threads count")
	gnuflag.Int64Var(&params.threadsCount, "t", 1,
		"Threads count on each osd")
	gnuflag.BoolVar(&params.parallel, "parallel", false,
		"Do test all osd in parallel mode")
	gnuflag.Parse(true)

	if params.mode == "osd" && len(params.define) != 0 {
		if i := strings.HasPrefix(params.define, "osd."); i != true {
			log.Fatalln("Define correct osd in format osd.X")
		}
	}
	blocksize, err := bytefmt.ToBytes(params.bs)
	params.blocksize = int64(blocksize)
	if err != nil {
		log.Println("Can't convert defined block size. 4K block size will be used")
		params.blocksize = 4096
	}

	objsize, err := bytefmt.ToBytes(params.os)
	params.objectsize = int64(objsize)
	if err != nil {
		log.Println("Can't convert defined block size. 4K block size will be used")
		params.objectsize = 4194304
	}
	if params.objectsize/params.blocksize < 2 {
		log.Fatalf("Current block size: %v\nCurrent object size: %v\nObject size must be at least 2 times bigger than block size", params.blocksize, params.objectsize)
	}
	return params
}
