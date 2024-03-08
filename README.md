# logbus
a simple logger dedicated to JSON output

### 背景
我们在开发环境下的log一般输出到本地文件，或者stdout（再由supervisor写文件）。 在线上环境将log统一打到log服务器（logserver）上，然后运维解析tag分发到 s3 / elasticsearch/ thinkingdata 等。 从服务器输出到logserver有多种方法：

- 直接输出到本地fluentd agent的socket中
- 若机器部署在k8s上，则可以输出到stdout，然后由fluentd-bit采集文件和分发
- 若机器部署在ec2上，则可以输出到stdout，然后由supervisor重定向到文件，运维对文件进行采集

### 基础用法

```go
import "github.com/sandwich-go/logbus"
import "gopkg.in/natefinch/lumberjack.v2"

func MustInstallLogger() {
    defer logbus.Close()
	//初始化
    logbus.Init(logbus.NewConf(
        logbus.WithDev(true),  //是否开发模式 false会输出json格式 ture的话会打印代颜色的易读log
        logbus.WithLogLevel(lv), // 对应zapcore.Level的枚举值
        logbus.WithWriteSyncer(zapcore.AddSync(&lumberjack.Logger{  //可以自定义日志的输出方式 默认的话是os.stdout
            Filename:   "./app.log", // 日志文件路径
            MaxSize:    10,          // 单个日志文件的最大大小，单位 MB
            MaxBackups: 10,          // 保留的旧日志文件数量
            MaxAge:     14,          // 保留的旧日志文件最大天数
            Compress:   false,       // 是否压缩旧日志文件
        })))
    
	//设置全局字段
    logbus.SetGlobalFields([]logbus.Field{
        logbus.String("host_name", hostName),
        logbus.String("server_tag", tag),
        logbus.String("server_id", common.LocalNodeID),
        logbus.String("server_ip", xip.GetLocalIP()),
        logbus.Int64("server_birth", time.Now().Unix()),
    })
	
	// 基础使用
    logbus.Error("test log", logbus.String("key", "value"), ErrorField(errors.New("rrrrrr")))
    logbus.Info("test log", logbus.Int("key1", 123), zap.Bool("key2", true))

    // scope logger
    playerLogger := logbus.NewScopeLogger("Player", zap.String("playername", "zhangsong"), zap.Int("playerid", 123))
    guildLogger := logbus.NewScopeLogger("Guild", zap.String("guildname", "guild1"))
    playerLogger.Info("player gold", logbus.Int("money", 648))
    guildLogger.Info("guild gold", logbus.Int("money", 6480))
    
	// queue
    q := logbus.NewQueue()
    q.Push(zap.Int("i", 1))
    q.Push(zap.Int("j", 2))
    logbus.Debug("msg", q.Retrieve()...)
}

```


### 配置说明

```go
// WithLogLevel 日志级别，默认 zap.DebugLevel
func WithLogLevel(v zapcore.Level) ConfOption

// WithDev 是否输出带颜色的易读log，默认关闭
func WithDev(v bool) ConfOption

// WithDefaultChannel 设置默认的dd_meta_channel
func WithDefaultChannel(v string) ConfOption

// WithDefaultTag 设置默认的tag
func WithDefaultTag(v string) ConfOption

// WithCallerSkip 等于zap.CallerSkip
func WithCallerSkip(v int) ConfOption

// WithStackLogLevel 是否输出log_xid，默认开启,打印stack的最低级别，默认ErrorLevel stack if level >= StackLogLevel
func WithStackLogLevel(v zapcore.Level) ConfOption

// WithBufferedStdout 输出stdout时使用 logbus.BufferedWriteSyncer
func WithBufferedStdout(v bool) ConfOption

// WithMonitorOutput 监控输出 Logbus, Noop, Prometheus
func WithMonitorOutput(v MonitorOutput) ConfOption

// WithDefaultPrometheusListenAddress prometheus监控输出端口，k8s集群保持默认9158端口
func WithDefaultPrometheusListenAddress(v string)

// WithDefaultPrometheusPath prometheus监控输出接口path
func WithDefaultPrometheusPath(v string) ConfOption

// WithDefaultPercentiles 监控统计耗时的分位值，默认统计耗时的 50%, 75%, 99%, 100% 的分位数
func WithDefaultPercentiles(v ...float64) ConfOption

// WithDefaultLabel 监控额外添加的全局label，会在监控指标中显示
func WithDefaultLabel(v prometheus.Labels) ConfOption

// WithMonitorTimingMaxAge monitor.Timing数据的最大生命周期
func WithMonitorTimingMaxAge(v time.Duration) ConfOption

// WithPrintAsError glog输出field带error时，将日志级别提升到error
func WithPrintAsError(v bool) ConfOption

// WithWriteSyncer 输出日志的WriteSyncer，默认为os.Stdout
func WithWriteSyncer(v zapcore.WriteSyncer) ConfOption
```

