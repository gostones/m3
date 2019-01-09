package main

import (
	"fmt"
	"github.com/dhnt/m3/cmd/systray/icon"
	m3 "github.com/dhnt/m3/internal"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"os"
	"strconv"
)

func main() {
	onExit := func() {
		//fmt.Println("Starting onExit")
		//now := time.Now()
		//ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
		//fmt.Println("Finished onExit")
		m3.StopGPM()
	}

	go m3.StartGPM()

	// Should be called at the very beginning of main().
	systray.Run(onReady, onExit)
}

func getPort(env string, port int) int {
	if p := os.Getenv(env); p != "" {
		if port, err := strconv.Atoi(p); err == nil {
			return port
		}
	}
	return port
}

//
func onReady() {
	//
	systray.SetIcon(icon.Data)
	systray.SetTitle("M3")
	systray.SetTooltip("Inverse of world wide web")

	//mQuitOrig := systray.AddMenuItem("Quit", "Quit M3")
	// go func() {
	// 	<-mQuitOrig.ClickedCh
	// 	fmt.Println("Requesting quit")
	// 	systray.Quit()
	// 	fmt.Println("Finished quitting")
	// }()

	// We can manipulate the systray in other goroutines
	go func() {
		// systray.SetIcon(icon.Data)
		// systray.SetTitle("M3")
		// systray.SetTooltip("Inverse of world wide web")

		//mChange := systray.AddMenuItem("Start", "Start/stop proxy")

		//mChecked := systray.AddMenuItem("Unchecked", "Check Me")
		//mEnabled := systray.AddMenuItem("Enabled", "Enabled")
		//systray.AddMenuItem("Ignored", "Ignored")
		mHome := systray.AddMenuItem("Home", "M3 home")

		mGit := systray.AddMenuItem("Git", "Git repository")

		// Sets the icon of a menu item. Only available on Mac.
		// mQuit.SetIcon(icon.Data)
		mAbout := systray.AddMenuItem("About", "M3 webstie")

		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit M3")

		//shown := true
		//running := true
		for {
			select {
			//case <-mChange.ClickedCh:
			// if running {
			// 	mChange.SetTitle("Stop")
			// 	running = false
			// } else {
			// 	mChange.SetTitle("Start")
			// 	running = true
			// }
			// case <-mChecked.ClickedCh:
			// 	if mChecked.Checked() {
			// 		mChecked.Uncheck()
			// 		mChecked.SetTitle("Unchecked")
			// 	} else {
			// 		mChecked.Check()
			// 		mChecked.SetTitle("Checked")
			// 	}
			// case <-mEnabled.ClickedCh:
			// 	mEnabled.SetTitle("Disabled")
			// 	mEnabled.Disable()
			case <-mHome.ClickedCh:
				open.Run(fmt.Sprintf("http://localhost:%v/", getPort("M3_HOME_PORT", 8080)))

			case <-mGit.ClickedCh:
				open.Run(fmt.Sprintf("http://localhost:%v/", getPort("M3_GIT_PORT", 3000)))
			// case <-mToggle.ClickedCh:
			// 	if shown {
			// 		mQuitOrig.Hide()
			// 		mEnabled.Hide()
			// 		shown = false
			// 	} else {
			// 		mQuitOrig.Show()
			// 		mEnabled.Show()
			// 		shown = true
			// 	}
			case <-mAbout.ClickedCh:
				open.Run("https://github.com/dhnt/m3")
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("Quit now...")
				return
			}
		}
	}()
}
