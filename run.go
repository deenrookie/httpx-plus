package httpx_plus

import (
	"github.com/deenrookie/httpx-plus/runner"
	"github.com/projectdiscovery/gologger"
	"os"
	"os/signal"
)

func HttpDetectStart(target string) (rets []runner.Result) {
	// Parse the command line flags and read config files
	options := runner.ParseOptions()

	options.StatusCode = true
	options.ExtractTitle = true
	options.TechDetect = true
	options.OutputServerHeader = true
	options.Threads = 5
	options.OutputCName = true

	httpxRunner, err := runner.New(options)
	if err != nil {
		gologger.Fatal().Msgf("Could not create runner: %s\n", err)
	}

	// Setup graceful exits
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			gologger.Info().Msgf("CTRL+C pressed: Exiting\n")
			httpxRunner.Close()
			if options.ShouldSaveResume() {
				gologger.Info().Msgf("Creating resume file: %s\n", runner.DefaultResumeFile)
				err := httpxRunner.SaveResumeConfig()
				if err != nil {
					gologger.Error().Msgf("Couldn't create resume file: %s\n", err)
				}
			}
			os.Exit(1)
		}
	}()

	rets = httpxRunner.RunEnumeration(target)
	httpxRunner.Close()
	return
}
//
//func main() {
//	fmt.Println(HttpDetectStart("d33n.cn"))
//}
