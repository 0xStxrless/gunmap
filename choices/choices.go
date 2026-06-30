package choices

import "math/rand"

type SubCategory struct {
	Title         string
	IsMultiChoice bool
	// this bool just tells if command takes input or no
	Commands map[string]bool
}

func GetAllPages() []*SubCategory {
	pages := []*SubCategory{
		HostDiscovery(),
		ScanTechniques(),
		PortSelection(),
		ServiceDetection(),
		OSDetection(),
		TimingAndPerformance(),
		FirewallEvasionNSpoofing(),
		Output(),
		Misc(),
	}
	return pages
}

func HostDiscovery() *SubCategory {
	return &SubCategory{
		Title:         "HostDiscovery",
		IsMultiChoice: true,
		Commands: map[string]bool{
			"-sL List Scan just list targets":                             false,
			"-sn Ping Scan - disable port scan":                           false,
			"-Pn Treat all hosts as online, skip discovery":               false,
			"-PS optional string TCP SYN discovery":                       false,
			"-PA optional string TCP ACK discovery":                       false,
			"-PU optional string UDP discovery":                           false,
			"-PY optional string SCTP discovery":                          false,
			"-PE ICMP echo discovery":                                     false,
			"-PP ICMP timestamp discovery":                                false,
			"-PM ICMP netmask discovery":                                  false,
			"-PO optional string IP protocol ping":                        false,
			"-PR ARP ping (local network)":                                false,
			"-n Never do DNS resolution":                                  false,
			"-R Always resolve DNS":                                       false,
			"--dns-servers <srv1: false,srv2> specify custom DNS servers": true,
			"--system-dns use OS's DNS resolver":                          false,
			"--traceroute trace hop path to host":                         false,
		},
	}
}

func ScanTechniques() *SubCategory {
	return &SubCategory{
		Title:         "ScanTechniques",
		IsMultiChoice: false,
		Commands: map[string]bool{
			"-sS TCP SYN scan (default, needs root)": false,
			"-sT TCP Connect scan":                   false,
			"-sA TCP ACK scan":                       false,
			"-sW TCP Window scan":                    false,
			"-sM TCP Maimon scan":                    false,
			"-sU UDP scan":                           false,
			"-sN TCP Null scan":                      false,
			"-sF TCP FIN scan":                       false,
			"-sX TCP Xmas scan":                      false,
			"-sY SCTP INIT scan":                     false,
			"-sZ SCTP COOKIE-ECHO scan":              false,
			"-sO IP protocol scan":                   false,
		},
	}
}

func PortSelection() *SubCategory {
	return &SubCategory{
		Title:         "PortSelection",
		IsMultiChoice: false,
		Commands: map[string]bool{
			"-p <port ranges> e.g. 22, 1-65535, U:53,T:21-25":                         true,
			"-p- all 65535 ports (shorthand for -p1-65535)":                           false,
			"--exclude-ports <port ranges> ports to exclude":                          true,
			"-F fast scan fewer ports than default":                                   false,
			"-r scan ports sequentially, don't randomize":                             false,
			"--top-ports <num> scan N most common ports":                              true,
			"--port-ratio <ratio> float (0–1) scan ports above given frequency ratio": true,
		},
	}
}

func ServiceDetection() *SubCategory {
	return &SubCategory{
		Title:         "ServiceDetection",
		IsMultiChoice: false,
		Commands: map[string]bool{
			"-sV probe open ports for service/version info":              false,
			"--version-intensity int slider intensity of version probes": true,
			"--version-light shorthand for intensity 2":                  false,
			"--version-all shorthand for intensity 9 (try every probe)":  false,
			"--version-trace show detailed version scan activity":        false,
		},
	}
}

func OSDetection() *SubCategory {
	return &SubCategory{
		Title:         "OS Detection",
		IsMultiChoice: false,
		Commands: map[string]bool{
			"-O enable OS detection":                                 false,
			"--osscan-limit limit OS detection to promising targets": false,
			"--osscan-guess guess OS more aggressively":              false,
			"--max-os-tries number set max OS detection tries":       true,
		},
	}
}

func TimingAndPerformance() *SubCategory {
	return &SubCategory{
		Title:         "Timing And Performance",
		IsMultiChoice: true,
		Commands: map[string]bool{
			"-T<0-5> (Paranoid/Sneaky/Polite/Normal/Aggressive/Insane) timing template": true,
			"-min-hostgroup <size> num of parallel host scan group sizes":               true,
			"--max-hostgroup <size>":                                    true,
			"--min-parallelism <num> probe parallelization":             true,
			"--max-parallelism <num> int":                               true,
			"--min-rtt-timeout <time> duration (e.g. 100ms)":            true,
			"--max-rtt-timeout <time> duration":                         true,
			"--initial-rtt-timeout <time> duration":                     true,
			"--max-retries <num> max port scan probe retransmissions":   true,
			"--host-timeout <time> duration give up on slow hosts":      true,
			"--scan-delay <time> duration delay between probes":         true,
			"--max-scan-delay <time> duration max delay between probes": true,
			"--min-rate <num> num of min packets sent per second":       true,
			"--max-rate <num> int max packets sent per second":          true,
			"--defeat-rst-ratelimit ignore RST rate limiting":           false,
		},
	}
}

func FirewallEvasionNSpoofing() *SubCategory {
	return &SubCategory{
		Title:         "Firewall And Evasion Spoofing",
		IsMultiChoice: true,
		Commands: map[string]bool{
			"-f fragment packets ":                                                     false,
			"--mtu <val> set custom fragment MTU ":                                     true,
			"-D <decoy1,decoy2[,ME],...> string list decoy scan ":                      true,
			"-S <IP> string spoof source address ":                                     true,
			"-e <iface> string use specific interface ":                                true,
			"--source-port / -g <port>  int spoof source port ":                        true,
			"--proxies <url1,url2,...> string list relay through HTTP/SOCKS4 proxies ": true,
			"--data <hex string>string append custom binary data to packets ":          true,
			"--data-string <string> string append custom ASCII string ":                true,
			"--data-length <num> int append random data of given length ":              true,
			"--ip-options <options> string send packets with specified IP options ":    true,
			"--ttl <val> int (0–255) set IP time-to-live ":                             true,
			"--spoof-mac <mac/prefix/vendor> string spoof MAC address ":                true,
			"--badsum send packets with bogus TCP/UDP checksum ":                       false,
			"--adler32 use deprecated Adler32 checksum for SCTP ":                      false,
		},
	}
}

func Output() *SubCategory {
	return &SubCategory{
		Title:         "Output",
		IsMultiChoice: true,
		Commands: map[string]bool{
			"-oN <file> filepath normal output ":                           true,
			"-oX <file> filepath XML output":                               true,
			"-oS <file> filepath script kiddie output":                     true,
			"-oG <file> filepath grepable outpu":                           true,
			"-oA <basename> filepath all formats at onc":                   true,
			"-v (repeatable) counter increase verbosit":                    true,
			"-d (repeatable) counter increase debuggin":                    true,
			"--reason show reason for port stat":                           false,
			"--open only show open port":                                   false,
			"--packet-trace show all packets sent/receive":                 false,
			"--iflist print host interfaces/route":                         false,
			"--append-output append rather than overwrite output file":     false,
			"--noninteractive disable runtime keyboard interaction":        false,
			"--stylesheet <path/URL> string XSL stylesheet for XML output": true,
			"--webxml reference Nmap.org stylesheet":                       false,
			"--no-stylesheet omit XSL association":                         false,
		},
	}
}

func Misc() *SubCategory {
	return &SubCategory{
		Title:         "Misc",
		IsMultiChoice: true,
		Commands: map[string]bool{
			"-6 enable IPv6 scanning":                                   false,
			"-A aggressive: OS detect + version + scripts + traceroute": false,
			"--datadir <dir> filepath custom Nmap data file location":   true,
			"--send-eth send using raw ethernet frames":                 false,
			"--send-ip send using raw IP packets":                       false,
			"--privileged assume user is privileged":                    false,
			"--unprivileged assume user lacks raw socket privileges":    false,
		},
	}
}

func RandomErrFace() string {
	faces := []string{
		"( ꩜ ᯅ ꩜)⁭ ⁭",
		"(╥﹏╥)",
		"(｡•́︿•̀｡)",
		"՞߹ - ߹՞",
		"( • ᴖ • ｡)",
		"૮◞ ‸ ◟ ა",
		"｡°(°¯᷄◠¯᷅°)°｡",
		":‹",
		":(",
	}
	face := faces[rand.Intn(len(faces))]
	return face
}
