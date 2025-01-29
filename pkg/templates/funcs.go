package templates

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/sboon-gg/svctl/pkg/maplist"
)

func (r *Renderer) FuncMap() template.FuncMap {
	f := sprig.TxtFuncMap()

	extra := template.FuncMap{
		"pyBool":  pyBool,
		"quote":   quote,
		"env":     env,
		"maplist": r.maplist,
	}

	for k, v := range extra {
		f[k] = v
	}

	return f
}

func (t *Renderer) maplist(filterMap interface{}, rawMaplist string) (string, error) {
	fmt.Printf("filterMap: %T\n", filterMap)

	var filter maplist.MapInfo
	if f, ok := filterMap.(string); ok {
		filter = maplist.Parse(fmt.Sprintf("%s %s", maplist.MaplistAppendStr, f))[0]
	} else {
		c, err := json.Marshal(filterMap)
		if err != nil {
			return "", errors.Join(errors.New("failed to marshal filter map"), err)
		}

		if err := json.Unmarshal(c, &filter); err != nil {
			return "", errors.Join(errors.New("failed to unmarshal filter"), err)
		}
	}

	return maplist.Compose(
		maplist.Filter(maplist.Parse(rawMaplist), filter),
	), nil
}

func pyBool(b bool) (string, error) {
	if b {
		return "True", nil
	}

	return "False", nil
}

func quote(v any) string {
	return fmt.Sprintf("%q", v)
}

func env(s string) (string, error) {
	return os.Getenv(s), nil
}
