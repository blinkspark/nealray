
mod lib;
use lib::{new_xchacha20poly1305, new_xchacha_nonce, FileCrypt};

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let (chacha, _key) = new_xchacha20poly1305()?;
    let nonce = new_xchacha_nonce()?;
    chacha.encrypt_file(nonce, "tmp/test.png", "tmp/test.png.enc")?;
    chacha.decrypt_file("tmp/test.png.enc", "tmp/test.png.png")?;
    Ok(())
}
