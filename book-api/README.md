Books API (Go + Docker + Kubernetes)

Overview =>

This is a RESTful API for managing books, built with Go, containerized with Docker, and deployed on Kubernetes using Minikube. The API supports CRUD operations and search functionality with concurrency optimization.

How to run (Without Docker) =>

1. Install GO(v1.24+)

2. Clone the repository and navigate to the project folder.

3. Install dependencies (if applicable):

4. Run the server:

5. API available at: http://localhost:8000

How to run (With Docker) => 

1. Build the Docker image:

2. Run the container:

3. API available at: http://localhost:8000

Deploy on Minikube =>

1. Use Minikube's Docker environment:

2. Build the Docker image inside Minikube:

3. Deploy to Kubernetes:

4. Get Minikube service URL:

4. API available at the displayed URL. e: http://127.0.0.1:53332/books

Techs used => 

1. Golang for backend API

2. Docker for containerization

3. Kubernetes (Minikube) for deployment

Notes =>

1. Ensure Minikube is running before deploying to Kubernetes.

2. Use kubectl logs -f <podsName> eg : books-api-6b5b47867c-zbzf7 to check application logs.

AUTHOR 

DEVELOPED AND TESTED BY JINUKA WEERASINGHE