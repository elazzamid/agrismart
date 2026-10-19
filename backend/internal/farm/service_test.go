package farm

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestServiceGetScopesFarmThroughAuthenticatedUser(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatalf("pgxmock.NewPool() error = %v", err) }
	defer pool.Close()

	s := NewService(pool)
	farmID := "00000000-0000-0000-0000-000000000001"
	userID := "10000000-0000-0000-0000-000000000001"
	profileID := "20000000-0000-0000-0000-000000000002"

	pool.ExpectQuery(`SELECT f.id, f.farmer_id, f.name, COALESCE\(f.location_name, ''\), f.latitude, f.longitude\s+FROM farms f\s+JOIN farmer_profiles fp ON fp.id = f.farmer_id\s+WHERE f.id = \$1 AND fp.user_id = \$2`).
		WithArgs(farmID, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "farmer_id", "name", "location_name", "latitude", "longitude"}).
			AddRow(farmID, profileID, "Kebun Utama", "Desa A", nil, nil))

	got, err := s.Get(context.Background(), userID, farmID)
	if err != nil { t.Fatalf("Get() error = %v", err) }
	if got.ID != farmID || got.FarmerID != profileID { t.Fatalf("unexpected farm: %+v", got) }
	if got.FarmerID == userID { t.Fatal("farm foreign key must use farmer_profiles.id, not users.id") }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatalf("expectations: %v", err) }
}
