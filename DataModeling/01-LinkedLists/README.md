DM.01 Linked List
=================

In this section, we will implement a linked list in Neo4j.

To start with, we need an empty list with a `ROOT` node to anchor everything else to.

``` Go
stmt := `
	CREATE (root { name: 'ROOT' })-[:LINK]->(root)
`
cq := neoism.CypherQuery{
	Statement: stmt,
}
err := db.Cypher(&cq)
```

We need a function to insert nodes, otherwise we can only have an empty list.

``` Go
func insertNode(val int) {
	stmt := `
		MATCH (root)-[:LINK*0..]->(before),(after)-[:LINK*0..]->(root),(before)-[old:LINK]->(after)
		WHERE root.name = 'ROOT' AND (before.value <  {valSub} OR before = root) AND ( {valSub} < after.value OR after = root)
		CREATE UNIQUE (before)-[:LINK]->({ value: {valSub} })-[:LINK]->(after)
		DELETE old
	`
	params := neoism.Props{"valSub": val}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     nil,
	}

	// execute query
	err := db.Cypher(&cq)
	panicErr(err)
}
```

We probably want a function to delete nodes, so we can do more than grow.

``` Go
func deleteNode(val int) {
	stmt := `
		MATCH
		  (root)-[:LINK*0..]->(before),
		  (before)-[delBefore:LINK]->(del)-[delAfter:LINK]->(after),
		  (after)-[:LINK*0..]->(root)
		WHERE root.name = 'ROOT' AND del.value = {valSub}
		CREATE UNIQUE (before)-[:LINK]->(after)
		DELETE del, delBefore, delAfter
	`
	params := neoism.Props{"valSub": val}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     nil,
	}

	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

}
```

Likely, we will want to know what's in the list.

``` Go
func listListData() {
	// query results
	res := []struct {
		From neoism.Node
		Rel  neoism.Relationship
		To   neoism.Node
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement: `
			MATCH (n)-[r]->(m)
			RETURN n AS From, r AS Rel, m AS To
			ORDER BY n.value
		`,
		Result: &res,
	}
	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	fmt.Println("List Nodes: ", len(res))
	for i, _ := range res {
		n := res[i]
		fmt.Printf("%02d: %v\n", i, n.From.Data["value"])
	}
	fmt.Println()
}
```

Finally, here's a small drive program.

``` Go
func main() {

	insertNode(25)
	insertNode(10)
	insertNode(14)
	insertNode(30)
	insertNode(35)
	insertNode(20)
	insertNode(05)

	listListData()

	deleteNode(05)
	deleteNode(35)
	deleteNode(20)

	listListData()

}
```

I'll let you write a function for determining if a value is in the list.
