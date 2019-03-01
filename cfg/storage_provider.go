package cfg

import (
	"bufio"
	"context"
	"strings"

	"cloud.google.com/go/storage"
)

// StorageProvider describes GCP Storage based loader which loads the configuration
// from a bucket and file listed.
type StorageProvider struct {
	Bucketname string
	Filename   string
}

// Provide implements the Provider interface.
func (sp StorageProvider) Provide() (map[string]string, error) {
	ctx := context.Background()

	var config = make(map[string]string)

	// Creating Storage client
	// The client will use your default application credentials.
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	// Reading object from storage
	file, err := storageClient.Bucket(sp.Bucketname).Object(sp.Filename).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) < 3 {
			// the line doesn't have enough data
			continue
		}

		if line[0] == '#' {
			// the line starts with a comment character
			continue
		}

		// find the first equals sign
		index := strings.Index(line, "=")

		// if we couldn't find one
		if index <= 0 {
			// the line is invalid
			continue
		}

		if index == len(line)-1 {
			// the line is invalid
			continue
		}

		// add the item to the config
		config[line[:index]] = line[index+1:]
	}

	return config, nil
}
