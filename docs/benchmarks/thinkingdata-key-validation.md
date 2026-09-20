# ThinkingData 属性名重复校验优化

## 变更范围

通用打点的 ExtractFields（标量快速路径）和 ExtractEncoder（复杂字段回退）已校验或过滤属性名，后续通过私有 user / trackWithType 入口传递 checkKeys=false，只省略 formatProperties 中重复的属性名正则匹配。

公开 Track / TrackWithType / User 仍使用 checkKeys=true。第一次非法字段静默过滤、事件名和事件 ID 检查、协议字段提取、USER_ADD 数值检查、time.Time 属性格式化、复杂 Field 和重复 key 的处理均保持原有行为。未改变日志格式和输出字节，也未修改 Tracker 的创建行为。

## 测试环境与方法

- 基线生产代码：23fe722a8b87646e81221d155e0951731f413bc7。
- Go：go1.25.3；系统/架构：darwin/arm64；CPU：Apple M4 Pro。
- 前后使用相同测试代码，编译为独立测试二进制；先运行基线，再运行修改版。
- 每项 6 次，每次 benchtime=500ms，cpu=1；下表为中位数，不是置信区间或线上压测结果。
- 微基准使用固定时间；构造字段在计时之外；每次调用都重新提取，不复用已消费的 MapObjectEncoder。
- Scalar3/32/128：3/32/128 个整数业务属性；Object32：32 个整数属性加一个自定义对象，触发回退；Invalid32：32 个合法属性加一个非法属性，触发过滤回退；UserAdd32：32 个数值属性的 USER_ADD。
- DirectTrack32：公开 thinkingdata.Track 对照组，保留完整校验；每次补回被提取的固定 #time。
- 完整 Track 使用现有 thinkingdata_baseline_benchmark_test.go：Conventional 为 3 个业务属性；AttributeRich 额外增加 32 个 256 字节属性，共 35 个业务属性；LongText 额外增加一个 4 KiB 字符串。
- 完整 Track 默认写入持续被读取的本地 OS pipe；Discard 不包含 OS 写出。两者都不包含 Fluent Bit、Logserver 或 TGA。

## 结果

耗时单位为 µs/op；变化为修改后相对修改前，负数表示耗时减少。B/op 和 allocs/op 均按 6 次结果取中位数。

| 场景 | 修改前 | 修改后 | 耗时变化 | B/op 前 → 后 | allocs/op 前 → 后 |
|---|---:|---:|---:|---:|---:|
| Extract/Scalar3 | 1.723 | 1.309 | -24.0% | 416 → 368 | 7 → 4 |
| Extract/Scalar32 | 10.478 | 6.189 | -40.9% | 2936 → 2424 | 38 → 6 |
| Extract/Scalar128 | 39.900 | 21.945 | -45.0% | 11640 → 9592 | 134 → 6 |
| Extract/Object32 | 17.965 | 14.418 | -19.7% | 7864 → 7344 | 57 → 24 |
| Extract/Invalid32 | 17.863 | 13.941 | -22.0% | 7504 → 6992 | 53 → 21 |
| Extract/UserAdd32 | 10.223 | 6.072 | -40.6% | 2920 → 2408 | 37 → 5 |
| DirectTrack32 | 4.979 | 5.101 | 2.5% | 528 → 528 | 33 → 33 |
| Track/Conventional | 5.004 | 5.101 | 1.9% | 2104 → 2072 | 21 → 18 |
| Track/AttributeRich | 80.127 | 77.123 | -3.8% | 5381 → 4836.5 | 117 → 82 |
| Track/LongText | 36.657 | 36.782 | 0.3% | 2154 → 2106 | 24 → 20 |
| Track/ConventionalDiscard | 5.050 | 4.931 | -2.3% | 2104 → 2072 | 21 → 18 |
| Track/AttributeRichDiscard | 80.064 | 77.700 | -3.0% | 5381 → 4837 | 117 → 82 |

32 属性提取耗时下降 40.9%，分配减少 32 次；完整属性较多的 Track 耗时下降 3.8%，分配从 117 次降至 82 次。小事件完整 Track 的 pipe 结果增加 1.9%，Discard 结果减少 2.3%，暂未显示一致的耗时改善；公开入口对照组耗时增加 2.5%，内存分配保持不变。这些小幅差异应继续结合测量波动看待，不据此宣称稳定提速或回退。

本次优化不减少日志量，不代表零点延迟峰值能下降同样比例。

## 验证

- 新增固定期望的回归测试，在基线生产代码（Go overlay）与修改版均通过。
- 覆盖非法 key、协议字段冲突、重复 key、事件 ID 类型和首字符、空/非法事件名、时间/IP/UUID、复杂 Object/Array/Stringer/Inline/Namespace、USER_ADD 数值限制、公开 API 非法 key 拒绝。
- go test -race ./thinkingdata ./utils 通过。
- go test -race . -run 'TestThinkingData|TestDisabledLevel|TestEnabledLevel' 通过。
- go test ./... 仅 TestTruncateWriteSyncer_ExceedLimit_WithExtraFields 失败：期望 error，实际 warn。已用修改前的根包测试二进制独立复现同一失败，与本次改动无关。
- git diff --check 通过。

## 复现命令

在基线与修改版分别运行；两边都必须使用本次新增的 extract_benchmark_test.go：

```sh
go test ./thinkingdata -run '^$' \
  -bench '^(BenchmarkExtractFields|BenchmarkDirectTrack32)$' \
  -benchmem -benchtime=500ms -count=6 -cpu=1

go test . -run '^$' \
  -bench '^BenchmarkThinkingDataBaselineGenericTrack(Conventional|AttributeRich|LongText|ConventionalDiscard|AttributeRichDiscard)$' \
  -benchmem -benchtime=500ms -count=6 -cpu=1
```

测量时不要并行运行前后基准。

原始结果：[提取前](key-validation-before-extract.txt)、[提取后](key-validation-after-extract.txt)、[完整打点前](key-validation-before-track.txt)、[完整打点后](key-validation-after-track.txt)。
