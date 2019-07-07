fn main() {
	correct := '${rand.Intn(10)}'

	scanner := bufio.NewScanner(os.Stdin)
	print('Guess the randomly-chosen number between 0 and 9: ')
	mut user_guess := ""
	mut guesses := 0

	for scanner.Scan() {
		user_guess = strings.TrimSpace(scanner.Text())
		guesses++
		if user_guess == correct {
			break
        }
		println('Nope! Try again...')
	}

	if scanner.Err() != nil {
		eprintln('Error getting user input: ${scanner.Err()}')
		exit(1)
	}

	println('The correct answer is $user_guess -- great work!')
	println('It only took you $guesses tries.')
}
