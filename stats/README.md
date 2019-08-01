```
type WrkClient struct {
}

func (c *WrkClient)Call()(int64,bool){
    //To Do return true
	return 1024,false
}
```

```
	var wrkClients []stats.Client
	parallel:=1
	total_calls:=1000000
	wrkClients[0]= &WrkClient{}
	stats.StartStats(parallel,total_calls,wrkClients)
```

