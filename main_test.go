package main

import (
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"testing"
	"time"
)

var data_been_setup bool = false

func dataSetup() {
	gRoot = &IPTree{node: nil}
	gASN = make(map[int]*ASN)
	gKey = make(map[string]*Key)

	files := []string{
		"./data/delegated-afrinic-extended-latest",
		"./data/delegated-apnic-extended-latest",
		"./data/delegated-arin-extended-latest",
		"./data/delegated-lacnic-extended-latest",
		"./data/delegated-ripencc-extended-latest",
	}
	for _, filename := range files {
		asns, ipv4s := read(filename)
		add_asns(asns)
		add_ipv4s(ipv4s, gRoot)
	}
	data_been_setup = true
	fmt.Println("Data imported")
}

func BenchmarkRandomLookup(b *testing.B) {
	if !data_been_setup {
		dataSetup()
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	runtime.GC()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		findRandomNet(r)
	}

}

func BenchmarkOneLookup(b *testing.B) {
	if !data_been_setup {
		dataSetup()
	}

	ip := net.IPv4(217, 70, 119, 150).To4()
	runtime.GC()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		findNet(ip)
	}

}

func BenchmarkManyLookups(b *testing.B) {
	if !data_been_setup {
		dataSetup()
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ips := make([]net.IP, b.N)
	for n := 0; n < b.N; n++ {
		ips[n] = net.IPv4(byte(r.Int31n(220)), byte(r.Int31n(256)), byte(r.Int31n(256)), byte(r.Int31n(256))).To4()

	}

	runtime.GC()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		findNet(ips[n])
	}

}

func BenchmarkManyLookupsWithPointers(b *testing.B) {
	if !data_been_setup {
		dataSetup()
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ips := make([]*net.IP, b.N)
	for n := 0; n < b.N; n++ {
		ip := net.IPv4(byte(r.Int31n(220)), byte(r.Int31n(256)), byte(r.Int31n(256)), byte(r.Int31n(256))).To4()
		ips[n] = &ip

	}

	runtime.GC()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		findNet(*ips[n])
	}

}
