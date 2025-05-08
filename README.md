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

