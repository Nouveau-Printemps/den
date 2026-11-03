package den

func ExampleDial() {
	c, err := Dial("client test", "unix", "/run/dend.sock")
	if err != nil {
		panic(err)
	}
	defer c.Close()
}

func ExampleDialTLS() {
	c, err := DialTLS("client test", "unix", "/run/dend.sock", nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()
}

func ExampleClient_Authn() {
	c, err := Dial("client test", "unix", "/run/dend.sock")
	if err != nil {
		panic(err)
	}
	defer c.Close()
	res, err := c.Authn("foo", "bar")
	if err != nil {
		panic(err)
	}
	if res == nil {
		// auth failed
	}
	// auth succeed
}
