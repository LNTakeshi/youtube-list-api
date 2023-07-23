package entity

import "time"

type GetList struct {
	Data        []Data
	PrivateInfo PrivateInfo
	Info        Info
}

func (e *GetList) IsRemovable(index int, UUID string) bool {
	if len(e.Data) < index {
		return false
	}
	if e.Data[index].Deleted {
		return false
	}
	if e.IsMaster(UUID) {
		return true
	}
	return e.PrivateInfo.SenderUUIDArray[index] == UUID
}

type Data struct {
	Time     time.Time
	Url      string
	Title    string
	Username string
	Length   int
	Start    int
	End      int
	Deleted  bool
}

type PrivateInfo struct {
	MasterID        string
	CreateRoomDate  time.Time
	SenderUUIDArray []string
	LastUpdateDate  time.Time
	TTL             time.Time
}

type Info struct {
	CurrentIndex int
}

func NewGetList(masterID string) *GetList {
	return &GetList{
		PrivateInfo: PrivateInfo{
			MasterID:       masterID,
			CreateRoomDate: time.Now(),
			LastUpdateDate: time.Now(),
			TTL:            time.Now().Add(time.Hour * 24 * 3),
		},
		Info: Info{CurrentIndex: -1},
	}
}

func (e *GetList) IsMaster(UUID string) bool {
	return e.PrivateInfo.MasterID == UUID
}

func NewData(time time.Time, url string, title string, username string, start int, end int) Data {
	return Data{
		Time:     time,
		Url:      url,
		Title:    title,
		Username: username,
		Length:   end - start,
		Start:    start,
		End:      end,
		Deleted:  false,
	}
}
