### go test -v -run="none" -bench=. -benchtime=3s  -benchmem
```
type WrkClient struct {
}

func (c *WrkClient)Call()bool{
    //To Do return true
	return false
}

```

```
	var wrkClients []stats.Client
	parallel:=1
	total_calls:=1000000
	wrkClients[0]= &WrkClient{}
	stats.StartStats(parallel,total_calls,wrkClients)
```

