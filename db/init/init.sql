CREATE TABLE IF NOT EXISTS cats (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    years_of_experience INT NOT NULL,
    breed TEXT NOT NULL,
    salary NUMERIC(10, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS missions (
    id SERIAL PRIMARY KEY,
    cat_id BIGINT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE IF NOT EXISTS targets (
    id SERIAL PRIMARY KEY,
    mission_id BIGINT NOT NULL REFERENCES missions(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(100) NOT NULL,
    notes TEXT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    UNIQUE(mission_id, name)
);