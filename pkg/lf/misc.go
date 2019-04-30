/*
 * LF: Global Fully Replicated Key/Value Store
 * Copyright (C) 2018-2019  ZeroTier, Inc.  https://www.zerotier.com/
 *
 * Licensed under the terms of the MIT license (see LICENSE.txt).
 */

package lf

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"time"

	"github.com/tidwall/pretty"
)

// TimeMs returns the time in milliseconds since epoch.
func TimeMs() uint64 { return uint64(time.Now().UnixNano()) / uint64(1000000) }

// TimeSec returns the time in seconds since epoch.
func TimeSec() uint64 { return uint64(time.Now().UnixNano()) / uint64(1000000000) }

// TimeMsToTime converts a time in milliseconds since epoch to a Go native time.Time structure.
func TimeMsToTime(ms uint64) time.Time { return time.Unix(int64(ms/1000), int64((ms%1000)*1000000)) }

// TimeSecToTime converts a time in seconds since epoch to a Go native time.Time structure.
func TimeSecToTime(s uint64) time.Time { return time.Unix(int64(s), 0) }

// byteAndArrayReader combines Reader and ByteReader capabilities
type byteAndArrayReader [1]io.Reader

func (mr byteAndArrayReader) Read(p []byte) (int, error) { return mr[0].Read(p) }

func (mr byteAndArrayReader) ReadByte() (byte, error) {
	var tmp [1]byte
	_, err := io.ReadFull(mr[0], tmp[:])
	return tmp[0], err
}

// countingWriter is an io.Writer that increments an integer for each byte "written" to it.
type countingWriter uint

// Write implements io.Writer
func (cr *countingWriter) Write(b []byte) (n int, err error) {
	n = len(b)
	*cr += countingWriter(n)
	return
}

func writeUVarint(out io.Writer, v uint64) (int, error) {
	var tmp [10]byte
	return out.Write(tmp[0:binary.PutUvarint(tmp[:], v)])
}

// integerSqrtRounded computes the rounded integer square root of a 32-bit unsigned int.
// This is used for proof of work calculations since we don't want any inconsisency between nodes regardless of FPU behavior.
func integerSqrtRounded(op uint32) (res uint32) {
	// Translated from C at https://stackoverflow.com/questions/1100090/looking-for-an-efficient-integer-square-root-algorithm-for-arm-thumb2
	one := uint32(1 << 30)
	for one > op {
		one >>= 2
	}
	for one != 0 {
		if op >= (res + one) {
			op = op - (res + one)
			res = res + 2*one
		}
		res >>= 1
		one >>= 2
	}
	if op > res { // rounding
		res++
	}
	return
}

func sliceContainsUInt(s []uint, e uint) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ipIsGlobalPublicUnicast returns true if IP is global unicast and is not a private (10.x.x.x etc.) range.
func ipIsGlobalPublicUnicast(ip net.IP) bool {
	if ip.IsGlobalUnicast() {
		ip4 := ip.To4()
		if len(ip4) == 4 {
			return ip4[0] != 10 && (!((ip4[0] == 192) && (ip4[1] == 168))) && (!((ip4[0] == 172) && (ip4[1] == 16)))
		}
		if len(ip) == 16 {
			return ((ip[0] & 0xfe) != 0xfc)
		}
	}
	return false
}

var jsonPrettyOptions = pretty.Options{
	Width:    2147483647, // always put arrays on one line
	Prefix:   "",
	Indent:   "  ",
	SortKeys: false,
}

// PrettyJSON returns a "pretty" JSON string or the "null" string if something goes wrong.
// This formats things a little differently from MarshalIndent to make the sorts of JSON we generate easier to read.
func PrettyJSON(obj interface{}) string {
	j, err := json.Marshal(obj)
	if err != nil {
		return "null"
	}
	return string(pretty.PrettyOptions(j, &jsonPrettyOptions))
}
