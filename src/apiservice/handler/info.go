package handler

import (
	"apiservice/Errs"
	"apiservice/constants"
	"apiservice/service"
	"apiservice/utils"
	"datamodels"

	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Response struct {
	Data       interface{}           `json:"data"`
	Pagination datamodels.Pagination `json:"pagination"`
}

func (p *Provider) InsertNewInfo(rw http.ResponseWriter, r *http.Request) {

	infoIns := &datamodels.UserInfo{}
	if ok := utils.ParseJSON(rw, r.Body, infoIns); !ok {
		return
	}

	infoIns.Id = primitive.NewObjectID()
	//zkey := utils.FormKey(redis_prefix, infoIns.Id.Hex())

	byt, err := json.Marshal(infoIns)
	if err != nil {
		utils.RenderJson(rw, http.StatusInternalServerError, Errs.AppErr{
			Message: "Something went wrong, Please try after sometime!.",
			Err:     err.Error(),
		})
		return
	}

	var sno float64
	if ccmd := p.RCli.Get(constants.RecCounter); ccmd != nil {
		sno, _ = strconv.ParseFloat(ccmd.Val(), 64)
	}

	//Ensuring that the value is originally 0 if so take the length from the
	if sno == 0 {
		if c, _ := p.Db.Collection(constants.UserInfo).CountDocuments(context.Background(), bson.D{}); c != 0 {
			sno = float64(c)
			p.RCli.Set(constants.RecCounter, sno, 0)
		}
	}

	cmd := p.RCli.ZAdd(constants.RedisPrefixKey, &redis.Z{Score: sno + 1, Member: string(byt)})
	if cmd == nil {
		utils.RenderJson(rw, http.StatusInternalServerError, struct {
			Message string `json:"message"`
		}{
			"Error in adding members!.",
		})
		return
	}

	p.RCli.Set(constants.RecCounter, sno+1, 0)
	utils.RenderJson(rw, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		"Member added successfully!.",
	})
	return
}

func (p *Provider) ReadInfo(rw http.ResponseWriter, r *http.Request) {

	var (
		err   error
		pno   int
		limit int
		data  []*datamodels.UserInfo
	)

	pIns := &datamodels.Pagination{}

	if pno, err = strconv.Atoi(r.URL.Query().Get("pageNo")); err != nil && pno == 0 {
		log.Printf("Error: ReadInfo - %v", err)
		utils.RenderJson(rw, http.StatusBadRequest, Errs.AppErr{
			Message: "Page no should be greater that zero!.",
		})
		return
	}

	if limit, err = strconv.Atoi(r.URL.Query().Get("limit")); err != nil && limit == 0 {
		log.Printf("Error: ReadInfo - %v", err)
		utils.RenderJson(rw, http.StatusBadRequest, Errs.AppErr{
			Message: "limit should be greater that zero!.",
		})
		return
	}

	pIns.PageNum = pno
	pIns.Limit = limit

	iSrv := service.NewInfoService(p.RCli, p.Cfg)

	if data, err = iSrv.FetchDataFromCache(pIns); err != nil {
		log.Printf("Error: ReadInfo - %v", err)
		utils.RenderJson(rw, http.StatusInternalServerError, Errs.AppErr{
			Message: "Something went wrong, Please try after sometimes!.",
			Err:     err.Error(),
		})
		return
	}

	utils.RenderJson(rw, http.StatusOK, Response{
		Data:       data,
		Pagination: *pIns,
	})
}
