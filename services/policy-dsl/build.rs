fn main() {
    // Support both local builds (from repo root: ../../services/proto)
    // and container builds where we copy into ./services/proto
    let local_path = "../../services/proto";
    let vendored_path = "services/proto";
    let proto_root = if std::path::Path::new(vendored_path).exists() {
        vendored_path
    } else {
        local_path
    };
    println!("cargo:rerun-if-changed={}/wohnfair/policy/v1/policy.proto", proto_root);
    println!("cargo:rerun-if-changed={}/wohnfair/common/v1/types.proto", proto_root);

    tonic_build::configure()
        .build_server(true)
        .build_client(false)
        .out_dir("src/gen")
        .compile(
            &[
                &format!("{}/wohnfair/policy/v1/policy.proto", proto_root),
            ],
            &[proto_root],
        )
        .expect("Failed to compile protos");
}
