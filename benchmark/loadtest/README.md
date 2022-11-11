locust --master -f dummy.py --headless -u 5000 -r 10 // -cpu-profile cpu.pprof -cpu-profile-duration 90s

go run locust.go --max-rps 50000 --request-increase-rate 200

# 结论

机器4核8G，supervisor把输出重定向到文件，fluentd读文件发送到logserver    
单进程模拟5000个用户，每秒打印10个tga日志, 每个日志730 bytes    
在20000 messages/second (相当于14MB/s) 以下时很平稳，load大约在0.7，之后负载会快速爬升  
测试三次分别在 41876.70 msg/s 43271.90 msg/s 和 40669.90 msg/s 时，worker失联弹出下面的报错，此时load达到4.7(大约:进程3 ruby0.8 supervisord0.7)

```bash
[2021-01-13 07:45:20,072] ip-10-83-59-169.us-west-2.compute.internal/INFO/locust.runners: Worker ip-10-83-59-169.us-west-2.compute.internal_5b5e08d2ed3447df99098ad2a6538490 failed to send heartbeat, setting state to missing.
[2021-01-13 07:45:20,072] ip-10-83-59-169.us-west-2.compute.internal/INFO/locust.runners: The last worker went missing, stopping test.
```

查看cpu.pprof可以看到cpu时间主要消耗在poll fd write和runtime.findrunnable  
logserver侧接收一直平稳，没有异常，峰值io达到最大能力的一半  
https://github.com/uber-go/zap/pull/782 合入后可能性能会有提升  