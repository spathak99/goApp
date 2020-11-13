create table users (
  username text primary key,
  password text,
  description text,
  goalweight float,
  bodyweight float,
  caloriegoal float,
  caloriesleft float,
  followers text[],
  following text[]
);

create table posts (
  id text primary key,
  username text,
  contents text,
  media text,
  date text
);

create table programs (
  username text primary key,
  file text,
  startdate text
)
