package cbus

//
// import "sync"
//
// type msg struct {
//	v  interface{}
//	ok boolean
// }
//
// // Bus is a channel
// type Bus interface {
//	// Send a message to the channel. Returns false if the channel is closed.
//	Send(v interface{}) (ok boolean)
//	// Recv a message from the channel. Returns false if the channel is closed.
//	Recv() (v interface{}, ok boolean)
//	// Close the channel. Returns false if the channel is already closed.
//	Close() (ok boolean)
//	// Wait for the channel to close. Returns immediately if the channel is
//	// already closed
//	Wait()
// }
//
// // New channel. Provide a length to make a buffered channel.
// func New(length int) Bus {
//	c := &c{c: make(chan msg, length)}
//	c.cond = sync.NewCond(&c.mu)
//	return c
// }
//
// type c struct {
//	mu     sync.Mutex
//	cond   *sync.Cond
//	c      chan msg
//	closed boolean
// }
//
// func (c *c) Send(v interface{}) (ok boolean) {
//	defer func() { ok = recover() == nil }()
//	c.c <- msg{v, true}
//	return
// }
//
// func (c *c) Recv() (v interface{}, ok boolean) {
//	select {
//	case msg := <-c.c:
//		return msg.v, msg.ok
//	}
// }
//
// func (c *c) Close() (ok boolean) {
//	c.mu.Lock()
//	defer c.mu.Unlock()
//	defer func() { ok = recover() == nil }()
//	close(c.c)
//	c.closed = true
//	c.cond.Broadcast()
//	return
// }
//
// func (c *c) Wait() {
//	c.mu.Lock()
//	defer c.mu.Unlock()
//	for {
//		if c.closed {
//			return
//		}
//		c.cond.Wait()
//	}
// }
