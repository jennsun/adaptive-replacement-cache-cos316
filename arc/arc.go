/*****************************************************************************
 * arc.go
 * Name: Jenny Sun and Shanzay Waseem
 * NetId: jwsun and swaseem
 *****************************************************************************/
 package arc

 // package cache
 
 import (
	 "errors"
	 "fmt"
	 "reflect"
	 "strconv"
	 "sync"
 )
  
  type ARC struct {
	  // <= mid --> LRU
	  // > mid --> LFU
	  mid   int
	  order []string
	  cache map[string]*cacheEntry
	  B1    map[string]string // Evicted from LRU part of cache
	  B2    map[string]string // Evicted from LFU part of cache
	  size  int
	  lock  *sync.Mutex
  }
  
  type cacheEntry struct {
	  key   string
	  value interface{}
	  freq  int
  }
  
  func NewARC(size int) *ARC {
	  c := new(ARC)
	  c.order = make([]string, size)
	  c.cache = make(map[string]*cacheEntry, size)
	  c.B1 = make(map[string]string)
	  c.B2 = make(map[string]string)
	  c.size = size
	  c.mid = int(size/2) - 1
	  c.lock = new(sync.Mutex)
	  return c
  }
  
  func (c *ARC) SizeARC() int {
	  c.lock.Lock()
	  defer c.lock.Unlock()
	  return c.size
  }
  
  func (c *ARC) GetARC(key string) interface{} {
	  c.lock.Lock()
	  defer c.lock.Unlock()
	  if e, ok := c.cache[key]; ok {
		  c.increment(e)
		  return e.value
	  }
	  /* Need to look through the ghost entries to see if it is
	  there, because if it is, it should change mid but still return nil */
	  if _, ok := c.B1[key]; ok {
		  // Increase size of T1 and drop entry from T2
		  // Put the dropped entry into B2
		  c.B1Hit()
		  return nil
	  }
	  if _, ok := c.B2[key]; ok {
		  // Increase size of T2 and drop entry from T1
		  // Put the dropped entry into B1
		  c.B2Hit()
		  return nil
	  }
	  return nil
  }
  
  // Private Function to increase size of T1 and drop entry from T2
  func (c *ARC) B1Hit() error {
	  if c.mid < c.size-1 && c.mid > -1 {
		  key := c.order[c.mid+1]
		  c.order[c.mid+1] = ""
		  delete(c.cache, key)
		  // Put the dropped entry into B2
		  c.B2[key] = key
		  c.mid++
		  return nil
	  }
	  return errors.New("Cannot Increase Recent Cache Entries")
  }
  
  // Private function to increase size of T2 and drop entry from T1
  func (c *ARC) B2Hit() error {
	  if c.mid > 0 && c.mid < c.size {
		  key := c.order[c.mid]
		  c.order[c.mid] = ""
		  delete(c.cache, key)
		  // Put the dropped entry into B1
		  c.B1[key] = key
		  c.mid--
		  return nil
	  }
	  return errors.New("Cannot Increase Recent Cache Entries")
  }
  
  // ARC PROTOCOL - if the item is already in the cache
  func (c *ARC) increment(e *cacheEntry) {
	  for index := 0; index < c.size; index++ {
		  if c.order[index] == e.key {
			  if index < c.mid+1 {
				  c.LRUincrement(e, index)
			  } else if index < c.size-1 {
				  c.LFUincrement(e, index)
			  }
			  e.freq++
			  return
		  }
	  }
  }
  
  /* If it is on LFU side, just increase the frequency of access and placement
  if it is not the most frequently accessed
  */
  func (c *ARC) LFUincrement(e *cacheEntry, index int) {
	  nextKey := c.order[index+1]
	  nextEntry := c.cache[nextKey]
	  if e.freq+1 > nextEntry.freq {
		  c.order[index+1] = e.key
		  c.order[index] = nextKey
	  }
  }
  
  /* IF it is on LRU side, increase frequency and place on LFU side
   */
  func (c *ARC) LRUincrement(e *cacheEntry, index int) {
	  var temp []string
	  if index != 0 {
		  temp = append(c.order[:index-1], c.order[index+1:c.mid]...)
	  } else {
		  temp = c.order[index+1 : c.mid]
	  }
	  temp = append(temp, "")
	  fmt.Println("dies here")
	  temp = append(temp, c.order[index])
	  fmt.Println(len(temp))
	  fmt.Println(len(c.order))
	  /*if c.order[c.mid+2] != [] {
		  temp = append(temp, c.order[c.mid+2:]...)
	  }*/
	  fmt.Println("dies here")
	  fmt.Println(len(temp))
	  // droppedkey := c.order[c.mid+1]
	  //c.order = temp
	  // Delete droped key from T2
	  //delete(c.cache, droppedkey)
	  // Put the dropped entry into B2
	  //c.B2[droppedkey] = droppedkey
  }
  
  func (c *ARC) SetARC(key string, value interface{}) {
	  c.lock.Lock()
	  defer c.lock.Unlock()
	  if e, ok := c.cache[key]; ok {
		  // value already exists for key. overwrite
		  e.value = value
		  c.increment(e)
	  } else {
		  // value doesn't exist. insert
		  e := new(cacheEntry)
		  e.key = key
		  e.value = value
		  e.freq = 1
		  c.cache[key] = e
		  // function to add a new item to LRU
		  c.addLRU(e)
	  }
  }
  
  func (c *ARC) addLRU(e *cacheEntry) {
	  temp := make([]string, 1)
	  temp[0] = e.key
	  temp = append(temp, c.order[0:c.mid]...)
	  temp = append(temp, c.order[c.mid+1:]...)
	  droppedkey := c.order[c.mid]
	  c.order = temp
  
	  if _, ok := c.cache[droppedkey]; ok {
		  // Delete droped key from T1
		  delete(c.cache, droppedkey)
		  // Put the dropped entry into B1
		  c.B1[droppedkey] = droppedkey
	  }
  }
  
  // TESTING
  
  func main() {
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
  