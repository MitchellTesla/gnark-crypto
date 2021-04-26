// Copyright 2020 ConsenSys Software Inc.
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

// Code generated by consensys/gnark-crypto DO NOT EDIT

package fptower

import (
	"fmt"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

// ------------------------------------------------------------
// tests

func TestE24Serialization(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	genA := GenE24()

	properties.Property("[BLS24-315] SetBytes(Bytes()) should stay constant", prop.ForAll(
		func(a *E24) bool {
			var b E24
			buf := a.Bytes()
			if err := b.SetBytes(buf[:]); err != nil {
				return false
			}
			return a.Equal(&b)
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestE24ReceiverIsOperand(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	genA := GenE24()
	genB := GenE24()

	properties.Property("[BLS24-315] Having the receiver as operand (addition) should output the same result", prop.ForAll(
		func(a, b *E24) bool {
			var c, d E24
			d.Set(a)
			c.Add(a, b)
			a.Add(a, b)
			b.Add(&d, b)
			return a.Equal(b) && a.Equal(&c) && b.Equal(&c)
		},
		genA,
		genB,
	))

	properties.Property("[BLS24-315] Having the receiver as operand (sub) should output the same result", prop.ForAll(
		func(a, b *E24) bool {
			var c, d E24
			d.Set(a)
			c.Sub(a, b)
			a.Sub(a, b)
			b.Sub(&d, b)
			return a.Equal(b) && a.Equal(&c) && b.Equal(&c)
		},
		genA,
		genB,
	))

	properties.Property("[BLS24-315] Having the receiver as operand (mul) should output the same result", prop.ForAll(
		func(a, b *E24) bool {
			var c, d E24
			d.Set(a)
			c.Mul(a, b)
			a.Mul(a, b)
			b.Mul(&d, b)
			return a.Equal(b) && a.Equal(&c) && b.Equal(&c)
		},
		genA,
		genB,
	))

	properties.Property("[BLS24-315] Having the receiver as operand (square) should output the same result", prop.ForAll(
		func(a *E24) bool {
			var b E24
			b.Square(a)
			a.Square(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("[BLS24-315] Having the receiver as operand (neg) should output the same result", prop.ForAll(
		func(a *E24) bool {
			var b E24
			b.Neg(a)
			a.Neg(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("[BLS24-315] Having the receiver as operand (double) should output the same result", prop.ForAll(
		func(a *E24) bool {
			var b E24
			b.Double(a)
			a.Double(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("[BLS24-315] Having the receiver as operand (Inverse) should output the same result", prop.ForAll(
		func(a *E24) bool {
			var b E24
			b.Inverse(a)
			a.Inverse(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestE24Ops(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	genA := GenE24()
	genB := GenE24()

	properties.Property("[BLS24-315] sub & add should leave an element invariant", prop.ForAll(
		func(a, b *E24) bool {
			var c E24
			c.Set(a)
			c.Add(&c, b).Sub(&c, b)
			return c.Equal(a)
		},
		genA,
		genB,
	))

	properties.Property("[BLS24-315] mul & inverse should leave an element invariant", prop.ForAll(
		func(a, b *E24) bool {
			var c, d E24
			d.Inverse(b)
			c.Set(a)
			c.Mul(&c, b).Mul(&c, &d)
			return c.Equal(a)
		},
		genA,
		genB,
	))

	properties.Property("[BLS24-315] inverse twice should leave an element invariant", prop.ForAll(
		func(a *E24) bool {
			var b E24
			b.Inverse(a).Inverse(&b)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("[BLS24-315] neg twice should leave an element invariant", prop.ForAll(
		func(a *E24) bool {
			var b E24
			b.Neg(a).Neg(&b)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("[BLS24-315] square and mul should output the same result", prop.ForAll(
		func(a *E24) bool {
			var b, c E24
			b.Mul(a, a)
			c.Square(a)
			return b.Equal(&c)
		},
		genA,
	))

	properties.Property("[BLS24-315] Double and add twice should output the same result", prop.ForAll(
		func(a *E24) bool {
			var b E24
			b.Add(a, a)
			a.Double(a)
			return a.Equal(&b)
		},
		genA,
	))

	// test conjugate in cubic extension
	properties.Property("[BLS24-315] test conjugate", prop.ForAll(
		func(a *E24) bool {
			var b, c, d E24
			var e, f, g, h E4
			b.Conjugate(a)
			c.Add(a, &b)
			d.Sub(a, &b)

			e.Double(&a.D0.C0)
			f.Double(&a.D1.C1)
			g.Double(&a.D2.C0)

			return c.D0.C1.Equal(&h) && c.D0.C0.Equal(&e) && c.D1.C0.Equal(&h) && c.D1.C1.Equal(&f) && c.D2.C1.Equal(&h) && c.D2.C0.Equal(&g)
		},
		genA,
	))

	properties.Property("[BLS24-315] test MulBy012", prop.ForAll(
		func(a *E24) bool {
			var l, b E24
			var r [3]E4

			r[0].SetRandom()
			r[1].SetRandom()
			r[2].SetRandom()

			l.D0.C0.Set(&r[0])
			l.D0.C1.Set(&r[1])
			l.D1.C0.Set(&r[2])

			fmt.Printf("a = %v\n", a.String())
			fmt.Printf("l = %v\n", l.String())

			b.Mul(a, &l)
			fmt.Printf("a*l = %v\n", b.String())
			a.MulBy012(&r[0], &r[1], &r[2])
			fmt.Printf("mulby012 = %v\n", a.String())

			return a.Equal(&b)
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

}

// ------------------------------------------------------------
// benches

func BenchmarkE24Add(b *testing.B) {
	var a, c E24
	a.SetRandom()
	c.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Add(&a, &c)
	}
}

func BenchmarkE24Sub(b *testing.B) {
	var a, c E24
	a.SetRandom()
	c.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Sub(&a, &c)
	}
}

func BenchmarkE24Mul(b *testing.B) {
	var a, c E24
	a.SetRandom()
	c.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Mul(&a, &c)
	}
}

func BenchmarkE24Square(b *testing.B) {
	var a E24
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Square(&a)
	}
}

func BenchmarkE24Inverse(b *testing.B) {
	var a E24
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Inverse(&a)
	}
}
