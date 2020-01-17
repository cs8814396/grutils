package grnetwork

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
	"net/url"
)

// DNSServer 提供备用的DNS
var (
	DNSServer []string
)

// Init 初始化备用DNS
func Init(dnsServer []string) {
	DNSServer = dnsServer
}

// DNSDecoder 将URL转化为IP
func DNSDecoder(urlStr string) (urlString, host string, err error) {
	var u *url.URL
	u, err = url.Parse(urlStr)
	if err != nil {
		return
	}
	host = u.Host
	netIP := getIPFromDNSSlice(u.Host)
	if netIP != nil {
		u.Host = netIP.String()
		urlString = u.String()
		return
	}
	err = fmt.Errorf("dnsDecoder fail")
	return
}

func getIPFromDNSSlice(host string) (ip net.IP) {
	var err error
	ip = nil
	for _, dnsServer := range DNSServer {
		ip, err = getIP(host, dnsServer+":53")
		if err == nil && ip != nil {
			return
		}
	}
	return
}

func getIP(host, dnsServer string) (ip net.IP, err error) {
	ip = nil
	err = nil

	var addrs *dns.Msg
	addrs, err = lookupFromType("CNAME", host, dnsServer)
	if err != nil {
		//lgd.Warn("dns chame fail with the host[%s]. error: [%s]", host, err.Error())
		return
	}

	for {
		if len(addrs.Answer) == 0 {
			break
		}
		host = addrs.Answer[0].(*dns.CNAME).Target
		addrs, err = lookupFromType("CNAME", host, dnsServer)
		if err != nil {
			//lgd.Warn("dns chame fail with the host[%s]. error: [%s]", host, err.Error())
			return
		}
	}
	addrs, err = lookupFromType("A", host, dnsServer)
	if err != nil {
		//lgd.Warn("dns a fail with the host[%s]. error: [%s]", host, err.Error())
		return
	}
	for _, a := range addrs.Answer {
		if a.(*dns.A).A != nil {
			ip = a.(*dns.A).A
			return
		}
	}

	return
}

func lookupFromType(ctype, host, dnsServer string) (response *dns.Msg, err error) {

	itype, ok := dns.StringToType[ctype]
	if !ok {
		return nil, fmt.Errorf("invalid type %s", ctype)
	}

	host = dns.Fqdn(host)
	client := &dns.Client{}
	msg := &dns.Msg{}
	msg.SetQuestion(host, itype)
	response, err = lookup(msg, client, dnsServer, false)
	return
}

func lookup(msg *dns.Msg, client *dns.Client, server string, edns bool) (response *dns.Msg, err error) {
	if edns {
		opt := &dns.OPT{
			Hdr: dns.RR_Header{
				Name:   ".",
				Rrtype: dns.TypeOPT,
			},
		}
		opt.SetUDPSize(dns.DefaultMsgSize)
		msg.Extra = append(msg.Extra, opt)
	}
	response, _, err = client.Exchange(msg, server)
	if err != nil {
		return
	}

	if msg.Id != response.Id {
		err = fmt.Errorf("DNS ID mismatch, request: %d, response: %d", msg.Id, response.Id)
		return
	}

	if response.MsgHdr.Truncated {
		if client.Net == "tcp" {
			err = fmt.Errorf("got truncated message on tcp")
			return
		}
		if edns {
			client.Net = "tcp"
		}
		return lookup(msg, client, server, !edns)
	}
	return
}
