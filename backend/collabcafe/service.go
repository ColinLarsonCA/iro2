package collabcafe

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/ColinLarsonCA/iro2/backend/collabcafe/scrapers"
	"github.com/ColinLarsonCA/iro2/backend/collabcafe/translators"
	"github.com/ColinLarsonCA/iro2/backend/pb"
	"github.com/google/uuid"
)

type service struct {
	pb.UnimplementedCollabCafeServiceServer

	db       *sql.DB
	japanese *translators.JapaneseTranslator
	scraper  *scrapers.CollaboCafeEventScraper
}

func NewService(db *sql.DB) pb.CollabCafeServiceServer {
	scraper := &scrapers.CollaboCafeEventScraper{}
	japanese := translators.NewJapaneseTranslator(db)
	return &service{db: db, scraper: scraper, japanese: japanese}
}

func (s *service) GetCollab(ctx context.Context, req *pb.GetCollabRequest) (*pb.GetCollabResponse, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, source, source_url, source_posted_at, collab_ja, collab_en, created_at FROM collabs WHERE id = $1 LIMIT 1;", req.GetId())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	collabs, err := s.scanCollabRows(rows, req.Language)
	if err != nil {
		return nil, err
	}
	if len(collabs) == 0 {
		return &pb.GetCollabResponse{}, nil
	}
	return &pb.GetCollabResponse{Collab: collabs[0]}, nil
}

func (s *service) SearchCollabs(ctx context.Context, req *pb.SearchCollabsRequest) (*pb.SearchCollabsResponse, error) {
	ftsLanguage := "english"
	if req.GetLanguage() == "ja" {
		ftsLanguage = "japanese"
	}
	rows, err := s.db.QueryContext(ctx, "SELECT id, source, source_url, source_posted_at, collab_ja, collab_en, created_at FROM collabs WHERE fts_collab_en @@ websearch_to_tsquery($1, $2) ORDER BY source_posted_at DESC;", ftsLanguage, req.GetQuery())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	collabs, err := s.scanCollabRows(rows, req.Language)
	if err != nil {
		return nil, err
	}
	return &pb.SearchCollabsResponse{Collabs: collabs}, nil
}

func (s *service) ListCollabs(ctx context.Context, req *pb.ListCollabsRequest) (*pb.ListCollabsResponse, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, source, source_url, source_posted_at, collab_ja, collab_en, created_at FROM collabs ORDER BY source_posted_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	collabs, err := s.scanCollabRows(rows, req.Language)
	if err != nil {
		return nil, err
	}
	return &pb.ListCollabsResponse{Collabs: collabs}, nil
}

type collabPair struct {
	url string
	en  *pb.Collab
	ja  *pb.Collab
}

func (s *service) ScanSources(ctx context.Context, req *pb.ScanSourcesRequest) (*pb.ScanSourcesResponse, error) {
	categories := []string{"cafe", "convenience-store", "pop-up-store"}
	urlsToSummaries := map[string]scrapers.CollaboSummary{}
	for _, category := range categories {
		summaries, err := s.scraper.ScrapeCategory(category)
		if err != nil {
			log.Printf("error scraping %s: %v\n", category, err)
		}
		for url, summary := range summaries {
			urlsToSummaries[url] = summary
		}
	}
	collabos := []scrapers.Collabo{}
	for url, summary := range urlsToSummaries {
		alreadyHasCollab, _ := s.hasCollabFromSourceURL(ctx, url)
		if alreadyHasCollab {
			continue
		}
		collabo, err := s.scraper.ScrapeCollaboPage(url, summary)
		if err != nil {
			log.Println(err)
			continue
		}
		collabos = append(collabos, collabo)
	}

	collabPairs := []collabPair{}
	for _, collabo := range collabos {
		id := uuid.New().String()
		ja := &pb.Collab{
			Id:         id,
			Type:       collabo.Type,
			Slug:       getSlug(collabo.URL),
			PostedDate: collabo.PostedDate,
			Summary: &pb.CollabSummary{
				Thumbnail:   collabo.Summary.Thumbnail,
				Title:       collabo.Summary.Title,
				Description: collabo.Summary.Description,
			},
			Content: &pb.CollabContent{
				Series:     collabo.Content.Series,
				Title:      collabo.Content.Title,
				Categories: collabo.Content.Categories,
				Tags:       collabo.Content.Tags,
				OfficialWebsite: &pb.CollabOfficialWebsite{
					Url:  collabo.Content.OfficialWebsite.URL,
					Text: collabo.Content.OfficialWebsite.Text,
				},
				Schedule: &pb.CollabSchedule{
					Events: s.mapCollaboEvents(collabo.Content.Schedule.Events),
				},
			},
			Images: &pb.CollabImages{
				Header: collabo.Content.Images.Header,
			},
		}
		en := s.japanese.CollabToEnglish(ja)
		collabPairs = append(collabPairs, collabPair{url: collabo.URL, en: en, ja: ja})
	}
	for _, pair := range collabPairs {
		s.insertCollab(ctx, pair)
	}
	return &pb.ScanSourcesResponse{NumNewCollabs: int64(len(collabPairs))}, nil
}

func (s *service) mapCollaboEvents(events []scrapers.CollaboEvent) []*pb.CollabEvent {
	pbEvents := []*pb.CollabEvent{}
	for _, event := range events {
		pbEvents = append(pbEvents, s.mapCollaboEvent(event))
	}
	return pbEvents
}

func (s *service) mapCollaboEvent(event scrapers.CollaboEvent) *pb.CollabEvent {
	return &pb.CollabEvent{
		Location:  event.Location,
		Period:    event.Period,
		StartDate: event.StartDate,
		EndDate:   event.EndDate,
		MapLink:   event.MapLink,
	}
}

func (s *service) hasCollabFromSourceURL(ctx context.Context, url string) (bool, error) {
	err := s.db.QueryRowContext(ctx, "SELECT id FROM collabs WHERE source_url = $1 LIMIT 1;", url).Scan()
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *service) scanCollabRows(rows *sql.Rows, language string) ([]*pb.Collab, error) {
	var jaBytes, enBytes []byte
	var id, source, sourceURL, sourcePostedAt, createdAt string
	collabs := []*pb.Collab{}
	for rows.Next() {
		err := rows.Scan(&id, &source, &sourceURL, &sourcePostedAt, &jaBytes, &enBytes, &createdAt)
		if err != nil {
			return nil, err
		}
		ja := &pb.Collab{}
		en := &pb.Collab{}
		err = json.Unmarshal(jaBytes, ja)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(enBytes, en)
		if err != nil {
			return nil, err
		}
		if language == "ja" {
			collabs = append(collabs, ja)
		} else {
			collabs = append(collabs, en)
		}
	}
	return collabs, nil
}

func (s *service) insertCollab(ctx context.Context, collabPair collabPair) error {
	hasCollabAlready, err := s.hasCollabFromSourceURL(ctx, collabPair.url)
	if err != nil {
		return err
	}
	if hasCollabAlready {
		return nil
	}
	ja := collabPair.ja
	en := collabPair.en
	jaJSON, _ := json.Marshal(ja)
	enJSON, _ := json.Marshal(en)
	_, err = s.db.ExecContext(ctx, "INSERT INTO collabs (id, source, source_url, source_posted_at, collab_ja, collab_en) VALUES ($1, $2, $3, $4, $5, $6)",
		en.Id, "collabo-cafe", collabPair.url, en.PostedDate, jaJSON, enJSON)
	if err != nil {
		return err
	}
	return nil
}

func getSlug(url string) string {
	trimmed := strings.TrimSuffix(url, "/")
	parts := strings.Split(trimmed, "/")
	return parts[len(parts)-1]
}

func parseCorrectTime(layout string, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		log.Println(err)
		return time.Time{}
	}
	return t
}
