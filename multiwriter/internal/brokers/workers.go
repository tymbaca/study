package brokers

import (
	"errors"
	"fmt"
	"net"
)

type reciever struct {
	addr string
	c    net.Conn
}

func newReciever(addr string, conn net.Conn) *reciever {
	return &reciever{addr, conn}
}

func (r *reciever) Write(data []byte) (int, error) {
	n, err := r.c.Write(data)
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			return n, &RecieverClosedError{r: r, err: err}
		}
		return n, err
	}

	return n, nil
}

type RecieverClosedError struct {
	r   *reciever
	err error
}

func (e *RecieverClosedError) Error() string {
	return fmt.Errorf("write to reciever with addr '%s' but it closed: %w", e.r.addr, e.err).Error()
}
