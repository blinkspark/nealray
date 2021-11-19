#![allow(dead_code)]
// #![allow(unused_imports)]

use std::fs::File;

use chacha20::{Key, XNonce};

mod key;
mod nonce;
mod stream;
use key::NewKey;
use nonce::GenNonce;
use stream::FileCryptChacha20;
fn main() {
    // let (chacha, _key) = new_xchacha20poly1305()?;
    // let nonce = new_xchacha_nonce()?;
    // chacha.encrypt_file(nonce, "tmp/test.png", "tmp/test.png.enc")?;
    // chacha.decrypt_file("tmp/test.png.enc", "tmp/test.png.png")?;
    let key = Key::new_from_password("123456").unwrap();
    // let nonce = XNonce::gen().unwrap();

    let mut f = File::open("tmp/test.png.enc").unwrap();
    let mut of = File::create("tmp/test.png.png").unwrap();
    f.decrypt(&key, &mut of).unwrap();
}
