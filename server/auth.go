package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"server/db"

	"github.com/google/uuid"
)

func authenticate(ctx context.Context, w http.ResponseWriter, token string) bool {
	parsed, err := uuid.Parse(token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return false
	}

	conn, err := getDb(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get db connection: %v", err), http.StatusInternalServerError)
		return false
	}

	defer closeDb(conn, ctx)

	queries := db.New(conn)
	hashed := fmt.Sprintf("%x", sha256.Sum256([]byte(parsed.String())))
	_, err = queries.GetToken(ctx, hashed)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return false
		}

		http.Error(w, fmt.Sprintf("failed to get token: %v", err), http.StatusInternalServerError)

		return false
	}

	return true
}
