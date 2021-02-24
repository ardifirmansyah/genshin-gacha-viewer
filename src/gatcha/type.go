package gatcha

import "net/url"

const (
	noviceWish    = "100"
	permanentWish = "200"
	charEventWish = "301"
	weaponWish    = "302"

	mihoyoHost   = "https://hk4e-api-os.mihoyo.com"
	gachaLogPath = "/event/gacha_info/api/getGachaLog"

	sizePerPage = "6"
)

var wishes = []string{
	noviceWish,
	permanentWish,
	charEventWish,
	weaponWish,
}

type (
	ProcessRequest struct {
		LogURL string `json:"logURL"`

		AuthKeyVersion string
		SignType       string
		AuthAppID      string
		InitType       string
		Lang           string
		DeviceType     string
		Ext            string
		GameVersion    string
		Region         string
		AuthKey        string
		GameBiz        string
	}

	WishResult struct {
		UserID    string `json:"uid"`
		GachaType string `json:"gacha_type"`
		ItemID    string `json:"item_id"`
		Count     string `json:"count"`
		Time      string `json:"time"`
		Name      string `json:"name"`
		Lang      string `json:"lang"`
		ItemType  string `json:"item_type"`
		RankType  string `json:"rank_type"`
	}

	GachaLogResponse struct {
		ReturnCode int          `json:"retcode"`
		Message    string       `json:"message"`
		Data       GachaLogData `json:"data"`
	}

	GachaLogData struct {
		Page   string       `json:"page"`
		Size   string       `json:"size"`
		Total  string       `json:"total"`
		List   []WishResult `json:"list"`
		Region string       `json:"region"`
	}
)

func (r *ProcessRequest) bindURLValues(values url.Values) {
	r.AuthKeyVersion = values.Get("authkey_ver")
	r.SignType = values.Get("sign_type")
	r.AuthAppID = values.Get("auth_appid")
	r.InitType = values.Get("init_type")
	r.Lang = values.Get("lang")
	r.DeviceType = values.Get("device_type")
	r.Ext = values.Get("ext")
	r.GameVersion = values.Get("game_version")
	r.Region = values.Get("region")
	r.AuthKey = values.Get("authkey")
	r.GameBiz = values.Get("game_biz")
}

func (r *ProcessRequest) toURLValues() url.Values {
	return map[string][]string{
		"authkey_ver":  {r.AuthKeyVersion},
		"sign_type":    {r.SignType},
		"auth_appid":   {r.AuthAppID},
		"init_type":    {r.InitType},
		"lang":         {r.Lang},
		"device_type":  {r.DeviceType},
		"ext":          {r.Ext},
		"game_version": {r.GameVersion},
		"region":       {r.Region},
		"authkey":      {r.AuthKey},
		"game_biz":     {r.GameBiz},
	}
}
