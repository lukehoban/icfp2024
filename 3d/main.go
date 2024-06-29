package main

func do() error {
	return nil
}

func main() {
	err := do()
	if err != nil {
		panic(err)
	}
}
