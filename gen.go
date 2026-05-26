package logbus

import (
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sandwich-go/boost/xpanic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type FetchLogContext func() []Field

//go:generate optionGen  --option_return_previous=false
func _ConfOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// log
		"LogLevel":        (zapcore.Level)(zap.DebugLevel), //@MethodComment(日志级别，默认 zap.DebugLevel)
		"Dev":             false,                           // false 输出json格式， true 则输出带颜色的易读log @MethodComment(是否输出带颜色的易读log，默认关闭)
		"DefaultChannel":  string(SERVERLOG),               // 默认的dd_meta_channel @MethodComment(设置默认的dd_meta_channel)
		"DefaultTag":      string(DefaultTag),              // 默认打印的tag @MethodComment(设置默认的tag)
		"CallerSkip":      3,                               // zap logger callerSkip @MethodComment(等于zap.CallerSkip)
		"FetchLogContext": FetchLogContext(nil),            // 打印日志时，获取额外的field @MethodComment(获取额外的field)
		//"LogId":         true,                            // 输出 log id @MethodComment(是否输出log_xid，默认开启) // 日志规范要求必须要有xid 不作为配置放出
		"StackLogLevel": (zapcore.Level)(zap.ErrorLevel), //@MethodComment(打印stack的最低级别，默认ErrorLevel stack if level >= StackLogLevel)
		// stdout
		"BufferedStdout": false, // @MethodComment(输出stdout时使用 logbus.BufferedWriteSyncer)
		// WriteSyncer
		"WriteSyncer":    zapcore.WriteSyncer(os.Stdout), // @MethodComment(输出日志的WriteSyncer，默认为os.Stdout)
		"UseSystemClock": false,                          // @MethodComment(是否使用系统时钟, 默认使用offset时钟)
		// monitor
		"MonitorOutput": MonitorOutput(Noop), // [Logbus, Noop, Prometheus] @MethodComment(监控输出 Logbus, Noop, Prometheus)
		// The Prometheus metrics will be made available on this port: @MethodComment(prometheus监控输出端口，k8s集群保持默认9158端口)
		"DefaultPrometheusListenAddress": ":9158",
		// This is the endpoint where the Prometheus metrics will be made available ("/metrics" is the default with Prometheus):
		"DefaultPrometheusPath": "/metrics", // @MethodComment(prometheus监控输出接口path)
		// DefaultPercentiles is the default spread of percentiles/quantiles we maintain for timings / histogram metrics:
		"DefaultPercentiles":         []float64{0.5, 0.75, 0.99, 1},                                       //@MethodComment(监控统计耗时的分位值，默认统计耗时的 50%, 75%, 99%, 100% 的分位数)
		"DefaultLabel":               prometheus.Labels(map[string]string{}),                              //@MethodComment(监控额外添加的全局label，会在监控指标中显示)
		"MonitorTimingMaxAge":        time.Duration(time.Minute),                                          // @MethodComment(monitor.Timing数据的最大生命周期)
		"EncodeCaller":               zapcore.CallerEncoder(nil),                                          // @MethodComment(指定CallerEncoder)
		"EnableTraceLevel":           true,                                                                // @MethodComment(允许track level log输出)
		"TruncateWriteSyncerOption":  (*TruncateWriteSyncerOption)(newDefaultTruncateWriteSyncerOption()), // @MethodComment(TruncateWriteSyncer的配置项，默认启用当日志超过一定长度时会被截断以避免占满磁盘)
		"DisableTruncateWriteSyncer": false,                                                               // @MethodComment(禁用TruncateWriteSyncer，默认启用，启用后当日志超过一定长度时会被截断以避免占满磁盘)

		// TrackFileDir track 类日志写入的基础目录，TrackOutput 为 FileOnly 或 Both 时必须非空。
		// 文件路径格式（ops 规范）：{TrackFileDir}/{channel}/{yyyymmdd}/{channel}_{hostname}_{slot}.log
		// 其中 hostname 由进程自动从 os.Hostname 获取；文件落盘 channel 段名可通过 TrackChannelAlias 重写
		// （例如 thinkingdata 在 ops 规范中目录与文件前缀都用 tga）。
		// PMT 发布环境（sys_cd_env 非空）下，logbus_track_file_dir 环变会强制覆写本字段。
		"TrackFileDir": string("./tracklogs"), // @MethodComment(track 日志写入的基础目录，TrackOutput 为 FileOnly/Both 时必须配置)
		// TrackChannelAlias 仅用于文件输出：原始 channel 名 → 落盘段名（目录段 + 文件名前缀）映射。例如 thinkingdata→tga。
		// 不改变 dd_meta_channel，也不影响 stdout、业务路由或监控标签。
		// bigquery_xxx 形态的 channel 会按前缀 bigquery 命中 alias 后再拼回表名段。
		"TrackChannelAlias": map[string]string{"thinkingdata": "tga"}, // @MethodComment(文件输出时原始 channel 名到落盘段名映射，默认包含 thinkingdata→tga)
		// TrackFileRotation track 日志文件的时间切割粒度，支持 HourlyRotation（小时）和 MinuteRotation（分钟）
		"TrackFileRotation": TrackRotation(HourlyRotation), // @MethodComment(track 日志时间切割粒度，默认按小时)
		// TrackOutput 控制 track 日志的输出目标：
		//   TrackOutputStdout（默认）— 只写 stdout，与旧行为完全兼容
		//   TrackOutputFile          — 只写文件，需配置 TrackFileDir
		//   TrackOutputBoth          — 同时写 stdout 和文件，需配置 TrackFileDir
		// PMT 发布环境（sys_cd_env 非空）下，logbus_track_output 环变（值 stdout/file/both）会强制覆写本字段。
		"TrackOutput": TrackOutput(TrackOutputStdout), // @MethodComment(track 日志输出目标：TrackOutputStdout/TrackOutputFile/TrackOutputBoth)
		// TrackKeepaliveInterval 心跳间隔；>0 时启用，每个 channel 的 worker 在闲置 interval 内
		// 没有写入则补一条 keepalive 日志。<=0 表示禁用。
		"TrackKeepaliveInterval": time.Duration(time.Minute), // @MethodComment(track 心跳间隔，闲置超过该间隔时补写一条 keepalive 日志，<=0 禁用)
		// TrackKeepaliveMessage 心跳触发时落盘的 msg 内容；可为空字符串。
		"TrackKeepaliveMessage": string(""), // @MethodComment(track 心跳消息内容，默认空)

		// glog
		"PrintAsError":       true, //@MethodComment(glog输出field带error时，将日志级别提升到error)
		"IgnoreLogicalError": true, //@MethodComment(忽略逻辑错误, 逻辑错误使用StringField)
	}
}

func init() {
	InstallConfWatchDog(func(cc *Conf) {
		if cc.DefaultLabel == nil {
			panic("DefaultLabel is nil")
		}
		if cc.MonitorOutput != Prometheus && cc.MonitorOutput != Logbus && cc.MonitorOutput != Noop {
			panic("MonitorOutput not match")
		}
		// PMT 发布环境（sys_cd_env 非空）下：若代码侧已显式进入非 stdout track 输出模式，
		// track 文件输出参数必须由环境变量 logbus_track_file_dir 与 logbus_track_output 提供，
		// 并强制覆写代码侧 Conf。缺失或非法时直接 panic 终止启动，避免静默使用错误的目录或输出模式。
		if os.Getenv(envSysCdEnv) != "" {
			applyPMTTrackEnvOverrides(cc)
		}
		if (cc.TrackOutput == TrackOutputFile || cc.TrackOutput == TrackOutputBoth) && cc.TrackFileDir == "" {
			panic("TrackFileDir must be set when TrackOutput is TrackOutputFile or TrackOutputBoth")
		}
	})
}

// PMT 发布相关环境变量名
const (
	envSysCdEnv                     = "sys_cd_env"                      // PMT 发布环境标识
	envLogbusTrackFileDir           = "logbus_track_file_dir"           // 强制覆写 TrackFileDir
	envLogbusTrackOutput            = "logbus_track_output"             // 强制覆写 TrackOutput，取值 stdout/file/both
	envLogbusTrackKeepaliveMessage  = "logbus_track_keepalive_message"  // 强制覆写 TrackKeepaliveMessage（允许空字符串覆写，需配合"已设置"判断）
	envLogbusTrackKeepaliveInterval = "logbus_track_keepalive_interval" // 强制覆写 TrackKeepaliveInterval，Go time.Duration 语法，0 表示禁用
)

// parseTrackOutputFromEnv 把环境变量字符串映射到 TrackOutput 枚举，大小写不敏感
func parseTrackOutputFromEnv(s string) (TrackOutput, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "stdout":
		return TrackOutputStdout, true
	case "file":
		return TrackOutputFile, true
	case "both":
		return TrackOutputBoth, true
	}
	return 0, false
}

// applyPMTTrackEnvOverrides 在 PMT 环境（sys_cd_env 非空）下强制覆写 track 相关配置。
// 若代码侧 TrackOutput 仍是默认 stdout，表示业务未启用文件输出，此时忽略 PMT track 文件环变。
// 规则：
//  1. logbus_track_output 必填；非法值 panic
//  2. 若 logbus_track_output=stdout，仅覆写 TrackOutput 并直接返回，不要求文件目录
//  3. 若覆写后的 TrackOutput 需要写文件，则 logbus_track_file_dir 必须非空，否则 panic
//  4. logbus_track_keepalive_interval 若设置，按 time.Duration 解析；负值或解析失败 panic；0 视为禁用
//  5. logbus_track_keepalive_message 若设置（包括显式空串），整体覆写 TrackKeepaliveMessage
func applyPMTTrackEnvOverrides(cc *Conf) {
	if cc.TrackOutput == TrackOutputStdout {
		return // 业务明确要求stdout这种情况不开启输出到文件，还是要听从业务逻辑的
	}
	rawOutput := os.Getenv(envLogbusTrackOutput)
	xpanic.WhenTrue(rawOutput == "", envLogbusTrackOutput+" must be set under PMT environment")
	output, ok := parseTrackOutputFromEnv(rawOutput)
	xpanic.WhenFalse(ok, envLogbusTrackOutput+" must be one of stdout/file/both, got: "+rawOutput)
	cc.TrackOutput = output
	if cc.TrackOutput == TrackOutputStdout {
		return // PMT 环境下如果是 stdout 模式， 那么可能集群不支持或者什么原因需要回退配置
	}

	dir := os.Getenv(envLogbusTrackFileDir)
	xpanic.WhenTrue(dir == "", "logbus_track_file_dir must be set under PMT environment")
	cc.TrackFileDir = dir

	// keepalive interval：使用 LookupEnv 区分"未设置"与"设为 0"；后者表示显式禁用
	if rawInterval, present := os.LookupEnv(envLogbusTrackKeepaliveInterval); present {
		d, err := time.ParseDuration(strings.TrimSpace(rawInterval))
		xpanic.WhenErrorAsFmtFirst(err, envLogbusTrackKeepaliveInterval+" must be a valid Go duration (e.g. 1m, 30s), got: "+rawInterval)
		xpanic.WhenTrue(d < 0, envLogbusTrackKeepaliveInterval+" must be >= 0, got: "+rawInterval)
		cc.TrackKeepaliveInterval = d
	}

	// keepalive message：显式覆写（包括空串），未设置时保持代码侧默认
	if rawMsg, present := os.LookupEnv(envLogbusTrackKeepaliveMessage); present {
		cc.TrackKeepaliveMessage = rawMsg
	}
}

//go:generate optionGen  --option_return_previous=false
func TrackLoggerConfOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"tags":    []string(nil), // @MethodComment(打点日志标签，必须提供，至少一个标签)
		"BiAppID": "",            //@MethodComment(bi appid, 比如 "gof.global.prod")
	}
}

//go:generate optionGen  --option_return_previous=false
func TruncateWriteSyncerOptionOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"TruncateMaxSize": int(DefaultTruncateMaxSize),      // @MethodComment(超出长度限制时输出日志的标志字段，便于检索和过滤，默认800KB)
		"MsgPrefixLen":    int(DefaultTruncateMsgPrefixLen), // @MethodComment(截断时保留的前缀长度，默认4KB)
		"MsgSuffixLen":    int(DefaultTruncateMsgSuffixLen), // @MethodComment(截断时保留的后缀长度，默认4KB)
		// StripFields 日志超限时优先尝试剔除的字段路径列表（支持 "a.b.c" 嵌套路径），
		// 剔除后若满足大小限制则直接写出，否则继续走摘要截断流程。
		// 默认剔除 api_call.request_body 和 api_call.response_body。
		"StripFields": []string(DefaultStripFields()), // @MethodComment(日志超限时优先尝试剔除的字段路径列表，支持 a.b.c 嵌套路径)
	}
}
