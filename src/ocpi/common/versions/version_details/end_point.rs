use serde::{Serialize, Deserialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Endpoint {
    pub identifier: ModuleID,
    pub role: InterfaceRole,
    pub url: String,
}

impl Endpoint {
    pub fn new(identifier: ModuleID, role: InterfaceRole, url: &str) -> Self {
        Endpoint {
            identifier,
            role,
            url: url.to_string(),
        }
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub enum InterfaceRole {
    SENDER,
    RECEIVER,
}

#[derive(Debug, Serialize, Deserialize)]
pub enum ModuleID {
    cdrs,
    chargingprofiles,
    commands,
    credentials,
    hubclientinfo,
    locations,
    sessions,
    tariffs,
    tokens,
}
