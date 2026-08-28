INSERT INTO symptoms (name, description) VALUES
('daun menguning','Daun menunjukkan perubahan warna menjadi kuning.'),
('daun mengeriting','Daun menggulung atau mengalami perubahan bentuk.'),
('bercak gelap pada buah','Buah memiliki bercak berwarna gelap atau kecokelatan.'),
('buah membusuk','Jaringan buah mengalami pembusukan.'),
('serangga kecil pada daun','Terlihat serangga berukuran kecil pada permukaan daun.'),
('tanaman layu','Tanaman kehilangan turgor dan tampak layu.'),
('pertumbuhan terhambat','Pertumbuhan tanaman lebih lambat dari kondisi normal.')
ON CONFLICT (name) DO NOTHING;
