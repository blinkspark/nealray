use blake2::{
    digest::{InvalidOutputSize, Update, VariableOutput},
    VarBlake2b,
};
use getrandom::getrandom;

pub trait NewKey: Sized {
    fn new_random() -> Result<Self, std::io::Error>;
    fn new_from_passowrd(password: &str) -> Result<Self, InvalidOutputSize>;
}

impl NewKey for chacha20::Key {
    fn new_random() -> Result<Self, std::io::Error> {
        let mut key = chacha20::Key::default();
        getrandom(&mut key)?;
        Ok(key)
    }

    fn new_from_passowrd(password: &str) -> Result<Self, InvalidOutputSize> {
        let mut hasher = VarBlake2b::new(32).unwrap();
        hasher.update(password.as_bytes());
        let mut key = chacha20::Key::default();
        key.clone_from_slice(&hasher.finalize_boxed());
        Ok(key)
    }
}
