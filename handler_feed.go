package main

import (
	"context"
	"github.com/Elijah99Lil/rss_aggregator/internal/database"
	"time"
	"github.com/google/uuid"
	"fmt"
	"database/sql"
	"log"
	"strconv"
	"errors"
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

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		return err
	}

	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}

	for _, item := range rssFeed.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time: t,
				Valid: true,
			}
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:			uuid.New(),
			CreatedAt:	time.Now().UTC(),
			UpdatedAt:	time.Now().UTC(),
			FeedID:		feed.ID,
			Title:		item.Title,
			Description: sql.NullString{
				String:	item.Description,
				Valid:	true,
			},
			Url:			item.Link,
			PublishedAt:	publishedAt,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.args[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}