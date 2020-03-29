package gingle

import (
	"fmt"
	"strings"
)

// node of trie tree
type node struct {
	pattern  string  // consistent matching pattern (only a few nodes have), such as /cn/content, etc
	part     string  // current presenting part (all nodes have), such as /cn, /content, etc
	children []*node // children of the node
	isFuzzy  bool    // whether support fuzzy match
}

// String defines the format for printing node
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isFuzzy=%t}", n.pattern, n.part, n.isFuzzy)
}

// traverse collects all nodes in level order
func (n *node) traverse(nodes *([]*node)) {
	if n.pattern != "" {
		*nodes = append(*nodes, n)
	}
	for _, child := range n.children {
		child.traverse(nodes)
	}
}

// matchChild is for inserting nodes
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		// Static match or dynamic match
		if child.part == part || child.isFuzzy {
			return child
		}
	}
	return nil
}

// matchChildren is for searching nodes
func (n *node) matchChildren(part string) []*node {
	children := make([]*node, 0)
	for _, child := range n.children {
		// Static match or dynamic match
		if child.part == part || child.isFuzzy {
			children = append(children, child)
		}
	}
	return children
}

// insert parses the pattern to insert nodes to trie tree
func (n *node) insert(pattern string, parts []string, height int) {
	// When height is equal to len(parts), it comes to the deepest of trie tree
	// It is also time to set pattern which relating to a specific handler
	if height == len(parts) {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)

	// If node does not exist (the same level), then new and append one
	if child == nil {
		child = &node{
			part:    part,
			isFuzzy: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}

	// If node do exist (the same level), continue to search
	child.insert(pattern, parts, height+1)
}

// search parses the pattern to search the node from trie tree
func (n *node) search(parts []string, height int) *node {
	// When it comes to the deepest of trie tree, node matches exactly
	// When it encounters `*` symbol, node matches fuzzily
	if height == len(parts) || strings.HasPrefix(n.part, "*") {
		// Empty pattern indicates no such traversal in trie tree
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result // If succeed, return result
		}
	}

	return nil // If failed, return nil
}
