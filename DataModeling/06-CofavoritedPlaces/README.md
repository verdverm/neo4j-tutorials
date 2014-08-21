DM.06 Cofavorited Places
========================


### User functions

List a user's favorited places.

```Go
func getFavoriteList(user string) {
	stmt := `
		MATCH (user:User {name: {userSub}})-[:favorite]->(place:Place)
		RETURN place.name
		ORDER BY place.name
	`
	params := neoism.Props{"userSub": user}

	res := []struct {
		Favorite string `json:"place.name"`
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

	fmt.Println(user, "favorites:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Favorite)
		}
	} else {
		fmt.Println("  No results found")
	}
}
```

### Place functions

Co-favorited places - users who like x also like y.

```Go
func cofavoritedPlaces(place string) {
	stmt := `
		MATCH (place:Place)<-[:favorite]-(person:User)-[:favorite]->(other:Place)
		WHERE place.name = {placeSub}
		RETURN other.name, count(*) AS ocount
		ORDER BY ocount DESC, other.name
	`
	params := neoism.Props{"placeSub": place}

	res := []struct {
		Other string `json:"other.name"`
		Count int    `json:"ocount"`
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

	fmt.Println(place, "cofavorites:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s  %4d\n", n.Other, n.Count)
		}
	} else {
		fmt.Println("  No results found")
	}
}```

Co-tagged places - places related through tags

```Go
func cotaggedPlaces(place string) {
	stmt := `
		MATCH (place:Place)-[:tagged]->(tag:Tag)<-[:tagged]-(other:Place)
		WHERE place.name = {placeSub}
		RETURN other.name, collect(tag.name) as tags
		ORDER BY length(collect(tag.name)) DESC, other.name
	`
	params := neoism.Props{"placeSub": place}

	res := []struct {
		Other string   `json:"other.name"`
		Tags  []string `json:"tags"`
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

	fmt.Println(place, "cotagged:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s  %q\n", n.Other, n.Tags)
		}
	} else {
		fmt.Println("  No results found")
	}
}
```
