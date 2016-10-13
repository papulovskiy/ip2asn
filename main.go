package main

import (
	"bufio"
	"fmt"
	// "io"
	// "io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type ASN struct {
	number  int
	status  string
	date    int
	country string
	source  string
	key     string
}

type IPv4Network struct {
	address net.IP
	mask    int
	status  string
	date    int
	country string
	source  string
	key     string
}

type IPv6Network struct {
	address net.IP
	mask    int
	status  string
	date    int
	country string
	source  string
	key     string
}

func make_asn(parts []string) *ASN {
	number, err := strconv.ParseInt(parts[3], 10, 32)
	check(err)

	var key string
	var date int64
	if parts[6] == "allocated" {
		date, err = strconv.ParseInt(parts[5], 10, 32)
		check(err)
		key = parts[7]
	} else {
		date = 0
		key = ""
	}
	return &ASN{
		number:  int(number),
		status:  parts[6],
		date:    int(date),
		country: parts[1],
		source:  parts[0],
		key:     key,
	}
}

func make_ipv4(parts []string) *IPv4Network {
	number, err := strconv.ParseInt(parts[4], 10, 32)
	check(err)
	mask := 32 - int(math.Log(float64(number))/math.Log(2.0))

	var key string
	var date int64
	if parts[6] == "allocated" || parts[6] == "assigned" {
		date, err = strconv.ParseInt(parts[5], 10, 32)
		check(err)
		key = parts[7]
	} else {
		date = 0
		key = ""
	}

	return &IPv4Network{
		address: net.ParseIP(parts[3]).To4(),
		mask:    mask,
		status:  parts[6],
		date:    int(date),
		country: parts[1],
		source:  parts[0],
		key:     key,
	}
}

func read(filename string) {
	var asns []*ASN
	var ipv4s []*IPv4Network

	var asns_count, ipv4s_count int64
	var current_asn, current_ipv4 int = 0, 0

	file, err := os.Open(filename)
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		parts := strings.Split(scanner.Text(), "|")
		if len(parts) < 5 || parts[0] == "2" {
			continue
		}
		if parts[5] == "summary" {
			switch parts[2] {
			case "asn":
				asns_count, err = strconv.ParseInt(parts[4], 10, 64)
				check(err)
				asns = make([]*ASN, asns_count)
			case "ipv4":
				ipv4s_count, err = strconv.ParseInt(parts[4], 10, 64)
				check(err)
				ipv4s = make([]*IPv4Network, ipv4s_count)
			}
			continue
		}
		switch parts[2] {
		case "asn":
			asns[current_asn] = make_asn(parts)
			current_asn++
		case "ipv4":
			ipv4s[current_ipv4] = make_ipv4(parts)
			current_ipv4++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Reading")
	read("./data/delegated-afrinic-extended-latest")
	read("./data/delegated-apnic-extended-latest")
	read("./data/delegated-arin-extended-latest")
	read("./data/delegated-lacnic-extended-latest")
	read("./data/delegated-ripencc-extended-latest")

}

// ftp://ftp.ripe.net/ripe/stats/RIR-Statistics-Exchange-Format.txt

// 3.3 Record format:

// After the defined file header, and excluding any space or
// comments, each line in the file represents a single allocation
// (or assignment) of a specific range of Internet number resources
// (IPv4, IPv6 or ASN), made by the RIR identified in the record.

// In the case of IPv4 the records may represent non-CIDR ranges
// or CIDR blocks, and therefore the record format represents a
// beginning of range, and a count. This can be converted to
// prefix/length using simple algorithms.

// In the case of IPv6 the record format represents the prefix
// and the count of /128 instances under that prefix.
