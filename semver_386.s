// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !purego
// +build !go1.17

#include "go_asm.h"
#include "textflag.h"

TEXT ·compare(SB),NOSPLIT,$0-12
	MOVL	t+0(FP), SI
	MOVL	o+4(FP), DI
	XORL	CX, CX		// Index of the last examined element.

	MOVOU	(SI), X2
	MOVOU	(DI), X5
	PCMPEQL	X5, X2
	MOVMSKPS X2, AX
	CMPL	AX, $0x0f
	JNE	diff
	MOVL	$4, CX

	MOVOU	16(SI), X3
	MOVOU	16(DI), X6
	PCMPEQL	X6, X3
	MOVMSKPS X3, AX
	CMPL	AX, $0x0f
	JNE	diff
	MOVL	$8, CX

	MOVOU	32(SI), X4
	MOVOU	32(DI), X7
	PCMPEQL	X7, X4
	MOVMSKPS X4, AX
	CMPL	AX, $0x0f
	JNE	diff
	MOVL	$12, CX

	MOVOU	48(SI), X0
	MOVOU	48(DI), X1
	PCMPEQL	X1, X0
	MOVMSKPS X0, AX
	ORL	$0xc, AX // Mask undefined space, due to 'build' and then nothing.
	CMPL	AX, $0x0f
	JNE	diff

equal:
	MOVL	$0, ret+8(FP)
	RET

diff:
	ORL	$0xfff0, AX	// See step below. These are unrelated and will be zeros.
	XORL	$0xffff, AX	// Invert mask from "equal" to "differ".
	BSFL	AX, BX		// Number of the first bit 1 from LSB on counted.
	XORL	AX, AX
	ADDL	BX, CX
	// Now compare those diverging elements. (AX, BX, DX are free)
	MOVL	(DI)(CX*4), BX
	CMPL	BX, (SI)(CX*4)
	SETLT	AX
	LEAL	-1(AX*2), AX
	MOVL	AX, ret+8(FP)
	RET

TEXT ·less(SB),NOSPLIT,$0-9
	MOVL	t+0(FP), SI
	MOVL	o+4(FP), DI

	XORL	DX, DX
	XORL	BX, BX
less_loop:
	MOVOU	(DI)(DX*1), X1
	MOVOU	(SI)(DX*1), X0

	MOVAPS	X1, X3
	PCMPEQL X0, X1
	MOVMSKPS	X1, AX
	CMPB	AX, $0x0f
	JNE less_determine
	ADDB	$16, DX
	CMPB	DX, $64
	JE	less_eol
	JMP	less_loop

less_determine:
	MOVAPS	X3, X1
	PCMPGTL	X0, X3		// 3.0.1.0 |>| 2.1.0.0 -> 1.0.1.0
	PCMPGTL	X1, X0		// 2.1.0.0 |>| 3.0.1.0 -> 0.1.0.0
	PSHUFL	$27, X3, X3	// $27 is [0, 1, 2, 3], reverse order of elements to get a workable mask below.
	PSHUFL	$27, X0, X0
	MOVMSKPS X3, CX		// 1010
	MOVMSKPS X0, AX		// 0100
	CMPB	CX, AX
	SETGT	BX
less_eol:
	MOVB	BX, ret+8(FP)
	RET
