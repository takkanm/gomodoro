package main

import "time"
import "fmt"
import ui "github.com/gizak/termui"
import "os"

type Gomodoro struct {
	limit    int
	rest_sec int
	mode     string
}

func (g *Gomodoro) dec() {
	if g.rest_sec > 0 {
		g.rest_sec = g.rest_sec - 1
	}
}

func (g Gomodoro) show() {
	min := g.rest_sec / 60
	sec := g.rest_sec % 60

	gauge := ui.NewGauge()
	gauge.BorderLabel = fmt.Sprintf("%02d:%02d\n", min, sec)
	gauge.Width = 100
	gauge.Percent = (g.rest_sec * 100 / g.limit)

	if g.rest_sec < 10 {
		gauge.BarColor = ui.ColorYellow
	} else {
		if g.mode == "work" {
			gauge.BarColor = ui.ColorGreen
		} else {
			gauge.BarColor = ui.ColorBlue
		}
	}

	ui.Render(gauge)
}

func (g Gomodoro) isFinished() bool {
	return g.rest_sec <= 0
}

func (g *Gomodoro) start() {
	g.rest_sec = g.limit
	g.show()
	ticker := time.NewTicker(time.Second)

	go func() {
		for range ticker.C {
			if g.isFinished() {
				ticker.Stop()
				return
			}

			g.dec()
			g.show()
		}
	}()
}

func (g *Gomodoro) stop() {
	g.rest_sec = 0
}

func createGomodoro(mode string) *Gomodoro {
	var gomodoro *Gomodoro
	if mode == "work" {
		gomodoro = &Gomodoro{limit: 25 * 60, mode: "work"}
	} else {
		gomodoro = &Gomodoro{limit: 5 * 60, mode: "rest"}
	}
	return gomodoro
}

func runGomodoro() {
	mode := "work"

	for {
		gomodoro := createGomodoro(mode)
		gomodoro.start()

		for !gomodoro.isFinished() {
			time.Sleep(1)
		}
		if mode == "work" {
			mode = "rest"
		} else {
			mode = "work"
		}
	}
}

func main() {
	go func() {
		var input string
		fmt.Scan(&input)
		ui.Close()
		os.Exit(0)
	}()

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	runGomodoro()
}
