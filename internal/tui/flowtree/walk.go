package flowtree

// walkNodes traverses the flow nodes in depth-first order, calling the given function for each node.
// If the function returns false, the traversal stops for that branch.
func walkNodes(fn func(node *flowNode, idx int) bool, node *flowNode) {
	idx := 0
	var walkFn func(node *flowNode)
	walkFn = func(node *flowNode) {
		descend := fn(node, idx)
		idx++
		if descend {
			for _, innerFlow := range node.innerActions {
				walkFn(innerFlow)
			}
		}
		for _, branch := range node.branches {
			walkFn(branch)
		}
	}
	walkFn(node)
}

// getVisibleNodeCount returns the number of nodes that are currently visible (expanded).
func getVisibleNodeCount(node *flowNode) int {
	count := 0
	walkNodes(func(node *flowNode, idx int) bool {
		count++
		return node.isExpanded
	}, node)
	return count
}
