GS.2 Movie Database
===================

### Our first generic node query function:

Careful! This is not a safe function for making queries.
It is, however, a simple and flexible implementation for the tutorials.

``` Go
func queryNodes(MATCH, WHERE, CREATE, RETURN, ORDERBY string) []struct{ N neoism.Node } {
	stmt := ""
	if MATCH != "" {
		stmt += "MATCH " + MATCH + " "
	}
	if WHERE != "" {
		stmt += "WHERE " + WHERE + " "
	}
	if CREATE != "" {
		stmt += "CREATE " + CREATE + " "
	}
	if RETURN != "" {
		stmt += "RETURN " + RETURN + " "
	}
	if ORDERBY != "" {
		stmt += "ORDERBY " + ORDERBY + " "
	}
	stmt += ";"
	// params
	params := neoism.Props{
		"MATCH":   MATCH,
		"WHERE":   WHERE,
		"CREATE":  CREATE,
		"RETURN":  RETURN,
		"ORDERBY": ORDERBY,
	}

	// query results
	res := []struct {
		N neoism.Node // Column "n" gets automagically unmarshalled into field N
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	return res
}
```

### Counting things:

Let's count all of the nodes.

``` Go
func countNodes() {
	res := queryNodes("(n)", "", "", "n", "")
	fmt.Println("countNodes()", len(res))
}
```

Now let's count the nodes of a particular type.
``` Go
func countNodesByType(typ string) {
	match := "(n:" + typ + ")"
	res := queryNodes(match, "", "", "n", "")
	fmt.Println("countNodesByType()", len(res))
}

countNodesByType("Actor")
countNodesByType("Movie")
```

* Note, these are not the most efficient implementations because
a list of nodes is return and then the length is found in Golang.
Ideally, we would just return the count which requires
a different return type than the `queryNodes()` function uses.
The implementation is left as an exercise for the student.


### Listing things:

Here are some functions for listing actors and movies.

``` Go
func showAllActors() {
	res := queryNodes("(n:Actor)", "", "", "n", "")
	fmt.Println("Actors: ", len(res))
	for i, _ := range res {
		n := res[i].N
		fmt.Printf("  Actor[%d] %+v\n", i, n.Data)
	}
}

func getActorByName(name string) {
	res := queryNodes("(n:Actor)", "n.name = '"+name+"'", "", "n", "")
	fmt.Println("Actors: ", len(res))
	for i, _ := range res {
		n := res[i].N
		fmt.Printf("  Actor[%d] %+v\n", i, n.Data)
	}
}

func showAllMovies() {
	res := queryNodes("(n:Movie)", "", "", "n", "")
	fmt.Println("Movies: ", len(res))
	for i, _ := range res {
		n := res[i].N
		fmt.Printf("  Movie[%d] %+v\n", i, n.Data)
	}
}

func getMovieByName(title string) {
	res := queryNodes("(n:Movie)", "n.title = '"+title+"'", "", "n", "")
	fmt.Println("Actors: ", len(res))
	for i, _ := range res {
		n := res[i].N
		fmt.Printf("  Actor[%d] %+v\n", i, n.Data)
	}
}
```

Here is a function for listing out the entire graph.

``` Go
func listGraphData() {
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
			RETURN n AS From, r AS Rel, m AS To;
		`,
		Result: &res,
	}
	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	fmt.Println("Graph Data: ", len(res))
	for i, _ := range res {
		n := res[i]
		fmt.Printf("  [%d] %+v -> %+v -> %+v\n", i, n.From.Data, n.Rel.Data, n.To.Data)
	}
}
```


### What about relationships?

Well, I'm single and probably shouldn't speak on the subject ;]
As for this tutorial, I have left the relationship functions as
an exercise for the student.

The first task is to write a `queryRelationships(...)` function
which mirrors `queryNodes(...)` except for the return type.
From there, the remaining `count*()` and `list*()` functions
for relationhips should follow naturally.
