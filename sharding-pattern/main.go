package main

import (
	"log"

	"github.com/pvoliveira/go-study/sharding-pattern/sharding"
)

type Obj struct {
	X string
	Y int
}

func main() {
	s := sharding.NewShardedMap[Obj](5)

	s.Set("alpha", Obj{"val1", 1})
	s.Set("beta", Obj{"val2", 5})
	s.Set("gamma", Obj{"hi", 8})
	s.Set("a", Obj{"hello", 9})
	s.Set("b", Obj{"go", 123})
	s.Set("c", Obj{"lang", 333})

	log.Println(s.Get("alpha"))
	log.Println(s.Get("b"))
	log.Println(s.Get("gamma"))

	keys := s.Keys()
	for _, k := range keys {
		log.Println(k)
	}
}
