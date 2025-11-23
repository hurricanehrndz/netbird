//go:build dragonfly || freebsd || netbsd || openbsd || darwin

package networkmonitor

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"syscall"
	"unsafe"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/route"
	"golang.org/x/sys/unix"

	"github.com/netbirdio/netbird/client/internal/routemanager/systemops"
)

type nexthopPair struct {
	V4 systemops.Nexthop
	V6 systemops.Nexthop
}

func prepareFd() (int, error) {
	return unix.Socket(syscall.AF_ROUTE, syscall.SOCK_RAW, syscall.AF_UNSPEC)
}

func routeCheck(ctx context.Context, fd int, nexthopv4, nexthopv6 systemops.Nexthop) error {
	expectedNexthops := nexthopPair{V4: nexthopv4, V6: nexthopv6}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			buf := make([]byte, 2048)
			n, err := unix.Read(fd, buf)
			if err != nil {
				if !errors.Is(err, unix.EBADF) && !errors.Is(err, unix.EINVAL) {
					log.Warnf("Network monitor: failed to read from routing socket: %v", err)
				}
				continue
			}
			if n < unix.SizeofRtMsghdr {
				log.Debugf("Network monitor: read from routing socket returned less than expected: %d bytes", n)
				continue
			}

			msg := (*unix.RtMsghdr)(unsafe.Pointer(&buf[0]))

			isRouteChange := msg.Type == syscall.RTM_CHANGE || msg.Type == syscall.RTM_ADD || msg.Type == syscall.RTM_DELETE

			if isRouteChange {
				route, err := parseRouteMessage(buf[:n])
				if err != nil {
					log.Debugf("Network monitor: error parsing routing message: %v", err)
					continue
				}
				isDefaultRoute := route.Dst.Bits() == 0
				if !isDefaultRoute {
					continue
				}
				if hasDefaultRouteChanged(expectedNexthops) {
					return nil
				}
			}
		}
	}
}

func hasDefaultRouteChanged(expectedNexthops nexthopPair) bool {
	actualNexthopv4, errv4 := systemops.GetNextHop(netip.IPv4Unspecified())
	actualNexthopv6, errv6 := systemops.GetNextHop(netip.IPv6Unspecified())
	if errv4 != nil || errv6 != nil {
		err := errors.Join(errv4, errv6)
		log.Infof("Network monitor: failed to check next hop, assuming no network connection available: %s", err)
		return true
	}

	if !expectedNexthops.V4.Equal(actualNexthopv4) || !expectedNexthops.V6.Equal(actualNexthopv6) {
		log.Infof("Network monitor: default route changed v4: %s -> %s", expectedNexthops.V4, actualNexthopv4)
		log.Infof("Network monitor: default route changed v6: %s -> %s", expectedNexthops.V6, actualNexthopv6)
		return true
	}
	return false
}

func parseRouteMessage(buf []byte) (*systemops.Route, error) {
	msgs, err := route.ParseRIB(route.RIBTypeRoute, buf)
	if err != nil {
		return nil, fmt.Errorf("parse RIB: %v", err)
	}

	if len(msgs) != 1 {
		return nil, fmt.Errorf("unexpected RIB message msgs: %v", msgs)
	}

	msg, ok := msgs[0].(*route.RouteMessage)
	if !ok {
		return nil, fmt.Errorf("unexpected RIB message type: %T", msgs[0])
	}

	return systemops.MsgToRoute(msg)
}
