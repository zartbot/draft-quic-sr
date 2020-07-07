package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"log"
	"net"

	"github.com/sirupsen/logrus"
)

var (
	port      = flag.String("port", "", "server port")
	debugflag = flag.Bool("debug", false, "debug flags")
)

func main() {

	flag.Parse()
	// check if we have anything

	if "" == *port {
		flag.Usage()
		log.Fatalln("\nremote server is not specified")
	}
	pc, err := net.ListenPacket("udp", ":"+*port)
	if err != nil {
		logrus.Fatal(err)
	}
	for {
		buf := make([]byte, 1500)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}
		go serve(pc, addr, buf[:n])
	}

}

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	//non quic packet
	if buf[0] != 0 {
		return
	}

	srcAddrLoc := uint8(buf[2]) + 4

	af := uint8(buf[3])
	if af == 0 {
		return
	}

	/* update SRHeader field
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	|  SR Type      |   SRH Len     |    LastEntry  |  Segment Left |
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+*/

	var segLen uint8 = 6
	if af == 2 {
		segLen = 18
	}

	srhLoc := srcAddrLoc + segLen
	segmentLeft := uint8(buf[srhLoc+3])
	lastEntry := buf[srhLoc+2]

	if segmentLeft == lastEntry {
		//NAT ALG applied on 1st hop
		segment := StrToByte(addr.String())
		copy(buf[srcAddrLoc:srhLoc], segment)
	}
	//reduce segmentleft
	buf[srhLoc+3]--
	//logrus.Warn("SL|LE:", segmentLeft, "|", lastEntry)
	if segmentLeft >= 1 {
		start := srhLoc + 4 + (segmentLeft-1)*segLen
		dst := buf[start : start+segLen]
		addr1 := ByteToNetAddr(dst)
		//logrus.Warn("start:", start, "|", addr1.String())
		//fmt.Println(hex.Dump(buf))
		pc.WriteTo(buf, addr1)
	}
}

//StrToByte is used convert string uaddr to Net.addr format
func StrToByte(str string) []byte {
	uaddr, err := net.ResolveUDPAddr("udp", str)
	if err != nil {
		return nil
	}
	port := uint16(uaddr.Port)
	ipv4 := uaddr.IP.To4()

	//IPv6 address return [16]Byte Addr + [2]byte Port
	if ipv4 == nil {
		buf := bytes.NewBuffer([]byte(uaddr.IP))
		buf.WriteByte(byte(port >> 8))
		buf.WriteByte(byte(port & 0xff))
		return buf.Bytes()
	}

	//IPv4 address return [4]Byte Addr + [2]byte Port
	buf := bytes.NewBuffer([]byte(ipv4))
	buf.WriteByte(byte(port >> 8))
	buf.WriteByte(byte(port & 0xff))
	return buf.Bytes()
}

//ByteToNetAddr is used to parse segment to net.Addr format
func ByteToNetAddr(b []byte) net.Addr {
	blen := len(b)
	if blen == 6 {
		ip := net.IP(b[0:4])
		port := binary.BigEndian.Uint16(b[4:6])

		return &net.UDPAddr{
			IP:   ip,
			Port: int(port),
		}

	} else if blen == 18 {
		ip := net.IP(b[0:16])
		port := binary.BigEndian.Uint16(b[16:18])

		return &net.UDPAddr{
			IP:   ip,
			Port: int(port),
		}
	}
	return nil
}
