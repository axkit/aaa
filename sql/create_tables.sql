CREATE SCHEMA aaa;

SET search_path TO aaa;

CREATE TABLE users (
    id int8 NOT NULL,
    role_id int2 NOT NULL,
    login text NOT NULL,
    encrypted_password text NOT NULL,
    locked_at timestamptz,
    row_version int8 NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz,
    deleted_at timestamptz,
    CONSTRAINT users_pk PRIMARY KEY (id)
);

CREATE UNIQUE INDEX users_login_uidx ON users (upper(LOGIN));

CREATE TABLE roles (
    id int2 NOT NULL,
    name text NOT NULL,
    permissions_bit_set int8 NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz,
    deleted_at timestamptz,
    CONSTRAINT roles_pk PRIMARY KEY (id)
);

CREATE TABLE permissions (
    id text NOT NULL,
    bit_pos int8 NOT NULL,
    name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz,
    deleted_at timestamptz,
    CONSTRAINT permissions_pk PRIMARY KEY (id),
    CONSTRAINT permissions_uk UNIQUE (bit_pos)
);

INSERT INTO permissions (id, bit_pos, name)
    VALUES ('TestCreateEntity', 0, 'role_name1'), ('TestUpdateEntity', 1, 'role_name2'), ('TestDeleteEntity', 2, 'role_name2');

INSERT INTO roles (id, name, permissions_bit_set)
    VALUES (1, 'TestAdmin', 7), (2, 'TestUser', 3), (3, 'TestUserNoPermissions', 0);

INSERT INTO users (id, LOGIN, role_id, encrypted_password)
    VALUES (11, 'testadmin', 1, 'test'), (12, 'testuser', 2, 'test'), (13, 'testnoperms', 3, 'test');

INSERT INTO users (id, LOGIN, role_id, encrypted_password, locked_at)
    VALUES (14, 'testlocked', 1, 'test', now());

