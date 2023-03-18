package main

import (
	"runtime/debug"
	"time"
)

var Build = func() string {
	now := time.Now().Format("20060102150405")
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value[0:8] + "-" + now
			}
		}
	}
	return now
}()
