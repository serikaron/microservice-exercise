package internal_test

import (
	"errors"
	"mse/chat/internal"
	"testing"
)

func Test_notify_to_listener(t *testing.T) {
	t.Parallel()
	hub := internal.NewSafeHub()
	marry, mErr := addSpyToHub(hub, "Marry")
	tom, tErr := addSpyToHub(hub, "Tom")

	hub.Notify("Greetings")
	hub.Close()

	messageOf(marry).withT(t).shouldBe("Greetings")
	messageOf(tom).withT(t).shouldBe("Greetings")
	errorOf(mErr).withT(t).shouldBe(nil)
	errorOf(tErr).withT(t).shouldBe(nil)
}

func Test_not_notify_to_removed_listener(t *testing.T) {
	t.Parallel()
	hub := internal.NewSafeHub()
	marry, mErr := addSpyToHub(hub, "Marry")
	tom, tErr := addSpyToHub(hub, "Tom")

	hub.Notify("Greetings")
	hub.Remove("Marry")
	hub.Notify("Good-bye")
	hub.Close()

	messageOf(marry).withT(t).shouldBe("Greetings")
	messageOf(tom).withT(t).shouldBe("Greetings", "Good-bye")
	errorOf(mErr).withT(t).shouldBe(nil)
	errorOf(tErr).withT(t).shouldBe(nil)
}

func Test_failed_listener_will_return_error(t *testing.T) {
	t.Parallel()
	hub := internal.NewSafeHub()
	marry, mErr := addSpyToHub(hub, "Marry")
	failed, fErr := addFailedToHub(hub, "failed")
	tom, tErr := addSpyToHub(hub, "tom")

	hub.Notify("Greetings")
	hub.Close()

	messageOf(marry).withT(t).shouldBe("Greetings")
	messageOf(failed).withT(t).times(0)
	messageOf(tom).withT(t).shouldBe("Greetings")
	errorOf(mErr).withT(t).shouldBe(nil)
	errorOf(fErr).withT(t).shouldBe(errors.New("expect error"))
	errorOf(tErr).withT(t).shouldBe(nil)
}

type spyChan chan string

type spyResult struct {
	c     spyChan
	t     *testing.T
	count int
}

func messageOf(c spyChan) *spyResult {
	return &spyResult{c, nil, 0}
}

func (s *spyResult) withT(t *testing.T) *spyResult {
	s.t = t
	return s
}

func (s *spyResult) shouldBe(msg ...string) *spyResult {
	for m := range s.c {
		if s.count >= len(msg) {
			s.t.Fatalf("receive %d messages, wants %d", s.count, len(msg))
		}
		if m != msg[s.count] {
			s.t.Fatalf("received message is %s, want %s", m, msg[s.count])
		}
		s.count++
	}
	return s
}

func (s *spyResult) times(t int) *spyResult {
	if s.count != t {
		s.t.Fatalf("receive %d messages, wants %d", s.count, t)
	}
	return s
}

type errChan chan error

type errorResult struct {
	c errChan
	t *testing.T
}

func errorOf(c errChan) *errorResult {
	return &errorResult{c, nil}
}

func (r *errorResult) withT(t *testing.T) *errorResult {
	r.t = t
	return r
}

func (r *errorResult) shouldBe(err error) *errorResult {
	e := <-r.c
	got := "<nil>"
	if e != nil {
		got = e.Error()
	}
	want := "<nil>"
	if err != nil {
		want = err.Error()
	}
	if got != want {
		r.t.Fatalf("errorResult.err is %s, want %s", got, want)
	}
	return r
}

func addSpyToHub(hub *internal.SafeHub, name string) (spyChan, errChan) {
	return addToHub(hub, name, func(msg string, c spyChan) error {
		c <- msg
		return nil
	})
}

func addFailedToHub(hub *internal.SafeHub, name string) (spyChan, errChan) {
	return addToHub(hub, name, func(msg string, _ spyChan) error {
		return errors.New("expect error")
	})
}

func addToHub(hub *internal.SafeHub, name string, f func(msg string, c spyChan) error) (spyChan, errChan) {
	c := make(chan string)
	e := make(chan error, 1)
	l := internal.NewListener(name)
	hub.Add(l)
	go func() {
		defer close(c)
		defer close(e)
		err := l.Listen(func(msg string) error {
			return f(msg, c)
		})
		e <- err
	}()
	return c, e
}
