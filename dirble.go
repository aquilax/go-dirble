package dirble

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Version is the library version
const Version = "0.1"

// APIBase is the root URL of the API
const APIBase = "http://api.dirble.com/v2/"

type Dirble struct {
	client *http.Client
	token  string
}

type Thumb struct {
	URL *string
}

type Image struct {
	URL   *string
	Thumb Thumb
}

type StationItem struct {
	ID                 int
	Name               string
	Description        string
	Country            string
	Accepted           int
	Added              time.Time
	Website            string
	CurrentSongReverse *string `json:"currentsong_reverse"`
	Image              Image
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Slug               string
	DisableSongChecks  bool
}

type Stations []StationItem

type Stream struct {
	ID           int
	Stream       string
	Bitrate      int
	ContenType   string
	Status       int
	StationId    int       `json:"Station_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Timedout     bool
	EmptyCounter int
}

type Category struct {
	ID          int
	Title       string
	Description string
	URLId       interface{}
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Slug        string
	Ancestry    string
	Position    interface{}
}

type StationSong struct {
	ID        string
	Name      string
	Title     string
	Week      int
	Year      int
	StationID int `json:"station_id"`
	Info      *string
	Date      time.Time
}

type Station struct {
	ID           int
	Name         string
	Accepted     int
	Added        time.Time
	Country      string
	Description  string
	Image        Image
	Slug         string
	Website      string
	Streams      []Stream
	Categories   []Category
	StationSongs []StationSong `json:"station_songs"`
}

func New(rt http.RoundTripper, token string) *Dirble {
	return &Dirble{
		client: &http.Client{
			Transport: rt,
		},
		token: token,
	}
}

func (d *Dirble) Stations(page, perPage, offset *int) (*Stations, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "stations"
	q := u.Query()
	q.Set("token", d.token)
	if page != nil {
		q.Set("page", strconv.Itoa(*page))
	}
	if perPage != nil {
		q.Set("per_page", strconv.Itoa(*perPage))
	}
	if offset != nil {
		q.Set("offset", strconv.Itoa(*offset))
	}
	u.RawQuery = q.Encode()
	var content []byte
	if content, err = d.fetchURL(u.String()); err != nil {
		return nil, err
	}
	var s Stations
	if err = json.Unmarshal(content, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (d *Dirble) Station(id int) (*Station, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "station/" + strconv.Itoa(id)
	q := u.Query()
	q.Set("token", d.token)
	u.RawQuery = q.Encode()
	var content []byte
	if content, err = d.fetchURL(u.String()); err != nil {
		return nil, err
	}
	var s Station
	if err = json.Unmarshal(content, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (d *Dirble) StationSongHistory(id int)            {}
func (d *Dirble) StationSimilar(id int)                {}
func (d *Dirble) Categories()                          {}
func (d *Dirble) CategoriesPrimary()                   {}
func (d *Dirble) CategoriesTree()                      {}
func (d *Dirble) CategoryStations(id int)              {}
func (d *Dirble) CategoryChilds(id int)                {}
func (d *Dirble) Countries()                           {}
func (d *Dirble) CountriesStations(code string)        {}
func (d *Dirble) Continents()                          {}
func (d *Dirble) ContinentsCountries(continent string) {}
func (d *Dirble) Search(query string)                  {}

func (d *Dirble) fetchURL(url string) ([]byte, error) {
	var err error
	var resp *http.Response
	if resp, err = d.client.Get(url); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Error: %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
