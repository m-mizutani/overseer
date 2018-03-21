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
	ssnTable     *SessionTable
}

func New() Overseer {
	seer := Overseer{}
	seer.ssnTable = NewSessionTable()
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

func (x *Overseer) Preprocess() error {
	if x.pcapHandle == nil {
		return errors.New("No available device or pcap file, need to specify one of them")
	}

	if x.packetSource == nil {
		x.packetSource = gopacket.NewPacketSource(x.pcapHandle, x.pcapHandle.LinkType())
	}

	return nil
}

func (x *Overseer) Loop() error {
	return x.Read(0)
}

func (x *Overseer) Read(readCount int) error {
	err := x.Preprocess()
	if err != nil {
		return err
	}

	count := 0
	const timeout time.Duration = time.Second
	packets := x.packetSource.Packets()
	ticker := time.Tick(timeout)

	for {
		select {
		case packet := <-packets:
			if packet == nil {
				return nil
			}

			x.ssnTable.ReadPacket(packet)
			count += 1

			if readCount > 0 && count >= readCount {
				return nil
			}

		case <-ticker:
			log.Println("tick")
		}
	}
}

func (x *Overseer) Close() {
	if x.pcapHandle != nil {
		x.pcapHandle.Close()
	}
}
