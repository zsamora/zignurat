INSERT INTO accounts (username, password, module, acc_role, created_at, updated_at)
VALUES ('admin@identirat.com', '$2a$14$gFmofQMcPVnEF4ztNGF73OsNpbqN7Hi/pXwgjuLrnFON4t1wJq5km',
           1, 1, now(), now()) -- identiratadmin123
    ON CONFLICT (username) WHERE deleted_at IS NULL DO NOTHING;

INSERT INTO accounts (username, password, module, acc_role, created_at, updated_at)
VALUES ('admin@notariat.com', '$2a$14$Y0I7FdbI/2BcISooGSpagu6NjhBvAyn/3c6v/ZwLlhiYyOLIhvola',
           2, 1, now(), now()) -- notariatadmin123
    ON CONFLICT (username) WHERE deleted_at IS NULL DO NOTHING;

INSERT INTO accounts (username, password, module, acc_role, owner_uuid, created_at, updated_at)
VALUES ('org@notariat.com', '$2a$14$Qe5jToR.H.IuspKiydz9Heq6yyZPHIZ0ZpXsWiGA2BmXD23S.8uBG',
           2, 2, '67f484f7-ffe9-46ee-b6e6-0b47066afc5a', now(), now()) -- notariatorg123
    ON CONFLICT (username) WHERE deleted_at IS NULL DO NOTHING;

INSERT INTO accounts (username, password, module, acc_role, owner_uuid, created_at, updated_at)
VALUES ('validator@notariat.com', '$2a$14$a8o5EifSElSM84x78WJ0S.Rn8GMuqMlm1ukvGiEnti.5rjJhFnMJS',
           2, 3, 'f4e3b845-942b-4082-ab9e-aa045f15ee26', now(), now()) -- notariatvalidator123
    ON CONFLICT (username) WHERE deleted_at IS NULL DO NOTHING;

INSERT INTO accounts (username, password, module, acc_role, owner_uuid, created_at, updated_at)
VALUES ('user@notariat.com', '$2a$14$7BVE0SFxScfXsGv4QMjfR.uEMC7cYJW/Rys7vaSrJJAKX9aTa8xiG',
           2, 4, '830169b7-e3bd-44e8-8f3a-ce1e5b2a5ff8', now(), now()) -- notariatuser123
    ON CONFLICT (username) WHERE deleted_at IS NULL DO NOTHING;