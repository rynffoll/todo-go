-- +migrate Up
CREATE TABLE todos (
id serial PRIMARY KEY,
title varchar NOT NULL,
done boolean DEFAULT false,
date timestamp DEFAULT NOW()
);

-- +migrate Down
DROP TABLE todos;
