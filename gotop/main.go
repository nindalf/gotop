package main

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui"
	"github.com/nindalf/gotop/daemon"
)

func updateMemCPU(gcpu, gmem *ui.Gauge) {
	done := make(chan struct{})
	memInfoChan, errc := daemon.TotalMemory(done, daemon.Delay)
	cpuInfoChan, errc := daemon.TotalCPU(done, daemon.Delay)
	timeout := time.After(2 * daemon.Delay)
	defer func() {
		close(done)
		// Necessary to read from error channel to prevent sending goroutine going into deadlock
		<-errc
	}()
	for {
		select {
		case memInfo := <-memInfoChan:
			pc := int(100 * (float64(memInfo.MemTotal-memInfo.MemFree) / float64(memInfo.MemTotal)))
			gmem.Percent = pc
			ui.Render(ui.Body)
		case cpuInfo := <-cpuInfoChan:
			pc := int(cpuInfo.AverageUtilization)
			gcpu.Percent = pc
			ui.Render(ui.Body)
		case err := <-errc:
			if err != nil {
			}
			return
		case <-timeout:
		}
	}
}

func updateUptime(p *ui.Par) {
	for {
		<-time.After(2 * daemon.Delay)
		ut, _ := daemon.Uptime()
		p.Text = ut.String()
		ui.Render(ui.Body)
	}
}

func main() {
	err := ui.Init()
	fmt.Println(daemon.UpSince())
	if err != nil {
		fmt.Println("Could not initialise UI")
	}
	defer ui.Close()

	ut, _ := daemon.Uptime()
	p := ui.NewPar(ut.String())
	p.Height = 3
	p.Width = 50
	p.TextFgColor = ui.ColorWhite
	p.Border.Label = "Uptime"
	p.Border.FgColor = ui.ColorCyan

	g0 := ui.NewGauge()
	g0.Percent = 40
	g0.Width = 50
	g0.Height = 3
	g0.Border.Label = "Memory"
	g0.BarColor = ui.ColorRed
	g0.Border.FgColor = ui.ColorWhite
	g0.Border.LabelFgColor = ui.ColorCyan

	g2 := ui.NewGauge()
	g2.Percent = 60
	g2.Width = 50
	g2.Height = 3
	g2.PercentColor = ui.ColorBlue
	g2.Y = 3
	g2.Border.Label = "CPU"
	g2.BarColor = ui.ColorYellow
	g2.Border.FgColor = ui.ColorWhite

	ui.Body.AddRows(ui.NewRow(ui.NewCol(6, 0, g0), ui.NewCol(6, 0, p)),
		ui.NewRow(ui.NewCol(6, 0, g2)))
	ui.Body.Align()
	ui.Render(ui.Body)
	go updateMemCPU(g2, g0)
	go updateUptime(p)
	<-ui.EventCh()
}
