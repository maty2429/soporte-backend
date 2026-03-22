-- Departamentos de soporte
INSERT INTO departamentos_soporte (descripcion, cod_departamento)
VALUES ('Informática', 'IT'),
       ('Plomería', 'PLMB'),
       ('Carpintería', 'CARP'),
       ('Electricidad', 'ELEC'),
       ('Mantenimiento Gral', 'MANT'),
       ('Climatización', 'CLIM'),
       ('Redes/Telecom', 'NET'),
       ('Gases Clínicos', 'GAS')
ON CONFLICT (cod_departamento) DO NOTHING;

-- Niveles de prioridad
INSERT INTO niveles_prioridad (id, descripcion)
VALUES (1, 'CRÍTICA'),
       (2, 'ALTA'),
       (3, 'MEDIA'),
       (4, 'BAJA')
ON CONFLICT (id) DO NOTHING;

-- Tipos de ticket
INSERT INTO tipo_ticket (id, cod_tipo_ticket, descripcion)
VALUES (1, 'INC', 'INCIDENTE'),
       (2, 'REQ', 'REQUERIMIENTO')
ON CONFLICT (id) DO NOTHING;

-- Estados de ticket
INSERT INTO estado_ticket (id, descripcion, cod_estado_ticket)
VALUES (1, 'CREADO', 'CRE'),
       (2, 'ASIGNADO', 'ASI'),
       (3, 'EN PROGRESO', 'PRO'),
       (4, 'PAUSADO', 'PAU'),
       (5, 'RESUELTO', 'RES'),
       (6, 'CERRADO', 'CER'),
       (7, 'CANCELADO', 'CAN'),
       (8, 'TRABAJO TERMINADO', 'TER'),
       (9, 'REABIERTO', 'REA'),
       (10, 'VISTO POR EL TÉCNICO', 'VITEC')
ON CONFLICT (id) DO NOTHING;
