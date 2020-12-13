# util/reader

Sample usage:

`in := reader.NewLineReader("sample.txt")`
`for {
    n, line, ok := in.ReadLine()
    if !ok {
        break
    }
    fmt.Printf("%v: %v\n", n, line)
    if line[0] == "^" {
        in.UnreadLine("caret")
    }
}`
