package handlers

// Generate dependencies mocks for handlers
//go:generate mockery --name=IWagerService --structname=MockWagerService --dir ../services --filename generated_mock_wager_service_test.go --testonly --output . --outpkg handlers
