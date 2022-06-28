package main

import (
	"fmt"
	"time"

	_ "embed"

	"github.com/getlantern/systray"
	"github.com/hugolgst/rich-go/client"
)

const AppID = "991335093878673448"

//go:embed icons/tray.png
var icon []byte

func main() {
	systray.Run(onReady, nil)
}

func onReady() {
	systray.SetIcon(icon)

	state := systray.AddMenuItem("No state", "")
	state.Disable()

	systray.AddSeparator()
	quit := systray.AddMenuItem("Quit", "")

	// updater
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			result, err := executeScript()
			if err != nil {
				state.SetTitle("Script error")
				client.Logout()
				continue
			}

			if result.state == "" {
				state.SetTitle("Apple Music is not running")
				client.Logout()
				continue
			}

			if result.name == "" || result.artist == "" {
				state.SetTitle("No song")
				client.Logout()
				continue
			}

			if err = client.Login(AppID); err != nil {
				state.SetTitle("Discord RPC error")
				continue
			}

			song := result.artist + " – " + result.name
			state.SetTitle(fmt.Sprintf("%s (%d:%02d / %d:%02d)",
				song, int(result.position/60), int(result.position)%60,
				int(result.duration/60), int(result.duration)%60))

			activity := client.Activity{
				LargeImage: "music",
				LargeText:  "Apple Music logo lol",
				SmallImage: "pause",
				SmallText:  "Paused",
				Details:    song,
			}
			if result.state == StatePlaying {
				activity.SmallImage = "play"
				activity.SmallText = "Playing"

				// start := time.Now().Add(-time.Duration(info.Position * float64(time.Second)))
				end := time.Now().
					Add(time.Duration(result.duration * float64(time.Second))).
					Add(-time.Duration(result.position * float64(time.Second)))
				activity.Timestamps = &client.Timestamps{
					Start: &time.Time{},
					End:   &end,
				}
			}

			client.SetActivity(activity)
		}
	}()

	for {
		select {
		case <-quit.ClickedCh:
			systray.Quit()
		}
	}
}
