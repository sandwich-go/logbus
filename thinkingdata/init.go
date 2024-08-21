package thinkingdata

import (
	"os"
	"strconv"
	"time"
)

const (
	envTGATimeZoneOffset = "TGA_TIMEZONE_OFFSET"
)

func Init() {
	if os.Getenv(envTGATimeZoneOffset) != "" {
		offset, ok := strconv.ParseInt(os.Getenv(envTGATimeZoneOffset), 10, 64)
		if ok == nil {
			locationTGA = time.FixedZone("", int(offset))
		}
	}
}
