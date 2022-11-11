package fluentd

import (
	"log"
	"testing"
	"time"

	"bitbucket.org/funplus/sandwich/base/scmd"

	"github.com/fluent/fluent-logger-golang/fluent"
)

// https://docs.fluentd.org/installation/install-by-dmg
// /var/log/td-agent/td-agent.log
// /etc/td-agent/td-agent.conf
// https://docs.fluentd.org/input/unix#example-configuration
func TestSocketConnect(t *testing.T) {
	if scmd.IsFalse(scmd.GetOptWithEnv("sandwich_test_enable_fluentd")) {
		return
	}
	Init(&fluent.Config{
		FluentNetwork:    "unix",
		FluentSocketPath: "/tmp/fun-collector.sock",
		MarshalAsJSON:    true,
		Async:            false,
	})
	defer Close()
	tag := "debug.access"
	var data = map[string]string{
		"foo":  "bar",
		"hoge": "hoge"}
	for i := 0; i < 3; i++ {
		e := GetClient().Post(tag, data)
		if e != nil {
			log.Println("Error while posting log: ", e)
		} else {
			log.Println("Success to post log")
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func TestTcpConnect(t *testing.T) {
	if scmd.IsFalse(scmd.GetOptWithEnv("sandwich_test_enable_fluentd")) {
		return
	}
	Init(&fluent.Config{
		FluentNetwork: "tcp",
		FluentHost:    "127.0.0.1",
		FluentPort:    24224,
	})
	defer Close()
	tag := "debug.access"
	var data = map[string]string{
		"huge": "package",
		"big":  "queue"}
	for i := 0; i < 3; i++ {
		e := GetClient().Post(tag, data)
		if e != nil {
			log.Println("Error while posting log: ", e)
		} else {
			log.Println("Success to post log")
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
