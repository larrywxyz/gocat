package relay

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/palantir/stacktrace"

	"gocat/internal/relay"
)

func NewTCPToVsockRelay(tcpToVsockPort, tcpToVsockAddress string) error {
	tcpToUnixHealthCheckInterval := 30 * time.Second
	bufferSize := 16384

	relayer, err := relay.NewTCPtoVsock(
		tcpToUnixHealthCheckInterval,
		tcpToVsockPort,
		tcpToVsockAddress,
		bufferSize,
	)
	if err != nil {
		return stacktrace.Propagate(err, "couldn't create relay from TCP to vsock")
	}

	osSignalCh := make(chan os.Signal, 1)
	defer close(osSignalCh)

	signal.Notify(osSignalCh, os.Interrupt, syscall.SIGTERM)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// Ctrl+C handler
	go func() {
		<-osSignalCh
		signal.Stop(osSignalCh)
		cancelFunc()
	}()

	err = relayer.Relay(ctx)
	if err != nil {
		return stacktrace.Propagate(err, "couldn't relay from TCP to vsock")
	}

	return nil
}
