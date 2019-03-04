package main

import (
	"log"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/juju/gnuflag"
)

func route() params {
	params := params{}
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
	gnuflag.StringVar(&params.user, "user", "client.admin",
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
	gnuflag.Uint64Var(&params.threadsCount, "threads", 1,
		"Threads count on each osd")
	gnuflag.Uint64Var(&params.threadsCount, "t", 1,
		"Threads count on each osd")
	gnuflag.BoolVar(&params.parallel, "parallel", false,
		"Do test all osd in parallel mode")
	gnuflag.StringVar(&params.cpuprofile, "cpuprofile", "",
		"Name of cpuprofile")
	gnuflag.StringVar(&params.memprofile, "memprofile", "",
		"Name of memprofile")
	gnuflag.Parse(true)

	blocksize, err := bytefmt.ToBytes(params.bs)
	params.blocksize = blocksize
	if err != nil {
		log.Println("Can't convert defined block size. 4K block size will be used")
		params.blocksize = 4096
	}

	objsize, err := bytefmt.ToBytes(params.os)
	params.objectsize = objsize
	if err != nil {
		log.Println("Can't convert defined block size. 4K block size will be used")
		params.objectsize = 4194304
	}
	if params.blocksize > params.objectsize {
		log.Fatalf("Current block size: %v\nCurrent object size: %v\nObject must be bigger or equal than block size", params.blocksize, params.objectsize)
	}
	if params.objectsize/params.blocksize < 2 {
		log.Printf("Current block size: %v\nCurrent object size: %v\nOffset always will be 0", params.blocksize, params.objectsize)
	}
	return params
}
