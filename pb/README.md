protoc ./rpc.proto --go_out=./

protoc ./rpc.proto --js_out=import_style=commonjs,binary:.

npm install google-protobuf

browserify exports.js -o  rpc_main.js