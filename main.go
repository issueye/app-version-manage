package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"

	bootstrap "github.com/issueye/app-version-manage/bootstrap"
)

// Constants
const htmlAbout = `Welcome on <b>Astilectron</b> demo!<br>
This is using the bootstrap and the bundler.`

// Vars injected via ldflags by bundler
var (
	AppName            string = "程序版本管理系统"
	BuiltAt            string
	VersionAstilectron string
	VersionElectron    string
)

// Application Vars
var (
	fs    = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	debug = fs.Bool("d", true, "enables the debug mode")
	w     *astilectron.Window
)

func main() {
	// Create logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	// Parse flags
	fs.Parse(os.Args[1:])

	// Run bootstrap
	l.Printf("Running app built at %s\n", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
			SingleInstance:     true,
			VersionAstilectron: VersionAstilectron,
			VersionElectron:    VersionElectron,
		},
		Debug:  *debug,
		Logger: l,
		MenuOptions: []*astilectron.MenuItemOptions{{
			Label: astikit.StrPtr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("About"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						if err := bootstrap.SendMessage(w, "about", htmlAbout, func(m *bootstrap.MessageIn) {
							// Unmarshal payload
							var s string
							if err := json.Unmarshal(m.Payload, &s); err != nil {
								l.Println(fmt.Errorf("unmarshaling payload failed: %w", err))
								return
							}
							l.Printf("About modal has been displayed and payload is %s!\n", s)
						}); err != nil {
							l.Println(fmt.Errorf("sending about event failed: %w", err))
						}
						return
					},
				},
				{Role: astilectron.MenuItemRoleClose},
			},
		}},
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]

			go func() {
				time.Sleep(5 * time.Second)
				w.ExecuteJavaScript(`console.log('测试')`)
				w.SendMessage(`{code: 200}`, func(m *astilectron.EventMessage) {
					var s string
					m.UnmarshalJSON([]byte(s))
					fmt.Println("s", s)
				})
				if err := bootstrap.SendMessage(w, "check.out.menu", "Don't forget to check out the menu!"); err != nil {
					l.Println(fmt.Errorf("sending check.out.menu event failed: %w", err))
				}
			}()
			return nil
		},
		Windows: []*bootstrap.Window{
			{
				Homepage:       "http://localhost:8080/",
				MessageHandler: handleMessages,
				Options: &astilectron.WindowOptions{
					// 打开时居中
					Center: astikit.BoolPtr(true),
					// 禁止改变窗口大小
					Resizable: astikit.BoolPtr(true),
					// 窗口大小
					Height: astikit.IntPtr(720),
					Width:  astikit.IntPtr(1280),
					// 无边框窗口
					Frame:           astikit.BoolPtr(false),
					Transparent:     astikit.BoolPtr(true),
					BackgroundColor: astikit.StrPtr("#00000000"),

					WebPreferences: &astilectron.WebPreferences{
						NodeIntegrationInWorker: astikit.BoolPtr(false),
					},
				},
			},
		},
	}); err != nil {
		l.Fatal(fmt.Errorf("running bootstrap failed: %w", err))
	}
}
