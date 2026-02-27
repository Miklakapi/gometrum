package mqtt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type DryRunClient struct {
	availabilityTopic string
	onlinePayload     []byte
}

func NewDryRun() *DryRunClient {
	return &DryRunClient{}
}

func (d *DryRunClient) SetAvailability(topic string, onlinePayload []byte) {
	d.availabilityTopic = topic
	d.onlinePayload = onlinePayload
}

func (d *DryRunClient) Connect(timeout time.Duration) error {
	fmt.Println("DRY-RUN: connect skipped")

	if d.availabilityTopic != "" && len(d.onlinePayload) > 0 {
		fmt.Printf("AVAILABILITY %s => %s\n", d.availabilityTopic, formatPayload(d.onlinePayload))
	}

	return nil
}

func (d *DryRunClient) Publish(topic string, qos byte, retain bool, payload []byte) error {
	formatted := formatPayload(payload)

	if stringsHasNewline(formatted) {
		fmt.Printf("PUBLISH %s =>\n%s\n", topic, formatted)
		return nil
	}

	fmt.Printf("PUBLISH %s => %s\n", topic, formatted)
	return nil
}

func (d *DryRunClient) Subscribe(topic string, qos byte, handler func(topic string, payload []byte)) error {
	fmt.Printf("SUBSCRIBE %s\n", topic)
	return nil
}

func (d *DryRunClient) Close() {
	fmt.Println("DRY-RUN: close skipped")

	if d.availabilityTopic != "" {
		fmt.Printf("AVAILABILITY %s => offline\n", d.availabilityTopic)
	}
}

func formatPayload(b []byte) string {
	if len(b) == 0 {
		return "<empty>"
	}

	trimmed := bytes.TrimSpace(b)
	if len(trimmed) == 0 {
		return "<empty>"
	}

	if trimmed[0] == '{' || trimmed[0] == '[' {
		var v any
		if err := json.Unmarshal(trimmed, &v); err == nil {
			pretty, err := json.MarshalIndent(v, "", "  ")
			if err == nil {
				return string(pretty)
			}
		}
	}

	return string(trimmed)
}

func stringsHasNewline(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			return true
		}
	}
	return false
}
