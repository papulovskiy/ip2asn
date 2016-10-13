package main

import (
	"encoding/binary"
	// "fmt"
	"math/rand"
	"net"
)

type Key struct {
	asns  []int
	ipv4s []*IPv4Network
}

func add_asns(asns []*ASN) {
	for _, as := range asns {
		// Store as it is
		gASN[as.number] = as
		// Create relation with key
		_, ok := gKey[as.key]
		if ok {
		} else {
			gKey[as.key] = &Key{}
		}
		gKey[as.key].asns = append(gKey[as.key].asns, as.number)
	}
	// fmt.Printf("Total count of AS: %v\n", len(gASN))
}

type IPTree struct {
	zero *IPTree
	one  *IPTree
	node *IPv4Network
}

func ip2int(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip)
}

func add_ipv4s(ipv4s []*IPv4Network, root *IPTree) {
	for _, net := range ipv4s {
		current := root
		// Store in a tree
		intip := ip2int(net.address)
		var i uint
		// Traverse tree
		for i = 0; i < uint(net.mask); i++ {
			bit := uint32(1 << uint(32-i))
			if intip&bit != 0 {
				if current.one == nil {
					// fmt.Println("Creating one node")
					current.one = &IPTree{}
				}
				current = current.one
			} else {
				if current.zero == nil {
					// fmt.Println("Creating zero node")
					current.zero = &IPTree{}
				}
				current = current.zero
			}
		}
		current.node = net
		// Create relation with key
		_, ok := gKey[net.key]
		if ok {
		} else {
			gKey[net.key] = &Key{}
		}
		gKey[net.key].ipv4s = append(gKey[net.key].ipv4s, net)
	}
}

func findNet(ip net.IP) *IPv4Network {
	intip := ip2int(ip)
	current := gRoot

	for i := 0; i < 32; i++ {
		bit := uint32(1 << uint(32-i))
		if intip&bit != 0 {
			if current.one != nil {
				current = current.one
			} else {
				return current.node
			}
		} else {
			if current.zero != nil {
				current = current.zero
			} else {
				return current.node
			}
		}
	}

	// fmt.Println(intip)
	return &IPv4Network{}
}

func findRandomNet(r *rand.Rand) *IPv4Network {
	a := r.Int31n(220)
	b := r.Int31n(256)
	c := r.Int31n(256)
	d := r.Int31n(256)
	ip := net.IPv4(byte(a), byte(b), byte(c), byte(d))
	return findNet(ip.To4())
}
