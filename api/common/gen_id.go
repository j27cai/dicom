package common

import "github.com/lithammer/shortuuid"


func GenShortUUID() string {
	id := shortuuid.New()
	return id
}
