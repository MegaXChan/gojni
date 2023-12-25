package java

import "github.com/petermattis/goid"

var selfClass = map[int64]uintptr{}

func GetSelfClassOrObject() uintptr {
	gid := goid.Get()
	id, ok := selfClass[gid]
	if ok {
		return id
	}
	return 0
}

func setSelfClassOrObject(obj uintptr) {
	gid := goid.Get()
	selfClass[gid] = obj
}

func clearSelfClassOrObject() {
	gid := goid.Get()
	delete(selfClass, gid)
}
