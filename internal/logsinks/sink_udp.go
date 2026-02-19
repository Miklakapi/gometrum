package logsinks

import (
	"fmt"
	"net"
	"sort"
	"strings"
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

func encodeTextLine(ev LogEvent) string {
	var b strings.Builder

	ts := ev.Time.Format(time.RFC3339Nano)
	b.WriteString(ts)
	b.WriteString(" ")
	b.WriteString(ev.Level)
	b.WriteString(" ")
	b.WriteString(ev.Msg)

	if len(ev.Attrs) > 0 {
		keys := make([]string, 0, len(ev.Attrs))
		for k := range ev.Attrs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			b.WriteString(" ")
			b.WriteString(k)
			b.WriteString("=")
			fmt.Fprint(&b, ev.Attrs[k])
		}
	}

	return b.String()
}
