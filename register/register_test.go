package register

import "testing"

func TestNew(t *testing.T) {
	// abnormal case
	node := New(nil, "/lzw/20180715", "make it work", 5, 30)
	if node != nil {
		t.Errorf("Expected return nil while to have invalid parameters")
	}
	etcdAddrs := []string{"127.0.0.1:2379"}
	serviceName := "/lzw/2018/07/15"
	serviceValue := "make it work"
	dialTimeout := 5
	ttl := 30
	node = New(etcdAddrs, "", serviceValue, dialTimeout, ttl)
	if node != nil {
		t.Errorf("Expected return nil while to have invalid parameters")
	}
	node = New(etcdAddrs, serviceName, serviceValue, 0, ttl)
	if node != nil {
		t.Errorf("Expected return nil while to have invalid parameters")
	}
	// normal case
	node = New(etcdAddrs, serviceName, serviceValue, dialTimeout, ttl)
	if node == nil {
		t.Errorf("Expected NOT nil")
	}
	err := node.Register()
	if err != nil {
		t.Errorf("Expected nil but Got " + err.Error())
	}
	// TODO: check with etcdClient.Get()
}
