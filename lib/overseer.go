package overseer

import (
	"errors"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)

type Overseer struct {
	sourceName   string
	pcapHandle   *pcap.Handle
	packetSource *gopacket.PacketSource
}

func New() Overseer {
	seer := Overseer{}
	return seer
}

func (x *Overseer) SetPcapFile(fileName string) error {
	if x.pcapHandle != nil {
		return errors.New("Already set pcap handler, do not specify multiple capture soruce")
	}

	log.Println("read from ", fileName)
	handle, pcapErr := pcap.OpenOffline(fileName)

	if pcapErr != nil {
		return pcapErr
	}

	x.pcapHandle = handle
	return nil
}

func (x *Overseer) SetPcapDev(devName string) error {
	if x.pcapHandle != nil {
		return errors.New("Already set pcap handler, do not specify multiple capture soruce")
	}

	log.Println("capture from ", devName)

	var (
		snapshotLen int32         = 0xffff
		promiscuous bool          = true
		timeout     time.Duration = -1 * time.Second
	)

	handle, pcapErr := pcap.OpenLive(devName, snapshotLen, promiscuous, timeout)

	if pcapErr != nil {
		return pcapErr
	}

	x.pcapHandle = handle
	return nil
}
