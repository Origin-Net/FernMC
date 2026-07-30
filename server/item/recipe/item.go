package recipe

import (
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/world"
	"math"
)



type Item interface {
	
	
	Count() int
	
	Empty() bool
}


type inputItem struct {
	
	Name string `nbt:"name"`
	
	Meta int32 `nbt:"meta"`
	
	Count int32 `nbt:"count"`
	
	State struct {
		Name       string                 `nbt:"name"`
		Properties map[string]interface{} `nbt:"states"`
		Version    int32                  `nbt:"version"`
	} `nbt:"block"`
	
	Tag string `nbt:"tag"`
}


func (i inputItem) Item() (Item, bool) {
	if i.Tag != "" {
		return NewItemTag(i.Tag, int(i.Count)), true
	}

	it, ok := world.ItemByName(i.Name, int16(i.Meta))
	if !ok {
		return nil, false
	}
	st := item.NewStack(it, int(i.Count))
	if i.Meta == math.MaxInt16 {
		st = st.WithValue("variants", true)
	}

	return st, true
}


type inputItems []inputItem


func (d inputItems) Items() ([]Item, bool) {
	s := make([]Item, 0, len(d))
	for _, i := range d {
		itemInput, ok := i.Item()
		if !ok {
			return nil, false
		}
		s = append(s, itemInput)
	}
	return s, true
}


type outputItem struct {
	
	Name string `nbt:"name"`
	
	Meta int32 `nbt:"meta"`
	
	Count int16 `nbt:"count"`
	
	State struct {
		Name       string                 `nbt:"name"`
		Properties map[string]interface{} `nbt:"states"`
		Version    int32                  `nbt:"version"`
	} `nbt:"block"`
	
	NBTData map[string]interface{} `nbt:"data"`
}


func (o outputItem) Stack() (item.Stack, bool) {
	it, ok := world.ItemByName(o.Name, int16(o.Meta))
	if !ok {
		return item.Stack{}, false
	}
	if n, ok := it.(world.NBTer); ok {
		it = n.DecodeNBT(o.NBTData).(world.Item)
	}

	return item.NewStack(it, int(o.Count)), true
}


type outputItems []outputItem


func (d outputItems) Stacks() ([]item.Stack, bool) {
	s := make([]item.Stack, 0, len(d))
	for _, o := range d {
		itemOutput, ok := o.Stack()
		if !ok {
			return nil, false
		}
		s = append(s, itemOutput)
	}
	return s, true
}
