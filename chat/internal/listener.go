package internal

type Listener struct {
	name string
	c    chan string
}

func NewListener(name string) *Listener {
	return &Listener{
		name: name,
		c:    make(chan string, 10),
	}
}

func (l *Listener) close() {
	close(l.c)
}

func (l *Listener) Listen(f func(msg string) error) error {
	for msg := range l.c {
		if err := f(msg); err != nil {
			return err
		}
	}
	return nil
}
