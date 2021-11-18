use chacha20poly1305::aead::Aead;
use chacha20poly1305::{aead::NewAead, ChaCha20Poly1305, Key, Nonce, XChaCha20Poly1305, XNonce};
use getrandom::getrandom;
use std::fs::File;
use std::io::{BufReader, BufWriter, Read, Write};

pub const KEY_SIZE: usize = 32;
pub const NONCE_SIZE: usize = 12;
pub const XNONCE_SIZE: usize = 24;
pub const BUFF_BLOCK_SIZE: usize = 4096;
pub const TAG_SIZE: usize = 16;

pub fn new_chacha20poly1305() -> Result<(ChaCha20Poly1305, Key), Box<dyn std::error::Error>> {
    let mut key = [0u8; KEY_SIZE];
    getrandom(&mut key)?;
    let key = Key::from(key);
    Ok((ChaCha20Poly1305::new(&key), key))
}

pub fn new_xchacha20poly1305() -> Result<(XChaCha20Poly1305, Key), Box<dyn std::error::Error>> {
    let mut key = [0u8; KEY_SIZE];
    getrandom(&mut key)?;
    let key = Key::from(key);
    Ok((XChaCha20Poly1305::new(&key), key))
}

pub fn new_chacha_nonce() -> Result<Nonce, Box<dyn std::error::Error>> {
    let mut nonce = [0u8; NONCE_SIZE];
    getrandom(&mut nonce)?;
    Ok(Nonce::from(nonce))
}

pub fn new_xchacha_nonce() -> Result<XNonce, Box<dyn std::error::Error>> {
    let mut nonce = [0u8; XNONCE_SIZE];
    getrandom(&mut nonce)?;
    Ok(XNonce::from(nonce))
}

pub trait FileCrypt<NC> {
    fn encrypt_file(
        &self,
        nonce: NC,
        input_file: &str,
        output_file: &str,
    ) -> Result<(), Box<dyn std::error::Error>>;
    fn decrypt_file(
        &self,
        input_file: &str,
        output_file: &str,
    ) -> Result<(), Box<dyn std::error::Error>>;
}

impl FileCrypt<Nonce> for ChaCha20Poly1305 {
    fn encrypt_file(
        &self,
        nonce: Nonce,
        input_file: &str,
        output_file: &str,
    ) -> Result<(), Box<dyn std::error::Error>> {
        let mut input = File::open(input_file)?;
        let output = File::create(output_file)?;
        let mut output = BufWriter::new(output);
        let mut buffer = [0u8; BUFF_BLOCK_SIZE - TAG_SIZE];
        output.write(&nonce.as_ref())?;
        loop {
            if let Ok(read) = input.read(buffer.as_mut()) {
                if read == 0 {
                    break;
                }
                let enc = self.encrypt(&nonce, buffer[..read].as_ref()).unwrap();
                output.write_all(&enc)?;
            } else {
                println!("Error reading file");
                break;
            }
        }
        Ok(())
    }

    fn decrypt_file(
        &self,
        input_file: &str,
        output_file: &str,
    ) -> Result<(), Box<dyn std::error::Error>> {
        let mut input = File::open(input_file)?;
        let output = File::create(output_file)?;
        let mut output = BufWriter::new(output);
        let mut buffer = [0u8; BUFF_BLOCK_SIZE];
        let mut nonce = [0u8; NONCE_SIZE];
        input.read_exact(nonce.as_mut())?;
        let nonce = Nonce::from(nonce);
        loop {
            if let Ok(read) = input.read(buffer.as_mut()) {
                if read == 0 {
                    break;
                }
                // println!("read {} bytes", read);
                let dec = self.decrypt(&nonce, buffer[..read].as_ref()).unwrap();
                output.write_all(&dec)?;
            } else {
                break;
            }
        }
        Ok(())
    }
}

impl FileCrypt<XNonce> for XChaCha20Poly1305 {
    fn encrypt_file(
        &self,
        nonce: XNonce,
        input_file: &str,
        output_file: &str,
    ) -> Result<(), Box<dyn std::error::Error>> {
        let mut input = File::open(input_file)?;
        let output = File::create(output_file)?;
        let mut output = BufWriter::new(output);
        let mut buffer = [0u8; BUFF_BLOCK_SIZE - TAG_SIZE];
        output.write(&nonce.as_ref())?;
        loop {
            if let Ok(read) = input.read(buffer.as_mut()) {
                if read == 0 {
                    break;
                }
                let enc = self.encrypt(&nonce, buffer[..read].as_ref()).unwrap();
                output.write_all(&enc)?;
            } else {
                println!("Error reading file");
                break;
            }
        }
        Ok(())
    }

    fn decrypt_file(
        &self,
        input_file: &str,
        output_file: &str,
    ) -> Result<(), Box<dyn std::error::Error>> {
        let mut input = File::open(input_file)?;
        let output = File::create(output_file)?;
        let mut output = BufWriter::new(output);
        let mut buffer = [0u8; BUFF_BLOCK_SIZE];
        let mut nonce = [0u8; XNONCE_SIZE];
        input.read_exact(nonce.as_mut())?;
        let nonce = XNonce::from(nonce);
        loop {
            if let Ok(read) = input.read(buffer.as_mut()) {
                if read == 0 {
                    break;
                }
                // println!("read {} bytes", read);
                let dec = self.decrypt(&nonce, buffer[..read].as_ref()).unwrap();
                output.write_all(&dec)?;
            } else {
                break;
            }
        }
        Ok(())
    }
}
