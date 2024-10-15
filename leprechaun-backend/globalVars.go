package main

type allSteamGamesStruct struct {
	Response struct {
		GameCount int `json:"game_count"`
		Games     []struct {
			Appid                    int     `json:"appid"`
			Name                     string  `json:"name"`
			PlaytimeForever          float32 `json:"playtime_forever"`
			ImgIconURL               string  `json:"img_icon_url"`
			PlaytimeWindowsForever   int     `json:"playtime_windows_forever"`
			PlaytimeMacForever       int     `json:"playtime_mac_forever"`
			PlaytimeLinuxForever     int     `json:"playtime_linux_forever"`
			PlaytimeDeckForever      int     `json:"playtime_deck_forever"`
			RtimeLastPlayed          int     `json:"rtime_last_played"`
			PlaytimeDisconnected     int     `json:"playtime_disconnected"`
			HasCommunityVisibleStats bool    `json:"has_community_visible_stats,omitempty"`
			ContentDescriptorids     []int   `json:"content_descriptorids,omitempty"`
			HasLeaderboards          bool    `json:"has_leaderboards,omitempty"`
			Playtime2Weeks           int     `json:"playtime_2weeks,omitempty"`
		} `json:"games"`
	} `json:"response"`
}

type SteamGameMetadataStruct struct {
	Success bool `json:"success"`
	Data    struct {
		Type       string `json:"type"`
		Name       string `json:"name"`
		SteamAppid int    `json:"steam_appid"`
		//RequiredAge         string `json:"required_age"`
		IsFree              bool   `json:"is_free"`
		Dlc                 []int  `json:"dlc"`
		DetailedDescription string `json:"detailed_description"`
		AboutTheGame        string `json:"about_the_game"`
		ShortDescription    string `json:"short_description"`
		SupportedLanguages  string `json:"supported_languages"`
		Reviews             string `json:"reviews"`
		HeaderImage         string `json:"header_image"`
		CapsuleImage        string `json:"capsule_image"`
		CapsuleImagev5      string `json:"capsule_imagev5"`
		Website             string `json:"website"`
		/* PcRequirements      struct {
			Minimum     string `json:"minimum"`
			Recommended string `json:"recommended"`
		} `json:"pc_requirements"` */
		//MacRequirements   []interface{} `json:"mac_requirements"`
		//LinuxRequirements []interface{} `json:"linux_requirements"`
		LegalNotice   string   `json:"legal_notice"`
		Developers    []string `json:"developers"`
		Publishers    []string `json:"publishers"`
		PriceOverview struct {
			Currency         string `json:"currency"`
			Initial          int    `json:"initial"`
			Final            int    `json:"final"`
			DiscountPercent  int    `json:"discount_percent"`
			InitialFormatted string `json:"initial_formatted"`
			FinalFormatted   string `json:"final_formatted"`
		} `json:"price_overview"`
		/* Packages      []int `json:"packages"`
		PackageGroups []struct {
			Name                    string `json:"name"`
			Title                   string `json:"title"`
			Description             string `json:"description"`
			SelectionText           string `json:"selection_text"`
			SaveText                string `json:"save_text"`
			DisplayType             int    `json:"display_type"`
			IsRecurringSubscription string `json:"is_recurring_subscription"`
			Subs                    []struct {
				Packageid                int    `json:"packageid"`
				PercentSavingsText       string `json:"percent_savings_text"`
				PercentSavings           int    `json:"percent_savings"`
				OptionText               string `json:"option_text"`
				OptionDescription        string `json:"option_description"`
				CanGetFreeLicense        string `json:"can_get_free_license"`
				IsFreeLicense            bool   `json:"is_free_license"`
				PriceInCentsWithDiscount int    `json:"price_in_cents_with_discount"`
			} `json:"subs"`
		} `json:"package_groups"` */
		Platforms struct {
			Windows bool `json:"windows"`
			Mac     bool `json:"mac"`
			Linux   bool `json:"linux"`
		} `json:"platforms"`
		Metacritic struct {
			Score int    `json:"score"`
			URL   string `json:"url"`
		} `json:"metacritic"`
		Categories []struct {
			ID          int    `json:"id"`
			Description string `json:"description"`
		} `json:"categories"`
		Genres []struct {
			ID          string `json:"id"`
			Description string `json:"description"`
		} `json:"genres"`
		Screenshots []struct {
			ID            int    `json:"id"`
			PathThumbnail string `json:"path_thumbnail"`
			PathFull      string `json:"path_full"`
		} `json:"screenshots"`
		Movies []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Thumbnail string `json:"thumbnail"`
			Webm      struct {
				Num480 string `json:"480"`
				Max    string `json:"max"`
			} `json:"webm"`
			Mp4 struct {
				Num480 string `json:"480"`
				Max    string `json:"max"`
			} `json:"mp4"`
			Highlight bool `json:"highlight"`
		} `json:"movies"`
		Recommendations struct {
			Total int `json:"total"`
		} `json:"recommendations"`
		Achievements struct {
			Total       int `json:"total"`
			Highlighted []struct {
				Name string `json:"name"`
				Path string `json:"path"`
			} `json:"highlighted"`
		} `json:"achievements"`
		ReleaseDate struct {
			ComingSoon bool   `json:"coming_soon"`
			Date       string `json:"date"`
		} `json:"release_date"`
		SupportInfo struct {
			URL   string `json:"url"`
			Email string `json:"email"`
		} `json:"support_info"`
		Background         string `json:"background"`
		BackgroundRaw      string `json:"background_raw"`
		ContentDescriptors struct {
			Ids   []int  `json:"ids"`
			Notes string `json:"notes"`
		} `json:"content_descriptors"`
		Ratings struct {
			Esrb struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"esrb"`
			Pegi struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"pegi"`
			Usk struct {
				Rating string `json:"rating"`
			} `json:"usk"`
			Oflc struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"oflc"`
			Nzoflc struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"nzoflc"`
			Kgrb struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"kgrb"`
			Dejus struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"dejus"`
			Mda struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"mda"`
			Fpb struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"fpb"`
			Csrr struct {
				Rating string `json:"rating"`
			} `json:"csrr"`
			Crl struct {
				Rating string `json:"rating"`
			} `json:"crl"`
			SteamGermany struct {
				RatingGenerated string `json:"rating_generated"`
				Rating          string `json:"rating"`
				RequiredAge     string `json:"required_age"`
				Banned          string `json:"banned"`
				UseAgeGate      string `json:"use_age_gate"`
				Descriptors     string `json:"descriptors"`
			} `json:"steam_germany"`
		} `json:"ratings"`
	} `json:"data"`
}

type gameStruct []struct {
	ID                    int     `json:"id"`
	AlternativeNames      []int   `json:"alternative_names"`
	Category              int     `json:"category"`
	Cover                 int     `json:"cover"`
	CreatedAt             int     `json:"created_at"`
	ExternalGames         []int   `json:"external_games"`
	FirstReleaseDate      int     `json:"first_release_date,omitempty"`
	Franchises            []int   `json:"franchises,omitempty"`
	GameEngines           []int   `json:"game_engines,omitempty"`
	Genres                []int   `json:"genres,omitempty"`
	Hypes                 int     `json:"hypes,omitempty"`
	InvolvedCompanies     []int   `json:"involved_companies"`
	Name                  string  `json:"name"`
	Platforms             []int   `json:"platforms,omitempty"`
	ReleaseDates          []int   `json:"release_dates,omitempty"`
	Screenshots           []int   `json:"screenshots,omitempty"`
	SimilarGames          []int   `json:"similar_games,omitempty"`
	Slug                  string  `json:"slug"`
	Summary               string  `json:"summary"`
	Tags                  []int   `json:"tags,omitempty"`
	Themes                []int   `json:"themes,omitempty"`
	UpdatedAt             int     `json:"updated_at"`
	URL                   string  `json:"url"`
	Websites              []int   `json:"websites,omitempty"`
	Checksum              string  `json:"checksum"`
	Collections           []int   `json:"collections,omitempty"`
	AgeRatings            []int   `json:"age_ratings,omitempty"`
	AggregatedRating      float64 `json:"aggregated_rating,omitempty"`
	AggregatedRatingCount int     `json:"aggregated_rating_count,omitempty"`
	GameModes             []int   `json:"game_modes,omitempty"`
	Keywords              []int   `json:"keywords,omitempty"`
	PlayerPerspectives    []int   `json:"player_perspectives,omitempty"`
	Rating                float64 `json:"rating,omitempty"`
	RatingCount           int     `json:"rating_count,omitempty"`
	Storyline             string  `json:"storyline,omitempty"`
	TotalRating           float64 `json:"total_rating,omitempty"`
	TotalRatingCount      int     `json:"total_rating_count,omitempty"`
	Videos                []int   `json:"videos,omitempty"`
	Remakes               []int   `json:"remakes,omitempty"`
	ExpandedGames         []int   `json:"expanded_games,omitempty"`
	Ports                 []int   `json:"ports,omitempty"`
	GameLocalizations     []int   `json:"game_localizations,omitempty"`
	Artworks              []int   `json:"artworks,omitempty"`
	Bundles               []int   `json:"bundles,omitempty"`
	Dlcs                  []int   `json:"dlcs,omitempty"`
	Expansions            []int   `json:"expansions,omitempty"`
	LanguageSupports      []int   `json:"language_supports,omitempty"`
	Status                int     `json:"status,omitempty"`
	ParentGame            int     `json:"parent_game,omitempty"`
	Collection            int     `json:"collection,omitempty"`
	VersionParent         int     `json:"version_parent,omitempty"`
	VersionTitle          string  `json:"version_title,omitempty"`
}
type TagsStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type ImgStruct []struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

var playerPerspectiveStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
var themeStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
var genresStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
var gameModesStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
var involvedCompaniesStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
var gameEngineStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var coverStruct []struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}
var screenshotStruct []struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

var summary string
var releaseDateTime string
var AggregatedRating float64
var Name string

const clientID = "bg50w140115zmfq2pi0uc0wujj9pn6"
const clientSecret = "1nk95mh97tui5t1ct1q5i7sqyfmqvd"
