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

func NewHvsockToTCPRelay(hvsockToTcpPath, hvsockToTcpAddress string) error {
	unixToTCPHealthCheckDuration := 30 * time.Second
	bufferSize := 16384

	relayer, err := relay.NewHvsockTcp(
		unixToTCPHealthCheckDuration,
		hvsockToTcpPath,
		hvsockToTcpAddress,
		bufferSize,
	)
	if err != nil {
		return stacktrace.Propagate(err, "couldn't create relay from hvsock to TCP")
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
		return stacktrace.Propagate(err, "couldn't relay from hvsock to TCP")
	}
	return nil
}
