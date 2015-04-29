package dirble

import (
	"net/http"
)

type Dirble struct {
	transport *http.RoundTripper
}

func New(transport *http.RoundTripper) *Dirble {
	return &Dirble{transport}
}

func (d *Dirble) Stations()                            {}
func (d *Dirble) Station(id int)                       {}
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
