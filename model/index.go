package model

import (
	"gopkg.in/mgo.v2"
)

var session *mgo.Session

type SessionStore struct {
	session *mgo.Session
}

func init() {
	sess, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:     []string{"dockerhost"},
		Database:  "blog",
		PoolLimit: 4096,
	})
	if err != nil {
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
