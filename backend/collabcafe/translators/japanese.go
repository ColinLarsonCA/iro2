package translators

import (
	"regexp"
	"strings"

	"github.com/ColinLarsonCA/iro2/backend/pb"
)

func JPtoEN(text string) string {
	en := TranslateKnownJP(text)
	if en == "" {
		en = GoogleTranslate(text, "jp")
	}
	return en
}

func JPtoENs(texts []string) []string {
	ens := []string{}
	for _, text := range texts {
		ens = append(ens, JPtoEN(text))
	}
	return ens
}

func JPtoENsDetectMonths(texts []string) []string {
	jpYearMonthRegex := regexp.MustCompile("^[0-9]+年[0-9]+月$")
	ens := []string{}
	for _, text := range texts {
		if jpYearMonthRegex.MatchString(text) {
			ens = append(ens, TranslateJPYearMonth(text))
		} else {
			ens = append(ens, JPtoEN(text))
		}
	}
	return ens
}

func TranslateKnownJP(text string) string {
	dictionaries := []map[string]string{dictionary}
	for _, dict := range dictionaries {
		if en, ok := dict[text]; ok {
			return en
		}
	}
	return text
}

// Translates Japanese dates (2024年12月19日) to English dates (2024-12-19)
func TranslateJPDate(text string) string {
	text = strings.Replace(text, "年", "-", 1)
	text = strings.Replace(text, "月", "-", 1)
	text = strings.Replace(text, "日", "", 1)
	return text
}

func TranslateJPYearMonth(text string) string {
	text = strings.Replace(text, "年", "-", 1)
	text = strings.Replace(text, "月", "", 1)
	return text
}

func GoogleTranslate(text string, from string) string {
	if from == "jp" {
		return text
	}
	return text
}

func TranslateJPCollabToEN(jp *pb.Collab) *pb.Collab {
	return &pb.Collab{
		Id:         jp.Id,
		Type:       jp.Type,
		PostedDate: jp.PostedDate,
		Summary: &pb.CollabSummary{
			Thumbnail:   jp.Summary.Thumbnail,
			Title:       JPtoEN(jp.Summary.Title),
			Description: JPtoEN(jp.Summary.Description),
		},
		Content: &pb.CollabContent{
			Series:     JPtoEN(jp.Content.Series),
			Title:      JPtoEN(jp.Content.Title),
			Categories: JPtoENsDetectMonths(jp.Content.Categories),
			Tags:       JPtoENs(jp.Content.Tags),
			OfficialWebsite: &pb.CollabOfficialWebsite{
				Url:  jp.Content.OfficialWebsite.Url,
				Text: JPtoEN(jp.Content.OfficialWebsite.Text),
			},
			Schedule: &pb.CollabSchedule{
				Events: translateJPCollabEventsToEN(jp.Content.Schedule.Events),
			},
		},
	}
}

func translateJPCollabEventToEN(jp *pb.CollabEvent) *pb.CollabEvent {
	return &pb.CollabEvent{
		Location:  JPtoEN(jp.Location),
		Period:    JPtoEN(jp.Period),
		StartDate: TranslateJPDate(jp.StartDate),
		EndDate:   TranslateJPDate(jp.EndDate),
		MapLink:   jp.MapLink,
	}
}

func translateJPCollabEventsToEN(jps []*pb.CollabEvent) []*pb.CollabEvent {
	ens := []*pb.CollabEvent{}
	for _, jp := range jps {
		ens = append(ens, translateJPCollabEventToEN(jp))
	}
	return ens
}

var (
	dictionary = map[string]string{
		// Common words/phrases
		"ポップアップストア":    "Pop-up Store",
		"期間":           "Period",
		"特設ページ":        "Special Page",
		"コラボカフェ":       "Collab Cafe",
		"公式サイト":        "Official Website",
		"開催場所":         "Location",
		"開催期間":         "Period",
		"アクセス・地図":      "Access/Map",
		"描き下ろし (イラスト)": "Original drawing (Illustration)",
		"ニュース":         "News",
		"原画展・展示会":      "Original Art Exhibition",
		"展示会":          "Exhibition",
		"ゲーム":          "Game",

		// Locations
		"新宿":            "Shinjuku",
		"新宿マルイアネックス":    "Shinjuku Marui Annex",
		"東京":            "Tokyo",
		"アニメイトカフェ":      "Animate Cafe",
		"アニメイトカフェ池袋3号店": "Animate Cafe Ikebukuro 3rd Store",
		"池袋":            "Ikebukuro",
		"あべのハルカス近鉄本店":   "Abeno Harukas Kintetsu Main Store",
		"大阪":            "Osaka",
		"有楽町マルイ":        "Marui Yurakucho",
		"ファミリーマート":      "FamilyMart",
		"全国":            "Nationwide",
		"秋葉原":           "Akihabara",
		"名古屋":           "Nagoya",
		"福岡":            "Fukuoka",
		"仙台":            "Sendai",
		"札幌":            "Sapporo",

		// Publication and distribution
		"週刊少年ジャンプ": "Weekly Shonen Jump",
		"集英社":      "Shueisha",

		// Series
		"呪術廻戦": "Jujutsu Kaisen",
		"抱かれたい男1位に脅されています。": "Dakaretai Otoko 1-i ni Odosarete Imasu.",
		"ギャグマンガ日和":          "Gag Manga Biyori",
		"ゼンレスゾーンゼロ":         "Zenless Zone Zero",

		// Creators
		"芥見下々":   "Akutami Gege",
		"桜日梯子":   "Sakurabi Hashigo",
		"増田こうすけ": "Masuda Kosuke",

		// Test dictionary
		"呪術廻戦 新作Ani-Artストア in 新宿マルイアネックス 1月17日より開催!": "Jujutsu Kaisen New Ani-Art Store in Shinjuku Marui Annex opens on January 17th!",
		"だかいち × アニメイトカフェ池袋3号店 1月15日よりコラボ開催!":         "Dakaichi x Animate Cafe Ikebukuro 3rd store collaboration starts on January 15th!",
		"ギャグマンガ日和 25周年記念展 in 東京/大阪 1月10日より開催!":       "Gag Manga Biyori 25th Anniversary Exhibition in Tokyo/Osaka will be held from January 10th!",
		"ゼンレスゾーンゼロ × ファミマ 1月7日より限定ステッカープレゼント!":       "Zenless Zone Zero x FamilyMart Limited sticker giveaway starting January 7th!",
	}
)
