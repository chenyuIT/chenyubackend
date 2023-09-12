rem hz update -I idl -idl ./idl/user/user.proto --model_dir ./biz/api/models --unset_omitempty
hz update -I idl -idl ./idl/contact/contact.proto --model_dir ./biz/api/models --unset_omitempty
.\cmd\protoc-go-inject-tag.exe -input="./biz/api/models/*/*.pb.go"

