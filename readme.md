#  Cove

**Cove** is a lightweight, local-first secret management API written in Go. It allows your personal and home-lab projects to store, retrieve, update, and delete secrets securely through a simple RESTful interface.

##  Features

- Secure secret storage with encryption at rest
- REST API for managing secrets
- Client authentication via bearer token
- JSON file-based vault for quick and minimal setup (no external DB needed)
- Designed for internal/private use

## [CoveClient](https://github.com/LSariol/CoveClient)
CoveClient is a lightweight Go module that simplifies communication with the Cove secret management API.
It provides a clean and reusable interface for reaching Cove’s endpoints, making it easy for projects like Lighthouse and others to securely access and manage secrets without handling raw HTTP logic.
