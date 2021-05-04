package main

import (
	bls12377 "github.com/consensys/gnark-crypto/ecc/bls12-377"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	bls24315 "github.com/consensys/gnark-crypto/ecc/bls24-315"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	bw6761 "github.com/consensys/gnark-crypto/ecc/bw6-761"
)

// note: pairing API is not code generated, and don't use interfaces{} for performance reasons
// we end up having some API disparities -- this section ensures that we don't forget to update some APIs

var err error

var (
	gtbls12_377 bls12377.GT
	gtbls12_381 bls12381.GT
	gtbn254     bn254.GT
	gtbw6_761   bw6761.GT
	gtbls24_315 bls24315.GT
)

func init() {
	// Pair
	gtbls12_377, err = bls12377.Pair([]bls12377.G1Affine{}, []bls12377.G2Affine{})
	gtbls12_381, err = bls12381.Pair([]bls12381.G1Affine{}, []bls12381.G2Affine{})
	gtbn254, err = bn254.Pair([]bn254.G1Affine{}, []bn254.G2Affine{})
	gtbw6_761, err = bw6761.Pair([]bw6761.G1Affine{}, []bw6761.G2Affine{})
	gtbls24_315, err = bls24315.Pair([]bls24315.G1Affine{}, []bls24315.G2Affine{})

	// MillerLoop
	gtbls12_377, err = bls12377.MillerLoop([]bls12377.G1Affine{}, []bls12377.G2Affine{})
	gtbls12_381, err = bls12381.MillerLoop([]bls12381.G1Affine{}, []bls12381.G2Affine{})
	gtbn254, err = bn254.MillerLoop([]bn254.G1Affine{}, []bn254.G2Affine{})
	gtbw6_761, err = bw6761.MillerLoop([]bw6761.G1Affine{}, []bw6761.G2Affine{})
	gtbls24_315, err = bls24315.MillerLoop([]bls24315.G1Affine{}, []bls24315.G2Affine{})

	// FinalExp
	gtbls12_377 = bls12377.FinalExponentiation(&gtbls12_377)
	gtbls12_381 = bls12381.FinalExponentiation(&gtbls12_381)
	gtbn254 = bn254.FinalExponentiation(&gtbn254)
	gtbw6_761 = bw6761.FinalExponentiation(&gtbw6_761)
	gtbls24_315 = bls24315.FinalExponentiation(&gtbls24_315)
}
