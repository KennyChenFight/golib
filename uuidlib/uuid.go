// Package uuidlib is for encapsulating github.com/satori/go.uuid any operations
//
// As a quick start:
// 	// this is easy for unit test mock
//	uv4 := uuidlib.GOUUIDGenerator{}.NewV4()
//	// just for generate
//	uv4 = uuidlib.NewV4()
//	fmt.Println(uv4)
package uuidlib

import uuid "github.com/satori/go.uuid"

type GOUUIDGenerator struct {
}

func (g GOUUIDGenerator) NewV1() uuid.UUID {
	return uuid.NewV1()
}

func (g GOUUIDGenerator) NewV2(domain byte) uuid.UUID {
	return uuid.NewV2(domain)
}

func (g GOUUIDGenerator) NewV3(ns uuid.UUID, name string) uuid.UUID {
	return uuid.NewV3(ns, name)
}

func (g GOUUIDGenerator) NewV4() uuid.UUID {
	return uuid.NewV4()
}

func (g GOUUIDGenerator) NewV5(ns uuid.UUID, name string) uuid.UUID {
	return uuid.NewV5(ns, name)
}

func NewV4() uuid.UUID {
	return uuid.NewV4()
}
