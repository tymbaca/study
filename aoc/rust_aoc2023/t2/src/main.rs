#![allow(dead_code, unused)]

use lazy_static::lazy_static;
use regex::{Match, Regex};
use std::{
    fmt::Debug,
    fs,
    io::{self, Write},
};

lazy_static! {
    static ref GAME_NUM_REGEX: Regex = Regex::new(r"Game ([0-9]+)").unwrap();
    static ref ROUND_REGEX: Regex = Regex::new(
        r"(?:(?P<red>\d+)\s*red)?\s*(?:(?P<blue>\d+)\s*blue)?\s*(?:(?P<green>\d+)\s*green)?"
    )
    .unwrap();
}

enum Cubes {
    Red(i32),
    Green(i32),
    Blue(i32),
}

type ReqsCount = i32;

struct Reqs {
    red: ReqsCount,
    green: ReqsCount,
    blue: ReqsCount,
}

type GameNum = i32;

struct Round {
    red: Cubes,
    green: Cubes,
    blue: Cubes,
}

impl Round {
    fn to_round(text: &str) -> Result<Round, String> {
        let caps = match ROUND_REGEX.captures(text) {
            Some(v) => v,
            None => return Err(format!("can't find cubes in line")),
        };
        let red_opt = caps.name("red");
        let green_opt = caps.name("green");
        let blue_opt = caps.name("blue");

        match (red_opt, green_opt, blue_opt) {
            (None, None, None) => return Err("not found any cubes in round".to_string()),
            _ => (),
        }

        let round = Round {
            red: Cubes::Red(parse_int_from_opt_match(red_opt)?),
            green: Cubes::Green(parse_int_from_opt_match(green_opt)?),
            blue: Cubes::Blue(parse_int_from_opt_match(blue_opt)?),
        };

        return Ok(round);
    }
}

struct Game {
    num: GameNum,
    rounds: Vec<Round>,
}

impl Game {
    fn from_line(value: &str, line_num: i64) -> Result<Self, String> {
        let (left, right) = match value.split_once(":") {
            Some((left, right)) => (left, right),
            None => {
                return Err(format!(
                    "no game number part separated with ':' in line {}",
                    line_num
                ));
            }
        };

        let num = match Self::to_gamenum(left, line_num) {
            Ok(v) => v,
            Err(e) => return Err(format!("can't get game number in line {}: {}", line_num, e)),
        };
        let rounds = match Self::to_rounds(right) {
            Ok(v) => v,
            Err(e) => return Err(format!("can't parse rounds in line {}: {}", line_num, e)),
        };
        return Ok(Game { num, rounds });
    }

    fn to_gamenum(text: &str, line_num: i64) -> Result<GameNum, String> {
        let caps = match GAME_NUM_REGEX.captures(text) {
            Some(v) => v,
            None => return Err("game number not found".to_string()),
        };

        let Some(num_str) = caps.get(1) else {
            return Err("game number not found".to_string());
        };

        let num: i32 = match num_str.as_str().parse() {
            Ok(n) => n,
            Err(e) => {
                return Err(format!(
                    "cannot parse '{}' to integer, error: {}",
                    num_str.as_str(),
                    e.to_string(),
                ))
            }
        };

        return Ok(num);
    }

    fn to_rounds(text: &str) -> Result<Vec<Round>, String> {
        let round_texts = text.split(";");

        let mut rounds: Vec<Round> = vec![];
        let mut round_num = 1;
        for raw_round in round_texts {
            match Round::to_round(raw_round) {
                Ok(v) => rounds.push(v),
                Err(e) => return Err(format!("error in round {}: {}", round_num, e)),
            }
            round_num += 1;
        }
        return Ok(rounds);
    }
}

struct Games(Vec<Game>);

impl Games {
    fn read_from_file(path: &str) -> Result<Self, String> {
        let data = match fs::read_to_string(path) {
            Ok(v) => v,
            Err(e) => return Err(e.to_string()),
        };

        let mut games: Vec<Game> = vec![];

        let mut line_num = 1;
        for line in data.lines() {
            let game = Game::from_line(line, line_num)?;
            games.push(game);

            line_num += 1;
        }

        Ok(Games(games))
    }
}

fn check_if_game_possible(reqs: Reqs, game: Game) -> bool {
    // for cubes in game {
    //     match cubes {
    //         Cubes::Red(n) => {
    //             if n > reqs.reds {
    //                 return false;
    //             }
    //         }
    //         Cubes::Green(n) => {
    //             if n > reqs.greens {
    //                 return false;
    //             }
    //         }
    //         Cubes::Blue(n) => {
    //             if n > reqs.blues {
    //                 return false;
    //             }
    //         }
    //     }
    // }
    return true;
}

fn parse_int_from_opt_match(opt_match: Option<Match>) -> Result<i32, String> {
    let mtch = match opt_match {
        Some(m) => m,
        None => return Ok(0),
    };

    let text = mtch.as_str();

    let num: i32 = match text.parse() {
        Ok(n) => n,
        Err(e) => {
            return Err(format!(
                "can't parse '{}' to integer, err: {}",
                text,
                e.to_string()
            ))
        }
    };
    return Ok(num);
}

fn main() {
    println!("Hello, world!");
    let reqs = Reqs {
        red: 12,
        green: 13,
        blue: 14,
    };
    let games_res = Games::read_from_file("input.txt");
    match games_res {
        Ok(_) => println!("ok"),
        Err(e) => println!("{}", e),
    }
}
