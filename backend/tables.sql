create table users (
  username text primary key,
  password text,
  description text,
  goalweight float,
  bodyweight float,
  caloriegoal float,
  caloriesleft float,
  followers text[],
  following text[],
  program
);

create table posts (
  id text primary key,
  username text,
  contents text,
  media text,
  date text,
  likes text[]
);

create table programs (
  username text primary key,
  programfile text,
  startdate text
)

create table customprogram(
  username text primary key,
  programdict json
)