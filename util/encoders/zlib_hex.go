package encoders

/*
	Leonidas C2 Framework
	Copyright (C) 2026  Leonidas C2 Project

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// ZlibHex - A custom Leonidas C2 encoder that:
//   1. Compresses data with zlib (deflate level 6).
//   2. Hex-encodes the compressed bytes.
//
// The hex output looks like CSS/JS minified data at a glance, helping
// C2 traffic blend in with application responses that embed compressed
// resources as hex strings (e.g. some font/image loaders).

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"io"
)

// ZlibHex encoder — composite zlib + hex
type ZlibHex struct{}

// Encode - Zlib compress then hex encode.
func (e ZlibHex) Encode(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(data); err != nil {
		w.Close()
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	encoded := hex.EncodeToString(buf.Bytes())
	return []byte(encoded), nil
}

// Decode - Hex decode then zlib decompress.
func (e ZlibHex) Decode(data []byte) ([]byte, error) {
	compressed, err := hex.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	r, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
