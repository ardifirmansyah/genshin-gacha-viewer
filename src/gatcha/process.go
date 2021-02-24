package gatcha

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/render"

	"github.com/ardifirmansyah/genshin-gacha-viewer/src/common/api"
	"github.com/ardifirmansyah/genshin-gacha-viewer/src/common/curl"
)

func Process(w http.ResponseWriter, r *http.Request) {
	var (
		err error

		req  = ProcessRequest{}
		resp = api.NewResponse()
	)

	defer func() {
		if err != nil {
			log.Println(err.Error())
			resp.AddError(err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(resp.JSON())
	}()

	err = bindRequest(r, &req)
	if err != nil {
		return
	}

	data, err := getGachaInfo(req)
	if err != nil {
		return
	}

	resp.Success = true
	resp.Data = data
}

func bindRequest(r *http.Request, req *ProcessRequest) error {
	err := render.Decode(r, req)
	if err != nil {
		return err
	}

	b, err := base64.URLEncoding.DecodeString(req.LogURL)
	if err != nil {
		return err
	}
	req.LogURL = string(b)

	url, err := url.ParseRequestURI(req.LogURL)
	if err != nil {
		return err
	}
	req.bindURLValues(url.Query())

	return nil
}

func getGachaInfo(req ProcessRequest) (map[string][]WishResult, error) {
	mapWishes := make(map[string][]WishResult)

	for _, wish := range wishes {
		page := 1

		for {
			resp, err := fetchGachaLog(req, wish, page)
			if err != nil {
				break
			}

			if len(resp.Data.List) == 0 {
				break
			}

			for _, result := range resp.Data.List {
				mapWishes[result.GachaType] = append(mapWishes[result.GachaType], result)
			}
			page++
		}
	}

	return mapWishes, nil
}

func fetchGachaLog(r ProcessRequest, wishType string, page int) (GachaLogResponse, error) {
	resp := GachaLogResponse{}

	requestParam := r.toURLValues()

	requestParam.Add("gacha_type", wishType)
	requestParam.Add("page", strconv.Itoa(page))
	requestParam.Add("size", sizePerPage)

	req := curl.NewRequest()
	req.URL = mihoyoHost + gachaLogPath
	req.Method = http.MethodGet
	req.Params = requestParam

	_, body, err := req.DoRequest()
	if err != nil {
		return GachaLogResponse{}, err
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return GachaLogResponse{}, err
	}

	return resp, nil
}
