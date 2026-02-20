package logsinks

import (
	"bytes"
	"context"
	"net/http"
	"strings"
	"time"
)

type HttpSink struct {
	name    string
	url     string
	method  string
	timeout time.Duration
	headers map[string]string
	codec   HttpCodec

	queue *Queue[LogEvent]
	batch HttpBatch
}

type HttpBatch struct {
	MaxItems int
	MaxWait  time.Duration
}

func NewHttpSink(
	name string,
	url string,
	method string,
	timeout time.Duration,
	headers map[string]string,
	codec HttpCodec,
	queueSize int,
	batch HttpBatch,
) *HttpSink {
	if method == "" {
		method = http.MethodPost
	}
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	if batch.MaxItems <= 0 {
		batch.MaxItems = 20
	}
	if batch.MaxWait <= 0 {
		batch.MaxWait = 1 * time.Second
	}

	return &HttpSink{
		name:    name,
		url:     url,
		method:  method,
		timeout: timeout,
		headers: headers,
		codec:   codec,
		queue:   NewQueue[LogEvent](queueSize),
		batch:   batch,
	}
}

func (s *HttpSink) Push(ev LogEvent) bool {
	return s.queue.Push(ev)
}

func (s *HttpSink) Start() {
	go s.loop()
}

func (s *HttpSink) Close() error {
	s.queue.Close()
	return nil
}

func (s *HttpSink) loop() {
	client := &http.Client{Timeout: s.timeout}

	batchable := s.codec == CodecNDJSON || s.codec == CodecLoki
	if !batchable {
		for ev := range s.queue.Chan() {
			_ = s.flush(client, []LogEvent{ev})
		}
		return
	}

	batch := make([]LogEvent, 0, s.batch.MaxItems)

	timer := time.NewTimer(s.batch.MaxWait)
	defer timer.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		_ = s.flush(client, batch)
		batch = batch[:0]
	}

	for {
		select {
		case ev, ok := <-s.queue.Chan():
			if !ok {
				flush()
				return
			}

			batch = append(batch, ev)

			if len(batch) >= s.batch.MaxItems {
				flush()
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(s.batch.MaxWait)
			}

		case <-timer.C:
			flush()
			timer.Reset(s.batch.MaxWait)
		}
	}
}

func (s *HttpSink) flush(client *http.Client, batch []LogEvent) error {
	if len(batch) == 0 {
		return nil
	}

	var body []byte
	var err error
	var contentType string

	switch s.codec {
	case CodecText:
		body = []byte(encodeTextLine(batch[0]))
		contentType = "text/plain; charset=utf-8"

	case CodecEventJSON:
		body, err = encodeEventJSON(batch[0])
		if err != nil {
			return err
		}
		contentType = "application/json"

	case CodecNDJSON:
		body, err = encodeNDJSON(batch)
		if err != nil {
			return err
		}
		contentType = "application/x-ndjson"

	case CodecLoki:
		labels := defaultLokiLabels(s.name)
		body, err = encodeLoki(batch, labels)
		if err != nil {
			return err
		}
		contentType = "application/json"

	default:
		body, err = encodeNDJSON(batch)
		if err != nil {
			return err
		}
		contentType = "application/x-ndjson"
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, s.method, s.url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)
	for k, v := range s.headers {
		if strings.TrimSpace(k) == "" {
			continue
		}
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	return nil
}
