package main

import (
	"fmt"
	"github.com/dhnt/m3/cmd/systray/icon"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	onExit := func() {
		fmt.Println("Finished onExit")
	}

	// Should be called at the very beginning of main
	systray.Run(onReady, onExit)
}

//
func onReady() {
	//
	systray.SetIcon(icon.Data)
	systray.SetTitle("M3")
	systray.SetTooltip("Inverse of world wide web")

	go func() {

		mHome := systray.AddMenuItem("Home", "M3 home")
		mGit := systray.AddMenuItem("Git", "Git repository")
		mTerm := systray.AddMenuItem("Term", "Terminal")
		mAbout := systray.AddMenuItem("About", "M3 webstie")

		//
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit M3")
		for {
			select {
			case <-mHome.ClickedCh:
				//open.Run(fmt.Sprintf("http://localhost:%v/", internal.GetIntEnv("M3_HOME_PORT", 5001)))
				open.Run("http://home/")

			case <-mGit.ClickedCh:
				// open.Run(fmt.Sprintf("http://localhost:%v/", internal.GetIntEnv("M3_GIT_PORT", 3000)))
				open.Run("http://git.home/")

			case <-mTerm.ClickedCh:
				// open.Run(fmt.Sprintf("http://localhost:%v/", internal.GetIntEnv("M3_TERM_PORT", 50022)))
				open.Run("http://term.home/")

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
