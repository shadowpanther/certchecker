package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/miekg/dns"
)

// Commandline flags
var (
	verbose, help bool
	dnsServer     string
)

// Commandline flags declaration
func init() {
	flag.BoolVar(&verbose, "v", false, "output raw request results")
	flag.BoolVar(&help, "h", false, "show this help")
	flag.StringVar(&dnsServer, "dns", "1.1.1.1", "set the DNS server to query from")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "%s [-h] [-v] [-dns 1.2.3.4] host.name...\n", os.Args[0])
		flag.PrintDefaults()
	}
}

// A simple wrapper for making DNS queries
func dnsQuery(hostname string, recordType uint16) (result []string, err error) {
	m := new(dns.Msg)
	m.SetQuestion(hostname, recordType)
	in, err := dns.Exchange(m, dnsServer+":53")
	if err != nil {
		return
	}
	if len(in.Answer) == 0 {
		err = errors.New("record not found")
		return
	}
	switch recordType {
	case dns.TypeCNAME:
		if t, ok := in.Answer[0].(*dns.CNAME); ok {
			result = append(result, t.Target)
		}
	case dns.TypeTXT:
		if t, ok := in.Answer[0].(*dns.TXT); ok {
			result = t.Txt
		}
	default:
		//
	}
	return
}

func verbosePrintln(a ...any) {
	if verbose {
		fmt.Println(a...)
	}
}

func main() {
	flag.Parse()

	if help {
		flag.Usage()
		return
	}
	if len(flag.Args()) == 0 {
		fmt.Println("No name(s) specified!")
		flag.Usage()
		return
	}

	first := true
	for _, challengeName := range flag.Args() {
		var (
			txtCNAME, txtStraight, hasCNAME []string
			err                             error
		)

		if first {
			first = false
		} else {
			fmt.Println()
		}

		if !strings.HasSuffix(challengeName, ".") {
			challengeName = challengeName + "."
		}
		if !strings.HasPrefix(challengeName, "_acme-challenge.") {
			challengeName = "_acme-challenge." + challengeName
		}

		fmt.Println("Checking for DNS name:", challengeName)
		if hasCNAME, err = dnsQuery(challengeName, dns.TypeCNAME); err == nil {
			verbosePrintln("CNAME resolve:", hasCNAME[0])
			if txtCNAME, err = dnsQuery(hasCNAME[0], dns.TypeTXT); err == nil {
				verbosePrintln("TXTs at CNAME:", txtCNAME)
			}
		}
		if txtStraight, err = dnsQuery(challengeName, dns.TypeTXT); err == nil {
			verbosePrintln("Straight TXTs:", txtStraight)
		}
		fmt.Print("This name")
		if len(hasCNAME) > 0 {
			fmt.Print(color.GreenString(" has a CNAME record"), " that points to ", hasCNAME)
			if len(txtCNAME) > 0 {
				fmt.Println(color.GreenString(" and has these TXTs: "), txtCNAME)
			} else {
				fmt.Println(color.RedString(" but it doesn't have any TXTs"))
			}
			if len(txtStraight) > 0 {
				fmt.Println(color.YellowString("It also has straight TXT records: "), txtStraight)
			}
		} else {
			fmt.Print(color.RedString(" doesn't have a CNAME record"))
			if len(txtStraight) > 0 {
				fmt.Println(color.GreenString(" but it has straight TXT records: "), txtStraight)
			} else {
				fmt.Println(color.RedString(" and it doesn't have any straight TXT records"))
			}
		}
	}
}
