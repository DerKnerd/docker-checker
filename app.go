package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/go-version"
	"log"
	"os"
	"sort"
)

func main() {
	log.Println("Read config file")
	pwd, _ := os.Getwd()
	configFilePath := ""
	flag.StringVar(&configFilePath, "config-file", fmt.Sprintf("%s/docker-checker.yaml", pwd), "Specifies the config file to use")
	flag.Parse()

	configuration, err := ParseConfig(configFilePath)
	if err != nil {
		log.Fatal("Failed to read configuration")
	}

	log.Println("Start check for docker images")
	for _, image := range configuration.Images {
		log.Printf("Check for docker image %s", image.Name)
		tagList, err := GetVersions(image.Name)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("Found %d tagList for image %s", len(tagList.Tags), tagList.Name)
		versionConstraint, err := version.NewConstraint(image.Constraint)
		if err != nil {
			log.Printf("Failed to create version constraint for version %s", image.Constraint)
			log.Println(err)
			log.Printf("Use used version as constraint version")
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
		log.Printf("Latest version for %s is %s", tagList.Name, versions[0].String())

		for _, tag := range versions {
			if versionConstraint.Check(tag) && !tag.LessThanOrEqual(usedVersion) {
				log.Printf("Found newer version for image %s, newer version is %s", tagList.Name, tag.String())
				if err = SendMail(*usedVersion, *tag, image, configuration.Email); err != nil {
					log.Printf("Failed to send message for image %s", image.Name)
					log.Println(err.Error())
				}
				break
			}
		}
	}
}
