package funcs

import (
	"strconv"
	"strings"
	"sync"
)

//Functions used to search through the all artists data to find matchs to search query


func SearchNames(query string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, C := range all_artists {
		if strings.HasPrefix(strings.ToLower(C.Name), strings.ToLower(query)) {
			temp := SearchRes{C.ID, C.Name, "", ""}
			searchResult = append(searchResult, temp)
		}
	}
}

func SearchMem(query string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, C := range all_artists {
		for _, m := range C.Members {
			if strings.HasPrefix(strings.ToLower(m), strings.ToLower(query)) {
				temp := SearchRes{C.ID, C.Name, "", m}
				searchResult = append(searchResult, temp)
			}
		}
	}

}

func SearchLoc(query string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, C := range all_artists {
		for _, l := range C.LocArray {
			if strings.Contains(strings.ToLower(l), strings.ToLower(query)) {
				temp := SearchRes{C.ID, C.Name, l, ""}
				searchResult = append(searchResult, temp)
			}
		}
	}

}

func SearchFirstAlbum(query string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, C := range all_artists {
		if strings.ToLower(C.FirstAlbum) == strings.ToLower(query) {
			temp := SearchRes{C.ID, C.Name, "", ""}
			searchResult = append(searchResult, temp)
		}
	}

}

func SearchCreationDate(query string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, C := range all_artists {
		if strings.ToLower(strconv.Itoa(C.CreationDate)) == strings.ToLower(query) {
			temp := SearchRes{C.ID, C.Name, "", ""}
			searchResult = append(searchResult, temp)
		}
	}
}
