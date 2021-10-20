// Copyright 2021 The Go Darwin Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build darwin && go1.12
// +build darwin,go1.12

#include "textflag.h"

TEXT CFDataGetLength_trampoline<>(SB),NOSPLIT,$0-0
	JMP	CFDataGetLength(SB)

GLOBL	·CFDataGetLength_trampoline_addr(SB), RODATA, $8
DATA	·CFDataGetLength_trampoline_addr(SB)/8, $CFDataGetLength_trampoline<>(SB)

TEXT CFDataGetBytePtr_trampoline<>(SB),NOSPLIT,$0-0
	JMP	CFDataGetBytePtr(SB)

GLOBL	·CFDataGetBytePtr_trampoline(SB), RODATA, $8
DATA	·CFDataGetBytePtr_trampoline_addr(SB)/8, $CFDataGetBytePtr_trampoline<>(SB)

TEXT CFStringCreateWithBytes_trampoline<>(SB),NOSPLIT,$0-0
	JMP	CFStringCreateWithBytes(SB)

GLOBL	·CFStringCreateWithBytes_trampoline_addr(SB), RODATA, $8
DATA	·CFStringCreateWithBytes_trampoline_addr(SB)/8, $CFStringCreateWithBytes_trampoline<>(SB)

TEXT CFCopySystemVersionDictionary_trampoline<>(SB),NOSPLIT,$0-0
	JMP	CFCopySystemVersionDictionary(SB)

GLOBL	·CFCopySystemVersionDictionary_trampoline_addr(SB), RODATA, $8
DATA	·CFCopySystemVersionDictionary_trampoline_addr(SB)/8, $CFCopySystemVersionDictionary_trampoline<>(SB)

TEXT CFDictionaryGetValue_trampoline<>(SB),NOSPLIT,$0-0
	JMP	CFDictionaryGetValue(SB)

GLOBL	·CFDictionaryGetValue_trampoline_addr(SB), RODATA, $8
DATA	·CFDictionaryGetValue_trampoline_addr(SB)/8, $CFDictionaryGetValue_trampoline<>(SB)

TEXT CFRelease_trampoline<>(SB),NOSPLIT,$0-0
	JMP	CFRelease(SB)

GLOBL	·CFRelease_trampoline_addr(SB), RODATA, $8
DATA	·CFRelease_trampoline_addr(SB)/8, $CFRelease_trampoline<>(SB)

TEXT CFStringGetCString_trampoline<>(SB),NOSPLIT,$0-0
	JMP	CFStringGetCString(SB)

GLOBL	·CFStringGetCString_trampoline_addr(SB), RODATA, $8
DATA	·CFStringGetCString_trampoline_addr(SB)/8, $CFStringGetCString_trampoline<>(SB)

TEXT CFStringGetMaximumSizeForEncoding_trampoline<>(SB),NOSPLIT,$0-0
	JMP	CFStringGetMaximumSizeForEncoding(SB)

GLOBL	·CFStringGetMaximumSizeForEncoding_trampoline_addr(SB), RODATA, $8
DATA	·CFStringGetMaximumSizeForEncoding_trampoline_addr(SB)/8, $CFStringGetMaximumSizeForEncoding_trampoline<>(SB)

TEXT CFStringGetSystemEncoding_trampoline<>(SB),NOSPLIT,$0-0
	JMP	CFStringGetSystemEncoding(SB)

GLOBL	·CFStringGetSystemEncoding_trampoline_addr(SB), RODATA, $8
DATA	·CFStringGetSystemEncoding_trampoline_addr(SB)/8, $CFStringGetSystemEncoding_trampoline<>(SB)

TEXT CFStringGetLength_trampoline<>(SB),NOSPLIT,$0-0
	JMP	CFStringGetLength(SB)

GLOBL	·CFStringGetLength_trampoline_addr(SB), RODATA, $8
DATA	·CFStringGetLength_trampoline_addr(SB)/8, $CFStringGetLength_trampoline<>(SB)
