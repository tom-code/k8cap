package main

import (
  "bytes"
  "encoding/binary"
  "sync"
  "time"
)


type PcapHeader struct {
  magic    uint32
  major    uint16
  minor    uint16
  zone     int32
  sigfigs  uint32
  snaplen  uint32
  network  uint32
}

type PcapRecordHeader struct {
  ts_sec   uint32
  ts_usec  uint32
  incl_len uint32
  orig_len uint32
}

type Pcap struct {
  buffer bytes.Buffer
  mutex sync.Mutex
}

func (p *Pcap) init() {
  header := PcapHeader {
    magic: 0xa1b2c3d4,
    major: 2,
    minor: 4,
    zone: 0,
    sigfigs: 0,
    network: 1,
  }
  binary.Write(&p.buffer, binary.BigEndian, &header)
}


func (p *Pcap) write(raw []byte) {
  now := time.Now()
  header := PcapRecordHeader{
  	ts_sec:   uint32(now.Unix()),
  	ts_usec:  uint32(now.UnixNano()/1000 - (now.Unix() * 1000000)),
  	incl_len: uint32(len(raw)),
  	orig_len: uint32(len(raw)),
  }
  p.mutex.Lock()
  defer p.mutex.Unlock()
  binary.Write(&p.buffer, binary.BigEndian, &header)
  p.buffer.Write(raw)
}

func (p *Pcap) get() []byte {
  p.mutex.Lock()
  defer p.mutex.Unlock()
  return p.buffer.Bytes()
}

func (p *Pcap) reset() {
  p.mutex.Lock()
  defer p.mutex.Unlock()
  p.buffer.Reset()
  p.init()
}