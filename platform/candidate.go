//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
package main

import (
	"context"
	"net"
	"os"
	"strings"
	"sync"
	// From ossrs.
	"github.com/ossrs/go-oryx-lib/logger"
)

var candidateWorker *CandidateWorker

type CandidateWorker struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewCandidateWorker() *CandidateWorker {
	return &CandidateWorker{}
}

func (v *CandidateWorker) Close() error {
	if v.cancel != nil {
		v.cancel()
	}
	v.wg.Wait()
	return nil
}

func (v *CandidateWorker) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	v.cancel = cancel

	ctx = logger.WithContext(ctx)
	logger.Tf(ctx, "candidate start a worker")

	return nil
}

// Resolve host to ip. Return nil if ignore the host resolving, for example, user disable resolving by
// set env NAME_LOOKUP to false.
func (v *CandidateWorker) Resolve(host string) (net.IP, error) {
	// Ignore the resolving.
	if os.Getenv("NAME_LOOKUP") == "off" {
		return nil, nil
	}

	// Ignore the port.
	if strings.Contains(host, ":") {
		if hostname, _, err := net.SplitHostPort(host); err != nil {
			return nil, err
		} else {
			host = hostname
		}
	}

	// Resolve the localhost to possible IP address.
	if host == "localhost" {
		// If directly run in host, like debugging, use the private ipv4.
		if os.Getenv("PLATFORM_DOCKER") == "off" {
			return conf.ipv4, nil
		}

		// Return lo for OBS WHIP or native client to access it.
		return net.IPv4(127, 0, 0, 1), nil
	}

	// Directly use the ip if not name.
	if ip := net.ParseIP(host); ip != nil {
		return ip, nil
	}

	// Lookup the name to parse to ip.
	if ips, err := net.LookupIP(host); err != nil {
		return nil, err
	} else if len(ips) > 0 {
		return ips[0], nil
	}

	return nil, nil
}
