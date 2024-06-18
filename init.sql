create schema if not exists energy_usage;

create table if not exists energy_usage.c02
(
    id             uuid                     default gen_random_uuid() not null,
    from_timestamp timestamp with time zone,
    to_timestamp   timestamp with time zone,
    forecast       integer                                            not null,
    actual         integer                                            not null,
    created_at     timestamp with time zone default now()             not null,
    constraint c02_unique_date_time_slot unique (from_timestamp, to_timestamp)
);

create table if not exists energy_usage.mix
(
    id             uuid                     default gen_random_uuid() not null,
    from_timestamp timestamp with time zone,
    to_timestamp   timestamp with time zone,
    created_at     timestamp with time zone default now()             not null,
    constraint mix_unique_date_time_slot unique (from_timestamp, to_timestamp),
    CONSTRAINT mix_unique_id UNIQUE (id)

);

CREATE TABLE IF NOT EXISTS energy_usage.mix_generation
(
    id         UUID                     DEFAULT gen_random_uuid() NOT NULL,
    mix_id     UUID                                               NOT NULL REFERENCES energy_usage.mix (id) ON DELETE CASCADE,
    fuel       VARCHAR(255)                                       NOT NULL,
    perc       FLOAT                                              NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()             NOT NULL,
    CONSTRAINT mix_generation_unique_fuel UNIQUE (mix_id, fuel)
);

CREATE TABLE IF NOT EXISTS energy_usage.octopus_usage
(
    id             UUID                     DEFAULT gen_random_uuid() NOT NULL,
    consumption    float                                              NOT NULL,
    from_timestamp timestamp with time zone                           NOT NULL,
    to_timestamp   timestamp with time zone                           NOT NULL,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW()             NOT NULL
);

CREATE TABLE IF NOT EXISTS energy_usage.octopus_pricing
(
    id             UUID                     DEFAULT gen_random_uuid() NOT NULL,
    value_inc_vat  float                                              NOT NULL,
    value_exc_vat  float                                              NOT NULL,
    from_timestamp timestamp with time zone                           NOT NULL,
    to_timestamp   timestamp with time zone                           NOT NULL,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW()             NOT NULL
);

CREATE TABLE IF NOT EXISTS energy_usage.daily_summary
(
    id                        UUID DEFAULT gen_random_uuid() NOT NULL,
    date                      date                           NOT NULL,
    carbon_intensity          float                          NOT NULL,
    average_carbon_intensity  float                          NOT NULL,
    total_energy_used         float                          NOT NULL,
    total_energy_cost_inc_vat float                          NOT NULL,
    total_energy_cost_exc_vat float                          NOT NULL,
    mix_percentage            jsonb                          NOT NULL
);
