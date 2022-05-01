/*****************************************************************************
 * arc.go
 * Name: Jenny Sun and Shanzay Waseem
 * NetId: jwsun and swaseem
 *****************************************************************************/
package arc

// package cache

import (
	"errors"
	"log"
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
	if size < 2 {
		log.Fatal("Size of cache is too small")
	}
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
	if c.mid >= 0 && c.mid < c.size-1 {
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
	for i := index + 1; i < c.size; i++ {
		nextKey := c.order[i]
		nextEntry := c.cache[nextKey]
		if e.freq+1 > nextEntry.freq {
			c.order[i-1] = nextKey
			c.order[i] = e.key
		} else {
			break
		}
	}
}

/* IF it is on LRU side, increase frequency and place on LFU side
 */
func (c *ARC) LRUincrement(e *cacheEntry, index int) {
	droppedkey := c.order[c.mid+1]
	backtemp := c.order[c.mid+1:]
	entered := false
	for i := len(backtemp) - 1; i >= 0; i-- {
		if backtemp[i] == "" {
			backtemp[i] = e.key
			entered = true
			break
		}
	}
	if entered == false {
		backtemp[0] = e.key
	}

	var fronttemp []string
	space := make([]string, 1)
	if index == 0 {
		fronttemp = c.order[1 : c.mid+1]
	} else {
		fronttemp = append(c.order[:index], c.order[index+1:c.mid+1]...)
	}

	backtemp = append(space, backtemp...)

	c.order = append(fronttemp, backtemp...)

	if _, ok := c.cache[droppedkey]; ok {
		// Delete droped key from T2
		delete(c.cache, droppedkey)
		// Put the dropped entry into B2
		c.B2[droppedkey] = droppedkey
	}
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
