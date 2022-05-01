/******************************************************************************
 * arc_test.go
 * Author:
 * Usage:    `go test`  or  `go test -v`
  ******************************************************************************/

package arc

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestARC(t *testing.T) {

	fmt.Println("TEST 1: Initalizing the data structure")
	// TEST 1: Initalizing the data structure
	arc := NewARC(8)
	size := arc.SizeARC()
	if size != 8 {
		fmt.Println("Failed to initalize Cache - Size is wrong")
	} else {
		fmt.Println("Initalized Cache - Size is correct")
	}
	fmt.Println(arc.mid)
	fmt.Println()

	fmt.Println("TEST 2: Filling up the full LRU (length 4)")
	// TEST 2: Filling up the full LRU (length 4)
	for i := 0; i < 4; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	fmt.Println(len(arc.order))
	fmt.Println("Order:")
	for i := 0; i < len(arc.order); i++ {
		fmt.Println(arc.order[i])
	}
	fmt.Println()

	fmt.Println("Test 3: Adding entry to a full LRU")
	// Test 3: Adding entry to a full LRU
	key := fmt.Sprintf("__%s", "4")
	val := fmt.Sprintf("__%s", "4")
	arc.SetARC(key, val)

	fmt.Println("Order:")
	for i := 0; i < len(arc.order); i++ {
		fmt.Println(arc.order[i])
	}
	// seeing if fallen entry is in B1
	fmt.Println("B1:")
	for key, value := range arc.B1 {
		fmt.Println(key, "->", value)
	}
	fmt.Println()

	fmt.Println("Test 4: Calling Get on item in index and seeing if it moves to LFU")
	// Test 4: Calling Get on item in index and seeing if it moves to LFU
	fmt.Println("Value:")
	value := arc.GetARC("__4")
	fmt.Println(reflect.ValueOf(value))
	fmt.Println("Order:")
	for i := 0; i < len(arc.order); i++ {
		fmt.Println(arc.order[i])
	}
	fmt.Println()

	fmt.Println("Test 5: Filling up LRU, overfilling LFU to drop an item to B2")
	// Test 5: Filling up LRU, overfilling LFU to drop an item to B2
	_ = arc.GetARC("__3")
	_ = arc.GetARC("__2")
	_ = arc.GetARC("__1")
	fmt.Println("Order:")
	for i := 0; i < len(arc.order); i++ {
		fmt.Println(arc.order[i])
	}

	for i := 5; i < 9; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	_ = arc.GetARC("__5")

	fmt.Println("Order:")
	for i := 0; i < len(arc.order); i++ {
		fmt.Println(arc.order[i])
	}

	fmt.Println("B2:")
	for key, value := range arc.B2 {
		fmt.Println(key, "->", value)
	}
	fmt.Println()

	fmt.Println("Test 6: Getting items in LFU to change the order by frequency")
	// Test 6: Getting items in LFU to change the order by frequency
	_ = arc.GetARC("__5")

	fmt.Println("Order:")
	for i := 0; i < len(arc.order); i++ {
		fmt.Println(arc.order[i])
	}
	fmt.Println()

	fmt.Println("Test 7: Increasing the frequency of different items in the LFU")
	// Test 7: Increasing the frequency of different items in the LFU
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__5")
	_ = arc.GetARC("__3")

	fmt.Println("Order:")
	for i := 0; i < len(arc.order); i++ {
		fmt.Println(arc.order[i])
	}
	fmt.Println()

	fmt.Println("Test 8: Calling Get on an item in B1 (__0) to increase size of LRU")
	// Test 8: Calling Get on an item in B1 (__0) to increase size of LRU
	_ = arc.GetARC("__0")

	fmt.Println("Order:")
	for i := 0; i < len(arc.order); i++ {
		fmt.Println(arc.order[i])
	}
	fmt.Println()

	fmt.Println("Test 9: Calling Get on an item in B2 (__1) to increase size of LFU")
	// Test 9: Calling Get on an item in B2 (__1) to increase size of LFU

	// First increase LRU to be full which means adding two more values
	for i := 9; i < 11; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	_ = arc.GetARC("__1")

	fmt.Println("Order:")
	for i := 0; i < len(arc.order); i++ {
		fmt.Println(arc.order[i])
	}
	fmt.Println()

	fmt.Println("Test 10: Calling B2 until cache is all LRU, then calling B2 again")
	// Test 10: Calling B2 until cache is all LRU, then calling B2 again

	for i := 0; i < 4; i++ {
		_ = arc.GetARC("__1")
	}

	fmt.Println("Order:")
	for i := 0; i < len(arc.order); i++ {
		fmt.Println(arc.order[i])
	}
	fmt.Println()

	fmt.Println("Test 11: Calling B1 until cache is all LRU, then calling B2 again")
	// Test 11: Calling B1 until cache is all LRU, then calling B2 again

	for i := 0; i < 3; i++ {
		_ = arc.GetARC("__0")
	}
	for i := 11; i < 14; i++ {
		key := fmt.Sprintf("__%s", strconv.Itoa(i))
		val := fmt.Sprintf("__%s", strconv.Itoa(i))
		arc.SetARC(key, val)
	}
	for i := 0; i < 5; i++ {
		_ = arc.GetARC("__0")
		fmt.Println(arc.mid)
	}

	fmt.Println("Order:")
	for i := 0; i < len(arc.order); i++ {
		fmt.Println(arc.order[i])
	}
	fmt.Println()

	/*
		  value := arc.GetARC(key)
			  if reflect.TypeOf(value).String() != "string" {
				  fmt.Sprintf("Failed to add binding with key: %s", key)
				  /*  if reflect.ValueOf(value) != val {
					  t.Errorf("Failed to add binding with key: %s", key)
					  t.FailNow()
	*/
}
