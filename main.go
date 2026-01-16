package main

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nfnt/resize"
	"golang.org/x/crypto/bcrypt"
)

type Entry struct {
	ID            int
	Title         string
	Content       string
	PhotoPath     string
	MediaType     string
	ThumbnailPath string
	Slug          string
	CreatedAt     time.Time
	TimeAgo       string
}

type EntryDisplay struct {
	ID            int
	Title         string
	Content       template.HTML
	FullContent   template.HTML
	IsTruncated   bool
	Photo         template.URL
	HasPhoto      bool
	MediaType     string
	HasAudio      bool
	HasVideo      bool
	Thumbnail     template.URL
	HasThumbnail  bool
	Slug          string
	CreatedAt     time.Time
	TimeAgo       string
	InitialLetter string
}

type ViewerPageData struct {
	Entries          []EntryDisplay
	TotalEntries     int
	TodayEntries     int
	Last24Hours      int
	SiteTitle        string
	SiteSubtitle     string
	EnableSubtitle   bool
	UserInitial      string
	AvatarPath       string
	AvatarPreference string
	InitialCount     int
	HasMore          bool
	ThemeCSS         template.CSS
}

type EditorPageData struct {
	Message            string
	MessageType        string
	Entries            []Entry
	HasPrivacyPassword bool
	EnableAudioUploads bool
	// New fields for sidebar layout and pagination
	View         string
	PageTitle    string
	CurrentPage  int
	TotalPages   int
	TotalEntries int
	StartEntry   int
	EndEntry     int
	PageNumbers  []int
}

type SinglePostPageData struct {
	Entry            EntryDisplay
	SiteTitle        string
	SiteSubtitle     string
	EnableSubtitle   bool
	UserInitial      string
	AvatarPath       string
	AvatarPreference string
	ThemeCSS         template.CSS
}

type SiteSettings struct {
	SiteTitle          string
	SiteSubtitle       string
	UserInitial        string
	AvatarPath         string
	AvatarPreference   string // "avatar" or "initials"
	SiteTheme          string
	CustomBgColor      string
	CustomTextColor    string
	CustomAccentColor  string
	HasViewerPassword  bool
	HasAdminPassword   bool
}

type SettingsPageData struct {
	Message               string
	MessageType           string
	Settings              SiteSettings
	EnableSubtitle        bool
	CustomDomain          *CustomDomain
	InstanceHostname      string
	CustomDomainEnabled   bool
	CanEnableCustomDomain bool
	View                  string
	PageTitle             string
}

type NotFoundPageData struct {
	SiteTitle      string
	SiteSubtitle   string
	EnableSubtitle bool
	UserInitial    string
	RecentEntries  []EntryDisplay
	ThemeCSS       template.CSS
}

type APIResponse struct {
	Entries []EntryDisplay `json:"entries"`
	HasMore bool           `json:"hasMore"`
}

type EntriesResponse struct {
	Entries []EntryJSON `json:"entries"`
	HasMore bool        `json:"hasMore"`
}

type EntryJSON struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Photo     string `json:"photo"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"createdAt"`
}

type Session struct {
	Token     string
	ExpiresAt time.Time
}

// CustomDomain represents a custom domain configuration
type CustomDomain struct {
	ID                   int
	Domain               string
	VerificationToken    string
	VerifiedAt           sql.NullTime
	ActivatedAt          sql.NullTime
	LastVerifiedAt       sql.NullTime
	VerificationAttempts int
	CreatedAt            time.Time
}

var db *sql.DB
var urlRegex *regexp.Regexp
var uploadsDir string
var enableAudioUploads bool
var enableSubtitle bool
var sessions = make(map[string]*Session)
var sessionMutex sync.RWMutex

// Rate limiting for domain verification
var domainVerifyAttempts = make(map[string][]time.Time)
var domainVerifyMutex sync.Mutex

// Hostname auto-detection
var hostnameDetected sync.Once

func init() {
	urlRegex = regexp.MustCompile(`https?://[^\s<>"{}|\\^\[\]` + "`" + `]+`)
}

// detectHostnameFromRequest extracts and stores the hostname from the first HTTP request
// Only stores hostname if it's a *.postastiq.com subdomain (enables custom domain feature)
func detectHostnameFromRequest(r *http.Request) {
	hostnameDetected.Do(func() {
		// Check if hostname is already set
		currentHostname := getInstanceHostname()
		if currentHostname != "" {
			return
		}

		// Extract hostname from Host header (strip port if present)
		host := r.Host
		if colonIdx := strings.LastIndex(host, ":"); colonIdx != -1 {
			// Check if it's not an IPv6 address
			if !strings.Contains(host, "]") || colonIdx > strings.LastIndex(host, "]") {
				host = host[:colonIdx]
			}
		}

		if host == "" || host == "localhost" || net.ParseIP(host) != nil {
			return
		}

		// Only store postastiq.com subdomains (custom domain feature only for managed instances)
		if !isPostastiqSubdomain(host) {
			return
		}

		// Store the detected hostname
		err := setInstanceHostname(host)
		if err != nil {
			log.Printf("Failed to store hostname: %v", err)
			return
		}
		log.Printf("Detected postastiq.com subdomain: %s (custom domains enabled)", host)
	})
}

// isPostastiqSubdomain checks if the hostname is a *.postastiq.com subdomain
func isPostastiqSubdomain(hostname string) bool {
	hostname = strings.ToLower(strings.TrimSpace(hostname))
	return strings.HasSuffix(hostname, ".postastiq.com")
}

// getHostFromRequest extracts the hostname from the request's Host header (strips port)
func getHostFromRequest(r *http.Request) string {
	host := r.Host
	if colonIdx := strings.LastIndex(host, ":"); colonIdx != -1 {
		// Check if it's not an IPv6 address
		if !strings.Contains(host, "]") || colonIdx > strings.LastIndex(host, "]") {
			host = host[:colonIdx]
		}
	}
	return strings.ToLower(strings.TrimSpace(host))
}

func initDB() error {
	var err error

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/app/data/blog.db"
	}

	uploadsDir = os.Getenv("UPLOADS_DIR")
	if uploadsDir == "" {
		uploadsDir = "/app/data/uploads"
	}

	// Audio uploads disabled by default, enable via environment variable
	enableAudioUploads = os.Getenv("ENABLE_AUDIO_UPLOADS") == "true"

	// Subtitle disabled by default, enable via environment variable
	enableSubtitle = os.Getenv("ENABLE_SUBTITLE") == "true"

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return fmt.Errorf("failed to create uploads directory: %v", err)
	}

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		content TEXT NOT NULL,
		photo_path TEXT,
		media_type TEXT DEFAULT 'photo',
		slug TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	// Create unique index on slug column for new tables
	_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_entries_slug ON entries(slug)`)
	if err != nil {
		log.Printf("Warning: failed to create unique index on slug: %v", err)
	}

	// Check if we need to migrate from old schema (photo BLOB) to new schema (photo_path TEXT)
	var columnExists bool
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('entries') WHERE name='photo'").Scan(&columnExists)
	if err == nil && columnExists {
		log.Println("Migrating database schema from photo BLOB to photo_path TEXT...")

		// Create new table with correct schema
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS entries_new (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				title TEXT,
				content TEXT NOT NULL,
				photo_path TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)
		`)
		if err != nil {
			return fmt.Errorf("failed to create new table: %v", err)
		}

		// Copy data (drop photo BLOB column as we can't convert it)
		_, err = db.Exec(`
			INSERT INTO entries_new (id, content, created_at)
			SELECT id, content, created_at FROM entries
		`)
		if err != nil {
			return fmt.Errorf("failed to copy data: %v", err)
		}

		// Drop old table and rename new one
		_, err = db.Exec(`DROP TABLE entries`)
		if err != nil {
			return fmt.Errorf("failed to drop old table: %v", err)
		}

		_, err = db.Exec(`ALTER TABLE entries_new RENAME TO entries`)
		if err != nil {
			return fmt.Errorf("failed to rename table: %v", err)
		}

		log.Println("Database migration completed successfully")
	}

	// Check if we need to add title column to existing table
	var titleColumnExists bool
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('entries') WHERE name='title'").Scan(&titleColumnExists)
	if err == nil && !titleColumnExists {
		log.Println("Adding title column to entries table...")
		_, err = db.Exec(`ALTER TABLE entries ADD COLUMN title TEXT`)
		if err != nil {
			return fmt.Errorf("failed to add title column: %v", err)
		}
		log.Println("Title column added successfully")
	}

	// Check if we need to add slug column to existing table
	var slugColumnExists bool
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('entries') WHERE name='slug'").Scan(&slugColumnExists)
	if err == nil && !slugColumnExists {
		log.Println("Adding slug column to entries table...")
		// Add column without UNIQUE constraint first (SQLite limitation)
		_, err = db.Exec(`ALTER TABLE entries ADD COLUMN slug TEXT`)
		if err != nil {
			return fmt.Errorf("failed to add slug column: %v", err)
		}
		log.Println("Slug column added successfully")

		// Generate slugs for existing entries
		log.Println("Generating slugs for existing entries...")
		rows, err := db.Query("SELECT id, title, content, created_at FROM entries WHERE slug IS NULL OR slug = ''")
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var id int
				var title sql.NullString
				var content string
				var createdAt time.Time
				if err := rows.Scan(&id, &title, &content, &createdAt); err == nil {
					// Apply title logic for migration
					finalTitle := strings.TrimSpace(title.String)
					if finalTitle == "" && content != "" {
						// Extract first line, trim to 60 chars at last full word
						firstLine := strings.Split(content, "\n")[0]
						if len(firstLine) > 60 {
							truncated := firstLine[:60]
							lastSpace := strings.LastIndex(truncated, " ")
							if lastSpace > 0 {
								finalTitle = truncated[:lastSpace]
							} else {
								finalTitle = truncated
							}
						} else {
							finalTitle = firstLine
						}
					} else if finalTitle == "" && content == "" {
						finalTitle = fmt.Sprintf("Untitled Post %s", createdAt.Format("2006-01-02"))
					}
					if len(finalTitle) > 80 {
						finalTitle = finalTitle[:80]
					}

					slug := generateSlug(finalTitle, createdAt)
					// Update both title and slug for consistency
					_, _ = db.Exec("UPDATE entries SET title = ?, slug = ? WHERE id = ?", finalTitle, slug, id)
				}
			}
		}
		log.Println("Slug generation completed")

		// Create unique index on slug column
		log.Println("Creating unique index on slug column...")
		_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_entries_slug ON entries(slug)`)
		if err != nil {
			log.Printf("Warning: failed to create unique index on slug: %v", err)
		} else {
			log.Println("Unique index created successfully")
		}
	}

	// Check if we need to add media_type column to existing table
	var mediaTypeColumnExists bool
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('entries') WHERE name='media_type'").Scan(&mediaTypeColumnExists)
	if err == nil && !mediaTypeColumnExists {
		log.Println("Adding media_type column to entries table...")
		_, err = db.Exec(`ALTER TABLE entries ADD COLUMN media_type TEXT DEFAULT 'photo'`)
		if err != nil {
			return fmt.Errorf("failed to add media_type column: %v", err)
		}
		log.Println("Media_type column added successfully")
	}

	// Check if we need to add thumbnail_path column to existing table
	var thumbnailColumnExists bool
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('entries') WHERE name='thumbnail_path'").Scan(&thumbnailColumnExists)
	if err == nil && !thumbnailColumnExists {
		log.Println("Adding thumbnail_path column to entries table...")
		_, err = db.Exec(`ALTER TABLE entries ADD COLUMN thumbnail_path TEXT`)
		if err != nil {
			return fmt.Errorf("failed to add thumbnail_path column: %v", err)
		}
		log.Println("Thumbnail_path column added successfully")
	}

	// Regenerate slugs for entries with empty or date-based slugs (migration from old format)
	log.Println("Checking for entries that need slug regeneration...")

	// First, drop the unique index temporarily to allow updates
	_, _ = db.Exec(`DROP INDEX IF EXISTS idx_entries_slug`)

	rows, err := db.Query("SELECT id, title, content, created_at FROM entries")
	if err == nil {
		// Read all entries first
		type entryData struct {
			id        int
			title     string
			content   string
			createdAt time.Time
		}
		var entriesToUpdate []entryData

		for rows.Next() {
			var id int
			var title sql.NullString
			var content string
			var createdAt time.Time
			if err := rows.Scan(&id, &title, &content, &createdAt); err == nil {
				entriesToUpdate = append(entriesToUpdate, entryData{
					id:        id,
					title:     title.String,
					content:   content,
					createdAt: createdAt,
				})
			}
		}
		rows.Close()

		// Now update all entries
		updateCount := 0
		for _, entry := range entriesToUpdate {
			// Apply title logic
			finalTitle := strings.TrimSpace(entry.title)
			if finalTitle == "" && entry.content != "" {
				// Extract first line, trim to 60 chars at last full word
				firstLine := strings.Split(entry.content, "\n")[0]
				if len(firstLine) > 60 {
					truncated := firstLine[:60]
					lastSpace := strings.LastIndex(truncated, " ")
					if lastSpace > 0 {
						finalTitle = truncated[:lastSpace]
					} else {
						finalTitle = truncated
					}
				} else {
					finalTitle = firstLine
				}
			} else if finalTitle == "" && entry.content == "" {
				finalTitle = fmt.Sprintf("Untitled Post %s", entry.createdAt.Format("2006-01-02"))
			}
			if len(finalTitle) > 80 {
				finalTitle = finalTitle[:80]
			}

			slug := generateSlug(finalTitle, entry.createdAt)
			_, _ = db.Exec("UPDATE entries SET title = ?, slug = ? WHERE id = ?", finalTitle, slug, entry.id)
			updateCount++
		}
		if updateCount > 0 {
			log.Printf("Regenerated slugs for %d entries", updateCount)
		}

		// Recreate the unique index after updates
		log.Println("Recreating unique index on slug column...")
		_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_entries_slug ON entries(slug)`)
		if err != nil {
			log.Printf("Warning: failed to create unique index on slug: %v", err)
		} else {
			log.Println("Unique index recreated successfully")
		}
	}

	// Create privacy settings table
	privacyTableQuery := `
	CREATE TABLE IF NOT EXISTS privacy_settings (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		password_hash TEXT
	)`
	_, err = db.Exec(privacyTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create privacy_settings table: %v", err)
	}

	// Create site settings table
	siteSettingsTableQuery := `
	CREATE TABLE IF NOT EXISTS site_settings (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		site_title TEXT DEFAULT 'My Blog',
		site_subtitle TEXT DEFAULT 'A Personal Blog',
		user_initial TEXT DEFAULT 'AB',
		site_theme TEXT DEFAULT 'default',
		admin_password_hash TEXT,
		viewer_password_hash TEXT
	)`
	_, err = db.Exec(siteSettingsTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create site_settings table: %v", err)
	}

	// Initialize settings with defaults if not exists
	_, err = db.Exec(`
		INSERT OR IGNORE INTO site_settings (id, site_title, site_subtitle, user_initial, site_theme)
		VALUES (1, 'My Blog', 'A Personal Blog', 'AB', 'default')
	`)
	if err != nil {
		return fmt.Errorf("failed to initialize site_settings: %v", err)
	}

	// Add avatar_path column if it doesn't exist
	var avatarColumnExists bool
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('site_settings') WHERE name='avatar_path'").Scan(&avatarColumnExists)
	if err == nil && !avatarColumnExists {
		log.Println("Adding avatar_path column to site_settings table...")
		_, err = db.Exec(`ALTER TABLE site_settings ADD COLUMN avatar_path TEXT`)
		if err != nil {
			return fmt.Errorf("failed to add avatar_path column: %v", err)
		}
		log.Println("Avatar_path column added successfully")
	}

	// Add avatar_preference column if it doesn't exist
	var avatarPrefColumnExists bool
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('site_settings') WHERE name='avatar_preference'").Scan(&avatarPrefColumnExists)
	if err == nil && !avatarPrefColumnExists {
		log.Println("Adding avatar_preference column to site_settings table...")
		_, err = db.Exec(`ALTER TABLE site_settings ADD COLUMN avatar_preference TEXT DEFAULT 'initials'`)
		if err != nil {
			return fmt.Errorf("failed to add avatar_preference column: %v", err)
		}
		log.Println("Avatar_preference column added successfully")
	}

	// Add password_change_required column if it doesn't exist
	var passwordChangeColumnExists bool
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('site_settings') WHERE name='password_change_required'").Scan(&passwordChangeColumnExists)
	if err == nil && !passwordChangeColumnExists {
		log.Println("Adding password_change_required column to site_settings table...")
		_, err = db.Exec(`ALTER TABLE site_settings ADD COLUMN password_change_required INTEGER DEFAULT 0`)
		if err != nil {
			return fmt.Errorf("failed to add password_change_required column: %v", err)
		}
		log.Println("Password_change_required column added successfully")
	}

	// Handle initial admin password - use ADMIN_PASSWORD env var or default to "admin"
	// Always require password change on first login for security
	var existingHash sql.NullString
	err = db.QueryRow("SELECT admin_password_hash FROM site_settings WHERE id = 1").Scan(&existingHash)
	if err == nil && (!existingHash.Valid || existingHash.String == "") {
		// No password set yet - this is a fresh installation
		initialAdminPassword := os.Getenv("ADMIN_PASSWORD")
		if initialAdminPassword == "" {
			initialAdminPassword = "admin" // Default password
			log.Println("No ADMIN_PASSWORD set, using default password 'admin'")
		} else {
			log.Println("Using ADMIN_PASSWORD from environment variable")
		}

		// Hash and store the initial password
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(initialAdminPassword), bcrypt.DefaultCost)
		if hashErr != nil {
			return fmt.Errorf("failed to hash initial admin password: %v", hashErr)
		}
		_, err = db.Exec("UPDATE site_settings SET admin_password_hash = ?, password_change_required = 1 WHERE id = 1", string(hashedPassword))
		if err != nil {
			return fmt.Errorf("failed to store initial admin password: %v", err)
		}
		log.Println("Initial admin password configured - password change required on first login")
	}

	// Add instance_hostname column to site_settings if it doesn't exist
	var instanceHostnameExists bool
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('site_settings') WHERE name='instance_hostname'").Scan(&instanceHostnameExists)
	if err == nil && !instanceHostnameExists {
		log.Println("Adding instance_hostname column to site_settings table...")
		_, err = db.Exec(`ALTER TABLE site_settings ADD COLUMN instance_hostname TEXT`)
		if err != nil {
			return fmt.Errorf("failed to add instance_hostname column: %v", err)
		}
		log.Println("Instance_hostname column added successfully")
	}

	// Add custom theme color columns if they don't exist
	customColorColumns := []struct {
		name         string
		defaultValue string
	}{
		{"custom_bg_color", "#fafafa"},
		{"custom_text_color", "#262626"},
		{"custom_accent_color", "#0095f6"},
	}
	for _, col := range customColorColumns {
		var colExists bool
		err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('site_settings') WHERE name=?", col.name).Scan(&colExists)
		if err == nil && !colExists {
			log.Printf("Adding %s column to site_settings table...", col.name)
			_, err = db.Exec(fmt.Sprintf(`ALTER TABLE site_settings ADD COLUMN %s TEXT DEFAULT '%s'`, col.name, col.defaultValue))
			if err != nil {
				return fmt.Errorf("failed to add %s column: %v", col.name, err)
			}
			log.Printf("%s column added successfully", col.name)
		}
	}


	// Create custom_domain table for custom domain management
	customDomainTableQuery := `
	CREATE TABLE IF NOT EXISTS custom_domain (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain TEXT UNIQUE NOT NULL,
		verification_token TEXT NOT NULL,
		verified_at DATETIME DEFAULT NULL,
		activated_at DATETIME DEFAULT NULL,
		last_verified_at DATETIME DEFAULT NULL,
		verification_attempts INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	_, err = db.Exec(customDomainTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create custom_domain table: %v", err)
	}

	log.Println("Database connection established and table created")
	log.Printf("Uploads directory: %s", uploadsDir)
	return nil
}

// runDatabaseMigrations adds any missing columns to the database schema
// This is called both at startup and after restoring a backup
func runDatabaseMigrations() error {
	// Add avatar_path column if it doesn't exist
	var avatarColumnExists bool
	err := db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('site_settings') WHERE name='avatar_path'").Scan(&avatarColumnExists)
	if err == nil && !avatarColumnExists {
		log.Println("Migration: Adding avatar_path column...")
		_, err = db.Exec(`ALTER TABLE site_settings ADD COLUMN avatar_path TEXT`)
		if err != nil {
			return fmt.Errorf("failed to add avatar_path column: %v", err)
		}
	}

	// Add avatar_preference column if it doesn't exist
	var avatarPrefColumnExists bool
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('site_settings') WHERE name='avatar_preference'").Scan(&avatarPrefColumnExists)
	if err == nil && !avatarPrefColumnExists {
		log.Println("Migration: Adding avatar_preference column...")
		_, err = db.Exec(`ALTER TABLE site_settings ADD COLUMN avatar_preference TEXT DEFAULT 'initials'`)
		if err != nil {
			return fmt.Errorf("failed to add avatar_preference column: %v", err)
		}
	}

	// Add password_change_required column if it doesn't exist
	var passwordChangeColumnExists bool
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('site_settings') WHERE name='password_change_required'").Scan(&passwordChangeColumnExists)
	if err == nil && !passwordChangeColumnExists {
		log.Println("Migration: Adding password_change_required column...")
		_, err = db.Exec(`ALTER TABLE site_settings ADD COLUMN password_change_required INTEGER DEFAULT 0`)
		if err != nil {
			return fmt.Errorf("failed to add password_change_required column: %v", err)
		}
	}

	// Add instance_hostname column if it doesn't exist
	var instanceHostnameExists bool
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('site_settings') WHERE name='instance_hostname'").Scan(&instanceHostnameExists)
	if err == nil && !instanceHostnameExists {
		log.Println("Migration: Adding instance_hostname column...")
		_, err = db.Exec(`ALTER TABLE site_settings ADD COLUMN instance_hostname TEXT`)
		if err != nil {
			return fmt.Errorf("failed to add instance_hostname column: %v", err)
		}
	}

	// Add custom theme color columns if they don't exist
	customColorColumns := []struct {
		name         string
		defaultValue string
	}{
		{"custom_bg_color", "#fafafa"},
		{"custom_text_color", "#262626"},
		{"custom_accent_color", "#0095f6"},
	}
	for _, col := range customColorColumns {
		var colExists bool
		err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('site_settings') WHERE name=?", col.name).Scan(&colExists)
		if err == nil && !colExists {
			log.Printf("Migration: Adding %s column...", col.name)
			_, err = db.Exec(fmt.Sprintf(`ALTER TABLE site_settings ADD COLUMN %s TEXT DEFAULT '%s'`, col.name, col.defaultValue))
			if err != nil {
				return fmt.Errorf("failed to add %s column: %v", col.name, err)
			}
		}
	}

	return nil
}

// ============================================================================
// Custom Domain Functions
// ============================================================================

// validateDomainInput validates a domain name input
func validateDomainInput(domain string) error {
	domain = strings.TrimSpace(strings.ToLower(domain))

	if domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}

	// Must have at least one dot (no TLDs)
	if !strings.Contains(domain, ".") {
		return fmt.Errorf("must be a fully qualified domain name")
	}

	// No wildcards
	if strings.Contains(domain, "*") {
		return fmt.Errorf("wildcard domains not allowed")
	}

	// No IP addresses
	if net.ParseIP(domain) != nil {
		return fmt.Errorf("IP addresses not allowed")
	}

	// Block system domains
	blockedPatterns := []string{
		"postastiq.com",
		"postastiq.io",
		"localhost",
		"127.0.0.1",
		"0.0.0.0",
		"example.com",
		"example.org",
	}
	for _, pattern := range blockedPatterns {
		if strings.Contains(domain, pattern) {
			return fmt.Errorf("this domain is not allowed")
		}
	}

	// Basic hostname validation
	parts := strings.Split(domain, ".")
	for _, part := range parts {
		if len(part) == 0 || len(part) > 63 {
			return fmt.Errorf("invalid domain format")
		}
		// Check for valid characters (alphanumeric and hyphen)
		for i, c := range part {
			if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
				return fmt.Errorf("invalid characters in domain")
			}
			// Cannot start or end with hyphen
			if c == '-' && (i == 0 || i == len(part)-1) {
				return fmt.Errorf("domain parts cannot start or end with hyphen")
			}
		}
	}

	return nil
}

// generateVerificationToken generates a cryptographically secure token
func generateVerificationToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "pstq_" + base64.URLEncoding.EncodeToString(b), nil
}

// canAttemptDomainVerification checks rate limiting for domain verification
func canAttemptDomainVerification(domain string) bool {
	domainVerifyMutex.Lock()
	defer domainVerifyMutex.Unlock()

	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)

	// Filter to only recent attempts
	var recent []time.Time
	for _, t := range domainVerifyAttempts[domain] {
		if t.After(oneHourAgo) {
			recent = append(recent, t)
		}
	}
	domainVerifyAttempts[domain] = recent

	return len(recent) < 5
}

// recordDomainVerificationAttempt records a verification attempt for rate limiting
func recordDomainVerificationAttempt(domain string) {
	domainVerifyMutex.Lock()
	defer domainVerifyMutex.Unlock()

	domainVerifyAttempts[domain] = append(domainVerifyAttempts[domain], time.Now())
}

// getInstanceHostname retrieves the instance hostname from the database
func getInstanceHostname() string {
	var hostname sql.NullString
	err := db.QueryRow("SELECT instance_hostname FROM site_settings WHERE id = 1").Scan(&hostname)
	if err != nil || !hostname.Valid || hostname.String == "" {
		return ""
	}
	return hostname.String
}

// setInstanceHostname sets the instance hostname in the database
func setInstanceHostname(hostname string) error {
	_, err := db.Exec("UPDATE site_settings SET instance_hostname = ? WHERE id = 1", hostname)
	return err
}

// verifyDomainDNS verifies both TXT and CNAME records for a domain
func verifyDomainDNS(domain, expectedToken string) error {
	// Verify TXT record
	txtHost := "_postastiq-verify." + domain
	txtRecords, err := net.LookupTXT(txtHost)
	if err != nil {
		return fmt.Errorf("TXT record not found: please add a TXT record for %s", txtHost)
	}

	found := false
	for _, txt := range txtRecords {
		if txt == expectedToken {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("verification token mismatch: the TXT record value does not match")
	}

	// Verify CNAME record
	instanceHostname := getInstanceHostname()
	if instanceHostname == "" {
		return fmt.Errorf("instance hostname not configured")
	}

	cname, err := net.LookupCNAME(domain)
	if err != nil {
		return fmt.Errorf("CNAME record not found: please add a CNAME record pointing to %s", instanceHostname)
	}

	// Normalize CNAME (remove trailing dot if present)
	cname = strings.TrimSuffix(cname, ".")
	expectedCNAME := strings.TrimSuffix(instanceHostname, ".")

	if cname != expectedCNAME {
		return fmt.Errorf("CNAME does not point to this instance: expected %s, got %s", expectedCNAME, cname)
	}

	return nil
}

// getCustomDomain retrieves the current custom domain from the database
func getCustomDomain() (*CustomDomain, error) {
	var cd CustomDomain
	err := db.QueryRow(`
		SELECT id, domain, verification_token, verified_at, activated_at,
		       last_verified_at, verification_attempts, created_at
		FROM custom_domain
		LIMIT 1
	`).Scan(&cd.ID, &cd.Domain, &cd.VerificationToken, &cd.VerifiedAt,
		&cd.ActivatedAt, &cd.LastVerifiedAt, &cd.VerificationAttempts, &cd.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cd, nil
}

// addCustomDomain adds a new custom domain with a verification token
func addCustomDomain(domain string) (*CustomDomain, error) {
	// Validate domain
	if err := validateDomainInput(domain); err != nil {
		return nil, err
	}

	// Check if domain already exists
	existing, err := getCustomDomain()
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("a custom domain is already configured; remove it first")
	}

	// Generate verification token
	token, err := generateVerificationToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %v", err)
	}

	// Insert into database
	result, err := db.Exec(`
		INSERT INTO custom_domain (domain, verification_token) VALUES (?, ?)
	`, strings.ToLower(domain), token)
	if err != nil {
		return nil, fmt.Errorf("failed to save domain: %v", err)
	}

	id, _ := result.LastInsertId()
	return &CustomDomain{
		ID:                int(id),
		Domain:            strings.ToLower(domain),
		VerificationToken: token,
		CreatedAt:         time.Now(),
	}, nil
}

// verifyCustomDomain verifies DNS records and marks domain as verified
func verifyCustomDomain() error {
	cd, err := getCustomDomain()
	if err != nil {
		return err
	}
	if cd == nil {
		return fmt.Errorf("no custom domain configured")
	}

	// Check rate limiting
	if !canAttemptDomainVerification(cd.Domain) {
		return fmt.Errorf("too many verification attempts; please wait an hour")
	}
	recordDomainVerificationAttempt(cd.Domain)

	// Increment attempt counter in database
	_, err = db.Exec("UPDATE custom_domain SET verification_attempts = verification_attempts + 1 WHERE id = ?", cd.ID)
	if err != nil {
		log.Printf("Warning: failed to update verification attempts: %v", err)
	}

	// Verify DNS records
	if err := verifyDomainDNS(cd.Domain, cd.VerificationToken); err != nil {
		return err
	}

	// Mark as verified
	_, err = db.Exec("UPDATE custom_domain SET verified_at = CURRENT_TIMESTAMP WHERE id = ?", cd.ID)
	if err != nil {
		return fmt.Errorf("failed to update verification status: %v", err)
	}

	return nil
}

// activateCustomDomain activates the custom domain by updating Caddy config
func activateCustomDomain() error {
	cd, err := getCustomDomain()
	if err != nil {
		return err
	}
	if cd == nil {
		return fmt.Errorf("no custom domain configured")
	}
	if !cd.VerifiedAt.Valid {
		return fmt.Errorf("domain must be verified before activation")
	}

	// Re-verify DNS before activation
	if err := verifyDomainDNS(cd.Domain, cd.VerificationToken); err != nil {
		return fmt.Errorf("DNS verification failed: %v", err)
	}

	// Update Caddy configuration
	if err := updateCaddyConfig(cd.Domain); err != nil {
		return fmt.Errorf("failed to update Caddy config: %v", err)
	}

	// Mark as activated
	_, err = db.Exec("UPDATE custom_domain SET activated_at = CURRENT_TIMESTAMP, last_verified_at = CURRENT_TIMESTAMP WHERE id = ?", cd.ID)
	if err != nil {
		return fmt.Errorf("failed to update activation status: %v", err)
	}

	return nil
}

// removeCustomDomain removes the custom domain configuration
func removeCustomDomain() error {
	cd, err := getCustomDomain()
	if err != nil {
		return err
	}
	if cd == nil {
		return fmt.Errorf("no custom domain configured")
	}

	// If domain was activated, update Caddy config to remove it
	if cd.ActivatedAt.Valid {
		if err := updateCaddyConfig(""); err != nil {
			log.Printf("Warning: failed to update Caddy config: %v", err)
		}
	}

	// Delete from database
	_, err = db.Exec("DELETE FROM custom_domain WHERE id = ?", cd.ID)
	if err != nil {
		return fmt.Errorf("failed to remove domain: %v", err)
	}

	return nil
}

// generateCaddyfile generates the Caddyfile content (for reference/fallback only)
func generateCaddyfile(customDomain string) string {
	hostname := getInstanceHostname()
	if hostname == "" {
		hostname = "localhost:8080"
	}

	config := fmt.Sprintf(`# Postastiq Caddyfile - Auto-generated
# DO NOT EDIT MANUALLY

# Default instance hostname
%s {
	reverse_proxy localhost:8080
}
`, hostname)

	if customDomain != "" {
		config += fmt.Sprintf(`
# Custom domain (verified)
%s {
	reverse_proxy localhost:8080
}
`, customDomain)
	}

	return config
}

// getCaddyAPIURL returns the Caddy Admin API URL from environment or default
func getCaddyAPIURL() string {
	apiURL := os.Getenv("CADDY_API_URL")
	if apiURL == "" {
		apiURL = "http://caddy:2019"
	}
	return apiURL
}

// caddyRouteConfig represents a Caddy route configuration
type caddyRouteConfig struct {
	ID     string   `json:"@id"`
	Match  []match  `json:"match"`
	Handle []handle `json:"handle"`
}

type match struct {
	Host []string `json:"host"`
}

type handle struct {
	Handler   string     `json:"handler"`
	Upstreams []upstream `json:"upstreams,omitempty"`
}

type upstream struct {
	Dial string `json:"dial"`
}

// addCaddyRouteViaAPI adds a custom domain route via Caddy Admin API
func addCaddyRouteViaAPI(domain string) error {
	apiURL := getCaddyAPIURL()

	// Build the route config with @id for later DELETE operations
	route := caddyRouteConfig{
		ID: "custom-domain-" + domain,
		Match: []match{
			{Host: []string{domain}},
		},
		Handle: []handle{
			{
				Handler: "reverse_proxy",
				Upstreams: []upstream{
					{Dial: "blog:8080"},
				},
			},
		},
	}

	routeJSON, err := json.Marshal(route)
	if err != nil {
		return fmt.Errorf("failed to marshal route config: %v", err)
	}

	// POST to routes array to CREATE new route
	// The @id field in the payload registers it so DELETE /id/{id} works later
	req, err := http.NewRequest("POST", apiURL+"/config/apps/http/servers/srv0/routes", bytes.NewBuffer(routeJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call Caddy API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Caddy API error (%d): %s", resp.StatusCode, string(body))
	}

	log.Printf("Added custom domain route via Caddy API: %s", domain)
	return nil
}

// removeCaddyRouteViaAPI removes a custom domain route via Caddy Admin API
func removeCaddyRouteViaAPI(domain string) error {
	apiURL := getCaddyAPIURL()

	req, err := http.NewRequest("DELETE", apiURL+"/id/custom-domain-"+domain, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call Caddy API: %v", err)
	}
	defer resp.Body.Close()

	// 404 is okay - route might not exist
	if resp.StatusCode >= 400 && resp.StatusCode != 404 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Caddy API error (%d): %s", resp.StatusCode, string(body))
	}

	log.Printf("Removed custom domain route via Caddy API: %s", domain)
	return nil
}

// isCaddyAPIAvailable checks if Caddy Admin API is reachable
func isCaddyAPIAvailable() bool {
	apiURL := getCaddyAPIURL()
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(apiURL + "/config/")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

// syncCustomDomainWithCaddy re-registers any activated custom domain with Caddy on startup
// This ensures routes persist across container restarts since Caddy Admin API changes are in-memory only
func syncCustomDomainWithCaddy() {
	cd, err := getCustomDomain()
	if err != nil {
		log.Printf("Startup sync: No custom domain configured")
		return
	}

	if cd == nil || !cd.ActivatedAt.Valid {
		log.Printf("Startup sync: No activated custom domain to sync")
		return
	}

	// Wait briefly for Caddy to be ready (it may start after the blog container)
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		if isCaddyAPIAvailable() {
			break
		}
		if i < maxRetries-1 {
			log.Printf("Startup sync: Waiting for Caddy API (attempt %d/%d)...", i+1, maxRetries)
			time.Sleep(2 * time.Second)
		}
	}

	if !isCaddyAPIAvailable() {
		log.Printf("Startup sync: Caddy API not available, skipping route sync for %s", cd.Domain)
		return
	}

	err = addCaddyRouteViaAPI(cd.Domain)
	if err != nil {
		log.Printf("Startup sync: Failed to register route for %s: %v", cd.Domain, err)
	} else {
		log.Printf("Startup sync: Successfully registered route for custom domain %s", cd.Domain)
	}
}

// updateCaddyConfig updates Caddy configuration for custom domain (uses Admin API if available)
func updateCaddyConfig(customDomain string) error {
	// Try Caddy Admin API first (preferred for Docker setups)
	if isCaddyAPIAvailable() {
		if customDomain != "" {
			return addCaddyRouteViaAPI(customDomain)
		} else {
			// When customDomain is empty, we need to remove any existing route
			// Get the current custom domain from database
			cd, _ := getCustomDomain()
			if cd != nil && cd.Domain != "" {
				return removeCaddyRouteViaAPI(cd.Domain)
			}
		}
		return nil
	}

	// Fallback: Check for local Caddy file-based config (non-Docker setups)
	caddyfilePath := os.Getenv("CADDYFILE_PATH")
	if caddyfilePath == "" {
		caddyfilePath = "/etc/caddy/Caddyfile"
	}

	// Check if Caddy directory exists
	caddyDir := filepath.Dir(caddyfilePath)
	if _, err := os.Stat(caddyDir); os.IsNotExist(err) {
		log.Printf("Caddy not available (no API, no local config) - skipping Caddy configuration")
		return nil
	}

	// Check if caddy binary is available
	if _, err := exec.LookPath("caddy"); err != nil {
		log.Printf("Caddy binary not found - skipping Caddy configuration")
		return nil
	}

	// Generate new Caddyfile content
	config := generateCaddyfile(customDomain)

	// Write to temporary file first
	tmpFile := caddyfilePath + ".tmp"
	if err := os.WriteFile(tmpFile, []byte(config), 0644); err != nil {
		return fmt.Errorf("failed to write config: %v", err)
	}

	// Validate configuration
	cmd := exec.Command("caddy", "validate", "--config", tmpFile)
	if output, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("invalid caddy config: %s", string(output))
	}

	// Atomic move
	if err := os.Rename(tmpFile, caddyfilePath); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to install config: %v", err)
	}

	// Reload Caddy
	cmd = exec.Command("caddy", "reload", "--config", caddyfilePath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reload caddy: %s", string(output))
	}

	log.Printf("Caddy config updated successfully (custom domain: %s)", customDomain)
	return nil
}

// startDomainRevalidationCron starts a background job to re-verify active domains
func startDomainRevalidationCron() {
	ticker := time.NewTicker(7 * 24 * time.Hour) // Weekly
	go func() {
		for range ticker.C {
			revalidateActiveDomains()
		}
	}()
}

// revalidateActiveDomains checks that active domains still have valid DNS
func revalidateActiveDomains() {
	cd, err := getCustomDomain()
	if err != nil || cd == nil || !cd.ActivatedAt.Valid {
		return
	}

	log.Printf("Re-validating custom domain: %s", cd.Domain)

	if err := verifyDomainDNS(cd.Domain, cd.VerificationToken); err != nil {
		log.Printf("Domain %s failed re-verification: %v", cd.Domain, err)
		// Deactivate the domain
		_, err = db.Exec("UPDATE custom_domain SET activated_at = NULL WHERE id = ?", cd.ID)
		if err != nil {
			log.Printf("Failed to deactivate domain: %v", err)
		}
		// Update Caddy to remove the domain
		if err := updateCaddyConfig(""); err != nil {
			log.Printf("Failed to update Caddy config: %v", err)
		}
		return
	}

	// Update last verified timestamp
	_, err = db.Exec("UPDATE custom_domain SET last_verified_at = CURRENT_TIMESTAMP WHERE id = ?", cd.ID)
	if err != nil {
		log.Printf("Failed to update last_verified_at: %v", err)
	}
	log.Printf("Domain %s re-verification successful", cd.Domain)
}

// ============================================================================
// End Custom Domain Functions
// ============================================================================

func saveUploadedFile(file io.Reader, filename string) (string, error) {
	// Read file content
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Generate unique filename using hash
	hash := sha256.Sum256(data)
	ext := filepath.Ext(filename)
	newFilename := hex.EncodeToString(hash[:]) + ext
	filePath := filepath.Join(uploadsDir, newFilename)

	// Write file to disk
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %v", err)
	}

	// Return relative path for storage in database
	return newFilename, nil
}

func validateMediaAndSave(file io.Reader, filename string, title string, createdAt time.Time, mediaType string) (string, error) {
	// Check if audio uploads are enabled
	if mediaType == "audio" && !enableAudioUploads {
		return "", fmt.Errorf("audio uploads are disabled")
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Different size limits for different media types
	maxSize := 5 * 1024 * 1024 // 5MB default for images
	if mediaType == "audio" {
		maxSize = 20 * 1024 * 1024 // 20MB for audio
	} else if mediaType == "video" {
		maxSize = 50 * 1024 * 1024 // 50MB for video
	}

	if len(data) > maxSize {
		return "", fmt.Errorf("file size exceeds %dMB limit", maxSize/(1024*1024))
	}

	ext := strings.ToLower(filename[strings.LastIndex(filename, ".")+1:])

	// Define valid extensions by media type
	validExtensions := make(map[string]bool)
	if mediaType == "photo" {
		validExtensions = map[string]bool{
			"jpg":  true,
			"jpeg": true,
			"png":  true,
			"gif":  true,
			"webp": true,
		}
	} else if mediaType == "audio" {
		validExtensions = map[string]bool{
			"mp3":  true,
			"m4a":  true,
			"wav":  true,
			"ogg":  true,
			"aac":  true,
			"webm": true,
		}
	} else if mediaType == "video" {
		validExtensions = map[string]bool{
			"mp4":  true,
			"webm": true,
			"mov":  true,
			"avi":  true,
		}
	}

	if !validExtensions[ext] {
		return "", fmt.Errorf("invalid file type for %s", mediaType)
	}

	if len(data) < 4 {
		return "", fmt.Errorf("file is too small to be valid")
	}

	// Validate file magic numbers based on type
	validMagic := false
	if mediaType == "photo" {
		if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
			validMagic = true // JPEG
		} else if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
			validMagic = true // PNG
		} else if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 {
			validMagic = true // GIF
		} else if len(data) >= 12 && string(data[8:12]) == "WEBP" {
			validMagic = true // WebP
		}
	} else if mediaType == "audio" {
		if len(data) >= 3 && data[0] == 0xFF && (data[1]&0xE0) == 0xE0 {
			validMagic = true // MP3
		} else if len(data) >= 4 && string(data[0:4]) == "ftyp" {
			validMagic = true // M4A
		} else if len(data) >= 4 && string(data[0:4]) == "RIFF" {
			validMagic = true // WAV
		} else if len(data) >= 4 && string(data[0:4]) == "OggS" {
			validMagic = true // OGG
		} else if len(data) >= 4 && string(data[0:4]) == "\x1A\x45\xDF\xA3" {
			validMagic = true // WebM (audio recording from browser)
		} else {
			validMagic = true // Allow other audio formats (AAC, etc.)
		}
	} else if mediaType == "video" {
		if len(data) >= 12 && (string(data[4:8]) == "ftyp" || string(data[4:12]) == "ftypmp42") {
			validMagic = true // MP4
		} else if len(data) >= 4 && string(data[0:4]) == "\x1A\x45\xDF\xA3" {
			validMagic = true // WebM
		} else if len(data) >= 4 && string(data[4:8]) == "moov" {
			validMagic = true // MOV
		} else if len(data) >= 4 && string(data[0:4]) == "RIFF" {
			validMagic = true // AVI
		} else {
			validMagic = true // Allow other video formats
		}
	}

	if !validMagic {
		return "", fmt.Errorf("file content does not match expected %s format", mediaType)
	}

	// Generate filename from title (same pattern as slug)
	slug := strings.TrimSpace(title)
	slug = strings.ToLower(slug)

	reg := regexp.MustCompile(`[^a-z0-9\s]+`)
	slug = reg.ReplaceAllString(slug, "")

	slug = strings.ReplaceAll(slug, " ", "-")

	multiHyphen := regexp.MustCompile(`-+`)
	slug = multiHyphen.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	if len(slug) > 80 {
		slug = slug[:80]
		slug = strings.TrimRight(slug, "-")
	}

	if slug == "" {
		slug = "media"
	}

	// Append date in format: yyyy-mm-dd (same as page slug)
	dateStr := createdAt.Format("2006-01-02")
	slug = fmt.Sprintf("%s-%s", slug, dateStr)

	// Save file and return path
	newFilename := slug + filepath.Ext(filename)
	filePath := filepath.Join(uploadsDir, newFilename)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %v", err)
	}

	return newFilename, nil
}

// Backward compatibility wrapper
func validateImageAndSave(file io.Reader, filename string, title string, createdAt time.Time) (string, error) {
	return validateMediaAndSave(file, filename, title, createdAt, "photo")
}

// Thumbnail upload and save function for video/audio entries
func validateAndSaveThumbnail(file io.Reader, filename string, title string, createdAt time.Time) (string, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read thumbnail file: %v", err)
	}

	// Max size 5MB for thumbnails
	if len(data) > 5*1024*1024 {
		return "", fmt.Errorf("thumbnail size exceeds 5MB limit")
	}

	ext := strings.ToLower(filename[strings.LastIndex(filename, ".")+1:])

	// Only allow image formats for thumbnails
	validExtensions := map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"gif":  true,
		"webp": true,
	}

	if !validExtensions[ext] {
		return "", fmt.Errorf("invalid thumbnail type: only images allowed")
	}

	if len(data) < 4 {
		return "", fmt.Errorf("thumbnail file is too small to be valid")
	}

	// Validate file magic numbers for images
	validMagic := false
	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		validMagic = true // JPEG
	} else if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		validMagic = true // PNG
	} else if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 {
		validMagic = true // GIF
	} else if len(data) >= 12 && string(data[8:12]) == "WEBP" {
		validMagic = true // WebP
	}

	if !validMagic {
		return "", fmt.Errorf("thumbnail content does not match expected image format")
	}

	// Generate filename from title with "thumb-" prefix
	slug := strings.TrimSpace(title)
	slug = strings.ToLower(slug)

	reg := regexp.MustCompile(`[^a-z0-9\s]+`)
	slug = reg.ReplaceAllString(slug, "")

	slug = strings.ReplaceAll(slug, " ", "-")

	multiHyphen := regexp.MustCompile(`-+`)
	slug = multiHyphen.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	if len(slug) > 80 {
		slug = slug[:80]
		slug = strings.TrimRight(slug, "-")
	}

	if slug == "" {
		slug = "thumbnail"
	}

	// Append date in format: yyyy-mm-dd
	dateStr := createdAt.Format("2006-01-02")
	slug = fmt.Sprintf("thumb-%s-%s", slug, dateStr)

	// Save file and return path
	newFilename := slug + "." + ext
	filePath := filepath.Join(uploadsDir, newFilename)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write thumbnail: %v", err)
	}

	return newFilename, nil
}

// Avatar upload and resize function
func validateAndResizeAvatar(file io.Reader, filename string) (string, error) {
	// Read file content
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Check file size (max 5MB)
	if len(data) > 5*1024*1024 {
		return "", fmt.Errorf("file size exceeds 5MB limit")
	}

	// Get extension and validate
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", fmt.Errorf("only JPG and PNG files are supported for avatars")
	}

	// Decode image
	var img image.Image
	if ext == ".png" {
		img, err = png.Decode(bytes.NewReader(data))
	} else {
		img, err = jpeg.Decode(bytes.NewReader(data))
	}
	if err != nil {
		return "", fmt.Errorf("invalid image file: %v", err)
	}

	// Resize to 200x200 (standard avatar size)
	resizedImg := resize.Resize(200, 200, img, resize.Lanczos3)

	// Generate filename
	hash := sha256.Sum256(data)
	newFilename := "avatar-" + hex.EncodeToString(hash[:16]) + ".png"
	filePath := filepath.Join(uploadsDir, newFilename)

	// Create output file
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	// Encode as PNG (best quality for avatars)
	err = png.Encode(outFile, resizedImg)
	if err != nil {
		return "", fmt.Errorf("failed to encode image: %v", err)
	}

	return newFilename, nil
}

// Session management functions
func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func createSession() (string, error) {
	token, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	session := &Session{
		Token:     token,
		ExpiresAt: time.Now().Add(60 * time.Minute),
	}

	sessionMutex.Lock()
	sessions[token] = session
	sessionMutex.Unlock()

	return token, nil
}

func getSession(token string) (*Session, bool) {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()

	session, exists := sessions[token]
	if !exists {
		return nil, false
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, false
	}

	return session, true
}

func deleteSession(token string) {
	sessionMutex.Lock()
	delete(sessions, token)
	sessionMutex.Unlock()
}

func cleanupExpiredSessions() {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	now := time.Now()
	for token, session := range sessions {
		if now.After(session.ExpiresAt) {
			delete(sessions, token)
		}
	}
}

func isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false
	}

	_, valid := getSession(cookie.Value)
	return valid
}

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAuthenticated(r) {
			http.Redirect(w, r, "/login?redirect="+r.URL.Path, http.StatusSeeOther)
			return
		}

		// Check if password change is required (first login with auto-generated password)
		// Skip this check if already on the change-password page
		if r.URL.Path != "/change-password" && isPasswordChangeRequired() {
			http.Redirect(w, r, "/change-password", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}

func timeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		seconds := int(duration.Seconds())
		if seconds <= 1 {
			return "just now"
		}
		return fmt.Sprintf("%ds ago", seconds)
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}

func getInitialLetter(content string) string {
	if len(content) > 0 {
		return string([]rune(content)[0])
	}
	return "?"
}

func truncateContent(content string, maxLength int) string {
	runes := []rune(content)
	if len(runes) <= maxLength {
		return content
	}
	return string(runes[:maxLength]) + "..."
}

func linkifyContent(content string) template.HTML {
	escaped := html.EscapeString(content)
	linked := urlRegex.ReplaceAllStringFunc(escaped, func(url string) string {
		return fmt.Sprintf(`<a href="%s" target="_blank" rel="noopener noreferrer">%s</a>`, url, url)
	})
	return template.HTML(linked)
}

func generateSlug(title string, createdAt time.Time) string {
	// Trim leading/trailing whitespace
	slug := strings.TrimSpace(title)

	// Convert to lowercase
	slug = strings.ToLower(slug)

	// Remove punctuation (keep only alphanumeric and spaces)
	reg := regexp.MustCompile(`[^a-z0-9\s]+`)
	slug = reg.ReplaceAllString(slug, "")

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Collapse multiple hyphens into one
	multiHyphen := regexp.MustCompile(`-+`)
	slug = multiHyphen.ReplaceAllString(slug, "-")

	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	// Truncate to 80 characters if needed
	if len(slug) > 80 {
		slug = slug[:80]
		// Remove trailing hyphen if truncation created one
		slug = strings.TrimRight(slug, "-")
	}

	// If slug is empty after processing, use fallback
	if slug == "" {
		slug = "post"
	}

	// Append date in format: yyyy-mm-dd
	dateStr := createdAt.Format("2006-01-02")
	slug = fmt.Sprintf("%s-%s", slug, dateStr)

	return slug
}

func getEntries(offset, limit int) ([]EntryDisplay, bool, error) {
	query := `
		SELECT id, title, content, photo_path, media_type, thumbnail_path, slug, created_at
		FROM entries
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := db.Query(query, limit+1, offset)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var entries []EntryDisplay
	count := 0
	for rows.Next() {
		count++
		if count > limit {
			break
		}

		var entry Entry
		var title sql.NullString
		var photoPath sql.NullString
		var mediaType sql.NullString
		var thumbnailPath sql.NullString
		var slug sql.NullString

		err := rows.Scan(&entry.ID, &title, &entry.Content, &photoPath, &mediaType, &thumbnailPath, &slug, &entry.CreatedAt)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			continue
		}

		if title.Valid {
			entry.Title = title.String
		}

		if slug.Valid {
			entry.Slug = slug.String
		}

		if mediaType.Valid {
			entry.MediaType = mediaType.String
		} else {
			entry.MediaType = "photo" // default for old entries
		}

		var photoURL template.URL
		hasPhoto := false
		hasAudio := false
		hasVideo := false

		if photoPath.Valid && photoPath.String != "" {
			photoURL = template.URL("/uploads/" + photoPath.String)
			if entry.MediaType == "photo" {
				hasPhoto = true
			} else if entry.MediaType == "audio" {
				hasAudio = true
			} else if entry.MediaType == "video" {
				hasVideo = true
			}
		}

		var thumbnailURL template.URL
		hasThumbnail := false
		if thumbnailPath.Valid && thumbnailPath.String != "" {
			thumbnailURL = template.URL("/uploads/" + thumbnailPath.String)
			hasThumbnail = true
		}

		truncatedContent := truncateContent(entry.Content, 150)

		entries = append(entries, EntryDisplay{
			ID:            entry.ID,
			Title:         entry.Title,
			Content:       linkifyContent(truncatedContent),
			Slug:          entry.Slug,
			Photo:         photoURL,
			HasPhoto:      hasPhoto,
			MediaType:     entry.MediaType,
			HasAudio:      hasAudio,
			HasVideo:      hasVideo,
			Thumbnail:     thumbnailURL,
			HasThumbnail:  hasThumbnail,
			CreatedAt:     entry.CreatedAt,
			TimeAgo:       timeAgo(entry.CreatedAt),
			InitialLetter: getInitialLetter(entry.Content),
		})
	}

	hasMore := count > limit
	return entries, hasMore, nil
}

// Blog viewer handlers
func handleBlogFeed(w http.ResponseWriter, r *http.Request) {
	// Auto-detect hostname from first request
	detectHostnameFromRequest(r)

	tmpl, err := template.New("feed").Parse(viewerTemplate)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	entries, hasMore, err := getEntries(0, 10)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		log.Printf("Database query error: %v", err)
		return
	}

	var totalEntries, todayEntries int
	db.QueryRow("SELECT COUNT(*) FROM entries").Scan(&totalEntries)
	db.QueryRow("SELECT COUNT(*) FROM entries WHERE DATE(created_at) = DATE('now')").Scan(&todayEntries)

	// Get settings from database
	settings, err := getSiteSettings()
	if err != nil {
		log.Printf("Error getting site settings: %v", err)
		// Use defaults if error
		settings = SiteSettings{
			SiteTitle:    "My Blog",
			SiteSubtitle: "A Personal Blog",
			UserInitial:  "AB",
			SiteTheme:    "default",
		}
	}

	data := ViewerPageData{
		Entries:          entries,
		TotalEntries:     totalEntries,
		TodayEntries:     todayEntries,
		SiteTitle:        settings.SiteTitle,
		SiteSubtitle:     settings.SiteSubtitle,
		EnableSubtitle:   enableSubtitle,
		UserInitial:      settings.UserInitial,
		AvatarPath:       settings.AvatarPath,
		AvatarPreference: settings.AvatarPreference,
		InitialCount:     len(entries),
		HasMore:          hasMore,
		ThemeCSS:         template.CSS(getThemeCSS()),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func handleAPIEntries(w http.ResponseWriter, r *http.Request) {
	offset := 0
	limit := 10

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil {
			offset = val
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 && val <= 50 {
			limit = val
		}
	}

	entries, hasMore, err := getEntries(offset, limit)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		log.Printf("Database query error: %v", err)
		return
	}

	response := APIResponse{
		Entries: entries,
		HasMore: hasMore,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleSinglePost(w http.ResponseWriter, r *http.Request) {
	// Extract slug from URL path
	// Expected format: /posts/slug-here/
	path := r.URL.Path
	slug := strings.TrimPrefix(path, "/posts/")
	slug = strings.TrimSuffix(slug, "/")

	if slug == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Query database for entry by slug
	query := `
		SELECT id, title, content, photo_path, media_type, thumbnail_path, slug, created_at
		FROM entries
		WHERE slug = ?
		LIMIT 1
	`

	var entry Entry
	var title sql.NullString
	var photoPath sql.NullString
	var mediaType sql.NullString
	var thumbnailPath sql.NullString
	var entrySlug sql.NullString

	err := db.QueryRow(query, slug).Scan(&entry.ID, &title, &entry.Content, &photoPath, &mediaType, &thumbnailPath, &entrySlug, &entry.CreatedAt)
	if err == sql.ErrNoRows {
		handle404(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("Database error: %v", err)
		return
	}

	if title.Valid {
		entry.Title = title.String
	}
	if photoPath.Valid {
		entry.PhotoPath = photoPath.String
	}
	if mediaType.Valid {
		entry.MediaType = mediaType.String
	} else {
		entry.MediaType = "photo"
	}
	if thumbnailPath.Valid {
		entry.ThumbnailPath = thumbnailPath.String
	}
	if entrySlug.Valid {
		entry.Slug = entrySlug.String
	}

	var photoURL template.URL
	hasPhoto := false
	hasAudio := false
	hasVideo := false
	if entry.PhotoPath != "" {
		photoURL = template.URL("/uploads/" + entry.PhotoPath)
		if entry.MediaType == "photo" {
			hasPhoto = true
		} else if entry.MediaType == "audio" {
			hasAudio = true
		} else if entry.MediaType == "video" {
			hasVideo = true
		}
	}

	var thumbnailURL template.URL
	hasThumbnail := false
	if entry.ThumbnailPath != "" {
		thumbnailURL = template.URL("/uploads/" + entry.ThumbnailPath)
		hasThumbnail = true
	}

	// Check if content needs truncation
	truncatedContent := truncateContent(entry.Content, 150)
	isTruncated := len([]rune(entry.Content)) > 150

	entryDisplay := EntryDisplay{
		ID:            entry.ID,
		Title:         entry.Title,
		Content:       linkifyContent(truncatedContent),
		FullContent:   linkifyContent(entry.Content),
		IsTruncated:   isTruncated,
		Photo:         photoURL,
		HasPhoto:      hasPhoto,
		MediaType:     entry.MediaType,
		HasAudio:      hasAudio,
		HasVideo:      hasVideo,
		Thumbnail:     thumbnailURL,
		HasThumbnail:  hasThumbnail,
		Slug:          entry.Slug,
		CreatedAt:     entry.CreatedAt,
		TimeAgo:       timeAgo(entry.CreatedAt),
		InitialLetter: getInitialLetter(entry.Content),
	}

	// Get settings from database
	settings, err := getSiteSettings()
	if err != nil {
		log.Printf("Error getting site settings: %v", err)
		// Use defaults if error
		settings = SiteSettings{
			SiteTitle:    "My Blog",
			SiteSubtitle: "A Personal Blog",
			UserInitial:  "AB",
			SiteTheme:    "default",
		}
	}

	data := SinglePostPageData{
		Entry:            entryDisplay,
		SiteTitle:        settings.SiteTitle,
		SiteSubtitle:     settings.SiteSubtitle,
		EnableSubtitle:   enableSubtitle,
		UserInitial:      settings.UserInitial,
		AvatarPath:       settings.AvatarPath,
		AvatarPreference: settings.AvatarPreference,
		ThemeCSS:         template.CSS(getThemeCSS()),
	}

	tmpl, err := template.New("post").Parse(postTemplate)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func handle404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	// Get recent entries
	recentEntries, _, err := getEntries(0, 3)
	if err != nil {
		log.Printf("Error fetching recent entries for 404 page: %v", err)
		recentEntries = []EntryDisplay{}
	}

	// Get settings from database
	settings, err := getSiteSettings()
	if err != nil {
		log.Printf("Error getting site settings: %v", err)
		// Use defaults if error
		settings = SiteSettings{
			SiteTitle:    "My Blog",
			SiteSubtitle: "A Personal Blog",
			UserInitial:  "AB",
			SiteTheme:    "default",
		}
	}

	data := NotFoundPageData{
		SiteTitle:      settings.SiteTitle,
		SiteSubtitle:   settings.SiteSubtitle,
		EnableSubtitle: enableSubtitle,
		UserInitial:    settings.UserInitial,
		RecentEntries:  recentEntries,
		ThemeCSS:       template.CSS(getThemeCSS()),
	}

	tmpl, err := template.New("404").Parse(notFoundTemplate)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("404 Template error: %v", err)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("404 Template execution error: %v", err)
	}
}

// Admin service handlers
func handleAdminIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Determine which view to show (default to new post form)
	view := r.URL.Query().Get("view")
	if view == "" {
		view = "new"
	}

	// Template functions
	funcMap := template.FuncMap{
		"jsEscape": jsEscape,
		"add": func(a, b int) int {
			return a + b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
	}
	tmpl := template.Must(template.New("admin").Funcs(funcMap).Parse(adminTemplate))

	// Check if viewer password is set in site_settings
	hasPassword := false
	var passwordHash sql.NullString
	err := db.QueryRow("SELECT viewer_password_hash FROM site_settings WHERE id = 1").Scan(&passwordHash)
	if err == nil && passwordHash.Valid && passwordHash.String != "" {
		hasPassword = true
	}

	// Read flash messages from cookies
	var message, messageType string
	if cookie, err := r.Cookie("flash_message"); err == nil {
		message = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "flash_message",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}
	if cookie, err := r.Cookie("flash_type"); err == nil {
		messageType = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "flash_type",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}

	data := EditorPageData{
		HasPrivacyPassword: hasPassword,
		EnableAudioUploads: enableAudioUploads,
		Message:            message,
		MessageType:        messageType,
		View:               view,
		PageTitle:          "Posts",
	}

	// For posts view, fetch paginated entries
	if view == "posts" {
		// Pagination settings
		perPage := 20
		pageStr := r.URL.Query().Get("page")
		currentPage, err := strconv.Atoi(pageStr)
		if err != nil || currentPage < 1 {
			currentPage = 1
		}

		// Get total count
		var totalEntries int
		err = db.QueryRow("SELECT COUNT(*) FROM entries").Scan(&totalEntries)
		if err != nil {
			log.Printf("Error counting entries: %v", err)
			totalEntries = 0
		}

		// Calculate pagination
		totalPages := (totalEntries + perPage - 1) / perPage
		if totalPages < 1 {
			totalPages = 1
		}
		if currentPage > totalPages {
			currentPage = totalPages
		}

		offset := (currentPage - 1) * perPage

		// Fetch entries for current page
		entries, err := getPaginatedEntries(offset, perPage)
		if err != nil {
			log.Printf("Error fetching entries: %v", err)
			entries = []Entry{}
		}

		// Calculate page numbers to show (max 5 pages around current)
		var pageNumbers []int
		startPage := currentPage - 2
		if startPage < 1 {
			startPage = 1
		}
		endPage := startPage + 4
		if endPage > totalPages {
			endPage = totalPages
			startPage = endPage - 4
			if startPage < 1 {
				startPage = 1
			}
		}
		for i := startPage; i <= endPage; i++ {
			pageNumbers = append(pageNumbers, i)
		}

		// Calculate start/end entry numbers for display
		startEntry := offset + 1
		endEntry := offset + len(entries)
		if totalEntries == 0 {
			startEntry = 0
			endEntry = 0
		}

		data.Entries = entries
		data.CurrentPage = currentPage
		data.TotalPages = totalPages
		data.TotalEntries = totalEntries
		data.StartEntry = startEntry
		data.EndEntry = endEntry
		data.PageNumbers = pageNumbers
		data.PageTitle = "Posts"
	} else if view == "new" {
		data.PageTitle = "New Post"
	}

	tmpl.Execute(w, data)
}

func handleGetEntriesForAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	offsetStr := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	limit := 5
	entries, err := getPaginatedEntries(offset, limit)
	if err != nil {
		log.Printf("Error fetching entries: %v", err)
		http.Error(w, "Error fetching entries", http.StatusInternalServerError)
		return
	}

	hasMore := len(entries) == limit

	var entriesJSON []EntryJSON
	for _, entry := range entries {
		photoURL := ""
		if entry.PhotoPath != "" {
			photoURL = "/uploads/" + entry.PhotoPath
		}
		entriesJSON = append(entriesJSON, EntryJSON{
			ID:        entry.ID,
			Title:     entry.Title,
			Content:   entry.Content,
			Photo:     photoURL,
			Slug:      entry.Slug,
			CreatedAt: entry.CreatedAt.Format("Jan 2, 2006 3:04 PM"),
		})
	}

	response := EntriesResponse{
		Entries: entriesJSON,
		HasMore: hasMore,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if len(content) > 2000 {
		content = content[:2000]
	}

	// Auto-generate title from content (first line, trim to 60 chars at last full word)
	var finalTitle string
	if content != "" {
		firstLine := strings.Split(content, "\n")[0]
		if len(firstLine) > 60 {
			// Find last space before 60 chars
			truncated := firstLine[:60]
			lastSpace := strings.LastIndex(truncated, " ")
			if lastSpace > 0 {
				finalTitle = truncated[:lastSpace]
			} else {
				finalTitle = truncated
			}
		} else {
			finalTitle = firstLine
		}
	} else {
		// Fallback title
		finalTitle = fmt.Sprintf("Untitled Post %s", time.Now().Format("2006-01-02"))
	}

	// Generate timestamp for both slug and media
	now := time.Now()
	slug := generateSlug(finalTitle, now)

	// Get media type from form (defaults to photo for backward compatibility)
	mediaType := r.FormValue("media_type")
	if mediaType == "" {
		mediaType = "photo"
	}

	var photoPath string
	file, header, err := r.FormFile("media")
	if err == nil {
		defer file.Close()
		photoPath, err = validateMediaAndSave(file, header.Filename, finalTitle, now, mediaType)
		if err != nil {
			log.Printf("Error saving %s: %v", mediaType, err)
		}
	}

	// Handle thumbnail upload for video/audio
	var thumbnailPath string
	if mediaType == "video" || mediaType == "audio" {
		thumbFile, thumbHeader, thumbErr := r.FormFile("thumbnail")
		if thumbErr == nil {
			defer thumbFile.Close()
			thumbnailPath, thumbErr = validateAndSaveThumbnail(thumbFile, thumbHeader.Filename, finalTitle, now)
			if thumbErr != nil {
				log.Printf("Error saving thumbnail: %v", thumbErr)
			}
		}
	}

	_, err = db.Exec("INSERT INTO entries (title, content, photo_path, media_type, thumbnail_path, slug) VALUES (?, ?, ?, ?, ?, ?)", finalTitle, content, photoPath, mediaType, thumbnailPath, slug)
	if err != nil {
		log.Printf("Error inserting entry: %v", err)
		showMessage(w, r, "Failed to create entry", "error")
		return
	}

	// Redirect to the new post page
	http.Redirect(w, r, "/posts/"+slug+"/", http.StatusSeeOther)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if len(content) > 2000 {
		content = content[:2000]
	}

	// Get the created_at timestamp for slug generation
	var createdAt time.Time
	err = db.QueryRow("SELECT created_at FROM entries WHERE id = ?", id).Scan(&createdAt)
	if err != nil {
		log.Printf("Error fetching entry: %v", err)
		showMessage(w, r, "Failed to find entry", "error")
		return
	}

	// Auto-generate title from content (first line, trim to 60 chars at last full word)
	var finalTitle string
	if content != "" {
		firstLine := strings.Split(content, "\n")[0]
		if len(firstLine) > 60 {
			// Find last space before 60 chars
			truncated := firstLine[:60]
			lastSpace := strings.LastIndex(truncated, " ")
			if lastSpace > 0 {
				finalTitle = truncated[:lastSpace]
			} else {
				finalTitle = truncated
			}
		} else {
			finalTitle = firstLine
		}
	} else {
		// Fallback title
		finalTitle = fmt.Sprintf("Untitled Post %s", createdAt.Format("2006-01-02"))
	}

	// Regenerate slug with final title
	slug := generateSlug(finalTitle, createdAt)

	// Get media type from form (defaults to photo for backward compatibility)
	mediaType := r.FormValue("media_type")
	if mediaType == "" {
		mediaType = "photo"
	}

	// Handle thumbnail upload/removal for video/audio
	var thumbnailPath sql.NullString
	thumbnailUpdated := false
	removeThumbnail := r.FormValue("remove_thumbnail") == "1"

	if removeThumbnail {
		thumbnailPath = sql.NullString{String: "", Valid: false}
		thumbnailUpdated = true
	} else if mediaType == "video" || mediaType == "audio" {
		thumbFile, thumbHeader, thumbErr := r.FormFile("thumbnail")
		if thumbErr == nil {
			defer thumbFile.Close()
			newThumbPath, thumbErr := validateAndSaveThumbnail(thumbFile, thumbHeader.Filename, finalTitle, createdAt)
			if thumbErr != nil {
				log.Printf("Error saving thumbnail: %v", thumbErr)
			} else {
				thumbnailPath = sql.NullString{String: newThumbPath, Valid: true}
				thumbnailUpdated = true
			}
		}
	}

	file, header, err := r.FormFile("media")
	if err == nil {
		defer file.Close()
		photoPath, err := validateMediaAndSave(file, header.Filename, finalTitle, createdAt, mediaType)
		if err != nil {
			log.Printf("Error saving %s: %v", mediaType, err)
		} else {
			if thumbnailUpdated {
				_, err = db.Exec("UPDATE entries SET title = ?, content = ?, photo_path = ?, media_type = ?, thumbnail_path = ?, slug = ? WHERE id = ?", finalTitle, content, photoPath, mediaType, thumbnailPath, slug, id)
			} else {
				_, err = db.Exec("UPDATE entries SET title = ?, content = ?, photo_path = ?, media_type = ?, slug = ? WHERE id = ?", finalTitle, content, photoPath, mediaType, slug, id)
			}
			if err != nil {
				log.Printf("Error updating entry: %v", err)
				showMessage(w, r, "Failed to update entry", "error")
				return
			}
			http.Redirect(w, r, "/posts/"+slug+"/", http.StatusSeeOther)
			return
		}
	}

	if thumbnailUpdated {
		_, err = db.Exec("UPDATE entries SET title = ?, content = ?, thumbnail_path = ?, slug = ? WHERE id = ?", finalTitle, content, thumbnailPath, slug, id)
	} else {
		_, err = db.Exec("UPDATE entries SET title = ?, content = ?, slug = ? WHERE id = ?", finalTitle, content, slug, id)
	}
	if err != nil {
		log.Printf("Error updating entry: %v", err)
		showMessage(w, r, "Failed to update entry", "error")
		return
	}

	http.Redirect(w, r, "/posts/"+slug+"/", http.StatusSeeOther)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM entries WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting entry: %v", err)
		showMessage(w, r, "Failed to delete entry", "error")
		return
	}

	showMessage(w, r, "Entry deleted successfully!", "success")
}

func getAllEntries() ([]Entry, error) {
	rows, err := db.Query("SELECT id, title, content, photo_path, media_type, thumbnail_path, slug, created_at FROM entries ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var entry Entry
		var title sql.NullString
		var photoPath sql.NullString
		var mediaType sql.NullString
		var thumbnailPath sql.NullString
		var slug sql.NullString
		err := rows.Scan(&entry.ID, &title, &entry.Content, &photoPath, &mediaType, &thumbnailPath, &slug, &entry.CreatedAt)
		if err != nil {
			return nil, err
		}
		if title.Valid {
			entry.Title = title.String
		}
		if photoPath.Valid {
			entry.PhotoPath = photoPath.String
		}
		if mediaType.Valid {
			entry.MediaType = mediaType.String
		} else {
			entry.MediaType = "photo"
		}
		if thumbnailPath.Valid {
			entry.ThumbnailPath = thumbnailPath.String
		}
		if slug.Valid {
			entry.Slug = slug.String
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func getPaginatedEntries(offset, limit int) ([]Entry, error) {
	rows, err := db.Query("SELECT id, title, content, photo_path, media_type, thumbnail_path, slug, created_at FROM entries ORDER BY created_at DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var entry Entry
		var title sql.NullString
		var photoPath sql.NullString
		var mediaType sql.NullString
		var thumbnailPath sql.NullString
		var slug sql.NullString
		err := rows.Scan(&entry.ID, &title, &entry.Content, &photoPath, &mediaType, &thumbnailPath, &slug, &entry.CreatedAt)
		if err != nil {
			return nil, err
		}
		if title.Valid {
			entry.Title = title.String
		}
		if photoPath.Valid {
			entry.PhotoPath = photoPath.String
		}
		if mediaType.Valid {
			entry.MediaType = mediaType.String
		} else {
			entry.MediaType = "photo"
		}
		if thumbnailPath.Valid {
			entry.ThumbnailPath = thumbnailPath.String
		}
		if slug.Valid {
			entry.Slug = slug.String
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func showMessage(w http.ResponseWriter, r *http.Request, message, messageType string) {
	// Set flash message cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "flash_message",
		Value:    message,
		Path:     "/",
		MaxAge:   5, // 5 seconds
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "flash_type",
		Value:    messageType,
		Path:     "/",
		MaxAge:   5, // 5 seconds
		HttpOnly: true,
	})

	// Redirect to admin page
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func showSettingsMessage(w http.ResponseWriter, r *http.Request, message, messageType, section string) {
	// Set flash message cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "flash_message",
		Value:    message,
		Path:     "/",
		MaxAge:   5, // 5 seconds
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "flash_type",
		Value:    messageType,
		Path:     "/",
		MaxAge:   5, // 5 seconds
		HttpOnly: true,
	})

	// Redirect to appropriate settings section
	redirectURL := "/admin/settings"
	switch section {
	case "appearance":
		redirectURL = "/admin/settings/appearance"
	case "security":
		redirectURL = "/admin/settings/security"
	case "backup":
		redirectURL = "/admin/settings/backup"
	}
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func jsEscape(s string) string {
	return template.JSEscapeString(s)
}

func getThemeCSS() string {
	// Get theme from database settings
	settings, err := getSiteSettings()
	if err != nil {
		log.Printf("Error getting site settings for theme: %v", err)
		return defaultThemeCSS
	}

	switch settings.SiteTheme {
	case "dark":
		return darkThemeCSS
	case "custom":
		return generateCustomThemeCSS(settings.CustomBgColor, settings.CustomTextColor, settings.CustomAccentColor)
	default:
		return defaultThemeCSS
	}
}

// generateCustomThemeCSS creates a custom theme CSS based on 3 base colors
func generateCustomThemeCSS(bgColor, textColor, accentColor string) string {
	// Parse colors and calculate derived colors
	isLight := isLightHexColor(bgColor)

	// Derive secondary colors from the base colors
	headerBg := adjustHexBrightness(bgColor, ifThen(isLight, -5, 10))
	cardBg := adjustHexBrightness(bgColor, ifThen(isLight, 3, 5))
	borderColor := adjustHexBrightness(bgColor, ifThen(isLight, -15, 20))
	secondaryText := blendColors(textColor, bgColor, 0.5)
	hoverShadow := "rgba(0, 0, 0, 0.15)"
	if !isLight {
		hoverShadow = "rgba(255, 255, 255, 0.1)"
	}
	subtleBorder := adjustHexBrightness(bgColor, ifThen(isLight, -8, 12))

	return fmt.Sprintf(`
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: %[1]s;
            color: %[2]s;
            line-height: 1.6;
        }

        .container {
            max-width: 614px;
            margin: 0 auto;
            background-color: %[1]s;
            min-height: 100vh;
        }

        .header {
            background: %[3]s;
            padding: 20px 32px;
            color: %[2]s;
            border-bottom: 1px solid %[4]s;
            position: sticky;
            top: 0;
            z-index: 100;
        }

        .header-content {
            display: flex;
            align-items: center;
            justify-content: space-between;
        }

        .header h1 {
            font-size: 28px;
            font-weight: 600;
            display: flex;
            align-items: center;
            gap: 12px;
            font-family: 'Segoe UI', Roboto, sans-serif;
        }

        .subtitle {
            font-size: 14px;
            color: %[5]s;
            margin-top: 4px;
            font-weight: 400;
        }

        .stats {
            padding: 16px 0;
            background-color: %[1]s;
            display: none;
        }

        .stat {
            text-align: center;
            padding: 16px;
            background-color: %[6]s;
            border: 1px solid %[4]s;
            border-radius: 8px;
            transition: transform 0.2s, box-shadow 0.2s;
        }

        .stat:hover {
            transform: translateY(-2px);
            box-shadow: 0 2px 8px %[7]s;
        }

        .stat-value {
            font-size: 32px;
            font-weight: 700;
            color: %[2]s;
            margin-bottom: 4px;
        }

        .stat-label {
            font-size: 13px;
            color: %[5]s;
            text-transform: uppercase;
            font-weight: 500;
            letter-spacing: 0.5px;
        }

        .feed {
            padding: 24px 0 80px;
        }

        .entry {
            padding: 0;
            margin-bottom: 24px;
            background-color: %[6]s;
            border: 1px solid %[4]s;
            border-radius: 8px;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
            display: flex;
            flex-direction: column;
        }

        .entry:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px %[7]s;
        }

        .entry-header {
            display: flex;
            align-items: center;
            padding: 14px 16px;
            border-bottom: 1px solid %[8]s;
        }

        .avatar {
            width: 56px;
            height: 56px;
            border-radius: 50%%;
            background: %[9]s;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 600;
            font-size: 20px;
            color: #ffffff;
            margin-right: 16px;
            flex-shrink: 0;
            padding: 3px;
        }

        .avatar-inner {
            width: 100%%;
            height: 100%%;
            background: %[6]s;
            border-radius: 50%%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: %[2]s;
        }

        .entry-info {
            flex: 1;
        }

        .username {
            display: none;
        }

        .timestamp {
            color: %[5]s;
            font-size: 12px;
            font-weight: 400;
            margin-top: 4px;
            display: block;
        }

        .entry-photo-container {
            width: 100%%;
            max-width: 1024px;
            height: 0;
            padding-bottom: 100%%;
            position: relative;
            overflow: hidden;
            background-color: #000000;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .entry-photo {
            position: absolute;
            top: 50%%;
            left: 50%%;
            transform: translate(-50%%, -50%%);
            max-width: 100%%;
            max-height: 100%%;
            width: auto;
            height: auto;
            object-fit: contain;
            display: block;
        }

        .entry-actions {
            display: none;
        }

        .action-btn {
            background: none;
            border: none;
            cursor: pointer;
            font-size: 24px;
            padding: 8px;
            line-height: 1;
            color: %[2]s;
            transition: opacity 0.2s;
        }

        .action-btn:hover {
            opacity: 0.6;
        }

        .entry-content {
            font-size: 15px;
            color: %[2]s;
            line-height: 20px;
            padding: 16px 16px 4px;
            word-wrap: break-word;
            overflow-wrap: break-word;
        }

        .entry-content h2 {
            font-size: 18px;
            font-weight: 600;
            margin-bottom: 8px;
            color: %[2]s;
        }

        .entry-content p {
            font-size: 15px;
            line-height: 1.6;
            color: %[2]s;
        }

        .entry-content a {
            color: %[9]s;
            text-decoration: none;
            font-weight: 500;
        }

        .entry-content a:hover {
            text-decoration: underline;
        }

        .entry-content-username {
            display: none;
        }

        .entry-timestamp {
            padding: 4px 16px 16px;
            color: %[5]s;
            font-size: 10px;
            text-transform: uppercase;
            letter-spacing: 0.2px;
        }

        .entry-footer {
            display: none;
        }

        .entry-id {
            color: %[5]s;
            font-size: 12px;
            font-weight: 500;
        }

        .entry-badge {
            background-color: %[9]s;
            color: #ffffff;
            padding: 4px 12px;
            border-radius: 4px;
            font-size: 11px;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        .empty-state {
            padding: 80px 32px;
            text-align: center;
        }

        .empty-state-icon {
            width: 120px;
            height: 120px;
            margin: 0 auto 24px;
            background: %[6]s;
            border-radius: 50%%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 64px;
            border: 2px solid %[4]s;
        }

        .empty-state-title {
            font-size: 28px;
            font-weight: 600;
            margin-bottom: 12px;
            color: %[2]s;
        }

        .empty-state-text {
            font-size: 16px;
            color: %[5]s;
            max-width: 400px;
            margin: 0 auto;
        }

        .loading-indicator {
            text-align: center;
            padding: 24px;
            color: %[5]s;
            font-size: 14px;
            display: none;
        }

        .loading-indicator.show {
            display: block;
        }

        .loading-spinner {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 2px solid %[4]s;
            border-top-color: %[2]s;
            border-radius: 50%%;
            animation: spin 0.8s linear infinite;
            margin-right: 8px;
            vertical-align: middle;
        }

        @keyframes spin {
            to { transform: rotate(360deg); }
        }

        .end-message {
            text-align: center;
            padding: 24px;
            color: %[5]s;
            font-size: 14px;
            display: none;
        }

        .end-message.show {
            display: block;
        }

        .loading {
            text-align: center;
            padding: 40px;
            color: %[5]s;
        }

        .rss-link {
            text-align: center;
            padding: 20px;
        }

        .rss-link a {
            color: %[9]s;
            text-decoration: none;
            font-size: 14px;
        }

        @media (max-width: 768px) {
            .container {
                background-color: %[6]s;
            }

            .header {
                padding: 16px 20px;
            }

            .header h1 {
                font-size: 24px;
            }

            .stats {
                padding: 12px 16px;
                gap: 12px;
            }

            .stat {
                padding: 12px;
            }

            .stat-value {
                font-size: 24px;
            }

            .feed {
                padding: 0 0 80px;
            }

            .entry {
                margin-bottom: 12px;
                border-radius: 0;
                border-left: none;
                border-right: none;
            }
        }
`, bgColor, textColor, headerBg, borderColor, secondaryText, cardBg, hoverShadow, subtleBorder, accentColor)
}

// Helper functions for color manipulation
func isLightHexColor(hex string) bool {
	if len(hex) != 7 || hex[0] != '#' {
		return true // default to light
	}
	r, _ := strconv.ParseInt(hex[1:3], 16, 64)
	g, _ := strconv.ParseInt(hex[3:5], 16, 64)
	b, _ := strconv.ParseInt(hex[5:7], 16, 64)
	luminance := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 255
	return luminance > 0.5
}

func adjustHexBrightness(hex string, percent int) string {
	if len(hex) != 7 || hex[0] != '#' {
		return hex
	}
	r, _ := strconv.ParseInt(hex[1:3], 16, 64)
	g, _ := strconv.ParseInt(hex[3:5], 16, 64)
	b, _ := strconv.ParseInt(hex[5:7], 16, 64)

	r = clamp(r + int64(float64(r)*float64(percent)/100), 0, 255)
	g = clamp(g + int64(float64(g)*float64(percent)/100), 0, 255)
	b = clamp(b + int64(float64(b)*float64(percent)/100), 0, 255)

	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func blendColors(color1, color2 string, ratio float64) string {
	if len(color1) != 7 || color1[0] != '#' || len(color2) != 7 || color2[0] != '#' {
		return color1
	}
	r1, _ := strconv.ParseInt(color1[1:3], 16, 64)
	g1, _ := strconv.ParseInt(color1[3:5], 16, 64)
	b1, _ := strconv.ParseInt(color1[5:7], 16, 64)
	r2, _ := strconv.ParseInt(color2[1:3], 16, 64)
	g2, _ := strconv.ParseInt(color2[3:5], 16, 64)
	b2, _ := strconv.ParseInt(color2[5:7], 16, 64)

	r := int64(float64(r1)*(1-ratio) + float64(r2)*ratio)
	g := int64(float64(g1)*(1-ratio) + float64(g2)*ratio)
	b := int64(float64(b1)*(1-ratio) + float64(b2)*ratio)

	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func clamp(value, min, max int64) int64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func ifThen(condition bool, trueVal, falseVal int) int {
	if condition {
		return trueVal
	}
	return falseVal
}

// Authentication handlers
func handleLogin(w http.ResponseWriter, r *http.Request) {
	type LoginData struct {
		Error    string
		Redirect string
	}

	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "/admin"
	}

	if r.Method == http.MethodGet {
		tmpl := template.Must(template.New("login").Parse(loginTemplate))
		tmpl.Execute(w, LoginData{Redirect: redirect})
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		password := r.FormValue("password")
		redirectPath := r.FormValue("redirect")
		if redirectPath == "" {
			redirectPath = "/admin"
		}

		// Get admin password hash from database
		var passwordHash sql.NullString
		err = db.QueryRow("SELECT admin_password_hash FROM site_settings WHERE id = 1").Scan(&passwordHash)

		// Check if password is valid using bcrypt
		var passwordValid bool
		if err == nil && passwordHash.Valid && passwordHash.String != "" {
			err = bcrypt.CompareHashAndPassword([]byte(passwordHash.String), []byte(password))
			passwordValid = (err == nil)
		}

		if !passwordValid {
			tmpl := template.Must(template.New("login").Parse(loginTemplate))
			tmpl.Execute(w, LoginData{
				Error:    "Invalid password",
				Redirect: redirectPath,
			})
			return
		}

		// Create session
		token, err := createSession()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Set cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Path:     "/",
			MaxAge:   3600, // 60 minutes
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})

		http.Redirect(w, r, redirectPath, http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func handleSetPrivacyPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	password := r.FormValue("password")
	if password == "" {
		showMessage(w, r, "Password cannot be empty", "error")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		showMessage(w, r, "Failed to set password", "error")
		return
	}

	// Update viewer password in site_settings
	_, err = db.Exec(`UPDATE site_settings SET viewer_password_hash = ? WHERE id = 1`, string(hashedPassword))

	if err != nil {
		log.Printf("Error saving viewer password: %v", err)
		showMessage(w, r, "Failed to save viewer password", "error")
		return
	}

	showMessage(w, r, "Viewer password set successfully!", "success")
}

func handleRemovePrivacyPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Remove the viewer password from site_settings
	_, err := db.Exec("UPDATE site_settings SET viewer_password_hash = NULL WHERE id = 1")
	if err != nil {
		log.Printf("Error removing viewer password: %v", err)
		showMessage(w, r, "Failed to remove viewer password", "error")
		return
	}

	showMessage(w, r, "Your blog is now public! Password protection removed successfully.", "success")
}

func getSiteSettings() (SiteSettings, error) {
	var settings SiteSettings
	var viewerPasswordHash, adminPasswordHash, avatarPath, avatarPreference sql.NullString
	var customBgColor, customTextColor, customAccentColor sql.NullString

	err := db.QueryRow(`
		SELECT site_title, site_subtitle, user_initial, avatar_path, avatar_preference, site_theme,
		       viewer_password_hash, admin_password_hash, custom_bg_color, custom_text_color, custom_accent_color
		FROM site_settings WHERE id = 1
	`).Scan(&settings.SiteTitle, &settings.SiteSubtitle, &settings.UserInitial, &avatarPath, &avatarPreference,
		&settings.SiteTheme, &viewerPasswordHash, &adminPasswordHash, &customBgColor, &customTextColor, &customAccentColor)

	if err != nil {
		// If query fails, it might be due to missing columns - try running migrations
		if migErr := runDatabaseMigrations(); migErr == nil {
			// Retry the query after migrations
			err = db.QueryRow(`
				SELECT site_title, site_subtitle, user_initial, avatar_path, avatar_preference, site_theme,
				       viewer_password_hash, admin_password_hash, custom_bg_color, custom_text_color, custom_accent_color
				FROM site_settings WHERE id = 1
			`).Scan(&settings.SiteTitle, &settings.SiteSubtitle, &settings.UserInitial, &avatarPath, &avatarPreference,
				&settings.SiteTheme, &viewerPasswordHash, &adminPasswordHash, &customBgColor, &customTextColor, &customAccentColor)
		}
		if err != nil {
			return settings, err
		}
	}

	if avatarPath.Valid {
		settings.AvatarPath = avatarPath.String
	}

	if avatarPreference.Valid {
		settings.AvatarPreference = avatarPreference.String
	} else {
		settings.AvatarPreference = "initials" // default
	}

	// Set custom colors with defaults
	if customBgColor.Valid && customBgColor.String != "" {
		settings.CustomBgColor = customBgColor.String
	} else {
		settings.CustomBgColor = "#fafafa"
	}
	if customTextColor.Valid && customTextColor.String != "" {
		settings.CustomTextColor = customTextColor.String
	} else {
		settings.CustomTextColor = "#262626"
	}
	if customAccentColor.Valid && customAccentColor.String != "" {
		settings.CustomAccentColor = customAccentColor.String
	} else {
		settings.CustomAccentColor = "#0095f6"
	}

	settings.HasViewerPassword = viewerPasswordHash.Valid && viewerPasswordHash.String != ""
	settings.HasAdminPassword = adminPasswordHash.Valid && adminPasswordHash.String != ""

	return settings, nil
}

func handleSettings(w http.ResponseWriter, r *http.Request) {
	handleSettingsWithView(w, r, "site-info")
}

func handleSettingsAppearance(w http.ResponseWriter, r *http.Request) {
	handleSettingsWithView(w, r, "appearance")
}

func handleSettingsSecurity(w http.ResponseWriter, r *http.Request) {
	handleSettingsWithView(w, r, "security")
}

func handleSettingsBackup(w http.ResponseWriter, r *http.Request) {
	handleSettingsWithView(w, r, "backup")
}

func handleSettingsWithView(w http.ResponseWriter, r *http.Request, view string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	settings, err := getSiteSettings()
	if err != nil {
		log.Printf("Error fetching settings: %v", err)
		http.Error(w, "Error fetching settings", http.StatusInternalServerError)
		return
	}

	// Read flash messages from cookies
	var message, messageType string
	if cookie, err := r.Cookie("flash_message"); err == nil {
		message = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "flash_message",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}
	if cookie, err := r.Cookie("flash_type"); err == nil {
		messageType = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "flash_type",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}

	// Also check URL query params for messages (used by domain handlers)
	if msg := r.URL.Query().Get("success"); msg != "" {
		message = msg
		messageType = "success"
	} else if msg := r.URL.Query().Get("error"); msg != "" {
		message = msg
		messageType = "error"
	}

	// Get custom domain info
	customDomain, _ := getCustomDomain()
	instanceHostname := getInstanceHostname()

	// Custom domains only available for *.postastiq.com subdomains
	customDomainEnabled := isPostastiqSubdomain(instanceHostname)

	// Check if current request is from a postastiq.com subdomain (for showing enable button)
	canEnableCustomDomain := false
	if !customDomainEnabled {
		currentHost := getHostFromRequest(r)
		canEnableCustomDomain = isPostastiqSubdomain(currentHost)
	}

	// Set page title based on view
	pageTitles := map[string]string{
		"site-info":  "Site Info",
		"appearance": "Appearance",
		"security":   "Security",
		"backup":     "Backup",
	}
	pageTitle := pageTitles[view]
	if pageTitle == "" {
		pageTitle = "Settings"
	}

	data := SettingsPageData{
		Settings:              settings,
		EnableSubtitle:        enableSubtitle,
		Message:               message,
		MessageType:           messageType,
		CustomDomain:          customDomain,
		InstanceHostname:      instanceHostname,
		CustomDomainEnabled:   customDomainEnabled,
		CanEnableCustomDomain: canEnableCustomDomain,
		View:                  view,
		PageTitle:             pageTitle,
	}

	tmpl, err := template.New("settings").Parse(settingsTemplate)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func handleSettingsUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the section being updated
	section := r.FormValue("section")
	if section == "" {
		section = "site-info"
	}

	switch section {
	case "site-info":
		handleSiteInfoUpdate(w, r)
	case "appearance":
		handleAppearanceUpdate(w, r)
	case "security":
		handleSecurityUpdate(w, r)
	default:
		showSettingsMessage(w, r, "Invalid section", "error", section)
	}
}

func handleSiteInfoUpdate(w http.ResponseWriter, r *http.Request) {
	siteTitle := r.FormValue("site_title")
	siteSubtitle := r.FormValue("site_subtitle")
	userInitial := r.FormValue("user_initial")

	// Validate inputs
	if siteTitle == "" || siteSubtitle == "" || userInitial == "" {
		showSettingsMessage(w, r, "Site title, subtitle, and user initial are required", "error", "site-info")
		return
	}

	if len(userInitial) > 3 {
		showSettingsMessage(w, r, "User initial must be 1-3 characters", "error", "site-info")
		return
	}

	// Handle avatar upload
	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		avatarPath, err := validateAndResizeAvatar(file, header.Filename)
		if err != nil {
			log.Printf("Error saving avatar: %v", err)
			showSettingsMessage(w, r, "Failed to save avatar: "+err.Error(), "error", "site-info")
			return
		}

		// Update avatar path in database
		_, err = db.Exec("UPDATE site_settings SET avatar_path = ? WHERE id = 1", avatarPath)
		if err != nil {
			log.Printf("Error updating avatar path: %v", err)
			showSettingsMessage(w, r, "Failed to update avatar", "error", "site-info")
			return
		}
	}

	// Update site info settings
	_, err = db.Exec(`
		UPDATE site_settings
		SET site_title = ?, site_subtitle = ?, user_initial = ?
		WHERE id = 1
	`, siteTitle, siteSubtitle, userInitial)

	if err != nil {
		log.Printf("Error updating settings: %v", err)
		showSettingsMessage(w, r, "Failed to update settings", "error", "site-info")
		return
	}

	showSettingsMessage(w, r, "Site info updated successfully!", "success", "site-info")
}

func handleAppearanceUpdate(w http.ResponseWriter, r *http.Request) {
	siteTheme := r.FormValue("site_theme")
	avatarPreference := r.FormValue("avatar_preference")
	customBgColor := r.FormValue("custom_bg_color")
	customTextColor := r.FormValue("custom_text_color")
	customAccentColor := r.FormValue("custom_accent_color")

	// Set default avatar preference if not provided
	if avatarPreference == "" {
		avatarPreference = "initials"
	}

	// Set default custom colors if not provided
	if customBgColor == "" {
		customBgColor = "#fafafa"
	}
	if customTextColor == "" {
		customTextColor = "#262626"
	}
	if customAccentColor == "" {
		customAccentColor = "#0095f6"
	}

	// Validate color format (basic hex color validation)
	colorRegex := regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)
	if !colorRegex.MatchString(customBgColor) || !colorRegex.MatchString(customTextColor) || !colorRegex.MatchString(customAccentColor) {
		showSettingsMessage(w, r, "Invalid color format. Please use hex colors (e.g., #ffffff)", "error", "appearance")
		return
	}

	// Update appearance settings including custom colors
	_, err := db.Exec(`
		UPDATE site_settings
		SET site_theme = ?, avatar_preference = ?, custom_bg_color = ?, custom_text_color = ?, custom_accent_color = ?
		WHERE id = 1
	`, siteTheme, avatarPreference, customBgColor, customTextColor, customAccentColor)

	if err != nil {
		log.Printf("Error updating appearance settings: %v", err)
		showSettingsMessage(w, r, "Failed to update appearance settings", "error", "appearance")
		return
	}

	showSettingsMessage(w, r, "Appearance updated successfully!", "success", "appearance")
}

func handleSecurityUpdate(w http.ResponseWriter, r *http.Request) {
	adminPassword := r.FormValue("admin_password")
	adminPasswordConfirm := r.FormValue("admin_password_confirm")
	viewerPassword := r.FormValue("viewer_password")
	removeViewerPassword := r.FormValue("remove_viewer_password")

	// Handle admin password update
	if adminPassword != "" {
		if adminPassword != adminPasswordConfirm {
			showSettingsMessage(w, r, "Admin passwords do not match", "error", "security")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing admin password: %v", err)
			showSettingsMessage(w, r, "Failed to update admin password", "error", "security")
			return
		}

		// Update password and clear password_change_required flag
		_, err = db.Exec("UPDATE site_settings SET admin_password_hash = ?, password_change_required = 0 WHERE id = 1", string(hashedPassword))
		if err != nil {
			log.Printf("Error updating admin password: %v", err)
			showSettingsMessage(w, r, "Failed to update admin password", "error", "security")
			return
		}
	}

	// Handle viewer password
	if removeViewerPassword == "true" {
		_, err := db.Exec("UPDATE site_settings SET viewer_password_hash = NULL WHERE id = 1")
		if err != nil {
			log.Printf("Error removing viewer password: %v", err)
			showSettingsMessage(w, r, "Failed to remove viewer password", "error", "security")
			return
		}
	} else if viewerPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(viewerPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing viewer password: %v", err)
			showSettingsMessage(w, r, "Failed to update viewer password", "error", "security")
			return
		}

		_, err = db.Exec("UPDATE site_settings SET viewer_password_hash = ? WHERE id = 1", string(hashedPassword))
		if err != nil {
			log.Printf("Error updating viewer password: %v", err)
			showSettingsMessage(w, r, "Failed to update viewer password", "error", "security")
			return
		}
	}

	showSettingsMessage(w, r, "Security settings updated successfully!", "success", "security")
}

func handleBackup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get database path
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/app/data/blog.db"
	}

	// Set headers for ZIP download
	timestamp := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("postastiq-backup-%s.zip", timestamp)
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Create ZIP writer
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	// Add database file to ZIP
	dbFile, err := os.Open(dbPath)
	if err != nil {
		log.Printf("Error opening database for backup: %v", err)
		http.Error(w, "Failed to create backup", http.StatusInternalServerError)
		return
	}
	defer dbFile.Close()

	dbInfo, err := dbFile.Stat()
	if err != nil {
		log.Printf("Error getting database info: %v", err)
		http.Error(w, "Failed to create backup", http.StatusInternalServerError)
		return
	}

	dbHeader, err := zip.FileInfoHeader(dbInfo)
	if err != nil {
		log.Printf("Error creating ZIP header for database: %v", err)
		http.Error(w, "Failed to create backup", http.StatusInternalServerError)
		return
	}
	dbHeader.Name = "blog.db"
	dbHeader.Method = zip.Deflate

	dbWriter, err := zipWriter.CreateHeader(dbHeader)
	if err != nil {
		log.Printf("Error creating ZIP entry for database: %v", err)
		http.Error(w, "Failed to create backup", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(dbWriter, dbFile)
	if err != nil {
		log.Printf("Error writing database to ZIP: %v", err)
		return
	}

	// Add uploads directory to ZIP
	err = filepath.Walk(uploadsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the uploads directory itself
		if path == uploadsDir {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(uploadsDir, path)
		if err != nil {
			return err
		}

		// Create ZIP header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.Join("uploads", relPath)

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("Error adding uploads to ZIP: %v", err)
	}

	log.Printf("Backup created successfully: %s", filename)
}

func handleRestore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Use small memory buffer (4MB) - excess streams to temp files automatically
	// This keeps memory usage low for large backup files
	err := r.ParseMultipartForm(4 << 20)
	if err != nil {
		log.Printf("Error parsing restore form: %v", err)
		showSettingsMessage(w, r, "Failed to parse upload", "error", "backup")
		return
	}

	// Get uploaded file
	file, header, err := r.FormFile("backup_file")
	if err != nil {
		log.Printf("Error getting backup file: %v", err)
		showSettingsMessage(w, r, "No backup file provided", "error", "backup")
		return
	}

	// Validate file extension
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".zip") {
		file.Close()
		showSettingsMessage(w, r, "Invalid file type. Please upload a ZIP file.", "error", "backup")
		return
	}

	// Stream upload to temp file using small buffer to minimize memory
	tempFile, err := os.CreateTemp("", "postastiq-restore-*.zip")
	if err != nil {
		file.Close()
		log.Printf("Error creating temp file: %v", err)
		showSettingsMessage(w, r, "Failed to process upload", "error", "backup")
		return
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath)

	// Use small copy buffer (32KB) to minimize memory usage
	copyBuf := make([]byte, 32*1024)
	written, err := io.CopyBuffer(tempFile, file, copyBuf)
	tempFile.Close()
	file.Close()

	// Release multipart form data to free memory
	r.MultipartForm.RemoveAll()
	runtime.GC()

	if err != nil {
		log.Printf("Error writing to temp file: %v", err)
		showSettingsMessage(w, r, "Failed to save upload", "error", "backup")
		return
	}
	log.Printf("Received backup file: %s (%d bytes)", header.Filename, written)

	// Get paths
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/app/data/blog.db"
	}
	tempDBPath := dbPath + ".restore-temp"

	// Process ZIP file with minimal memory: extract database first
	// We open/close the ZIP reader for each phase to allow GC between phases
	err = extractDatabaseFromZip(tempPath, tempDBPath, copyBuf)
	if err != nil {
		if err.Error() == "missing blog.db" {
			showSettingsMessage(w, r, "Invalid backup: missing blog.db", "error", "backup")
		} else {
			log.Printf("Error extracting database: %v", err)
			showSettingsMessage(w, r, "Failed to extract database", "error", "backup")
		}
		return
	}

	// Close current database connection before replacing file
	if db != nil {
		db.Close()
		db = nil
	}

	// Move temp database to final location (atomic on same filesystem)
	if err := os.Rename(tempDBPath, dbPath); err != nil {
		log.Printf("Error replacing database file: %v", err)
		os.Remove(tempDBPath)
		reconnectDB(dbPath)
		showSettingsMessage(w, r, "Failed to replace database", "error", "backup")
		return
	}

	// Reconnect to the restored database immediately
	if err := reconnectDB(dbPath); err != nil {
		log.Printf("Error reconnecting to restored database: %v", err)
		showSettingsMessage(w, r, "Database restored but reconnection failed. Please restart the server.", "error", "backup")
		return
	}

	// Run database migrations to add any missing columns from newer versions
	if err := runDatabaseMigrations(); err != nil {
		log.Printf("Warning: some database migrations failed after restore: %v", err)
		// Continue anyway - the core restore worked
	}

	// Force GC before extracting uploads
	runtime.GC()

	// Extract upload files (database is already reconnected)
	err = extractUploadsFromZip(tempPath, uploadsDir, copyBuf)
	if err != nil {
		log.Printf("Warning: some upload files may not have been restored: %v", err)
	}

	log.Printf("Backup restored successfully from: %s", header.Filename)
	showSettingsMessage(w, r, "Backup restored successfully!", "success", "backup")
}

// extractDatabaseFromZip extracts only the blog.db file from the ZIP
func extractDatabaseFromZip(zipPath, destPath string, buf []byte) error {
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("invalid ZIP file: %w", err)
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		if f.Name == "blog.db" {
			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("failed to open blog.db in ZIP: %w", err)
			}

			outFile, err := os.Create(destPath)
			if err != nil {
				rc.Close()
				return fmt.Errorf("failed to create database file: %w", err)
			}

			_, err = io.CopyBuffer(outFile, rc, buf)
			outFile.Close()
			rc.Close()

			if err != nil {
				os.Remove(destPath)
				return fmt.Errorf("failed to write database file: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("missing blog.db")
}

// extractUploadsFromZip extracts upload files from the ZIP one at a time
func extractUploadsFromZip(zipPath, uploadsDir string, buf []byte) error {
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open ZIP: %w", err)
	}
	defer zipReader.Close()

	var lastErr error
	for _, f := range zipReader.File {
		if !strings.HasPrefix(f.Name, "uploads/") {
			continue
		}

		relPath := strings.TrimPrefix(f.Name, "uploads/")
		if relPath == "" {
			continue
		}

		targetPath := filepath.Join(uploadsDir, relPath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(targetPath, 0755)
			continue
		}

		// Ensure parent directory exists
		os.MkdirAll(filepath.Dir(targetPath), 0755)

		// Extract file with minimal memory
		if err := extractSingleFile(f, targetPath, buf); err != nil {
			log.Printf("Error extracting %s: %v", f.Name, err)
			lastErr = err
		}
	}

	return lastErr
}

// extractSingleFile extracts a single file from ZIP using provided buffer
func extractSingleFile(f *zip.File, destPath string, buf []byte) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	outFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.CopyBuffer(outFile, rc, buf)
	return err
}

// reconnectDB closes the existing connection and opens a new one
func reconnectDB(dbPath string) error {
	// Close existing connection if any
	if db != nil {
		db.Close()
		db = nil
	}

	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Configure connection pool for better reliability
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	return db.Ping()
}

func handleViewerAuth(w http.ResponseWriter, r *http.Request) {
	// Get settings from database
	settings, err := getSiteSettings()
	if err != nil {
		log.Printf("Error getting site settings: %v", err)
		settings = SiteSettings{
			SiteTitle: "My Blog",
		}
	}

	if r.Method == http.MethodGet {
		data := struct {
			SiteTitle string
			Error     string
			Redirect  string
		}{
			SiteTitle: settings.SiteTitle,
			Redirect:  r.URL.Query().Get("redirect"),
		}

		tmpl, err := template.New("viewer-password").Parse(viewerPasswordTemplate)
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
		return
	}

	if r.Method == http.MethodPost {
		password := r.FormValue("password")
		redirect := r.FormValue("redirect")

		// Get stored password hash from site_settings (not privacy_settings)
		var passwordHash sql.NullString
		err := db.QueryRow("SELECT viewer_password_hash FROM site_settings WHERE id = 1").Scan(&passwordHash)
		if err != nil || !passwordHash.Valid {
			// No password set, redirect to home or specified page
			if redirect != "" {
				http.Redirect(w, r, redirect, http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			return
		}

		// Verify password
		err = bcrypt.CompareHashAndPassword([]byte(passwordHash.String), []byte(password))
		if err != nil {
			data := struct {
				SiteTitle string
				Error     string
				Redirect  string
			}{
				SiteTitle: settings.SiteTitle,
				Error:     "Incorrect password",
				Redirect:  redirect,
			}

			tmpl, _ := template.New("viewer-password").Parse(viewerPasswordTemplate)
			tmpl.Execute(w, data)
			return
		}

		// Password correct, set viewer session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "viewer_authenticated",
			Value:    "true",
			Path:     "/",
			MaxAge:   86400 * 30, // 30 days
			HttpOnly: true,
		})

		// Redirect to original page or home
		if redirect != "" && strings.HasPrefix(redirect, "/") {
			http.Redirect(w, r, redirect, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func requireViewerAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if viewer password is set in site_settings
		var passwordHash sql.NullString
		err := db.QueryRow("SELECT viewer_password_hash FROM site_settings WHERE id = 1").Scan(&passwordHash)
		if err != nil || !passwordHash.Valid || passwordHash.String == "" {
			// No password set, allow access
			next(w, r)
			return
		}

		// Check if user is authenticated
		cookie, err := r.Cookie("viewer_authenticated")
		if err != nil || cookie.Value != "true" {
			// Not authenticated, redirect to password page with return URL
			redirectURL := "/viewer-auth?redirect=" + template.URLQueryEscaper(r.URL.Path)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		deleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func isPasswordChangeRequired() bool {
	var required sql.NullInt64
	err := db.QueryRow("SELECT password_change_required FROM site_settings WHERE id = 1").Scan(&required)
	if err != nil {
		return false
	}
	return required.Valid && required.Int64 == 1
}

func handleChangePassword(w http.ResponseWriter, r *http.Request) {
	// Verify user is authenticated
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only allow access if password change is required (first login)
	// Otherwise redirect to settings page for password changes
	if !isPasswordChangeRequired() {
		http.Redirect(w, r, "/admin/settings", http.StatusSeeOther)
		return
	}

	type ChangePasswordData struct {
		Error   string
		Success string
	}

	if r.Method == http.MethodGet {
		tmpl := template.Must(template.New("change-password").Parse(changePasswordTemplate))
		tmpl.Execute(w, ChangePasswordData{})
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		newPassword := r.FormValue("new_password")
		confirmPassword := r.FormValue("confirm_password")

		// Validate passwords
		if newPassword == "" {
			tmpl := template.Must(template.New("change-password").Parse(changePasswordTemplate))
			tmpl.Execute(w, ChangePasswordData{Error: "Password cannot be empty"})
			return
		}

		if len(newPassword) < 8 {
			tmpl := template.Must(template.New("change-password").Parse(changePasswordTemplate))
			tmpl.Execute(w, ChangePasswordData{Error: "Password must be at least 8 characters long"})
			return
		}

		if newPassword != confirmPassword {
			tmpl := template.Must(template.New("change-password").Parse(changePasswordTemplate))
			tmpl.Execute(w, ChangePasswordData{Error: "Passwords do not match"})
			return
		}

		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			tmpl := template.Must(template.New("change-password").Parse(changePasswordTemplate))
			tmpl.Execute(w, ChangePasswordData{Error: "Failed to update password"})
			return
		}

		// Update password and clear the password_change_required flag
		_, err = db.Exec("UPDATE site_settings SET admin_password_hash = ?, password_change_required = 0 WHERE id = 1", string(hashedPassword))
		if err != nil {
			log.Printf("Error updating password: %v", err)
			tmpl := template.Must(template.New("change-password").Parse(changePasswordTemplate))
			tmpl.Execute(w, ChangePasswordData{Error: "Failed to update password"})
			return
		}

		log.Println("Admin password changed successfully, password_change_required cleared")

		// Redirect to admin panel
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// HealthResponse represents the JSON response for the health endpoint
type HealthResponse struct {
	Status    string `json:"status"`
	Bootstrap bool   `json:"bootstrap"`
}

// adminRouter handles all /admin routes with a single prefix handler
// This ensures all admin paths bypass viewer auth and only require admin auth
func adminRouter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/admin" || path == "/admin/":
		handleAdminIndex(w, r)
	case path == "/admin/create":
		handleCreate(w, r)
	case path == "/admin/update":
		handleUpdate(w, r)
	case path == "/admin/delete":
		handleDelete(w, r)
	case path == "/admin/entries":
		handleGetEntriesForAdmin(w, r)
	case path == "/admin/privacy/set":
		handleSetPrivacyPassword(w, r)
	case path == "/admin/privacy/remove":
		handleRemovePrivacyPassword(w, r)
	case path == "/admin/settings":
		handleSettings(w, r)
	case path == "/admin/settings/appearance":
		handleSettingsAppearance(w, r)
	case path == "/admin/settings/security":
		handleSettingsSecurity(w, r)
	case path == "/admin/settings/backup":
		handleSettingsBackup(w, r)
	case path == "/admin/settings/update":
		handleSettingsUpdate(w, r)
	case path == "/admin/backup":
		handleBackup(w, r)
	case path == "/admin/restore":
		handleRestore(w, r)
	case path == "/admin/domain/add":
		handleDomainAdd(w, r)
	case path == "/admin/domain/verify":
		handleDomainVerify(w, r)
	case path == "/admin/domain/activate":
		handleDomainActivate(w, r)
	case path == "/admin/domain/remove":
		handleDomainRemove(w, r)
	case path == "/admin/domain/status":
		handleDomainStatus(w, r)
	case path == "/admin/domain/set-hostname":
		handleSetInstanceHostname(w, r)
	case path == "/admin/domain/enable":
		handleEnableCustomDomain(w, r)
	default:
		http.NotFound(w, r)
	}
}

// ============================================================================
// Custom Domain HTTP Handlers
// ============================================================================

func handleDomainAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	domain := strings.TrimSpace(r.FormValue("domain"))
	if domain == "" {
		http.Redirect(w, r, "/admin/settings?error=Domain+cannot+be+empty", http.StatusSeeOther)
		return
	}

	cd, err := addCustomDomain(domain)
	if err != nil {
		http.Redirect(w, r, "/admin/settings?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	log.Printf("Custom domain added: %s (token: %s)", cd.Domain, cd.VerificationToken[:20]+"...")
	http.Redirect(w, r, "/admin/settings?success=Domain+added.+Please+configure+DNS+records.", http.StatusSeeOther)
}

func handleDomainVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := verifyCustomDomain()
	if err != nil {
		http.Redirect(w, r, "/admin/settings?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	// Automatically activate after verification
	err = activateCustomDomain()
	if err != nil {
		log.Printf("Domain verified but activation failed: %v", err)
		http.Redirect(w, r, "/admin/settings?success=Domain+verified.+Activation+pending.", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/settings?success=Domain+verified+and+activated!", http.StatusSeeOther)
}

func handleDomainActivate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := activateCustomDomain()
	if err != nil {
		http.Redirect(w, r, "/admin/settings?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/settings?success=Domain+activated!", http.StatusSeeOther)
}

func handleDomainRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := removeCustomDomain()
	if err != nil {
		http.Redirect(w, r, "/admin/settings?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/settings?success=Domain+removed.", http.StatusSeeOther)
}

func handleDomainStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cd, err := getCustomDomain()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if cd == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"configured": false,
		})
		return
	}

	response := map[string]interface{}{
		"configured":    true,
		"domain":        cd.Domain,
		"verified":      cd.VerifiedAt.Valid,
		"activated":     cd.ActivatedAt.Valid,
		"attempts":      cd.VerificationAttempts,
		"attemptsLeft":  5 - cd.VerificationAttempts,
	}

	if cd.VerifiedAt.Valid {
		response["verifiedAt"] = cd.VerifiedAt.Time.Format(time.RFC3339)
	}
	if cd.ActivatedAt.Valid {
		response["activatedAt"] = cd.ActivatedAt.Time.Format(time.RFC3339)
	}
	if cd.LastVerifiedAt.Valid {
		response["lastVerifiedAt"] = cd.LastVerifiedAt.Time.Format(time.RFC3339)
	}

	json.NewEncoder(w).Encode(response)
}

func handleSetInstanceHostname(w http.ResponseWriter, r *http.Request) {
	// Deprecated - redirect to enable endpoint
	http.Redirect(w, r, "/admin/domain/enable", http.StatusSeeOther)
}

// handleEnableCustomDomain enables custom domain feature by capturing hostname from current request
func handleEnableCustomDomain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if already enabled
	existingHostname := getInstanceHostname()
	if existingHostname != "" {
		http.Redirect(w, r, "/admin/settings?error=Custom+domains+already+enabled", http.StatusSeeOther)
		return
	}

	// Get hostname from current request
	currentHost := getHostFromRequest(r)
	if currentHost == "" {
		http.Redirect(w, r, "/admin/settings?error=Could+not+detect+hostname", http.StatusSeeOther)
		return
	}

	// Only allow postastiq.com subdomains
	if !isPostastiqSubdomain(currentHost) {
		http.Redirect(w, r, "/admin/settings?error=Custom+domains+only+available+for+postastiq.com+subdomains", http.StatusSeeOther)
		return
	}

	// Store the hostname
	err := setInstanceHostname(currentHost)
	if err != nil {
		log.Printf("Failed to enable custom domains: %v", err)
		http.Redirect(w, r, "/admin/settings?error=Failed+to+enable+custom+domains", http.StatusSeeOther)
		return
	}

	log.Printf("Custom domains enabled for: %s", currentHost)
	http.Redirect(w, r, "/admin/settings?success=Custom+domains+enabled+for+"+url.QueryEscape(currentHost), http.StatusSeeOther)
}

// ============================================================================
// End Custom Domain HTTP Handlers
// ============================================================================

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Check if admin password has been set (bootstrap complete)
	var passwordHash sql.NullString
	err := db.QueryRow("SELECT admin_password_hash FROM site_settings WHERE id = 1").Scan(&passwordHash)

	bootstrapComplete := err == nil && passwordHash.Valid && passwordHash.String != ""

	if !bootstrapComplete {
		// Still initializing - admin password not yet set
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(HealthResponse{
			Status:    "initializing",
			Bootstrap: false,
		})
		return
	}

	// App is ready
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HealthResponse{
		Status:    "ready",
		Bootstrap: true,
	})
}

func handleRSSFeed(w http.ResponseWriter, r *http.Request) {
	// Get site settings
	settings, err := getSiteSettings()
	if err != nil {
		log.Printf("Error getting site settings for RSS: %v", err)
		settings = SiteSettings{
			SiteTitle:    "My Blog",
			SiteSubtitle: "A Personal Blog",
		}
	}

	// Get latest 20 entries for RSS feed
	entries, _, err := getEntries(0, 20)
	if err != nil {
		http.Error(w, "Error generating RSS feed", http.StatusInternalServerError)
		log.Printf("Error getting entries for RSS: %v", err)
		return
	}

	// Build RSS feed XML
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")

	// Get the base URL from the request
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, r.Host)

	// Start RSS feed
	fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintf(w, `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">`)
	fmt.Fprintf(w, `<channel>`)
	fmt.Fprintf(w, `<title>%s</title>`, html.EscapeString(settings.SiteTitle))
	fmt.Fprintf(w, `<link>%s</link>`, html.EscapeString(baseURL))
	fmt.Fprintf(w, `<description>%s</description>`, html.EscapeString(settings.SiteSubtitle))
	fmt.Fprintf(w, `<language>en-us</language>`)
	fmt.Fprintf(w, `<atom:link href="%s/rss" rel="self" type="application/rss+xml" />`, html.EscapeString(baseURL))

	// Add lastBuildDate (current time)
	fmt.Fprintf(w, `<lastBuildDate>%s</lastBuildDate>`, time.Now().UTC().Format(time.RFC1123Z))

	// Add items
	for _, entry := range entries {
		fmt.Fprintf(w, `<item>`)

		// Title - use content preview if no title
		title := entry.Title
		if title == "" {
			// Use first 60 characters of content as title
			contentRunes := []rune(string(entry.Content))
			if len(contentRunes) > 60 {
				title = string(contentRunes[:60]) + "..."
			} else {
				title = string(entry.Content)
			}
		}
		fmt.Fprintf(w, `<title>%s</title>`, html.EscapeString(title))

		// Link to individual post
		postURL := fmt.Sprintf("%s/posts/%s/", baseURL, entry.Slug)
		fmt.Fprintf(w, `<link>%s</link>`, html.EscapeString(postURL))

		// GUID (unique identifier)
		fmt.Fprintf(w, `<guid isPermaLink="true">%s</guid>`, html.EscapeString(postURL))

		// Description (text content only - media handled via enclosure tag)
		fmt.Fprintf(w, `<description><![CDATA[%s]]></description>`, string(entry.Content))

		// Enclosure for photos (standard RSS 2.0 media handling)
		if entry.HasPhoto && entry.Photo != "" {
			photoURL := fmt.Sprintf("%s%s", baseURL, string(entry.Photo))
			// Determine MIME type from file extension
			mimeType := "image/jpeg" // default
			photoStr := string(entry.Photo)
			if strings.HasSuffix(strings.ToLower(photoStr), ".png") {
				mimeType = "image/png"
			} else if strings.HasSuffix(strings.ToLower(photoStr), ".gif") {
				mimeType = "image/gif"
			} else if strings.HasSuffix(strings.ToLower(photoStr), ".webp") {
				mimeType = "image/webp"
			}
			fmt.Fprintf(w, `<enclosure url="%s" type="%s" length="0" />`, html.EscapeString(photoURL), mimeType)
		}

		// Publication date
		fmt.Fprintf(w, `<pubDate>%s</pubDate>`, entry.CreatedAt.UTC().Format(time.RFC1123Z))

		fmt.Fprintf(w, `</item>`)
	}

	fmt.Fprintf(w, `</channel>`)
	fmt.Fprintf(w, `</rss>`)
}

// resetAdminPassword resets the admin password from CLI
func resetAdminPassword(newPassword string) {
	if newPassword == "" {
		fmt.Println("Error: Password cannot be empty")
		os.Exit(1)
	}

	// Initialize database connection
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/app/data/blog.db"
	}

	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		os.Exit(1)
	}

	// Update password and require change on next login
	_, err = database.Exec("UPDATE site_settings SET admin_password_hash = ?, password_change_required = 1 WHERE id = 1", string(hashedPassword))
	if err != nil {
		fmt.Printf("Error updating password: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Admin password has been reset successfully.")
	fmt.Println("Password change will be required on next login.")
}

// enableCustomDomainCLI enables custom domain feature via CLI (platform admin use)
func enableCustomDomainCLI(hostname string) {
	hostname = strings.ToLower(strings.TrimSpace(hostname))

	if hostname == "" {
		fmt.Println("Error: hostname cannot be empty")
		os.Exit(1)
	}

	if !isPostastiqSubdomain(hostname) {
		fmt.Printf("Error: '%s' is not a postastiq.com subdomain\n", hostname)
		fmt.Println("Custom domains can only be enabled for *.postastiq.com subdomains")
		os.Exit(1)
	}

	// Initialize database
	if err := initDB(); err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Check if already set
	existing := getInstanceHostname()
	if existing != "" {
		fmt.Printf("Custom domains already enabled for: %s\n", existing)
		os.Exit(0)
	}

	// Set the hostname
	err := setInstanceHostname(hostname)
	if err != nil {
		fmt.Printf("Error enabling custom domains: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Custom domains enabled for: %s\n", hostname)
}

// printHelp prints CLI usage information
func printHelp() {
	fmt.Println("Postastiq - Self-hosted micro-blogging platform")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  postastiq                                   Start the web server")
	fmt.Println("  postastiq --reset-password <pwd>            Reset admin password")
	fmt.Println("  postastiq --enable-custom-domain <hostname> Enable custom domains (platform admin)")
	fmt.Println("  postastiq --help                            Show this help message")
	fmt.Println("")
	fmt.Println("Environment Variables:")
	fmt.Println("  PORT              Server port (default: 8080)")
	fmt.Println("  DB_PATH           Database path (default: /app/data/blog.db)")
	fmt.Println("  UPLOADS_DIR       Uploads directory (default: /app/data/uploads)")
	fmt.Println("  ADMIN_PASSWORD    Initial admin password (default: admin)")
}

func main() {
	// Check for CLI commands
	if len(os.Args) >= 3 && os.Args[1] == "--reset-password" {
		resetAdminPassword(os.Args[2])
		return
	}

	if len(os.Args) >= 3 && os.Args[1] == "--enable-custom-domain" {
		enableCustomDomainCLI(os.Args[2])
		return
	}

	if len(os.Args) >= 2 && os.Args[1] == "--help" {
		printHelp()
		return
	}

	if err := initDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()

	// Sync custom domain routes with Caddy on startup (runs in background)
	go syncCustomDomainWithCaddy()

	// Start session cleanup goroutine
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			cleanupExpiredSessions()
		}
	}()

	// Start domain re-validation cron (weekly)
	startDomainRevalidationCron()

	// Serve uploaded files
	fs := http.FileServer(http.Dir(uploadsDir))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", fs))

	// Authentication routes
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/viewer-auth", handleViewerAuth)
	http.HandleFunc("/change-password", requireAuth(handleChangePassword))

	// Public endpoints (no auth required)
	http.HandleFunc("/health", handleHealth)

	// RSS feed (protected by privacy password if set)
	http.HandleFunc("/rss", requireViewerAuth(handleRSSFeed))

	// Blog viewer routes (protected by privacy password if set)
	http.HandleFunc("/", requireViewerAuth(handleBlogFeed))
	http.HandleFunc("/api/entries", requireViewerAuth(handleAPIEntries))
	http.HandleFunc("/posts/", requireViewerAuth(handleSinglePost))

	// Protected admin routes - use prefix pattern to catch all /admin* paths
	// This ensures admin routes bypass viewer auth (only require admin auth)
	http.HandleFunc("/admin", requireAuth(adminRouter))
	http.HandleFunc("/admin/", requireAuth(adminRouter))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Blog server starting on port %s", port)
	log.Printf("Viewer: http://localhost:%s/", port)
	log.Printf("Admin: http://localhost:%s/admin", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
