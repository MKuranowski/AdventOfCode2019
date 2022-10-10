package day23

import (
	"io"

	"github.com/MKuranowski/AdventOfCode2019/intcode"
	"github.com/MKuranowski/AdventOfCode2019/util/deque"
	"github.com/MKuranowski/AdventOfCode2019/util/set"
)

type Packet struct{ Dest, X, Y int }

type Network struct {
	deviceByAddress map[int]*intcode.SyncInterpreter

	readyDevices   deque.Deque[*intcode.SyncInterpreter]
	blockedDevices set.Set[*intcode.SyncInterpreter]

	lastNatPacket Packet
}

func NewNetwork(nicCode *intcode.SyncInterpreter, devices int) *Network {
	n := &Network{
		deviceByAddress: make(map[int]*intcode.SyncInterpreter),
		readyDevices:    deque.NewDeque[*intcode.SyncInterpreter](),
		blockedDevices:  make(set.Set[*intcode.SyncInterpreter]),
	}

	for i := 0; i < devices; i++ {
		d := nicCode.Clone()
		d.Input.PushBack(i)

		n.deviceByAddress[i] = d
		n.readyDevices.PushBack(d)
	}

	return n
}

func (n *Network) SendPacket(p Packet) {
	if p.Dest == 255 {
		n.lastNatPacket = p
	} else {
		d := n.deviceByAddress[p.Dest]
		d.Input.PushBack(p.X)
		d.Input.PushBack(p.Y)

		if n.blockedDevices.Has(d) {
			n.blockedDevices.Remove(d)
			n.readyDevices.PushBack(d)
		}
	}
}

func (n *Network) RunUntilBlocked() Packet {
	for n.readyDevices.Len() > 0 {
		// Get a device that's ready
		d := n.readyDevices.PopFront()

		// Execute it until it's blocked on input
		state := d.ExecAll()

		// Ensure state is BlockedOnInput
		if state == intcode.SyncExecutionStateHalted {
			panic("nic has turned itself off")
		} else if state == intcode.SyncExecutionStateReady {
			panic("intcode.SyncInterpreter.ExecAll() signalled ready")
		}

		// Move the device to blocked queue
		n.blockedDevices.Add(d)

		// Process device outputs
		for d.Output.Len() > 0 {
			p := Packet{}
			p.Dest = d.Output.PopFront()
			p.X = d.Output.PopFront()
			p.Y = d.Output.PopFront()

			n.SendPacket(p)
		}
	}

	return n.lastNatPacket
}

func (n *Network) RunUntilNATPacket() Packet {
	for {
		p := n.RunUntilBlocked()
		if p.Dest == 255 {
			return p
		}

		// Unblock the network by sending -1 to every receiver
		for d := range n.blockedDevices {
			d.Input.PushBack(-1)
			n.readyDevices.PushBack(d)
			delete(n.blockedDevices, d)
		}
	}
}

func (n *Network) RunWithNAT() Packet {
	lastY := -1

	for {
		p := n.RunUntilNATPacket()
		p.Dest = 0

		if lastY == p.Y {
			return p
		} else {
			lastY = p.Y
			n.SendPacket(p)
		}
	}
}

func SolveA(r io.Reader) any {
	nicCode := intcode.NewSyncInterpreter(r)
	network := NewNetwork(nicCode, 50)

	p := network.RunUntilNATPacket()
	return p.Y
}

func SolveB(r io.Reader) any {
	nicCode := intcode.NewSyncInterpreter(r)
	network := NewNetwork(nicCode, 50)

	p := network.RunWithNAT()
	return p.Y
}
