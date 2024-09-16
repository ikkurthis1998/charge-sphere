use serde::{Serialize, Deserialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Version {
    pub version: String,
    pub url: String,
}

impl Version {
    pub fn new(version: &str, url: &str) -> Self {
        Version {
            version: version.to_string(),
            url: url.to_string(),
        }
    }
}
