package update

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type version struct {
	Major int
	Minor int
	Patch int
}

// var (
// 	errNoNeed = errors.New("no need to update")
// )

func (v *version) Newer(o *version) bool {
	if v.Major > o.Major {
		return true
	} else if v.Major < o.Major {
		return false
	} else if v.Minor > o.Minor {
		return true
	} else if v.Minor < o.Minor {
		return false
	} else if v.Patch > o.Patch {
		return true
	}
	return false
}

func (v *version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

const (
	repo = "https://github.com/v2fly/v2ray-core/releases/"
)

func init() {
	os.Setenv("http_proxy", "http://127.0.0.1:22331")
	os.Setenv("https_proxy", "http://127.0.0.1:22331")
}

func UpdateCore(v *version) error {
	err := downloadCore(v)
	if err != nil {
		return err
	}
	return nil
}

func downloadCore(v *version) error {
	// https://github.com/v2fly/v2ray-core/releases/download/v4.36.2/v2ray-windows-64.zip
	downURL := buildDownloadURL(v)
	res, err := http.Get(downURL)
	if err != nil {
		return err
	}
	durls := strings.Split(downURL, "/")
	fname := durls[len(durls)-1]
	bodyByts, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return os.WriteFile(fname, bodyByts, 0666)
}

func buildDownloadURL(v *version) string {
	arch := runtime.GOARCH
	goos := runtime.GOOS

	return repo + "/download/v" + v.String() + "/" + "v2ray-" + goos + "-" + archConv(arch) + ".zip"
}

func archConv(arch string) string {
	if arch == "amd64" {
		return "64"
	} else if arch == "386" {
		return "32"
	}
	return arch
}

// TODO
func GetCoreVersion(corePath string) (*version, error) {
	return strToVersion(corePath)
}

func GetLatestVersion(url string) (*version, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	ele := doc.Find(".release-header a")
	vStr := ele.First().Text()

	v, err := strToVersion(vStr)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func strToVersion(s string) (*version, error) {
	log.Println(s)
	s = strings.TrimPrefix(s, "v")
	log.Println(s)
	ss := strings.Split(s, ".")

	major, err := strconv.Atoi(ss[0])
	if err != nil {
		return nil, err
	}

	minor, err := strconv.Atoi(ss[1])
	if err != nil {
		return nil, err
	}

	patch, err := strconv.Atoi(ss[2])
	if err != nil {
		return nil, err
	}

	return &version{major, minor, patch}, nil
}
