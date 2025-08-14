package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	"fyne.io/systray"
	"github.com/shirou/gopsutil/cpu"
)

var icons [][]byte
var showCPUUsage bool
var mu sync.Mutex

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTooltip("CPU Usage")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	mShowUsage := systray.AddMenuItemCheckbox("Show CPU Usage", "Toggle CPU usage in menu bar", false)

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			case <-mShowUsage.ClickedCh:
				mu.Lock()
				// Toggle the internal state
				showCPUUsage = !showCPUUsage
				// Update the UI to match the new state
				if showCPUUsage {
					mShowUsage.Check()
				} else {
					mShowUsage.Uncheck()
				}
				mu.Unlock()
			}
		}
	}()

	// Load icons
	for i := 0; i < 5; i++ {
		iconPath := filepath.Join("assets", fmt.Sprintf("%d.png", i))
		icon, err := os.ReadFile(iconPath)
		if err != nil {
			log.Fatalf("Could not load icon: %v", err)
		}
		icons = append(icons, icon)
	}

	speedUpdateChan := make(chan time.Duration)
	go animateCat(speedUpdateChan)
	go monitorCPU(speedUpdateChan)
}

func onExit() {
	// clean up here
	log.Println("Exiting")
}

func animateCat(speedUpdateChan <-chan time.Duration) {
	sleepDuration := 500 * time.Millisecond
	ticker := time.NewTicker(sleepDuration)
	iconIndex := 0

	for {
		select {
		case newDuration := <-speedUpdateChan:
			if newDuration != sleepDuration {
				ticker.Reset(newDuration)
				sleepDuration = newDuration
			}
		case <-ticker.C:
			systray.SetIcon(icons[iconIndex])
			iconIndex = (iconIndex + 1) % len(icons)
		}
	}
}

func monitorCPU(speedUpdateChan chan<- time.Duration) {
	for {
		percentages, err := cpu.Percent(time.Second, false)
		if err != nil {
			log.Printf("Error getting CPU usage: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		if len(percentages) > 0 {
			cpuUsage := percentages[0]

			mu.Lock()
			show := showCPUUsage
			mu.Unlock()

			usageStr := fmt.Sprintf("CPU Usage: %.1f%%", cpuUsage)
			systray.SetTooltip(usageStr)

			if show {
				systray.SetTitle(fmt.Sprintf(" %.1f%%", cpuUsage))
			} else {
				systray.SetTitle("")
			}

			divisor := math.Max(1.0, math.Min(20.0, cpuUsage/5.0))
			intervalSeconds := 0.2 / divisor
			newSpeed := time.Duration(intervalSeconds * float64(time.Second))

			speedUpdateChan <- newSpeed
		}
		time.Sleep(1 * time.Second)
	}
}
