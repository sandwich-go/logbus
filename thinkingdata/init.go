package thinkingdata

import (
	"fmt"
	"github.com/sandwich-go/boost/xpanic"
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
		xpanic.WhenTrue(err != nil, fmt.Sprintf("TGA_TIMEZONE_OFFSET %s is not number", os.Getenv(envTGATimeZoneOffset)))
		xpanic.WhenTrue(offset > 12*60*60 || offset < -12*60*60, fmt.Sprintf("TGA_TIMEZONE_OFFSET %d is out of range", offset))
		locationTGA = time.FixedZone("", int(offset))
	}
}
