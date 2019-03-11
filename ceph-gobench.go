package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func bench(cephconn *cephconnection, osddevice Device, buff *[]byte, startbuff *[]byte, params *params,
	wg *sync.WaitGroup, result chan string, totalLats chan avgLatencies, osdStatsChan chan osdStatLine, objectnames []string) {
	defer wg.Done()
	threadresult := make(chan []time.Duration, params.threadsCount)
	var osdlatencies []time.Duration
	defer func() {
		for _, object := range objectnames {
			cephconn.ioctx.Delete(object)
		}
	}()
	// Create and truncate each object
	for _, object := range objectnames {
		if err := cephconn.ioctx.WriteFull(object, *startbuff); err != nil {
			log.Printf("Can't write object: %v, osd: %v", object, osddevice.Name)
		}
		if err := cephconn.ioctx.Truncate(object, uint64(params.objectsize)); err != nil {
			log.Printf("Can't truncate object: %v, osd: %v", object, osddevice.Name)
		}
	}

	for i := 0; i < int(params.threadsCount); i++ {
		go benchthread(cephconn, osddevice, params, buff, threadresult, objectnames[i*16:i*16+16])
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

	totalLats <- avgLatencies{latencytotal: latencytotal, len: int64(len(osdlatencies))}

	latencytotal = latencytotal / int64(len(osdlatencies))
	// iops = 1s / latency
	iops := 1000000 / latencytotal * int64(params.threadsCount)
	// avg speed = iops * block size / 1 MB
	avgspeed := float64(iops) * float64(params.blocksize) / 1024 / 1024
	avgline := fmt.Sprintf("Avg iops: %-5v    Avg speed: %.3f MB/s    Total writes count: %-5v    Total writes (MB): %-5v\n\n",
		iops, avgspeed, len(osdlatencies), uint64(len(osdlatencies))*params.blocksize/1024/1024)
	osdavgline := fmt.Sprintf("%-8v Avg iops: %-5v    Avg speed: %.3f MB/s    Total writes count: %-5v    Total writes (MB): %-5v",
		osddevice.Name, iops, avgspeed, len(osdlatencies), uint64(len(osdlatencies))*params.blocksize/1024/1024)

	switch {
	case iops < 80:
		buffer.WriteString(darkred(avgline))
		osdStatsChan <- osdStatLine{osddevice.ID, darkred(osdavgline)}
	case iops < 200:
		buffer.WriteString(red(avgline))
		osdStatsChan <- osdStatLine{osddevice.ID, red(osdavgline)}
	case iops < 500:
		buffer.WriteString(yellow(avgline))
		osdStatsChan <- osdStatLine{osddevice.ID, yellow(osdavgline)}
	default:
		buffer.WriteString(green(avgline))
		osdStatsChan <- osdStatLine{osddevice.ID, green(osdavgline)}
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

func benchthread(cephconn *cephconnection, osddevice Device, params *params, buff *[]byte,
	result chan []time.Duration, objnames []string) {
	var latencies []time.Duration
	starttime := time.Now()
	endtime := starttime.Add(params.duration)
	for {
		offset := rand.Int63n(int64(params.objectsize/params.blocksize)) * int64(params.blocksize)
		objname := objnames[rand.Int31n(int32(len(objnames)))]
		startwritetime := time.Now()
		if startwritetime.After(endtime) {
			break
		}
		err := cephconn.ioctx.Write(objname, *buff, uint64(offset))
		endwritetime := time.Now()
		if err != nil {
			log.Printf("Can't write object: %v, osd: %v", objname, osddevice.Name)
			continue
		}
		latencies = append(latencies, endwritetime.Sub(startwritetime))
	}
	result <- latencies
}

func main() {
	params := route()
	if params.cpuprofile != "" {
		f, err := os.Create(params.cpuprofile)
		if err != nil {
			log.Fatal("Could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if params.memprofile != "" {
		f, err := os.Create(params.memprofile)
		if err != nil {
			log.Fatal("Could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("Could not write memory profile: ", err)
		}
	}
	cephconn := connectioninit(params)
	defer cephconn.conn.Shutdown()

	// https://tracker.ceph.com/issues/24114
	time.Sleep(time.Millisecond * 100)

	startbuff := make([]byte, 4096)
	osddevices := getOsds(cephconn, params)
	buff := make([]byte, params.blocksize)
	rand.Read(buff)

	var wg sync.WaitGroup
	results := make(chan string, len(osddevices)*int(params.threadsCount))
	totalLats := make(chan avgLatencies, len(osddevices))
	avgLats := []avgLatencies{}
	osdStatsChan := make(chan osdStatLine, len(osddevices))
	osdStats := map[int64]string{}

	log.Println("Calculating objects")
	objectnames := map[int64][]string{}
	// calculate object for each thread
	for suffix := 0; ; suffix++ {
		name := "bench_" + strconv.Itoa(suffix)
		osdid := getObjActingPrimary(cephconn, params, name)

		objectsdone := 0
		for _, osddevice := range osddevices {
			if osddevice.ID == osdid {
				if len(objectnames[osdid]) < int(params.threadsCount)*16 {
					objectnames[osdid] = append(objectnames[osdid], name)
				} else {

				}
			}
			if len(objectnames[osddevice.ID]) >= int(params.threadsCount)*16 {
				objectsdone++
			}
		}
		if objectsdone >= len(osddevices) {
			break
		}
	}

	log.Println("Benchmark started")

	for _, osd := range osddevices {
		wg.Add(1)
		if params.parallel == true {
			go bench(cephconn, osd, &buff, &startbuff, &params, &wg, results, totalLats, osdStatsChan, objectnames[osd.ID])
		} else {
			bench(cephconn, osd, &buff, &startbuff, &params, &wg, results, totalLats, osdStatsChan, objectnames[osd.ID])
			avgLats = append(avgLats, <-totalLats)
			osdStat := <-osdStatsChan
			osdStats[osdStat.num] = osdStat.line
			log.Println(<-results)
		}

	}

	if params.parallel == true {
		go func() {
			wg.Wait()
			close(results)
			close(totalLats)
			close(osdStatsChan)
		}()

		for message := range results {
			log.Println(message)
		}
		for lat := range totalLats {
			avgLats = append(avgLats, lat)
		}
		for osdStat := range osdStatsChan {
			osdStats[osdStat.num] = osdStat.line
		}
	}

	//print sorted stats for all osd
	var keys []int64
	for k := range osdStats {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, k := range keys {
		fmt.Println(osdStats[k])
	}
	fmt.Println()

	sumLat := int64(0)
	countLat := int64(0)
	for _, avgLat := range avgLats {
		sumLat += avgLat.latencytotal
		countLat += avgLat.len
	}

	//count avg statistics
	sumLat = sumLat / int64(countLat)
	avgIops := 1000000 / sumLat * int64(params.threadsCount)
	sumIops := 1000000 / sumLat * int64(params.threadsCount) * int64(len(osddevices))
	avgSpeed := float64(avgIops) * float64(params.blocksize) / 1024 / 1024
	sumSpeed := float64(sumIops) * float64(params.blocksize) / 1024 / 1024

	color.Set(color.FgHiYellow)
	defer color.Unset()

	fmt.Printf("Average iops per osd:%5d    Average speed per osd: %.3f MB/s\n"+
		"Total writes count:%11d    Total writes (MB): %v\n",
		avgIops, avgSpeed, countLat, uint64(countLat)*params.blocksize/1024/1024)
	if params.parallel {
		fmt.Printf("Summary avg iops:%13d    Summary avg speed: %.3f MB/s\n", sumIops, sumSpeed)
	}
}
