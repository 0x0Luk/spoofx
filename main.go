package main

import (
	"bufio"
	"encoding/csv"
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

func logToFile(domain string, dmarc DMARCResult, spf SPFResult) {
	f, err := os.OpenFile("log.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		color.Red("[!] Failed to write log: %v", err)
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	entry := []string{
		time.Now().Format("2006-01-02 15:04:05"),
		domain,
		dmarc.Policy,
		spf.Record,
		spf.Rigor,
	}
	writer.Write(entry)
}

func scan(domain string) {
	mag := color.New(color.FgMagenta).SprintFunc()
	gray := color.New(color.FgHiWhite).SprintFunc()
	alert := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s %s %s\n", mag("[*] Domain:"), gray(domain), "")
	dmarc := checkDMARC(domain)
	spf := checkSPF(domain)

	// DMARC
	if dmarc.Policy != "" {
		fmt.Printf("    %s %s\n", mag("DMARC Policy:"), gray(dmarc.Policy))

		if dmarc.Policy == "none" {
			fmt.Println(alert("    [!] DMARC policy is weak: none"))
		}

		if dmarc.Rua != "" || dmarc.Ruf != "" {
			fmt.Printf("    %s %s\n", mag("DMARC Reporting:"), gray(fmt.Sprintf("rua=%s, ruf=%s", dmarc.Rua, dmarc.Ruf)))
		}
		fmt.Printf("    %s %s\n", mag("Full DMARC Record:"), gray(dmarc.FullRecord))
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

	logToFile(domain, dmarc, spf)
	fmt.Println(strings.Repeat("-", 50))
}

func banner() {
	fmt.Printf("\033[31m" + `


	.▄▄ ·  ▄▄▄·            ·▄▄▄▐▄• ▄ 
	▐█ ▀. ▐█ ▄█▪     ▪     ▐▄▄· █▌█▌▪
	▄▀▀▀█▄ ██▀· ▄█▀▄  ▄█▀▄ ██▪  ·██· 
	▐█▄▪▐█▐█▪·•▐█▌.▐▌▐█▌.▐▌██▌.▪▐█·█▌
	 ▀▀▀▀ .▀    ▀█▄▀▪ ▀█▄▀▪▀▀▀ •▀▀ ▀▀	
	    ` + "\033[35mhttps://github.com/0x0Luk\n\n")
}

func main() {
	flag.Parse()
	args := flag.Args()
	var domains []string
	banner()

	if len(args) > 0 {
		filePath := args[0]
		f, err := os.Open(filePath)
		if err != nil {
			color.Red("[!] Failed to open file: %v", err)
			return
		}
		defer f.Close()
	
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			domains = append(domains, strings.TrimSpace(scanner.Text()))
		}
	} else {
		info, _ := os.Stdin.Stat()
		if (info.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				domains = append(domains, strings.TrimSpace(scanner.Text()))
			}
		} else {
			color.Red("Usage: spoofx domains.txt or tool | spoofx")
			return
		}
	}
	

	for _, domain := range domains {
		scan(domain)
	}
}
