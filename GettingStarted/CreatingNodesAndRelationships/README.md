GS.1 Creating nodes and relationships
=====================================

## Create a node for the actor Tom Hanks:

``` Go
	// query statemunt
	stmt := `
		CREATE (actor:Actor { name:{actorSub}})
		RETURN actor
	`
	// query params
	actor := "Tom Hanks"
	params := neoism.Props{"actorSub": actor}

	// query results
	res := []struct {
		Actor neoism.Node
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

	// check results
	if len(res) != 1 {
		panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 1\n", len(res)))
	}

	n := res[0].Actor // Only one row of data returned
	fmt.Println("createNode()", n.Data)
```

## Find the node we just created:

``` Go
	// query statemunt
	stmt := `
		MATCH (actor:Actor)
		WHERE actor.name = {actorSub}
		RETURN actor
	`
	// query params
	actor := "Tom Hanks"
	params := neoism.Props{"actorSub": actor}

	// query results
	res := []struct {
		Actor neoism.Node
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

	// check results
	if len(res) != 1 {
		panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 1\n", len(res)))
	}

	n := res[0].Actor // Only one row of data returned
	fmt.Printf("queryNode() -> %+v\n", n.Data)
```

## Create a movie and connect it to an actor in one query
``` Go
	actor := "Tom Hanks"
	movie := "Sleepless in Seattle"

	// query statemunt
	stmt := `
		MATCH (actor:Actor)
		WHERE actor.name = {actorSub}
		CREATE (movie:Movie {title: {movieSub}})
		CREATE (actor)-[:ACTED_IN]->(movie);
	`
	// query params
	params := neoism.Props{
		"actorSub": actor,
		"movieSub": movie,
	}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     nil,
	}
	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	fmt.Println("createMovie()")
```

## Same as the last, with less of a query statement:
also note the use of ``json:"tags"``
``` Go
	actor := "Tom Hanks"
	movie := "Forrest Gump"

	// query statemunt
	stmt := `
		MATCH (actor:Actor {name: {actorSub}})
		CREATE UNIQUE (actor)-[r:ACTED_IN]->(movie:Movie {title: {movieSub}})
		RETURN actor.name, type(r), movie.title;
	`
	// query params
	params := neoism.Props{
		"actorSub": actor,
		"movieSub": movie,
	}

	// query results
	res := []struct {
		// `json` tag matches column name in query
		Name  string `json:"actor.name"`
		Rel   string `json:"type(r)"`
		Movie string `json:"movie.title"`
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

	r := res[0]
	fmt.Println("createUnique()", r.Name, r.Rel, r.Movie)
```

## Set a property on a node:

``` Go
	actor := "Tom Hanks"
	dob := 1944

	// query statemunt
	stmt := `
		MATCH (actor:Actor {name: {actorSub}})
		SET actor.DoB = {dobSub}
		RETURN actor.name, actor.DoB;
	`
	// query params
	params := neoism.Props{
		"actorSub": actor,
		"dobSub":   dob,
	}

	// query results
	res := []struct {
		Name string `json:"actor.name"` // `json` tag matches column name in query
		DoB  string `json:"actor.DoB"`
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

	r := res[0]
	fmt.Println("setNodeProperty()", r.Name, r.DoB)
```

## List all Movies

``` Go
	// query statemunt
	stmt := `
		MATCH (movie:Movie)
		RETURN movie;
	`
	// query params
	actor := "Tom Hanks"
	params := neoism.Props{"actorSub": actor}

	// query results
	res := []struct {
		Movie neoism.Node
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

	// check results
	if len(res) != 2 {
		panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 2\n", len(res)))
	}

	fmt.Printf("queryMovies()\n")
	for i, _ := range res {
		n := res[i].Movie // Only one row of data returned
		fmt.Printf("  Node[%d] %+v\n", i, n.Data)
	}
```
