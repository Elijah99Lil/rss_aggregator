package main

import (
	"context"
	"rss_aggregator_mod/internal/database"
	"time"
	"github.com/google/uuid"
	"fmt"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("name and url are required")
	}

	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	params := database.CreateFeedParams {
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		Name:		cmd.args[0],	
		Url:		cmd.args[1],
		UserID:		user.ID,
	}

	feed, err := s.db.CreateFeed(ctx, params)
	if err != nil {
		return err
	}

	fmt.Println("Feed added: ")
	fmt.Println(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()

	feedsWithUsers, err := s.db.GetFeedsWithUsers(ctx)
	if err != nil {
		return err
	}
	
	for _, row := range feedsWithUsers {
		fmt.Printf("- %s\n- %s\n", row.Name, row.UserName)
	}
	return nil
}