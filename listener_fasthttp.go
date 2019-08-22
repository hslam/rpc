package rpc

import (
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
	"hslam.com/mgit/Mort/rpc/log"
)

type FASTHTTPListener struct {
	address			string
	fastrouter		*fasthttprouter.Router
}
func ListenFASTHTTP(address string) (Listener, error) {
	router := fasthttprouter.New()
	router.POST("/", func (ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		var RemoteAddr=ctx.RemoteAddr().String()
		log.AllInfof("new client %s comming",RemoteAddr)
		if useWorkerPool{
			workerPool.Process(func(obj interface{}, args ...interface{}) interface{} {
				var requestCtx = obj.(*fasthttp.RequestCtx)
				return ServeHTTP(requestCtx,requestCtx.Request.Body())
			},ctx)
		}else {
			ServeHTTP(ctx,ctx.Request.Body())
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


