CONFIG_NAME = videoservice

include ../../makefile

## Generate mocks (Specific to video-service)
mocks:
	$(MOCKGEN) -source="./internal/app/video_service.go" -destination="./internal/app/mock/video_service_mock.go" -mock_names=VideoService=MockVideoService
	$(MOCKGEN) -source="./internal/ports/video_repo.go" -destination="./internal/ports/mock/video_repo_mock.go" -mock_names=VideoRepository=MockVideoRepository
	$(MOCKGEN) -source="./internal/ports/id_gen.go" -destination="./internal/ports/mock/id_gen_mock.go" -mock_names=IDGen=MockIDGen
