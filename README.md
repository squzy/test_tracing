# Example of application monitoring on golang

That project have 2 service:

1. GIN - http server run on 7878
2. GRPC - grpc server run on 7879

## Tracing:

On that page you can see every request on every service together in one transaction

[**LINK**](Screenshot 2020-07-24 at 22.33.36)

![Page](https://github.com/squzy/test_tracing/blob/master/example.png)

### Gin

That server have one route which did request on grpc server:

```go
engine.GET("hello", func(context *gin.Context) {
    res, err := clint.Echo(context, &service.EchoMsg{})
    if err != nil {
        context.AbortWithError(200, err)
        return
    }
    context.JSON(200, res)
})
```

### GRPC

GRPC server accept that request and did http call to external domain:

```go
func (s *server) Echo(ctx context.Context, msg *service.EchoMsg) (*service.EchoMsg, error) {
	trx := core.GetTransactionFromContext(ctx).CreateTransaction("calcaulate time", api.TransactionType_TRANSACTION_TYPE_INTERNAL, nil)

	client := http.Client{Transport: sHttp.NewRoundTripper(s.application, nil)}

	req, err := http.NewRequest("GET", "https://api.exchangeratesapi.io/latest?base=USD", nil)

	if err != nil {
		return nil, err
	}
	_, err = client.Do(sHttp.NewRequest(trx, req))

	trx.End(nil)
	return &service.EchoMsg{}, nil
}
```

## Cronjob

On examples service we have cron with:

```bash
0 * * * * /usr/bin/curl http://localhost:7878/hello
```

It produce transaction every hour all transaction [here](https://demo.squzy.app/applications/5eef71dcaac3ab3dc67a4ef3/list)



