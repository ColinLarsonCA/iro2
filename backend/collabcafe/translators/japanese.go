package translators

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

	translate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
	"github.com/ColinLarsonCA/iro2/backend/pb"
)

type JapaneseTranslator struct {
	db *sql.DB
}

func NewJapaneseTranslator(db *sql.DB) *JapaneseTranslator {
	return &JapaneseTranslator{db: db}
}

func (ja *JapaneseTranslator) ToEnglish(text string) string {
	if text == "" {
		return ""
	}
	en := ja.GetStoredTranslation(text)
	if en == "" {
		var err error
		en, err = ja.GoogleTranslate(text, "ja")
		if err != nil {
			log.Printf("%v\n", err)
		}
		if en == "" {
			en = text
		} else {
			ja.StoreTranslation(text, en)
		}
	}
	return en
}

func (ja *JapaneseTranslator) ManyToEnglish(texts []string) []string {
	ens := []string{}
	for _, text := range texts {
		ens = append(ens, ja.ToEnglish(text))
	}
	return ens
}

func (ja *JapaneseTranslator) ManyJapaneseToEnglishDetectMonths(texts []string) []string {
	jpYearMonthRegex := regexp.MustCompile("^[0-9]+年[0-9]+月$")
	ens := []string{}
	for _, text := range texts {
		if jpYearMonthRegex.MatchString(text) {
			ens = append(ens, ja.YearPlusMonth(text))
		} else {
			ens = append(ens, ja.ToEnglish(text))
		}
	}
	return ens
}

// Translates Japanese dates (2024年12月19日) to English dates (2024-12-19)
func (ja *JapaneseTranslator) Date(text string) string {
	text = strings.Replace(text, "年", "-", 1)
	text = strings.Replace(text, "月", "-", 1)
	text = strings.Replace(text, "日", "", 1)
	return text
}

func (ja *JapaneseTranslator) YearPlusMonth(text string) string {
	text = strings.Replace(text, "年", "-", 1)
	text = strings.Replace(text, "月", "", 1)
	return text
}

func (ja *JapaneseTranslator) GoogleTranslate(text string, from string) (string, error) {
	if from != "ja" {
		return "", fmt.Errorf("language not supported: %s", from)
	}
	ctx := context.Background()
	client, err := translate.NewTranslationClient(ctx)
	if err != nil {
		return "", fmt.Errorf("could not create GoogleTranslate client: %v", err)
	}
	defer client.Close()

	req := &translatepb.TranslateTextRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", "iro2-448003"),
		SourceLanguageCode: "ja",
		TargetLanguageCode: "en-US",
		MimeType:           "text/plain", // Mime types: "text/plain", "text/html"
		Contents:           []string{text},
	}

	resp, err := client.TranslateText(ctx, req)
	if err != nil {
		return "", fmt.Errorf("could not Google Translate text: %v", err)
	}

	for _, translation := range resp.GetTranslations() {
		return translation.GetTranslatedText(), nil
	}
	return "", fmt.Errorf("no translations found for text: %s", text)
}

func (ja *JapaneseTranslator) CollabToEnglish(jp *pb.Collab) *pb.Collab {
	return &pb.Collab{
		Id:         jp.Id,
		Type:       jp.Type,
		Slug:       jp.Slug,
		PostedDate: jp.PostedDate,
		Summary: &pb.CollabSummary{
			Thumbnail:   jp.Summary.Thumbnail,
			Title:       ja.ToEnglish(jp.Summary.Title),
			Description: ja.ToEnglish(jp.Summary.Description),
		},
		Content: &pb.CollabContent{
			Series:     ja.ToEnglish(jp.Content.Series),
			Title:      ja.ToEnglish(jp.Content.Title),
			Categories: ja.ManyJapaneseToEnglishDetectMonths(jp.Content.Categories),
			Tags:       ja.ManyToEnglish(jp.Content.Tags),
			OfficialWebsite: &pb.CollabOfficialWebsite{
				Url:  jp.Content.OfficialWebsite.Url,
				Text: ja.ToEnglish(jp.Content.OfficialWebsite.Text),
			},
			Schedule: &pb.CollabSchedule{
				Events: ja.Events(jp.Content.Schedule.Events),
			},
		},
	}
}

func (ja *JapaneseTranslator) Event(event *pb.CollabEvent) *pb.CollabEvent {
	return &pb.CollabEvent{
		Location:  ja.ToEnglish(event.Location),
		Period:    ja.ToEnglish(event.Period),
		StartDate: ja.Date(event.StartDate),
		EndDate:   ja.Date(event.EndDate),
		MapLink:   event.MapLink,
	}
}

func (ja *JapaneseTranslator) Events(events []*pb.CollabEvent) []*pb.CollabEvent {
	ens := []*pb.CollabEvent{}
	for _, event := range events {
		ens = append(ens, ja.Event(event))
	}
	return ens
}

func (ja *JapaneseTranslator) GetStoredTranslation(text string) string {
	var en string
	ja.db.QueryRow("SELECT en FROM ja_to_en_lookup WHERE ja = $1", text).Scan(&en)
	return en
}

func (ja *JapaneseTranslator) StoreTranslation(japanese, english string) {
	_, err := ja.db.Exec("INSERT INTO ja_to_en_lookup (ja, en) VALUES ($1, $2) ON CONFLICT (ja) DO UPDATE SET en = $3", japanese, english, english)
	if err != nil {
		log.Printf("Failed to store translation: %v\n", err)
	}
}
