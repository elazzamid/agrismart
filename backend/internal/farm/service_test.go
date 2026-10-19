package farm

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestServiceGetScopesFarmThroughAuthenticatedUserProfile(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock.NewPool() error = %v", err)
	}
	defer pool.Close()

	s := NewService(pool)
	farmID := "00000000-0000-0000-0000-000000000001"
	userID := "00000000-0000-0000-0000-000000000002"
	profileID := "00000000-0000-0000-0000-000000000003"

	pool.ExpectQuery(`SELECT f.id, f.farmer_id, f.name, COALESCE\(f.location_name, ''\), f.latitude, f.longitude\s+FROM farms f\s+JOIN farmer_profiles fp ON fp.id = f.farmer_id\s+WHERE f.id = \$1 AND fp.user_id = \$2`).
		WithArgs(farmID, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "farmer_id", "name", "location_name", "latitude", "longitude"}).
			AddRow(farmID, profileID, "Kebun Utama", "Desa A", nil, nil))

	got, err := s.Get(context.Background(), userID, farmID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.ID != farmID || got.FarmerID != profileID {
		t.Fatalf("unexpected farm: %+v", got)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}


func TestServiceCreateResolvesAuthenticatedUserToFarmerProfile(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	s := NewService(pool)
	userID := "00000000-0000-0000-0000-000000000001"
	profileID := "00000000-0000-0000-0000-000000000002"
	farmID := "00000000-0000-0000-0000-000000000003"
	pool.ExpectQuery(`INSERT INTO farms \(farmer_id, name, location_name, latitude, longitude\)\s+SELECT fp.id, \$2, NULLIF\(\$3, ''\), \$4, \$5\s+FROM farmer_profiles fp\s+WHERE fp.user_id = \$1\s+RETURNING id, farmer_id, name, COALESCE\(location_name, ''\), latitude, longitude`).
		WithArgs(userID, "Kebun Utama", "Desa A", (*float64)(nil), (*float64)(nil)).
		WillReturnRows(pgxmock.NewRows([]string{"id", "farmer_id", "name", "location_name", "latitude", "longitude"}).
			AddRow(farmID, profileID, "Kebun Utama", "Desa A", nil, nil))

	got, err := s.Create(context.Background(), userID, CreateInput{Name: "Kebun Utama", LocationName: "Desa A"})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != farmID || got.FarmerID != profileID {
		t.Fatalf("unexpected farm: %+v", got)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestServiceListScopesFarmsThroughAuthenticatedUserProfile(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	s := NewService(pool)
	userID := "00000000-0000-0000-0000-000000000001"
	profileID := "00000000-0000-0000-0000-000000000002"
	farmID := "00000000-0000-0000-0000-000000000003"
	pool.ExpectQuery(`SELECT f.id, f.farmer_id, f.name, COALESCE\(f.location_name, ''\), f.latitude, f.longitude\s+FROM farms f\s+JOIN farmer_profiles fp ON fp.id = f.farmer_id\s+WHERE fp.user_id = \$1\s+ORDER BY f.created_at DESC`).
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "farmer_id", "name", "location_name", "latitude", "longitude"}).
			AddRow(farmID, profileID, "Kebun Utama", "Desa A", nil, nil))

	got, err := s.List(context.Background(), userID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != farmID || got[0].FarmerID != profileID {
		t.Fatalf("unexpected farms: %+v", got)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
