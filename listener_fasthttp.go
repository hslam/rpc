package rpc

import (
	"sync"
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
	"hslam.com/mgit/Mort/rpc/log"
)

type FASTHTTPListener struct {
	reqMutex 		sync.Mutex
	address			string
	serialize		string
	fastrouter		*fasthttprouter.Router
}
func ListenFASTHTTP(address string) (Listener, error) {
	router := fasthttprouter.New()
	router.PUT("/", func (ctx *fasthttp.RequestCtx) {
		var RemoteAddr=ctx.RemoteAddr().String()
		log.AllInfof("new client %s comming",RemoteAddr)
		if useWorkerPool{
			workerPool.Process(func(obj interface{}, args ...interface{}) interface{} {
				var requestCtx = obj.(*fasthttp.RequestCtx)
				return ServeFASTHTTP(requestCtx,requestCtx.Request.Body())
			},ctx)
		}else {
			ServeFASTHTTP(ctx,ctx.Request.Body())
		}
		log.AllInfof("client %s exiting",RemoteAddr)
	})
	listener:=  &FASTHTTPListener{address:address,fastrouter:router}
	return listener,nil
}

func (l *FASTHTTPListener)Serve() (error) {
	log.Allf( "%s", "Waiting for clients")
	err:=fasthttp.ListenAndServe(l.address, l.fastrouter.Handler)
	if err!=nil{
		log.Fatalf("fatal error: %s", err)
	}
	return nil
}
func (l *FASTHTTPListener)Addr() (string) {
	return l.address
}
func ServeFASTHTTP(ctx *fasthttp.RequestCtx,data []byte)error {
	_,res_bytes, _ := ServeRPC(data)
	if res_bytes!=nil{
		_,err:=ctx.Write(res_bytes)
		return err
	}else {
		return nil
	}
	return ErrConnExit
}

