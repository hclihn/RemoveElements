package main

import (
	"fmt"
	"reflect"
	"strings"
)

func RemoveElements(a interface{}, rm_idx []int) error {
  if len(rm_idx) == 0 {
		return nil
	}
	ta := reflect.TypeOf(a)
	if ta.Kind() != reflect.Ptr {
		return fmt.Errorf("type mismatch: expected a pointer to slice, got %s", ta.Kind())
	}
    	va := reflect.ValueOf(a).Elem() // slice value
	if va.Kind() != reflect.Slice {
		return fmt.Errorf("type mismatch: expected a pointer to slice, got a pointer to %s", va.Kind())
	} 
    	zero := reflect.Zero(ta.Elem().Elem()) // zero value of the slice element type
	la, lr, ir, dest, ln := va.Len(), len(rm_idx), 0, 0, 0
	nextRmIdx := func() (int, error) {
		i := la
		if ir < lr {
			i = rm_idx[ir]
      if i < 0 || i >= la {
				return -1, fmt.Errorf("index (%d) specified in rm_idx[%d] out of range 0..%d", i, ir, la-1)
			}
			ir++
		}
		return i, nil
	}
	idx, err := nextRmIdx()
  if err != nil {
    return fmt.Errorf("bad index specified: %w", err)
  }
	for src := 0; src < la; src++ {
		if src < idx {
			if src != dest {
				va.Index(dest).Set(va.Index(src))
			}
			dest++
			ln++
			continue
		}
		if idx, err = nextRmIdx(); err != nil {
      return fmt.Errorf("bad index specified: %w", err)
    }
	}
	for ; dest < la; dest++ {
		va.Index(dest).Set(zero) // clear the removed elements
	}
	va.SetLen(ln) // adjust the length, required a to be assignable (a pointer to slice)
	return nil
}

type t struct {
	name string
	age int
}

func main() {
	a := make([]t, 10)
	for i := range a {
		a[i] = t{name: strings.Repeat("*", (i % 4) +1), age: i % 4}
	}
	rm_idx := make([]int, 0)
	for i, v := range a {
		if v.age >= 3  {
			rm_idx = append(rm_idx, i)
		}
	}
	fmt.Println(a, rm_idx)
	fmt.Println(RemoveElements(&a, rm_idx))
	fmt.Println(a, rm_idx)
}
