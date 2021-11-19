use chacha20::XNonce;

pub trait GenNonce {
    fn gen() -> Result<Box<Self>, std::io::Error>;
}

impl GenNonce for XNonce {
    fn gen() -> Result<Box<XNonce> , std::io::Error> {
      let mut ret = XNonce::default();
      getrandom::getrandom(ret.as_mut_slice())?;
      Ok(Box::new(ret))
    }
}

