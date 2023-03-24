package dump

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"strings"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Peer struct {
	PublicKey                   wgtypes.Key
	PresharedKey                wgtypes.Key
	Endpoint                    *net.UDPAddr
	AllowedIPs                  []net.IPNet
	LastHandshakeTime           time.Time
	ReceiveBytes                uint64
	TransmitBytes               uint64
	PersistentKeepaliveInterval time.Duration
}

func Parse(raw []byte) ([]Peer, error) {
	var out []Peer

	scanner := bufio.NewScanner(bytes.NewBuffer(raw))
	for scanner.Scan() {
		splits := strings.SplitN(scanner.Text(), "\t", 8)
		if len(splits) != 8 {
			return nil, errors.New("invalid line")
		}

		raw := splits[0]
		publicKey, err := wgtypes.ParseKey(raw)
		if nil != err {
			return nil, fmt.Errorf("invalid public key (%s): %v", raw, err)
		}

		raw = splits[1]
		presharedKey, err := wgtypes.ParseKey(raw)
		if nil != err {
			return nil, fmt.Errorf("invalid preshared key (%s): %v", raw, err)
		}

		endpoint, err := func() (*net.UDPAddr, error) {
			raw := splits[2]
			if raw == "(none)" {
				return nil, nil
			}

			addrPort, err := netip.ParseAddrPort(raw)
			if nil != err {
				return nil, fmt.Errorf("invalid endpoint (%s): %v", raw, err)
			}
			endpoint := net.UDPAddrFromAddrPort(addrPort)
			return endpoint, nil
		}()
		if nil != err {
			return nil, err
		}

		var allowedIPs []net.IPNet
		raw = splits[3]
		rawAllowedIPs := strings.Split(raw, ",")
		for _, v := range rawAllowedIPs {
			_, ipnet, err := net.ParseCIDR(v)
			if nil != err {
				return nil, fmt.Errorf("invalid allowed ip (%s): %v", v, err)
			}

			allowedIPs = append(allowedIPs, *ipnet)
		}

		raw = splits[4]
		var latestHandshakeTime time.Time
		if raw != "0" {
			i, err := strconv.ParseInt(raw, 10, 64)
			if nil != err {
				return nil, fmt.Errorf("invalid latest handshake unix timestamp (%s): %v", raw, err)
			}
			if i < 0 {
				return nil, fmt.Errorf("unexpected negative latest handshake unix timestamp value (%s)", raw)
			}
			latestHandshakeTime = time.Unix(i, 0)
		}

		raw = splits[5]
		receiveBytes, err := strconv.ParseInt(raw, 10, 64)
		if nil != err {
			return nil, fmt.Errorf("invalid number of received bytes (%s): %v", raw, err)
		}

		raw = splits[6]
		transmitBytes, err := strconv.ParseInt(raw, 10, 64)
		if nil != err {
			return nil, fmt.Errorf("invalid number of transmitted bytes (%s): %v", raw, err)
		}

		persistentKeepaliveInterval, err := func() (time.Duration, error) {
			raw := splits[7]
			if raw == "off" {
				return 0, nil
			}

			i, err := strconv.ParseInt(raw, 10, 32)
			if nil != err {
				return 0, fmt.Errorf("invalid persistent keepalive interval (%s): %v", raw, err)
			}
			if i < 0 {
				return 0, fmt.Errorf("unexpected negative persistent keepalive value (%s)", raw)
			}

			return time.Duration(i), nil
		}()
		if nil != err {
			return nil, err
		}

		out = append(out, Peer{
			PublicKey:                   publicKey,
			PresharedKey:                presharedKey,
			Endpoint:                    endpoint,
			AllowedIPs:                  allowedIPs,
			LastHandshakeTime:           latestHandshakeTime,
			ReceiveBytes:                uint64(receiveBytes),
			TransmitBytes:               uint64(transmitBytes),
			PersistentKeepaliveInterval: time.Duration(persistentKeepaliveInterval),
		})
	}

	return out, nil
}
