package main

import (
	"github.com/ceph/go-ceph/rados"
	"log"
	"os"
)

func connectioninit(params params) *cephconnection {
	cephConn := &cephconnection{}
	var err error
	if _, err := os.Stat(params.config); os.IsNotExist(err) {
		log.Fatalf("Congif file not exists. Error: %v\n", err)
	}

	if _, err := os.Stat(params.keyring); os.IsNotExist(err) {
		log.Fatalf("Keyring file not exists. Error: %v\n", err)
	}

	cephConn.conn, err = rados.NewConnWithClusterAndUser(params.cluster, params.user)
	if err != nil {
		log.Fatalf("Can't create connection with cluster:%v and user:%v. Error: %v\n", params.cluster, params.user, err)
	}

	if err = cephConn.conn.ReadConfigFile(params.config); err != nil {
		log.Fatalf("Can't read config file. Error: %v\n", err)
	}

	if err = cephConn.conn.SetConfigOption("keyring", params.keyring); err != nil {
		log.Fatalf("Can't set config option. Error: %v\n", err)
	}

	if err = cephConn.conn.Connect(); err != nil {
		log.Fatalf("Failed to connect cluster. Error: %v\n", err)
	}

	cephConn.ioctx, err = cephConn.conn.OpenIOContext(params.pool)
	if err != nil {
		log.Fatalf("Can't open pool %v. Error: %v\n", params.pool, err)
	}
	return cephConn
}
