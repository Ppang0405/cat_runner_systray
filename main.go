package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
	"time"

	"fyne.io/systray"
	"github.com/shirou/gopsutil/cpu"
)

var icons [][]byte
var cpuUsage float64
var mu sync.Mutex

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	// systray.SetTitle("Cat Runner")
	systray.SetTooltip("CPU Usage")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	// Load icons
	for i := 0; i < 5; i++ {
		iconPath := filepath.Join("assets", fmt.Sprintf("%d.png", i))
		icon, err := ioutil.ReadFile(iconPath)
		if err != nil {
			log.Fatalf("Could not load icon: %v", err)
		}
		icons = append(icons, icon)
	}

	go animateCat()
	go monitorCPU()
}

func onExit() {
	// clean up here
	log.Println("Exiting")
}

func animateCat() {
	iconIndex := 0
	for {
		mu.Lock()
		currentCPUUsage := cpuUsage
		mu.Unlock()

		var sleepDuration time.Duration
		if currentCPUUsage > 50 {
			sleepDuration = 100 * time.Millisecond
		} else if currentCPUUsage > 20 {
			sleepDuration = 300 * time.Millisecond
		} else {
			sleepDuration = 500 * time.Millisecond
		}

		systray.SetIcon(icons[iconIndex])
		iconIndex = (iconIndex + 1) % len(icons)
		time.Sleep(sleepDuration)
	}
}

func monitorCPU() {
	for {
		percentages, err := cpu.Percent(time.Second, false)
		if err != nil {
			log.Printf("Error getting CPU usage: %v", err)
			continue
		}
		if len(percentages) > 0 {
			mu.Lock()
			cpuUsage = percentages[0]
			mu.Unlock()
			systray.SetTooltip(fmt.Sprintf("CPU Usage: %.2f%%", cpuUsage))
		}
	}
}
