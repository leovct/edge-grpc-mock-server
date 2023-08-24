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

var (
	// log is the package-level variable used for logging messages and errors.
	log zerolog.Logger

	// saveEndpoint is the URL path for the save endpoint.
	saveEndpoint string

	// proofsDir is the repository in which proofs are saved to the disk.
	proofsDir string

	// count the number of proofs saved on disk.
	proofCount = 1
)

type ServerConfig struct {
	LogLevel        zerolog.Level
	Port            int
	SaveEndpoint    string
	ProofsOutputDir string
}

// StartHTTPServer starts an HTTP server on the specified port and sets up the necessary endpoints.
// The server listens for incoming requests and handles them accordingly.
// The `/save` endpoint allows clients to save data to a file in the specified output directory.
func StartHTTPServer(config ServerConfig) error {
	// Set up the logger.
	lc := logger.LoggerConfig{
		Level:       config.LogLevel,
		CallerField: "http",
	}
	log = logger.NewLogger(lc)

	// Create the proofs directory if it doesn't exist.
	if _, err := os.Stat(config.ProofsOutputDir); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(config.ProofsOutputDir, 0755); err != nil {
				log.Error().Err(err).Msg("Unable to create the proofs directory")
				return err
			}
		} else {
			log.Error().Err(err).Msg("Unable to check if the proofs directory exists")
			return err
		}
	}
	proofsDir = config.ProofsOutputDir

	// Start the HTTP server.
	saveEndpoint = config.SaveEndpoint
	http.HandleFunc(saveEndpoint, saveHandler)
	log.Info().Msgf("HTTP server save endpoint: %s ready", saveEndpoint)
	log.Info().Msgf("HTTP server is starting on port %d", config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil); err != nil {
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
