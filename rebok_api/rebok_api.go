package rebok_api

import (
	h "bet25-calendar-sync/helpers"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/s00500/env_logger"
)

type RebokApi struct {
	Client  *http.Client
	cookies map[string]*http.Cookie // name as key

	baseUrl  string
	calendar string

	events      map[int]RebokApiEvent
	cookiesPath string
}

type RebokApiEvent struct {
	Date        string
	Name        string
	Description string
}

type RebookApiRequest struct {
	Startdate      string     `json:"startdate"`
	Days           int        `json:"days"`
	SkipRollOut    string     `json:"skipRollOut"`
	ImageLinkStyle string     `json:"imageLinkStyle"`
	Filter         Filter     `json:"filter"`
	BasicScema     BasicScema `json:"basicScema"`
}

type Filter struct {
	Users           string `json:"Users"`
	UserCompetences string `json:"UserCompetences"`
	Events          string `json:"Events"`
	Labels          string `json:"Labels"`
	Competences     string `json:"Competences"`
	Groups          string `json:"Groups"`
	Customers       string `json:"Customers"`
	EventClass      string `json:"EventClass"`
	Absence         string `json:"Absence"`
}

type BasicScema struct {
	IsActive          string `json:"IsActive"`
	MasterRepSettings string `json:"masterRepSettings"`
}

func NewRebokApi() *RebokApi {
	path, err := os.Getwd()
	log.Should(err)

	// // Load client certificate
	// cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Load CA cert
	// caCert, err := ioutil.ReadFile(caFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM(caCert)

	// // Setup HTTPS client
	// tlsConfig := &tls.Config{
	// 	Certificates: []tls.Certificate{cert},
	// 	RootCAs:      caCertPool,
	// }
	// tlsConfig.BuildNameToCertificate()
	// transport := &http.Transport{TLSClientConfig: tlsConfig}
	// client := &http.Client{Transport: transport, Timeout: time.Second * 5}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	return &RebokApi{
		baseUrl:     "https://rebok.se",
		calendar:    "kanal75",
		Client:      client,
		events:      make(map[int]RebokApiEvent),
		cookiesPath: path + "/cookies",
		cookies:     make(map[string]*http.Cookie),
	}
}

type RebokApiUser struct {
	Username string
	Password string
}

func (r *RebokApi) Login(user RebokApiUser) (string, error) {
	client := r.Client
	// Login to rebok
	if user.Username == "" {
		panic("No username provided")
	}
	if user.Password == "" {
		panic("No password provided")
	}

	reqUrl := fmt.Sprintf("%s/default.aspx?c=kanal75", r.baseUrl)

	// formData := &url.Values{}
	// formData.Add("__EVENTTARGET", url.QueryEscape("ctl00$cpMain$ctl00$ctrlLogin1$btnLogin"))
	// formData.Add("__EVENTARGUMENT", url.QueryEscape(""))
	// formData.Add(url.QueryEscape("ctl00$CtrlMenuTop1$ctrlUserPanel1$activeComp"), url.QueryEscape(r.calendar))
	// formData.Add(url.QueryEscape("ctl00$cpMain$ctl00$ctrlLogin1$tbUsername"), url.QueryEscape(user.Username))
	// formData.Add(url.QueryEscape("ctl00$cpMain$ctl00$ctrlLogin1$tbPassword"), url.QueryEscape(user.Password))
	// formData.Add("__VIEWSTATEGENERATOR", "CA0B0334")
	// formData.Add("__EVENTVALIDATION", url.QueryEscape("/wEdAAYkbw8ORALStEkzBaUjFNDklyAyffgGGOqUedA65u5urP+YDvfolijYlnGq8jIinlI24RKjr/nhpGJlowvOMRw2766HccErQapcOPH22LIzzF3h532NjFUlv0EIV03chIDFoWyGMHxOLzE9ydaiBGS5J1/uyZApv7zrJbZee4Q5kw=="))
	// formData.Add("__VIEWSTATE", url.QueryEscape("/wEPDwUKMTMzMTM1NTAzNA8WAh4TVmFsaWRhdGVSZXF1ZXN0TW9kZQIBFgJmD2QWBmYPZBYCAgUPFgIeBGhyZWYFHi4uL3N0YXRpYy9jc3MvY3NzXzEuNi41LjBfLmNzc2QCAQ9kFgRmD2QWBgIFD2QWAgIBDxYCHgtfIUl0ZW1Db3VudAL/////D2QCBw8WAh8CAv////8PZAIJD2QWAmYPFgIeB1Zpc2libGVoZAIBD2QWAgIBD2QWAmYPZBYCZg9kFghmD2QWAgIBD2QWAgINDw8WAh4EVGV4dGVkZAIBDxYCHwNoZAICD2QWAgIBDw8WAh8DaGRkAgMPZBYCAgEPZBYCAgcPDxYCHwRkZGQCAg8WAh8DaGRk3ZxCe0AQ4RjNnrzUJf//NTz/ue4BqUWNQY96OjXV9Ow="))

	// body := "__EVENTTARGET=ctl00%24cpMain%24ctl00%24ctrlLogin1%24btnLogin&__EVENTARGUMENT=&ctl00%24CtrlMenuTop1%24ctrlUserPanel1%24activeComp=kanal75&ctl00%24cpMain%24ctl00%24ctrlLogin1%24tbUsername=bet25%40kanal75.se&ctl00%24cpMain%24ctl00%24ctrlLogin1%24tbPassword=kanal75&__VIEWSTATEGENERATOR=CA0B0334&__EVENTVALIDATION=%2FwEdAAYkbw8ORALStEkzBaUjFNDklyAyffgGGOqUedA65u5urP%2BYDvfolijYlnGq8jIinlI24RKjr%2FnhpGJlowvOMRw2766HccErQapcOPH22LIzzF3h532NjFUlv0EIV03chIDFoWyGMHxOLzE9ydaiBGS5J1%2FuyZApv7zrJbZee4Q5kw%3D%3D&__VIEWSTATE=%2FwEPDwUKMTMzMTM1NTAzNA8WAh4TVmFsaWRhdGVSZXF1ZXN0TW9kZQIBFgJmD2QWBmYPZBYCAgUPFgIeBGhyZWYFHi4uL3N0YXRpYy9jc3MvY3NzXzEuNi41LjBfLmNzc2QCAQ9kFgRmD2QWBgIFD2QWAgIBDxYCHgtfIUl0ZW1Db3VudAL%2F%2F%2F%2F%2FD2QCBw8WAh8CAv%2F%2F%2F%2F8PZAIJD2QWAmYPFgIeB1Zpc2libGVoZAIBD2QWAgIBD2QWAmYPZBYCZg9kFghmD2QWAgIBD2QWAgINDw8WAh4EVGV4dGVkZAIBDxYCHwNoZAICD2QWAgIBDw8WAh8DaGRkAgMPZBYCAgEPZBYCAgcPDxYCHwRkZGQCAg8WAh8DaGRk3ZxCe0AQ4RjNnrzUJf%2F%2FNTz%2Fue4BqUWNQY96OjXV9Ow%3D"

	// req, err := http.NewRequest(http.MethodPost, reqUrl, strings.NewReader(formData.Encode()))
	// req, err := http.NewRequest(http.MethodPost, reqUrl, nil)
	// if err != nil {
	// 	return "", err
	// }

	// multipartFormData := &multipart.Form{
	// 	Value: url.Values{
	// 		"__EVENTTARGET":   {"ctl00$cpMain$ctl00$ctrlLogin1$btnLogin"},
	// 		"__EVENTARGUMENT": {""},
	// 		"ctl00$CtrlMenuTop1$ctrlUserPanel1$activeComp": {r.calendar},
	// 		"ctl00$cpMain$ctl00$ctrlLogin1$tbUsername":     {user.Username},
	// 		"ctl00$cpMain$ctl00$ctrlLogin1$tbPassword":     {user.Password},
	// 		"__VIEWSTATEGENERATOR":                         {"CA0B0334"},
	// 		"__EVENTVALIDATION":                            {"/wEdAAYkbw8ORALStEkzBaUjFNDklyAyffgGGOqUedA65u5urP+YDvfolijYlnGq8jIinlI24RKjr/nhpGJlowvOMRw2766HccErQapcOPH22LIzzF3h532NjFUlv0EIV03chIDFoWyGMHxOLzE9ydaiBGS5J1/uyZApv7zrJbZee4Q5kw=="},
	// 		"__VIEWSTATE":                                  {"/wEPDwUKMTMzMTM1NTAzNA8WAh4TVmFsaWRhdGVSZXF1ZXN0TW9kZQIBFgJmD2QWBmYPZBYCAgUPFgIeBGhyZWYFHi4uL3N0YXRpYy9jc3MvY3NzXzEuNi41LjBfLmNzc2QCAQ9kFgRmD2QWBgIFD2QWAgIBDxYCHgtfIUl0ZW1Db3VudAL/////D2QCBw8WAh8CAv////8PZAIJD2QWAmYPFgIeB1Zpc2libGVoZAIBD2QWAgIBD2QWAmYPZBYCZg9kFghmD2QWAgIBD2QWAgINDw8WAh4EVGV4dGVkZAIBDxYCHwNoZAICD2QWAgIBDw8WAh8DaGRkAgMPZBYCAgEPZBYCAgcPDxYCHwRkZGQCAg8WAh8DaGRk3ZxCe0AQ4RjNnrzUJf//NTz/ue4BqUWNQY96OjXV9Ow="},
	// 	},
	// }
	// req.MultipartForm = multipartFormData
	// req.ParseMultipartForm(48000)

	// req.Host = "rebok.se"
	// req.Header.Set("Accept", "*/*")
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("Content-Type", "text/plain")
	// req.Header.Set("Content-Length", fmt.Sprintf("%d", len(formData.Encode())))
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")
	// req.Header.Set("Refer", "https://rebok.se/default.aspx")

	// req.Form = formData
	// req.ParseForm()
	// log.Should(err)
	// log.Warn(log.Indent(req.PostForm))

	form := map[string]string{
		"__EVENTTARGET":   "ctl00$cpMain$ctl00$ctrlLogin1$btnLogin",
		"__EVENTARGUMENT": "",
		"ctl00$CtrlMenuTop1$ctrlUserPanel1$activeComp": r.calendar,
		"ctl00$cpMain$ctl00$ctrlLogin1$tbUsername":     user.Username,
		"ctl00$cpMain$ctl00$ctrlLogin1$tbPassword":     user.Password,
		"__VIEWSTATEGENERATOR":                         "CA0B0334",
		"__EVENTVALIDATION":                            "/wEdAAYkbw8ORALStEkzBaUjFNDklyAyffgGGOqUedA65u5urP+YDvfolijYlnGq8jIinlI24RKjr/nhpGJlowvOMRw2766HccErQapcOPH22LIzzF3h532NjFUlv0EIV03chIDFoWyGMHxOLzE9ydaiBGS5J1/uyZApv7zrJbZee4Q5kw==",
		"__VIEWSTATE":                                  "/wEPDwUKMTMzMTM1NTAzNA8WAh4TVmFsaWRhdGVSZXF1ZXN0TW9kZQIBFgJmD2QWBmYPZBYCAgUPFgIeBGhyZWYFHi4uL3N0YXRpYy9jc3MvY3NzXzEuNi41LjBfLmNzc2QCAQ9kFgRmD2QWBgIFD2QWAgIBDxYCHgtfIUl0ZW1Db3VudAL/////D2QCBw8WAh8CAv////8PZAIJD2QWAmYPFgIeB1Zpc2libGVoZAIBD2QWAgIBD2QWAmYPZBYCZg9kFghmD2QWAgIBD2QWAgINDw8WAh4EVGV4dGVkZAIBDxYCHwNoZAICD2QWAgIBDw8WAh8DaGRkAgMPZBYCAgEPZBYCAgcPDxYCHwRkZGQCAg8WAh8DaGRk3ZxCe0AQ4RjNnrzUJf//NTz/ue4BqUWNQY96OjXV9Ow=",
	}
	ct, length, body, err := createForm(form)
	if err != nil {
		panic(err)
	}
	// resp, err := http.Post(reqUrl, ct, body)
	req, err := http.NewRequest(http.MethodPost, reqUrl, body)
	if err != nil {
		return "", err
	}
	req.Host = "rebok.se"
	req.Header.Set("Host", "rebok.se")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", ct)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", length))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/")

	resp, err := client.Do(req)
	// resp, err := http.PostForm(reqUrl, formData)
	if err != nil {
		return "", err
	}


	// req2Url := fmt.Sprintf("%s/adminStart.aspx", r.baseUrl)
	// req2, err := http.NewRequest(http.MethodGet, req2Url, strings.NewReader(formData.Encode()))
	// log.Should(err)

	// req2.Header.Set("Accept", "*/*")
	// req2.Header.Set("Authority", strings.Replace(r.baseUrl, "https://", "", -1))
	// req2.Header.Set("Host", r.baseUrl)
	// req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")
	// // req2.Header.Set("Refer", "https://rebok.se/default.aspx")
	// req2.Header.Set("Content-Length", fmt.Sprintf("%d", len(formData.Encode())))

	// resp2, err := client.Do(req2)
	// // resp, err := http.PostForm(reqUrl, formData)
	// if err != nil {
	// 	return "", err
	// }

	if resp == nil {
		return "", err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	for _, c := range resp.Cookies() {
		r.cookies[c.Name] = c
	}

	req.TLS = &tls.ConnectionState{}

	// log.Warn(log.Indent(resp.Cookies()))
	log.Warn(log.Indent(resp.Header))
	log.Warn(string(respBody))
	log.Warn(log.Indent(req.Host))
	log.Warn(log.Indent(req.Header))
	location := resp.Header.Get("Location")
	log.Warn(log.Indent(log.Indent(location)))
	log.Warn(log.Indent(resp.Request.RequestURI))
	// log.Warn(log.Indent(resp.))

	// log.Warn(log.Indent(req.Header))
	// log.Warn(string(respBody))

	// r.saveCookies()
	r.Client = client
	log.Warn(resp.Status)
	log.Debug(resp.Status)
	return string(respBody), nil
}

func (r *RebokApi) GetEventsFromRebok(users string) (string, error) {
	client := r.Client
	now := time.Now()
	days := 7
	if users == "" {
		return "", fmt.Errorf("No users provided")
	}

	if !strings.HasPrefix(users, ",") {
		users = "," + users
	}
	if !strings.HasSuffix(users, ",") {
		users = users + ","
	}

	// Get events from rebok
	url := fmt.Sprintf("%s%s%d", r.baseUrl, "/MethodProxy/Schema2.aspx/GetEventsFromDate?rnd=", rand.IntN(40000))
	data := RebookApiRequest{
		Startdate:      now.Format(time.RFC3339),
		Days:           days,
		SkipRollOut:    "false",
		ImageLinkStyle: "0",
		Filter: Filter{
			Users:           users,
			UserCompetences: ",All,",
			Events:          ",All,",
			Labels:          ",All,",
			Competences:     ",All,",
			Groups:          ",All,",
			Customers:       ",All,",
			EventClass:      ",All,",
			Absence:         ",All,",
		},
		BasicScema: BasicScema{
			IsActive:          "false",
			MasterRepSettings: fmt.Sprintf("-1,%s,1", now.Format("2006-01-02")),
		},
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", r.baseUrl)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")
	req.Header.Set("Refer", "https://rebok.se/Schema.aspx")
	req.Header.Set("Cookie", fmt.Sprintf("SelectedSchema=7045f041-abf2-409b-9efe-6cb436cbd793; SchemaSetup=Tools=1&AllSchema=0&Wish=0&Available=1&Needs=1&Economy=0&Comment=0&Image=0&Mini=0&Period=%d&StartDate=%s;", days, now.Format("2006-01-02")))

	// cookies := []*http.Cookie{}
	// err = h.ReadJson(&cookies, r.cookiesPath, "cookies.json")
	// log.Should(err)

	for _, c := range r.cookies {
		log.Warnln(log.Indent(c))
		req.AddCookie(c)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp == nil {
		return "", err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// log.Warn(log.Indent(string(respBody)))

	r.Client = client
	log.Debug(resp.Status)
	return string(respBody), nil
}

func (r *RebokApi) SetEvents(events map[int]RebokApiEvent) {
	r.events = events
}

func (r *RebokApi) GetEvents() map[int]RebokApiEvent {
	return r.events
}

func (r *RebokApi) GetEvent(id int) RebokApiEvent {
	return r.events[id]
}

func (r *RebokApi) AddEvent(event RebokApiEvent) {
	r.events[len(r.events)] = event
}

func (r *RebokApi) UpdateEvent(id int, event RebokApiEvent) {
	r.events[id] = event
}

func (r *RebokApi) DeleteEvent(id int) {
	delete(r.events, id)
}

func (r *RebokApi) ClearEvents() {
	r.events = make(map[int]RebokApiEvent)
}

func (r *RebokApi) GetEventCount() int {
	return len(r.events)
}

func (r *RebokApi) saveCookies() error {
	// client := r.Client
	// url, err := url.Parse(r.baseUrl)
	// if err != nil {
	// 	return err
	// }

	cookies := []*http.Cookie{}

	// cookies := client.Jar.Cookies(url)
	// log.Warn(log.Indent(cookies))

	for _, c := range r.cookies {
		log.Warn(log.Indent(c))
		cookies = append(cookies, c)
		// if c.Name == "auth" {
		// 	r.AuthToken = c.Value
		// }
	}

	// store cookies
	err := h.UpdateJson(cookies, r.cookiesPath, "cookies.json")
	log.Should(err)

	// log.Warn(log.Indent(cookies))
	return nil
}

func (r *RebokApi) loadCookies() error {
	client := r.Client
	url, err := url.Parse(r.baseUrl)
	if err != nil {
		return err
	}

	cookies := []*http.Cookie{}
	err = h.ReadJson(&cookies, r.cookiesPath, "cookies.json")
	log.Should(err)

	// for _, c := range cookies {
	// 	if c.Name == "auth" {
	// 		r.AuthToken = c.Value
	// 	}
	// }

	client.Jar.SetCookies(url, cookies)
	r.Client = client

	return nil
}

func createForm(form map[string]string) (string, int, io.Reader, error) {
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	defer mp.Close()
	for key, val := range form {
		key = url.PathEscape(key)
		val = url.PathEscape(val)

		if strings.HasPrefix(val, "@") {
			val = val[1:]
			file, err := os.Open(val)
			if err != nil {
				return "", 0, nil, err
			}
			defer file.Close()
			part, err := mp.CreateFormFile(key, val)
			if err != nil {
				return "", 0, nil, err
			}
			io.Copy(part, file)
		} else {
			mp.WriteField(key, val)
		}
	}
	return mp.FormDataContentType(), len(body.Bytes()), body, nil
}
