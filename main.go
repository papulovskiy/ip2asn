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
	var date int64
	if parts[6] == "allocated" {
		date, err = strconv.ParseInt(parts[5], 10, 32)
		check(err)
	} else {
		date = 0
	}
	return &ASN{
		number:  int(number),
		status:  parts[6],
		date:    int(date),
		country: parts[1],
		source:  parts[0],
		key:     parts[7],
	}
}

func make_ipv4(parts []string) *IPv4Network {
	number, err := strconv.ParseInt(parts[4], 10, 32)
	check(err)
	mask := 32 - int(math.Log(float64(number))/math.Log(2.0))

	var date int64
	if parts[6] == "allocated" {
		date, err = strconv.ParseInt(parts[5], 10, 32)
		check(err)
	} else {
		date = 0
	}

	return &IPv4Network{
		address: net.ParseIP(parts[3]).To4(),
		mask:    mask,
		status:  parts[6],
		date:    int(date),
		country: parts[1],
		source:  parts[0],
		key:     parts[7],
	}
}

// func make_ipv6(parts []string) *IPv6Network {
// 	number, err := strconv.ParseInt(parts[4], 10, 32)
// 	check(err)
// 	mask := int(math.Log(float64(number))/math.Log(2.0))

// 	var date int64
// 	if parts[6] == "allocated" {
// 		date, err = strconv.ParseInt(parts[5], 10, 32)
// 		check(err)
// 	} else {
// 		date = 0
// 	}

// 	return &IPv6Network{
// 		address: net.ParseIP(parts[3]).To16(),
// 		mask:    mask,
// 		status:  parts[6],
// 		date:    int(date),
// 		country: parts[1],
// 		source:  parts[0],
// 		key:     parts[7],
// 	}
// }

func read(filename string) {
	var asns []*ASN
	var ipv4s []*IPv4Network
	// var ipv6s []*IPv6Network

	var asns_count, ipv4s_count, ipv6s_count int64
	var current_asn, current_ipv4, current_ipv6 int = 0, 0, 0

	file, err := os.Open(filename)
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		parts := strings.Split(scanner.Text(), "|")
		if parts[0] == "2" {
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
			case "ipv6":
				ipv6s_count, err = strconv.ParseInt(parts[4], 10, 64)
				check(err)
				// ipv6s = make([]*IPv6Network, ipv6s_count)
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
		case "ipv6":
			// ipv6s[current_ipv6] = make_ipv6(parts)
			current_ipv6++
		}
	}

	// Just to avoid compiler complaints
	fmt.Println(asns_count)
	fmt.Println(ipv4s_count)
	fmt.Println(ipv6s_count)
	fmt.Println(asns[0])
	fmt.Println(ipv4s[0])
	// fmt.Println(ipv6s[0])

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	read("./data/delegated-afrinic-extended-latest")

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
