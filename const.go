package logbus

var (
	THINKINGDATACENTRALIZATION = "tgac" // 所有tga数据集中处理，再根据appid分类
	THINKINGDATA               = "thinkingdata"
	SERVERLOG                  = "server"
	BIGQUERY                   = "bigquery"
	BI                         = "bi"
	Monitor                    = "monitor"
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
