INSERT INTO accounts (username, password, module, acc_role, created_at, updated_at)
VALUES ('admin@identirat.com', '$2a$14$gFmofQMcPVnEF4ztNGF73OsNpbqN7Hi/pXwgjuLrnFON4t1wJq5km',
           1, 1, now(), now()) -- identiratadmin123
    ON CONFLICT (username) WHERE deleted_at IS NULL DO NOTHING;