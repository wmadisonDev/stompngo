//
// Copyright © 2011-2013 Guy M. Allard
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package stompngo

import (
	"os"
	"testing"
)

/*
	Test Subscribe, no destination.
*/
func TestSubNoSub(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	h := empty_headers
	// Subscribe, no dest
	_, e := c.Subscribe(h)
	if e == nil {
		t.Errorf("Expected subscribe error, got [nil]\n")
	}
	if e != EREQDSTSUB {
		t.Errorf("Subscribe error, expected [%v], got [%v]\n", EREQDSTSUB, e)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test subscribe, no ID.
*/
func TestSubNoIdOnce(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	d := "/queue/subunsub.genl.01"
	h := Headers{"destination", d}
	//
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if s == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}
	select {
	case v := <-c.MessageData:
		t.Errorf("Unexpected frame received, got [%v]\n", v)
	default:
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test subscribe, no ID, twice to same destination.
*/
func TestSubNoIdTwice(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	d := "/queue/subunsub.genl.02"
	h := Headers{"destination", d}
	// First time
	c.protocol = SPL_10 // force
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if s == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}
	select {
	case v := <-c.MessageData:
		t.Errorf("Unexpected frame received, got [%v]\n", v)
	default:
	}
	// Second time
	s, e = c.Subscribe(h)
	if e == EDUPSID {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if s == nil {
		t.Errorf("Expected subscribe channel, got nil\n")
	}

	c.protocol = SPL_11
	s, e = c.Subscribe(h) // This will work
	s, e = c.Subscribe(h) // This should not
	// 1.1+ errors
	if e == nil {
		t.Errorf("Expected subscribe twice  error, got [nil]\n")
	}
	if e != EDUPSID {
		t.Errorf("Subscribe twice error, expected [%v], got [%v]\n", EDUPSID, e)
	}
	if s != nil {
		t.Errorf("Expected nil subscribe channel, got [%v]\n", s)
	}

	select {
	case v := <-c.MessageData:
		t.Errorf("Unexpected frame received, got [%v]\n", v)
	default:
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test send, subscribe, read, unsubscribe.
*/
func TestSubUnsubBasic(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	m := "A message"
	d := "/queue/subunsub.basic.01"
	h := Headers{"destination", d}
	_ = c.Send(h, m)
	//
	h = h.Add("id", d)
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if s == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}
	md := <-s // Read message data
	//
	if md.Error != nil {
		t.Errorf("Expected no message data error, got [%v]\n", md.Error)
	}
	msg := md.Message
	rd := msg.Headers.Value("destination")
	if rd != d {
		t.Errorf("Expected destination [%v], got [%v]\n", d, rd)
	}
	ri := msg.Headers.Value("subscription")
	if ri != d {
		t.Errorf("Expected subscription [%v], got [%v]\n", d, ri)
	}
	//
	e = c.Unsubscribe(h)
	if e != nil {
		t.Errorf("Expected no unsubscribe error, got [%v]\n", e)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test send, subscribe, read, unsubscribe, 1.0 only, no sub id.
*/
func TestSubUnsubBasic10(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") != "" {
		println("TestSubUnsubBasic10 norun")
		return
	}
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	m := "A message"
	d := "/queue/subunsub.basic.r10.01"
	h := Headers{"destination", d}
	_ = c.Send(h, m)
	//
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if s == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}
	md := <-s // Read message data
	//
	if md.Error != nil {
		t.Errorf("Expected no message data error, got [%v]\n", md.Error)
	}
	msg := md.Message
	rd := msg.Headers.Value("destination")
	if rd != d {
		t.Errorf("Expected destination [%v], got [%v]\n", d, rd)
	}
	//
	e = c.Unsubscribe(h)
	if e != nil {
		t.Errorf("Expected no unsubscribe error, got [%v]\n", e)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test establishSubscription.
*/
func TestSubestablishSubscription(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	d := "/queue/estabsub.01"
	h := Headers{"destination", d}
	// First time
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if s == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Test unsubscribe, set subscribe channel capacity.
*/
func TestSubSetCap(t *testing.T) {
	if os.Getenv("STOMP_TEST11p") == "" {
		println("TestSubSetCap norun")
		return
	}
	//
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	c, _ := Connect(n, ch)
	//
	p := 25
	c.SetSubChanCap(p)
	r := c.SubChanCap()
	if r != p {
		t.Errorf("Expected get capacity [%v], got [%v]\n", p, r)
	}
	//
	d := "/queue/subsetcap.basic.01"
	h := Headers{"destination", d, "id", d}
	s, e := c.Subscribe(h)
	if e != nil {
		t.Errorf("Expected no subscribe error, got [%v]\n", e)
	}
	if s == nil {
		t.Errorf("Expected subscribe channel, got [nil]\n")
	}
	if cap(s) != p {
		t.Errorf("Expected subchan capacity [%v], got [%v]\n", p, cap(s))
	}
	//
	e = c.Unsubscribe(h)
	if e != nil {
		t.Errorf("Expected no unsubscribe error, got [%v]\n", e)
	}
	//
	_ = c.Disconnect(empty_headers)
	_ = closeConn(t, n)
}
