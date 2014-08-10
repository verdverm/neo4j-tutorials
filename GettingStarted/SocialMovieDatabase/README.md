GS.3 Social Movie Database
==========================

In this section, we expand the movie database with features typical of a social website.
This means our graph will now have more node types and multiple relationship types.
We will also write some generic functions for adding information to the graph.

### Creating things:

``` Go
func createUser(name string) {
	queryNodes("", "", "(n:User {name: '"+name+"'})", "n", "")
}

func createFriendship(user, friend string) {
	match := "(u:User {name: '" + user + "'}),(f:User {name: '" + friend + "'})"
	create := "(u)-[:FRIEND]->(f)"
	queryNodes(match, "", create, "", "")
}

func createRating(user, title, stars, comment string) {
	match := "(u:User {name: '" + user + "'}),(m:Movie {title: '" + title + "'})"
	create := "(u)-[:RATED { stars: " + stars + ", comment: '" + comment + "'}]->(m)"
	queryNodes(match, "", create, "", "")
}
```

### Listing user information:

``` Go
func getRatingsByUser(user string) {
	stmt := `
		MATCH (u:User {name: {userSub}}),(u)-[rating:RATED]->(movie)
	    RETURN movie, rating;
	`
	params := neoism.Props{"userSub": user}

	res := []struct {
		Movie  neoism.Node
		Rating neoism.Relationship
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

	fmt.Println("User Ratings: ", user, len(res))
	for i, _ := range res {
		m := res[i].Movie.Data
		r := res[i].Rating.Data.(map[string]interface{})
		fmt.Printf("  [%d] %v    %v    %v\n",
			i, m["title"], r["stars"], r["comment"])
	}
}

func getFriendsByUser(user string) {
	stmt := `
		MATCH (u:User {name: {userSub}}),(u)-[r:FRIEND]->(f)
	    RETURN type(r) AS T, f.name AS F;
	`
	params := neoism.Props{"userSub": user}

	// query results
	res := []struct {
		T string
		F string
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

	fmt.Println("User Friends: ", user, len(res))
	for i, _ := range res {
		n := res[i]
		fmt.Printf("  [%d] %q  %q\n", i, n.T, n.F)
	}
}
```

### Limitations and exercises:

The above functions make use of our `queryNodes()` function,
but do not make use of Neo4j's parameter capabilities.
To do so requires a more complicated `queryNodes()` function,
which makes `neoism.Props{}` or a `map[string]string` a parameter of the function.
We leave this as an exercise for the reader.
