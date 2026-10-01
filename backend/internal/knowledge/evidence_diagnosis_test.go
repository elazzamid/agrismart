package knowledge

import (
    "context"
    "testing"

    "github.com/pashagolub/pgxmock/v4"
)

func TestDiagnoseBySymptomsRequiresLatestApprovedVersion(t *testing.T) {
    pool, err := pgxmock.NewPool()
    if err != nil {
        t.Fatal(err)
    }
    defer pool.Close()
    s := NewService(pool)

    pool.ExpectQuery(`(?s)SELECT CASE WHEN ps\.pest_id IS NOT NULL THEN 'pest' ELSE 'disease' END,.*FROM problem_symptoms ps.*problem_symptom_sources pss.*d\.status='published'.*v\.version_no = \(SELECT MAX\(v2\.version_no\) FROM knowledge_versions v2 WHERE v2\.document_id=d\.id\).*kv\.decision='approved'.*AND \(ps\.pest_id IS NOT NULL\) <> \(ps\.disease_id IS NOT NULL\)`).
        WithArgs([]string{"symptom-1"}, "crop-1", 10).
        WillReturnRows(pgxmock.NewRows([]string{"problem_type", "problem_id", "matched_symptoms", "score", "evidence"}).
            AddRow("disease", "disease-1", 1, 2, []string{"leaf spot"}))

    items, err := s.DiagnoseBySymptoms(context.Background(), []string{" symptom-1 "}, "crop-1", 10)
    if err != nil {
        t.Fatal(err)
    }
    if len(items) != 1 || items[0].ProblemID != "disease-1" {
        t.Fatalf("unexpected diagnosis candidates: %+v", items)
    }
    if err := pool.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}

func TestDiagnoseBySymptomsRequiresExplicitEvidenceMapping(t *testing.T) {
    pool, err := pgxmock.NewPool()
    if err != nil {
        t.Fatal(err)
    }
    defer pool.Close()
    s := NewService(pool)

    pool.ExpectQuery(`(?s)FROM problem_symptoms ps.*EXISTS \(.*problem_symptom_sources pss.*pss\.symptom_id=ps\.symptom_id.*d\.status='published'.*v\.version_no = \(SELECT MAX\(v2\.version_no\) FROM knowledge_versions v2 WHERE v2\.document_id=d\.id\).*kv\.decision='approved'`).
        WithArgs([]string{"symptom-1"}, "crop-1", 10).
        WillReturnRows(pgxmock.NewRows([]string{"problem_type", "problem_id", "matched_symptoms", "score", "evidence"}))

    items, err := s.DiagnoseBySymptoms(context.Background(), []string{"symptom-1"}, "crop-1", 10)
    if err != nil {
        t.Fatal(err)
    }
    if len(items) != 0 {
        t.Fatalf("expected no candidates without explicit evidence mapping: %+v", items)
    }
    if err := pool.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}

func TestDiagnoseBySymptomsKeepsEvidenceMappingWhenCropFilterOmitted(t *testing.T) {
    pool, err := pgxmock.NewPool()
    if err != nil {
        t.Fatal(err)
    }
    defer pool.Close()
    s := NewService(pool)

    pool.ExpectQuery(`(?s)FROM problem_symptoms ps.*WHERE ps\.symptom_id::text = ANY\(\$1::text\[\]\).*AND EXISTS \(.*problem_symptom_sources pss.*d\.status='published'.*AND \(\$2='' OR EXISTS \(SELECT 1 FROM knowledge_crop_links k WHERE k\.document_id=d\.id AND k\.crop_id::text=\$2\)\).*v\.version_no = \(SELECT MAX\(v2\.version_no\) FROM knowledge_versions v2 WHERE v2\.document_id=d\.id\).*kv\.decision='approved'`).
        WithArgs([]string{"symptom-1"}, "", 10).
        WillReturnRows(pgxmock.NewRows([]string{"problem_type", "problem_id", "matched_symptoms", "score", "evidence"}).
            AddRow("disease", "disease-1", 1, 2, []string{"leaf spot"}))

    items, err := s.DiagnoseBySymptoms(context.Background(), []string{"symptom-1"}, "", 10)
    if err != nil {
        t.Fatal(err)
    }
    if len(items) != 1 || items[0].ProblemID != "disease-1" {
        t.Fatalf("unexpected diagnosis candidates without crop filter: %+v", items)
    }
    if err := pool.ExpectationsWereMet(); err != nil {
        t.Fatal(err)
    }
}
