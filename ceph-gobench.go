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

func bench(cephconn *Cephconnection, osddevice Device, buffs *[][]byte, startbuff *[]byte, params *Params,
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
			if err := cephconn.ioctx.WriteFull(name, *startbuff); err != nil {
				log.Printf("Can't write object: %v, osd: %v", name, osddevice.Name)
			}
			if err := cephconn.ioctx.Truncate(name, uint64(params.objectsize)); err != nil {
				log.Printf("Can't truncate object: %v, osd: %v", name, osddevice.Name)
			}
		}
	}
	for i := 0; i < int(params.threadsCount); i++ {
		go BenchThread(cephconn, osddevice, (*buffs)[i*2:i*2+2], params, threadresult, objectnames[i*16:i*16+16])
	}
	for i := uint64(0); i < params.threadsCount; i++ {
		for _, lat := range <-threadresult {
			osdlatencies = append(osdlatencies, lat)
		}
	}
	close(threadresult)
	latencygrade := map[int64]int{}
	latencytotal := int64(0)
	for _, lat := range osdlatencies {
		micro := lat.Nanoseconds() / 1000
		rounded := micro
		switch {
		case micro < 1000: // 0-1ms round to 0.1ms
			rounded = (micro / 100) * 100
		case micro < 10000: // 2-10ms round to 1ms
			rounded = (micro / 1000) * 1000
		case micro < 100000: // 10-100ms round to 10ms
			rounded = (micro / 10000) * 10000
		case micro < 1000000: // 100-1000ms round to 100ms
			rounded = (micro / 100000) * 100000
		default: // 1000+ms round to 1s
			rounded = (micro / 1000000) * 1000000
		}
		latencytotal += micro
		latencygrade[rounded]++
	}

	var buffer bytes.Buffer

	//color info
	yellow := color.New(color.FgHiYellow).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()
	darkred := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()
	buffer.WriteString(fmt.Sprintf("Bench result for %v\n", osddevice.Name))
	infos := map[string]string{
		"front_addr":           strings.Split(osddevice.Info.FrontAddr, "/")[0],
		"ceph_release/version": osddevice.Info.CephRelease + "/" + osddevice.Info.CephVersionShort,
		"cpu":                  osddevice.Info.CPU,
		"hostname":             osddevice.Info.Hostname,
		"default_device_class": osddevice.Info.DefaultDeviceClass,
		"devices":              osddevice.Info.Devices,
		"distro_description":   osddevice.Info.DistroDescription,
		"journal_rotational":   osddevice.Info.JournalRotational,
		"rotational":           osddevice.Info.Rotational,
		"kernel_version":       osddevice.Info.KernelVersion,
		"mem_swap_kb":          osddevice.Info.MemSwapKb,
		"mem_total_kb":         osddevice.Info.MemTotalKb,
		"osd_data":             osddevice.Info.OsdData,
		"osd_objectstore":      osddevice.Info.OsdObjectstore,
	}
	var infokeys []string
	width := []int{0, 0, 0, 0, 0, 0}
	for k := range infos {
		infokeys = append(infokeys, k)
	}
	sort.Strings(infokeys)
	for n, key := range infokeys {
		if width[n%3] < len(key) {
			width[n%3] = len(key)
		}
		if width[3+n%3] < len(infos[key]) {
			width[3+n%3] = len(infos[key])
		}
	}
	buffer.WriteString(
		green("osdname") + strings.Repeat(" ", width[2]-len("osdname")+2) +
			red(osddevice.Name) + strings.Repeat(" ", width[5]-len(osddevice.Name)+2))
	for infonum, key := range infokeys {
		if (infonum % 3) == 2 {
			buffer.WriteString("\n")
		}
		buffer.WriteString(
			green(key) + strings.Repeat(" ", width[infonum%3]-len(key)+2) +
				yellow(infos[key]) + strings.Repeat(" ", width[3+infonum%3]-len(infos[key])+2))
	}
	buffer.WriteString("\n\n")

	latencytotal = latencytotal / int64(len(osdlatencies))
	// iops = 1s / latency
	iops := 1000000 / latencytotal * int64(params.threadsCount)
	// avg speed = iops * block size / 1 MB
	avgspeed := 1000000 / float64(latencytotal) * float64(params.blocksize) / 1024 / 1024
	avgline := fmt.Sprintf("Avg iops: %-5v    Avg speed: %.3f MB/s    Total writes count: %-5v    Total writes (MB): %-5v\n\n",
		iops, avgspeed, len(osdlatencies), uint64(len(osdlatencies))*params.blocksize/1024/1024)
	switch {
	case iops < 80:
		buffer.WriteString(darkred(avgline))
	case iops < 200:
		buffer.WriteString(red(avgline))
	case iops < 500:
		buffer.WriteString(yellow(avgline))
	default:
		buffer.WriteString(green(avgline))
	}

	//sort latencies
	var keys []int64
	for k := range latencygrade {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, k := range keys {
		var blocks bytes.Buffer
		var mseconds string
		switch {
		case k < 1000:
			mseconds = green(fmt.Sprintf("[%.1f-%.1f)", float64(k)/1000, 0.1+float64(k)/1000))
		case k < 2000:
			mseconds = yellow(fmt.Sprintf("[%.1f-%.1f)", float64(k)/1000, 0.1+float64(k)/1000))
		case k < 9000:
			mseconds = yellow(fmt.Sprintf("[%.1f-%.1f)", float64(k/1000), float64(1+k/1000)))
		case k < 10000:
			mseconds = color.YellowString(fmt.Sprintf("[%.1f-%v)", float64(k/1000), 1+k/1000))
		case k < 100000:
			mseconds = red(fmt.Sprintf("[%v-%v)", k/1000, 10+k/1000))
		case k < 1000000:
			mseconds = darkred(fmt.Sprintf("[%v-%v]", k/1000, 99+k/1000))
		default:
			mseconds = darkred(fmt.Sprintf("[%vs-%vs]", k/1000000, 1+k/1000000))
		}
		for i := 0; i < 50*(latencygrade[k]*100/len(osdlatencies))/100; i++ {
			blocks.WriteString("#")
		}
		megabyteswritten := (float64(latencygrade[k]) * float64(params.blocksize)) / 1024 / 1024
		buffer.WriteString(fmt.Sprintf("%-20v ms: [%-50v]    Count: %-5v    Total written: %6.3f MB\n",
			mseconds, blocks.String(), latencygrade[k], megabyteswritten))
	}
	result <- buffer.String()
}

func BenchThread(cephconn *Cephconnection, osddevice Device, buffs [][]byte, params *Params,
	result chan []time.Duration, objnames []string) {

	starttime := time.Now()
	var latencies []time.Duration
	endtime := starttime.Add(params.duration)
	n := 0
	for {
		offset := rand.Int63n(int64(params.objectsize/params.blocksize)) * int64(params.blocksize)
		objname := objnames[rand.Int31n(int32(len(objnames)))]
		startwritetime := time.Now()
		if startwritetime.After(endtime) {
			break
		}
		err := cephconn.ioctx.Write(objname, buffs[n], uint64(offset))
		endwritetime := time.Now()
		if err != nil {
			log.Printf("Can't write object: %v, osd: %v", objname, osddevice.Name)
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
	for i := uint64(0); i < 2*params.threadsCount; i++ {
		buffs = append(buffs, make([]byte, params.blocksize))
	}
	startbuff := make([]byte, 4096)
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
			go bench(cephconn, osd, &buffs, &startbuff, &params, &wg, results)
		} else {
			bench(cephconn, osd, &buffs, &startbuff, &params, &wg, results)
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
