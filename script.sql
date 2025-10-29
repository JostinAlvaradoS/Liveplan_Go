INSERT INTO "categoria_costos" (id, nombre) VALUES
(1, 'Mano de obra'),(2, 'Materia prima'),(3, 'Gastos indirectos'),(4,'Otros');

insert into "tipo_inversion_inicials" (id,tipo) values (1, 'Activos fijos'),
                                                   (2, 'Gastos preoperativos y de constitucion'),
                                                   (3, 'Capital de trabajo inicial');

INSERT INTO gastos_operacion_bases (descripcion, valor) VALUES
('Sueldos gerenciales (incluye prestaciones sociales)', 1050),
('Sueldos de colaboradores (incluye prestaciones)', 1950),
('Uniformes', 30),
('Honorarios', 50),
('Publicidad', 150),
('Teléfono', 50),
('Energía eléctrica', 120),
('Agua', 42),
('Gas', 180),
('Gasolina', 120),
('Mantenimiento de vehículos', 45),
('Mantenimiento de planta y equipo', 63),
('Seguros contra daños', 65),
('Papelería y útiles para oficina', 23),
('Gastos de viaje/representación', 0),
('Renta de locales', 500),
('Otros', 0);