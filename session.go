package gin_session_redis

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)
var key = "akjsdsadsasd#@34dfd"

type Store interface {
	sessions.Store

}

type Session interface {
	ID() string
	Get(key interface{}) interface{}
	Set(key interface{}) interface{}
	Delete(key interface{}, val interface{})
	Clear()

	// Save save all keys
	Save()
}

type session struct {
	name    string
	request *http.Request
	store   Store
	session *sessions.Session
	written bool
	writer  http.ResponseWriter
}

func Sessions(name string, store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := &session{
			name: name,
			request: c.Request,
			store: store,
			session: nil,
			written: false,
			writer: c.Writer,
		}

		c.Set(key, s)

		defer context.Clear(c.Request)

		c.Next()
	}
}

func (s *session) ID() string {
	return s.Session().ID
}

func (s *session) Get(key interface{}) interface{} {
	return s.Session().Values[key]
}

func (s *session) Set(key interface{}, val interface{}) {
	s.Session().Values[key] = val
	s.written = true
}

func (s *session) Delete(key interface{}) {
	delete(s.Session().Values, key)
	s.written = true
}

func (s *session) Clear() {
	for key := range s.Session().Values {
		s.Delete(key)
	}
}

func (s *session) Session() *sessions.Session {
	if s.session == nil {
		var err error

		s.session, err = s.store.Get(s.request, s.name)
		if err != nil {
			log.Println(err.Error())
		}
	}

	return s.session
}

func (s *session) Written() bool {
	return s.written
}

func Default(c *gin.Context) Session {
	return c.MustGet(key).(Session)
}

func DefaultMany(c *gin.Context, name string) Session {
	return c.MustGet(key).(map[string]Session)[name]
}
