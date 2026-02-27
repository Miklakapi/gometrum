package buttons

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/Miklakapi/gometrum/internal/config"
)

type Button interface {
	Key() string
	Name() string
	Command() []string
	Timeout() time.Duration
	Icon() string

	Execute(ctx context.Context) ([]byte, error)
}

type base struct {
	key     string
	name    string
	command []string
	timeout time.Duration
	icon    string
}

func (b base) Key() string            { return b.key }
func (b base) Name() string           { return b.name }
func (b base) Command() []string      { return b.command }
func (b base) Timeout() time.Duration { return b.timeout }
func (b base) Icon() string           { return b.icon }

type commandButton struct {
	base
}

func Build(cfg config.Config) ([]Button, error) {
	out := make([]Button, 0, len(cfg.Buttons))
	for key, bc := range cfg.Buttons {
		out = append(out, newCommandButton(key, bc))
	}
	return out, nil
}

func (b *commandButton) Execute(parent context.Context) ([]byte, error) {
	ctx, cancel := context.WithTimeout(parent, b.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, b.command[0], b.command[1:]...)
	out, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return out, fmt.Errorf("timeout exceeded (%s)", b.timeout)
	}

	return out, err
}

func newCommandButton(key string, cfg config.ButtonConfig) Button {
	return &commandButton{
		base: base{
			key:     key,
			name:    cfg.Name,
			command: cfg.Command,
			timeout: cfg.Timeout,
			icon:    cfg.HA.Icon,
		},
	}
}
