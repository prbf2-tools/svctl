package templates

import (
	"os"
	"testing"

	"github.com/sboon-gg/svctl/pkg/maplist"
	"github.com/stretchr/testify/assert"
)

func TestTemplatesRender(t *testing.T) {
	tests := []struct {
		name     string
		values   Values
		expected []RenderOutput
		env      map[string]string
	}{
		{
			name:   "defaults",
			values: Values{},
			expected: []RenderOutput{
				{
					Template: Template{
						Source:      "file.yaml.tpl",
						Destination: "file.yaml",
					},
					Content: []byte("pyBool: True\nanother: test-string\nquoted: \"should-be-quoted\"\nnegativeBool: false\nenvVal: \nzeroValue: 1\n"),
				},
				{
					Template: Template{
						Source:      "maplist.con.tpl",
						Destination: "maplist.con",
					},
					Content: []byte(`mapList.append kashan_desert gpm_cq 16
mapList.append kashan_desert gpm_cq 32
mapList.append kashan_desert gpm_cq 64
mapList.append kashan_desert gpm_cq 128
mapList.append sahel gpm_cq 64
mapList.append sahel gpm_insurgency 64
mapList.append sahel gpm_skirmish 64
mapList.append sahel gpm_coop 64
`),
				},
			},
		},
		{
			name: "overwrite",
			values: Values{
				"boolTest":              false,
				"test":                  "changed-string",
				"quoted":                "but different",
				"overriddenByZeroValue": 0,
				"maps": []map[string]interface{}{
					{
						"name": "saaremaa",
					},
				},
			},
			env: map[string]string{
				"ENV_VAL": "env-val",
			},
			expected: []RenderOutput{
				{
					Template: Template{
						Source:      "file.yaml.tpl",
						Destination: "file.yaml",
					},
					Content: []byte("pyBool: False\nanother: changed-string\nquoted: \"but different\"\nnegativeBool: false\nenvVal: env-val\nzeroValue: 0\n"),
				},
				{
					Template: Template{
						Source:      "maplist.con.tpl",
						Destination: "maplist.con",
					},
					Content: []byte(`mapList.append saaremaa gpm_cq 16
mapList.append saaremaa gpm_cq 32
mapList.append saaremaa gpm_cq 64
mapList.append saaremaa gpm_cq 128
mapList.append saaremaa gpm_skirmish 16
mapList.append saaremaa gpm_cnc 16
mapList.append saaremaa gpm_cnc 32
mapList.append saaremaa gpm_cnc 64
mapList.append saaremaa gpm_cnc 128
mapList.append saaremaa gpm_coop 32
mapList.append saaremaa gpm_coop 64
`),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.env {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			tmpl, err := NewFromPath("./testdata/example", WithMaps(maplist.DefaultMapList))
			assert.NoError(t, err)

			out, err := tmpl.Render(test.values)
			assert.NoError(t, err)

			assert.Len(t, out, len(test.expected))

			for i, v := range test.expected {
				println(string(out[i].Content))
				assert.Equal(t, v, out[i])
			}
		})
	}
}
