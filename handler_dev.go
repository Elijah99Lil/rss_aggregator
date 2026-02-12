package main

import (
	"context"
)

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()

	err := s.db.DeleteAllUsers(ctx)
	if err != nil {
		return err
	}
	println("Database successfully reset")
	return nil
}