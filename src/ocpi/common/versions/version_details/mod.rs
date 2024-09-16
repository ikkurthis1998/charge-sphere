mod end_point;

use serde::{Serialize, Deserialize};
pub use end_point::*;

#[derive(Debug, Serialize, Deserialize)]
pub struct VersionDetails {
    pub version: String,
    pub endpoints: Vec<Endpoint>,
}

impl VersionDetails {
    pub fn new(version: &str, endpoints: Vec<Endpoint>) -> Self {
        VersionDetails {
            version: version.to_string(),
            endpoints,
        }
    }
}