package d_dingtalk

import (
	"context"
	"encoding/json"
	"errors"
	d "github.com/Etpmls/devtool"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var Dingtalk dingtalk

type dingtalk struct {
	AgentId string
	AppKey string
	AppSecret string
	Optional optionalDingtalk
	enable bool
}

type optionalDingtalk struct {
	CorpId string
}

func (this *dingtalk) Init() error {
	if this.AgentId == "" {
		return errors.New("agentId is required")
	}

	if this.AppKey == "" {
		return errors.New("appKey is required")
	}

	if this.AppSecret == "" {
		return errors.New("appSecret is required")
	}

	this.enable = true
	return nil
}

// 获取启动的状态
func (this *dingtalk) GetEnabledStatus() bool {
	return this.enable
}

const (
	baseUri                      = "oapi.dingtalk.com"
	PathSendNotificationsOfWork  = "/topapi/message/corpconversation/asyncsend_v2"
	PathGetUserInfo              = "user/getuserinfo"
	FieldAccessToken             = "DingtalkAccessToken"	// 储存access token的自定义字段的键名
	DefaultFieldValueAccessToken = "DingtalkAccessToken"	// 默认储存access token的自定义字段的值
)

// 获取accessToken
type GetAccessTokenResponse struct {
	AccessToken string	`json:"access_token"`
	ExpiresIn time.Duration	`json:"expires_in"`
	Errmsg string `json:"errmsg"`
	Errcode int `json:"errcode"`
}
func (this *dingtalk) GetAccessToken() (s string, err error) {
	if d.Cache.GetEnabledStatus() {
		return this.getAccessTokenCache()
	} else {
		return this.getAccessTokenNoCache()
	}
}
func (this *dingtalk) getAccessTokenCache() (s string, err error)  {
	k, err := d.GetField(FieldAccessToken)
	if err != nil {
		k = DefaultFieldValueAccessToken
	}
	s, err = d.CacheClient.Get(context.Background(), k).Result()
	if err != nil {
		if err == redis.Nil {
			return this.getAccessTokenNoCache()
		}
		return "", err
	}
	return s, err
}
func (this *dingtalk) getAccessTokenNoCache() (s string, err error)  {
	// 设置参数
	v := url.Values{}
	v.Set("appkey", this.AppKey)
	v.Set("appsecret", this.AppSecret)

	u := url.URL{}
	u.Host = baseUri
	u.Scheme = "https"
	u.Path = "gettoken"
	u.RawQuery = v.Encode()
	// 请求路径
	resp, err := http.Get(u.String())

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// 解析json
	var t GetAccessTokenResponse
	err = json.Unmarshal(body, &t)
	if err != nil {
		return "", err
	}

	// 如果开启了缓存
	if d.Cache.GetEnabledStatus() {
		d.CacheClient.Set(context.Background(), DefaultFieldValueAccessToken, t.AccessToken, time.Second * t.ExpiresIn).Err()
	}

	return t.AccessToken, err
}

//
// 发送工作通知
// https://developers.dingtalk.com/document/app/asynchronous-sending-of-enterprise-session-messages
//
type SendNotificationsOfWorkRequest struct {
	AgentId    string                             `json:"agent_id"`
	UseridList string                             `json:"userid_list"`
	DeptIdList string	`json:"dept_id_list"`
	ToAllUser bool 	`json:"to_all_user"`
	Msg        interface{} 						  `json:"msg"`
}

// 发送文字型工作通知
// https://developers.dingtalk.com/document/app/message-types-and-data-format/title-dfs-oxn-29n
type TextMessage struct {
	Msgtype string	`json:"msgtype"`
	Text struct{
		Content string	`json:"content"`
	}	`json:"text"`
}
func (this *dingtalk) SendTextNotificationsOfWork(userID string, content string) (err error){
	// 获取access token
	token, err := this.GetAccessToken()
	if err != nil {
		return err
	}

	// 设置发送内容
	var m TextMessage
	m.Msgtype = "text"
	m.Text.Content = content

	// 设置Body参数
	ctx := SendNotificationsOfWorkRequest{
		AgentId:    this.AgentId,
		UseridList: userID,
		Msg:        m,
	}

	v := url.Values{}
	v.Set("access_token", token)	// Query参数

	// 解析URL到结构体
	u := url.URL{}
	u.Host = baseUri
	u.Scheme = "https"
	// u.PathGetUserInfo = "/topapi/TextMessage/corpconversation/asyncsend_v2"
	u.Path = PathSendNotificationsOfWork
	u.RawQuery = v.Encode()

	tmp, err := json.Marshal(ctx)
	if err != nil {
		return err
	}
	resp, err := http.Post(u.String(), "application/json", strings.NewReader(string(tmp)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return err
}

// 发送OA型工作通知
type OaMessage struct {
	Msgtype string	`json:"msgtype"`
	Oa struct{
		MessageUrl	string	`json:"message_url"`
		PcMessageUrl string	`json:"pc_message_url"`
		Head struct{
			Bgcolor string	`json:"bgcolor"`	// 头部标题颜色
			Text	string	`json:"text"`	// 这个字段实际没什么用，没有显示
		}	`json:"head"`
		StatusBar struct{
			StatusValue	string `json:"status_value"`
			StatusBg string `json:"status_bg"`
		}	`json:"status_bar,omitempty"`
		Body struct{
			Title	string `json:"title"`	// 正文标题
			Form []struct{
				Key	string	`json:"key"`
				Value string	`json:"value"`
			}	`json:"form,omitempty"`
			Rich	struct{
				Num string	`json:"num"`
				Unit string `json:"unit"`
			}	`json:"rich,omitempty"`
			Content string	`json:"content"`
			Image	string	`json:"image"`
			FileCount string 	`json:"file_count"`
			Author	string	`json:"author"`
		}	`json:"body"`
	}	`json:"oa"`
}

// SendOaNotificationsOfWork Example:
// // 设置字段内容
// // 跳转钉钉工作台需要的URL前缀
// jumpWorkbenchUrlPrefix := "dingtalk://dingtalkclient/action/openapp?corpid=" + this.CorpId + "&container_type=work_platform&app_id=0_" + this.AgentId + "&redirect_type=jump&redirect_url="
// 	var m OaMessage
//	m.Msgtype = "oa"
//	m.Oa.MessageUrl = appUrl
//	m.Oa.PcMessageUrl = jumpWorkbenchUrlPrefix + appUrl
//	m.Oa.Head.Bgcolor = "FFFF0000"
//	m.Oa.Head.Text = "xxx通知"	// 这个字段实际没什么用，没有显示
//	m.Oa.Body.Title = "xxx通知"
//	m.Oa.Body.Content = content
//	m.Oa.Body.Image = dingtalk_imageid
//	m.Oa.Body.Author = "Author Name"
//	for k, v := range form {
//		m.Oa.Body.Form = append(m.Oa.Body.Form, message_OA_form{
//			Key:   k,
//			Value: v,
//		})
//	}
func (this *dingtalk) SendOaNotificationsOfWork(userID string, m OaMessage) (err error) {
	// 获取access token
	token, err := this.GetAccessToken()
	if err != nil {
		return err
	}

	// 设置Body参数
	ctx := SendNotificationsOfWorkRequest{
		AgentId:    this.AgentId,
		UseridList: userID,
		Msg:        m,
	}

	v := url.Values{}
	v.Set("access_token", token)

	// 解析URL到结构体
	u := url.URL{}
	u.Host = baseUri
	u.Scheme = "https"
	// u.PathGetUserInfo = "/topapi/TextMessage/corpconversation/asyncsend_v2"
	u.Path = PathSendNotificationsOfWork
	u.RawQuery = v.Encode()

	tmp, err := json.Marshal(ctx)
	if err != nil {
		return err
	}
	resp, err := http.Post(u.String(), "application/json", strings.NewReader(string(tmp)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return err
}

// 免密登录————获取用户信息
// https://developers.dingtalk.com/document/app/get-user-userid-through-login-free-code
type GetUserInfoResponse struct {
	Errcode int `json:"errcode"`
	Errmsg string `json:"errmsg"`
	UserId string	`json:"userid"`
	Name string `json:"name"`
	DeviceId string `json:"deviceId"`
	IsSys bool `json:"is_sys"`
	SysLevel int	`json:"sys_level"`
}
func (this *dingtalk) GetUserInfo(access_token string, code string) (userInfo GetUserInfoResponse, err error) {
	// 设置参数
	v := url.Values{}
	v.Set("access_token", access_token)	//调用服务端API的应用凭证
	v.Set("code", code)	// 免登授权码

	u := url.URL{}
	u.Host = baseUri
	u.Scheme = "https"
	u.Path = PathGetUserInfo
	u.RawQuery = v.Encode()

	// 请求路径
	resp, err := http.Get(u.String())

	if err != nil {
		return GetUserInfoResponse{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// 解析json
	var jsn GetUserInfoResponse
	err = json.Unmarshal(body, &jsn)
	if err != nil {
		return GetUserInfoResponse{}, err
	}

	return jsn, nil
}

// 免密登录————获取用户Id
func (this *dingtalk) GetUserId(access_token string, code string) (userId string, err error) {
	userInfo, err := this.GetUserInfo(access_token, code)
	if err != nil {
		return "", err
	}
	return userInfo.UserId, nil
}

// 设置字段内容
// 跳转钉钉工作台需要的URL前缀
func (this *dingtalk) GenerateJumpableUrl(appurl string) (string, error) {
	if this.Optional.CorpId == "" {
		return "", errors.New("corpId of dingtalk is not set")
	}

	u := "dingtalk://dingtalkclient/action/openapp?corpid=" + this.Optional.CorpId + "&container_type=work_platform&app_id=0_" + this.AgentId + "&redirect_type=jump&redirect_url=" + appurl

	return u, nil
}
