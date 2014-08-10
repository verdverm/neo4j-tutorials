GS.4 Finding Paths
==================

In this section, we implement a few functions for exploring paths within the graph.
Notice the use of:
 - aggregate methods in the query statement
 - parameters are passed as `neoism.Props{}`
 - results are returned as a slice of anonymous structs
 - `RETURN` makes use of `AS` instead of `json:tags`


### Find other movies for actors in a movie:

Pick a movie.
Then find all the other movies that the actors were in.
For each secondary movie, include how many actors there are
and a list of their names.

``` Go

func getOtherMoviesViaActors(movie string) {
	stmt := `
		MATCH (:Movie { title: {movieSub} })<-[:ACTS_IN]-(actor)-[:ACTS_IN]->(movie)
		RETURN movie.title AS Title, collect(actor.name) AS Actors, count(*) AS Count
		ORDER BY count(*) DESC ;
	`
	params := neoism.Props{"movieSub": movie}

	// query results
	res := []struct {
		Title  string
		Actors []string
		Count  int
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

	fmt.Println("Movies: ", movie, len(res))
	for i, _ := range res {
		n := res[i]
		fmt.Printf("  [%d] %24q  %d  %v\n", i, n.Title, n.Count, n.Actors)
	}
}

### Find co-actors of actors in a movie:

Pick a movie.
For each actor,
return a unique list of actors they have worked with
in any movie.

``` Go
func getCoActingFromMovie(movie string) {
	stmt := `
		MATCH (:Movie { title: {movieSub} })<-[:ACTS_IN]-(actor)-[:ACTS_IN]->(movie)<-[:ACTS_IN]-(colleague)
		RETURN actor.name AS Actor, collect(DISTINCT colleague.name) AS Actors;
	`
	params := neoism.Props{"movieSub": movie}

	// query results
	res := []struct {
		Actor  string
		Actors []string
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

	fmt.Println("Movies: ", movie, len(res))
	for i, _ := range res {
		n := res[i]
		fmt.Printf("  [%d] %24q  %v\n", i, n.Actor, n.Actors)
	}
}
```

### Paths between actors:

Pick two actors.
Find the ten shortest paths between these actors,
through movies and co-actors.
A.k.a the bacon path.

``` Go
func getActorPaths(actor1, actor2 string) {
	stmt := `
		MATCH p =(:Actor { name: {actor1Sub} })-[:ACTS_IN*0..5]-(:Actor { name: {actor2Sub} })
		RETURN extract(n IN nodes(p)| coalesce(n.title,n.name)) AS List, length(p) AS Len
		ORDER BY length(p)
		LIMIT 10;
	`
	params := neoism.Props{"actor1Sub": actor1, "actor2Sub": actor2}

	// query results
	res := []struct {
		List []string
		Len  int
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

	fmt.Println("Paths: ", actor1, actor2)
	for i, _ := range res {
		n := res[i]
		fmt.Printf("  [%d] %d  %v\n", i, n.Len, n.List)
	}
}
```
