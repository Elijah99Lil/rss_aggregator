package main

import (
	"context"
	"rss_aggregator_mod/internal/database"
	"time"
	"github.com/google/uuid"
	"fmt"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("name and url are required")
	}

	ctx := context.Background()

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

	feedFollowParams := database.CreateFeedFollowParams {
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		UserID:		user.ID,
		FeedID:		feed.ID,
	}

	followFeed, err := s.db.CreateFeedFollow(ctx, feedFollowParams)
	if err != nil {
		return err
	}

	fmt.Println("Feed added: ")
	fmt.Println(feed)

	fmt.Println("Feed follow added: ")
	fmt.Println(followFeed)
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

func handlerFollowFeeds(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	feed, err := s.db.GetFeedByURL(ctx, cmd.args[0])
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams {
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		UserID:		user.ID,
		FeedID:		feed.ID,
	}

	followFeed, err := s.db.CreateFeedFollow(ctx, params)
	if err != nil {
		return err
	}
	fmt.Println(followFeed.FeedName, followFeed.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	feed, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}

	for _, row := range feed {
		fmt.Printf("- %s\n", row.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: unfollow <url>")
	}

	ctx := context.Background()

	feed, err := s.db.GetFeedByURL(ctx, cmd.args[0])
	if err != nil {
		return err
	}

	arg := database.UnfollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.Unfollow(ctx, arg) 
	if err != nil {
		return err
	}

	fmt.Println("Feed unfollowed")
	return nil
}