// Copyright 2021 The Go Darwin Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build darwin
// +build darwin

package availability

import (
	"bytes"
	"os/exec"
	"testing"
)

func swvers() (ver int32) {
	cmd := exec.Command("sw_vers", "-productVersion")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0
	}
	parts := bytes.Split(out, []byte("."))

	const step = 100 // decrease 100 division
	base := 10000    // started at major version
	for _, part := range parts {
		ver += int32(atoi(string(part)) * base)
		base = base / step
	}

	return ver
}

const (
	maxUint = ^uint(0)
	maxInt  = int(maxUint >> 1)
)

// atoi parses an int from a string s.
// The bool result reports whether s is a number
// representable by a value of type int.
func atoi(s string) int {
	if s == "" {
		return 0
	}

	neg := false
	if s[0] == '-' {
		neg = true
		s = s[1:]
	}

	un := uint(0)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0
		}
		if un > maxUint/10 {
			// overflow
			return 0
		}
		un *= 10
		un1 := un + uint(c) - '0'
		if un1 < un {
			// overflow
			return 0
		}
		un = un1
	}

	if !neg && un > uint(maxInt) {
		return 0
	}
	if neg && un > uint(maxInt)+1 {
		return 0
	}

	n := int(un)
	if neg {
		n = -n
	}

	return n
}

func TestDarwinVersion(t *testing.T) {
	tests := map[string]struct {
		want int32
	}{
		"host": {
			want: swvers(),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := DarwinVersion()
			if got != tt.want {
				t.Fatalf("DarwinVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

var n int32

func BenchmarkDarwinVersion(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n += DarwinVersion()
	}
}
