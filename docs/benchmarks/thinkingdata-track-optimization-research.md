# ThinkingData 通用 Track 优化调研

> 调研日期：2026-09-03
>
> 基线版本：`89224ab`；P0 修复提交：`cc47e5b`
>
> 适用对象：业务通过 `logbus.Tracker(logbus.THINKINGDATA).Track(fields...)` 写 ThinkingData 日志的路径

## 结论

当前瓶颈不是一个点，而是“字段转换 → ThinkingData JSON → Logbus 外层 JSON → 未缓冲写出”的串行链路。

1. **ThinkingData 写完当前日志后会归还内部临时内存。** 其他模块通过 `Zap2Json` 取得日志内容时，会拿到独立副本，因此函数返回后仍可继续保存和使用。三组基准的分配字节分别降低 27.7%、85.5% 和 86.1%。
2. **先控制日志量和属性量，再优化其余库内代码。** 属性较多的输入比常规属性慢约 14 倍、分配约 14 倍；减少可选属性或对非关键事件采样，会按比例降低成本。
3. **默认未缓冲写出是最先值得在线上试验的配置项。** 本地 Pipe 采样中，系统写入调用占 95.4% CPU 样本。现有 `WithBufferedStdout(true)` 可作为低风险验证入口，但必须先确认可接受的日志延迟和进程异常时的丢失边界。
4. **通用 Track 的标量字段已不再构造完整中间 map。** `ExtractFields` 直接拆分协议字段和业务属性；复杂 `zap.Field` 自动回退旧路径。高频、字段固定的事件仍可迁移到 `TrackWithTGAData`，跳过整个通用适配层。
5. **不要优先做 MapObjectEncoder 对象池。** `Data.Properties` 持有其 map，直接复用会造成并发数据串写；复制后再复用也可能抵消收益。
6. **双 JSON 和字段校验去重有优化空间，但属于协议/语义敏感改动。** 应排在事件/属性治理、写出缓冲和类型化路径验证之后。

## 基线事实

完整原始数据、复现命令和输出校验见 [ThinkingData 通用 Track 基线](thinkingdata-generic-track-baseline.md)。以下为 5 次运行中位数，环境为 Apple M4 Pro、Go 1.25.3、单线程。

|输入形态|耗时|分配字节|分配次数|说明|
|---|---:|---:|---:|---|
|常规属性|6.11 µs/条|3,816 B|32|6 个 ThinkingData 协议字段 + 3 个业务属性|
|属性较多|85.96 µs/条|53,942 B|140|常规属性外增加 32 个 256 B 业务属性|
|长文本属性|40.15 µs/条|20,252 B|40|常规属性外增加 1 个 4 KiB 业务属性|
|创建 Tracker|63.67 ns/次|104 B|4|仅创建 `Tracker(THINKINGDATA)`|

属性较多相对常规属性：耗时约 **14.1 倍**、分配字节约 **14.1 倍**、分配次数约 **4.4 倍**。这直接支持“减少非必要属性”优先于细节微调。

8 路并发的总吞吐只达到单路的约 1.4 倍（常规属性）和 2.5 倍（属性较多），说明单进程加并发并不能线性解决压力问题。

## P0 修复结果

修复提交为 `cc47e5b fix(thinkingdata): release zap JSON buffers`。在两个干净工作目录分别检出 `89224ab` 与 `cc47e5b`，使用完全相同的输入和 Pipe 输出路径运行：

```bash
go test . -run '^$' -bench '^BenchmarkThinkingDataBaselineGenericTrack(Conventional|AttributeRich|LongText)$' -benchmem -benchtime=1s -count=5 -cpu=1
```

环境为 Apple M4 Pro、Go 1.25.3。下表为 5 次运行中位数；耗时变化也单独列出，避免与分配数据混在一起。

|输入形态|修复前中位数|修复后中位数|耗时变化|分配字节变化|分配次数变化|
|---|---:|---:|---:|---:|---:|
|常规属性|6.036 µs · 3,816 B · 32 次|5.840 µs · 2,760 B · 30 次|-3.2%|-27.7%|-6.3%|
|属性较多|89.565 µs · 53,942 B · 140 次|82.087 µs · 7,824 B · 130 次|-8.3%|-85.5%|-7.1%|
|长文本属性|41.381 µs · 20,252 B · 40 次|39.042 µs · 2,809 B · 33 次|-5.7%|-86.1%|-17.5%|

5 次 `ns/op` 原始结果如下。常规属性的 5 次修复后结果均低于修复前；属性较多与长文本各有一次 Pipe 写出抖动，因此以中位数作比较。此次复测三组中位耗时均下降，但结果仅代表本机 Pipe 路径，不能直接等同于 GM15 线上收益。

|输入形态|修复前 5 次（ns/op）|修复后 5 次（ns/op）|
|---|---|---|
|常规属性|5,951 · 5,999 · 6,036 · 6,441 · 6,497|5,839 · 5,892 · 5,840 · 5,841 · 5,828|
|属性较多|89,379 · 89,565 · 89,257 · 89,755 · 89,586|82,007 · 82,279 · 81,323 · 91,142 · 82,087|
|长文本属性|41,058 · 42,511 · 41,090 · 41,381 · 41,526|39,042 · 39,051 · 39,039 · 42,703 · 38,920|

此前记录的 8 路并发采样中，常规属性从 4.31 µs 降至 4.09 µs；属性较多从 33.98 µs 降至 17.49 µs。该并发数据尚未与本次串行复测在同一批次重跑，作为历史参考而非本次结论。

## P1 修复结果（已完成）

P1 在 `cc47e5b` 的基础上比较。两边均在同一台 Apple M4 Pro、Go 1.25.3、5 次运行中取中位数；耗时变化低于 1% 的项目只说明没有可确认的 CPU 改善，不作为性能收益宣称。

|输入形态|P0 中位数|P1 中位数|耗时变化|分配字节变化|分配次数变化|
|---|---:|---:|---:|---:|---:|
|常规属性，单线程|7,016 ns · 2,768 B · 30 次|6,702 ns · 2,383 B · 23 次|-4.5%|-13.9%|-23.3%|
|属性较多，单线程|84,369 ns · 7,898 B · 130 次|84,905 ns · 5,704 B · 119 次|+0.6%|-27.8%|-8.5%|
|长文本属性，单线程|40,229 ns · 2,823 B · 33 次|40,482 ns · 2,437 B · 26 次|+0.6%|-13.7%|-21.2%|
|常规属性，并发|4,795 ns · 2,773 B · 30 次|4,675 ns · 2,387 B · 23 次|-2.5%|-13.9%|-23.3%|
|属性较多，并发|19,721 ns · 8,198 B · 130 次|17,596 ns · 5,873 B · 119 次|-10.8%|-28.4%|-8.5%|

P1 的确定收益是减少通用适配层分配；属性较多、单线程的耗时持平，说明内外两层 JSON 编码和写出仍主导这一路径。

## 当前调用链与成本位置

```text
业务 fields
  ↓
trackLogger.Track
  ├─ thinkingdata.ExtractFields
  │    ├─ 标量 / Reflect：直接拆分协议字段和 Properties
  │    └─ Object、Array、Namespace、Error、非法字段等：回退 MapObjectEncoder + ExtractEncoder
  ├─ thinkingdata.Data.WithJSONV2
  │    └─ ThinkingData 内层 JSON
  └─ StdLogger.PrintThingkingData
       └─ TrackWithChannel（同步外层 JSON）
          └─ utils.WithZapJSON 归还内层 buffer
```

对应代码：

- `track_logger.go:46`：通用 `Track` 调用 `thinkingdata.ExtractFields`。
- `thinkingdata/extract.go:15,88`：标量快速路径直接拆分协议字段和业务属性；不适用的字段回退 `MapObjectEncoder`。
- `thinkingdata/data.go:97`：在回调期间提供内层 JSON，回调结束即失效。
- `utils/zap2json.go:27`：编码后在 `defer` 中归还 Zap buffer。
- `stdlogger_itracker_impl.go:10`：在同步 `TrackWithChannel` 调用内消费内层 JSON。
- `baisc_encoder.go:14`：每条外层日志生成 `log_xid`，拼接全局字段和动态字段。
- `truncate_write_syncer.go:48`：默认长度未超限时直接调用底层 `Write`。

去除系统写入后的分配采样显示：Zap buffer 约占 27.4%，`Data.MarshalAsJsonV2` 约占 19.6%，`MapObjectEncoder.AddInt64` 约占 16.1%，`MapObjectEncoder.AddString` 约占 10.7%，`glsEncoder.EncodeEntry` 直接占 10.7%（累计 38.7%）。

这说明 `NewMapObjectEncoder` 的创建本身不是主要问题；它在样本中约占 1.8%，更大成本来自逐字段填充 map、两层 JSON 和外层日志字段处理。

`AddInt64` 和 `AddString` 会把值存入 `map[string]interface{}`。编译器逃逸分析确认这些值会逃逸；但属性较多的额外成本还包括 map 扩容、字段转换和两层 JSON 编码，不能把 4.4 倍分配次数完全归因为装箱。

## 问题、代码与解决方案对应表

以下百分比来自 `pprof -alloc_space`：它表示采样期间**分配字节**归属，不是 CPU 时间比例，也不能简单相加为 100%。常规属性事件总分配为 3,816 B；例如 Zap buffer 的 27.4% 约对应 1 KiB 临时内存。

### 问题 1：内层 JSON 的 Zap buffer 未归还对象池（P0，已完成）

**修复内容**

`utils.WithZapJSON` 将刚生成的 JSON 交给 ThinkingData。当前日志的外层 JSON 写完后，立即执行 `buffer.Free()` 归还临时内存；之后不会再继续使用它。

```go
func WithZapJSON(data []zap.Field, consume func([]byte)) error {
    buffer, err := jsonEncoder.EncodeEntry(entry, data)
    if err != nil {
        if buffer != nil { buffer.Free() }
        return err
    }
    defer buffer.Free()
    encoded := buffer.Bytes()
    encoded = encoded[:len(encoded)-len(zapcore.DefaultLineEnding)]
    consume(encoded)
    return nil
}
```

BI、BigQuery、Elasticsearch 等模块通过 `Zap2Json` 取得 JSON 时，会先拿到一份独立副本。函数返回后它们仍可继续保存和使用该内容；原来的临时内存可以立即回收。ThinkingData 的当前写日志过程不需要这次复制。

`Zap2Json` 有 6 个生产调用点。在当前默认配置下，只要 BigQuery 日志带业务属性，就会先生成业务属性 JSON，再生成整条日志 JSON，共两次；两次临时内存都会归还。开启 `UseRecord` 或没有业务属性时，则只有一次。

编码器会在 JSON 末尾添加换行。代码按当前换行配置的实际长度去除它；日后调整换行配置时，需要重新运行对应测试确认结果。

**验证结果**

- 三组 5 次中位数基准见“P0 修复结果”；分配字节明显下降，日志协议未变化。
- `TestThinkingDataOutputSurvivesInnerBufferReuse` 连续写出 128 条不同事件，逐条校验内层 `sequence`。
- `TestZap2JsonReturnedBytesRemainValid` 验证其他模块拿到的 JSON 在后续 128 次编码后仍保持正确。
- 相关单元测试与 `-race` 通过；全量测试仍有既存的 `TestTruncateWriteSyncer_ExceedLimit_WithExtraFields` 断言失败，不属于本次修改。

### 问题 2：通用 Track 为每条日志构造中间 map（P1，已完成）

**问题与根因**

业务已持有 `[]zap.Field`，旧通用路径仍会先逐项写入 `MapObjectEncoder.Fields`，再拆出协议字段和 `Properties`。数值字段写入约占分配 16.1%，字符串字段约占 10.7%；`NewMapObjectEncoder` 自身仅约 1.8%，所以只给 encoder 加对象池不是有效方向。

**修复内容**

`track_logger.Track` 现在调用 `thinkingdata.ExtractFields(fields)`。纯标量、作为属性的 `Reflect` 和 `Skip` 字段直接处理：协议字段存入局部变量，只有业务属性写入最终要交给 JSON 的 `Properties` map。这样避免了“先把 6 个协议字段放入 map，再删除”的中间转换。

`Object`、`Array`、`Inline`、`Namespace`、`Error`、`Stringer`、未知类型或非法字段名会自动回退原有的 `MapObjectEncoder → ExtractEncoder` 路径。原因是这些字段可能延迟执行用户的 marshaler、产生附加错误字段或改变嵌套结构，不能靠简单类型断言替代。

```go
// track_logger.go
data, err := thinkingdata.ExtractFields(fields)
logger.PrintThingkingData(data)
```

**行为边界与验证**

- 保留字段仍按最后一个同名字段生效；`#time`、`#ip`、`#uuid` 和 `USER_ADD` 仍交由现有 ThinkingData 校验处理；`#first_check_id` 保持在 `Properties`，没有改变协议。
- `TestExtractFieldsMatchesMapObjectEncoder` 将标量、重复字段、非法键、`Reflect`、`Object`、`Array`、`Namespace`、`Inline`、`Error`、`Stringer` 与旧路径逐项对照。
- `TestScalarFieldValueMatchesMapObjectEncoder` 覆盖基础数值、字节、时间、浮点、复数和 `Reflect` 的实际存储类型，避免把 `zap.Int` 等误变为不同类型。
- 端到端 JSON 输出和 `-race` 测试通过；全量测试仍有既存的 `TestTruncateWriteSyncer_ExceedLimit_WithExtraFields` 断言失败，不属于本次修改。

**收益边界**

这项修复不能消除所有属性装箱：`Data.Properties` 的公开结构仍是 `map[string]interface{}`，业务属性仍须进入该 map。因此它的主要收益是减少适配层内存，而不是替代高频事件使用 `TrackWithTGAData` 的价值。

### 问题 3：字段名可能被正则校验两次（P2）

**现象与影响**

属性较多时，标量快速路径会逐个匹配字段名正则；复杂字段回退路径则在 `ExtractEncoder` 中匹配。保留下来的属性又会进入 `TrackWithType → formatProperties`，再次匹配相同正则。32 个额外属性会放大这部分重复工作。

**对应代码**

```go
// thinkingdata/extract.go:97-99（标量路径）
if !KeyPattern.MatchString(field.Key) {
    return emptyData, nil, false // 回退后由旧路径静默删除
}

// thinkingdata/extract.go:28-32（回退路径）
for k := range memoryEncoder.Fields {
    if !KeyPattern.MatchString(k) {
        delete(memoryEncoder.Fields, k)
    }
}

// thinkingdata/utils.go:74-79
for k, v := range d.Properties {
    if !checkPattern([]byte(k)) {
        return emptyData, errors.New("Invalid property key: " + k)
    }
    // ...
}
```

**解决方案**

保留第一次“过滤非法字段”的语义，为通用路径增加仅跳过第二次 key 正则校验的内部函数；时间格式化和 `USER_ADD` 数值校验仍必须执行。不要简单删除任一循环，否则可能把当前“忽略非法字段”改成“整条日志报错”。

**验证方式**

Go 允许在 `for range` 中删除当前 map 的键；这里没有可依赖的遍历顺序或并发语义。优化时要保持“非法字段静默删除”、保留字段拆分和重复 key 最终值的结果。

覆盖合法/非法 key、账户或访客 ID 缺失、`#event_id` 非法、`#time`、`#uuid`、`USER_ADD`、重复 key、Object/Array/Reflect 字段。特别加入“非法字段 + 保留字段冲突”的组合用例；属性较多基准应是主要观察对象。

### 问题 4：ThinkingData JSON 被作为字符串再次包进 Logbus JSON（P4）

**现象与影响**

`PrintThingkingData` 先获得内层 JSON，再把它作为 `zap.ByteString(MsgBody, bytes)` 交给外层 JSON encoder。常规属性的内层 JSON 是 269 B，完整外层日志是 595 B；内层的引号需要在外层字符串中转义。属性较多时内层 9,037 B、外层 9,491 B，额外比例较小但编码仍发生两次。

**对应代码**

```go
// stdlogger_itracker_impl.go:10-13
err := data.WithJSONV2(func(bytes []byte) {
    s.TrackWithChannel(THINKINGDATA, zap.ByteString(MsgBody, bytes))
})
```

**解决方案与边界**

只有在 Logserver、Fluent Bit、查询脚本和所有消费者确认可以把 `msg` 从 JSON 字符串改为嵌套 JSON 对象时，才能考虑改变外层格式。这个改动可减少一次转义/编码，但属于日志协议修改；在下游契约未确认前，不应实施。当前更安全的先后顺序是先修复 buffer 生命周期，再处理字段和类型化入口。

### 问题 5：每条外层日志都会构造追踪字段并立即写出（P2 / P3）

**对应代码**

- `baisc_encoder.go:14-28`：生成 `log_xid`，创建 `head` 切片并拼接全局/动态字段。
- `export_setting.go:47-51`：`BufferedWriteSyncer` 固定为 256 KiB、30 秒刷出。
- `config_snapshot.go:47-59`：仅在 `WithBufferedStdout(true)` 时启用该 writer；默认未启用。

**问题与方案**

`glsEncoder` 的临时字段切片和 XID 是全日志公共成本，分配采样中其直接占 10.7%、累计占 38.7%。但它影响所有日志类型，且动态全局字段、字段覆盖顺序与追踪 ID 都是兼容约束；应先有独立基准再改。

未缓冲写出在本地 Pipe CPU 采样中由系统调用主导。可以在单个 GM15 压测实例启用 `WithBufferedStdout(true)` 进行对照，但这只减少写系统调用：不减少单条 JSON 编码、分配或消息大小。必须验证日志可见性延迟、Fluent Bit 丢弃、SIGTERM 优雅退出与 `Close()` 刷出完整性。

### 低优先级：缓存 Tracker（P5）

`Tracker(THINKINGDATA)` 每次约 63.67 ns、104 B、4 次分配；相对常规 `Track` 的 6.11 µs、3,816 B、32 次分配，耗时约 1%。可以将只读 tracker 缓存在包级或事件专用对象中，但它不是压测主矛盾。

### 为什么不直接给 MapObjectEncoder 加对象池

`NewMapObjectEncoder` 只占约 1.8%，而 `ExtractEncoder` 仍然把它的 map 交给 `Data.Properties`。对象池要么过早归还导致并发日志串写，要么复制/清空 map 后再归还；后者会新增遍历和复制。相比之下，buffer 生命周期修复和类型化入口的收益边界更明确、风险更小。

## 优化建议与优先级

|优先级|建议|收益依据|风险与前置条件|建议决策|
|---|---|---|---|---|
|P0（完成）|日志写完后归还内层 JSON 的临时内存|三组事件分配字节降低 27.7%、85.5%、86.1%；常规路径少 2 次分配|ThinkingData 写完一条日志后立即归还临时内存。其他模块通过 `Zap2Json` 取得内容时会拿到独立副本。BigQuery 处理带业务属性的日志时会生成两段 JSON，两次临时内存都会归还。去除 JSON 末尾换行时按当前换行配置的实际长度计算；日后调整配置需重新运行测试|已由 `cc47e5b` 完成；后续若引入异步 CoreWrapper，再检查其写日志完成后才回收临时内存|
|P0|治理事件与属性：去掉无消费方的字段；对非关键事件按业务规则采样|属性较多单条成本是常规属性约 14 倍；减少一条日志可同时减少编码、分配和写出|会改变分析数据，需要 BI/运营确认保留字段、采样率和补偿口径|优先做；先产出事件/字段白名单|
|P3|验证 `WithBufferedStdout(true)` 的效果|未缓冲 Pipe 采样中系统写入占主导；现有配置已支持缓冲输出|本地 Pipe 不能外推线上；缓冲会引入可见性延迟，异常退出前可能遗失未刷出的日志；当前全局缓冲配置需确认所有同进程日志的影响|在 GM15 压测以 1 个实例试验，不直接全量开启|
|P1（已完成）|通用 `Track` 的标量字段快速路径|常规属性少 7 次分配；属性较多少 11 次；并发属性较多吞吐提升约 12%|复杂或非法 `zap.Field` 回退旧路径；`Properties` 仍是 `map[string]interface{}`，不能消除全部属性装箱|保持 `ExtractFields` 与旧路径的对照测试|
|P1|高频核心事件迁移到 `thinkingdata.Track` + `TrackWithTGAData`|可进一步绕过 `ExtractFields` 和通用协议适配层|直接 API 对非法属性键返回错误；通用路径会先删除非法键，语义不完全相同。UUID、时间和保留字段也要逐项对齐|先补齐等价基准和 JSON 结构对比，再只迁移 1～2 个高频事件|
|P5|为业务调用方缓存 `Tracker(THINKINGDATA)`|每次创建约 64 ns、104 B、4 次分配；`trackLogger` 配置创建后只读|改善约占常规路径 1% 的耗时，但可减少 12.5% 的分配次数；收益有限|可随业务改动顺手完成，不应单独作为压测主方案|
|P2|在通用路径消除重复的字段名校验，并用轻量校验替代正则|`ExtractEncoder` 先校验每个字段名，`formatProperties` 又校验保留下来的属性名|map 遍历中删除键是合法写法；必须保持非法字段静默删除、保留字段拆分和重复 key 的结果|仅为通用路径增加“已校验”内部标记；补充非法键与保留字段冲突、`USER_ADD`、`#time`、`#uuid`、重复 key 回归用例后再做|
|P2|减少 `glsEncoder` 外层临时字段切片与每条 `log_xid` 的构造成本|其累计分配占比 38.7%；每条都会构造头部切片并生成 XID|影响所有日志类型，且日志追踪字段是平台约束；对象池需要证明 Encoder 不保留切片|先做独立微基准和全日志兼容性测试；不与 ThinkingData 优化绑定发布|
|P4|评估将内层 ThinkingData JSON 作为原始 JSON 而非转义字符串写入外层日志|常规属性下内层 269 B、外层 595 B，双层编码造成额外 CPU 与字节；属性较多时增量相对较小|改变 `msg` 的类型和 Logserver/Fluent Bit/查询链路协议，风险最高|先确认下游契约；未确认前不得改动|

## 不建议的方向

### 直接复用 MapObjectEncoder / map

不建议。`ExtractEncoder` 将 `memoryEncoder.Fields` 直接作为 `Data.Properties` 继续传给 JSON 编码。若在编码完成前归还对象池，多个请求会共享同一 map，可能产生字段串写、并发竞态或错误日志。先复制再回收会重新引入一份 map 拷贝，未必优于类型化路径。

### 只靠提高并发

不建议作为主方案。8 路并发没有线性提升吞吐，且不会减少每条 3.8 KiB～53 KiB 的分配和下游日志量。并发只能作为吞吐配置，不是成本治理手段。

### 直接删除 TruncateWriteSyncer

不建议。常规日志在长度阈值以内只经过一次长度判断，几乎不是热点；删除会失去超大日志保护。

## 推荐实施顺序

### 阶段 1：修复 buffer 生命周期并建立兼容性用例（已完成）

目标：消除确定的 buffer 池复用缺口，不改变日志协议。

- ThinkingData 写完当前日志后归还临时内存；其他模块通过 `Zap2Json` 拿到独立副本。
- 已覆盖 JSON 等价、连续 128 条写出、其他模块拿到的 JSON 在后续编码后保持正确，以及 `-race`。`BufferedStdout`、错误路径和自定义异步 CoreWrapper 仍应在引入相关实现时补测。
- 用同一基线复测，收益主要体现在 `B/op`、`allocs/op` 和 GC，而不改变事件内容。

### 阶段 2：压测侧验证写出缓冲

目标：确认系统写入是否是 GM15 压测 CPU 的主要来源，并确认日志延迟和退出刷盘行为。

- 在单个压测实例开启 `WithBufferedStdout(true)`，保持事件内容不变。
- 同时采集：进程 CPU、GC、标准输出字节/秒、Fluent Bit 输入/丢弃、日志从产生到可检索的延迟、优雅退出后的尾部日志完整性。
- 使用常规属性、属性较多、长文本三组基准复测；不要将本地 Pipe 数据直接当作线上收益承诺。

### 阶段 3：减少字段和迁移高频事件

目标：从源头降低每条事件成本。

- 以事件名统计 QPS、平均属性数、平均/95 分位消息字节数、消费方和保留期限。
- 先移除无消费方的属性；需要采样的事件明确采样规则和 BI 口径。
- 选择 1～2 个高频且字段固定的事件，使用类型化 ThinkingData 构造并调用 `TrackWithTGAData`。
- 对同一输入解析后比较内层 JSON：协议字段、属性、错误行为必须一致；时间字段需使用固定时钟或排除动态时间后比较。

### 阶段 4：库内通用路径优化

目标：不改变通用 `Track(fields...)` 对业务的行为。

- 为“合法字段、非法字段、账户/访客 ID、事件 ID、`#time`、`#uuid`、`USER_ADD`”建立回归测试。
- 先处理重复校验，再评估 `glsEncoder` 临时切片；每次只改一个变量。
- 不改外层 `msg` 格式，除非 Logserver 和所有消费者确认可以接受原始 JSON 类型。

## 验收标准

每项改动都必须同时满足以下条件：

1. 使用现有基线的相同输入、相同命令运行 5 次，以中位数比较 `ns/op`、`B/op`、`allocs/op`。
2. 通过 `TestThinkingDataBaselineOutput` 和连续写出用例，保证外层 Logbus JSON、内层 ThinkingData JSON 和多条事件内容均合法。
3. 新增或通过对应语义用例，覆盖非法字段、保留字段和时间/UUID 行为。
4. 在压测环境确认：日志接收无新增丢弃、关键事件可检索、CPU/GC 改善不以不可接受的日志延迟为代价。

## 当前证据边界

本调研已覆盖 logbus 进程内的通用 ThinkingData 路径及本地 Pipe 写出。尚未验证 GM15 实例实际事件 QPS、属性分布、stdout 到 Fluent Bit 的背压、Logserver 接收性能和下游消费规则。因此，文中的优先级代表“最值得验证的顺序”，不是对线上收益的保证。

补充验证：8 路常规属性和属性较多基准均通过 `-race` 检测，未发现共享编码器的数据竞争。该模式的耗时和分配会被 race 插桩放大，不纳入性能基线。
