# ceph-gobench
ceph-gobench is benchmark for ceph which allows you to measure the speed/iops of each osd.
- [ceph-gobench](#ceph-gobench)
	- [Feature](#user-content-feature)
	- [Quickstart](#user-content-quickstart)
	- [How it works](#user-content-how-it-works)
	- [Example](#user-content-example)

Feature
--------------
 - Static build
 - Multithreading
 - Parallel mode
 - Custom parametersh such as block size, object size, test duration
 - Credentials support
 
Quickstart
--------------
1. Build or download latest release
2. Create pool with size=1 min_size=1 (default pool name "bench")
3. Use
```
Usage of ./ceph-gobench:
-c, --config (= "/etc/ceph/ceph.conf")
    Ceph config
-d, --duration  (= 30s)
    Time limit for each test in seconds
--define (= "")
    Define specifically osd or host. osd.X or ceph-host-X
-k, --keyring (= "/etc/ceph/ceph.client.admin.keyring")
    Ceph user keyring
-n, --cluster (= "ceph")
    Ceph cluster name
-o, --objectsize (= "4M")
    Object size in format  KB = K = KiB = 1024 MB = M = MiB = 1024 * K GB = G = GiB = 1024 * M TB = T = TiB = 1024 * G
-p, --pool (= "bench")
    Ceph pool
--parallel  (= false)
    Do test all osd in parallel mode
-s, --blocksize (= "4K")
    Block size in format  KB = K = KiB  = 1024 MB = M = MiB = 1024 * K GB = G = GiB = 1024 * M TB = T = TiB = 1024 * G
-t, --threads  (= 1)
    Threads count on each osd
-u, --user (= "client.admin")
    Ceph user (cephx)
```    

How it works
--------------
In order to measure the performance of a specific OSD, you must create a pool with a size/min_size 1.
Benchmark creates 16 objects of size defined with parameter objectsize, create threads that's will work parallel and write random buffer to aligned to
block size offset with measuring write time. According to the received data builds a report.

Tests run by one if not defined parameter --parallel.
You can also define specifix osd or host for test.


Example
--------------

```
[root@host-1]# ./ceph-gobench -d 10s -p testpool --define osd.61
2019/02/28 18:26:19 Bench result for osd.61
osdname               osd.61                     ceph_release/version  mimic/13.2.2        cpu                 Intel(R) Xeon(R) CPU E5-2680 v4 @ 2.40GHz
default_device_class  ssd                        devices               dm-19,dm-24,sdv     distro_description  CentOS Linux 7 (Core)
front_addr            10.0.0.1:6831              hostname              host-1              journal_rotational  0
kernel_version        3.10.0-957.1.3.el7.x86_64  mem_swap_kb           0                   mem_total_kb        263921472
osd_data              /var/lib/ceph/osd/ceph-61  osd_objectstore       bluestore           rotational          0

Avg iops: 1410     Avg speed: 5.510 MB/s    Total writes count: 14066    Total writes (MB): 54

[0.4-0.5)   ms: [                                                  ]    Count: 63       Total written:  0.246 MB
[0.5-0.6)   ms: [##########                                        ]    Count: 2830     Total written: 11.055 MB
[0.6-0.7)   ms: [#############                                     ]    Count: 3857     Total written: 15.066 MB
[0.7-0.8)   ms: [################                                  ]    Count: 4578     Total written: 17.883 MB
[0.8-0.9)   ms: [#######                                           ]    Count: 2082     Total written:  8.133 MB
[0.9-1.0)   ms: [#                                                 ]    Count: 391      Total written:  1.527 MB
[1.0-1.1)   ms: [                                                  ]    Count: 255      Total written:  0.996 MB
[2.0-3.0)   ms: [                                                  ]    Count: 8        Total written:  0.031 MB
[3.0-4.0)   ms: [                                                  ]    Count: 1        Total written:  0.004 MB
[4.0-5.0)   ms: [                                                  ]    Count: 1        Total written:  0.004 MB
```
