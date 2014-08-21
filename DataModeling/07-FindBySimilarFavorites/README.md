DM.07 Find By Similar Favorites
===============================


### List a user's favorited places.

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

### List the users and number of similar favorites to a user

``` Go
func findSameFavoriteUsersOfUser(user string) {
	stmt := `
		MATCH (user:User)-[:favorite]->(thing)<-[:favorite]-(other:User)
		WHERE user.name = {userSub}
		RETURN other.name, count(thing) AS ocount
		ORDER BY ocount DESC, other.name
	`
	params := neoism.Props{"userSub": user}

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

	fmt.Println(user, "friends    #cofavorites:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s  %4d\n", n.Other, n.Count)
		}
	} else {
		fmt.Println("  No results found")
	}
}
```

### List the friends and number of similar favorites to a user

Same as last, but restricted to friends

``` Go
func findSameFavoriteUsersOfUser(user string) {
	stmt := `
		MATCH (user:User)-[:favorite]->(thing)<-[:favorite]-(other:User)
		WHERE user.name = {userSub}
		RETURN other.name, count(thing) AS ocount
		ORDER BY ocount DESC, other.name
	`
	params := neoism.Props{"userSub": user}

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

	fmt.Println(user, "friends    #cofavorites:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s  %4d\n", n.Other, n.Count)
		}
	} else {
		fmt.Println("  No results found")
	}
}
```


### List of users who favorited a thing

```Go
func findThingFavoriteList(thing string) {
	stmt := `
		MATCH (user:User)-[:favorite]->(thing:Thing {name: {thingSub}})
		RETURN user.name
		ORDER BY user.name
	`
	params := neoism.Props{"thingSub": thing}

	res := []struct {
		Name string `json:"user.name"`
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

	fmt.Println(thing, "favorited by:", len(res))
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Name)
		}
	} else {
		fmt.Println("  No results found")
	}

}
```
