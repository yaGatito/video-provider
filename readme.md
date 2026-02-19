# VIDEO PROVIDER

Lightweight video service domain and application logic. This repository focuses
on the core domain model, validation and service layer for videos (publisher
scoped operations, global search and pagination). Infrastructure (HTTP API,
storage implementation, auth, upload handling) is intentionally left out and
planned as separate adapters.

## Goals

- Provide a small, testable video domain and service layer that can be
	composed with different storage and transport adapters.
- Support publisher-scoped operations (create, list) and global/publisher
	search with safe validation and pagination.
- Keep business rules (validation, search/pagination policy) centralized and
	covered by unit tests.

## Tasks / Planned work

- Implement concrete `VideoRepository` adapters (SQL, in-memory, etc.).
- Add HTTP API and CLI adapters to expose service endpoints.
- Add authentication/authorization and upload handling (storage + streaming).
- Add user, comments, playlists and social features as separate domain modules.
- Add CI pipeline and containerization for delivery (Docker image, Makefile).

## Implemented features (what's already in this repo)

- Domain model `Video` with fields: ID, PublisherID, Topic, Description,
	CreatedAt and Status (draft/published) — [internal/video-service/domain/video.go](internal/video-service/domain/video.go#L1-L50)
- Validation rules for `Video` (topic and description formats, publisher ID)
	using centralized policy regexes — [internal/video-service/domain/video_validation.go](internal/video-service/domain/video_validation.go#L1-L60)
- Business interactor `VideoService` (`VideoInteractor`) with methods:
	`Create`, `GetByID`, `GetByPublisher`, `SearchPublisher`, `SearchGlobal` —
	[internal/video-service/app/video_service.go](internal/video-service/app/video_service.go#L1-L200)
- Search and pagination policy with constants and regex helpers
	(`policy` package) — [internal/video-service/policy/constants.go](internal/video-service/policy/constants.go#L1-L80)
- Repository boundary `VideoRepository` interface and request types
	(`PageRequest`, `VideoSearch`) — [internal/video-service/ports/video_repo.go](internal/video-service/ports/video_repo.go#L1-L80)
- Generated GoMock mocks for `VideoRepository` to simplify unit testing
	— [internal/video-service/ports/mock/video_repo_mock.go](internal/video-service/ports/mock/video_repo_mock.go#L1-L220)
- Unit tests covering the interactor and search/pagination logic —
	[internal/video-service/app/video_service_test.go](internal/video-service/app/video_service_test.go#L1-L400)

## How to run tests

Run all tests in the module:

```bash
go test ./... -v
```

## Notes

- This repository provides only the core domain and application logic. To run
	a complete system you will need to implement adapters (HTTP server, DB
	repository, object storage for uploads) and wire them to the `VideoService`.
- If you want, I can add a minimal in-memory repository and a small HTTP
	server example to demonstrate end-to-end usage.


