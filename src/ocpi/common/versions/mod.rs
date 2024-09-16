mod version;
mod version_details;

pub use version::*;
pub use version_details::*;

pub trait VersionsModule {
    fn get_versions(&self) -> Result<Vec<Version>, Box<dyn std::error::Error>>;
    fn get_version_details(&self, version: &str) -> Result<VersionDetails, Box<dyn std::error::Error>>;
}
