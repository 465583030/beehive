/*
 *    Copyright (C) 2014 Christian Muehlhaeuser
 *
 *    This program is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU Affero General Public License as published
 *    by the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *
 *    This program is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU Affero General Public License for more details.
 *
 *    You should have received a copy of the GNU Affero General Public License
 *    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *    Authors:
 *      Christian Muehlhaeuser <muesli@gmail.com>
 *      Johannes Fürmann <johannes@weltraumpflege.org>
 */

// beehive's SpaceAPI module.
package spaceapibee

import (
	"encoding/json"
	"github.com/muesli/beehive/bees"
	"io/ioutil"
	"net/http"
)

type SpaceApiBee struct {
	bees.Bee

	url string

	evchan chan bees.Event
}

func (mod *SpaceApiBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "get_status":
		type SpaceApiResult struct {
			Status string `json: "status"`
			Open   bool   `json: "open"`
		}
		api_state := new(SpaceApiResult)

		// get json data
		resp, err := http.Get(mod.url)
		var text string
		if err != nil {
			text = "Error: SpaceAPI instance @ " + mod.url + " not reachable"
		} else {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			err = json.Unmarshal(body, api_state)

			if err != nil {
				text = "Sorry, couldn't unmarshal the JSON data from SpaceAPI Instance @ " + mod.url
			} else {
				text = api_state.Status
			}
		}

		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "query_result",
			Options: []bees.Placeholder{
				bees.Placeholder{
					Name:  "status",
					Type:  "bool",
					Value: api_state.Open,
				},
				bees.Placeholder{
					Name:  "text",
					Type:  "string",
					Value: text,
				},
			},
		}
		mod.evchan <- ev

	default:
		panic("Unknown action triggered in " +mod.Name()+": "+action.Name)
	}

	return outs
}

func (mod *SpaceApiBee) Run(eventChan chan bees.Event) {
	mod.evchan = eventChan
}
