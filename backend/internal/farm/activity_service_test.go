package farm

import (
	"context"
	"strings"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestFarmActivityCreateScopesToOwnedCropCycle(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer pool.Close()

	service := NewFarmActivityService(pool)
	userID := "00000000-0000-0000-0000-000000000001"
	cycleID := "00000000-0000-0000-0000-000000000002"
	activityID := "00000000-0000-0000-0000-000000000003"
	pool.ExpectQuery(`INSERT INTO farm_activities.*WHERE cc.id = \$1 AND fp.user_id = \$2 RETURNING`).
		WithArgs(cycleID, userID, "fertilizing", "2026-10-19", "NPK").
		WillReturnRows(pgxmock.NewRows([]string{"id", "crop_cycle_id", "activity_type", "activity_date", "notes"}).
			AddRow(activityID, cycleID, "fertilizing", "2026-10-19", "NPK"))

	got, err := service.Create(context.Background(), userID, cycleID, CreateFarmActivityInput{
		ActivityType: " Fertilizing ", ActivityDate: "2026-10-19", Notes: " NPK ",
	})
	if err != nil { t.Fatal(err) }
	if got.ID != activityID || got.CropCycleID != cycleID || got.ActivityType != "fertilizing" || got.Notes != "NPK" {
		t.Fatalf("unexpected activity: %+v", got)
	}
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

func TestFarmActivityCreateRejectsInvalidInputBeforeDatabaseCall(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer pool.Close()

	service := NewFarmActivityService(pool)
	validDate := "2026-10-19"
	for _, in := range []CreateFarmActivityInput{
		{ActivityDate: validDate},
		{ActivityType: "watering", ActivityDate: "2026-02-30"},
		{ActivityType: "watering", ActivityDate: validDate, Notes: strings.Repeat("x", 2001)},
	} {
		if _, err := service.Create(context.Background(), "user", "cycle", in); err != ErrInvalidFarmActivityInput {
			t.Fatalf("Create(%+v) error = %v, want ErrInvalidFarmActivityInput", in, err)
		}
	}
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

func TestFarmActivityListHidesUnownedCropCycle(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer pool.Close()

	service := NewFarmActivityService(pool)
	pool.ExpectQuery(`SELECT EXISTS.*WHERE cc.id = \$1 AND fp.user_id = \$2`).
		WithArgs("cycle-foreign", "user-1").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	_, err = service.List(context.Background(), "user-1", "cycle-foreign")
	if err != ErrFarmActivityNotFound { t.Fatalf("List error = %v, want ErrFarmActivityNotFound", err) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

func TestFarmActivityListReturnsActivitiesForOwnedCropCycle(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer pool.Close()

	service := NewFarmActivityService(pool)
	cycleID := "00000000-0000-0000-0000-000000000002"
	activityID := "00000000-0000-0000-0000-000000000003"
	pool.ExpectQuery(`SELECT EXISTS.*WHERE cc.id = \$1 AND fp.user_id = \$2`).
		WithArgs(cycleID, "user-1").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
	pool.ExpectQuery(`SELECT id, crop_cycle_id.*WHERE crop_cycle_id = \$1 ORDER BY activity_date DESC, created_at DESC`).
		WithArgs(cycleID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "crop_cycle_id", "activity_type", "activity_date", "notes"}).
			AddRow(activityID, cycleID, "irrigation", "2026-10-19", ""))

	got, err := service.List(context.Background(), "user-1", cycleID)
	if err != nil { t.Fatal(err) }
	if len(got) != 1 || got[0].ID != activityID || got[0].ActivityType != "irrigation" {
		t.Fatalf("unexpected activities: %+v", got)
	}
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}
