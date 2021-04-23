package main

import (
  "fmt"
  "sort"
  "net/http"
  "log"
  "encoding/json"
  "time"
)


// GithubRelease is the response from the Github API
type GithubRelease struct {
  CreatedAt string `json:"created_at"`
  Name string `json:"name"`
}
// Release is the modified response from the Github API/GithubRelease
type Release struct {
  CreatedAt time.Time
  Name string
}

func githubRepoCall() *http.Response {
  githubAPI := "https://api.github.com"
  req, err := http.NewRequest("GET", githubAPI + "/repos/Lilja/smap/releases", nil)
  req.Header.Add("Accept", "application/vnd.github.v3+json")

  if err != nil {
    log.Fatal("Error creating request object", err)
  }
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    log.Fatal("Error sending request to github", err)
  }
  log.Println("Github resp", resp, req.URL)
  return resp
}

func getReleases(githubReleases []GithubRelease) []Release {
  var releases []Release

  for _, githubRelease := range githubReleases {
    time, err := time.Parse("2006-01-02T15:04:05Z", githubRelease.CreatedAt)
    if err != nil {
      log.Fatal("Cannot decode time from string", err, githubRelease.CreatedAt)
      panic("Github changed time formating on API")
    }
    release := Release {
      CreatedAt: time,
      Name: githubRelease.Name,
    }
    releases = append(releases, release)
  }

  return releases
}

// GetLatestRepoVersion returns the latest release via Github's API
func GetLatestRepoVersion() Release {
  call := githubRepoCall()
  var githubReleases []GithubRelease
  err := json.NewDecoder(call.Body).Decode(&githubReleases)
  if err != nil {
    log.Fatal("Error creating json decoder", err)
  }
  log.Println("Github releases", githubReleases)
  releases := getReleases(githubReleases)
  sort.Slice(
    releases, func(i int, j int) bool {
      return releases[i].CreatedAt.Before(releases[j].CreatedAt)
    },
  )
  log.Println("Releases", releases)
  // Go is stupid
  var release Release
  if len(releases) >= 1 {
    release = releases[0]
  } else {
    log.Fatal("Could not find version from github")
  }
  return release
}

// CheckForUpdates doesn't do much at the moment
func CheckForUpdates(config Config) {
  if config.version != "dev" {
    release := GetLatestRepoVersion()
    fmt.Println(release)
    fmt.Println(config.builtBy)
    fmt.Println(config.builtDate)
    fmt.Println(config.commit)
    fmt.Println(config.version)
  } else {
    fmt.Println("Running dev/nightly. Can only update from release to release")
  }
}
