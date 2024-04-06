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

func monitorAndSpawn(template string, targetLoad float64, maxProcesses int) {
	failed := make(chan string)
	durations := make(chan time.Duration)

	processes := 0
	completed := 0
	totalDuration := time.Duration(0)
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		if processes < maxProcesses && getSystemLoad() < targetLoad {
			go spawnProcess(template, failed, durations)
			processes++
			continue
		}

		select {
		case duration := <-durations:
			totalDuration += duration
			processes--
			completed++
			// Avoid division by zero
			avgDuration := totalDuration / time.Duration(completed)
			log.Printf("Process completed. Running: %d, Completed %d, Avg Time/Process: %s, Last Process Time: %s",
				processes, completed, avgDuration, duration)
		case output := <-failed:
			// stop spawning new processes and exit
			log.Printf("Process failed with output:\n%s\n", output)
			return
		case <-ticker.C:
			// Skip waiting for ticker if maxProcesses reached
			if processes >= maxProcesses {
				continue
			}
		}
	}
}

func main() {
	var targetLoad float64
	var maxProcesses int

	flag.Float64Var(&targetLoad, "l", 0.5, "target system load")
	flag.IntVar(&maxProcesses, "m", 10, "maximum number of processes to spawn")
	flag.Parse()

	// Concatenate all remaining command-line arguments into a single template string
	template := strings.Join(flag.Args(), " ")
	if template == "" {
		fmt.Println("Usage: go-lemmings -l [target load] -m [max processes] [template]")
		return
	}

	if targetLoad <= 0 || maxProcesses <= 0 {
		fmt.Println("target load and max processes must be greater than 0")
		return
	}

	monitorAndSpawn(template, targetLoad, maxProcesses)
}
