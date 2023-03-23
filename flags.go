package main

import (
	"code.cloudfoundry.org/bytefmt"
	"fmt"
	"github.com/juju/gnuflag"
	"log"
	"os"
	"time"
)

func route() params {
	params := params{}
	var showversion bool
	const version = 1.3
	gnuflag.DurationVar(&params.duration, "duration", 30*time.Second,
		"Time limit for each test in seconds")
	gnuflag.DurationVar(&params.duration, "d", 30*time.Second,
		"Time limit for each test in seconds")
	gnuflag.StringVar(&params.bs, "blockSize", "4K",
		"Block size in format  KB = K = KiB	= 1024 MB = M = MiB = 1024 * K GB = G = GiB = 1024 * M TB = T = TiB = 1024 * G")
	gnuflag.StringVar(&params.bs, "s", "4K",
		"Block size in format  KB = K = KiB	= 1024 MB = M = MiB = 1024 * K GB = G = GiB = 1024 * M TB = T = TiB = 1024 * G")
	gnuflag.StringVar(&params.os, "objectSize", "4M",
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
		"Define specifically osd or host. Example: osd.X, ceph-host-X")
	gnuflag.StringVar(&params.rdefine, "rdefine", "",
		"Rdefine specifically osd or host in Posix Regex (replaces define). Example: osd.X, ceph-host-X, osd.[0-9]1?$, ceph-host-[1-2]~hdd")
	gnuflag.Uint64Var(&params.threadsCount, "threads", 1,
		"Threads count")
	gnuflag.Uint64Var(&params.threadsCount, "t", 1,
		"Threads count on each osd")
	gnuflag.BoolVar(&params.parallel, "parallel", false,
		"Do test all osd in parallel mode")
	gnuflag.BoolVar(&params.disableCheck, "disablepoolsizecheck", false,
		"Do test all osd in parallel mode")
	gnuflag.StringVar(&params.cpuprofile, "cpuprofile", "",
		"Name of cpuprofile")
	gnuflag.StringVar(&params.memprofile, "memprofile", "",
		"Name of memprofile")
	gnuflag.BoolVar(&showversion, "version", false, "Show sversion")
	gnuflag.BoolVar(&showversion, "v", false, "Show version")

	gnuflag.Parse(true)

	if showversion {
		fmt.Printf("go-bench version v%v\n", version)
		os.Exit(0)
	}

	blocksize, err := bytefmt.ToBytes(params.bs)
	params.blockSize = blocksize
	if err != nil {
		log.Println("Can't convert defined block size. 4K block size will be used")
		params.blockSize = 4096
	}

	objsize, err := bytefmt.ToBytes(params.os)
	params.objectSize = objsize
	if err != nil {
		log.Println("Can't convert defined block size. 4K block size will be used")
		params.objectSize = 4194304
	}
	if params.blockSize > params.objectSize {
		log.Fatalf("Current block size: %v\nCurrent object size: %v\nObject must be bigger or equal than block size", params.blockSize, params.objectSize)
	}
	if params.objectSize/params.blockSize < 2 {
		log.Printf("Current block size: %v\nCurrent object size: %v\nOffset always will be 0", params.blockSize, params.objectSize)
	}
	return params
}
