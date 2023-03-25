package dump_test

import (
	"bufio"
	"bytes"
	_ "embed"
	"net"
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"github.com/xeptore/wireuse/pkg/dump"
)

var (
	//go:embed testdata/dump.txt
	rawValidDump []byte
	//go:embed testdata/dump-invalid.txt
	rawInvalidDump []byte
)

func mustParseKey(t *testing.T, s string) wgtypes.Key {
	key, err := wgtypes.ParseKey(s)
	require.Nil(t, err)

	return key
}

func mustParseAllowedIPs(t *testing.T, s ...string) []net.IPNet {
	out := make([]net.IPNet, len(s))
	for i := 0; i < len(s); i++ {
		_, ipNet, err := net.ParseCIDR(s[i])
		require.Nil(t, err)
		out[i] = *ipNet
	}

	return out
}

func mustParseUDPAddr(t *testing.T, s string) *net.UDPAddr {
	addrPort, err := netip.ParseAddrPort(s)
	require.Nil(t, err)

	return net.UDPAddrFromAddrPort(addrPort)
}

func mustParseUnixTime(t *testing.T, u int64) time.Time {
	return time.Unix(u, 0)
}

func TestParseValidDumpFile(t *testing.T) {
	peers, err := dump.Parse(bytes.NewBuffer(rawValidDump))
	require.Nil(t, err)
	require.NotNil(t, peers)
	require.Equal(t, []dump.Peer{
		{
			PublicKey:                   mustParseKey(t, "yLk+kQwd2Zm4I6f6Aeg/OGlT3GLtHsNmY+kUXK2sfD0="),
			PresharedKey:                mustParseKey(t, "QMXzIPq+o5zNvecAIQZn61rB9SmolhW0fgj7Zwylzvg="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.19/32", "fdd0:438e:19ba:5069::13/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "XuyQmI/KXTBvwdpfjUzULj4PfBW15mrMPXm8pS1OnhQ="),
			PresharedKey:                mustParseKey(t, "kT0tuNL/rhPoWjWiQhSCHFJuaSKLxlqg2aAuxwwpz4I="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.20/32", "fdd0:438e:19ba:5069::14/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "UDrZG/thxkzPbow3Lr1gGHYeKbNlQF8Znob9hujRYjo="),
			PresharedKey:                mustParseKey(t, "fkKqj8RmLjHbo8T63OtLn6JM8MnCIp+LZLpvxiIl2HU="),
			Endpoint:                    mustParseUDPAddr(t, "83.123.9.201:42944"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.21/32", "fdd0:438e:19ba:5069::15/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679600629),
			ReceiveBytes:                7138444,
			TransmitBytes:               103247936,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "qwQBlkUITgLsddSi9kbKM14HCkvqmVYi1VT/8n9u6kY="),
			PresharedKey:                mustParseKey(t, "MFToe6QumdelJu3whvgSDniVNM+jmuTNMXP9FY07F8g="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.29/32", "fdd0:438e:19ba:5069::1d/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "5jNQhl3WZJBOI4+M5e4Pk7+9SCMSg6GyG+Dkc+wpCVA="),
			PresharedKey:                mustParseKey(t, "l7itiN0Pr0fQMD7Ae8EMpDHmQ0rL4KyTNuL/5J9+O6o="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.30/32", "fdd0:438e:19ba:5069::1e/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "OafCe4tLL8rShnOiQod35G+NdDPQhwIN4JKRoyKbNkI="),
			PresharedKey:                mustParseKey(t, "mJBHO390BWD3GxHMIdWFey70N7c6x5EDrVuzcROh/Ao="),
			Endpoint:                    mustParseUDPAddr(t, "5.125.49.118:32738"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.41/32", "fdd0:438e:19ba:5069::29/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679608693),
			ReceiveBytes:                69519144,
			TransmitBytes:               475482968,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "R86vOf/VdN9m0Ey4AD80D4i5UdTefMsKJw/+BuOSeS4="),
			PresharedKey:                mustParseKey(t, "+V9iBlvEB7FCZcUMcY40c1FazbleACJY/qButRjkUUg="),
			Endpoint:                    mustParseUDPAddr(t, "95.38.119.159:52491"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.42/32", "fdd0:438e:19ba:5069::2a/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679597842),
			ReceiveBytes:                8385684,
			TransmitBytes:               188015208,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "FaibxGztHHy3pm0kiiBfI6GMM6UpdiD6Kinekh0ARDk="),
			PresharedKey:                mustParseKey(t, "8R5QEYDbJw/VYsNw62LX9p+RGtn3NKeVyDeCol1zzT8="),
			Endpoint:                    mustParseUDPAddr(t, "188.229.11.180:62680"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.43/32", "fdd0:438e:19ba:5069::2b/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679610586),
			ReceiveBytes:                91786228,
			TransmitBytes:               109298580,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "gOZgWvD4pVpGNx8ZBJOrEFbUvdOmIqyJ5xaWkb0QkwY="),
			PresharedKey:                mustParseKey(t, "4fpiAvyjMmUPGDsi5IV7EPWmq9ANOCsuBU81s1KKvwg="),
			Endpoint:                    mustParseUDPAddr(t, "5.123.27.254:10852"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.44/32", "fdd0:438e:19ba:5069::2c/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679622245),
			ReceiveBytes:                38969328,
			TransmitBytes:               500779248,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "nRhGvfiDIdUts1yJn4Bt2QVURBK73ks++SrqdseUZjA="),
			PresharedKey:                mustParseKey(t, "DK+KyyjvtAHARlSkvH6ogNYeCb//GSVfp/SaYGCkTvc="),
			Endpoint:                    mustParseUDPAddr(t, "83.122.114.201:53992"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.45/32", "fdd0:438e:19ba:5069::2d/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679614805),
			ReceiveBytes:                49627932,
			TransmitBytes:               1698468224,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "cM6cjxjXEfbV/MmdIaqnI3iG6NX0eTGrxi2xdr8+XlA="),
			PresharedKey:                mustParseKey(t, "yHHspi9idt69QY3EawOQrEVdi2kMQ4CdW4jadVHDWHo="),
			Endpoint:                    mustParseUDPAddr(t, "5.126.154.201:32527"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.46/32", "fdd0:438e:19ba:5069::2e/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679622324),
			ReceiveBytes:                41659912,
			TransmitBytes:               387348224,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "gsqgOUyV3wl2QjGemvO73O1MD9PB/umq6lH18Px/1RU="),
			PresharedKey:                mustParseKey(t, "+NjgCZQ3224mpbDZUgIj24mzX6NcovotLuDOqjTHLw8="),
			Endpoint:                    mustParseUDPAddr(t, "5.125.223.142:36784"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.47/32", "fdd0:438e:19ba:5069::2f/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679611750),
			ReceiveBytes:                60529560,
			TransmitBytes:               1421720032,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "mZNGrAccdCJ9a/7UOevRHKxjxOy7fPvgW3ev+bcROWk="),
			PresharedKey:                mustParseKey(t, "/Km1HGiYK+O/jJJJWdPZIEC6pNBzsghhYm9hsc61xFg="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.48/32", "fdd0:438e:19ba:5069::30/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "mYDlBMbMZKEjcX8nTjYWmOGQlguyLomQVm+q9yc9eAM="),
			PresharedKey:                mustParseKey(t, "MXSRqHNEMM48WGXqyUmCC2YxsueGZvyxIW8UnwFibt0="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.49/32", "fdd0:438e:19ba:5069::31/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "P7wDzaPcNHAfR2MVuQV//qD+YfXHR3Ad/7uLm61xUXA="),
			PresharedKey:                mustParseKey(t, "VuYrRbWESPD7kFpXdSs6BJ7fMppBM4IBxOURx/lYv9Q="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.50/32", "fdd0:438e:19ba:5069::32/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "rauZLam+acOP+/mC0C/vAjIweVoWZNQ0jK5iBmpAT0w="),
			PresharedKey:                mustParseKey(t, "dQQA+bYhLgtXwnPug9rQVok1R+ud9RRAqKgo04K/ljU="),
			Endpoint:                    mustParseUDPAddr(t, "5.119.75.31:47338"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.51/32", "fdd0:438e:19ba:5069::33/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679622318),
			ReceiveBytes:                49183124,
			TransmitBytes:               946856032,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "xwlGtnl3nmV+uR+nN8No7mMQiv1OTTOKMdgrwaDBMFk="),
			PresharedKey:                mustParseKey(t, "NKeDh2sJkweEbC7ZspA5fPBL+9ahk/eqqKFcl0IemBg="),
			Endpoint:                    mustParseUDPAddr(t, "37.129.244.184:43698"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.52/32", "fdd0:438e:19ba:5069::34/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679613152),
			ReceiveBytes:                57819876,
			TransmitBytes:               687229524,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "VF3fBaIiqEAlPdYWQO/XNy1cXxAzJOnHRWgJh4knglw="),
			PresharedKey:                mustParseKey(t, "SBX4QINGYs7wCF0GCIyqt1doc50FXgbQ5CYQz2X0YZk="),
			Endpoint:                    mustParseUDPAddr(t, "5.115.179.111:36227"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.54/32", "fdd0:438e:19ba:5069::36/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679622359),
			ReceiveBytes:                147183980,
			TransmitBytes:               1465251556,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "4ZhiIszWe4sPsnCf+bHDs+O2pM+SSDI92urE9gO1AyY="),
			PresharedKey:                mustParseKey(t, "oxOYXuZsIZ1cvmN/r61l63XG7I4/ZeqJQk1Eizmslps="),
			Endpoint:                    mustParseUDPAddr(t, "5.125.254.72:63259"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.55/32", "fdd0:438e:19ba:5069::37/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679604971),
			ReceiveBytes:                30436016,
			TransmitBytes:               800796928,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "dwwb9pB9b43e2mDF1OAF4ViKbRy8MhNog2/5zfsbHxM="),
			PresharedKey:                mustParseKey(t, "T+uhSXH45y+SAB/74swecdmKwbzoQkAk3QubTCQU4Xs="),
			Endpoint:                    mustParseUDPAddr(t, "5.126.41.5:17955"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.56/32", "fdd0:438e:19ba:5069::38/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679622289),
			ReceiveBytes:                193450932,
			TransmitBytes:               1615951448,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "Sgz2K9TQGaDPPuMPgjdPnlrKCutJqUPkf6GrO9FdAgQ="),
			PresharedKey:                mustParseKey(t, "aFkkiUHklssOZ0oCGpG+Oj4Oj/P0FWf5Hei8GKRlCGc="),
			Endpoint:                    mustParseUDPAddr(t, "5.124.107.243:18896"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.57/32", "fdd0:438e:19ba:5069::39/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679621178),
			ReceiveBytes:                47751636,
			TransmitBytes:               323048336,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "g0KUtRZ0tJG6YfDgVzfoY6pkmbe1fkeqsp8qu0CrjFk="),
			PresharedKey:                mustParseKey(t, "jg9Hor8qi3AOkBeRDV51l19Z0UQU+Ytrcap7Zb9v/9c="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.59/32", "fdd0:438e:19ba:5069::3b/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "sfIcIpMqJ1+RL34PDS56Y4m2f4P5Qj+KzhiWtJCjDVs="),
			PresharedKey:                mustParseKey(t, "DMcD4y+fKUyBEXW/DzZFBcr28hpcd2GSqMZQvg7Fmiw="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.60/32", "fdd0:438e:19ba:5069::3c/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "yNnsydGMEqEa3M3pi+cXtuWHlHKP7s5NMGReU+Jt81Q="),
			PresharedKey:                mustParseKey(t, "ZyJ35QWf+SKrnsLjHEKCUECLH6VPbwn5p+av3bFjJLE="),
			Endpoint:                    mustParseUDPAddr(t, "5.125.201.34:11978"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.61/32", "fdd0:438e:19ba:5069::3d/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679608317),
			ReceiveBytes:                9238972,
			TransmitBytes:               141261976,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "bpvh8xyCKBZ5wjQRTKI7UxAWAFYXscV2ffsDhWyB+Ek="),
			PresharedKey:                mustParseKey(t, "1ZWVE8w0bk6q6Md4jhNmp8kDaadN+QBZZ8pZSZd/JjY="),
			Endpoint:                    mustParseUDPAddr(t, "37.129.171.116:54772"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.62/32", "fdd0:438e:19ba:5069::3e/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679602156),
			ReceiveBytes:                25262308,
			TransmitBytes:               758835208,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "YESGUd6GzN4sldhXn80Nz3UeX/Tv1nANyTW98gHQTEY="),
			PresharedKey:                mustParseKey(t, "tXKZHknajC+JeytG8BIPjyir0D+fyctLz9CVG7FwZGU="),
			Endpoint:                    mustParseUDPAddr(t, "5.232.199.164:18016"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.64/32", "fdd0:438e:19ba:5069::40/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679609091),
			ReceiveBytes:                53371808,
			TransmitBytes:               828476932,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "pw8kkrzQ762vv3fp3z0DSMpNtmLwRYj0m/iyzVMaYXM="),
			PresharedKey:                mustParseKey(t, "m3CxYl7p2snW00B+VIx/vLDG/le93DwfWMPnPC9EpL4="),
			Endpoint:                    mustParseUDPAddr(t, "83.122.106.35:33994"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.65/32", "fdd0:438e:19ba:5069::41/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679620210),
			ReceiveBytes:                54243332,
			TransmitBytes:               1026011872,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "538nb9hYbOk+idgQ6qfY4qigMY4lKRZyDKkcaQ1E0mw="),
			PresharedKey:                mustParseKey(t, "StqYgoU+2ckQvK4AVVQYHQA5GJvj/KbMRbQM6P62I1A="),
			Endpoint:                    mustParseUDPAddr(t, "5.124.60.6:59688"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.67/32", "fdd0:438e:19ba:5069::43/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679622348),
			ReceiveBytes:                38488680,
			TransmitBytes:               694107612,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "Zc7k5xA70Nfbo2nLJziFFBVY/2OZu5HTsGLQH+KUBRU="),
			PresharedKey:                mustParseKey(t, "bMO4+pL4mnQjLpPSgXNDUUMhjtXmQAyHqc1sZadb+S0="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.68/32", "fdd0:438e:19ba:5069::44/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "t4LWLF/J7BpPtS8p0WX8+oGpQoIO31eE2Uwm2uCoOBo="),
			PresharedKey:                mustParseKey(t, "OpTNAEpQHp+C0SM6AMSr3/keabLhdq3vI9KWZNU8Htc="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.69/32", "fdd0:438e:19ba:5069::45/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "KgIUwSPETcpI5f1tL6Zcsvvpr/sCbaCVVXeWgdlupAM="),
			PresharedKey:                mustParseKey(t, "RU+OfbsJnOqzLeQkwmdrMYgmkv43XKJ52bhsYK7yXdU="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.70/32", "fdd0:438e:19ba:5069::46/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "Ik8ivyv+wRAagtGRJP/kjBIfj8Gp0HKH4YW1c886THw="),
			PresharedKey:                mustParseKey(t, "tSWKQWD3632XTURXkqRLt18iTz2RxK2L617S04b08Mg="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.71/32", "fdd0:438e:19ba:5069::47/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "9chK4Nm6kokH7DKLBkwDtiu1ohUi2IKjHvUQ7ryYilQ="),
			PresharedKey:                mustParseKey(t, "8qCqx2Jc8b41bqJgsUbceIFPZK+tiLscjtDoE7wQ7WA="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.72/32", "fdd0:438e:19ba:5069::48/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "3uXe1350QPqYkCFti5P1p8QNb9JG5Q1aiMWTvCqSai0="),
			PresharedKey:                mustParseKey(t, "OnFaqQUryKomwWEUU4YNxMshiZTY006VzS/8HSahEaE="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.73/32", "fdd0:438e:19ba:5069::49/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "htNim8Y//7BuWZcMW9RmEDl4jrQrfx3/l4GOmp6vLk4="),
			PresharedKey:                mustParseKey(t, "eW+z9QOsTGzrIQG8mnqUw54j0C5mWkY7LqgvwIbvkDM="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.74/32", "fdd0:438e:19ba:5069::4a/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "4rBI2mkO/tXLLbAgerKzEllf5AnmKQO0Yaz9NZBnwE4="),
			PresharedKey:                mustParseKey(t, "Xw1IoyiR1q9g9jgkgRR+0AgDQbYnjCuOfxjz/b4uLdQ="),
			Endpoint:                    mustParseUDPAddr(t, "2.181.107.248:10574"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.76/32", "fdd0:438e:19ba:5069::4c/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679620078),
			ReceiveBytes:                60474420,
			TransmitBytes:               1988225048,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "qvbdv+lPHuGDWUjrojSJSIHY/tjg8e3N3NTFAPp7zDY="),
			PresharedKey:                mustParseKey(t, "RktzPDRjLXNdytEPSXy0cVme2K2u0GPzzVU8SbzXrDU="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.78/32", "fdd0:438e:19ba:5069::4e/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "QZKGqgNRF2cS8g5qS3wcx+/CUpMQOCjAqsOKRaoiNS4="),
			PresharedKey:                mustParseKey(t, "kjZdEwHD4eV8je1EdypnRTs1ipZxiseYQe2cpE/SzdM="),
			Endpoint:                    mustParseUDPAddr(t, "5.125.118.44:7585"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.79/32", "fdd0:438e:19ba:5069::4f/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679612230),
			ReceiveBytes:                58485180,
			TransmitBytes:               549983952,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "tcaDGgjy0/MHmnNaBpYV10vEb+HIE8A66y1q63XtKD8="),
			PresharedKey:                mustParseKey(t, "qxXReWqCnJOfEreDRSiggrUBAsvSL8wEgnZPiV4Lmhk="),
			Endpoint:                    mustParseUDPAddr(t, "5.125.241.99:14167"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.80/32", "fdd0:438e:19ba:5069::50/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679621917),
			ReceiveBytes:                91208368,
			TransmitBytes:               1594288628,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "SxSm/S+Sr/VSptW0O2WTCRLjz6XUFvC6KvAWdHwjMB4="),
			PresharedKey:                mustParseKey(t, "Es9JNg13wuC5YvEchEvT6bS281ozKM8nzdnUu3XBOkU="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.81/32", "fdd0:438e:19ba:5069::51/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "XWMRIDUA0SjDF6pUFPOpa4ffgkCW2luMGGnK2nuL3X8="),
			PresharedKey:                mustParseKey(t, "ZLU0MTM7of2f244My92DY8LfCLyg0BMQQSa+992NlZM="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.82/32", "fdd0:438e:19ba:5069::52/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "Ol/ET8HoIOqK6DvisCbPDUS2WWXlBarKmpZ0y+dsCys="),
			PresharedKey:                mustParseKey(t, "fh8U59aCSPtj7+rXT87hbnQjwRiRKpip0DIfXEFFL9I="),
			Endpoint:                    mustParseUDPAddr(t, "83.123.96.52:49644"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.83/32", "fdd0:438e:19ba:5069::53/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679602712),
			ReceiveBytes:                57556000,
			TransmitBytes:               864674044,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "9llEVDcvyhKrsKUoKkfXJuigLBcRl4U/tPtcUfVMdyU="),
			PresharedKey:                mustParseKey(t, "PfZX2T/KiVC2fRiZ+JzvXUppa50UQ/6mZiA8W27mIgg="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.84/32", "fdd0:438e:19ba:5069::54/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "L7Cmw+UIFKbT/ATkhQjPTrZzp4yphHuv+0oYW9EZdHk="),
			PresharedKey:                mustParseKey(t, "/eM2JPyOP5rEejZYP8iXZesrnlkD7S5Gc9MQuW9Lit8="),
			Endpoint:                    mustParseUDPAddr(t, "5.124.231.78:17780"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.85/32", "fdd0:438e:19ba:5069::55/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679622312),
			ReceiveBytes:                96118376,
			TransmitBytes:               2048227908,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "HQGRHBO8qFMwRs6R4VIcM7zrNRJH/hO19vFu3gWVyk8="),
			PresharedKey:                mustParseKey(t, "d+O833BjLsxc2WwBmGsdcB8ogDtBYGaXn57xTrJGaAw="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.86/32", "fdd0:438e:19ba:5069::56/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "sPgKlh9EnLc1UKG93ibfpIAovwpn8YjYlOBbQT/sjSM="),
			PresharedKey:                mustParseKey(t, "6qsQOM1KgdWLrFOY5u6QPep8TfUHcInsc3jv6FyIuGE="),
			Endpoint:                    mustParseUDPAddr(t, "5.126.94.232:59749"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.87/32", "fdd0:438e:19ba:5069::57/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679609810),
			ReceiveBytes:                46650644,
			TransmitBytes:               320625632,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "+cOve1HTMMq42V4HBb3MIIH4paGdmBT2AJ/ejV0ZcH4="),
			PresharedKey:                mustParseKey(t, "NWmkGO1bQJ7q894n24JG+PMv9Wq9sHJ63vocV7LAIdo="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.88/32", "fdd0:438e:19ba:5069::58/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "4ATnN6sLBFrmnFbEDDGLdaxtokPX3kSJ4dI+NWDGozA="),
			PresharedKey:                mustParseKey(t, "WVlJrNwx76dcr4SH94ircmCR/LfPGRdNwG2n3Mvq2CE="),
			Endpoint:                    mustParseUDPAddr(t, "5.125.220.175:17516"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.89/32", "fdd0:438e:19ba:5069::59/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679608542),
			ReceiveBytes:                95543172,
			TransmitBytes:               328747492,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "k//CV10qBAipId1soU4e8+J+lsBY47+g9oeuyYQf1UQ="),
			PresharedKey:                mustParseKey(t, "WxEXQe2/HyV+mAXKzOzUFBBLRrbXRvCY7mNKYRvAcos="),
			Endpoint:                    mustParseUDPAddr(t, "83.123.33.104:58065"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.90/32", "fdd0:438e:19ba:5069::5a/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679601734),
			ReceiveBytes:                17100408,
			TransmitBytes:               193644292,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "a5H1GKKO+9dYr1EjvI+bMKSsNfQz9ghO22oBLbCV1W8="),
			PresharedKey:                mustParseKey(t, "St9drSEDjTXMkBRFPVYXDOa6vmJ4f7dtM5A6zCbaF9M="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.91/32", "fdd0:438e:19ba:5069::5b/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "dZd8wgMH+AdaYiAbb/emUsieNYSu3VLh+R5I8FGmCCE="),
			PresharedKey:                mustParseKey(t, "Wf8tdgYDJ6R6BsyAptSwtoPDgdaFykexjPRdX/7TZJQ="),
			Endpoint:                    mustParseUDPAddr(t, "172.80.252.9:40934"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.92/32", "fdd0:438e:19ba:5069::5c/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679607825),
			ReceiveBytes:                38538968,
			TransmitBytes:               479650036,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "WSxAe8GiFWTRwA+NZXYHG2fy5CVSCyt6H00iE4gb/QM="),
			PresharedKey:                mustParseKey(t, "R1wPt+WT/66PR3zGv59EjbCvwDMLQvTXEAuhFEh86Y4="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.93/32", "fdd0:438e:19ba:5069::5d/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "vf/YD28YgyaitJiGyOFw9P9dEj3m4XI7NBEhmsQtuS4="),
			PresharedKey:                mustParseKey(t, "OMKwZuXGTrTPmIIBkA9leNTXiceeTBXYtmM3E2dDHFM="),
			Endpoint:                    mustParseUDPAddr(t, "83.122.177.202:63738"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.94/32", "fdd0:438e:19ba:5069::5e/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679622119),
			ReceiveBytes:                214899792,
			TransmitBytes:               403954636,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "Fgyh9YI3H0BMXpIwuzGsPNuqZBpNX6MHjL9bn2Ao2mc="),
			PresharedKey:                mustParseKey(t, "HPk6anM0vNx7NBle4LqoJip2yMd8RmQzJMDyG7gBOtU="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.95/32", "fdd0:438e:19ba:5069::5f/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "S/VgSky31TeNufxjySvoX71kKhJ6c20gTx3qedTa7Hc="),
			PresharedKey:                mustParseKey(t, "Ad/h2P5vN2A5D3irqe+a6OoZp9+feGkxdsbZv06z0J4="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.96/32", "fdd0:438e:19ba:5069::60/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "UM0khRt0dMUguG+S/mPoqY0QlGnPpFWtpPo6z7BGQWc="),
			PresharedKey:                mustParseKey(t, "GjqKzSvi+1M86YTuZuDwjZ/kN4UHk4m1JqNkbRUpmGk="),
			Endpoint:                    mustParseUDPAddr(t, "95.162.238.226:15423"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.99/32", "fdd0:438e:19ba:5069::63/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679622287),
			ReceiveBytes:                200415180,
			TransmitBytes:               1479932240,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "oxN6Cn4EPglTKBCwgyTiHZ1bU540zD/h8li1laXCbj0="),
			PresharedKey:                mustParseKey(t, "eS7dvTHZF5T60vF1AQ8S/u8jauj6uVz+HMEFOfK85zs="),
			Endpoint:                    mustParseUDPAddr(t, "5.126.57.129:51519"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.100/32", "fdd0:438e:19ba:5069::64/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679613535),
			ReceiveBytes:                61072720,
			TransmitBytes:               1065353052,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "4xcn/eRHubjW9Nytc0AwoN6WSRaTKtjPesz/7ImuCzA="),
			PresharedKey:                mustParseKey(t, "yhWvanfIx4spsiCXH5HT27rvlmL6fjjm2VgxQ4V7/m4="),
			Endpoint:                    mustParseUDPAddr(t, "37.129.11.127:32809"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.101/32", "fdd0:438e:19ba:5069::65/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679619234),
			ReceiveBytes:                27569848,
			TransmitBytes:               539436024,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "nkAGr/kNARuXwSHSKCm7XdpMhrimPLzeWHcsKJcko1g="),
			PresharedKey:                mustParseKey(t, "NPT8zU1f+zm9OTeIODh7ba9JGOoOmpE2un6OPBYz5Xw="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.102/32", "fdd0:438e:19ba:5069::66/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "zMsFGX2HA5AVK2NzOjPXZz0tL6cb70Uwt41YxLEogTA="),
			PresharedKey:                mustParseKey(t, "Xx9kbd2Vtzm6INWM4irctZjdCA3vXwPJfEkj4DpCX9Y="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.103/32", "fdd0:438e:19ba:5069::67/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "WhIPaOV1QuE2Wqgf3kKIzAS/W9rX81cBwFbMFATv7Sg="),
			PresharedKey:                mustParseKey(t, "R01ZFMDTCvB2n+2YCySVowy8v0kXAxSScUi1/ELuVok="),
			Endpoint:                    mustParseUDPAddr(t, "5.119.170.37:61162"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.104/32", "fdd0:438e:19ba:5069::68/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679617644),
			ReceiveBytes:                21586800,
			TransmitBytes:               286196660,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "tCXkXzYf1FXFXNOUP7J0gGEbLYa2+OtmN3xGMjHvn2Y="),
			PresharedKey:                mustParseKey(t, "02znarJ3rTH03+7ZGMORYkfWtMNpDUXv/39be7QQxBY="),
			Endpoint:                    mustParseUDPAddr(t, "5.114.158.124:5085"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.105/32", "fdd0:438e:19ba:5069::69/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679610971),
			ReceiveBytes:                92314284,
			TransmitBytes:               1539250464,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "pyoJ/0UqULHxtXDVbVJTQsOxj0bkrhcxv5QbiAo/2yA="),
			PresharedKey:                mustParseKey(t, "oE5g05yXJhiY19WMCyorKURphqG7/0vfIYTu21bqdGU="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.106/32", "fdd0:438e:19ba:5069::6a/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "2VWoryzJr403WxW/Q0g7d+akXniBOdH7xMzSBkBz4gk="),
			PresharedKey:                mustParseKey(t, "iO61JNmmi4BdL2MWwpIsVBzgO/YuEOXn/5VhIOGQTfE="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.107/32", "fdd0:438e:19ba:5069::6b/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "NehSfEfBgEIMmMWKJc6dNe+Kk+eowsokeUTJlioPOWM="),
			PresharedKey:                mustParseKey(t, "h/76X+Cuwo123DQETPWEDqVTZwiZXHM0zCS9b000WaM="),
			Endpoint:                    mustParseUDPAddr(t, "83.122.28.199:49330"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.108/32", "fdd0:438e:19ba:5069::6c/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679599088),
			ReceiveBytes:                28062836,
			TransmitBytes:               626552728,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "1wgfj6gR/fhrfW909Ejnu9xPS+gCMDhvqCQFzH8tP0Y="),
			PresharedKey:                mustParseKey(t, "Eqekgl1wqY4Gkf8VaApuXVw0354BWbCelcpZv9WzYHc="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.109/32", "fdd0:438e:19ba:5069::6d/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "WtNgBST4LflxP2IR/oaNkULvQ5909BhZaizrtODXFQs="),
			PresharedKey:                mustParseKey(t, "zoR7N5mGIR1/A9PNC08DK2eOFHCBdG+SapPNyXyNGIo="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.110/32", "fdd0:438e:19ba:5069::6e/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "kPuHjkT58MsNP55oQFJszXhpuN6ukG++WUDsR4DJQjo="),
			PresharedKey:                mustParseKey(t, "G4Z57Ihjk3Et+O68iFj8I3CT2SgbjMDnpa+/lFuH6gU="),
			Endpoint:                    nil,
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.111/32", "fdd0:438e:19ba:5069::6f/128"),
			LastHandshakeTime:           time.Time{},
			ReceiveBytes:                0,
			TransmitBytes:               0,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "4ejeXAAOX6UrBnaAivqQ11JSanLzXxBpVd+A1Aozl3A="),
			PresharedKey:                mustParseKey(t, "77EuIYC1+S1i4LP5j8TibwL8didMy+mI+5bOfqKC+sA="),
			Endpoint:                    mustParseUDPAddr(t, "5.232.204.179:10610"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.112/32", "fdd0:438e:19ba:5069::70/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1679607401),
			ReceiveBytes:                20695096,
			TransmitBytes:               349602820,
			PersistentKeepaliveInterval: 0,
		},
		{
			PublicKey:                   mustParseKey(t, "RFrjTKnZWr/CmI/XC+zgPeGJzbDGQhA9eOBb+ge90wA="),
			PresharedKey:                mustParseKey(t, "uox3nW4W81djWoJEN/lgdrvjqtLnl36lhkJZcbSkwOM="),
			Endpoint:                    mustParseUDPAddr(t, "5.232.204.179:10611"),
			AllowedIPs:                  mustParseAllowedIPs(t, "10.0.0.113/32", "fdd0:438e:19ba:5069::71/128"),
			LastHandshakeTime:           mustParseUnixTime(t, 1672307401),
			ReceiveBytes:                10695096,
			TransmitBytes:               349602820,
			PersistentKeepaliveInterval: 10,
		},
	}, peers)
}

func TestParseInvalidDumpFile(t *testing.T) {
	scanner := bufio.NewScanner(bytes.NewBuffer(rawInvalidDump))
	for scanner.Scan() {
		peers, err := dump.Parse(bytes.NewBuffer(scanner.Bytes()))
		require.Nil(t, peers)
		require.NotNil(t, err)
	}
}