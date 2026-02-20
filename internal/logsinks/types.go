package logsinks

import "time"

type HttpCodec string

const (
	CodecEventJSON HttpCodec = "event_json"
	CodecText      HttpCodec = "text"
	CodecNDJSON    HttpCodec = "ndjson"
	CodecLoki      HttpCodec = "loki"
)

type LogEvent struct {
	Time  time.Time      `json:"time"`
	Level string         `json:"level"`
	Msg   string         `json:"msg"`
	Attrs map[string]any `json:"attrs,omitempty"`
}

type lokiPush struct {
	Streams []lokiStream `json:"streams"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"`
}
