// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "textflag.h"

TEXT ·compare(SB),NOSPLIT,$0-24
	MOVQ	t+0(FP), SI
	MOVQ	o+8(FP), DI
	XORQ	CX, CX		// Index of the last examined element.

	MOVOU	(SI), X2
	MOVOU	(DI), X5
	PCMPEQL	X5, X2
	MOVMSKPS X2, AX
	CMPL	AX, $0x0f
	JNE	diff
	MOVQ	$4, CX

	MOVOU	16(SI), X3
	MOVOU	16(DI), X6
	PCMPEQL	X6, X3
	MOVMSKPS X3, AX
	CMPL	AX, $0x0f
	JNE	diff
	MOVQ	$8, CX

	MOVOU	32(SI), X4
	MOVOU	32(DI), X7
	PCMPEQL	X7, X4
	MOVMSKPS X4, AX
	CMPL	AX, $0x0f
	JNE	diff
	MOVQ	$12, CX

	MOVOU	48(SI), X0
	MOVOU	48(DI), X1
	PCMPEQL	X1, X0
	MOVMSKPS X0, AX
	ORQ	$0xc, AX // Mask undefined space, due to 'build' and then nothing.
	CMPL	AX, $0x0f
	JNE	diff

equal:
	MOVQ	$0, ret+16(FP)
	RET

diff:
	ORQ	$0xfff0, AX	// See step below. These are unrelated and will be zeros.
	XORQ	$0xffff, AX	// Invert mask from "equal" to "differ".
	BSFQ	AX, BX		// Number of the first bit 1 from LSB on counted.
	XORQ	AX, AX
	ADDQ	BX, CX
	// Now compare those diverging elements. (AX, BX, DX are free)
	MOVL	(DI)(CX*4), BX
	CMPL	BX, (SI)(CX*4)
	SETLT	AX
	LEAQ	-1(AX*2), AX
	MOVQ	AX, ret+16(FP)
	RET

TEXT ·less(SB),NOSPLIT,$0-17
	MOVQ	t+0(FP), SI
	MOVQ	o+8(FP), DI

	XORQ	DX, DX
less_loop:
	MOVOU	(DI)(DX*1), X4
	MOVOU	(SI)(DX*1), X5
	PSHUFL	$27, X4, X1	// $27 is [0, 1, 2, 3], reverse order of elements to get a workable mask below.
	PSHUFL	$27, X5, X0
	MOVAPS	X1, X3
	PCMPGTL	X0, X3		// 3.0.1.0 |>| 2.1.0.0 -> 1.0.1.0
	PCMPGTL	X1, X0		// 2.1.0.0 |>| 3.0.1.0 -> 0.1.0.0
	MOVMSKPS X3, BX		// 1010
	MOVMSKPS X0, AX		// 0100
	// !(AX == 0 && BX == 0)
	MOVQ	BX, R8
	ORQ	AX, R8		// R8 will be 0 if both masks are 0
	JNE	less_determine	// yes, one is not 0
	// no, both masks are 0
	ADDQ	$16, DX
	CMPQ	DX, $64
	JE	less_eol
	JMP	less_loop

less_determine:
	XORQ	R8, R8
	CMPQ	BX, AX
	SETGT	R8
less_eol:
	MOVB	R8, ret+16(FP)
	RET
