package logsinks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func encodeLoki(batch []LogEvent, labels map[string]string) ([]byte, error) {
	values := make([][2]string, 0, len(batch))

	for _, ev := range batch {
		ns := ev.Time.UTC().UnixNano()
		line := encodeTextLine(ev)
		values = append(values, [2]string{strconv.FormatInt(ns, 10), line})
	}

	p := lokiPush{
		Streams: []lokiStream{
			{
				Stream: labels,
				Values: values,
			},
		},
	}

	return json.Marshal(p)
}

func defaultLokiLabels(sinkName string) map[string]string {
	return map[string]string{
		"job":  "gometrum",
		"sink": sinkName,
	}
}

func encodeEventJSON(ev LogEvent) ([]byte, error) {
	return json.Marshal(ev)
}

func encodeNDJSON(batch []LogEvent) ([]byte, error) {
	var buf bytes.Buffer
	for i := range batch {
		b, err := json.Marshal(batch[i])
		if err != nil {
			return nil, err
		}
		buf.Write(b)
		buf.WriteByte('\n')
	}
	return buf.Bytes(), nil
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
