package ws

import (
	"net/http"
)

// Echo sends every message back unchanged, keeping its type.
func (c *Controller) Echo(w http.ResponseWriter, r *http.Request) {
	s := c.accept(w, r)
	if s == nil {
		return
	}
	defer s.done()

	s.echo()
}

// echo reads until the client goes away, writing each message back.
func (s *session) echo() {
	for {
		typ, data, err := s.conn.Read(s.ctx)
		if err != nil {
			return
		}
		if err := s.write(typ, data); err != nil {
			return
		}
	}
}
