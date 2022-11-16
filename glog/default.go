package glog

/*const (
	CallDepth = 3
)

type OutputProvider interface {
	Output(calldepth int, s string) error
}

type DefaultLogger struct {
	o OutputProvider
}

type bytesFiledWriter struct {
	bb *bytes.Buffer
}

func (b *bytesFiledWriter) InternalError(key string, err error) {
	if err == nil {
		return
	}
	_, err2 := b.bb.WriteString(fmt.Sprintf("LogbusInternalError_%s: %s", key, err.Error()))
	if err2 != nil {
		//todo : 托底逻辑
		glogInternalError(fmt.Sprintf("[logbus] can not handler error when write field:%s error:%s error2:%s ", key, err.Error(), err2.Error()))
	}
}
func (b *bytesFiledWriter) WriteString(s string, key string) {
	_, err := b.bb.WriteString(s)
	b.InternalError(key, err)
}
func (b *bytesFiledWriter) WriteByteSeparator(c byte) {
	err := b.bb.WriteByte(c)
	b.InternalError("ByteSeparator", err)
}

func (b *bytesFiledWriter) Reset() {
	b.bb.Reset()
}
func (b *bytesFiledWriter) String() string {
	return b.bb.String()
}

func newBytesFiledWriter() *bytesFiledWriter {
	return &bytesFiledWriter{bb: new(bytes.Buffer)}
}

var bufPool = sync.Pool{
	New: func() interface{} {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return newBytesFiledWriter()
	},
}

func getFieldWriter() *bytesFiledWriter {
	b := bufPool.Get().(*bytesFiledWriter)
	b.bb.Reset()
	return b
}

func putFieldWriter(x *bytesFiledWriter) {
	bufPool.Put(x)
}

func NewDefaultLogger(o OutputProvider) *DefaultLogger {
	return &DefaultLogger{o: o}
}

func (d *DefaultLogger) transferField(level string, msg string, v ...Field) string {
	b := getFieldWriter()
	b.WriteString(level, "level")
	b.WriteByteSeparator('\t')
	b.WriteString(fmt.Sprintf("%s: %s\t", "msg", msg), "msg")
	var ret string
	for i := 0; i < len(v); i++ {
		v[i].AddTo(b)
		if i != len(v)-1 {
			b.WriteByteSeparator(' ')
		}
	}
	ret = b.String()
	putFieldWriter(b)
	return ret
}

func (d *DefaultLogger) printAsErr(msg string, fields ...Field) bool {
	hasErr := false
	for _, v := range fields {
		if v.Type == ErrorType {
			hasErr = true
			break
		}
	}
	if hasErr {
		_ = d.o.Output(CallDepth+1, d.transferField("ERROR", msg, fields...))
	}
	return hasErr
}

func (d *DefaultLogger) GDebug(msg string, v ...Field) {
	if d.printAsErr(msg, v...) {
		return
	}
	_ = d.o.Output(CallDepth, d.transferField("DEBUG", msg, v...))
}
func (d *DefaultLogger) GInfo(msg string, v ...Field) {
	if d.printAsErr(msg, v...) {
		return
	}
	_ = d.o.Output(CallDepth, d.transferField("INFO", msg, v...))
}
func (d *DefaultLogger) GWarn(msg string, v ...Field) {
	if d.printAsErr(msg, v...) {
		return
	}
	_ = d.o.Output(CallDepth, d.transferField("WARN", msg, v...))
}
func (d *DefaultLogger) GError(msg string, v ...Field) {
	_ = d.o.Output(CallDepth, d.transferField("ERROR", msg, v...))
}
func (d *DefaultLogger) GPanic(msg string, v ...Field) {
	s := d.transferField("PANIC", msg, v...)
	_ = d.o.Output(CallDepth, s)
	panic(s)
}
func (d *DefaultLogger) GFatal(msg string, v ...Field) {
	_ = d.o.Output(CallDepth, d.transferField("FATAL", msg, v...))
	os.Exit(1)
}

func (d *DefaultLogger) GDebugDepth(depth int, msg string, v ...Field) {
	if d.printAsErr(msg, v...) {
		return
	}
	_ = d.o.Output(CallDepth, d.transferField("DEBUG", msg, v...))
}
func (d *DefaultLogger) GInfoDepth(depth int, msg string, v ...Field) {
	if d.printAsErr(msg, v...) {
		return
	}
	_ = d.o.Output(CallDepth+depth, d.transferField("INFO", msg, v...))
}
func (d *DefaultLogger) GWarnDepth(depth int, msg string, v ...Field) {
	if d.printAsErr(msg, v...) {
		return
	}
	_ = d.o.Output(CallDepth+depth, d.transferField("WARN", msg, v...))
}
func (d *DefaultLogger) GErrorDepth(depth int, msg string, v ...Field) {
	_ = d.o.Output(CallDepth+depth, d.transferField("ERROR", msg, v...))
}
func (d *DefaultLogger) GFatalDepth(depth int, msg string, v ...Field) {
	_ = d.o.Output(CallDepth+depth, d.transferField("FATAL", msg, v...))
	os.Exit(1)
}*/
