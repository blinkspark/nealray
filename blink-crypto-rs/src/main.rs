#![allow(dead_code)]
// #![allow(unused_imports)]

mod nonce;
use blake2::{
    digest::{Update, VariableOutput},
    VarBlake2b,
};
use chacha20::{Key, XNonce};
use nonce::GenNonce;

mod key;
use key::NewKey;

use clap::{App, Arg};
use std::{
    fs::File,
    io::{Read, Write},
};

mod stream;
use stream::FileCryptXChacha20;

const BUFF_SIZE: usize = 4096;
const OUTPUT_PLACE_HOLDER: &str = "tmp.enc";

fn main() {
    let app = App::new("blink file crypto")
        .version("0.1.0")
        .about("encrypt and decrypt file")
        .subcommand(
            App::new("encrypt")
                .alias("enc")
                .about("encrypt -i INPUT -o OUTPUT [-p password] [-k key_path] [-g]")
                .arg(
                    Arg::new("input")
                        .short('i')
                        .long("input")
                        .about("-i FILE_PATH")
                        .takes_value(true),
                )
                .arg(
                    Arg::new("output")
                        .short('o')
                        .long("output")
                        .about("-o FILE_PATH")
                        .takes_value(true),
                )
                .arg(
                    Arg::new("password")
                        .about("-p password")
                        .short('p')
                        .long("password")
                        .takes_value(true),
                )
                .arg(
                    Arg::new("key")
                        .about("-k key")
                        .short('k')
                        .long("key")
                        .takes_value(true),
                )
                .arg(
                    Arg::new("gen")
                        .about("-g #flag represent to generate a new key")
                        .short('g')
                        .long("gen"),
                )
                .arg(
                    Arg::new("dir")
                        .about("-d DIR")
                        .short('d')
                        .long("dir")
                        .takes_value(true),
                ),
        )
        .subcommand(
            App::new("decrypt")
                .alias("dec")
                .about("decrypt -i INPUT -o OUTPUT [-p password] [-k key_path] [-g]")
                .arg(
                    Arg::new("input")
                        .short('i')
                        .long("input")
                        .about("-i FILE_PATH")
                        .takes_value(true),
                )
                .arg(
                    Arg::new("output")
                        .short('o')
                        .long("output")
                        .about("-o FILE_PATH")
                        .takes_value(true),
                )
                .arg(
                    Arg::new("password")
                        .about("-p password")
                        .short('p')
                        .long("password")
                        .takes_value(true),
                )
                .arg(
                    Arg::new("key")
                        .about("-k key")
                        .short('k')
                        .long("key")
                        .takes_value(true),
                ),
        )
        .get_matches();

    if let Some(enc_cmd) = app.subcommand_matches("encrypt") {
        let input = enc_cmd.value_of("input").expect("need a input file");
        let mut output = enc_cmd.value_of("output").unwrap_or("");
        let gen = enc_cmd.is_present("gen");
        let dir = enc_cmd.value_of("dir").unwrap_or("");
        let key = if let Some(password) = enc_cmd.value_of("password") {
            Key::new_from_password(password).expect("can not generate key from password!")
        } else if let Some(key_path) = enc_cmd.value_of("key") {
            let inner_key = if gen {
                let key = Key::new_rand_key().expect("can not generate key");
                let mut key_file = File::create(key_path).expect("can not create key file");
                key_file
                    .write_all(key.as_slice())
                    .expect("can not write key file");
                key
            } else {
                let mut key_file = File::open(key_path).expect("can not open key file");
                let mut buf = Vec::new();
                key_file
                    .read_to_end(&mut buf)
                    .expect("can not read key file");
                Box::new(*Key::from_slice(buf.as_slice()))
            };
            inner_key
        } else {
            panic!("need a password or key file!")
        };

        let need_rename = true;
        if output == "" {
            output = OUTPUT_PLACE_HOLDER;
        }

        // need to close files here
        {
            let mut input = File::open(input).expect("can not open input file");
            let mut output = File::create(output).expect("can not create output file");
            let nonce = XNonce::gen().expect("can not generate nonce");

            input
                .encrypt(&key, &nonce, &mut output)
                .expect("can not encrypt file");
        }

        if need_rename {
            // rename output file to its hash
            let mut hasher = VarBlake2b::new(32).expect("can not create hasher");
            let mut output = File::open(OUTPUT_PLACE_HOLDER).expect("can not open output file");
            let mut buffer = [0u8; BUFF_SIZE];
            while let Ok(read) = output.read(buffer.as_mut()) {
                if read == 0 {
                    break;
                }
                hasher.update(&buffer[..read]);
            }
            let mut hash = Vec::new();
            hasher.finalize_variable(|res| {
                hash = res.to_vec();
            });
            let hash = hex::encode(hash);
            let path = std::path::PathBuf::new().join(dir).join(hash + ".enc");
            std::fs::rename(OUTPUT_PLACE_HOLDER, path.to_str().unwrap())
                .expect("can not rename output file");
        }
    }
    if let Some(dec_cmd) = app.subcommand_matches("decrypt") {
        let input = dec_cmd.value_of("input").expect("need a input file");
        let output = dec_cmd.value_of("output").expect("need a output file");
        let mut input = File::open(input).expect("can not open input file");
        let mut output = File::create(output).expect("can not create output file");
        let password = dec_cmd.value_of("password");
        let key_path = dec_cmd.value_of("key");

        let key = if let Some(password) = password {
            Key::new_from_password(password).expect("can not generate key from password!")
        } else if let Some(key_path) = key_path {
            let mut key_file = File::open(key_path).expect("can not open key file");
            let mut buf = Vec::new();
            key_file
                .read_to_end(&mut buf)
                .expect("can not read key file");
            Box::new(*Key::from_slice(buf.as_slice()))
        } else {
            panic!("need a password or key file!");
        };

        input
            .decrypt(&key, &mut output)
            .expect("can not decrypt file");
    }
}
