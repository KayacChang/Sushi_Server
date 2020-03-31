package gameattach

import "github.com/YWJSonic/ServerUtility/attach"

// UserAttach ...
type UserAttach struct {
	userID  int64
	dataMap map[int64]map[int64]*attach.Info
}

// LoadData ...
func (ga *UserAttach) LoadData() {
	// redis load data

	// if fail sql load data
}

// Get ...
func (ga *UserAttach) Get(attachkind int64, attachtype int64) *attach.Info {
	if _, ok := (*ga.GetType(attachkind))[attachtype]; !ok {
		ga.SetValue(attachkind, attachtype, "", 0)
	}
	return ga.dataMap[attachkind][attachtype]
}

// GetType ...
func (ga *UserAttach) GetType(attachkind int64) *map[int64]*attach.Info {
	if _, ok := ga.dataMap[attachkind]; !ok {
		ga.dataMap[attachkind] = make(map[int64]*attach.Info)
	}
	result := ga.dataMap[attachkind]
	return &result
}

// SetDBValue ...
func (ga *UserAttach) SetDBValue(attachKind, attachType int64, SValue string, IValue int64) {

	if att, ok := (*ga.GetType(attachKind))[attachType]; !ok {
		att = attach.NewInfo(attachKind, attachType, true)
		att.SetSValue(SValue)
		att.SetIValue(IValue)
		ga.dataMap[attachKind][attachType] = att
	} else {
		att.SetSValue(SValue)
		att.SetIValue(IValue)
	}
}

// SetValue ...
func (ga *UserAttach) SetValue(attachKind, attachType int64, SValue string, IValue int64) {

	if att, ok := (*ga.GetType(attachKind))[attachType]; !ok {
		att = attach.NewInfo(attachKind, attachType, false)
		att.SetSValue(SValue)
		att.SetIValue(IValue)
		ga.dataMap[attachKind][attachType] = att
	} else {
		att.SetSValue(SValue)
		att.SetIValue(IValue)
	}
}

// SetAttach ...
func (ga *UserAttach) SetAttach(info *attach.Info) {
	ga.dataMap[info.GetKind()][info.GetTypes()] = info
}

// Save ...
func (ga *UserAttach) Save() {

}

// InitData(userID int64)
// Get(attachkind string, attachtype string) Info
// GetType(attachkind string) map[string]Info
// SetValue(attachKind, attachType, SValue string, IValue int64)
// SetAttach(attach Info)
// Save()
