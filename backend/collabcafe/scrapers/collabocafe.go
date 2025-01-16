package scrapers

import (
	"log"
	"regexp"
	"strings"

	neturl "net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type CollaboCafeEventScraper struct {
}

type OfficialWebsite struct {
	URL  string
	Text string
}

type CollaboEvent struct {
	Location  string
	Period    string
	StartDate string
	EndDate   string
	MapLink   string
}

type CollaboSchedule struct {
	Events []CollaboEvent
}

type CollaboImages struct {
	Header string
}

type CollaboSummary struct {
	Thumbnail   string
	Title       string
	Description string
}

type CollaboContent struct {
	Series          string
	Categories      []string
	Tags            []string
	Title           string
	OfficialWebsite OfficialWebsite
	Schedule        CollaboSchedule
	Images          CollaboImages
}

type Collabo struct {
	URL        string
	Type       string
	PostedDate string
	Summary    CollaboSummary
	Content    CollaboContent
}

var (
	sublocationRegex = regexp.MustCompile("【.*】")
)

func (s *CollaboCafeEventScraper) ScrapeHomepage() (map[string]CollaboSummary, error) {
	collector := colly.NewCollector(
		colly.AllowedDomains("collabo-cafe.com"),
	)
	summaries := map[string]CollaboSummary{}
	collector.OnHTML(".top-post-list", func(e *colly.HTMLElement) {
		log.Println("found collabo list")
		e.ForEach("article", func(i int, e *colly.HTMLElement) {
			url := e.ChildAttr("a", "href")
			summaries[url] = CollaboSummary{
				Thumbnail:   e.ChildAttr("img", "src"),
				Title:       e.ChildText(".entry-title"),
				Description: e.ChildText(".description"),
			}
		})
	})
	err := collector.Visit("https://collabo-cafe.com/")
	if err != nil {
		return nil, err
	}
	collector.Wait()
	return summaries, nil
}

func (s *CollaboCafeEventScraper) ScrapeCollaboPage(url string, summary CollaboSummary) (Collabo, error) {
	collector := colly.NewCollector(
		colly.AllowedDomains("collabo-cafe.com"),
	)
	c := Collabo{
		URL:     url,
		Summary: summary,
		Content: CollaboContent{},
	}
	collector.OnHTML("#main", func(e *colly.HTMLElement) {
		log.Println(url, "found collabo page content")
		c.PostedDate = e.ChildAttr("time", "datetime")
		c.Content.Title = e.ChildText(".entry-title")
		c.Content.Images.Header = e.DOM.Find(".eyecatch").First().Find("img").First().AttrOr("src", "")
		e.DOM.Find(".eo__event__category").First().Find("a").Each(func(i int, s *goquery.Selection) {
			c.Content.Categories = append(c.Content.Categories, s.Text())
		})
		e.DOM.Find(".eo__event__tags").First().Find("a").Each(func(i int, s *goquery.Selection) {
			c.Content.Tags = append(c.Content.Tags, s.Text())
		})
		c.Type = findEventTypeInCategories(c.Content.Categories)
		c.Content.Series = guessSeries(c.Content.Title, c.Content.Categories)
		eventDetailsTable := e.DOM.Find(".table__container").First()
		events := []CollaboEvent{}
		eventDetailsTable.Find("tr").Each(func(i int, s *goquery.Selection) {
			events = ensureAtLeastNEvents(events, 1)
			keyNode := s.Find("th").First()
			valueNode := s.Find("td").First()
			key := keyNode.Text()
			switch key {
			case "公式サイト": // Official Website
				c.Content.OfficialWebsite = OfficialWebsite{
					Text: strings.TrimSpace(valueNode.Text()),
					URL:  stripUTMParams(valueNode.Find("a").AttrOr("href", "")),
				}
			case "開催場所": // Location
				if nodeContainsMultipleDatums(valueNode) {
					events = ensureAtLeastNEvents(events, countNodeDatums(valueNode))
					eventIndex := 0
					valueNode.Contents().Each(func(_ int, s *goquery.Selection) {
						if !s.Is("br") {
							if eventIndex < len(events) {
								events[eventIndex].Location = strings.TrimSpace(s.Text())
								eventIndex++
							} else {
								log.Println("eventIndex out of bounds when scraping locations:", url)
							}
						}
					})
				} else {
					events[0].Location = valueNode.Text()
				}
			case "開催期間": // Period
				if nodeContainsMultipleDatums(valueNode) {
					events = ensureAtLeastNEvents(events, countNodeDatums(valueNode))
					eventIndex := 0
					valueNode.Contents().Each(func(_ int, s *goquery.Selection) {
						if !s.Is("br") {
							if eventIndex < len(events) {
								period := s.Text()
								period = sublocationRegex.ReplaceAllString(period, "")
								events[eventIndex].Period = strings.TrimSpace(period)
								events[eventIndex].StartDate, events[eventIndex].EndDate = parseStartDateAndEndDate(events[eventIndex].Period)
								eventIndex++
							} else {
								log.Println("eventIndex out of bounds when scraping periods:", url)
							}
						}
					})
				} else {
					events[0].Period = valueNode.Text()
					events[0].StartDate, events[0].EndDate = parseStartDateAndEndDate(events[0].Period)
				}
			case "アクセス・地図": // Map Link
				if nodeContainsMultipleDatums(valueNode) {
					events = ensureAtLeastNEvents(events, countNodeDatums(valueNode))
					eventIndex := 0
					valueNode.Find("a").Each(func(_ int, s *goquery.Selection) {
						if eventIndex < len(events) {
							events[eventIndex].MapLink = s.AttrOr("href", "")
							eventIndex++
						} else {
							log.Println("eventIndex out of bounds when scraping map links:", url)
						}
					})
				} else {
					events[0].MapLink = valueNode.Find("a").AttrOr("href", "")
				}
			}
		})
		c.Content.Schedule.Events = events
	})
	err := collector.Visit(url)
	if err != nil {
		return c, err
	}
	collector.Wait()
	return c, nil
}

func nodeContainsMultipleDatums(node *goquery.Selection) bool {
	return countNodeDatums(node) > 1
}

func countNodeDatums(node *goquery.Selection) int {
	return node.Find("br").Length() + 1
}

func ensureAtLeastNEvents(events []CollaboEvent, n int) []CollaboEvent {
	if len(events) >= n {
		return events
	}
	for i := len(events); i < n; i++ {
		events = append(events, CollaboEvent{})
	}
	return events
}

func stripUTMParams(url string) string {
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return ""
	}
	queryParams := parsedURL.Query()
	for key := range queryParams {
		if strings.HasPrefix(key, "utm_") {
			queryParams.Del(key)
		}
	}
	parsedURL.RawQuery = queryParams.Encode()
	return parsedURL.String()
}

func findEventTypeInCategories(categories []string) string {
	for _, category := range categories {
		if category == "コラボカフェ" {
			return "Collab Cafe"
		}
		if category == "ポップアップストア" {
			return "Pop-up Store"
		}
		if category == "原画展・展示会" {
			return "Original Art Exhibition"
		}
		if category == "コンビニ" {
			return "Convenience Store"
		}
	}
	return "Unknown"
}

func guessSeries(title string, categories []string) string {
	dateCategoryRegex := regexp.MustCompile("^[0-9]+年[0-9]+月$")
	possibleSeries := []string{}
	for _, category := range categories {
		// filter out common categories
		if category == "ポップアップストア" { // pop-up store
			continue
		}
		if category == "コラボカフェ" { // collab cafe
			continue
		}
		if category == "期間" { // period
			continue
		}

		// filter out date categories
		if dateCategoryRegex.MatchString(category) {
			continue
		}

		// if the category is in the title, it's probably the series
		if strings.Contains(title, category) {
			return category
		}
		possibleSeries = append(possibleSeries, category)
	}
	// if there is only one possible series, return it
	if len(possibleSeries) == 1 {
		return possibleSeries[0]
	}
	return "Unknown"
}

func parseStartDateAndEndDate(period string) (string, string) {
	// 2025年1月17日〜2月2日
	dateRange := strings.Split(period, "〜")
	if len(dateRange) != 2 {
		return "", ""
	}
	startDate := strings.TrimSpace(dateRange[0])
	endDate := strings.TrimSpace(dateRange[1])
	// if the end date is missing the year, add it by copying the year from the start date
	if regexp.MustCompile("^[0-9]+月[0-9]+日$").MatchString(endDate) {
		year := regexp.MustCompile("^[0-9]+年").FindString(startDate)
		endDate = year + endDate
	}
	return startDate, endDate
}
