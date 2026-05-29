package main

import (
	"strings"
	"testing"
)

func TestIsPrintable(t *testing.T) {
	tests := []struct {
		input byte
		want  bool
	}{
		{0x00, false},
		{0x1F, false}, // 31, just below printable range
		{0x20, true},  // space
		{0x41, true},  // 'A'
		{0x7E, true},  // '~'
		{0x7F, false}, // DEL
		{0xFF, false},
	}
	for _, tt := range tests {
		if got := isPrintable(tt.input); got != tt.want {
			t.Errorf("isPrintable(0x%02X) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestPrintableChar(t *testing.T) {
	tests := []struct {
		input byte
		want  byte
	}{
		{0x41, 0x41}, // 'A' stays 'A'
		{0x20, 0x20}, // space stays space
		{0x7E, 0x7E}, // '~' stays '~'
		{0x00, '.'},  // null → dot
		{0x1F, '.'},  // 31 → dot
		{0x7F, '.'},  // DEL → dot
		{0xFF, '.'},  // high byte → dot
	}
	for _, tt := range tests {
		if got := printableChar(tt.input); got != tt.want {
			t.Errorf("printableChar(0x%02X) = 0x%02X, want 0x%02X", tt.input, got, tt.want)
		}
	}
}

func TestDumpLine(t *testing.T) {
	tests := []struct {
		name    string
		offset  int
		numread int
		data    []byte
		want    string
	}{
		{
			name:    "full line of printable bytes",
			offset:  0,
			numread: 16,
			data:    []byte("AAAAAAAAAAAAAAAA"),
			want:    "00000000   41 41 41 41 41 41 41 41  41 41 41 41 41 41 41 41  |AAAAAAAAAAAAAAAA|\n",
		},
		{
			name:    "full line of null bytes",
			offset:  0,
			numread: 16,
			data:    make([]byte, 16),
			want:    "00000000   00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|\n",
		},
		{
			name:    "partial line",
			offset:  16,
			numread: 4,
			data:    append([]byte{0x41, 0x42, 0x43, 0x44}, make([]byte, 12)...),
			want:    "00000010   41 42 43 44                                       |ABCD            |\n",
		},
		{
			name:    "nonzero offset",
			offset:  0x20,
			numread: 16,
			data:    []byte("0123456789ABCDEF"),
			want:    "00000020   30 31 32 33 34 35 36 37  38 39 41 42 43 44 45 46  |0123456789ABCDEF|\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			dumpLine(&sb, tt.offset, tt.numread, tt.data, false)
			if got := sb.String(); got != tt.want {
				t.Errorf("\ngot  %q\nwant %q", got, tt.want)
			}
		})
	}
}
