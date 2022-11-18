package logbus

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// https://gist.github.com/dtjm/c6ebc86abe7515c988ec
// go test -v -run=BENCH -bench=. -benchtime 5s -benchmem

var (
	testData = []string{"a", "b", "c", "d", "e"}
)

func BenchmarkJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := strings.Join(testData, ":")
		_ = s
	}
}

func BenchmarkSprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := fmt.Sprintf("%s:%s:%s:%s:%s", testData[0], testData[1], testData[2], testData[3], testData[4])
		_ = s
	}
}

func BenchmarkConcat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := testData[0] + ":"
		s += testData[1] + ":"
		s += testData[2] + ":"
		s += testData[3] + ":"
		s += testData[4]
		_ = s
	}
}

func BenchmarkConcatOneLine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := testData[0] + ":" +
			testData[1] + ":" +
			testData[2] + ":" +
			testData[3] + ":" +
			testData[4]
		_ = s
	}
}

func BenchmarkBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var b bytes.Buffer
		b.WriteString(testData[0])
		b.WriteByte(':')
		b.WriteString(testData[1])
		b.WriteByte(':')
		b.WriteString(testData[2])
		b.WriteByte(':')
		b.WriteString(testData[3])
		b.WriteByte(':')
		b.WriteString(testData[4])
		s := b.String()
		_ = s
	}
}

func BenchmarkBufferWithReset(b *testing.B) {
	var buf bytes.Buffer

	for i := 0; i < b.N; i++ {
		buf.Reset()

		buf.WriteString(testData[0])
		buf.WriteByte(':')
		buf.WriteString(testData[1])
		buf.WriteByte(':')
		buf.WriteString(testData[2])
		buf.WriteByte(':')
		buf.WriteString(testData[3])
		buf.WriteByte(':')
		buf.WriteString(testData[4])
		s := buf.String()
		_ = s
	}
}

func BenchmarkBufferFprintf(b *testing.B) {
	buf := &bytes.Buffer{}

	for i := 0; i < b.N; i++ {
		buf.Reset()

		fmt.Fprintf(buf, "%s:%s:%s:%s:%s", testData[0], testData[1], testData[2], testData[3], testData[4])
		s := buf.String()
		_ = s
	}

}

func BenchmarkBufferStringBuilder(b *testing.B) {
	var buf strings.Builder

	for i := 0; i < b.N; i++ {
		buf.Reset()

		buf.WriteString(testData[0])
		buf.WriteByte(':')
		buf.WriteString(testData[1])
		buf.WriteByte(':')
		buf.WriteString(testData[2])
		buf.WriteByte(':')
		buf.WriteString(testData[3])
		buf.WriteByte(':')
		buf.WriteString(testData[4])
		s := buf.String()
		_ = s
	}
}

func BenchmarkDebug(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Debug("", String("a", "b"), String("C", "D"))
	}
}
