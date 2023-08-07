// Package http provides functionalities to start an HTTP server used to save POST data to files.
package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"zero-provers/server/logger"

	"github.com/rs/zerolog"
)

// saveEndpoint is the URL path for the save endpoint.
const saveEndpoint = "/save"

var (
	// log is the package-level variable used for logging messages and errors.
	log zerolog.Logger

	// proofsDir is the repository in which proofs are saved to the disk.
	proofsDir string

	// count the number of proofs saved on disk.
	proofCount = 1
)

// StartHTTPServer starts an HTTP server on the specified port and sets up the necessary endpoints.
// The server listens for incoming requests and handles them accordingly.
// The `/save` endpoint allows clients to save data to a file in the specified output directory.
func StartHTTPServer(port int, outputDir string) error {
	// Set up the logger.
	lc := logger.LoggerConfig{
		Level:       zerolog.InfoLevel,
		CallerField: "http",
	}
	log = logger.NewLogger(lc)

	// Create proof directory.
	err := os.Mkdir(outputDir, 0755)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create the proofs directory")
		return err
	}
	proofsDir = outputDir

	// Start the HTTP server.
	log.Info().Msgf("HTTP server is starting on port %d", port)
	http.HandleFunc(saveEndpoint, saveHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Error().Err(err).Msg("Unable to start the HTTP server")
		return err
	}
	return nil
}

// saveHandler is the handler function for the `/save` endpoint.
// It processes incoming POST requests containing JSON data and displays the content of the proof.
func saveHandler(w http.ResponseWriter, r *http.Request) {
	// Only handle POST requests.
	if r.Method == http.MethodPost {
		log.Info().Msgf("POST request received on %s endpoint", saveEndpoint)

		// Decode the incoming JSON data into an interface{}.
		var data interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			log.Error().Err(err).Msg("Unable to decode JSON")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Marshal the data to JSON format with indentation.
		indentedData, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			log.Error().Err(err).Msg("Unable to marshal data to JSON format")
			return
		}

		// Save proof to disk.
		proofPath := fmt.Sprintf("%s/%d.json", proofsDir, proofCount)
		if err = os.WriteFile(proofPath, indentedData, 0644); err != nil {
			log.Error().Err(err).Msg("Unable to write to file")
			return
		}
		proofCount++
		log.Info().Msg("Proof saved to disk")
	} else {
		log.Info().Msgf("Invalid request method on %s endpoint", saveEndpoint)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
