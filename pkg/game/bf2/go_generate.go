//go:build ignore

package bf2

//go:generate mockgen -source=../common.go -destination=bf2_mock_test.go -package=$GOPACKAGE -write_package_comment=false
