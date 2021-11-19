use blake2::{
    digest::{InvalidOutputSize, Update, VariableOutput},
    VarBlake2b,
};
use std::io::Error;

pub trait NewKey {
    fn new_rand_key() -> Result<Box<Self>, Error>;
    fn new_from_password(password: &str) -> Result<Box<Self>, InvalidOutputSize>;
}

impl NewKey for chacha20::Key {
    fn new_rand_key() -> Result<Box<Self>, Error> {
        let mut ret = Self::default();
        getrandom::getrandom(ret.as_mut_slice())?;
        Ok(Box::new(ret))
    }

    fn new_from_password(password: &str) -> Result<Box<Self>, InvalidOutputSize> {
        let mut ret = Vec::new();
        let mut hasher = VarBlake2b::new(32)?;
        hasher.update(password.as_bytes());
        hasher.finalize_variable(|res| {
            ret = res.to_vec();
        });
        let ret = chacha20::Key::from_slice(ret.as_slice());
        Ok(Box::new(*ret))
    }
}
