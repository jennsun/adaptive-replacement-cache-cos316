/******************************************************************************
 * arc_test.go
 * Author:
 * Usage: `go test`  or  `go test -v`
  ******************************************************************************/

package arc

import (
	"fmt"
	"strconv"
	"testing"
)

// TEST 1: Initalizing the data structure and test size retrieval - cache even size
func TestNewARC(t *testing.T) {
	arc, _ := NewARC(8)
	size := arc.SizeARC()
	if size != 8 {
		t.Errorf("Size of ARC cache is %d, should be 8", size)
	}
}


// TEST 1: Initalizing the data structure and test size retrieval - empty cache
func TestNewEmptyARC(t *testing.T) {
	_, err := NewARC(0)
	if err == nil {
		t.Errorf("Error should be thrown because cache is empty")
	}
}


// TEST 1: Initalizing the data structure and test size retrieval - cache odd size
func TestNewOddARC(t *testing.T) {
	_, err := NewARC(1)
	if err == nil {
		t.Errorf("Error should be thrown because cache size is <2, too small")
	}
}

// TEST 2: Filling up the full LRU (LRU length = 4)
func TestFillLRU(t *testing.T) {
	arc, err := NewARC(8)
	if err != nil {
		t.Error("Error was thrown")
	}
	for i := 0; i < 4; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	orderSize := len(arc.order)
	if orderSize != 8 {
		t.Errorf("Size of ARC cache is %d, should be 8", orderSize)
	}

	LRUsize := arc.mid + 1
	if LRUsize != 4 {
		t.Errorf("Size of LRU cache is %d, should be 4", LRUsize)
	}
	// test LRU items
	for i := 0; i < LRUsize; i++ {
		expectedVal := fmt.Sprintf("__%s", strconv.Itoa(LRUsize-i-1))
		if arc.order[i] != expectedVal {
			t.Errorf("Item in LRU cache is %s, should be %s", arc.order[i], expectedVal)
		}
	}

	// test LFU items (should all be empty)
	for i := LRUsize; i < orderSize; i++ {
		if arc.order[i] != "" {
			t.Error("Item in LFU cache should be nil")
		}
	}
}

// Test 3: Adding entry to a full LRU
func TestAddEntryToFullLRU(t *testing.T) {
	arc, err := NewARC(8)
	if err != nil {
		t.Error("Error was thrown")
	}
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}

	// current items in LRU don't include 0
	LRUsize := arc.mid + 1
	if LRUsize != 4 {
		t.Errorf("Size of LRU cache is %d, should be 4", LRUsize)
	}
	// test LRU items
	for i := 0; i < LRUsize; i++ {
		expectedVal := fmt.Sprintf("__%s", strconv.Itoa(LRUsize-i))
		if arc.order[i] != expectedVal {
			t.Errorf("Item in LRU cache is %s, should be %s", arc.order[i], expectedVal)
		}
	}

	// see if fallen entry is in B1
	if arc.B1 == nil {
		t.Error("B1 should not be nil")
	}
	expectedVal := fmt.Sprintf("__%s", strconv.Itoa(0))
	if arc.B1[expectedVal] != expectedVal {
		t.Errorf("Fallen entry is %s, should be %s", arc.B1[expectedVal], expectedVal)
	}

}

// Test 4: Calling Get on item in index and seeing if it moves to LFU
func TestGetItemMovestoLFU(t *testing.T) {
	arc, err := NewARC(8)
	if err != nil {
		t.Error("Error was thrown")
	}
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	_ = arc.GetARC("__4")
	orderSize := len(arc.order)
	expectedVal := fmt.Sprintf("__%s", strconv.Itoa(4))
	if arc.order[orderSize-1] != expectedVal {
		t.Errorf(" Item in LFU cache is %s, should be %s", arc.order[orderSize-1], expectedVal)
	}
}

// Test 5: Filling up LRU, overfilling LFU to drop an item to B2
func TestOverfillLFU(t *testing.T) {
	arc, err := NewARC(8)
	if err != nil {
		t.Error("Error was thrown")
	}
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	_ = arc.GetARC("__4")
	_ = arc.GetARC("__3")
	_ = arc.GetARC("__2")
	_ = arc.GetARC("__1")

	// check LFU items
	testCache1 := [8]string{"", "", "", "", "__1", "__2", "__3", "__4"}
	for i := 0; i <= arc.mid; i++ {
		if arc.order[i] != testCache1[i] {
			t.Errorf(" Item in LFU cache is %s, should be %s", arc.order[arc.mid+i], testCache1[i])
		}
	}

	// fill in cache
	for i := 5; i < 9; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	_ = arc.GetARC("__5")

	// check order of current cache items
	testCache := [8]string{"__8", "__7", "__6", "", "__5", "__2", "__3", "__4"}

	orderSize := len(arc.order)
	for i := 0; i < orderSize; i++ {
		if arc.order[i] != testCache[i] {
			t.Errorf("Item in ARC cache is %s, should be %s", arc.order[i], testCache[i])
		}
	}

	// see if fallen entry is in B2
	if arc.B2 == nil {
		t.Error("B2 should not be nil")
	}
	expectedVal := fmt.Sprintf("__%s", strconv.Itoa(1))
	if arc.B2[expectedVal] != expectedVal {
		t.Errorf("Fallen entry is %s, should be %s", arc.B2[expectedVal], expectedVal)
	}
}

// Test 6: Getting items in LFU to change the order by frequency
func TestLFUOrderChange(t *testing.T) {
	arc, err := NewARC(8)
	if err != nil {
		t.Error("Error was thrown")
	}
	// fill in cache
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	_ = arc.GetARC("__4")
	_ = arc.GetARC("__3")
	_ = arc.GetARC("__2")
	_ = arc.GetARC("__1")

	// fill in cache
	for i := 5; i < 9; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	// change order of LFU by getting same value twice
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")

	// check order of current cache items
	testCache := [8]string{"__8", "__7", "__6", "", "__2", "__3", "__4", "__5"}
	orderSize := len(arc.order)
	for i := 0; i < orderSize; i++ {
		if arc.order[i] != testCache[i] {
			t.Errorf("Item in ARC cache is %s, should be %s", arc.order[i], testCache[i])
		}
	}
}

// Test 7: Increasing the frequency of different items in the LFU
func TestIncItemFreqinLFU(t *testing.T) {
	arc, err := NewARC(8)
	if err != nil {
		t.Error("Error was thrown")
	}
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	_ = arc.GetARC("__4")
	_ = arc.GetARC("__3")
	_ = arc.GetARC("__2")
	_ = arc.GetARC("__1")

	// fill in cache
	for i := 5; i < 9; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	// change order of LFU by getting same value twice
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__3")

	testCache := [8]string{"__8", "__7", "__6", "", "__2", "__4", "__3", "__5"}
	orderSize := len(arc.order)
	for i := 0; i < orderSize; i++ {
		if arc.order[i] != testCache[i] {
			t.Errorf("Item in ARC cache is %s, should be %s", arc.order[i], testCache[i])
		}
	}
}

// Test 8: Calling Get on an item in B1 (__0) to increase size of LRU
func TestGetItemInB1(t *testing.T) {
	arc, err := NewARC(8)
	if err != nil {
		t.Error("Error was thrown")
	}
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	_ = arc.GetARC("__4")
	_ = arc.GetARC("__3")
	_ = arc.GetARC("__2")
	_ = arc.GetARC("__1")

	// fill in cache
	for i := 5; i < 9; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	// change order of LFU by getting same value twice
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__3")
	_ = arc.GetARC("__0")

	testCache := [8]string{"__8", "__7", "__6", "", "", "__4", "__3", "__5"}
	orderSize := len(arc.order)
	for i := 0; i < orderSize; i++ {
		if arc.order[i] != testCache[i] {
			t.Errorf("Item in ARC cache is %s, should be %s", arc.order[i], testCache[i])
		}
	}
}

// Test 9: Calling Get on an item in B2 (__1) to increase size of LFU
func TestGetItemInB2(t *testing.T) {
	arc, err := NewARC(8)
	if err != nil {
		t.Error("Error was thrown")
	}
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	_ = arc.GetARC("__4")
	_ = arc.GetARC("__3")
	_ = arc.GetARC("__2")
	_ = arc.GetARC("__1")

	// fill in cache
	for i := 5; i < 9; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	// change order of LFU by getting same value twice
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__3")
	_ = arc.GetARC("__0")

	// Increase LRU to be full which means adding two more values
	for i := 9; i < 11; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	_ = arc.GetARC("__1")

	testCache2 := [8]string{"__10", "__9", "__8", "__7", "", "__4", "__3", "__5"}
	orderSize := len(arc.order)
	for i := 0; i < orderSize; i++ {
		if arc.order[i] != testCache2[i] {
			t.Errorf("Item in ARC cache is %s, should be %s", arc.order[i], testCache2[i])
		}
	}
}

// Test 10: Calling B2 until cache is all LRU, then calling B2 again
func TestCallB2ThenB2(t *testing.T) {
	arc, err := NewARC(8)
	if err != nil {
		t.Error("Error was thrown")
	}
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	_ = arc.GetARC("__4")
	_ = arc.GetARC("__3")
	_ = arc.GetARC("__2")
	_ = arc.GetARC("__1")

	// fill in cache
	for i := 5; i < 9; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	// change order of LFU by getting same value twice
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__3")
	_ = arc.GetARC("__0")

	// Increase LRU to be full which means adding two more values
	for i := 9; i < 11; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	_ = arc.GetARC("__1")

	for i := 0; i < 4; i++ {
		_ = arc.GetARC("__1")
	}

	testCache := [8]string{"__10", "", "", "", "", "__4", "__3", "__5"}
	orderSize := len(arc.order)
	for i := 0; i < orderSize; i++ {
		if arc.order[i] != testCache[i] {
			t.Errorf("Item in ARC cache is %s, should be %s", arc.order[i], testCache[i])
		}
	}

}

// Test 11: Calling B1 until cache is all LRU, then calling B2 again
func TestCallB1ThenB2(t *testing.T) {
	arc, err := NewARC(8)
	if err != nil {
		t.Error("Error was thrown")
	}
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	_ = arc.GetARC("__4")
	_ = arc.GetARC("__3")
	_ = arc.GetARC("__2")
	_ = arc.GetARC("__1")

	// fill in cache
	for i := 5; i < 9; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	// change order of LFU by getting same value twice
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__3")
	_ = arc.GetARC("__0")

	// Increase LRU to be full which means adding two more values
	for i := 9; i < 11; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	_ = arc.GetARC("__1")

	for i := 0; i < 4; i++ {
		_ = arc.GetARC("__1")
	}

	for i := 0; i < 3; i++ {
		_ = arc.GetARC("__0")
	}
	for i := 11; i < 14; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}

	// check midpoints
	testMid := [5]int{4, 5, 6, 7, 7}
	for i := 0; i < 5; i++ {
		_ = arc.GetARC("__0")
		if arc.mid != testMid[i] {
			t.Errorf("Midpoint is %d, should be %d", arc.mid, testMid[i])
		}
	}

	testCache := [8]string{"__13", "__12", "__11", "__10", "", "", "", ""}
	orderSize := len(arc.order)
	for i := 0; i < orderSize; i++ {
		if arc.order[i] != testCache[i] {
			t.Errorf("Item in ARC cache is %s, should be %s", arc.order[i], testCache[i])
		}
	}

}
