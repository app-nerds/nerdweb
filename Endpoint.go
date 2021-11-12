package nerdweb

import (
	"net/http"
	"strings"
)

/*
Endpoint defines a single HTTP endpoint. Each endpoint is used
to configure a Gorilla Mux route.
*/
type Endpoint struct {
	Path        string
	Methods     []string
	HandlerFunc http.HandlerFunc
	Handler     http.Handler
}

/*
Endpoints represents an Endpoint slice.
*/
type Endpoints []*Endpoint

func (a Endpoints) Len() int {
	return len(a)
}

func (a Endpoints) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Endpoints) Less(i, j int) bool {

	firstDynamic := isDynamic(a[i].Path)
	secondDynamic := isDynamic(a[j].Path)

	if firstDynamic && !secondDynamic {
		return false
	}

	if !firstDynamic && secondDynamic {
		return true
	}

	if len(a[i].Path) != len(a[j].Path) {
		return len(a[i].Path) > len(a[j].Path)
	}

	// if a[i].Method != a[j].Method {
	// 	return a[i].Method != "GET"
	// }

	if a[i].Path == a[j].Path {
		panic("Two endpoints can't be same")
	}

	return true
}

func isDynamic(url string) bool {
	return strings.Contains(url, "{") && strings.Contains(url, "}")
}
