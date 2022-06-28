package main

import (
	"fmt"
	"net/url"
	"time"

	_ "embed"

	"github.com/Code-Hex/go-generics-cache/policy/fifo"
	"github.com/getlantern/systray"
	"github.com/hugolgst/rich-go/client"
)

const AppID = "991335093878673448"

//go:embed assets/tray.png
var icon []byte

func main() {
	systray.Run(onReady, nil)
}

func onReady() {
	systray.SetIcon(icon)

	state := systray.AddMenuItem("No state", "")
	state.Disable()

	systray.AddSeparator()
	restart := systray.AddMenuItem("Restart Discord RPC", "")
	quit := systray.AddMenuItem("Quit", "")

	go func() {
		cache := fifo.NewCache[string, string](fifo.WithCapacity(100))
		ticker := time.NewTicker(5 * time.Second)
		// hacky trick to force first tick
		for ; true; <-ticker.C {
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
				SmallImage: "pause",
				Details:    song,
				Buttons: []*client.Button{
					{
						Label: "Search on YouTube",
						Url: "https://www.youtube.com/results?" + url.Values{
							"search_query": {song},
						}.Encode(),
					},
				},
			}
			if result.state == StatePlaying {
				activity.SmallImage = "play"

				end := time.Now().
					Add(time.Duration(result.duration * float64(time.Second))).
					Add(-time.Duration(result.position * float64(time.Second)))
				activity.Timestamps = &client.Timestamps{
					Start: &time.Time{},
					End:   &end,
				}
			}

			artwork, _ := cache.Get(song)
			if artwork == "" {
				if artwork = artworkSearch(result.artist + " " + result.name); artwork != "" {
					cache.Set(song, artwork)
				}
			}
			if artwork != "" {
				activity.LargeImage = artwork
			}

			client.SetActivity(activity)
		}
	}()

	for {
		select {
		case <-restart.ClickedCh:
			client.Logout()
		case <-quit.ClickedCh:
			systray.Quit()
		}
	}
}
