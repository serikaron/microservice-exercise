package pkg

import (
	"flag"
	"fmt"
)

type FlagItem interface {
	Attach()
}

var (
	AuthAddr  = newAddrItem("auth")
	ChatAddr  = newAddrItem("chat")
	RedisAddr = newAddrItem("redis")
	//ServerCertPath = &StringItem{name: "server-cert-path", def: "server.cert", desc: "server-cert-path"}
	//ServerKeyPath  = &StringItem{name: "server-key-path", def: "server.key", desc: "server-key-path"}
	CertsPath         = &CertsPathItem{path: StringItem{name: "cert-path", def: "", desc: "cert-path"}}
	IntegrationKey    = &StringItem{name: "integration-test-key", def: "", desc: "integration test key"}
	IntegrationEnable = &BoolItem{name: "integration-test-enable", def: false, desc: "integration test enable"}
)

func ParseItem(itemList []FlagItem) {
	for _, item := range itemList {
		item.Attach()
	}
	flag.Parse()
}

type UintItem struct {
	name string
	def  uint
	desc string
	Val  uint
}

func (i *UintItem) Attach() {
	flag.UintVar(&i.Val, i.name, i.def, i.desc)
}

type StringItem struct {
	name string
	def  string
	desc string
	Val  string
}

func (i *StringItem) Attach() {
	flag.StringVar(&i.Val, i.name, i.def, i.desc)
}

type BoolItem struct {
	name string
	def  bool
	desc string
	Val  bool
}

func (i *BoolItem) Attach() {
	flag.BoolVar(&i.Val, i.name, i.def, i.desc)
}

type AddrItem struct {
	host StringItem
	port UintItem
}

func (ai *AddrItem) Attach() {
	ai.host.Attach()
	ai.port.Attach()
}

func (ai *AddrItem) Addr() string {
	return fmt.Sprintf("%s:%v", ai.host.Val, ai.port.Val)
}

func newAddrItem(prefix string) *AddrItem {
	return &AddrItem{
		host: StringItem{
			name: prefix + "-host",
			def:  prefix,
			desc: prefix + "-host",
			Val:  "",
		},
		port: UintItem{
			name: prefix + "-port",
			def:  0,
			desc: prefix + "-port",
			Val:  0,
		},
	}
}

type CertsPathItem struct {
	path StringItem
}

func (cp *CertsPathItem) Attach() {
	cp.path.Attach()
}

func (cp *CertsPathItem) Key() string {
	return cp.path.Val + "/server.key"
}

func (cp *CertsPathItem) Pem() string {
	return cp.path.Val + "/server.pem"
}
