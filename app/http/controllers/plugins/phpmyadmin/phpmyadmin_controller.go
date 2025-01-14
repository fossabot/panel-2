package phpmyadmin

import (
	"os"
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"panel/app/http/controllers"
	"panel/pkg/tools"
)

type PhpMyAdminController struct {
}

func NewPhpMyAdminController() *PhpMyAdminController {
	return &PhpMyAdminController{}
}

func (c *PhpMyAdminController) Info(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "phpmyadmin")
	if check != nil {
		return check
	}

	files, err := os.ReadDir("/www/server/phpmyadmin")
	if err != nil {
		return controllers.Error(ctx, http.StatusBadRequest, "找不到 phpMyAdmin 目录")
	}

	var phpmyadmin string
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "phpmyadmin_") {
			phpmyadmin = f.Name()
		}
	}
	if len(phpmyadmin) == 0 {
		return controllers.Error(ctx, http.StatusBadRequest, "找不到 phpMyAdmin 目录")
	}

	conf := tools.Read("/www/server/vhost/phpmyadmin.conf")
	match := regexp.MustCompile(`listen\s+(\d+);`).FindStringSubmatch(conf)
	if len(match) == 0 {
		return controllers.Error(ctx, http.StatusBadRequest, "找不到 phpMyAdmin 端口")
	}

	return controllers.Success(ctx, http.Json{
		"phpmyadmin": phpmyadmin,
		"port":       match[1],
	})
}

func (c *PhpMyAdminController) SetPort(ctx http.Context) http.Response {
	check := controllers.Check(ctx, "phpmyadmin")
	if check != nil {
		return check
	}

	port := ctx.Request().Input("port")
	if len(port) == 0 {
		return controllers.Error(ctx, http.StatusBadRequest, "端口不能为空")
	}

	conf := tools.Read("/www/server/vhost/phpmyadmin.conf")
	conf = regexp.MustCompile(`listen\s+(\d+);`).ReplaceAllString(conf, "listen "+port+";")
	tools.Write("/www/server/vhost/phpmyadmin.conf", conf, 0644)

	if tools.IsRHEL() {
		tools.Exec("firewall-cmd --zone=public --add-port=" + port + "/tcp --permanent")
		tools.Exec("firewall-cmd --reload")
	} else {
		tools.Exec("ufw allow " + port + "/tcp")
		tools.Exec("ufw reload")
	}
	tools.Exec("systemctl reload openresty")

	return controllers.Success(ctx, nil)
}
