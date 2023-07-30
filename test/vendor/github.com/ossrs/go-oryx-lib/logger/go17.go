// The MIT License (MIT)
//
// Copyright (c) 2013-2017 Oryx(ossrs)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// +build go1.7

package logger

import (
	"context"
	"fmt"
	"os"
)

func (v *loggerPlus) Println(ctx Context, a ...interface{}) {
	args := v.contextFormat(ctx, a...)
	v.doPrintln(args...)
}

func (v *loggerPlus) Printf(ctx Context, format string, a ...interface{}) {
	format, args := v.contextFormatf(ctx, format, a...)
	v.doPrintf(format, args...)
}

func (v *loggerPlus) contextFormat(ctx Context, a ...interface{}) []interface{} {
	if ctx, ok := ctx.(context.Context); ok {
		if cid, ok := ctx.Value(cidKey).(int); ok {
			return append([]interface{}{fmt.Sprintf("[%v][%v]", os.Getpid(), cid)}, a...)
		}
	} else {
		return v.format(ctx, a...)
	}
	return a
}

func (v *loggerPlus) contextFormatf(ctx Context, format string, a ...interface{}) (string, []interface{}) {
	if ctx, ok := ctx.(context.Context); ok {
		if cid, ok := ctx.Value(cidKey).(int); ok {
			return "[%v][%v] " + format, append([]interface{}{os.Getpid(), cid}, a...)
		}
	} else {
		return v.formatf(ctx, format, a...)
	}
	return format, a
}

// User should use context with value to pass the cid.
type key string

var cidKey key = "cid.logger.ossrs.org"

var gCid int = 999

// Create context with value.
func WithContext(ctx context.Context) context.Context {
	gCid += 1
	return context.WithValue(ctx, cidKey, gCid)
}

// Create context with value from parent, copy the cid from source context.
// @remark Create new cid if source has no cid represent.
func AliasContext(parent context.Context, source context.Context) context.Context {
	if source != nil {
		if cid, ok := source.Value(cidKey).(int); ok {
			return context.WithValue(parent, cidKey, cid)
		}
	}
	return WithContext(parent)
}
