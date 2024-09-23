package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"kotunnel/base"
	"kotunnel/cli"
	"os"
)

func main() {

	started := false

	myApp := app.New()
	myWindow := myApp.NewWindow("KoTunnel GUI")
	myWindow.Resize(fyne.Size{
		Width: 400,
	})

	var args []string
	cache, _ := os.ReadFile("./gui_cache")
	_ = json.Unmarshal(cache, &args)

	protocol := widget.NewEntry()
	protocol.SetPlaceHolder("tcp")
	if len(args) > 0 {
		protocol.SetText(args[0])
	}

	secret := widget.NewEntry()
	secret.SetPlaceHolder("123456")
	if len(args) > 1 {
		secret.SetText(args[1])
	}

	tunnelAddr := widget.NewEntry()
	tunnelAddr.SetPlaceHolder("127.0.0.1:8080")
	if len(args) > 2 {
		tunnelAddr.SetText(args[2])
	}

	localPort := widget.NewEntry()
	localPort.SetPlaceHolder("9090")
	if len(args) > 3 {
		localPort.SetText(args[3])
	}

	idleNum := widget.NewEntry()
	idleNum.SetPlaceHolder("5")
	if len(args) > 4 {
		idleNum.SetText(args[4])
	}

	tip := widget.NewLabel("")

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Protocol", Widget: protocol},
			{Text: "Secret", Widget: secret},
			{Text: "Tunnel addr", Widget: tunnelAddr},
			{Text: "Local port", Widget: localPort},
			{Text: "Idle num", Widget: idleNum},
			{Text: "", Widget: tip},
		},
		OnSubmit: func() { // optional, handle form submission
			if started {
				tip.SetText("It's already started, cannot be restarted.")
				return
			}
			bytes, _ := json.Marshal([]string{protocol.Text, secret.Text, tunnelAddr.Text, localPort.Text, idleNum.Text})
			_ = os.WriteFile("./gui_cache", bytes, 0644)
			// 配置加载
			base.InitConfig([]string{"./main", "client", protocol.Text, secret.Text, tunnelAddr.Text, localPort.Text, idleNum.Text})
			marshal, _ := json.Marshal(base.Config().App)
			// 日志加载
			base.InitLog()
			base.Println(33, 40, "config: "+string(marshal))
			started = true
			tip.SetText("It's already started.")
			go client(base.Config().App)
			// myWindow.Close()
		},
		SubmitText: "Run",
	}

	myWindow.SetContent(form)
	myWindow.ShowAndRun()
}

func client(opts base.AppOptions) {
	if opts.Protocol == "udp" {
		// TODO
	} else {
		for _, v := range opts.Clients {
			bytes, _ := json.Marshal(v)
			base.Println(36, 40, fmt.Sprintf("tcp client: %s", string(bytes)))
			for i := 0; i < v.IdleNum-1; i++ {
				go cli.TCP(v.TunnelAddr, v.LocalPort, opts.Secret)
			}
			cli.TCP(v.TunnelAddr, v.LocalPort, opts.Secret)
		}
	}
}
