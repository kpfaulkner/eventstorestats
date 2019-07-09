package main

import (
  "github.com/kpfaulkner/eventstorestats/models"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)
func main() {

	fmt.Printf("So it begins......\n\n\n")

	filename := flag.String("f", "", "the filename to process")
	cpu := flag.Bool("cpu", false, "CPU usage")
	mem := flag.Bool("mem", false, "Memory usage")
	queues := flag.Bool("queues", false, "Current Queue lengths")

	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		log.Fatalf("Unable to process file: %s\n", err.Error())
	}

	fileSize := stats.Size()

	// will just read last 10000 bytes of file
	if fileSize > 20000 {
		file.Seek(int64(fileSize-int64(20000)), 0)
	}
	scanner := bufio.NewReader(file)

	for {
		line, err := scanner.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// EOF of file... just sleep a bit and continue again
				time.Sleep(time.Duration(1) * time.Second)
			} else {
				fmt.Printf("break %s\n", err.Error())
				break
			}
		} else {

			statsLine := string(line)
			proc := models.ProcStats{}
			json.Unmarshal([]byte(statsLine), &proc)

			displayStats(proc, *cpu, *mem, *queues)
		}
	}
}

// displayStats just writes to stdout with various statistics about the EV node.
func displayStats( proc models.ProcStats, cpu bool, mem bool, queues bool) {

	fmt.Printf("==========================================================\n")
	if cpu {
		fmt.Printf("Proc CPU: %f\n", proc.Proc.CPU)
		fmt.Printf("Sys  CPU: %f\n", proc.Sys.CPU)
	}

	if mem {
		fmt.Printf("Proc MEM: %d\n", proc.Proc.Mem)
		fmt.Printf("Sys Free MEM: %d\n", proc.Sys.FreeMem)
	}

	// if want queues, loop through all and display details.
	if queues {
    displayQueues( proc.Es.Queue.IndexCommitter)
		displayQueues( proc.Es.Queue.MainQueue)
		displayQueues( proc.Es.Queue.MasterReplicationService)
		displayQueues( proc.Es.Queue.MonitoringQueue)
		displayQueues( proc.Es.Queue.PersistentSubscriptions)
		displayQueues( proc.Es.Queue.ProjectionCore0)
		displayQueues( proc.Es.Queue.ProjectionCore1)
		displayQueues( proc.Es.Queue.ProjectionCore2)
		displayQueues( proc.Es.Queue.ProjectionsMaster)
	}
	fmt.Printf("==========================================================\n")
}

// displayQueues writes out queue name and some interesting stats.
func displayQueues( queue models.ESQueue) {
	fmt.Printf("Queue: %-30s Avg IPS: %d\tLength: %d\tTotal Items Processed %d\n", queue.QueueName, queue.AvgItemsPerSecond, queue.Length, queue.TotalItemsProcessed)
}
