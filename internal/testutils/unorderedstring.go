package testutils

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// UnorderedStringParts allows easy comparison of strings where the order of
// their delimited parts needs to be ignored
type UnorderedStringParts struct {
	seperator string
	parts     []string
}

// NewUnorderedStringParts creates a new *UnorderedStringParts by breaking it
// into parts and sorting them
func NewUnorderedStringParts(value, seperator string) *UnorderedStringParts {
	result := UnorderedStringParts{
		seperator: seperator,
		parts:     strings.Split(value, seperator),
	}
	sort.Strings(result.parts)
	return &result
}

// String converts the parts back into a comparable string.
func (unordered *UnorderedStringParts) String() string {
	return strings.Join(unordered.parts, unordered.seperator)
}

// Equals reports if two *UnorderedStringParts are equal. They do not
// to have the same seperator to be equal, as long as the reconstituted
// string is equal:
// ie. NewUnorderedStringParts("c;d,f,e,b,a", ",").Equals(
//             NewUnorderedStringParts("d,e,f;a,b,c", ";")
//     ) == true
func (unordered *UnorderedStringParts) Equals(b *UnorderedStringParts) bool {
	return unordered.String() == b.String()
}

// EqualsString same as Equals, but breaks and sorts its argument (assuming the
// same seperator) before comparison.
func (unordered *UnorderedStringParts) EqualsString(b string) bool {
	innerB := NewUnorderedStringParts(b, unordered.seperator)
	return unordered.Equals(innerB)
}

// GetUnorderedStringMaker is a convienece method allowing to create multiple
// *UnorderedStringParts with the same seperator.
func GetUnorderedStringMaker(seperator string) func(string) *UnorderedStringParts {
	return func(value string) *UnorderedStringParts {
		return NewUnorderedStringParts(value, seperator)
	}
}

// AssertEqual asserts that two *UnorderedStringParts are equal. They
// do not to have the same seperator to be equal, as long as the
// reconstituted string is equal:
// ie. NewUnorderedStringParts("c;d,f,e,b,a", ",").AssertEqual(
//             t,
//             NewUnorderedStringParts("d,e,f;a,b,c", ";")
//     ) // passes
func (unordered *UnorderedStringParts) AssertEqual(t *testing.T, actual *UnorderedStringParts, msgAndArgs ...interface{}) bool {
	return assert.Equal(t, unordered.String(), actual.String(), msgAndArgs...)
}

// AssertEqualString same as AssertEqual, but breaks and sorts actual
// (assuming the same seperator) before comparison.
func (unordered *UnorderedStringParts) AssertEqualString(t *testing.T, actual string, msgAndArgs ...interface{}) bool {
	innerActual := NewUnorderedStringParts(actual, unordered.seperator)
	return unordered.AssertEqual(t, innerActual)
}

// AssertEqualf asserts that two *UnorderedStringParts are equal. They
// do not to have the same seperator to be equal, as long as the
// reconstituted string is equal:
// ie. NewUnorderedStringParts("c;d,f,e,b,a", ",").AssertEqual(
//             t,
//             NewUnorderedStringParts("d,e,f;a,b,c", ";")
//     ) // passes
func (unordered *UnorderedStringParts) AssertEqualf(t *testing.T, actual *UnorderedStringParts, msg string, args ...interface{}) bool {
	return assert.Equalf(t, unordered.String(), actual.String(), msg, args...)
}

// AssertEqualStringf same as AssertEqualf, but breaks and sorts actual
// (assuming the same seperator) before comparison.
func (unordered *UnorderedStringParts) AssertEqualStringf(t *testing.T, actual string, msg string, args ...interface{}) bool {
	innerActual := NewUnorderedStringParts(actual, unordered.seperator)
	return unordered.AssertEqualf(t, innerActual, msg, args...)
}

// AssertNotEqual asserts that two *UnorderedStringParts are not equal.
// They do not to have the same seperator to be equal, as long as the
// reconstituted string is equal:
// ie. NewUnorderedStringParts("c;d,f,e,b,a", ",").AssertNotEqual(
//             t,
//             NewUnorderedStringParts("d,e,f;a,b,c", ";")
//     ) // fails
func (unordered *UnorderedStringParts) AssertNotEqual(t *testing.T, actual *UnorderedStringParts, msgAndArgs ...interface{}) bool {
	return assert.NotEqual(t, unordered.String(), actual.String(), msgAndArgs...)
}

// AssertNotEqualString same as AssertNotEqual, but breaks and sorts actual
// (assuming the same seperator) before comparison.
func (unordered *UnorderedStringParts) AssertNotEqualString(t *testing.T, actual string, msgAndArgs ...interface{}) bool {
	innerActual := NewUnorderedStringParts(actual, unordered.seperator)
	return unordered.AssertNotEqual(t, innerActual)
}

// AssertNotEqualf asserts that two *UnorderedStringParts are not equal.
// They do not to have the same seperator to be equal, as long as the
// reconstituted string is equal:
// ie. NewUnorderedStringParts("c;d,f,e,b,a", ",").AssertNotEqual(
//             t,
//             NewUnorderedStringParts("d,e,f;a,b,c", ";")
//     ) // fails
func (unordered *UnorderedStringParts) AssertNotEqualf(t *testing.T, actual *UnorderedStringParts, msg string, args ...interface{}) bool {
	return assert.NotEqualf(t, unordered.String(), actual.String(), msg, args...)
}

// AssertNotEqualStringf same as AssertNotEqualf, but breaks and sorts actual
// (assuming the same seperator) before comparison.
func (unordered *UnorderedStringParts) AssertNotEqualStringf(t *testing.T, actual string, msg string, args ...interface{}) bool {
	innerActual := NewUnorderedStringParts(actual, unordered.seperator)
	return unordered.AssertNotEqualf(t, innerActual, msg, args...)
}
