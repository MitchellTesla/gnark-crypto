/*
Copyright © 2020 ConsenSys

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package twistededwards

import (
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bls24-315/fr"
)

func TestMarshal(t *testing.T) {

	var point, unmarshalPoint PointAffine
	point.Set(&edwards.Base)
	for i := 0; i < 20; i++ {
		b := point.Marshal()
		unmarshalPoint.Unmarshal(b)
		if !point.Equal(&unmarshalPoint) {
			t.Fatal("error unmarshal(marshal(point))")
		}
		point.Add(&point, &edwards.Base)
	}
}

func TestAdd(t *testing.T) {

	var p1, p2 PointAffine

	p1.X.SetString("606504702213383506549967602976986366481287558001216693316618222045568485627")
	p1.Y.SetString("127651159908496619690870959004301413053185147667621912551672989511587272066")

	p2.X.SetString("7662037931514052354271565197730264559843515866448459536468212544166952037402")
	p2.Y.SetString("5374285597877801050393191367017366954304415358605622399798479904245876287656")

	var expectedX, expectedY fr.Element

	expectedX.SetString("4909234445028952914550300988482203785951133624557087480905450649944663760912")
	expectedY.SetString("611136336403416894346286896787687479722638935668225997874514807048441098127")

	p1.Add(&p1, &p2)

	if !p1.X.Equal(&expectedX) {
		t.Fatal("wrong x coordinate")
	}
	if !p1.Y.Equal(&expectedY) {
		t.Fatal("wrong y coordinate")
	}

}

func TestAddProj(t *testing.T) {

	var p1, p2 PointAffine
	var p1proj, p2proj PointProj

	p1.X.SetString("606504702213383506549967602976986366481287558001216693316618222045568485627")
	p1.Y.SetString("127651159908496619690870959004301413053185147667621912551672989511587272066")

	p2.X.SetString("7662037931514052354271565197730264559843515866448459536468212544166952037402")
	p2.Y.SetString("5374285597877801050393191367017366954304415358605622399798479904245876287656")

	p1proj.FromAffine(&p1)
	p2proj.FromAffine(&p2)

	var expectedX, expectedY fr.Element

	expectedX.SetString("4909234445028952914550300988482203785951133624557087480905450649944663760912")
	expectedY.SetString("611136336403416894346286896787687479722638935668225997874514807048441098127")

	p1proj.Add(&p1proj, &p2proj)
	p1.FromProj(&p1proj)

	if !p1.X.Equal(&expectedX) {
		t.Fatal("wrong x coordinate")
	}
	if !p1.Y.Equal(&expectedY) {
		t.Fatal("wrong y coordinate")
	}

}

func TestDouble(t *testing.T) {

	var p PointAffine

	p.X.SetString("7549511283461560995748751740971458842863896160950410295399102497784871995152")
	p.Y.SetString("8431724913862871887763243257724106664766884684772085757470778116377697070220")

	p.Double(&p)

	var expectedX, expectedY fr.Element

	expectedX.SetString("7315562943575520667466083922507663735652170740276006185986921942734506037136")
	expectedY.SetString("1549739767261155844537395769333564653013978414582915869918154407683216747833")

	if !p.X.Equal(&expectedX) {
		t.Fatal("wrong x coordinate")
	}
	if !p.Y.Equal(&expectedY) {
		t.Fatal("wrong y coordinate")
	}
}

func TestDoubleProj(t *testing.T) {

	var p PointAffine
	var pproj PointProj

	p.X.SetString("7549511283461560995748751740971458842863896160950410295399102497784871995152")
	p.Y.SetString("8431724913862871887763243257724106664766884684772085757470778116377697070220")

	pproj.FromAffine(&p).Double(&pproj)

	p.FromProj(&pproj)

	var expectedX, expectedY fr.Element

	expectedX.SetString("7315562943575520667466083922507663735652170740276006185986921942734506037136")
	expectedY.SetString("1549739767261155844537395769333564653013978414582915869918154407683216747833")

	if !p.X.Equal(&expectedX) {
		t.Fatal("wrong x coordinate")
	}
	if !p.Y.Equal(&expectedY) {
		t.Fatal("wrong y coordinate")
	}
}

func TestScalarMul(t *testing.T) {

	// set curve parameters
	ed := GetEdwardsCurve()

	var scalar big.Int
	scalar.SetUint64(23902374)

	var p PointAffine
	p.ScalarMul(&ed.Base, &scalar)

	var expectedX, expectedY fr.Element

	expectedX.SetString("4154591871721798907960491452113634764565764905318523916560046077883612108719")
	expectedY.SetString("5027805281744385394708323878140623534818171762863978182421480579085966004942")

	if !expectedX.Equal(&p.X) {
		t.Fatal("wrong x coordinate")
	}
	if !expectedY.Equal(&p.Y) {
		t.Fatal("wrong y coordinate")
	}

	// test consistancy with negation
	var expected, base PointAffine
	expected.Set(&ed.Base).Neg(&expected)
	scalar.Set(&ed.Order).Lsh(&scalar, 3) // multiply by cofactor=8
	scalar.Sub(&scalar, big.NewInt(1))
	base.Set(&ed.Base)
	base.ScalarMul(&base, &scalar)
	if !base.Equal(&expected) {
		t.Fatal("Mul by order-1 not consistant with neg")
	}
}
