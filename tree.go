// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package i3ipc

import (
	"encoding/json"
)

// I3Node represents a Node in the i3 tree. For documentation of the fields,
// refer to http://i3wm.org/docs/ipc.html#_tree_reply.
type I3Node struct {
	ID                 int32
	Name               string
	Type               string
	Border             string
	CurrentBorderWidth int32 `json:"current_border_width"`
	Layout             string
	Orientation        string
	Percent            float64
	Rect               Rect
	WindowRect         Rect
	DecoRect           Rect `json:"deco_rect"`
	Geometry           Rect
	Window             int32
	Urgent             bool
	Focused            bool
	Nodes              []I3Node
	FloatingNodes      []I3Node `json:"floating_nodes"`
	Parent             *I3Node

	// Properties, not listed in docs:
	Sticky            bool
	Floating          string
	Last_Split_Layout string `json:"last_split_layout"`
	Fullscreen_Mode   int32  `json:"fullscreen_mode"`
	Scratchpad_State  string `json:"scratchpad_state"`
	Workspace_Layout  string `json:"workspace_layout"`

	WindowProperties struct {
		Title    string
		Instance string
		Class    string
	} `json:"window_properties"`
}

// Traverses the tree setting correct reference to a parent node.
func setParent(node, parent *I3Node) {

	node.Parent = parent

	for i := range node.Nodes {
		setParent(&node.Nodes[i], node)
	}
	for i := range node.FloatingNodes {
		setParent(&node.FloatingNodes[i], node)
	}
}

// GetTree fetches the layout tree.
func (socket *IPCSocket) GetTree() (root I3Node, err error) {
	jsonReply, err := socket.Raw(I3GetTree, "")
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonReply, &root)
	defer setParent(&root, nil)

	if err == nil {
		return
	}
	// For an explanation of this error silencing, see GetOutputs().
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		err = nil
	}
	return
}
