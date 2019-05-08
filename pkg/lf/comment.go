/*
 * LF: Global Fully Replicated Key/Value Store
 * Copyright (C) 2018-2019  ZeroTier, Inc.  https://www.zerotier.com/
 *
 * Licensed under the terms of the MIT license (see LICENSE.txt).
 */

package lf

import (
	"fmt"
)

// These are protocol and database constants and can't be changed.
const (
	commentAssertionNil                         byte = 0
	commentAssertionRecordCollidesWithClaimedID byte = 1

	commentReasonNone                 byte = 0 // No reason given
	commentReasonAutomaticallyFlagged byte = 1 // Issue detected automatically
	commentReasonManuallyFlagged      byte = 2 // A meat sack said so
)

// comment describes a record datum in a commentary record generated by a node.
type comment struct {
	subject           []byte // subject/target of comment (max 255 bytes)
	assertion, reason byte
}

func (c *comment) string() string {
	var reason string
	switch c.reason {
	case commentReasonNone:
		reason = "no reason given"
	case commentReasonAutomaticallyFlagged:
		reason = "automatically flagged"
	case commentReasonManuallyFlagged:
		reason = "manually flagged"
	default:
		reason = fmt.Sprintf("unknown reason %.2x", c.reason)
	}

	switch c.assertion {
	case commentAssertionNil:
		return "nil"
	case commentAssertionRecordCollidesWithClaimedID:
		return fmt.Sprintf("=%s collides with previously claimed ID (%s)", Base62Encode(c.subject), reason)
	}

	return fmt.Sprintf("unknown assertion %.2x subject %x reason %.2x", c.assertion, c.subject, c.reason)
}

func (c *comment) sizeBytes() int {
	return 3 + len(c.subject)
}

func (c *comment) appendTo(b []byte) ([]byte, error) {
	if len(c.subject) > 255 {
		return nil, ErrInvalidObject
	}
	b = append(b, c.assertion, c.reason, byte(len(c.subject)))
	return append(b, c.subject...), nil
}

func (c *comment) readFrom(b []byte) ([]byte, error) {
	if len(b) < 3 {
		return nil, ErrInvalidObject
	}

	c.assertion = b[0]
	c.reason = b[1]
	subLen := int(b[2])
	if len(b) < subLen+3 {
		return nil, ErrInvalidObject
	}
	if subLen > 0 {
		c.subject = make([]byte, subLen)
		copy(c.subject, b[3:subLen+3])
	}

	return b[subLen+3:], nil
}
