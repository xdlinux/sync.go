package mirrorz

type MirrorZ struct {
	Version float32   `json:"version"`
	Site    Site      `json:"site"`
	Info    []Info    `json:"info"`
	Mirrors []*Mirror `json:"mirrors"`
}

type Site struct {
	Url            string `json:"url"`
	Logo           string `json:"logo"`
	LogoDark       string `json:"logo_darkmode"`
	Abbreviation   string `json:"abbr"`
	Name           string `json:"name"`
	Homepage       string `json:"homepage"`
	IssueTracker   string `json:"issue"`
	RequestTracker string `json:"request"`
	Email          string `json:"email"`
	Group          string `json:"group"`
	DiskInfo       string `json:"disk"`
	Note           string `json:"note"`
	BigFilePath    string `json:"big"`
}

type Info struct {
	Distro   string    `json:"distro"`
	Category string    `json:"category"`
	Urls     []InfoURL `json:"urls"`
}

type InfoURL struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Mirror struct {
	Cname    string       `json:"cname"`
	Desc     string       `json:"desc"`
	URL      string       `json:"url"`
	Status   MirrorStatus `json:"status"`
	Help     string       `json:"help,omitempty"`
	Upstream string       `json:"upstream,omitempty"`
	Size     string       `json:"size,omitempty"`
}
