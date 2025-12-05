package main

import "C"

//export Run
func Run() {
	if err := rootCmd.Execute(); err != nil {
		// In a C-shared library, we probably shouldn't os.Exit(1) directly if we want to handle errors gracefully in Python,
		// but since the original main calls os.Exit, we'll keep it simple for now or just let it return.
		// However, cobra's Execute prints errors to stderr by default.
	}
}
