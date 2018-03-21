package overseer

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/tcpassembly"
	"log"
)

type SessionTable struct {
	streamFactory *SessionFactory
	streamPool    *tcpassembly.StreamPool
	assembler     *tcpassembly.Assembler
}

type SessionFactory struct {
}

func (f *SessionFactory) New(netFlow, tcpFlow gopacket.Flow) tcpassembly.Stream {
	// Create a new stream.
	s := &Session{}
	s.netFlow = netFlow
	s.tcpFlow = tcpFlow
	log.Println("New", netFlow, tcpFlow)
	return s
}

func NewSessionTable() *SessionTable {
	t := &SessionTable{}
	t.streamFactory = &SessionFactory{}
	t.streamPool = tcpassembly.NewStreamPool(t.streamFactory)
	t.assembler = tcpassembly.NewAssembler(t.streamPool)
	return t
}

func (x *SessionTable) ReadPacket(pkt gopacket.Packet) {
	if pkt.NetworkLayer() == nil || pkt.TransportLayer() == nil ||
		pkt.TransportLayer().LayerType() != layers.LayerTypeTCP {
		return
	}

	tcp := pkt.TransportLayer().(*layers.TCP)
	x.assembler.AssembleWithTimestamp(pkt.NetworkLayer().NetworkFlow(), tcp,
		pkt.Metadata().Timestamp)
}

type Session struct {
	netFlow, tcpFlow gopacket.Flow
	dataSize         int
}

func (s *Session) Reassembled(rs []tcpassembly.Reassembly) {
	for _, r := range rs {
		s.dataSize += len(r.Bytes)
	}
}

func (s *Session) ReassemblyComplete() {
	log.Println("comlete", s.netFlow, s.tcpFlow, s.dataSize)
}
