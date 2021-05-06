package twistededwards

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bls24-315/fr"
)

// CurveParams curve parameters: ax^2 + y^2 = 1 + d*x^2*y^2
type CurveParams struct {
	A, D     fr.Element // in Montgomery form
	Cofactor fr.Element // not in Montgomery form
	Order    big.Int
	Base     PointAffine
}

var edwards CurveParams

// GetEdwardsCurve returns the twisted Edwards curve on BLS24-315's Fr
func GetEdwardsCurve() CurveParams {

	// copy to keep Order private
	var res CurveParams

	res.A.Set(&edwards.A)
	res.D.Set(&edwards.D)
	res.Cofactor.Set(&edwards.Cofactor)
	res.Order.Set(&edwards.Order)
	res.Base.Set(&edwards.Base)

	return res
}

func init() {

	edwards.A.SetUint64(257732)
	edwards.D.SetUint64(257728)
	edwards.Cofactor.SetUint64(8).FromMont()
	edwards.Order.SetString("1437753473921907580703509300571927811987591765799164617677716990775193563777", 10)

	edwards.Base.X.SetString("2939212147167698761989877129329620056992201626344573600805933973699038335332")
	edwards.Base.Y.SetString("1210739767513185331118744674165833946943116652645479549122735386298364723201")
}
