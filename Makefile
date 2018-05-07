proto: auth.pb.go

%.pb.go: proto/%.proto
	protoc -Iproto --go_out=common $<
