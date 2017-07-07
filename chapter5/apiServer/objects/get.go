package objects

import (
	"../../lib/es"
	"../locate"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	version := r.URL.Query()["version"]
	var hash string
	if len(version) == 0 {
		_, hash, _ = es.SearchLatestVersion(name)
	} else {
		v, e := strconv.Atoi(version[0])
		if e != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		hash, _ = es.GetHash(name, v)
	}
	if hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	object := url.PathEscape(hash)
	s := locate.Locate(object)
	if len(s) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	request, e := http.NewRequest("GET", "http://"+s[0].Addr+"/objects/"+object, r.Body)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := http.Client{}
	nr, e := client.Do(request)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(nr.StatusCode)
	io.Copy(w, nr.Body)
}
