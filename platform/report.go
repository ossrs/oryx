// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
package main

import (
	"context"
)

// queryLatestVersion is to query the latest and stable version from Oryx API.
func queryLatestVersion(ctx context.Context) (*Versions, error) {
	return &Versions{
		Version: version,
		Stable:  "v1.0.193",
		Latest:  "v1.0.307",
	}, nil
}
