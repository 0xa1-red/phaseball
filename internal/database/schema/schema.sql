CREATE TYPE POSITION AS ENUM('C','P','1B','2B','3B','SS','LF','CF','RF');
CREATE TYPE HANDEDNESS AS ENUM('L', 'R', 'S');

CREATE TABLE teams (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING UNIQUE
);

CREATE TABLE players (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    idteam UUID,
    name STRING UNIQUE,
    pos POSITION,
    hand HANDEDNESS,
    batter_pow SMALLINT,
    batter_con SMALLINT,
    batter_eye SMALLINT,
    batter_spd SMALLINT,
    batter_def SMALLINT,
    pitcher_fb SMALLINT,
    pitcher_ch SMALLINT,
    pitcher_bb SMALLINT,
    pitcher_ctl SMALLINT,
    pitcher_bat SMALLINT,

    FOREIGN KEY (idteam) REFERENCES teams (id) ON DELETE SET NULL
);

CREATE TABLE games (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    idaway UUID NOT NULL,
    idhome UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (idaway) REFERENCES teams (id),
    FOREIGN KEY (idhome) REFERENCES teams (id)
);

CREATE TABLE game_logs (
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    idgame UUID NOT NULL,
    entry JSONB NOT NULL,

    FOREIGN KEY (idgame) REFERENCES games (id)
);