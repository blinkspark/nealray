use blake2::{
    digest::{InvalidOutputSize, Update, VariableOutput},
    VarBlake2b,
};
use chacha20::XNonce;


pub fn from_password(passwd: &str) -> Result<Vec<u8>, InvalidOutputSize> {
    // let ret = null

    let mut ret: Vec<u8> = Vec::new();
    let mut hasher = VarBlake2b::new(32)?;
    hasher.update(passwd.as_bytes());
    hasher.finalize_variable(|res| {
        ret = res.to_vec();
    });
    Ok(ret)
}

pub fn from_rand() -> Result<Vec<u8>, std::io::Error> {
    let mut ret = [0u8; 32];
    getrandom::getrandom(&mut ret)?;
    Ok(ret.to_vec())
}

mod tests {
    #[test]
    fn test_from_rand() {
        let ret = super::from_rand();
        assert!(ret.is_ok());
    }
    #[test]
    fn test_from_password() {
        let ret = super::from_password("password");
        assert!(ret.is_ok());
    }
}
