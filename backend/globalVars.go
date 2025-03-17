package main

import "time"

type FilterStruct struct {
	Tags      []string `json:"tags"`
	Name      []string `json:"name"`
	Platforms []string `json:"platforms"`
	Devs      []string `json:"devs"`
}

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

type SteamScreenshotStruct struct {
	ID            int    `json:"id"`
	PathThumbnail string `json:"path_thumbnail"`
	PathFull      string `json:"path_full"`
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

type igdbSearchResult []struct {
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

type SteamWishlistStruct struct {
	Response struct {
		Items []struct {
			Appid     int `json:"appid"`
			Priority  int `json:"priority"`
			DateAdded int `json:"date_added"`
		} `json:"items"`
	} `json:"response"`
}

var clientID string
var clientSecret string

type igdbMetaData struct {
	Name               string
	UID                string
	Summary            string
	ReleaseDateTime    string
	AggregatedRating   float64
	ScreenshotPaths    ImgStruct
	CoverArtPath       ImgStruct
	InvolvedCompanies  TagsStruct
	Themes             TagsStruct
	PlayerPerspectives TagsStruct
	Genres             TagsStruct
	GameModes          TagsStruct
}

var PsGameStruct struct {
	Titles []struct {
		TitleID           string `json:"titleId"`
		Name              string `json:"name"`
		LocalizedName     string `json:"localizedName"`
		ImageURL          string `json:"imageUrl"`
		LocalizedImageURL string `json:"localizedImageUrl"`
		Category          string `json:"category"`
		Service           string `json:"service"`
		PlayCount         int    `json:"playCount"`
		Concept           struct {
			ID       int      `json:"id"`
			TitleIds []string `json:"titleIds"`
			Name     string   `json:"name"`
			Media    struct {
				Audios []interface{} `json:"audios"`
				Videos []interface{} `json:"videos"`
				Images []struct {
					URL    string `json:"url"`
					Format string `json:"format"`
					Type   string `json:"type"`
				} `json:"images"`
			} `json:"media"`
			Genres        []string `json:"genres"`
			LocalizedName struct {
				DefaultLanguage string `json:"defaultLanguage"`
				Metadata        struct {
					FiFI   string `json:"fi-FI"`
					UkUA   string `json:"uk-UA"`
					DeDE   string `json:"de-DE"`
					EnUS   string `json:"en-US"`
					KoKR   string `json:"ko-KR"`
					PtBR   string `json:"pt-BR"`
					EsES   string `json:"es-ES"`
					ArAE   string `json:"ar-AE"`
					NoNO   string `json:"no-NO"`
					FrCA   string `json:"fr-CA"`
					ItIT   string `json:"it-IT"`
					PlPL   string `json:"pl-PL"`
					RuRU   string `json:"ru-RU"`
					ZhHans string `json:"zh-Hans"`
					NlNL   string `json:"nl-NL"`
					PtPT   string `json:"pt-PT"`
					ZhHant string `json:"zh-Hant"`
					SvSE   string `json:"sv-SE"`
					DaDK   string `json:"da-DK"`
					TrTR   string `json:"tr-TR"`
					FrFR   string `json:"fr-FR"`
					EnGB   string `json:"en-GB"`
					Es419  string `json:"es-419"`
					JaJP   string `json:"ja-JP"`
				} `json:"metadata"`
			} `json:"localizedName"`
			Country  string `json:"country"`
			Language string `json:"language"`
		} `json:"concept"`
		Media struct {
			Audios []interface{} `json:"audios"`
			Videos []interface{} `json:"videos"`
			Images []struct {
				URL    string `json:"url"`
				Format string `json:"format"`
				Type   string `json:"type"`
			} `json:"images"`
		} `json:"media"`
		FirstPlayedDateTime time.Time `json:"firstPlayedDateTime"`
		LastPlayedDateTime  time.Time `json:"lastPlayedDateTime"`
		PlayDuration        string    `json:"playDuration"`
	} `json:"titles"`
}
var PSTrophyStruct struct {
	TrophyTitles []struct {
		NpServiceName       string `json:"npServiceName"`
		NpCommunicationID   string `json:"npCommunicationId"`
		TrophySetVersion    string `json:"trophySetVersion"`
		TrophyTitleName     string `json:"trophyTitleName"`
		TrophyTitleIconURL  string `json:"trophyTitleIconUrl"`
		TrophyTitlePlatform string `json:"trophyTitlePlatform"`
		HasTrophyGroups     bool   `json:"hasTrophyGroups"`
		TrophyGroupCount    int    `json:"trophyGroupCount"`
		DefinedTrophies     struct {
			Bronze   int `json:"bronze"`
			Silver   int `json:"silver"`
			Gold     int `json:"gold"`
			Platinum int `json:"platinum"`
		} `json:"definedTrophies"`
		Progress       int `json:"progress"`
		EarnedTrophies struct {
			Bronze   int `json:"bronze"`
			Silver   int `json:"silver"`
			Gold     int `json:"gold"`
			Platinum int `json:"platinum"`
		} `json:"earnedTrophies"`
		HiddenFlag          bool      `json:"hiddenFlag"`
		LastUpdatedDateTime time.Time `json:"lastUpdatedDateTime"`
		TrophyTitleDetail   string    `json:"trophyTitleDetail,omitempty"`
	} `json:"trophyTitles"`
	TotalItemCount int `json:"totalItemCount"`
}

type IGDBInsertGameReturn struct {
	Title             string `json:"title"`
	ReleaseDate       string `json:"releaseDate"`
	SelectedPlatforms []struct {
		Value string `json:"value"`
		Label string `json:"label"`
	} `json:"selectedPlatforms"`
	TimePlayed   string `json:"timePlayed"`
	Rating       string `json:"rating"`
	SelectedDevs []struct {
		Value string `json:"value"`
		Label string `json:"label"`
	} `json:"selectedDevs"`
	SelectedTags []struct {
		Value string `json:"value"`
		Label string `json:"label"`
	} `json:"selectedTags"`
	Description string   `json:"description"`
	CoverImage  string   `json:"coverImage"`
	SSImage     []string `json:"ssImage"`
	IsWishlist  int      `json:"isWishlist"`
}
