// Copyright (c) 2018 NEC Laboratories Europe GmbH.
//
// Authors: Wenting Li     <wenting.li@neclab.eu>
//          Sergey Fedorov <sergey.fedorov@neclab.eu>
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

package usig

import (
	"bytes"
	"encoding/binary"
)

// USIG (Unique Sequential Identifier Generator) is a tamper-proof
// component in MinBFT that assigns unique, monotonic, and sequential
// counter to messages and signs it
type USIG interface {
	// CreateUI returns a unique identifier for the specified
	// message. A unique, monotonic, and sequential counter is
	// incremented on each invocation to produce the UI
	CreateUI(message []byte) (*UI, error)

	// VerifyUI verifies if the UI is valid for the message and
	// was generated by the specified USIG identity
	VerifyUI(message []byte, ui *UI, usigID []byte) error

	// ID returns the identity of this USIG instance
	ID() []byte
}

// UI is a unique identifier assigned to a message by a USIG
type UI struct {
	// Unique value for each USIG instance
	Epoch uint64

	// Unique, monotonic, and sequential counter
	Counter uint64

	// Certificate created by a tamper-proof component of the USIG
	// that certifies the counter assigned to a particular message
	Cert []byte
}

// MarshalBinary marshals UI to byte array
func (ui *UI) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// First, marshal the epoch and counter
	err := binary.Write(buf, binary.LittleEndian, ui.Epoch)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, ui.Counter)
	if err != nil {
		return nil, err
	}

	// Then, append the USIG certificate bytes
	err = binary.Write(buf, binary.LittleEndian, ui.Cert)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary unmarshals byte array to UI
func (ui *UI) UnmarshalBinary(in []byte) error {
	buf := bytes.NewBuffer(in)

	// First, unmarshal the epoch and counter
	err := binary.Read(buf, binary.LittleEndian, &ui.Epoch)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.LittleEndian, &ui.Counter)
	if err != nil {
		return err
	}

	// The rest are the USIG certificate bytes
	ui.Cert = buf.Bytes()

	return nil
}
