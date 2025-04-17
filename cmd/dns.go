/*
Copyright © 2025 Ulrich Wisser

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"net"
	"strings"


	"github.com/apex/log"

	"github.com/spf13/viper"

	"github.com/miekg/dns"

)

func getNS(domain string) (nslist []string) {
	nslist = make([]string, 0)
	msg := resolve(domain, dns.TypeNS)
	if msg == nil {
		log.Errorf("No name servers for %s", domain)
		return
	}
	if msg.Rcode != dns.RcodeSuccess {
		log.Errorf("NS resolution error for %s: %s (Rcode %d)", domain, dns.RcodeToString[msg.Rcode], msg.Rcode)
		return
	}
	for _,rr := range msg.Answer {
		if rr.Header().Rrtype != dns.TypeNS {
			continue
		}
		nslist = append(nslist, rr.(*dns.NS).Ns)
	}

	nslist = unique(nslist)

	return
}

func getIP4(domain string) (ip4list []string) {
	ip4list = make([]string, 0)
	msg := resolve(domain, dns.TypeA)
	if msg == nil {
		log.Errorf("No IPv4 for %s", domain)
		return
	}
	if msg.Rcode != dns.RcodeSuccess {
		log.Errorf("IPv4 resolution error %s: %s (Rcode %d)", domain, dns.RcodeToString[msg.Rcode], msg.Rcode)
		return
	}
	for _,rr := range msg.Answer {
		if rr.Header().Rrtype != dns.TypeA {
			continue
		}
		ip4list = append(ip4list, rr.(*dns.A).A.String())
	}

	ip4list = unique(ip4list)

	return
}

func getIP6(domain string) (ip6list []string) {
	ip6list = make([]string, 0)
	msg := resolve(domain, dns.TypeAAAA)
	if msg == nil {
		log.Errorf("No IPv6 for %s", domain)
		return
	}
	if msg.Rcode != dns.RcodeSuccess {
		log.Errorf("IPv6 resolution error %s: %s (Rcode %d)", domain, dns.RcodeToString[msg.Rcode], msg.Rcode)
		return
	}
	for _,rr := range msg.Answer {
		if rr.Header().Rrtype != dns.TypeAAAA {
			continue
		}
		ip6list = append(ip6list, rr.(*dns.AAAA).AAAA.String())
	}

	ip6list = unique(ip6list)

	return
}




// getResolvers will read the list of resolvers from /etc/resolv.conf
func getResolver() string {

	resolver := viper.GetString(RESOLVER)

	ip := net.ParseIP(resolver)
	if ip == nil {
		log.Fatalf("Could not parse resolver ip: %s", resolver)
	}

	ipstr := ip.String()
	if strings.ContainsAny(":", ipstr) {
		// IPv6 address
		ipstr = "[" + ipstr + "]:53"
	} else {
		// IPv4 address
		ipstr = ipstr + ":53"
	}
	return ipstr
}

// resolv will send a query and save the result
func resolve(domain string, qtype uint16) *dns.Msg {
	log.Debugf("CALLED resolve(%s, %d)", domain, qtype)
	
	server := getResolver()

	// Setting up query
	query := new(dns.Msg)
	query.RecursionDesired = true
	query.Question = make([]dns.Question, 1)
	//query.SetEdns0(1232, false)
	//query.IsEdns0().SetDo()

	// Setting up resolver
	client := new(dns.Client)
	client.ReadTimeout = TIMEOUT * 1e9
	client.Net = "tcp"

	query.SetQuestion(dns.Fqdn(domain), qtype)

	// limit repeats
	var repeat int = 0
	// query until we get an answer
	for {

		// limit repeats
		repeat++
		log.Debugf("%-30s: %d repeats reached (server %s, %s)", domain, repeat, server, dns.TypeToString[qtype])
		if repeat > 10 {
			log.Errorf("%-30s: 10 repeats reached (server %s)", domain, server)
			break
		}

		// make the query and wait for answer
		r, _, err := client.Exchange(query, server)

		// check for errors
		if err != nil {
			log.Errorf("%-30s: Error resolving %s (server %s)", domain, err, server)
			continue
		}
		if r == nil {
			log.Errorf("%-30s: No answer (Server %s)", domain, server)
			continue
		}
		if r.Rcode != dns.RcodeSuccess {
			log.Errorf("%-30s: %s (Rcode %d, Server %s)", domain, dns.RcodeToString[r.Rcode], r.Rcode, server)
			break
		}

		// we got an answer
		return r
	}

	return nil
}

