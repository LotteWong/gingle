package gingle

import (
	"fmt"
	"strings"
)

type node struct {
	pattern string // TODO: 全局
	part string // TODO: 局部
	children []*node // TODO: 局部
	isFuzzy bool // TODO: 局部
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isFuzzy=%t}", n.pattern, n.part, n.isFuzzy)
}

func (n *node) traverse(nodes *([]*node)) {
	if n.pattern != "" {
		*nodes = append(*nodes, n)
	}
	for _, child := range n.children {
		child.traverse(nodes)
	}
}

// TODO: 用于插入（不需要考虑模糊匹配的情况）
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isFuzzy {
			return child
		}
	}
	return nil
}

// TODO: 用于查询（需要考虑模糊匹配情况）
func (n *node) matchChildren(part string) []*node {
	children := make([]*node, 0)
	for _, child :=  range n.children {
		if child.part == part || child.isFuzzy {
			children = append(children, child)
		}
	}
	return children
}

// TODO: 重复插入问题
func (n *node) insert(pattern string, parts []string, height int) {
	if height == len(parts) {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part: part,
			isFuzzy: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}

	child.insert(pattern, parts, height+1)
}

// TODO: 查询不到问题
func (n *node) search(parts []string, height int) *node {
	if height == len(parts) || strings.HasPrefix(n.part, "*") {
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
			return result
		}
	}

	return nil
}