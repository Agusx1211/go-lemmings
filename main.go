package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"github.com/shirou/gopsutil/load"
)

func getSystemLoad() float64 {
	avg, err := load.Avg()
	if err != nil {
		log.Printf("Failed to get system load: %v\n", err)
		return -1
	}
	return avg.Load1
}

func replaceTemplate(template string) string {
	randomValue := rand.Uint32()
	return strings.Replace(template, "{random}", fmt.Sprint(randomValue), -1)
}

func spawnProcess(template string, failed chan string, durations chan time.Duration) {
	start := time.Now()
	cmd := exec.Command("sh", "-c", replaceTemplate(template))
	output, err := cmd.CombinedOutput()
	duration := time.Since(start)

	if err != nil {
		// Error message includes the command, the error and the output
		errMsg := fmt.Sprintf("Failed to run command: %v\nError: %v\nOutput: %s\n", cmd.Args, err, output)
		failed <- errMsg
	} else {
		durations <- duration
	}
}

func monitorAndSpawn(template string, targetLoad float64, maxProcesses, delayMs int, ignoreErr string) {
	failed := make(chan string)
	durations := make(chan time.Duration)
	ignoredErrors := 0

	processes := 0
	completed := 0
	totalDuration := time.Duration(0)
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		if processes < maxProcesses && getSystemLoad() < targetLoad {
			go spawnProcess(template, failed, durations)
			processes++
			time.Sleep(time.Millisecond * time.Duration(delayMs))
			continue
		}

		select {
		case duration := <-durations:
			totalDuration += duration
			processes--
			completed++
			avgDuration := totalDuration / time.Duration(completed)
			log.Printf("Process completed. Errors: %d, Running: %d, Completed %d, Avg Time/Process: %s, Last Process Time: %s",
				ignoredErrors, processes, completed, avgDuration, duration)
		case errMsg := <-failed:
			if ignoreErr != "" && strings.Contains(errMsg, ignoreErr) {
				ignoredErrors++
				log.Printf("Ignoring error")
				continue
			}

			log.Printf("Process failed with output:\n%s\n", errMsg)
			return
		case <-ticker.C:
			if processes >= maxProcesses {
				continue
			}
		}
	}
}

func main() {
	var targetLoad float64
	var maxProcesses int
	var delayMs int
	var ignoreErr string

	flag.Float64Var(&targetLoad, "l", 0.5, "target system load")
	flag.IntVar(&maxProcesses, "m", 10, "maximum number of processes to spawn")
	flag.IntVar(&delayMs, "d", 0, "delay in milliseconds after executing a process")
	flag.StringVar(&ignoreErr, "i", "", "error string to ignore")
	flag.Parse()

	template := strings.Join(flag.Args(), " ")
	if template == "" {
		fmt.Println("Usage: go-lemmings -l [target load] -m [max processes] -d [delay ms] -i [ignore err] [template]")
		return
	}

	if targetLoad <= 0 || maxProcesses <= 0 {
		fmt.Println("target load and max processes must be greater than 0")
		return
	}

	monitorAndSpawn(template, targetLoad, maxProcesses, delayMs, ignoreErr)
}
