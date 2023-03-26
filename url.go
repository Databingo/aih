package main

import (
 "fmt"
 "github.com/rocketlaunchr/google-search"
)

func main() {
  
 query := "Geoffrey Hinton's AI Lectures video"
 ops :=  googlesearch.SearchOptions{ProxyAddr:"socks5://127.0.0.1:7890"}
 results, _ := googlesearch.Search(nil, query, ops)
 //js := fmt.Sprintf("%+v", results)
 fmt.Println("------up-to-date------")
 for index, i := range results{
  fmt.Print("[", index, "] ")
  fmt.Println(i.URL)
  fmt.Println(i.Title)
 }








}

