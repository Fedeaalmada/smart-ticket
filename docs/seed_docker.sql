SET NAMES utf8mb4;
USE ticketek_db;

SET SQL_SAFE_UPDATES = 0;
SET FOREIGN_KEY_CHECKS = 0;

DELETE FROM transferencias;
DELETE FROM entradas;
DELETE FROM sectores;
DELETE FROM eventos;
DELETE FROM usuarios;

ALTER TABLE eventos AUTO_INCREMENT = 1;
ALTER TABLE sectores AUTO_INCREMENT = 1;
ALTER TABLE usuarios AUTO_INCREMENT = 1;

SET FOREIGN_KEY_CHECKS = 1;

INSERT INTO usuarios (id, nombre, email, password_hash, rol, activo) VALUES
(1, 'Admin', 'admin@smartticket.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'administrador', 1),
(2, 'Juan Perez', 'juan@email.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'cliente', 1),
(3, 'Maria Garcia', 'maria@email.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'cliente', 1);

INSERT INTO eventos (id, titulo, descripcion, foto_url, fecha, horario, duracion_minutos, categoria, capacidad_total, estado, creado_por) VALUES
(1, 'Tan Bionica - La Ultima Noche Magica Tour', 'La banda mas convocante del rock nacional regresa con su gira de despedida. Una noche unica e irrepetible.', 'https://cdn.efstatic.com/VistasNew/Contenido/tanbioonica.jpeg', '2026-08-10', '21:00:00', 150, 'Concierto', 5000, 'activo', 1),
(2, 'River vs Boca - Superclasico', 'El partido mas esperado del ano en el Monumental.', 'https://media.elpatagonico.com/p/9b3c2a0383ab1a7187a54ea268c24ac7/adjuntos/193/imagenes/036/334/0036334699/rivjpeg.jpeg', '2026-09-05', '17:00:00', 120, 'Deporte', 3000, 'activo', 1),
(3, 'Comic Con Argentina 2026', 'La convencion de cultura pop mas grande del pais. Celebridades, cosplay, gaming y mucho mas.', 'https://turismo.buenosaires.gob.ar/sites/turismo/files/field/image/comic-con-arg-1500x610_0.png', '2026-10-09', '10:00:00', 480, 'Feria', 2000, 'activo', 1);

INSERT INTO sectores (id, evento_id, nombre, capacidad_maxima, capacidad_disponible, precio) VALUES
(1, 1, 'Campo', 2000, 1999, 15000.00),
(2, 1, 'Platea Baja', 1500, 1499, 25000.00),
(3, 1, 'Platea Alta', 1000, 1000, 18000.00),
(4, 1, 'VIP', 500, 500, 45000.00),
(5, 2, 'Popular', 1500, 1499, 8000.00),
(6, 2, 'Platea', 1200, 1198, 20000.00),
(7, 2, 'Platea Preferencial', 300, 299, 35000.00),
(8, 3, 'General', 1500, 1500, 5000.00),
(9, 3, 'VIP Pass', 500, 500, 12000.00);
