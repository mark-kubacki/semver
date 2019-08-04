// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "textflag.h"

TEXT ·less(SB),NOSPLIT,$0-17
	MOVQ	t+0(FP), DI // Flip SI and DI because we use "greater than".
	MOVQ	o+8(FP), SI

	MOVOU	(SI), X2
	MOVOU	(DI), X5
	// PCMPGTD X5, X2
	BYTE $0x66; BYTE $0x0f; BYTE $0x66; BYTE $0xd5
	PMOVMSKB X2, AX
	// TEST	AX, AX
	BYTE $0x85; BYTE $0xc0
	JNE	is_less

	MOVOU	16(SI), X3
	MOVOU	16(DI), X6
	// PCMPGTD X6, X3
	BYTE $0x66; BYTE $0x0f; BYTE $0x66; BYTE $0xde
	PMOVMSKB X3, AX
	// TEST	AX, AX
	BYTE $0x85; BYTE $0xc0
	JNE	is_less

	MOVOU	32(SI), X4
	MOVOU	32(DI), X7
	// PCMPGTD X7, X4
	BYTE $0x66; BYTE $0x0f; BYTE $0x66; BYTE $0xe7
	PMOVMSKB X4, AX
	// TEST	AX, AX
	BYTE $0x85; BYTE $0xc0
	JNE	is_less

	MOVOU	48(SI), X0
	MOVOU	48(DI), X1
	// PCMPGTD X1, X0
	BYTE $0x66; BYTE $0x0f; BYTE $0x66; BYTE $0xc1
	PMOVMSKB X0, AX
// Now comes an exception: We over-read to catch the adjacent 'build'; which lives in t[len-1]+1
// because Go allocates and aligns everything on lines of 8 byte.
// Good thing is we need to compare that anyway, but have to filter out the unclaimed space.
// TEST is bitwise-and anyway, so run with 0x0fff.
	// TEST	0x0fff, AX
	BYTE $0xa9; BYTE $0xff; BYTE $0x0f; BYTE $0x00; BYTE $0x00
	JNE	is_less

not_less:
	MOVB	$0, ret+16(FP)
	RET

is_less:
	MOVB	$1, ret+16(FP)
	RET
