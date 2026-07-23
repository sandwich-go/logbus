# C2：移除日志热路径对 `Setting` 的依赖

## 结论

本次只做核心 C2：`Init` 接收的配置会在初始化时复制，Logger 创建时取得所需配置，此后打印日志不再读取全局 `Setting`。

保留动态日志级别。`SetLogLevel` 和 PMT 配置监听仍然可以在运行中调整所有 Logger 的日志级别。

## 为什么改

以前 `Setting` 同时承担两种职责：

- 保存初始化配置；
- 在每次打印日志时提供 `FetchLogContext`、默认 channel 等值。

这会带来两个问题：

- 调用方在 `Init` 后修改原始 `Conf` 或 `Setting`，可能悄悄改变已创建 Logger 的行为；
- 日志打印与配置写入并发时，全局可变对象会成为数据竞争入口。

核心 C2 把边界放在 `Init` 和 Logger 构造阶段。配置先复制，再交给 Logger；日志热路径只读 Logger 自己持有的值。

## 修改内容

### 1. 增加内部配置快照

新增 `configSnapshot`，复制 `Conf` 以及其中可变的 slice、map 和截断配置。

`Setting` 仍然保留，避免破坏公开 API，但它只是当前配置的一份独立副本。初始化完成后再修改 `Setting`，不会影响已创建 Logger。

`currentConfig` 的惰性默认路径也会填充默认 encoder 和 writer，保证未手动调用 `Init` 时创建的 Logger 能输出完整的日志键。

### 2. Logger 持有构造期配置

`StdLogger` 保存：

- `FetchLogContext`
- 默认 channel 和 tag
- encoder
- writer
- 其他构造 Logger 所需配置

`FetchLogContext` 仍然在每次日志调用时执行，因此请求上下文字段可以动态变化；固定的是“使用哪个回调函数”。

### 3. 全局 Logger 使用原子指针

包级 `StdLogger` 和 `GLogger` 改为原子替换。`Init` 更新全局 Logger 时，正在打印日志的调用不会读到一个正在被修改的对象。

### 4. 动态日志级别保持不变

日志级别没有固化进单个 Logger。所有 Logger 仍然读取进程级的原子日志级别：

- `SetLogLevel` 可即时生效；
- PMT 的 `ops_config.json` 监听可继续动态调整；
- `EnableTraceLevel` 仍是进程级开关。

### 5. 初始化前创建的 Logger

初始化前创建的 scope Logger 会在 `Init` 时刷新，并保留：

- 自定义 `FetchLogContext`
- `PrintAsError` 模式
- channel

旧的 depth Logger 缓存会清空，避免继续引用刷新前的 Logger。

## 明确不包含的改动

以下内容不属于核心 C2，本次没有改变其原有语义：

- `IgnoreLogicalError` 仍是进程级策略，只改为原子读取，避免数据竞争；
- `BufferedWriteSyncer` 仍使用原有公共对象和刷新方式；
- 没有增加 `configCore` 或逐 Logger 的错误字段转换；
- 没有改变 `Close` 的资源生命周期；
- 没有调整 tracking、ES logger、截断 writer 的业务行为。

`Init` 仍应只在启动阶段调用，不支持与另一次 `Init` 或 `NewScopeLogger` 并发执行。

## 主要文件

| 文件 | 修改 |
| --- | --- |
| `config_snapshot.go` | 配置复制、当前配置的原子发布、writer 准备 |
| `logger_init.go` | 使用配置快照初始化各模块 |
| `stdlogger.go` | Logger 保存自己的配置和 context 回调 |
| `global_stdloggers.go` | 构造 Logger 时注入配置；刷新早期 Logger |
| `simple.go` / `setgloballoger.go` | 原子读取和替换全局 Logger |
| `field.go` | `IgnoreLogicalError` 改为进程级原子值 |
| `config_snapshot_test.go` | 验证配置隔离、早期 Logger 刷新和并发读取 |

## 行为示例

```go
conf := logbus.NewConf(
    logbus.WithDefaultChannel("game"),
    logbus.WithFetchLogContext(fetchRequestFields),
)
logbus.Init(conf)

conf.DefaultChannel = "changed"
logbus.Setting.DefaultChannel = "also-changed"
```

上面的两次修改都不会改变已经创建的 Logger；它仍使用 `game` 和初始化时取得的 `fetchRequestFields`。

动态级别仍然有效：

```go
logbus.SetLogLevel(zapcore.DebugLevel)
```

现有 Logger 会立即开始输出 debug 日志。

## 验证结果

- 核心 C2 定向测试：9 项全部通过；
- 核心 C2 竞态测试：9 项全部通过，没有数据竞争；
- 全仓普通测试：67 项通过、1 项既有失败；
- 既有失败：`TestTruncateWriteSyncer_ExceedLimit_WithExtraFields` 期望日志级别为 `error`，当前实现输出 `warn`，与本次 C2 无关；
- `go vet ./...`：保留 1 个既有提示，`stdlogger_itracker_impl.go:31` 存在无参数 `append`；
- `git diff --check`：通过。
