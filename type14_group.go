/*
* File Name:	type14_group.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-08-19
 */
package godmi

import (
	"fmt"
)

type GroupAssociationsItem struct {
	Type   SMBIOSStructureType
	Handle SMBIOSStructureHandle
}

type GroupAssociations struct {
	infoCommon
	GroupName string
	Item      []GroupAssociationsItem
}

func (g GroupAssociations) String() string {
	return fmt.Sprintf("Group Associations:\n"+
		"\tGroup Name: %s\n"+
		"\tItem: %#v\n",
		g.GroupName,
		g.Item)
}

func newGroupAssociations(h dmiHeader) dmiTyper {
	var ga GroupAssociations
	data := h.data
	ga.GroupName = h.FieldString(int(data[0x04]))
	cnt := (h.Length - 5) / 3
	items := data[5:]
	var i byte
	for i = 0; i < cnt; i++ {
		var gai GroupAssociationsItem
		gai.Type = SMBIOSStructureType(items[i*3])
		gai.Handle = SMBIOSStructureHandle(u16(items[i*3+1:]))
		ga.Item = append(ga.Item, gai)
	}
	return &ga
}

func GetGroupAssociations() *GroupAssociations {
	if d, ok := gdmi[SMBIOSStructureTypeGroupAssociations]; ok {
		return d.(*GroupAssociations)
	}
	return nil
}

func init() {
	addTypeFunc(SMBIOSStructureTypeGroupAssociations, newGroupAssociations)
}
