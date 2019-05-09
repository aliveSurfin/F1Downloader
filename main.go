package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// butchered from https://github.com/SoMuchForSubtlety/F1viewer , needs refactored
const urlStart = "https://f1tv.formula1.com"
const sessionURLstart = "https://f1tv.formula1.com/api/session-occurrence/?fields=uid,nbc_status,status,editorial_start_time,live_sources_path,data_source_id,available_for_user,global_channel_urls,global_channel_urls__uid,global_channel_urls__slug,global_channel_urls__self,channel_urls,channel_urls__ovps,channel_urls__slug,channel_urls__name,channel_urls__uid,channel_urls__self,channel_urls__driver_urls,channel_urls__driver_urls__driver_tla,channel_urls__driver_urls__driver_racingnumber,channel_urls__driver_urls__first_name,channel_urls__driver_urls__last_name,channel_urls__driver_urls__image_urls,channel_urls__driver_urls__image_urls__image_type,channel_urls__driver_urls__image_urls__url,channel_urls__driver_urls__team_url,channel_urls__driver_urls__team_url__name,channel_urls__driver_urls__team_url__colour,eventoccurrence_url,eventoccurrence_url__slug,eventoccurrence_url__circuit_url,eventoccurrence_url__circuit_url__short_name,session_type_url,session_type_url__name&fields_to_expand=global_channel_urls,channel_urls,channel_urls__driver_urls,channel_urls__driver_urls__image_urls,channel_urls__driver_urls__team_url,eventoccurrence_url,eventoccurrence_url__circuit_url,session_type_url&slug="

const tagsURL = "https://f1tv.formula1.com/api/tags/"
const vodTypesURL = "http://f1tv.formula1.com/api/vod-type-tag/"
const seriesListURL = "https://f1tv.formula1.com/api/series/"
const seriesF1URL = "https://f1tv.formula1.com/api/series/seri_436bb431c3a24d7d8e200a74e1d11de4/"
const teamsURL = "https://f1tv.formula1.com/api/episodes/"

type episodeStruct struct {
	Subtitle               string    `json:"subtitle"`
	UID                    string    `json:"uid"`
	ScheduleUrls           []string  `json:"schedule_urls"`
	SessionoccurrenceUrls  []string  `json:"sessionoccurrence_urls"`
	Stats                  string    `json:"stats"`
	Title                  string    `json:"title"`
	Self                   string    `json:"self"`
	DriverUrls             []string  `json:"driver_urls"`
	CircuitUrls            []string  `json:"circuit_urls"`
	VodTypeTagUrls         []string  `json:"vod_type_tag_urls"`
	DataSourceFields       []string  `json:"data_source_fields"`
	ParentURL              string    `json:"parent_url"`
	DataSourceID           string    `json:"data_source_id"`
	Tags                   []string  `json:"tags"`
	ImageUrls              []string  `json:"image_urls"`
	SeriesUrls             []string  `json:"series_urls"`
	TeamUrls               []string  `json:"team_urls"`
	HierarchyURL           string    `json:"hierarchy_url"`
	SponsorUrls            []string  `json:"sponsor_urls"`
	PlanUrls               []string  `json:"plan_urls"`
	EpisodeNumber          string    `json:"episode_number"`
	Slug                   string    `json:"slug"`
	LastDataIngest         time.Time `json:"last_data_ingest"`
	Talent                 []string  `json:"talent"`
	Language               string    `json:"language"`
	Created                time.Time `json:"created"`
	Items                  []string  `json:"items"`
	RatingUrls             []string  `json:"rating_urls"`
	Modified               time.Time `json:"modified"`
	RecommendedContentUrls []string  `json:"recommended_content_urls"`
	Synopsis               string    `json:"synopsis"`
	Editability            string    `json:"editability"`
}

type assetStruct struct {
	MaxDevices             interface{}   `json:"max_devices"`
	UID                    string        `json:"uid"`
	ScheduleUrls           []string      `json:"schedule_urls"`
	Self                   string        `json:"self"`
	SessionoccurrenceUrls  []string      `json:"sessionoccurrence_urls"`
	Duration               string        `json:"duration"`
	Stats                  interface{}   `json:"stats"`
	Title                  string        `json:"title"`
	Guidance               bool          `json:"guidance"`
	AssetTypeURL           string        `json:"asset_type_url"`
	DriverUrls             []string      `json:"driver_urls"`
	CircuitUrls            []string      `json:"circuit_urls"`
	DurationInSeconds      int           `json:"duration_in_seconds"`
	Subtitles              bool          `json:"subtitles"`
	DataSourceFields       []string      `json:"data_source_fields"`
	ParentURL              string        `json:"parent_url"`
	DataSourceID           string        `json:"data_source_id"`
	VodTypeTagUrls         []string      `json:"vod_type_tag_urls"`
	StatsLastUpdated       interface{}   `json:"stats_last_updated"`
	Tags                   []interface{} `json:"tags"`
	GuidanceText           string        `json:"guidance_text"`
	AccountUrls            []string      `json:"account_urls"`
	SeriesUrls             []string      `json:"series_urls"`
	TeamUrls               []string      `json:"team_urls"`
	HierarchyURL           string        `json:"hierarchy_url"`
	SponsorUrls            []string      `json:"sponsor_urls"`
	ImageUrls              []string      `json:"image_urls"`
	PlanUrls               []string      `json:"plan_urls"`
	Slug                   string        `json:"slug"`
	LastDataIngest         time.Time     `json:"last_data_ingest"`
	Sound                  bool          `json:"sound"`
	Talent                 []interface{} `json:"talent"`
	Language               string        `json:"language"`
	Created                time.Time     `json:"created"`
	URL                    string        `json:"url"`
	ReleaseDate            interface{}   `json:"release_date"`
	RatingUrls             []string      `json:"rating_urls"`
	Modified               time.Time     `json:"modified"`
	RecommendedContentUrls []string      `json:"recommended_content_urls"`
	Ovps                   []struct {
		AccountURL string `json:"account_url"`
		StreamURL  string `json:"stream_url"`
	} `json:"ovps"`
	Licensor    string `json:"licensor"`
	Editability string `json:"editability"`
}

type seriesStruct struct {
	Name                  string    `json:"name"`
	Language              string    `json:"language"`
	Created               time.Time `json:"created"`
	Self                  string    `json:"self"`
	Modified              time.Time `json:"modified"`
	ImageUrls             []string  `json:"image_urls"`
	ContentUrls           []string  `json:"content_urls"`
	LastDataIngest        time.Time `json:"last_data_ingest"`
	DataSourceFields      []string  `json:"data_source_fields"`
	SessionoccurrenceUrls []string  `json:"sessionoccurrence_urls"`
	Editability           string    `json:"editability"`
	DataSourceID          string    `json:"data_source_id"`
	UID                   string    `json:"uid"`
}

type vodTypesStruct struct {
	Objects []struct {
		Name             string    `json:"name"`
		Language         string    `json:"language"`
		Created          time.Time `json:"created"`
		Self             string    `json:"self"`
		Modified         time.Time `json:"modified"`
		ImageUrls        []string  `json:"image_urls"`
		ContentUrls      []string  `json:"content_urls"`
		LastDataIngest   time.Time `json:"last_data_ingest"`
		DataSourceFields []string  `json:"data_source_fields"`
		Editability      string    `json:"editability"`
		DataSourceID     string    `json:"data_source_id"`
		UID              string    `json:"uid"`
	} `json:"objects"`
}

type driverStruct struct {
	LastName                     string    `json:"last_name"`
	UID                          string    `json:"uid"`
	EventoccurrenceAsWinner1Urls []string  `json:"eventoccurrence_as_winner_1_urls"`
	NationURL                    string    `json:"nation_url"`
	ChannelUrls                  []string  `json:"channel_urls"`
	LastSeason                   int       `json:"last_season"`
	FirstName                    string    `json:"first_name"`
	DriverReference              string    `json:"driver_reference"`
	Self                         string    `json:"self"`
	FirstSeason                  int       `json:"first_season"`
	DriverTla                    string    `json:"driver_tla"`
	DataSourceFields             []string  `json:"data_source_fields"`
	EventoccurrenceAsWinner2Urls []string  `json:"eventoccurrence_as_winner_2_urls"`
	DataSourceID                 string    `json:"data_source_id"`
	DriveroccurrenceUrls         []string  `json:"driveroccurrence_urls"`
	ImageUrls                    []string  `json:"image_urls"`
	LastDataIngest               time.Time `json:"last_data_ingest"`
	EventoccurrenceAsWinner3Urls []string  `json:"eventoccurrence_as_winner_3_urls"`
	Language                     string    `json:"language"`
	Created                      time.Time `json:"created"`
	Modified                     time.Time `json:"modified"`
	ContentUrls                  []string  `json:"content_urls"`
	TeamURL                      string    `json:"team_url"`
	Editability                  string    `json:"editability"`
	DriverRacingnumber           int       `json:"driver_racingnumber"`
}

type teamStruct struct {
	Name                 string    `json:"name"`
	Language             string    `json:"language"`
	Created              time.Time `json:"created"`
	Colour               string    `json:"colour"`
	DriveroccurrenceUrls []string  `json:"driveroccurrence_urls"`
	DriverUrls           []string  `json:"driver_urls"`
	Modified             time.Time `json:"modified"`
	ImageUrls            []string  `json:"image_urls"`
	NationURL            string    `json:"nation_url"`
	ContentUrls          []string  `json:"content_urls"`
	LastDataIngest       time.Time `json:"last_data_ingest"`
	DataSourceFields     []string  `json:"data_source_fields"`
	Self                 string    `json:"self"`
	Editability          string    `json:"editability"`
	DataSourceID         string    `json:"data_source_id"`
	UID                  string    `json:"uid"`
}

type seasonStruct struct {
	Name                     string        `json:"name"`
	Language                 string        `json:"language"`
	Created                  time.Time     `json:"created"`
	ScheduleUrls             []string      `json:"schedule_urls"`
	Self                     string        `json:"self"`
	HasContent               bool          `json:"has_content"`
	ImageUrls                []string      `json:"image_urls"`
	Modified                 time.Time     `json:"modified"`
	ScheduleAfterNextYearURL string        `json:"schedule_after_next_year_url"`
	LastDataIngest           time.Time     `json:"last_data_ingest"`
	DataSourceFields         []interface{} `json:"data_source_fields"`
	Year                     int           `json:"year"`
	EventoccurrenceUrls      []string      `json:"eventoccurrence_urls"`
	Editability              string        `json:"editability"`
	DataSourceID             string        `json:"data_source_id"`
	UID                      string        `json:"uid"`
}

type allSeasonStruct struct {
	Seasons []seasonStruct `json:"objects"`
}

type eventStruct struct {
	EventURL              string    `json:"event_url"`
	UID                   string    `json:"uid"`
	RaceSeasonURL         string    `json:"race_season_url"`
	ScheduleUrls          []string  `json:"schedule_urls"`
	Winner3URL            string    `json:"winner_3_url"`
	OfficialName          string    `json:"official_name"`
	NationURL             string    `json:"nation_url"`
	SessionoccurrenceUrls []string  `json:"sessionoccurrence_urls"`
	CircuitURL            string    `json:"circuit_url"`
	Self                  string    `json:"self"`
	DataSourceFields      []string  `json:"data_source_fields"`
	StartDate             string    `json:"start_date"`
	DataSourceID          string    `json:"data_source_id"`
	EndDate               string    `json:"end_date"`
	ImageUrls             []string  `json:"image_urls"`
	Slug                  string    `json:"slug"`
	LastDataIngest        time.Time `json:"last_data_ingest"`
	Winner2URL            string    `json:"winner_2_url"`
	Name                  string    `json:"name"`
	Language              string    `json:"language"`
	Created               time.Time `json:"created"`
	Modified              time.Time `json:"modified"`
	SponsorURL            string    `json:"sponsor_url"`
	Winner1URL            string    `json:"winner_1_url"`
	Editability           string    `json:"editability"`
}

type sessionStruct struct {
	UID                      string        `json:"uid"`
	ScheduleAfterMidnightURL string        `json:"schedule_after_midnight_url"`
	ScheduleUrls             []string      `json:"schedule_urls"`
	SessionExpiredTime       time.Time     `json:"session_expired_time"`
	ChannelUrls              []string      `json:"channel_urls"`
	GlobalChannelUrls        []string      `json:"global_channel_urls"`
	AvailableForUser         bool          `json:"available_for_user"`
	ScheduleAfter7DaysURL    string        `json:"schedule_after_7_days_url"`
	NbcStatus                string        `json:"nbc_status"`
	Self                     string        `json:"self"`
	ReplayStartTime          time.Time     `json:"replay_start_time"`
	DataSourceFields         []interface{} `json:"data_source_fields"`
	DataSourceID             string        `json:"data_source_id"`
	Status                   string        `json:"status"`
	ScheduleAfter14DaysURL   string        `json:"schedule_after_14_days_url"`
	EventoccurrenceURL       string        `json:"eventoccurrence_url"`
	DriveroccurrenceUrls     []interface{} `json:"driveroccurrence_urls"`
	StartTime                time.Time     `json:"start_time"`
	ImageUrls                []string      `json:"image_urls"`
	LiveSourcesPath          string        `json:"live_sources_path"`
	StatusOverride           interface{}   `json:"status_override"`
	NbcPid                   int           `json:"nbc_pid"`
	LiveSourcesMd5           string        `json:"live_sources_md5"`
	Slug                     string        `json:"slug"`
	LastDataIngest           time.Time     `json:"last_data_ingest"`
	Name                     string        `json:"name"`
	SessionTypeURL           string        `json:"session_type_url"`
	EditorialStartTime       time.Time     `json:"editorial_start_time"`
	EventConfigMd5           string        `json:"event_config_md5"`
	EditorialEndTime         interface{}   `json:"editorial_end_time"`
	Language                 string        `json:"language"`
	Created                  time.Time     `json:"created"`
	Modified                 time.Time     `json:"modified"`
	ContentUrls              []string      `json:"content_urls"`
	ScheduleAfter24HURL      string        `json:"schedule_after_24h_url"`
	EndTime                  time.Time     `json:"end_time"`
	SeriesURL                string        `json:"series_url"`
	SessionName              string        `json:"session_name"`
	Editability              string        `json:"editability"`
}

type sessionStreamsStruct struct {
	Objects []struct {
		Status         string `json:"status"`
		SessionTypeURL struct {
			Name string `json:"name"`
		} `json:"session_type_url"`
		EditorialStartTime time.Time `json:"editorial_start_time"`
		NbcStatus          string    `json:"nbc_status"`
		EventoccurrenceURL struct {
			CircuitURL struct {
				ShortName string `json:"short_name"`
			} `json:"circuit_url"`
			Slug string `json:"slug"`
		} `json:"eventoccurrence_url"`
		LiveSourcesPath   string              `json:"live_sources_path"`
		UID               string              `json:"uid"`
		ChannelUrls       []channelUrlsStruct `json:"channel_urls"`
		GlobalChannelUrls []struct {
			Self string `json:"self"`
			Slug string `json:"slug"`
			UID  string `json:"uid"`
		} `json:"global_channel_urls"`
		DataSourceID     string `json:"data_source_id"`
		AvailableForUser bool   `json:"available_for_user"`
	} `json:"objects"`
}

type channelUrlsStruct struct {
	UID        string             `json:"uid"`
	Self       string             `json:"self"`
	DriverUrls []driverUrlsStruct `json:"driver_urls"`
	Ovps       []struct {
		AccountURL    string `json:"account_url"`
		Path          string `json:"path"`
		Domain        string `json:"domain"`
		FullStreamURL string `json:"full_stream_url"`
	} `json:"ovps"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type driverUrlsStruct struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ImageUrls []struct {
		ImageType string `json:"image_type"`
		URL       string `json:"url"`
	} `json:"image_urls"`
	DriverTla string `json:"driver_tla"`
	TeamURL   struct {
		Colour string `json:"colour"`
		Name   string `json:"name"`
	} `json:"team_url"`
	DriverRacingnumber int `json:"driver_racingnumber"`
}
type homepageContentStruct struct {
	Objects []struct {
		Items []struct {
			Position   int `json:"position"`
			ContentURL struct {
				Items []struct {
					Position    int    `json:"position"`
					ContentType string `json:"content_type"`
					ContentURL  struct {
						Self string `json:"self"`
						UID  string `json:"uid"`
					} `json:"content_url"`
				} `json:"items"`
				Self        string `json:"self"`
				UID         string `json:"uid"`
				SetTypeSlug string `json:"set_type_slug"`
				Title       string `json:"title"`
			} `json:"content_url"`
			ContentType string `json:"content_type"`
			DisplayType string `json:"display_type,omitempty"`
		} `json:"items"`
		Slug        string `json:"slug"`
		SetTypeSlug string `json:"set_type_slug"`
	} `json:"objects"`
}

func main() {
	println("Hello World")
	downloadAllSeasons()
}
func downloadAllSeasons() {
	var allFiles fileStruct
	allSeasons := getSeasons()
	for _, season := range allSeasons.Seasons { // for every season
		println(season.Name)
		allFiles.Seasons = append(allFiles.Seasons, fileSeason{
			Name: season.Name,
		})

		for _, eventID := range season.EventoccurrenceUrls { // for every event in the season
			event := getEvent(eventID) // get event struct
			var curSeason = len(allFiles.Seasons) - 1
			allFiles.Seasons[curSeason].Events = append(allFiles.Seasons[curSeason].Events, fileEvent{
				Name: event.Name,
			})

			println("    " + event.Name)
			fmt.Println("    " + event.EndDate)
			date, err := time.Parse("2006-01-02", event.EndDate)
			if err != nil {
				return
			}
			curDate := time.Now()
			println("current date : " + curDate.String() + `|` + "event date : " + date.String())
			days := curDate.Sub(date).Hours()
			days = days / 24
			fmt.Printf("difference = %v\n", days)
			if days < -7 {
				continue
			}

			for _, sessionID := range event.SessionoccurrenceUrls {
				var curEvent = len(allFiles.Seasons[curSeason].Events) - 1

				session := getSession(sessionID)
				allFiles.Seasons[curSeason].Events[curEvent].Sessions = append(allFiles.Seasons[curSeason].Events[curEvent].Sessions, fileSession{
					Name: session.Name,
				})
				println("    " + "    " + session.Name)
				sessionStreams := getSessionStreams(session.Slug)

				perspectives := sessionStreams.Objects[0].ChannelUrls
				for _, per := range perspectives {
					switch per.Name {
					case "WIF":
						per.Name = "Main Feed"
					case "pit lane":
						per.Name = "Pit Lane"
					case "driver":
						per.Name = "Driver View"
					case "data":
						per.Name = "Data Feed"

					}

					assetURL := getPlayableURL(per.Self)
					//download // need checks for date
					// (url/assetURL) (title/per.Name) (session/session.Name) (track/event.Name) (season/season.Name)

					path := `downloaded/` + season.Name + `/` + event.Name + `/` + session.Name + `/` + per.Name + ".m3u8"
					if _, err := os.Stat(path); os.IsNotExist(err) { // if file does not exist
						path = downloadAsset(assetURL, per.Name, session.Name, event.Name, season.Name)
						//println(path)
						print("Downloading ... ")

					}

					println("    " + "    " + "    " + per.Name)
					var curSession = len(allFiles.Seasons[curSeason].Events[curEvent].Sessions) - 1
					allFiles.Seasons[curSeason].Events[curEvent].Sessions[curSession].Perspectives = append(allFiles.Seasons[curSeason].Events[curEvent].Sessions[curSession].Perspectives, filePerspective{
						Name: per.Name,
						Path: path,
					})
					// JSON FILLING
					// newFile := createFileStruct(event.Name, session.Name, per.Name, path, season.Name)
					// allFiles = append(allFiles, newFile)

				}
			}

		}
	}
	fmt.Printf("%+v\n", allFiles)
	//JSON WRITING
	var jsonData []byte
	jsonData, err := json.Marshal(allFiles)
	if err != nil {
		log.Println(err)
	}
	println("**********************************")
	println(len(jsonData))
	//fmt.Printf("%+v\n", jsonData)
	// var outputFile allFileStruct
	// outputFile.Files = allFiles
	// var jsonData []byte
	// jsonData, err := json.Marshal(outputFile)
	// if err != nil {
	// 	log.Println(err)
	// }
	fmt.Println(string(jsonData))
	f, err := os.Create("test.JSON")
	w := bufio.NewWriter(f)
	n4, err := w.WriteString(string(jsonData))
	fmt.Printf("wrote %d bytes\n", n4)
	w.Flush()
}

//downloads json from URL and returns the json as string and whether it's valid as bool
func getJSON(url string) (bool, string) {
	resp, err := http.Get(url)
	if err != nil {
		debugPrint(err.Error())
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	response := buf.String()
	return isJSON(response), response
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func debugPrint(s string, x ...string) {
	y := s
	for _, str := range x {
		y += " " + str
	}
	println(y)
}
func getEvent(eventID string) eventStruct {
	var event eventStruct
	_, jsonString := getJSON(urlStart + eventID)
	json.Unmarshal([]byte(jsonString), &event)
	return event
}
func getHomepageContent() homepageContentStruct {
	var home homepageContentStruct
	_, jsonString := getJSON("https://f1tv.formula1.com/api/sets/?slug=home&fields=slug,set_type_slug,items,items__position,items__content_type,items__display_type,items__content_url,items__content_url__uid,items__content_url__self,items__content_url__set_type_slug,items__content_url__display_type_slug,items__content_url__title,items__content_url__items,items__content_url__items__set_type_slug,items__content_url__items__position,items__content_url__items__content_type,items__content_url__items__content_url,items__content_url__items__content_url__self,items__content_url__items__content_url__uid&fields_to_expand=items__content_url,items__content_url__items__content_url")
	json.Unmarshal([]byte(jsonString), &home)
	return home
}
func getLive() (ok bool) {
	home := getHomepageContent()
	firstContent := home.Objects[0].Items[0].ContentURL.Items[0].ContentURL.Self
	if strings.Contains(firstContent, "/api/event-occurrence/") {
		event := getEvent(firstContent)
		for _, sessionID := range event.SessionoccurrenceUrls {
			session := getSession(sessionID)
			if session.Status == "live" {
				println(session.Name + " is live")
				return true

			}

		}
	}
	println("not live")
	return false

}

func testing() {
	home := getHomepageContent()
	firstContent := home.Objects[0].Items[0].ContentURL.Items[0].ContentURL.Self
	println(firstContent)
}
func isLive() bool {
	home := getHomepageContent()
	firstContent := home.Objects[0].Items[0].ContentURL.Items[0].ContentURL.Self
	if strings.Contains(firstContent, "/api/event-occurrence/") {
		event := getEvent(firstContent)
		//var fileJSON []fileStruct
		println(event.Name)
		for _, sessionID := range event.SessionoccurrenceUrls {
			session := getSession(sessionID)
			println("**********")

			println("Session :" + session.Name)

			sessionStreams := getSessionStreams(session.Slug)
			perspectives := sessionStreams.Objects[0].ChannelUrls
			for _, per := range perspectives {
				//streamPath := per.Ovps[0]
				println(per.Name)
				//assetURL := getPlayableURL(per.Self)
				//title := per.Name
				//path := downloadAsset(assetURL, title, session.Name, event.Name)
				//newFile := createFileStruct(event.Name, session.Name, title, path)
				//fileJSON = append(fileJSON, newFile)
				//fmt.Printf("%+v\n", streamPath)
				//println(path)
			}
			println("**********")
			if session.Status == "live" {
				println(session.Name + " is live")
				return true

			}

		}
	}
	println("not live")
	return false

}
func getSession(sessionID string) sessionStruct {
	var session sessionStruct
	_, jsonString := getJSON(urlStart + sessionID)
	json.Unmarshal([]byte(jsonString), &session)
	return session
}
func getSessionStreams(sessionSlug string) sessionStreamsStruct {
	var sessionStreams sessionStreamsStruct
	_, jsonString := getJSON(sessionURLstart + sessionSlug)
	json.Unmarshal([]byte(jsonString), &sessionStreams)
	return sessionStreams
}
func getPlayableURL(assetID string) string {
	formattedID := ""
	isChannel := false
	if strings.Contains(assetID, "/api/channels/") {
		isChannel = true
		formattedID = `{"channel_url":"` + assetID + `"}`
	} else {
		formattedID = `{"asset_url":"` + assetID + `"}`
	}
	//make request
	body := strings.NewReader(formattedID)
	req, err := http.NewRequest("POST", "https://f1tv.formula1.com/api/viewings/", body)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//converts response body to string
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	repsAsString := buf.String()

	//extract url form json
	type urlStruct struct {
		Objects []struct {
			Tata struct {
				TokenisedURL string `json:"tokenised_url"`
			} `json:"tata"`
		} `json:"objects"`
	}

	type channelURLstruct struct {
		TokenisedURL string `json:"tokenised_url"`
	}

	var urlString = ""
	if isChannel {
		var finalURL channelURLstruct
		err = json.Unmarshal([]byte(repsAsString), &finalURL)
		if err != nil {
			fmt.Println(err)
		}
		urlString = finalURL.TokenisedURL

	} else {
		var finalURL urlStruct
		json.Unmarshal([]byte(repsAsString), &finalURL)
		urlString = finalURL.Objects[0].Tata.TokenisedURL
	}
	//debugPrint(urlString)
	return strings.Replace(urlString, "&", "\x26", -1)
}

// DOWNLOADING

func downloadAsset(url string, title string, session string, track string, season string) string {
	//sanitize title
	title = strings.Replace(title, ":", "", -1)
	//abort if the URL is not valid
	if len(url) < 10 {
		return ""
	}
	//download and patch .m3u8 file
	data := downloadData(url)
	fixedLineArray := fixData(data, url)
	path := writeToFile(fixedLineArray, title+".m3u8", session, track, season)
	return strings.Replace(path, " ", "\x20", -1)
}
func downloadData(url string) []string {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// convert body to string array
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	lineArray := strings.Split(buf.String(), "\n")
	return lineArray
}

func fixData(lines []string, url string) []string {
	var newLines []string
	//trim url
	var re1 = regexp.MustCompile(`[^\/]*$`)
	url = re1.ReplaceAllString(url, "")

	//fix URLs in m3u8
	for _, line := range lines {
		if strings.Contains(line, "https") {
		} else if len(line) > 6 && (line[:5] == "layer" || line[:4] == "clip" || line[:3] == "OTT") {
			line = url + line
		} else {
			var re = regexp.MustCompile(`[^"]*m3u8"`)
			tempString := re.FindString(line)
			line = re.ReplaceAllString(line, url+tempString)
		}
		var re2 = regexp.MustCompile(`https:\/\/f1tv-cdn[^\.]*\.formula1\.com`)
		line = re2.ReplaceAllString(line, "https://f1tv.secure.footprint.net")
		newLines = append(newLines, line)
	}
	return newLines
}

//write slice of lines to file and return the full file path
func writeToFile(lines []string, path string, session string, track string, season string) string {
	//create downloads folder if it doesnt exist
	if _, err := os.Stat(`/downloaded/` + season + `/` + track + `/` + session); os.IsNotExist(err) {
		os.MkdirAll(`./downloaded/`+season+`/`+track+`/`+session, os.ModePerm)
	}
	path = `./downloaded/` + season + `/` + track + `/` + session + `/` + path
	file, err := os.Create(path)
	if err != nil {
		debugPrint(err.Error())
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	err = w.Flush()
	if err != nil {
		debugPrint(err.Error())
	}
	return path
}

type fileStruct struct {
	Seasons []fileSeason `json:"seasons"`
}
type fileSeason struct {
	Name   string      `json:"Name"`
	Events []fileEvent `json:"events"`
}
type fileEvent struct {
	Name     string        `json:"Name"`
	Sessions []fileSession `json:"sessions"`
}
type fileSession struct {
	Name         string            `json:"Name"`
	Perspectives []filePerspective `json:"perspectives"`
}
type filePerspective struct {
	Name string `json:"Name"`
	Path string `json:"path"`
}

type allFileStruct struct {
	Files []fileStruct `json:"files"`
}

var listOfSeasons allSeasonStruct

func getSeasons() allSeasonStruct {
	if len(listOfSeasons.Seasons) < 1 {
		_, jsonString := getJSON("https://f1tv.formula1.com/api/race-season/?fields=year,name,self,has_content,eventoccurrence_urls&year__gt=2017&order=year")
		json.Unmarshal([]byte(jsonString), &listOfSeasons)
	}
	return listOfSeasons
}
