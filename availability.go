// Copyright 2021 The Go Darwin Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build darwin
// +build darwin

// package availability provides the darwin Availability for Go.
package availability

import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"

	"github.com/go-darwin/sys"
)

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

//go:linkname syscall crypto/x509/internal/macos.syscall
func syscall(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1 uintptr)

// CFRef is an opaque reference to a Core Foundation object. It is a pointer,
// but to memory not owned by Go, so not an unsafe.Pointer.
type CFRef uintptr

type CFString CFRef

const kCFAllocatorDefault = 0
const kCFStringEncodingUTF8 = 0x08000100
const kCFStringEncodingASCII = 0x0600

// CFDataToSlice returns a copy of the contents of data as a bytes slice.
func CFDataToSlice(data CFRef) []byte {
	length := CFDataGetLength(data)
	ptr := CFDataGetBytePtr(data)
	src := (*[1 << 20]byte)(unsafe.Pointer(ptr))[:length:length]
	out := make([]byte, length)
	copy(out, src)
	return out
}

func CFDataGetLength(data CFRef) int {
	ret := syscall(CFDataGetLength_trampoline_addr, uintptr(data), 0, 0, 0, 0, 0)
	return int(ret)
}

//go:cgo_import_dynamic CFDataGetLength CFDataGetLength "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var CFDataGetLength_trampoline_addr uintptr

func CFDataGetBytePtr(data CFRef) uintptr {
	ret := syscall(CFDataGetBytePtr_trampoline_addr, uintptr(data), 0, 0, 0, 0, 0)
	return ret
}

//go:cgo_import_dynamic CFDataGetBytePtr_trampoline CFDataGetBytePtr_trampoline "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var CFDataGetBytePtr_trampoline_addr uintptr

// StringToCFString returns a copy of the UTF-8 contents of s as a new CFString.
func StringToCFString(s string) CFString {
	p := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data)
	ret := syscall(CFStringCreateWithBytes_trampoline_addr, kCFAllocatorDefault, uintptr(p),
		uintptr(len(s)), uintptr(kCFStringEncodingUTF8), 0 /* isExternalRepresentation */, 0)
	runtime.KeepAlive(p)
	return CFString(ret)
}

var CFStringCreateWithBytes_trampoline_addr uintptr

//go:cgo_import_dynamic CFStringCreateWithBytes CFStringCreateWithBytes "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

func CFCopySystemVersionDictionary() (CFRef, error) {
	ret, _, errno := sys.Syscall(CFCopySystemVersionDictionary_trampoline_addr, 0, 0, 0)
	if errno != 0 {
		return 0, sys.Errno(errno)
	}

	return CFRef(ret), nil
}

//go:cgo_import_dynamic CFCopySystemVersionDictionary _CFCopySystemVersionDictionary "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var CFCopySystemVersionDictionary_trampoline_addr uintptr

// CFDictionaryGetValue retrieves the value associated with the given key.
func CFDictionaryGetValue(dict CFRef, key CFString) (CFString, error) {
	ret, _, errno := sys.Syscall(CFDictionaryGetValue_trampoline_addr, uintptr(dict), uintptr(key), 0)
	fmt.Printf("errno: %#v\n", errno)
	if errno != 0 {
		return 0, sys.Errno(errno)
	}

	return CFString(ret), nil
}

//go:cgo_import_dynamic CFDictionaryGetValue CFDictionaryGetValue "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var CFDictionaryGetValue_trampoline_addr uintptr

func CFRelease(ref CFRef) {
	syscall(CFRelease_trampoline_addr, uintptr(ref), 0, 0, 0, 0, 0)
}

//go:cgo_import_dynamic CFRelease CFRelease "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var CFRelease_trampoline_addr uintptr

func CFStringGetCString(str CFRef) ([]byte, uintptr, error) {
	var cstr [256]byte
	_, _, errno := sys.RawSyscall6(CFStringGetCString_trampoline_addr, uintptr(str), uintptr(unsafe.Pointer(&cstr[0])), unsafe.Sizeof(cstr), uintptr(kCFStringEncodingUTF8), 0, 0)
	if errno != 0 {
		return nil, 0, sys.Errno(errno)
	}
	out := CFDataToSlice(str)

	cstrlen := CFStringGetMaximumSizeForEncoding(CFStringGetLength(str), CFStringEncoding(kCFStringEncodingUTF8))
	cstrlen = cstrlen + CFIndex(1)

	return out, uintptr(cstrlen), nil
}

//go:cgo_import_dynamic CFStringGetCString CFStringGetCString "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var CFStringGetCString_trampoline_addr uintptr

type CFIndex CFRef

func CFStringGetMaximumSizeForEncoding(length CFIndex, encoding CFStringEncoding) CFIndex {
	ret, _, errno := sys.Syscall(CFStringGetMaximumSizeForEncoding_trampoline_addr, uintptr(length), uintptr(encoding), 0)
	if errno != 0 {
		return 0
	}
	fmt.Printf("ret: %#v\n", CFIndex(unsafe.Pointer(&ret)))

	return CFIndex(ret)
}

//go:cgo_import_dynamic CFStringGetMaximumSizeForEncoding CFStringGetMaximumSizeForEncoding "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var CFStringGetMaximumSizeForEncoding_trampoline_addr uintptr

type CFStringEncoding CFRef

func CFStringGetSystemEncoding() CFStringEncoding {
	ret, _, errno := sys.Syscall(CFStringGetSystemEncoding_trampoline_addr, 0, 0, 0)
	if errno != 0 {
		return 0
	}

	return CFStringEncoding(ret)
}

//go:cgo_import_dynamic CFStringGetSystemEncoding CFStringGetSystemEncoding "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var CFStringGetSystemEncoding_trampoline_addr uintptr

func CFStringGetLength(str CFRef) CFIndex {
	ret, _, errno := sys.Syscall(CFStringGetLength_trampoline_addr, uintptr(str), 0, 0)
	if errno != 0 {
		return 0
	}

	return CFIndex(ret)
}

//go:cgo_import_dynamic CFStringGetLength CFStringGetLength "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var CFStringGetLength_trampoline_addr uintptr

//go:cgo_import_dynamic _kCFSystemVersionProductVersionKey _kCFSystemVersionProductVersionKey "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

//go:linkname _kCFSystemVersionProductVersionKey _kCFSystemVersionProductVersionKey
var _kCFSystemVersionProductVersionKey byte

// On Darwin, we set MacOSXProductVersion to the corresponding OS X release.
// This is for compatibility with scripts that set MACOSX_DEPLOYMENT_TARGET
// based on sw_vers -productVersion
var MacOSXProductVersionKey = StringToCFString("MacOSXProductVersion")

func Version() (ver int32) {
	dict, err := CFCopySystemVersionDictionary()
	if err != nil {
		return 0
	}
	defer CFRelease(dict)

	cfnum, err := CFDictionaryGetValue(dict, MacOSXProductVersionKey)
	if err != nil {
		return 0
	}

	out, nout, err := CFStringGetCString(CFRef(cfnum))
	if err != nil || nout == 0 {
		return 0
	}

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
