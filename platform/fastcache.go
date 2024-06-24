// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
package main

import "context"

var fastCache *FastCache

type FastCache struct {
	// Whether delivery HLS in high performance mode.
	HLSHighPerformance bool
	// Whether deliver HLS in low latency mode.
	HLSLowLatency bool
}

func NewFastCache() *FastCache {
	return &FastCache{}
}

func (v *FastCache) Refresh(ctx context.Context) error {
	if vs, _ := rdb.HGet(ctx, SRS_LL_HLS, "hlsLowLatency").Result(); vs == "true" {
		v.HLSLowLatency = true
	} else {
		v.HLSLowLatency = false
	}

	if vs, _ := rdb.HGet(ctx, SRS_HP_HLS, "noHlsCtx").Result(); vs == "true" {
		v.HLSHighPerformance = true
	} else {
		v.HLSHighPerformance = false
	}

	return nil
}
