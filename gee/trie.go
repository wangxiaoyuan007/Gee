package gee

import "strings"

type node struct {
	pattern string //整个路径
	part string
	children []*node
	isWild bool //part 含有 : 或 * 时为true
}

//第一个匹配成功的节点
func (n *node)matchCild(part string) *node  {
	for _,n := range n.children {
		if part == n.part || n.isWild {
			return n
		}
	}
	return nil
}

//第所有匹配成功的节点
func (n *node)matchCildren(part string) []*node  {
	nodes := make([]*node,0)
	for _,n := range n.children {
		if part == n.part || n.isWild {
			nodes = append(nodes,n)
		}
	}
	return nodes
}

func (n *node) insert(pattern string,parts []string, i int)  {
	if i == len(parts) {
		n.pattern = pattern
		return
	}
	part := parts[i]
	firstNode := n.matchCild(part)
	if firstNode == nil {
		firstNode = &node{
			part: part,
			isWild: part[0] == ':' || part[0] == '*',
			children: make([]*node,0),
		}
		n.children = append(n.children,firstNode)
	}
	firstNode.insert(pattern, parts, i + 1)

}

func (n * node) search(parts []string, i int) *node  {
	if len(parts) == i || strings.HasPrefix(n.part,"*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[i]
	children := n.matchCildren(part)
	for _, child := range children {
		result := child.search(parts,i + 1)
		if result != nil {
			return result
		}
	}
	return nil
}

