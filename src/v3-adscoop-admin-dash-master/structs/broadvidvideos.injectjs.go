package structs

type InjectJss []InjectJs

func (i *InjectJss) FindAll() error {
	return BroadvidDB.Table("inject_js").Find(&i).Error
}

type InjectJs struct {
	ID            int64  `form:"id"`
	Code          string `form:"code"`
	DefaultLander int64  `form:"default_lander"`
	Name          string `form:"name"`
}

func (v InjectJs) TableName() string {
	return "inject_js"
}

func (i *InjectJs) Find(id string) error {
	return BroadvidDB.Find(&i, id).Error
}

func (i *InjectJs) Save() error {
	return BroadvidDB.Save(&i).Error
}

type InjectJsOptions struct {
	ID          int64  `form:"id"`
	Key         string `form:"key"`
	Value       string `form:"value"`
	Lander      int64  `form:"lander"`
	LanderLabel string `form:"-" sql:"-"`
	InjectJsID  int64  `form:"inject_js_id"`
}

func (v InjectJsOptions) TableName() string {
	return "inject_js_options"
}
