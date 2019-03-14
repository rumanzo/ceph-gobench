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
Usage of ./ceph-gobench-v1.3:
-c, --config (= "/etc/ceph/ceph.conf")
    Ceph config
--cpuprofile (= "")
    Name of cpuprofile
-d, --duration  (= 30s)
    Time limit for each test in seconds
--define (= "")
    Define specifically osd or host. Example: osd.X, ceph-host-X
-k, --keyring (= "/etc/ceph/ceph.client.admin.keyring")
    Ceph user keyring
--memprofile (= "")
    Name of memprofile
-n, --cluster (= "ceph")
    Ceph cluster name
-o, --objectsize (= "4M")
    Object size in format  KB = K = KiB = 1024 MB = M = MiB = 1024 * K GB = G = GiB = 1024 * M TB = T = TiB = 1024 * G
-p, --pool (= "bench")
    Ceph pool
--parallel  (= false)
    Do test all osd in parallel mode
--rdefine (= "")
    Rdefine specifically osd or host in Posix Regex (replaces define). Example: osd.X, ceph-host-X, osd.[0-9]1?$, ceph-host-[1-2]~hdd
-s, --blocksize (= "4K")
    Block size in format  KB = K = KiB  = 1024 MB = M = MiB = 1024 * K GB = G = GiB = 1024 * M TB = T = TiB = 1024 * G
-t, --threads  (= 1)
    Threads count on each osd
-u, --user (= "client.admin")
    Ceph user (cephx)
-v, --version  (= false)
    Show version
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
[root@host1 ~]# ./ceph-gobench-v1.3 -t 10 -d 10s --rdefine "osd.[0-2]$" --parallel
2019/03/11 16:53:58 Calculating objects
2019/03/11 16:54:07 Benchmark started
2019/03/11 16:54:21 Bench result for osd.2
osdname               osd.2                      ceph_release/version  nautilus/14.1.0           cpu                 Intel(R) Xeon(R) CPU E5-2603 v4 @ 1.70GHz  
default_device_class  hdd                        devices               sdc                       distro_description  CentOS Linux 7 (Core)                      
front_addr            [v2:10.0.0.10:6887         hostname              host1                     journal_rotational  1                                          
kernel_version        3.10.0-957.5.1.el7.x86_64  mem_swap_kb           32964604                  mem_total_kb        65693020                                   
osd_data              /var/lib/ceph/osd/ceph-2   osd_objectstore       bluestore                 rotational          1                                          

Avg iops: 330      Avg speed: 1.289 MB/s    Total writes count: 3381     Total writes (MB): 13   

[8.0-9.0)   ms: [                                                  ]    Count: 9        Total written:  0.035 MB
[10-20)     ms: [#####################                             ]    Count: 1460     Total written:  5.703 MB
[20-30)     ms: [#####################                             ]    Count: 1421     Total written:  5.551 MB
[40-50)     ms: [                                                  ]    Count: 1        Total written:  0.004 MB
[50-60)     ms: [                                                  ]    Count: 10       Total written:  0.039 MB
[60-70)     ms: [                                                  ]    Count: 10       Total written:  0.039 MB
[70-80)     ms: [                                                  ]    Count: 20       Total written:  0.078 MB
[80-90)     ms: [#                                                 ]    Count: 71       Total written:  0.277 MB
[90-100)    ms: [##                                                ]    Count: 150      Total written:  0.586 MB
[100-199]   ms: [###                                               ]    Count: 229      Total written:  0.895 MB

2019/03/11 16:54:21 Bench result for osd.1
osdname               osd.1                      ceph_release/version  nautilus/14.1.0           cpu                 Intel(R) Xeon(R) CPU E5-2603 v4 @ 1.70GHz
default_device_class  hdd                        devices               sdb                       distro_description  CentOS Linux 7 (Core)
front_addr            [v2:10.0.0.10:6824         hostname              host1                     journal_rotational  1
kernel_version        3.10.0-957.5.1.el7.x86_64  mem_swap_kb           32964604                  mem_total_kb        65693020
osd_data              /var/lib/ceph/osd/ceph-1   osd_objectstore       bluestore                 rotational          1

Avg iops: 330      Avg speed: 1.289 MB/s    Total writes count: 3397     Total writes (MB): 13

[10-20)     ms: [###################                               ]    Count: 1303     Total written:  5.090 MB
[20-30)     ms: [#####################                             ]    Count: 1480     Total written:  5.781 MB
[40-50)     ms: [                                                  ]    Count: 14       Total written:  0.055 MB
[50-60)     ms: [                                                  ]    Count: 64       Total written:  0.250 MB
[60-70)     ms: [                                                  ]    Count: 66       Total written:  0.258 MB
[70-80)     ms: [##                                                ]    Count: 158      Total written:  0.617 MB
[80-90)     ms: [##                                                ]    Count: 168      Total written:  0.656 MB
[90-100)    ms: [#                                                 ]    Count: 68       Total written:  0.266 MB
[100-199]   ms: [#                                                 ]    Count: 76       Total written:  0.297 MB

2019/03/11 16:54:21 Bench result for osd.0
osdname               osd.0                      ceph_release/version  nautilus/14.1.0           cpu                 Intel(R) Xeon(R) CPU E5-2603 v4 @ 1.70GHz
default_device_class  hdd                        devices               sda                       distro_description  CentOS Linux 7 (Core)
front_addr            [v2:10.0.0.10:6867         hostname              host1                     journal_rotational  1
kernel_version        3.10.0-957.5.1.el7.x86_64  mem_swap_kb           32964604                  mem_total_kb        65693020
osd_data              /var/lib/ceph/osd/ceph-0   osd_objectstore       bluestore                 rotational          1

Avg iops: 320      Avg speed: 1.250 MB/s    Total writes count: 3260     Total writes (MB): 12

[10-20)     ms: [####################                              ]    Count: 1330     Total written:  5.195 MB
[20-30)     ms: [####################                              ]    Count: 1359     Total written:  5.309 MB
[30-40)     ms: [                                                  ]    Count: 19       Total written:  0.074 MB
[40-50)     ms: [                                                  ]    Count: 21       Total written:  0.082 MB
[50-60)     ms: [                                                  ]    Count: 60       Total written:  0.234 MB
[60-70)     ms: [                                                  ]    Count: 1        Total written:  0.004 MB
[70-80)     ms: [                                                  ]    Count: 57       Total written:  0.223 MB
[80-90)     ms: [##                                                ]    Count: 135      Total written:  0.527 MB
[90-100)    ms: [##                                                ]    Count: 138      Total written:  0.539 MB
[100-199]   ms: [##                                                ]    Count: 140      Total written:  0.547 MB

osd.0    Avg iops: 320      Avg speed: 1.250 MB/s    Total writes count: 3260     Total writes (MB): 12
osd.1    Avg iops: 330      Avg speed: 1.289 MB/s    Total writes count: 3397     Total writes (MB): 13
osd.2    Avg iops: 330      Avg speed: 1.289 MB/s    Total writes count: 3381     Total writes (MB): 13

Average iops per osd:      330    Average speed per osd: 1.289 MB/s
Total writes count:      10038    Total writes (MB): 39
Summary avg iops:          990    Summary avg speed: 3.867 MB/s
```
