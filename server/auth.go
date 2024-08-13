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

func authenticate(ctx context.Context, w http.ResponseWriter, hashedToken string) bool {
	conn, err := getDb(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get db connection: %v", err), http.StatusInternalServerError)
		return false
	}

	defer closeDb(conn, ctx)

	queries := db.New(conn)
	_, err = queries.GetToken(ctx, hashedToken)

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

func hashToken(token string) string {
	parsed, _ := uuid.Parse(token)

	return fmt.Sprintf("%x", sha256.Sum256([]byte(parsed.String())))
}
