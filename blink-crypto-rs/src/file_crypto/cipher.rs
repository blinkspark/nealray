use chacha20::{cipher::StreamCipher, XChaCha20};
use std::{
    fs::File,
    io::{BufReader, BufWriter, Read, Write},
};

const BUFF_SIZE: usize = 4096;

pub trait FileCrypto: Read {
    fn encrypt_file_stream<T: StreamCipher>(
        &mut self,
        output_path: &str,
        cipher: &mut T,
    ) -> Result<(), std::io::Error>;
    // fn decrypt_file(&mut self, input_file: &str, output_file: &str) -> Result<(), std::io::Error>;
}

impl FileCrypto for File {
    fn encrypt_file_stream<T: StreamCipher>(
        &mut self,
        output_path: &str,
        cipher: &mut T,
    ) -> Result<(), std::io::Error> {
        let mut output_file = BufWriter::new(File::create(output_path)?);
        let mut buffer = [0u8; BUFF_SIZE];
        loop {
            let read_size = self.read(&mut buffer)?;
            if read_size == 0 {
                break;
            }
            cipher.apply_keystream(&mut buffer[..read_size]);
            output_file.write_all(&buffer[..read_size])?;
        }
        output_file.flush()?;
        Ok(())
    }
}
