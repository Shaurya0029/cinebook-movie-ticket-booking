// Command seed populates the database with demo cities, theaters, movies,
// shows, and seat maps. It is safe to run multiple times: every insert is
// an upsert keyed on the table's natural unique constraint.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"movieticketbooking/internal/config"
	"movieticketbooking/internal/db"
)

type cityDef struct {
	name, state string
	lat, lng    float64
}

type theaterDef struct {
	cityName, name, address string
	lat, lng                float64
}

type movieDef struct {
	title       string
	genres      []string
	language    string
	duration    int
	description string
	posterSlug  string
	status      string // now_showing | coming_soon
	releaseDay  int    // days offset from today (negative = past, positive = future)
}

var cities = []cityDef{
	{"Mumbai", "Maharashtra", 19.0760, 72.8777},
	{"Delhi", "Delhi", 28.7041, 77.1025},
	{"Bengaluru", "Karnataka", 12.9716, 77.5946},
	{"Hyderabad", "Telangana", 17.3850, 78.4867},
	{"Pune", "Maharashtra", 18.5204, 73.8567},
}

var theaters = []theaterDef{
	{"Mumbai", "PVR Phoenix Mills", "Lower Parel, Mumbai", 19.0000, 72.8300},
	{"Mumbai", "INOX R City", "Ghatkopar, Mumbai", 19.0860, 72.9080},
	{"Delhi", "PVR Select Citywalk", "Saket, Delhi", 28.5280, 77.2190},
	{"Delhi", "INOX Nehru Place", "Nehru Place, Delhi", 28.5490, 77.2510},
	{"Bengaluru", "PVR Forum Mall", "Koramangala, Bengaluru", 12.9350, 77.6110},
	{"Bengaluru", "INOX Garuda Mall", "Magrath Road, Bengaluru", 12.9720, 77.6070},
	{"Hyderabad", "PVR Forum Sujana Mall", "Kukatpally, Hyderabad", 17.4930, 78.4000},
	{"Hyderabad", "AMB Cinemas Gachibowli", "Gachibowli, Hyderabad", 17.4400, 78.3490},
	{"Pune", "PVR Phoenix Marketcity", "Viman Nagar, Pune", 18.5620, 73.9170},
	{"Pune", "INOX Bund Garden", "Bund Garden Road, Pune", 18.5390, 73.8830},
}

var movies = []movieDef{
	{"Interstellar Reloaded", []string{"Sci-Fi", "Drama"}, "English", 168, "A crew of explorers travels through a newly discovered wormhole to save humanity.", "interstellar-reloaded", "now_showing", -10},
	{"The Last Heist", []string{"Action", "Thriller"}, "English", 132, "A retired thief is pulled back for one final, impossible job.", "the-last-heist", "now_showing", -7},
	{"Dil Se Dosti", []string{"Romance", "Drama"}, "Hindi", 145, "Two childhood friends discover love was there all along.", "dil-se-dosti", "now_showing", -14},
	{"Shadow Protocol", []string{"Action", "Spy"}, "English", 128, "An intelligence officer races to stop a global conspiracy.", "shadow-protocol", "now_showing", -3},
	{"Chennai Express Nights", []string{"Comedy", "Action"}, "Hindi", 150, "A chaotic train journey turns into the adventure of a lifetime.", "chennai-express-nights", "now_showing", -21},
	{"The Silent Witness", []string{"Mystery", "Thriller"}, "English", 118, "A lone witness holds the key to unraveling a city's darkest secret.", "the-silent-witness", "now_showing", -5},
	{"Kollywood Kings", []string{"Action", "Drama"}, "Tamil", 155, "Rival families clash over power, pride, and legacy.", "kollywood-kings", "now_showing", -30},
	{"Laughter Unlimited", []string{"Comedy"}, "Hindi", 120, "A stand-up comedian's chaotic week before his biggest show.", "laughter-unlimited", "now_showing", -2},
	{"Galaxy's Edge", []string{"Sci-Fi", "Adventure"}, "English", 140, "Explorers venture past the edge of known space.", "galaxys-edge", "coming_soon", 14},
	{"Monsoon Melody", []string{"Romance", "Musical"}, "Hindi", 135, "A musician finds love and inspiration during the monsoon season.", "monsoon-melody", "coming_soon", 21},
	{"The Final Verdict", []string{"Drama"}, "English", 125, "A young lawyer takes on the case that could end her career.", "the-final-verdict", "coming_soon", 28},
	{"Warriors of Kalinga", []string{"Historical", "Action"}, "Hindi", 160, "An epic retelling of one of history's greatest battles.", "warriors-of-kalinga", "coming_soon", 35},
	{"Ocean's Whisper", []string{"Adventure", "Family"}, "English", 110, "A young sailor discovers a hidden world beneath the waves.", "oceans-whisper", "coming_soon", 45},
}

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := db.EnsureSchema(ctx, pool); err != nil {
		log.Fatalf("failed to ensure schema: %v", err)
	}

	cityIDs, err := seedCities(ctx, pool)
	if err != nil {
		log.Fatalf("seed cities: %v", err)
	}
	log.Printf("seeded %d cities", len(cityIDs))

	theaterIDs, err := seedTheaters(ctx, pool, cityIDs)
	if err != nil {
		log.Fatalf("seed theaters: %v", err)
	}
	log.Printf("seeded %d theaters", len(theaterIDs))

	movieIDs, err := seedMovies(ctx, pool)
	if err != nil {
		log.Fatalf("seed movies: %v", err)
	}
	log.Printf("seeded %d movies", len(movieIDs))

	showIDs, err := seedShows(ctx, pool, movieIDs, theaterIDs)
	if err != nil {
		log.Fatalf("seed shows: %v", err)
	}
	log.Printf("seeded %d shows", len(showIDs))

	seatCount, err := seedSeats(ctx, pool, showIDs)
	if err != nil {
		log.Fatalf("seed seats: %v", err)
	}
	log.Printf("seeded %d seats", seatCount)

	log.Println("seed complete")
}

func seedCities(ctx context.Context, pool *pgxpool.Pool) (map[string]int64, error) {
	ids := make(map[string]int64, len(cities))
	for _, c := range cities {
		var id int64
		err := pool.QueryRow(ctx, `
			INSERT INTO cities (name, state, latitude, longitude)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (name) DO UPDATE SET state = EXCLUDED.state
			RETURNING id`, c.name, c.state, c.lat, c.lng).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("city %q: %w", c.name, err)
		}
		ids[c.name] = id
	}
	return ids, nil
}

func seedTheaters(ctx context.Context, pool *pgxpool.Pool, cityIDs map[string]int64) ([]int64, error) {
	var ids []int64
	for _, t := range theaters {
		cityID, ok := cityIDs[t.cityName]
		if !ok {
			return nil, fmt.Errorf("theater %q references unknown city %q", t.name, t.cityName)
		}
		var id int64
		err := pool.QueryRow(ctx, `
			INSERT INTO theaters (city_id, name, address, latitude, longitude)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (city_id, name) DO UPDATE SET address = EXCLUDED.address
			RETURNING id`, cityID, t.name, t.address, t.lat, t.lng).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("theater %q: %w", t.name, err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func seedMovies(ctx context.Context, pool *pgxpool.Pool) (map[string]int64, error) {
	ids := make(map[string]int64, len(movies))
	now := time.Now()
	for _, m := range movies {
		posterURL := fmt.Sprintf("https://picsum.photos/seed/%s/400/600", m.posterSlug)
		releaseDate := now.AddDate(0, 0, m.releaseDay)

		var id int64
		err := pool.QueryRow(ctx, `
			INSERT INTO movies (title, genres, language, duration_minutes, description, poster_url, status, release_date)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (title) DO UPDATE SET status = EXCLUDED.status, release_date = EXCLUDED.release_date
			RETURNING id`,
			m.title, m.genres, m.language, m.duration, m.description, posterURL, m.status, releaseDate).
			Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("movie %q: %w", m.title, err)
		}
		ids[m.title] = id
	}
	return ids, nil
}

// seedShows generates showtimes for every now_showing movie across a rotating
// pair of theaters, over the next few days, at a handful of times per day.
func seedShows(ctx context.Context, pool *pgxpool.Pool, movieIDs map[string]int64, theaterIDs []int64) ([]int64, error) {
	showTimes := []string{"13:00", "17:00", "20:30"}
	daysAhead := 5
	pricesCents := map[string]int{"13:00": 15000, "17:00": 22000, "20:30": 30000}

	var showIDs []int64
	theaterCursor := 0

	for _, m := range movies {
		if m.status != "now_showing" {
			continue
		}
		movieID := movieIDs[m.title]

		// Each movie plays at 2 theaters, picked round-robin across the seeded list.
		movieTheaters := []int64{
			theaterIDs[theaterCursor%len(theaterIDs)],
			theaterIDs[(theaterCursor+1)%len(theaterIDs)],
		}
		theaterCursor += 2

		for dayOffset := 0; dayOffset < daysAhead; dayOffset++ {
			day := time.Now().AddDate(0, 0, dayOffset)
			for _, theaterID := range movieTheaters {
				for _, t := range showTimes {
					startsAt, err := time.ParseInLocation("2006-01-02 15:04", day.Format("2006-01-02")+" "+t, time.Local)
					if err != nil {
						return nil, err
					}

					var showID int64
					err = pool.QueryRow(ctx, `
						INSERT INTO shows (movie_id, theater_id, starts_at, price_cents)
						VALUES ($1, $2, $3, $4)
						ON CONFLICT (movie_id, theater_id, starts_at) DO UPDATE SET price_cents = EXCLUDED.price_cents
						RETURNING id`,
						movieID, theaterID, startsAt, pricesCents[t]).Scan(&showID)
					if err != nil {
						return nil, fmt.Errorf("show for %q: %w", m.title, err)
					}
					showIDs = append(showIDs, showID)
				}
			}
		}
	}
	return showIDs, nil
}

// seedSeats generates an 8x10 seat grid (A1..H10) for every show.
func seedSeats(ctx context.Context, pool *pgxpool.Pool, showIDs []int64) (int, error) {
	rows := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	count := 0

	for _, showID := range showIDs {
		for _, row := range rows {
			for col := 1; col <= 10; col++ {
				label := fmt.Sprintf("%s%d", row, col)
				_, err := pool.Exec(ctx, `
					INSERT INTO seats (show_id, seat_label, is_booked)
					VALUES ($1, $2, false)
					ON CONFLICT (show_id, seat_label) DO NOTHING`, showID, label)
				if err != nil {
					return count, fmt.Errorf("seat %s for show %d: %w", label, showID, err)
				}
				count++
			}
		}
	}
	return count, nil
}
