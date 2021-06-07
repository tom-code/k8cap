package main

import (
  "encoding/hex"
  "fmt"
  "log"
  "runtime"
  "syscall"

  "github.com/containernetworking/plugins/pkg/ns"
  "github.com/google/gopacket"
  "github.com/google/gopacket/layers"
)

var pcap Pcap

type Decoder struct {
  eth layers.Ethernet
  ip4 layers.IPv4
  ip6 layers.IPv6
  tcp layers.TCP
  udp layers.UDP
  payload gopacket.Payload
  parser *gopacket.DecodingLayerParser
}

func (d *Decoder) init() {
  d.parser = gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet,
                                             &d.eth, &d.ip4, &d.ip6, &d.tcp, &d.udp, &d.payload)
}

func (d *Decoder) decode(raw []byte) {
  decoded := []gopacket.LayerType{}
  err := d.parser.DecodeLayers(raw, &decoded)
  if err != nil {
    log.Println(err)
    return
  }
  log.Printf("%s->%s %s->%s %d->%d %s", d.eth.SrcMAC, d.eth.DstMAC, d.ip4.SrcIP, d.ip4.DstIP,
                d.tcp.SrcPort, d.tcp.DstPort, hex.EncodeToString(d.payload.LayerContents()))
}

func captureWorker() error {
  fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, 0x300)
  if err != nil {
    return err
  }
  defer syscall.Close(fd)

  decoder := new(Decoder)
  decoder.init()
  for {
    buffer := make([]byte, 1024*64)
    n, _, err := syscall.Recvfrom(fd, buffer, 0)
    if err != nil {
      log.Println(err)
      continue
    }
    pcap.write(buffer)
    hex := hex.EncodeToString(buffer[:n])
    log.Println(hex)
    decoder.decode(buffer[:n])
  }
}

func capture(pid int) {
  pcap.init()
  runtime.LockOSThread()
  defer runtime.UnlockOSThread()
  nspath := fmt.Sprintf("/proc/%d/ns/net", pid)
  nsx, err := ns.GetNS(nspath)
  if err != nil {
    log.Println(err)
    return
  }
  nsx.Do(func(_ ns.NetNS) error {
    captureWorker()
    return nil
  })
}