<h1 align="center"> ✉️ SpoofX </h1>

**spoofx** is a lightweight, high-performance CLI tool built in Go for identifying potential email spoofing vectors in domains.  

It was designed with **offensive security workflows** in mind — particularly for:

- 🔍 **Bug bounty recon**
- 🔥 **Red team infrastructure mapping**
- 🧰 **Security audits for SPF/DMARC compliance**

## Install

From go:

``` bash
go install github.com/0x0Luk/spoofx@latest
```

## Usage

Scan domains from file:

``` bash
spoofx domains.txt
```

Or via pipe:
``` bash
subfinder -d www.test.com -all -silent | spoofx
```

## 🔍 What it does

- ✅ Fetches and parses SPF records
- ✅ Extracts and classifies SPF strictness: strict, soft, neutral, unknown
- ✅ Looks up DMARC policy, rua, ruf, and full TXT content
- 🚨 Highlights weak or missing configurations in red
- 🗂 Logs every result to log.csv

## Output log 
A CSV file (log.csv) is generated automatically with:

``` bash
Timestamp,Domain,DMARC Policy,SPF Record,SPF Strictness
2025-05-08 14:21:00,badmail.com,,,
```

# 💥 Why use spoofx?

- 🔥 Fast CLI tool for mass-scanning SPF/DMARC
- 🛡️ Helps find email spoofing vectors during recon
- 🐞 Perfect for Bug Bounty, Red Team and Pentest
- 🧰 Can be piped into toolchains with cat, httpx, dnsx, etc

