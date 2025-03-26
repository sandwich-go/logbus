package thinkingdata

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sandwich-go/boost/xpanic"
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
