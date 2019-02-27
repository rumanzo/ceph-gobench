package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func bench(cephconn *Cephconnection, osddevice Device, buffs *[][]byte, bs int64, objsize int64, params *Params,
	wg *sync.WaitGroup, result chan string) {
	defer wg.Done()
	threadresult := make(chan []time.Duration, params.threadsCount)
	var objectnames []string
	var osdlatencies []time.Duration
	defer func() {
		for _, object := range objectnames {
			cephconn.ioctx.Delete(object)
		}
	}()
	// calculate object for each thread
	for suffix := 0; len(objectnames) < int(params.threadsCount)*16; suffix++ {
		name := "bench_" + strconv.Itoa(suffix)
		if osddevice.ID == GetObjActingPrimary(cephconn, *params, name) {
			objectnames = append(objectnames, name)
		}
	}
	for i := 0; i < int(params.threadsCount); i++ {
		go bench_thread(cephconn, osddevice, (*buffs)[i*2:i*2+2], bs, objsize, params, threadresult, objectnames[i*16:i*16+16])
	}
	for i := int64(0); i < params.threadsCount; i++ {
		for _, lat := range <-threadresult {
			osdlatencies = append(osdlatencies, lat)
		}
	}
	close(threadresult)
	latencygrade := map[int]int{}
	for _, lat := range osdlatencies {
		switch {
		case lat < time.Millisecond*10:
			latencygrade[int(lat.Round(time.Millisecond).Nanoseconds()/1000000)]++
		case lat < time.Millisecond*100:
			latencygrade[int(lat.Round(time.Millisecond*10)/1000000)]++
		default:
			latencygrade[int(lat.Round(time.Millisecond*100)/1000000)]++
		}
	}

	var buffer bytes.Buffer

	//color info
	yellow := color.New(color.FgHiYellow).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()
	darkred := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()
	darkgreen := color.New(color.FgGreen).SprintFunc()
	buffer.WriteString(fmt.Sprintf("Bench result for %v\n", osddevice.Name))
	infos := map[string]string{"front_addr": strings.Split(osddevice.Info.FrontAddr, "/")[0],
		"ceph_release/version": osddevice.Info.CephRelease + "/" + osddevice.Info.CephVersionShort, "cpu": osddevice.Info.CPU,
		"hostname": osddevice.Info.Hostname, "default_device_class": osddevice.Info.DefaultDeviceClass, "devices": osddevice.Info.Devices,
		"distro_description": osddevice.Info.DistroDescription, "journal_rotational": osddevice.Info.JournalRotational,
		"rotational": osddevice.Info.Rotational, "kernel_version": osddevice.Info.KernelVersion, "mem_swap_kb": osddevice.Info.MemSwapKb,
		"mem_total_kb": osddevice.Info.MemTotalKb, "osd_data": osddevice.Info.OsdData, "osd_objectstore": osddevice.Info.OsdObjectstore}
	infonum := 1
	var infokeys []string
	for k := range infos {
		infokeys = append(infokeys, k)
	}
	sort.Strings(infokeys)
	buffer.WriteString(fmt.Sprintf("%-30v %-45v", darkgreen("osdname"), red(osddevice.Name)))
	for _, key := range infokeys {
		infonum++
		buffer.WriteString(fmt.Sprintf("%-30v %-45v", darkgreen(key), yellow(infos[key])))
		if (infonum % 3) == 0 {
			buffer.WriteString("\n")
		}
	}

	//sort latencies
	var keys []int
	for k := range latencygrade {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		var blocks bytes.Buffer
		var mseconds string
		switch {
		case 10 <= k && k < 20:
			mseconds = green(fmt.Sprintf("[%v-%v]", k, k+9))
		case 20 <= k && k < 100:
			mseconds = red(fmt.Sprintf("[%v-%v]", k, k+9))
		case k >= 100:
			mseconds = darkred(fmt.Sprintf("[%v-%v]", k, k+99))
		default:
			mseconds = green(k)
		}
		for i := 0; i < 50*(latencygrade[k]*100/len(osdlatencies))/100; i++ {
			blocks.WriteString("#")
		}
		iops := latencygrade[k] / int(params.duration.Seconds())
		avgspeed := (float64(latencygrade[k]) * float64(params.blocksize) / float64(params.duration.Seconds())) / 1024 / 1024 //mb/sec
		megabyteswritten := (float64(latencygrade[k]) * float64(params.blocksize)) / 1024 / 1024
		buffer.WriteString(fmt.Sprintf("%+9v ms: [%-50v]    Count: %-5v    IOPS: %-5v    Avg speed: %-6.3f Mb/Sec    Summary written: %6.3f Mb\n",
			mseconds, blocks.String(), latencygrade[k], iops, avgspeed, megabyteswritten))
	}
	result <- buffer.String()
}

func bench_thread(cephconn *Cephconnection, osddevice Device, buffs [][]byte, bs int64, objsize int64, params *Params,
	result chan []time.Duration, objnames []string) {

	starttime := time.Now()
	var latencies []time.Duration
	endtime := starttime.Add(params.duration)
	n := 0
	for {
		offset := rand.Int63n(objsize/bs) * bs
		objname := objnames[rand.Int31n(int32(len(objnames)))]
		startwritetime := time.Now()
		if startwritetime.After(endtime) {
			break
		}
		err := cephconn.ioctx.Write(objname, buffs[n], uint64(offset))
		endwritetime := time.Now()
		if err != nil {
			log.Printf("Can't write obj: %v, osd: %v", objname, osddevice.Name)
			continue
		}
		latencies = append(latencies, endwritetime.Sub(startwritetime))
		if n == 0 {
			n++
		} else {
			n = 0
		}
	}
	result <- latencies
}

func main() {
	params := Route()
	cephconn := connectioninit(params)
	defer cephconn.conn.Shutdown()

	// https://tracker.ceph.com/issues/24114
	time.Sleep(time.Millisecond * 100)

	var buffs [][]byte
	for i := int64(0); i < 2*params.threadsCount; i++ {
		buffs = append(buffs, make([]byte, params.blocksize))
	}
	for num := range buffs {
		_, err := rand.Read(buffs[num])
		if err != nil {
			log.Fatalln(err)
		}
	}
	osddevices := GetOsds(cephconn, params)

	var wg sync.WaitGroup
	results := make(chan string, len(osddevices)*int(params.threadsCount))
	for _, osd := range osddevices {
		wg.Add(1)
		if params.parallel == true {
			go bench(cephconn, osd, &buffs, params.blocksize, params.objectsize, &params, &wg, results)
		} else {
			bench(cephconn, osd, &buffs, params.blocksize, params.objectsize, &params, &wg, results)
			log.Println(<-results)
		}

	}

	if params.parallel == true {
		go func() {
			wg.Wait()
			close(results)
		}()

		for message := range results {
			log.Println(message)

		}
	}

}
