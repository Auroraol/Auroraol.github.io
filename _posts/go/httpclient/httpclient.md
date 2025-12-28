## 多客户端管理

```go
type ClientManager struct {
	clients map[string]httpclient.Client
}

var spiClientMgr = &ClientManager{
	clients: make(map[string]httpclient.Client),
}

func InitSpiClientsFromConfig(ctx context.Context) error {
	httpClientCfg := config.GetHttpClientConfig()
	if httpClientCfg == nil {
		return fmt.Errorf("http client config is nil")
	}

	for _, item := range httpClientCfg.Clients {
		if item.Name == "" {
			logger.WithContext(ctx).Warnf("Skipping http client with empty name, address=%s", item.Address)
			continue
		}

		opts := []httpclient.Option{
			httpclient.WithAddress(item.Address),
			httpclient.WithTimeout(item.Timeout * time.Second),
			httpclient.WithPreRequestHooks(hooks.LoggingRequest()),
			httpclient.WithAfterResponseHooks(hooks.LoggingResponse()),
		}

		client, err := httpclient.NewClient(opts...)
		if err != nil {
			logger.WithContext(ctx).Errorf("create http client failed, name=%s, err=%v", item.Name, err)
			return err
		}

		spiClientMgr.clients[item.Name] = client
	}

	logger.WithContext(ctx).Infof("InitSpiClientsFromConfig done, total=%d", len(spiClientMgr.clients))
	return nil
}

func GetSpiClient(name string) (httpclient.Client, error) {
	if c, ok := spiClientMgr.clients[name]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("client(%s) not found", name)
}
```

## Get

```
			_, err := client.NewRequest(ctx).SetQueryParams(map[string]string{
				"code":  authMsg.Code,
				"state": "",
			}).Get(targetPath)
			if err != nil {
				logger.WithContext(ctx).Errorf("failed OpenMsgTypeSellerAuth%s: %v", targetPath, err)
				return err
			}
```

## Post

```
		_, err = client.NewRequest(ctx).SetBody(data.Body).SetQueryParam("secret", appSecret).Post(targetPath)
		if err != nil {
		// 返回的是超时
			logger.WithContext(ctx).Errorf("failed OpenMsgTypePaymentSuccess%s: %v", targetPath, err)
		}
		return nil
```







## 例子

```go

func (c *Client) doPostRequest(ctx context.Context, appName, method string, param map[string]interface{}, ret interface{}, platShopId string) error {
	req := map[string]string{
		"method":      method,
		"app_key":     conf.AppId,
		"timestamp":   strconv.FormatInt(timestamp, 10),
		"sign":        GetSha256Sign(ctx, conf.AppId, conf.AppSecret, method, timestamp, paramJson),
		"sign_method": "hmac-sha256",
		"v":           "2",
	}

	client, err := c.GetClient(appName)
	if err != nil {
		return fmt.Errorf("不支持的应用链接:%s", appName)
	}

	// 限流重试
	var commonApiRes proto.CommonApiRes
	const maxRetries = 1
	var retryCount int
	for {

		if err := ctx.Err(); err != nil {
			return errors.Wrap(err, "context canceled")
		}

       // 请求
		_, err = client.NewRequest(ctx).
			SetQueryParams(req).
			SetBody(paramJson).
			SetResult(ret).
			Post(conf.ApiURI + MethodToUri(ctx, method))
		if err != nil {
			return errors.Wrap(err, "call dy api error")
		}

		if err = copier.Copy(&commonApiRes, ret); err != nil {
			logger.Infof(ctx, "copy commonApiRes err: %v", err)
			return nil
		}

		// 限流重试或应用限流重试
		if (commonApiRes.Code == proto.CODE_BUSINESS_PROCESSING_FAILURE && commonApiRes.SubCode == proto.SUB_CODE_GOODS_REQUEST_LIMITED) ||
			commonApiRes.Code == proto.CODE_TRIGGER_LIMITED {
			logger.Infof(ctx, "call dy api trigger limited error, Code=%d, Msg=%s, SubCode=%s, SubMsg=%s", commonApiRes.Code, commonApiRes.Msg, commonApiRes.SubCode, commonApiRes.SubMsg)
			retryCount++
			if retryCount > maxRetries {
				return errors.New("call dy api trigger limited error")
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Second):
				continue
			}
		}

		return nil
	}
}

func (c *Client) doRequestWithoutToken(ctx context.Context, appName, method string, param map[string]string, ret interface{}) error {
	conf, err := c.GetConf(appName)
	if err != nil {
		return fmt.Errorf("不支持的应用名称:%s", appName)
	}
	req := map[string]string{
		"method":     method,
		"app_key":    conf.AppId,
		"param_json": GetSortJSONString(param),
		"timestamp":  time.Now().Format(proto.TIME_LAYOUT),
		"v":          "2",
	}
	req["sign"] = GetSign(ctx, req, conf.AppSecret)
	client, err := c.GetClient(appName)
	if err != nil {
		return fmt.Errorf("不支持的应用链接:%s", appName)
	}
	_, err = client.NewRequest(ctx).SetQueryParams(req).SetResult(ret).Get(conf.ApiURI + MethodToUri(ctx, method))
	if err != nil {
		return errors.Wrap(err, "call dy api error")
	}
	return nil
}		
```

