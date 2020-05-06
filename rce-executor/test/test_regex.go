package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func main() {

	s := `	
[0m etcd.service - ETCD Service
   Loaded: loaded (/etc/systemd/system/etcd.service; disabled; vendor preset: disabled)
   Active: active (running) since Tue 2018-05-15 17:55:04 CST; 2h 9min ago
 Main PID: 13701 (etcd)
   CGroup: /system.slice/etcd.service
           13701 /usr/bin/etcd --name=mysql1 --initial-advertise-peer-urls=http://192.168.0.185:2380 --listen-peer-urls=http://192.168.0.185:2380 --listen-client-urls=http://192.168.0.185:2379...

May 15 17:55:44 mysql1 etcd[13701]: updated the cluster version from 3.0 to 3.2
May 15 17:55:44 mysql1 etcd[13701]: enabled capabilities for version 3.2
May 15 17:58:11 mysql1 etcd[13701]: failed to send out heartbeat on time (exceeded the 100ms timeout for 169.793744ms)
May 15 17:58:11 mysql1 etcd[13701]: server is likely overloaded
May 15 17:58:11 mysql1 etcd[13701]: failed to send out heartbeat on time (exceeded the 100ms timeout for 169.852712ms)
May 15 17:58:11 mysql1 etcd[13701]: server is likely overloaded
May 15 17:58:11 mysql1 etcd[13701]: failed to send out heartbeat on time (exceeded the 100ms timeout for 79.890441ms)
May 15 17:58:11 mysql1 etcd[13701]: server is likely overloaded
May 15 17:58:11 mysql1 etcd[13701]: failed to send out heartbeat on time (exceeded the 100ms timeout for 79.945134ms)
May 15 17:58:11 mysql1 etcd[13701]: server is likely overloaded
`
	lines := strings.Split(s, "\n")
	re1 := regexp.MustCompile(`^\s*Main PID:\s*(\d+)[\s\w]*`)
	for _, v := range lines {
		ss := re1.FindStringSubmatch(v)
		if len(ss) > 0 {
			for _, v2 := range ss {
				fmt.Printf("%s\n", v2)
			}
		}
	}
	fmt.Printf("%q\n", strings.SplitN("a,b,c", ",", 2))
	re2 := regexp.MustCompile(`\$\{\{[\w]+\}\}`)
	ss := `execute script statement error at line 1: connect target failed:ssh: handshake failed: ssh: unable to authenticate, attempted methods [none password], no supported methods remain, statement='connect ssh target target185 host=enmo.wicp.net username=root password=${password} port=8185'`
	tokens := re2.FindAllString(ss, -1)
	log.Printf("tokens=%q\n", tokens)
	for _, v := range tokens {
		log.Printf("var=%s", v[3:len(v)-2])
	}
	re3 := regexp.MustCompile(`(?i)password[\s|\S]*`)
	log.Printf(re3.ReplaceAllString(ss, "password ********"))
}
