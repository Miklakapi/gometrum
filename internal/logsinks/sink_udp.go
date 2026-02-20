package logsinks

import (
	"net"
	"time"
)

type UdpSink struct {
	name  string
	addr  string
	queue *Queue[LogEvent]
}

func NewUdpSink(name, addr string, queueSize int) *UdpSink {
	return &UdpSink{
		name:  name,
		addr:  addr,
		queue: NewQueue[LogEvent](queueSize),
	}
}

func (s *UdpSink) Push(ev LogEvent) bool {
	return s.queue.Push(ev)
}

func (s *UdpSink) Start() {
	go s.loop()
}

func (s *UdpSink) Close() error {
	s.queue.Close()
	return nil
}

func (s *UdpSink) loop() {
	var err error

	dial := func() net.Conn {
		c, err := net.Dial("udp", s.addr)
		if err != nil {
			return nil
		}
		return c
	}

	conn := dial()
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	for ev := range s.queue.Chan() {
		line := encodeTextLine(ev)

		if conn == nil {
			conn = dial()
			if conn == nil {
				time.Sleep(250 * time.Millisecond)
				continue
			}
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			_ = conn.Close()
			conn = nil
			continue
		}
	}
}
