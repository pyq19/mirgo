package mirtcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc"
)

func init() {
	proc.RegisterProcessor("mir.client.tcp", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {
		bundle.SetTransmitter(new(ClientTCPMessageTransmitter))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
	proc.RegisterProcessor("mir.server.tcp", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {
		bundle.SetTransmitter(new(ServerTCPMessageTransmitter))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
