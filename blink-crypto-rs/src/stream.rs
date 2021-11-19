use std::{fs::File, io::Read, io::Write};

use chacha20::{
    cipher::{NewCipher, StreamCipher},
    Key, XChaCha20, XNonce,
};
const BUFF_SIZE: usize = 4096;
const XNONCE_SIZE: usize = 24;

pub trait FileCrypt<K, N> {
    fn encrypt(&mut self, key: &K, nonce: &N, output: &File) -> Result<(), std::io::Error>;
    fn decrypt(&mut self, key: &K, output: &File) -> Result<(), std::io::Error>;
}

impl FileCrypt<Key, XNonce> for File {
    fn encrypt(&mut self, key: &Key, nonce: &XNonce, mut output: &File) -> Result<(), std::io::Error> {
        let mut xchacha = XChaCha20::new(key, nonce);
        let mut buffer = [0u8; BUFF_SIZE];
        output.write_all(nonce.as_slice())?;
        loop {
            let read = self.read(&mut buffer)?;
            if read == 0 {
                break;
            }
            xchacha.apply_keystream(&mut buffer[..read]);
            output.write_all(&buffer[..read])?;
        }
        Ok(())
    }

    fn decrypt(&mut self, key: &Key, mut output: &File) -> Result<(), std::io::Error> {
        let mut nonce = [0u8; XNONCE_SIZE];
        self.read_exact(&mut nonce)?;
        let nonce = XNonce::from_slice(&nonce);

        let mut xchacha = XChaCha20::new(key, nonce);
        let mut buffer = [0u8; BUFF_SIZE];
        loop {
            let read = self.read(&mut buffer)?;
            if read == 0 {
                break;
            }
            xchacha.apply_keystream(&mut buffer[..read]);
            output.write_all(&buffer[..read])?;
        }
        Ok(())
    }
}
