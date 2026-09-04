# ThinkingData 通用 Track 基线

这组基准固定覆盖业务当前使用的 `logbus.Tracker(logbus.THINKINGDATA).Track(fields...)`。
它用于比较后续优化，不能代表 Fluent Bit、Logserver 或网络接收端的时延。

## 输入定义

|名称|含义|
|---|---|
|常规属性|6 个 ThinkingData 协议字段（账号、UUID、类型、事件、事件 ID、AppID）和 3 个业务属性（区服、结果、分数），共 9 个字段。|
|属性较多|常规属性外，再增加 32 个业务属性，每个值为 256 B；额外源数据约 8 KiB。|
|长文本属性|常规属性外，再增加 1 个 4 KiB 业务文本属性。|

“常规属性”“属性较多”“长文本属性”描述本基准的输入形态。

## 边界

- 包含：`Track(fields...)` 内的字段转换、ThinkingData 提取和校验、内层 ThinkingData JSON、外层 Logbus JSON、默认未缓冲写入。
- 不包含：业务侧拼接 `fields`、生成 UUID、Fluent Bit / Logserver / 网络、磁盘，以及下游消费。
- 输出写入本地 Pipe 并持续读取，避免终端吞吐成为变量；仍保留未缓冲 WriteSyncer 的调用方式。
- `Discard` 项移除操作系统写入，只用作进程内开销下界，不能替代默认输出链路的结果。
- 基准启动时会关闭进程初始化阶段可能遗留的 PMT 动态日志级别轮询，避免后台文件监听进入采样。
- 并发项由 `RunParallel` 运行，使用 `-cpu=8`；它测单进程的共享 logger 路径，不代表线上容器的总吞吐。

## 复现命令

```sh
rtk proxy go test . -run '^TestThinkingDataBaselineOutput$' -v
rtk proxy go test . -run '^$' -bench '^BenchmarkThinkingDataBaseline(TrackerCreation|GenericTrack(Conventional|AttributeRich|LongText)(Discard)?)$' -benchmem -benchtime=1s -count=5 -cpu=1
rtk proxy go test . -run '^$' -bench '^BenchmarkThinkingDataBaselineGenericTrack(Conventional|AttributeRich)Parallel$' -benchmem -benchtime=2s -count=5 -cpu=8
```

使用中位数作优化前后对比，并同时比较 `ns/op`、`B/op`、`allocs/op` 和测试输出字节数。优化后必须通过 `TestThinkingDataBaselineOutput`，以确认外层日志和内层 ThinkingData JSON 仍然有效。

## 本次结果

采集时间：2026-09-01（Asia/Shanghai）<br>
代码版本：`89224ab`<br>
环境：Apple M4 Pro，`go1.25.3 darwin/arm64`

单线程命令使用 `-cpu=1 -benchtime=1s -count=5`。下表保留全部原始耗时；对比时使用中位数。

|项目|5 次 `ns/op`|中位数|`B/op`|`allocs/op`|
|---|---:|---:|---:|---:|
|创建 `Tracker(THINKINGDATA)`|61.60, 61.20, 63.69, 63.67, 63.98|63.67 ns|104|4|
|常规属性，默认未缓冲 Pipe 输出|5,987, 6,092, 6,904, 6,111, 6,242|6,111 ns|3,816|32|
|常规属性，Discard（进程内下界）|6,000, 6,072, 6,776, 6,084, 6,003|6,072 ns|3,816|32|
|属性较多，默认未缓冲 Pipe 输出|85,304, 100,467, 95,145, 85,151, 85,959|85,959 ns|53,942|140|
|属性较多，Discard（进程内下界）|113,923, 91,675, 97,273, 105,837, 93,501|97,273 ns|53,943|140|
|长文本属性，默认未缓冲 Pipe 输出|40,101, 39,819, 40,725, 40,149, 40,325|40,149 ns|20,252|40|

Pipe 和 Discard 的差异处于这台开发机的微基准波动范围，不能据此估算生产输出端的耗时；两者保留的意义是让后续改动可以分别复测“完整写出”和“仅进程内处理”。

并发命令使用 `-cpu=8 -benchtime=2s -count=5`。`ns/op` 是该进程的总吞吐折算值，不是单个请求的延迟。

|项目|5 次 `ns/op`|中位数|总吞吐（约）|`B/op` / `allocs/op`|
|---|---:|---:|---:|---:|
|常规属性，并发 8 路|4,239, 4,309, 4,365, 4,350, 4,299|4,309 ns|232,000 条/秒|3,826 / 32|
|属性较多，并发 8 路|32,481, 32,796, 33,979, 34,304, 34,583|33,979 ns|29,400 条/秒|54,586 / 141|

单线程中位数换算后，常规属性约 163,600 条/秒、属性较多约 11,600 条/秒。因此 8 路并发相对单路的总吞吐增幅约为 1.4 倍和 2.5 倍，未达到线性扩展。

输出完整性测试得到的单条字节数如下：

|输入|外层 Logbus JSON|内层 ThinkingData JSON|
|---|---:|---:|
|常规属性|595 B|269 B|
|属性较多|9,491 B|9,037 B|
|长文本属性|4,713 B|4,383 B|

### 采样定位

- 默认未缓冲 Pipe 路径的 4 秒 CPU 采样中，`syscall.syscall` 占 95.4% 的样本，调用链来自 `TruncateWriteSyncer.Write` 到 Pipe 写入。这只说明一次日志写出的系统调用是独立的大头；线上 stdout、Fluent Bit 和接收端必须另行压测。
- 去除操作系统写入后的 4 秒分配采样中，Zap buffer 约占 27.4%，`thinkingdata.Data.MarshalAsJsonV2` 约占 19.6%，`MapObjectEncoder.AddInt64` 约占 16.1%，`MapObjectEncoder.AddString` 约占 10.7%，外层 `glsEncoder.EncodeEntry` 直接占 10.7%（累计 38.7%）。
- `NewMapObjectEncoder` 本身在该样本中约占 1.8%；更大的 MapObjectEncoder 成本来自逐字段写入及后续对 map 的编码，而不只是创建 map。

后续优化应在同一输入、同一命令下比较中位数；同时保持输出字节数合理，并通过 `TestThinkingDataBaselineOutput`。
