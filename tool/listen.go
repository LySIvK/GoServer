// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package netutil provides network tool functions, complementing the more
// common ones in the net package.
package tool

import (
	"net"
	"net/http"
	"sync"
)

// LimitListener returns a Listener that accepts at most n simultaneous
// connections from the provided Listener.
func LimitListener(l net.Listener, n int) net.Listener {
	return &limitListener{l, make(chan struct{}, n)}
}

type limitListener struct {
	net.Listener
	sem chan struct{}
}

func (l *limitListener) acquire() { l.sem <- struct{}{} }
func (l *limitListener) release() { <-l.sem }

func (l *limitListener) Accept() (net.Conn, error) {
	l.acquire()
	c, err := l.Listener.Accept()
	if err != nil {
		l.release()
		return nil, err
	}
	return &limitListenerConn{Conn: c, release: l.release}, nil
}

type limitListenerConn struct {
	net.Conn
	releaseOnce sync.Once
	release     func()
}

func (l *limitListenerConn) Close() error {
	err := l.Conn.Close()
	l.releaseOnce.Do(l.release)
	return err
}

func HttpLimitListen(addr string, max int) error {
	if max <= 0 {
		return http.ListenAndServe(addr, nil)
	} else {
		l, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}
		defer l.Close()
		l = LimitListener(l, max)
		return http.Serve(l, nil)
	}

	return nil
}
