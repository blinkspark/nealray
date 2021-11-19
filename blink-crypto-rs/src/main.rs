#![allow(dead_code)]

use std::ops::Deref;

use chacha20::{Key, XNonce};

mod stream;
mod key;
mod nonce;
use nonce::Gen;

fn main() {
    // let (chacha, _key) = new_xchacha20poly1305()?;
    // let nonce = new_xchacha_nonce()?;
    // chacha.encrypt_file(nonce, "tmp/test.png", "tmp/test.png.enc")?;
    // chacha.decrypt_file("tmp/test.png.enc", "tmp/test.png.png")?;
    let nonce = XNonce::gen().unwrap();
    println!("{:?}", nonce);
}
