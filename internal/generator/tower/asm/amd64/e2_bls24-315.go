// Copyright 2021 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package amd64

import "github.com/consensys/bavard/amd64"

func (fq2 *Fq2Amd64) generateMulByNonResidueE2BLS24_315() {
	// // MulByNonResidue multiplies a E2 by (0,1)
	// func (z *E2) MulByNonResidue(x *E2) *E2 {
	// 	z.A0, z.A1 = x.A1, x.A0
	// 	fp.MulBy13(&z.A0)
	// 	return z
	// }

	registers := fq2.FnHeader("mulNonResE2", 0, 16)

	// a *= 13
	x := fq2.Pop(&registers)
	t := fq2.PopN(&registers)
	s := fq2.PopN(&registers)
	u := fq2.PopN(&registers)

	fq2.MOVQ("x+8(FP)", x)

	fq2.Mov(x, t, fq2.NbWords)

	fq2.Add(t, t)
	fq2.ReduceElement(t, s)
	fq2.Add(t, t)
	fq2.ReduceElement(t, u)

	fq2.Mov(t, u) // u == 4

	fq2.Add(t, t) // t == 8
	fq2.ReduceElement(t, s)

	fq2.Add(u, t) // t == 12
	fq2.ReduceElement(t, s)

	fq2.Add(x, t, fq2.NbWords) // t == 13
	fq2.ReduceElement(t, s)

	fq2.Push(&registers, u...)

	z := fq2.Pop(&registers)
	fq2.MOVQ("res+0(FP)", z)

	// 	z.A0 = t
	//  z.A1 = x.A0
	fq2.Mov(x, s)
	fq2.Mov(s, z, 0, fq2.NbWords)
	fq2.Mov(t, z)

	fq2.RET()

	fq2.Push(&registers, x)
	fq2.Push(&registers, z)
	fq2.Push(&registers, s...)
	fq2.Push(&registers, t...)
}

func (fq2 *Fq2Amd64) generateMulE2BLS24_315(forceCheck bool) {
	// func (z *E2) Mul(x, y *E2) *E2 {
	// 	var a, b, c fp.Element
	// 	a.Add(&x.A0, &x.A1)
	// 	b.Add(&y.A0, &y.A1)
	// 	a.Mul(&a, &b)
	// 	b.Mul(&x.A0, &y.A0)
	// 	c.Mul(&x.A1, &y.A1)
	// 	z.A1.Sub(&a, &b).Sub(&z.A1, &c)
	// 	fp.MulBy13(&c)
	// 	z.A0.Add(&c, &b)
	// 	return z

	// we need a bit of stack space to store the results of the xA0yA0 and xA1yA1 multiplications
	const argSize = 24
	minStackSize := 0
	if forceCheck {
		minStackSize = argSize
	}
	stackSize := fq2.StackSize(fq2.NbWords*4+2, 2, minStackSize)
	registers := fq2.FnHeader("mulAdxE2", stackSize, argSize, amd64.DX, amd64.AX)
	defer fq2.AssertCleanStack(stackSize, minStackSize)

	fq2.WriteLn("NO_LOCAL_POINTERS")

	fq2.WriteLn(`
	// 	var a, b, c fp.Element
	// 	a.Add(&x.A0, &x.A1)
	// 	b.Add(&y.A0, &y.A1)
	// 	a.Mul(&a, &b)
	// 	b.Mul(&x.A0, &y.A0)
	// 	c.Mul(&x.A1, &y.A1)
	// 	z.A1.Sub(&a, &b).Sub(&z.A1, &c)
	// 	fp.MulBy13(&c)
	// 	z.A0.Add(&c, &b)
	`)

	lblNoAdx := fq2.NewLabel()
	// check ADX instruction support
	if forceCheck {
		fq2.CMPB("·supportAdx(SB)", 1)
		fq2.JNE(lblNoAdx)
	}

	// used in the mul operation
	op1 := registers.PopN(fq2.NbWords)
	res := registers.PopN(fq2.NbWords)

	xat := func(i int) string {
		return string(op1[i])
	}

	ax := amd64.AX
	dx := amd64.DX

	aStack := fq2.PopN(&registers, true)
	cStack := fq2.PopN(&registers, true)

	fq2.MOVQ("x+8(FP)", ax)

	// c = x.A1 * y.A1
	fq2.Mov(ax, op1, fq2.NbWords)
	fq2.MulADX(&registers, xat, func(i int) string {
		fq2.MOVQ("y+16(FP)", dx)
		return dx.At(i + fq2.NbWords)
	}, res)
	fq2.ReduceElement(res, op1)
	// res = x.A1 * y.A1
	// pushing on stack for later use.
	fq2.Mov(res, cStack)

	fq2.MOVQ("x+8(FP)", ax)
	fq2.MOVQ("y+16(FP)", dx)

	// a = x.a0 + x.a1
	fq2.Mov(ax, op1, fq2.NbWords)
	fq2.Add(ax, op1)
	fq2.Mov(op1, aStack)

	// b = y.a0 + y.a1
	fq2.Mov(dx, op1)
	fq2.Add(dx, op1, fq2.NbWords)
	// --> note, we don't reduce, as this is used as input to the mul which accept input of size D-1/2 -1

	// a = 	a * b = (x.a0 + x.a1) *  (y.a0 + y.a1)
	fq2.MulADX(&registers, xat, func(i int) string {
		return string(aStack[i])
	}, res)
	fq2.ReduceElement(res, op1)

	// moving result to the stack.
	fq2.Mov(res, aStack)

	// b = x.A0 * y.AO
	fq2.MOVQ("x+8(FP)", ax)

	fq2.Mov(ax, op1)

	fq2.MulADX(&registers, xat, func(i int) string {
		fq2.MOVQ("y+16(FP)", dx)
		return dx.At(i)
	}, res)
	fq2.ReduceElement(res, op1)

	zero := dx
	fq2.XORQ(zero, zero)

	// a = a - b -c
	fq2.Mov(aStack, op1)
	fq2.Sub(res, op1) // a -= b
	fq2.Mov(res, aStack)
	fq2.modReduceAfterSubScratch(zero, op1, res)

	fq2.Sub(cStack, op1) // a -= c
	fq2.modReduceAfterSubScratch(zero, op1, res)

	// z.A1 = a - b - c
	fq2.MOVQ("res+0(FP)", ax)
	fq2.Mov(op1, ax, 0, fq2.NbWords)

	// c *= 13
	fq2.Mov(cStack, res)
	fq2.Mov(aStack, op1)
	fq2.Add(res, op1)
	fq2.ReduceElement(op1, res)
	fq2.Mov(cStack, res)
	fq2.Mov(op1, cStack) // cStack = c + b

	fq2.Add(res, res)
	fq2.ReduceElement(res, op1)
	fq2.Add(res, res)
	fq2.ReduceElement(res, op1)

	// we could avoid stack here
	fq2.Mov(res, aStack) // u == 4

	fq2.Add(res, res) // t == 8
	fq2.ReduceElement(res, op1)

	fq2.Add(aStack, res) // t == 12
	fq2.ReduceElement(res, op1)

	fq2.Add(cStack, res) // t == 13
	fq2.ReduceElement(res, op1)

	// z.A0 = b + c
	fq2.Mov(res, ax)

	fq2.RET()

	// No adx
	if forceCheck {
		fq2.LABEL(lblNoAdx)
		fq2.MOVQ("res+0(FP)", amd64.AX)
		fq2.MOVQ(amd64.AX, "(SP)")
		fq2.MOVQ("x+8(FP)", amd64.AX)
		fq2.MOVQ(amd64.AX, "8(SP)")
		fq2.MOVQ("y+16(FP)", amd64.AX)
		fq2.MOVQ(amd64.AX, "16(SP)")
		fq2.WriteLn("CALL ·mulGenericE2(SB)")
		fq2.RET()

	}

	fq2.Push(&registers, aStack...)
	fq2.Push(&registers, cStack...)

}
