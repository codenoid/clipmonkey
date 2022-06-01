package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"golang.design/x/clipboard"
	"golang.org/x/exp/slices"
)

var clipboardHistory = make([]string, 0)
var clipboardLock sync.RWMutex

func init() {
	// load previous stored items
	loadItems()

	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	go watchForChange()
}

func watchForChange() {
	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
	for data := range ch {
		if slices.Contains(clipboardHistory, string(data)) {
			continue
		}

		showNewItem(string(data))

		clipboardLock.Lock()
		clipboardHistory = append([]string{string(data)}, clipboardHistory...)
		if len(clipboardHistory) >= 20 {
			clipboardHistory = clipboardHistory[0:20]
		}
		clipboardLock.Unlock()

		saveItems()
	}
}

func loadItems() {
	clipboardLock.Lock()
	defer clipboardLock.Unlock()

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		return
	}

	cachePath := dirname + string(os.PathSeparator) + ".clipmonkey"
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		log.Println(err)
		return
	}

	data, err := ioutil.ReadFile(cachePath)
	if err != nil {
		log.Println(err)
		return
	}

	dec := gob.NewDecoder(bytes.NewBuffer(data))
	if err := dec.Decode(&clipboardHistory); err != nil {
		log.Println(err)
		return
	}
}

func saveItems() {
	clipboardLock.RLock()
	defer clipboardLock.RUnlock()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(clipboardHistory); err != nil {
		log.Println(err)
		return
	}
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		return
	}

	cachePath := dirname + string(os.PathSeparator) + ".clipmonkey"
	err = ioutil.WriteFile(cachePath, buf.Bytes(), 0644)
	if err != nil {
		log.Println(err)
		return
	}
}
