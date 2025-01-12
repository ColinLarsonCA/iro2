package collabcafe

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/ColinLarsonCA/iro2/backend/collabcafe/scrapers"
	"github.com/ColinLarsonCA/iro2/backend/collabcafe/translators"
	"github.com/ColinLarsonCA/iro2/backend/pb"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
)

type service struct {
	pb.UnimplementedCollabCafeServiceServer

	db      *sql.DB
	scraper *scrapers.CollaboCafeEventScraper
}

func NewService(db *sql.DB) pb.CollabCafeServiceServer {
	return &service{db: db, scraper: &scrapers.CollaboCafeEventScraper{}}
}

func (s *service) GetCollab(ctx context.Context, req *pb.GetCollabRequest) (*pb.GetCollabResponse, error) {
	var id string
	err := s.db.QueryRow("SELECT id FROM collabs WHERE id = %s LIMIT 1", req.Id).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &pb.GetCollabResponse{Collab: &pb.Collab{Id: id}}, nil
}

func (s *service) SearchCollabs(ctx context.Context, req *pb.SearchCollabsRequest) (*pb.SearchCollabsResponse, error) {
	collabs := []*pb.Collab{}
	// TODO(Colin): Implement search
	return &pb.SearchCollabsResponse{Collabs: collabs}, nil
}

type collabPair struct {
	en *pb.Collab
	jp *pb.Collab
}

func (s *service) ScanSources(ctx context.Context, req *pb.ScanSourcesRequest) (*pb.ScanSourcesResponse, error) {
	// TODO(Colin): Implement top-level scraping
	collaboCafeURLs := []string{
		"https://collabo-cafe.com/events/collabo/jujutsukaisen-pop-up-store-shinjuku-marui-annex-2025/",          // pop-up store
		"https://collabo-cafe.com/events/collabo/dakaretai-1st-cafe-animate-ikebukuro2025/",                      // cafe
		"https://collabo-cafe.com/events/collabo/gyagumanga-biyori-25th-anniversary-exhibition-tokyo-osaka2025/", // art exhibit
		"https://collabo-cafe.com/events/collabo/zenless-campaign-family-mart2024-add-info-dry/",                 // konbini
	}
	collabos := []scrapers.Collabo{}
	for _, url := range collaboCafeURLs {
		collabo, err := s.scraper.Scrape(url, "")
		if err != nil {
			log.Println(err)
			continue
		}
		collabos = append(collabos, collabo)
	}

	collabPairs := []collabPair{}
	for _, collabo := range collabos {
		id := uuid.New().String()
		jp := &pb.Collab{
			Id:         id,
			Type:       collabo.Type,
			Slug:       getSlug(collabo.URL),
			PostedDate: collabo.PostedDate,
			Summary: &pb.CollabSummary{
				Thumbnail:   "",
				Title:       collabo.Content.Title,
				Description: "",
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
					Events: []*pb.CollabEvent{},
				},
			},
		}
		en := translators.TranslateJPCollabToEN(jp)
		collabPairs = append(collabPairs, collabPair{en: en, jp: jp})
	}
	for _, pair := range collabPairs {
		spew.Dump(pair.en)
		log.Println()
	}
	return &pb.ScanSourcesResponse{NumNewCollabs: int64(len(collaboCafeURLs))}, nil
}

func getSlug(url string) string {
	url = strings.TrimSuffix(url, "/")
	parts := strings.Split(url, "/")
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
