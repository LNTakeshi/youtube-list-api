package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"youtubelist/domain/entity"
	"youtubelist/errors"

	"cloud.google.com/go/firestore"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/morikuni/failure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	roomID string
	uUID   string
	fsCli  *firestore.Client
	redis  *redis.Client
}

func NewService(roomID string, uUID string, fsCli *firestore.Client, redis *redis.Client) *Service {
	if uUID == "" {
		uUID = uuid.NewString()
	}
	return &Service{roomID: roomID, uUID: uUID, fsCli: fsCli, redis: redis}
}

func (s *Service) GetRoom(ctx context.Context) (*entity.GetList, string, error) {
	var getList *entity.GetList
	b, err := s.redis.Get(ctx, s.roomID).Bytes()
	if err == nil {
		err = json.Unmarshal(b, &getList)
		if err != nil {
			failure.Wrap(err)
		}
		return getList, s.uUID, nil
	}
	if err != nil {
		doc, err := s.fsCli.Collection("Room").Doc(s.roomID).Get(ctx)
		if status.Code(err) == codes.NotFound {
			if err = s.fsCli.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
				doc, txErr := t.Get(s.fsCli.Collection("Room").Doc(s.roomID))
				if status.Code(txErr) == codes.NotFound {
					getList, txErr = s.createRoom(ctx, t)
					if txErr != nil {
						return txErr
					}
				} else {
					txErr = mapstructure.Decode(doc.Data(), &getList)
					if txErr != nil {
						return failure.Wrap(txErr)
					}
				}

				if getList == nil {
					getList, txErr = s.createRoom(ctx, t)
					if txErr != nil {
						return txErr
					}
				}
				return nil
			}); err != nil {
				return nil, "", failure.Wrap(err)
			}
		} else {
			err = mapstructure.Decode(doc.Data(), &getList)
			if err != nil {
				return nil, "", failure.Wrap(err)
			}
		}
		b, err := json.Marshal(getList)
		if err != nil {
			return nil, "", failure.Wrap(err)
		}
		s.redis.Set(ctx, s.roomID, b, time.Hour*24)
	}

	return getList, s.uUID, nil
}

func (s *Service) createRoom(ctx context.Context, t *firestore.Transaction) (*entity.GetList, error) {
	getList := entity.NewGetList(s.uUID)
	if err := t.Create(s.fsCli.Collection("Room").Doc(s.roomID), getList); err != nil {
		return nil, failure.Wrap(err)
	}
	return getList, nil
}

func (s *Service) Add(ctx context.Context, fetchResult *FetchResult, username string, startStr string, endStr string) error {
	startTime, err := time.Parse("15:04:05", startStr)
	if err != nil {
		return failure.New(errors.ErrInvalidTime)
	}
	endTime, err := time.Parse("15:04:05", endStr)
	if err != nil {
		return failure.New(errors.ErrInvalidTime)
	}
	println(startTime.String())
	println(endTime.String())
	println(fetchResult.Length)

	zeroTime, err := time.Parse("15:04:05", "00:00:00")
	if err != nil {
		return failure.New(errors.ErrInvalidTime)
	}
	println(zeroTime.String())
	if endTime.Equal(zeroTime) {
		endTime = zeroTime.Add(time.Second * time.Duration(fetchResult.Length))
	}
	println(endTime.String())
	if startTime.After(endTime) {
		return failure.New(errors.ErrInvalidTime)
	}
	start := startTime.Second() + startTime.Minute()*60 + startTime.Hour()*60*60
	end := endTime.Second() + endTime.Minute()*60 + endTime.Hour()*60*60
	if end == 0 {
		end = fetchResult.Length
	}
	if end > fetchResult.Length {
		return failure.New(errors.ErrInvalidTime)
	}
	println(start)
	println(end)

	e := entity.NewData(time.Now(), fetchResult.Url, fetchResult.Title, username, start, end)
	if err := s.fsCli.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
		doc, txErr := t.Get(s.fsCli.Collection("Room").Doc(s.roomID))
		var getList *entity.GetList
		if status.Code(txErr) == codes.NotFound {
			getList, txErr = s.createRoom(ctx, t)
			if txErr != nil {
				return txErr
			}
		}
		txErr = mapstructure.Decode(doc.Data(), &getList)
		if txErr != nil {
			return failure.Wrap(txErr)
		}

		getList.Data = append(getList.Data, e)
		getList.PrivateInfo.SenderUUIDArray = append(getList.PrivateInfo.SenderUUIDArray, s.uUID)
		getList.PrivateInfo.LastUpdateDate = time.Now()

		txErr = t.Set(s.fsCli.Collection("Room").Doc(s.roomID), getList)
		if txErr != nil {
			return failure.Wrap(txErr)
		}
		s.redis.Del(ctx, s.roomID)
		return nil
	}); err != nil {
		return failure.Wrap(err)
	}
	return nil
}

func (s *Service) Remove(ctx context.Context, indexStr string) error {
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return failure.Wrap(err)
	}
	if err = s.fsCli.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
		doc, txErr := t.Get(s.fsCli.Collection("Room").Doc(s.roomID))
		if txErr != nil {
			return failure.Wrap(txErr)
		}
		var getList *entity.GetList
		txErr = mapstructure.Decode(doc.Data(), &getList)
		if txErr != nil {
			return failure.Wrap(txErr)
		}

		if len(getList.Data) < index || (getList.PrivateInfo.SenderUUIDArray[index] != s.uUID && !getList.IsMaster(s.uUID)) {
			return failure.New(errors.ErrBadRequest)
		}
		getList.Data[index].Deleted = true
		getList.PrivateInfo.LastUpdateDate = time.Now()

		txErr = t.Set(s.fsCli.Collection("Room").Doc(s.roomID), getList)
		if txErr != nil {
			return failure.Wrap(txErr)
		}
		s.redis.Del(ctx, s.roomID)
		return nil
	}); err != nil {
		return failure.Wrap(err)
	}
	return nil
}

func (s *Service) SetCurrentIndex(ctx context.Context, indexStr string) error {
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return failure.Wrap(err)
	}
	if err = s.fsCli.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
		doc, txErr := t.Get(s.fsCli.Collection("Room").Doc(s.roomID))
		if txErr != nil {
			return failure.Wrap(txErr)
		}
		var getList *entity.GetList
		txErr = mapstructure.Decode(doc.Data(), &getList)
		if txErr != nil {
			return failure.Wrap(txErr)
		}

		if len(getList.Data) < index || !getList.IsMaster(s.uUID) {
			return failure.New(errors.ErrBadRequest)
		}
		getList.Info.CurrentIndex = index
		getList.PrivateInfo.LastUpdateDate = time.Now()
		txErr = t.Set(s.fsCli.Collection("Room").Doc(s.roomID), getList)
		if txErr != nil {
			return failure.Wrap(txErr)
		}
		s.redis.Del(ctx, s.roomID)
		return nil
	}); err != nil {
		return failure.Wrap(err)
	}
	return nil
}
