-- Default demo data: two baptism register books for Santo Espíritu Sagrado,
-- with a handful of entries spread across them, plus a smaller set for
-- Capilla San José (used to test that books/registers/validators stay
-- scoped to their own organization).

INSERT INTO books (org_id, book_nr, date_initial, date_final, book_type, created_at, updated_at)
SELECT id, 12, '1978-01-01', '1988-12-31', 2, now(), now()
FROM organizations WHERE name = 'Santo Espíritu Sagrado';

INSERT INTO books (org_id, book_nr, date_initial, date_final, book_type, created_at, updated_at)
SELECT id, 13, '1989-01-01', '1999-12-31', 2, now(), now()
FROM organizations WHERE name = 'Santo Espíritu Sagrado';

-- Book 12, page 45, reg 112
WITH next_index_id AS MATERIALIZED (
    SELECT nextval(pg_get_serial_sequence('index_baptisms', 'id')) AS id
),
reg AS (
    INSERT INTO register_baptisms (book_id, page_number, reg_number, index_id, org_baptism,
        baptizer, date_baptism, baptized_name_f, baptized_name_s, rut, date_birth, place_birth,
        father_name, father_surname, mother_name, mother_surname, godfather, godmother,
        created_at, updated_at)
    SELECT b.id, 45, 112, next_index_id.id, o.id,
        'Padre Pedro', '1982-05-16', 'María', 'Fernanda', '9.876.543-2', '1982-04-02', 'Puente Alto',
        'Juan', 'Soto Pardo', 'Carmen', 'Rojas Díaz', 'Manuel Herrera López', 'Isidora Contreras Muñoz',
        now(), now()
    FROM books b, organizations o, next_index_id
    WHERE b.book_nr = 12 AND o.name = 'Santo Espíritu Sagrado'
    RETURNING id, book_id, org_baptism
)
INSERT INTO index_baptisms (id, org_id, book_id, reg_id, user_surname_f, user_surname_m,
    user_name_first, user_name_second, page_number, created_at, updated_at)
SELECT next_index_id.id, reg.org_baptism, reg.book_id, reg.id, 'Soto Pardo', 'Rojas Díaz', 'María', 'Fernanda', 45, now(), now()
FROM reg, next_index_id;

-- Book 12, page 78, reg 145
WITH next_index_id AS MATERIALIZED (
    SELECT nextval(pg_get_serial_sequence('index_baptisms', 'id')) AS id
),
reg AS (
    INSERT INTO register_baptisms (book_id, page_number, reg_number, index_id, org_baptism,
        baptizer, date_baptism, baptized_name_f, baptized_name_s, rut, date_birth, place_birth,
        father_name, father_surname, mother_name, mother_surname, godfather, godmother,
        created_at, updated_at)
    SELECT b.id, 78, 145, next_index_id.id, o.id,
        'Padre Pedro', '1985-11-23', 'Cristóbal', '', '10.234.567-8', '1985-10-11', 'Santiago',
        'Ricardo', 'Fuentes Silva', 'Patricia', 'Vergara Muñoz', 'José Muñoz Castro', 'Teresa Silva Vega',
        now(), now()
    FROM books b, organizations o, next_index_id
    WHERE b.book_nr = 12 AND o.name = 'Santo Espíritu Sagrado'
    RETURNING id, book_id, org_baptism
)
INSERT INTO index_baptisms (id, org_id, book_id, reg_id, user_surname_f, user_surname_m,
    user_name_first, user_name_second, page_number, created_at, updated_at)
SELECT next_index_id.id, reg.org_baptism, reg.book_id, reg.id, 'Fuentes Silva', 'Vergara Muñoz', 'Cristóbal', '', 78, now(), now()
FROM reg, next_index_id;

-- Book 13, page 8, reg 201
WITH next_index_id AS MATERIALIZED (
    SELECT nextval(pg_get_serial_sequence('index_baptisms', 'id')) AS id
),
reg AS (
    INSERT INTO register_baptisms (book_id, page_number, reg_number, index_id, org_baptism,
        baptizer, date_baptism, baptized_name_f, baptized_name_s, rut, date_birth, place_birth,
        father_name, father_surname, mother_name, mother_surname, godfather, godmother,
        created_at, updated_at)
    SELECT b.id, 8, 201, next_index_id.id, o.id,
        'Padre Pedro', '1993-08-01', 'Martín', 'Ignacio', '13.456.789-0', '1993-06-19', 'Puente Alto',
        'Álvaro', 'Espinoza Reyes', 'Daniela', 'Contreras López', 'Francisco Reyes Ortiz', 'Valentina López Cerda',
        now(), now()
    FROM books b, organizations o, next_index_id
    WHERE b.book_nr = 13 AND o.name = 'Santo Espíritu Sagrado'
    RETURNING id, book_id, org_baptism
)
INSERT INTO index_baptisms (id, org_id, book_id, reg_id, user_surname_f, user_surname_m,
    user_name_first, user_name_second, page_number, created_at, updated_at)
SELECT next_index_id.id, reg.org_baptism, reg.book_id, reg.id, 'Espinoza Reyes', 'Contreras López', 'Martín', 'Ignacio', 8, now(), now()
FROM reg, next_index_id;

INSERT INTO books (org_id, book_nr, date_initial, date_final, book_type, created_at, updated_at)
SELECT id, 5, '2005-01-01', '2015-12-31', 2, now(), now()
FROM organizations WHERE name = 'Capilla San José';

-- Book 5, page 12, reg 33
WITH next_index_id AS MATERIALIZED (
    SELECT nextval(pg_get_serial_sequence('index_baptisms', 'id')) AS id
),
reg AS (
    INSERT INTO register_baptisms (book_id, page_number, reg_number, index_id, org_baptism,
        baptizer, date_baptism, baptized_name_f, baptized_name_s, rut, date_birth, place_birth,
        father_name, father_surname, mother_name, mother_surname, godfather, godmother,
        created_at, updated_at)
    SELECT b.id, 12, 33, next_index_id.id, o.id,
        'Hermano Francisco', '2008-09-14', 'Valentina', 'Isidora', '15.678.901-2', '2008-08-03', 'Providencia',
        'Sebastián', 'Muñoz Vera', 'Camila', 'Torres Bravo', 'Andrés Bravo Soto', 'Francisca Vera Lagos',
        now(), now()
    FROM books b, organizations o, next_index_id
    WHERE b.book_nr = 5 AND o.name = 'Capilla San José'
    RETURNING id, book_id, org_baptism
)
INSERT INTO index_baptisms (id, org_id, book_id, reg_id, user_surname_f, user_surname_m,
    user_name_first, user_name_second, page_number, created_at, updated_at)
SELECT next_index_id.id, reg.org_baptism, reg.book_id, reg.id, 'Muñoz Vera', 'Torres Bravo', 'Valentina', 'Isidora', 12, now(), now()
FROM reg, next_index_id;

-- Book 5, page 27, reg 41
WITH next_index_id AS MATERIALIZED (
    SELECT nextval(pg_get_serial_sequence('index_baptisms', 'id')) AS id
),
reg AS (
    INSERT INTO register_baptisms (book_id, page_number, reg_number, index_id, org_baptism,
        baptizer, date_baptism, baptized_name_f, baptized_name_s, rut, date_birth, place_birth,
        father_name, father_surname, mother_name, mother_surname, godfather, godmother,
        created_at, updated_at)
    SELECT b.id, 27, 41, next_index_id.id, o.id,
        'Hermano Francisco', '2011-03-20', 'Tomás', '', '17.345.678-9', '2011-02-11', 'Santiago',
        'Rodrigo', 'Salinas Peña', 'Javiera', 'Cortés Núñez', 'Felipe Núñez Ramos', 'Constanza Peña Ortiz',
        now(), now()
    FROM books b, organizations o, next_index_id
    WHERE b.book_nr = 5 AND o.name = 'Capilla San José'
    RETURNING id, book_id, org_baptism
)
INSERT INTO index_baptisms (id, org_id, book_id, reg_id, user_surname_f, user_surname_m,
    user_name_first, user_name_second, page_number, created_at, updated_at)
SELECT next_index_id.id, reg.org_baptism, reg.book_id, reg.id, 'Salinas Peña', 'Cortés Núñez', 'Tomás', '', 27, now(), now()
FROM reg, next_index_id;
