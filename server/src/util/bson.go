package util

import "gopkg.in/mgo.v2/bson"

func StructToBsonM(obj interface{}) (m bson.M) {
	if obj == nil {
		return
	}
	return
}
