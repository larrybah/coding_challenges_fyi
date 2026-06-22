package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// DATA STRUCTURES
// ============================================================================

// Note represents a personal note with metadata.
// This is the core data structure for our API.
type Note struct {
	ID        int       `json:"id"`        // Unique identifier for the note
	Title     string    `json:"title"`     // Title of the note
	Content   string    `json:"content"`   // Main content of the note
	CreatedAt time.Time `json:"created_at"` // Timestamp when note was created
	UpdatedAt time.Time `json:"updated_at"` // Timestamp when note was last updated
}

// CreateNoteRequest represents the JSON payload for creating a new note.
// Used when clients POST a new note to the API.
type CreateNoteRequest struct {
	Title   string `json:"title"`   // Title is required
	Content string `json:"content"` // Content is required
}

// UpdateNoteRequest represents the JSON payload for updating an existing note.
// Both fields are optional - clients can update just title, just content, or both.
type UpdateNoteRequest struct {
	Title   *string `json:"title"`   // Pointer allows detecting if field was provided
	Content *string `json:"content"` // Pointer allows detecting if field was provided
}

// ErrorResponse represents a standard error response from the API.
// Ensures consistent error formatting across all endpoints.
type ErrorResponse struct {
	Error   string `json:"error"`   // Error message
	Details string `json:"details"` // Additional details (optional)
}

// SuccessResponse represents a standard success response from the API.
// Used for operations that return a single note.
type SuccessResponse struct {
	Data Note `json:"data"` // The note data
}

// NotesStore manages in-memory storage of notes.
// Includes a mutex for thread-safe concurrent access.
type NotesStore struct {
	mu    sync.RWMutex // RWMutex allows multiple readers but exclusive writers
	notes map[int]*Note // Map stores notes with ID as key for O(1) access
	nextID int           // Counter for generating unique IDs
}

// ============================================================================
// GLOBAL STATE
// ============================================================================

// store is the global in-memory database for our notes.
// In production, this would be replaced with a real database (PostgreSQL, MongoDB, etc.)
var store = &NotesStore{
	notes: make(map[int]*Note),
	nextID: 1,
}

// ============================================================================
// MIDDLEWARE
// ============================================================================

// loggingMiddleware logs information about each HTTP request.
// This is a middleware pattern - it wraps other handlers and can execute code before/after them.
// Learn about middleware for handling cross-cutting concerns like logging, authentication, etc.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the start time to calculate request duration
		startTime := time.Now()

		// Log the incoming request
		log.Printf("[%s] %s %s - User-Agent: %s", 
			r.Method, 
			r.URL.Path, 
			r.RemoteAddr, 
			r.UserAgent(),
		)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// Log the request duration after it's processed
		duration := time.Since(startTime)
		log.Printf("Request completed in %v", duration)
	})
}

// jsonContentTypeMiddleware ensures all responses have the correct Content-Type header.
// This tells clients to parse the response as JSON.
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Content-Type header to JSON for all responses
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// ============================================================================
// HANDLER FUNCTIONS
// ============================================================================

// handleCreateNote handles POST requests to create a new note.
// REST principle: POST is used for creating new resources.
// 
// Expected request body (JSON):
// {
//   "title": "My Note Title",
//   "content": "Note content here"
// }
//
// Returns: 201 Created with the new note, or 400 Bad Request if validation fails.
func handleCreateNote(w http.ResponseWriter, r *http.Request) {
	// Verify the HTTP method is POST
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", "Only POST is allowed")
		return
	}

	// Decode the JSON request body into our CreateNoteRequest struct
	var req CreateNoteRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	// Validate the input
	if err := validateCreateNoteRequest(req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	// Create a new note
	newNote := Note{
		ID:        store.getNextID(),
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store the note
	store.addNote(&newNote)

	// Respond with 201 Created status and the new note
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SuccessResponse{Data: newNote})
}

// handleGetNotes handles GET requests to retrieve all notes.
// REST principle: GET is used for retrieving resources without modifying them.
//
// Returns: 200 OK with an array of all notes.
func handleGetNotes(w http.ResponseWriter, r *http.Request) {
	// Verify the HTTP method is GET
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", "Only GET is allowed")
		return
	}

	// Extract the note ID from the URL path
	// The route is /notes/:id, so we parse the ID parameter
	parts := strings.Split(r.URL.Path, "/")
	
	// If there are exactly 3 parts ("/", "notes", "id"), they want a specific note
	if len(parts) == 3 && parts[2] != "" {
		// Handle GET /notes/:id - retrieve a single note
		handleGetNoteByID(w, r, parts[2])
		return
	}

	// Handle GET /notes - retrieve all notes
	notes := store.getAllNotes()
	w.WriteHeader(http.StatusOK)
	// Return as an array wrapped in JSON
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": notes,
	})
}

// handleGetNoteByID retrieves a single note by its ID.
// This is a helper function called from handleGetNotes.
func handleGetNoteByID(w http.ResponseWriter, r *http.Request, idStr string) {
	// Parse the ID from string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", "ID must be an integer")
		return
	}

	// Retrieve the note from storage
	note, found := store.getNoteByID(id)
	if !found {
		respondWithError(w, http.StatusNotFound, "Not found", fmt.Sprintf("Note with ID %d not found", id))
		return
	}

	// Return the note
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{Data: *note})
}

// handleUpdateNote handles PATCH requests to update an existing note.
// REST principle: PATCH is used for partial updates to existing resources.
//
// Expected request body (JSON) - all fields optional:
// {
//   "title": "Updated Title",
//   "content": "Updated content"
// }
//
// Returns: 200 OK with the updated note, or 404 Not Found if note doesn't exist.
func handleUpdateNote(w http.ResponseWriter, r *http.Request) {
	// Verify the HTTP method is PATCH
	if r.Method != http.MethodPatch {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", "Only PATCH is allowed")
		return
	}

	// Extract the ID from the URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		respondWithError(w, http.StatusBadRequest, "Missing ID", "Note ID is required")
		return
	}

	// Parse the ID
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", "ID must be an integer")
		return
	}

	// Decode the request body
	var req UpdateNoteRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	// Retrieve the existing note
	note, found := store.getNoteByID(id)
	if !found {
		respondWithError(w, http.StatusNotFound, "Not found", fmt.Sprintf("Note with ID %d not found", id))
		return
	}

	// Update the fields that were provided in the request
	// Using pointers allows us to distinguish between "not provided" and "empty string"
	if req.Title != nil {
		if err := validateTitle(*req.Title); err != nil {
			respondWithError(w, http.StatusBadRequest, "Validation failed", err.Error())
			return
		}
		note.Title = *req.Title
	}

	if req.Content != nil {
		if err := validateContent(*req.Content); err != nil {
			respondWithError(w, http.StatusBadRequest, "Validation failed", err.Error())
			return
		}
		note.Content = *req.Content
	}

	// Update the timestamp
	note.UpdatedAt = time.Now()

	// Save the updated note
	store.updateNote(note)

	// Return the updated note
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{Data: *note})
}

// handleDeleteNote handles DELETE requests to remove a note.
// REST principle: DELETE is used for removing resources.
//
// Returns: 204 No Content on success, or 404 Not Found if note doesn't exist.
func handleDeleteNote(w http.ResponseWriter, r *http.Request) {
	// Verify the HTTP method is DELETE
	if r.Method != http.MethodDelete {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", "Only DELETE is allowed")
		return
	}

	// Extract the ID from the URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		respondWithError(w, http.StatusBadRequest, "Missing ID", "Note ID is required")
		return
	}

	// Parse the ID
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", "ID must be an integer")
		return
	}

	// Check if the note exists before attempting deletion
	_, found := store.getNoteByID(id)
	if !found {
		respondWithError(w, http.StatusNotFound, "Not found", fmt.Sprintf("Note with ID %d not found", id))
		return
	}

	// Delete the note
	store.deleteNote(id)

	// Return 204 No Content (successful deletion, no response body)
	w.WriteHeader(http.StatusNoContent)
}

// handleHealthCheck handles GET requests to /health.
// This is a simple endpoint to verify the API is running.
// Useful for load balancers and monitoring systems.
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", "Only GET is allowed")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().String(),
	})
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// respondWithError sends a standardized error response.
// This ensures all error responses follow the same format.
func respondWithError(w http.ResponseWriter, statusCode int, message, details string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   message,
		Details: details,
	})
}

// validateCreateNoteRequest validates the fields in a CreateNoteRequest.
// Input validation is crucial for API security and reliability.
func validateCreateNoteRequest(req CreateNoteRequest) error {
	if err := validateTitle(req.Title); err != nil {
		return err
	}
	if err := validateContent(req.Content); err != nil {
		return err
	}
	return nil
}

// validateTitle validates the title field.
// Rules:
// - Title cannot be empty
// - Title must not exceed 200 characters
func validateTitle(title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if len(title) > 200 {
		return fmt.Errorf("title must not exceed 200 characters (current: %d)", len(title))
	}
	return nil
}

// validateContent validates the content field.
// Rules:
// - Content cannot be empty
// - Content must not exceed 5000 characters
func validateContent(content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	if len(content) > 5000 {
		return fmt.Errorf("content must not exceed 5000 characters (current: %d)", len(content))
	}
	return nil
}

// ============================================================================
// NOTESSTORE METHODS (Thread-safe operations)
// ============================================================================

// getNextID generates the next unique ID for a note.
// Must be called with store lock held for safety.
func (s *NotesStore) getNextID() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID
	s.nextID++
	return id
}

// addNote adds a new note to the store.
// Thread-safe: uses mutex to ensure exclusive access during write.
func (s *NotesStore) addNote(note *Note) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notes[note.ID] = note
}

// getNoteByID retrieves a note by its ID.
// Thread-safe: uses RLock for concurrent reads.
func (s *NotesStore) getNoteByID(id int) (*Note, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	note, found := s.notes[id]
	return note, found
}

// getAllNotes retrieves all notes.
// Returns a copy to prevent external modification of internal state.
func (s *NotesStore) getAllNotes() []Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	notes := make([]Note, 0, len(s.notes))
	for _, note := range s.notes {
		notes = append(notes, *note)
	}
	return notes
}

// updateNote updates an existing note.
// Thread-safe: uses mutex for exclusive write access.
func (s *NotesStore) updateNote(note *Note) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notes[note.ID] = note
}

// deleteNote removes a note from the store.
// Thread-safe: uses mutex for exclusive write access.
func (s *NotesStore) deleteNote(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.notes, id)
}

// ============================================================================
// MAIN FUNCTION
// ============================================================================

func main() {
	// Initialize HTTP routes
	// The path parameter "/:id" is parsed manually in handlers.
	// In production, use a router library like gorilla/mux or chi for cleaner routing.

	// Health check endpoint
	http.HandleFunc("/health", handleHealthCheck)

	// Notes endpoints
	// Note: This simple routing uses the same handler for all CRUD operations.
	// The handler determines which operation to perform based on the HTTP method.
	http.HandleFunc("/notes", handleCreateNote)
	http.HandleFunc("/notes/", handleGetNotes)

	// Create a new router to apply middleware
	// Middleware wraps the default mux to add logging and JSON headers to all requests.
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealthCheck)
	mux.HandleFunc("/notes", handleCreateNote)
	mux.HandleFunc("/notes/", handleGetNotes)

	// Apply middleware to the mux
	// Middleware is applied in reverse order (right-to-left execution)
	handler := jsonContentTypeMiddleware(mux)
	handler = loggingMiddleware(handler)

	// Start the HTTP server
	// The server listens on port 8080
	port := ":8080"
	log.Printf("Starting Personal Notes API server on http://localhost:8080")
	log.Printf("Endpoints:")
	log.Printf("  GET    /health                  - Health check")
	log.Printf("  POST   /notes                   - Create a note")
	log.Printf("  GET    /notes                   - Get all notes")
	log.Printf("  GET    /notes/:id               - Get a specific note")
	log.Printf("  PATCH  /notes/:id               - Update a note")
	log.Printf("  DELETE /notes/:id               - Delete a note")

	// Start listening for incoming HTTP requests
	// ListenAndServe blocks until the server is shut down
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
