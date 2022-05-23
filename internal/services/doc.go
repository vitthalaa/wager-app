package services

// Generate dependencies mocks for services
//go:generate mockery --name=IWagerRepo --structname=MockWagerRepo --dir ../repo --filename generated_mock_wager_repo_test.go --testonly --output . --outpkg services
