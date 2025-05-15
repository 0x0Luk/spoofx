package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type DMARCResult struct {
	Policy     string
	Rua        string
	Ruf        string
	FullRecord string
}

type SPFResult struct {
	Record string
	Rigor  string
}

var verbose bool

func checkDMARC(domain string) DMARCResult {
	var result DMARCResult
	records, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		return result
	}
	for _, record := range records {
		if strings.Contains(record, "v=DMARC1") {
			parts := strings.Split(record, ";")
			result.FullRecord = record
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "p=") {
					result.Policy = strings.TrimPrefix(part, "p=")
				} else if strings.HasPrefix(part, "rua=") {
					result.Rua = strings.TrimPrefix(part, "rua=")
				} else if strings.HasPrefix(part, "ruf=") {
					result.Ruf = strings.TrimPrefix(part, "ruf=")
				}
			}
			break
		}
	}
	return result
}

func checkSPF(domain string) SPFResult {
	var result SPFResult
	records, err := net.LookupTXT(domain)
	if err != nil {
		return result
	}
	for _, record := range records {
		if strings.HasPrefix(record, "v=spf1") {
			result.Record = record
			if strings.Contains(record, "-all") {
				result.Rigor = "strict"
			} else if strings.Contains(record, "~all") {
				result.Rigor = "soft"
			} else if strings.Contains(record, "?all") {
				result.Rigor = "neutral"
			} else {
				result.Rigor = "unknown"
			}
			break
		}
	}
	return result
}

func scan(domain string) {
	mag := color.New(color.FgMagenta).SprintFunc()
	gray := color.New(color.FgHiWhite).SprintFunc()
	alert := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s %s\n", mag("[*] Domain:"), gray(domain))

	dmarc := checkDMARC(domain)
	spf := checkSPF(domain)

	if verbose {
		// DMARC
		if dmarc.Policy != "" {
			fmt.Printf("    %s %s\n", mag("DMARC Policy:"), gray(dmarc.Policy))
			if dmarc.Rua != "" || dmarc.Ruf != "" {
				fmt.Printf("    %s %s\n", mag("DMARC Reporting:"), gray(fmt.Sprintf("rua=%s, ruf=%s", dmarc.Rua, dmarc.Ruf)))
			}
			fmt.Printf("    %s %s\n", mag("Full DMARC Record:"), gray(dmarc.FullRecord))
			if dmarc.Policy == "none" {
				fmt.Println(alert("    [!] DMARC policy is weak: none"))
			}
		} else {
			fmt.Println(alert("    [!] No DMARC policy found"))
		}

		// SPF
		if spf.Record != "" {
			fmt.Printf("    %s %s\n", mag("SPF Record:"), gray(spf.Record))
			fmt.Printf("    %s %s\n", mag("SPF Strictness:"), gray(spf.Rigor))
			if spf.Rigor == "neutral" || spf.Rigor == "unknown" || spf.Rigor == "" {
				fmt.Println(alert(fmt.Sprintf("    [!] SPF is weak or undefined: %s", spf.Rigor)))
			}
		} else {
			fmt.Println(alert("    [!] No SPF record found"))
		}

		fmt.Println(strings.Repeat("-", 50))
	}
	time.Sleep(300 * time.Millisecond)
}

func banner() {
	fmt.Printf("\033[31m" + `
	.▄▄ ·  ▄▄▄·            ·▄▄▄▐▄• ▄ 
	▐█ ▀. ▐█ ▄█▪     ▪     ▐▄▄· █▌█▌▪
	▄▀▀▀█▄ ██▀· ▄█▀▄  ▄█▀▄ ██▪  ·██· 
	▐█▄▪▐█▐█▪·•▐█▌.▐▌▐█▌.▐▌██▌.▪▐█·█▌
	 ▀▀▀▀ .▀    ▀█▄▀▪ ▀█▄▀▪▀▀▀ •▀▀ ▀▀	
	    ` + "\033[35mhttps://github.com/luq0x/spoofx\n\n")
}

func printUsage() {
	red := color.New(color.FgRed).SprintFunc()
	mag := color.New(color.FgMagenta).SprintFunc()

	fmt.Println(red("Usage:"))
	fmt.Println(mag("  spoofx [-v] -d domain.com"))
	fmt.Println(mag("  spoofx [-v] domains.txt"))
	fmt.Println(mag("  cat list.txt | spoofx [-v]"))

	fmt.Println(red("\n\nParameters: "))
	fmt.Println(mag("  -h : help"))
	fmt.Println(mag("  -d : single domain to scan"))
	fmt.Println(mag("  -v : enable verbose output"))
}

func main() {
	flag.Usage = func() {
		banner()
		printUsage()
	}

	var domain string
	flag.StringVar(&domain, "d", "", "single domain to scan")
	flag.BoolVar(&verbose, "v", false, "enable verbose output")

	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	args := flag.Args()
	var domains []string
	banner()

	if domain != "" {
		domains = append(domains, domain)

	} else if len(args) > 0 {
		filePath := args[0]
		f, err := os.Open(filePath)
		if err != nil {
			color.Red("[!] Failed to open file: %v", err)
			return
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			d := strings.TrimSpace(scanner.Text())
			if d != "" {
				domains = append(domains, d)
			}
		}
	} else {
		info, _ := os.Stdin.Stat()
		if (info.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				d := strings.TrimSpace(scanner.Text())
				if d != "" {
					domains = append(domains, d)
				}
			}
		} else {
			printUsage()
			return
		}
	}

	for _, domain := range domains {
		scan(domain)
	}
}
