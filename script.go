package main

//#cgo CFLAGS: -x objective-c
//#cgo LDFLAGS: -framework Foundation
//#include "script.h"
import "C"
import "errors"

const (
	// эплскрипт кринж, оно запускает прилу при отправке/получении данных, и это никак нельзя изменить
	source = `
if application "Music" is running then
	tell application "Music"
		set a to ""
		set n to ""
		set p to 0
		set d to 0
			try
				set tr to current track
				set a to artist of tr
				set n to name of tr
				set p to player position
				set d to duration of tr
			end try
		return {player state as string, a, n, p, d}
	end tell
end if
`
	StatePlaying = "playing"
	StatePaused  = "paused"
	StateStopped = "stopped"
)

var script *C.NSAppleScript

func init() {
	if script = C.compileScript(C.CString(source)); script == nil {
		panic("failed to compile script")
	}
}

type result struct {
	state    string
	artist   string
	name     string
	position float64
	duration float64
}

func executeScript() (*result, error) {
	descriptor := C.executeScript(script)
	if descriptor == nil {
		return nil, errors.New("failed to execute script")
	}
	return &result{
		state:    getStringFromDescriptor(descriptor, 1),
		artist:   getStringFromDescriptor(descriptor, 2),
		name:     getStringFromDescriptor(descriptor, 3),
		position: getFloatFromDescriptor(descriptor, 4),
		duration: getFloatFromDescriptor(descriptor, 5),
	}, nil
}

func getStringFromDescriptor(descriptor *C.NSAppleEventDescriptor, index int) string {
	return C.GoString(C.getStringFromDescriptor(descriptor, C.int(index)))
}

func getIntFromDescriptor(descriptor *C.NSAppleEventDescriptor, index int) int {
	return int(C.getIntFromDescriptor(descriptor, C.int(index)))
}

func getFloatFromDescriptor(descriptor *C.NSAppleEventDescriptor, index int) float64 {
	return float64(C.getFloatFromDescriptor(descriptor, C.int(index)))
}
