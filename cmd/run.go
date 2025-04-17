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
	"fmt"
	"time"
	"strings"

	"io/ioutil"

	"github.com/apex/log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

)

type RPKIstat struct {
	Domain string
	Date time.Time
	Names int
	NamesFull int
	NamesPartial int
	IPv4 int
	IPv4roas int
	IPv6 int
	IPv6roas int
	TAs4 int
	TAs6 int
	AS4 int
	AS6 int
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Get RPKI statistics",
	Long: `Get RPKI statistics for a single domain name or a list of domain names and optionally save to MariaDB`,
	Run: execRun,
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")
	runCmd.Flags().StringP(DOMAIN, DOMAIN_SHORT, "", "domain name")
	runCmd.Flags().StringP(DOMAIN_FILE, DOMAIN_FILE_SHORT, "", "file with a list of domain names")
	runCmd.Flags().StringP(ROUTINATOR, ROUTINATOR_SHORT, "", "address (including port) of the routinator instance to use")
	runCmd.Flags().StringP(RESOLVER, RESOLVER_SHORT, "", "address of the resolver to use")

	// Use flags for viper values
	viper.BindPFlags(runCmd.Flags())
}

func execRun(cmd *cobra.Command, args []string) {
	fmt.Println("run called")

	if viper.GetString(ROUTINATOR) == "" {
		cmd.Help();
		log.Fatal("Routinator must be given.")
	}

	if viper.GetString(RESOLVER) == "" {
		cmd.Help();
		log.Fatal("Resolver must be given.")
	}

	if viper.GetString(DOMAIN) == "" && viper.GetString(DOMAIN_FILE) == "" {
		cmd.Help();
		log.Fatal("Domain or domain list must be given.")
	}

	// single domain gets only printed on command line, not saved to database
	if viper.GetString(DOMAIN) != "" {
		viper.Set(VERBOSE, VERBOSE_INFO)
		domain := viper.GetString(DOMAIN)
		log.Debugf("Single domain statistics (no db): %s", domain)
		domainStat(domain)
		return
	}

	// open database
	//db := openDB()

	
	domainfile := viper.GetString(DOMAIN_FILE)
	log.Debugf("Using domain file: %s", domainfile)
	handleDomainList(domainfile)
}

func handleDomainList(filename string) (stats []*RPKIstat) {
	stats = make([]*RPKIstat, 0)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading Domain file %s: %s", filename, err)
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		 line= strings.TrimSpace(line)
		
		 // jump over empty lines
		if line == "" {
			continue
		}
		
		// jump over comments
		if strings.HasPrefix(line, "#") {
			continue
		}
		
		// now we should have a valid domain 
		domain := strings.ToLower(line)

		log.Debugf("Running domain: %s", domain)
		stats = append(stats, domainStat(domain))
	}
	return
}

func domainStat(domain string) (stat *RPKIstat) {
	stat = &RPKIstat{Domain: domain, Date: time.Now()}

	nameservers := getNS(domain)

	name2ip4 := make(map[string][]string, 0)
	name2ip6 := make(map[string][]string, 0)

	ip4list := make([]string, 0)
	ip6list := make([]string, 0)

	for _,ns := range nameservers {
		name2ip4[ns] = getIP4(ns)
		name2ip6[ns] = getIP6(ns)
	}

	ip4roas := make(map[string]*ROA, 0)
	ip6roas := make(map[string]*ROA, 0)
	for ns := range name2ip4 {
		for _,ip4 := range name2ip4[ns] {
			ip4list = append(ip4list, ip4)
			if _,ok:=ip4roas[ip4]; ok {
				// already done
				continue
			}
			roa := getROA(ip4)
			if roa != nil {
				ip4roas[ip4] = roa
			}
		}
	}
	ip4list = unique(ip4list)

	for ns := range name2ip6 {
		for _,ip6 := range name2ip6[ns] {
			ip6list = append(ip6list, ip6)
			if _,ok:=ip6roas[ip6]; ok {
				// already done
				continue
			}
			roa := getROA(ip6)
			if roa != nil {
				ip6roas[ip6] = roa
			}
		}
	}
	ip6list = unique(ip6list)

	ta4 := make([]string, 0)
	ta6 := make([]string, 0)
	asn4 := make([]string, 0)
	asn6 := make([]string, 0)

	for ip4 := range ip4roas {
		if ip4roas[ip4] == nil {
			continue
		}
		for _,ta := range ip4roas[ip4].Ta {
			ta4 = append(ta4, ta)
		}
		for _, asn := range ip4roas[ip4].Asn {
			asn4 = append(asn4, asn)
		}
	}
	ta4 = unique(ta4)
	asn4 = unique(asn4)

	for ip6 := range ip6roas {
		if ip6roas[ip6] == nil {
			continue
		}
		for _,ta := range ip6roas[ip6].Ta {
			ta6 = append(ta6, ta)
		}	
		for _, asn := range ip6roas[ip6].Asn {
			asn6 = append(asn6, asn)
		}
	}
	ta6 = unique(ta6)
	asn6 = unique(asn6)

	names_full := 0
	names_partial := 0

	for _,ns := range nameservers {
    	roas := 0
    	for _,ip4 := range name2ip4[ns] {
			if _,ok := ip4roas[ip4]; ok {
        		roas++ 
			}
    	}
    	for _,ip6 := range name2ip6[ns] {
			if _,ok := ip6roas[ip6]; ok {
        		roas++ 
			}
    	}
	    if roas == len(name2ip4[ns])+len(name2ip6[ns]) {
        	log.Debugf("%s is full", ns);
        	names_full++
    	} else if roas > 0 {
        	log.Debugf("%s is partial", ns)
        	names_partial++
		}
    }



	stat.Names = len(nameservers)
	stat.NamesFull = names_full
	stat.NamesPartial = names_partial
	stat.IPv4 = len(ip4list)
	stat.IPv4roas = len(ip4roas)
	stat.IPv6 = len(ip6list)
	stat.IPv6roas = len(ip6roas)
	stat.TAs4 = len(ta4)
	stat.TAs6 = len(ta6)
	stat.AS4 = len(asn4)
	stat.AS6 = len(asn6)

	log.Debugf("Result for %s: %v", domain, stat)

	return
}
