-- Structural OPT catalog for Cabai. No pesticide dose or spray schedule is encoded here.
INSERT INTO pests (name, scientific_name, description)
VALUES
 ('Thrips', 'Thrips parvispinus', 'Hama yang dapat merusak daun dan jaringan tanaman cabai.'),
 ('Kutu daun persik', 'Myzus persicae', 'Kutu daun yang dapat mengisap cairan tanaman dan berperan dalam penyebaran beberapa virus.'),
 ('Kutu kebul', 'Bemisia tabaci', 'Serangga pengisap yang berasosiasi dengan gangguan pada tanaman dan penularan virus.'),
 ('Ulat grayak', 'Spodoptera litura', 'Larva pemakan daun yang dapat menyebabkan kerusakan jaringan daun.'),
 ('Lalat buah', NULL, 'Kelompok lalat buah yang dapat menyebabkan kerusakan pada buah cabai.')
ON CONFLICT DO NOTHING;

INSERT INTO diseases (name, scientific_name, description)
VALUES
 ('Antraknosa', NULL, 'Penyakit penting pada buah cabai yang dapat menyebabkan bercak dan pembusukan.'),
 ('Layu bakteri', NULL, 'Penyakit layu yang disebabkan oleh bakteri dan dapat menurunkan kesehatan tanaman.'),
 ('Layu Fusarium', 'Fusarium spp.', 'Penyakit layu yang berkaitan dengan patogen Fusarium.'),
 ('Virus kuning', NULL, 'Gangguan virus pada cabai yang dapat menyebabkan gejala perubahan warna dan pertumbuhan tanaman.'),
 ('Bercak bakteri', NULL, 'Penyakit bakteri yang dapat menimbulkan bercak pada jaringan tanaman.')
ON CONFLICT DO NOTHING;
