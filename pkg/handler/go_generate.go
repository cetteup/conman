//go:build ignore

package handler

//go:generate mockgen -source=handler.go -destination=handler_mock_test.go -package=$GOPACKAGE -write_package_comment=false
//go:generate mockgen -destination=os_mock_test.go -package=$GOPACKAGE os DirEntry
