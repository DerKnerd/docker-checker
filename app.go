package main

import (
	"docker-checker/configuration"
	"docker-checker/dockerApi"
	"docker-checker/mailing"
	"flag"
	"fmt"
	"github.com/hashicorp/go-version"
	"log"
	"os"
	"runtime"
	"sort"
	"sync"
)

func checkForUpdate(wg *sync.WaitGroup, channel chan configuration.Image, cpu int, email *configuration.EmailConfig) {
	for image := range channel {
		log.Printf("CPU %d: Check for docker image %s", cpu, image.Name)
		tagList, err := dockerApi.GetVersions(image.Name, cpu)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("CPU %d: Found %d tags for image %s", cpu, len(tagList.Tags), tagList.Name)
		versionConstraint, err := version.NewConstraint(image.Constraint)
		if err != nil {
			log.Printf("CPU %d: Failed to create version constraint for version %s", cpu, image.Constraint)
			log.Printf("CPU %d: %s", cpu, err.Error())
			log.Printf("CPU %d: Use used version as constraint version", cpu)
			versionConstraint, _ = version.NewConstraint("> " + image.UsedVersion)
		}

		versions := make([]*version.Version, 0)
		for _, raw := range tagList.Tags {
			v, _ := version.NewVersion(raw)
			if v != nil {
				versions = append(versions, v)
			}
		}

		usedVersion, _ := version.NewVersion(image.UsedVersion)

		sort.Sort(sort.Reverse(version.Collection(versions)))
		log.Printf("CPU %d: Latest version for %s is %s", cpu, tagList.Name, versions[0].String())

		for _, tag := range versions {
			if versionConstraint.Check(tag) && !tag.LessThanOrEqual(usedVersion) {
				log.Printf("CPU %d: Found newer version for image %s:%s, newer version is %s", cpu, tagList.Name, usedVersion.String(), tag.String())
				if err = mailing.SendMail(usedVersion, tag, &image, email); err != nil {
					log.Printf("CPU %d: Failed to send message for image %s", cpu, image.Name)
					log.Printf("CPU %d: %s", cpu, err.Error())
				}
				break
			}
		}
	}
	wg.Done()
}

func main() {
	log.Println("MAIN:  Read config file")
	pwd, _ := os.Getwd()
	configFilePath := ""
	flag.StringVar(&configFilePath, "config-file", fmt.Sprintf("%s/docker-checker.yaml", pwd), "Specifies the config file to use")
	flag.Parse()

	config, err := configuration.ParseConfig(configFilePath)
	if err != nil {
		log.Fatal("MAIN:  Failed to read config")
	}

	wg := &sync.WaitGroup{}
	wg.Add(runtime.NumCPU())
	imageChan := make(chan configuration.Image, len(config.Images))

	for i := 0; i < runtime.NumCPU(); i++ {
		go checkForUpdate(wg, imageChan, i, &config.Email)
	}
	log.Println("MAIN:  Start check for docker images")
	for _, image := range config.Images {
		imageChan <- image
	}

	close(imageChan)
	wg.Wait()
}
