package main

import (
	"context"
	"fmt"
	"github.com/Elijah99Lil/rss_aggregator/internal/database"
	"strings"
	"time"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("username is required")
	}

	username := cmd.args[0]

	ctx := context.Background()
	_, err := s.db.GetUser(ctx, username)
	if err != nil {
		if strings.Contains(err.Error(), `no rows`) {
			return fmt.Errorf("username does not exist")
		}
		return err
	}
	
	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Println("user has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("must provide a command")
	}

	params := database.CreateUserParams {
		ID: 		uuid.New(),
		CreatedAt: 	time.Now(),
		UpdatedAt: 	time.Now(),
		Name: 		cmd.args[0],
	}

	ctx := context.Background()
	_, err := s.db.CreateUser(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_name_key"`) {
			return fmt.Errorf("username already exists")	
		}
		return err
	}


	err = s.cfg.SetUser(params.Name)
	if err != nil {
		return err
	}

	fmt.Println("user created")
	return nil	
}

func handlerGet(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Println("*", user.Name, "(current)")
		} else {
			fmt.Println("*",user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: go run . agg <time duration like 1s>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)

	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return err
		}
	}
}

func handlerHelp(s *state, cmd command) error {
	fmt.Println(
		`Here are the commands available:
The prefix is always "go run . <command>"

	- login <user>
	- register <user>
	- reset
	- users
	- agg (Use this in a different terminal to auto-scrape feeds in the background)
	- addfeed <url>
	- feeds
	- follow <url>
	- following
	- unfollow <url>
	- browse`,
	)
	return nil
}


