// Copyright 2021-2022 The jdh99 Authors. All rights reserved.
// UDP通信模块
// Authors: jdh99 <jdh821@163.com>

package udp

import (
	"github.com/jdhxyy/lagan"
	"net"
)

const (
	Version = 1

	tag = "udp"
)

// RxCallback 接收回调函数
type RxCallback func(data []uint8, ip uint32, port uint16)

var listener *net.UDPConn
var observers []RxCallback

// Load 模块载入
func Load(ip uint32, port uint16, frameMaxLen int) error {
	lagan.Info(tag, "init")

	var err error
	addr := &net.UDPAddr{IP: []uint8{uint8(ip >> 24), uint8(ip >> 16), uint8(ip >> 8), uint8(ip)}, Port: int(port)}
	listener, err = net.ListenUDP("udp", addr)
	if err != nil {
		lagan.Error(tag, "bind pipe net failed")
		return err
	}

	go func() {
		data := make([]uint8, frameMaxLen)
		for {
			num, addr, err := listener.ReadFromUDP(data)
			if err != nil {
				lagan.Error(tag, "listen net failed:%v", err)
				continue
			}
			if num <= 0 {
				continue
			}
			lagan.Debug(tag, "udp rx:%v len:%d", addr, num)
			lagan.PrintHex(tag, lagan.LevelDebug, data[:num])
			notifyObservers(data[:num], addr)
		}
	}()

	return nil
}

func notifyObservers(data []uint8, addr *net.UDPAddr) {
	ipValue := addr.IP.To4()
	ip := (uint32(ipValue[0]) << 24) + (uint32(ipValue[1]) << 16) + (uint32(ipValue[2]) << 8) + uint32(ipValue[3])

	n := len(observers)
	for i := 0; i < n; i++ {
		observers[i](data, ip, uint16(addr.Port))
	}
}

// RegisterObserver 注册观察者
func RegisterObserver(callback RxCallback) {
	observers = append(observers, callback)
}

// Send 发送数据
func Send(data []uint8, ip uint32, port uint16) {
	addr := &net.UDPAddr{IP: []uint8{uint8(ip >> 24), uint8(ip >> 16), uint8(ip >> 8), uint8(ip)}, Port: int(port)}
	lagan.Debug(tag, "udp send:addr:%v len:%d", addr, len(data))
	lagan.PrintHex(tag, lagan.LevelDebug, data)

	_, err := listener.WriteToUDP(data, addr)
	if err != nil {
		lagan.Error(tag, "udp send error:%v addr:%v", err, addr)
		return
	}
}
