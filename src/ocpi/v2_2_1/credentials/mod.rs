use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Credentials {
    pub token: String,
    pub url: String,
    pub roles: Vec<CredentialsRole>,
}

impl Credentials {
    pub fn new(token: &str, url: &str, roles: Vec<CredentialsRole>) -> Self {
        Credentials {
            token: token.to_string(),
            url: url.to_string(),
            roles,
        }
    }

    pub fn default() -> Self {
        Credentials {
            token: String::new(),
            url: String::new(),
            roles: Vec::new(),
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CredentialsRole {
    pub role: Role,
    pub party_id: String,
    pub country_code: String,
    pub business_details: BusinessDetails,
}

impl CredentialsRole {
    pub fn new(role: Role, party_id: &str, country_code: &str, business_details: BusinessDetails) -> Self {
        CredentialsRole {
            role,
            party_id: party_id.to_string(),
            country_code: country_code.to_string(),
            business_details,
        }
    }

    pub fn default() -> Self {
        CredentialsRole {
            role: Role::Cpo,
            party_id: String::new(),
            country_code: String::new(),
            business_details: BusinessDetails::default(),
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BusinessDetails {
    pub name: String,
    pub logo: Option<Logo>,
    pub website: Option<String>,
}

impl BusinessDetails {
    pub fn new(name: &str, logo: Option<Logo>, website: Option<String>) -> Self {
        BusinessDetails {
            name: name.to_string(),
            logo,
            website,
        }
    }

    pub fn default() -> Self {
        BusinessDetails {
            name: String::new(),
            logo: None,
            website: None,
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Logo {
    pub url: String,
    pub thumbnail: Option<String>,
    pub category: LogoCategory,
    pub r#type: String,
    pub width: u32,
    pub height: u32,
}

impl Logo {
    pub fn new(url: &str, thumbnail: Option<String>, category: LogoCategory, r#type: &str, width: u32, height: u32) -> Self {
        Logo {
            url: url.to_string(),
            thumbnail,
            category,
            r#type: r#type.to_string(),
            width,
            height,
        }
    }

    pub fn default() -> Self {
        Logo {
            url: String::new(),
            thumbnail: None,
            category: LogoCategory::Operator,
            r#type: String::new(),
            width: 0,
            height: 0,
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum LogoCategory {
    #[serde(rename = "OPERATOR")]
    Operator,
    #[serde(rename = "PROVIDER")]
    Provider,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum Role {
    #[serde(rename = "CPO")]
    Cpo,
    #[serde(rename = "EMSP")]
    Emsp,
    #[serde(rename = "Hub")]
    Hub,
}