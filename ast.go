package main

import "github.com/pointlander/peg/tree"

func peg2ast(pegGrammar string) (*node32, error) {
	p := &Peg{
		Tree:   tree.New(false, false, false),
		Buffer: pegGrammar,
	}
	p.Init(Pretty(true), Size(1<<15))
	if err := p.Parse(); err != nil {
		return nil, err
	}
	p.Execute()

	return p.AST(), nil
}

func getRawNodeValue(node *node32, buffer string) string {
	return buffer[node.begin:node.end]
}

// Reverses linked list
func reverse(node *node32) *node32 {
	var head *node32
	for node != nil {
		next := node.next
		node.next = head
		head = node
		node = next
	}
	return head
}

// DFS prefix
func traverse(root *node32, f func(node *node32)) {
	for root != nil {
		f(root)
		if root.up != nil {
			traverse(root.up, f)
		}
		root = root.next
	}
}

// Returns all nodes matching the predicate
func collect(root *node32, predicate func(node *node32) bool) []*node32 {
	output := make([]*node32, 0, 10)
	traverse(root, func(node *node32) {
		if predicate(node) {
			output = append(output, node)
		}
	})
	return output
}

func remove(root *node32, ruleType pegRule) *node32 {
	if root == nil {
		return nil
	} else if root.pegRule == ruleType {
		return nil
	} else {
		root.up = remove(root.up, ruleType)
		root.next = remove(root.next, ruleType)
		return root
	}
}
