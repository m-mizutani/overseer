package overseer

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/tcpassembly"
	"log"
	"time"
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
	s.init = false
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

func (x *SessionTable) Timeout(trimTime time.Time) {
	log.Println("---- FLUSHING ----")
	x.assembler.FlushOlderThan(trimTime)
}

type Session struct {
	netFlow, tcpFlow gopacket.Flow
	dataSize         int
	init             bool
	first            time.Time
	last             time.Time
}

func (s *Session) Reassembled(rs []tcpassembly.Reassembly) {
	for _, r := range rs {
		s.dataSize += len(r.Bytes)
		if !s.init {
			s.first = r.Seen
			s.last = r.Seen
			s.init = true
		} else {
			if s.last.Before(r.Seen) {
				s.last = r.Seen
			}
		}
	}
}

func (s *Session) ReassemblyComplete() {
	log.Println("comlete", s.netFlow, s.tcpFlow, s.dataSize, s.last.Sub(s.first))
}
