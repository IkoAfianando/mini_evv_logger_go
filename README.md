# Mini EVV Logger - Backend (Go)

This repository contains the backend source code for the "Mini EVV Logger – Caregiver Shift Tracker" home assignment by Blue Horn Tech.

The backend is a RESTful API built with Go, designed to manage caregiver schedules, visits, and care tasks. It uses an in-memory data store for simplicity and is fully documented with Swagger UI.

---

### Live Demo URL

**Backend API Base URL:** `https://evv.ikoafianando.my.id/`

**Swagger API Documentation:** `https://evv.ikoafianando.my.id/swagger/index.html`

---

### Tech Stack & Key Decisions

* **Language:** **Go (Golang)**
    * **Reasoning:** Chosen for its performance, simplicity, strong concurrency model, and static typing, making it an excellent choice for building robust and scalable APIs.

* **Framework:** **Fiber v2**
    * **Reasoning:** A high-performance web framework built on top of Fasthttp. Its Express.js-like API design allows for rapid development while maintaining excellent speed. It's lightweight and has great support for middleware.

* **Database:** **In-Memory Store**
    * **Reasoning:** To fulfill the project requirement of a simple, zero-dependency setup, an in-memory data store was implemented using Go maps and a `sync.Mutex`. The mutex ensures data consistency by safely handling concurrent API requests. This approach makes the application easy to clone and run instantly.

* **API Specification:** **Swagger (OpenAPI)**
    * **Reasoning:** To meet the bonus requirement for API documentation, Swagger UI is integrated using `fiber-swagger`. It provides interactive, self-documenting API endpoints, which makes it easy for developers (and reviewers) to explore and test the API directly from the browser.

* **Testing:** **Go's built-in `testing` package + `stretchr/testify`**
    * **Reasoning:** To fulfill the bonus requirement for unit tests, a comprehensive test suite was developed. The tests are written as end-to-end HTTP tests to validate the complete request-response cycle for every endpoint. The `stretchr/testify` library is used for its fluent and readable assertion API (`assert`), which makes test cases cleaner and more maintainable.

* **Logging & Error Handling:**
    * **Reasoning:** Basic structured logging is implemented using the standard `log` package. For API responses, errors are returned as structured JSON objects (e.g., `{"error": "message"}`) with appropriate HTTP status codes, which is a standard practice for REST APIs.

---

### Setup and Local Installation

Follow these instructions to run the backend server locally.

**Prerequisites:**
* **Go:** Version 1.23 or newer.

**Instructions:**
1.  **Clone the repository:**
    ```sh
    git clone https://github.com/IkoAfianando/mini_evv_logger_go
    cd evv-logger-backend
    ```

2.  **Install dependencies:**
    This command will download all the necessary modules defined in `go.mod`.
    ```sh
    go mod tidy
    ```

3.  **Run the server:**
    ```sh
    go run main.go
    ```

4.  **Access the application:**
    * The server will start on `http://localhost:8080`.
    * You will see a log message confirming the server is running.
    * The interactive Swagger UI documentation is available at:
      **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

**Running Tests:**
To run the complete test suite, execute the following command from the root directory:
```sh
go test -v ./...