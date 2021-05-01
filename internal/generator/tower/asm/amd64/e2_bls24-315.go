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
