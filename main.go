package main

import (
	"fmt"
	"net"
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

func main() {
	files := []string{
		"./data/delegated-afrinic-extended-latest",
		"./data/delegated-apnic-extended-latest",
		"./data/delegated-arin-extended-latest",
		"./data/delegated-lacnic-extended-latest",
		"./data/delegated-ripencc-extended-latest",
	}
	for _, filename := range files {
		fmt.Println("Reading", filename)
		asns, ipv4s := read(filename)
		fmt.Printf("Got %v ASNs and %v IPv4 networks\n", len(asns), len(ipv4s))
	}

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
