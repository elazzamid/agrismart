INSERT INTO nutrients (code, name, description) VALUES
('N', 'Nitrogen', 'Unsur makro yang berperan penting dalam pertumbuhan vegetatif dan pembentukan biomassa.'),
('P', 'Phosphorus', 'Unsur makro yang berperan dalam perkembangan akar, energi, dan proses reproduktif tanaman.'),
('K', 'Potassium', 'Unsur makro yang berperan dalam regulasi air, kekuatan jaringan, dan berbagai proses fisiologis tanaman.'),
('Ca', 'Calcium', 'Unsur yang berperan dalam struktur dinding sel dan pertumbuhan jaringan.'),
('Mg', 'Magnesium', 'Unsur penyusun inti molekul klorofil dan berperan dalam fotosintesis.'),
('S', 'Sulfur', 'Unsur yang berperan dalam pembentukan asam amino dan metabolisme tanaman.'),
('Fe', 'Iron', 'Mikronutrien penting dalam proses fisiologis dan pembentukan klorofil.'),
('Zn', 'Zinc', 'Mikronutrien yang berperan dalam aktivitas enzim dan regulasi pertumbuhan.'),
('B', 'Boron', 'Mikronutrien yang berperan dalam pertumbuhan jaringan dan proses reproduktif.'),
('Mn', 'Manganese', 'Mikronutrien yang berperan dalam fotosintesis dan aktivitas enzim.'),
('Cu', 'Copper', 'Mikronutrien yang berperan dalam aktivitas enzim dan metabolisme.'),
('Mo', 'Molybdenum', 'Mikronutrien yang berperan dalam metabolisme nitrogen.'),
('Cl', 'Chlorine', 'Unsur yang dibutuhkan tanaman dalam jumlah kecil dan berperan dalam keseimbangan ion.'),
('Ni', 'Nickel', 'Mikronutrien yang dibutuhkan dalam jumlah sangat kecil dan terkait metabolisme nitrogen.')
ON CONFLICT (code) DO NOTHING;
