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
