package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	ylog "github.com/Joe-Degs/yocki/internal/log"
)

var (
	test = "the string to write LOL!"
)

func openfile(typ string) *os.File {
	f, err := os.OpenFile("./test."+typ, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func testStore() {
	fmt.Printf("\n+++++++++++++++++++++store++++++++++++++++++++++++++++\n")
	store, err := ylog.NewStore(openfile("store"))
	defer store.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("len of test -> %d\n", len(test))

	// store some thing\
	var poss []uint64
	for i := 0; i < 5; i++ {
		n, pos, err := store.Append([]byte(test))
		if err != nil {
			log.Fatal(err)
		}
		poss = append(poss, pos)
		fmt.Printf("Append -> n: %d, pos: %d\n", n, pos)
	}

	b, err := store.Read(poss[3])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Read -> %s\n", b)
}

func testIndex() {
	fmt.Printf("\n+++++++++++++++++++++index++++++++++++++++++++++++++++\n")
	config := ylog.Config{}
	config.Segment.MaxIndexBytes = 1024
	idx, err := ylog.NewIndex(openfile("index"), config)
	if err != nil {
		log.Fatal(err)
	}
	defer idx.Close()
	a, b, err := idx.Read(-1)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			log.Fatal(err)
		}
	}
	fmt.Printf("out: %d, pos: %d\n", a, b)
	entries := []struct {
		Off uint32
		Pos uint64
	}{
		{Off: 0, Pos: 0},
		{Off: 1, Pos: 10},
	}
	for _, ent := range entries {
		if err := idx.Write(ent.Off, ent.Pos); err != nil {
			log.Fatal(err)
		}

		out, pos, err := idx.Read(int64(ent.Off))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("out: %d, pos: %d\n", out, pos)
	}
}

func main() {
	testStore()
	testIndex()
}
