#[derive(Debug, Serialize, Deserialize)]
pub struct Version {
    pub version: VersionNumber,
    pub url: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Endpoint {
    pub identifier: ModuleID,
    pub role: InterfaceRole,
    pub url: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct VersionDetails {
    pub version: VersionNumber,
    pub endpoints: Vec<Endpoint>,
}

#[derive(Debug, Serialize, Deserialize, PartialEq, Eq, Clone, Copy)]
#[serde(rename_all = "snake_case")]
pub enum InterfaceRole {
    Sender,
    Receiver,
}

#[derive(Debug, Serialize, Deserialize, PartialEq, Eq, Clone, Copy)]
#[serde(rename_all = "snake_case")]
pub enum ModuleID {
    Cdrs,
    ChargingProfiles,
    Commands,
    Credentials,
    HubClientInfo,
    Locations,
    Sessions,
    Tariffs,
    Tokens,
}
