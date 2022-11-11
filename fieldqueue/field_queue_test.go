package fieldqueue

import (
	"github.com/sandwich-go/logbus"
	"testing"

	"go.uber.org/zap"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFieldQueue(t *testing.T) {
	Convey("test multi NewQueue\n", t, func() {
		q := NewQueue()
		So(q, ShouldNotBeNil)
		for i := 0; i < 10; i++ {
			q.Push(zap.Int("1", i))
		}
		data := q.Retrieve()
		So(len(data), ShouldEqual, 10)
		q = NewQueue()
		So(q, ShouldNotBeNil)
		q.Push(zap.Int("1", 11))
		data = q.Retrieve()
		So(len(data), ShouldEqual, 1)
		So(data[0].Integer, ShouldEqual, 11)
		logbus.Debug(q.Retrieve()...)
	})
}
