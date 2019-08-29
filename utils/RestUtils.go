package utils

import (
	"github.com/goplume/core/utils/stereotypes"
)

func EntityPath(i stereotypes.ApiRestInteface) string {
	return i.GetRoot() + "/" + i.GetEntityPath()
}

func EntityIdPath(i stereotypes.ApiRestInteface) string {
	return EntityPath(i) + "/:" + i.GetEntityId()
}

func ClientPathEntities(ri stereotypes.ApiRestInteface) string {
	return ri.GetClientRoot() + "/" + ri.GetEntityPath()
}

func ClientPathEntityId(ri stereotypes.ApiRestInteface) string {
	return ClientPathEntities(ri) + "/{" + ri.GetEntityId() + "}"
}
