fn main() {
    let version = charge_sphere::ocpi::versions::Version::new("2.2.1", "https://example.com");
    println!("{:?}", version);
}
