package logbus

var (
	THINKINGDATA = "thinkingdata"
	SERVERLOG    = "server"
	BIGQUERY     = "bigquery"
	BI           = "bi"
	Monitor      = "monitor"
)

var (
	Meta     = "dd_meta_channel"
	MsgBody  = "msg"
	LevelKey = "log_level"
	TimeKey  = "date"
	Tags     = "tags"
)

type MonitorOutput int

const (
	Prometheus MonitorOutput = iota
	Logbus
	Noop
)

var (
	DefaultTag = "logbus"
	MonitorTag = "monitor"
	LogId      = "log_xid"
)
