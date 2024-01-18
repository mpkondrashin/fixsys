
.PHONY: clean tidy

ifeq ($(OS),Windows_NT)
define zip
	\$compress = @{
  		Path = $(2), $(3), $(4)
		DestinationPath = $(1)
	}
	Compress-Archive @compress
endef
else
define zip
	zip $(1) $(2) $(3) $(4)
endef
endif
#powershell Compress-Archive  -Force $(2) $(3) $(4) $(1)
fix.zip: fixsys.exe README.txt
	$(call zip, "fix.zip" , "fixsys.exe", "README.txt")

fixsys.exe: cmd/fixsys/main.go
#	echo $(wildcard cmd/install/*.go)
	GOOS=windows GOARCH=amd64 go build ./cmd/fixsys

#deploy.exe: cmd/deploy/main.go
#	echo $(wildcard cmd/install/*.go)
#	GOOS=windows GOARCH=amd64 go build ./cmd/deploy
