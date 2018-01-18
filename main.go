package main

func main() {
	bc := NewBlockchain("-1")
	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()
}
