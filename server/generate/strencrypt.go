package generate

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

// Build-time string encryption for Leonidas C2 implants.
//
// This pass scans rendered Go source for raw string literals annotated with
// the magic comment `//leon:encrypt` on the same line and replaces them with
// an XOR-decrypted expression using a per-build random key embedded as a
// byte array. The result is invisible to static-string scanners.
//
// Example (in template source):
//
//	agentName := "LEONIDAS_AGENT" //leon:encrypt
//
// After the pass the compiled implant will not contain the literal string
// "LEONIDAS_AGENT" — it reconstructs it at runtime.
//
// Limitations:
//   - Only processes single-line string literals (no multi-line backtick strings).
//   - Does not handle escape sequences inside strings (keeps them as-is).
//   - Must run before garble so obfuscated symbol names don't break parsing.

import (
	"crypto/rand"
	"fmt"
	"go/scanner"
	"go/token"
	"regexp"
	"strings"
)

var encryptAnnotation = regexp.MustCompile(`"([^"]*)"(\s*//leon:encrypt)`)

// EncryptStrings walks a rendered Go source string and replaces annotated
// string literals with XOR runtime-decode expressions.
// It returns the rewritten source and a map of (literal → key) for debugging.
func EncryptStrings(src string) (string, error) {
	if !strings.Contains(src, "//leon:encrypt") {
		return src, nil
	}

	lines := strings.Split(src, "\n")
	for i, line := range lines {
		if !strings.Contains(line, "//leon:encrypt") {
			continue
		}
		rewritten, err := rewriteLine(line)
		if err != nil {
			buildLog.Warnf("[strencrypt] line %d: %v — skipping", i+1, err)
			continue
		}
		lines[i] = rewritten
	}
	return strings.Join(lines, "\n"), nil
}

// rewriteLine replaces the first annotated string literal on a line with
// an XOR-decode expression.
func rewriteLine(line string) (string, error) {
	m := encryptAnnotation.FindStringSubmatchIndex(line)
	if m == nil {
		return line, nil
	}

	// m[2]:m[3] is the content inside the quotes
	plaintext := line[m[2]:m[3]]

	key, err := randomKey(len(plaintext) + 1)
	if err != nil {
		return "", err
	}

	// Build XOR-encrypted byte array literal
	encrypted := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i++ {
		encrypted[i] = plaintext[i] ^ key[i%len(key)]
	}

	encLiteral := byteSliceLiteral(encrypted)
	keyLiteral := byteSliceLiteral(key)

	// Runtime decode expression: func() string { ... }()
	decodeExpr := fmt.Sprintf(
		`func() string { _k := %s; _e := %s; for _i := range _e { _e[_i] ^= _k[_i%%len(_k)] }; return string(_e) }()`,
		keyLiteral, encLiteral,
	)

	// Replace the matched "literal" //leon:encrypt with the inline decryption
	return line[:m[0]] + decodeExpr, nil
}

// randomKey generates a random XOR key of the given length.
func randomKey(n int) ([]byte, error) {
	if n < 1 {
		n = 16
	}
	key := make([]byte, n)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	// Ensure no zero bytes (XOR with 0 is a no-op and looks suspicious)
	for i, b := range key {
		if b == 0 {
			key[i] = 0xFF
		}
	}
	return key, nil
}

// byteSliceLiteral formats a byte slice as a Go []byte literal.
func byteSliceLiteral(b []byte) string {
	parts := make([]string, len(b))
	for i, v := range b {
		parts[i] = fmt.Sprintf("0x%02x", v)
	}
	return "[]byte{" + strings.Join(parts, ",") + "}"
}

// Ensure go/scanner is used (imported for future AST-based extension).
var _ = scanner.Scanner{}
var _ = token.STRING
