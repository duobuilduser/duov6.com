package session

import (
	//"duov6.com/applib"
	"duov6.com/common"
	//"duov6.com/email"
	//"duov6.com/config"
	"duov6.com/objectstore/client"
	"duov6.com/term"
	"encoding/json"
	"strings"
	//"fmt"
)

type AuthCertificate struct {
	UserID, Username, Name, Email, SecurityToken, Domain, DataCaps, ClientIP string
	Otherdata                                                                map[string]string
}

type TenantAutherized struct {
	ID            string
	UserID        string
	TenantID      string
	SecurityLevel string
	Autherized    bool
}

func AddSession(a AuthCertificate) {
	client.Go(a.SecurityToken, "s.duosoftware.auth", "sessions").StoreObject().WithKeyField("SecurityToken").AndStoreOne(a).Ok()
	term.Write("AddSession for "+a.Name+" with SecurityToken :"+a.SecurityToken, term.Debug)
}

func RemoveSession(SecurityToken string) {
	client.Go("ignore", "s.duosoftware.auth", "sessions").DeleteObject().ByUniqueKey(SecurityToken)
	//client.Go("ignore", "s.duosoftware.auth", "sessions").StoreObject().WithKeyField("SecurityToken").AndStoreOne(a).Ok()
	term.Write("LogOut for SecurityToken :"+SecurityToken, term.Debug)
	//return true
}

func AutherizedUser(TenantID, UserID string) (bool, TenantAutherized) {
	bytes, err := client.Go("ignore", "com.duosoftware.tenant", "authorized").GetOne().ByUniqueKey(common.GetHash(UserID + "-" + TenantID)).Ok()
	if err == "" {
		var uList TenantAutherized
		err := json.Unmarshal(bytes, &uList)
		if err == nil {
			return uList.Autherized, uList
		} else {
			return false, TenantAutherized{}
		}
	} else {
		return false, TenantAutherized{}
	}

	bytes1, err1 := client.Go("ignore", "com.duosoftware.tenant", "authorized").GetOne().ByUniqueKey(TenantID).Ok()
	if err1 == "" {
		var uList TenantAutherized
		err := json.Unmarshal(bytes1, &uList)
		if err == nil {
			return uList.Autherized, uList
		} else {
			return false, TenantAutherized{}
		}
	} else {
		return false, TenantAutherized{}
	}
}

func GetSession(key, Domain string) (AuthCertificate, string) {
	bytes, err := client.Go(key, "s.duosoftware.auth", "sessions").GetOne().ByUniqueKey(key).Ok()
	//bytes, err := client.Go(key, "s.duosoftware.auth", "sessions").GetOne().ByUniqueKeyCache(key, 3600).Ok()
	term.Write("GetSession For SecurityToken "+key, term.Debug)
	//term.Write("GetSession For SecurityToken "+string(bytes), term.Debug)

	var c AuthCertificate
	if err == "" {
		if bytes != nil {
			var uList AuthCertificate
			err := json.Unmarshal(bytes, &uList)
			if err == nil {
				if Domain == "Nil" {
					return uList, ""
				} else {

					if strings.ToLower(uList.Domain) != strings.ToLower(Domain) {
						x, _ := AutherizedUser(uList.UserID, Domain)
						if x {
							uList.Domain = strings.ToLower(Domain)
							uList.SecurityToken = common.GetGUID()
							uList.Otherdata = make(map[string]string)
							uList.Otherdata["unused"] = "sss"
							AddSession(uList)
							return uList, ""
						} else {
							return c, Domain + " Session Cound not be Created "
						}
					} else {
						return uList, ""
					}
				}

			} else {
				term.Write("GetSession Error "+err.Error(), term.Error)
			}
		}
	} else {
		term.Write("GetSession Error "+err, term.Error)
	}
	term.Write("GetSession No Session for SecurityToken "+key, term.Debug)

	return c, "Error Session Not Found"
}
