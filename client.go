package gologix

import (
	"bytes"
	"net"
	"sync"
	"time"
)

type Client struct {
	IPAddress     string
	Path          *bytes.Buffer
	SocketTimeout time.Duration

	// you have to change this read sequencer every time you make a new tag request.  If you don't, you
	// won't get an error but it will return the last value you requested again.
	// You don't even have to keep incrementing it.  just going back and forth between 1 and 0 works OK.
	// Use Sequencer() instead of accessing this directly to achieve that.
	sequencerValue uint16

	KnownTags map[string]KnownTag

	Mutex                  sync.Mutex
	Conn                   net.Conn
	SessionHandle          uint32
	OTNetworkConnectionID  uint32
	HeaderSequenceCounter  uint16
	Connected              bool
	ConnectionSize         int
	ConnectionSerialNumber uint16
	Context                uint64 // fun fact - rockwell PLCs don't mind being rickrolled.
}

func (client *Client) Sequencer() uint16 {
	client.sequencerValue++
	return client.sequencerValue
}

type KnownTag struct {
	Name        string
	Type        CIPType
	Class       CIPClass
	Instance    CIPInstance
	Array_Order []int
}

func (t KnownTag) Bytes() []byte {
	ins := CIPInstance(t.Instance)
	b := bytes.Buffer{}
	b.Write(CIPObject_Symbol.Bytes()) // 0x20 0x6B
	b.Write(ins.Bytes())
	return b.Bytes()
}
