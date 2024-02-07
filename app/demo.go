package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"pg-intro/website"
)

func RunDemoRepository(ctx context.Context, repository website.Repository) {
	fmt.Println("1. MIGRATE REPOSITORY")

	if err := repository.Migrate(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("2. CREATE RECORDS OF REPOSITORY")
	gosamples := website.Website{
		Name: "GOSAMPLES",
		URL:  "https://gosamples.dev",
		Rank: 2,
	}

	golang := website.Website{
		Name: "Golang official website",
		URL:  "https:://golang.org",
		Rank: 1,
	}

	createdGosamples, err := repository.Create(ctx, gosamples)
	if errors.Is(err, website.ErrDuplicate) {
		log.Printf("record: %+v already exists.\n", gosamples)
	} else if err != nil {
		log.Fatal(err)
	}

	createdGolang, err := repository.Create(ctx, golang)
	if errors.Is(err, website.ErrDuplicate) {
		log.Printf("record: %+v already exists.\n", golang)
	} else if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n%+v\n", createdGosamples, createdGolang)

	fmt.Println("3. GET RECORD BY NAME")
	gotGosamples, err := repository.GetByName(ctx, "GOSAMPLES")
	if errors.Is(err, website.ErrNotExist) {
		log.Println("record: GOSAMPLES does not exits in the repository")
	} else if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", gotGosamples)

	fmt.Println("4. UPDATE RECORD")
	createdGosamples.Rank = 1
	if _, err := repository.Update(ctx, createdGosamples.ID, *createdGosamples); err != nil {
		if errors.Is(err, website.ErrDuplicate) {
			fmt.Printf("record: %+v already exists\n", createdGosamples)
		} else if errors.Is(err, website.ErrUpdatedFailed) {
			fmt.Printf("update record: %+v failed", createdGosamples)
		} else {
			log.Fatal(err)
		}
	}

	fmt.Println("5. GET ALL")
	websites, err := repository.All(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, website := range websites {
		fmt.Printf("%+v\n", website)
	}

	fmt.Println("6. DELETE RECORD")
	if err := repository.Delete(ctx, createdGolang.ID); err != nil {
		if errors.Is(err, website.ErrDeletedFailed) {
			fmt.Printf("delete of record: %d failed", createdGolang.ID)
		} else {
			log.Fatal(err)
		}
	}

	fmt.Println("7. GET ALL")
	websites, err = repository.All(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, website := range websites {
		fmt.Printf("%+v\n", website)
	}
}
