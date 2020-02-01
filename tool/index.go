package tool

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"unsafe"
)

func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func ObjectID2Str(id bson.ObjectId) string {
	return fmt.Sprintf("%x", string(id))
}

func MarkdownToHtml(md string) template.HTML {
	b := Str2Bytes(md)
	b = blackfriday.Run(b, blackfriday.WithNoExtensions())
	return template.HTML(Bytes2Str(b))
}