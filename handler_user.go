package main

import (
	"context"
	"fmt"
	"rss_aggregator_mod/internal/database"
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
	ctx := context.Background()
	feedUrl := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(ctx, feedUrl)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", feed)
	return nil
}


