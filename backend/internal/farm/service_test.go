package farm

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestServiceGetScopesFarmByFarmer(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatalf("pgxmock.NewPool() error = %v", err) }
	defer pool.Close()

	s := NewService(pool)
	farmID := "00000000-0000-0000-0000-000000000001"
	farmerID := "00000000-0000-0000-0000-000000000002"

	pool.ExpectQuery(`SELECT id, farmer_id, name, COALESCE\(location_name, ''\), latitude, longitude\s+FROM farms WHERE id = \$1 AND farmer_id = \$2`).
		WithArgs(farmID, farmerID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "farmer_id", "name", "location_name", "latitude", "longitude"}).
			AddRow(farmID, farmerID, "Kebun Utama", "Desa A", nil, nil))

	got, err := s.Get(context.Background(), farmerID, farmID)
	if err != nil { t.Fatalf("Get() error = %v", err) }
	if got.ID != farmID || got.FarmerID != farmerID { t.Fatalf("unexpected farm: %+v", got) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatalf("expectations: %v", err) }
}
