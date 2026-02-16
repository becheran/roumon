package ui

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/becheran/roumon/internal/model"
	"github.com/gizak/termui/v3/widgets"

	termui "github.com/gizak/termui/v3"
)

const (
	padding         = 1
	keepRoutineHist = 100
)

// statusAbbreviations maps common goroutine status names to consistent abbreviations
var statusAbbreviations = map[string]string{
	"running":         "run",
	"runnable":        "rbl",
	"waiting":         "wai",
	"IO wait":         "IO",
	"chan receive":    "cha",
	"chan send":       "ch1",
	"select":          "sel",
	"sync.Mutex.Lock": "syn",
	"sync.Cond.Wait":  "syw",
	"syscall":         "sys",
	"sleep":           "slp",
	"idle":            "idl",
	"dead":            "ded",
	"copystack":       "cps",
	"preempted":       "pre",
	"GC assist wait":  "gca",
	"GC sweep wait":   "gcs",
	"GC scavenge wait": "gcv",
}

// UI contains all user interface elements
type UI struct {
	list           *widgets.List
	filter         *widgets.Paragraph
	details        *widgets.Paragraph
	routineHist    *widgets.Plot
	barchart       *widgets.BarChart
	barchartLegend *widgets.Paragraph
	paused         *widgets.Paragraph
	legend         *widgets.Paragraph
	help           *widgets.Paragraph

	grid          *termui.Grid
	filtered      bool
	origData      []model.Goroutine
	filteredData  []model.Goroutine
	minGoRoutines int
	maxGoRoutines int
	avgGoRoutines float64
}

// NewUI creates a new console user interface
func NewUI() *UI {
	if err := termui.Init(); err != nil {
		log.Fatalf("Failed to initialize termui: %v", err)
	}

	filter := widgets.NewParagraph()
	filter.Text = "TYPE TO FILTER"
	filter.TextStyle.Fg = termui.ColorWhite
	filter.BorderStyle.Fg = termui.ColorGreen
	filter.Title = "Filter"
	filter.PaddingTop = padding
	filter.PaddingRight = padding
	filter.PaddingLeft = padding
	filter.PaddingBottom = padding

	plot := widgets.NewPlot()
	plot.Data = make([][]float64, 1)
	plot.Data[0] = make([]float64, 2, keepRoutineHist)
	plot.AxesColor = termui.ColorWhite
	plot.LineColors[0] = termui.ColorGreen
	plot.HorizontalScale = 2
	plot.PaddingTop = padding
	plot.PaddingRight = padding
	plot.PaddingLeft = padding
	plot.PaddingBottom = padding

	routineList := widgets.NewList()
	routineList.PaddingTop = padding
	routineList.PaddingRight = padding
	routineList.PaddingLeft = padding
	routineList.PaddingBottom = padding
	routineList.Rows = []string{}
	routineList.TextStyle.Fg = termui.ColorGreen
	routineList.SelectedRowStyle.Fg = termui.ColorWhite
	routineList.SelectedRowStyle.Bg = termui.ColorGreen

	details := widgets.NewParagraph()
	details.PaddingTop = padding
	details.PaddingRight = padding
	details.PaddingLeft = padding
	details.PaddingBottom = padding
	details.Title = "Details"
	details.TextStyle = termui.NewStyle(termui.ColorWhite)
	details.SetRect(0, 0, 60, 10)

	barchart := widgets.NewBarChart()
	barchart.Title = "Status"
	barchart.BarWidth = 3
	barchart.BarGap = 1
	barchart.BarColors = []termui.Color{termui.ColorGreen}
	barchart.LabelStyles = []termui.Style{termui.NewStyle(termui.ColorWhite)}
	barchart.PaddingTop = padding
	barchart.PaddingRight = padding
	barchart.PaddingLeft = padding
	barchart.PaddingBottom = padding
	barchart.BorderRight = false

	barchartLabel := widgets.NewParagraph()
	barchartLabel.BorderLeft = false
	barchartLabel.PaddingTop = padding
	barchartLabel.PaddingRight = padding
	barchartLabel.PaddingLeft = padding
	barchartLabel.PaddingBottom = padding
	barchartLabel.Text = ""

	help := widgets.NewParagraph()
	help.TextStyle.Fg = termui.ColorGreen
	help.Text = "Help\n\nArrows up/down: Select from list\nText input: Filter results\nF10: Quit\nF2: Pause\n\nPress any key to continue"
	help.PaddingBottom = 2
	help.PaddingLeft = 2
	help.PaddingRight = 2
	help.PaddingTop = 2

	paused := widgets.NewParagraph()
	paused.TextStyle.Fg = termui.ColorGreen
	paused.Text = "Paused. Press any key to continue"
	paused.PaddingBottom = 2
	paused.PaddingLeft = 2
	paused.PaddingRight = 2
	paused.PaddingTop = 2

	legend := widgets.NewParagraph()
	legend.Text = "F1 Help | F2 Pause | F10 Quit"
	legend.TextStyle.Fg = termui.ColorGreen
	legend.Border = false

	grid := termui.NewGrid()

	ui := UI{
		filter:         filter,
		list:           routineList,
		details:        details,
		routineHist:    plot,
		barchart:       barchart,
		barchartLegend: barchartLabel,
		help:           help,
		paused:         paused,
		legend:         legend,
		grid:           grid,
	}

	grid.Set(
		termui.NewRow(3.0/10,
			termui.NewCol(3.0/10,
				termui.NewCol(5.0/8, ui.barchart),
				termui.NewCol(3.0/8, ui.barchartLegend)),
			termui.NewCol(7.0/10, ui.routineHist),
		),
		termui.NewRow(7.0/10,
			termui.NewCol(1.0/6,
				termui.NewRow(1.5/10, ui.filter),
				termui.NewRow(8.5/10, ui.list)),
			termui.NewCol(5.0/6, ui.details),
		),
	)

	ui.updatePlotTitle()

	return &ui
}

func (ui *UI) updatePlotTitle() {
	ui.routineHist.Title = fmt.Sprintf("History # goroutines (Min: %d Avg: %0.2f Max: %d)",
		ui.minGoRoutines, ui.avgGoRoutines, ui.maxGoRoutines)
}

func (ui *UI) updateStatus() {
	typeCount := make(map[string]float64)
	for i := 0; i < len(ui.origData); i++ {
		num := typeCount[ui.origData[i].Status]
		typeCount[ui.origData[i].Status] = num + 1
	}

	types := make([]string, 0, len(typeCount))
	for key := range typeCount {
		types = append(types, key)
	}
	sort.Strings(types)
	data := make([]float64, len(types))
	labels := make([]string, len(types))
	legend := ""
	usedAbbrevs := make(map[string]bool)

	for idx, t := range types {
		data[idx] = typeCount[t]

		// Use predefined abbreviation or fallback to custom logic
		abbrev, exists := statusAbbreviations[t]
		if !exists {
			// Fallback: use first 3 characters
			abbrev = t[:min(3, len(t))]
		}
		
		// Handle collisions by adding a number suffix
		originalAbbrev := abbrev
		counter := 2
		for usedAbbrevs[abbrev] {
			abbrev = fmt.Sprintf("%s%d", originalAbbrev[:min(2, len(originalAbbrev))], counter)
			counter++
		}
		usedAbbrevs[abbrev] = true

		labels[idx] = abbrev
		// Add to legend with count
		legend = fmt.Sprintf("%s%s: %s (%.0f)\n", legend, abbrev, t, typeCount[t])
	}
	ui.barchart.Data = data
	ui.barchart.Labels = labels
	ui.barchartLegend.Text = legend
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (ui *UI) updateList() {
	if ui.filter.Text == "" || !ui.filtered {
		ui.filteredData = ui.origData
	} else {
		ui.filteredData = make([]model.Goroutine, 0)
		for _, d := range ui.origData {
			filterText := strings.ToLower(ui.filter.Text)
			matchID := strings.Contains(strings.ToLower(fmt.Sprintf("%d", d.ID)), filterText)
			matchStatus := strings.Contains(strings.ToLower(d.Status), filterText)
			matchCreatedBy := d.CratedBy != nil && strings.Contains(strings.ToLower(d.CratedBy.String()), filterText)
			matchStackTrace := model.StackContains(d.StackTrace, filterText)
			matchLockedToThread := d.LockedToThread && strings.Contains("locked to thread", filterText)
			if matchStatus || matchID || matchCreatedBy || matchStackTrace || matchLockedToThread {
				ui.filteredData = append(ui.filteredData, d)
			}
		}
	}

	// Update list
	ui.list.Rows = make([]string, len(ui.filteredData))
	for i := 0; i < len(ui.filteredData); i++ {
		ui.list.Rows[i] = fmt.Sprintf("%05d %s ", ui.filteredData[i].ID, ui.filteredData[i].Status)
	}

	if len(ui.filteredData) == 0 {
		ui.list.SelectedRow = 0
		ui.details.Text = ""
		ui.list.Title = "Routines (0/0)"
		return
	}

	if ui.list.SelectedRow >= len(ui.filteredData) {
		ui.list.SelectedRow = len(ui.filteredData) - 1
	} else if ui.list.SelectedRow < 0 {
		ui.list.SelectedRow = 0
	}

	selectedData := ui.filteredData[ui.list.SelectedRow]
	trace := ""
	for _, t := range selectedData.StackTrace {
		trace += fmt.Sprintf("  %s\n", t.String())
	}
	createdBy := ""
	if selectedData.CratedBy != nil {
		createdBy = fmt.Sprintf("Created by:\n  %s\n\n", selectedData.CratedBy.String())
	}
	lockedToThread := ""
	if selectedData.LockedToThread {
		lockedToThread = " [locked to thread](mod:bold)"
	}
	ui.details.Text = fmt.Sprintf("ID: [%d](mod:bold)\n\nStatus: [%s](mod:bold)\n\nWait Since: [%d min](mod:bold)%s\n\n%sTrace:\n%s",
		selectedData.ID,
		selectedData.Status,
		selectedData.WaitSinceMin,
		lockedToThread,
		createdBy,
		trace)

	ui.list.Title = fmt.Sprintf("Routines (%d/%d)", ui.list.SelectedRow+1, len(ui.list.Rows))
}

// Stop UI and close all event listeners
func (ui *UI) Stop() {
	termui.Close()
}

func (ui *UI) resize(width, height int) {
	log.Printf("Resize to: (%d,%d)", width, height)
	ui.paused.SetRect(width/2.0-25, height/4.0-4, width/2.0+25, height/4.0+4)
	ui.help.SetRect(width/2.0-20, height/4.0-10, width/2.0+20, height/4.0+10)
	ui.legend.SetRect(width-35, height-4, width-1, height-1)
	ui.grid.SetRect(0, 0, width, height)
}

// Run UI in fullscreen mode
func (ui *UI) Run(terminate chan<- error, routinesUpdate <-chan []model.Goroutine) {
	ui.updateList()

	termWidth, termHeight := termui.TerminalDimensions()
	ui.resize(termWidth, termHeight)

	termui.Render(ui.grid, ui.legend)

	pollEvents := termui.PollEvents()
	for {
		select {
		case evt := <-pollEvents:
			switch evt.Type {
			case termui.MouseEvent:
				continue
			case termui.ResizeEvent:
				resized, ok := evt.Payload.(termui.Resize)
				if !ok {
					log.Printf("Failed to parse payload for resize. %v", evt)
				} else {
					ui.resize(resized.Width, resized.Height)
				}
			case termui.KeyboardEvent:
				terminateEvent := ui.handleKeyEvent(evt.ID, pollEvents)
				if terminateEvent {
					terminate <- nil
					return
				}
			}
		case routines := <-routinesUpdate:
			// History data size cannot be limited in termui. This is a workaround
			var keepRoutineHist = (ui.routineHist.Dx() - 10) >> 1
			ui.origData = routines
			if len(ui.routineHist.Data[0]) >= keepRoutineHist {
				ui.routineHist.Data[0] = ui.routineHist.Data[0][1:]
			}
			ui.routineHist.Data[0] = append(ui.routineHist.Data[0], float64(len(routines)))

			if ui.minGoRoutines == 0 || len(routines) < ui.minGoRoutines {
				ui.minGoRoutines = len(routines)
			}
			if len(routines) > ui.maxGoRoutines {
				ui.maxGoRoutines = len(routines)
			}
			if ui.avgGoRoutines > 0 {
				ui.avgGoRoutines = (ui.avgGoRoutines + float64(len(routines))) / 2.0
			} else {
				ui.avgGoRoutines = float64(len(routines))
			}
			ui.updatePlotTitle()
			ui.updateList()
			ui.updateStatus()
		}

		termui.Render(ui.grid, ui.legend)
	}
}

func (ui *UI) handleKeyEvent(keyID string, pollEvents <-chan termui.Event) (terminate bool) {
	switch keyID {
	case "<C-c>", "<F10>":
		return true
	case "<F1>":
		termui.Render(ui.grid, ui.legend, ui.help)
		e := <-pollEvents
		if e.ID == "<C-c>" || e.ID == "<F10>" {
			return true
		}
		termui.Render(ui.grid, ui.legend)
	case "<F2>":
		// Pause
		termui.Render(ui.grid, ui.legend, ui.paused)
		e := <-pollEvents
		if e.ID == "<C-c>" || e.ID == "<F10>" {
			return true
		}
		termui.Render(ui.grid, ui.legend)
	case "<Down>":
		ui.list.ScrollDown()
		ui.updateList()
	case "<Up>":
		ui.list.ScrollUp()
		ui.updateList()
	case "<PageDown>":
		ui.list.ScrollPageDown()
		ui.updateList()
	case "<PageUp>":
		ui.list.ScrollPageUp()
		ui.updateList()
	case "<Home>":
		ui.list.ScrollTop()
		ui.updateList()
	case "<End>":
		ui.list.ScrollBottom()
		ui.updateList()
	case "<Backspace>", "<C-<Backspace>>":
		if len(ui.filter.Text) > 0 {
			ui.filter.Text = ui.filter.Text[:len(ui.filter.Text)-1]
		}
	case "<Space>":
		if !ui.filtered {
			ui.filter.Text = ""
		}
		ui.filtered = true
		ui.filter.Text += " "
		ui.updateList()
	default:
		// < sign
		if keyID[0] != 0x3C {
			if !ui.filtered {
				ui.filter.Text = ""
			}
			ui.filtered = true
			ui.filter.Text += keyID
		}
		ui.updateList()
	}
	return false
}
