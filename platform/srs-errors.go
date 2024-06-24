//
// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
//
package main

// General error code for Oryx, it uses 100~999 for all errors.
type SrsStackError int

// Error code for callback, 100 ~ 200.
const (
	// Error for callback module, about the record events.
	SrsStackErrorCallbackRecord SrsStackError = 100
)
