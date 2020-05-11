package handler

import (
	"apiservice/Errs"
	"apiservice/constants"
	"apiservice/service"
	"apiservice/utils"
	"datamodels"

	"log"
	"net/http"
	"strconv"

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

	if err := service.NewInfoService(p.RCli, p.Db.Collection(constants.COLUserInfo)).AddDataToCache(infoIns); err != nil {
		utils.RenderJson(rw, http.StatusOK, struct {
			Message string `json:"message"`
		}{
			"Something went wrong, Please try after sometime!.",
		})
		return
	}

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
		data  interface{}
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

	iSrv := service.NewInfoService(p.RCli, p.Db.Collection(constants.COLUserInfo))

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
