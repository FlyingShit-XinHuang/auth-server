package auth

import (
	"github.com/RangelReale/osin"

	"bytes"
	"fmt"
	"net/http"

	"whispir/auth-server/storage"
)

type Service interface {
	GetAccessToken(req *http.Request) *osin.Response
	Info(req *http.Request) *osin.Response
	GetAuthCode(req *http.Request) (*osin.Response, []byte)
}

func NewBasicService(storage storage.OAuth2Storage) Service {
	cfg := osin.NewServerConfig()
	cfg.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE, osin.PASSWORD, osin.CLIENT_CREDENTIALS}
	return &authService{
		osin.NewServer(cfg, storage),
		storage,
	}
}

type authService struct{
	osinServer *osin.Server
	storage storage.OAuth2Storage
}

func (a *authService) GetAccessToken(req *http.Request) (resp *osin.Response) {
	resp = a.osinServer.NewResponse()
	defer func() {
		if resp.IsError && resp.InternalError != nil {
			fmt.Printf("ERROR: %s\n", resp.InternalError)
		}
	}()

	if ar := a.osinServer.HandleAccessRequest(resp, req); ar != nil {
		if osin.PASSWORD == ar.Type {
			ok, err := a.validateUser(ar.Username, ar.Password)
			if nil != err {
				resp.SetError(osin.E_SERVER_ERROR, "")
				return
			}
			if !ok {
				resp.SetError(osin.E_INVALID_GRANT, "invalid username or password")
				return
			}
		}
		ar.GenerateRefresh = false
		ar.Authorized = true
		a.osinServer.FinishAccessRequest(resp, req, ar)
	}
	return
}

func (a *authService) Info(req *http.Request) *osin.Response {
	resp := a.osinServer.NewResponse()
	if ir := a.osinServer.HandleInfoRequest(resp, req); ir != nil {
		a.osinServer.FinishInfoRequest(resp, req, ir)
	}
	if resp.IsError && resp.InternalError != nil {
		fmt.Printf("ERROR: %s\n", resp.InternalError)
	}
	return resp
}

func (a *authService) GetAuthCode(req *http.Request) (*osin.Response, []byte) {
	resp := a.osinServer.NewResponse()

	if ar := a.osinServer.HandleAuthorizeRequest(resp, req); nil != ar {
		if http.MethodGet == req.Method {
			return nil, authPage(req.URL.RawQuery, ar.Client.GetUserData().(string))
		}
		req.ParseForm()
		ok, err := a.validateUser(req.Form.Get("user"), req.Form.Get("password"))
		if nil != err {
			resp.SetError(osin.E_SERVER_ERROR, "")
		} else if !ok {
			resp.SetError(osin.E_ACCESS_DENIED, "invalid username or password")
		} else {
			ar.Authorized = true
			a.osinServer.FinishAuthorizeRequest(resp, req, ar)
		}
	}
	if resp.IsError && resp.InternalError != nil {
		fmt.Printf("ERROR: %s\n", resp.InternalError)
	}
	return resp, nil
}

func authPage(query, name string) []byte {
	buf := bytes.NewBuffer([]byte("<html><body>"))
	fmt.Fprintf(buf, "Application %s want to query your resources, please login to authorize", name)
	fmt.Fprintf(buf, "<form action=\"%s?%s\" method=\"POST\">", AuthPath, query)
	buf.Write([]byte("User: <input type=\"text\" name=\"user\" /><br/>"))
	buf.Write([]byte("Password: <input type=\"password\" name=\"password\" /><br/>"))
	buf.Write([]byte("<input type=\"submit\" value=\"Allow\"/></input>"))
	buf.Write([]byte("</form></body></html>"))
	return buf.Bytes()
}

func (a *authService) validateUser(name, password string) (bool, error) {
	user, err := a.storage.GetUserByNameAndPassword(name, password)
	if nil != err {
		return false, err
	}
	return nil != user, nil
}
