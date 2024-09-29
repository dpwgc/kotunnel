package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"kotunnel/base"
	"kotunnel/cli"
	"os"
	"slices"
	"strings"
	"sync"
)

const HisMax = 6

type ButtonLock struct {
	m map[string]bool
	l sync.Mutex
}

var buttonLock = &ButtonLock{
	m: make(map[string]bool),
	l: sync.Mutex{},
}

func (bl *ButtonLock) Lock(key string) bool {
	bl.l.Lock()
	defer bl.l.Unlock()
	if bl.m[key] {
		return false
	}
	bl.m[key] = true
	return true
}

func (bl *ButtonLock) Unlock(key string) {
	bl.l.Lock()
	defer bl.l.Unlock()
	bl.m[key] = false
}

type Cache struct {
	Last    []string
	History [][]string
}

func main() {

	myApp := app.New()
	myWindow := myApp.NewWindow("KoTunnel GUI")
	myWindow.Resize(fyne.Size{
		Width: 700,
	})

	var cache Cache
	data, _ := os.ReadFile("./gui_cache")
	_ = json.Unmarshal(data, &cache)

	protocol := widget.NewSelect([]string{"TCP", "UDP"}, func(s string) {})
	protocol.SetSelected("TCP")

	secret := widget.NewPasswordEntry()
	secret.SetPlaceHolder("123456")

	tunnelAddr := widget.NewEntry()
	tunnelAddr.SetPlaceHolder("127.0.0.1:8080")

	localPort := widget.NewEntry()
	localPort.SetPlaceHolder("9090")

	idleNum := widget.NewEntry()
	idleNum.SetPlaceHolder("5")

	tip := widget.NewLabel("")

	runButton := widget.NewButton("Run", func() {
		go func() {
			newAdd := []string{strings.ToLower(protocol.Selected), secret.Text, tunnelAddr.Text, localPort.Text, idleNum.Text}
			newAddToString := strings.Join(newAdd, ",")
			if !buttonLock.Lock(newAddToString) {
				tip.SetText("This tunnel has been started, do not repeat the operation!")
				return
			}
			defer func() {
				buttonLock.Unlock(newAddToString)
			}()
			tip.SetText("This tunnel has been started, do not close the window!")
			cache.Last = newAdd
			if len(cache.History) <= 0 {
				cache.History = [][]string{newAdd}
			} else {
				// 检查一下，历史记录里有一样的，就不往里加了
				var check []string
				for _, v := range cache.History {
					check = append(check, strings.Join(v, ","))
				}
				if !slices.Contains(check, newAddToString) {
					cache.History = append(cache.History, newAdd)
				}
			}
			if len(cache.History) > HisMax {
				cache.History = cache.History[len(cache.History)-HisMax:]
			}

			bytes, _ := json.Marshal(cache)
			_ = os.WriteFile("./gui_cache", bytes, 0644)

			base.InitConfig([]string{"", "client", newAdd[0], newAdd[1], newAdd[2], newAdd[3], newAdd[4]})
			base.InitLog()
			client(base.Config().App)
		}()
	})

	if len(cache.Last) > 0 {
		if cache.Last[0] == "tcp" || cache.Last[0] == "udp" {
			protocol.SetSelected(strings.ToUpper(cache.Last[0]))
		}
		secret.SetText(cache.Last[1])
		tunnelAddr.SetText(cache.Last[2])
		localPort.SetText(cache.Last[3])
		idleNum.SetText(cache.Last[4])
	}

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Widget: widget.NewLabel("")},
			{Text: "  Protocol", Widget: protocol},
			{Text: "  Secret", Widget: secret},
			{Text: "  Tunnel", Widget: tunnelAddr},
			{Text: "  Local", Widget: localPort},
			{Text: "  Idle", Widget: idleNum},
			{Widget: runButton},
			{Text: "  ", Widget: tip},
		},
	}

	// 创建一个容器来模拟表格
	table := container.NewVBox()

	table.Add(widget.NewLabel(""))
	for index, item := range cache.History {

		line := container.NewHBox()

		// 创建一个按钮
		button := widget.NewButton("Use", func(i int) func() {
			return func() {
				row := cache.History[i]
				if row[0] == "tcp" || row[0] == "udp" {
					protocol.SetSelected(strings.ToUpper(row[0]))
				}
				secret.SetText(row[1])
				tunnelAddr.SetText(row[2])
				localPort.SetText(row[3])
				idleNum.SetText(row[4])
			}
		}(index))

		line.Add(widget.NewLabel(" "))
		line.Add(button)
		line.Add(widget.NewLabel(strings.ToUpper(item[0])))
		addr := item[2]
		if len(addr) > 30 {
			addr = addr[:30] + "..."
		}
		line.Add(widget.NewLabel(addr))
		line.Add(widget.NewLabel(item[3]))
		line.Add(widget.NewLabel(item[4]))
		line.Add(widget.NewLabel(" "))
		table.Add(line)
	}
	table.Add(widget.NewLabel(""))

	myWindow.SetContent(container.New(layout.NewGridLayout(2), form, table))
	myWindow.ShowAndRun()
}

func client(opts base.AppOptions) {
	if opts.Protocol == "udp" {
		// TODO
	} else {
		conf := opts.Clients[0]
		bytes, _ := json.Marshal(conf)
		base.Println(36, 40, fmt.Sprintf("tcp client: %s", string(bytes)))
		cli.TCP(conf.TunnelAddr, conf.LocalPort, opts.Secret)
	}
}
