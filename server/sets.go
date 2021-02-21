package server

import (
	"fmt"
	"sort"
	"strings"

	"github.com/albatross-org/sergeant"
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
		response = append(response, setToJSON(name, set))
	}

	sort.Sort(response)

	return response
}

// SetHeatmapJSON is the response sent to the client when it's asked for heatmap data.
type SetHeatmapJSON struct {
	Day   string `json:"day"`
	Value int    `json:"value"`
}

// getSetHeatmapJSON returns the SetHeatmapJSON for a specific set.
func getSetHeatmapJSON(set *sergeant.Set) []SetHeatmapJSON {
	days := map[string]int{}

	for _, card := range set.Cards {
		for _, completion := range card.CompletionsPerfect {
			days[completion.Date.Format("2006-01-02")] += int(completion.Duration)
		}
		for _, completion := range card.CompletionsMinor {
			days[completion.Date.Format("2006-01-02")] += int(completion.Duration)
		}
		for _, completion := range card.CompletionsMajor {
			days[completion.Date.Format("2006-01-02")] += int(completion.Duration)
		}

	}

	list := []SetHeatmapJSON{}

	for day, value := range days {
		list = append(list, SetHeatmapJSON{
			Day:   day,
			Value: value / (1000 * 1000 * 1000), // Convert nanoseconds to seconds
		})
	}

	return list
}
