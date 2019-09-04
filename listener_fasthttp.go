package rpc

import (
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
	"hslam.com/mgit/Mort/rpc/log"
)

type FASTHTTPListener struct {
	server			*Server
	address			string
	fastrouter		*fasthttprouter.Router
}
func ListenFASTHTTP(address string,server *Server) (Listener, error) {
	router := fasthttprouter.New()
	router.POST("/", func (ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		var RemoteAddr=ctx.RemoteAddr().String()
		log.AllInfof("new client %s comming",RemoteAddr)
		if server.useWorkerPool{
			server.workerPool.Process(func(obj interface{}, args ...interface{}) interface{} {
				var requestCtx = obj.(*fasthttp.RequestCtx)
				var server = args[0].(*Server)
				return ServeHTTP(server,requestCtx,requestCtx.Request.Body())
			},ctx,server)
		}else {
			ServeHTTP(server,ctx,ctx.Request.Body())
		}
		log.AllInfof("client %s exiting",RemoteAddr)
	})
	listener:=  &FASTHTTPListener{address:address,fastrouter:router,server:server}
	return listener,nil
}

func (l *FASTHTTPListener)Serve() (error) {
	log.Allf( "%s", "Waiting for clients")
	err:=fasthttp.ListenAndServe(l.address, l.fastrouter.Handler)
	if err!=nil{
		log.Errorf("fatal error: %s", err)
		return err
	}
	return nil
}
func (l *FASTHTTPListener)Addr() (string) {
	return l.address
}


