package funcs

import (
	"net/http"
	"strconv"
	"strings"
)

// Called to filter results when user submits the form
func Filter(artists []Artist, r *http.Request, w http.ResponseWriter) {
	err := r.ParseForm()
	if err != nil {
		res := Result{}
		res.Code = 500
		res.Status = "Internal Server Error"
		ErrorHandler(w, r, &res)
		return
	} else {
		//Retrieve filter values from the form
		minCreate := r.FormValue("minC")
		maxCreate := r.FormValue("maxC")
		minA := r.FormValue("minA")
		maxA := r.FormValue("maxA")
		numMembers := r.Form["mem[]"]
		location := r.FormValue("location-filter")

		artists := all_artists

		result := ArtistArray{}
		result.Empty = false
		if minCreate != "" && maxCreate != "" && minA != "" && maxA != "" {
			minC, err1 := strconv.Atoi(minCreate)
			maxC, err2 := strconv.Atoi(maxCreate)
			minAl, err3 := strconv.Atoi(minA)
			maxAl, err4 := strconv.Atoi(maxA)

			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				res := Result{}
				res.Code = 500
				res.Status = "Internal Server Error"
				ErrorHandler(w, r, &res)
				return
			}

			for _, c := range artists {
				firstDate := strings.Split(c.FirstAlbum, "-")
				album := firstDate[2]
				albumInt, err5 := strconv.Atoi(album)

				if err5 != nil {
					res := Result{}
					res.Code = 500
					res.Status = "Internal Server Error"
					ErrorHandler(w, r, &res)
					return
				}

				if c.CreationDate >= minC && c.CreationDate <= maxC && albumInt >= minAl && albumInt <= maxAl {
					if len(numMembers) != 0 {
						for _, n := range numMembers {
							num, err6 := strconv.Atoi(n)

							if err6 != nil {
								res := Result{}
								res.Code = 500
								res.Status = "Internal Server Error"
								ErrorHandler(w, r, &res)
								return
							}

							if len(c.Members) == num {
								if location != "" {
									locas := c.LocArray
									for _, l := range locas {
										if l == location {
											result.Artists = append(result.Artists, c)
											break
										}
									}
								} else {
									result.Artists = append(result.Artists, c)
									break
								}
							}
						}
					} else {
						if location != "" {
							locas := c.LocArray
							for _, l := range locas {
								if l == location {
									result.Artists = append(result.Artists, c)
									break
								}
							}
						} else {
							result.Artists = append(result.Artists, c)
						}
					}
				}
			}
		}

		if len(result.Artists) == 0 {
			result.Empty = true
		}

		//Incase there is no error in fetching the data, execute the index.html file with the artists array
		err := tpl.ExecuteTemplate(w, "index.html", &result)
		if err != nil {
			res := Result{}
			res.Code = 500
			res.Status = "Internal Server Error"
			ErrorHandler(w, r, &res)
			return
		}

	}
}
