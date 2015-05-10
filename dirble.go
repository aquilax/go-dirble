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

type Timestamped struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	Slug               string
	DisableSongChecks  bool
	Timestamped
}

type Stations []StationItem

type Stream struct {
	ID           int
	Stream       string
	Bitrate      int
	ContenType   string
	Status       int
	StationId    int `json:"Station_id"`
	Timedout     bool
	EmptyCounter int
	Timestamped
}

type Category struct {
	ID          int
	Title       string
	Description string
	URLId       interface{}
	Slug        string
	Ancestry    string
	Position    interface{}
	Timestamped
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

type StationSongs []StationSong
type Categories []Category

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
	Categories   Categories
	StationSongs StationSongs `json:"station_songs"`
}

type CategoryTree struct {
	Children Categories
	Category
}

type CategoryStations []Station

type Country struct {
	ID          int
	CountryCode string `json:"country_code"`
	ContinentID int    `json:"Continent_id"`
	Timestamped
}

type Countries []Country

type Continent struct {
	ID       int
	Name     string
	Slug     string
	LegacyID int `json:"legacy_id"`
	Timestamped
}

type Continents []Continent

type SearchResults []Station

func New(rt http.RoundTripper, token string) *Dirble {
	return &Dirble{
		client: &http.Client{
			Transport: rt,
		},
		token: token,
	}
}

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

func (d *Dirble) getStations(url string) (*Stations, error) {
	var err error
	var content []byte
	if content, err = d.fetchURL(url); err != nil {
		return nil, err
	}
	var s Stations
	if err = json.Unmarshal(content, &s); err != nil {
		return nil, err
	}
	return &s, nil
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
	return d.getStations(u.String())
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

func (d *Dirble) StationSongHistory(id int) (*StationSongs, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "stations/" + strconv.Itoa(id) + "/song_history"
	q := u.Query()
	q.Set("token", d.token)
	u.RawQuery = q.Encode()
	var content []byte
	if content, err = d.fetchURL(u.String()); err != nil {
		return nil, err
	}
	var ss StationSongs
	if err = json.Unmarshal(content, &ss); err != nil {
		return nil, err
	}
	return &ss, nil
}

func (d *Dirble) StationSimilar(id int) (*Stations, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "stations/" + strconv.Itoa(id) + "/similar"
	q := u.Query()
	q.Set("token", d.token)
	u.RawQuery = q.Encode()
	return d.getStations(u.String())
}

func (d *Dirble) Categories() (*Categories, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "categories"
	q := u.Query()
	q.Set("token", d.token)
	u.RawQuery = q.Encode()
	var content []byte
	if content, err = d.fetchURL(u.String()); err != nil {
		return nil, err
	}
	var c Categories
	if err = json.Unmarshal(content, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (d *Dirble) CategoriesPrimary() (*Categories, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "categories/primary"
	q := u.Query()
	q.Set("token", d.token)
	u.RawQuery = q.Encode()
	var content []byte
	if content, err = d.fetchURL(u.String()); err != nil {
		return nil, err
	}
	var c Categories
	if err = json.Unmarshal(content, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (d *Dirble) CategoriesTree() (*CategoryTree, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "categories/tree"
	q := u.Query()
	q.Set("token", d.token)
	u.RawQuery = q.Encode()
	var content []byte
	if content, err = d.fetchURL(u.String()); err != nil {
		return nil, err
	}
	var ct CategoryTree
	if err = json.Unmarshal(content, &ct); err != nil {
		return nil, err
	}
	return &ct, nil
}

func (d *Dirble) CategoryStations(id int, all bool, page, perPage, offset *int) (*CategoryStations, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "category/" + strconv.Itoa(id) + "/stations"
	q := u.Query()
	q.Set("token", d.token)
	if all {
		q.Set("all", "1")
	}
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
	var cs CategoryStations
	if err = json.Unmarshal(content, &cs); err != nil {
		return nil, err
	}
	return &cs, nil
}

func (d *Dirble) CategoryChilds(id int) (*Categories, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "category/" + strconv.Itoa(id) + "/childs"
	q := u.Query()
	q.Set("token", d.token)
	u.RawQuery = q.Encode()
	var content []byte
	if content, err = d.fetchURL(u.String()); err != nil {
		return nil, err
	}
	var c Categories
	if err = json.Unmarshal(content, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (d *Dirble) Countries() (*Countries, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "countries"
	q := u.Query()
	q.Set("token", d.token)
	u.RawQuery = q.Encode()
	var content []byte
	if content, err = d.fetchURL(u.String()); err != nil {
		return nil, err
	}
	var c Countries
	if err = json.Unmarshal(content, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (d *Dirble) CountriesStations(code string, all bool, page, perPage, offset *int) (*CategoryStations, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "countries/" + code + "/stations"
	q := u.Query()
	q.Set("token", d.token)
	if all {
		q.Set("all", "1")
	}
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
	var cs CategoryStations
	if err = json.Unmarshal(content, &cs); err != nil {
		return nil, err
	}
	return &cs, nil
}

func (d *Dirble) Continents() (*Continents, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "continents"
	q := u.Query()
	q.Set("token", d.token)
	u.RawQuery = q.Encode()
	var content []byte
	if content, err = d.fetchURL(u.String()); err != nil {
		return nil, err
	}
	var c Continents
	if err = json.Unmarshal(content, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (d *Dirble) ContinentsCountries(continentId int) (*Countries, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "continents/" + strconv.Itoa(continentId) + "/countries"
	q := u.Query()
	q.Set("token", d.token)
	u.RawQuery = q.Encode()
	var content []byte
	if content, err = d.fetchURL(u.String()); err != nil {
		return nil, err
	}
	var c Countries
	if err = json.Unmarshal(content, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (d *Dirble) Search(query string, page *int) (*SearchResults, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "search/" + query
	q := u.Query()
	q.Set("token", d.token)

	u.RawQuery = q.Encode()
	var content []byte
	if content, err = d.fetchURL(u.String()); err != nil {
		return nil, err
	}
	var sr SearchResults
	if err = json.Unmarshal(content, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}
