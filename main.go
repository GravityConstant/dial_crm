package main

import (
	_ "zq/callout_crm/routers"
	_ "zq/callout_crm/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
