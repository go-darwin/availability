// SPDX-FileCopyrightText: 2021 The Go Darwin Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build darwin && amd64 && gc
// +build darwin,amd64,gc

// package availability provides the darwin Availability for Go.
package availability

import (
	"unsafe"
)

//go:noescape
//go:linkname sysctl runtime.sysctl
func sysctl(mib *uint32, miblen uint32, oldp *byte, oldlenp *uintptr, newp *byte, newlen uintptr) int32

const (
	_CTL_KERN       = 1
	_KERN_OSRELEASE = 2
)

// DarwinVersion gets darwin kernel version by using sysctl.
//
// darwin version and macOS Codename table: https://en.wikipedia.org/wiki/MacOS#Release_history.
func DarwinVersion() int {
	// Use sysctl to fetch kern.osrelease
	mib := [2]uint32{_CTL_KERN, _KERN_OSRELEASE}
	var out [32]byte
	nout := unsafe.Sizeof(out)

	ret := sysctl(&mib[0], 2, (*byte)(unsafe.Pointer(&out)), &nout, nil, 0)
	if ret >= 0 {
		ver := 0
		for i := 0; i < int(nout) && out[i] >= '0' && out[i] <= '9'; i++ {
			ver *= 10
			ver += int(out[i] - '0')
		}

		return ver
	}

	return 0 // unreachable
}
