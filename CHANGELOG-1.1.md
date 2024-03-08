### v1.1.1 🌈 (2024-03-08 18:01:13)

#### 🐛  Bug Fixed
  * tga日志缺少properties字段 ([ecab5d2](https://github.com/sandwich-go/logbus/commit/ecab5d2efa7eb02453e5e7ee56097f11c3ba92c3)) (<small>[Daming Yang](daming.yang@centurygame.com)@2024-03-08 18:01:13 &#43;0800 &#43;0800</small>)

#### 🤖  Tools
  * **sem**: make changelog ([7587e39](https://github.com/sandwich-go/logbus/commit/7587e39cfab844dcc846a77aea45c4dd9236a06f)) (<small>[zhengyang.zhu](zhengyang.zhu@centurygame.com)@2024-03-08 11:47:05 &#43;0800 &#43;0800</small>)

### v1.1.0 (2024-03-08 11:28:30)

#### 🚀  New Feature
  * add GetZapLogger() ([3a5adc3](https://github.com/sandwich-go/logbus/commit/3a5adc37b874f24413373996e4bd5681ea71ccc6)) (<small>[zhengyang.zhu](zhengyang.zhu@centurygame.com)@2024-02-04 15:12:55 &#43;0800 &#43;0800</small>)

#### 🛠  Refactor
  * MarshalAsJson使用栈空间 ([9815e56](https://github.com/sandwich-go/logbus/commit/9815e565f6fafe24d3467e8d7ae6e19c88c77211)) (<small>[Daming Yang](daming.yang@centurygame.com)@2024-01-30 14:34:19 &#43;0800 &#43;0800</small>)
  * track判断是否输出逻辑前置 ([5151d83](https://github.com/sandwich-go/logbus/commit/5151d836762e89c7277c3a719080e2ab8abf9f2c)) (<small>[Daming Yang](daming.yang@centurygame.com)@2024-01-30 14:19:50 &#43;0800 &#43;0800</small>)
  * metrics family use LoadOrStore instead of singleflight ([3865ee1](https://github.com/sandwich-go/logbus/commit/3865ee1799a7ebf3d98f9cfa570eca733639511f)) (<small>[Daming Yang](daming.yang@centurygame.com)@2024-01-29 13:01:55 &#43;0800 &#43;0800</small>)
  * metrics family使用sync.Map和singleflight 替换全局锁 ([49feb6e](https://github.com/sandwich-go/logbus/commit/49feb6e636d35809adb97ef00735c5f4c7457cf3)) (<small>[Daming Yang](daming.yang@centurygame.com)@2024-01-29 12:33:51 &#43;0800 &#43;0800</small>)
  * update go runtime collectors all ([5a88474](https://github.com/sandwich-go/logbus/commit/5a88474bcb077609c5a9f4f39c88b967f31929fa)) (<small>[Daming Yang](daming.yang@centurygame.com)@2024-01-26 17:41:49 &#43;0800 &#43;0800</small>)

#### 💪  Commit
  * Merge pull request [#2](https://github.com/sandwich-go/2/issues/%!s(MISSING)) from ayamzh/1.1/develop ([c233794](https://github.com/sandwich-go/logbus/commit/c233794a0acc252d2f11ae2fa57fd2bcc5fe135a)) (<small>[Zhu Zhengyang](buptzhuzhengyang@foxmail.com)@2024-03-08 11:28:30 &#43;0800 &#43;0800</small>)
  * 简单文档 ([5fbb22c](https://github.com/sandwich-go/logbus/commit/5fbb22c192f70e652737ce70be05129989b9c77a)) (<small>[zheng.ma](zheng.ma@centurygame.com)@2024-03-06 13:41:51 &#43;0800 &#43;0800</small>)
  * Merge pull request [#1](https://github.com/sandwich-go/1/issues/%!s(MISSING)) from ayamzh/1.1/develop ([561d852](https://github.com/sandwich-go/logbus/commit/561d8522f5cbbe3f7117a66276e7ceffdb0a7c82)) (<small>[Zhu Zhengyang](buptzhuzhengyang@foxmail.com)@2024-03-06 11:56:23 &#43;0800 &#43;0800</small>)
  * 增加writeSyncer配置 ([dbb0fdc](https://github.com/sandwich-go/logbus/commit/dbb0fdcc3dfd55e93b808b088a5638ad61914420)) (<small>[zheng.ma](zheng.ma@centurygame.com)@2024-03-05 18:21:00 &#43;0800 &#43;0800</small>)
  * Merge branch '1.1/develop' of github.com:sandwich-go/logbus into 1.1/develop ([baec373](https://github.com/sandwich-go/logbus/commit/baec373d0d62255cde980b25bc77df13529d0479)) (<small>[Daming Yang](daming.yang@centurygame.com)@2024-01-30 14:20:08 &#43;0800 &#43;0800</small>)
  * tga event: 支持提供#event_id, 用于支持tga event update ([9bad4ed](https://github.com/sandwich-go/logbus/commit/9bad4ed2c32beb9ef447c88e92210aa08015f963)) (<small>[biyongze](yongze.bi@centurygame.com)@2024-01-29 16:53:04 &#43;0800 &#43;0800</small>)



