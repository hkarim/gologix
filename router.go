package gologix

import (
	"fmt"
	"log"
	"sync"
)

type PathRouter struct {
	Path map[string]TagProvider
}

func NewRouter() *PathRouter {
	p := new(PathRouter)
	p.Path = make(map[string]TagProvider)
	return p
}

func (router *PathRouter) AddHandler(path []byte, p TagProvider) {
	if router.Path == nil {
		router.Path = make(map[string]TagProvider)
	}
	router.Path[string(path)] = p
}

// find the tag provider for a given cip path
func (router *PathRouter) Resolve(path []byte) (TagProvider, error) {
	tp, ok := router.Path[string(path)]
	if !ok {
		return nil, fmt.Errorf("path %v not recognized", path)
	}
	return tp, nil
}

type TagProvider interface {
	Read(tag string, qty int16) (any, error)
	Write(tag string, value any) error
}

type MapTagProvider struct {
	Mutex sync.Mutex
	Data  map[string]any
}

func (p *MapTagProvider) Read(tag string, qty int16) (any, error) {
	log.Printf("Trying to read %v from MapTagProvider", tag)
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	if p.Data == nil {
		p.Data = make(map[string]any)
	}

	val, ok := p.Data[tag]
	if !ok {
		return nil, fmt.Errorf("tag %v not in map", tag)
	}
	return val, nil
}

func (p *MapTagProvider) Write(tag string, value any) error {
	log.Printf("Trying to set %v=%v from MapTagProvider", tag, value)
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	if p.Data == nil {
		p.Data = make(map[string]any)
	}
	p.Data[tag] = value
	return nil
}