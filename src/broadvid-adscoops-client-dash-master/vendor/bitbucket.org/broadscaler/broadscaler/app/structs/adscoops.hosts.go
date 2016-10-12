package structs

type Host struct {
	ID   int64
	Host string
}

func (h Host) TableName() string {
	return "adscoop_hosts"
}

type Hosts []Host

func (h *Hosts) FindAll() error {
	return AdscoopsDB.Table("adscoop_hosts").Find(&h).Error
}
