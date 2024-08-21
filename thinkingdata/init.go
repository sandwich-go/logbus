package thinkingdata

import (
	"os"
	"strconv"
	"time"
)

const (
	envTGATimeZoneOffset = "TGA_TIMEZONE_OFFSET"
)

func init() {
	if os.Getenv(envTGATimeZoneOffset) != "" {
		offset, err := strconv.ParseInt(os.Getenv(envTGATimeZoneOffset), 10, 64)
		if err == nil {
			locationTGA = time.FixedZone("", int(offset))
		}
	}
}
