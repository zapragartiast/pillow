package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"pillow/audit"
	"pillow/auth"
	"pillow/database"
	"pillow/routes"
)

// NOTE: These integration tests require a test PostgreSQL database. Set DATABASE_URL
// to a disposable test DB. Tests will create fixtures (user, role, permission) and
// generate a JWT for an admin user so protected routes can be exercised.

const (
	// Defaults used when running tests locally. These are safe defaults for CI/local
	// but you should replace TEST_DATABASE_URL with your disposable test DB.
	testJWTSecret   = "testsecret_integration"
	testDatabaseURL = "postgres://postgres:postgres@localhost:5432/pillow_test?sslmode=disable"
)

func TestMain(m *testing.M) {
	// Provide sensible defaults for local test runs if env vars not set.
	if os.Getenv("JWT_SECRET") == "" {
		_ = os.Setenv("JWT_SECRET", "testsecret_integration")
	}
	if os.Getenv("DATABASE_URL") == "" {
		_ = os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/pillow_test?sslmode=disable")
	}
	os.Exit(m.Run())
}

func mustDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set; skipping integration tests")
	}
	db := database.Connect(dsn)
	if db == nil {
		t.Fatal("failed to connect to database")
	}
	return db
}

func cleanupAuditTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`DELETE FROM "Audit_Log"`)
	if err != nil {
		t.Fatalf("failed to cleanup Audit_Log: %v", err)
	}
}

// helper to count audit rows by action
func countAuditRows(t *testing.T, db *sql.DB, action string) int {
	t.Helper()
	var n int
	err := db.QueryRow(`SELECT COUNT(*) FROM "Audit_Log" WHERE action = $1`, action).Scan(&n)
	if err != nil {
		t.Fatalf("countAuditRows query failed: %v", err)
	}
	return n
}

// createAdminForTest creates a user, role, permission and assigns manage_users permission to the user.
// Returns userID and JWT token string. Caller should clean up created rows.
func createAdminForTest(t *testing.T, db *sql.DB) (string, string) {
	t.Helper()

	// create permission
	permID := "00000000-0000-0000-0000-000000000101"
	_, err := db.Exec(`INSERT INTO "Permissions" (id, name, description, scope_level) VALUES ($1,$2,$3,$4) ON CONFLICT DO NOTHING`,
		permID, "manage_users", "manage users", "system")
	if err != nil {
		t.Fatalf("failed to ensure permission: %v", err)
	}

	// create role
	roleID := "00000000-0000-0000-0000-000000000201"
	_, err = db.Exec(`INSERT INTO "Roles" (id, name, description) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING`,
		roleID, "test-admin", "test admin")
	if err != nil {
		t.Fatalf("failed to ensure role: %v", err)
	}

	// assign permission to role
	_, err = db.Exec(`INSERT INTO "Role_Permissions" (role_id, permission_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
		roleID, permID)
	if err != nil {
		t.Fatalf("failed to ensure role_permission: %v", err)
	}

	// create user
	userID := "00000000-0000-0000-0000-00000000a001"
	hashed, err := auth.HashPassword("secretpwd")
	if err != nil {
		t.Fatalf("failed to hash password for test admin: %v", err)
	}
	_, err = db.Exec(`INSERT INTO "Users" (id, username, password_hash, email, is_active, created_at) VALUES ($1,$2,$3,$4,true,NOW()) ON CONFLICT DO NOTHING`,
		userID, "integ-admin", hashed, "integ-admin@example.com")
	if err != nil {
		t.Fatalf("failed to create admin user: %v", err)
	}

	// assign role to user
	_, err = db.Exec(`INSERT INTO "User_Roles" (user_id, role_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
		userID, roleID)
	if err != nil {
		t.Fatalf("failed to assign role to user: %v", err)
	}

	// generate JWT for this user
	// generate JWT for this user
	uid, _ := uuid.Parse(userID)
	token, err := auth.GenerateJWT(uid, "integ-admin")
	if err != nil {
		t.Fatalf("failed to generate jwt: %v", err)
	}

	return userID, token
}

func cleanupAdminForTest(t *testing.T, db *sql.DB) {
	t.Helper()
	_, _ = db.Exec(`DELETE FROM "User_Roles" WHERE user_id = '00000000-0000-0000-0000-00000000a001'`)
	_, _ = db.Exec(`DELETE FROM "Users" WHERE id = '00000000-0000-0000-0000-00000000a001'`)
	_, _ = db.Exec(`DELETE FROM "Role_Permissions" WHERE role_id = '00000000-0000-0000-0000-000000000201'`)
	_, _ = db.Exec(`DELETE FROM "Roles" WHERE id = '00000000-0000-0000-0000-000000000201'`)
	_, _ = db.Exec(`DELETE FROM "Permissions" WHERE id = '00000000-0000-0000-0000-000000000101'`)
}

// TestSuccessfulUpdateProducesSingleAuditRow performs an end-to-end request that
// should succeed and verifies exactly one audit row is created for the action.
func TestSuccessfulUpdateProducesSingleAuditRow(t *testing.T) {
	db := mustDB(t)
	defer db.Close()

	// Start global audit queue used by middleware
	audit.StartAuditQueue(db, 50)
	defer audit.StopAuditQueue()

	// Clean audit table
	cleanupAuditTable(t, db)

	// Setup router
	r := routes.SetupRoutes(db)

	// create admin user + JWT
	_, token := createAdminForTest(t, db)
	defer cleanupAdminForTest(t, db)

	// Pre-create a user to update
	userID := "634f0557-5597-4477-9dda-b077a5e286c9"
	_, err := db.Exec(`INSERT INTO "Users" (id, username, password_hash, email, is_active, created_at) VALUES ($1,$2,$3,$4,true,NOW())`,
		userID, "testuser", "fakehash", "testuser@example.com")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	defer db.Exec(`DELETE FROM "Users" WHERE id = $1`, userID)

	// Build request payload to update username
	payload := map[string]interface{}{"username": "zefriks", "email": "zefriks@example.com"}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/api/users/"+userID, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Wait a short while for the queue worker to persist audit
	time.Sleep(500 * time.Millisecond)

	// Expect 200 or 204
	if rr.Code < 200 || rr.Code >= 300 {
		t.Fatalf("expected 2xx response, got %d: body=%s", rr.Code, rr.Body.String())
	}

	// Verify exactly one audit row for action "USER_UPDATED"
	n := countAuditRows(t, db, "USER_UPDATED")
	if n != 1 {
		t.Fatalf("expected 1 audit row for USER_UPDATED, got %d", n)
	}
}

// TestFailedUpdateProducesNoAuditRow ensures that if the handler returns a 4xx/5xx,
// no audit row is created.
func TestFailedUpdateProducesNoAuditRow(t *testing.T) {
	db := mustDB(t)
	defer db.Close()

	audit.StartAuditQueue(db, 50)
	defer audit.StopAuditQueue()

	cleanupAuditTable(t, db)

	r := routes.SetupRoutes(db)

	// create admin user + JWT
	_, token := createAdminForTest(t, db)
	defer cleanupAdminForTest(t, db)

	// Use an invalid UUID to trigger a 400 in UpdateUser handler
	req := httptest.NewRequest(http.MethodPut, "/api/users/invalid-uuid", bytes.NewReader([]byte(`{"username":"x"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Wait a short while for any potential queue activity
	time.Sleep(250 * time.Millisecond)

	// Expect 400
	if rr.Code < 400 || rr.Code >= 500 {
		t.Fatalf("expected 4xx response, got %d", rr.Code)
	}

	// Verify no audit row for method/path fallback action
	// The middleware fallback action uses method+path; ensure no rows exist for this path prefix.
	rows := 0
	err := db.QueryRow(`SELECT COUNT(*) FROM "Audit_Log" WHERE action LIKE 'PUT /api/users/%'`).Scan(&rows)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if rows != 0 {
		t.Fatalf("expected 0 audit rows for failed update, got %d", rows)
	}
}
