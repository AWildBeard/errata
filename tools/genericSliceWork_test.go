package tools

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

// TestType is some struct that can't implement comparable
type TestType struct {
	Elem *int
}

func (tt *TestType) Equal(ott *TestType) bool {
	return *(tt.Elem) == *(ott.Elem)
}

const (
	oddTestSize  = 1000
	evenTestSize = 1000
	allTestSize  = 5000
)

var (
	oddTest  = make([]*TestType, oddTestSize)
	evenTest = make([]*TestType, evenTestSize)
	allTest  = make([]*TestType, allTestSize)
)

func init() {
	{
		odd := 0
		for index := range oddTest {
			for odd++; odd%3 != 0; odd++ {
			}

			oddCopy := odd
			oddTest[index] = &TestType{Elem: &oddCopy}
		}
	}

	{
		even := 0
		for index := range evenTest {
			for even++; even%2 != 0; even++ {
			}

			evenCopy := even
			evenTest[index] = &TestType{Elem: &evenCopy}
		}
	}

	for index := range allTest {
		i := index
		allTest[index] = &TestType{Elem: &i}
	}
}

func TestSliceContains(t *testing.T) {
	assert.LessOrEqual(t, len(oddTest), len(allTest))
	assert.LessOrEqual(t, len(evenTest), len(allTest))

	oddTestElem := oddTest[rand.Intn(len(oddTest))]
	assert.True(t, SliceContains(allTest, oddTestElem))

	evenTestElem := evenTest[rand.Intn(len(evenTest))]
	assert.True(t, SliceContains(allTest, evenTestElem))

	highOrderTestElem := func() *TestType {
		targetNum := *(evenTest[len(evenTest)-1].Elem)
		targetNum *= 2
		return &TestType{Elem: &targetNum}
	}()

	assert.False(t, SliceContains(evenTest, highOrderTestElem))
}

func TestSliceContainsSet(t *testing.T) {
	assert.LessOrEqual(t, len(oddTest), len(allTest))
	assert.LessOrEqual(t, len(evenTest), len(allTest))

	// Solve the number generation complexity issue not being contained within allTest
	assert.True(t, len(oddTest) >= 100)
	assert.True(t, len(evenTest) >= 100)
	assert.True(t, len(allTest) >= 4500)

	assert.True(t, SliceContainsSet(allTest, oddTest[:100]))
}

func TestSlicesEqual(t *testing.T) {
	assert.True(t, SlicesEqual(evenTest, evenTest))
	assert.True(t, SlicesEqual(oddTest, oddTest))
	assert.True(t, SlicesEqual(allTest, allTest))
	assert.False(t, SlicesEqual(evenTest, oddTest))
	assert.False(t, SlicesEqual(oddTest, evenTest))
	assert.False(t, SlicesEqual(oddTest, allTest))
}

func TestGetDuplicatesFromSlice(t *testing.T) {
	evenTestDuplicates := GetDuplicatesFromSlice(evenTest, evenTest)
	assert.True(t, SlicesEqual(evenTest, evenTestDuplicates))
	assert.False(t, SlicesEqual(oddTest, evenTestDuplicates))

	allTestDuplicates := GetDuplicatesFromSlice(allTest, evenTest)

	assert.True(t, SliceContainsSet(allTest, allTestDuplicates))

	assert.False(t, SliceContainsSet(allTest[:(len(allTest)-1)/2], allTest[(len(allTest)-1)/2:]))
}

func TestGetUniqueFromSlice(t *testing.T) {
	expected := allTest[:(len(allTest)-1)/2]

	assert.True(t, SlicesEqual(expected, GetUniqueFromSlice(allTest[(len(allTest)-1)/2:], expected)))

	assert.Len(t, GetUniqueFromSlice(expected, expected), 0)
}
