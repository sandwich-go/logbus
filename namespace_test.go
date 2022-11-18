package logbus

import (
	"testing"
)

func TestNameSpace(t *testing.T) {
	Warn("namespace", Int("int1", 2), NameSpace("namespace1"), String("str1", "str"), Int("int1", 1))
	//{"log_level":"warn","date":"2022-11-18T17:26:14.595+0800","dd_meta_channel":"Game","tags":"Advance","log_xid":"cdrkvdlvqc7lehqluhag",
	//"server_id":"cdrkvdlvqc7lehqluh70","server_ip":"10.0.49.62","server_birth":1668763574,"host_name":"biyongzedeMacBook-Pro.local","int1":2,
	//"namespace1":{"str1":"str","int1":1,"glog-msg":"namespace"}} // todo glog-msg的位置不正确
}
