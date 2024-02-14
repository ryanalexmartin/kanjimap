# KanjiMap
This application was made because I couldn't find anything else like it.
I wanted to be able to see at a glance how many Chinese characters I've learned, and how far I have to go.  I hope you'll find it useful.

I'll accept any pull requests, and comments/suggestions are encouraged.  Feel free to email me, this is just a fun side project of mine
but I should have some time to work on it.

## Project Structure

- `go/`: Contains the Go backend which handles user login and tracks user progress.
- `vue/`: Contains the Vue.js frontend which displays the Chinese characters and user progress.
- `sql/`: Contains the MySQL database setup and Python scripts for updating the characters table.
- `docker-compose.yml`: Used to manage the application services.

## Getting Started

### Prerequisites

- Node.js and npm
- Go
- Docker
- Python

### Installation

1. Clone the repository.
2. Install the dependencies for the Vue.js application by running `npm install` in the `vue/` directory.
3. Install the dependencies for the Python scripts by running `pip install -r requirements.txt` in the `sql/` directory.

### Running the Application

1. Start the Vue.js application by running `npm run serve` in the `vue/` directory.
2. Start the Go backend by running `go run main.go` in the `go/` directory.
3. Start the MySQL database by running `docker-compose up` in the root directory.
4. On first run, you should run `python3 update_characters_table.py` to populate the database with the necessary Chinese characters.

## Usage

Once the application is running, you can visit `http://localhost:8080` in your web browser.  Register and log in (it's an old school registration system, no anti-patterns here!) and start clicking on Chinese characters to learn them.  

## Testing

You can run the tests for the Vue.js application by running `npm run test` in the `vue/` directory. The Go backend tests can be run with `go test` in the `go/` directory.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details# Chinese Characters App

This application is a multi-component system for learning Chinese characters. It includes a Vue.js frontend, a Go backend, and a MySQL database managed with Python scripts.

## Project Structure

- `go/`: Contains the Go backend which handles user login and tracks user progress.
- `vue/`: Contains the Vue.js frontend which displays the Chinese characters and user progress.
- `sql/`: Contains the MySQL database setup and Python scripts for updating the characters table.
- `docker-compose.yml`: Used to manage the application services.

## Getting Started

### Prerequisites

- Node.js and npm
- Go
- Docker
- Python

### Installation

1. Clone the repository.
2. Install the dependencies for the Vue.js application by running `npm install` in the `vue/` directory.
3. Install the dependencies for the Python scripts by running `pip install -r requirements.txt` in the `sql/` directory.

### Running the Application

1. Start the Vue.js application by running `npm run serve` in the `vue/` directory.
2. Start the Go backend by running `go run main.go` in the `go/` directory.
3. Start the MySQL database by running `docker-compose up` in the root directory.

## Usage

Once the application is running, you can visit `http://localhost:8080` in your web browser to start learning Chinese characters. Click on a character to mark it as learned, and click again to mark it as unlearned.

## Testing

You can run the tests for the Vue.js application by running `npm run test` in the `vue/` directory. The Go backend tests can be run with `go test` in the `go/` directory.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details