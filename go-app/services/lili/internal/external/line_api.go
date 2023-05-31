package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/80andCo/LiLi-LABO/services/lili/internal/config"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/constant"
	"github.com/80andCo/LiLi-LABO/services/lili/internal/util"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/pkg/errors"
)

type (
	LineAPI interface {
		FetchAccessToken(ctx context.Context, authCode string) (token string, err error)
		FetchLineProfileByAccessToken(ctx context.Context, accessToken string) (*LineProfile, error)
	}

	lineAPI struct {
		cl  *linebot.Client
		cfg config.Config
	}

	LineProfile struct {
		LineID        string `json:"line_id"`
		DisplayName   string `json:"display_name"`
		IconURL       string `json:"icon_url"`
		StatusMessage string `json:"status_message"`
	}
)

func newLineAPIClient(msgAPIChannelSecret, msgAPIChannelAccessToken string) (cl *linebot.Client, err error) {
	defer func() {
		err = errors.Wrap(err, "newLineAPIClient error")
	}()

	return linebot.New(
		msgAPIChannelSecret,
		msgAPIChannelAccessToken,
	)
}

func NewLineAPI(cfg config.Config) (l LineAPI, err error) {
	defer func() {
		err = errors.Wrap(err, "NewLineAPI error")
	}()

	cl, err := newLineAPIClient(
		cfg.LineMsgAPIChannelSecret,
		cfg.LineMsgAPIChannelAccessToken,
	)
	if err != nil {
		return nil, err
	}

	return &lineAPI{
		cl:  cl,
		cfg: cfg,
	}, nil
}

func (l *lineAPI) FetchAccessToken(ctx context.Context, authCode string) (token string, err error) {
	defer func() {
		err = errors.Wrap(err, "FetchAccessToken error")
	}()

	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("code", authCode)
	values.Set("redirect_uri", config.LineLoginCallbackURL(l.cfg.FrontURL))
	values.Set("client_id", fmt.Sprint(l.cfg.LineLoginChannelID))
	values.Set("client_secret", l.cfg.LineLoginChannelSecret)
	reqBody := strings.NewReader(values.Encode())

	req, err := http.NewRequest("POST", constant.LineFetchAccessTokenBaseURL, reqBody)
	if err != nil {
		return "", err
	}

	util.SetReqContentTypeURLEncoded(req)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("faild to fetch access token\nstatus code: %d", res.StatusCode)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var resp struct {
		AccessToken string `json:"access_token"`
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		return "", err
	}

	return resp.AccessToken, nil
}

func (l *lineAPI) FetchLineProfileByAccessToken(ctx context.Context, accessToken string) (p *LineProfile, err error) {
	defer func() {
		err = errors.Wrap(err, "FetchLineProfileByAccessToken error")
	}()

	req, err := http.NewRequest("GET", constant.LineProfileBaseURL, nil)
	if err != nil {
		return nil, err
	}

	util.SetBearerToken(req, accessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("faild to fetch LINE user profile. status code: %d", res.StatusCode)
	}

	var resp struct {
		UserID        string `json:"userId"`
		DisplayName   string `json:"displayName"`
		PictureUrl    string `json:"pictureUrl"`
		StatusMessage string `json:"statusMessage"`
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &LineProfile{
		LineID:        resp.UserID,
		DisplayName:   resp.DisplayName,
		IconURL:       resp.PictureUrl,
		StatusMessage: resp.StatusMessage,
	}, nil
}
