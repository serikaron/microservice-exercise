package pkg

import (
	"flag"
	"fmt"
)

type FlagItem interface {
	Attach()
}

var (
	AuthAddr = newAddrItem("auth")
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

func (i *UintItem) Get() interface{} {
	return i.Val
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

func (i *StringItem) Get() interface{} {
	return i.Val
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
