package rpc

import (
	bytesTool "bytes"
	"github.com/hslam/rpc/examples/service/bytes"
	"github.com/hslam/rpc/examples/service/code"
	"github.com/hslam/rpc/examples/service/gob"
	"github.com/hslam/rpc/examples/service/json"
	service "github.com/hslam/rpc/examples/service/pb"
	"github.com/hslam/rpc/examples/service/xml"
	"testing"
	"time"
)

func TestIPC(t *testing.T) {
	{
		go func() {
			Register(new(service.Arith))
			ListenAndServe("ipc", "/tmp/ipc")
		}()
	}
	time.Sleep(time.Second)
	{
		conn, err := Dial("ipc", "/tmp/ipc", "pb") //tcp|ws|quic|http
		if err != nil {
			t.Fatalf("dailing error: %s", err.Error())
		}
		req := &service.ArithRequest{A: 9, B: 2}
		var res service.ArithResponse
		err = conn.Call("Arith.Multiply", req, &res)
		if err != nil {
			t.Fatalf("arith error: %s", err.Error())
		}
		if res.Pro != 18 {
			t.Errorf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
		}
		conn.Close()
	}
	Close()
	time.Sleep(time.Second)
}
func TestCodec(t *testing.T) {
	{
		{
			go func() {
				Register(new(service.Arith))
				ListenAndServe("ipc", "/tmp/ipc")
			}()
		}
		time.Sleep(time.Second)
		{
			conn, err := Dial("ipc", "/tmp/ipc", "pb") //tcp|ws|quic|http
			if err != nil {
				t.Fatalf("dailing error: %s", err.Error())
			}
			req := &service.ArithRequest{A: 9, B: 2}
			var res service.ArithResponse
			err = conn.Call("Arith.Multiply", req, &res)
			if err != nil {
				t.Fatalf("arith error: %s", err.Error())
			}
			if res.Pro != 18 {
				t.Errorf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
			}
			conn.Close()
		}
		Close()
		time.Sleep(time.Second)
	}

	{
		{
			go func() {
				Register(new(json.Arith))
				ListenAndServe("ipc", "/tmp/ipc")
			}()
		}
		time.Sleep(time.Second)
		{
			conn, err := Dial("ipc", "/tmp/ipc", "json") //tcp|ws|quic|http
			if err != nil {
				t.Fatalf("dailing error: %s", err.Error())
			}
			req := &json.ArithRequest{A: 9, B: 2}
			var res json.ArithResponse
			err = conn.Call("Arith.Multiply", req, &res)
			if err != nil {
				t.Fatalf("arith error: %s", err.Error())
			}
			if res.Pro != 18 {
				t.Errorf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
			}
			conn.Close()
		}
		Close()
		time.Sleep(time.Second)
	}
	{
		{
			go func() {
				Register(new(xml.Arith))
				ListenAndServe("ipc", "/tmp/ipc")
			}()
		}
		time.Sleep(time.Second)
		{
			conn, err := Dial("ipc", "/tmp/ipc", "xml") //tcp|ws|quic|http
			if err != nil {
				t.Fatalf("dailing error: %s", err.Error())
			}
			req := &xml.ArithRequest{A: 9, B: 2}
			var res xml.ArithResponse
			err = conn.Call("Arith.Multiply", req, &res)
			if err != nil {
				t.Fatalf("arith error: %s", err.Error())
			}
			if res.Pro != 18 {
				t.Errorf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
			}
			conn.Close()
		}
		Close()
		time.Sleep(time.Second)
	}
	{
		{
			go func() {
				Register(new(gob.Arith))
				ListenAndServe("ipc", "/tmp/ipc")
			}()
		}
		time.Sleep(time.Second)
		{
			conn, err := Dial("ipc", "/tmp/ipc", "gob") //tcp|ws|quic|http
			if err != nil {
				t.Fatalf("dailing error: %s", err.Error())
			}
			req := &gob.ArithRequest{A: 9, B: 2}
			var res gob.ArithResponse
			err = conn.Call("Arith.Multiply", req, &res)
			if err != nil {
				t.Fatalf("arith error: %s", err.Error())
			}
			if res.Pro != 18 {
				t.Errorf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
			}
			conn.Close()
		}
		Close()
		time.Sleep(time.Second)
	}
	{
		{
			go func() {
				Register(new(code.Arith))
				ListenAndServe("ipc", "/tmp/ipc")
			}()
		}
		time.Sleep(time.Second)
		{
			conn, err := Dial("ipc", "/tmp/ipc", "code") //tcp|ws|quic|http
			if err != nil {
				t.Fatalf("dailing error: %s", err.Error())
			}
			req := &code.ArithRequest{A: 9, B: 2}
			var res code.ArithResponse
			err = conn.Call("Arith.Multiply", req, &res)
			if err != nil {
				t.Fatalf("arith error: %s", err.Error())
			}
			if res.Pro != 18 {
				t.Errorf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
			}
			conn.Close()
		}
		Close()
		time.Sleep(time.Second)
	}
	{
		{
			go func() {
				Register(new(bytes.Echo))
				ListenAndServe("ipc", "/tmp/ipc")
			}()
		}
		time.Sleep(time.Second)
		{
			conn, err := Dial("ipc", "/tmp/ipc", "bytes") //tcp|ws|quic|http
			if err != nil {
				t.Fatalf("dailing error: %s", err.Error())
			}
			var req = []byte("Hello World")
			var res []byte
			err = conn.Call("Echo.ToLower", &req, &res)
			if err != nil {
				t.Fatalf("Echo error: %s", err.Error())
			}
			if !bytesTool.Equal(res, []byte("hello world")) {
				t.Errorf("Echo.ToLower : %s\n", string(res))

			}
			conn.Close()
		}
		Close()
		time.Sleep(time.Second)
	}

}

func TestPipelining(t *testing.T) {
	networks := []string{"tcp", "quic", "ws", "udp", "http"}
	for _, network := range networks {
		{
			go func() {
				Register(new(service.Arith))
				SetPipelining(true)
				ListenAndServe(network, ":8080")
			}()
			time.Sleep(time.Second)
		}
		{
			opts := DefaultOptions()
			opts.SetPipelining(true)
			conn, err := DialWithOptions(network, "127.0.0.1:8080", "pb", opts)
			if err != nil {
				t.Fatalf("dailing error: %s", err.Error())
			}
			req := &service.ArithRequest{A: 9, B: 2}
			var res service.ArithResponse
			err = conn.Call("Arith.Multiply", req, &res)
			if err != nil {
				t.Fatalf("arith error: %s", err.Error())
			}
			if res.Pro != 18 {
				t.Errorf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
			}
			conn.Close()
		}
		Close()
		time.Sleep(time.Second)
	}
}

func TestPipeliningAndBatching(t *testing.T) {
	networks := []string{"tcp", "quic", "ws", "udp", "http"}
	for _, network := range networks {
		{
			go func() {
				Register(new(service.Arith))
				SetPipelining(true)
				SetBatching(true)
				ListenAndServe(network, ":8080")
			}()
			time.Sleep(time.Second)
		}
		{
			opts := DefaultOptions()
			opts.SetPipelining(true)
			opts.SetBatching(true)
			conn, err := DialWithOptions(network, "127.0.0.1:8080", "pb", opts)
			if err != nil {
				t.Fatalf("dailing error: %s", err.Error())
			}
			req := &service.ArithRequest{A: 9, B: 2}
			var res service.ArithResponse
			err = conn.Call("Arith.Multiply", req, &res)
			if err != nil {
				t.Fatalf("arith error: %s", err.Error())
			}
			if res.Pro != 18 {
				t.Errorf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
			}
			conn.Close()
		}
		Close()
		time.Sleep(time.Second)
	}
}

func TestMultiplexing(t *testing.T) {
	networks := []string{"tcp", "quic", "ws", "udp", "http"}
	for _, network := range networks {
		{
			go func() {
				Register(new(service.Arith))
				SetMultiplexing(true)
				ListenAndServe(network, ":8080")
			}()
			time.Sleep(time.Second)
		}
		{
			opts := DefaultOptions()
			opts.SetMultiplexing(true)
			conn, err := DialWithOptions(network, "127.0.0.1:8080", "pb", opts)
			if err != nil {
				t.Fatalf("dailing error: %s", err.Error())
			}
			req := &service.ArithRequest{A: 9, B: 2}
			var res service.ArithResponse
			err = conn.Call("Arith.Multiply", req, &res)
			if err != nil {
				t.Fatalf("arith error: %s", err.Error())
			}
			if res.Pro != 18 {
				t.Errorf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
			}
			conn.Close()
		}
		Close()
		time.Sleep(time.Second)
	}
}

func TestMultiplexingAndBatching(t *testing.T) {
	networks := []string{"tcp", "quic", "ws", "udp", "http"}
	for _, network := range networks {
		{
			go func() {
				Register(new(service.Arith))
				SetMultiplexing(true)
				SetBatching(true)
				ListenAndServe(network, ":8080")
			}()
			time.Sleep(time.Second)
		}
		{
			opts := DefaultOptions()
			opts.SetMultiplexing(true)
			opts.SetBatching(true)
			conn, err := DialWithOptions(network, "127.0.0.1:8080", "pb", opts)
			if err != nil {
				t.Fatalf("dailing error: %s", err.Error())
			}
			req := &service.ArithRequest{A: 9, B: 2}
			var res service.ArithResponse
			err = conn.Call("Arith.Multiply", req, &res)
			if err != nil {
				t.Fatalf("arith error: %s", err.Error())
			}
			if res.Pro != 18 {
				t.Errorf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
			}
			conn.Close()
		}
		Close()
		time.Sleep(time.Second)
	}
}
