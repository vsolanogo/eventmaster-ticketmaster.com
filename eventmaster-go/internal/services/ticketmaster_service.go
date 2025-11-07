package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"eventmaster-go/internal/models"
	"eventmaster-go/internal/repositories"
)

type TicketmasterService struct {
	eventRepo          repositories.EventRepository
	imageService       ImageService
	participantService ParticipantService
	apiKey             string
	systemUserID       string
}

func NewTicketmasterService(
	eventRepo repositories.EventRepository,
	imageService ImageService,
	participantService ParticipantService,
	apiKey string,
	systemUserID string,
) *TicketmasterService {
	return &TicketmasterService{
		eventRepo:          eventRepo,
		imageService:       imageService,
		participantService: participantService,
		apiKey:             apiKey,
		systemUserID:       systemUserID,
	}
}

type TicketmasterEvent struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Images      []struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"images"`
	Embedded struct {
		Venues []struct {
			Name string `json:"name"`
			City struct {
				Name string `json:"name"`
			} `json:"city"`
			Country struct {
				Name string `json:"name"`
			} `json:"country"`
			Location struct {
				Latitude  string `json:"latitude"`
				Longitude string `json:"longitude"`
			} `json:"location"`
		} `json:"venues"`
		Attractions []struct {
			Images []struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"images"`
		} `json:"attractions"`
	} `json:"_embedded"`
	Dates struct {
		Start struct {
			DateTime  string `json:"dateTime"`
			LocalDate string `json:"localDate"`
			LocalTime string `json:"localTime"`
		} `json:"start"`
		Status struct {
			Code string `json:"code"`
		} `json:"status"`
	} `json:"dates"`
	Classifications []struct {
		Segment struct {
			Name string `json:"name"`
		} `json:"segment"`
		Genre struct {
			Name string `json:"name"`
		} `json:"genre"`
		SubGenre struct {
			Name string `json:"name"`
		} `json:"subGenre"`
	} `json:"classifications"`
}

type TicketmasterResponse struct {
	Embedded struct {
		Events []TicketmasterEvent `json:"events"`
	} `json:"_embedded"`
}


func (s *TicketmasterService) FetchAndSaveEvents() error {
	if s.apiKey == "" {
		return errors.New("ticketmaster API key not configured")
	}

	url := fmt.Sprintf(
		"https://app.ticketmaster.com/discovery/v2/events.json?countryCode=US&size=100&apikey=%s",
		s.apiKey,
	)

	log.Printf("Ticketmaster fetch started: url=%s", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching events from Ticketmaster: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	var data TicketmasterResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("error parsing JSON response: %w", err)
	}

	if len(data.Embedded.Events) == 0 {
		log.Printf("Ticketmaster fetch completed: no events returned")
		return nil
	}

	var (
		newEvents     int
		skippedEvents int
	)
	for _, tmEvent := range data.Embedded.Events {
		// Check if event already exists
		existing, _ := s.eventRepo.FindByExternalID(tmEvent.ID)
		if existing != nil {
			skippedEvents++
			continue // Skip if event already exists
		}

		event := s.mapToEvent(tmEvent)
		if s.systemUserID != "" {
			event.UserID = s.systemUserID
		}

		// Attach images from Ticketmaster payload
		imageLinks := collectImageLinks(tmEvent)
		if len(imageLinks) > 0 {
			images, err := s.imageService.CreateImagesWithLinks(imageLinks)
			if err != nil {
				log.Printf("Ticketmaster image creation failed: id=%s err=%v", tmEvent.ID, err)
			} else if len(images) > 0 {
				event.Images = make([]models.Image, len(images))
				for i, img := range images {
					event.Images[i] = *img
				}
			}
		}

		err := s.eventRepo.Create(event)
		if err != nil {
			// Log error but continue with next event
			// In a production environment, you might want to implement retry logic
			log.Printf("Ticketmaster event save failed: id=%s err=%v", tmEvent.ID, err)
			continue
		}

		// Generate fake participants similar to NestJS implementation
		participantCount := 2
		if err := s.participantService.GenerateFakeParticipants(event, participantCount); err != nil {
			log.Printf("Ticketmaster participant generation failed: id=%s err=%v", tmEvent.ID, err)
		}
		newEvents++
	}

 	log.Printf("Ticketmaster fetch completed: total=%d new=%d skipped=%d", len(data.Embedded.Events), newEvents, skippedEvents)

	return nil
}

func (s *TicketmasterService) mapToEvent(tmEvent TicketmasterEvent) *models.Event {
	event := &models.Event{
		Title:       tmEvent.Name,
		Description: generateEventDescription(tmEvent),
		EventType:   tmEvent.Type,
		ExternalURL: tmEvent.URL,
		IsExternal:  true,
		ExternalID:  tmEvent.ID,
	}

	// Parse event date with fallbacks to local date/time information
	var eventDate time.Time
	if tmEvent.Dates.Start.DateTime != "" {
		if parsed, err := time.Parse(time.RFC3339, tmEvent.Dates.Start.DateTime); err == nil {
			eventDate = parsed
		}
	}

	if eventDate.IsZero() && tmEvent.Dates.Start.LocalDate != "" {
		layout := "2006-01-02"
		value := tmEvent.Dates.Start.LocalDate
		if tmEvent.Dates.Start.LocalTime != "" {
			layout = "2006-01-02 15:04:05"
			value = fmt.Sprintf("%s %s", tmEvent.Dates.Start.LocalDate, tmEvent.Dates.Start.LocalTime)
		}
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			eventDate = parsed
		}
	}

	if eventDate.IsZero() {
		eventDate = time.Now()
	}
	event.EventDate = &eventDate

	// Add location information if available
	if len(tmEvent.Embedded.Venues) > 0 {
		venue := tmEvent.Embedded.Venues[0]
		event.Location = fmt.Sprintf("%s, %s, %s",
			venue.Name,
			venue.City.Name,
			venue.Country.Name)

		if lat, err := strconv.ParseFloat(venue.Location.Latitude, 64); err == nil {
			event.Latitude = lat
		}
		if lng, err := strconv.ParseFloat(venue.Location.Longitude, 64); err == nil {
			event.Longitude = lng
		}

		if venue.Name != "" {
			event.Organizer = venue.Name
		}
	}

	if event.Organizer == "" {
		event.Organizer = "Ticketmaster"
	}

	return event
}

// StartScheduler triggers Ticketmaster fetch on the provided interval until the context is cancelled.
func (s *TicketmasterService) StartScheduler(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := s.FetchAndSaveEvents(); err != nil {
					log.Printf("Ticketmaster scheduled fetch error: %v", err)
				}
			}
		}
	}()
}

func collectImageLinks(event TicketmasterEvent) []string {
	seen := make(map[string]struct{})
	links := make([]string, 0)

	addLink := func(url string) {
		if url == "" {
			return
		}
		if _, exists := seen[url]; exists {
			return
		}
		seen[url] = struct{}{}
		links = append(links, url)
	}

	// Prefer the highest-resolution image (by height) for each attraction, mirroring NestJS behaviour.
	for _, attraction := range event.Embedded.Attractions {
		if len(attraction.Images) == 0 {
			continue
		}
		best := attraction.Images[0]
		for _, candidate := range attraction.Images[1:] {
			if candidate.Height > best.Height {
				best = candidate
			}
		}
		addLink(best.URL)
	}

	// Fallback to the largest event-level image if no attraction images were added.
	if len(links) == 0 && len(event.Images) > 0 {
		best := event.Images[0]
		for _, candidate := range event.Images[1:] {
			if candidate.Height > best.Height {
				best = candidate
			}
		}
		addLink(best.URL)
	}

	return links
}

func generateEventDescription(event TicketmasterEvent) string {
	var builder strings.Builder

	eventType := "Unknown"
	if len(event.Classifications) > 0 {
		class := event.Classifications[0]
		segment := class.Segment.Name
		genre := class.Genre.Name
		subGenre := class.SubGenre.Name
		eventType = fmt.Sprintf("%s - %s (%s)", segment, genre, subGenre)
	}

	localDate := event.Dates.Start.LocalDate
	if localDate == "" {
		localDate = "Unknown Date"
	}

	localTime := event.Dates.Start.LocalTime
	if localTime == "" {
		localTime = "Unknown Time"
	}

	status := event.Dates.Status.Code
	if status == "" {
		status = "Unknown"
	}

	var venueName, venueAddress string
	if len(event.Embedded.Venues) > 0 {
		venue := event.Embedded.Venues[0]
		venueName = venue.Name
		venueAddress = fmt.Sprintf("%s, %s", venue.City.Name, venue.Country.Name)
	}

	builder.WriteString(fmt.Sprintf("**Event Type:** %s<br />", eventType))
	builder.WriteString(fmt.Sprintf("**Date and Time:** %s at %s<br />", localDate, localTime))
	builder.WriteString(fmt.Sprintf("**Event Status:** %s<br />", status))
	if venueName != "" {
		builder.WriteString(fmt.Sprintf("**Venue:** %s, %s", venueName, venueAddress))
	}

	return builder.String()
}

func randomParticipantCount() int {
	return rand.Intn(81) + 20 // 20-100 inclusive
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
