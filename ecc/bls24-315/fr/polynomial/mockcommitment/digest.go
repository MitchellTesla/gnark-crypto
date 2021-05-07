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

package mockcommitment

import (
	"io"

	"github.com/consensys/gnark-crypto/ecc/bls24-315/fr"
)

// MockDigest contains the polynomial itself
type MockDigest struct {
	d fr.Element
}

// WriteTo mock impementation
func (md *MockDigest) WriteTo(w io.Writer) (n int64, err error) {
	return 0, nil
}

// ReadFrom mock impementation
func (md *MockDigest) ReadFrom(r io.Reader) (n int64, err error) {
	return 0, nil
}

// Bytes serializes the digest
func (md *MockDigest) Bytes() []byte {
	b := md.d.Bytes()
	return b[:]
}
