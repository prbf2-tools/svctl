package maplist

import (
	"bufio"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

const (
	MaplistAppendStr = "mapList.append"
)

type MapInfo struct {
	Name  string `json:"name"`
	Mode  string `json:"mode"`
	Layer int    `json:"layer"`
}

func Parse(maplistContent string) []MapInfo {
	reader := strings.NewReader(maplistContent)
	scanner := bufio.NewScanner(reader)

	var maps []MapInfo

	for scanner.Scan() {
		text := scanner.Text()

		if !strings.HasPrefix(text, MaplistAppendStr) {
			continue
		}

		parts := strings.Split(text, " ")

		layer, _ := strconv.Atoi(parts[3])

		maps = append(maps, MapInfo{
			Name:  parts[1],
			Mode:  parts[2],
			Layer: layer,
		})
	}

	return maps
}

func Compose(maps []MapInfo) string {
	builder := strings.Builder{}

	for _, m := range maps {
		builder.WriteString(fmt.Sprintf("mapList.append %s %s %d\n", m.Name, m.Mode, m.Layer))
	}

	return builder.String()
}

func Filter(allMaps []MapInfo, filter MapInfo) []MapInfo {
	var maps []MapInfo

	for _, m := range allMaps {
		if strings.Contains(m.Name, filter.Name) && strings.Contains(m.Mode, filter.Mode) {
			if filter.Layer == 0 || m.Layer == filter.Layer {
				maps = append(maps, m)
			}
		}
	}

	return maps
}
