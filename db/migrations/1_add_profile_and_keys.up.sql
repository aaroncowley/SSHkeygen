DROP TABLE IF EXISTS profiles;
CREATE TABLE profiles(
    id SERIAL PRIMARY KEY,
    code_name VARCHAR (50) NOT NULL
);

DROP TABLE IF EXISTS keys;
CREATE TABLE keys(
    id SERIAL PRIMARY KEY,
    profile_id SERIAL,
    public_key VARCHAR (4096) NOT NULL,
    private_key VARCHAR (4096) NOT NULL,
    creation_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles (id)
);

