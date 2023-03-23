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

func bench(cephConn *cephconnection, osdDevice Device, buff *[]byte, startbuff *[]byte, params *params,
	wg *sync.WaitGroup, result chan string, totalLats chan avgLatencies, osdStatsChan chan osdStatLine, objectNames []string) {
	defer wg.Done()
	threadResult := make(chan []time.Duration, params.threadsCount)
	var osdLatencies []time.Duration
	defer func() {
		for _, object := range objectNames {
			cephConn.ioctx.Delete(object)
		}
	}()
	// Create and truncate each object
	for _, object := range objectNames {
		if err := cephConn.ioctx.WriteFull(object, *startbuff); err != nil {
			log.Printf("Can't write object: %v, osd: %v", object, osdDevice.Name)
		}
		if err := cephConn.ioctx.Truncate(object, uint64(params.objectSize)); err != nil {
			log.Printf("Can't truncate object: %v, osd: %v", object, osdDevice.Name)
		}
	}

	for i := 0; i < int(params.threadsCount); i++ {
		go benchthread(cephConn, osdDevice, params, buff, threadResult, objectNames[i*16:i*16+16])
	}
	for i := uint64(0); i < params.threadsCount; i++ {
		for _, lat := range <-threadResult {
			osdLatencies = append(osdLatencies, lat)
		}
	}
	close(threadResult)
	latencyGrade := map[int64]int{}
	latencyTotal := int64(0)
	for _, lat := range osdLatencies {
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
		latencyTotal += micro
		latencyGrade[rounded]++
	}

	var buffer bytes.Buffer

	//color info
	yellow := color.New(color.FgHiYellow).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()
	darkred := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()
	buffer.WriteString(fmt.Sprintf("Bench result for %v\n", osdDevice.Name))
	infos := map[string]string{
		"front_addr":           strings.Split(osdDevice.Info.FrontAddr, "/")[0],
		"ceph_release/version": osdDevice.Info.CephRelease + "/" + osdDevice.Info.CephVersionShort,
		"cpu":                  osdDevice.Info.CPU,
		"hostname":             osdDevice.Info.Hostname,
		"default_device_class": osdDevice.Info.DefaultDeviceClass,
		"devices":              osdDevice.Info.Devices,
		"distro_description":   osdDevice.Info.DistroDescription,
		"journal_rotational":   osdDevice.Info.JournalRotational,
		"rotational":           osdDevice.Info.Rotational,
		"kernel_version":       osdDevice.Info.KernelVersion,
		"mem_swap_kb":          osdDevice.Info.MemSwapKb,
		"mem_total_kb":         osdDevice.Info.MemTotalKb,
		"osd_data":             osdDevice.Info.OsdData,
		"osd_objectstore":      osdDevice.Info.OsdObjectstore,
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
			red(osdDevice.Name) + strings.Repeat(" ", width[5]-len(osdDevice.Name)+2))
	for infoNum, key := range infokeys {
		if (infoNum % 3) == 2 {
			buffer.WriteString("\n")
		}
		buffer.WriteString(
			green(key) + strings.Repeat(" ", width[infoNum%3]-len(key)+2) +
				yellow(infos[key]) + strings.Repeat(" ", width[3+infoNum%3]-len(infos[key])+2))
	}
	buffer.WriteString("\n\n")

	totalLats <- avgLatencies{latencyTotal: latencyTotal, len: int64(len(osdLatencies))}

	latencyTotal = latencyTotal / int64(len(osdLatencies))
	// iops = 1s / latency
	iops := 1000000 / latencyTotal * int64(params.threadsCount)
	// avg speed = iops * block size / 1 MB
	avgSpeed := float64(iops) * float64(params.blockSize) / 1024 / 1024
	avgLine := fmt.Sprintf("Avg iops: %-5v    Avg speed: %.3f MB/s    Total writes count: %-5v    Total writes (MB): %-5v\n\n",
		iops, avgSpeed, len(osdLatencies), uint64(len(osdLatencies))*params.blockSize/1024/1024)
	osdAvgLine := fmt.Sprintf("%-8v Avg iops: %-5v    Avg speed: %.3f MB/s    Total writes count: %-5v    Total writes (MB): %-5v",
		osdDevice.Name, iops, avgSpeed, len(osdLatencies), uint64(len(osdLatencies))*params.blockSize/1024/1024)

	switch {
	case iops < 80:
		buffer.WriteString(darkred(avgLine))
		osdStatsChan <- osdStatLine{osdDevice.ID, darkred(osdAvgLine)}
	case iops < 200:
		buffer.WriteString(red(avgLine))
		osdStatsChan <- osdStatLine{osdDevice.ID, red(osdAvgLine)}
	case iops < 500:
		buffer.WriteString(yellow(avgLine))
		osdStatsChan <- osdStatLine{osdDevice.ID, yellow(osdAvgLine)}
	default:
		buffer.WriteString(green(avgLine))
		osdStatsChan <- osdStatLine{osdDevice.ID, green(osdAvgLine)}
	}

	//sort latencies
	var keys []int64
	for k := range latencyGrade {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, k := range keys {
		var blocks bytes.Buffer
		var mSeconds string
		switch {
		case k < 1000:
			mSeconds = green(fmt.Sprintf("[%.1f-%.1f)", float64(k)/1000, 0.1+float64(k)/1000))
		case k < 2000:
			mSeconds = yellow(fmt.Sprintf("[%.1f-%.1f)", float64(k)/1000, 0.1+float64(k)/1000))
		case k < 9000:
			mSeconds = yellow(fmt.Sprintf("[%.1f-%.1f)", float64(k/1000), float64(1+k/1000)))
		case k < 10000:
			mSeconds = color.YellowString(fmt.Sprintf("[%.1f-%v)", float64(k/1000), 1+k/1000))
		case k < 100000:
			mSeconds = red(fmt.Sprintf("[%v-%v)", k/1000, 10+k/1000))
		case k < 1000000:
			mSeconds = darkred(fmt.Sprintf("[%v-%v]", k/1000, 99+k/1000))
		default:
			mSeconds = darkred(fmt.Sprintf("[%vs-%vs]", k/1000000, 1+k/1000000))
		}
		for i := 0; i < 50*(latencyGrade[k]*100/len(osdLatencies))/100; i++ {
			blocks.WriteString("#")
		}
		megabytesWritten := (float64(latencyGrade[k]) * float64(params.blockSize)) / 1024 / 1024
		buffer.WriteString(fmt.Sprintf("%-20v ms: [%-50v]    Count: %-5v    Total written: %6.3f MB\n",
			mSeconds, blocks.String(), latencyGrade[k], megabytesWritten))
	}
	result <- buffer.String()
}

func benchthread(cephConn *cephconnection, osdDevice Device, params *params, buff *[]byte,
	result chan []time.Duration, objNames []string) {
	var latencies []time.Duration
	startTime := time.Now()
	endTime := startTime.Add(params.duration)
	for {
		offset := rand.Int63n(int64(params.objectSize/params.blockSize)) * int64(params.blockSize)
		objName := objNames[rand.Int31n(int32(len(objNames)))]
		startWriteTime := time.Now()
		if startWriteTime.After(endTime) {
			break
		}
		err := cephConn.ioctx.Write(objName, *buff, uint64(offset))
		endWriteTime := time.Now()
		if err != nil {
			log.Printf("Can't write object: %v, osd: %v", objName, osdDevice.Name)
			continue
		}
		latencies = append(latencies, endWriteTime.Sub(startWriteTime))
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
	cephConn := connectioninit(params)
	defer cephConn.conn.Shutdown()

	// https://tracker.ceph.com/issues/24114
	time.Sleep(time.Millisecond * 100)

	startBuff := make([]byte, 4096)
	osdDevices := getOsds(cephConn, params)
	buff := make([]byte, params.blockSize)
	rand.Read(buff)

	var wg sync.WaitGroup
	results := make(chan string, len(osdDevices)*int(params.threadsCount))
	totalLats := make(chan avgLatencies, len(osdDevices))
	avgLats := []avgLatencies{}
	osdStatsChan := make(chan osdStatLine, len(osdDevices))
	osdStats := map[int64]string{}

	log.Println("Calculating objects")
	objectNames := map[int64][]string{}
	// calculate object for each thread
	for suffix := 0; ; suffix++ {
		name := "bench_" + strconv.Itoa(suffix)
		osdId := getObjActingPrimary(cephConn, params, name)

		objectsDone := 0
		for _, osdDevice := range osdDevices {
			if osdDevice.ID == osdId {
				if len(objectNames[osdId]) < int(params.threadsCount)*16 {
					objectNames[osdId] = append(objectNames[osdId], name)
				} else {

				}
			}
			if len(objectNames[osdDevice.ID]) >= int(params.threadsCount)*16 {
				objectsDone++
			}
		}
		if objectsDone >= len(osdDevices) {
			break
		}
	}

	log.Println("Benchmark started")

	for _, osd := range osdDevices {
		wg.Add(1)
		if params.parallel == true {
			go bench(cephConn, osd, &buff, &startBuff, &params, &wg, results, totalLats, osdStatsChan, objectNames[osd.ID])
		} else {
			bench(cephConn, osd, &buff, &startBuff, &params, &wg, results, totalLats, osdStatsChan, objectNames[osd.ID])
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
		sumLat += avgLat.latencyTotal
		countLat += avgLat.len
	}

	//count avg statistics
	sumLat = sumLat / int64(countLat)
	avgIops := 1000000 / sumLat * int64(params.threadsCount)
	sumIops := 1000000 / sumLat * int64(params.threadsCount) * int64(len(osdDevices))
	avgSpeed := float64(avgIops) * float64(params.blockSize) / 1024 / 1024
	sumSpeed := float64(sumIops) * float64(params.blockSize) / 1024 / 1024

	color.Set(color.FgHiYellow)
	defer color.Unset()

	fmt.Printf("Average iops per osd:%9d    Average speed per osd: %.3f MB/s\n"+
		"Total writes count:%11d    Total writes (MB): %v\n",
		avgIops, avgSpeed, countLat, uint64(countLat)*params.blockSize/1024/1024)
	if params.parallel {
		fmt.Printf("Summary avg iops:%13d    Summary avg speed: %.3f MB/s\n", sumIops, sumSpeed)
	}
}
