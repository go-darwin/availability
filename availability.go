// Copyright 2021 The Go Darwin Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build darwin
// +build darwin

// package availability provides the darwin Availability for Go.
package availability

import "unsafe"

//go:noescape
//go:linkname sysctlbyname runtime.sysctlbyname
func sysctlbyname(name *byte, oldp *byte, oldlenp *uintptr, newp *byte, newlen uintptr) int32

// byte slice of "kern.osproductversion".
var kernOSProductVersion = []byte{0x6b, 0x65, 0x72, 0x6e, 0x2e, 0x6f, 0x73, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e}

// DarwinVersion gets darwin kernel version by using sysctl.
//
// The darwin version and macOS Codename table: https://en.wikipedia.org/wiki/MacOS#Release_history.
//go:nosplit
func DarwinVersion() (ver int32) {
	var out [16]byte
	nout := unsafe.Sizeof(out)
	sysctlbyname(&kernOSProductVersion[0], (*byte)(unsafe.Pointer(&out)), &nout, nil, 0)

	const step = int32(10)
	for i := 0; i < int(nout) && (out[i] >= '0' && out[i] <= '9' || out[i] == '.'); i++ {
		if out[i] == '.' {
			continue
		}
		ver *= step
		ver += int32(out[i] - '0')
	}

	n := 7 - int(nout) // 7 is Availability digits
	for n >= 0 {
		ver *= step
		n--
	}

	return ver
}
