package http

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/gateway/permission"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	scas "github.com/qiangmzsx/string-adapter"
)

//NewPermissioner creates a gin Handle Func to controll permission
//currently there is only 2 permission for POST/GET requests
func NewPermissioner(readKeyID, writeKeyID, confirmKeyID, rebalanceKeyID string) (gin.HandlerFunc, error) {
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
m = g(r.sub, p.sub)  && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
`
	)

	pol := fmt.Sprintf(`
p, %[1]s, /*, GET

p, %[2]s, /*, GET
p, %[2]s, /create-asset, POST
p, %[2]s, /update-asset, POST 
p, %[2]s, /create-asset-exchange, POST
p, %[2]s, /update-exchange, POST
p, %[2]s, /create-trading-pair, POST
p, %[2]s, /update-trading-pair, POST

p, %[3]s, /*, GET
p, %[3]s, /create-asset/:id, (PUT)|(DELETE)
p, %[3]s, /update-asset/:id, (PUT)|(DELETE)
p, %[3]s, /create-asset-exchange/:id, (PUT)|(DELETE)
p, %[3]s, /update-exchange/:id, (PUT)|(DELETE)
p, %[3]s, /create-trading-pair/:id, (PUT)|(DELETE) 
p, %[3]s, /update-trading-pair/:id, (PUT)|(DELETE) 

p, %[4]s, /*, GET
`, readKeyID, writeKeyID, confirmKeyID, rebalanceKeyID)
	sa := scas.NewAdapter(pol)
	e := casbin.NewEnforcer(casbin.NewModel(conf), sa)
	if err := e.LoadPolicy(); err != nil {
		return nil, err
	}

	p := permission.NewPermissioner(e)
	return p, nil
}
