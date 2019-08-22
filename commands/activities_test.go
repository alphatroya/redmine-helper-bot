package commands

import "testing"

func TestActivities_Handle(t *testing.T) {
}

func TestActivities_IsCompleted(t *testing.T) {
	redmineMock := &RedmineMock{}
	storageMock := &StorageMock{}
	sut := newActivitiesCommand(redmineMock, storageMock, 1)
	if sut.IsCompleted() != true {
		t.Error("activities command should always be completed")
	}
}
