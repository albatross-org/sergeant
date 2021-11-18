package server

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/albatross-org/sergeant"
	"github.com/gin-gonic/gin"
)

// SetJSON is the response returned when a client asks for info about a set.
type SetJSON struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`

	Background string `json:"background"`
	Color      string `json:"color"`
}

// SetListJSON is the response returned when a Client asks for the available sets.
// It is a sort.Interface and the sorting method is alphabetical but "All" always comes first.
type SetListJSON []SetJSON

// Len is the number of elements in the SetList. This method is used to implement sort.Interface.
func (setList SetListJSON) Len() int {
	return len(setList)
}

// Less reports whether the leement with index i short short before the leement with index j.
// This method is used to implement sort.Interface.
func (setList SetListJSON) Less(i, j int) bool {
	iName := (setList)[i].Name
	jName := (setList)[j].Name

	if iName == "All" {
		return true
	}

	return iName[0] < jName[0]
}

// Swap swaps the elements with indexes i and j.
// This method is used to implement sort.Interface.
func (setList SetListJSON) Swap(i, j int) {
	(setList)[i], (setList)[j] = (setList)[j], (setList)[i]
}

// setToJSON returns the SetJSON representation of a sergeant.ConfigSet.
func setToJSON(name string, set sergeant.ConfigSet) SetJSON {
	if set.Name == "" {
		set.Name = strings.ReplaceAll(set.Name, "-", " ")
		set.Name = strings.Title(set.Name)
	}

	if set.Description == "" {
		set.Description = fmt.Sprintf("This is a custom set containing all %s cards.", set.Name)
	}

	var color colorPair

	if set.Background != "" {
		color.background = set.Background
	}

	if set.Color != "" {
		color.color = set.Color
	}

	if set.Background == "" && set.Color == "" {
		color = gradientColorPair(name)
	}

	return SetJSON{
		Name:        name,
		DisplayName: set.Name,
		Description: set.Description,
		Background:  color.background,
		Color:       color.color,
	}
}

// getSetListJSON returns the SetListJSON for all sets in the store.
func getSetListJSON() SetListJSON {
	response := SetListJSON{}

	for name, set := range store.Sets {
		if !set.Hidden {
			response = append(response, setToJSON(name, set))
		}
	}

	sort.Sort(response)

	return response
}

// SetRequest is a request for a certain type of set. This can include things like the Set
// name And paths that it should contain.
type SetRequest struct {
	Name string

	Paths          []string
	Tags           []string
	BeforeDuration time.Duration
	AfterDuration  time.Duration
	BeforeDate     time.Time
	AfterDate      time.Time

	Color      string
	Background string
}

// setConfigFromRequest returns a new ConfigSet struct from a *gin.C by looking at query parameteres.
// It will return an error if the set request isn't valid.
func setConfigFromRequest(c *gin.Context) (sergeant.ConfigSet, error) {
	config := sergeant.ConfigSet{}

	var exists bool
	var err error

	name, exists := c.GetQuery("setName")
	if !exists {
		name = "all"
	}

	existingConfig, exists := store.Sets[name]
	if !exists {
		return sergeant.ConfigSet{}, fmt.Errorf("the existing set %q doesn't exist", name)
	}

	config = existingConfig

	background, exists := c.GetQuery("setBackground")
	if exists {
		config.Background = background
	} else if !exists && config.Background == "" {
		config.Background = "linear-gradient(315deg, #11998e, #38ef7d)"
	}

	color, exists := c.GetQuery("setColor")
	if exists {
		config.Color = color
	} else if !exists && config.Color == "" {
		config.Color = "#11998e"
	}

	pathsOr, exists := c.GetQueryArray("setPathsOr")
	if exists {
		config.PathsOr = append(config.PathsOr, pathsOr...)
	}

	pathsAnd, exists := c.GetQueryArray("setPathsAnd")
	if exists {
		config.PathsAnd = append(config.PathsOr, pathsAnd...)
	}

	tagsOr, exists := c.GetQueryArray("setTagsOr")
	if exists {
		config.TagsOr = append(config.TagsOr, tagsOr...)
	}

	tagsAnd, exists := c.GetQueryArray("setTagsAnd")
	if exists {
		config.TagsAnd = append(config.TagsOr, tagsAnd...)
	}

	rawBeforeDuration, exists := c.GetQuery("setBeforeDuration")
	if exists {
		config.BeforeDuration, err = time.ParseDuration(rawBeforeDuration)
		if err != nil {
			return sergeant.ConfigSet{}, fmt.Errorf("invalid before duration %q specified: %w", rawBeforeDuration, err)
		}
	}

	rawAfterDuration, exists := c.GetQuery("setAfterDuration")
	if exists {
		config.AfterDuration, err = time.ParseDuration(rawAfterDuration)
		if err != nil {
			return sergeant.ConfigSet{}, fmt.Errorf("invalid after duration %q specified: %w", rawAfterDuration, err)
		}
	}

	rawBeforeDate, exists := c.GetQuery("setBeforeDate")
	if exists {
		config.BeforeDate, err = time.Parse("2006-01-02 15:04", rawBeforeDate)
		if err != nil {
			return sergeant.ConfigSet{}, fmt.Errorf("invalid before date %q specified: %w", rawBeforeDate, err)
		}
	}

	rawAfterDate, exists := c.GetQuery("setAfterDate")
	if exists {
		config.AfterDate, err = time.Parse("2006-01-02 15:04", rawAfterDate)
		if err != nil {
			return sergeant.ConfigSet{}, fmt.Errorf("invalid after date %q specified: %w", rawAfterDate, err)
		}
	}

	return config, nil
}

// toConfig will turn a SetRequest to a set

// SetHeatmapJSON is the response sent to the client when it's asked for heatmap data.
type SetHeatmapJSON struct {
	Day     string `json:"day"`
	Value   int    `json:"value"`
	Perfect int    `json:"perfect"`
	Minor   int    `json:"minor"`
	Major   int    `json:"major"`
}

// getSetHeatmapJSON returns the SetHeatmapJSON for a specific set and user. If user is blank, it uses all completions.
func getSetHeatmapJSON(set *sergeant.Set, user string) []SetHeatmapJSON {
	days := map[string]*SetHeatmapJSON{}

	for _, card := range set.Cards {
		for _, completion := range card.CompletionsPerfect {
			if user == "" || user == completion.User {
				date := completion.Date.Format("2006-01-02")

				if days[date] == nil {
					days[date] = &SetHeatmapJSON{}
				}

				days[date].Perfect++
				days[date].Day = date
				days[date].Value += int(completion.Duration) / (1000 * 1000 * 1000)
			}
		}

		for _, completion := range card.CompletionsMinor {
			if user == "" || user == completion.User {
				date := completion.Date.Format("2006-01-02")

				if days[date] == nil {
					days[date] = &SetHeatmapJSON{}
				}

				days[date].Minor++
				days[date].Day = date
				days[date].Value += int(completion.Duration) / (1000 * 1000 * 1000)
			}
		}

		for _, completion := range card.CompletionsMajor {
			if user == "" || user == completion.User {
				date := completion.Date.Format("2006-01-02")

				if days[date] == nil {
					days[date] = &SetHeatmapJSON{}
				}

				days[date].Major++
				days[date].Day = date
				days[date].Value += int(completion.Duration) / (1000 * 1000 * 1000)
			}
		}

	}

	list := []SetHeatmapJSON{}

	for _, data := range days {
		list = append(list, *data)
	}

	return list
}
