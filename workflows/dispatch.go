package workflows

import (
	"context"
	"fmt"
	"time"

	"github.com/IAmRiteshKoushik/gravemind/cmd"
	"github.com/IAmRiteshKoushik/gravemind/db"
)

func DispatchBadge(username string, count int, category string) error {
	var badgeName string
	switch category {
	case "pull-request":
		if count == 1 {
			badgeName = "Hello, Tinkerer"
		} else if count == 5 {
			badgeName = "Ninja Contributor"
		} else if count == 10 {
			badgeName = "Engineer Overclocked"
		} else if count == 20 {
			badgeName = "Doomguy"
		}
	case "bug":
		if count == 1 {
			badgeName = "Sanitizer"
		} else if count == 5 {
			badgeName = "Pest Control"
		} else if count == 10 {
			badgeName = "Planet Cleanser"
		}
	case "doc":
		if count == 2 {
			badgeName = "Doc Champ"
		}
	case "feat":
		if count == 2 {
			badgeName = "High Charity"
		}
	case "help":
		if count == 1 {
			badgeName = "The Scholar"
		} else if count == 3 {
			badgeName = "The Gulliver"
		} else if count == 5 {
			badgeName = "The Oracle"
		}
	case "test":
		if count == 1 {
			badgeName = "Lab Assistant"
		} else if count == 5 {
			badgeName = "Quality Assurer"
		} else if count == 10 {
			badgeName = "Full-Coded Alchemist"
		}
	case "impact":
		badgeName = "Zeppelin of Mighty Gargantuaness (ZOMG)"
	case "stack":
		if count == 3 {
			badgeName = "Polygot"
		} else if count == 5 {
			badgeName = "Jack of All Stacks"
		}
	case "issue":
		badgeName = "Pirate of Issue-bians"
	case "streak":
		badgeName = "Enamoured"
	default:
		return fmt.Errorf("invalid category")
	}

	// No badge awarded
	if badgeName == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	q := db.New()
	tx, err := cmd.DBPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = q.InsertBadgeQuery(ctx, tx, db.InsertBadgeQueryParams{
		Ghusername: username,
		BadgeName:  badgeName,
	})
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
