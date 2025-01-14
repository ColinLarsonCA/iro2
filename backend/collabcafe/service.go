package collabcafe

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

	db      *sql.DB
	scraper *scrapers.CollaboCafeEventScraper
}

func NewService(db *sql.DB) pb.CollabCafeServiceServer {
	return &service{db: db, scraper: &scrapers.CollaboCafeEventScraper{}}
}

func (s *service) GetCollab(ctx context.Context, req *pb.GetCollabRequest) (*pb.GetCollabResponse, error) {
	var id string
	err := s.db.QueryRow("SELECT id FROM collabs WHERE id = $1 LIMIT 1", req.Id).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &pb.GetCollabResponse{Collab: &pb.Collab{Id: id}}, nil
}

func (s *service) SearchCollabs(ctx context.Context, req *pb.SearchCollabsRequest) (*pb.SearchCollabsResponse, error) {
	// TODO(Colin): Implement real search
	res, err := s.ListCollabs(ctx, &pb.ListCollabsRequest{})
	if err != nil {
		return nil, err
	}
	return &pb.SearchCollabsResponse{Collabs: res.GetCollabs()}, nil
}

func (s *service) ListCollabs(ctx context.Context, req *pb.ListCollabsRequest) (*pb.ListCollabsResponse, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, source, source_url, source_posted_at, collab_jp, collab_en, created_at FROM collabs ORDER BY source_posted_at DESC")
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
	jp  *pb.Collab
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
		alreadyHasCollab, _ := s.hasCollabFromSourceURL(ctx, url)
		if alreadyHasCollab {
			continue
		}
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
		collabPairs = append(collabPairs, collabPair{url: collabo.URL, en: en, jp: jp})
	}
	for _, pair := range collabPairs {
		s.insertCollab(ctx, pair)
	}
	return &pb.ScanSourcesResponse{NumNewCollabs: int64(len(collaboCafeURLs))}, nil
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
	var jpBytes, enBytes []byte
	var id, source, sourceURL, sourcePostedAt, createdAt string
	collabs := []*pb.Collab{}
	for rows.Next() {
		err := rows.Scan(&id, &source, &sourceURL, &sourcePostedAt, &jpBytes, &enBytes, &createdAt)
		if err != nil {
			return nil, err
		}
		jp := &pb.Collab{}
		en := &pb.Collab{}
		fmt.Println(jpBytes)
		err = json.Unmarshal(jpBytes, jp)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(enBytes, en)
		if err != nil {
			return nil, err
		}
		if language == "jp" {
			collabs = append(collabs, jp)
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
	jp := collabPair.jp
	en := collabPair.en
	jpJSON, _ := json.Marshal(jp)
	enJSON, _ := json.Marshal(en)
	_, err = s.db.ExecContext(ctx, "INSERT INTO collabs (id, source, source_url, source_posted_at, collab_jp, collab_en) VALUES ($1, $2, $3, $4, $5, $6)",
		en.Id, "collabo-cafe", collabPair.url, en.PostedDate, jpJSON, enJSON)
	if err != nil {
		fmt.Println(err)
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
