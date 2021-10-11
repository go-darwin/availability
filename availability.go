// SPDX-FileCopyrightText: 2021 The Go Darwin Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build darwin && amd64 && gc
// +build darwin,amd64,gc

// package availability provides the darwin Availability for Go.
package availability

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

//go:noescape
//go:nosplit
//go:linkname sysctl runtime.sysctl
func sysctl(mib *uint32, miblen uint32, oldp *byte, oldlenp *uintptr, newp *byte, newlen uintptr) int32

// Darwin C02XL055HX8F 21.0.0 Darwin Kernel Version 21.0.0: Thu May 27 21:01:58 PDT 2021; root:xnu-7938.0.0.111.2~2/RELEASE_X86_64 x86_64 i386 iMacPro1,1 Darwin
// OUT: 200112

// DarwinVersion gets darwin kernel version by using sysctl.
//
// darwin version and macOS Codename table: https://en.wikipedia.org/wiki/MacOS#Release_history.
func DarwinVersion() int {
	ret1, err := unix.Sysctl("kern.osproductversion")
	if err != nil {
		return 0
	}
	fmt.Printf("ret1: %#v\n", ret1)
	return 0

	// Use sysctl to fetch kern.osrelease
	mib := [2]uint32{unix.CTL_KERN, unix.KERN_OSRELEASE}
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
