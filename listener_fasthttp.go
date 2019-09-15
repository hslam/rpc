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
	maxConnNum		int
	connNum			int
}
func ListenFASTHTTP(address string,server *Server) (Listener, error) {
	router := fasthttprouter.New()
	listener:=  &FASTHTTPListener{address:address,fastrouter:router,server:server,maxConnNum:DefaultMaxConnNum*server.asyncMax}
	return listener,nil
}

func (l *FASTHTTPListener)Serve() (error) {
	log.Allf( "%s", "Waiting for clients")
	workerChan := make(chan bool,l.maxConnNum)
	connChange := make(chan int)
	go func() {
		for conn_change := range connChange {
			l.connNum += conn_change
		}
	}()
	l.fastrouter.POST("/", func (ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		var RemoteAddr=ctx.RemoteAddr().String()
		workerChan<-true
		func() {
			defer func() {
				if err := recover(); err != nil {
				}
				<-workerChan
			}()
			connChange <- 1
			log.Infof("new client %s comming",RemoteAddr)
			ServeHTTP(l.server,ctx,ctx.Request.Body())
			log.Infof("client %s exiting",RemoteAddr)
			connChange <- -1
		}()
	})
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


