// Copyright 2018 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// +build linux

package mtu

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

const (
	externalProbeIPv4 = "1.1.1.1"
	externalProbeIPv6 = "2606:4700:4700::1111"
)

func getRoute(externalProbe string) ([]netlink.Route, error) {
	ip := net.ParseIP(externalProbe)
	if ip == nil {
		return nil, fmt.Errorf("unable to parse IP %s", externalProbe)
	}

	routes, err := netlink.RouteGet(ip)
	if err != nil {
		return nil, fmt.Errorf("unable to lookup route to %s: %s", externalProbe, err)
	}

	if len(routes) == 0 {
		return nil, fmt.Errorf("no route to %s", externalProbe)
	}

	return routes, nil
}

func autoDetect() (int, error) {
	var routes []netlink.Route
	var err error

	routes, err = getRoute(externalProbeIPv4)
	if err != nil {
		prevErr := err
		routes, err = getRoute(externalProbeIPv6)
		if err != nil {
			return 0, fmt.Errorf("%v, %v", err.Error(), prevErr.Error())
		}
	}

	link, err := netlink.LinkByIndex(routes[0].LinkIndex)
	if err != nil {
		return 0, fmt.Errorf("unable to find interface of default route: %s", err)
	}

	if mtu := link.Attrs().MTU; mtu != 0 {
		log.Infof("Detected MTU %d", mtu)
		return mtu, nil
	}

	return EthernetMTU, nil
}
