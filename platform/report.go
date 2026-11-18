// Copyright (c) 2022-2026 Winlin
//
// SPDX-License-Identifier: MIT
package main

import (
	"context"
)

func queryLatestVersion(ctx context.Context) (*Versions, error) {
	return &Versions{
		Version: version,
		Stable:  "v1.0.193",
		Latest:  "v1.0.307",
	}, nil
}
