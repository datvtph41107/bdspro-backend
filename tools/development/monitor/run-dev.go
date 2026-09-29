package main

import (
	"os/exec"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func checkProcess(name string) string {
	cmd := exec.Command("pgrep", "-f", name)
	if err := cmd.Run(); err != nil {
		return "❌ Stopped"
	}
	return "✅ Running"
}

func main() {
	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)

	services := []string{"user-service", "order-service", "payment-service"}

	// Goroutine update bảng mỗi 2 giây
	go func() {
		for {
			app.QueueUpdateDraw(func() {
				// Header
				headers := []string{"Service", "Status"}
				for i, h := range headers {
					table.SetCell(0, i,
						tview.NewTableCell(h).
							SetSelectable(false).
							SetAlign(tview.AlignCenter).
							SetTextColor(tcell.ColorYellow).
							SetAttributes(tcell.AttrBold).
							SetBackgroundColor(tcell.ColorDefault),
					)
				}

				// Rows
				for i, s := range services {
					status := checkProcess(s)
					color := tcell.ColorGreen
					if status == "❌ Stopped" {
						color = tcell.ColorRed
					}

					table.SetCell(i+1, 0,
						tview.NewTableCell(s).
							SetAlign(tview.AlignLeft).
							SetTextColor(tcell.ColorWhite).
							SetBackgroundColor(tcell.ColorDefault),
					)

					table.SetCell(i+1, 1,
						tview.NewTableCell(status).
							SetAlign(tview.AlignCenter).
							SetTextColor(color).
							SetBackgroundColor(tcell.ColorDefault),
					)
				}
			})
			time.Sleep(2 * time.Second)
		}
	}()

	// false = không fill hết màn hình → giữ nền terminal
	if err := app.SetRoot(table, false).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
