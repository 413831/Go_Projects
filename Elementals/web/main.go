package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", render)

	http.ListenAndServe(":80", nil)
}

func render(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("q") != "" {
		fmt.Fprintf(w, renderHello(), r.FormValue("q"))
		return
	}

	if r.Method == "POST" {
		count, _ := strconv.Atoi(r.FormValue("counter"))
		count++

		fmt.Fprintf(w, renderCounterForm(), count)
		return
	}

	fmt.Fprintf(w, renderForm())
}

func renderHello() string {
	return `
			<body>
				<h1>Hola %s !</h1>
			</body>
			`
}

func renderForm() string {
	return `
			<body>
				<form action="/" method="GET">
					<label>Ingresa tu nombre</label>
					<input type="text" name="q">
					<button type="submit">Enviar</button>
				</form>

				<form action="/" method="POST">
					<label>Contador</label>
					<input name="counter" value="1" readonly>
					<button type="submit">Sumar</button>
				</form>
			</body>
			<a href="/">Reset</a>
			`
}

func renderCounterForm() string {
	return `<body>
				<form action="/" method="POST">
				<label>Contador</label>
				<input name="counter" value="%d" readonly>
				<button type="submit">Sumar</button>
				</form>
				<a href="/">Reset</a>
			</body>
			`
}
