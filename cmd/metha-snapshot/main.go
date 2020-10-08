// Download all known endpoints, generate a single JSON file.
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"

	"github.com/miku/metha"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var (
	filename    = flag.String("f", "", fmt.Sprintf("filename with endpoints, defaults to list of %d sites", len(metha.Endpoints)))
	baseDir     = flag.String("base-dir", metha.GetBaseDir(), "base dir for harvested files")
	format      = flag.String("format", "oai_dc", "metadata format")
	maxRequests = flag.Int("max", 1048576, "maximum number of token loops")
	quiet       = flag.Bool("q", false, "suppress all output")
	numWorkers  = flag.Int("w", 64, "workers")
	shuffle     = flag.Bool("S", false, "shuffle hosts")
)

func main() {
	flag.Parse()
	var (
		endpoints = metha.Endpoints
		failed    []string
	)
	if *filename != "" {
		b, err := ioutil.ReadFile(*filename)
		if err != nil {
			log.Fatal(err)
		}
		endpoints = strings.Split(string(b), "\n")
	}
	if *shuffle {
		rand.Shuffle(len(endpoints), func(i, j int) {
			endpoints[i], endpoints[j] = endpoints[j], endpoints[i]
		})
	}
	if *quiet {
		log.SetOutput(ioutil.Discard)
	}
	g := new(errgroup.Group)
	urlC := make(chan string) // produce URL
	g.Go(func() error {
		defer close(urlC)
		for _, endpoint := range endpoints {
			urlC <- endpoint
		}
		return nil
	})
	for i := 0; i < *numWorkers; i++ {
		name := fmt.Sprintf("%04d", i)
		// Each worker starts a go routine that fetches urls off the urls
		// channel and runs the harvest.
		g.Go(func() error {
			var j int
			for u := range urlC {
				j++
				log.Printf("[%s@%d] %s", name, j, u)
				harvest, err := metha.NewHarvest(u)
				if err != nil {
					failed = append(failed, u)
					log.Printf("failed (init): %s, %v", u, err)
					continue
				}
				harvest.MaxRequests = *maxRequests
				harvest.CleanBeforeDecode = true
				harvest.Format = *format
				if err := harvest.Run(); err != nil {
					switch err {
					case metha.ErrAlreadySynced:
					default:
						harvest.DisableSelectiveHarvesting = true
						if err := harvest.Run(); err != nil {
							failed = append(failed, u)
							log.Printf("failed (harvest): %s, %v", u, err)
							continue
						}
					}
				}
			}
			return nil
		})
	}
	g.Wait()

	// Iterate over all successful endpoints and output JSON.

	// for i, u := range endpoints {
	// 	log.Printf("%d/%d", i, len(endpoints))
	// 	harvest, err := metha.NewHarvest(u)
	// 	if err != nil {
	// 		failed = append(failed, u)
	// 		log.Printf("failed (init): %s", u)
	// 		continue
	// 	}
	// 	harvest.MaxRequests = *maxRequests
	// 	harvest.CleanBeforeDecode = true
	// 	harvest.Format = *format
	// 	if err := harvest.Run(); err != nil {
	// 		switch err {
	// 		case metha.ErrAlreadySynced:
	// 			log.Println("this repository is up-to-date")
	// 		default:
	// 			harvest.DisableSelectiveHarvesting = true
	// 			if err := harvest.Run(); err != nil {
	// 				failed = append(failed, u)
	// 				log.Printf("failed (harvest): %s", u)
	// 				continue
	// 			}
	// 		}
	// 	}
	// }
	f, err := ioutil.TempFile("", "metha-snapshot-")
	if err != nil {
		for _, f := range failed {
			fmt.Println(f)
		}
		os.Exit(1)
	}
	defer f.Close()
	for _, u := range failed {
		if _, err := io.WriteString(f, u); err != nil {
			log.Println(err)
		}
	}
	log.Printf("wrote failed endpoints to %s", f.Name())
}