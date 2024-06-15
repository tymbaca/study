package brokers

import (
	"io"
	"net"
)

type Broker struct {
	addr string
	cfg  Config
	// recievers
	rs []io.Writer
	w  io.Writer
	l  net.Listener
}

func (b *Broker) Write(data []byte) (int, error) {
	return len(data), nil
}

type Config struct {
	Recievers []string
}

func New(addr string, cfg Config) (*Broker, error) {
	rs := make([]io.Writer, 0, len(cfg.Recievers))

	for _, rAddr := range cfg.Recievers {
		conn, err := net.Dial("tcp", rAddr)
		if err != nil {
			return nil, err
		}
		r := newReciever(rAddr, conn)
		rs = append(rs, r)
	}

	w := io.MultiWriter(rs...)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Broker{
		addr: addr,
		cfg:  cfg,
		w:    w,
		l:    l,
	}, nil
}

func (b *Broker) Reconnect() {

}
