package models

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestUserBeforeCreateGeneratesUUID(t *testing.T) {
	db := newTestDB(t)

	user := &User{Username: "alice", PasswordHash: "hash"}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	if _, err := uuid.Parse(user.ID); err != nil {
		t.Fatalf("expected generated ID to be a UUID, got %q: %v", user.ID, err)
	}
}

func TestUserBeforeCreatePreservesExplicitID(t *testing.T) {
	db := newTestDB(t)

	id := uuid.NewString()
	user := &User{ID: id, Username: "bob", PasswordHash: "hash"}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	if user.ID != id {
		t.Fatalf("expected ID %q to be preserved, got %q", id, user.ID)
	}
}
