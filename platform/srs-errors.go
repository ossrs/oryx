//
// Copyright (c) 2022-2026 Winlin
//
// SPDX-License-Identifier: MIT
//
package main

// Oryx uses 100 to 999 for its error codes.
type SrsStackError int

// The callback errors use 100 to 200.
const (
	SrsStackErrorCallbackRecord SrsStackError = 100
)
