package main

import ("github.com/mhernan88/dag-bisect/dag")
// import (
//   "github.com/urfave/cli/v2"
// )
//
// var flags = []cli.Flag{
//   &cli.StringFlag{
//     Name:  "input",
//     Value: "dag.json",
//     Usage: "serialization of dag",
//   },
// },

func main() {
  p := dag.Load("dag2.json")
  p.Link()
  p.Distance(name)
}
