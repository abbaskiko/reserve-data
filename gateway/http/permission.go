package http

import (
	"fmt"

	"github.com/KyberNetwork/reserve-data/gateway/permission"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	scas "github.com/qiangmzsx/string-adapter"

	"github.com/KyberNetwork/httpsign-utils/authenticator"
)

func addKeyReadPolicy(key string) string {
	return fmt.Sprintf(`
p, %[1]s, /*, GET`, key)
}

func addKeyWritePolicy(key string) string {
	return fmt.Sprintf(`
p, %[1]s, /*, GET
p, %[1]s, /v3/create-asset, POST
p, %[1]s, /v3/update-asset, POST 
p, %[1]s, /v3/create-asset-exchange, POST
p, %[1]s, /v3/update-exchange, POST
p, %[1]s, /v3/create-trading-pair, POST
p, %[1]s, /v3/update-trading-pair, POST
p, %[1]s, /v3/setting-change, POST
p, %[1]s, /v3/change-asset-address, POST`, key)
}

func addKeyConfirmPolicy(key string) string {
	return fmt.Sprintf(`
p, %[1]s, /*, GET
p, %[1]s, /v3/create-asset/:id, (PUT)|(DELETE)
p, %[1]s, /v3/update-asset/:id, (PUT)|(DELETE)
p, %[1]s, /v3/create-asset-exchange/:id, (PUT)|(DELETE)
p, %[1]s, /v3/update-exchange/:id, (PUT)|(DELETE)
p, %[1]s, /v3/create-trading-pair/:id, (PUT)|(DELETE) 
p, %[1]s, /v3/update-trading-pair/:id, (PUT)|(DELETE)
p, %[1]s, /v3/setting-change/:id, (PUT)|(DELETE)
p, %[1]s, /v3/change-asset-address, (PUT)|(DELETE)`, key)
}

func addKeyRebalancePolicy(key string) string {
	return fmt.Sprintf(`
p, %[1]s, /*, GET
p, %[1]s, /price-factor, POST,
p, %[1]s, /cancelorder/:exchangeid, POST
p, %[1]s, /deposit/:exchangeid, POST
p, %[1]s, /withdraw/:exchangeid, POST
p, %[1]s, /trade, POST
p, %[1]s, /setrates, POST
p, %[1]s, /holdrebalance, POST
p, %[1]s, /enablesetrate, POST
p, %[1]s, /holdsetrate, POST
p, %[1]s, /enablesetrate, POST`, key)
}

//NewPermissioner creates a gin Handle Func to controll permission
//currently there is only 2 permission for POST/GET requests
func NewPermissioner(readKeys, writeKeys, confirmKeys, rebalanceKeys []authenticator.KeyPair) (gin.HandlerFunc, error) {
	const (
		conf = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _ , _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub)  && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
`
	)
	var pol string
	for _, key := range readKeys {
		pol += addKeyReadPolicy(key.AccessKeyID)
	}
	for _, key := range writeKeys {
		pol += addKeyWritePolicy(key.AccessKeyID)
	}
	for _, key := range confirmKeys {
		pol += addKeyConfirmPolicy(key.AccessKeyID)
	}
	for _, key := range rebalanceKeys {
		pol += addKeyRebalancePolicy(key.AccessKeyID)
	}

	sa := scas.NewAdapter(pol)
	e := casbin.NewEnforcer(casbin.NewModel(conf), sa)
	if err := e.LoadPolicy(); err != nil {
		return nil, err
	}

	p := permission.NewPermissioner(e)
	return p, nil
}
