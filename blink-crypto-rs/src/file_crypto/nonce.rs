use getrandom::getrandom;
pub trait NewNonce: Sized {
    fn new_nonce() -> Result<Self, getrandom::Error>;
}

impl NewNonce for chacha20::XNonce {
    fn new_nonce() -> Result<Self, getrandom::Error> {
        let mut nonce = chacha20::XNonce::default();
        getrandom(&mut nonce)?;
        Ok(nonce)
    }
}
