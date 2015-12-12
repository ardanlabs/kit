package udp_test

import (
	"bytes"
	"io"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ardanlabs/kit/tests"
	"github.com/ardanlabs/kit/udp"
)

func init() {
	tests.Init("KIT")
}

//==============================================================================

// TestUDP provide a test of listening for a connection and
// echoing the data back.
func TestUDP(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to listen and process UDP data.")
	{
		// Create a configuration.
		cfg := udp.Config{
			NetType: "udp4",
			Addr:    ":0",

			ConnHandler: udpConnHandler{},
			ReqHandler:  udpReqHandler{},
			RespHandler: udpRespHandler{},

			OptIntPool: udp.OptIntPool{
				RecvMinPoolSize: func() int { return 2 },
				RecvMaxPoolSize: func() int { return 1000 },
				SendMinPoolSize: func() int { return 2 },
				SendMaxPoolSize: func() int { return 1000 },
			},
		}

		// Create a new UDP value.
		u, err := udp.New("TEST", "TEST", &cfg)
		if err != nil {
			t.Fatal("\tShould be able to create a new UDP listener.", tests.Failed, err)
		}
		t.Log("\tShould be able to create a new UDP listener.", tests.Success)

		// Start accepting client data.
		if err := u.Start("TEST"); err != nil {
			t.Fatal("\tShould be able to start the UDP listener.", tests.Failed, err)
		}
		t.Log("\tShould be able to start the UDP listener.", tests.Success)

		defer u.Stop("TEST")

		// Let's connect back and send a UDP package
		conn, err := net.Dial("udp4", u.Addr().String())
		if err != nil {
			t.Fatal("\tShould be able to dial a new UDP connection.", tests.Failed, err)
		}
		t.Log("\tShould be able to dial a new UDP connection.", tests.Success)

		// Send some know data to the udp listener.
		b := bytes.NewBuffer([]byte{0x01, 0x3D, 0x06, 0x00, 0x58, 0x68, 0x9b, 0x9d, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFC, 0x00, 0x01})
		b.WriteTo(conn)

		// Setup a limit reader to extract the response.
		lr := io.LimitReader(conn, 6)

		// Let's read the response.
		data := make([]byte, 6)
		if _, err := lr.Read(data); err != nil {
			t.Fatal("\tShould be able to read the response from the connection.", tests.Failed, err)
		}
		t.Log("\tShould be able to read the response from the connection.", tests.Success)

		response := string(data)

		if response == "GOT IT" {
			t.Log("\tShould receive the string \"GOT IT\".", tests.Success)
		} else {
			t.Error("\tShould receive the string \"GOT IT\".", tests.Failed, response)
		}

		d := atomic.LoadInt64(&dur)
		duration := time.Duration(d)

		if duration <= 2*time.Second {
			t.Log("\tShould be less that 2 seconds.", tests.Success)
		} else {
			t.Error("\tShould be less that 2 seconds.", tests.Failed, duration)
		}
	}
}

// Test udp.Addr works correctly.
func TestUDPAddr(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to listen on any port and know that bound UDP address.")
	{
		// Create a configuration.
		cfg := udp.Config{
			NetType: "udp4",
			Addr:    ":0", // Defer port assignment to OS.

			ConnHandler: udpConnHandler{},
			ReqHandler:  udpReqHandler{},
			RespHandler: udpRespHandler{},

			OptIntPool: udp.OptIntPool{
				RecvMinPoolSize: func() int { return 2 },
				RecvMaxPoolSize: func() int { return 1000 },
				SendMinPoolSize: func() int { return 2 },
				SendMaxPoolSize: func() int { return 1000 },
			},
		}

		// Create a new UDP value.
		u, err := udp.New("TEST", "TEST", &cfg)
		if err != nil {
			t.Fatal("\tShould be able to create a new UDP listener.", tests.Failed, err)
		}
		t.Log("\tShould be able to create a new UDP listener.", tests.Success)

		// Addr should be nil before start.
		if addr := u.Addr(); addr != nil {
			t.Fatalf("\tAddr() should be nil before Start; Addr() = %q. %s", addr, tests.Failed)
		}
		t.Log("\tAddr() should be nil before Start.", tests.Success)

		// Start accepting client data.
		if err := u.Start("TEST"); err != nil {
			t.Fatal("\tShould be able to start the UDP listener.", tests.Failed, err)
		}
		defer u.Stop("TEST")

		// Addr should be non-nil after Start.
		addr := u.Addr()
		if addr == nil {
			t.Fatal("\tAddr() should be not be nil after Start.", tests.Failed)
		}
		t.Log("\tAddr() should be not be nil after Start.", tests.Success)

		// The OS should assign a random open port, which shouldn't be 0.
		_, port, err := net.SplitHostPort(addr.String())
		if err != nil {
			t.Fatalf("\tSplitHostPort should not fail. tests.Failed %v. %s", err, tests.Failed)
		}
		if port == "0" {
			t.Fatalf("\tAddr port should not be %q. %s", port, tests.Failed)
		}
		t.Logf("\tAddr() should be not be 0 after Start (port = %q). %s", port, tests.Success)
	}
}

// Test generic UDP write timeout.
func TestUDPWriteTimeout(t *testing.T) {
	t.Log("Given the need to get a timeout error on UDP write.")
	{
		localAddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("Should not have tests.Failed to resolve UDP address. Err[%v] %s", err, tests.Failed)
		}
		remoteAddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:1234")
		if err != nil {
			t.Fatalf("Should not have tests.Failed to resolve UDP address. Err[%v] %s", err, tests.Failed)
		}

		conn, err := net.ListenUDP("udp4", localAddr)
		if err != nil {
			t.Fatalf("Should not have tests.Failed to create UDP connection. Err[%v] %s", err, tests.Failed)
		}

		if err := conn.SetWriteDeadline(time.Now().Add(1)); err != nil {
			t.Fatalf("Should not have tests.Failed to set write deadline. Err[%v] %s", err, tests.Failed)
		}

		const str = "String to send via UDP socket for testing purposes."

		_, err = conn.WriteToUDP([]byte(str), remoteAddr)
		if err == nil {
			t.Fatalf("Should not gotten an error %s", tests.Failed)
		}

		t.Logf("Got error [%T: %v] %s", err, err, tests.Success)

		opError, ok := err.(*net.OpError)
		if !ok {
			t.Fatalf("Should have gotten *net.OpError, got [%T: %v] %s", err, err, tests.Failed)
		}

		if !opError.Timeout() {
			t.Fatalf("Should have gotten a timeout error %s", tests.Failed)
		}

		t.Logf("Got timeout error %s", tests.Success)
	}
}
