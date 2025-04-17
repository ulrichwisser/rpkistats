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
//	"fmt"
	"net/http"
	"net/netip"
	"encoding/json"
	
	"github.com/apex/log"

	"github.com/spf13/viper"
)

type ROA struct {
	Ip string
	Prefix string
	Asn []string
	Ta []string
}

func getROA(ip string) (roa *ROA) {

	prefix := ip2prefix(ip)
	url := viper.GetString(ROUTINATOR) + prefix
	log.Debugf("Routinator URL: %s", url)

	roa = &ROA{Ip: ip, Prefix: prefix, Asn: make([]string, 0), Ta: make([]string, 0)}

	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Error contacting routinator: %s", err)
		return nil
	}
	
	var response map[string]interface{}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Errorf("Error decoding received data: %s", err)
		return nil
	}

	roaraw := response["roas"].([]interface{})

	if len(roaraw) == 0 {
		log.Debugf("No ROA found for %s", prefix)
		return nil
	}

	ta := make(map[string]bool)
	asn := make(map[string]bool)

	for _,r := range roaraw {
		vrp := r.(map[string]interface{})
		ta[vrp["ta"].(string)] = true
		asn[vrp["asn"].(string)] = true
	}

	for t := range ta {
		roa.Ta = append(roa.Ta, t) 
	}
	for a := range asn {
		roa.Asn = append(roa.Asn, a)
	}

	return
}

func ip2prefix(ipstr string) string {

	ip,err := netip.ParseAddr(ipstr)
	if err != nil {
		log.Fatalf("Could not parse ip %s: %s", ipstr, err)
	}

	mask := 24
	if ip.Is6() {
		mask = 64
	} 
	var prefix netip.Prefix
	prefix,err = ip.Prefix(mask)
	if err != nil {
		log.Fatalf("Could mask ip %s mask %d: %s", ip, mask, err)
	}
	return prefix.String()
}
