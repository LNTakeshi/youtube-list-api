package status

import (
	"net/http"
	"strings"
	"time"
	"youtubelist/domain/entity"
)

type GetList struct {
	Data        []Data      `json:"data"`
	PrivateInfo PrivateInfo `json:"privateInfo"`
	Info        Info        `json:"info"`
}

type Data struct {
	Time      string `json:"time"`
	Url       string `json:"url"`
	Title     string `json:"title"`
	Username  string `json:"username"`
	Length    string `json:"length"`
	Start     int    `json:"start"`
	End       int    `json:"end"`
	Deleted   bool   `json:"deleted"`
	Removable bool   `json:"removable"`
}

type PrivateInfo struct {
	MasterID *string `json:"masterId"`
	UUID     string  `json:"uuid"`
}

type Info struct {
	CurrentIndex int  `json:"currentIndex"`
	NeedUpdate   bool `json:"needUpdate"`
}

func NewGetList(r *http.Request, e *entity.GetList, uuid string, lastUpdateDate time.Time) GetList {
	if len(e.Data) > 0 && e.PrivateInfo.LastUpdateDate.Before(lastUpdateDate) {
		return GetList{Data: []Data{}, Info: NewInfo(e.Info, false)}
	}
	return GetList{
		Data:        NewData(r, e, uuid),
		PrivateInfo: NewPrivateInfo(e, uuid),
		Info:        NewInfo(e.Info, true),
	}
}

func NewData(r *http.Request, getList *entity.GetList, uuid string) []Data {
	res := make([]Data, 0)
	loc, _ := time.LoadLocation("Asia/Tokyo")
	for _, v := range getList.Data {
		if strings.HasPrefix(v.Title, "[HIDDEN]") && !strings.Contains(r.Header.Get("User-Agent"), "Unity") {
			v.Url = "https://www.youtube.com/"
		}

		res = append(res, Data{
			Time:      v.Time.In(loc).Format("2006-01-02 15:04:05"),
			Url:       v.Url,
			Title:     v.Title,
			Username:  v.Username,
			Length:    time.Time{}.Add(time.Duration(v.Length) * time.Second).Format("15:04:05"),
			Start:     v.Start,
			End:       v.End,
			Deleted:   v.Deleted,
			Removable: getList.IsRemovable(len(res), uuid),
		})
	}
	return res
}

func NewPrivateInfo(getList *entity.GetList, uuid string) PrivateInfo {
	st := PrivateInfo{UUID: uuid}
	if !getList.IsMaster(uuid) {
		return st
	}
	st.MasterID = &getList.PrivateInfo.MasterID
	return st
}

func NewInfo(e entity.Info, needUpdate bool) Info {
	return Info{CurrentIndex: e.CurrentIndex, NeedUpdate: needUpdate}
}
