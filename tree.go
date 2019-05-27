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
	//int32 isn't large enough to hold all the ids
	ID                 int64
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
	Window             int64
	Urgent             bool
	Focused            bool
	Floating_Nodes     []I3Node
	Nodes              []I3Node
	Parent             *I3Node

	// Properties, not listed in docs:
	WindowProps struct {
		// Transient_for ?
		Title    string
		Instance string
		Class    string
	} `json:"window_properties"`
	// Swallows []I3Node ?
	Sticky            bool
	Floating          string
	Last_Split_Layout string
	Focus             []int64
	FocusOrder        []*I3Node
	Fullscreen_Mode   int32
	Scratchpad_State  string
	Workspace_Layout  string
}

// Traverses the tree setting correct reference to a parent node.
func setParent(node, parent *I3Node) {

	node.Parent = parent

	for i := range node.Nodes {
		setParent(&node.Nodes[i], node)
	}
	for i := range node.Floating_Nodes {
		setParent(&node.Floating_Nodes[i], node)
	}
}

func setFocusList(node *I3Node) {
	node.FocusOrder = make([]*I3Node, len(node.Nodes))

	for i := range node.Focus {
		for n := range node.Nodes {
			curr := node.Nodes[n]
			if curr.ID == node.Focus[i] {
				node.FocusOrder[i] = &curr
			}
		}
	}

	for i := range node.Nodes {
		setFocusList(&node.Nodes[i])
	}
}

func parseTree(root *I3Node) {
	setParent(root, nil)
	setFocusList(root)
}

// GetTree fetches the layout tree.
func (socket *IPCSocket) GetTree() (root I3Node, err error) {
	jsonReply, err := socket.Raw(I3GetTree, "")
	if err != nil {
		return
	}
	// fmt.Printf("%s\n", string(jsonReply))
	// panic("")

	defer parseTree(&root)

	err = json.Unmarshal(jsonReply, &root)
	if err == nil {
		return
	}
	// For an explanation of this error silencing, see GetOutputs().
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		err = nil
	}
	return
}
