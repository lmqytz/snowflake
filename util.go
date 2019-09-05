package main

import (
	"net"
	"strings"
)

func getWorkId(netCard string) (workerId string, err error) {
	ipv4, err := getNetCardIpv4(netCard)
	if err != nil {
		return "", err
	}

	return strings.Split(ipv4, ".")[3], nil
}

func getNetCardIpv4(name string) (ipv4 string, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range interfaces {
		if i.Name != name {
			continue
		}

		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return "", err
		}

		addresses, err := byName.Addrs()
		if err != nil {
			return "", err
		}

		for _, v := range addresses {
			if ip4 := v.(*net.IPNet).IP.To4(); ip4 != nil {
				return ip4.String(), nil
			}
		}
	}

	return "", nil
}
