package main

import (
	"code.cloudfoundry.org/bytefmt"
	"github.com/ceph/go-ceph/rados"
	"github.com/juju/gnuflag"
	"log"
	"os"
)

func main() {
	var duration int
	var bs, cluster, user, keyring, config, pool string
	gnuflag.IntVar(&duration, "duration", 30,
		"Time limit for each test in seconds")
	gnuflag.IntVar(&duration, "d", 30,
		"Time limit for each test in seconds")
	gnuflag.StringVar(&bs, "blocksize", "4K",
		"Block size in format  KB = K = KiB	= 1024 MB = M = MiB = 1024 * K GB = G = GiB = 1024 * M TB = T = TiB = 1024 * G")
	gnuflag.StringVar(&bs, "s", "4K",
		"Block size in format  KB = K = KiB	= 1024 MB = M = MiB = 1024 * K GB = G = GiB = 1024 * M TB = T = TiB = 1024 * G")
	gnuflag.StringVar(&user, "user", "admin",
		"Ceph user (cephx)")
	gnuflag.StringVar(&user, "u", "client.admin",
		"Ceph user (cephx)")
	gnuflag.StringVar(&cluster, "cluster", "ceph",
		"Ceph cluster name")
	gnuflag.StringVar(&cluster, "n", "ceph",
		"Ceph cluster name")
	gnuflag.StringVar(&keyring, "keyring", "/etc/ceph/ceph.client.admin.keyring",
		"Ceph user keyring")
	gnuflag.StringVar(&keyring, "k", "/etc/ceph/ceph.client.admin.keyring",
		"Ceph user keyring")
	gnuflag.StringVar(&config, "config", "/etc/ceph/ceph.conf",
		"Ceph config")
	gnuflag.StringVar(&config, "c", "/etc/ceph/ceph.conf",
		"Ceph config")
	gnuflag.StringVar(&pool, "pool", "bench",
		"Ceph pool")
	gnuflag.StringVar(&pool, "p", "bench",
		"Ceph pool")
	gnuflag.Parse(true)

	if _, err := os.Stat(config); os.IsNotExist(err) {
		log.Fatalf("Congif file not exists. Error: %v\n", err)
	}

	if _, err := os.Stat(keyring); os.IsNotExist(err) {
		log.Fatalf("Keyring file not exists. Error: %v\n", err)
	}

	conn, err := rados.NewConnWithClusterAndUser(cluster, user)
	if err != nil {
		log.Fatalf("Can't create connection with cluster:%v and user:%v. Error: %v\n", cluster, user, err)
	}

	if err := conn.ReadConfigFile(config); err != nil {
		log.Fatalf("Can't read config file. Error: err\n", err)
	}

	if err := conn.SetConfigOption("keyring", keyring); err != nil {
		log.Fatalf("Can't set config option. Error: err\n", err)
	}

	if err := conn.Connect(); err != nil {
		log.Fatalf("Failed to connect cluster. Error: %v\n", err)
	}

	ioctx, err := conn.OpenIOContext(pool)
	if err != nil {
		log.Fatalf("Can't open pool %v. Error: %v\n", pool, err)
	}
	stats, _ := ioctx.GetPoolStats()
	conn.Shutdown()
	log.Println(stats)
	blocksize, err := bytefmt.ToBytes(bs)
	if err != nil {
		log.Println("Can't convert defined block size. 4K block size will be used\n")
		blocksize = 4096
	}
	log.Printf("%v\n", blocksize)
}
