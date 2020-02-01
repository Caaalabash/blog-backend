package model

import (
	"blog-go/config"
	"fmt"
	"gopkg.in/mgo.v2"
	"time"
)

var session *mgo.Session

type SessionStore struct {
	session *mgo.Session
}

func init() {
	sess, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:     []string{config.MongoURL},
		Database:  "blog",
		PoolLimit: 4096,
		Timeout:   time.Second * 5,
	})
	if err != nil {
		fmt.Println("初始化mongodb连接失败，请检查Addrs")
		panic(err)
	}
	session = sess
	session.SetMode(mgo.Monotonic, true)
}

func (s *SessionStore) C(name string) *mgo.Collection {
	return s.session.DB("blog").C(name)
}

func (s *SessionStore) Close() {
	s.session.Close()
}

func GetConn() *SessionStore {
	return &SessionStore{
		session: session.Copy(),
	}
}
