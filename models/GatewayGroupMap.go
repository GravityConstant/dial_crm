package models

type GatewayGroupMap struct {
    Id           int
    Gateway      *Gateway      `orm:"rel(fk)"`
    GatewayGroup *GatewayGroup `orm:"rel(fk)"`
}

func (a *GatewayGroupMap) TableName() string {
    return GatewayGroupMapTBName()
}