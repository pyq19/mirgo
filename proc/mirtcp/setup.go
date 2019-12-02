package mirtcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc"
)

func init() {
	proc.RegisterProcessor("mir.tcp.ltv", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {
		bundle.SetTransmitter(new(TCPMessageTransmitter))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
