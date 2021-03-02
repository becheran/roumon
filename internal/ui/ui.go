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
	padding = 1
)

var keepRoutineHist = 100

type UI struct {
	list           *widgets.List
	filter         *widgets.Paragraph
	details        *widgets.Paragraph
	routineHist    *widgets.Plot
	barchart       *widgets.BarChart
	barchartLegend *widgets.Paragraph

	filtered     bool
	origData     []model.Goroutine
	filteredData []model.Goroutine
}

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
	plot.Title = "History # goroutines"
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
	barchart.NumStyles = []termui.Style{termui.NewStyle(termui.ColorBlack)}
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

	ui := UI{
		filter:         filter,
		list:           routineList,
		details:        details,
		routineHist:    plot,
		barchart:       barchart,
		barchartLegend: barchartLabel,
	}

	ui.updateList()

	return &ui
}

func sliceContains(ss []string, subString string) bool {
	for _, s := range ss {
		if strings.Contains(s, subString) {
			return true
		}
	}
	return false
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
	label := ""
	uniqueID := 1
	for idx, t := range types {
		data[idx] = typeCount[t]
		newLabel := t[:3]
		if sliceContains(labels, newLabel) {
			newLabel = fmt.Sprintf("%s%d", t[:2], uniqueID)
			uniqueID++
		}
		labels[idx] = newLabel
		label = fmt.Sprintf("%s%s: %s\n", label, newLabel, t)
	}
	ui.barchart.Data = data
	ui.barchart.Labels = labels
	ui.barchartLegend.Text = label
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

	// TODO: prevent always doing this!
	routineList := make([]string, len(ui.filteredData))
	for i := 0; i < len(ui.filteredData); i++ {
		routineList[i] = fmt.Sprintf("%05d %s ", ui.filteredData[i].ID, ui.filteredData[i].Status)
	}
	ui.list.Rows = routineList

	if len(ui.filteredData) == 0 {
		ui.details.Text = ""
		ui.list.Title = "Routines (0/0)"
		return
	}
	if ui.list.SelectedRow >= len(ui.filteredData) {
		ui.list.SelectedRow = len(ui.filteredData) - 1
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

func (ui *UI) Stop() {
	termui.Close()
}

func (ui *UI) Run(terminate chan<- error, routinesUpdate <-chan []model.Goroutine) {
	grid := termui.NewGrid()

	paused := widgets.NewParagraph()
	paused.TextStyle.Fg = termui.ColorGreen
	paused.Text = "Paused. Press F2 to continue"
	paused.PaddingBottom = 2
	paused.PaddingLeft = 2
	paused.PaddingRight = 2
	paused.PaddingTop = 2

	help := widgets.NewParagraph()
	help.TextStyle.Fg = termui.ColorGreen
	help.Text = "Help\n\nArrows up/down: Select from list\nText input: Filter results\nF10: Quit\nF2: Pause\n\nPress any key to continue"
	help.PaddingBottom = 2
	help.PaddingLeft = 2
	help.PaddingRight = 2
	help.PaddingTop = 2

	legend := widgets.NewParagraph()
	legend.Text = "F1 Help | F2 Pause | F10 Quit"
	legend.TextStyle.Fg = termui.ColorGreen
	legend.Border = false

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

	resize := func(width, height int) {
		log.Printf("Resize to: (%d,%d)", width, height)
		paused.SetRect(width/2.0-17, height/4.0-4, width/2.0+17, height/4.0+4)
		help.SetRect(width/2.0-20, height/4.0-10, width/2.0+20, height/4.0+10)
		legend.SetRect(width-35, height-4, width-1, height-1)
		grid.SetRect(0, 0, width, height)
	}

	termWidth, termHeight := termui.TerminalDimensions()
	resize(termWidth, termHeight)

	termui.Render(grid, legend)

	pollEvents := termui.PollEvents()
loop:
	for {
		select {
		case evt := <-pollEvents:
			if evt.Type == termui.MouseEvent {
				continue loop
			} else if evt.Type == termui.ResizeEvent {
				resized, ok := evt.Payload.(termui.Resize)
				if !ok {
					log.Printf("Failed to parse payload for resize. %v", evt)
					continue loop
				}
				resize(resized.Width, resized.Height)
			} else {
				// Handle keyboard events
				switch evt.ID {
				case "<C-c>", "<F10>":
					terminate <- nil
					return
				case "<F1>":
					termui.Render(grid, legend, help)
					e := <-termui.PollEvents()
					if e.ID == "<C-c>" || e.ID == "<F10>" {
						terminate <- nil
						return
					}
				case "<F2>":
					// Pause
					termui.Render(grid, legend, paused)
				paused:
					for {
						e := <-termui.PollEvents()
						switch e.ID {
						case "<F2>":
							break paused
						case "<C-c>", "<F10>":
							terminate <- nil
							return
						}
					}
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
					if evt.ID[0] != 0x3C {
						if !ui.filtered {
							ui.filter.Text = ""
						}
						ui.filtered = true
						ui.filter.Text += evt.ID
					}
					ui.updateList()
				}

			}
		case routines := <-routinesUpdate:
			// Hacky, but the length is needed here!
			var keepRoutineHist = (ui.routineHist.Dx() - 10) >> 1
			ui.origData = routines
			if len(ui.routineHist.Data[0]) >= keepRoutineHist {
				ui.routineHist.Data[0] = ui.routineHist.Data[0][1:]
			}
			ui.routineHist.Data[0] = append(ui.routineHist.Data[0], float64(len(routines)))

			ui.updateList()
			ui.updateStatus()
		}

		termui.Render(grid, legend)
	}
}
