package bigquery

import (
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegexp(t *testing.T) {
	Convey("pattern", t, func() {
		a := KeyPatternReg.MatchString("aaaa1")
		So(a, ShouldBeTrue)
		b := KeyPatternReg.MatchString("#bbb2")
		So(b, ShouldBeFalse)
		c := KeyPatternReg.MatchString("$cccc_5")
		So(c, ShouldBeTrue)
	})
	Convey("column", t, func() {
		var columnPattern, _ = regexp.Compile("^[$]")
		a := columnPattern.MatchString("aaaa1")
		So(a, ShouldBeFalse)
		b := columnPattern.MatchString("#bbb2")
		So(b, ShouldBeFalse)
		c := columnPattern.MatchString("$cccc_5")
		So(c, ShouldBeTrue)
	})
}
