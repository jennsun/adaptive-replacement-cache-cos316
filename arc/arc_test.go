/******************************************************************************
 * arc_test.go
 * Name: Jenny Sun and Shanzay Waseem
 * NetId: jwsun and swaseem
 * Usage:    `go test`  or  `go test -v`
  ******************************************************************************/

  package arc

  import (
	  "fmt"
	  "reflect"
	  "strconv"
	  "testing"
  )
  
  func TestNewARC(t *testing.T) {
  
	  // TEST 1: Initalizing the data structure
	  arc := NewARC(8)
	  size := arc.SizeARC()
	  if size != 8 {
		  fmt.Println("Failed to initalize Cache - Size is wrong")
	  } else {
		  fmt.Println("Initalized Cache - Size is correct")
	  }
	  fmt.Println(arc.mid)
  
	  // TEST 2: Filling up the full LRU (length 4)
	  for i := 0; i < 4; i++ {
		  key := fmt.Sprintf("__%s", strconv.Itoa(i))
		  val := fmt.Sprintf("__%s", strconv.Itoa(i))
		  arc.SetARC(key, val)
	  }
	  fmt.Println(len(arc.order))
	  for i := 0; i < len(arc.order); i++ {
		  fmt.Println(arc.order[i])
	  }
  
	  // Test 3: Adding entry to a full LRU
	  key := fmt.Sprintf("__%s", "4")
	  val := fmt.Sprintf("__%s", "4")
	  arc.SetARC(key, val)
  
	  fmt.Println(len(arc.order))
	  for i := 0; i < len(arc.order); i++ {
		  fmt.Println(arc.order[i])
	  }
	  // seeing if fallen entry is in B1
	  for key, value := range arc.B1 {
		  fmt.Println(key, "->", value)
	  }
  
	  // Test 4: Calling Get on item in index and seeing if it moves to LFU
	  value := arc.GetARC("__4")
	  fmt.Println(reflect.ValueOf(value))
  
	  fmt.Println(len(arc.order))
	  for i := 0; i < len(arc.order); i++ {
		  fmt.Println(arc.order[i])
	  }
  
	  // Test 4: Filling up LRU, overfilling LFU to drop an item to B2
  
	  /*
		  value := arc.GetARC("__3")
		  fmt.Println(reflect.ValueOf(value))
  
		  for i := 0; i < len(arc.order); i++ {
			  fmt.Println(arc.order[i])
		  }*/
  
	  /*	mid   int
		  order []string
		  cache map[string]*cacheEntry
		  B1    map[string]string // Evicted from LRU part of cache
		  B2    map[string]string // Evicted from LFU part of cache
		  size  int
		  lock  *sync.Mutex */
  
	  /*
		  value := arc.GetARC(key)
			  if reflect.TypeOf(value).String() != "string" {
				  fmt.Sprintf("Failed to add binding with key: %s", key)
				  /*  if reflect.ValueOf(value) != val {
					  t.Errorf("Failed to add binding with key: %s", key)
					  t.FailNow()*/
  }